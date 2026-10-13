package mage

import (
	"fmt"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// Frostfire Bolt is Forever's trained spell at 40/50/60 (401502, 1237312, 1237313). Client 70009: a 3 sec
// frostfire bolt, .814 spell power, plus a 9 sec dot with no coefficient and a 40% slow the sim has no
// use for. School 20 is fire and frost at once, so fire and frost talents both reach it; the client's
// class masks agree (Fire Power, Piercing Ice, Critical Mass, Ice Shards, Improved Fireball, Presence of
// Mind, and the chill bit Frostbolt and Cone of Cold carry).
func (mage *Mage) registerFrostfireBoltSpell() {
	spellData.FrostfireBolt.Each(func(_ int32, rank *spelldata.Spell) { mage.registerFrostfireBoltRank(rank) })
}

func (mage *Mage) registerFrostfireBoltRank(rank *spelldata.Spell) {
	tick := rank.PeriodicEffect()
	tickLength := tick.Period()

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellFrostfireBolt,
		MaxRange:       float64(rank.MaxRange),
		Rank:           rank.RankNumber(),
		MissileSpeed:   float64(rank.Speed),

		ManaCost: core.ManaCostOptions{
			FlatCost: int32(rank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      rank.GCD(),
				CastTime: rank.CastTime(),
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("FrostfireBoltDoT-%d", rank.RankNumber()),
			},
			NumberOfTicks:    int32(rank.Duration() / tickLength),
			TickLength:       tickLength,
			BonusCoefficient: tick.Coeff(),
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, tick.Average(core.CharacterLevel))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, periodicTickOutcome(rank, dot))
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: rank.DamageEffect().Coeff(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, rank.DamageEffect().Roll(sim, core.CharacterLevel), spell.OutcomeMagicHitAndCrit)
			if mage.procFingersOfFrost != nil {
				mage.procFingersOfFrost(sim, spell, result)
			}
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			})
		},
	})
}
