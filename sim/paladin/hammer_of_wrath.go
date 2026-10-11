package paladin

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

var HammerOfWrathRankMap = spellData.HammerOfWrath

// Hammer of Wrath
// https://www.wowhead.com/forever/spell=24239
//
// Hurls a hammer that strikes an enemy for 498 Holy damage. Only usable on enemies that have 20%
// or less health.
func (paladin *Paladin) registerHammerOfWrath(_ int32, rank *spelldata.Spell) {
	damage := rank.DamageEffect()

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskHammerOfWrath,
		Rank:           rank.RankNumber(),
		MaxRange:       float64(rank.MaxRange),
		MissileSpeed:   float64(rank.Speed),

		ManaCost: rank.ManaCost(),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				// The client's 1s GCD sits at core's floor, so it is named as the floor too or
				// GCDTime clamps it.
				GCDMin:   rank.GCD(),
				GCD:      rank.GCD(),
				CastTime: rank.CastTime(),
			},
			CD: core.Cooldown{
				Timer:    paladin.sharedTimer(&paladin.hammerOfWrathTimer),
				Duration: cooldown(rank),
			},
			// Only a hard cast restarts the swing timer (holdMeleeForCast). Once Instrument of Law makes
			// it instant, the swing timer is left alone.
			ModifyCast: paladin.holdMeleeForCast,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: damage.Coeff(),

		ExtraCastCondition: func(sim *core.Simulation, _ *core.Unit) bool {
			return sim.IsExecutePhase20()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, damage.Roll(sim, core.CharacterLevel), spell.OutcomeRangedHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
