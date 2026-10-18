package retribution

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
	"github.com/wowsims/forever/sim/paladin"
)

// MythicSim patch 159: Champion of the Light adds 20% of Intellect a rank to spell damage (the 17 October
// hotfix); ChampionOfTheLightReverted false gives the 16 October notes' 10% a rank.
func TestChampionOfTheLightShare(t *testing.T) {
	if !paladin.ChampionOfTheLightReverted {
		t.Fatal("ChampionOfTheLightReverted is off by default; patch 159 turns it on")
	}

	// Final intellect and spell damage, and the talent's rank, of a Retribution paladin with a flat
	// 1000 bonus Intellect and no gear.
	sheet := func(talents string) (intellect, spellDamage float64, rank int32) {
		items := make([]*proto.ItemSpec, 15)
		for i := range items {
			items[i] = &proto.ItemSpec{}
		}
		var bonus stats.Stats
		bonus[stats.Intellect] = 1000
		player := core.WithSpec(&proto.Player{
			Race:          proto.Race_RaceHuman,
			Class:         proto.Class_ClassPaladin,
			Equipment:     &proto.EquipmentSpec{Items: items},
			BonusStats:    &proto.UnitStats{Stats: bonus.ToProtoArray()},
			TalentsString: talents,
		}, DefaultOptions)
		raid := core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{})
		env, raidStats, _ := core.NewEnvironment(raid, core.MakeSingleTargetEncounter(0), false, true)
		final := raidStats.Parties[0].Players[0].FinalStats.Stats
		rank = env.Raid.Parties[0].Players[0].(paladin.PaladinAgent).GetPaladin().Talents.ChampionOfTheLight
		return final[stats.Intellect], final[stats.SpellDamage], rank
	}

	_, baseSpellDamage, baseRank := sheet("")
	if baseRank != 0 {
		t.Fatalf("an untalented paladin has Champion of the Light rank %d", baseRank)
	}

	for _, tc := range []struct {
		reverted bool
		perRank  float64
	}{
		{true, 0.20},
		{false, 0.10},
	} {
		withSwitch(&paladin.ChampionOfTheLightReverted, tc.reverted, func() {
			if got := paladin.ChampionOfTheLightIntellectShare(); got != tc.perRank {
				t.Errorf("reverted=%v: share a rank %.2f, want %.2f", tc.reverted, got, tc.perRank)
			}
			intellect, spellDamage, rank := sheet(RetTalents)
			if rank == 0 {
				t.Fatal("RetTalents has no Champion of the Light")
			}
			got := (spellDamage - baseSpellDamage) / intellect
			want := tc.perRank * float64(rank)
			if math.Abs(got-want) > 1e-3 {
				t.Errorf("reverted=%v: rank %d adds %.4f of Intellect as spell damage, want %.2f", tc.reverted, rank, got, want)
			}
		})
	}
}
