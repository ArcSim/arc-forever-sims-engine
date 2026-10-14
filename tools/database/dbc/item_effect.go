package dbc

import (
	"fmt"
	"maps"
	"math"
	"slices"
	"strings"

	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

// Kept apart from the item's stats so the sim can apply them per encounter.
func addAreaStats(opts *proto.ScalingItemProperties, areaType proto.AreaType, statMap map[int32]float64) {
	for _, existing := range opts.AreaStats {
		if existing.AreaType == areaType {
			for stat, value := range statMap {
				existing.Stats[stat] += value
			}
			return
		}
	}
	opts.AreaStats = append(opts.AreaStats, &proto.AreaStats{AreaType: areaType, Stats: maps.Clone(statMap)})
}

// ItemEffect represents an item effect in the game.
type ItemEffect struct {
	ID                   int // Effect ID
	LegacySlotIndex      int // Legacy slot index
	TriggerType          int // Trigger type
	Charges              int // Number of charges
	CoolDownMSec         int // Cooldown in milliseconds
	CategoryCoolDownMSec int // Category cooldown in milliseconds
	SpellCategoryID      int // Spell category ID
	MaxCumulativeStacks  int // Max cumulative stacks
	SpellID              int // Spell ID
	ChrSpecializationID  int // Character specialization ID
	ParentItemID         int // Parent item ID
}

// ToMap returns a generic representation of the effect.
func (e *ItemEffect) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"ID":                   e.ID,
		"LegacySlotIndex":      e.LegacySlotIndex,
		"TriggerType":          e.TriggerType,
		"Charges":              e.Charges,
		"CoolDownMSec":         e.CoolDownMSec,
		"CategoryCoolDownMSec": e.CategoryCoolDownMSec,
		"MaxCumulativeStacks":  e.MaxCumulativeStacks,
		"SpellCategoryID":      e.SpellCategoryID,
		"SpellID":              e.SpellID,
		"ChrSpecializationID":  e.ChrSpecializationID,
		"ParentItemID":         e.ParentItemID,
	}
}

func GetItemEffect(effectId int) ItemEffect {
	return dbcInstance.ItemEffects[effectId]
}

func makeBaseProto(e *ItemEffect, statsSpellID int) *proto.ItemEffect {
	sp := dbcInstance.Spells[e.SpellID]
	base := &proto.ItemEffect{
		BuffId:           int32(e.SpellID),
		BuffName:         fmt.Sprintf("%s (%d)", sp.NameLang, e.SpellID),
		EffectDurationMs: int32(sp.Duration),
		ScalingOptions:   make(map[int32]*proto.ScalingItemEffectProperties),
	}
	// override duration if stats spell defines its own
	if dur := dbcInstance.Spells[statsSpellID].Duration; dur > 0 {
		base.EffectDurationMs = int32(dur)
	}
	return base
}

// The internal cooldown of a proc. Most spells carry it on the trigger's SpellAuraOptions, but a
// few leave that at zero and put it on the buff's spell category instead: Bulwark of Azzinoth's
// armor buff sits in category 30 with a 60s CategoryRecoveryTime while its trigger has nothing.
//
// Only a buff that is a *separate* spell from the trigger counts. Where the two are the same spell
// the recovery being read is that spell's own cast throttle rather than a gate on re-applying the
// buff, and in TBC that is almost always the generic 1s bucket - category 99 holds 170 spells, 101
// of them with exactly 1000, most of them quest-reward trinkets. Reading those as internal
// cooldowns perturbs every fixture that equips one and states nothing the game does not already do.
func procIcdMs(trigger Spell, buff Spell) int32 {
	if trigger.ProcCategoryRecovery > 0 {
		return trigger.ProcCategoryRecovery
	}

	if buff.ID != trigger.ID {
		return buff.CategoryRecoveryTime
	}

	return 0
}

