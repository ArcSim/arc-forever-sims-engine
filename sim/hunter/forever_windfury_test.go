package hunter

import (
	"testing"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

// Windfury Totem's proc is an extra auto attack (client data effect 19): a white main-hand swing that glances against
// the level 63 dummy as often as the hunter's own swings. Rolled as a special attack, it never glanced.
func TestWindfuryTotemSwingsGlance(t *testing.T) {
	rotation := &proto.APLRotation{}
	if err := protojson.Unmarshal([]byte(`{"type":"TypeAPL","priorityList":[{"action":{"castSpell":{"spellId":{"spellId":14266}}}}]}`), rotation); err != nil {
		t.Fatal(err)
	}
	player := core.WithSpec(&proto.Player{
		Race:               proto.Race_RaceHuman,
		Class:              proto.Class_ClassHunter,
		Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
		TalentsString:      "-00531500113-500230230050122151",
		Rotation:           rotation,
		DistanceFromTarget: 3,
	}, abOptions)
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{WindfuryTotem: true}, &proto.Debuffs{}),
		Encounter:  abEncounter(),
		SimOptions: &proto.SimOptions{Iterations: 100, RandomSeed: 20261022},
	})
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
	}
	glanceRate := func(match func(id *proto.ActionID) bool) (float64, float64) {
		var glances, total int32
		for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
			if !match(a.Id) {
				continue
			}
			for _, tm := range a.Targets {
				glances += tm.Glances
				total += tm.Hits + tm.Crits + tm.Glances + tm.Dodges + tm.Misses + tm.Parries + tm.Blocks
			}
		}
		if total == 0 {
			return 0, 0
		}
		return float64(glances) / float64(total), float64(total)
	}
	wf, wfSwings := glanceRate(func(id *proto.ActionID) bool { return id.GetSpellId() == 10610 && id.Tag == 1 })
	mh, mhSwings := glanceRate(func(id *proto.ActionID) bool {
		return id.GetOtherId() == proto.OtherAction_OtherActionAttack && id.Tag == 1
	})
	t.Logf("Windfury swings glance %.1f%% (%.0f swings), main-hand swings %.1f%% (%.0f swings)", wf*100, wfSwings, mh*100, mhSwings)
	if wfSwings < 300 {
		t.Fatalf("only %.0f Windfury swings", wfSwings)
	}
	if wf < mh-0.05 || wf > mh+0.05 {
		t.Errorf("Windfury swings glance %.1f%% of the time, main-hand swings %.1f%%; want the same rate", wf*100, mh*100)
	}
}
