package database

import (
	"regexp"

	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
	"github.com/wowsims/forever/tools/database/dbc"
)

// Allows you to ignore certain Spell Effects that the Sim does not support.
// This prevents them from being added to the missing effects list to prevent confusion towards users.
// Empty array means ignore all effects of that type, otherwise it will be ignored
// based on EffectMiscValue_0
var IgnoreSpellEffectByAuraType = map[dbc.EffectAuraType][]int{
	dbcenums.A_MOD_MECHANIC_RESISTANCE: {},
	dbcenums.A_MOD_STEALTH:             {},
	dbcenums.A_MOD_STEALTH_DETECT:      {},
	dbcenums.A_MOD_STEALTH_LEVEL:       {},
	dbcenums.A_MOD_INVISIBILITY:        {},
	dbcenums.A_MOD_INVISIBILITY_DETECT: {},
	dbcenums.A_MOD_SKILL: {
		356, // Fishing Skill
		393, // Skinning Skill
	},
	dbcenums.A_MOD_INCREASE_MOUNTED_SPEED:        {},
	dbcenums.A_MOD_MOUNTED_SPEED_ALWAYS:          {},
	dbcenums.A_MOD_MOUNTED_SPEED_NOT_STACK:       {},
	dbcenums.A_MOD_INCREASE_MOUNTED_FLIGHT_SPEED: {},
	dbcenums.A_MOD_MOUNTED_FLIGHT_SPEED_ALWAYS:   {},
	dbcenums.A_TRANSFORM:                         {},
	dbcenums.A_MECHANIC_IMMUNITY:                 {},
	dbcenums.A_TRACK_CREATURES:                   {},
	dbcenums.A_TRACK_RESOURCES:                   {},
	dbcenums.A_FAR_SIGHT:                         {},
}

// Auras the sim has nothing to simulate for, which put a spell out of scope only where they are all
// it carries: Frostguard's Chilled 16927 slows movement beside a melee slow.
var IgnoreSpellEffectAloneByAuraType = []dbc.EffectAuraType{
	dbcenums.A_MOD_DECREASE_SPEED,
}

var IgnoreSpellEffectBySpellEffectType = map[dbc.SpellEffectType][]int{
	dbcenums.E_CREATE_ITEM:    {},
	dbcenums.E_SUMMON:         {},
	dbcenums.E_TELEPORT_UNITS: {},
}

// Spells that are flavour rather than mechanics, keyed by spell ID so that reporting them as
// missing effects can be suppressed without hiding anything real.
//
// Keyed per spell on purpose. The obvious shortcut - treat an effect made only of A_DUMMY auras
// as noise, which is what MoP does - does not hold in TBC, where A_DUMMY is how a whole class of
// real item effects is encoded: every relic, idol, libram and totem bonus, plus Zandalarian Hero
// Medallion, Thick Obsidian Breastplate and Ashtongue Talisman of Zeal. That rule removed 48
// entries here and only a handful of them were noise.
var IgnoreMissingEffectBySpellID = map[int]string{
	16372: "Seal of Ascension - no tooltip and no mechanic",
}

// Absorbs a server script restricts to spells the client does not list, keyed by spell ID with the
// reason the generated file gives. 1287808, Onyxia Blood Talisman's, states 10000000000 beside the
// A_DUMMY that names the Dragon Breath spells.
var UnsupportedAbsorbBySpellID = map[int32]string{
	1287808: "the absorb of 10000000000 beside an A_DUMMY absorbs only the spells a script names, which the client does not list",
}

var OtherItemIdsToFetch = []string{}
var ConsumableOverrides = []*proto.Consumable{}

// Empty: all 26 entries were TBC items absent from this client.
var ItemOverrides = []*proto.UIItem{}

// Keep these sorted by item ID.
var ItemAllowList = map[int32]struct{}{
	1168:  {}, // Skullflame Shield
	2140:  {},
	2505:  {},
	8345:  {}, // Wolfshead Helm
	11815: {}, // Hand of Justice
	18168: {}, // Force Reactive Disk
	19337: {}, // The Black Book
}

