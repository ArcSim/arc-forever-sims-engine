//go:build with_db

package core

import (
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// The weapon skill the embedded database gives each item: Forever's client items from their stat
// row (90-103) or equip spell, the unreleased items master filled in from Wowhead Classic's planner
// (tools/database/gen_db/classic_weapon_skill.go). The ranged skill items Forever rebuilt as ranged
// hit keep none.
func TestItemWeaponSkills(t *testing.T) {
	for id, want := range map[int32]map[proto.PseudoStat]float64{
		21126:  {proto.PseudoStat_PseudoStatDaggersSkill: 3},         // Death's Sting, client stat 96
		7951:   {proto.PseudoStat_PseudoStatDaggersSkill: 4},         // Hands of Thero-shan
		12775:  {proto.PseudoStat_PseudoStatTwoHandedAxesSkill: 2},   // Huge Thorium Battleaxe, Classic's +10 cut to +2
		4548:   {proto.PseudoStat_PseudoStatTwoHandedMacesSkill: 1},  // Servomechanic Sledgehammer
		276102: {proto.PseudoStat_PseudoStatTwoHandedSwordsSkill: 1}, // Bladewind
		274929: {proto.PseudoStat_PseudoStatTwoHandedAxesSkill: 1},   // Hurley's Felling Axe
		275669: {proto.PseudoStat_PseudoStatUnarmedSkill: 1},         // Punchy's Punchers, fist weapons
		13056:  {proto.PseudoStat_PseudoStatPolearmsSkill: 1},        // Frenzied Striker, client stat 101
		20005:  {proto.PseudoStat_PseudoStatUnarmedSkill: 4},         // Devilsaur Claws, client equip spell 24362
		18823:  {proto.PseudoStat_PseudoStatDaggersSkill: 5},         // Aged Core Leather Gloves, Classic
		18392:  {proto.PseudoStat_PseudoStatDaggersSkill: 6},         // Distracting Dagger, Classic
		18505:  {proto.PseudoStat_PseudoStatDaggersSkill: 5},         // Mugger's Belt, Classic
		18822:  {proto.PseudoStat_PseudoStatTwoHandedSwordsSkill: 8}, // Obsidian Edged Blade, Classic
		14551: { // Edgemaster's Handguards, Classic
			proto.PseudoStat_PseudoStatAxesSkill:    7,
			proto.PseudoStat_PseudoStatDaggersSkill: 7,
			proto.PseudoStat_PseudoStatSwordsSkill:  7,
		},
		13255: {}, // Trueaim Gauntlets: 3% ranged hit instead
		22207: {}, // Sash of the Grand Hunt: 1% ranged hit instead
		2907:  {}, // Dwarven Tree Chopper: expertise instead
	} {
		item := NewItem(ItemSpec{ID: id})
		for _, skill := range WeaponSkillPseudoStats {
			if got := stats.PseudoStatValue(item.PseudoStats, skill); got != want[skill] {
				t.Errorf("%d %s %s = %v, want %v", id, item.Name, skill, got, want[skill])
			}
		}
	}
}

// Each hand attacks with the skill of its own weapon type, summed over the equipment.
func TestCharacterWeaponSkillPerHand(t *testing.T) {
	for _, tc := range []struct {
		name              string
		items             map[proto.ItemSlot]int32
		mainHand, offHand float64
		ranged            float64
	}{
		// Redfall's test: Alcor's Sunrazor and Distracting Dagger with Aged Core Leather Gloves.
		{"daggers", map[proto.ItemSlot]int32{
			proto.ItemSlot_ItemSlotMainHand: 14555, proto.ItemSlot_ItemSlotOffHand: 18392, proto.ItemSlot_ItemSlotHands: 18823,
		}, 11, 11, 0},
		// Edgemaster's raises swords and daggers alike; Distracting Dagger only daggers.
		{"sword and dagger", map[proto.ItemSlot]int32{
			proto.ItemSlot_ItemSlotMainHand: 18832, proto.ItemSlot_ItemSlotOffHand: 18392, proto.ItemSlot_ItemSlotHands: 14551,
		}, 7, 13, 0},
		// Dagger skill does nothing for an axe.
		{"dagger gloves with an axe", map[proto.ItemSlot]int32{
			proto.ItemSlot_ItemSlotMainHand: 17068, proto.ItemSlot_ItemSlotHands: 18823,
		}, 0, 0, 0},
	} {
		character := &Character{Class: proto.Class_ClassWarrior}
		for slot, id := range tc.items {
			character.Equipment[slot] = NewItem(ItemSpec{ID: id})
		}
		character.updateWeaponSkills()
		if got := character.weaponSkill; got != [weaponSkillHands]float64{tc.mainHand, tc.offHand, tc.ranged} {
			t.Errorf("%s: %v, want %v %v %v", tc.name, got, tc.mainHand, tc.offHand, tc.ranged)
		}
	}
}
