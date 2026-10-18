package forever

import (
	"github.com/wowsims/forever/sim/common/shared"
)

func RegisterAllProcs() {

	// Equip

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces your damage taken by 5% in the Battle Ring and The Maul.
	// https://www.wowhead.com/forever/spell=1317285
	// unsupported: the row applies only in area group 9337, which names no area type
	// equip: 1317285 (A_MOD_DAMAGE_PERCENT_TAKEN)
	// shared.NewSpellDataEquipAura(shared.SpellDataProc{TriggerSpellID: 1317285},
	//	[]shared.ItemVariant{
	//	{ItemID: 18706, ItemName: "Arena Master"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases your pet's armor by 10%.
	// https://www.wowhead.com/forever/spell=27225
	// unsupported: effect 1 A_MOD_BASE_RESISTANCE_PCT misc 1 is not parsed on a pet
	// equip: 27225 keeps 27208 up (A_MOD_BASE_RESISTANCE_PCT)
	// shared.NewSpellDataEquipAura(shared.SpellDataProc{TriggerSpellID: 27225},
	//	[]shared.ItemVariant{
	//	{ItemID: 22060, ItemName: "Beastmaster's Tunic"},
	//	{ItemID: 226886, ItemName: "Beastmaster's Tunic"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken and chance to be critically hit by 5% in Warsong Gulch.
	// https://www.wowhead.com/forever/spell=430432
	// unsupported: the row applies only in area group 6588, which names no area type; effect 2 A_MOD_ATTACKER_SPELL_AND_WEAPON_CRIT_CHANCE misc 0 is not parsed on the wearer
	// equip: 430432 (A_MOD_DAMAGE_PERCENT_TAKEN, A_MOD_ATTACKER_SPELL_AND_WEAPON_CRIT_CHANCE)
	// shared.NewSpellDataEquipAura(shared.SpellDataProc{TriggerSpellID: 430432},
	//	[]shared.ItemVariant{
	//	{ItemID: 211500, ItemName: "Resilient Cloth Headband"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken and chance to be critically hit by 5% in Warsong Gulch.
	// https://www.wowhead.com/forever/spell=430432
	// unsupported: the row applies only in area group 6588, which names no area type; effect 2 A_MOD_ATTACKER_SPELL_AND_WEAPON_CRIT_CHANCE misc 0 is not parsed on the wearer
	// equip: 430432 (A_MOD_DAMAGE_PERCENT_TAKEN, A_MOD_ATTACKER_SPELL_AND_WEAPON_CRIT_CHANCE)
	// shared.NewSpellDataEquipAura(shared.SpellDataProc{TriggerSpellID: 430432},
	//	[]shared.ItemVariant{
	//	{ItemID: 211856, ItemName: "Resilient Mail Coif"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken and chance to be critically hit by 5% in Warsong Gulch.
	// https://www.wowhead.com/forever/spell=430432
	// unsupported: the row applies only in area group 6588, which names no area type; effect 2 A_MOD_ATTACKER_SPELL_AND_WEAPON_CRIT_CHANCE misc 0 is not parsed on the wearer
	// equip: 430432 (A_MOD_DAMAGE_PERCENT_TAKEN, A_MOD_ATTACKER_SPELL_AND_WEAPON_CRIT_CHANCE)
	// shared.NewSpellDataEquipAura(shared.SpellDataProc{TriggerSpellID: 430432},
	//	[]shared.ItemVariant{
	//	{ItemID: 211857, ItemName: "Resilient Leather Mask"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases your Armor contribution from items by 4%.
	// https://www.wowhead.com/forever/spell=1270490
	// unsupported: effect 2 A_MOD_BONUS_ARMOR_PCT misc 0 is not parsed on the wearer
	// equip: 1270490 (A_MOD_BASE_RESISTANCE_PCT, A_MOD_BONUS_ARMOR_PCT)
	// shared.NewSpellDataEquipAura(shared.SpellDataProc{TriggerSpellID: 1270490},
	//	[]shared.ItemVariant{
	//	{ItemID: 263435, ItemName: "Mark of Urs'endris"},
	// })

	// After entering combat, reduce the next instance of Physical damage taken within 20s by 15.
	// https://www.wowhead.com/forever/spell=1291908
	// equip: 1291908 (A_MOD_DAMAGE_TAKEN)
	shared.NewSpellDataEquipAura(shared.SpellDataProc{TriggerSpellID: 1291908},
		[]shared.ItemVariant{
			{ItemID: 5443, ItemName: "Gold-plated Buckler"},
		})

	// Spell Damage received is reduced by 5.
	// https://www.wowhead.com/forever/spell=1292692
	// equip: 1292692 (A_MOD_DAMAGE_TAKEN)
	shared.NewSpellDataEquipAura(shared.SpellDataProc{TriggerSpellID: 1292692},
		[]shared.ItemVariant{
			{ItemID: 6907, ItemName: "Tortoise Armor"},
		})

	// Fire damage taken increased by 10%.
	// https://www.wowhead.com/forever/spell=1293003
	// equip: 1293003 (A_MOD_DAMAGE_PERCENT_TAKEN)
	shared.NewSpellDataEquipAura(shared.SpellDataProc{TriggerSpellID: 1293003},
		[]shared.ItemVariant{
			{ItemID: 9509, ItemName: "Petrolspill Leggings"},
		})

	// When struck in combat has a chance to reduce all damage taken by 51 for 10s.
	// https://www.wowhead.com/forever/spell=15595
	// equip: 15595 (A_MOD_DAMAGE_TAKEN)
	shared.NewSpellDataEquipAura(shared.SpellDataProc{TriggerSpellID: 15595},
		[]shared.ItemVariant{
			{ItemID: 11810, ItemName: "Force of Will"},
		})

	// After entering combat, reduce the next instance of Physical damage taken within 20s by 60.
	// https://www.wowhead.com/forever/spell=1298231
	// equip: 1298231 (A_MOD_DAMAGE_TAKEN)
	shared.NewSpellDataEquipAura(shared.SpellDataProc{TriggerSpellID: 1298231},
		[]shared.ItemVariant{
			{ItemID: 13955, ItemName: "Stoneform Shoulders"},
		})

	// Spell Damage received is reduced by 5.
	// https://www.wowhead.com/forever/spell=1292692
	// equip: 1292692 (A_MOD_DAMAGE_TAKEN)
	shared.NewSpellDataEquipAura(shared.SpellDataProc{TriggerSpellID: 1292692},
		[]shared.ItemVariant{
			{ItemID: 18351, ItemName: "Magically Sealed Bracers"},
		})

	// Increases damage dealt by your pet by 3%.
	// https://www.wowhead.com/forever/spell=27206
	// equip: 27206 keeps 27205 up (A_MOD_DAMAGE_PERCENT_DONE)
	shared.NewSpellDataEquipAura(shared.SpellDataProc{TriggerSpellID: 27206},
		[]shared.ItemVariant{
			{ItemID: 22061, ItemName: "Beastmaster's Boots"},
		})

	// Spell Damage received is reduced by 10.
	// https://www.wowhead.com/forever/spell=27518
	// equip: 27518 (A_MOD_DAMAGE_TAKEN)
	shared.NewSpellDataEquipAura(shared.SpellDataProc{TriggerSpellID: 27518},
		[]shared.ItemVariant{
			{ItemID: 22191, ItemName: "Obsidian Mail Tunic"},
		})

	// Increases damage dealt by your pet by 3%.
	// https://www.wowhead.com/forever/spell=27206
	// equip: 27206 keeps 27205 up (A_MOD_DAMAGE_PERCENT_DONE)
	shared.NewSpellDataEquipAura(shared.SpellDataProc{TriggerSpellID: 27206},
		[]shared.ItemVariant{
			{ItemID: 226881, ItemName: "Beastmaster's Treads"},
		})

	// Increases Spirit by 5%.
	// https://www.wowhead.com/forever/spell=1248751
	// equip: 1248751 (A_MOD_TOTAL_STAT_PERCENTAGE)
	shared.NewSpellDataEquipAura(shared.SpellDataProc{TriggerSpellID: 1248751},
		[]shared.ItemVariant{
			{ItemID: 249396, ItemName: "Mystic Mushroom"},
		})

	// Procs

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases Strength by 200 for 12s.
	// https://www.wowhead.com/forever/spell=17152
	// unsupported: states no rate
	// trigger 17152 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataProc(shared.SpellDataProc{TriggerSpellID: 17152, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 647, ItemName: "Destiny"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Smites an enemy for 30 Holy damage.
	// https://www.wowhead.com/forever/spell=13519
	// unsupported: states no rate
	// trigger 13519 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 13519, BuffSpellID: 13519, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 754, ItemName: "Shortsword of Vengeance"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target causing them to bleed for 155 damage over 15s.
	// https://www.wowhead.com/forever/spell=17504
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 809, ItemName: "Bloodrazor"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Launches a bolt of frost at the enemy causing 50 Frost damage and slowing movement speed by 50% for 5s.
	// https://www.wowhead.com/forever/spell=13439
	// unsupported: states no rate
	// trigger 13439 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 13439, BuffSpellID: 13439, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 810, ItemName: "Hammer of the Northern Wind"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 126 Nature damage.
	// https://www.wowhead.com/forever/spell=18104
	// unsupported: states no rate
	// trigger 18104 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18104, BuffSpellID: 18104, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 811, ItemName: "Axe of the Deep Woods"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Smite your target for 112 Holy damage. Damage is tripled against Undead.
	// https://www.wowhead.com/forever/spell=1325229
	// unsupported: states no rate
	// trigger 1325229 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1325229, BuffSpellID: 1325229, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 867, ItemName: "Gloves of Holy Might"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Decrease the armor of the target by 100 for 30s. While affected, the target cannot stealth or turn invisible.
	// https://www.wowhead.com/forever/spell=13752
	// unsupported: states no rate
	// trigger 13752 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 13752, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 869, ItemName: "Dazzling Longsword"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hurls a fiery ball that causes 262 Fire damage and an additional 60 damage over 6s.
	// https://www.wowhead.com/forever/spell=18796
	// unsupported: states no rate
	// trigger 18796 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18796, BuffSpellID: 18796, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 870, ItemName: "Fiery War Axe"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Damage caused by the target is reduced by 5 for 2min.
	// https://www.wowhead.com/forever/spell=8552
	// unsupported: states no rate
	// trigger 8552 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 8552, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 880, ItemName: "Staff of Horrors"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Poisons target for 3 Nature damage every 3.0 sec for 15s.
	// https://www.wowhead.com/forever/spell=18077
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 899, ItemName: "Venom Web Fang"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Lowers the Attack Power of the target by 40 for 1min.
	// https://www.wowhead.com/forever/spell=13524
	// unsupported: states no rate
	// trigger 13524 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 13524, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 934, ItemName: "Stalvan's Reaper"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 125 Shadow damage.
	// https://www.wowhead.com/forever/spell=18138
	// unsupported: states no rate
	// trigger 18138 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18138, BuffSpellID: 18138, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 937, ItemName: "Black Duskwood Staff"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck in combat has a 2% chance of dealing 115 Fire damage to all targets around you.
	// https://www.wowhead.com/forever/spell=18818
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 1168, ItemName: "Skullflame Shield -  - "},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck in combat, inflicts 3 Physical damage to the attacker.
	// https://www.wowhead.com/forever/spell=1325228
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 1204, ItemName: "The Green Tower - NatureResistance - "},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target's head, dealing 250 Physical damage, and increasing the casting time of all spells by
	// 30% for 30s.
	// https://www.wowhead.com/forever/spell=17148
	// unsupported: states no rate
	// trigger 17148 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 17148, BuffSpellID: 17148, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 1263, ItemName: "Brain Hacker"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Poisons target for 13 Nature damage every 5.0 sec for 25s.
	// https://www.wowhead.com/forever/spell=18208
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 1265, ItemName: "Scorpion Sting"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 84 Shadow damage.
	// https://www.wowhead.com/forever/spell=13480
	// unsupported: states no rate
	// trigger 13480 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 13480, BuffSpellID: 13480, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 1318, ItemName: "Night Reaver"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy, causing 100 Shadow damage.
	// https://www.wowhead.com/forever/spell=1312438
	// unsupported: states no rate
	// trigger 1312438 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1312438, BuffSpellID: 1312438, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 1387, ItemName: "Ghoulfang"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 42 Shadow damage.
	// https://www.wowhead.com/forever/spell=1292158
	// unsupported: states no rate
	// trigger 1292158 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1292158, BuffSpellID: 1292158, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 1481, ItemName: "Grimclaw"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 42 Shadow damage.
	// https://www.wowhead.com/forever/spell=1292158
	// unsupported: states no rate
	// trigger 1292158 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1292158, BuffSpellID: 1292158, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 1482, ItemName: "Shadowfang"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck in combat inflicts 2 Shadow damage to the attacker.
	// https://www.wowhead.com/forever/spell=1292160
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 1489, ItemName: "Gloomshroud Armor"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Poisons target for 30 Nature damage every 6.0 sec for 30s.
	// https://www.wowhead.com/forever/spell=16401
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 1726, ItemName: "Poison-tipped Bone Spear"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces target enemy's attack power by 84 for 30s.
	// https://www.wowhead.com/forever/spell=13528
	// unsupported: states no rate
	// trigger 13528 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 13528, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 1727, ItemName: "Sword of Decay"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 140 Fire damage.
	// https://www.wowhead.com/forever/spell=1300753
	// unsupported: states no rate
	// trigger 1300753 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1300753, BuffSpellID: 1300753, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 1728, ItemName: "Teebu's Blazing Longsword"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 231 Shadow damage.
	// https://www.wowhead.com/forever/spell=18211
	// unsupported: states no rate
	// trigger 18211 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18211, BuffSpellID: 18211, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 1982, ItemName: "Nightblade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target for 102 damage.
	// https://www.wowhead.com/forever/spell=18090
	// unsupported: states no rate
	// trigger 18090 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18090, BuffSpellID: 18090, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 1986, ItemName: "Gutrender"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 85 Arcane damage.
	// https://www.wowhead.com/forever/spell=18091
	// unsupported: states no rate
	// trigger 18091 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18091, BuffSpellID: 18091, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 2000, ItemName: "Archeus"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Movement speed increased by 2% in Elwynn Forest, Westfall, Redridge Mountains, and the Deadmines.
	// https://www.wowhead.com/forever/spell=1292011
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 2042, ItemName: "Staff of Westfall"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Fires a Flaming Cannonball at your target for 49 Fire damage.
	// https://www.wowhead.com/forever/spell=29639
	// unsupported: states no rate
	// trigger 29639 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 29639, BuffSpellID: 29639, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 2099, ItemName: "Dwarven Hand Cannon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 125 Shadow damage.
	// https://www.wowhead.com/forever/spell=18138
	// unsupported: states no rate
	// trigger 18138 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18138, BuffSpellID: 18138, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 2163, ItemName: "Shadowblade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Delivers a fatal wound for 105 Physical damage. Deals 100% increased damage to targets below 25% health
	// https://www.wowhead.com/forever/spell=18107
	// unsupported: states no rate
	// trigger 18107 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18107, BuffSpellID: 18107, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 2164, ItemName: "Gut Ripper"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 98 Shadow damage.
	// https://www.wowhead.com/forever/spell=18217
	// unsupported: states no rate
	// trigger 18217 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18217, BuffSpellID: 18217, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 2205, ItemName: "Duskbringer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Next spell cast within 4s will cast instantly.
	// https://www.wowhead.com/forever/spell=18803
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 2243, ItemName: "Hand of Edward the Odd"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 30 Shadow damage.
	// https://www.wowhead.com/forever/spell=13440
	// unsupported: states no rate
	// trigger 13440 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 13440, BuffSpellID: 13440, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 2256, ItemName: "Skeletal Club"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 35 Nature damage.
	// https://www.wowhead.com/forever/spell=14119
	// unsupported: states no rate
	// trigger 14119 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 14119, BuffSpellID: 14119, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 2263, ItemName: "Phytoblade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the duration of all movement impairing effects on you by 10%.
	// https://www.wowhead.com/forever/spell=1292546
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 2280, ItemName: "Kam's Walking Stick"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target causing them to bleed for 500 damage over 15s.
	// https://www.wowhead.com/forever/spell=17153
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 2291, ItemName: "Kang the Decapitator"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hurls a fiery ball that causes 98 Fire damage and an additional 18 damage over 6s.
	// https://www.wowhead.com/forever/spell=18199
	// unsupported: states no rate
	// trigger 18199 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18199, BuffSpellID: 18199, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 2299, ItemName: "Burning War Axe"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Fires a Searing Arrow at the target for 28 Fire damage.
	// https://www.wowhead.com/forever/spell=29638
	// unsupported: states no rate
	// trigger 29638 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 29638, BuffSpellID: 29638, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 2825, ItemName: "Bow of Searing Arrows"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 35 Shadow damage.
	// https://www.wowhead.com/forever/spell=16409
	// unsupported: states no rate
	// trigger 16409 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 16409, BuffSpellID: 16409, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 2912, ItemName: "Claw of the Shadowmancer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hurls a fiery ball that causes 300 Fire damage and an additional 100 damage over 8s.
	// https://www.wowhead.com/forever/spell=16415
	// unsupported: states no rate
	// trigger 16415 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 16415, BuffSpellID: 16415, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 2915, ItemName: "Taran Icebreaker"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Pummel the target for 25 damage and interrupt the spell being cast for 5s.
	// https://www.wowhead.com/forever/spell=13491
	// unsupported: states no rate
	// trigger 13491 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 13491, BuffSpellID: 13491, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 2942, ItemName: "Iron Knuckles"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target for 42 damage. Naga are Stunned for 2s.
	// https://www.wowhead.com/forever/spell=1292647
	// unsupported: states no rate
	// trigger 1292647 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1292647, BuffSpellID: 1292647, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 3078, ItemName: "Naga Heartpiercer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 70 Shadow damage.
	// https://www.wowhead.com/forever/spell=1292155
	// unsupported: states no rate
	// trigger 1292155 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1292155, BuffSpellID: 1292155, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 3194, ItemName: "Black Malice"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target causing them to bleed for 30 damage over 30s.
	// https://www.wowhead.com/forever/spell=18078
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 3336, ItemName: "Flesh Piercer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Dying inflicts 420 Holy damage to all enemies within 15 yds.
	// https://www.wowhead.com/forever/spell=1292746
	// unsupported: ProcTypeMask DEATH; no callback in the proc mask
	// trigger 1292749 (every time, core.CallbackEmpty, core.ProcMaskUnknown) -> buff 1292746
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1292749, BuffSpellID: 1292746},
	//	[]shared.ItemVariant{
	//	{ItemID: 3416, ItemName: "Martyr's Chain"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=1293501
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 3456, ItemName: "Dog Whistle"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deals 7 Fire damage to anyone who strikes you with a melee attack.
	// https://www.wowhead.com/forever/spell=21142
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 3475, ItemName: "Cloak of Flames"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces damage from falling.
	// https://www.wowhead.com/forever/spell=1292150
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 3748, ItemName: "Feline Mantle"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 35 Shadow damage.
	// https://www.wowhead.com/forever/spell=16409
	// unsupported: states no rate
	// trigger 16409 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 16409, BuffSpellID: 16409, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 3822, ItemName: "Runic Darkblade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Launches a bolt of frost at the enemy causing 50 Frost damage and slowing movement speed by 50% for 5s.
	// https://www.wowhead.com/forever/spell=13439
	// unsupported: states no rate
	// trigger 13439 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 13439, BuffSpellID: 13439, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 3854, ItemName: "Frost Tiger Blade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases your lockpicking skill slightly.
	// https://www.wowhead.com/forever/spell=9133
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 4248, ItemName: "Dark Leather Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Poisons target for 5 Nature damage every 3.0 sec for 15s.
	// https://www.wowhead.com/forever/spell=13518
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 4446, ItemName: "Blackvenom Blade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Poisons target for 6 Nature damage every 3.0 sec for 15s.
	// https://www.wowhead.com/forever/spell=16400
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 4449, ItemName: "Naraxis' Fang"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 50 Frost damage.
	// https://www.wowhead.com/forever/spell=18092
	// unsupported: states no rate
	// trigger 18092 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18092, BuffSpellID: 18092, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 5182, ItemName: "Shiver Blade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the chance your Pick Pocket and Distract abilities are successful by 1%
	// https://www.wowhead.com/forever/spell=1291987
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 5192, ItemName: "Thief's Blade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Tenderize the target, increasing physical damage taken by 1 for 1min.
	// https://www.wowhead.com/forever/spell=1292008
	// unsupported: states no rate
	// trigger 1292008 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataAuraProc(shared.SpellDataProc{TriggerSpellID: 1292008, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 5197, ItemName: "Cookie's Tenderizer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases your Cooking speed by 25%.
	// https://www.wowhead.com/forever/spell=1322312
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 5198, ItemName: "Cookie's Stirring Rod"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the crafting time of Smelt Copper, Smelt Tin, Smelt Bronze, and Smelt Silver by 25%.
	// https://www.wowhead.com/forever/spell=1291915
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 5199, ItemName: "Smelting Pants"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance on hit to blast the target for 28 Fire damage.
	// https://www.wowhead.com/forever/spell=1291662
	// unsupported: states no rate
	// trigger 1291686 (no stated rate, core.CallbackOnSpellHitDealt, core.ProcMaskRangedAuto) -> buff 1291662
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1291686, BuffSpellID: 1291662},
	//	[]shared.ItemVariant{
	//	{ItemID: 5243, ItemName: "Firebelcher"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Poisons target for 9 Nature damage every 3.0 sec for 15s.
	// https://www.wowhead.com/forever/spell=18197
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 5426, ItemName: "Serpent's Kiss"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target causing them to bleed for 77 damage over 7s.
	// https://www.wowhead.com/forever/spell=1294426
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 5616, ItemName: "Gutwrencher"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Poisons target for 6 Nature damage every 3.0 sec for 15s.
	// https://www.wowhead.com/forever/spell=16400
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 5752, ItemName: "Wyvern Tailspike"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 63 Frost damage.
	// https://www.wowhead.com/forever/spell=1293653
	// unsupported: states no rate
	// trigger 1293653 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1293653, BuffSpellID: 1293653, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 5756, ItemName: "Sliverblade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 75 Frost damage.
	// https://www.wowhead.com/forever/spell=20869
	// unsupported: states no rate
	// trigger 20869 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 20869, BuffSpellID: 20869, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 5815, ItemName: "Glacial Stone"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 42 Fire damage.
	// https://www.wowhead.com/forever/spell=1292152
	// unsupported: states no rate
	// trigger 1292152 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1292152, BuffSpellID: 1292152, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 6220, ItemName: "Meteor Shard"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases movement speed by 2% in Tirisfal Glades, Silverpine Forest, and Shadowfang Keep.
	// https://www.wowhead.com/forever/spell=1292121
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 6321, ItemName: "Silverlaine's Family Seal"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Cause party members within 30 yds to regenerate 4 health every 5 sec.
	// https://www.wowhead.com/forever/spell=1292142
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 6323, ItemName: "Baron's Scepter"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces target's attack power by 30 for 30s.
	// https://www.wowhead.com/forever/spell=13490
	// unsupported: states no rate
	// trigger 13490 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 13490, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 6331, ItemName: "Howling Blade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck in combat, inflicts 1 Nature damage to the attacker.
	// https://www.wowhead.com/forever/spell=1291698
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 6629, ItemName: "Sporid Cape"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Heals wielder of 78 damage over 12s.
	// https://www.wowhead.com/forever/spell=8348
	// unsupported: states no rate
	// trigger 8348 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataHealProc(shared.SpellDataProc{TriggerSpellID: 8348, BuffSpellID: 8348, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 6660, ItemName: "Julie's Dagger"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hurl spikes at the target, dealing 45 Physical damage and jumping to 1 additional nearby enemy. Each jump
	// reduces the damage by 50%.
	// https://www.wowhead.com/forever/spell=1293183
	// unsupported: states no rate
	// trigger 1293183 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1293183, BuffSpellID: 1293183, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 6681, ItemName: "Thornspike"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases movement speed by 3% in The Barrens, Thousand Needles, Razorfen Kraul, and Razorfen Downs.
	// https://www.wowhead.com/forever/spell=1293199
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 6693, ItemName: "Agamaggan's Clutch"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target causing them to bleed for 85 damage over 30s.
	// https://www.wowhead.com/forever/spell=16403
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 6738, ItemName: "Bleeding Crescent"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 77 Shadow damage.
	// https://www.wowhead.com/forever/spell=1293407
	// unsupported: states no rate
	// trigger 1293407 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1293407, BuffSpellID: 1293407, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 6831, ItemName: "Black Menace"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Poisons target for 14 Nature damage every 3.0 sec for 9s.
	// https://www.wowhead.com/forever/spell=8313
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 6904, ItemName: "Bite of Serra'kis"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Inflict 42 Fire and 14 Frost damage, slow the target by 50% for 5s, reduce armor by 50, and inflict 7
	// Nature damage every 3.0 sec for 12s.
	// https://www.wowhead.com/forever/spell=1292657
	// unsupported: states no rate
	// trigger 1292657 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 1292657, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 6909, ItemName: "Strike of the Hydra"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Ignite the target in the flames of the Great Forge, dealing 49 Holyfire damage. Deals 3 times as much
	// damage to Fire Elementals and Earth Elementals.
	// https://www.wowhead.com/forever/spell=1322218
	// unsupported: states no rate
	// trigger 1322218 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1322218, BuffSpellID: 1322218, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 6953, ItemName: "Verigan's Fist"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases swim speed by 15%.
	// https://www.wowhead.com/forever/spell=8747
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 7052, ItemName: "Azure Silk Belt"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck by a melee attacker, that attacker has a 5% chance of being put to sleep for 10s. Only affects
	// enemies level 50 and below.
	// https://www.wowhead.com/forever/spell=9159
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 7375, ItemName: "Green Whelp Armor"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Melee attacks with this weapon torture the target, dealing 1 Fire damage and increasing your damage dealt
	// to that target with this effect by 1, stacking up to 10 times.
	// https://www.wowhead.com/forever/spell=1293433
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 7682, ItemName: "Torturing Poker"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Discipline the target, dealing 91 Physical damage. Canines and Hyenas are Stunned and sit for 3s.
	// https://www.wowhead.com/forever/spell=1293482
	// unsupported: states no rate
	// trigger 1293482 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1293482, BuffSpellID: 1293482, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 7710, ItemName: "Loksey's Training Stick"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Attack all nearby enemies, dealing weapon damage plus 5 every 3.0 sec for 9s, and reducing your own movement
	// speed by 80%.
	// https://www.wowhead.com/forever/spell=9632
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 7717, ItemName: "Ravager"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the cast time of spells and effects that revive targets by 20%.
	// https://www.wowhead.com/forever/spell=1293540
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 7720, ItemName: "Whitemane's Chapeau"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deals 14 Holy damage every time you block.
	// https://www.wowhead.com/forever/spell=1293539
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 7726, ItemName: "Aegis of the Scarlet Commander"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 133 Frost damage.
	// https://www.wowhead.com/forever/spell=18204
	// unsupported: states no rate
	// trigger 18204 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18204, BuffSpellID: 18204, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 7730, ItemName: "Cobalt Crusher"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Silences an enemy, preventing it from casting spells for 6s.
	// https://www.wowhead.com/forever/spell=1293654
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 7736, ItemName: "Fight Club"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target causing them to bleed for 250 damage over 10s.
	// https://www.wowhead.com/forever/spell=18200
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 7753, ItemName: "Bloodspiller"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Disarm target's weapon for 10s.
	// https://www.wowhead.com/forever/spell=13534
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 7954, ItemName: "The Shatterer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Diseases a target for 100 Nature damage and an additional 360 damage over 1min.
	// https://www.wowhead.com/forever/spell=9796
	// unsupported: states no rate
	// trigger 9796 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 9796, BuffSpellID: 9796, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 7959, ItemName: "Blight"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Protects the caster with a holy shield.
	// https://www.wowhead.com/forever/spell=9800
	// unsupported: states no rate
	// trigger 9800 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataAbsorbProc(shared.SpellDataProc{TriggerSpellID: 9800, BuffSpellID: 9800, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 7960, ItemName: "Truesilver Champion"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Decrease the armor of the target by 100 for 20s. While affected, the target cannot stealth or turn invisible.
	// https://www.wowhead.com/forever/spell=9806
	// unsupported: states no rate
	// trigger 9806 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 9806, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 7961, ItemName: "Phantom Blade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 15 Nature damage.
	// https://www.wowhead.com/forever/spell=13482
	// unsupported: states no rate
	// trigger 13482 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 13482, BuffSpellID: 13482, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 8006, ItemName: "The Ziggler"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target for 75 damage.
	// https://www.wowhead.com/forever/spell=16405
	// unsupported: states no rate
	// trigger 16405 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 16405, BuffSpellID: 16405, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 8190, ItemName: "Hanzo Sword"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Defense +48 for 5s.
	// https://www.wowhead.com/forever/spell=10351
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 8223, ItemName: "Blade of the Basilisk"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target causing them to bleed for 85 damage over 30s.
	// https://www.wowhead.com/forever/spell=16403
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 8224, ItemName: "Silithid Ripper"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Corrupts the target, causing 63 damage over 3s.
	// https://www.wowhead.com/forever/spell=13530
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 8225, ItemName: "Tainted Pierce"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deal 112 Fire damage. Deals 3 times as much damage to Earth Elementals, Fire Elementals, and Mountain
	// Giants.
	// https://www.wowhead.com/forever/spell=1317432
	// unsupported: states no rate
	// trigger 1317432 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1317432, BuffSpellID: 1317432, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 8708, ItemName: "Hammer of Expertise"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Lash the enemy with the rage of Sul'thraze, reducing their attack power by 180 for 15s and dealing 159
	// Shadow damage. After 5s, the target suffers an additional 159 Shadow damage. This additional Shadow damage
	// deals 2 times as much damage to Trolls.
	// https://www.wowhead.com/forever/spell=1294430
	// unsupported: states no rate
	// trigger 1294430 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1294430, BuffSpellID: 1294430, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 9372, ItemName: "Sul'thraze the Lasher"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases your chance to Parry attacks by 2% for 9s.
	// https://www.wowhead.com/forever/spell=1294243
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 9379, ItemName: "Sang'thraze the Deflector"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Shields the wielder from physical damage, absorbing 84 damage. Lasts 20s.
	// https://www.wowhead.com/forever/spell=11657
	// unsupported: states no rate
	// trigger 11657 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataAbsorbProc(shared.SpellDataProc{TriggerSpellID: 11657, BuffSpellID: 11657, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 9380, ItemName: "Jang'thraze the Protector"},
	//	{ItemID: 11086, ItemName: "Jang'thraze the Protector"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hurls a fiery ball that causes 51 Fire damage and an additional 12 damage over 6s.
	// https://www.wowhead.com/forever/spell=13438
	// unsupported: states no rate
	// trigger 13438 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 13438, BuffSpellID: 13438, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 9386, ItemName: "Excavator's Brand"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Strike your ranged target with a Fire Blast for 28 Fire damage.
	// https://www.wowhead.com/forever/spell=29644
	// unsupported: states no rate
	// trigger 29644 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 29644, BuffSpellID: 29644, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 9412, ItemName: "Galgann's Fireblaster"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases Feral Combat +1.
	// https://www.wowhead.com/forever/spell=1320726
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 9414, ItemName: "Oilskin Leggings"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Rend earth and stone, dealing 196 Physical damage. Deals 3 times as much damage to Earth Elementals, Golems,
	// and Titan-forged.
	// https://www.wowhead.com/forever/spell=12731
	// unsupported: states no rate
	// trigger 12731 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 12731, BuffSpellID: 12731, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 9418, ItemName: "Stoneslayer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance on harmful spell cast to blast a target for 77 Fire damage.
	// https://www.wowhead.com/forever/spell=18083
	// unsupported: states no rate
	// trigger 1293980 (no stated rate, core.CallbackOnCastComplete, core.ProcMaskSpellDamage) -> buff 18083
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1293980, BuffSpellID: 18083},
	//	[]shared.ItemVariant{
	//	{ItemID: 9419, ItemName: "Galgann's Firehammer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases your attack speed by 15% for 10s.
	// https://www.wowhead.com/forever/spell=13533
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 9423, ItemName: "The Jackhammer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Delivers a fatal wound for 311 Physical damage. Deals 50% increased damage to targets below 25% health.
	// https://www.wowhead.com/forever/spell=10373
	// unsupported: states no rate
	// trigger 10373 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 10373, BuffSpellID: 10373, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 9425, ItemName: "Pendulum of Doom -  - "},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Each swing alternates your attack speed between 100% faster and 50% slower.
	// https://www.wowhead.com/forever/spell=1293881
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 9425, ItemName: "Pendulum of Doom -  - "},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 28 Nature damage. Mechanical targets take 2 x damage.
	// https://www.wowhead.com/forever/spell=1292851
	// unsupported: states no rate
	// trigger 1292851 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1292851, BuffSpellID: 1292851, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 9446, ItemName: "Electrocutioner Leg"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Underwater Breath lasts 50% longer than normal.
	// https://www.wowhead.com/forever/spell=11789
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 9452, ItemName: "Hydrocane"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deals 14 Nature damage every 3.0 sec to any enemy in an 8 yard radius around the caster for 12s.
	// https://www.wowhead.com/forever/spell=11790
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 9453, ItemName: "Toxic Revenger"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Punctures target's armor, lowering it by 100.
	// https://www.wowhead.com/forever/spell=1293005
	// unsupported: states no rate
	// trigger 1293005 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 1293005, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 9465, ItemName: "Digmaster 5000"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deal 84 Nature damage, increased by 2% per stack of Static Electricity. Deals 2 times as much damage to
	// Aquatic enemies.
	// https://www.wowhead.com/forever/spell=3742
	// unsupported: states no rate
	// trigger 3742 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 3742, BuffSpellID: 3742, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 9467, ItemName: "Gahz'rilla Fang -  - "},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Movement can generate Static Electricity.
	// https://www.wowhead.com/forever/spell=1294477
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 9467, ItemName: "Gahz'rilla Fang -  - "},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Delivers a fatal wound for 174 Physical damage. Deals 50% increased damage to targets below 25% health.
	// https://www.wowhead.com/forever/spell=18206
	// unsupported: states no rate
	// trigger 18206 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18206, BuffSpellID: 18206, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 9475, ItemName: "Diabolic Skiver"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deal 294 Physical damage and stun for 1s. Trolls take 2 times as much damage and stun duration.
	// https://www.wowhead.com/forever/spell=1294339
	// unsupported: states no rate
	// trigger 1294339 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1294339, BuffSpellID: 1294339, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 9477, ItemName: "The Chief's Enforcer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target for 70 Physical damage.
	// https://www.wowhead.com/forever/spell=1292573
	// unsupported: states no rate
	// trigger 1292573 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1292573, BuffSpellID: 1292573, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 9478, ItemName: "Ripsaw"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance on hit to light the target ablaze, dealing 12 Fire damage over 3s.
	// https://www.wowhead.com/forever/spell=1294537
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 9483, ItemName: "Flaming Incinerator"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Punctures target's armor, lowering it by 100.
	// https://www.wowhead.com/forever/spell=1293005
	// unsupported: states no rate
	// trigger 1293005 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 1293005, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 9485, ItemName: "Vibroblade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 105 Nature damage.
	// https://www.wowhead.com/forever/spell=13527
	// unsupported: states no rate
	// trigger 13527 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 13527, BuffSpellID: 13527, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 9486, ItemName: "Supercharger Battle Axe"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target for 77 Physical damage.
	// https://www.wowhead.com/forever/spell=13486
	// unsupported: states no rate
	// trigger 13486 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 13486, BuffSpellID: 13486, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 9511, ItemName: "Bloodletter Scalpel"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Disarm target's weapon for 5s.
	// https://www.wowhead.com/forever/spell=11879
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 9608, ItemName: "Shoni's Disarming Tool"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts nearby enemies with thunder increasing the time between their attacks by 11% for 10s and doing
	// 42 Nature damage to them. Will affect up to 4 targets.
	// https://www.wowhead.com/forever/spell=13532
	// unsupported: states no rate
	// trigger 13532 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 13532, BuffSpellID: 13532, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 9639, ItemName: "The Hand of Antu'sul"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 91 Nature damage.
	// https://www.wowhead.com/forever/spell=18081
	// unsupported: states no rate
	// trigger 18081 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18081, BuffSpellID: 18081, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 9651, ItemName: "Gryphon Rider's Stormhammer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Water breathing lasts 100% longer than normal.
	// https://www.wowhead.com/forever/spell=1320643
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 10506, ItemName: "Deepdive Helmet"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Strike your ranged target with a Quill Shot for 33 Nature damage.
	// https://www.wowhead.com/forever/spell=29646
	// unsupported: states no rate
	// trigger 29646 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 29646, BuffSpellID: 29646, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 10567, ItemName: "Quillshooter"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Launches a bolt of frost at the enemy causing 50 Frost damage and slowing movement speed by 50% for 5s.
	// https://www.wowhead.com/forever/spell=13439
	// unsupported: states no rate
	// trigger 13439 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 13439, BuffSpellID: 13439, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 10623, ItemName: "Winter's Bite"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases Attack Power by 102 and melee attack speed by 6% for 15s.
	// https://www.wowhead.com/forever/spell=12686
	// unsupported: states no rate
	// trigger 12686 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataProc(shared.SpellDataProc{TriggerSpellID: 12686, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 10626, ItemName: "Ragehammer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Bludgeon the target, dealing 126 Physical damage. Canines and Hyaenidae are stunned and sit for 3s.
	// https://www.wowhead.com/forever/spell=1300209
	// unsupported: states no rate
	// trigger 1300209 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1300209, BuffSpellID: 1300209, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 10627, ItemName: "Bludgeon of the Grinning Dog"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Delivers a fatal wound for 159 Physical damage. Deals 75% increased damage to targets below 25% health.
	// https://www.wowhead.com/forever/spell=1300185
	// unsupported: states no rate
	// trigger 1300185 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1300185, BuffSpellID: 1300185, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 10628, ItemName: "Deathblow"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Heals wielder of 42 damage.
	// https://www.wowhead.com/forever/spell=1300356
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 10750, ItemName: "Lifeforce Dirk"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Unleash a mighty swing, dealing 42 Holy damage. Deals 3 times as much damage to Undead and Swine.
	// https://www.wowhead.com/forever/spell=1293787
	// unsupported: states no rate
	// trigger 1293787 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1293787, BuffSpellID: 1293787, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 10758, ItemName: "X'caliboar"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target causing them to bleed for 70 damage over 7s.
	// https://www.wowhead.com/forever/spell=18075
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 10772, ItemName: "Glutton's Cleaver"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the duration of all movement impairing effects on you by 7%.
	// https://www.wowhead.com/forever/spell=1293819
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 10776, ItemName: "Silky Spider Cape"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hurls a fiery ball that causes 77 Fire damage and an additional 14 damage over 7s.
	// https://www.wowhead.com/forever/spell=16413
	// unsupported: states no rate
	// trigger 16413 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 16413, BuffSpellID: 16413, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 10797, ItemName: "Firebreather"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Corrupts the target, causing 90 damage over 3s.
	// https://www.wowhead.com/forever/spell=18088
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 10803, ItemName: "Blade of the Wretched"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Steals 56 life from target enemy.
	// https://www.wowhead.com/forever/spell=18084
	// unsupported: states no rate
	// trigger 18084 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18084, BuffSpellID: 18084, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 10804, ItemName: "Fist of the Damned"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Steals 21 life from the target enemy. Deals 3 times as much damage to Undead.
	// https://www.wowhead.com/forever/spell=1299869
	// unsupported: states no rate
	// trigger 1299869 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1299869, BuffSpellID: 1299869, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 10805, ItemName: "Eater of the Dead"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Harmful spell casts have a chance to steal 112 life from the target enemy over 8s.
	// https://www.wowhead.com/forever/spell=1299943
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 10844, ItemName: "Spire of Hakkar"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 28 Shadow damage.
	// https://www.wowhead.com/forever/spell=16408
	// unsupported: states no rate
	// trigger 16408 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 16408, BuffSpellID: 16408, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 11121, ItemName: "Darkwater Talwar"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target for 225 damage.
	// https://www.wowhead.com/forever/spell=1320808
	// unsupported: states no rate
	// trigger 1320808 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1320808, BuffSpellID: 1320808, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 11603, ItemName: "Vilerend Slicer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces targets armor by 300 for 20s.
	// https://www.wowhead.com/forever/spell=15280
	// unsupported: states no rate
	// trigger 15280 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 15280, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 11607, ItemName: "Dark Iron Sunderer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Stuns target for 8s.
	// https://www.wowhead.com/forever/spell=15283
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 11608, ItemName: "Dark Iron Pulverizer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Corrosive acid that deals 13 Nature damage every 2.0 sec and lowers target's armor by 50 for 14s.
	// https://www.wowhead.com/forever/spell=13526
	// unsupported: states no rate
	// trigger 13526 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 13526, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 11635, ItemName: "Hookfang Shanker"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Harmful spell casts have a chance to deal 70 Fire damage. Deals 3 times as much damage to Plants.
	// https://www.wowhead.com/forever/spell=1300782
	// unsupported: states no rate
	// trigger 1300781 (no stated rate, core.CallbackOnCastComplete, core.ProcMaskSpellDamage) -> buff 1300782
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1300781, BuffSpellID: 1300782},
	//	[]shared.ItemVariant{
	//	{ItemID: 11750, ItemName: "Kindling Stave"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 182 Fire damage.
	// https://www.wowhead.com/forever/spell=18086
	// unsupported: states no rate
	// trigger 18086 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18086, BuffSpellID: 18086, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 11803, ItemName: "Force of Magma"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deal 109 Physical damage. Trolls and Orcs take 2 times as much damage and are horrified for 4s.
	// https://www.wowhead.com/forever/spell=1300783
	// unsupported: states no rate
	// trigger 1300783 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1300783, BuffSpellID: 1300783, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 11816, ItemName: "Angerforge's Battle Axe"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 84 for 10s.
	// https://www.wowhead.com/forever/spell=15602
	// unsupported: states no rate
	// trigger 15602 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataProc(shared.SpellDataProc{TriggerSpellID: 15602, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 11817, ItemName: "Lord General's Sword"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 98 Nature damage.
	// https://www.wowhead.com/forever/spell=18089
	// unsupported: states no rate
	// trigger 18089 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18089, BuffSpellID: 18089, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 11902, ItemName: "Linken's Sword of Mastery"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Steals 49 life from target enemy. Deals 2 times as much damage to Undead.
	// https://www.wowhead.com/forever/spell=16414
	// unsupported: states no rate
	// trigger 16414 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 16414, BuffSpellID: 16414, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 11920, ItemName: "Wraith Scythe"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Dying heals all party members within 40 yds for 1176.
	// https://www.wowhead.com/forever/spell=1300891
	// unsupported: ProcTypeMask DEATH; no callback in the proc mask; the heal lands on implicit target 22, not the wearer
	// trigger 1300890 (every time, core.CallbackEmpty, core.ProcMaskUnknown) -> buff 1300891
	// shared.NewSpellDataHealProc(shared.SpellDataProc{TriggerSpellID: 1300890, BuffSpellID: 1300891},
	//	[]shared.ItemVariant{
	//	{ItemID: 11923, ItemName: "The Hammer of Grace"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hurls a fiery ball that causes 171 Fire damage and an additional 18 damage over 6s.
	// https://www.wowhead.com/forever/spell=15662
	// unsupported: states no rate
	// trigger 15662 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 15662, BuffSpellID: 15662, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 12243, ItemName: "Smoldering Claw"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 30 Shadow damage.
	// https://www.wowhead.com/forever/spell=13440
	// unsupported: states no rate
	// trigger 13440 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 13440, BuffSpellID: 13440, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 12250, ItemName: "Midnight Axe"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=16372
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12344, ItemName: "Seal of Ascension"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target causing them to bleed for 161 damage over 7s. Deals 2 times as much damage to Dragonkin.
	// https://www.wowhead.com/forever/spell=14118
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12463, ItemName: "Drakefang Butcher"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Leave lacerations on the target, increasing physical damage taken by 1 for 1min. Stacks up to 5 times.
	// https://www.wowhead.com/forever/spell=10370
	// unsupported: states no rate
	// trigger 10370 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataAuraProc(shared.SpellDataProc{TriggerSpellID: 10370, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 12469, ItemName: "Mutilator"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Stuns target for 3s.
	// https://www.wowhead.com/forever/spell=56
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12528, ItemName: "The Judge's Gavel"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 56 Fire damage and increases Fire damage done to target by 10 for 30s.
	// https://www.wowhead.com/forever/spell=16454
	// unsupported: states no rate
	// trigger 16454 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 16454, BuffSpellID: 16454, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 12531, ItemName: "Searing Needle"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Inflicts numbing pain that deals 10 Nature damage every 1.0 sec and increases time between target's attacks
	// by 11% for 10s.
	// https://www.wowhead.com/forever/spell=16528
	// unsupported: states no rate
	// trigger 16528 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 16528, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 12582, ItemName: "Keris of Zul'Serak"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Delivers a fatal wound for 448 Physical damage. Deals 50% increased damage to targets below 25% health.
	// https://www.wowhead.com/forever/spell=16549
	// unsupported: states no rate
	// trigger 16549 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 16549, BuffSpellID: 16549, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 12583, ItemName: "Blackhand Doomsaw"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deals 75 damage when you are the victim of a critical melee strike.
	// https://www.wowhead.com/forever/spell=16550
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12588, ItemName: "Bonespike Shoulder"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// All attacks are guaranteed to land and will be critical strikes for the next 3s.
	// https://www.wowhead.com/forever/spell=16551
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12590, ItemName: "Felstriker"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Transfers 10 health every 1.0 seconds from the target to the caster for 5s.
	// https://www.wowhead.com/forever/spell=16603
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12621, ItemName: "Demonfork"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck has a 3% chance of stealing 240 life from the attacker over 4s.
	//
	// https://www.wowhead.com/forever/spell=16608
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12628, ItemName: "Demon Forged Breastplate"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck in combat has a {UNK: H}% chance to make you invulnerable to melee damage for 3s. This effect
	// can only occur once every 10 sec.
	// https://www.wowhead.com/forever/spell=16621
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12641, ItemName: "Invulnerable Mail"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Enemy is inflicted with the Bleakwood Curse that reduces their magic resistances by 25. Can be applied
	// up to 3 times.
	// https://www.wowhead.com/forever/spell=16871
	// unsupported: states no rate
	// trigger 16871 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 16871, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 12769, ItemName: "Bleakwood Hew"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Burns the enemy for 200 damage over 30s.
	// https://www.wowhead.com/forever/spell=16898
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12777, ItemName: "Blazing Rapier"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Dispels a magic effect on the current foe.
	// https://www.wowhead.com/forever/spell=16908
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12781, ItemName: "Serenity"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Heal yourself for 360 and increase your Strength by 120 for 30s.
	// https://www.wowhead.com/forever/spell=16916
	// unsupported: states no rate
	// trigger 16916 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataProc(shared.SpellDataProc{TriggerSpellID: 16916, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 12790, ItemName: "Arcanite Champion"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target causing them to bleed for 98 damage over 7s. Deals 2 times as much damage if applied
	// while behind the target.
	// https://www.wowhead.com/forever/spell=1300787
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12791, ItemName: "Barman Shanker"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hurls a fiery ball that causes 228 Fire damage and an additional 36 damage over 6s.
	// https://www.wowhead.com/forever/spell=18082
	// unsupported: states no rate
	// trigger 18082 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18082, BuffSpellID: 18082, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 12792, ItemName: "Volcanic Hammer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target causing them to bleed for 100 damage over 30s.
	// https://www.wowhead.com/forever/spell=13318
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12795, ItemName: "Blood Talon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Stuns target for 3s.
	// https://www.wowhead.com/forever/spell=56
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12796, ItemName: "Hammer of the Titans"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Slows the target's movement speed by 30% and increases the time between their attacks by 25% for 5s.
	// https://www.wowhead.com/forever/spell=16927
	// unsupported: states no rate
	// trigger 16927 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 16927, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 12797, ItemName: "Frostguard"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance on hit to grant nearby party members 4% increased critical strike chance.
	// https://www.wowhead.com/forever/spell=16939
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12802, ItemName: "Darkspear"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Unable to be tracked.
	// https://www.wowhead.com/forever/spell=1301707
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12930, ItemName: "Briarwood Reed"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deals 51 Nature damage every 2.0 sec to all enemies within an 8 yard radius of the caster for 10s.
	// https://www.wowhead.com/forever/spell=17196
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12969, ItemName: "Seeping Willow"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 40 Shadow damage.
	// https://www.wowhead.com/forever/spell=14106
	// unsupported: states no rate
	// trigger 14106 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 14106, BuffSpellID: 14106, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 12974, ItemName: "The Black Knight"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hurls a fiery ball that causes 80 Fire damage and an additional 18 damage over 6s.
	// https://www.wowhead.com/forever/spell=1299850
	// unsupported: states no rate
	// trigger 1299850 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1299850, BuffSpellID: 1299850, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 12992, ItemName: "Searing Blade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target causing them to bleed for 100 damage over 30s.
	// https://www.wowhead.com/forever/spell=13318
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13016, ItemName: "Killmaim"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Corrupts the target, causing 30 damage over 3s.
	// https://www.wowhead.com/forever/spell=17510
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13032, ItemName: "Sword of Corruption"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Poisons target for 8 Nature damage every 2.0 sec for 20s.
	// https://www.wowhead.com/forever/spell=17511
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13035, ItemName: "Serpent Slicer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance to strike your ranged target with a Shadowbolt for 16 Shadow damage.
	// https://www.wowhead.com/forever/spell=29640
	// unsupported: states no rate
	// trigger 29626 (no stated rate, core.CallbackOnSpellHitDealt, core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial) -> buff 29640
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 29626, BuffSpellID: 29640},
	//	[]shared.ItemVariant{
	//	{ItemID: 13040, ItemName: "Heartseeking Crossbow"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 150 Shadow damage.
	// https://www.wowhead.com/forever/spell=18214
	// unsupported: states no rate
	// trigger 18214 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18214, BuffSpellID: 18214, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 13051, ItemName: "Witchfury"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 231 Shadow damage.
	// https://www.wowhead.com/forever/spell=18211
	// unsupported: states no rate
	// trigger 18211 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18211, BuffSpellID: 18211, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 13053, ItemName: "Doombringer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target for 130 damage.
	// https://www.wowhead.com/forever/spell=14126
	// unsupported: states no rate
	// trigger 14126 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 14126, BuffSpellID: 14126, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 13054, ItemName: "Grim Reaper"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target causing them to bleed for 98 damage over 14s.
	// https://www.wowhead.com/forever/spell=18202
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13057, ItemName: "Bloodpike"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target for 75 damage.
	// https://www.wowhead.com/forever/spell=16405
	// unsupported: states no rate
	// trigger 16405 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 16405, BuffSpellID: 16405, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 13060, ItemName: "The Needler"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Strike the target with a Flaming Shell for 38 Fire damage.
	// https://www.wowhead.com/forever/spell=29647
	// unsupported: states no rate
	// trigger 29647 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 29647, BuffSpellID: 29647, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 13146, ItemName: "Shell Launcher Shotgun"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 205 Frost damage.
	// https://www.wowhead.com/forever/spell=19260
	// unsupported: states no rate
	// trigger 19260 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 19260, BuffSpellID: 19260, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 13148, ItemName: "Chillpike"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Poisons target for 7 Nature damage every 1.0 sec for 15s.
	// https://www.wowhead.com/forever/spell=18203
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13183, ItemName: "Venomspitter"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Knocks target silly for 1s. Orcs and Ogres are stunned for 2 times as long.
	// https://www.wowhead.com/forever/spell=17308
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13198, ItemName: "Hurd Smasher"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Slows target enemy's casting speed and increases the time between melee and ranged attacks by 9% for 10s.
	// https://www.wowhead.com/forever/spell=17331
	// unsupported: states no rate
	// trigger 17331 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 17331, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 13218, ItemName: "Fang of the Crystal Spider"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Has a chance when struck in combat to increase your chance to block by 18% for 10s. Undead that strike
	// you are 2 times as likely to activate this effect.
	// https://www.wowhead.com/forever/spell=17351
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13243, ItemName: "Argent Defender"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 700 Fire damage.
	// https://www.wowhead.com/forever/spell=18112
	// unsupported: states no rate
	// trigger 18112 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18112, BuffSpellID: 18112, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 13262, ItemName: "Ashbringer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target for 100 damage and deals an additional 13 damage every 1.0 sec for 10s.
	// https://www.wowhead.com/forever/spell=17407
	// unsupported: states no rate
	// trigger 17407 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 17407, BuffSpellID: 17407, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 13285, ItemName: "The Blackrock Slicer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 150 Shadow damage and dealing 34 damage every 2.0 sec for 6s.
	// https://www.wowhead.com/forever/spell=17483
	// unsupported: states no rate
	// trigger 17483 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 17483, BuffSpellID: 17483, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 13348, ItemName: "Demonshear"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=1298517
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13353, ItemName: "Book of the Dead"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deals 22 Physical damage every time you block.
	// https://www.wowhead.com/forever/spell=17496
	// unsupported: an outcome the proc mask has no bit for
	// trigger 17495 (every time, core.CallbackOnSpellHitTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial) -> buff 17496
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 17495, BuffSpellID: 17496},
	//	[]shared.ItemVariant{
	//	{ItemID: 13375, ItemName: "Crest of Retribution"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=1300668
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13382, ItemName: "Cannonball Runner"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Knocks target silly for 2s and increases Strength by 53 for 20s.
	// https://www.wowhead.com/forever/spell=17500
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13393, ItemName: "Malown's Slam"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target, causing them to bleed for 140 damage over 14s.
	// https://www.wowhead.com/forever/spell=1320824
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13399, ItemName: "Gargoyle Shredder Talons"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Target enemy loses 10 health and mana every 1.0 sec for 10s.
	// https://www.wowhead.com/forever/spell=17506
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13408, ItemName: "Soul Breaker"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the effects that healing and mana potions have on the wearer by 20%.
	// https://www.wowhead.com/forever/spell=17619
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13503, ItemName: "Alchemist's Stone"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Silences an enemy, preventing it from casting spells for 6s.
	// https://www.wowhead.com/forever/spell=1293654
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13953, ItemName: "Silent Fang"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=1318945
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13965, ItemName: "Blackhand's Breadth"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=1318846
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13968, ItemName: "Eye of the Beast"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Diseases target enemy for 42 Nature damage every 3.0 sec for 12s.
	// https://www.wowhead.com/forever/spell=18289
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13983, ItemName: "Gravestone War Axe"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=1298508
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 14022, ItemName: "Barov Peasant Caller"},
	//	{ItemID: 14023, ItemName: "Barov Peasant Caller"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increase the duration of Fear effects applied to the target by 25% for 1min.
	// https://www.wowhead.com/forever/spell=19755
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 14024, ItemName: "Frightalon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces target enemy's attack power by 24 for 30s.
	// https://www.wowhead.com/forever/spell=18381
	// unsupported: states no rate
	// trigger 18381 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 18381, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 14145, ItemName: "Cursed Felblade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the cooldown of your Fade ability by -2.0 sec.
	// https://www.wowhead.com/forever/spell=18388
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 14154, ItemName: "Truefaith Vestments"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deals 15 Shadow damage every 2.0 sec for 14s and lowers their attack power by 92 for 14s.
	// https://www.wowhead.com/forever/spell=18633
	// unsupported: states no rate
	// trigger 18633 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 18633, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 14531, ItemName: "Frightskull Shaft"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deals 49 Shadow damage every 1.0 sec for 5s. All damage done is then transferred to the caster.
	// https://www.wowhead.com/forever/spell=18652
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 14541, ItemName: "Barovian Family Sword"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Protects the wearer from being fully engulfed by Shadow Flame.
	// https://www.wowhead.com/forever/spell=22683
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 15138, ItemName: "Onyxia Scale Cloak"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 252 Nature damage.
	// https://www.wowhead.com/forever/spell=19874
	// unsupported: states no rate
	// trigger 19874 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 19874, BuffSpellID: 19874, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 15418, ItemName: "Shimmering Platinum Warhammer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target causing them to bleed for 80 damage over 30s.
	// https://www.wowhead.com/forever/spell=16406
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 15814, ItemName: "Hameya's Slayer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Inflicts 21 Nature damage every 2.0 sec for 10s.
	// https://www.wowhead.com/forever/spell=20586
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 15853, ItemName: "Windreaper"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Fires a Shadow Shot at the target for 26 Shadow damage.
	// https://www.wowhead.com/forever/spell=29641
	// unsupported: states no rate
	// trigger 29641 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 29641, BuffSpellID: 29641, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 16004, ItemName: "Dark Iron Rifle"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=1318325
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16022, ItemName: "Arcanite Dragonling"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage absorbed by your Mana Shield by 285.
	// https://www.wowhead.com/forever/spell=23037
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16391, ItemName: "Knight-Lieutenant's Silk Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Sprint ability by 3.0 sec.
	// https://www.wowhead.com/forever/spell=23049
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16392, ItemName: "Knight-Lieutenant's Leather Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the Holy damage bonus of your Judgement of the Crusader by 20.
	// https://www.wowhead.com/forever/spell=23300
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16410, ItemName: "Knight-Lieutenant's Lamellar Gauntlets"},
	//	{ItemID: 23274, ItemName: "Knight-Lieutenant's Lamellar Gauntlets"},
	//	{ItemID: 227147, ItemName: "Knight-Lieutenant's Lamellar Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage absorbed by your Mana Shield by 285.
	// https://www.wowhead.com/forever/spell=23037
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16440, ItemName: "Marshal's Silk Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Sprint ability by 3.0 sec.
	// https://www.wowhead.com/forever/spell=23049
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16446, ItemName: "Marshal's Leather Footguards"},
	//	{ItemID: 231546, ItemName: "Marshal's Leather Footguards"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the Holy damage bonus of your Judgement of the Crusader by 20.
	// https://www.wowhead.com/forever/spell=23300
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16471, ItemName: "Marshal's Lamellar Gloves"},
	//	{ItemID: 231643, ItemName: "Marshal's Lamellar Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage absorbed by your Mana Shield by 285.
	// https://www.wowhead.com/forever/spell=23037
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16487, ItemName: "Blood Guard's Silk Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Sprint ability by 3.0 sec.
	// https://www.wowhead.com/forever/spell=23049
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16498, ItemName: "Blood Guard's Leather Treads"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hamstring Rage cost reduced by -3.0.
	// https://www.wowhead.com/forever/spell=22778
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16510, ItemName: "Blood Guard's Plate Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16518, ItemName: "Blood Guard's Mail Walkers"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage absorbed by your Mana Shield by 285.
	// https://www.wowhead.com/forever/spell=23037
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16540, ItemName: "General's Silk Handguards"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Sprint ability by 3.0 sec.
	// https://www.wowhead.com/forever/spell=23049
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16558, ItemName: "General's Leather Treads"},
	//	{ItemID: 231552, ItemName: "General's Leather Treads"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16573, ItemName: "General's Mail Boots"},
	//	{ItemID: 231667, ItemName: "General's Mail Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Poisons target for 8 Nature damage every 2.0 sec for 20s.
	// https://www.wowhead.com/forever/spell=17511
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17002, ItemName: "Ichor Spitter"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 84 Arcane damage.
	// https://www.wowhead.com/forever/spell=20883
	// unsupported: states no rate
	// trigger 20883 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 20883, BuffSpellID: 20883, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 17054, ItemName: "Joonho's Mercy"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 235 Mana when you kill a Dragonkin that gives experience. This effect cannot occur more than
	// once every 10 sec.
	// https://www.wowhead.com/forever/spell=1305426
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17064, ItemName: "Shard of the Scale"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=1305405
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17067, ItemName: "Ancient Cornerstone Grimoire"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy, causing 140 Shadow damage.
	// https://www.wowhead.com/forever/spell=1305398
	// unsupported: states no rate
	// trigger 1305398 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1305398, BuffSpellID: 1305398, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 17068, ItemName: "Deathbringer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 75 Shadow damage and lowering all stats by 25 for 30s.
	// https://www.wowhead.com/forever/spell=21151
	// unsupported: states no rate
	// trigger 21151 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 21151, BuffSpellID: 21151, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 17071, ItemName: "Gutgore Ripper"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Knocks down all nearby enemies for 3s.
	// https://www.wowhead.com/forever/spell=21152
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17073, ItemName: "Earthshaker"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Delivers a fatal wound for 236 Physical damage. Deals 75% increased damage to targets below 25% health.
	// https://www.wowhead.com/forever/spell=1305394
	// unsupported: states no rate
	// trigger 1305394 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1305394, BuffSpellID: 1305394, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 17075, ItemName: "Vis'kag the Bloodletter"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 150 mana or 20 rage when you kill a target that gives experience; this effect cannot occur more
	// than once every 10 seconds.
	// https://www.wowhead.com/forever/spell=21186
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17104, ItemName: "Spinal Reaper"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases your attack speed by 20% for 10s.
	// https://www.wowhead.com/forever/spell=21165
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17112, ItemName: "Empyrean Demolisher"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hurls a fiery ball that causes 92 Fire damage and an additional 16 damage over 8s.
	// https://www.wowhead.com/forever/spell=21159
	// unsupported: states no rate
	// trigger 21159 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 21159, BuffSpellID: 21159, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 17193, ItemName: "Sulfuron Hammer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts up to 3 targets for 200 Nature damage. Each target after the first takes less damage.
	// https://www.wowhead.com/forever/spell=21179
	// unsupported: states no rate
	// trigger 21179 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 21179, BuffSpellID: 21179, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 17223, ItemName: "Thunderstrike"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Searing Pain.
	// https://www.wowhead.com/forever/spell=23046
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17564, ItemName: "Knight-Lieutenant's Dreadweave Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Searing Pain.
	// https://www.wowhead.com/forever/spell=23046
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17577, ItemName: "Blood Guard's Dreadweave Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Searing Pain.
	// https://www.wowhead.com/forever/spell=23046
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17584, ItemName: "Marshal's Dreadweave Gloves"},
	//	{ItemID: 231586, ItemName: "Marshal's Dreadweave Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Searing Pain.
	// https://www.wowhead.com/forever/spell=23046
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17588, ItemName: "General's Dreadweave Gloves"},
	//	{ItemID: 231589, ItemName: "General's Dreadweave Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=23043
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17596, ItemName: "Knight-Lieutenant's Satin Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=23043
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17608, ItemName: "Marshal's Satin Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=23043
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17617, ItemName: "Blood Guard's Satin Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=23043
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17620, ItemName: "General's Satin Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 30 Frost damage.
	// https://www.wowhead.com/forever/spell=16407
	// unsupported: states no rate
	// trigger 16407 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 16407, BuffSpellID: 16407, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 17704, ItemName: "Edge of Winter"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Grants an extra attack on your next swing.
	// https://www.wowhead.com/forever/spell=21919
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17705, ItemName: "Thrash Blade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target causing them to bleed for 210 damage over 7s.
	// https://www.wowhead.com/forever/spell=21949
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17730, ItemName: "Gatorbite Axe"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 32 mana. Restore 1 times as much if attacking an Earth Elementals, Golems, or Titan-forged enermy.
	//
	// https://www.wowhead.com/forever/spell=21951
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17733, ItemName: "Fist of Stone"},
	//	{ItemID: 17943, ItemName: "Fist of Stone"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Poisons target for 13 Nature damage every 1.0 sec for 7s.
	// https://www.wowhead.com/forever/spell=21952
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17738, ItemName: "Claw of Celebras"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the duration of all poison effects applied to you by 10%.
	// https://www.wowhead.com/forever/spell=1294777
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17744, ItemName: "Heart of Noxxion"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 70 Shadow damage.
	// https://www.wowhead.com/forever/spell=1292155
	// unsupported: states no rate
	// trigger 1292155 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1292155, BuffSpellID: 1292155, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 17752, ItemName: "Satyr's Lash"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Strike your target with Keeper's Sting for 28 Nature damage.
	// https://www.wowhead.com/forever/spell=29655
	// unsupported: states no rate
	// trigger 29655 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 29655, BuffSpellID: 29655, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 17753, ItemName: "Verdant Keeper's Aim"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target for 168 damage and lowers their armor by 100.
	// https://www.wowhead.com/forever/spell=21961
	// unsupported: states no rate
	// trigger 21961 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 21961, BuffSpellID: 21961, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 17766, ItemName: "Princess Theradras' Scepter"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Has a 2% chance when struck in combat of increasing all stats by 21 for 1min.
	// https://www.wowhead.com/forever/spell=21970
	// unsupported: ProcTypeMask TAKE_HELPFUL_SPELL
	// trigger 21969 (2%, core.CallbackOnSpellHitTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage) -> buff 21970
	// shared.NewSpellDataProc(shared.SpellDataProc{TriggerSpellID: 21969, BuffSpellID: 21970},
	//	[]shared.ItemVariant{
	//	{ItemID: 17774, ItemName: "Mark of the Chosen"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance on landing a damaging spell to deal 112 Shadow damage and restore 104 mana to you.
	// https://www.wowhead.com/forever/spell=27860
	// unsupported: states no rate
	// trigger 21978 (no stated rate, core.CallbackOnCastComplete, core.ProcMaskSpellDamage) -> buff 27860
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 21978, BuffSpellID: 27860},
	//	[]shared.ItemVariant{
	//	{ItemID: 17780, ItemName: "Blade of Eternal Darkness"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Slows enemy's movement by 60% and causes them to bleed for 105 damage over 5s.
	// https://www.wowhead.com/forever/spell=22639
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18202, ItemName: "Eskhandar's Left Claw"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases your Attack Speed by 10% for 6s.
	// https://www.wowhead.com/forever/spell=22640
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18203, ItemName: "Eskhandar's Right Claw"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Decreases all threat generated by 4%.
	// https://www.wowhead.com/forever/spell=1298505
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18308, ItemName: "Clever Hat"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance on landing a damaging spell to restore 48 mana to you.
	// https://www.wowhead.com/forever/spell=1302191
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18321, ItemName: "Energetic Rod"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Slice the target, dealing 70 Nature damage. Deals 3 times as much damage to Aquatic enemies.
	// https://www.wowhead.com/forever/spell=1302166
	// unsupported: states no rate
	// trigger 1302166 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1302166, BuffSpellID: 1302166, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 18324, ItemName: "Waveslicer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Grants the wielder 13 Defense, 300 Armor, 30 Attack Power, and 36 additional Attack Power against Dragonkin
	// for 10s.
	// https://www.wowhead.com/forever/spell=22850
	// unsupported: states no rate
	// trigger 22850 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataProc(shared.SpellDataProc{TriggerSpellID: 22850, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 18348, ItemName: "Quel'Serrar"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 9 mana per 5 sec to your Imp.
	// https://www.wowhead.com/forever/spell=22855
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18354, ItemName: "Pimgib's Collar"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage done by your pets by 4%.
	// https://www.wowhead.com/forever/spell=22854
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18355, ItemName: "Ferra's Collar"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Shatter the target, dealing 49 Physical damage. Deals 3 times as much damage to Earth Elementals, Golems,
	// and Titan-forged.
	// https://www.wowhead.com/forever/spell=1302279
	// unsupported: states no rate
	// trigger 1302279 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1302279, BuffSpellID: 1302279, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 18388, ItemName: "Stoneshatter"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance on landing a damaging spell to assault the target's mind, increasing the casting time of all spells
	// by 10% for 30s. Stacks up to 5 times.
	// https://www.wowhead.com/forever/spell=1302342
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18396, ItemName: "Mind Carver"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases run speed by 35% for 10s.
	// https://www.wowhead.com/forever/spell=22863
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18410, ItemName: "Sprinter's Sword"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Minor increase to running and swimming speed. Does not stack with similar effects.
	// https://www.wowhead.com/forever/spell=24090
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18411, ItemName: "Spry Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the chance your Pick Pocket and Distract abilities are successful by 2%.
	// https://www.wowhead.com/forever/spell=1302345
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18463, ItemName: "Ogre Pocket Knife"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases Spirit by 16 in Cavernous or Underground areas.
	// https://www.wowhead.com/forever/spell=1318480
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18469, ItemName: "Royal Seal of Eldre'Thalas"},
	//	{ItemID: 18470, ItemName: "Royal Seal of Eldre'Thalas"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance on hit to restore 65 mana.
	// https://www.wowhead.com/forever/spell=1302350
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18483, ItemName: "Mana Channeling Wand"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Slice the target with magically enhanced precision, dealing 42 Arcane damage. Deals 3 times as much damage
	// to Plants.
	// https://www.wowhead.com/forever/spell=1302344
	// unsupported: states no rate
	// trigger 1302344 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1302344, BuffSpellID: 1302344, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 18498, ItemName: "Hedgecutter"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deal 175 Fire damage. Deals 3 times as much damage to Plants.
	// https://www.wowhead.com/forever/spell=1302354
	// unsupported: states no rate
	// trigger 1302354 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1302354, BuffSpellID: 1302354, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 18538, ItemName: "Treant's Bane"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the critical effect chance of your Holy spells by 2%.
	// https://www.wowhead.com/forever/spell=23236
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18608, ItemName: "Benediction"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 90 Fire damage.
	// https://www.wowhead.com/forever/spell=13442
	// unsupported: states no rate
	// trigger 13442 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 13442, BuffSpellID: 13442, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 18671, ItemName: "Baron Charr's Sceptre"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deal 42 Fire damage over 3s.
	// https://www.wowhead.com/forever/spell=1302343
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18755, ItemName: "Xorothian Firestick"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 48 Fire damage.
	// https://www.wowhead.com/forever/spell=23267
	// unsupported: states no rate
	// trigger 23267 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 23267, BuffSpellID: 23267, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 18816, ItemName: "Perdition's Blade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// While equipped, the wearer suffers less damage from falls.
	// https://www.wowhead.com/forever/spell=23409
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18951, ItemName: "Evonice's Landin' Pilla"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces your damage taken by 6% in the Battle Ring and The Maul.
	// https://www.wowhead.com/forever/spell=1318318
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19024, ItemName: "Arena Grand Master"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduce your threat to the current target making them less likely to attack you.
	// https://www.wowhead.com/forever/spell=23604
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19166, ItemName: "Black Amnesty"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Spell damage taken by target increased by 15% for 5s.
	// https://www.wowhead.com/forever/spell=23605
	// unsupported: states no rate
	// trigger 23605 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataAuraProc(shared.SpellDataProc{TriggerSpellID: 23605, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 19169, ItemName: "Nightfall"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 231 Shadow damage.
	// https://www.wowhead.com/forever/spell=18211
	// unsupported: states no rate
	// trigger 18211 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18211, BuffSpellID: 18211, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 19170, ItemName: "Ebon Hand"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives the wearer a 10% chance of being able to resurrect with 20% health and mana.
	// https://www.wowhead.com/forever/spell=23701
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19290, ItemName: "Darkmoon Card: Twisting Nether"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Stuns target for 1s.
	// https://www.wowhead.com/forever/spell=23454
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19323, ItemName: "The Unstoppable Force"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of Hammer of Justice by 0.5 sec.
	// https://www.wowhead.com/forever/spell=24188
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19588, ItemName: "Hero's Brand"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the amount of damage absorbed by Power Word: Shield by 35.
	// https://www.wowhead.com/forever/spell=24191
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19594, ItemName: "The All-Seeing Eye of Zuldazar"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the cooldown of Counterspell by -2.0 sec.
	// https://www.wowhead.com/forever/spell=24429
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19601, ItemName: "Jewel of Kajaro"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the radius of Rain of Fire and Hellfire by 1 yard.
	// https://www.wowhead.com/forever/spell=24430
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19605, ItemName: "Kezan's Unstoppable Taint"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Decreases the mana cost of your Healing Stream and Mana Spring totems by 20.
	// https://www.wowhead.com/forever/spell=24436
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19609, ItemName: "Unmarred Vision of Voodress"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the critical hit chance of Wrath and Starfire by 2%.
	// https://www.wowhead.com/forever/spell=24433
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19613, ItemName: "Pristine Enchanted South Seas Kelp"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Decreases the cooldown of Kick by -0.5 sec.
	// https://www.wowhead.com/forever/spell=24434
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19617, ItemName: "Zandalarian Shadow Mastery Talisman"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Decreases the cooldown of Feign Death by -2.0 sec.
	// https://www.wowhead.com/forever/spell=24432
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19621, ItemName: "Maelstrom's Wrath"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces damage from falling.
	// https://www.wowhead.com/forever/spell=1300663
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19982, ItemName: "Duskbat Drape"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your Frostbolt spells have a 6% chance to restore 50 mana when cast.
	// https://www.wowhead.com/forever/spell=24392
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20035, ItemName: "Glacial Spike"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20048, ItemName: "Highlander's Plate Greaves"},
	//	{ItemID: 20127, ItemName: "Highlander's Plate Greaves"},
	//	{ItemID: 20128, ItemName: "Highlander's Plate Greaves"},
	//	{ItemID: 20129, ItemName: "Highlander's Plate Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20049, ItemName: "Highlander's Lamellar Greaves"},
	//	{ItemID: 20109, ItemName: "Highlander's Lamellar Greaves"},
	//	{ItemID: 20110, ItemName: "Highlander's Lamellar Greaves"},
	//	{ItemID: 20111, ItemName: "Highlander's Lamellar Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20050, ItemName: "Highlander's Chain Greaves"},
	//	{ItemID: 20091, ItemName: "Highlander's Chain Greaves"},
	//	{ItemID: 20092, ItemName: "Highlander's Chain Greaves"},
	//	{ItemID: 20093, ItemName: "Highlander's Chain Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20051, ItemName: "Highlander's Mail Greaves"},
	//	{ItemID: 20121, ItemName: "Highlander's Mail Greaves"},
	//	{ItemID: 20122, ItemName: "Highlander's Mail Greaves"},
	//	{ItemID: 20123, ItemName: "Highlander's Mail Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20052, ItemName: "Highlander's Leather Boots"},
	//	{ItemID: 20112, ItemName: "Highlander's Leather Boots"},
	//	{ItemID: 20113, ItemName: "Highlander's Leather Boots"},
	//	{ItemID: 20114, ItemName: "Highlander's Leather Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20053, ItemName: "Highlander's Lizardhide Boots"},
	//	{ItemID: 20100, ItemName: "Highlander's Lizardhide Boots"},
	//	{ItemID: 20101, ItemName: "Highlander's Lizardhide Boots"},
	//	{ItemID: 20102, ItemName: "Highlander's Lizardhide Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20054, ItemName: "Highlander's Cloth Boots"},
	//	{ItemID: 20094, ItemName: "Highlander's Cloth Boots"},
	//	{ItemID: 20095, ItemName: "Highlander's Cloth Boots"},
	//	{ItemID: 20096, ItemName: "Highlander's Cloth Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20154, ItemName: "Defiler's Chain Greaves"},
	//	{ItemID: 20155, ItemName: "Defiler's Chain Greaves"},
	//	{ItemID: 20156, ItemName: "Defiler's Chain Greaves"},
	//	{ItemID: 20157, ItemName: "Defiler's Chain Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20159, ItemName: "Defiler's Cloth Boots"},
	//	{ItemID: 20160, ItemName: "Defiler's Cloth Boots"},
	//	{ItemID: 20161, ItemName: "Defiler's Cloth Boots"},
	//	{ItemID: 20162, ItemName: "Defiler's Cloth Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20167, ItemName: "Defiler's Lizardhide Boots"},
	//	{ItemID: 20168, ItemName: "Defiler's Lizardhide Boots"},
	//	{ItemID: 20169, ItemName: "Defiler's Lizardhide Boots"},
	//	{ItemID: 20170, ItemName: "Defiler's Lizardhide Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20181, ItemName: "Defiler's Lamellar Greaves"},
	//	{ItemID: 20182, ItemName: "Defiler's Lamellar Greaves"},
	//	{ItemID: 20183, ItemName: "Defiler's Lamellar Greaves"},
	//	{ItemID: 20185, ItemName: "Defiler's Lamellar Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20186, ItemName: "Defiler's Leather Boots"},
	//	{ItemID: 20187, ItemName: "Defiler's Leather Boots"},
	//	{ItemID: 20188, ItemName: "Defiler's Leather Boots"},
	//	{ItemID: 20189, ItemName: "Defiler's Leather Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20199, ItemName: "Defiler's Mail Greaves"},
	//	{ItemID: 20200, ItemName: "Defiler's Mail Greaves"},
	//	{ItemID: 20201, ItemName: "Defiler's Mail Greaves"},
	//	{ItemID: 20202, ItemName: "Defiler's Mail Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20208, ItemName: "Defiler's Plate Greaves"},
	//	{ItemID: 20209, ItemName: "Defiler's Plate Greaves"},
	//	{ItemID: 20210, ItemName: "Defiler's Plate Greaves"},
	//	{ItemID: 20211, ItemName: "Defiler's Plate Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=1300667
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20503, ItemName: "Enamored Water Spirit"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck in combat, inflicts 3 Physical damage to the attacker.
	// https://www.wowhead.com/forever/spell=1300680
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20517, ItemName: "Razorsteel Shoulders"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Disorients the target, causing it to wander aimlessly for up to 3s.
	// https://www.wowhead.com/forever/spell=26108
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 21134, ItemName: "Dark Edge of Insanity"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=1318470
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 21326, ItemName: "Defender of the Timbermaw"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Underwater Breath lasts 50% longer than normal.
	// https://www.wowhead.com/forever/spell=11789
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 21526, ItemName: "Band of Icy Depths"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases your pet's maximum health by 3%.
	// https://www.wowhead.com/forever/spell=27038
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22013, ItemName: "Beastmaster's Cap"},
	//	{ItemID: 226887, ItemName: "Beastmaster's Cap"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases your pet's critical strike chance by 2%.
	// https://www.wowhead.com/forever/spell=27043
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22015, ItemName: "Beastmaster's Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// On successful melee or ranged attack gain 8 mana and if possible drain 8 mana from the target.
	// https://www.wowhead.com/forever/spell=18350
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22194, ItemName: "Black Grasp of the Destroyer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck by a non-periodic damage spell you have a 30% chance of getting a 6s spell shield that absorbs
	// 400 of that school of damage.
	// https://www.wowhead.com/forever/spell=27539
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22196, ItemName: "Thick Obsidian Breastplate"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck by a harmful spell, the caster of that spell has a 5% chance to be silenced for 3s.
	// https://www.wowhead.com/forever/spell=27559
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22198, ItemName: "Jagged Obsidian Shield"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases healing done by Lesser Healing Wave by up to 80.
	// https://www.wowhead.com/forever/spell=27855
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22396, ItemName: "Totem of Life"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the initial healing of Regrowth by 5% when cast on a target affected by your Rejuvenation.
	// https://www.wowhead.com/forever/spell=27853
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22398, ItemName: "Idol of Rejuvenation"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the casting time of your Healing Touch spell by 0.15 sec.
	// https://www.wowhead.com/forever/spell=27846
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22399, ItemName: "Idol of Health"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the Mana cost of Purify and Cleanse by 10%.
	// https://www.wowhead.com/forever/spell=27850
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22400, ItemName: "Libram of Truth"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the mana cost of your Cleanse spell by 25.
	// https://www.wowhead.com/forever/spell=27847
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22402, ItemName: "Libram of Grace"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Harmful spell casts have a chance to deal 91 Arcane or Shadow damage.
	// https://www.wowhead.com/forever/spell=1300435
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22458, ItemName: "Moonshadow Stave"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the spell critical chance of all party members within 30 yards by 2%.
	// https://www.wowhead.com/forever/spell=28142
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22589, ItemName: "Atiesh, Greatstaff of the Guardian"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Sprint ability by 3.0 sec.
	// https://www.wowhead.com/forever/spell=23049
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22856, ItemName: "Blood Guard's Leather Walkers"},
	//	{ItemID: 227062, ItemName: "Blood Guard's Leather Walkers"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22857, ItemName: "Blood Guard's Mail Greaves"},
	//	{ItemID: 227158, ItemName: "Blood Guard's Mail Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Searing Pain.
	// https://www.wowhead.com/forever/spell=23046
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22865, ItemName: "Blood Guard's Dreadweave Handwraps"},
	//	{ItemID: 227099, ItemName: "Blood Guard's Dreadweave Handwraps"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hamstring Rage cost reduced by -3.0.
	// https://www.wowhead.com/forever/spell=22778
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22868, ItemName: "Blood Guard's Plate Gauntlets"},
	//	{ItemID: 227050, ItemName: "Blood Guard's Plate Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=23043
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22869, ItemName: "Blood Guard's Satin Handwraps"},
	//	{ItemID: 227126, ItemName: "Blood Guard's Satin Handwraps"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage absorbed by your Mana Shield by 285.
	// https://www.wowhead.com/forever/spell=23037
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22870, ItemName: "Blood Guard's Silk Handwraps"},
	//	{ItemID: 227111, ItemName: "Blood Guard's Silk Handwraps"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gain up to 25 mana each time you cast Healing Touch.
	// https://www.wowhead.com/forever/spell=28847
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23004, ItemName: "Idol of Longevity"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Regain up to 10 mana each time you cast Lesser Healing Wave.
	// https://www.wowhead.com/forever/spell=28849
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23005, ItemName: "Totem of Flowing Water"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Causes Judgement to heal all nearby party members for 25.
	// https://www.wowhead.com/forever/spell=1302540
	// unsupported: the heal lands on implicit target 20, not the wearer
	// trigger 28853 (every time, core.CallbackOnCastComplete, core.ProcMaskSpellDamage) -> buff 1302540
	// shared.NewSpellDataHealProc(shared.SpellDataProc{TriggerSpellID: 28853, BuffSpellID: 1302540},
	//	[]shared.ItemVariant{
	//	{ItemID: 23201, ItemName: "Libram of Divinity"},
	//	{ItemID: 23202, ItemName: "Libram of Divinity"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Searing Pain.
	// https://www.wowhead.com/forever/spell=23046
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23282, ItemName: "Knight-Lieutenant's Dreadweave Handwraps"},
	//	{ItemID: 227100, ItemName: "Knight-Lieutenant's Dreadweave Handwraps"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Sprint ability by 3.0 sec.
	// https://www.wowhead.com/forever/spell=23049
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23285, ItemName: "Knight-Lieutenant's Leather Walkers"},
	//	{ItemID: 227064, ItemName: "Knight-Lieutenant's Leather Walkers"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hamstring Rage cost reduced by -3.0.
	// https://www.wowhead.com/forever/spell=22778
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23286, ItemName: "Knight-Lieutenant's Plate Gauntlets"},
	//	{ItemID: 227053, ItemName: "Knight-Lieutenant's Plate Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=23043
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23288, ItemName: "Knight-Lieutenant's Satin Handwraps"},
	//	{ItemID: 227128, ItemName: "Knight-Lieutenant's Satin Handwraps"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage absorbed by your Mana Shield by 285.
	// https://www.wowhead.com/forever/spell=23037
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23290, ItemName: "Knight-Lieutenant's Silk Handwraps"},
	//	{ItemID: 227113, ItemName: "Knight-Lieutenant's Silk Handwraps"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Prevents an attack that would otherwise kill you. Triggering this effect also grants you 3s of damage
	// immunity and shatters the phylactery.
	// https://www.wowhead.com/forever/spell=370391
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 191312, ItemName: "Failsafe Phylactery"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gain the Quick Draw ability:
	//
	// |Tinterface/icons/inv_musket_02.blp:24|t
	// Draw your ranged weapon and fire a quick shot at an enemy, causing normal ranged weapon damage and reducing
	// the target's movement speed by 50% for 6s. Awards 1 combo point.
	//
	// Quick Draw benefits from all talents and effects that trigger from or modify Sinister Strike.
	//
	// https://www.wowhead.com/forever/spell=398197
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 202256, ItemName: "Privateer's Ornate Pistol"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 205218, ItemName: "Libram of Discovery"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 205420, ItemName: "Libram of Judgement"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 206376, ItemName: "Bulwark Icon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 206381, ItemName: "Dyadic Icon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 206382, ItemName: "Tempest Icon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 206383, ItemName: "Lithic Icon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 206384, ItemName: "Astral Icon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 206385, ItemName: "Clastic Icon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 206386, ItemName: "Galvanic Icon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 206387, ItemName: "Kajaric Icon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 206388, ItemName: "Sulfurous Icon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 206389, ItemName: "Brimstone Icon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 206390, ItemName: "Ichthyan Icon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 206391, ItemName: "Syzygic Icon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=414827
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 206954, ItemName: "Idol of Ursine Rage"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// On landing a killing blow that grants experience or honor, your next attack will critically strike.
	// https://www.wowhead.com/forever/spell=418509
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 208222, ItemName: "Old Guard Retaliator"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 208414, ItemName: "Lunar Idol"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=463001
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 208424, ItemName: "Sun Shades"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 208689, ItemName: "Ferocious Idol"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 208849, ItemName: "Libram of Blessings"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 208851, ItemName: "Libram of Justice"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=423191
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 210195, ItemName: "Unbalanced Idol"},
	//	{ItemID: 210195, ItemName: "Unbalanced Idol"},
	//	{ItemID: 210195, ItemName: "Unbalanced Idol"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 210534, ItemName: "Idol of the Wild"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 211472, ItemName: "Libram of Banishment"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts enemies in front of you with the power of wind, fire, all that kind of thing!
	// https://www.wowhead.com/forever/spell=14537
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 211941, ItemName: "Windwalker's Yari"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 213513, ItemName: "Libram of Deliverance"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 213594, ItemName: "Idol of the Heckler"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of Rip by 2 sec.
	// https://www.wowhead.com/forever/spell=446212
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 220606, ItemName: "Idol of the Dream"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Minor increase to running and swimming speed. Does not stack with similar effects.
	// https://www.wowhead.com/forever/spell=24090
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 220835, ItemName: "First Sergeant's Mail Sabatons"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 220840, ItemName: "First Sergeant's Inscribed Sabatons"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 220846, ItemName: "First Sergeant's Pulsing Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 220915, ItemName: "Idol of the Raging Shambler"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Minor increase to running and swimming speed. Does not stack with similar effects.
	// https://www.wowhead.com/forever/spell=24090
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 223077, ItemName: "Sergeant Major's Mail Sabatons"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 225838, ItemName: "Voltaic Icon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases your pet's critical strike chance by 2%.
	// https://www.wowhead.com/forever/spell=27043
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 226883, ItemName: "Beastmaster's Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage done by your Multi-Shot by 4%.
	// https://www.wowhead.com/forever/spell=459593
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 227075, ItemName: "Blood Guard's Chain Vices"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage done by your Multi-Shot by 4%.
	// https://www.wowhead.com/forever/spell=459593
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 227077, ItemName: "Knight-Lieutenant's Chain Vices"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage done by your Raptor Strike by 4%.
	// https://www.wowhead.com/forever/spell=459598
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 227081, ItemName: "Blood Guard's Chain Grips"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage done by your Raptor Strike by 4%.
	// https://www.wowhead.com/forever/spell=459598
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 227087, ItemName: "Knight-Lieutenant's Chain Grips"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Arcane Blast.
	// https://www.wowhead.com/forever/spell=459600
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 227114, ItemName: "Knight-Lieutenant's Silk Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Arcane Blast.
	// https://www.wowhead.com/forever/spell=459600
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 227115, ItemName: "Blood Guard's Silk Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=459604
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 227133, ItemName: "Blood Guard's Satin Grips"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=459604
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 227139, ItemName: "Knight-Lieutenant's Satin Grips"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the critical strike chance of your Holy Shock by 2%.
	// https://www.wowhead.com/forever/spell=459602
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 227152, ItemName: "Knight-Lieutenant's Lamellar Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=459606
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 227164, ItemName: "Blood Guard's Mail Sabatons"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=459606
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 227170, ItemName: "Blood Guard's Mail Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the mana cost of your shapeshifts by 150.
	// https://www.wowhead.com/forever/spell=459594
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 227180, ItemName: "Blood Guard's Dragonhide Grips"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the mana cost of your shapeshifts by 150.
	// https://www.wowhead.com/forever/spell=459594
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 227183, ItemName: "Knight-Lieutenant's Dragonhide Grips"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Wrath.
	// https://www.wowhead.com/forever/spell=459595
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 227187, ItemName: "Blood Guard's Dragonhide Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Wrath.
	// https://www.wowhead.com/forever/spell=459595
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 227193, ItemName: "Knight-Lieutenant's Dragonhide Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Barkskin by 3 sec.
	// https://www.wowhead.com/forever/spell=459596
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 227198, ItemName: "Knight-Lieutenant's Dragonhide Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Barkskin by 3 sec.
	// https://www.wowhead.com/forever/spell=459596
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 227204, ItemName: "Blood Guard's Dragonhide Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=408953
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 227444, ItemName: "Idol of the Huntress"},
	//	{ItemID: 227444, ItemName: "Idol of the Huntress"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hamstring Rage cost reduced by -3.0.
	// https://www.wowhead.com/forever/spell=459608
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231532, ItemName: "General's Plate Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hamstring Rage cost reduced by -3.0.
	// https://www.wowhead.com/forever/spell=459608
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231541, ItemName: "Marshal's Plate Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage done by your Raptor Strike by 4%.
	// https://www.wowhead.com/forever/spell=459598
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231560, ItemName: "Marshal's Chain Grips"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage done by your Raptor Strike by 4%.
	// https://www.wowhead.com/forever/spell=459598
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231569, ItemName: "General's Chain Grips"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage done by your Multi-Shot by 4%.
	// https://www.wowhead.com/forever/spell=459593
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231575, ItemName: "General's Chain Vices"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage done by your Multi-Shot by 4%.
	// https://www.wowhead.com/forever/spell=459593
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231578, ItemName: "Marshal's Chain Vices"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Arcane Blast.
	// https://www.wowhead.com/forever/spell=459600
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231599, ItemName: "General's Silk Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage absorbed by your Mana Shield by 285.
	// https://www.wowhead.com/forever/spell=459599
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231600, ItemName: "General's Silk Handwraps"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Arcane Blast.
	// https://www.wowhead.com/forever/spell=459600
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231608, ItemName: "Marshal's Silk Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage absorbed by your Mana Shield by 285.
	// https://www.wowhead.com/forever/spell=459599
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231609, ItemName: "Marshal's Silk Handwraps"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=459604
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231613, ItemName: "General's Satin Grips"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=459604
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231617, ItemName: "Marshal's Satin Grips"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the duration of your Weakened Soul by 2 sec.
	// https://www.wowhead.com/forever/spell=459603
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231623, ItemName: "Marshal's Satin Handwraps"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the duration of your Weakened Soul by 2 sec.
	// https://www.wowhead.com/forever/spell=459603
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231633, ItemName: "General's Satin Handwraps"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your damaging Judgements deal 20 additional damage.
	// https://www.wowhead.com/forever/spell=459601
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231650, ItemName: "Marshal's Lamellar Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=459606
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231656, ItemName: "General's Mail Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=459606
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231661, ItemName: "General's Mail Sabatons"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Barkskin by 3 sec.
	// https://www.wowhead.com/forever/spell=459596
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231676, ItemName: "General's Dragonhide Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Wrath.
	// https://www.wowhead.com/forever/spell=459595
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231677, ItemName: "General's Dragonhide Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the mana cost of your shapeshifts by 150.
	// https://www.wowhead.com/forever/spell=459594
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231688, ItemName: "General's Dragonhide Grips"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the mana cost of your shapeshifts by 150.
	// https://www.wowhead.com/forever/spell=459594
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231694, ItemName: "Marshal's Dragonhide Grips"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Wrath.
	// https://www.wowhead.com/forever/spell=459595
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231700, ItemName: "Marshal's Dragonhide Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Barkskin by 3 sec.
	// https://www.wowhead.com/forever/spell=459596
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231706, ItemName: "Marshal's Dragonhide Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=469141
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 231877, ItemName: "Tyler's Balanced Bobble"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the casting speed of your spells by 2% per piece of Timeworn armor equipped.
	// https://www.wowhead.com/forever/spell=1213398
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234016, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234017, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234017, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234018, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234019, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234020, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234021, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234021, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234022, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234023, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234024, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234025, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234026, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234026, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234027, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234028, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234029, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234030, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234030, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234031, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234032, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234033, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234034, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234034, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234035, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234198, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234198, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234199, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234199, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234200, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234200, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234201, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234201, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234202, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234202, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234436, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234437, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234438, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234439, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234440, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234964, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234965, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234966, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234967, ItemName: "Signet Ring of the Bronze Dragonflight"},
	//	{ItemID: 234968, ItemName: "Signet Ring of the Bronze Dragonflight"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234542, ItemName: "High Warlord's Greatsword"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234543, ItemName: "High Warlord's Battle Axe"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234545, ItemName: "High Warlord's Pulverizer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234546, ItemName: "High Warlord's Destroyer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234547, ItemName: "High Warlord's Pig Sticker"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234548, ItemName: "High Warlord's Pig Poker"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234549, ItemName: "High Warlord's War Staff"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234550, ItemName: "High Warlord's Spellblade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234551, ItemName: "High Warlord's Battle Mace"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234552, ItemName: "High Warlord's Blade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234553, ItemName: "High Warlord's Quickblade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234554, ItemName: "High Warlord's Cleaver"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234555, ItemName: "High Warlord's Bludgeon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234556, ItemName: "High Warlord's Razor"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234557, ItemName: "High Warlord's Right Claw"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234558, ItemName: "High Warlord's Left Claw"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234559, ItemName: "High Warlord's Recurve"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234560, ItemName: "High Warlord's Crossbow"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234561, ItemName: "High Warlord's Street Sweeper"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234562, ItemName: "High Warlord's Shield Wall -  - "},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234563, ItemName: "High Warlord's Tome of Destruction"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234564, ItemName: "High Warlord's Tome of Mending"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234565, ItemName: "Grand Marshal's Claymore"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234566, ItemName: "Grand Marshal's Sunderer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234567, ItemName: "Grand Marshal's Battle Hammer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234568, ItemName: "Grand Marshal's Demolisher"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234569, ItemName: "Grand Marshal's Glaive"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234570, ItemName: "Grand Marshal's Polearm"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234571, ItemName: "Grand Marshal's Stave"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234574, ItemName: "Grand Marshal's Mageblade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234576, ItemName: "Grand Marshal's Warhammer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234578, ItemName: "Grand Marshal's Longsword"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234579, ItemName: "Grand Marshal's Swiftblade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234580, ItemName: "Grand Marshal's Handaxe"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234581, ItemName: "Grand Marshal's Punisher"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234582, ItemName: "Grand Marshal's Dirk"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234583, ItemName: "Grand Marshal's Right Hand Blade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234584, ItemName: "Grand Marshal's Left Hand Blade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234585, ItemName: "Grand Marshal's Bullseye"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234586, ItemName: "Grand Marshal's Repeater"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234587, ItemName: "Grand Marshal's Hand Cannon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234588, ItemName: "Grand Marshal's Aegis -  - "},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234589, ItemName: "Grand Marshal's Tome of Power"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 234590, ItemName: "Grand Marshal's Tome of Restoration"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 235473, ItemName: "Grand Marshal's Barricade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 235474, ItemName: "High Warlord's Barricade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 235476, ItemName: "High Warlord's Hacker"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 235477, ItemName: "High Warlord's Bonecracker"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 235478, ItemName: "High Warlord's Shiv"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 235479, ItemName: "Grand Marshal's Shiv"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 235480, ItemName: "Grand Marshal's Bonecracker"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces all damage taken from players and their pets by 1%. This effect cannot be combined from multiple
	// sources.
	// https://www.wowhead.com/forever/spell=1216997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 235481, ItemName: "Grand Marshal's Hacker"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage dealt by your damage over time spells by 2%.
	// https://www.wowhead.com/forever/spell=1222994
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 239565, ItemName: "Garb of Revelation (Sanctified)"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage dealt by your damage over time spells by 2%.
	// https://www.wowhead.com/forever/spell=1222994
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 239572, ItemName: "Boots of Revelation (Sanctified)"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage dealt by your damage over time spells by 2%.
	// https://www.wowhead.com/forever/spell=1222994
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 239574, ItemName: "Hands of Revelation (Sanctified)"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage dealt by your damage over time spells by 2%.
	// https://www.wowhead.com/forever/spell=1222994
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 239575, ItemName: "Crown of Revelation (Sanctified)"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage dealt by your damage over time spells by 2%.
	// https://www.wowhead.com/forever/spell=1222994
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 239577, ItemName: "Pants of Revelation (Sanctified)"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage dealt by your damage over time spells by 4%.
	// https://www.wowhead.com/forever/spell=1222997
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 239581, ItemName: "Mantle of Revelation (Sanctified)"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage dealt by your damage over time spells by 2%.
	// https://www.wowhead.com/forever/spell=1222994
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 239582, ItemName: "Girdle of Revelation (Sanctified)"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage dealt by your damage over time spells by 2%.
	// https://www.wowhead.com/forever/spell=1222994
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 239583, ItemName: "Wrists of Revelation (Sanctified)"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Allows 8% of your Mana regeneration to continue while casting.
	// https://www.wowhead.com/forever/spell=1248756
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 249398, ItemName: "Polished Driftwood Icon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Causes Wrath to have a 50% chance to restore 35 Mana.
	// https://www.wowhead.com/forever/spell=1302521
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 249441, ItemName: "Talons of Wrath"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the cooldown of Grounding Totem by 1 sec.
	// https://www.wowhead.com/forever/spell=1249011
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 249443, ItemName: "Totem of Ancestral Protection"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your casts of Greater Heal in combat grant up to 40 increased healing and up to 13 increased damage for
	// 15s.
	// https://www.wowhead.com/forever/spell=1249119
	// unsupported: named ability
	// trigger 1249118 (every time, core.CallbackOnHealDealt, core.ProcMaskSpellHealing) -> buff 1249119
	// shared.NewSpellDataProc(shared.SpellDataProc{TriggerSpellID: 1249118, BuffSpellID: 1249119},
	//	[]shared.ItemVariant{
	//	{ItemID: 249473, ItemName: "Dormant Heart of the Mountain"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases movement speed by 2% in Silverpine Forest and Hillsbrad Foothills.
	// https://www.wowhead.com/forever/spell=1309369
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 251533, ItemName: "Forsaken Greataxe"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases movement speed by 2% in Silverpine Forest and Hillsbrad Foothills.
	// https://www.wowhead.com/forever/spell=1309369
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 251534, ItemName: "Gnarled Necromancer's Staff"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sends a shadowy bolt at the enemy causing 140 Shadow damage.
	// https://www.wowhead.com/forever/spell=16784
	// unsupported: states no rate
	// trigger 16784 (every time, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 16784, BuffSpellID: 16784, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 260186, ItemName: "Darkblade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance on hit to deal 75 Arcane damage. Deals 2 times as much damage to Naga and Satyrs.
	// https://www.wowhead.com/forever/spell=1265634
	// unsupported: states no rate
	// trigger 1318159 (no stated rate, core.CallbackOnSpellHitDealt, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage) -> buff 1265634
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1318159, BuffSpellID: 1265634},
	//	[]shared.ItemVariant{
	//	{ItemID: 260205, ItemName: "Highborne Research Tablet"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the mana cost of shapeshifting spells by 40.
	// https://www.wowhead.com/forever/spell=1270470
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 263411, ItemName: "Idol of Shifting Tides"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the mana cost of Flame Shock and Lightning Bolt by 5 Mana.
	// https://www.wowhead.com/forever/spell=1270478
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 263412, ItemName: "Totem of Charged Flames"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the mana cost of your totems by 5%.
	// https://www.wowhead.com/forever/spell=1270492
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 263436, ItemName: "Firestorm Totem"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sear certain enchanted targets with blazing light causing them to take additional damage from Holy attacks
	// and spells.
	// https://www.wowhead.com/forever/spell=1282482
	// unsupported: states no rate; effect 1 lands on implicit target 38
	// trigger 1282482 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataAuraProc(shared.SpellDataProc{TriggerSpellID: 1282482, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 268484, ItemName: "Moonsilver Blade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=1318514
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 268873, ItemName: "Defender of the Barkskin"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Harmful spell casts and attacks against Furbolg have a chance to deal 105 Nature damage.
	// https://www.wowhead.com/forever/spell=1318324
	// unsupported: states no rate
	// trigger 1318323 (no stated rate, core.CallbackOnSpellHitDealt, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage) -> buff 1318324
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1318323, BuffSpellID: 1318324},
	//	[]shared.ItemVariant{
	//	{ItemID: 269741, ItemName: "Scented Runewood Brooch"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your Enrage ability generates an additional 10 Rage over its duration.
	// https://www.wowhead.com/forever/spell=1291060
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272428, ItemName: "Enraged Idol"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Each of your heal over time effects on the target reduces the remaining cooldown on Swiftmend by 1 sec
	// when you cast it.
	// https://www.wowhead.com/forever/spell=1291075
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272429, ItemName: "Idol of Synthesis"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Insect Swarm ability by 2 sec.
	// https://www.wowhead.com/forever/spell=1291061
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272430, ItemName: "Swarming Idol"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the Mana cost of your Healing Wave ability by 5%.
	// https://www.wowhead.com/forever/spell=1291076
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272431, ItemName: "Tidal Totem"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your Lightning Bolt ability can now also trigger the Maelstrom Weapon talent, but with a 50% reduced chance.
	// https://www.wowhead.com/forever/spell=1291078
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272432, ItemName: "Totem of the Storm"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Flame Shock ability by 3 sec.
	// https://www.wowhead.com/forever/spell=1291077
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272433, ItemName: "Burning Totem"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage absorbed by your Mana Shield by 285.
	// https://www.wowhead.com/forever/spell=23037
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272455, ItemName: "Premier Magus Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hamstring Rage cost reduced by -3.0.
	// https://www.wowhead.com/forever/spell=22778
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272478, ItemName: "Premier Mortarplate Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272484, ItemName: "Premier Mail Walkers"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage done by your Multi-Shot by 4%.
	// https://www.wowhead.com/forever/spell=28539
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272495, ItemName: "Premier Chain Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage absorbed by your Mana Shield by 285.
	// https://www.wowhead.com/forever/spell=23037
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272505, ItemName: "Premier Champion's Magus Handguards"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hamstring Rage cost reduced by -3.0.
	// https://www.wowhead.com/forever/spell=22778
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272513, ItemName: "Premier Champion's Mortarplate Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Sprint ability by 3.0 sec.
	// https://www.wowhead.com/forever/spell=23049
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272523, ItemName: "Premier Centurion's Shadowhide Treads"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage done by your Multi-Shot by 4%.
	// https://www.wowhead.com/forever/spell=28539
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272536, ItemName: "Premier Champion's Chain Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272538, ItemName: "Premier Centurion's Linked Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Searing Pain.
	// https://www.wowhead.com/forever/spell=23046
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272553, ItemName: "Premier Felweave Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Searing Pain.
	// https://www.wowhead.com/forever/spell=23046
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272556, ItemName: "Premier Champion's Felweave Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=23043
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272569, ItemName: "Premier Satin Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=23043
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272572, ItemName: "Premier Champion's Satin Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Sprint ability by 3.0 sec.
	// https://www.wowhead.com/forever/spell=23049
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272608, ItemName: "Premier Shadowhide Walkers"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272609, ItemName: "Premier Linked Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Flash Heal.
	// https://www.wowhead.com/forever/spell=1293542
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272619, ItemName: "Premier Silken Handwraps"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage absorbed by your Mana Shield by 285.
	// https://www.wowhead.com/forever/spell=23037
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272703, ItemName: "Premier Silk Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Sprint ability by 3.0 sec.
	// https://www.wowhead.com/forever/spell=23049
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272704, ItemName: "Premier Leather Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage done by your Multi-Shot by 4%.
	// https://www.wowhead.com/forever/spell=28539
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272715, ItemName: "Premier Chainmail Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hamstring Rage cost reduced by -3.0.
	// https://www.wowhead.com/forever/spell=22778
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272717, ItemName: "Premier Plate Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Repentance by 1 sec.
	// https://www.wowhead.com/forever/spell=1316101
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272719, ItemName: "Premier Chevalier Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage absorbed by your Mana Shield by 285.
	// https://www.wowhead.com/forever/spell=23037
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272749, ItemName: "Premier Lieutenant Commander's Silk Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Sprint ability by 3.0 sec.
	// https://www.wowhead.com/forever/spell=23049
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272755, ItemName: "Premier Knight-Champion's Leather Footguards"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage done by your Multi-Shot by 4%.
	// https://www.wowhead.com/forever/spell=28539
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272772, ItemName: "Premier Lieutenant Commander's Chainmail Grips"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Repentance by 1 sec.
	// https://www.wowhead.com/forever/spell=1316101
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272780, ItemName: "Premier Lieutenant Commander's Chevalier Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hamstring Rage cost reduced by -3.0.
	// https://www.wowhead.com/forever/spell=22778
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272793, ItemName: "Premier Lieutenant Commander's Plate Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Searing Pain.
	// https://www.wowhead.com/forever/spell=23046
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272796, ItemName: "Premier Dreadweave Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Searing Pain.
	// https://www.wowhead.com/forever/spell=23046
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272808, ItemName: "Premier Lieutenant Commander's Dreadweave Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=23043
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272812, ItemName: "Premier Voidcloth Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=23043
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 272824, ItemName: "Premier Lieutenant Commander's Voidcloth Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck in combat, drain 2 mana from the attacker.
	// https://www.wowhead.com/forever/spell=1292706
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 273041, ItemName: "Slitherwrap Girdle"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When damaged, has a chance to reduce threat level on all enemies by a small amount for 10s.
	// https://www.wowhead.com/forever/spell=12685
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 273044, ItemName: "Violet Sorcerer's Robes"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Swim speed increased by 33%.
	// https://www.wowhead.com/forever/spell=1291749
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 273084, ItemName: "Cloak of Hermitic Bliss"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deals 6 Physical damage every time you block.
	// https://www.wowhead.com/forever/spell=1292039
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 273293, ItemName: "Bandsaw Wristbands"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the critical strike chance of your Flash of Light by 2%.
	//
	// https://www.wowhead.com/forever/spell=1293544
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 273329, ItemName: "Premier Lieutenant Commander's Luminous Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 273336, ItemName: "Premier Centurion's Mail Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 273344, ItemName: "Premier Centurion's Ringmail Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 273355, ItemName: "Premier Ringmail Walkers"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the critical strike chance of your Flash of Light by 2%.
	//
	// https://www.wowhead.com/forever/spell=1293544
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 273363, ItemName: "Premier Luminous Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Flash Heal.
	// https://www.wowhead.com/forever/spell=1293542
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 273427, ItemName: "Premier Lieutenant Commander's Mooncloth Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Flash Heal.
	// https://www.wowhead.com/forever/spell=1293542
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 273431, ItemName: "Premier Champion's Silken Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Flash Heal.
	// https://www.wowhead.com/forever/spell=1293542
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 273445, ItemName: "Premier Mooncloth Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Wounds the target for 70 Physical damage.
	// https://www.wowhead.com/forever/spell=1292573
	// unsupported: states no rate
	// trigger 1292573 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1292573, BuffSpellID: 1292573, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 273819, ItemName: "Boneslicer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases movement speed by 2% in Duskwood.
	// https://www.wowhead.com/forever/spell=1292575
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 273820, ItemName: "Nightskulker Ring"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Lockpicking +4.
	// https://www.wowhead.com/forever/spell=1292581
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 273824, ItemName: "Defias Jailbreakers"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Grants 6 health and 2 mana every 5 sec for 15s when one of your spells is resisted.
	// https://www.wowhead.com/forever/spell=1292598
	// unsupported: an outcome the proc mask has no bit for
	// trigger 1292594 (every time, core.CallbackOnCastComplete, core.ProcMaskSpellDamage) -> buff 1292598
	// shared.NewSpellDataProc(shared.SpellDataProc{TriggerSpellID: 1292594, BuffSpellID: 1292598},
	//	[]shared.ItemVariant{
	//	{ItemID: 273827, ItemName: "Debt Collector"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck in combat, inflicts 2 Nature damage to the attacker.
	// https://www.wowhead.com/forever/spell=1292670
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 273839, ItemName: "Spiked Shell Band"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 70 Nature damage.
	// https://www.wowhead.com/forever/spell=1292679
	// unsupported: states no rate
	// trigger 1292679 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1292679, BuffSpellID: 1292679, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 273841, ItemName: "Twilight Maul"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases all threat generated by 1%.
	// https://www.wowhead.com/forever/spell=1292683
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 273842, ItemName: "Treacherous Treads"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 274176, ItemName: "Premier Scalemail Walkers"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 274183, ItemName: "Premier Knight-Champion's Flatmail Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 274192, ItemName: "Premier Flatmail Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 274198, ItemName: "Premier Knight-Champion's Scalemail Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 274206, ItemName: "Premier Knight-Champion's Linkmail Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 274217, ItemName: "Premier Linkmail Walkers"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Repentance by 1 sec.
	// https://www.wowhead.com/forever/spell=1316101
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 274227, ItemName: "Premier Scaled Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Repentance by 1 sec.
	// https://www.wowhead.com/forever/spell=1316101
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 274236, ItemName: "Premier Champion's Scaled Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the critical strike chance of your Flash of Light by 2%.
	//
	// https://www.wowhead.com/forever/spell=1293544
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 274250, ItemName: "Premier Champion's Lamellar Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the critical strike chance of your Flash of Light by 2%.
	//
	// https://www.wowhead.com/forever/spell=1293544
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 274256, ItemName: "Premier Lamellar Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Launch a magic severing bolt at the target, dealing 36 Physical damage and interrupting the spell being
	// cast for 5s.
	// https://www.wowhead.com/forever/spell=1293714
	// unsupported: states no rate
	// trigger 1293714 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1293714, BuffSpellID: 1293714, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 274293, ItemName: "Spellsever Crossbow"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance on hit to reduce the mana cost of your next spell cast within 10s by 24.
	// https://www.wowhead.com/forever/spell=1294058
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 274425, ItemName: "White Obsidian Wand"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// While in combat, your healing spells have a chance to grant an ally 300 armor for 10s.
	// https://www.wowhead.com/forever/spell=1294065
	// unsupported: states no rate
	// trigger 1294064 (no stated rate, core.CallbackOnHealDealt, core.ProcMaskSpellHealing) -> buff 1294065
	// shared.NewSpellDataProc(shared.SpellDataProc{TriggerSpellID: 1294064, BuffSpellID: 1294065},
	//	[]shared.ItemVariant{
	//	{ItemID: 274430, ItemName: "Rock Sentinel Slicer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=1295267
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 274747, ItemName: "Soggy Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=1318317
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 274761, ItemName: "Parachute-Priest Pager"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Diseases target for 24 Nature damage every 1.0 sec for 7s. Deals 2 times as much damage to Aquatic enemies.
	// https://www.wowhead.com/forever/spell=1295744
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 274963, ItemName: "Rot-Covered Harpoon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blunderflame the target for 56 Fire damage.
	// https://www.wowhead.com/forever/spell=1296730
	// unsupported: states no rate
	// trigger 1296730 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1296730, BuffSpellID: 1296730, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 275439, ItemName: "Blunderflame Bone Bow"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Cremate the target for 56 Fire damage.
	// https://www.wowhead.com/forever/spell=1297921
	// unsupported: states no rate
	// trigger 1297921 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1297921, BuffSpellID: 1297921, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 275441, ItemName: "Cremation Longbow"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance on harmful spell cast to reduce target's attack power by 60 for 30s.
	//
	// https://www.wowhead.com/forever/spell=1297082
	// unsupported: states no rate
	// trigger 1297085 (no stated rate, core.CallbackOnCastComplete, core.ProcMaskSpellDamage) -> buff 1297082
	// shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 1297085, BuffSpellID: 1297082},
	//	[]shared.ItemVariant{
	//	{ItemID: 275630, ItemName: "Depleted Eye of Influence"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Heals the wielder for 182 damage over 14s.
	// https://www.wowhead.com/forever/spell=1297357
	// unsupported: states no rate
	// trigger 1297357 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataHealProc(shared.SpellDataProc{TriggerSpellID: 1297357, BuffSpellID: 1297357, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 275645, ItemName: "Reforged Spear"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your healing spells have a 5% chance to remove 1 Poison effect from the target.
	// https://www.wowhead.com/forever/spell=1297369
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 275647, ItemName: "Cleansed Ritual Kris"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck in combat, inflicts 3 Nature damage to the attacker.
	// https://www.wowhead.com/forever/spell=1297378
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 275649, ItemName: "Bramblebark Barrier"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck in combat, inflicts 3 Nature damage to the attacker.
	// https://www.wowhead.com/forever/spell=1297910
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 275833, ItemName: "Bristlecone Cloak"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=1316049
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 276895, ItemName: "Transformative Cocoon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=1317283
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 277910, ItemName: "Mysticbloom Seed"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your Well Fed buffs metabolize 50% slower. This effect is increased by an additional 50% in Wasteland
	// areas.
	// https://www.wowhead.com/forever/spell=1302968
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 277950, ItemName: "Wastewanderer Rations"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the critical strike chance of your Lesser Healing Wave spell by 4%.
	// https://www.wowhead.com/forever/spell=1306448
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 279249, ItemName: "Totem of Urgency"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the cooldown of your Swiftmend spell by 3 sec.
	// https://www.wowhead.com/forever/spell=1306488
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 279250, ItemName: "Idol of Swiftness"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your Lacerate hits have a 10% chance to reset the cooldown on Primal Bite.
	// https://www.wowhead.com/forever/spell=1306483
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 279251, ItemName: "Idol of the Ursine Twins"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts a target for 42 Fire damage.
	// https://www.wowhead.com/forever/spell=1322295
	// unsupported: states no rate
	// trigger 1322295 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1322295, BuffSpellID: 1322295, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 279895, ItemName: "Ironforge Greathammer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage dealt by your Stormstrike ability by 10%.
	// https://www.wowhead.com/forever/spell=1309422
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 280604, ItemName: "Rage of the Storm"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the critical strike chance of your Healing Wave and Lesser Healing Wave spells by 2%.
	// https://www.wowhead.com/forever/spell=1309539
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 280605, ItemName: "Tidebringer's Claw"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the critical strike chance of your Lightning Bolt spell by 1%.
	// https://www.wowhead.com/forever/spell=1309442
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 280607, ItemName: "Lightning's Grasp"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blast the target with the fires of a forge, inflicting 51 Fire damage.
	// https://www.wowhead.com/forever/spell=1312176
	// unsupported: states no rate
	// trigger 1312176 (every time, core.CallbackOnSpellHitDealt, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeMHSpecial)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1312176, BuffSpellID: 1312176, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 281665, ItemName: "Twilight Forgehammer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck in combat, inflicts 2 Fire damage to the attacker.
	// https://www.wowhead.com/forever/spell=1312955
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 282080, ItemName: "Flame Seared Signet"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Smash the target for 85 damage.
	// https://www.wowhead.com/forever/spell=1314011
	// unsupported: states no rate
	// trigger 1314011 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1314011, BuffSpellID: 1314011, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 282551, ItemName: "Golem Fist"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Strike the target with a Flaming Shell, inflicting 7 Fire damage every 3 sec for 9 sec.
	// https://www.wowhead.com/forever/spell=1314040
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 282561, ItemName: "Siegebreaker's Blaster"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Puncture the target, inflicting 21 Physical damage.
	// https://www.wowhead.com/forever/spell=1314305
	// unsupported: states no rate
	// trigger 1314305 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1314305, BuffSpellID: 1314305, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 282703, ItemName: "Needletooth's Needletooth"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases intellect by 5 at night.
	// https://www.wowhead.com/forever/spell=1315339
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 284193, ItemName: "Moongazer's Wand"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the target enemy's attack power by 60 for 30s.
	// https://www.wowhead.com/forever/spell=1315767
	// unsupported: states no rate
	// trigger 1315767 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 1315767, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 284262, ItemName: "Howling Hide"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sear your target's flesh, dealing 84 Fire damage.
	// https://www.wowhead.com/forever/spell=1316039
	// unsupported: states no rate
	// trigger 1316039 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1316039, BuffSpellID: 1316039, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 284386, ItemName: "Whipfang's Skinsearer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Burns the target for 33 Fire damage every 3.0 sec for 15s.
	// https://www.wowhead.com/forever/spell=1316865
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 284667, ItemName: "Flame Seared Sword"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Underwater breath lasts 50% longer than normal.
	// https://www.wowhead.com/forever/spell=1316928
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 284699, ItemName: "Still Water Band"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases movement speed by 3% in Dustwallow Marsh and Swamp of Sorrows.
	// https://www.wowhead.com/forever/spell=1317486
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 285089, ItemName: "Swamp Roamer's Band"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your attacks have a high chance to poison the target for 8 Nature damage every 3.0 sec for 9s.
	// https://www.wowhead.com/forever/spell=1317583
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 285102, ItemName: "Darkmist Dirk"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Stun the target for 1s.
	// https://www.wowhead.com/forever/spell=1318178
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 285281, ItemName: "Arcanite Blacksmith Hammer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Puncture the target, inflicting 90 Physical damage.
	// https://www.wowhead.com/forever/spell=1318250
	// unsupported: states no rate
	// trigger 1318250 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1318250, BuffSpellID: 1318250, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 285332, ItemName: "Puncturing Spear"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Thunderstomp nearby enemies, inflicting 23 Nature damage.
	// https://www.wowhead.com/forever/spell=1320498
	// unsupported: states no rate
	// trigger 1320498 (0%, core.CallbackEmpty, core.ProcMaskUnknown)
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1320498, BuffSpellID: 1320498, IsWeaponProc: true},
	//	[]shared.ItemVariant{
	//	{ItemID: 286537, ItemName: "Thunderstomp's Horn"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deals 22 Nature damage every time you block. Damage is doubled in swamps.
	// https://www.wowhead.com/forever/spell=1324685
	// unsupported: an outcome the proc mask has no bit for
	// trigger 1324686 (every time, core.CallbackOnSpellHitTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial) -> buff 1324685
	// shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1324686, BuffSpellID: 1324685},
	//	[]shared.ItemVariant{
	//	{ItemID: 286982, ItemName: "Unearthed Barrier"},
	// })

	// When struck in combat has a 1% chance of inflicting 100 Shadow damage to the attacker.
	// https://www.wowhead.com/forever/spell=16783
	// trigger 7617 (1%, core.CallbackOnSpellHitTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial) -> buff 16783
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 7617, BuffSpellID: 16783},
		[]shared.ItemVariant{
			{ItemID: 1131, ItemName: "Totem of Infliction"},
		})

	// When struck in melee combat, has a 3% chance of stealing 270 life from target enemy.
	// https://www.wowhead.com/forever/spell=18817
	// trigger 18815 (effect 1's chance, core.CallbackOnSpellHitTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial) -> buff 18817
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 18815, BuffSpellID: 18817},
		[]shared.ItemVariant{
			{ItemID: 1168, ItemName: "Skullflame Shield -  - "},
		})

	// When struck in combat has a 1% chance of raising a thorny shield that inflicts 5 Nature damage to attackers
	// when hit and increases Nature resistance by 50 for 30s.
	//
	// https://www.wowhead.com/forever/spell=17154
	// trigger 18097 (1%, core.CallbackOnSpellHitTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial) -> buff 17154
	shared.NewSpellDataProc(shared.SpellDataProc{TriggerSpellID: 18097, BuffSpellID: 17154},
		[]shared.ItemVariant{
			{ItemID: 1204, ItemName: "The Green Tower - NatureResistance - "},
		})

	// When struck in combat has a 3% chance to encase the caster in bone, increasing armor by 150 for 20s.
	// https://www.wowhead.com/forever/spell=18828
	// trigger 19409 (3%, core.CallbackOnSpellHitTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial) -> buff 18828
	shared.NewSpellDataProc(shared.SpellDataProc{TriggerSpellID: 19409, BuffSpellID: 18828},
		[]shared.ItemVariant{
			{ItemID: 1979, ItemName: "Wall of the Dead"},
		})

	// Hurricane (2824) left the client in the 17 October 2026 hotfix cache (upstream #817), so its Frost
	// Arrow proc (trigger 1316319 -> 29502) is no longer registered. MythicSim patch 164.

	// 5% chance of dealing 40 Fire damage on a successful melee attack.
	// https://www.wowhead.com/forever/spell=9057
	// trigger 9233 (10%, core.CallbackOnSpellHitDealt, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial) -> buff 9057
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 9233, BuffSpellID: 9057},
		[]shared.ItemVariant{
			{ItemID: 7284, ItemName: "Red Whelp Gloves"},
		})

	// When struck in combat has a 5% chance of inflicting 140 Shadow damage to the attacker.
	// https://www.wowhead.com/forever/spell=1293421
	// trigger 7619 (5%, core.CallbackOnSpellHitTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial) -> buff 1293421
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 7619, BuffSpellID: 1293421},
		[]shared.ItemVariant{
			{ItemID: 7747, ItemName: "Vile Protector"},
		})

	// When struck in combat has a 3% chance to heal you for 80.
	// https://www.wowhead.com/forever/spell=9777
	// trigger 9778 (3%, core.CallbackOnSpellHitTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial) -> buff 9777
	shared.NewSpellDataHealProc(shared.SpellDataProc{TriggerSpellID: 9778, BuffSpellID: 9777},
		[]shared.ItemVariant{
			{ItemID: 7939, ItemName: "Truesilver Breastplate"},
		})

	// When struck in combat has a 5% chance of inflicting 93 Nature damage to the attacker.
	// https://www.wowhead.com/forever/spell=1292879
	// trigger 1292880 (5%, core.CallbackOnSpellHitTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial) -> buff 1292879
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1292880, BuffSpellID: 1292879},
		[]shared.ItemVariant{
			{ItemID: 9458, ItemName: "Thermaplugg's Central Core"},
		})

	// Spells and attacks against Swine deal 2 Nature damage.
	// https://www.wowhead.com/forever/spell=1293782
	// trigger 1293783 (every time, core.CallbackOnSpellHitDealt, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage) -> buff 1293782
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1293783, BuffSpellID: 1293782},
		[]shared.ItemVariant{
			{ItemID: 10760, ItemName: "Swine Fists"},
		})

	// 2% chance when struck in melee to gain a holy shield, absorbing 216 damage for 15s. This chance is doubled
	// in Wasteland and Haunted areas.
	// https://www.wowhead.com/forever/spell=10368
	// trigger 8397 (4%, the column over effect 1's 2%, core.CallbackOnSpellHitTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial) -> buff 10368
	shared.NewSpellDataAbsorbProc(shared.SpellDataProc{TriggerSpellID: 8397, BuffSpellID: 10368},
		[]shared.ItemVariant{
			{ItemID: 11302, ItemName: "Uther's Strength"},
		})

	// When struck in combat has a 1% chance of inflicting 100 Shadow damage to the attacker.
	// https://www.wowhead.com/forever/spell=16783
	// trigger 7617 (1%, core.CallbackOnSpellHitTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial) -> buff 16783
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 7617, BuffSpellID: 16783},
		[]shared.ItemVariant{
			{ItemID: 11861, ItemName: "Girdle of Reprisal"},
		})

	// Adds 4 Fire damage to your weapon attack.
	// https://www.wowhead.com/forever/spell=7714
	// trigger 7721 (every time, core.CallbackOnSpellHitDealt, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial) -> buff 7714
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 7721, BuffSpellID: 7714},
		[]shared.ItemVariant{
			{ItemID: 12631, ItemName: "Fiery Plate Gauntlets"},
		})

	// Adds 3 Lightning damage to your melee attacks.
	// https://www.wowhead.com/forever/spell=16614
	// trigger 16615 (every time, core.CallbackOnSpellHitDealt, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial) -> buff 16614
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 16615, BuffSpellID: 16614},
		[]shared.ItemVariant{
			{ItemID: 12632, ItemName: "Storm Gauntlets"},
		})

	// Reduces an enemy's armor by 165. Stacks up to 3 times.
	// https://www.wowhead.com/forever/spell=16928
	// trigger 16928 (1 ppm, core.CallbackEmpty, core.ProcMaskUnknown)
	shared.NewSpellDataDebuffProc(shared.SpellDataProc{TriggerSpellID: 16928, IsWeaponProc: true},
		[]shared.ItemVariant{
			{ItemID: 12798, ItemName: "Annihilator"},
		})

	// Grants a chance on striking the enemy for 60 Fire damage.
	// https://www.wowhead.com/forever/spell=13441
	// trigger 16982 (1%, core.CallbackOnSpellHitDealt, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial) -> buff 13441
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 16982, BuffSpellID: 13441},
		[]shared.ItemVariant{
			{ItemID: 12805, ItemName: "Orb of Fire"},
		})

	// Adds 2 fire damage to your melee attacks.
	// https://www.wowhead.com/forever/spell=7712
	// trigger 7711 (every time, core.CallbackOnSpellHitDealt, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial) -> buff 7712
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 7711, BuffSpellID: 7712},
		[]shared.ItemVariant{
			{ItemID: 17111, ItemName: "Blazefury Medallion"},
		})

	// Adds 2 Arcane damage to your melee attacks.
	// https://www.wowhead.com/forever/spell=1302247
	// trigger 1302248 (every time, core.CallbackOnSpellHitDealt, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial) -> buff 1302247
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1302248, BuffSpellID: 1302247},
		[]shared.ItemVariant{
			{ItemID: 18383, ItemName: "Force Imbued Gauntlets"},
		})

	// When struck in combat has a 5% chance of inflicting 50 Nature damage to the attacker.
	// https://www.wowhead.com/forever/spell=16782
	// trigger 13959 (5%, core.CallbackOnSpellHitTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial) -> buff 16782
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 13959, BuffSpellID: 16782},
		[]shared.ItemVariant{
			{ItemID: 18825, ItemName: "Grand Marshal's Aegis - "},
		})

	// When struck in combat has a 5% chance of inflicting 50 Nature damage to the attacker.
	// https://www.wowhead.com/forever/spell=16782
	// trigger 13959 (5%, core.CallbackOnSpellHitTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial) -> buff 16782
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 13959, BuffSpellID: 16782},
		[]shared.ItemVariant{
			{ItemID: 18826, ItemName: "High Warlord's Shield Wall - "},
		})

	// 2% chance on successful spellcast to increase your Spirit by 150 for 15s.
	// https://www.wowhead.com/forever/spell=23684
	// trigger 23688 (2%, core.CallbackOnCastComplete, core.ProcMaskSpellDamage | core.ProcMaskSpellHealing) -> buff 23684
	shared.NewSpellDataProc(shared.SpellDataProc{TriggerSpellID: 23688, BuffSpellID: 23684},
		[]shared.ItemVariant{
			{ItemID: 19288, ItemName: "Darkmoon Card: Blue Dragon"},
		})

	// Gives a chance when your harmful spells land to increase the damage of your spells and effects by 132
	// for 10s.
	// https://www.wowhead.com/forever/spell=25907
	// trigger 25906 (5%, core.CallbackOnCastComplete, core.ProcMaskSpellDamage) -> buff 25907
	shared.NewSpellDataProc(shared.SpellDataProc{TriggerSpellID: 25906, BuffSpellID: 25907},
		[]shared.ItemVariant{
			{ItemID: 21190, ItemName: "Wrath of Cenarius"},
		})

	// Chance on harmful spellcast to increase your spell damage and healing by up to 35 for 10s. This spell
	// damage increase is doubled against Dragonkin.
	// https://www.wowhead.com/forever/spell=1318930
	// trigger 1318931 (every time, core.CallbackOnCastComplete, core.ProcMaskSpellDamage) -> buff 1318930
	shared.NewSpellDataProc(shared.SpellDataProc{TriggerSpellID: 1318931, BuffSpellID: 1318930},
		[]shared.ItemVariant{
			{ItemID: 22268, ItemName: "Draconic Infused Emblem"},
		})

	// When struck in combat has a 20% chance of inflicting 50 Nature damage to the attacker.
	// https://www.wowhead.com/forever/spell=16782
	// trigger 1216968 (20%, core.CallbackOnSpellHitTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial) -> buff 16782
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1216968, BuffSpellID: 16782},
		[]shared.ItemVariant{
			{ItemID: 234562, ItemName: "High Warlord's Shield Wall -  - "},
		})

	// When struck in combat has a 20% chance of inflicting 50 Nature damage to the attacker.
	// https://www.wowhead.com/forever/spell=16782
	// trigger 1216968 (20%, core.CallbackOnSpellHitTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial) -> buff 16782
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1216968, BuffSpellID: 16782},
		[]shared.ItemVariant{
			{ItemID: 234588, ItemName: "Grand Marshal's Aegis -  - "},
		})

	// When struck in combat has a 5% chance of inflicting 50 Nature damage to the attacker.
	// https://www.wowhead.com/forever/spell=16782
	// trigger 13959 (5%, core.CallbackOnSpellHitTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial) -> buff 16782
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 13959, BuffSpellID: 16782},
		[]shared.ItemVariant{
			{ItemID: 272591, ItemName: "Premier High Warlord's Shield Wall"},
		})

	// When struck in combat has a 5% chance of inflicting 50 Nature damage to the attacker.
	// https://www.wowhead.com/forever/spell=16782
	// trigger 13959 (5%, core.CallbackOnSpellHitTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial) -> buff 16782
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 13959, BuffSpellID: 16782},
		[]shared.ItemVariant{
			{ItemID: 272838, ItemName: "Premier Grand Marshal's Aegis"},
		})

	// Harmful spell casts sear the target for 7 Fire damage.
	// https://www.wowhead.com/forever/spell=1291570
	// trigger 1291568 (every time, core.CallbackOnCastComplete, core.ProcMaskSpellDamage) -> buff 1291570
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1291568, BuffSpellID: 1291570},
		[]shared.ItemVariant{
			{ItemID: 273003, ItemName: "Searing Dagger"},
		})

	// When damaged, has a chance to grant you 1080 armor for 12s.
	// https://www.wowhead.com/forever/spell=1293345
	// trigger 1293342 (33%, core.CallbackOnSpellHitTaken | core.CallbackOnPeriodicDamageTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage) -> buff 1293345
	shared.NewSpellDataProc(shared.SpellDataProc{TriggerSpellID: 1293342, BuffSpellID: 1293345},
		[]shared.ItemVariant{
			{ItemID: 273046, ItemName: "Guardian's Dualblade"},
		})

	// Spells and attacks deal 2 Shadow damage. Deals 4 times as much damage to Murlocs.
	// https://www.wowhead.com/forever/spell=1292675
	// trigger 1292674 (every time, core.CallbackOnSpellHitDealt, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage) -> buff 1292675
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1292674, BuffSpellID: 1292675},
		[]shared.ItemVariant{
			{ItemID: 273840, ItemName: "Cursed Murloc Eye"},
		})

	// Your melee weapon attacks deal 2 Nature damage to both you and your target.
	//
	// https://www.wowhead.com/forever/spell=1293333
	// trigger 1293331 (every time, core.CallbackOnSpellHitDealt, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial) -> buff 1293333
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1293331, BuffSpellID: 1293333},
		[]shared.ItemVariant{
			{ItemID: 274159, ItemName: "Thorncursed Grips"},
		})

	// 17% chance when struck in combat to gain 12 mana per 5 sec for 10s.
	// https://www.wowhead.com/forever/spell=1293700
	// trigger 1293701 (17%, core.CallbackOnSpellHitTaken, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage) -> buff 1293700
	shared.NewSpellDataProc(shared.SpellDataProc{TriggerSpellID: 1293701, BuffSpellID: 1293700},
		[]shared.ItemVariant{
			{ItemID: 274290, ItemName: "Painwalker Buckler"},
		})

	// Melee attacks deal 98 additional Fire damage against Frozen targets.
	// https://www.wowhead.com/forever/spell=1300130
	// trigger 1300128 (every time, core.CallbackOnSpellHitDealt, core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto) -> buff 1300130
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1300128, BuffSpellID: 1300130},
		[]shared.ItemVariant{
			{ItemID: 276631, ItemName: "Coldflame Saber"},
		})

	// Thrown attacks explode on impact, causing 6 Fire damage to nearby enemies.
	// https://www.wowhead.com/forever/spell=1318031
	// trigger 1318034 (every time, core.CallbackOnSpellHitDealt, core.ProcMaskRangedAuto) -> buff 1318031
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1318034, BuffSpellID: 1318031},
		[]shared.ItemVariant{
			{ItemID: 285275, ItemName: "Satchel of Copper Bombs"},
		})

	// Thrown attacks explode on impact, causing 11 Fire damage to nearby enemies.
	// https://www.wowhead.com/forever/spell=1318061
	// trigger 1318062 (every time, core.CallbackOnSpellHitDealt, core.ProcMaskRangedAuto) -> buff 1318061
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1318062, BuffSpellID: 1318061},
		[]shared.ItemVariant{
			{ItemID: 285276, ItemName: "Satchel of Bronze Bombs"},
		})

	// Thrown attacks explode on impact, causing 17 Fire damage to nearby enemies.
	// https://www.wowhead.com/forever/spell=1318068
	// trigger 1318069 (every time, core.CallbackOnSpellHitDealt, core.ProcMaskRangedAuto) -> buff 1318068
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1318069, BuffSpellID: 1318068},
		[]shared.ItemVariant{
			{ItemID: 285277, ItemName: "Satchel of Iron Bombs"},
		})

	// Thrown attacks explode on impact, causing 21 Fire damage to nearby enemies.
	// https://www.wowhead.com/forever/spell=1318121
	// trigger 1318123 (every time, core.CallbackOnSpellHitDealt, core.ProcMaskRangedAuto) -> buff 1318121
	shared.NewSpellDataDamageProc(shared.SpellDataProc{TriggerSpellID: 1318123, BuffSpellID: 1318121},
		[]shared.ItemVariant{
			{ItemID: 285278, ItemName: "Satchel of Dark Iron Bombs"},
		})

	// Skipped
	// Not simulated: Rod of the Sleepwalker: "Resist Sleep 05" (1292652) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=1292652
	// Not simulated: Mug O' Hurt: "Dazed" (13496) - ignored aura type 33
	// https://www.wowhead.com/forever/spell=13496
	// Not simulated: Girdle of the Blindwatcher: "Stealth Detection 05" (1292149) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=1292149
	// Not simulated: Medal of Courage: "Resist Fear 02" (1324578) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=1324578
	// Not simulated: Mask of Thero-shan: "Stealth 06" (17746) - ignored aura type 154
	// https://www.wowhead.com/forever/spell=17746
	// Not simulated: Nightscape Boots: "Stealth 06" (17746) - ignored aura type 154
	// https://www.wowhead.com/forever/spell=17746
	// Not simulated: Catseye Ultra Goggles: "Stealth Detection 15" (12418) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=12418
	// Not simulated: Stealthblade: "Stealth 08" (27037) - ignored aura type 154
	// https://www.wowhead.com/forever/spell=27037
	// Not simulated: Tooth of Eranikus: "Resist Sleep 03" (1300001) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=1300001
	// Not simulated: Dragon's Call: "Dragon's Call" (13049) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=13049
	// Not simulated: Carrot on a Stick: "Mount Speed" (13587) - ignored aura type 130
	// https://www.wowhead.com/forever/spell=13587
	// Not simulated: Savage Gladiator Chain: "Resist Fear 04" (1292574) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=1292574
	// Not simulated: High Priestess Boots: "Resist Charm 10" (1300925) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=1300925
	// Not simulated: Blackblade of Shahram: "Shahram" (16602) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=16602
	// Not simulated: Stronghold Gauntlets: "Immune to Disarm" (7219) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=7219
	// Not simulated: Tome of Knowledge: "Stealth 05" (1293068) - ignored aura type 154
	// https://www.wowhead.com/forever/spell=1293068
	// Not simulated: Spectral Essence: "Visions of the Past" (17623) - ignored aura type 19
	// https://www.wowhead.com/forever/spell=17623
	// Not simulated: Voice Amplification Modulator: "Resist Silence 07" (19786) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=19786
	// Not simulated: Knight-Lieutenant's Dragonhide Gloves: "Stealth Detection 10" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Marshal's Dragonhide Gauntlets: "Stealth Detection 10" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Blood Guard's Dragonhide Gauntlets: "Stealth Detection 10" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: General's Dragonhide Gloves: "Stealth Detection 10" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Shard of the Defiler: "Echo of Archimonde" (21079) - ignored aura type 56
	// https://www.wowhead.com/forever/spell=21079
	// Not simulated: Grovekeeper's Drape: "Resist Charm 05" (1294786) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=1294786
	// Not simulated: Mark of Resolution: "Stout Heart" (21958) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=21958
	// Not simulated: Zum'rah's Vexing Cane: "Ward of Zum'rah" (1294258) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=1294258
	// Not simulated: Helm of Awareness: "Stealth Detection 05" (1292149) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=1292149
	// Not simulated: Murmuring Ring: "Resist Silence 07" (19786) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=19786
	// Not simulated: Kreeg's Mug: "Resist Fear 10" (1302346) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=1302346
	// Not simulated: The Eye of Divinity: "Eye of Divinity" (23101) - ignored aura type 19
	// https://www.wowhead.com/forever/spell=23101
	// Not simulated: Death Grips: "Immune to Disarm" (7219) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=7219
	// Not simulated: Pirate's Eye Patch: "Fear Resistance 4" (24351) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=24351
	// Not simulated: Bloodvine Lens: "Stealth Detection 10" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Whisperwalk Boots: "Stealth 06" (17746) - ignored aura type 154
	// https://www.wowhead.com/forever/spell=17746
	// Not simulated: Gnomish Turban of Psychic Might: "Resist Silence 10" (26208) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=26208
	// Not simulated: Darkmantle Boots: "Stealth 08" (27037) - ignored aura type 154
	// https://www.wowhead.com/forever/spell=27037
	// Not simulated: Bloodclot Band: "Resist Bleed 10" (1300755) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=1300755
	// Not simulated: Blade of Necromancy: "Necromancy" (1302109) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=1302109
	// Not simulated: Blood Guard's Dragonhide Grips: "Stealth Detection 10" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Knight-Lieutenant's Dragonhide Grips: "Stealth Detection 10" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Sergeant Major's Leather Boots: "Stealth 06" (17746) - ignored aura type 154
	// https://www.wowhead.com/forever/spell=17746
	// Not simulated: First Sergeant's Leather Boots: "Stealth 06" (17746) - ignored aura type 154
	// https://www.wowhead.com/forever/spell=17746
	// Not simulated: Darkmantle Cap: "Stealth 08" (27037) - ignored aura type 154
	// https://www.wowhead.com/forever/spell=27037
	// Not simulated: Darkmantle Footpads: "Stealth 05" (1293068) - ignored aura type 154
	// https://www.wowhead.com/forever/spell=1293068
	// Not simulated: Premier Kodohide Gauntlets: "Stealth Detection 10" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Premier Champion's Dragonhide Gloves: "Stealth Detection 10" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Premier Dragonhide Grips: "Stealth Detection 10" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Premier Beasthide Gloves: "Stealth Detection 10" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Premier Lieutenant Commander's Beasthide Gauntlets: "Stealth Detection 10" (23217) - ignored aura type
	// Not simulated: 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Cloak of Hermitic Bliss: "-25% Movement Speed" (1291748) - ignored aura type 33
	// https://www.wowhead.com/forever/spell=1291748
	// Not simulated: Premier Lunarhide Gloves: "Stealth Detection 10" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Premier Dreamhide Gloves: "Stealth Detection 10" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Premier Wyrmhide Gauntlets: "Stealth Detection 10" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Premier Lieutenant Commander's Lunarhide Gauntlets: "Stealth Detection 10" (23217) - ignored aura type
	// Not simulated: 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Premier Lieutenant Commander's Dreamhide Gauntlets: "Stealth Detection 10" (23217) - ignored aura type
	// Not simulated: 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Premier Champion's Wyrmhide Gloves: "Stealth Detection 10" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Premier Champion's Kodohide Gloves: "Stealth Detection 10" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Sorcerer Collar: "Resist Silence 04" (1292222) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=1292222
	// Not simulated: Blindwatcher's Sight: "Resist Disorient 07" (1292268) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=1292268
	// Not simulated: Nightskulker Ring: "Resist Fear 04" (1292574) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=1292574
	// Not simulated: Skullduggery Belt: "Stealth 05" (1293068) - ignored aura type 154
	// https://www.wowhead.com/forever/spell=1293068
	// Not simulated: Field Agent Beverage: "Resist Charm 06" (1297455) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=1297455
	// Not simulated: Wintersaber Hide Lined Gloves: "Mount Speed" (1315778) - ignored aura type 130
	// https://www.wowhead.com/forever/spell=1315778
	// Not simulated: Mithril Blacksmith Hammer: "Concussed" (1318163) - ignored aura type 33
	// https://www.wowhead.com/forever/spell=1318163
}