// Keep these sorted by item ID.
var ItemDenyList = map[int32]struct{}{
	17782: {}, // talisman of the binding shard
	17783: {}, // talisman of the binding fragment
	17802: {}, // Deprecated version of Thunderfury
	18582: {},
	18583: {},
	18584: {},
	22736: {},
	23363: {}, // Titanic Breastplate
	24265: {},
}

// Item icons to include in the DB, so they don't need to be separately loaded in the UI.
var ExtraItemIcons = []int32{
	// Demonic Rune
	12662,

	// Flask IDs
	13512,

	// Elixer IDs
	9224,
	13452,
	13454,

	// Potions / In Battle Consumes
	13442,
	9421,

	// Thistle Tea
	7676,

	// Scrolls
	10305,
	10306,
	10308,
	10309,
	10310,

	// Weapon Imbues
	12404,
	12643,
	18262,
	20750,
	23123,

	// Explosives
	15993,
	18641,
}

// Item Ids of consumables to allow
var ConsumableAllowList = []int32{
	7676,  // Thisle Tea
	9088,  // Gift of Arthas
	9155,  // Arcane Elixir
	9224,  // Elixir of Demonslaying
	9421,  // Major Healthstone
	13442, // Migty Rage Potion
	13452, // Elixir of the Mongoose
	13454, // Greater Arcane Elixir
	12662, // Demonic Rune
	5206,  // Bogling Root
}

