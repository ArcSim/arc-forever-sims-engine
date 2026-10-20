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
	// No Agility crit: the pet's physical crit is its owner's (patch 136, makeStatInheritance).
	hp.AddStatDependency(stats.Intellect, stats.SpellCritPercent, core.CritPerIntMaxLevel[proto.Class_ClassWarrior])
	hp.trackPetCritOffset()

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

	if !tryCast(hp.extraAbility) && !tryCast(hp.specialAbility) && !hp.tryFocusDump(sim, tryCast) && hp.GCD.IsReady(sim) {
		hp.WaitUntil(sim, sim.CurrentTime+time.Millisecond*500)
	}
}

// PetPoolsFocusForCooldowns (MythicSim patch 118) makes a pet hold its focus dump while that would leave
// it short of a higher priority ability with a cooldown: the Cat Claws only when the focus left after Claw,
// plus what it regenerates before Bite is ready, still pays for Bite, so Bite goes on cooldown. Off by
// default: the pet's autocast spends Claw (25 focus, no cooldown) the moment it has 25, so focus never
// climbs to Bite's 35 once Bite is off cooldown and Bite lands once, off the full bar the pet starts
// with (Kragtor, Discord, 12 October 2026: one Bite in 240 seconds), and upstream modelled that from beta
// logs (TestPetStrikes: long-cooldown strikes land about once a fight, Claw spending the focus first).
// The Raptor's Savage Rend (50) and the Gorilla's Thunderstomp (60) starve the same way. Set it true for
// a pet whose Claw autocast the Hunter manages around Bite.
var PetPoolsFocusForCooldowns = false

// tryFocusDump casts the focus dump unless it would starve the pet's extra or special ability.
func (hp *HunterPet) tryFocusDump(sim *core.Simulation, tryCast func(*core.Spell) bool) bool {
	if hp.focusDump == nil {
		return false
	}
	if PetPoolsFocusForCooldowns {
		left := hp.CurrentFocus() - hp.focusDump.Cost.GetCurrentCost()
		for _, ability := range []*core.Spell{hp.extraAbility, hp.specialAbility} {
			if hp.starves(sim, ability, left) {
				return false
			}
		}
	}
	return tryCast(hp.focusDump)
}

// starves reports whether a pet left with focus left could not pay for ability when its cooldown ends.
// Abilities without a cooldown (Dust Cloud), or that could not be cast now for another reason (Swipe on
// fewer than three targets, Furious Howl while it is up), are not pooled for.
func (hp *HunterPet) starves(sim *core.Simulation, ability *core.Spell, left float64) bool {
	if ability == nil || ability == hp.focusDump || ability.CD.Timer == nil {
		return false
	}
	if ability.ExtraCastCondition != nil && !ability.ExtraCastCondition(sim, hp.CurrentTarget) {
		return false
	}
	regen := hp.FocusRegenPerSecond() * max(ability.TimeToReady(sim), hp.GCD.TimeToReady(sim)+PetGCD).Seconds()
	return left+regen < ability.Cost.GetCurrentCost()
}

// Forever is a Classic-era realm: the pet takes none of its owner's stats, where TBC's does, but two.
// From 16 October 2026 the pet's crit is its master's crit (MythicSim patch 136; Forever dev to Bae in
// #hunter: "Pets will use Ranged Crit starting Thursday"). It replaces the pet's own crit, which came from
// its 100 base Agility at the Warrior's 20 Agility a point (5%) plus any Agility it gained.
//
// 16 October 2026 beta notes: "Pets now inherit whichever crit is higher, melee or ranged (fix of pets
// only inheriting melee crit)", and "Pets now benefit from Expertise". The Hunter's Auto Shot rolls
// PhysicalCritPercent (Agility, gear, buffs, Lethal Attacks) plus RangedCritPercent (the ranged-only
// offset: scopes, less a melee-only stone) plus the Auto Shot spell's own crit mods; the melee swing rolls
// PhysicalCritPercent plus the main-hand auto's crit mods. The pet takes PhysicalCritPercent plus the
// higher of the two offsets, so an Elemental Sharpening Stone (+2% physical, -2% ranged) leaves the pet
// on the melee value. It also takes the owner's ExpertisePercent (rating already converted).
//
// The owner's stat changes reach the pet on the pet heartbeat (core.PetUpdateInterval), as every dynamic
// inheritance does. The core passes stat deltas, which suits the linear parts; the max is not linear, so
// petCritOffset holds the offset currently on the pet and each update moves it to the offset the owner's
// current stats and spell mods give (a spell mod gained or lost is picked up at the next owner stat
// change). The target's side (crit suppression, ReducedCritTakenPercent) stays on the pet's own attack
// table, as for any attacker. Crit auras on the pet itself still add: Ferocity ("Increases the critical
// strike chance of your pets and hawks", 19598) and Intimidation's next attack.
// PetsInheritHit (mythicsim patch 178, ported from wowsims/forever a205077834): the pet takes the hunter's hit,
// the higher of melee hit and melee plus ranged hit. Hunter Pet Scaling (415429) lists hit among the stats a
// pet inherits. Ranged hit only ever adds (scopes, gear, talents), so the higher of the two is melee plus
// ranged, which keeps the inheritance linear: the owner's stats reach the pet as changes, and the Rust engine
// carries the inheritance as linear terms. False leaves the pet on its own hit, as before.
var PetsInheritHit = true

func (hunter *Hunter) makeStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		offset := hunter.petCritOffsetNow()
		change := offset - hunter.petCritOffset
		hunter.petCritOffset = offset
		inherited := stats.Stats{
			stats.PhysicalCritPercent: ownerStats[stats.PhysicalCritPercent] + change,
			stats.ExpertisePercent:    ownerStats[stats.ExpertisePercent],
		}
		if PetsInheritHit {
			inherited[stats.PhysicalHitPercent] = ownerStats[stats.PhysicalHitPercent] + ownerStats[stats.RangedHitPercent]
		}
		return inherited
	}
}

// petCritOffsetNow is the crit the pet takes on top of the owner's PhysicalCritPercent: the higher of the
// Auto Shot's ranged-only crit and the main-hand swing's melee-only crit (16 October 2026 beta notes).
func (hunter *Hunter) petCritOffsetNow() float64 {
	ranged := hunter.GetStat(stats.RangedCritPercent)
	if autoShot := hunter.AutoAttacks.RangedAuto(); autoShot != nil {
		ranged += autoShot.BonusCritPercent
	}
	melee := 0.0
	if swing := hunter.AutoAttacks.MHAuto(); swing != nil {
		melee = swing.BonusCritPercent
	}
	return max(ranged, melee)
}

// trackPetCritOffset clears petCritOffset once the pet is dismissed: the core has taken back everything
// the pet inherited, offset included, so the next summon starts from nothing.
func (hp *HunterPet) trackPetCritOffset() {
	hp.OnPetDisable = func(sim *core.Simulation) {
		hp.hunterOwner.petCritOffset = 0
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
			if !tryCast(hp.specialAbility) && !hp.tryFocusDump(sim, tryCast) && hp.GCD.IsReady(sim) {
				hp.WaitUntil(sim, sim.CurrentTime+time.Millisecond*500)
			}
		},
	},
	proto.HunterOptions_WindSerpent: {
		Name: "Wind Serpent", SpecialAbility: Bite, FocusDump: LightningBreath,
		Health: 1.00, Armor: 1.00, Damage: 1.07,
	},
	proto.HunterOptions_Wolf: {
		Name: "Wolf", FocusDump: Bite, ExtraAbility: FuriousHowl,
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