func assignTrigger(e *ItemEffect, statsSpellID int, pe *proto.ItemEffect) {
	spTop := dbcInstance.Spells[e.SpellID]
	statsSP := dbcInstance.Spells[statsSpellID]

	switch resolveTriggerType(e.TriggerType, e.SpellID) {
	case ITEM_SPELLTRIGGER_ON_USE:
		pe.Effect = &proto.ItemEffect_OnUse{OnUse: &proto.OnUseEffect{
			CooldownMs:         int32(e.CoolDownMSec),
			CategoryId:         int32(e.SpellCategoryID),
			CategoryCooldownMs: int32(e.CategoryCoolDownMSec),
		}}
	case ITEM_SPELLTRIGGER_CHANCE_ON_HIT:
		proc := &proto.ProcEffect{
			IcdMs: procIcdMs(spTop, statsSP),
		}
		// The column is the roll only where it is one: 0 means the client states no rate at all and
		// anything above 100 is its "the rate lives elsewhere" sentinel. A proc whose rate is not in
		// the spell data reaches the sim through the spell store, where an override answers it.
		if spTop.ProcChance > 0 && spTop.ProcChance <= 100 {
			proc.ProcRate = &proto.ProcEffect_ProcChance{
				ProcChance: float64(spTop.ProcChance) / 100,
			}
		}
		pe.BuffId = statsSP.ID
		pe.BuffName = fmt.Sprintf("%s (%d)", statsSP.NameLang, e.SpellID)
		pe.Effect = &proto.ItemEffect_Proc{Proc: proc}
	}

	// Stacks are a property of the buff aura, not of how it is triggered, so this sits outside
	// the switch: Pendant of the Violet Eye is an on-use whose buff, 35095 Enlightenment, has
	// CumulativeAura 20, and assigning only on the chance-on-hit path skipped it entirely.
	//
	// The count lives on whichever spell carries the aura, usually the stats spell rather than
	// the trigger - Blackened Naaru Sliver's trigger has none while its buff, 45041 Combat
	// Insight, has 10.
	maxCumulativeStacks := max(statsSP.MaxCumulativeStacks, spTop.MaxCumulativeStacks)
	if maxCumulativeStacks > 0 {
		pe.MaxCumulativeStacks = maxCumulativeStacks
	}
}

// The stats an effect grants. A stacking effect grants nothing through the aura its trigger
// applies - the stats it is named after live on the aura that accumulates - so an empty scaling
// option is not the same as no stats, and no caller may read ScalingOptions directly to answer
// this question.
func EffectStats(effect *proto.ItemEffect) map[int32]float64 {
	if stats := effect.GetScalingOptions()[int32(0)].GetStats(); len(stats) > 0 {
		return stats
	}

	return effect.GetStackingAura().GetScalingOptions()[int32(0)].GetStats()
}

// On a stacking effect the aura the trigger applies grants nothing itself, so its scaling options
// resolve to an empty message. ScalingItemEffectProperties holds only a stats map, so such an entry
// carries no information at all - drop the whole map rather than ship an empty entry per state.
// Readers reach the map through GetStats(), which is nil-safe.
//
// A no-op when anything did resolve, which keeps this honest if an effect ever carries stats on
// both the trigger aura and the aura it accumulates.
func dropEmptyScalingOptions(pe *proto.ItemEffect) {
	for _, opt := range pe.ScalingOptions {
		if len(opt.GetStats()) > 0 {
			return
		}
	}

	pe.ScalingOptions = nil
}

func (e *ItemEffect) ToProto(itemLevel int) (*proto.ItemEffect, bool) {
	statsSpellID := resolveStatsSpell(e.SpellID)

	pe := makeBaseProto(e, statsSpellID)
	assignTrigger(e, statsSpellID, pe)

	pe.ScalingOptions[int32(0)] = buildItemEffectScalingProps(e.SpellID, itemLevel)

	// The stats may live on the accumulating aura rather than on the one the trigger applies, in
	// which case the effect is real even though the scaling options above resolved to nothing.
	if stacking := buildStackingAura(e.SpellID, statsSpellID, itemLevel); stacking != nil {
		stacking.Aura.ScalingOptions[int32(0)] = buildItemEffectScalingProps(int(stacking.Aura.BuffId), itemLevel)
		applyStackingAura(pe, stacking)
	}

	dropEmptyScalingOptions(pe)

	if len(EffectStats(pe)) == 0 {
		return nil, false
	}

	return pe, true
}

func resolveStatsSpell(spellID int) int {
	return newChainWalker().resolveStatsSpell(spellID)
}

func (w *chainWalker) resolveStatsSpell(spellID int) int {
	effects := w.effects(spellID)
	for _, se := range effects {
		if se.GrantsStats() {
			return spellID
		}
	}

	// If we cant resolve the spell in the first loop, we follow proc triggers downwards
	for _, se := range effects {
		if se.IsProcTrigger() {
			return w.resolveStatsSpell(se.EffectTriggerSpell)
		}
	}
	return spellID
}

