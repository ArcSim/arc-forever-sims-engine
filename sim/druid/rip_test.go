package druid

import "testing"

// Wowhead Forever 9896: 243 damage over 12 sec at 1 point, 855 at 5, before attack power.
func TestRipReadsItsRow(t *testing.T) {
	for cp, want := range map[float64]float64{1: 243, 5: 855} {
		if got := 6 * ripTickDamage(cp, 0); got != want {
			t.Errorf("Rip at %v points deals %v, want %v", cp, got, want)
		}
	}
}

// Wowhead Forever 9904: 61 on the hit and 102 over 9 sec before attack power; 5% of it on each (#769).
func TestRakeAddsFivePercentOfAttackPower(t *testing.T) {
	if hit, tick := rakeHitDamage(0), 3*rakeTickDamage(0); hit != 61 || tick != 102 {
		t.Errorf("Rake rank 4 deals %v + %v, want 61 + 102", hit, tick)
	}
	if got := rakeTickDamage(1000) - rakeTickDamage(0); got != 50 {
		t.Errorf("1000 attack power adds %v a tick, want 50", got)
	}
}
