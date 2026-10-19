package hunter

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

// Generator gap: the channel's own rank rows carry only the area-trigger dummy, so the per-tick
// damage has to come from the separate spells Forever fires each tick (1279721, 1279719, 1279715):
// 70/91/112.
var volleyTickDamage = [4]float64{0, 70, 91, 112}

func (hunter *Hunter) registerVolleySpell() {
	rank := spellData.Volley.Highest()
	baseDamage := volleyTickDamage[rank.RankNumber()]

	hunter.Volley = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ClassSpellMask: HunterSpellVolley,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,
		MaxRange:       float64(rank.MaxRange),

		ManaCost: core.ManaCostOptions{
			FlatCost: int32(rank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD(),
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Volley",
			},
			NumberOfTicks: 6,
			TickLength:    time.Second * 1,
			// The tick spells (1279721, 1279719, 1279715) have no coefficient and the channel's dummy
			// effect carries .03, the same placeholder Blizzard and Rain of Fire carry. Bae's in-game
			// tests find no Hunter ability scaling with spell power, so the ticks take none
			// (HunterAbilitiesIgnoreSpellPower, patch 173); with the switch off, Classic's .056 a tick.
			BonusCoefficient: hunterSpellPowerCoefficient(.056),

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, baseDamage)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// Each tick is its own ranged damage spell (client 1.60.1.70205: DefenseType 3, no
				// Cannot Crit), so it rolls ranged hit and crit like Multi-Shot.
				for _, aoeTarget := range sim.Encounter.ActiveTargetUnits {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.Spell.OutcomeRangedHitAndCrit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// A channel is a non-melee cast: it holds the melee swing too, which restarts when it ends.
			hunter.AutoAttacks.DelayRangedUntil(sim, sim.CurrentTime+rank.Duration())
			hunter.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+rank.Duration())
			spell.AOEDot().Apply(sim)
		},
	})
}
