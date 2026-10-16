package druid

import (
	"github.com/wowsims/forever/sim/core"
)

var rakeRank = spellData.Rake.Highest()
var rakeTick = rakeRank.PeriodicEffect()

// The client states a flat hit and a flat tick with no BonusCoefficientFromAP, but Blizzard's Druid
// deep dive (30 Sep) says Rake "now gains increased damage from Attack Power". Beta logs put it at
// about 5% of attack power on the hit and on each tick (issue #769): Moondre (foreverlogs 2717,
// ~310 AP) and Axl (2771, ~405 AP, both read off their Rip ticks) land rank 1 at 43/36 and 51/43
// against the client's 23/16 before Savage Fury. The hit ignores armor too: each druid's hit is the
// same on every mob while their white swings differ by armor. A tick reads its share at the tick, like Rip.
const rakeAttackPowerShare = 0.05

func rakeHitDamage(attackPower float64) float64 {
	return rakeRank.DamageEffect().Average(core.CharacterLevel) + rakeAttackPowerShare*attackPower
}

func rakeTickDamage(attackPower float64) float64 {
	return rakeTick.Average(core.CharacterLevel) + rakeAttackPowerShare*attackPower
}

func (druid *Druid) registerRakeSpell() {
	druid.Rake = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rakeRank.ID},
		SpellSchool:    rakeRank.SpellSchool(),
		DefenseType:    rakeRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: DruidSpellRake,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagIgnoreResists,
		Rank:           rakeRank.RankNumber(),

		EnergyCost: core.EnergyCostOptions{
			Cost:   int32(rakeRank.Cost()),
			Refund: rakeRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rakeRank.GCD(),
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		MaxRange:         core.MaxMeleeRange,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Rake",
				Duration: rakeRank.Duration(),
			},
			NumberOfTicks: int32(rakeRank.Duration() / rakeTick.Period()),
			TickLength:    rakeTick.Period(),

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.SnapshotPhysical(target, rakeTickDamage(dot.Spell.MeleeAttackPower(target)))
				dot.SnapshotAttackPowerShare(target, rakeAttackPowerShare, false)
				druid.UpdateBleedPower(druid.Rake, sim, target, true, true)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, periodicTickOutcome(rakeRank, dot))
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealDamage(sim, target, rakeHitDamage(spell.MeleeAttackPower(target)), spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				druid.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				spell.Dot(target).Apply(sim)
			} else {
				spell.IssueRefund(sim)
			}
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, dot.OutcomeTick)
			}
			ticks := spell.CalcPeriodicDamage(sim, target, rakeTickDamage(spell.MeleeAttackPower(target)), spell.OutcomeExpectedMagicAlwaysHit)
			attackTable := spell.Unit.AttackTables[target.UnitIndex]
			critChance := spell.PhysicalCritChance(attackTable)
			ticks.Damage *= 1 + critChance*(spell.CritDamageMultiplier(attackTable)-1)
			return ticks
		},
	})

	druid.Rake.ShortName = "Rake"
}

func (druid *Druid) CurrentRakeCost() float64 {
	return druid.Rake.Cost.GetCurrentCost()
}
