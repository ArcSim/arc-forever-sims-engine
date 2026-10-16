package mage

import (
	"github.com/wowsims/forever/sim/core"
)

// Damage bonus against a target the mage counts as frozen, i.e. one held by Fingers of Frost.
// +200% (was +300%) since the beta build of 2026-10-16 (Blizzard's beta notes, forums t/2360696 post 6).
const IceLanceFrozenMultiplier = 3.0

func (mage *Mage) registerIceLanceSpell() {
	if !mage.Talents.IceLance {
		return
	}

	iceLanceRank := spellData.IceLance.Highest()

	// The client's damage effect carries no spell power coefficient (the row reads 0), as Season of
	// Discovery's reworked Ice Lance (400640) does, whose 2024-12-04 hotfix raised it from .143 to .572.
	// Fitted from the beta logs with each mage's spell power read off their other spells in the same
	// log (client base + coefficient): Mana (report 2727, ~80 spell power from Fire Blast, Frostbolt,
	// Cone of Cold, Frost Nova, Blizzard) hit rank 1 (base 30.1 at its max level 26) for 38.8 on average
	// unfrozen over 52 hits and 38.2 frozen (/4, 32 hits), .10-.11; Toma (report 2687, ~33 spell power,
	// level 20, base 28) for 31.4, .09-.11. .143 would put Mana's hits at 38.8-43.8, but 13 land at 36-37.
	iceLanceCoefficient := 0.1

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: iceLanceRank.ID},
		SpellSchool:    iceLanceRank.SpellSchool(),
		DefenseType:    iceLanceRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagBinary,
		ClassSpellMask: MageSpellIceLance,
		MaxRange:       float64(iceLanceRank.MaxRange),
		MissileSpeed:   float64(iceLanceRank.Speed),

		ManaCost: core.ManaCostOptions{
			FlatCost: int32(iceLanceRank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: iceLanceRank.GCD(),
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: iceLanceCoefficient,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, iceLanceRank.DamageEffect().Roll(sim, core.CharacterLevel), spell.OutcomeMagicHitAndCrit)
			// A bonus on the whole hit rather than the base roll, so spell power is multiplied too.
			if mage.IsTargetFrozen() {
				result.Damage *= IceLanceFrozenMultiplier
			}
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
