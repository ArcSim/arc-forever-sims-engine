package dbc

import (
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
)

// Each index was read off a client item of that weapon type whose Wowhead Forever tooltip names the
// skill; the ones no item carries stay unmapped.
func TestMapBonusStatIndexToWeaponSkill(t *testing.T) {
	for index, want := range map[int]proto.PseudoStat{
		90:  proto.PseudoStat_PseudoStatTwoHandedAxesSkill,   // Huge Thorium Battleaxe
		91:  proto.PseudoStat_PseudoStatTwoHandedMacesSkill,  // Servomechanic Sledgehammer
		92:  proto.PseudoStat_PseudoStatTwoHandedSwordsSkill, // Bladewind
		96:  proto.PseudoStat_PseudoStatDaggersSkill,         // Death's Sting
		98:  proto.PseudoStat_PseudoStatUnarmedSkill,         // Punchy's Punchers
		101: proto.PseudoStat_PseudoStatPolearmsSkill,        // Frenzied Striker
		103: proto.PseudoStat_PseudoStatSwordsSkill,          // Test Sword Chest
	} {
		if got, ok := MapBonusStatIndexToWeaponSkill(index); !ok || got != want {
			t.Errorf("%d: %v %v, want %v", index, got, ok, want)
		}
	}
	for _, index := range []int{89, 93, 94, 95, 97, 99, 100, 102, 104, 112} {
		if got, ok := MapBonusStatIndexToWeaponSkill(index); ok {
			t.Errorf("%d mapped to %v", index, got)
		}
	}
	// None of them is a Stat as well, so GetStats skips them.
	for index := 90; index <= 103; index++ {
		if stat, ok := MapBonusStatIndexToStat(index); ok {
			t.Errorf("%d is also %v", index, stat)
		}
	}
}
