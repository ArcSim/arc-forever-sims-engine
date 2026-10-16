package feralbear

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

// The damage of the first plain Swipe hit (no crit, block or glance) on a Bear with this much attack power
// added, and the attack power Swipe then read. Both runs seed the same random numbers, so a given cast
// rolls the same outcome in each.
func firstPlainSwipeHit(t *testing.T, extraAttackPower float64) (damage, attackPower float64) {
	t.Helper()
	sim, d := newBearInForm(t)
	d.AddStatDynamic(sim, stats.AttackPower, extraAttackPower)
	target := d.CurrentTarget
	swipe := d.Swipe.Spell

	for range 200 {
		sim.CurrentTime += 2 * time.Second
		d.GCD.Set(sim.CurrentTime)
		d.AddRage(sim, 40, d.RageRefundMetrics)

		before := swipe.SpellMetrics[0]
		if !d.Swipe.Cast(sim, target) {
			t.Fatal("Swipe failed to cast")
		}
		after := swipe.SpellMetrics[0]
		if after.Hits > before.Hits && after.Crits == before.Crits && after.Blocks == before.Blocks &&
			after.Glances == before.Glances && after.Crushes == before.Crushes {
			return after.TotalDamage - before.TotalDamage, swipe.MeleeAttackPower(target)
		}
	}
	t.Fatal("no plain Swipe hit in 200 casts")
	return 0, 0
}

// Swipe gains 5% of attack power, a fit to beta logs (upstream #787; client 1.60.1.70170's patch note
// said 3%). The client rows carry no coefficient, so a hit is the rank's flat damage plus 0.05 attack
// power, times whatever damage modifiers apply. Two casts of the same roll at different attack power give the coefficient without
// naming the modifiers: d1 / d2 = (flat + c AP1) / (flat + c AP2).
func TestSwipeGainsFivePercentOfAttackPower(t *testing.T) {
	flat := spelldata.MustFind(9908).DamageEffect().Average(core.CharacterLevel)
	if flat <= 0 {
		t.Fatalf("Swipe flat damage %v", flat)
	}

	d1, ap1 := firstPlainSwipeHit(t, 0)
	d2, ap2 := firstPlainSwipeHit(t, 1500)

	if ap2-ap1 < 1499 || ap2-ap1 > 1501 {
		t.Fatalf("attack power moved by %v, want 1500", ap2-ap1)
	}
	if d2 <= d1 {
		t.Fatalf("1500 more attack power moved a Swipe hit from %v to %v: it does not scale", d1, d2)
	}

	coefficient := flat * (d2 - d1) / (d1*ap2 - d2*ap1)
	if !core.WithinToleranceFloat64(0.05, coefficient, 0.0001) {
		t.Errorf("Swipe scales %0.5f of attack power, want 0.05 (flat %v, hits %v at %v AP and %v at %v AP)",
			coefficient, flat, d1, ap1, d2, ap2)
	}
}
