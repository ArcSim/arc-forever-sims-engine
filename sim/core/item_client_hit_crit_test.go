//go:build with_db

package core

import (
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// Forever's revised items take the hit and crit their client equip spells state, of the type those
// spells name, in place of master's Classic rating (dbc.TakeClientHitAndCrit). The percents are not
// unified, so Blackhand's Breadth gives no spell crit and Eye of the Beast no melee crit. Checks the
// embedded binary used by real simulations.
func TestForeverRevisedItemsTakeClientHitAndCrit(t *testing.T) {
	type sheet struct {
		physicalHit, rangedHit, spellHit, physicalCrit, rangedCrit, spellCrit float64
	}
	for _, tc := range []struct {
		id   int32
		slot proto.ItemSlot
		want sheet
	}{
		{13965, proto.ItemSlot_ItemSlotTrinket1, sheet{physicalCrit: 1}}, // Blackhand's Breadth, 1318954: melee and ranged
		{13968, proto.ItemSlot_ItemSlotTrinket1, sheet{spellCrit: 1}},    // Eye of the Beast, 1318956: spells
		{13255, proto.ItemSlot_ItemSlotHands, sheet{rangedHit: 3}},       // Trueaim Gauntlets, 1301127: bows, guns and crossbows
		{22207, proto.ItemSlot_ItemSlotWaist, sheet{rangedHit: 1}},       // Sash of the Grand Hunt, 1294769
		{18388, proto.ItemSlot_ItemSlotRanged, sheet{rangedHit: 1}},      // Stoneshatter
		{18296, proto.ItemSlot_ItemSlotWrist, sheet{rangedHit: 1}},       // Marksman Bands
		{17754, proto.ItemSlot_ItemSlotLegs, sheet{rangedHit: 1}},        // Infernal Trickster Leggings
	} {
		var equipment Equipment
		equipment[tc.slot] = NewItem(ItemSpec{ID: tc.id})
		got := equipment.Stats(proto.Spec_SpecWarlock)
		ratings := [4]float64{got[stats.MeleeHitRating], got[stats.SpellHitRating], got[stats.MeleeCritRating], got[stats.SpellCritRating]}
		if ratings != [4]float64{} {
			t.Errorf("%d %s keeps hit and crit rating %v", tc.id, equipment[tc.slot].Name, ratings)
		}
		if sheet := (sheet{got[stats.PhysicalHitPercent], got[stats.RangedHitPercent], got[stats.SpellHitPercent],
			got[stats.PhysicalCritPercent], got[stats.RangedCritPercent], got[stats.SpellCritPercent]}); sheet != tc.want {
			t.Errorf("%d %s percents %+v, want %+v", tc.id, equipment[tc.slot].Name, sheet, tc.want)
		}
	}
}

// RandPropPoints prices no legendary item, so these took their client stats from the Forever
// planner (tools/database/gen_db/planner_legendary_stats.go). Generic hit and crit are rating,
// paid to every pool.
func TestForeverLegendaryItemStats(t *testing.T) {
	for id, want := range map[int32]map[stats.Stat]float64{
		22589: {stats.Stamina: 31, stats.Intellect: 32, stats.Spirit: 24, stats.MeleeHitRating: 20, stats.SpellCritRating: 0}, // Atiesh, mage
		22630: {stats.Stamina: 30, stats.Intellect: 29, stats.MeleeCritRating: 28},                                            // Atiesh, warlock
		22631: {stats.Stamina: 28, stats.Intellect: 28, stats.Spirit: 27},                                                     // Atiesh, priest
		22632: {stats.Stamina: 28, stats.Intellect: 28, stats.Spirit: 27},                                                     // Atiesh, druid
		19019: {stats.Agility: 5, stats.Stamina: 8, stats.FireResistance: 8, stats.NatureResistance: 9},                       // Thunderfury
		17182: {stats.Strength: 12, stats.Stamina: 12, stats.FireResistance: 30},                                              // Sulfuras
	} {
		item := NewItem(ItemSpec{ID: id})
		for stat, value := range want {
			if item.Stats[stat] != value {
				t.Errorf("%d %s %s = %v, want %v", id, item.Name, stat.StatName(), item.Stats[stat], value)
			}
		}
	}
}