// Classic consumables and the picker each goes in. The client has no battle/guardian split for
// them, and most sit below the loader's level-50 floor, so without this list the sim silently
// treats them as empty. The set is master's consumables picker plus its world buffs.
var ClassicConsumableTypes = map[int32]proto.ConsumableType{
	// Potions
	1710: proto.ConsumableType_ConsumableTypePotion, 3928: proto.ConsumableType_ConsumableTypePotion,
	13446: proto.ConsumableType_ConsumableTypePotion, 3827: proto.ConsumableType_ConsumableTypePotion,
	6149: proto.ConsumableType_ConsumableTypePotion, 13443: proto.ConsumableType_ConsumableTypePotion,
	13444: proto.ConsumableType_ConsumableTypePotion, 13442: proto.ConsumableType_ConsumableTypePotion,
	5633: proto.ConsumableType_ConsumableTypePotion, 5631: proto.ConsumableType_ConsumableTypePotion,
	9036: proto.ConsumableType_ConsumableTypePotion, 13461: proto.ConsumableType_ConsumableTypePotion,
	13457: proto.ConsumableType_ConsumableTypePotion, 13456: proto.ConsumableType_ConsumableTypePotion,
	13460: proto.ConsumableType_ConsumableTypePotion, 13458: proto.ConsumableType_ConsumableTypePotion,
	13459: proto.ConsumableType_ConsumableTypePotion, 13455: proto.ConsumableType_ConsumableTypePotion,
	4623: proto.ConsumableType_ConsumableTypePotion,
	// Forever's Merchant's Favor potions below the loader's level-50 floor (Nina Surefire, the
	// alchemy vendor; client 1.60.1.70450): Potion of Venomous Blood (level 30), Dragonfire Potion
	// (35), Potion of Elemental Purging (40) and Caustic Smog Potion (45). The level 55 Major
	// Spellblasting, Frenzy and Mender's Potions load on their own (mythicsim patch 171).
	250950: proto.ConsumableType_ConsumableTypePotion, 250953: proto.ConsumableType_ConsumableTypePotion,
	250954: proto.ConsumableType_ConsumableTypePotion, 250951: proto.ConsumableType_ConsumableTypePotion,
	// Flasks
	13510: proto.ConsumableType_ConsumableTypeFlask, 13511: proto.ConsumableType_ConsumableTypeFlask,
	13512: proto.ConsumableType_ConsumableTypeFlask, 13513: proto.ConsumableType_ConsumableTypeFlask,
	// Food
	21023: proto.ConsumableType_ConsumableTypeFood, 13928: proto.ConsumableType_ConsumableTypeFood,
	20452: proto.ConsumableType_ConsumableTypeFood, 18254: proto.ConsumableType_ConsumableTypeFood,
	13810: proto.ConsumableType_ConsumableTypeFood, 13813: proto.ConsumableType_ConsumableTypeFood,
	13931: proto.ConsumableType_ConsumableTypeFood, 18045: proto.ConsumableType_ConsumableTypeFood,
	21217: proto.ConsumableType_ConsumableTypeFood, 13851: proto.ConsumableType_ConsumableTypeFood,
	21072: proto.ConsumableType_ConsumableTypeFood,
	// Battle elixir slot: agility elixirs
	13452: proto.ConsumableType_ConsumableTypeBattleElixir, 9187: proto.ConsumableType_ConsumableTypeBattleElixir,
	8949: proto.ConsumableType_ConsumableTypeBattleElixir, 3390: proto.ConsumableType_ConsumableTypeBattleElixir,
	// Forever's Elixir of Ferocity (18 Strength and Agility) and Elixir of Cunning (25 Agility and
	// Intellect). Forever's level 55 elixirs carry no Battle or Guardian Elixir text, so each goes
	// in the slot of the Classic elixir whose stat it shares (mythicsim patch 171).
	250350: proto.ConsumableType_ConsumableTypeBattleElixir, 250328: proto.ConsumableType_ConsumableTypeBattleElixir,
	// Spell power elixirs
	13454: proto.ConsumableType_ConsumableTypeSpellPowerElixir, 9155: proto.ConsumableType_ConsumableTypeSpellPowerElixir,
	250333: proto.ConsumableType_ConsumableTypeSpellPowerElixir, // Forever's Greater Cleric's Elixir (+40 healing)
	// Guardian elixir slot: health and mana regen (Forever adds the Greater Mageblood Elixir)
	3825: proto.ConsumableType_ConsumableTypeGuardianElixir, 2458: proto.ConsumableType_ConsumableTypeGuardianElixir,
	20007: proto.ConsumableType_ConsumableTypeGuardianElixir, 250341: proto.ConsumableType_ConsumableTypeGuardianElixir,
	// Forever's Elixir of Greater Fortitude (400 health), Elixir of Wicked Regeneration (9 health
	// and 10 mana every 5 sec), Elixir of the Owl (25 Intellect, 2% crit), Elixir of Sages and
	// Elixir of the Whale (25 Spirit, 2% crit each), and the Elixir of Greater Spirit (18 Spirit).
	250335: proto.ConsumableType_ConsumableTypeGuardianElixir, 250336: proto.ConsumableType_ConsumableTypeGuardianElixir,
	250337: proto.ConsumableType_ConsumableTypeGuardianElixir, 250338: proto.ConsumableType_ConsumableTypeGuardianElixir,
	250348: proto.ConsumableType_ConsumableTypeGuardianElixir, 250347: proto.ConsumableType_ConsumableTypeGuardianElixir,
	// Armor elixirs
	13445: proto.ConsumableType_ConsumableTypeDefenseElixir, 8951: proto.ConsumableType_ConsumableTypeDefenseElixir,
	3389: proto.ConsumableType_ConsumableTypeDefenseElixir, 5997: proto.ConsumableType_ConsumableTypeDefenseElixir,
	250329: proto.ConsumableType_ConsumableTypeDefenseElixir, // Forever's Elixir of the Phalanx (400 health, 500 armor)
	// School power elixirs
	21546: proto.ConsumableType_ConsumableTypeSchoolElixir, 6373: proto.ConsumableType_ConsumableTypeSchoolElixir,
	17708: proto.ConsumableType_ConsumableTypeSchoolElixir, 9264: proto.ConsumableType_ConsumableTypeSchoolElixir,
	250343: proto.ConsumableType_ConsumableTypeSchoolElixir, // Forever's Elixir of Nature Power
	// Strength: Juju Power, Elixir of Giants, Elixir of Ogre's Strength
	12451: proto.ConsumableType_ConsumableTypeStrengthBuff, 9206: proto.ConsumableType_ConsumableTypeStrengthBuff,
	3391:   proto.ConsumableType_ConsumableTypeStrengthBuff,
	250351: proto.ConsumableType_ConsumableTypeStrengthBuff, // Forever's Elixir of the Grizzly (25 Strength, 2% crit)
	// Attack power: Juju Might, Winterfall Firewater
	12460: proto.ConsumableType_ConsumableTypeAttackPowerBuff, 12820: proto.ConsumableType_ConsumableTypeAttackPowerBuff,
	// Blasted Lands and Zanza buffs
	8410: proto.ConsumableType_ConsumableTypeZanza, 8412: proto.ConsumableType_ConsumableTypeZanza,
	8423: proto.ConsumableType_ConsumableTypeZanza, 8424: proto.ConsumableType_ConsumableTypeZanza,
	8411: proto.ConsumableType_ConsumableTypeZanza, 20079: proto.ConsumableType_ConsumableTypeZanza,
	// Alcohol: Rumsey Rum Black Label/Dark/Light, Gordok Green Grog, Kreeg's Stout Beatdown
	21151: proto.ConsumableType_ConsumableTypeAlcohol, 21114: proto.ConsumableType_ConsumableTypeAlcohol,
	20709: proto.ConsumableType_ConsumableTypeAlcohol, 18269: proto.ConsumableType_ConsumableTypeAlcohol,
	18284: proto.ConsumableType_ConsumableTypeAlcohol,
}
var ConsumableDenyList = []int32{}

