package mage

import (
	"github.com/wowsims/forever/sim/core"
)

// Arcane Wrath (1326079, ranks to 1326082) is an Arcane talent of client 70425 (16 October 2026 beta):
// a 3.5 sec Arcane bolt on a 6 sec cooldown that "gains additional Spell Damage equal to 5% of your
// maximum Mana". The row's dummy effect holds the 5; the bonus is spell damage, so it scales with the
// damage effect's coefficient like the mage's own spell damage does, and is read at cast from the
// mage's maximum mana at that moment.
func (mage *Mage) registerArcaneWrathSpell() {
	if !mage.Talents.ArcaneWrath {
		return
	}

	rank := spellData.ArcaneWrath.Highest()
	damage := rank.DamageEffect()
	coefficient := damage.Coeff()
	manaShare := rank.EffectN(2).BasePoints / 100

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellArcaneWrath,
		MaxRange:       float64(rank.MaxRange),
		MissileSpeed:   float64(rank.Speed),

		ManaCost: core.ManaCostOptions{
			FlatCost: int32(rank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      rank.GCD(),
				CastTime: rank.CastTime(),
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: rank.Cooldown(),
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: coefficient,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := damage.Roll(sim, core.CharacterLevel) + coefficient*manaShare*mage.MaxMana()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
