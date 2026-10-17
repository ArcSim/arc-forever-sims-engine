package hunter

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
)

type HunterPet struct {
	core.Pet

	config PetConfig

	hunterOwner *Hunter

	// the family's abilities the player left on, highest priority first
	abilities []*petAbility

	uptimePercent float64
}

// A pet ability in the rotation. Abilities with a cooldown save focus for themselves when they come up;
// fillers are only used with focus to spare.
type petAbility struct {
	spell  *core.Spell
	filler bool
	// extra condition, e.g. a debuff that is still running
	wanted func(sim *core.Simulation) bool
}

func (hunter *Hunter) NewHunterPet() *HunterPet {
	if hunter.Talents.LoneWolf {
		return nil
	}
	if hunter.Options.PetType == proto.Hunter_Options_PetNone {
		return nil
	}
	if hunter.Options.PetUptime <= 0 {
		return nil
	}
	petConfig := PetConfigs[hunter.Options.PetType]

	hunterPetBaseStats := stats.Stats{}

	baseMinDamage := 0.0
	baseMaxDamage := 0.0
	attackSpeed := 2.0

	// every Forever pet starts at 2.0; the speed a pet shows when tamed already includes its
	// Faster Attack (1.7 = I ... 1.0 = VII) or Slower Attack (2.4, 2.5) passive
	switch hunter.Options.PetAttackSpeed {
	case proto.Hunter_Options_One:
		attackSpeed = 1.0
	case proto.Hunter_Options_OneTwo:
		attackSpeed = 1.2
	case proto.Hunter_Options_OneThree:
		attackSpeed = 1.3
	case proto.Hunter_Options_OneFour:
		attackSpeed = 1.4
	case proto.Hunter_Options_OneFive:
		attackSpeed = 1.5
	case proto.Hunter_Options_OneSix:
		attackSpeed = 1.6
	case proto.Hunter_Options_OneSeven:
		attackSpeed = 1.7
	case proto.Hunter_Options_Two:
		attackSpeed = 2
	case proto.Hunter_Options_TwoFour:
		attackSpeed = 2.4
	case proto.Hunter_Options_TwoFive:
		attackSpeed = 2.5
	}

	baseMinDamage = 18.17 * attackSpeed
	baseMaxDamage = 27.66 * attackSpeed

	hunterPetBaseStats = stats.Stats{
		stats.Strength:  136,
		stats.Agility:   100,
		stats.Stamina:   274,
		stats.Intellect: 50,
		stats.Spirit:    80,

		stats.AttackPower: -20,
	}

	hp := &HunterPet{
		Pet:         core.NewPet(petConfig.Name, &hunter.Character, hunterPetBaseStats, hunter.makeStatInheritance(), true, false),
		config:      petConfig,
		hunterOwner: hunter,
	}

	hp.Pet.MobType = petConfig.MobType

	// AB_PET_BUG=1 (harness only) makes every pet hit as hard as a 2.0-speed pet, so fast pets deal 2.0/speed times their
	// proper DPS: the Forever beta bug, which the game does not have, kept for comparisons.
	apPerDPS := core.DefaultAttackPowerPerDPS
	if os.Getenv("AB_PET_BUG") != "" {
		baseMinDamage = 18.17 * 2.0
		baseMaxDamage = 27.66 * 2.0
		apPerDPS = core.DefaultAttackPowerPerDPS * attackSpeed / 2.0
	}

	hp.EnableAutoAttacks(hp, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:     baseMinDamage,
			BaseDamageMax:     baseMaxDamage,
			SwingSpeed:        attackSpeed,
			AttackPowerPerDPS: apPerDPS,
		},
		AutoSwingMelee: true,
	})

	// Happiness: always happy (+25% damage)
	hp.PseudoStats.DamageDealtMultiplier *= 1.25

	// Family scalars
	hp.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= hp.config.Damage
	hp.PseudoStats.ArmorMultiplier *= hp.config.Armor
	hp.MultiplyStat(stats.Health, hp.config.Health)

	hp.AddStatDependency(stats.Strength, stats.AttackPower, 2)

	// Warrior crit scaling
	hp.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[proto.Class_ClassWarrior]*core.CritRatingPerCritChance)
	hp.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[proto.Class_ClassWarrior]*core.SpellCritRatingPerCritChance)

	core.ApplyPetConsumeEffects(&hp.Character, hunter.Consumes)

	hunter.AddPet(hp)

	return hp
}

