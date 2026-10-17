package core

import (
	"os"
	"strconv"
	"time"
)

const CharacterMaxLevel = 60

const GCDMin = time.Second * 1
const GCDDefault = time.Millisecond * 1500
const SpellBatchWindow = time.Millisecond * 10

const DefaultAttackPowerPerDPS = 14.0
const ArmorPenPerPercentArmor = 13.99

const MaxMeleeAttackDistance = 5
// Ranged minimum 8 yd: Forever client data ("Hunter Range" 8-35 for Auto Shot and every shot). Knob AB_MIN_RANGED for experiments.
var MinRangedAttackDistance = 8.0

func init() {
	if v, err := strconv.ParseFloat(os.Getenv("AB_MIN_RANGED"), 64); err == nil && v > 0 {
		MinRangedAttackDistance = v
	}
}

const MissDodgeParryBlockCritChancePerDefense = 0.04

const DefenseRatingToChanceReduction = (1.0 / DefenseRatingPerDefense) * MissDodgeParryBlockCritChancePerDefense / 100

const ResilienceRatingPerCritDamageReductionPercent = ResilienceRatingPerCritReductionChance / 2.2

// Updated based on formulas supplied by InDebt on WoWSims Discord
const EnemyAutoAttackAPCoefficient = 1.0 / (14.0 * 177.0)

const AverageMagicPartialResistPerLevelMultiplier = 0.02

// IDs for items used in core
const (
	ItemIDBraidedEterniumChain  = 24114
	ItemIDChainOfTheTwilightOwl = 24121
	ItemIDEyeOfTheNight         = 24116
	ItemIDJadePendantOfBlasting = 20966
	ItemIDTheLightningCapacitor = 28785
)

type Hand bool

const MainHand Hand = true
const OffHand Hand = false

type DefenseType byte

const (
	DefenseTypeNone DefenseType = iota
	DefenseTypeMagic
	DefenseTypeMelee
	DefenseTypeRanged

	DefenseTypeLen
)
