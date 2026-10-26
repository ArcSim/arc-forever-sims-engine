package hunter

import (
	"testing"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

// Goblin Sapper Charge's damage to its user is not damage done: one hit per charge on the one target, and the charge's
// damage stays within its own 450 to 750.
func TestSapperSelfDamageIsNotDamageDone(t *testing.T) {
	rotation := &proto.APLRotation{}
	if err := protojson.Unmarshal([]byte(`{"type":"TypeAPL","priorityList":[{"action":{"castSpell":{"spellId":{"itemId":10646}}}}]}`), rotation); err != nil {
		t.Fatal(err)
	}
	player := core.WithSpec(&proto.Player{
		Race:               proto.Race_RaceHuman,
		Class:              proto.Class_ClassHunter,
		Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
		TalentsString:      "-00531500113-500230230050122151",
		Rotation:           rotation,
		Profession1:        proto.Profession_Engineering,
		Consumes:           &proto.Consumes{SapperExplosive: proto.SapperExplosive_SapperGoblinSapper},
		DistanceFromTarget: 11,
	}, abOptions)
	const iterations = 200
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  abEncounter(),
		SimOptions: &proto.SimOptions{Iterations: iterations, RandomSeed: 20261026},
	})
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
	}
	pm := res.RaidMetrics.Parties[0].Players[0]
	var casts, landed int32
	var damage float64
	for _, a := range pm.Actions {
		for _, tm := range a.Targets {
			if a.Id.GetItemId() == 10646 && a.Id.Tag == 0 {
				casts += tm.Casts
				landed += tm.Hits + tm.Crits
				damage += tm.Damage
			}
		}
	}
	if casts < iterations {
		t.Fatalf("%d sapper charges in %d fights", casts, iterations)
	}
	if landed > casts {
		t.Errorf("%d sapper charges landed %d hits on one target; the extra ones are the charge hitting the hunter", casts, landed)
	}
	if perCharge := damage / float64(casts); perCharge > 750*1.5*1.1 {
		t.Errorf("a sapper charge counts %.0f damage, more than its 450 to 750 with crits and modifiers", perCharge)
	}
}
