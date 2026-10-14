package druid

import (
	"github.com/wowsims/forever/sim/core"
)

var swipeRank = spellData.Swipe.Highest()

// Swipe scales with attack power through a server script the client rows don't carry: dev
// AidanZMoon removed a Season of Discovery early exit that skipped it (forever-bugs #125, 'Will be
// in next build', 26 Sep). The 0.05 is fitted to beta logs, not SoD's 0.1: Thebig (report 2756,
// rank 2) hit the same mobs for 32.8 without and 34.4 with Blessing of Might rank 3 (+40 AP), 126
// non-crits, and his Bear swings there put armor near 0.69, so ~0.047 an AP point at Feral Instinct 3/3.
const swipeAttackPowerCoeff = 0.05

func (druid *Druid) registerSwipeBearSpell() {
	druid.Swipe = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: swipeRank.ID},
		SpellSchool:    swipeRank.SpellSchool(),
		DefenseType:    swipeRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: DruidSpellSwipe,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   int32(swipeRank.Cost()),
			Refund: swipeRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: swipeRank.GCD(),
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		// Season of Discovery's "Modifies Threat +101%", which the client does not carry.
		ThreatMultiplier: 2,
		MaxRange:         core.MaxMeleeRange,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			numHits := min(3, len(druid.Env.Encounter.AllTargetUnits))
			for i := 0; i < numHits; i++ {
				aoeTarget := druid.Env.Encounter.AllTargetUnits[i]
				baseDamage := swipeRank.DamageEffect().Average(core.CharacterLevel) + swipeAttackPowerCoeff*spell.MeleeAttackPower(aoeTarget)
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			}
		},
	})
}
