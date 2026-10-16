package buffmanifest

import "github.com/wowsims/forever/sim/core/proto"

// One row per proto field of RaidBuffs, PartyBuffs, IndividualBuffs and Debuffs, in the order the
// settings tab lists them, with the fields that have no input after the last one that does. A
// field's proto number is its position in its slice.

var Raid = []BuffSpec{
	{
		Field:    "arcane_brilliance",
		SpellID:  23028,
		Category: "StatBuff",
		Stats: []proto.Stat{
			proto.Stat_StatIntellect,
		},
	},
	{
		Field:    "prayer_of_spirit",
		SpellID:  27681,
		Category: "StatBuff",
		Stats: []proto.Stat{
			proto.Stat_StatSpirit,
			proto.Stat_StatSpellDamage,
		},
	},
	{
		Field:   "gift_of_the_wild",
		SpellID: 21850,
		Stats: []proto.Stat{
			proto.Stat_StatArmor,
			proto.Stat_StatStrength,
			proto.Stat_StatAgility,
			proto.Stat_StatIntellect,
			proto.Stat_StatSpirit,
			proto.Stat_StatStamina,
		},
	},
	{
		Field:      "thorns",
		SpellID:    9910,
		Category:   "Thorns",
		SingleAura: true,
		Stats: []proto.Stat{
			proto.Stat_StatDefenseRating,
			proto.Stat_StatStamina,
		},
	},
	{
		Field:   "prayer_of_fortitude",
		SpellID: 21564,
		Stats: []proto.Stat{
			proto.Stat_StatStamina,
		},
	},
	{
		Field:   "prayer_of_shadow_protection",
		SpellID: 27683,
		Stats: []proto.Stat{
			proto.Stat_StatShadowResistance,
			proto.Stat_StatStamina,
		},
	},
	{
		Field:          "fire_resistance_aura",
		SpellID:        19900,
		Category:       "FireResistanceAura",
		SharedCategory: "PaladinAura",
		SingleAura:     true,
		Stats: []proto.Stat{
			proto.Stat_StatFireResistance,
		},
	},
	{
		Field:          "frost_resistance_aura",
		SpellID:        19898,
		Category:       "FrostResistanceAura",
		SharedCategory: "PaladinAura",
		SingleAura:     true,
		Stats: []proto.Stat{
			proto.Stat_StatFrostResistance,
		},
	},
	{
		Field:          "shadow_resistance_aura",
		SpellID:        19896,
		Category:       "ShadowResistanceAura",
		SharedCategory: "PaladinAura",
		SingleAura:     true,
		Stats: []proto.Stat{
			proto.Stat_StatShadowResistance,
		},
	},
	{
		Field:   "fire_resistance_totem",
		SpellID: 10535,
		CastID:  10538,
		Stats: []proto.Stat{
			proto.Stat_StatFireResistance,
		},
	},
	{
		Field:   "frost_resistance_totem",
		SpellID: 10477,
		CastID:  10479,
		Stats: []proto.Stat{
			proto.Stat_StatFrostResistance,
		},
	},
	{
		Field:   "nature_resistance_totem",
		SpellID: 10599,
		CastID:  10601,
		Stats: []proto.Stat{
			proto.Stat_StatNatureResistance,
		},
	},
	{
		Field:   "aspect_of_the_wild",
		SpellID: 20190,
		Stats: []proto.Stat{
			proto.Stat_StatNatureResistance,
		},
	},
}

