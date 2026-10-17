package hunter

import (
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
)

type PetAbilityType int

// Pet AI doesn't use abilities immediately, so model this with a 1.6s GCD.
const PetGCD = time.Millisecond * 1600

const (
	Unknown PetAbilityType = iota
	Bite
	Claw
	LightningBreath
	ScorpidPoison
	SavageRend
	FuriousHowl
	Thunderstomp
	Dismember
	Pinch
	TendonRip
	Web
	Swipe
	Mine
	Charge
	DemoralizingScreech
	LavaBreath
	DustCloud
	TrickstersDance
	SonicBlast
)

// option is the ability's value in the "abilities the pet does not use" setting.
func (a PetAbilityType) option() proto.Hunter_Options_PetAbility {
	return map[PetAbilityType]proto.Hunter_Options_PetAbility{
		Bite:                proto.Hunter_Options_PetAbilityBite,
		Claw:                proto.Hunter_Options_PetAbilityClaw,
		LightningBreath:     proto.Hunter_Options_PetAbilityLightningBreath,
		ScorpidPoison:       proto.Hunter_Options_PetAbilityScorpidPoison,
		SavageRend:          proto.Hunter_Options_PetAbilitySavageRend,
		FuriousHowl:         proto.Hunter_Options_PetAbilityFuriousHowl,
		Thunderstomp:        proto.Hunter_Options_PetAbilityThunderstomp,
		Dismember:           proto.Hunter_Options_PetAbilityDismember,
		Pinch:               proto.Hunter_Options_PetAbilityPinch,
		TendonRip:           proto.Hunter_Options_PetAbilityTendonRip,
		Web:                 proto.Hunter_Options_PetAbilityWeb,
		Swipe:               proto.Hunter_Options_PetAbilitySwipe,
		Mine:                proto.Hunter_Options_PetAbilityMine,
		Charge:              proto.Hunter_Options_PetAbilityCharge,
		DemoralizingScreech: proto.Hunter_Options_PetAbilityDemoralizingScreech,
		LavaBreath:          proto.Hunter_Options_PetAbilityLavaBreath,
		DustCloud:           proto.Hunter_Options_PetAbilityDustCloud,
		TrickstersDance:     proto.Hunter_Options_PetAbilityTrickstersDance,
		SonicBlast:          proto.Hunter_Options_PetAbilitySonicBlast,
	}[a]
}

// Top ranks per WoW Forever client data (1.60.1, with hotfixes).
func (hp *HunterPet) newPetAbility(abilityType PetAbilityType) *petAbility {
	switch abilityType {
	case Bite:
		return &petAbility{spell: hp.newBite()}
	case Claw:
		return &petAbility{spell: hp.newClaw(), filler: true}
	case LightningBreath:
		return &petAbility{spell: hp.newLightningBreath(), filler: true}
	case ScorpidPoison:
		spell := hp.newScorpidPoison()
		// keep 5 stacks up, refresh near the end
		return &petAbility{spell: spell, wanted: func(sim *core.Simulation) bool {
			dot := spell.Dot(hp.CurrentTarget)
			return dot.GetStacks() < dot.MaxStacks || dot.RemainingDuration(sim) < time.Second*3
		}}
	case SavageRend:
		return &petAbility{spell: hp.newPetBleed(1265069, "Savage Rend", core.SpellSchoolPhysical, 50, time.Second*60, 26, 6, time.Second*3)}
	case TendonRip:
		return &petAbility{spell: hp.newPetBleed(1265042, "Tendon Rip", core.SpellSchoolPhysical, 25, time.Second*30, 20, 3, time.Second*3)}
	case Web:
		return &petAbility{spell: hp.newPetBleed(1265883, "Web", core.SpellSchoolNature, 20, time.Second*40, 13, 4, time.Second)}
	case SonicBlast:
		// rank 5 (1264482): 77-91 Nature damage, 80 focus, 1 min cooldown (its cast slow does not matter here)
		return &petAbility{spell: hp.newPetStrike(1264482, core.SpellSchoolNature, 80, time.Minute, 77, 91)}
	case Thunderstomp:
		return &petAbility{spell: hp.newPetStrike(1264455, core.SpellSchoolNature, 60, time.Second*60, 122, 142)}
	case Pinch:
		return &petAbility{spell: hp.newPetStrike(1264742, core.SpellSchoolPhysical, 50, time.Second*30, 88, 102)}
	case Mine:
		return &petAbility{spell: hp.newPetStrike(1265058, core.SpellSchoolPhysical, 20, time.Second*90, 41, 47)}
	case Dismember:
		return &petAbility{spell: hp.newPetStrike(1264933, core.SpellSchoolPhysical, 35, time.Second*15, 50, 58)}
	case Swipe:
		return &petAbility{spell: hp.newPetStrike(1264502, core.SpellSchoolPhysical, 20, time.Second*5, 20, 22), filler: true}
	case DemoralizingScreech:
		return &petAbility{spell: hp.newPetStrike(24579, core.SpellSchoolPhysical, 20, time.Second*10, 24, 42), filler: true}
	case LavaBreath:
		// rank 2 (444681): 50 focus, no cooldown, 106 + 0.442 per level up to 54 (about 130), up to 2 targets
		return &petAbility{spell: hp.newPetStrike(444681, core.SpellSchoolFire, 50, 0, 123, 137), filler: true}
	case FuriousHowl:
		return &petAbility{spell: hp.newFuriousHowl()}
	case TrickstersDance:
		return &petAbility{spell: hp.newTrickstersDance()}
	case Charge:
		return &petAbility{spell: hp.newCharge()}
	case DustCloud:
		spell, aura := hp.newDustCloud()
		return &petAbility{spell: spell, wanted: func(sim *core.Simulation) bool {
			debuff := aura.Get(hp.CurrentTarget)
			return !debuff.IsActive() || debuff.RemainingDuration(sim) < time.Second*3
		}}
	}
	return nil
}

