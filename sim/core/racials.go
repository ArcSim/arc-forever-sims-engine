package core

import (
	"os"
	"slices"
	"time"

	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
)

// Racials follow the WoW Forever client data (build 1.60.1.70334) with Blizzard's racial tuning pass of 26 October
// 2026 and its 27 October additions on top: there is no weapon or Bow/Gun/Thrown Specialization, Blood Fury and
// Berserking are flat 10%, Night Elves have Elune's Light and Dwarves Big Game Hunter. Skyborne is a Forever-only race.
func applyRaceEffects(agent Agent) {
	character := agent.GetCharacter()

	switch character.Race {
	case proto.Race_RaceDwarf:
		character.AddStat(stats.FrostResistance, 10)
		// Artisan's Precision: +1% hit with spells and attacks (replaced Earth Fury's critical damage)
		character.AddStat(stats.MeleeHit, 1*MeleeHitRatingPerHitChance)
		character.AddStat(stats.SpellHit, 1*SpellHitRatingPerHitChance)
		character.applyBeastDamageRacial()

		actionID := ActionID{SpellID: 20594}

		statDep := character.NewDynamicMultiplyStat(stats.Armor, 1.1)
		stoneFormAura := character.NewTemporaryStatsAuraWrapped("Stoneform", actionID, stats.Stats{}, time.Second*8, func(aura *Aura) {
			aura.ApplyOnGain(func(aura *Aura, sim *Simulation) {
				aura.Unit.EnableDynamicStatDep(sim, statDep)
			})
			aura.ApplyOnExpire(func(aura *Aura, sim *Simulation) {
				aura.Unit.DisableDynamicStatDep(sim, statDep)
			})
		})

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				stoneFormAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeSurvival,
			ShouldActivate: func(s *Simulation, c *Character) bool {
				return false
			},
		})
	case proto.Race_RaceGnome:
		character.AddStat(stats.ArcaneResistance, 10)
		if character.HasRageBar() {
			character.AddMaxRage(5)
		} else {
			character.MultiplyStat(stats.Mana, 1.05)
		}
		makeEurekaCooldown(character)
	case proto.Race_RaceHuman:
		// The Human Spirit: +10% Spirit
		character.MultiplyStat(stats.Spirit, 1.10)
		// Ruthless: +1% crit with any weapon
		character.addCritRacial(1)
	case proto.Race_RaceNightElf:
		character.AddStat(stats.NatureResistance, 10)
		character.AddStat(stats.Dodge, 1)

		// Quickness: +2% run speed, stacking with other speed bonuses such as Strider Kick
		quicknessID := ActionID{SpellID: 20582}
		character.RegisterAura(Aura{
			Label:    "Quickness",
			ActionID: quicknessID,
			Duration: NeverExpires,
			OnReset: func(aura *Aura, sim *Simulation) {
				aura.Activate(sim)
			},
			OnGain: func(aura *Aura, sim *Simulation) {
				aura.Unit.MultiplyStackingMoveSpeed(1.02)
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				aura.Unit.MultiplyStackingMoveSpeed(1 / 1.02)
			},
		})

		actionID := ActionID{SpellID: 1259799}
		eluneAura := character.NewTemporaryStatsAura("Elune's Light", actionID, stats.Stats{
			stats.MeleeCrit: 10 * CritRatingPerCritChance,
			stats.SpellCrit: 10 * SpellCritRatingPerCritChance,
		}, time.Second*15)

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete | SpellFlagAPL,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				eluneAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeDPS,
		})
	case proto.Race_RaceOrc:
		// Command: Hunter and Warlock pets deal 5% more damage (Savage's crit is gone)
		if character.Class == proto.Class_ClassHunter || character.Class == proto.Class_ClassWarlock {
			for _, pet := range character.Pets {
				pet.PseudoStats.DamageDealtMultiplier *= 1.05
			}
		}

		actionID := ActionID{SpellID: 20572}
		// Blood Fury: +10% attack power and spell power (Forever 20572 raises both); single-school power ("+X fire
		// damage") cannot take a multiplier in the engine and stays as is
		var bloodFuryDeps []*stats.StatDependency
		for _, stat := range []stats.Stat{stats.AttackPower, stats.RangedAttackPower, stats.SpellPower, stats.SpellDamage, stats.HealingPower} {
			bloodFuryDeps = append(bloodFuryDeps, character.NewDynamicMultiplyStat(stat, 1.1))
		}
		bloodFuryAura := character.RegisterAura(Aura{
			Label:    "Blood Fury",
			ActionID: actionID,
			Duration: time.Second * 15,
			OnGain: func(aura *Aura, sim *Simulation) {
				for _, dep := range bloodFuryDeps {
					aura.Unit.EnableDynamicStatDep(sim, dep)
				}
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				for _, dep := range bloodFuryDeps {
					aura.Unit.DisableDynamicStatDep(sim, dep)
				}
			},
		})

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete | SpellFlagAPL,
			Cast: CastConfig{
				DefaultCast: Cast{GCD: 0},
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				bloodFuryAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeDPS,
		})
	case proto.Race_RaceTauren:
		character.AddStat(stats.NatureResistance, 10)
		// Endurance: +5% Health and nothing else. Plainsrunning (run speed stacks) is not modelled: how its stacks build
		// is not known yet
		character.MultiplyStat(stats.Health, 1.05)
	case proto.Race_RaceTroll:
		character.applyBeastDamageRacial()
		makeBerserkingCooldown(character)
	case proto.Race_RaceUndead:
		character.AddStat(stats.ShadowResistance, 10)
		applyTouchOfTheGrave(character)
	case proto.Race_RaceSkyborneAlliance, proto.Race_RaceSkyborneHorde:
		// Wind Blessed: +2% melee, ranged and casting speed
		character.PseudoStats.MeleeSpeedMultiplier *= 1.02
		character.PseudoStats.RangedSpeedMultiplier *= 1.02
		character.PseudoStats.CastSpeedMultiplier *= 1.02
		// Elemental Insight
		character.applyMobTypeDamageRacial(proto.MobType_MobTypeElemental)
		character.applySkyborneBlessing()
	}
}

