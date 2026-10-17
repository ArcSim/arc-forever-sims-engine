package hunter

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

const RaptorStrikeRanks = 8

var RaptorStrikeSpellId = [RaptorStrikeRanks + 1]int32{0, 2973, 14260, 14261, 14262, 14263, 14264, 14265, 14266}
var RaptorStrikeSpellIdMeleeSpecialist = [RaptorStrikeRanks + 1]int32{0, 415335, 415336, 415337, 415338, 415340, 415341, 415342, 415343}
// Forever client data (the 34..140 values are the Season of Discovery rune version, not Forever's)
var RaptorStrikeBaseDamage = [RaptorStrikeRanks + 1]float64{0, 5, 11, 21, 30, 35, 40, 55, 70}
var RaptorStrikeManaCost = [RaptorStrikeRanks + 1]float64{0, 15, 25, 35, 45, 55, 70, 80, 100}
var RaptorStrikeLevel = [RaptorStrikeRanks + 1]int{0, 1, 8, 16, 24, 32, 40, 48, 56}

func (hunter *Hunter) TryRaptorStrike(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
	if hunter.curQueueAura == nil || !hunter.curQueueAura.IsActive() || hunter.RaptorStrikeHit == nil {
		return mhSwingSpell
	}
	if !hunter.RaptorStrikeHit.CanCast(sim, hunter.CurrentTarget) {
		hunter.curQueueAura.Deactivate(sim)
		return mhSwingSpell
	}
	return hunter.RaptorStrikeHit
}

func (hunter *Hunter) unqueueRaptorStrike(sim *core.Simulation) {
	if hunter.curQueueAura != nil {
		hunter.curQueueAura.Deactivate(sim)
	}
}

func (hunter *Hunter) registerRaptorStrikeSpell() {
	rank := map[int32]int{
		25: 4,
		40: 6,
		50: 7,
		60: 8,
	}[hunter.Level]
	if rank == 0 {
		rank = 8
	}

	spellID := RaptorStrikeSpellId[rank]
	baseDamage := RaptorStrikeBaseDamage[rank]
	manaCost := RaptorStrikeManaCost[rank]

	cd := hunter.NewTimer()

	hunter.RaptorStrikeHit = hunter.RegisterSpell(core.SpellConfig{
		SpellCode:   SpellCode_HunterRaptorStrikeHit,
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeMHAuto,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},

		DamageMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := baseDamage + hunter.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			if hunter.curQueueAura != nil {
				hunter.curQueueAura.Deactivate(sim)
			}
			cd.Set(sim.CurrentTime + 6*time.Second)
			if result.Landed() {
				hunter.AutoAttacks.RestartRangedSwing(sim)
			}
		},
	})

	hunter.raptorQueueAura = hunter.RegisterAura(core.Aura{
		Label:    "Raptor Strike Queued",
		ActionID: core.ActionID{SpellID: spellID}.WithTag(2),
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.curQueueAura = aura
			hunter.curQueuedAutoSpell = hunter.RaptorStrikeHit
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.curQueueAura = nil
			hunter.curQueuedAutoSpell = nil
		},
	})

	// any ranged attack (a shot, a sting, the Auto Shot) switches you out of melee mode, which unqueues Raptor Strike
	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: "Raptor Strike Unqueued By Ranged",
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.DefenseType == core.DefenseTypeRanged {
				hunter.unqueueRaptorStrike(sim)
			}
		},
	}))

	hunter.RaptorStrike = hunter.RegisterSpell(core.SpellConfig{
		SpellCode: SpellCode_HunterRaptorStrike,
		ActionID:  core.ActionID{SpellID: spellID}.WithTag(1),
		Flags:     core.SpellFlagAPL | core.SpellFlagCastTimeNoGCD,
		Cast:      core.CastConfig{},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Forever: it queues at any range and waits for the next melee swing
			return hunter.curQueueAura == nil &&
				cd.IsReady(sim) &&
				hunter.CurrentMana() >= hunter.RaptorStrikeHit.Cost.GetCurrentCost()
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			hunter.raptorQueueAura.Activate(sim)
		},
	})
	// the hit starts the 6 s cooldown on cd; pointing the castable spell at it lets rotations read Raptor Strike's
	// cooldown (a cast only happens when cd is ready, so its zero-length Use changes nothing)
	hunter.RaptorStrike.CD.Timer = cd
}
