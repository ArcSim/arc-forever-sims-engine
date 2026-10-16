package mage

import (
	"fmt"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// Every rank is registered: a rotation can drop to a cheaper rank when mana runs short.
func (mage *Mage) registerArcaneMissilesSpell() {
	spellData.ArcaneMissiles.Each(func(_ int32, rank *spelldata.Spell) { mage.registerArcaneMissilesRank(rank) })
}

func (mage *Mage) registerArcaneMissilesRank(arcaneMissilesRank *spelldata.Spell) {
	missileRank := spellData.ArcaneMissilesTriggered.Rank(arcaneMissilesRank.RankNumber())

	// One missile a second for the channel; the row states the channel's length, not its period.
	tickLength := time.Second
	numTicks := int32(arcaneMissilesRank.Duration() / tickLength)

	arcaneMissilesTickSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: missileRank.ID},
		SpellSchool:    missileRank.SpellSchool(),
		DefenseType:    missileRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagNoOnCastComplete,
		ClassSpellMask: MageSpellArcaneMissilesTick,
		MissileSpeed:   float64(missileRank.Speed),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: missileRank.DamageEffect().Coeff(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, missileRank.DamageEffect().Roll(sim, core.CharacterLevel), spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: arcaneMissilesRank.ID},
		SpellSchool:    arcaneMissilesRank.SpellSchool(),
		DefenseType:    arcaneMissilesRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,
		ClassSpellMask: MageSpellArcaneMissilesCast,
		MaxRange:       float64(arcaneMissilesRank.MaxRange),
		Rank:           arcaneMissilesRank.RankNumber(),

		ManaCost: core.ManaCostOptions{
			FlatCost: int32(arcaneMissilesRank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: arcaneMissilesRank.GCD(),
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("ArcaneMissiles-%d", arcaneMissilesRank.RankNumber()),
			},
			NumberOfTicks: numTicks,
			TickLength:    tickLength,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				arcaneMissilesTickSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})
}
