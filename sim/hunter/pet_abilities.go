package hunter

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/buffs"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
)

type PetAbilityType int

// Pet AI doesn't use abilities immediately, so model this with a 1.6s GCD.
const PetGCD = time.Millisecond * 1600

const (
	Unknown PetAbilityType = iota

	Bite
	Claw
	LightningBreath
	Screech
	ScorpidPoison
	SavageRend
	Pinch
	Dismember
	Mine
	TendonRip
	DustCloud
	Thunderstomp
	Swipe
	Web
	FuriousHowl
)

func (hp *HunterPet) NewPetAbility(abilityType PetAbilityType) *core.Spell {
	switch abilityType {
	case Bite:
		return hp.newBite()
	case Claw:
		return hp.newClaw()
	case LightningBreath:
		return hp.newLightningBreath()
	case Screech:
		// Demoralizing Screech: a single melee hit off the client row (rank 4 24579: 24-42, 20 focus,
		// 10 sec cooldown). 24582 is the learn spell. The attack power reduction is left out.
		return hp.newPetStrike(spellData.DemoralizingScreechTriggered.Highest())
	case ScorpidPoison:
		return hp.newScorpidPoison()
	case SavageRend:
		return hp.newPetBleed(spellData.SavageRendTriggered.Highest())
	case TendonRip:
		return hp.newPetBleed(spellData.TendonRipTriggered.Highest())
	case Web:
		return hp.newPetBleed(spellData.WebTriggered.Highest())
	case Pinch:
		return hp.newPetStrike(spellData.PinchTriggered.Highest())
	case Dismember:
		return hp.newPetStrike(spellData.DismemberTriggered.Highest())
	case Mine:
		return hp.newPetStrike(spellData.MineTriggered.Highest())
	case DustCloud:
		return hp.newDustCloud()
	case Thunderstomp:
		return hp.newThunderstomp()
	case Swipe:
		return hp.newSwipe()
	case FuriousHowl:
		if !WolfCastsFuriousHowl {
			return nil
		}
		return hp.newFuriousHowl()
	case Unknown:
		return nil
	default:
		panic("Invalid pet ability type")
	}
}

func (hp *HunterPet) newBite() *core.Spell {
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 17261},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ClassSpellMask: HunterPetDamage,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		MaxRange:       core.MaxMeleeRange,

		FocusCost: core.FocusCostOptions{
			Cost: 35,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: 10 * time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, sim.Roll(81, 99), spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

func (hp *HunterPet) newClaw() *core.Spell {
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 3009},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ClassSpellMask: HunterPetDamage,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		MaxRange:       core.MaxMeleeRange,

		FocusCost: core.FocusCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, sim.Roll(43, 59), spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

// Beta client: every rank lower than Classic's, and no more growth per level. Rank 6 is 86-98.
func (hp *HunterPet) newLightningBreath() *core.Spell {
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 25012},
		SpellSchool:    core.SpellSchoolNature,
		DefenseType:    core.DefenseTypeMagic,
		ClassSpellMask: HunterPetDamage,
		ProcMask:       core.ProcMaskSpellDamage,
		MaxRange:       20,

		FocusCost: core.FocusCostOptions{
			Cost: 50,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, sim.Roll(86, 98), spell.OutcomeMagicHitAndCrit)
		},
	})
}

// Beta client: 5 a tick at rank 4, down from 8, every 2 sec for 10 sec, stacking to 5 (24587 MaxStack).
func (hp *HunterPet) newScorpidPoison() *core.Spell {
	const baseDamageTick = 5.0

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 24587},
		SpellSchool:    core.SpellSchoolNature,
		DefenseType:    core.DefenseTypeMelee,
		ClassSpellMask: HunterPetDamage,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagPassiveSpell | core.SpellFlagPoison,
		MaxRange:       core.MaxMeleeRange,

		FocusCost: core.FocusCostOptions{
			Cost: 30,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: time.Second * 4,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Scorpid Poison",
				MaxStacks: 5,
				Duration:  time.Second * 10,
			},
			NumberOfTicks: 5,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, baseDamageTick*float64(dot.GetStacks()))
			},
			// 24587 carries Periodic Can Crit: the ticks roll the pet's melee crit.
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, spellData.ScorpidPoisonTriggered.Rank(4).TickOutcomeHitRolled(dot))
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if !result.Landed() {
				return
			}

			// Apply would wipe the stacks, so a landed poison on a poisoned target refreshes it instead.
			dot := spell.Dot(target)
			if dot.IsActive() {
				dot.Refresh(sim)
				if dot.GetStacks() < dot.MaxStacks {
					dot.AddStack(sim)
				}
			} else {
				dot.Apply(sim)
				dot.SetStacks(sim, 1)
			}
			dot.TakeSnapshot(sim)
		},
	})
}

