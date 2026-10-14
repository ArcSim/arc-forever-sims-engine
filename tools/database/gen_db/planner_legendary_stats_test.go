package main

import (
	"maps"
	"slices"
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/tools/database"
)

func TestFillPlannerLegendaryStats(t *testing.T) {
	const (
		strength  = int32(proto.Stat_StatStrength)
		agility   = int32(proto.Stat_StatAgility)
		stamina   = int32(proto.Stat_StatStamina)
		intellect = int32(proto.Stat_StatIntellect)
		spirit    = int32(proto.Stat_StatSpirit)
		meleeHit  = int32(proto.Stat_StatMeleeHitRating)
		meleeCrit = int32(proto.Stat_StatMeleeCritRating)
		spellCrit = int32(proto.Stat_StatSpellCritRating)
		fireRes   = int32(proto.Stat_StatFireResistance)
		natureRes = int32(proto.Stat_StatNatureResistance)
	)
	legendary := func(id int32, stats map[int32]float64) *proto.UIItem {
		item := statItem(id, stats)
		item.Quality = proto.ItemQuality_ItemQualityLegendary
		return item
	}
	db := database.NewWowDatabase()
	db.Items[22589] = legendary(22589, map[int32]float64{healingPower: 150, spellDamage: 150}) // Atiesh, mage
	db.Items[22630] = legendary(22630, map[int32]float64{healingPower: 183, spellDamage: 183}) // Atiesh, warlock
	db.Items[19019] = legendary(19019, nil)                                                    // Thunderfury
	db.Items[17182] = legendary(17182, map[int32]float64{strength: 15})                        // Sulfuras, a stat already set stays
	db.Items[18608] = statItem(18608, map[int32]float64{3: 31})                                // Benediction, epic: left alone

	changed := fillPlannerLegendaryStats(db, foreverPlanner)

	if want := []int32{17182, 19019, 22589, 22630}; !slices.Equal(changed, want) {
		t.Errorf("changed %v, want %v", changed, want)
	}
	for id, want := range map[int32]map[int32]float64{
		// The mage staff's splcritstrkpct is Power of the Guardian, not a stat of the wearer.
		22589: {healingPower: 150, spellDamage: 150, stamina: 31, intellect: 32, spirit: 24, meleeHit: 20},
		22630: {healingPower: 183, spellDamage: 183, stamina: 30, intellect: 29, meleeCrit: 28},
		19019: {agility: 5, stamina: 8, fireRes: 8, natureRes: 9},
		17182: {strength: 15, stamina: 12, fireRes: 30},
		18608: {3: 31},
	} {
		if got := db.Items[id].ScalingOptions[0].Stats; !maps.Equal(got, want) {
			t.Errorf("item %d stats %v, want %v", id, got, want)
		}
	}
	if db.Items[22589].ScalingOptions[0].Stats[spellCrit] != 0 {
		t.Error("the mage staff took its party aura as spell crit")
	}
}