func resolveTriggerType(topType, spellID int) int {
	if topType == ITEM_SPELLTRIGGER_ON_USE || topType == ITEM_SPELLTRIGGER_CHANCE_ON_HIT {
		return topType
	}

	// Carrying A_PROC_TRIGGER_SPELL is not on its own a proc. The client also uses that aura to hang
	// a permanent sub-aura off an equip effect: Leggings of Beast Mastery and Void Star Talisman
	// grant pet stats that way, with no proc mask, no chance and no duration, and promoting them
	// emitted a proc with no rate at all. Something has to say when it would fire.
	//
	// The chance-on-hit weapons are unaffected: the database types those CHANCE_ON_HIT itself, so
	// they return above without being promoted here.
	sp := dbcInstance.Spells[spellID]
	statesAProc := slices.ContainsFunc(sp.ProcTypeMask, func(bits int) bool { return bits != 0 }) || sp.ProcChance != 0
	if !statesAProc {
		return topType
	}

	for _, se := range dbcInstance.SpellEffectsInOrder(spellID) {
		if se.IsProcTrigger() {
			return ITEM_SPELLTRIGGER_CHANCE_ON_HIT
		}
	}

	return topType
}

func buildItemEffectScalingProps(spellID int, itemLevel int) *proto.ScalingItemEffectProperties {
	return &proto.ScalingItemEffectProperties{Stats: collectStats(spellID, itemLevel).ToProtoMap()}
}

func buildBaseStatScalingProps(spellID int, itemSpellID int) *proto.ScalingItemEffectProperties {
	var total stats.Stats

	// check if spell is procced by a SPELL_WITH_VALUE
	if effects := dbcInstance.SpellEffectsInOrder(itemSpellID); len(effects) > 0 {
		for _, se := range effects {
			if se.HitsAnEnemy() {
				continue
			}
			// TBC ANNI: Items can have "static" ItemEffects that don't have a duration.
			// We need to parse these into stats just as is done for ItemSparse data.
			stat := ConvertEffectAuraToStatIndex(se.EffectAura, se.EffectMiscValues[0])
			if stat >= 0 || stat == -2 {
				value := float64(se.EffectBasePoints)
				// Make sure it's not Feral AP
				if strings.Contains(dbcInstance.Spells[se.SpellID].Description, "forms only") {
					stat = proto.Stat_StatFeralAttackPower
				}
				if stat == proto.Stat_StatArmorPenetration || stat == proto.Stat_StatSpellPiercing {
					// Make these not negative
					value = math.Abs(value)
				}

				if se.EffectAura == dbcenums.A_MOD_RESISTANCE && stat == -2 {
					// All Resists
					total[stats.ArcaneResistance] += value
					total[stats.FireResistance] += value
					total[stats.FrostResistance] += value
					total[stats.NatureResistance] += value
					total[stats.ShadowResistance] += value
				} else {
					total[int32(stat)] += value
				}

				continue
			}

			if se.EffectAura == dbcenums.A_PROC_TRIGGER_SPELL_WITH_VALUE && spellID == se.EffectTriggerSpell {
				for idx := range total {
					if total[idx] == 0 {
						continue
					}

					total[idx] = float64(se.EffectBasePoints)
				}
			}
		}
	}

	return &proto.ScalingItemEffectProperties{Stats: total.ToProtoMap()}
}

func collectStats(spellID, itemLevel int) stats.Stats {
	var total stats.Stats
	newChainWalker().collectStats(spellID, itemLevel, &total)
	return total
}

func (w *chainWalker) collectStats(spellID, itemLevel int, total *stats.Stats) {
	sp := dbcInstance.Spells[spellID]
	for _, se := range w.effects(spellID) {
		if se.HitsAnEnemy() {
			continue
		}
		if s, resolved := se.ParseStatEffect(sp.ScalesWithItemLevel(), itemLevel); resolved {
			total.AddInplace(&s)
		} else if se.EffectAura == dbcenums.A_PROC_TRIGGER_SPELL {
			// Deliberately narrower than IsProcTrigger: descending through an
			// A_PROC_TRIGGER_SPELL_WITH_VALUE would collect the triggered spell's own amounts,
			// past the point where the caller can still override them with the value the
			// trigger carries.
			w.collectStats(se.EffectTriggerSpell, itemLevel, total)
		}
	}
}

func ParseItemEffects(itemID, itemLevel int) []*proto.ItemEffect {
	raw := dbcInstance.ItemEffectsByParentID[itemID]
	out := make([]*proto.ItemEffect, 0, len(raw))
	for _, ie := range raw {
		if pe, ok := ie.ToProto(itemLevel); ok {
			out = append(out, pe)
		}
	}
	return out
}