// Savage Rend (Raptor, skill line 217) and Tendon Rip (Hyena, 654) are new in Forever: a melee hit that
// lands a bleed, everything read off the client row (rank 5: Savage Rend 26 every 3 sec for 18 sec, 50
// focus, 1 min; Tendon Rip 20 every 3 sec for 9 sec, 25 focus, 30 sec). Beta logs (Tynman's and
// Consumer's raptors, foreverlogs 2650/2669/2673/2674) put a Savage Rend tick at 6.4-7.0 at rank 1,
// which is the 5 base times the pet's happiness and Raptor damage scalars, so no attack power share.
// The 5% more bleed damage Savage Rend adds and Tendon Rip's snare are left out.
//
// Web (Spider, skill line 203) is new in Forever too and rides the same code: Nature damage every sec
// for 4 sec off a ranged-defense hit (rank 5: 13 a tick, 20 focus, 40 sec cooldown); the root is left
// out, bosses are immune. Beta logs (Consumer's spider, foreverlogs 2682/2683) tick rank 1's 3 at 4,
// four ticks a cast, one crit in 28: the base times the pet's happiness and Spider damage scalars.
func (hp *HunterPet) newPetBleed(rank *spelldata.Spell) *core.Spell {
	tick := rank.PeriodicEffect()
	tickLength := tick.Period()

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ClassSpellMask: HunterPetDamage,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		MaxRange:       core.MaxMeleeRange,

		FocusCost: core.FocusCostOptions{
			Cost: int32(rank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: rank.Cooldown(),
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: rank.Name,
			},
			NumberOfTicks: int32(rank.Duration() / tickLength),
			TickLength:    tickLength,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, tick.Average(core.CharacterLevel))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, rank.TickOutcomeHitRolled(dot))
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			outcome := spell.OutcomeMeleeSpecialHit
			if spell.DefenseType == core.DefenseTypeRanged {
				outcome = spell.OutcomeRangedHit
			}
			if spell.CalcAndDealOutcome(sim, target, outcome).Landed() {
				spell.Dot(target).Apply(sim)
			}
		},
	})
}

// Pinch (Crab, skill line 214), Dismember (Crocolisk, 212) and Mine! (Owl, 655) are new in Forever: a
// single melee hit with the damage, focus cost and cooldown read off the client row (rank 5: Pinch 95
// for 50 focus on 30 sec, Dismember 54 for 35 focus on 6 sec, Mine! 44 for 20 focus on 1 min). Pinch's
// snare, Dismember's healing reduction and Mine!'s disarm are left out. Beta logs: Consumer's crab (foreverlogs 2674) landed rank 1 Pinch (20)
// 17 times at ~17 a hit through level 20 mob armor, so no attack power share.
func (hp *HunterPet) newPetStrike(rank *spelldata.Spell) *core.Spell {
	damage := rank.DamageEffect()

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ClassSpellMask: HunterPetDamage,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		MaxRange:       core.MaxMeleeRange,

		FocusCost: core.FocusCostOptions{
			Cost: int32(rank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: rank.Cooldown(),
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, damage.Roll(sim, core.CharacterLevel), spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

// DustCloudYieldsToSunderAndExpose is the switch for Dust Cloud against the other physical armor
// reductions. Players report (Redfall, MythicSim Discord, October 2026) that in Forever Sunder
// Armor, Expose Armor and Dust Cloud don't stack, only the strongest applies; no combat log confirms
// it yet. With it set (the default) the Tallstrider holds Dust Cloud while Sunder Armor or Expose
// Armor is on the target. Both outweigh its 505 (Sunder 450 a stack, 2250 at five; Expose 2250), so
// the target never carries Dust Cloud's armor on top of them and the pet keeps its focus for Bite.
// Faerie Fire and Curse of Recklessness are magic and don't compete with it. Set it false to let
// Dust Cloud stack with everything, as before.
//
// Dust Cloud doesn't join Sunder and Expose's single-aura category itself: Sunder's raid copy ramps
// from 0 stacks, so there a Dust Cloud would shut Sunder out for the fight.
var DustCloudYieldsToSunderAndExpose = true

// Dust Cloud is new in Forever and the Tallstrider's alone (client SkillLineAbility, skill line 218):
// the target's armor down by the client row's amount for 30 sec, 10 focus, no cooldown (rank 5 505,
// rank 1 at level 12 65). The pet recasts it when it falls off, unless Sunder Armor or Expose Armor
// is up (DustCloudYieldsToSunderAndExpose). Faerie Fire and Curse of Recklessness stack with it.
func (hp *HunterPet) newDustCloud() *core.Spell {
	rank := spellData.DustCloudTriggered.Highest()
	auras := hp.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		aura := target.GetOrRegisterAura(core.Aura{
			Label:    rank.Name,
			ActionID: core.ActionID{SpellID: rank.ID},
			Duration: rank.Duration(),
		})
		spelldata.ParseEffects(nil, aura, rank, spelldata.Level(core.CharacterLevel), spelldata.IgnoreStacks(),
			spelldata.Exclusive("DustCloud", true))
		return aura
	})

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: rank.ID},
		SpellSchool: rank.SpellSchool(),
		DefenseType: rank.DefenseTypeCore(),
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,
		MaxRange:    core.MaxMeleeRange,

		FocusCost: core.FocusCostOptions{
			Cost: int32(rank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
		},

		ThreatMultiplier: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			if DustCloudYieldsToSunderAndExpose && target.GetExclusiveEffectCategory(buffs.SunderArmorCategory).AnyActive() {
				return false
			}
			return hp.IsEnabled() && !auras.Get(target).IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit).Landed() {
				auras.Get(target).Activate(sim)
			}
		},

		RelatedAuraArrays: auras.ToMap(),
	})
}

