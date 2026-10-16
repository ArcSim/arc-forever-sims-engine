package mage

import (
	"github.com/wowsims/forever/sim/core"
)

func (mage *Mage) registerArcaneBlastSpell() {
	if !mage.Talents.ArcaneBlast {
		return
	}

	arcaneBlastRank := spellData.ArcaneBlast.Highest()

	mage.ArcaneBlast = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: arcaneBlastRank.ID},
		SpellSchool:    arcaneBlastRank.SpellSchool(),
		DefenseType:    arcaneBlastRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellArcaneBlast,
		MaxRange:       float64(arcaneBlastRank.MaxRange),

		ManaCost: arcaneBlastRank.ManaCost(),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      arcaneBlastRank.GCD(),
				CastTime: arcaneBlastRank.CastTime(),
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: arcaneBlastRank.DamageEffect().Coeff(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, arcaneBlastRank.DamageEffect().Roll(sim, core.CharacterLevel), spell.OutcomeMagicHitAndCrit)
			// The buff runs 15 sec from its first stack (client 70425): a stack does not refresh it.
			if !mage.ArcaneBlastAura.IsActive() {
				mage.ArcaneBlastAura.Activate(sim)
			}
			mage.ArcaneBlastAura.AddStack(sim)
		},
	})
}
