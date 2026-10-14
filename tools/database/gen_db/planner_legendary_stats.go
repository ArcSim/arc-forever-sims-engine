package main

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/tools/database"
)

// The planner keys a legendary item's client stats are read from, and the stat each fills.
var plannerLegendaryStats = []struct {
	key   string
	stat  proto.Stat
	scale float64
}{
	{"str", proto.Stat_StatStrength, 1},
	{"agi", proto.Stat_StatAgility, 1},
	{"sta", proto.Stat_StatStamina, 1},
	{"int", proto.Stat_StatIntellect, 1},
	{"spi", proto.Stat_StatSpirit, 1},
	// Generic in the client (ItemSparse 31 and 32), which the database reads as melee rating
	// like every other item's; UnifyGearHitAndCrit pays it into the spell pool too.
	{"hitpct", proto.Stat_StatMeleeHitRating, core.PhysicalHitRatingPerHitPercent},
	{"critstrkpct", proto.Stat_StatMeleeCritRating, core.PhysicalCritRatingPerCritPercent},
	{"arcres", proto.Stat_StatArcaneResistance, 1},
	{"firres", proto.Stat_StatFireResistance, 1},
	{"frores", proto.Stat_StatFrostResistance, 1},
	{"natres", proto.Stat_StatNatureResistance, 1},
	{"shares", proto.Stat_StatShadowResistance, 1},
}

// The Forever client prices an item's stats from its ItemSparse allocations and the RandPropPoints
// budget of its quality and item level (dbc.Item.GetScaledStat). RandPropPoints has columns for
// epic, rare and uncommon only, so a legendary item's budget is 0 and every allocated stat came out
// as nothing: the four Atiesh staves had no Intellect, Stamina or Spirit, the mage's no 2% hit
// (ItemSparse 31) and the warlock's no 2% crit (ItemSparse 32); Thunderfury and Sulfuras had none of
// their primary stats or resistances. Nor does the epic column fit: Atiesh's allocations come out
// at a budget of 59 at item level 90 (epic 64), Sulfuras' at 54 at item level 80 (epic 57).
//
// Wowhead's Forever planner states the priced stats, and each is the client's allocation at that
// budget: Atiesh (mage) 5254, 5424 and 4068 of 59 are its 31 Stamina, 32 Intellect and 24 Spirit, and
// 3390 is 20 hit rating, the planner's 2%. This fills a legendary item's primary stats, resistances,
// and generic hit and crit from the planner where the item has none. The planner's splcritstrkpct on
// the mage's staff is Power of the Guardian (28142), a party aura the item already carries as an
// effect, so it is not read. Returns the ids it changed.
func fillPlannerLegendaryStats(db *database.WowDatabase, path string) []int32 {
	return applyPlannerLegendaryStats(db, readPlannerStats(path, "Forever gear planner"))
}

func applyPlannerLegendaryStats(db *database.WowDatabase, rows map[string]map[string]json.RawMessage) []int32 {
	var changed []int32
	for id, item := range db.Items {
		if item.Quality != proto.ItemQuality_ItemQualityLegendary {
			continue
		}
		key := fmt.Sprint(id)
		stats, ok := rows[key]
		opt := item.ScalingOptions[0]
		if !ok || opt == nil {
			continue
		}
		filled := false
		for _, entry := range plannerLegendaryStats {
			value, ok := plannerStat(key, stats, entry.key)
			if !ok || value <= 0 || opt.Stats[int32(entry.stat)] != 0 {
				continue
			}
			if opt.Stats == nil {
				opt.Stats = map[int32]float64{}
			}
			opt.Stats[int32(entry.stat)] = value * entry.scale
			filled = true
		}
		if filled {
			changed = append(changed, id)
		}
	}
	slices.Sort(changed)
	return changed
}