var Party = []BuffSpec{
	{
		Field:   "blood_pact",
		SpellID: 11767,
		Stats: []proto.Stat{
			proto.Stat_StatStamina,
		},
	},
	// The improved state is a warrior wearing three pieces of Battlegear of Wrath, whose set spell
	// 23563 adds 30 to every effect of the shout. No set spell is resolved, so driveBattleShout adds it.
	{
		Field:      "battle_shout",
		SpellID:    25289,
		Category:   "BattleShout",
		SingleAura: true,
		Driver:     true,
		ImpAction:  &ActionRef{SpellID: 23563},
		Stats: []proto.Stat{
			proto.Stat_StatAttackPower,
		},
	},
	{
		Field:          "devotion_aura",
		SpellID:        10293,
		Category:       "DevotionAura",
		SharedCategory: "PaladinAura",
		SingleAura:     true,
		Stats: []proto.Stat{
			proto.Stat_StatArmor,
		},
	},
	// The client row is A_MOD_CRIT_PCT, every kind of crit, and 17007 calls it exclusive with Moonkin
	// Aura. The 16 October 2026 beta notes keep it to melee and ranged crit and let the two stack, so
	// it has a category of its own and is driven: the driver reads the row's crit as physical crit.
	{
		Field:      "leader_of_the_pack",
		SpellID:    24932,
		CastID:     17007,
		Category:   "LeaderOfThePack",
		SingleAura: true,
		Driver:     true,
		Stats: []proto.Stat{
			proto.Stat_StatMeleeCritRating,
		},
	},
	{
		Field:    "mana_spring_totem",
		SpellID:  10494,
		CastID:   10497,
		Talent:   16187,
		Category: "ManaSpringTotem",
		Stats: []proto.Stat{
			proto.Stat_StatMP5,
		},
	},
	// The cast states the totem's 13 s life and the 5 min cooldown, the aura the mana.
	{
		Field:    "mana_tide_totems",
		SpellID:  17360,
		CastID:   17359,
		Kind:     KindExternalCD,
		Category: "ManaTideTotem",
		Stats: []proto.Stat{
			proto.Stat_StatMP5,
		},
		Label: "Mana Tide Totem",
	},
	// See leader_of_the_pack: since the 16 October 2026 beta notes Moonkin Aura is spell crit alone,
	// in a category of its own, and driven so the row's crit is read as spell crit.
	{
		Field:      "moonkin_aura",
		SpellID:    24907,
		Category:   "MoonkinAura",
		SingleAura: true,
		Driver:     true,
		Stats: []proto.Stat{
			proto.Stat_StatSpellCritRating,
		},
	},
	// Driven because the damage scales with the providing paladin's Holy spell power.
	{
		Field:          "retribution_aura",
		SpellID:        10301,
		Category:       "RetributionAura",
		SharedCategory: "PaladinAura",
		SingleAura:     true,
		Driver:         true,
		Stats: []proto.Stat{
			proto.Stat_StatArmor,
			proto.Stat_StatDefenseRating,
		},
	},
	{
		Field:  "retribution_aura_spell_power",
		Kind:   KindFlag,
		Proto:  ProtoDouble,
		Reason: "retired: Retribution Aura reads the spell power of the unit it is on, not the paladin's (beta logs 2729-2736), so nothing reads this; the row stays so the party fields after it keep their proto numbers.",
	},
	{
		Field:          "concentration_aura",
		SpellID:        19746,
		Category:       "ConcentrationAura",
		SharedCategory: "PaladinAura",
		SingleAura:     true,
		Stats: []proto.Stat{
			proto.Stat_StatDefenseRating,
		},
	},
	// Rank 4 (20905, level 50) states 75 ranged attack power and rank 5 (20906, level 60) only 50,
	// where Era's rank 5 was 100. Lower ranks stay castable, so a level 60 hunter keeps rank 4 up:
	// 75 is the floor whatever rank 5 turns out to be on the server.
	{
		Field:     "trueshot_aura",
		SpellID:   20905,
		LowerRank: true,
		Stats: []proto.Stat{
			proto.Stat_StatRangedAttackPower,
		},
	},
	{
		Field:   "atiesh_mage",
		SpellID: 28142,
		Owner:   proto.Class_ClassMage,
		Kind:    KindItemCount,
		Label:   "Atiesh - Mage",
		Stats: []proto.Stat{
			proto.Stat_StatSpellDamage,
			proto.Stat_StatHealingPower,
		},
	},
	{
		Field:   "atiesh_warlock",
		SpellID: 28143,
		Owner:   proto.Class_ClassWarlock,
		Kind:    KindItemCount,
		Label:   "Atiesh - Warlock",
		Stats: []proto.Stat{
			proto.Stat_StatSpellDamage,
			proto.Stat_StatHealingPower,
		},
	},
	{
		Field:    "strength_of_earth_totem",
		SpellID:  25362,
		CastID:   25361,
		Category: "StrengthOfEarthTotem",
		Stats: []proto.Stat{
			proto.Stat_StatStrength,
		},
	},
	// Driven because totem_twisting shortens the aura to the 9 s of every 10 a twisting shaman keeps it up.
	{
		Field:    "grace_of_air_totem",
		SpellID:  25360,
		CastID:   25359,
		Category: "GraceOfAirTotem",
		Driver:   true,
		Stats: []proto.Stat{
			proto.Stat_StatAgility,
		},
	},
	{
		Field:    "windfury_totem",
		SpellID:  10610,
		CastID:   10614,
		Kind:     KindProc,
		Category: "WindfuryTotem",
		Stats: []proto.Stat{
			proto.Stat_StatAttackPower,
		},
	},
	{
		Field:   "atiesh_druid",
		SpellID: 28145,
		Owner:   proto.Class_ClassDruid,
		Kind:    KindItemCount,
		Label:   "Atiesh - Druid",
	},
	{
		Field:   "atiesh_priest",
		SpellID: 28144,
		Owner:   proto.Class_ClassPriest,
		Kind:    KindItemCount,
		Label:   "Atiesh - Priest",
	},
	{
		Field:  "totem_twisting",
		Kind:   KindFlag,
		Reason: "sim behaviour toggle with no spell source; rendered under Other Inputs.",
	},
	// Another shaman's Flametongue Totem (patch 70). The aura is the totem's party aura (15036): 100% on a
	// landed melee auto attack, triggering 16389's fire hit. driveFlametongueTotem adds the hit on main-hand
	// autos, and a main-hand Flametongue Weapon turns it off (the category). A fire totem, so unlike the air
	// totems it joins no slot. No spec's default turns it on.
	{
		Field:    "flametongue_totem",
		SpellID:  15036,
		CastID:   16387,
		Kind:     KindManual,
		Category: "FlametongueTotem",
	},
}

