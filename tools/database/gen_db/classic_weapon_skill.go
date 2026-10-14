package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/wowsims/forever/sim/core/stats"
	"github.com/wowsims/forever/tools/database"
	"github.com/wowsims/forever/tools/database/dbc"
)

// Weapon skill for the items mergeForeverSimDB fills in from master's database, whose stats leave
// it out (master's sim had no weapon skill to give it to).
//
// What Forever did to an item's weapon skill shows in the client's equip spells where it ships
// them: Trueaim Gauntlets, Sash of the Grand Hunt and the other ranged skill items carry new ranged
// hit spells and no skill, and MergeItemEffectsForAllStates takes whatever skill the client's spells
// state. An item with no client equip spell at all (Aged Core Leather Gloves, Distracting Dagger,
// Edgemaster's Handguards: Forever has not released them, and has dropped their Classic spells, as
// it does when it moves an item's equip spells into its stat row) keeps its Classic skill, read
// from Wowhead Classic's planner (skillBuff, by skill line), as it keeps its other Classic stats.
func fillClassicWeaponSkill(db *database.WowDatabase, filled []int32, itemEffects map[int][]dbc.ItemEffect, classicPlannerPath string) {
	rows := readPlannerStats(classicPlannerPath, "Classic gear planner")
	for _, id := range filled {
		item := db.Items[id]
		if item == nil || len(itemEffects[int(id)]) > 0 {
			continue
		}
		skills := plannerWeaponSkills(fmt.Sprint(id), rows[fmt.Sprint(id)])
		if len(skills) == 0 {
			continue
		}
		pseudoStats := make([]float64, stats.PseudoStatsLen)
		copy(pseudoStats, item.PseudoStats)
		for skill, value := range skills {
			pseudoStats[skill] += value
		}
		item.PseudoStats = pseudoStats
	}
}

// An item's weapon skill in Wowhead's planner, skillBuff {skill line: amount}, as pseudo stat
// indices. Skill lines that are not weapons (Lockpicking, Defense) are left out.
func plannerWeaponSkills(id string, stats map[string]json.RawMessage) map[int]float64 {
	raw, ok := stats["skillBuff"]
	if !ok {
		return nil
	}
	var buffs map[string]float64
	if err := json.Unmarshal(raw, &buffs); err != nil {
		panic(fmt.Sprintf("gear planner item %s: skillBuff is not a map of numbers: %s", id, raw))
	}
	out := map[int]float64{}
	for line, value := range buffs {
		skillLine, err := strconv.Atoi(line)
		if err != nil {
			panic(fmt.Sprintf("gear planner item %s: skill line %q", id, line))
		}
		if skill, ok := dbc.MapSkillLineToWeaponSkill(skillLine); ok {
			out[int(skill)] += value
		}
	}
	return out
}
