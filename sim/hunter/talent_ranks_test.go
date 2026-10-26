package hunter

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/wowsims/classic/sim/core/proto"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// The rank table matches the site's talent tree talent by talent, so the sim caps a string where the talent picker
// does.
func TestTalentRanksMatchSiteTree(t *testing.T) {
	data, err := os.ReadFile("../../ui/core/talents/trees/hunter.json")
	if err != nil {
		t.Fatal(err)
	}
	var trees []struct {
		Name    string `json:"name"`
		Talents []struct {
			FieldName string `json:"fieldName"`
			MaxPoints int    `json:"maxPoints"`
		} `json:"talents"`
	}
	if err := json.Unmarshal(data, &trees); err != nil {
		t.Fatal(err)
	}
	if len(trees) != 3 {
		t.Fatalf("%d trees on the site, want 3", len(trees))
	}
	fields := (&proto.HunterTalents{}).ProtoReflect().Descriptor().Fields()
	offset := 0
	for i, tree := range trees {
		if len(tree.Talents) != TalentTreeSizes[i] || len(TalentMaxRanks[i]) != TalentTreeSizes[i] {
			t.Errorf("%s: %d talents on the site, tree size %d, %d ranks", tree.Name, len(tree.Talents), TalentTreeSizes[i], len(TalentMaxRanks[i]))
			continue
		}
		for j, talent := range tree.Talents {
			fd := fields.ByNumber(protowire.Number(offset + j + 1))
			if fd == nil || fd.JSONName() != talent.FieldName {
				t.Errorf("%s talent %d is %s on the site, field %v in the sim", tree.Name, j, talent.FieldName, fd)
				continue
			}
			if TalentMaxRanks[i][j] != talent.MaxPoints {
				t.Errorf("%s: %d ranks, the site has %d", talent.FieldName, TalentMaxRanks[i][j], talent.MaxPoints)
			}
			if fd.Kind() == protoreflect.BoolKind && talent.MaxPoints != 1 {
				t.Errorf("%s is a yes/no field with %d ranks on the site", talent.FieldName, talent.MaxPoints)
			}
		}
		offset += TalentTreeSizes[i]
	}
	if offset != fields.Len() {
		t.Errorf("trees hold %d talents, the sim has %d fields", offset, fields.Len())
	}
}

// An imported Classic Marksmanship string gives Improved Stings and Barrage 5 points in the Forever tree, which has 3
// of each, and puts a 5 on Lone Wolf. The sim reads it as the talent picker shows it (3/3, 3/3, Lone Wolf taken)
// instead of failing on the rank tables.
func TestTalentsCappedAtForeverRanks(t *testing.T) {
	_, hunter := newPatchTestHunter(t, patchTestSetup{talents: "-05451002503051-33400023023"})
	if got := hunter.Talents.ImprovedStings; got != 3 {
		t.Errorf("Improved Stings %d, want 3", got)
	}
	if got := hunter.Talents.Barrage; got != 3 {
		t.Errorf("Barrage %d, want 3", got)
	}
	if !hunter.Talents.LoneWolf {
		t.Error("Lone Wolf not taken, the talent picker shows it 1/1")
	}
	if got := hunter.SerpentSting.DamageMultiplier; got != 1.20 {
		t.Errorf("Serpent Sting damage multiplier %.2f, want 1.20 (Improved Stings 3/3)", got)
	}

	// a 17th Beast Mastery digit is not a Marksmanship talent
	_, hunter = newPatchTestHunter(t, patchTestSetup{talents: "00000000000000003"})
	if got := hunter.Talents.HawkEye; got != 0 {
		t.Errorf("Hawk Eye %d from a Beast Mastery tree one digit too long, want 0", got)
	}
}
