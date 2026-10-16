package shaman

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

var lavaBurstRank = spellData.LavaBurst.Highest()

// Lava Burst is new in Forever. Its second effect is the bonus it gains against a target already
// burning with Flame Shock.
func (shaman *Shaman) registerLavaBurstSpell() {
	if !shaman.Talents.LavaBurst {
		return
	}

	flameShockBonus := 1 + lavaBurstRank.EffectN(2).Percent()

	spell := shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: lavaBurstRank.ID},
		SpellSchool:    lavaBurstRank.SpellSchool(),
		DefenseType:    lavaBurstRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | SpellFlagShamanSpell | SpellFlagFocusable,
		ClassSpellMask: SpellMaskLavaBurst,
		MaxRange:       float64(lavaBurstRank.MaxRange),
		MissileSpeed:   float64(lavaBurstRank.Speed),

		ManaCost: core.ManaCostOptions{
			FlatCost: int32(lavaBurstRank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      lavaBurstRank.GCD(),
				CastTime: lavaBurstRank.CastTime(),
			},
			ModifyCast: shaman.holdMeleeForCast,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: max(lavaBurstRank.Cooldown(), lavaBurstRank.CategoryCooldown()),
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: lavaBurstRank.DamageEffect().Coeff(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if shaman.FlameShock.RelatedDotSpell.Dot(target).IsActive() {
				spell.DamageMultiplier *= flameShockBonus
				defer func() { spell.DamageMultiplier /= flameShockBonus }()
			}

			result := spell.CalcDamage(sim, target, lavaBurstRank.DamageEffect().Roll(sim, core.CharacterLevel), spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
	// Client 70425 added rank 4 (1326047). Rotations saved with an older rank keep casting it.
	spellData.LavaBurst.Each(func(_ int32, rank *spelldata.Spell) {
		if rank.ID != lavaBurstRank.ID {
			spell.RankAliases = append(spell.RankAliases, core.ActionID{SpellID: rank.ID})
		}
	})
}