// SkysightStacks: AB_SKYSIGHT_STACKS=1 lets Skysight stack with other speed bonuses, for comparing; only Night Elf
// Quickness is known to stack, so by default only the strongest speed bonus counts.
var SkysightStacks = os.Getenv("AB_SKYSIGHT_STACKS") != ""

// Skysight (Horde Skyborne): +15% run speed for 15 s (a regular speed bonus: only the strongest counts). Read Leyline
// (Alliance Skyborne): +50% health and mana regeneration for 15 s. Each is on a 2 min cooldown and the GCD, with a cast
// (client data: Skysight 0.5 s, Read Leyline 2 s), so neither goes off while moving. Charged and Attuned, which raise
// later casts after one at an elemental convergence or a leyline, are removed on entering a raid, so they never count.
// Neither is cast in combat on its own: only a rotation that asks for it casts it.
func (character *Character) applySkyborneBlessing() {
	horde := character.Race == proto.Race_RaceSkyborneHorde
	spellID, buffID, label, castTime := int32(1259705), int32(1270842), "Read Leyline", 2*time.Second
	if horde {
		spellID, buffID, label, castTime = 1259686, 1259688, "Skysight", 500*time.Millisecond
	}
	speed, regen := 1.15, 0.5

	buffActionID := ActionID{SpellID: buffID}
	aura := character.RegisterAura(Aura{
		Label:    label,
		ActionID: buffActionID,
		Duration: time.Second * 15,
		OnGain: func(aura *Aura, sim *Simulation) {
			if horde && SkysightStacks {
				aura.Unit.MultiplyStackingMoveSpeed(speed)
			} else if horde {
				aura.Unit.AddMoveSpeedModifier(&buffActionID, speed)
			} else {
				aura.Unit.PseudoStats.SpiritRegenMultiplier += regen
				aura.Unit.UpdateManaRegenRates()
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			if horde && SkysightStacks {
				aura.Unit.MultiplyStackingMoveSpeed(1 / speed)
			} else if horde {
				aura.Unit.RemoveMoveSpeedModifier(&buffActionID)
			} else {
				aura.Unit.PseudoStats.SpiritRegenMultiplier -= regen
				aura.Unit.UpdateManaRegenRates()
			}
		},
	})

	character.RegisterSpell(SpellConfig{
		ActionID: ActionID{SpellID: spellID},
		Flags:    SpellFlagNoOnCastComplete | SpellFlagAPL,
		Cast: CastConfig{
			DefaultCast: Cast{
				GCD:      GCDDefault,
				CastTime: castTime,
			},
			CD: Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Minute * 2,
			},
		},
		ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
			aura.Activate(sim)
		},
	})

}

// +pct% crit with all spells and attacks.
func (character *Character) addCritRacial(pct float64) {
	character.AddStat(stats.MeleeCrit, pct*CritRatingPerCritChance)
	character.AddStat(stats.SpellCrit, pct*SpellCritRatingPerCritChance)
}

// Beast Slaying (Troll) and Big Game Hunter (Dwarf): +5% damage dealt to Beasts.
func (character *Character) applyBeastDamageRacial() {
	character.applyMobTypeDamageRacial(proto.MobType_MobTypeBeast)
}

func (character *Character) applyMobTypeDamageRacial(mobType proto.MobType) {
	character.Env.RegisterPostFinalizeEffect(func() {
		for _, t := range character.Env.Encounter.Targets {
			if t.MobType == mobType {
				for _, at := range character.AttackTables[t.UnitIndex] {
					at.DamageDealtMultiplier *= 1.05
				}
			}
		}
	})
}

