package core

import (
	"strings"
	"testing"

	"github.com/wowsims/classic/sim/core/proto"
)

// A talent string with more trees than the class has, or digits past its last talent, fills what fits instead of
// crashing.
func TestFillTalentsProtoSkipsWhatDoesNotFit(t *testing.T) {
	talents := &proto.HunterTalents{}
	FillTalentsProto(talents.ProtoReflect(), "1-2-"+strings.Repeat("0", 16)+"515-3", [3]int{16, 16, 18})
	if talents.DeadlyAspects != 1 || talents.HawkEye != 2 || talents.LightningReflexes != 5 || !talents.LaceratingStrikes {
		t.Errorf("filled %v, want Deadly Aspects 1, Hawk Eye 2, Lightning Reflexes 5, Lacerating Strikes", talents)
	}
}

// Each tree keeps its own talents, each capped at its rank count, as the site's talent picker shows the string.
func TestCapTalentsString(t *testing.T) {
	ranks := [3][]int{{5, 1}, {3}, {2, 2}}
	for in, want := range map[string]string{
		"55-3-22":   "51-3-22",
		"19-5":      "11-3",
		"557-34-29": "51-3-22",
		"5-3-2-5":   "5-3-2",
		"x1--":      "01--",
		"":          "",
	} {
		if got := CapTalentsString(in, ranks); got != want {
			t.Errorf("%q read as %q, want %q", in, got, want)
		}
	}
}