func (hp *HunterPet) GetPet() *core.Pet {
	return &hp.Pet
}

func (hp *HunterPet) Initialize() {
	// pets take the higher of the hunter's melee and ranged sheet crit; only the main-hand stone sets them apart.
	// Added once here: the inheritance above also runs on every change of the hunter's stats.
	hp.AddStat(stats.MeleeCrit, max(hp.hunterOwner.MainHandCritBonus, 0))

	disabled := hp.hunterOwner.Options.PetDisabledAbilities
	for _, ability := range hp.config.Abilities {
		if slices.Contains(disabled, ability.option()) {
			continue
		}
		if pa := hp.newPetAbility(ability); pa != nil {
			hp.abilities = append(hp.abilities, pa)
		}
	}

	// Forever pets regenerate focus twice as fast (logs: about 10 focus per second spent on Claw)
	hp.EnableFocusBar(2, func(sim *core.Simulation) {
		if hp.GCD.IsReady(sim) {
			hp.OnGCDReady(sim)
		}
	})
}

func (hp *HunterPet) Reset(_ *core.Simulation) {
	hp.uptimePercent = min(1, max(0, hp.hunterOwner.Options.PetUptime))
}

func (hp *HunterPet) ExecuteCustomRotation(sim *core.Simulation) {
	percentRemaining := sim.GetRemainingDurationPercent()
	if percentRemaining < 1.0-hp.uptimePercent { // once fight is % completed, disable pet.
		hp.Disable(sim)
		return
	}

	// no abilities left on: the pet only auto attacks
	if len(hp.abilities) == 0 {
		return
	}

	target := hp.CurrentTarget

	// using Cast() directly is very expensive, since cast failures are logged, involving string operations
	tryCast := func(spell *core.Spell) bool {
		if !spell.CanCast(sim, target) {
			return false
		}
		if !spell.Cast(sim, target) {
			panic(fmt.Sprintf("Cast failed after CanCast() for spell %d", spell.SpellID))
		}
		return true
	}

	for _, pa := range hp.abilities {
		if !pa.spell.IsReady(sim) || (pa.wanted != nil && !pa.wanted(sim)) {
			continue
		}
		if tryCast(pa.spell) {
			return
		}
		if !pa.filler {
			// a cooldown ability is up: save focus for it
			break
		}
	}
	if hp.GCD.IsReady(sim) {
		hp.WaitUntil(sim, sim.CurrentTime+time.Millisecond*500)
	}
}

func (hunter *Hunter) makeStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		ap := ownerStats[stats.AttackPower]
		rap := ownerStats[stats.RangedAttackPower]
		if rap > ap {
			ap = rap
		}
		return stats.Stats{
			stats.AttackPower: ap * 0.10,
			stats.MeleeCrit:   ownerStats[stats.MeleeCrit],
			stats.Expertise:   ownerStats[stats.Expertise],
			stats.SpellCrit:   ownerStats[stats.SpellCrit],
			stats.MeleeHit:    ownerStats[stats.MeleeHit],
			stats.SpellHit:    ownerStats[stats.MeleeHit],
		}
	}
}

type PetConfig struct {
	Name    string
	MobType proto.MobType

	// the family's damage abilities (Forever), highest priority first
	Abilities []PetAbilityType

	Health float64
	Armor  float64
	Damage float64
}

