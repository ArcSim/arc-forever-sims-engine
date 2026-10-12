package rogue

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/common/forever"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// Frost Oil (mythicsim patch 126) casts its Frostbolt (350025) off 10% of the oiled weapon's landed
// hits and none of the other hand's.
func TestFrostOilProcsOffItsHand(t *testing.T) {
	run := func(enabled bool) (bolts, mhLanded, ohLanded float64) {
		saved := forever.EnableFrostOil
		forever.EnableFrostOil = enabled
		defer func() { forever.EnableFrostOil = saved }()

		player := core.WithSpec(&proto.Player{
			Race:          proto.Race_RaceHuman,
			Class:         proto.Class_ClassRogue,
			Equipment:     daggersOnly(),
			Consumables:   &proto.ConsumesSpec{MhImbueId: 3595}, // Frost Oil's use spell
			TalentsString: SubtletyTalents,
			Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, DefaultOptions)
		res := core.RunRaidSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{Iterations: 200, RandomSeed: 1},
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
		})
		if res.Error != nil {
			t.Fatal(res.Error.Message)
		}
		for _, action := range res.RaidMetrics.Parties[0].Players[0].Actions {
			for _, target := range action.Targets {
				switch {
				case action.Id.GetSpellId() == 350025:
					// A proc's spell records its outcomes, not casts.
					bolts += float64(target.Hits + target.Crits + target.Misses)
				case action.Id.GetOtherId() == proto.OtherAction_OtherActionAttack && action.Id.Tag == 1:
					mhLanded += float64(target.Hits + target.Crits + target.Glances + target.Blocks)
				case action.Id.GetOtherId() == proto.OtherAction_OtherActionAttack && action.Id.Tag == 2:
					ohLanded += float64(target.Hits + target.Crits + target.Glances + target.Blocks)
				}
			}
		}
		return bolts, mhLanded, ohLanded
	}

	bolts, mh, oh := run(true)
	if mh == 0 || oh == 0 {
		t.Fatalf("no landed swings: main hand %v, off hand %v", mh, oh)
	}
	rate := bolts / mh
	// 200 fights give a few thousand main-hand hits; four standard errors either side of 10%.
	if tolerance := 4 * math.Sqrt(0.1*0.9/mh); math.Abs(rate-0.1) > tolerance {
		t.Errorf("Frost Oil cast %v Frostbolts off %v main-hand hits (%.3f), want 10%% +/- %.3f", bolts, mh, rate, tolerance)
	}
	t.Logf("%v Frostbolts off %v main-hand and %v off-hand landed swings: %.4f of the main hand's", bolts, mh, oh, rate)
	if rate > 0.1*(mh+oh)/mh*0.75 {
		t.Errorf("Frost Oil's %v Frostbolts look like both hands' %v hits", bolts, mh+oh)
	}

	if off, _, _ := run(false); off != 0 {
		t.Errorf("with the switch off Frost Oil still cast %v Frostbolts", off)
	}
}
