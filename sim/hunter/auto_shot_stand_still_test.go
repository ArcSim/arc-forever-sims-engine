package hunter

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

func autoShotStandStillSim(t *testing.T) (*core.Simulation, *Hunter) {
	t.Helper()
	player := &proto.Player{
		Name: "stand still", Class: proto.Class_ClassHunter, Race: proto.Race_RaceNightElf, TalentsString: MarksmanshipTalents,
		Equipment: WeaponsOnly, Consumables: &proto.ConsumesSpec{}, DistanceFromTarget: 25,
		Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{ClassOptions: &proto.HunterOptions{
			Ammo: proto.HunterOptions_Doomshot, PetType: proto.HunterOptions_PetNone}}}},
		Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	sim.PrePull() // starts the autos at 0s
	return sim, sim.Raid.Parties[0].Players[0].(HunterAgent).GetHunter()
}

// 16 October 2026 beta notes: "Auto Shot now requires 0.5 seconds of standing still before shooting".
// The second Auto Shot is due at T; a move from moveStart to moveEnd (relative to T) holds it until the
// Hunter has stood still for half a second, where before it fired at T or, for a shot that came due
// mid-move, at the first 500ms retry after the move.
func TestAutoShotNeedsHalfASecondStandingStill(t *testing.T) {
	ms := time.Millisecond
	for _, tc := range []struct {
		name               string
		move               bool
		moveStart, moveEnd time.Duration
		want               time.Duration
	}{
		{"no move", false, 0, 0, 0},
		{"move ends 0.6s before", true, -1000 * ms, -600 * ms, 0},
		{"move ends exactly 0.5s before", true, -1000 * ms, -500 * ms, 0},
		{"move ends 0.3s before", true, -1000 * ms, -300 * ms, 200 * ms},
		{"move ends 0.2s after", true, -500 * ms, 200 * ms, 700 * ms},
		{"move ends 0.4s after", true, -500 * ms, 400 * ms, 900 * ms},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sim, hunter := autoShotStandStillSim(t)
			aa := &hunter.AutoAttacks

			// Let the first shot go, so the next is a full swing out.
			for i := 0; i < 100000 && aa.PreviousRangedAttack() <= 0; i++ {
				sim.Step()
			}
			due := aa.NextRangedAttackAt()
			if due-sim.CurrentTime < 1500*ms {
				t.Fatalf("test setup: the next Auto Shot is due %v after now, want at least 1.5s", due-sim.CurrentTime)
			}
			if tc.move {
				sim.AddPendingAction(&core.PendingAction{NextActionAt: due + tc.moveStart, OnAction: func(sim *core.Simulation) {
					hunter.MoveDuration(tc.moveEnd-tc.moveStart, sim)
				}})
			}
			for i := 0; i < 100000 && aa.PreviousRangedAttack() < due; i++ {
				sim.Step()
			}
			if got := aa.PreviousRangedAttack() - due; got != tc.want {
				t.Errorf("Auto Shot fired %v after it was due, want %v", got, tc.want)
			}
		})
	}
}
