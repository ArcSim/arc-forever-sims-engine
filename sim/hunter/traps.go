package hunter

import (
	"math"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// Each trap's cooldown is its client category's: Immolation and Explosive Trap share 411, Freezing
// and Frost Trap 2183, 30 sec each, in client 1.60.1.70205 (Classic Era put every trap in 411 at
// 15 sec), so a Freezing Trap no longer locks the fire traps.
func trapCooldown(hunter *Hunter, rank *spelldata.Spell) core.Cooldown {
	return core.Cooldown{Timer: hunter.CategoryTimer(int32(rank.Category)), Duration: rank.CategoryCooldown()}
}

// Trap damage takes 13% of the hunter's ranged or melee attack power, whichever is higher, and no
// spell power, as Hmm measured in game (MythicSim Discord, October 2026, client 70425). The client
// carries neither share (SpellEffect 691047), so this is server side. It covers all of Immolation
// Trap, split over the burn's ticks, and Explosive Trap's explosion, but not the burning ground the
// explosion leaves.
const trapAttackPowerShare = 0.13

// trapAttackPower is the attack power a trap reads, and whether it is the ranged one.
func trapAttackPower(spell *core.Spell, target *core.Unit) (float64, bool) {
	rap, ap := spell.RangedAttackPower(target), spell.MeleeAttackPower(target)
	return math.Max(rap, ap), rap >= ap
}

// explosiveTrapHitDamage is the explosion's damage before modifiers: the client's roll plus the
// attack power share.
func explosiveTrapHitDamage(sim *core.Simulation, spell *core.Spell, target *core.Unit, damageRange [2]float64) float64 {
	attackPower, _ := trapAttackPower(spell, target)
	return sim.Roll(damageRange[0], damageRange[1]) + trapAttackPowerShare*attackPower
}

// Generator gap: the trap rows hold only the area trigger, and the effect spell's damage is stored
// as the average of the client's range. Explosive Trap rolls 104-135 / 145-193 / 208-265.
var explosiveTrapRange = [4][2]float64{{}, {104, 135}, {145, 193}, {208, 265}}

func (hunter *Hunter) registerExplosiveTrapSpell() {
	rank := spellData.ExplosiveTrap.Highest()
	effect := spellData.ExplosiveTrapEffect.Rank(rank.RankNumber())
	// The dot sits on a persistent area aura, which PeriodicEffect does not answer for.
	tick := effect.Effect(dbcenums.A_PERIODIC_DAMAGE, 0)

	damageRange := explosiveTrapRange[rank.RankNumber()]
	numHits := hunter.Env.ActiveTargetCount()

	hunter.ExplosiveTrap = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    core.SpellSchoolFire,
		DefenseType:    core.DefenseTypeMagic,
		ClassSpellMask: HunterSpellExplosiveTrap,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		MaxRange:       core.MaxMeleeRange,

		ManaCost: core.ManaCostOptions{
			FlatCost: int32(rank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: trapCooldown(hunter, rank),
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Explosive Trap",
				Tag:   "ExplosiveTrap",
			},
			NumberOfTicks: 10,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, tick.Average(core.CharacterLevel))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.ActiveTargetUnits {
					// The Explosive Trap dot only ticks where no Immolation Trap is already burning.
					if !aoeTarget.HasActiveAuraWithTag("ImmolationTrap") {
						dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, effect.TickOutcomeHitRolled(dot))
					}
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := explosiveTrapHitDamage(sim, spell, curTarget, damageRange) * sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, curTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
				curTarget = sim.Environment.NextActiveTargetUnit(curTarget)
			}
			spell.AOEDot().Apply(sim)
		},
	})
}

func (hunter *Hunter) registerImmolationTrapSpell() {
	rank := spellData.ImmolationTrap.Highest()
	effect := spellData.ImmolationTrapTriggered.Rank(rank.RankNumber())
	tick := effect.PeriodicEffect()
	numberOfTicks := int32(effect.Duration() / tick.Period())
	apPerTick := trapAttackPowerShare / float64(numberOfTicks)

	hunter.ImmolationTrap = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    core.SpellSchoolFire,
		DefenseType:    core.DefenseTypeMagic,
		ClassSpellMask: HunterSpellImmolationTrap,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		MaxRange:       core.MaxMeleeRange,

		ManaCost: core.ManaCostOptions{
			FlatCost: int32(rank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: trapCooldown(hunter, rank),
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Immolation Trap",
				Tag:   "ImmolationTrap",
			},
			// 5 ticks 3 sec apart in both clients (13797, 14298-14301); the sim had Season of
			// Discovery's 1.5 sec.
			NumberOfTicks: numberOfTicks,
			TickLength:    tick.Period(),

			// The attack power picked at the snapshot (ranged or melee, whichever is higher then) is the
			// one each tick reads as it stands.
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				attackPower, ranged := trapAttackPower(dot.Spell, target)
				dot.Snapshot(target, tick.Average(core.CharacterLevel)+apPerTick*attackPower)
				dot.SnapshotAttackPowerShare(target, apPerTick, ranged)
			},
			// Every effect rank carries Periodic Can Crit (client 1.60.1.70205), as Explosive Trap's do.
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, effect.TickOutcomeHitRolled(dot))
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
		},
	})
}

// Freezing Trap deals no damage; it is registered so the trap talents and the shared cooldown have
// something to act on, and so an APL can press it.
func (hunter *Hunter) registerFreezingTrapSpell() {
	rank := spellData.FreezingTrap.Rank(1)

	hunter.FreezingTrap = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    core.SpellSchoolFrost,
		DefenseType:    core.DefenseTypeMagic,
		ClassSpellMask: HunterSpellFreezingTrap,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		MaxRange:       core.MaxMeleeRange,

		ManaCost: core.ManaCostOptions{
			FlatCost: int32(rank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: trapCooldown(hunter, rank),
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		},
	})
}
