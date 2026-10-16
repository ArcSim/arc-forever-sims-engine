package mage

import (
	"math"
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

// Fire Power (client 70425): 3% a Fire effect on the target at 3 points, up to four effects. Frostfire
// Bolt's DoT mixes schools and does not count; Frostfire Bolt itself is a Fire spell and takes it.
func TestFirePowerCountsFireEffectsOnTheTarget(t *testing.T) {
	sim, mage := frostfireSim(frostfireHybridTalents)
	if mage.Talents.FirePower != 3 || mage.Talents.Pyroblast != true || mage.Talents.ImprovedScorch == 0 {
		t.Fatalf("frostfireHybridTalents no longer takes Fire Power 3, Pyroblast and Improved Scorch: %v", mage.Talents)
	}
	target := mage.CurrentTarget
	table := mage.AttackTables[target.UnitIndex]
	fireball := mage.GetSpell(core.ActionID{SpellID: spellData.Fireball.Highest().ID})
	frostfire := mage.GetSpell(core.ActionID{SpellID: 1237313})
	frostbolt := mage.GetSpell(core.ActionID{SpellID: spellData.Frostbolt.Highest().ID})
	want := func(step string, fire float64) {
		t.Helper()
		if got := fireball.TargetDamageMultiplier(sim, table, false); math.Abs(got-fire) > 1e-9 {
			t.Errorf("%s: Fireball takes %v, want %v", step, got, fire)
		}
		if got := frostfire.TargetDamageMultiplier(sim, table, false); math.Abs(got-fire) > 1e-9 {
			t.Errorf("%s: Frostfire Bolt takes %v, want %v", step, got, fire)
		}
		if got := frostbolt.TargetDamageMultiplier(sim, table, false); got != 1 {
			t.Errorf("%s: Frostbolt takes %v, want 1", step, got)
		}
	}
	want("nothing on the target", 1)
	frostfire.Dot(target).Apply(sim)
	want("Frostfire Bolt's DoT", 1)
	mage.GetSpell(core.ActionID{SpellID: spellData.Pyroblast.Highest().ID}).Dot(target).Apply(sim)
	want("Pyroblast", 1.03)
	fireball.Dot(target).Apply(sim)
	want("Pyroblast and Fireball", 1.06)
	mage.Ignite.Dot(target).Apply(sim)
	want("three DoTs", 1.09)
	mage.ImprovedScorchAura.Activate(sim)
	want("three DoTs and Fire Vulnerability", 1.12)
	mage.GetSpell(core.ActionID{SpellID: FlameStrikeRankMap.Highest().ID}).AOEDot().Apply(sim)
	want("five effects, capped at four", 1.12)
}

// Missile Barrage (client 70425) has two ranks, 20% and 40% on Arcane Blast; an Arcane Missiles missile
// rolls a quarter of it and Arcane Wrath half.
func TestMissileBarrageRanks(t *testing.T) {
	for _, row := range []struct {
		talents string
		chance  float64
	}{
		{"00010000001", 0.20},
		{"00020000001", 0.40},
	} {
		sim, mage := newThreeTargetMage(row.talents)
		if mage.MissileBarrageAura == nil || mage.ArcaneBlast == nil {
			t.Fatalf("%s: Missile Barrage or Arcane Blast missing", row.talents)
		}
		const trials = 4000
		procs := 0
		for i := 0; i < trials; i++ {
			mage.ArcaneBlast.CalcAndDealDamage(sim, mage.CurrentTarget, 1, mage.ArcaneBlast.OutcomeAlwaysHit)
			if mage.MissileBarrageAura.IsActive() {
				procs++
				mage.MissileBarrageAura.Deactivate(sim)
			}
		}
		if got := float64(procs) / trials; math.Abs(got-row.chance) > 0.025 {
			t.Errorf("%s: Arcane Blast procced %.3f, want %.2f", row.talents, got, row.chance)
		}
	}
}

// Arcane Wrath (client 70425): a 3.5 sec, 6 sec cooldown Arcane bolt whose spell damage includes 5% of
// the mage's maximum mana.
func TestArcaneWrathScalesWithMaximumMana(t *testing.T) {
	sim, mage := newThreeTargetMage("000000000000001")
	wrath := mage.GetSpell(core.ActionID{SpellID: spellData.ArcaneWrath.Highest().ID})
	if wrath == nil {
		t.Fatal("Arcane Wrath not registered")
	}
	if wrath.DefaultCast.CastTime != 3500*time.Millisecond || wrath.CD.Duration != 6*time.Second || !wrath.Matches(MageSpellArcaneWrath) {
		t.Fatalf("Arcane Wrath: cast %v, cooldown %v", wrath.DefaultCast.CastTime, wrath.CD.Duration)
	}
	wrath.BonusCritPercent -= 1000
	wrath.Flags |= core.SpellFlagIgnoreResists
	mean := func() float64 {
		metrics := &wrath.SpellMetrics[mage.CurrentTarget.UnitIndex]
		damage, hits := metrics.TotalDamage, metrics.Hits
		for i := 0; i < 400; i++ {
			wrath.CD.Reset()
			wrath.SkipCastAndApplyEffects(sim, mage.CurrentTarget)
			for j := 0; j < 20 && sim.Step() == false; j++ {
			}
		}
		return (metrics.TotalDamage - damage) / float64(metrics.Hits-hits)
	}
	before := mean()
	mage.AddStatDynamic(sim, stats.Mana, 10000)
	after := mean()
	want := spellData.ArcaneWrath.Highest().DamageEffect().Coeff() * 0.05 * 10000
	if got := after - before; math.Abs(got-want) > 0.1*want {
		t.Errorf("10000 more maximum mana added %.1f damage a hit, want about %.1f", got, want)
	}
}

// Winter's Chill (client 70425) empowers Cone of Cold, 5% a stack, and no longer adds crit to Frostbolt.
func TestWintersChillEmpowersConeOfCold(t *testing.T) {
	sim, mage := newThreeTargetMage(FrostTalents)
	if mage.WintersChillAura == nil {
		t.Fatal("FrostTalents no longer take Winter's Chill; pick a build that does")
	}
	cone := mage.GetSpell(core.ActionID{SpellID: spellData.ConeOfCold.Highest().ID})
	frostbolt := mage.GetSpell(core.ActionID{SpellID: spellData.Frostbolt.Highest().ID})
	coneBefore, crit := cone.DamageMultiplierAdditive, frostbolt.BonusCritPercent
	mage.WintersChillAura.Activate(sim)
	mage.WintersChillAura.SetStacks(sim, 3)
	if got := cone.DamageMultiplierAdditive - coneBefore; math.Abs(got-0.15) > 1e-9 {
		t.Errorf("three stacks add %v to Cone of Cold, want 0.15", got)
	}
	if frostbolt.BonusCritPercent != crit {
		t.Errorf("Winter's Chill changed Frostbolt's crit from %v to %v", crit, frostbolt.BonusCritPercent)
	}
}
