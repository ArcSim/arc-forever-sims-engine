package core

import (
	"fmt"
	"time"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// Scrolls share StatBuffCategory with the raid buffs granting the same stat, so a
// scroll doesn't stack with e.g. Arcane Brilliance or Divine Spirit; only the
// strongest source of that stat applies.
func registerScrollAura(character *Character, label string, itemID int32, stat stats.Stat, amount float64) *Aura {
	aura := character.GetOrRegisterAura(Aura{
		Label:      label,
		ActionID:   ActionID{ItemID: itemID},
		Duration:   NeverExpires,
		BuildPhase: CharacterBuildPhaseConsumes,
	})
	makeExclusiveFlatStatBuff(aura, stat, amount, StatBuffCategory)
	return MakePermanent(aura)
}

// Registers all consume-related effects to the Agent.
func applyConsumeEffects(agent Agent, _ *proto.PartyBuffs) {
	character := agent.GetCharacter()
	consumables := character.Consumables
	if consumables == nil {
		return
	}

	if consumables.FlaskId != 0 {
		flask := GetConsumableByID(consumables.FlaskId)
		character.AddStats(consumableStats(flask.Stats))
	}

	if consumables.BattleElixirId != 0 {
		// Elixir of Demonslaying
		if consumables.BattleElixirId == 9224 {
			character.Env.RegisterPostFinalizeEffect(func() {
				for _, at := range character.AttackTables {
					at.MobTypeBonusStats[proto.MobType_MobTypeDemon] = at.MobTypeBonusStats[proto.MobType_MobTypeDemon].Add(stats.Stats{
						stats.AttackPower:       265,
						stats.RangedAttackPower: 265,
					})
				}
			})
		} else {
			elixir := GetConsumableByID(consumables.BattleElixirId)
			character.AddStats(consumableStats(elixir.Stats))
		}
	}

	if consumables.GuardianElixirId != 0 {
		// Gift of Arthas
		if consumables.GuardianElixirId == 9088 {
			character.AddStat(stats.ShadowResistance, 10)
			auras := character.NewEnemyAuraArray(func(target *Unit) *Aura {
				return registeredBuffs().GiftOfArthasAura(target)
			})
			procSpell := character.RegisterSpell(SpellConfig{
				ActionID:    ActionID{SpellID: 11374},
				SpellSchool: SpellSchoolNature,
				ProcMask:    ProcMaskEmpty,

				FlatThreatBonus: 90,

				ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
					spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
					auras.Get(target).Activate(sim)
				},
			})

			character.MakeProcTriggerAura(ProcTrigger{
				Name:       "Gift of Arthas - Trigger",
				ICD:        time.Second * 3,
				ProcChance: 0.3,
				Outcome:    OutcomeLanded,
				Callback:   CallbackOnSpellHitTaken,
				Handler: func(sim *Simulation, spell *Spell, _ *SpellResult) {
					procSpell.Cast(sim, spell.Unit)
				},
			})
		} else {
			elixir := GetConsumableByID(consumables.GuardianElixirId)
			character.AddStats(consumableStats(elixir.Stats))
		}
	}
	if consumables.FoodId != 0 {
		food := GetConsumableByID(consumables.FoodId)
		character.AddStats(consumableStats(food.Stats))
	}

	// Classic buffs that stack beside the elixirs: jujus, Blasted Lands/Zanza, alcohol, the
	// school power and armor elixirs. Their stats come from the client like every other consumable.
	for _, id := range []int32{consumables.StrengthBuffId, consumables.AttackPowerBuffId, consumables.ZanzaId,
		consumables.AlcoholId, consumables.SpellPowerElixirId, consumables.SchoolElixirId, consumables.DefenseElixirId} {
		if id != 0 {
			character.AddStats(consumableStats(GetConsumableByID(id).Stats))
		}
	}
	if consumables.DragonbreathChili {
		registerDragonbreathChili(character)
	}

	// Static Imbues. Forever's Windfury Totem is a party aura, not Era's main-hand enchant (#550), so
	// it no longer displaces a main-hand stone or oil.
	// A stone or oil of its own on the item swap's main-hand weapon (mythicsim patch 133).
	if !registerPerWeaponMainHandImbue(agent, consumables) && consumables.MhImbueId != 0 {
		registerStaticImbue(agent, consumables.MhImbueId, character.AutoAttacks.MH(), ProcMaskMeleeMH)
	}
	if consumables.OhImbueId != 0 {
		registerStaticImbue(agent, consumables.OhImbueId, character.AutoAttacks.OH(), ProcMaskMeleeOH)
	}
	// A stone or oil beside a class imbue on the same hand (mythicsim patch 100).
	if consumables.MhTempEnchantId != 0 {
		registerStaticImbue(agent, consumables.MhTempEnchantId, character.AutoAttacks.MH(), ProcMaskMeleeMH)
	}
	if consumables.OhTempEnchantId != 0 {
		registerStaticImbue(agent, consumables.OhTempEnchantId, character.AutoAttacks.OH(), ProcMaskMeleeOH)
	}

	// Scrolls
	if consumables.ScrollAgi {
		registerScrollAura(character, "Scroll of Agility IV", 10309, stats.Agility, 17)
	}
	if consumables.ScrollStr {
		registerScrollAura(character, "Scroll of Strength IV", 10310, stats.Strength, 17)
	}
	if consumables.ScrollInt {
		registerScrollAura(character, "Scroll of Intellect IV", 10308, stats.Intellect, 16)
	}
	if consumables.ScrollSpi {
		registerScrollAura(character, "Scroll of Spirit IV", 10306, stats.Spirit, 15)
	}
	if consumables.ScrollArm {
		registerScrollAura(character, "Scroll of Protection IV", 10305, stats.Armor, 240)
	}

	// Bogling Root: +1 physical damage for 10 min (item 5206, spell 5665).
	if consumables.BoglingRoot {
		character.AddStat(stats.PhysicalDamage, 1)
	}

	explosivesSharedTimer := character.NewTimer()

	registerPotionCD(agent, consumables)
	registerConjuredCD(agent, consumables)
	registerExplosivesCD(agent, consumables, explosivesSharedTimer)
}

