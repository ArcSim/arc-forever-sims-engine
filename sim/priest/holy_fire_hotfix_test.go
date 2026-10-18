package priest

import (
	"testing"

	"github.com/wowsims/forever/sim/core/spelldata"
)

// MythicSim patch 165: upstream's 17 October 2026 hotfix cache (#817) cuts Holy Fire rank 1's base
// damage from 58 to 54 (52-66 at its level-24 cap). The other ranks keep the 70450 rows.
func TestHolyFireRank1Hotfix20261017(t *testing.T) {
	if got := spelldata.MustFind(14914).Effects[0].BasePoints; got != 54 {
		t.Errorf("Holy Fire rank 1 base points %v, want 54", got)
	}
	if got := spelldata.MustFind(15261).Effects[0].BasePoints; got == 54 {
		t.Errorf("Holy Fire rank 8 base points %v, moved with rank 1", got)
	}
}
