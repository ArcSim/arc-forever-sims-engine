package hunter

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

type HunterPet struct {
	core.Pet

	config PetConfig

	hunterOwner *Hunter

	BestialWrathAura *core.Aura

	specialAbility *core.Spell
	focusDump      *core.Spell
	extraAbility   *core.Spell

	uptimePercent float64
}

var petAttackSpeeds = [...]float64{
	proto.HunterOptions_One:      1.0,
	proto.HunterOptions_OneTwo:   1.2,
	proto.HunterOptions_OneThree: 1.3,
	proto.HunterOptions_OneFour:  1.4,
	proto.HunterOptions_OneFive:  1.5,
	proto.HunterOptions_OneSix:   1.6,
	proto.HunterOptions_OneSeven: 1.7,
	proto.HunterOptions_Two:      2.0,
	proto.HunterOptions_TwoFour:  2.4,
	proto.HunterOptions_TwoFive:  2.5,
}

func (hunter *Hunter) NewHunterPet() *HunterPet {
	if hunter.Options.PetType == proto.HunterOptions_PetNone {
		return nil
	}

	if hunter.Options.PetUptime <= 0 {
		return nil
	}

	petConfig := DefaultPetConfigs[hunter.Options.PetType]
	attackSpeed := petAttackSpeeds[hunter.Options.PetAttackSpeed]

	conf := core.PetConfig{
		Name:  petConfig.Name,
		Owner: &hunter.Character,
		BaseStats: stats.Stats{
			stats.Strength:  136,
			stats.Agility:   100,
			stats.Stamina:   274,
			stats.Intellect: 50,
			stats.Spirit:    80,

			// Apparently pets and warriors have an AP penalty.
			stats.AttackPower: -20,
		},
		StatInheritance:       hunter.makeStatInheritance(),
		EnabledOnStart:        true,
		IsDynamic:             true,
		IsGuardian:            false,
		StartsAtOwnerDistance: true,
	}

	hp := &HunterPet{
		Pet:         core.NewPet(conf),
		config:      petConfig,
		hunterOwner: hunter,
	}

	hp.Pet.MobType = proto.MobType_MobTypeBeast

	hp.AddStatDependency(stats.Strength, stats.AttackPower, 2.0)
	// No Agility crit: the pet's physical crit is its owner's ranged crit (makeStatInheritance).
	hp.AddStatDependency(stats.Intellect, stats.SpellCritPercent, core.CritPerIntMaxLevel[proto.Class_ClassWarrior])

	// Bestial Discipline buys the pet 10% focus regen a rank (19590 effect 1).
	hp.EnableFocusBar(spellData.BestialDiscipline.EffectAt(1).MultiplierAt(hunter.Talents.BestialDiscipline))

	hp.EnableAutoAttacks(hp, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin: 18.17 * attackSpeed,
			BaseDamageMax: 27.66 * attackSpeed,
			SwingSpeed:    attackSpeed,
			MaxRange:      core.MaxMeleeRange,
		},
		AutoSwingMelee: true,
	})

	// Happiness.
	hp.PseudoStats.DamageDealtMultiplier *= 1.25

	// Family scalars.
	hp.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= petConfig.Damage
	hp.PseudoStats.ArmorMultiplier *= petConfig.Armor
	hp.MultiplyStat(stats.Health, petConfig.Health)

	hunter.AddPet(hp)
	return hp
}

// Cobra Reflexes is TBC's; Forever's pets do not have it.
func (hp *HunterPet) ApplyTalents() {
}

func (hp *HunterPet) GetPet() *core.Pet {
	return &hp.Pet
}

func (hp *HunterPet) Initialize() {
	hp.Pet.Initialize()

	hp.specialAbility = hp.NewPetAbility(hp.config.SpecialAbility)
	hp.focusDump = hp.NewPetAbility(hp.config.FocusDump)
	hp.extraAbility = hp.NewPetAbility(hp.config.ExtraAbility)
}

func (hp *HunterPet) Reset(_ *core.Simulation) {
	hp.uptimePercent = min(1, max(0, hp.hunterOwner.Options.PetUptime))
}

func (hp *HunterPet) OnEncounterStart(_ *core.Simulation) {
}

func (hp *HunterPet) ExecuteCustomRotation(sim *core.Simulation) {
	percentRemaining := sim.GetRemainingDurationPercent()
	if percentRemaining < 1.0-hp.uptimePercent { // once fight is % completed, disable pet.
		hp.Disable(sim)
		return
	}

	if hp.DistanceFromTarget > core.MaxMeleeRange {
		if !hp.Moving {
			hp.MoveTo(core.MaxMeleeRange-1, sim)
		}
		return
	}

	target := hp.CurrentTarget

	// Using Cast() directly is expensive, since cast failures are logged, involving string
	// operations.
	tryCast := func(spell *core.Spell) bool {
		if spell == nil || !spell.CanCast(sim, target) {
			return false
		}
		spell.Cast(sim, target)
		return true
	}

	if hp.config.CustomRotation != nil {
		hp.config.CustomRotation(sim, hp, tryCast)
		return
	}

	if !tryCast(hp.extraAbility) && !tryCast(hp.specialAbility) && !tryCast(hp.focusDump) && hp.GCD.IsReady(sim) {
		hp.WaitUntil(sim, sim.CurrentTime+time.Millisecond*500)
	}
}