// ConsumableCritIsAllCrit (mythicsim patch 166, default true): a consumable's "chance to critically
// hit" counts for spells as well as attacks. Forever's client moved the Elixir of the Mongoose from
// Era's A_MOD_CRIT_PERCENT (52, weapon crit; client 1.15.9) to A_MOD_CRIT_PCT (290), and every
// Forever crit consumable states 290: the Elixirs of the Owl, Sages, Whale and Grizzly and the
// Grilled Squid's well fed buff. The engine reads 290 as crit with all spells and attacks
// elsewhere (Elune's Light, racials.go), while the database reads it as weapon crit only, so the
// spell half is added here. False keeps weapon crit only.
var ConsumableCritIsAllCrit = true

// consumableStats returns a stat consumable's stats as they apply: its crit, which the database
// holds as weapon crit, also counts as spell crit (ConsumableCritIsAllCrit).
func consumableStats(s stats.Stats) stats.Stats {
	if ConsumableCritIsAllCrit && s[stats.MeleeCritRating] != 0 && s[stats.SpellCritRating] == 0 {
		s[stats.SpellCritRating] = s[stats.MeleeCritRating] / PhysicalCritRatingPerCritPercent * SpellCritRatingPerCritPercent
	}
	return s
}

// Dragonbreath Chili (12217): its aura (15852) has a 5% chance, 10 s cooldown, on landed melee
// hits to cast 15851, 65 Fire damage (+-12.3%, SP coefficient 1) on every enemy - all client values.
func registerDragonbreathChili(character *Character) {
	procSpell := character.RegisterSpell(SpellConfig{
		ActionID:         ActionID{SpellID: 15851},
		SpellSchool:      SpellSchoolFire,
		DefenseType:      DefenseTypeMagic,
		ProcMask:         ProcMaskSpellDamageProc,
		Flags:            SpellFlagNoOnCastComplete | SpellFlagPassiveSpell,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,
		ApplyEffects: func(sim *Simulation, _ *Unit, spell *Spell) {
			for _, aoeTarget := range sim.Encounter.ActiveTargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, sim.Roll(57, 73), spell.OutcomeMagicHitAndCrit)
			}
		},
	})
	character.MakeProcTriggerAura(ProcTrigger{
		Name:       "Dragonbreath Chili",
		ActionID:   ActionID{SpellID: 15852},
		Callback:   CallbackOnSpellHitDealt,
		ProcMask:   ProcMaskMelee,
		Outcome:    OutcomeLanded,
		ProcChance: 0.05,
		ICD:        time.Second * 10,
		Handler: func(sim *Simulation, _ *Spell, result *SpellResult) {
			procSpell.Cast(sim, result.Target)
		},
	})
}

