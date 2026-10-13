package mage

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// Every rank is registered: a rotation can drop to a cheaper rank when mana runs short.
func (mage *Mage) registerFrostboltSpell() {
	spellData.Frostbolt.Each(func(_ int32, rank *spelldata.Spell) { mage.registerFrostboltRank(rank) })
}

func (mage *Mage) registerFrostboltRank(frostboltRank *spelldata.Spell) {
	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: frostboltRank.ID},
		SpellSchool:    frostboltRank.SpellSchool(),
		DefenseType:    frostboltRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagBinary,
		ClassSpellMask: MageSpellFrostbolt,
		MaxRange:       float64(frostboltRank.MaxRange),
		Rank:           frostboltRank.RankNumber(),
		MissileSpeed:   float64(frostboltRank.Speed),

		ManaCost: core.ManaCostOptions{
			FlatCost: int32(frostboltRank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      frostboltRank.GCD(),
				CastTime: frostboltRank.CastTime(),
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: frostboltRank.DamageEffect().Coeff(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, frostboltRank.DamageEffect().Roll(sim, core.CharacterLevel), spell.OutcomeMagicHitAndCrit)
			if mage.procFingersOfFrost != nil {
				mage.procFingersOfFrost(sim, spell, result)
			}
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
