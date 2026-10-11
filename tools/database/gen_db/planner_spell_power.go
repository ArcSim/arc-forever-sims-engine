package main

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/tools/database"
)

// Wowhead's Forever planner states an item's spell power in three parts, and they add up. splpwr
// is the client's ItemSparse 45, damage and healing alike (spellPowerStats). spldmg is damage only
// and splheal healing only, on top of it. The client rows that carry all three agree: Earthen Silk
// Slippers, splpwr 18 and spldmg 6, is 24 SpellDamage and 18 HealingPower; Ogre Loincloth, splpwr
// 6 and splheal 12, is 6 and 18. The school keys (shasplpwr and the rest) are the school stats.
var plannerSpellPowerKeys = []struct {
	stat proto.Stat
	keys []string
}{
	{proto.Stat_StatSpellDamage, []string{"splpwr", "spldmg"}},
	{proto.Stat_StatHealingPower, []string{"splpwr", "splheal"}},
	{proto.Stat_StatArcaneDamage, []string{"arcsplpwr"}},
	{proto.Stat_StatFireDamage, []string{"firsplpwr"}},
	{proto.Stat_StatFrostDamage, []string{"frosplpwr"}},
	{proto.Stat_StatHolyDamage, []string{"holsplpwr"}},
	{proto.Stat_StatNatureDamage, []string{"natsplpwr"}},
	{proto.Stat_StatShadowDamage, []string{"shasplpwr"}},
}

// Forever gives caster weapons spell power the client's item rows do not carry: Crackling Staff
// shows 99 spell damage in game, its row 25 (ItemSparse 45) and Wowhead's Forever planner 25 + 74.
// The 74 scales with item level and quality, and a healer's weapon gets a third of it as damage and
// three times that as healing (Benediction 65 and 195). About 160 weapons have it, and none of them
// carried it here. Master's database, which fills in the items the client lacks, states "damage and
// healing" as SpellDamage alone, and the engine heals off HealingPower only, so those gave none.
//
// This raises each item's spell power stats to the planner's where the planner's are higher, and
// never lowers one. Items with area stats are left alone: the planner adds those in (Jewel of the
// Guard Captain shows 24, 8 everywhere and 16 more in the open world). Returns the ids it changed.
func fillPlannerSpellPower(db *database.WowDatabase, path string) []int32 {
	return applyPlannerSpellPower(db, readPlannerStats(path, "Forever gear planner"))
}

func applyPlannerSpellPower(db *database.WowDatabase, rows map[string]map[string]json.RawMessage) []int32 {
	var changed []int32
	for id, item := range db.Items {
		key := fmt.Sprint(id)
		stats, ok := rows[key]
		opt := item.ScalingOptions[0]
		if !ok || opt == nil || len(opt.AreaStats) > 0 {
			continue
		}
		raised := false
		for _, entry := range plannerSpellPowerKeys {
			var want float64
			for _, k := range entry.keys {
				value, _ := plannerStat(key, stats, k)
				want += value
			}
			if want <= opt.Stats[int32(entry.stat)] {
				continue
			}
			if opt.Stats == nil {
				opt.Stats = map[int32]float64{}
			}
			opt.Stats[int32(entry.stat)] = want
			raised = true
		}
		if raised {
			changed = append(changed, id)
		}
	}
	slices.Sort(changed)
	return changed
}

// The items mergeForeverSimDB filled in that the Forever planner does not list keep master's
// SpellDamage-only "damage and healing" (Rod of the Ogre Magi: 23 damage, no healing). Their
// healing is read from the same Wowhead Classic planner dropClassicBonusDamage reads, where
// splheal is the item's whole healing. Only an item with no HealingPower at all is filled, so
// master's healing-only rows (Papal Fez, 22) stay as they are. Returns the ids it changed.
func fillClassicHealing(db *database.WowDatabase, filled []int32, classicPlannerPath, foreverPlannerPath string) []int32 {
	return applyClassicHealing(db, filled,
		readPlannerStats(classicPlannerPath, "Classic gear planner"),
		readPlannerStats(foreverPlannerPath, "Forever gear planner"))
}

func applyClassicHealing(db *database.WowDatabase, filled []int32, classic, forever map[string]map[string]json.RawMessage) []int32 {
	healing := int32(proto.Stat_StatHealingPower)
	var changed []int32
	for _, id := range filled {
		key := fmt.Sprint(id)
		item := db.Items[id]
		if item == nil || item.ScalingOptions[0] == nil {
			continue
		}
		if _, listed := forever[key]; listed {
			continue
		}
		opt := item.ScalingOptions[0]
		heal, _ := plannerStat(key, classic[key], "splheal")
		if heal <= 0 || opt.Stats[healing] != 0 {
			continue
		}
		if opt.Stats == nil {
			opt.Stats = map[int32]float64{}
		}
		opt.Stats[healing] = heal
		changed = append(changed, id)
	}
	slices.Sort(changed)
	return changed
}