// Thunderstomp is the Gorilla's (Classic ranks 1-3 at 30/40/50; rank 4 at 60 is new in Forever): Nature
// damage to up to the row's 4 enemies around the pet, read off the client row (rank 4 132 +-7%, so
// 122-142, 60 focus, 1 min cooldown). Magic, so it rolls spell hit and crit and no attack power share.
func (hp *HunterPet) newThunderstomp() *core.Spell {
	rank := spellData.ThunderstompTriggered.Highest()
	damage := rank.DamageEffect()

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ClassSpellMask: HunterPetDamage,
		ProcMask:       core.ProcMaskSpellDamage,
		MaxRange:       core.MaxMeleeRange,

		FocusCost: core.FocusCostOptions{
			Cost: int32(rank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: rank.Cooldown(),
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealCleaveDamageWithVariance(sim, target, int32(rank.MaxTargets), spell.OutcomeMagicHitAndCrit,
				func(sim *core.Simulation, _ *core.Spell) float64 { return damage.Roll(sim, core.CharacterLevel) })
		},
	})
}

// Swipe is new in Forever and the Bear's alone: a melee hit on up to 3 enemies (client row ChainTargets 3;
// rank 5 1264502 21 +-7%, so 20-22, 20 focus, 5 sec cooldown; Wowhead env 16 agrees). Claw's 43-59 for 25 focus
// beats it on one or two targets (sim: 2 targets 193.9 -> 193.2 hunter DPS with Swipe, 3 targets 194.8 ->
// 197.2), so the bear swipes only when 3 or more enemies are up.
func (hp *HunterPet) newSwipe() *core.Spell {
	rank := spellData.SwipeTriggered.Highest()
	damage := rank.DamageEffect()

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ClassSpellMask: HunterPetDamage,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		MaxRange:       core.MaxMeleeRange,

		FocusCost: core.FocusCostOptions{
			Cost: int32(rank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: rank.Cooldown(),
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled() && sim.Environment.ActiveTargetCount() >= 3
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealCleaveDamageWithVariance(sim, target, int32(damage.ChainTargets), spell.OutcomeMeleeSpecialHitAndCrit,
				func(sim *core.Simulation, _ *core.Spell) float64 { return damage.Roll(sim, core.CharacterLevel) })
		},
	})
}

// WolfCastsFuriousHowl (MythicSim patch 119) gives the Wolf its Furious Howl. Without it the Wolf only
// Bites, and its party buff, the reason to pick a Wolf for a melee Hunter (Kragtor, Discord, 12 October
// 2026), is missing. Set it false for the Bite-only Wolf.
var WolfCastsFuriousHowl = true

// Furious Howl is the Wolf's: in Forever an attack power aura on the party around the pet (client row
// 24597, rank 4: A_MOD_ATTACK_POWER 82 + 0.4 a level from 56, so 84 at 60; implicit target 20, party
// within 15 yards; 60 sec duration, 30 sec cooldown, 60 focus). The Sep 24 beta notes cut it 40%
// ("Furious Howl Attack Power bonus reduced 40%"); 70334's row is the cut value. Melee attack power only,
// so it pays a melee Hunter and the Wolf itself, not Auto Shot. The Wolf howls when the buff is not up on
// it; the Hunter takes it while within 15 yards of the pet, which stands at the target. Other party
// members are not simulated.
func (hp *HunterPet) newFuriousHowl() *core.Spell {
	rank := spellData.FuriousHowlTriggered.Highest()
	radius := 15.0
	if e := rank.Effect(dbcenums.A_MOD_ATTACK_POWER, 0); e != nil && e.RadiusMax > 0 {
		radius = float64(e.RadiusMax)
	}
	register := func(unit *core.Unit) *core.Aura {
		aura := unit.GetOrRegisterAura(core.Aura{
			Label:    rank.Name,
			ActionID: core.ActionID{SpellID: rank.ID},
			Duration: rank.Duration(),
		})
		spelldata.ParseEffects(nil, aura, rank, spelldata.Level(core.CharacterLevel))
		return aura
	}
	petAura := register(&hp.Unit)
	ownerAura := register(&hp.hunterOwner.Unit)

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: rank.ID},
		SpellSchool: rank.SpellSchool(),
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,

		FocusCost: core.FocusCostOptions{
			Cost: int32(rank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: rank.Cooldown(),
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled() && !petAura.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			petAura.Activate(sim)
			if hp.hunterOwner.DistanceFromTarget <= radius {
				ownerAura.Activate(sim)
			}
		},
	})
}