// Families and their abilities per WoW Forever (wow-petopia.com/forever, checked against client data).
// Health / armor / damage are the family scalars; Core Hound and Fox list damage 1.00 and nothing else yet.
var PetConfigs = map[proto.Hunter_Options_PetType]PetConfig{
	proto.Hunter_Options_Bat:         {Name: "Bat", Abilities: []PetAbilityType{SonicBlast, Bite}, Health: 1.00, Armor: 1.00, Damage: 1.07},
	proto.Hunter_Options_Bear:        {Name: "Bear", Abilities: []PetAbilityType{Bite, Claw, Swipe}, Health: 1.08, Armor: 1.05, Damage: 0.91},
	proto.Hunter_Options_BirdOfPrey:  {Name: "Bird of Prey", Abilities: []PetAbilityType{Mine, Claw}, Health: 1.00, Armor: 1.00, Damage: 1.07},
	proto.Hunter_Options_Owl:         {Name: "Owl", Abilities: []PetAbilityType{Mine, Claw}, Health: 1.00, Armor: 1.00, Damage: 1.07},
	proto.Hunter_Options_Boar:        {Name: "Boar", Abilities: []PetAbilityType{Charge, Bite}, Health: 1.04, Armor: 1.09, Damage: 0.90},
	proto.Hunter_Options_CarrionBird: {Name: "Carrion Bird", Abilities: []PetAbilityType{Bite, Claw, DemoralizingScreech}, Health: 1.00, Armor: 1.05, Damage: 1.00},
	proto.Hunter_Options_Cat:         {Name: "Cat", Abilities: []PetAbilityType{Bite, Claw}, Health: 0.98, Armor: 1.00, Damage: 1.10},
	proto.Hunter_Options_CoreHound:   {Name: "Core Hound", Abilities: []PetAbilityType{LavaBreath, Bite}, Health: 1.00, Armor: 1.00, Damage: 1.00},
	proto.Hunter_Options_Crab:        {Name: "Crab", Abilities: []PetAbilityType{Pinch, Claw}, Health: 0.96, Armor: 1.13, Damage: 0.95},
	proto.Hunter_Options_Crocolisk:   {Name: "Crocolisk", Abilities: []PetAbilityType{Bite, Dismember}, Health: 0.95, Armor: 1.10, Damage: 1.00},
	proto.Hunter_Options_Fox:         {Name: "Fox", Abilities: []PetAbilityType{Bite, TrickstersDance}, Health: 1.00, Armor: 1.00, Damage: 1.00},
	proto.Hunter_Options_Gorilla:     {Name: "Gorilla", Abilities: []PetAbilityType{Thunderstomp, Bite}, Health: 1.04, Armor: 1.00, Damage: 1.02},
	proto.Hunter_Options_Hyena:       {Name: "Hyena", Abilities: []PetAbilityType{TendonRip, Bite}, Health: 1.00, Armor: 1.05, Damage: 1.00},
	proto.Hunter_Options_Raptor:      {Name: "Raptor", Abilities: []PetAbilityType{SavageRend, Bite, Claw}, Health: 0.95, Armor: 1.03, Damage: 1.10},
	proto.Hunter_Options_Scorpid:     {Name: "Scorpid", Abilities: []PetAbilityType{ScorpidPoison, Claw}, Health: 1.00, Armor: 1.10, Damage: 0.94},
	proto.Hunter_Options_Spider:      {Name: "Spider", Abilities: []PetAbilityType{Web, Bite}, Health: 1.00, Armor: 1.00, Damage: 1.07},
	proto.Hunter_Options_Tallstrider: {Name: "Tallstrider", Abilities: []PetAbilityType{DustCloud, Bite}, Health: 1.05, Armor: 1.00, Damage: 1.00},
	proto.Hunter_Options_Turtle:      {Name: "Turtle", Abilities: []PetAbilityType{Bite}, Health: 1.00, Armor: 1.13, Damage: 0.90},
	proto.Hunter_Options_WindSerpent: {Name: "Wind Serpent", Abilities: []PetAbilityType{Bite, LightningBreath}, Health: 1.00, Armor: 1.00, Damage: 1.07},
	proto.Hunter_Options_Wolf:        {Name: "Wolf", Abilities: []PetAbilityType{FuriousHowl, Bite}, Health: 1.00, Armor: 1.05, Damage: 1.00},
}

func init() {
	for family, config := range PetConfigs {
		config.MobType = proto.MobType_MobTypeBeast
		PetConfigs[family] = config
	}
}