// Forever is a Classic-era realm: the pet takes none of its owner's stats, where TBC's does, but one.
// From 16 October 2026 its physical crit is its master's ranged crit and it has no crit of its own (dev
// AidanZMoon, forever-bugs #331: "It currently uses Melee Crit on Beta but will use Ranged Crit starting
// Thursday"; pets "don't have erroneous base crit"). Auto Shot rolls PhysicalCritPercent plus the
// ranged-only RangedCritPercent; the weapon racials are melee spell mods, so they stay out as they do
// for Auto Shot. Ferocity and Intimidation still add on the pet.
func (hunter *Hunter) makeStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.PhysicalCritPercent: ownerStats[stats.PhysicalCritPercent] + ownerStats[stats.RangedCritPercent],
		}
	}
}

type PetConfig struct {
	Name string

	SpecialAbility PetAbilityType
	FocusDump      PetAbilityType
	// Cast first, whenever it is ready: a cooldown the pet would never skip.
	ExtraAbility PetAbilityType

	Health float64
	Armor  float64
	Damage float64

	CustomRotation func(*core.Simulation, *HunterPet, func(*core.Spell) bool)
}

var DefaultPetConfigs = [...]PetConfig{
	proto.HunterOptions_PetNone: {},
	proto.HunterOptions_Cat: {
		Name: "Cat", SpecialAbility: Bite, FocusDump: Claw,
		Health: 0.98, Armor: 1.00, Damage: 1.10,
		// Bite is worth more per Focus than Claw, so it goes on cooldown and Claw fills in.
		CustomRotation: func(sim *core.Simulation, hp *HunterPet, tryCast func(*core.Spell) bool) {
			if !tryCast(hp.specialAbility) && !tryCast(hp.focusDump) && hp.GCD.IsReady(sim) {
				hp.WaitUntil(sim, sim.CurrentTime+time.Millisecond*500)
			}
		},
	},
	proto.HunterOptions_WindSerpent: {
		Name: "Wind Serpent", SpecialAbility: Bite, FocusDump: LightningBreath,
		Health: 1.00, Armor: 1.00, Damage: 1.07,
	},
	proto.HunterOptions_Wolf: {
		Name: "Wolf", FocusDump: Bite,
		Health: 1.00, Armor: 1.05, Damage: 1.00,
	},
	proto.HunterOptions_Bat: {
		Name: "Bat", SpecialAbility: Bite, FocusDump: Screech,
		Health: 1.00, Armor: 1.00, Damage: 1.07,
	},
	proto.HunterOptions_Bear: {
		Name: "Bear", SpecialAbility: Bite, FocusDump: Claw, ExtraAbility: Swipe,
		Health: 1.08, Armor: 1.05, Damage: 0.91,
	},
	proto.HunterOptions_Owl: {
		Name: "Owl", FocusDump: Claw, ExtraAbility: Mine,
		Health: 1.00, Armor: 1.00, Damage: 1.07,
	},
	proto.HunterOptions_Boar: {
		Name: "Boar", FocusDump: Bite,
		Health: 1.04, Armor: 1.09, Damage: 0.90,
	},
	proto.HunterOptions_CarrionBird: {
		Name: "Carrion Bird", SpecialAbility: Bite, FocusDump: Claw,
		Health: 1.00, Armor: 1.05, Damage: 1.00,
	},
	proto.HunterOptions_Crab: {
		Name: "Crab", FocusDump: Claw, ExtraAbility: Pinch,
		Health: 0.96, Armor: 1.13, Damage: 0.95,
	},
	proto.HunterOptions_Crocolisk: {
		Name: "Crocolisk", SpecialAbility: Bite, FocusDump: Dismember,
		Health: 0.95, Armor: 1.10, Damage: 1.00,
	},
	proto.HunterOptions_Gorilla: {
		Name: "Gorilla", FocusDump: Bite, ExtraAbility: Thunderstomp,
		Health: 1.04, Armor: 1.00, Damage: 1.02,
	},
	proto.HunterOptions_Hyena: {
		Name: "Hyena", FocusDump: Bite, ExtraAbility: TendonRip,
		Health: 1.00, Armor: 1.05, Damage: 1.00,
	},
	proto.HunterOptions_Raptor: {
		Name: "Raptor", SpecialAbility: Bite, FocusDump: Claw, ExtraAbility: SavageRend,
		Health: 0.95, Armor: 1.03, Damage: 1.10,
	},
	proto.HunterOptions_Scorpid: {
		Name: "Scorpid", SpecialAbility: ScorpidPoison, FocusDump: Claw,
		Health: 1.00, Armor: 1.10, Damage: 0.94,

		CustomRotation: func(sim *core.Simulation, hp *HunterPet, tryCast func(*core.Spell) bool) {
			target := hp.CurrentTarget
			dot := hp.specialAbility.Dot(target)

			if dot.GetStacks() < dot.MaxStacks || dot.RemainingDuration(sim) < time.Second*3 {
				if !tryCast(hp.specialAbility) && hp.GCD.IsReady(sim) {
					hp.WaitUntil(sim, sim.CurrentTime+time.Millisecond*500)
				}
			} else if !tryCast(hp.focusDump) && hp.GCD.IsReady(sim) {
				hp.WaitUntil(sim, sim.CurrentTime+time.Millisecond*500)
			}
		},
	},
	proto.HunterOptions_Spider: {
		Name: "Spider", FocusDump: Bite, ExtraAbility: Web,
		Health: 1.00, Armor: 1.00, Damage: 1.07,
	},
	proto.HunterOptions_Tallstrider: {
		Name: "Tallstrider", FocusDump: Bite, ExtraAbility: DustCloud,
		Health: 1.05, Armor: 1.00, Damage: 1.00,
	},
	proto.HunterOptions_Turtle: {
		Name: "Turtle", FocusDump: Bite,
		Health: 1.00, Armor: 1.13, Damage: 0.90,
	},
}