// Berserking: +10% attack and casting speed for 15 s (27 October additions, was 10 s; client 70634 still says 10), no
// cost, off the GCD. Spell 20554 in Forever (the Classic 26297 is not in its client data).
func makeBerserkingCooldown(character *Character) {
	actionID := ActionID{SpellID: 20554}

	berserkingAura := character.RegisterAura(Aura{
		Label:    "Berserking",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.MultiplyCastSpeed(1.1)
			character.MultiplyAttackSpeed(sim, 1.1)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.MultiplyCastSpeed(1 / 1.1)
			character.MultiplyAttackSpeed(sim, 1/1.1)
		},
	})

	berserkingSpell := character.RegisterSpell(SpellConfig{
		ActionID: actionID,
		Flags:    SpellFlagNoOnCastComplete | SpellFlagAPL,
		Cast: CastConfig{
			CD: Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Minute * 3,
			},
		},
		ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
			berserkingAura.Activate(sim)
		},
	})

	character.AddMajorCooldown(MajorCooldown{
		Spell: berserkingSpell,
		Type:  CooldownTypeDPS,
	})
}

// Eureka! (Forever): for 8 s, damaging abilities cost 10% less and deal 10% more damage, 2 min cooldown. For Rage and
// Energy users only the non-periodic abilities get the damage (one spell per resource in the client data).
func makeEurekaCooldown(character *Character) {
	actionID := ActionID{SpellID: 1259817}
	periodicToo := true
	if character.HasRageBar() {
		actionID, periodicToo = ActionID{SpellID: 1259813}, false
	} else if character.HasEnergyBar() {
		actionID, periodicToo = ActionID{SpellID: 1259812}, false
	}

	// the class's own damaging abilities: castable, paid for, and dealing melee, ranged or spell damage
	var affected []*Spell
	character.OnSpellRegistered(func(spell *Spell) {
		if !spell.Flags.Matches(SpellFlagAPL) || spell.Cost == nil || !spell.ProcMask.Matches(ProcMaskMeleeOrRangedSpecial|ProcMaskSpellDamage) {
			return
		}
		if !periodicToo && (len(spell.dots) > 0 || spell.aoeDot != nil) {
			return
		}
		affected = append(affected, spell)
	})

	eurekaAura := character.RegisterAura(Aura{
		Label:    "Eureka!",
		ActionID: actionID,
		Duration: time.Second * 8,
		OnGain: func(aura *Aura, sim *Simulation) {
			for _, spell := range affected {
				spell.DamageMultiplierAdditive += 0.1
				spell.Cost.Multiplier -= 10
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			for _, spell := range affected {
				spell.DamageMultiplierAdditive -= 0.1
				spell.Cost.Multiplier += 10
			}
		},
	})

	spell := character.RegisterSpell(SpellConfig{
		ActionID: actionID,
		Flags:    SpellFlagNoOnCastComplete | SpellFlagAPL,
		Cast: CastConfig{
			CD: Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Minute * 2,
			},
		},
		ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
			eurekaAura.Activate(sim)
		},
	})

	character.AddMajorCooldown(MajorCooldown{
		Spell: spell,
		Type:  CooldownTypeDPS,
	})
}

// Touch of the Grave (Forever): spells and attacks that land have a 30% chance (20% for Paladins, Rogues and
// Warriors) to drain 5% of your maximum health from the target, 15 s cooldown. The drain ignores damage bonuses on
// either side (Forever patch notes); the heal on yourself is not modelled.
func touchOfTheGraveChance(class proto.Class) float64 {
	switch class {
	case proto.Class_ClassPaladin, proto.Class_ClassRogue, proto.Class_ClassWarrior:
		return 0.2
	}
	return 0.3
}

func applyTouchOfTheGrave(character *Character) {
	chance := touchOfTheGraveChance(character.Class)

	drain := character.RegisterSpell(SpellConfig{
		ActionID:    ActionID{SpellID: 1260189},
		SpellSchool: SpellSchoolShadow,
		DefenseType: DefenseTypeMagic,
		ProcMask:    ProcMaskEmpty,
		Flags:       SpellFlagNoOnCastComplete | SpellFlagPassiveSpell | SpellFlagIgnoreModifiers,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			spell.CalcAndDealDamage(sim, target, 0.05*spell.Unit.MaxHealth(), spell.OutcomeAlwaysHit)
		},
	})

	MakeProcTriggerAura(&character.Unit, ProcTrigger{
		Name:       "Touch of the Grave",
		Callback:   CallbackOnSpellHitDealt,
		ProcMask:   ProcMaskMeleeOrRanged | ProcMaskSpellDamage,
		Outcome:    OutcomeLanded,
		ProcChance: chance,
		ICD:        time.Second * 15,
		Handler: func(sim *Simulation, _ *Spell, result *SpellResult) {
			drain.Cast(sim, result.Target)
		},
	})
}

func (character *Character) GetFaction() proto.Faction {
	if slices.Contains([]proto.Race{proto.Race_RaceHuman, proto.Race_RaceDwarf, proto.Race_RaceGnome, proto.Race_RaceNightElf, proto.Race_RaceSkyborneAlliance}, character.Race) {
		return proto.Faction_Alliance
	} else if slices.Contains([]proto.Race{proto.Race_RaceOrc, proto.Race_RaceTroll, proto.Race_RaceTauren, proto.Race_RaceUndead, proto.Race_RaceSkyborneHorde}, character.Race) {
		return proto.Faction_Horde
	}
	return proto.Faction_Unknown
}