func (hp *HunterPet) petCast(cost float64, cd time.Duration) core.CastConfig {
	cast := core.CastConfig{
		DefaultCast: core.Cast{GCD: PetGCD},
		IgnoreHaste: true,
	}
	if cd > 0 {
		cast.CD = core.Cooldown{Timer: hp.NewTimer(), Duration: cd}
	}
	return cast
}

// A single hit: physical abilities are melee specials, other schools are spells.
func (hp *HunterPet) newPetStrike(spellID int32, school core.SpellSchool, cost float64, cd time.Duration, minDamage, maxDamage float64) *core.Spell {
	physical := school == core.SpellSchoolPhysical
	config := core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: school,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,

		FocusCost: core.FocusCostOptions{Cost: cost},
		Cast:      hp.petCast(cost, cd),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			outcome := spell.OutcomeMagicHitAndCrit
			if physical {
				outcome = spell.OutcomeMeleeSpecialHitAndCrit
			}
			spell.CalcAndDealDamage(sim, target, sim.Roll(minDamage, maxDamage), outcome)
		},
	}
	if physical {
		config.DefenseType = core.DefenseTypeMelee
		config.ProcMask = core.ProcMaskMeleeMHSpecial
		config.Flags = core.SpellFlagMeleeMetrics
	}
	return hp.RegisterSpell(config)
}

// Damage over time: an attack roll to land, then ticks that cannot crit.
func (hp *HunterPet) newPetBleed(spellID int32, label string, school core.SpellSchool, cost float64, cd time.Duration, tick float64, ticks int32, tickLength time.Duration) *core.Spell {
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: school,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		FocusCost: core.FocusCostOptions{Cost: cost},
		Cast:      hp.petCast(cost, cd),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura:          core.Aura{Label: label},
			NumberOfTicks: ticks,
			TickLength:    tickLength,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.Snapshot(target, tick, false)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
		},
	})
}

// Trickster's Dance (fox, 1310612): 30% melee haste for 12 s on a 3 min cooldown, 10 focus (the dodge half does not matter here).
func (hp *HunterPet) newTrickstersDance() *core.Spell {
	actionID := core.ActionID{SpellID: 1310612}
	buff := hp.RegisterAura(core.Aura{
		Label:    "Trickster's Dance",
		ActionID: actionID,
		Duration: time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyMeleeSpeed(sim, 1.3)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyMeleeSpeed(sim, 1/1.3)
		},
	})
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:  actionID,
		FocusCost: core.FocusCostOptions{Cost: 10},
		Cast:      hp.petCast(10, time.Minute*3),
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			buff.Activate(sim)
		},
	})
}

// Furious Howl: +melee attack power for the hunter (the howl does not reach the rest of the party).
func (hp *HunterPet) newFuriousHowl() *core.Spell {
	// Forever's value is the client's 83.6 less 40%
	const ap = 83.6 * 0.6
	actionID := core.ActionID{SpellID: 24597}
	owner := hp.hunterOwner
	buff := owner.RegisterAura(core.Aura{
		Label:    "Furious Howl",
		ActionID: actionID,
		Duration: time.Second * 60,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.AttackPower, ap)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.AttackPower, -ap)
		},
	})
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:  actionID,
		FocusCost: core.FocusCostOptions{Cost: 60},
		Cast:      hp.petCast(60, time.Second*30),
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			buff.Activate(sim)
		},
	})
}

// Charge: the boar opens with it (it needs room to charge), +520 attack power on its next attack.
func (hp *HunterPet) newCharge() *core.Spell {
	const ap = 520
	actionID := core.ActionID{SpellID: 27685}
	buff := hp.RegisterAura(core.Aura{
		Label:    "Boar Charge",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.AttackPower, ap)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.AttackPower, -ap)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMelee) {
				aura.Deactivate(sim)
			}
		},
	})
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:  actionID.WithTag(1),
		FocusCost: core.FocusCostOptions{Cost: 35},
		Cast:      hp.petCast(35, time.Second*25),
		ExtraCastCondition: func(sim *core.Simulation, _ *core.Unit) bool {
			return sim.CurrentTime < time.Second
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			buff.Activate(sim)
		},
	})
}

