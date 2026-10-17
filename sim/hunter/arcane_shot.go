package hunter

import (
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
)

const ArcaneShotAPCoeff = 0.15

func (hunter *Hunter) getArcaneShotConfig(rank int, timer *core.Timer) core.SpellConfig {
	spellId := [9]int32{0, 3044, 14281, 14282, 14283, 14284, 14285, 14286, 14287}[rank]
	// Forever client data (ranks 2-7 are lower than Classic)
	baseDamage := [9]float64{0, 20, 26, 39, 65, 94, 134, 170, 217}[rank]
	manaCost := [9]float64{0, 25, 35, 50, 80, 105, 135, 160, 190}[rank]
	level := [9]int{0, 6, 12, 20, 28, 36, 44, 52, 60}[rank]

	return core.SpellConfig{
		SpellCode:     SpellCode_HunterArcaneShot,
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolArcane,
		DefenseType:   core.DefenseTypeRanged,
		ProcMask:      core.ProcMaskRangedSpecial,
		Flags:         core.SpellFlagMeleeMetrics | core.SpellFlagAPL | SpellFlagShot,
		CastType:      proto.CastType_CastTypeRanged,
		Rank:          rank,
		RequiredLevel: level,
		MissileSpeed:  24,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second*6 - time.Millisecond*300*time.Duration(hunter.Talents.ImprovedArcaneShot),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget >= core.MinRangedAttackDistance
		},

		CritDamageBonus: hunter.mortalShots(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			ap := spell.RangedAttackPower(target, false)
			damage := baseDamage + ap*ArcaneShotAPCoeff
			result := spell.CalcDamage(sim, target, damage, spell.OutcomeRangedHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	}
}

func (hunter *Hunter) registerArcaneShotSpell(timer *core.Timer) {
	maxRank := 8

	for i := 1; i <= maxRank; i++ {
		config := hunter.getArcaneShotConfig(i, timer)

		if config.RequiredLevel <= int(hunter.Level) {
			hunter.ArcaneShot = hunter.GetOrRegisterSpell(config)
		}
	}
}