var Individual = []BuffSpec{
	{
		Field:   "greater_blessing_of_kings",
		SpellID: 25898,
		Stats: []proto.Stat{
			proto.Stat_StatAgility,
			proto.Stat_StatIntellect,
			proto.Stat_StatSpirit,
			proto.Stat_StatStamina,
			proto.Stat_StatStrength,
		},
	},
	{
		Field:   "greater_blessing_of_might",
		SpellID: 25916,
		Stats: []proto.Stat{
			proto.Stat_StatAttackPower,
		},
	},
	{
		Field:   "greater_blessing_of_wisdom",
		SpellID: 25918,
		Stats: []proto.Stat{
			proto.Stat_StatMP5,
		},
	},
	{
		Field:   "greater_blessing_of_salvation",
		SpellID: 25895,
	},
	// Raises the healing Holy Light and Flash of Light do through flat modifiers on those families.
	// The paladin's heals read the bonus off their own rows and only look for the aura, which the
	// category tags.
	{
		Field:    "greater_blessing_of_light",
		SpellID:  25890,
		Kind:     KindManual,
		Category: "BlessingOfLight",
		Stats: []proto.Stat{
			proto.Stat_StatHealingPower,
		},
	},
	// Auras 134 and 110 are the spirit regen the driver applies.
	{
		Field:    "innervates",
		SpellID:  29166,
		Kind:     KindExternalCD,
		Category: "Innervate",
		Stats: []proto.Stat{
			proto.Stat_StatMP5,
		},
		Label: "Innervates",
	},
	{
		Field:    "power_infusions",
		SpellID:  10060,
		Kind:     KindExternalCD,
		Category: "PowerInfusion",
		Stats: []proto.Stat{
			proto.Stat_StatSpellHasteRating,
		},
		Label: "Power Infusions",
	},
}