var PotionAuraTag = "Potion"

func registerPotionCD(agent Agent, consumes *proto.ConsumesSpec) {
	character := agent.GetCharacter()
	defaultPotion := consumes.PotId

	for _, potionId := range consumes.Potions {
		potion := GetConsumableByID(potionId)
		if potion.Type == proto.ConsumableType_ConsumableTypePotion {
			potMCD := makePotionActivationSpell(potion.Id, character)
			if defaultPotion == potion.Id {
				potMCD.Spell.Flags |= SpellFlagCombatPotion
				character.AddMajorCooldown(potMCD)
			}
		}
	}
}

// Empty: the eight ids this carried are MoP-era alchemist stones and none of them is in
// this client's item database, so HasAlchStone was always false. The lookup is kept for
// whatever Forever's equivalent turns out to be.
var AlchStoneItemIDs = []int32{}

func (character *Character) HasAlchStone() bool {
	alchStoneEquipped := false
	for _, itemID := range AlchStoneItemIDs {
		alchStoneEquipped = alchStoneEquipped || character.HasTrinketEquipped(itemID)
	}
	return character.HasProfession(proto.Profession_Alchemy) && alchStoneEquipped
}

func makePotionActivationSpell(potionId int32, character *Character) MajorCooldown {
	potion := GetConsumableByID(potionId)
	categoryCooldownDuration := TernaryDuration(potion.CategoryCooldownDuration > 0, potion.CategoryCooldownDuration, time.Minute*2)
	mcd := makePotionActivationSpellInternal(potion, character)

	if mcd.Spell != nil {
		// Mark as 'Encounter Only' so that users are forced to select the generic Potion
		// placeholder action instead of specific potion spells, in APL prepull. This
		// prevents a mismatch between Consumes and Rotation settings.
		mcd.Spell.Flags |= SpellFlagEncounterOnly | SpellFlagPotion | SpellFlagAPL
		oldApplyEffects := mcd.Spell.ApplyEffects
		mcd.Spell.ApplyEffects = func(sim *Simulation, target *Unit, spell *Spell) {
			oldApplyEffects(sim, target, spell)
			if sim.CurrentTime < 0 {
				spell.SharedCD.Set(sim.CurrentTime + categoryCooldownDuration)
				character.UpdateMajorCooldowns()
			}
		}
	}

	return mcd
}

type resourceGainConfig struct {
	resType  proto.ResourceType
	min      float64
	spread   float64
	period   time.Duration // Duration between ticks; 0 for one-shot gains.
	duration time.Duration // Total duration of periodic gains.
}