// Raid buffs / debuffs
var SharedSpellsIcons = []int32{
	// Registered CD's
	10060,
	16190,
	29166,
	2825,

	17051,

	25898,

	20140,
	8071,

	14767,

	8075,

	20045,

	30808,
	19506,

	12861,
	18696,

	20245,
	5675,
	16206,

	17007,

	8512,
	29193,

	24907,

	3738,
	8227,

	16840,

	// Raid Debuffs
	8647,

	770,
	702,
	18180,

	12879,
	16862,

	20337,

	3043,

	1490,

	20271,

	11374,
	15235,

	// Raid buffs, debuffs and shared consumables.
	603,   // Curse of Doom
	688,   // Summon Imp
	691,   // Summon Felhunter
	697,   // Summon Voidwalker
	704,   // Curse of Recklessness
	706,   // Demon Armor
	712,   // Summon Succubus
	6117,  // Mage Armor
	7302,  // Ice Armor
	8024,  // Flametongue Weapon
	8033,  // Frostbrand Weapon
	8232,  // Windfury Weapon
	10538, // Fire Resistance Totem
	12470, // Fire Nova
	13339, // Fire Blast
	13376, // Fire Shield
	13889, // Minor Speed
	14325, // Hunter's Mark
	17768, // Wolfshead Helm
	18803, // Focus
	19615, // Frenzy Effect
	20554, // Berserking
	20572, // Blood Fury
	20574, // Axe Specialization
	20594, // Stoneform
	20597, // Sword Specialization
	23110, // Dash
	23563, // Enhanced Battle Shout
	25076, // Cobra Reflexes
	21564, // Prayer of Fortitude
	25895, // Greater Blessing of Salvation
	25916, // Greater Blessing of Might
	25918, // Greater Blessing of Wisdom
	27681, // Prayer of Spirit
	27683, // Prayer of Shadow Protection
	26654, // Sweeping Strikes
	28142, // Power of the Guardian
	28143, // Power of the Guardian
	29414, // Haste

	1259686, // Skysight
	1259688, // Elemental Blessing
	1259705, // Read Ley Line
	1259719, // Mace Specialization
	1259799, // Elune's Light
	1259812, // Eureka!
	1259813, // Eureka!
	1259817, // Eureka!
	1259821, // Eureka!
	1259823, // Eureka!
	1260198, // Touch of the Grave
	1270842, // Energized
	1299026, // Shatter Curse
}

