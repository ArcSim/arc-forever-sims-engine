package hunter

import (
	"os"
	"strconv"
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/stats"
)

// ImmolationTrapAPCoefficient: the share of attack power the burn adds over its 15 s.
const ImmolationTrapAPCoefficient = 0.132

// trapAPBonus: T.N.T. (428717) adds a share of the higher of melee and ranged AP to Immolation Trap's total burn and
// to Explosive Trap's first hit. 0.132 from live Immolation Trap logs (the client's 25% is not what is live).
// AB_TRAP_AP sets another coefficient for comparisons.
func trapAPBonus(spell *core.Spell, target *core.Unit) float64 {
	coef := ImmolationTrapAPCoefficient
	if v, err := strconv.ParseFloat(os.Getenv("AB_TRAP_AP"), 64); err == nil {
		coef = v
	}
	return coef * max(spell.RangedAttackPower(target, false), spell.MeleeAttackPower(target))
}

// trapReach: how far from the target a trap can be dropped and still catch it. In game it lands about a yard in front
// of you: from 6 yd it hits, from 8 yd it does not (7 yd untested).
func (hunter *Hunter) trapReach() float64 {
	// AB_TRAP_REACH sets the reach in yards for comparisons
	if v, err := strconv.ParseFloat(os.Getenv("AB_TRAP_REACH"), 64); err == nil {
		return v
	}
	return core.MaxMeleeAttackDistance + 1
}

// trapArmDelay: time from the cast to the burn landing on the target. In Forever logs it lands 0.8 to 1.6 s after the
// cast (65 traps in two reports, median 1.2 s in each), then ticks every 3 s; the sim takes the 1.2 s middle.
// AB_TRAP_ARM sets it in seconds for comparisons.
func trapArmDelay() time.Duration {
	if v, err := strconv.ParseFloat(os.Getenv("AB_TRAP_ARM"), 64); err == nil {
		return time.Duration(v * float64(time.Second))
	}
	return 1200 * time.Millisecond
}

func (hunter *Hunter) getImmolationTrapConfig(rank int, timer *core.Timer) core.SpellConfig {
	// the Forever trap spells (the 4095xx ids are not in the Forever data)
	spellId := [6]int32{0, 13795, 14302, 14303, 14304, 14305}[rank]
	dotDamage := [6]float64{0, 105, 215, 340, 510, 690}[rank]
	manaCost := [6]float64{0, 50, 90, 135, 190, 245}[rank]
	level := [6]int{0, 16, 26, 36, 46, 56}[rank]

	return core.SpellConfig{
		SpellCode:   SpellCode_HunterImmolationTrap,
		ActionID:    core.ActionID{SpellID: spellId},
		SpellSchool: core.SpellSchoolFire,
		// a magic spell in the Forever data: its ticks crit for 1.5x (logs: 65 -> 98)
		DefenseType: core.DefenseTypeMagic,
		// a magic effect in Forever: it does not trigger melee procs (Windfury, Hand of Justice, Chili)
		ProcMask: core.ProcMaskSpellDamage,
		// a real cast on the GCD, not passive: the results count its casts, so it gets damage per cast and per execute time
		Flags:         core.SpellFlagAPL | SpellFlagTrap,
		Rank:          rank,
		RequiredLevel: level,
		MissileSpeed:  24,
		// a trap at your feet only catches a target within its reach (it does nothing from further out)
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceNow(sim) <= hunter.trapReach()
		},

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second * 30,
			},
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "ImmolationTrap" + hunter.Label + strconv.Itoa(rank),
				Tag:   "ImmolationTrap",
			},
			NumberOfTicks: 5,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				tickDamage := (dotDamage + trapAPBonus(dot.Spell, target)) / float64(dot.NumberOfTicks)
				dot.Snapshot(target, tickDamage, isRollover)
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickPhysicalCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if hunter.DistanceNow(sim) > hunter.trapReach() {
				return
			}
			// Traps gain no benefit from hit bonuses except for the Trap Mastery talent, since this is a unique interaction this is my workaround
			spellHit := spell.Unit.GetStat(stats.SpellHit) + target.PseudoStats.BonusSpellHitRatingTaken
			spell.Unit.AddStatDynamic(sim, stats.SpellHit, spellHit*-1)
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			spell.Unit.AddStatDynamic(sim, stats.SpellHit, spellHit)
			core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt: sim.CurrentTime + max(spell.TravelTime(), trapArmDelay()),
				OnAction: func(s *core.Simulation) {
					spell.DealOutcome(sim, result)
					if result.Landed() {
						spell.Dot(target).Apply(sim)
					}
				},
			})
		},
	}
}

func (hunter *Hunter) registerImmolationTrapSpell(timer *core.Timer) {
	maxRank := 5
	for i := 1; i <= maxRank; i++ {
		config := hunter.getImmolationTrapConfig(i, timer)

		if config.RequiredLevel <= int(hunter.Level) {
			hunter.ImmolationTrap = hunter.GetOrRegisterSpell(config)
		}
	}
}
