package mage

import (
	"github.com/wowsims/forever/sim/core"
)

func (mage *Mage) registerPyroblastSpell() {
	if !mage.Talents.Pyroblast {
		return
	}

	pyroblastRank := spellData.Pyroblast.Highest()
	pyroblastTick := pyroblastRank.PeriodicEffect()
	tickLength := pyroblastTick.Period()

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: pyroblastRank.ID},
		SpellSchool:    pyroblastRank.SpellSchool(),
		DefenseType:    pyroblastRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellPyroblast,
		MaxRange:       float64(pyroblastRank.MaxRange),
		MissileSpeed:   float64(pyroblastRank.Speed),

		ManaCost: core.ManaCostOptions{
			FlatCost: int32(pyroblastRank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      pyroblastRank.GCD(),
				CastTime: pyroblastRank.CastTime(),
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "PyroblastDoT",
			},
			NumberOfTicks:    int32(pyroblastRank.Duration() / tickLength),
			TickLength:       tickLength,
			BonusCoefficient: pyroblastTick.Coeff(),
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, pyroblastTick.Average(core.CharacterLevel))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, periodicTickOutcome(pyroblastRank, dot))
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: pyroblastRank.DamageEffect().Coeff(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, pyroblastRank.DamageEffect().Roll(sim, core.CharacterLevel), spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			})
		},
	})
}