var Debuffs = []BuffSpec{
	{
		Field:      "hunters_mark",
		SpellID:    14325,
		Category:   "HuntersMark",
		SingleAura: true,
		Stats: []proto.Stat{
			proto.Stat_StatRangedAttackPower,
			proto.Stat_StatAttackPower,
		},
	},
	// The Holy damage taken is applied by JudgementOfTheCrusaderAura in debuffs.go, which the
	// paladin's own ranks share.
	{
		Field:      "judgement_of_the_crusader",
		SpellID:    20303,
		Category:   "Judgement of the Crusader",
		SingleAura: true,
		Stats: []proto.Stat{
			proto.Stat_StatHolyDamage,
			proto.Stat_StatSpellDamage,
		},
	},
	{
		Field:   "judgement_of_light",
		SpellID: 20346,
		Kind:    KindProc,
		Stats: []proto.Stat{
			proto.Stat_StatDefenseRating,
		},
	},
	{
		Field:   "judgement_of_wisdom",
		SpellID: 20355,
		Kind:    KindProc,
		Stats: []proto.Stat{
			proto.Stat_StatMP5,
		},
	},
	{
		Field:      "curse_of_elements",
		SpellID:    1311680,
		Category:   "CurseOfElements",
		SingleAura: true,
		Stats: []proto.Stat{
			proto.Stat_StatSpellDamage,
		},
	},
	// Forever gives Curse of Recklessness and Faerie Fire the same 505 armor and stops them stacking
	// it (Wowhead's Forever class guides; the client rows carry no category that says so). They share
	// a category per stat, so only the armor competes and the curse keeps the rest.
	{
		Field:    "curse_of_recklessness",
		SpellID:  11717,
		Category: "MinorArmorReduction",
		PerStat:  true,
		Stats: []proto.Stat{
			proto.Stat_StatAttackPower,
		},
	},
	{
		Field:    "faerie_fire",
		SpellID:  9907,
		Category: "MinorArmorReduction",
		PerStat:  true,
		Stats: []proto.Stat{
			proto.Stat_StatAttackPower,
			proto.Stat_StatMeleeHitRating,
		},
	},
	// The -450 armor is per combo point, so the raid config's debuff is the five-point finisher. A
	// rogue spending fewer needs a driver that prices the aura, and there is none.
	{
		Field:      "expose_armor",
		SpellID:    11198,
		Category:   "MajorArmorReduction",
		SingleAura: true,
		Stats: []proto.Stat{
			proto.Stat_StatAttackPower,
		},
	},
	// Five stacks are worth nothing until they are there, so a driver ramps them.
	{
		Field:      "sunder_armor",
		SpellID:    11597,
		Category:   "MajorArmorReduction",
		SingleAura: true,
		Driver:     true,
		Stats: []proto.Stat{
			proto.Stat_StatAttackPower,
		},
	},
	{
		Field:      "gift_of_arthas",
		SpellID:    11374,
		Category:   "GiftOfArthasAura",
		SingleAura: true,
		Stats: []proto.Stat{
			proto.Stat_StatAttackPower,
			proto.Stat_StatDefenseRating,
		},
	},
	{
		Field:      "demoralizing_roar",
		SpellID:    9898,
		Category:   "Demoralizing",
		SingleAura: true,
		Stats: []proto.Stat{
			proto.Stat_StatStamina,
			proto.Stat_StatDefenseRating,
		},
	},
	{
		Field:      "demoralizing_shout",
		SpellID:    11556,
		Category:   "Demoralizing",
		SingleAura: true,
		Stats: []proto.Stat{
			proto.Stat_StatStamina,
			proto.Stat_StatDefenseRating,
		},
	},
	{
		Field:    "thunder_clap",
		SpellID:  11581,
		Category: "AtkSpdReduction",
		Stats: []proto.Stat{
			proto.Stat_StatStamina,
			proto.Stat_StatDefenseRating,
		},
	},
	{
		Field:   "insect_swarm",
		SpellID: 24977,
		Stats: []proto.Stat{
			proto.Stat_StatStamina,
			proto.Stat_StatDefenseRating,
		},
	},
	{
		Field:   "scorpid_sting",
		SpellID: 3043,
		Stats: []proto.Stat{
			proto.Stat_StatStamina,
			proto.Stat_StatDefenseRating,
		},
	},
}