func makePotionActivationSpellInternal(potion Consumable, character *Character) MajorCooldown {
	stoneMul := TernaryFloat64(character.HasAlchStone(), 1.4, 1.0)
	cooldownDuration := TernaryDuration(potion.CooldownDuration > 0, potion.CooldownDuration, time.Minute*2)

	potionCast := CastConfig{
		CD: Cooldown{
			Timer:    character.NewTimer(),
			Duration: cooldownDuration,
		},
		SharedCD: Cooldown{
			Timer:    character.GetPotionCD(),
			Duration: cooldownDuration,
		},
	}

	actionID := ActionID{ItemID: potion.Id}
	var aura *StatBuffAura
	mcd := MajorCooldown{
		Spell: character.GetOrRegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast:     potionCast,
		}),
	}
	if effect, relatedAura := character.foreverPotionEffect(potion.Id); effect != nil {
		mcd.Type = CooldownTypeDPS
		mcd.Spell.RelatedSelfBuff = relatedAura
		mcd.Spell.ApplyEffects = func(sim *Simulation, target *Unit, _ *Spell) {
			if target == nil || target == &character.Unit {
				target = character.CurrentTarget
			}
			effect(sim, target)
		}
		return mcd
	}
	if potion.BuffDuration > 0 {
		// Add stat buff aura if applicable
		aura = character.NewTemporaryStatsAura(potion.Name, actionID, potion.Stats, potion.BuffDuration)
		mcd.Spell.RelatedSelfBuff = aura.Aura
		mcd.Type = aura.InferCDType()
		mcd.BuffAura = aura
	}
	var gains []resourceGainConfig
	resourceMetrics := make(map[proto.ResourceType]*ResourceMetrics)

	// Stats applied by triggered auras (e.g. Fel Mana Potion's spell damage reduction).
	// These may be positive or negative.
	var auraStats stats.Stats
	var auraDuration time.Duration
	var auraSpellId int32

	for _, effectID := range potion.EffectIds {
		e := GetSpellEffectByID(effectID)
		resourceType := e.GetResourceType()
		isPeriodic := e.AuraPeriodMs > 0
		if resourceType != 0 && (isPeriodic || e.Type == proto.EffectType_EffectTypeResourceGain) {
			if resourceType == proto.ResourceType_ResourceTypeMana && mcd.Type != CooldownTypeSurvival {
				mcd.Type = CooldownTypeMana
			} else if resourceType == proto.ResourceType_ResourceTypeHealth {
				mcd.Type = CooldownTypeSurvival
			} else {
				mcd.Type = CooldownTypeDPS
			}
			gains = append(gains, resourceGainConfig{
				resType:  resourceType,
				min:      e.MinEffectSize,
				spread:   e.EffectSpread,
				period:   time.Duration(e.AuraPeriodMs) * time.Millisecond,
				duration: time.Duration(e.DurationMs) * time.Millisecond,
			})
			if _, exists := resourceMetrics[resourceType]; !exists {
				resourceMetrics[resourceType] = character.Metrics.NewResourceMetrics(actionID, resourceType)
			}
			// Preload resource types that are found on this item
			if resourceMetrics[resourceType] == nil {
				resourceMetrics[resourceType] = character.Metrics.NewResourceMetrics(actionID, resourceType)
			}
			continue
		}
		if effectStats := stats.FromProtoArray(e.Stats); effectStats != (stats.Stats{}) {
			auraStats = auraStats.Add(effectStats)
			auraDuration = max(auraDuration, time.Duration(e.DurationMs)*time.Millisecond)
			if auraSpellId == 0 {
				auraSpellId = e.SpellId
			}
		}
	}

	var debuffAura *StatBuffAura
	if auraDuration > 0 && auraStats != (stats.Stats{}) {
		debuffAura = character.NewTemporaryStatsAura(fmt.Sprintf("%s Debuff (%d)", potion.Name, auraSpellId), ActionID{SpellID: auraSpellId}, auraStats, auraDuration)
	}

	mcd.Spell.ApplyEffects = func(sim *Simulation, _ *Unit, _ *Spell) {
		if aura != nil {
			aura.Activate(sim)
		}
		if debuffAura != nil {
			debuffAura.Activate(sim)
		}
		for _, config := range gains {
			gain := config.min + sim.RandomFloat(potion.Name)*config.spread
			gain *= stoneMul
			if config.period > 0 {
				// Periodic gains (e.g. Fel Mana Potion) tick over the effect's duration.
				startPeriodicResourceGain(sim, character, config, gain, resourceMetrics[config.resType])
			} else {
				if config.resType == proto.ResourceType_ResourceTypeHealth {
					gain *= character.PseudoStats.HealingTakenMultiplier
				}
				character.ExecuteResourceGain(sim, config.resType, gain, resourceMetrics[config.resType])
			}
		}
	}

	mcd.ShouldActivate = func(sim *Simulation, character *Character) bool {
		shouldActivate := true
		for _, config := range gains {
			switch config.resType {
			case proto.ResourceType_ResourceTypeMana:
				totalRegen := character.ManaRegenPerSecondWhileCasting() * 5
				manaGain := config.min + config.spread
				manaGain *= stoneMul
				if config.period > 0 {
					manaGain *= float64(startPeriodicResourceGainTicks(config))
				}
				shouldActivate = character.MaxMana()-(character.CurrentMana()+totalRegen) >= manaGain
			}
		}
		return shouldActivate
	}

	return mcd

}

// startPeriodicResourceGain ticks a periodic resource gain (e.g. Fel Mana Potion's mana over
// time) across the effect's duration.
func startPeriodicResourceGain(sim *Simulation, character *Character, config resourceGainConfig, gainPerTick float64, metrics *ResourceMetrics) {
	if config.resType != proto.ResourceType_ResourceTypeMana && config.resType != proto.ResourceType_ResourceTypeHealth {
		return
	}
	if config.resType == proto.ResourceType_ResourceTypeHealth {
		gainPerTick *= character.PseudoStats.HealingTakenMultiplier
	}
	StartPeriodicAction(sim, PeriodicActionOptions{
		Period:   config.period,
		NumTicks: int(startPeriodicResourceGainTicks(config)),
		Priority: ActionPriorityDOT,
		OnAction: func(sim *Simulation) {
			character.ExecuteResourceGain(sim, config.resType, gainPerTick, metrics)
		},
	})
}