func GetItemEffectSpellTooltip(itemID int, buffId int) (string, int) {
	raw := dbcInstance.ItemEffectsByParentID[itemID]
	var spellID int

	for _, effect := range raw {
		spellID = effect.SpellID
		if effect.SpellID == buffId {
			spellID = effect.SpellID
			break
		} else {
			triggerEffects := dbcInstance.SpellEffects[effect.SpellID]
			if len(triggerEffects) == 0 {
				continue
			}
			if spellEffect := GetSpellEffectRecursive(buffId, triggerEffects); spellEffect != nil {
				if spellEffect.EffectTriggerSpell == buffId {
					spellID = effect.SpellID
				}
				break
			}
		}
	}
	spell := dbcInstance.Spells[spellID]
	return spell.Description, spellID
}

func GetItemEffectForBuffID(itemID int, buffId int) *ItemEffect {
	raw := dbcInstance.ItemEffectsByParentID[itemID]
	var itemEffect *ItemEffect
	for _, effect := range raw {
		if effect.SpellID == buffId {
			itemEffect = &effect
			break
		} else {
			triggerEffects := dbcInstance.SpellEffects[effect.SpellID]
			if len(triggerEffects) == 0 {
				continue
			}
			if spellEffect := GetSpellEffectRecursive(buffId, triggerEffects); spellEffect != nil {
				if spellEffect.EffectTriggerSpell == buffId {
					return &effect
				}
				break
			}
		}
	}
	return itemEffect
}

// Walks the trigger chain below the given effects and returns the effect that triggers
// spellIDToMatch, or nil. Effects are visited in index order and a branch that does not contain
// the match no longer aborts the search of its siblings, which is what hid any buff reachable
// only through a later effect.
//
// Takes the effects rather than a spell ID because that is what the call sites hold; everything
// below the first level goes through the chain walker.
func GetSpellEffectRecursive(spellIDToMatch int, spellEffects map[int]SpellEffect) *SpellEffect {
	w := newChainWalker()

	indices := make([]int, 0, len(spellEffects))
	for idx := range spellEffects {
		indices = append(indices, idx)
	}
	slices.Sort(indices)

	for _, idx := range indices {
		spellEffect := spellEffects[idx]
		if spellEffect.EffectTriggerSpell == 0 {
			continue
		}

		if spellEffect.EffectTriggerSpell == spellIDToMatch {
			return &spellEffect
		}

		if match := w.findTriggerOf(spellIDToMatch, spellEffect.EffectTriggerSpell); match != nil {
			return match
		}
	}
	return nil
}

func (w *chainWalker) findTriggerOf(spellIDToMatch int, spellID int) *SpellEffect {
	for _, spellEffect := range w.effects(spellID) {
		if spellEffect.EffectTriggerSpell == 0 {
			continue
		}

		if spellEffect.EffectTriggerSpell == spellIDToMatch {
			return &spellEffect
		}

		if match := w.findTriggerOf(spellIDToMatch, spellEffect.EffectTriggerSpell); match != nil {
			return match
		}
	}
	return nil
}

