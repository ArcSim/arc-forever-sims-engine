package shared

import (
	"strings"
	"testing"

	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// Client 1.60.1.70170 moved Mystic Mushroom's Increased Spirit to A_MOD_TOTAL_STAT_PERCENTAGE with no stat
// named, and client 70425 names Spirit again in the stat mask alone (misc 0, mask 16). Either way the
// parser reads Spirit, not the Strength misc 0 would mean on its own.
func TestPercentStatRowsNamingNoStatReadAsSpirit(t *testing.T) {
	for id, want := range map[int32]float64{1248751: 1.05} {
		row := spelldata.MustFind(id)
		effect := row.EffectN(1)
		if effect.Aura != dbcenums.A_MOD_TOTAL_STAT_PERCENTAGE {
			t.Fatalf("%s %d is aura %d now; revisit the percent stat parse", row.Name, id, effect.Aura)
		}
		parsed := spelldata.DryRun(row)
		if len(parsed.Applied) == 0 {
			t.Fatalf("%s %d: the parse applied nothing, skipped %v", row.Name, id, parsed.Skipped)
		}
		for _, applied := range parsed.Applied {
			if !strings.Contains(applied.Kind, "Spirit") || strings.Contains(applied.Kind, "Strength") {
				t.Errorf("%s %d: applied %q, want a Spirit multiplier", row.Name, id, applied.Kind)
			}
			if applied.Value < want-1e-9 || applied.Value > want+1e-9 {
				t.Errorf("%s %d: multiplier %v, want %v", row.Name, id, applied.Value, want)
			}
		}
	}
}

// A row that does name a stat is not touched: Spirit Tap is A_MOD_TOTAL_STAT_PERCENTAGE with misc 4 and mask 16.
func TestPercentStatRowsThatNameAStatAreUnchanged(t *testing.T) {
	for _, id := range []int32{15271} {
		row := spelldata.MustFind(id)
		parsed := spelldata.DryRun(row)
		if len(parsed.Applied) == 0 {
			t.Fatalf("%s %d: the parse applied nothing", row.Name, id)
		}
		for _, applied := range parsed.Applied {
			if !strings.Contains(applied.Kind, "Spirit") {
				t.Errorf("%s %d: applied %q, want a Spirit multiplier", row.Name, id, applied.Kind)
			}
		}
	}
}
