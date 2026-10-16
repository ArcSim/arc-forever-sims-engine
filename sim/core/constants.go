package core

import (
	"time"

	"github.com/wowsims/forever/sim/core/proto"
)

const CharacterLevel = 60
const MinIlvl = 60
const MaxIlvl = 600

const GCDMin = time.Second * 1
const GCDDefault = time.Millisecond * 1500
const BossGCD = time.Millisecond * 1620
const MaxSpellQueueWindow = time.Millisecond * 400
const SpellBatchWindow = time.Millisecond * 10
const PetUpdateInterval = time.Millisecond * 5250
const SpellPushbackDuration = time.Millisecond * 500

// How often a ranged auto that came due while moving checks whether it can fire.
const RangedAutoRetryInterval = time.Millisecond * 500

// How long a unit has to stand still before its ranged auto can fire (16 October 2026 beta notes: "Auto
// Shot now requires 0.5 seconds of standing still before shooting").
const RangedAutoStandStill = time.Millisecond * 500
const MaxMeleeRange = 5.0  // in yards
const MinRangedRange = 8.0 // in yards; bows, guns, crossbows and thrown weapons cannot fire inside this, leaving a deadzone above melee range

const DefaultAttackPowerPerDPS = 14.0

const ArmorPenPerPercentArmor = 5.92
const MissDodgeParryBlockCritChancePerDefense = 0.04
const DefenseRatingPerAvoidancePercent = DefenseRatingPerDefenseLevel / MissDodgeParryBlockCritChancePerDefense

const EnemyAutoAttackAPCoefficient = 0.00052

// IDs for items used in core
// const ()

type Hand bool

const MainHand Hand = true
const OffHand Hand = false

const CombatTableCoverageCap = 1.024 // 102.4% chance to avoid an attack

const NumItemSlots = proto.ItemSlot_ItemSlotRanged + 1

func TrinketSlots() []proto.ItemSlot {
	return []proto.ItemSlot{proto.ItemSlot_ItemSlotTrinket1, proto.ItemSlot_ItemSlotTrinket2}
}

func AllWeaponSlots() []proto.ItemSlot {
	return []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand, proto.ItemSlot_ItemSlotOffHand, proto.ItemSlot_ItemSlotRanged}
}

func AllMeleeWeaponSlots() []proto.ItemSlot {
	return []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand, proto.ItemSlot_ItemSlotOffHand}
}

//go:generate stringer -type=DefenseType

// Which hit table a spell rolls on and which crit multiplier it takes. Values match the
// DefenseType column of the client's SpellCategories table, so a spell's value is looked up
// there rather than inferred from its school: Thunder Clap is Physical but rolls as Magic,
// Seal of Command procs are Holy but roll as Melee.
type DefenseType byte

const (
	DefenseTypeNone DefenseType = iota
	DefenseTypeMagic
	DefenseTypeMelee
	DefenseTypeRanged

	DefenseTypeLen
)
