package main

import (
	"slices"
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/tools/database"
)

const foreverPlanner = "../../../assets/db_inputs/wowhead_forever_gearplanner.txt"

func statItem(id int32, stats map[int32]float64) *proto.UIItem {
	return &proto.UIItem{Id: id, ScalingOptions: map[int32]*proto.ScalingItemProperties{0: {Stats: stats}}}
}

const (
	spellDamage  = int32(proto.Stat_StatSpellDamage)
	healingPower = int32(proto.Stat_StatHealingPower)
	shadowDamage = int32(proto.Stat_StatShadowDamage)
)

func TestFillPlannerSpellPower(t *testing.T) {
	db := database.NewWowDatabase()
	db.Items[19102] = statItem(19102, map[int32]float64{3: 25, healingPower: 25, spellDamage: 25}) // Crackling Staff, client row
	db.Items[18608] = statItem(18608, map[int32]float64{3: 31})                                    // Benediction
	db.Items[20536] = statItem(20536, map[int32]float64{3: 16})                                    // Soul Harvester: shadow only
	db.Items[9415] = statItem(9415, map[int32]float64{spellDamage: 10})                            // Grimlok's Tribal Vestments
	db.Items[254013] = statItem(254013, map[int32]float64{healingPower: 18, spellDamage: 24})      // Earthen Silk Slippers: agrees
	area := statItem(270225, map[int32]float64{healingPower: 8, spellDamage: 8})                   // Jewel of the Guard Captain
	area.ScalingOptions[0].AreaStats = []*proto.AreaStats{{Stats: map[int32]float64{healingPower: 16, spellDamage: 16}}}
	db.Items[270225] = area

	changed := fillPlannerSpellPower(db, foreverPlanner)

	if want := []int32{9415, 18608, 19102, 20536}; !slices.Equal(changed, want) {
		t.Errorf("changed %v, want %v", changed, want)
	}
	for id, want := range map[int32]map[int32]float64{
		19102:  {3: 25, healingPower: 99, spellDamage: 99},
		18608:  {3: 31, healingPower: 195, spellDamage: 65},
		20536:  {3: 16, shadowDamage: 73},
		9415:   {healingPower: 29, spellDamage: 10},
		254013: {healingPower: 18, spellDamage: 24},
		270225: {healingPower: 8, spellDamage: 8},
	} {
		got := db.Items[id].ScalingOptions[0].Stats
		for stat, value := range want {
			if got[stat] != value {
				t.Errorf("item %d stat %d = %v, want %v", id, stat, got[stat], value)
			}
		}
	}
}

// Never lowers a stat the item already has more of.
func TestFillPlannerSpellPowerKeepsHigherValues(t *testing.T) {
	db := database.NewWowDatabase()
	db.Items[19102] = statItem(19102, map[int32]float64{healingPower: 120, spellDamage: 120})

	if changed := fillPlannerSpellPower(db, foreverPlanner); len(changed) != 0 {
		t.Errorf("changed %v", changed)
	}
	if got := db.Items[19102].ScalingOptions[0].Stats[spellDamage]; got != 120 {
		t.Errorf("spell damage lowered to %v", got)
	}
}

func TestFillClassicHealing(t *testing.T) {
	db := database.NewWowDatabase()
	db.Items[18534] = statItem(18534, map[int32]float64{spellDamage: 23}) // Rod of the Ogre Magi
	db.Items[9431] = statItem(9431, map[int32]float64{healingPower: 22})  // Papal Fez: healing only already
	db.Items[9415] = statItem(9415, map[int32]float64{spellDamage: 10})   // Grimlok's: the Forever planner lists it
	db.Items[11665] = statItem(11665, map[int32]float64{spellDamage: 13}) // Ogreseer Fists, not one of the filled-in items
	filled := []int32{18534, 9431, 9415}

	changed := fillClassicHealing(db, filled, classicPlanner, foreverPlanner)

	if want := []int32{18534}; !slices.Equal(changed, want) {
		t.Errorf("changed %v, want %v", changed, want)
	}
	for id, want := range map[int32]float64{18534: 23, 9431: 22, 9415: 0, 11665: 0} {
		if got := db.Items[id].ScalingOptions[0].Stats[healingPower]; got != want {
			t.Errorf("item %d healing = %v, want %v", id, got, want)
		}
	}
}
