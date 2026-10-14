package dbc

import (
	"maps"
	"slices"
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
)

func TestAddAreaStatsKeepsOneEntryPerArea(t *testing.T) {
	ap := int32(proto.Stat_StatAttackPower)
	sp := int32(proto.Stat_StatSpellDamage)
	source := map[int32]float64{ap: 29}
	opts := &proto.ScalingItemProperties{Stats: map[int32]float64{ap: 15}}

	addAreaStats(opts, proto.AreaType_AreaTypeForestGrassland, source)
	addAreaStats(opts, proto.AreaType_AreaTypeForestGrassland, map[int32]float64{ap: 1, sp: 8})
	addAreaStats(opts, proto.AreaType_AreaTypeMountainous, map[int32]float64{sp: 6})

	if got := opts.Stats[ap]; got != 15 {
		t.Errorf("the unconditional stats changed: attack power = %v, want 15", got)
	}
	if len(opts.AreaStats) != 2 {
		t.Fatalf("got %d area entries, want one per area", len(opts.AreaStats))
	}
	forest, mountain := opts.AreaStats[0], opts.AreaStats[1]
	if forest.AreaType != proto.AreaType_AreaTypeForestGrassland || forest.Stats[ap] != 30 || forest.Stats[sp] != 8 {
		t.Errorf("forest entry = %v, want attack power 30 and spell damage 8", forest)
	}
	if mountain.AreaType != proto.AreaType_AreaTypeMountainous || mountain.Stats[sp] != 6 {
		t.Errorf("mountain entry = %v, want spell damage 6", mountain)
	}
	if source[ap] != 29 {
		t.Errorf("the caller's map was written through: %v", source)
	}
}

func TestTakeClientHitAndCrit(t *testing.T) {
	pseudo := func(values map[proto.PseudoStat]float64) []float64 {
		out := make([]float64, proto.PseudoStat_PseudoStatRangedCritPercent+8)
		for stat, value := range values {
			out[stat] = value
		}
		return out
	}
	meleeHit, spellHit := int32(proto.Stat_StatMeleeHitRating), int32(proto.Stat_StatSpellHitRating)
	meleeCrit, spellCrit := int32(proto.Stat_StatMeleeCritRating), int32(proto.Stat_StatSpellCritRating)
	agility := int32(proto.Stat_StatAgility)

	for _, tc := range []struct {
		name      string
		stats     map[int32]float64
		spell     map[proto.PseudoStat]float64
		wantStats map[int32]float64
		wantTook  bool
	}{
		{"Blackhand's Breadth: 1318954's melee and ranged crit replaces Classic's 2% crit rating",
			map[int32]float64{meleeCrit: 28},
			map[proto.PseudoStat]float64{proto.PseudoStat_PseudoStatMeleeCritPercent: 1, proto.PseudoStat_PseudoStatRangedCritPercent: 1},
			map[int32]float64{}, true},
		{"Eye of the Beast: 1318956's spell crit replaces Classic's spell crit rating",
			map[int32]float64{spellCrit: 28},
			map[proto.PseudoStat]float64{proto.PseudoStat_PseudoStatSpellCritPercent: 1},
			map[int32]float64{}, true},
		{"Trueaim Gauntlets: ranged hit replaces the hit and leaves the crit and the rest",
			map[int32]float64{meleeHit: 10, spellHit: 10, meleeCrit: 14, agility: 9},
			map[proto.PseudoStat]float64{proto.PseudoStat_PseudoStatRangedHitPercent: 3},
			map[int32]float64{meleeCrit: 14, agility: 9}, true},
		{"a spell with neither leaves the item alone",
			map[int32]float64{meleeHit: 10},
			map[proto.PseudoStat]float64{proto.PseudoStat_PseudoStatDodgePercent: 1},
			map[int32]float64{meleeHit: 10}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			item := &proto.UIItem{ScalingOptions: map[int32]*proto.ScalingItemProperties{0: {Stats: tc.stats}}}
			into := pseudo(nil)
			into[proto.PseudoStat_PseudoStatBlockPercent] = 2
			equip := pseudo(tc.spell)

			took := TakeClientHitAndCrit(item, equip, into)

			if took != tc.wantTook {
				t.Errorf("took = %v, want %v", took, tc.wantTook)
			}
			if got := item.ScalingOptions[0].Stats; !maps.Equal(got, tc.wantStats) {
				t.Errorf("stats = %v, want %v", got, tc.wantStats)
			}
			want := pseudo(nil)
			want[proto.PseudoStat_PseudoStatBlockPercent] = 2
			if tc.wantTook {
				for stat, value := range tc.spell {
					want[stat] += value
				}
			}
			if !slices.Equal(into, want) {
				t.Errorf("pseudo stats = %v, want %v", into, want)
			}
		})
	}
}
