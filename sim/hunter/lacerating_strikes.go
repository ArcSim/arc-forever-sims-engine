package hunter

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

// The bleed lasts 21 sec and carries 40% of the Mongoose Bite that applied it. The tooltip gives no
// tick interval, so it ticks every 3 sec like every other bleed. A bite on a running bleed rolls the
// unticked damage into the new one and keeps the tick timer, as Forever's Deep Wounds does (Zirene on
// forever-bugs #234, patch 89): Mongoose Bite comes round every 5 to 13 sec, so a bleed that replaced
// the last one dealt under half of its 40%. Its ticks crit where the client marks
// Periodic Can Crit (1310536, set in build 70009). It reports under that bleed's own id, so results
// name it Lacerating Strikes rather than a second Mongoose Bite.
func (hunter *Hunter) registerLaceratingStrikesSpell() {
	if !hunter.Talents.LaceratingStrikes {
		return
	}

	hunter.LaceratingStrikes = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellData.LaceratingStrikesTriggered.Highest().ID},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ClassSpellMask: HunterSpellLaceratingStrikes,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Lacerating Strikes" + hunter.Label,
			},
			NumberOfTicks: 7,
			TickLength:    time.Second * 3,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, spellData.LaceratingStrikesTriggered.Highest().TickOutcomeHitRolled(dot))
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if LaceratingStrikesRollsOver {
				spell.Dot(target).ApplyKeepingTickTimer(sim)
			} else {
				spell.Dot(target).Apply(sim)
			}
		},
	})
}

// The bleed is written after the cast: the expiry of a running bleed zeroes its snapshot, so a value
// written before the cast was lost (patch 87). What the running bleed still owes is read first and
// spread with the new bite's 40% over the ticks left (LaceratingStrikesRollsOver).
func (hunter *Hunter) procLaceratingStrikes(sim *core.Simulation, result *core.SpellResult) {
	dot := hunter.LaceratingStrikes.Dot(result.Target)
	owed := 0.0
	if LaceratingStrikesRollsOver {
		owed = dot.OutstandingDmg()
	}
	hunter.LaceratingStrikes.Cast(sim, result.Target)

	dot.SnapshotBaseDamage = (owed + result.Damage*0.4) / float64(dot.RemainingTicks())
	dot.SnapshotAttackerMultiplier = 1
}

// LaceratingStrikesRollsOver carries a running bleed's unticked damage into the next one (patch 132).
// False replaces the running bleed and restarts its tick timer, the behaviour before the patch.
var LaceratingStrikesRollsOver = true