// Parses a UIItem and loops through Scaling Options for that item.
func MergeItemEffectsForAllStates(parsed *proto.UIItem) []*proto.ItemEffect {
	var effects []*proto.ItemEffect
	pseudoStats := make([]float64, stats.PseudoStatsLen)

	itemEffects := dbcInstance.ItemEffectsByParentID[int(parsed.Id)]
	for idx := range itemEffects {
		baseEff := &itemEffects[idx]

		if baseEff.TriggerType == ITEM_SPELLTRIGGER_ON_EQUIP {
			props := buildBaseStatScalingProps(resolveStatsSpell(baseEff.SpellID), baseEff.SpellID)
			equipStats := stats.FromProtoMap(props.Stats)
			equipPseudoStats := make([]float64, stats.PseudoStatsLen)
			addedStats, addedPseudoStats := AddEquipSpellStats(&equipStats, equipPseudoStats, baseEff.SpellID)

			if len(props.Stats) > 0 || addedStats || addedPseudoStats {
				// An item the client ships no stat row for came from our sim's database
				// (forever_sim_db.json), whose stats already include its equip spells. Its hit
				// and crit are the exception: where the client's own equip spell states them,
				// that spell replaces master's Classic value (TakeClientHitAndCrit).
				if _, fromClient := dbcInstance.Items[int(parsed.Id)]; !fromClient {
					TakeClientHitAndCrit(parsed, equipPseudoStats, pseudoStats)
					continue
				}
				if areaType := spelldata.AreaTypeOfGroup(dbcInstance.Spells[baseEff.SpellID].RequiredAreasID); areaType != proto.AreaType_AreaTypeUnknown {
					addAreaStats(parsed.ScalingOptions[0], areaType, equipStats.ToProtoMap())
					continue
				}
				for pseudoStat, value := range equipPseudoStats {
					pseudoStats[pseudoStat] += value
				}
				for stat, value := range equipStats.ToProtoMap() {
					parsed.ScalingOptions[0].Stats[stat] += value
				}
				continue
			}
		} else if baseEff.TriggerType != ITEM_SPELLTRIGGER_CHANCE_ON_HIT && baseEff.CoolDownMSec <= 0 {
			continue
		}

		statsSpellID := resolveStatsSpell(baseEff.SpellID)
		pe := makeBaseProto(baseEff, statsSpellID)
		assignTrigger(baseEff, statsSpellID, pe)

		// add scaling for each saved state
		for state, opt := range parsed.ScalingOptions {
			ilvl := int(opt.Ilvl)
			pe.ScalingOptions[state] = buildItemEffectScalingProps(baseEff.SpellID, ilvl)
		}

		// A container aura that accumulates a separate stat aura resolves its amounts at every
		// state too, and per stack rather than in total.
		if stacking := buildStackingAura(baseEff.SpellID, statsSpellID, int(parsed.ScalingOptions[0].Ilvl)); stacking != nil {
			for state, opt := range parsed.ScalingOptions {
				stacking.Aura.ScalingOptions[state] = buildItemEffectScalingProps(int(stacking.Aura.BuffId), int(opt.Ilvl))
			}
			applyStackingAura(pe, stacking)
		}

		dropEmptyScalingOptions(pe)

		// Appended once per effect. Inside the loop above it appended the same pointer once
		// per scaling state, which only stays invisible while items carry a single state.
		effects = append(effects, pe)
	}

	parsed.PseudoStats = NullFloat(pseudoStats)
	return effects
}

// The two halves of an item's hit and crit: the pseudo stats a client equip spell states them in,
// and the rating stats master's database states them in.
var clientHitAndCrit = []struct {
	pseudoStats []proto.PseudoStat
	ratings     []proto.Stat
}{
	{
		[]proto.PseudoStat{proto.PseudoStat_PseudoStatMeleeHitPercent, proto.PseudoStat_PseudoStatRangedHitPercent, proto.PseudoStat_PseudoStatSpellHitPercent},
		[]proto.Stat{proto.Stat_StatMeleeHitRating, proto.Stat_StatSpellHitRating},
	},
	{
		[]proto.PseudoStat{proto.PseudoStat_PseudoStatMeleeCritPercent, proto.PseudoStat_PseudoStatRangedCritPercent, proto.PseudoStat_PseudoStatSpellCritPercent},
		[]proto.Stat{proto.Stat_StatMeleeCritRating, proto.Stat_StatSpellCritRating},
	},
}

// Forever revised some of the Classic-era items the client ships no stat row for yet, and gave them
// new equip spells that state their attack type: Blackhand's Breadth's 1318954 is "Increased Critical
// 1 - Melee and Ranged" (A_MOD_WEAPON_CRIT_PERCENT 1), Eye of the Beast's 1318956 "Increased
// Critical 1 - Spell" (A_MOD_SPELL_CRIT_CHANCE 1), and Trueaim Gauntlets' 1301127 +3% hit with bows,
// guns and crossbows. Master's database still has these items at Classic's size and as rating
// (Blackhand's Breadth 28 crit rating, 2%), which UnifyGearHitAndCrit then pays into every pool.
//
// Where an equip spell of such an item states hit, or crit, the item's rating of that kind (Classic
// states all hit and crit on an equip spell, never in the item row) is dropped and the spell's
// pseudo stats are added to into, the item's pseudo stats: percents of the type the spell names,
// which the sim does not unify. Returns whether the spell stated either.
func TakeClientHitAndCrit(item *proto.UIItem, equipPseudoStats []float64, into []float64) bool {
	took := false
	for _, kind := range clientHitAndCrit {
		states := false
		for _, pseudoStat := range kind.pseudoStats {
			if equipPseudoStats[pseudoStat] != 0 {
				states = true
			}
		}
		if !states {
			continue
		}
		took = true
		for _, opt := range item.ScalingOptions {
			for _, rating := range kind.ratings {
				delete(opt.Stats, int32(rating))
			}
		}
		for _, pseudoStat := range kind.pseudoStats {
			into[pseudoStat] += equipPseudoStats[pseudoStat]
		}
	}
	return took
}
