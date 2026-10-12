package rogue

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// Battle Squawk (mythicsim patch 127): each stack the party assumes is 5% melee attack speed, the
// stacks multiply, and they last the squawk's 4 minutes.
func TestBattleSquawkStacksMeleeSpeed(t *testing.T) {
	swings := func(stacks int32, seconds float64) float64 {
		player := core.WithSpec(&proto.Player{
			Race:          proto.Race_RaceHuman,
			Class:         proto.Class_ClassRogue,
			Equipment:     daggersOnly(),
			Consumables:   &proto.ConsumesSpec{},
			TalentsString: SubtletyTalents,
			Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, DefaultOptions)
		raid := core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{})
		raid.Parties[0].BattleSquawks = stacks
		encounter := core.MakeSingleTargetEncounter(0)
		encounter.Duration = seconds
		res := core.RunRaidSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{Iterations: 20, RandomSeed: 1},
			Raid:       raid,
			Encounter:  encounter,
		})
		if res.Error != nil {
			t.Fatal(res.Error.Message)
		}
		total := 0.0
		for _, action := range res.RaidMetrics.Parties[0].Players[0].Actions {
			if action.Id.GetOtherId() == proto.OtherAction_OtherActionAttack && action.Id.Tag == 1 {
				for _, target := range action.Targets {
					total += float64(target.Casts)
				}
			}
		}
		return total
	}

	base := swings(0, 120)
	for _, tc := range []struct {
		stacks  int32
		seconds float64
		want    float64
	}{
		{1, 120, 1.05},
		{5, 120, math.Pow(1.05, 5)},
		{9, 120, math.Pow(1.05, 5)}, // a party holds five
		// Half of an eight minute fight runs without the squawks.
		{5, 480, (math.Pow(1.05, 5) + 1) / 2},
	} {
		got := swings(tc.stacks, tc.seconds) / (base * tc.seconds / 120)
		if math.Abs(got-tc.want) > 0.01 {
			t.Errorf("%d Battle Squawks over %.0f s: %.4f times the main-hand swings, want %.4f", tc.stacks, tc.seconds, got, tc.want)
		}
	}

	saved := core.EnableBattleSquawk
	core.EnableBattleSquawk = false
	defer func() { core.EnableBattleSquawk = saved }()
	if off := swings(5, 120); off != base {
		t.Errorf("with the switch off 5 Battle Squawks still changed the swings: %v, want %v", off, base)
	}
}
