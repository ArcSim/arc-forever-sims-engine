package warlock

import (
	"github.com/wowsims/forever/sim/core"
)

// AmplifyCurseWholeTick (mythicsim patch 177, ported from wowsims/forever e2f38aef70): Amplify Curse raises
// Bane of Agony's whole tick, the ramping base and the spell power share alike. False restores the base-only
// reading, which left the spell power share unamplified.
var AmplifyCurseWholeTick = true

// Forever renamed Curse of Agony to Bane of Agony and moved it onto the bane slot, so it holds the
// target beside a curse. It ramps: the snapshot pays half the tick, and the other half is added
// back every four ticks, leaving the last four at 150%. Amplify Curse adds half again to the whole
// tick (AmplifyCurseWholeTick) and is spent on the snapshot. The spell power share is read at each
// tick (patch 22), so the multiplier is kept per target and applied to the tick as it lands.
func (warlock *Warlock) registerCurseOfAgony() {
	rank := spellData.BaneOfAgony.Highest()
	tick := rank.PeriodicEffect()
	amplify := 1 + spellData.AmplifyCurse.EffectAt(1).FractionAt(1)

	rampStep := 0.0
	wholeTickAmp := map[int32]float64{}

	warlock.CurseOfAgony = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellCurseOfAgony,
		MaxRange:       float64(rank.MaxRange),

		ManaCost: core.ManaCostOptions{FlatCost: int32(rank.Cost())},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD(),
			},
		},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,
		BonusCoefficient:         tick.Coeff(),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Bane of Agony",
				Tag:   "Affliction",
			},
			NumberOfTicks:    int32(rank.Duration() / tick.Period()),
			TickLength:       tick.Period(),
			BonusCoefficient: tick.Coeff(),

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				base := tick.Average(core.CharacterLevel)
				wholeTickAmp[target.UnitIndex] = 1
				if warlock.AmplifyCurseAura.IsActive() {
					if AmplifyCurseWholeTick {
						wholeTickAmp[target.UnitIndex] = amplify
					} else {
						base *= amplify
					}
					warlock.AmplifyCurseAura.Deactivate(sim)
				}

				rampStep = base * 0.5
				dot.Snapshot(target, rampStep)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.CalcSnapshotDamage(sim, target, periodicTickOutcome(rank, dot))
				result.Damage *= wholeTickAmp[target.UnitIndex]
				dot.Spell.DealPeriodicDamage(sim, result)
				if dot.TickCount()%4 == 0 {
					dot.SnapshotBaseDamage += rampStep
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				dot := spell.Dot(target)
				warlock.takeBaneSlot(sim, target, dot.Aura)
				dot.Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			// The ramp averages out to the plain tick over the full duration.
			dot := spell.Dot(target)
			if useSnapshot {
				result := dot.CalcSnapshotDamage(sim, target, spell.OutcomeExpectedMagicAlwaysHit)
				if amp, ok := wholeTickAmp[target.UnitIndex]; ok {
					result.Damage *= amp
				}
				return result
			}
			return spell.CalcPeriodicDamage(sim, target, tick.Average(core.CharacterLevel), spell.OutcomeExpectedMagicAlwaysHit)
		},
	})
}
