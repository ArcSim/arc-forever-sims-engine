package paladin

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

var TalentTreeSizes = [3]int{17, 16, 17}

type Paladin struct {
	core.Character

	Talents *proto.PaladinTalents

	Forbearance       *core.Aura
	RighteousFuryAura *core.Aura

	Judgement *core.Spell

	// The seal the paladin is under, and the ranks it has cast so far this iteration.
	currentSeal *sealConfig

	// The judgement effects this paladin can put on enemies; its melee strikes refresh them.
	JudgementAuras []core.AuraArray

	// Twist of Light: one Echo per seal that leaves one, keyed by the Echo's spell id.
	echoes map[int32]*sealEcho

	// Consecrated Ground: the targets the talent's bonus is active on, one aura per enemy.
	consecratedGroundAuras core.AuraArray

	// Light's Vigil: one entry per rank, so Holy Shock can find the vigil it consumes.
	lightsVigils []*lightsVigil

	// What gear adds to numbers the spells read as they register. Item effects and set bonuses
	// apply before Initialize, so the spells pick these up.
	sealOfTheCrusaderBonusAttackPower float64
	judgementOfTheCrusaderBonus       float64
	flashOfLightBonusHealing          float64
	holyShieldBlockValueMultiplier    float64
	forbearanceReduction              time.Duration

	// Timers shared by the ranks of one ability.
	judgementTimer     *core.Timer
	holyStrikeTimer    *core.Timer
	consecrationTimer  *core.Timer
	exorcismTimer      *core.Timer
	hammerOfWrathTimer *core.Timer
	holyWrathTimer     *core.Timer
	holyShockTimer     *core.Timer
	holyShieldTimer    *core.Timer
	lightsVigilTimer   *core.Timer
}

// Implemented by each Paladin spec.
type PaladinAgent interface {
	GetPaladin() *Paladin
}

func (paladin *Paladin) GetCharacter() *core.Character {
	return &paladin.Character
}

func (paladin *Paladin) GetPaladin() *Paladin {
	return paladin
}

func (paladin *Paladin) AddRaidBuffs(_ *proto.RaidBuffs) {
}

func (paladin *Paladin) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (paladin *Paladin) Initialize() {
	paladin.registerForbearance()
	paladin.registerRighteousFury()

	paladin.registerJudgement()
	paladin.registerSeals()
	paladin.registerAuras()

	HolyStrikeRankMap.Each(paladin.registerHolyStrike)
	paladin.registerHammerOfTheRighteous()
	ConsecrationRankMap.Each(paladin.registerConsecration)
	ExorcismRankMap.Each(paladin.registerExorcism)
	HammerOfWrathRankMap.Each(paladin.registerHammerOfWrath)
	HolyWrathRankMap.Each(paladin.registerHolyWrath)

	HolyLightRankMap.Each(paladin.registerHolyLight)
	FlashOfLightRankMap.Each(paladin.registerFlashOfLight)
	LayOnHandsRankMap.Each(paladin.registerLayOnHands)
}

func (paladin *Paladin) Reset(_ *core.Simulation) {
	paladin.currentSeal = nil
	for _, echo := range paladin.echoes {
		echo.seal = nil
	}
}

func (paladin *Paladin) OnEncounterStart(_ *core.Simulation) {
}

func (paladin *Paladin) GetMainHandType() proto.HandType {
	mh := paladin.GetMHWeapon()

	if mh != nil && (mh.HandType == proto.HandType_HandTypeTwoHand) {
		return proto.HandType_HandTypeTwoHand
	}

	return proto.HandType_HandTypeOneHand
}

func (paladin *Paladin) sharedTimer(timer **core.Timer) *core.Timer {
	if *timer == nil {
		*timer = paladin.NewTimer()
	}
	return *timer
}

// A rank's own cooldown, or the one it shares with its category: the client keeps every paladin
// ability cooldown but Judgement's and the talent cooldowns on the category.
func cooldown(rank *spelldata.Spell) time.Duration {
	return max(rank.Cooldown(), rank.CategoryCooldown())
}

func NewPaladin(character *core.Character, talentsStr string, _ *proto.PaladinOptions) *Paladin {
	paladin := &Paladin{
		Character: *character,
		Talents:   &proto.PaladinTalents{},

		holyShieldBlockValueMultiplier: 1,
	}

	core.FillTalentsProto(paladin.Talents.ProtoReflect(), talentsStr, TalentTreeSizes)

	// The attack table already holds the base 5% parry and block (sim/core/target.go); adding them
	// here too gave the paladin 10% of each. Base dodge 0.7% and 1% dodge per 19.8 Agility are the
	// Classic paladin's, the same as its crit, as on master.
	paladin.PseudoStats.CanParry = true
	paladin.PseudoStats.BaseDodgeChance += 0.007

	paladin.EnableManaBar()

	paladin.EnableAutoAttacks(paladin, core.AutoAttackOptions{
		MainHand:       paladin.WeaponFromMainHand(),
		AutoSwingMelee: true,
	})

	paladin.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	// Block value from Strength is Classic's Str/20 - 1, as the warrior's and master's.
	paladin.AddStatDependency(stats.Strength, stats.BlockValue, 1/20.0)
	paladin.AddStat(stats.BlockValue, -1)
	paladin.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[character.Class])
	paladin.AddStatDependency(stats.Agility, stats.DodgeRating, core.CritPerAgiMaxLevel[character.Class]*core.DodgeRatingPerDodgePercent)
	paladin.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	return paladin
}

// holdMeleeForCast is the Forever hard-cast rule for a Paladin's cast-time spells
// (core.AutoAttacks.HoldMeleeForCast): the next swing comes a full weapon speed after the cast
// completes. A cast made instant (Instrument of Law, a proc) leaves the swing timer alone.
func (paladin *Paladin) holdMeleeForCast(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
	if castTime := paladin.ApplyCastSpeedForSpell(cast.CastTime, spell); castTime > 0 {
		paladin.AutoAttacks.HoldMeleeForCast(sim, sim.CurrentTime+castTime)
	}
}
