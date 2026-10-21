package hunter

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

func rangedMinimumSim(t *testing.T, distance float64) (*core.Simulation, *Hunter) {
	t.Helper()
	player := &proto.Player{
		Name: "ranged minimum", Class: proto.Class_ClassHunter, Race: proto.Race_RaceNightElf, TalentsString: MarksmanshipTalents,
		Equipment: WeaponsOnly, Consumables: &proto.ConsumesSpec{}, DistanceFromTarget: distance,
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

// Distances are from the target's center. The client's 8 yd ("Hunter Range", 8-35) counts from the hitbox edges, so
// the shot stops well inside a mob's center: two logs put Auto Shot releases at 10.5 and 10.8 yd from a wolf's center
// (warcraftlogs 7CbTgWdMmYQ1qwJN) and 37 ranged casts in a dungeon run no closer than 11.4 yd (T7GB4KhFXyaNVR91).
// Standing at 10 yd the Auto Shot never fires; at 11 yd it does (MinRangedRange 10.5: shots go off from beyond it).
func TestRangedDeadzoneEndsBetweenTenAndElevenYards(t *testing.T) {
	for _, tc := range []struct {
		distance float64
		fires    bool
	}{{10, false}, {11, true}} {
		sim, hunter := rangedMinimumSim(t, tc.distance)
		for i := 0; i < 100000 && sim.CurrentTime < 5*time.Second; i++ {
			sim.Step()
		}
		if fired := hunter.AutoAttacks.PreviousRangedAttack() > 0; fired != tc.fires {
			t.Errorf("standing at %.0f yd for 5s: Auto Shot fired %v, want %v", tc.distance, fired, tc.fires)
		}
	}
}
