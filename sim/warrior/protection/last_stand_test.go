package protection

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
)

// Client 1.60.1.70291 moved Last Stand's buff (12976) from A_MOD_MAX_HEALTH to A_MOD_INCREASE_HEALTH,
// still 30; the talent tooltip reads "30% of your maximum health". Casting it panicked on the old aura.
func TestLastStandGrantsThirtyPercentOfMaximumHealth(t *testing.T) {
	sim, war, _ := shieldTestSim(t, DefaultProtectionTalents, 0)
	lastStand := war.GetSpell(core.ActionID{SpellID: 12975})
	if lastStand == nil {
		t.Fatal("the warrior has no Last Stand")
	}

	before := war.MaxHealth()
	lastStand.Cast(sim, &war.Unit)
	if got, want := war.MaxHealth(), before*1.3; !core.WithinToleranceFloat64(want, got, 0.01) {
		t.Fatalf("max health under Last Stand: got %.1f, want %.1f", got, want)
	}
}
