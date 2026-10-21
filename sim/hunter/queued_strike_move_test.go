package hunter

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

// A weave steps out of melee as its strike goes off. Here the hunter shoots from 11 yd, walks in during the
// shot's GCD, and the Strider Kick it presses on the way (in a strict sequence, as weaves put Immolation Trap) is
// queued onto the end of that GCD (1.5 s). The walk back out has to start as the queued kick goes off, not when
// the kick's own GCD ends 1.5 s later; with the switch off it waits for that GCD.
func TestMoveFollowsAQueuedStrike(t *testing.T) {
	const rotationJSON = `{"type": "TypeAPL", "priorityList": [
		{"action": {"castSpell": {"spellId": {"spellId": 14287}}}},
		{"action": {"condition": {"and": {"vals": [
			{"cmp": {"op": "OpGt", "lhs": {"unitDistance": {"sourceUnit": {"type": "Self"}}}, "rhs": {"const": {"val": "5"}}}},
			{"cmp": {"op": "OpEq", "lhs": {"spellTimeToReady": {"spellId": {"spellId": 1317257}}}, "rhs": {"const": {"val": "0s"}}}}
		]}}, "move": {"rangeFromTarget": {"const": {"val": "3"}}}}},
		{"action": {"strictSequence": {"actions": [{"castSpell": {"spellId": {"spellId": 1317257}}}]}}},
		{"action": {"condition": {"and": {"vals": [
			{"cmp": {"op": "OpLe", "lhs": {"unitDistance": {"sourceUnit": {"type": "Self"}}}, "rhs": {"const": {"val": "5"}}}},
			{"cmp": {"op": "OpGt", "lhs": {"spellTimeToReady": {"spellId": {"spellId": 1317257}}}, "rhs": {"const": {"val": "0s"}}}}
		]}}, "move": {"rangeFromTarget": {"const": {"val": "11"}}}}}
	]}`
	firstAt := func(logs string, pattern string) float64 {
		t.Helper()
		m := regexp.MustCompile(`(?m)^\[([0-9.]+)\] .*` + pattern).FindStringSubmatch(logs)
		if m == nil {
			t.Fatalf("no log line matching %q", pattern)
		}
		at, _ := strconv.ParseFloat(m[1], 64)
		return at
	}
	kickAndMoveOut := func() (kick float64, moveOut float64, queued bool) {
		t.Helper()
		rotation := &proto.APLRotation{}
		if err := protojson.Unmarshal([]byte(rotationJSON), rotation); err != nil {
			t.Fatalf("rotation: %v", err)
		}
		player := &proto.Player{
			Name: "queued strike", Class: proto.Class_ClassHunter, Race: proto.Race_RaceNightElf, TalentsString: SurvivalTalents,
			Equipment: WeaponsOnly, Consumables: &proto.ConsumesSpec{}, DistanceFromTarget: 11, ReactionTimeMs: 100,
			Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{ClassOptions: &proto.HunterOptions{
				Ammo: proto.HunterOptions_Doomshot, PetType: proto.HunterOptions_PetNone}}}},
			Rotation: rotation,
		}
		result := core.RunRaidSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
			SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 3, Debug: true},
		})
		if result.Error != nil {
			t.Fatalf("sim failed: %s", result.Error.Message)
		}
		kick = firstAt(result.Logs, `Casting \{SpellID: 1317257\}`)
		moveOut = firstAt(result.Logs, `Moving to 11\.0 yards`)
		queued = regexp.MustCompile(`Queueing up \{SpellID: 1317257\}`).MatchString(result.Logs)
		return
	}

	kick, moveOut, queued := kickAndMoveOut()
	if !queued {
		t.Fatal("the Strider Kick was not queued; the test no longer covers a queued cast")
	}
	if wait := moveOut - kick; wait > 0.05 {
		t.Errorf("queued Strider Kick at %.2f s, walk out at %.2f s: %.2f s standing in melee, want the walk to start as the kick goes off",
			kick, moveOut, wait)
	}

	core.RotationRunsAfterQueuedCast = false
	defer func() { core.RotationRunsAfterQueuedCast = true }()
	kick, moveOut, _ = kickAndMoveOut()
	if wait := moveOut - kick; wait < 1 {
		t.Error(fmt.Sprintf("switch off: walk out %.2f s after the queued kick, want it to wait for the kick's GCD", wait))
	}
}
