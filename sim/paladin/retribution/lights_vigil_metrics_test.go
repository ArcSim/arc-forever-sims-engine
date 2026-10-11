package retribution

import (
	"fmt"
	"os"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

// The shared Shockadin report combined Vigil's cost and refund under one
// action ID. Its net-negative row hid the actual refund from mana sources.
func TestLightsVigilReportsCostAndRefundSeparately(t *testing.T) {
	for rank, id := range []int32{1310911, 1311590, 1311595} {
		t.Run(fmt.Sprintf("rank-%d", rank+1), func(t *testing.T) {
			testLightsVigilMetrics(t, id)
		})
	}
}

func testLightsVigilMetrics(t *testing.T, spellID int32) {
	t.Helper()
	raw, err := os.ReadFile("testdata/lights_vigil_request.json")
	if err != nil {
		t.Fatal(err)
	}
	request := &proto.RaidSimRequest{}
	if err := protojson.Unmarshal(raw, request); err != nil {
		t.Fatal(err)
	}
	// The action regression also runs without the optional embedded item database.
	if !core.WITH_DB {
		request.Raid.Parties[0].Players[0].Equipment = &proto.EquipmentSpec{}
	}
	// Isolate the refund from unrelated cooldown and aura conditions in the report.
	rotation := &proto.APLRotation{}
	if err := protojson.Unmarshal([]byte(fmt.Sprintf(`{"type":"TypeAPL","priorityList":[{"action":{"castSpell":{"spellId":{"spellId":%d}}}},{"action":{"castSpell":{"spellId":{"spellId":20930}}}}]}`, spellID)), rotation); err != nil {
		t.Fatal(err)
	}
	request.Raid.Parties[0].Players[0].Rotation = rotation
	result := core.RunRaidSim(request)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	p := result.RaidMetrics.Parties[0].Players[0]
	var spent, refunded *proto.ResourceMetrics
	for _, r := range p.Resources {
		if r.Type != proto.ResourceType_ResourceTypeMana || r.Id.GetSpellId() != spellID {
			continue
		}
		if r.Id.Tag == 0 {
			spent = r
		}
		if r.Id.Tag == 1 {
			refunded = r
		}
	}
	if spent == nil || refunded == nil {
		t.Fatalf("missing separate cost/refund: cost=%v refund=%v", spent, refunded)
	}
	if spent.Gain >= 0 || spent.ActualGain != spent.Gain {
		t.Fatalf("invalid spend row: %v", spent)
	}
	if refunded.Gain <= 0 || refunded.ActualGain <= 0 || refunded.ActualGain > refunded.Gain {
		t.Fatalf("invalid refund row: %v", refunded)
	}
	if refunded.Events > spent.Events {
		t.Fatal("more refunds than applications")
	}
	if refunded.Gain > -spent.Gain*0.750001 {
		t.Fatal("refund exceeds 75% of the actual mana paid")
	}
}