func startPeriodicResourceGainTicks(config resourceGainConfig) int32 {
	return int32(config.duration / config.period)
}

var ConjuredAuraTag = "Conjured"

func registerConjuredCD(agent Agent, consumes *proto.ConsumesSpec) {
	character := agent.GetCharacter()

	for _, conjuredId := range consumes.ConjuredItems {
		// The UI sends its whole eligible list, unfiltered by the consumable database.
		if GetConsumableByID(conjuredId).Id == 0 {
			continue
		}

		conjuredMCD := makeConjuredActivationSpell(conjuredId, character)

		if conjuredMCD.Spell != nil {
			oldShouldActivate := conjuredMCD.ShouldActivate
			conjuredMCD.ShouldActivate = func(sim *Simulation, character *Character) bool {
				return oldShouldActivate(sim, character) && consumes.ConjuredId == conjuredId
			}
			character.AddMajorCooldown(conjuredMCD)
		}
	}

}

func makeConjuredActivationSpell(conjuredId int32, character *Character) MajorCooldown {
	conjured := GetConsumableByID(conjuredId)
	categoryCooldownDuration := TernaryDuration(conjured.CategoryCooldownDuration > 0, conjured.CategoryCooldownDuration, time.Minute*2)
	mcd := makeConjuredActivationSpellInternal(conjured, character)

	if mcd.Spell != nil {
		mcd.Spell.Flags |= SpellFlagConjured | SpellFlagAPL
		oldApplyEffects := mcd.Spell.ApplyEffects
		mcd.Spell.ApplyEffects = func(sim *Simulation, target *Unit, spell *Spell) {
			oldApplyEffects(sim, target, spell)
			if sim.CurrentTime < 0 {
				spell.SharedCD.Set(sim.CurrentTime + categoryCooldownDuration)
				character.UpdateMajorCooldowns()
			}
		}
	}

	return mcd
}

func makeConjuredActivationSpellInternal(conjured Consumable, character *Character) MajorCooldown {
	cooldownDuration := TernaryDuration(conjured.CooldownDuration > 0, conjured.CooldownDuration, time.Minute*2)

	conjuredCast := CastConfig{
		CD: Cooldown{
			Timer:    character.NewTimer(),
			Duration: cooldownDuration,
		},
		SharedCD: Cooldown{
			Timer:    character.GetOrInitSpellCategoryTimer(conjured.CategoryId),
			Duration: time.Minute * 2,
		},
	}

	actionID := ActionID{ItemID: conjured.Id}
	var aura *StatBuffAura
	mcd := MajorCooldown{
		Spell: character.GetOrRegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast:     conjuredCast,
		}),
	}
	if conjured.BuffDuration > 0 {
		// Add stat buff aura if applicable
		aura = character.NewTemporaryStatsAura(conjured.Name, actionID, conjured.Stats, conjured.BuffDuration)
		mcd.Spell.RelatedSelfBuff = aura.Aura
		mcd.Type = aura.InferCDType()
	}
	var gains []resourceGainConfig
	resourceMetrics := make(map[proto.ResourceType]*ResourceMetrics)

	for _, effectID := range conjured.EffectIds {
		e := GetSpellEffectByID(effectID)
		resourceType := e.GetResourceType()
		if (e.Type == proto.EffectType_EffectTypeResourceGain || e.Type == proto.EffectType_EffectTypeHeal) && resourceType != 0 {
			if resourceType == proto.ResourceType_ResourceTypeMana && mcd.Type != CooldownTypeSurvival {
				mcd.Type = CooldownTypeMana
			} else if resourceType == proto.ResourceType_ResourceTypeHealth {
				mcd.Type = CooldownTypeSurvival
			} else {
				mcd.Type = CooldownTypeDPS
			}
			gains = append(gains, resourceGainConfig{
				resType: resourceType,
				min:     e.MinEffectSize,
				spread:  e.EffectSpread,
			})

			if _, exists := resourceMetrics[resourceType]; !exists {
				resourceMetrics[resourceType] = character.Metrics.NewResourceMetrics(actionID, resourceType)
			}
			// Preload resource types that are found on this item
			if resourceMetrics[resourceType] == nil {
				resourceMetrics[resourceType] = character.Metrics.NewResourceMetrics(actionID, resourceType)
			}
		}
	}

	mcd.Spell.ApplyEffects = func(sim *Simulation, _ *Unit, _ *Spell) {
		if aura != nil {
			aura.Activate(sim)
		}

		for _, config := range gains {
			gain := config.min + TernaryFloat64(config.spread > 1, sim.RandomFloat(conjured.Name)*config.spread, config.spread)
			switch config.resType {
			case proto.ResourceType_ResourceTypeHealth:
				gain *= character.PseudoStats.HealingTakenMultiplier
			}
			character.ExecuteResourceGain(sim, config.resType, gain, resourceMetrics[config.resType])
		}
	}

	mcd.ShouldActivate = func(sim *Simulation, character *Character) bool {
		shouldActivate := true
		for _, config := range gains {
			switch config.resType {
			case proto.ResourceType_ResourceTypeMana:
				totalRegen := character.ManaRegenPerSecondWhileCasting() * 5
				manaGain := config.min + config.spread
				shouldActivate = character.MaxMana()-(character.CurrentMana()+totalRegen) >= manaGain
			case proto.ResourceType_ResourceTypeEnergy:
				// Thistle Tea (9512) restores a flat 100: the client row has no level scaling. Waiting for
				// room for all of it would mean waiting for an empty bar, so some may spill. Energy refills
				// at 10 a second, so a bar is rarely under 10 when a 1 s global cooldown ends: with 10 allowed
				// to spill a 100-energy Combat Rogue drank in 0.3% of fights and a Gnome's 105 in all of them.
				// 15 lets every Rogue drink once a fight, and measured best of 10, 15, 20 and 30.
				shouldActivate = character.MaximumEnergy()-character.CurrentEnergy() >= config.min+config.spread-ThistleTeaSpill
			}
		}
		return shouldActivate
	}

	return mcd

}