// If any of these match the item name, don't include it.
var DenyListNameRegexes = []*regexp.Regexp{
	regexp.MustCompile(`30 Epic`),
	regexp.MustCompile(`130 Epic`),
	regexp.MustCompile(`63 Blue`),
	regexp.MustCompile(`63 Green`),
	regexp.MustCompile(`66 Epic`),
	regexp.MustCompile(`90 Epic`),
	regexp.MustCompile(`90 Green`),
	regexp.MustCompile(`Boots 1`),
	regexp.MustCompile(`Boots 2`),
	regexp.MustCompile(`Boots 3`),
	regexp.MustCompile(`Bracer 1`),
	regexp.MustCompile(`Bracer 2`),
	regexp.MustCompile(`Bracer 3`),
	regexp.MustCompile(`DB\d`),
	regexp.MustCompile(`DEPRECATED`),
	regexp.MustCompile(`OLD`),
	regexp.MustCompile(`Deprecated`),
	regexp.MustCompile(`Deprecated: Keanna`),
	regexp.MustCompile(`Indalamar`),
	regexp.MustCompile(`Monster -`),
	regexp.MustCompile(`NEW`),
	regexp.MustCompile(`PH`),
	regexp.MustCompile(`QR XXXX`),
	regexp.MustCompile(`TEST`),
	regexp.MustCompile(`Test`),
	regexp.MustCompile(`Enchant Template`),
	regexp.MustCompile(`Arcane Amalgamation`),
	regexp.MustCompile(`Deleted`),
	regexp.MustCompile(`DELETED`),
	regexp.MustCompile(`zOLD`),
	regexp.MustCompile(`Archaic Spell`),
	regexp.MustCompile(`Well Repaired`),
	regexp.MustCompile(`Boss X`),
	regexp.MustCompile(`Adventurine`),
	regexp.MustCompile(`Sardonyx`),
	regexp.MustCompile(`Zyanite`),
	regexp.MustCompile(`zzold`),
	regexp.MustCompile(`Tom's`),
	regexp.MustCompile(`Stabilized Eternium Scope`),
}

// Allows manual overriding for Gem fields in case WowHead is wrong.
// Empty: this build ships no gem system -- GemProperties has zero rows.
var GemOverrides = []*proto.UIGem{}
var GemAllowList = map[int32]struct{}{}
var EnchantDenyListSpells = map[int32]struct{}{}
var EnchantDenyListItems = map[int32]struct{}{}
var GemDenyList = map[int32]struct{}{}

var EnchantDenyList = map[int32]struct{}{}

var EnchantAllowList = []int32{
	804,  // Enchant Cloak - Lesser Shadow Resistance
	684,  // Enchant Gloves - Major Strength
	963,  // Enchant Weapon - Major Striking
	1900, // Enchant Weapon - Crusader
	2564, // Weapon 15 Agi
	2583, // Presence of Might
	2588, // Presence of Sight
}

// Note: EffectId is required for all enchants, because they are
// used by various importers/exporters
// Note: ItemId, SpellId and Name are part of the enchant DB key, so they must
// match the generated enchant or the override is added as a separate entry.
var EnchantOverrides = []*proto.UIEnchant{
	// The head/leg resistance kits grant their stats through an equip spell that no
	// longer exists in the client data, so they have to be filled in by hand.
	{EffectId: 2681, ItemId: 22635, SpellId: 28162, Name: "Savage Guard", Stats: stats.Stats{stats.NatureResistance: 10}.ToProtoArray()},
	{EffectId: 2682, ItemId: 22636, SpellId: 28164, Name: "Ice Guard", Stats: stats.Stats{stats.FrostResistance: 10}.ToProtoArray()},
	{EffectId: 2683, ItemId: 22638, SpellId: 28166, Name: "Shadow Guard", Stats: stats.Stats{stats.ShadowResistance: 10}.ToProtoArray()},
}
