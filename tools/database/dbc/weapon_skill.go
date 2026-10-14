package dbc

import (
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// Forever moved the weapon skill of its items from Classic's equip spells into the item's stat
// array, with indices of its own. Each was read off the client's items, by the item's own weapon
// type and its Wowhead Forever tooltip: Huge Thorium Battleaxe (two-handed axe, 90, "Increased
// Two-Handed Axes +2"), Servomechanic Sledgehammer (two-handed mace, 91), Bladewind (two-handed
// sword, 92), Death's Sting (dagger, 96, "Increased Daggers +3"), Punchy's Punchers (98, fist
// weapons), Frenzied Striker (polearm, 101) and Test Sword Chest (103, "Increased Swords" in
// Classic). 93-95, 97, 99, 100 and 102 sit between them in what reads as an alphabetical order
// (axes, bows, crossbows, ..., staves) but no item carries one yet, so they stay unmapped. The
// amount is priced from the item's budget like any other stat (GetScaledStat).
var bonusStatWeaponSkills = map[int]proto.PseudoStat{
	90:  proto.PseudoStat_PseudoStatTwoHandedAxesSkill,
	91:  proto.PseudoStat_PseudoStatTwoHandedMacesSkill,
	92:  proto.PseudoStat_PseudoStatTwoHandedSwordsSkill,
	96:  proto.PseudoStat_PseudoStatDaggersSkill,
	98:  proto.PseudoStat_PseudoStatUnarmedSkill,
	101: proto.PseudoStat_PseudoStatPolearmsSkill,
	103: proto.PseudoStat_PseudoStatSwordsSkill,
}

func MapBonusStatIndexToWeaponSkill(index int) (proto.PseudoStat, bool) {
	skill, ok := bonusStatWeaponSkills[index]
	return skill, ok
}

// The weapon skill in an item's stat array, as pseudo stats.
func (item *Item) GetWeaponSkills(itemLevel int) []float64 {
	pseudoStats := make([]float64, stats.PseudoStatsLen)
	for i, index := range item.BonusStat {
		if skill, ok := MapBonusStatIndexToWeaponSkill(index); ok {
			pseudoStats[skill] += item.GetScaledStat(i, itemLevel)
		}
	}
	return pseudoStats
}

// The skill lines an A_MOD_SKILL equip spell can raise weapon skill on (EffectMiscValue[0]), as in
// Classic: Edgemaster's Handguards' 15874 is 44 (one-handed axes), Devilsaur Claws' 24362 162
// (unarmed, which fist weapons use).
var skillLineWeaponSkills = map[int]proto.PseudoStat{
	44:  proto.PseudoStat_PseudoStatAxesSkill,
	43:  proto.PseudoStat_PseudoStatSwordsSkill,
	54:  proto.PseudoStat_PseudoStatMacesSkill,
	173: proto.PseudoStat_PseudoStatDaggersSkill,
	162: proto.PseudoStat_PseudoStatUnarmedSkill,
	473: proto.PseudoStat_PseudoStatUnarmedSkill,
	172: proto.PseudoStat_PseudoStatTwoHandedAxesSkill,
	55:  proto.PseudoStat_PseudoStatTwoHandedSwordsSkill,
	160: proto.PseudoStat_PseudoStatTwoHandedMacesSkill,
	229: proto.PseudoStat_PseudoStatPolearmsSkill,
	136: proto.PseudoStat_PseudoStatStavesSkill,
	45:  proto.PseudoStat_PseudoStatBowsSkill,
	46:  proto.PseudoStat_PseudoStatGunsSkill,
	226: proto.PseudoStat_PseudoStatCrossbowsSkill,
	176: proto.PseudoStat_PseudoStatThrownSkill,
}

func MapSkillLineToWeaponSkill(skillLine int) (proto.PseudoStat, bool) {
	skill, ok := skillLineWeaponSkills[skillLine]
	return skill, ok
}

// Every weapon skill pseudo stat, for copying weapon skill alone between pseudo stat arrays.
func WeaponSkillPseudoStats() []proto.PseudoStat {
	out := make([]proto.PseudoStat, 0, len(skillLineWeaponSkills))
	for skill := proto.PseudoStat_PseudoStatAxesSkill; skill <= proto.PseudoStat_PseudoStatThrownSkill; skill++ {
		out = append(out, skill)
	}
	return out
}