// The energy Thistle Tea may waste when it is used without being asked for (MythicSim patch 102).
const ThistleTeaSpill = 15

var GoblinSapperActionID = ActionID{ItemID: 10646}
var EzThroDynamiteTwoActionID = ActionID{ItemID: 18588}
var CrystalChargeActionID = ActionID{ItemID: 11566}
var ThoriumGrenadeActionID = ActionID{ItemID: 15993}
var DenseDynamiteActionID = ActionID{ItemID: 18641}
var CryoblastActionID = ActionID{ItemID: 217495}

// Scroll of Cryoblast's use spell, the value the explosives picker saves for it.
const CryoblastSpellID = 440212

// Forever's SAF-T / EZ-Thro bombs, keyed on their use spell (the picker's value). ItemSparse
// (70009) gives them required level 1 and no RequiredSkill, so anyone can throw them. Each is a
// 1 s cast of Fire damage shared on the 1 min explosives category; damage is the client's base
// points +- half its variance, with no spell power coefficient. EZ-Thro Bronze Mortar (own 10 min
// cooldown) and EZ-Thro Mana Bomb (a mana burn) are left out.
var saftBombs = map[int32]struct {
	itemID          int32
	min, max, speed float64
}{
	1269161: {260793, 22, 28, 14},   // SAF-T Copper Bomb: 25, variance .24
	1269155: {260792, 26, 34, 14},   // SAF-T Dynamite: 30, .267
	1269192: {260795, 43, 57, 14},   // EZ-Thro Copper Bomb XL: 50, .28
	1269216: {260797, 73, 97, 14},   // SAF-T Bronze Bomb: 85, .282
	1269264: {260798, 128, 172, 14}, // SAF-T Jumbo Dynamite: 150, .293
	1269272: {260805, 149, 201, 14}, // SAF-T Bomb: 175, .297
	1269278: {260809, 149, 201, 14}, // Tru-Trigger Frag Bomb: 175, .297
	1269270: {260803, 213, 287, 14}, // EZ-Thro Grenade: 250, .296
	1269282: {260814, 340, 460, 14}, // SAF-T Clever Dynamite: 400, .3
	1269330: {260816, 300, 500, 25}, // EZ-Thro Thorium Grenade: 400, .5
	1269334: {260817, 225, 675, 14}, // EZ-Thro Dark Bomb: 450, 1
}

