package protection

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core"
)

// Client 1.60.1.70291 states Last Stand's 30% through A_MOD_INCREASE_HEALTH, where the reference
// asked for A_MOD_MAX_HEALTH and panicked as the aura was gained. The buff still grants 30% of
// maximum health for 20 seconds and takes it back on expiry.
func TestLastStandGrantsThirtyPercentOfMaximumHealth(t *testing.T) {
	sim, war, _ := shieldTestSim(t, DefaultProtectionTalents, 0)
	spell := war.GetSpell(core.ActionID{SpellID: 12975})
	if spell == nil {
		t.Fatal("Last Stand is not registered")
	}
	before := war.MaxHealth()
	spell.Cast(sim, &war.Unit)
	if got, want := war.MaxHealth(), before*1.3; math.Abs(got-want) > 1e-6 {
		t.Fatalf("maximum health with Last Stand = %.3f, want %.3f", got, want)
	}
	war.GetAura("Last Stand").Deactivate(sim)
	if got := war.MaxHealth(); math.Abs(got-before) > 1e-6 {
		t.Fatalf("maximum health after Last Stand = %.3f, want %.3f", got, before)
	}
}