// the Faerie Fire armor group (core debuffs.go, common item_effects.go use the same name)
const minorArmorReductionCategory = "MinorArmorReduction"

// Dust Cloud: -505 armor on the target for 30 s. It does not stack with Faerie Fire or Curse of Recklessness (tested
// in game: one -505 debuff), so it only helps when both are missing.
func (hp *HunterPet) newDustCloud() (*core.Spell, core.AuraArray) {
	actionID := core.ActionID{SpellID: 1265904}
	debuffs := hp.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		aura := target.GetOrRegisterAura(core.Aura{
			Label:    "Dust Cloud",
			ActionID: actionID,
			Duration: time.Second * 30,
		})
		aura.NewExclusiveEffect(minorArmorReductionCategory, true, core.ExclusiveEffect{
			Priority: 505,
			OnGain: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
				ee.Aura.Unit.AddStatDynamic(sim, stats.Armor, -505)
			},
			OnExpire: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
				ee.Aura.Unit.AddStatDynamic(sim, stats.Armor, 505)
			},
		})
		return aura
	})
	spell := hp.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		FocusCost:   core.FocusCostOptions{Cost: 10},
		Cast:        hp.petCast(10, 0),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			debuffs.Get(target).Activate(sim)
		},
	})
	return spell, debuffs
}

func (hp *HunterPet) newClaw() *core.Spell {
	baseDamageMin := map[int32]float64{
		25: 16,
		40: 26,
		50: 35,
		60: 43,
	}[hp.Owner.Level]

	baseDamageMax := map[int32]float64{
		25: 22,
		40: 36,
		50: 49,
		60: 59,
	}[hp.Owner.Level]

	spellID := map[int32]int32{
		25: 16830,
		40: 16832,
		50: 3010,
		60: 3009,
	}[hp.Owner.Level]

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellCode:   SpellCode_HunterPetClaw,
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

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

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageMin, baseDamageMax)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

func (hp *HunterPet) newBite() *core.Spell {
	baseDamageMin := map[int32]float64{
		25: 31,
		40: 49,
		50: 66,
		60: 81,
	}[hp.Owner.Level]

	baseDamageMax := map[int32]float64{
		25: 37,
		40: 59,
		50: 80,
		60: 99,
	}[hp.Owner.Level]

	spellID := map[int32]int32{
		25: 17257,
		40: 17259,
		50: 17260,
		60: 17261,
	}[hp.Owner.Level]

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellCode:   SpellCode_HunterPetBite,
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

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

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageMin, baseDamageMax)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

func (hp *HunterPet) newLightningBreath() *core.Spell {
	baseDamageMin := map[int32]float64{
		25: 36,
		40: 36,
		50: 78,
		60: 86,
	}[hp.Owner.Level]

	baseDamageMax := map[int32]float64{
		25: 41,
		40: 41,
		50: 91,
		60: 98,
	}[hp.Owner.Level]

	spellID := map[int32]int32{
		25: 25009,
		40: 25009, // rank 4 not available in SoD Phase 2
		50: 25011,
		60: 25012,
	}[hp.Owner.Level]

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellCode:   SpellCode_HunterPetLightningBreath,
		SpellSchool: core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,

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
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageMin, baseDamageMax)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (hp *HunterPet) newScorpidPoison() *core.Spell {
	baseDamageTick := map[int32]float64{
		25: 3,
		40: 6,
		50: 6,
		60: 5,
	}[hp.Owner.Level]
	spellID := map[int32]int32{
		25: 24583,
		40: 24586,
		50: 24586,
		60: 24587,
	}[hp.Owner.Level]

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellCode:   SpellCode_HunterPetScorpidPoison,
		SpellSchool: core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagPassiveSpell | core.SpellFlagPoison,

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
				Label:     "ScorpidPoison",
				MaxStacks: 5,
				Duration:  time.Second * 10,
			},
			NumberOfTicks: 5,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, applyStack bool) {
				if !applyStack {
					return
				}

				// only the first stack snapshots the multiplier
				if dot.GetStacks() == 1 {
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType]
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable, true)
					dot.SnapshotBaseDamage = 0
				}

				dot.SnapshotBaseDamage += baseDamageTick
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if !result.Landed() {
				return
			}

			dot := spell.Dot(target)
			dot.ApplyOrRefresh(sim)
			if dot.GetStacks() < dot.MaxStacks {
				dot.AddStack(sim)
				dot.TakeSnapshot(sim, true)
			}
		},
	})
}