func registerExplosivesCD(agent Agent, consumes *proto.ConsumesSpec, sharedTimer *Timer) {
	character := agent.GetCharacter()
	engineer := character.HasProfession(proto.Profession_Engineering)

	if consumes.GoblinSapper && engineer {
		character.AddMajorCooldown(MajorCooldown{
			Spell:    character.newGoblinSapperSpell(sharedTimer),
			Type:     CooldownTypeDPS | CooldownTypeExplosive,
			Priority: CooldownPriorityLow + 20,
		})
	}

	var filler *Spell
	bomb, isSaftBomb := saftBombs[consumes.ExplosiveId]
	switch {
	case isSaftBomb:
		filler = character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, ActionID{ItemID: bomb.itemID}, SpellSchoolFire, bomb.min, bomb.max, bomb.speed, time.Second, Cooldown{}))
	case consumes.ExplosiveId == CryoblastSpellID:
		// A mage's vendor scroll, not an engineer's bomb, but it shares their 1 min cooldown.
		if character.Class == proto.Class_ClassMage {
			filler = character.newCryoblastSpell(sharedTimer)
		}
	// Ez-Thro Dynamite II and Crystal Charge state no RequiredSkill in ItemSparse (70009): anyone can throw them.
	case consumes.ExplosiveId == 18588:
		filler = character.newEzThroDynamiteTwoSpell(sharedTimer)
	case consumes.ExplosiveId == 15239:
		filler = character.newCrystalChargeSpell(sharedTimer)
	case !engineer:
	case consumes.ExplosiveId == 19769:
		filler = character.newThoriumGrenadeSpell(sharedTimer)
	case consumes.ExplosiveId == 23063, consumes.ExplosiveId == 18641: // 18641: the item id Forever saved before the merge
		filler = character.newDenseDynamiteSpell(sharedTimer)
	}
	if filler != nil {
		character.AddMajorCooldown(MajorCooldown{
			Spell:    filler,
			Type:     CooldownTypeDPS | CooldownTypeExplosive,
			Priority: CooldownPriorityLow + 10,
		})
	}
}

// Creates a spell object for the common explosive case.
func (character *Character) newBasicExplosiveSpellConfig(sharedTimer *Timer, actionID ActionID, school SpellSchool, minDamage float64, maxDamage float64, speed float64, castTime time.Duration, cooldown Cooldown) SpellConfig {
	var selfDamage *Spell
	if actionID.SameAction(GoblinSapperActionID) {
		selfDamage = character.newSapperSelfDamageSpell(actionID, school)
	}

	return SpellConfig{
		ActionID:     actionID,
		SpellSchool:  school,
		DefenseType:  DefenseTypeMagic, // Every explosive's damage spell is Magic in SpellCategories, so they crit for 150%
		ProcMask:     ProcMaskEmpty,
		Flags:        SpellFlagExplosive,
		MissileSpeed: speed,

		Cast: CastConfig{
			DefaultCast: Cast{
				CastTime: castTime,
			},
			CD: cooldown,
			SharedCD: Cooldown{
				Timer:    sharedTimer,
				Duration: time.Minute,
			},
		},

		// Explosives always have 1% resist chance, so just give them hit cap.
		BonusHitPercent:  100,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			baseDamage := sim.Roll(minDamage, maxDamage) * sim.Encounter.AOECapMultiplier()
			spell.CalcAoeDamage(sim, baseDamage, spell.OutcomeMagicHitAndCrit)
			if speed > 0 {
				spell.WaitTravelTime(sim, func(sim *Simulation) {
					spell.DealBatchedAoeDamage(sim)
				})
			} else {
				spell.DealBatchedAoeDamage(sim)
			}

			if selfDamage != nil {
				baseDamage := sim.Roll(minDamage, maxDamage)
				selfDamage.CalcAndDealDamage(sim, &character.Unit, baseDamage, selfDamage.OutcomeMagicHitAndCrit)
			}
		},
	}
}

// The half of a sapper charge that goes off in the thrower's face. Its own spell so that the hit
// carries the kind it is - a harmful spell landing on the character - which is what a listener on
// spell damage taken hears. The tag keeps it apart from the charge's outgoing damage, whose own
// hits state no kind.
func (character *Character) newSapperSelfDamageSpell(actionID ActionID, school SpellSchool) *Spell {
	return character.GetOrRegisterSpell(SpellConfig{
		ActionID:    actionID.WithTag(1),
		SpellSchool: school,
		DefenseType: DefenseTypeMagic,
		ProcMask:    ProcMaskSpellDamage,
		Flags:       SpellFlagExplosive,

		BonusHitPercent:  100,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
	})
}
func (character *Character) newGoblinSapperSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, GoblinSapperActionID, SpellSchoolFire, 450, 750, 0, 0, Cooldown{Timer: character.NewTimer(), Duration: time.Minute * 5}))
}
func (character *Character) newCrystalChargeSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, CrystalChargeActionID, SpellSchoolFire, 383, 517, 0, 0, Cooldown{}))
}
func (character *Character) newEzThroDynamiteTwoSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, EzThroDynamiteTwoActionID, SpellSchoolFire, 213, 287, 14, time.Second, Cooldown{}))
}
func (character *Character) newThoriumGrenadeSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, ThoriumGrenadeActionID, SpellSchoolFire, 300, 500, 25, time.Second, Cooldown{}))
}

// Client 1.60.1.70009: 215 Frost damage with 0.2977 variance (183 - 247) in 5 yards and no spell
// power coefficient. Beta logs agree: 23 non-crit hits of two level 20 mages average 209.
func (character *Character) newCryoblastSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, CryoblastActionID, SpellSchoolFrost, 183, 247, 0, 0, Cooldown{}))
}
func (character *Character) newDenseDynamiteSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, DenseDynamiteActionID, SpellSchoolFire, 340, 460, 14, time.Second, Cooldown{}))
}

func imbueFlatWeaponDamage(imbueId int32) float64 {
	switch imbueId {
	case 16138, 16622: // Dense Sharpening Stone / Dense Weightstone
		return 8
	}
	return 0
}

// Flat weapon damage the main-hand imbue adds, for classes that build their
// main-hand weapon from the equipped item.
func (character *Character) MHImbueFlatWeaponDamage() float64 {
	return imbueFlatWeaponDamage(character.Consumables.MhImbueId) + imbueFlatWeaponDamage(character.Consumables.MhTempEnchantId)
}

// StaticImbueProcs are the temporary weapon enchantments whose effect is a combat proc, by imbue ID
// (the item's use spell, or its item ID). sim/common registers them, since the proc's spell is built
// from the client row there; hand is the hand the imbue is on (ProcMaskMeleeMH or ProcMaskMeleeOH).
var StaticImbueProcs = map[int32]func(agent Agent, hand ProcMask){}

func registerStaticImbue(agent Agent, imbueId int32, weapon *Weapon, hand ProcMask) {
	character := agent.GetCharacter()
	if proc, ok := StaticImbueProcs[imbueId]; ok {
		proc(agent, hand)
	}
	switch imbueId {
	case 25123: // Mana Oil
		character.AddStat(stats.HealingPower, 30)
		character.AddStat(stats.MP5, 15)
	case 25122, 20749: // Brilliant Wizard Oil (20749: the item id Forever saved before the merge)
		character.AddStat(stats.SpellDamage, 36)
		character.AddStat(stats.HealingPower, 36)
		character.AddStat(stats.SpellCritPercent, 1)
	case 25121: // Wizard Oil: 24 in the client (enchant 2627, spell 25111), reverted 2026-09-24
		character.AddStat(stats.SpellDamage, 24)
	case 22756, 18262: // Elemental Sharpening Stone (18262: the item id Forever saved before the merge)
		// RangedCritPercent is the ranged offset from PhysicalCritPercent, so the melee-only
		// crit has to be cancelled there.
		character.AddStat(stats.PhysicalCritPercent, 2)
		character.AddStat(stats.RangedCritPercent, -2)
	case 28898: // Blessed Wizard Oil
		character.Env.RegisterPostFinalizeEffect(func() {
			for _, at := range character.AttackTables {
				at.MobTypeBonusStats[proto.MobType_MobTypeUndead] = at.MobTypeBonusStats[proto.MobType_MobTypeUndead].Add(stats.Stats{
					stats.SpellDamage: 58,
				})
			}
		})
	case 28891: // Consecrated Sharpening Stone
		character.Env.RegisterPostFinalizeEffect(func() {
			for _, at := range character.AttackTables {
				at.MobTypeBonusStats[proto.MobType_MobTypeUndead] = at.MobTypeBonusStats[proto.MobType_MobTypeUndead].Add(stats.Stats{
					stats.AttackPower:       100,
					stats.RangedAttackPower: 100,
				})
			}
		})
	}

	if flat := imbueFlatWeaponDamage(imbueId); flat != 0 {
		weapon.BaseDamageMin += flat
		weapon.BaseDamageMax += flat
	}
}
