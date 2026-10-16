package core

import (
	"math"
	"time"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

func applyRaceEffects(agent Agent) {
	character := agent.GetCharacter()

	// Blizzard beta development notes, 16 Oct build (us.forums.blizzard.com/en/wow/t/2360696 post 6).
	switch character.Race {
	case proto.Race_RaceDwarf:
		// Earth Fury replaced Mace Specialization: 2% increased critical damage and healing.
		character.PseudoStats.CritDamageMultiplier *= 1.02
		applyCreatureTypeSlaying(character, proto.MobType_MobTypeBeast)

		actionID := ActionID{SpellID: 20594}

		stoneFormAura := character.RegisterAura(Aura{
			Label:    "Stoneform",
			ActionID: actionID,
			Duration: time.Second * 8,
		}).AttachMultiplicativePseudoStatBuff(&character.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical], 0.9)

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				DefaultCast: Cast{
					GCD: GCDDefault,
				},
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				stoneFormAura.Activate(sim)
			},

			RelatedSelfBuff: stoneFormAura,
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeSurvival,
		})
	case proto.Race_RaceGnome:
		applyExpansiveMind(character)
		applyEureka(character)
	case proto.Race_RaceHuman:
		character.MultiplyStat(stats.Spirit, 1.05)
		// Ruthless replaced Sword Specialization: 2% critical strike chance, any weapon.
		character.AddStats(stats.Stats{stats.PhysicalCritPercent: 2, stats.SpellCritPercent: 2})
	case proto.Race_RaceNightElf:
		character.PseudoStats.BaseDodgeChance += 0.01
		character.PseudoStats.MovementSpeedMultiplier *= 1.02

		RegisterTemporaryStatsOnUseCD(character,
			"Elune's Light",
			stats.Stats{
				stats.PhysicalCritPercent: 10,
				stats.SpellCritPercent:    10,
			},
			time.Second*15,
			SpellConfig{
				ActionID: ActionID{SpellID: 1259799},
				Cast: CastConfig{
					CD: Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 3,
					},
				},
			})
	case proto.Race_RaceOrc:
		// Savage replaced Axe Specialization: 1% critical strike chance, any weapon.
		character.AddStats(stats.Stats{stats.PhysicalCritPercent: 1, stats.SpellCritPercent: 1})

		bloodFuryID := ActionID{SpellID: 20572}

		bloodFuryAura := character.NewTemporaryStatMultiplierAura(Aura{
			Label:    "Blood Fury",
			ActionID: bloodFuryID,
			Duration: time.Second * 15,
		}, []StatMultiplier{
			{Stat: stats.AttackPower, Multiplier: 1.1},
			{Stat: stats.RangedAttackPower, Multiplier: 1.1},
			{Stat: stats.SpellDamage, Multiplier: 1.1},
			{Stat: stats.HealingPower, Multiplier: 1.1},
		})

		bloodFuryCD := Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Minute * 3,
		}
		bloodFuryAura.Icd = &bloodFuryCD

		bloodFury := character.RegisterSpell(SpellConfig{
			ActionID: bloodFuryID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: bloodFuryCD,
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				bloodFuryAura.Activate(sim)
			},

			RelatedSelfBuff: bloodFuryAura.Aura,
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell:    bloodFury,
			Type:     CooldownTypeDPS,
			BuffAura: bloodFuryAura,
		})

		shatterCurseID := ActionID{SpellID: 1299026}

		shatterCurseAura := character.RegisterAura(Aura{
			Label:    "Shatter Curse",
			ActionID: shatterCurseID,
			Duration: time.Second * 8,
		})
		for _, school := range []stats.SchoolIndex{
			stats.SchoolIndexArcane,
			stats.SchoolIndexFire,
			stats.SchoolIndexFrost,
			stats.SchoolIndexHoly,
			stats.SchoolIndexNature,
			stats.SchoolIndexShadow,
		} {
			shatterCurseAura.AttachMultiplicativePseudoStatBuff(&character.PseudoStats.SchoolDamageTakenMultiplier[school], 0.95)
		}

		shatterCurse := character.RegisterSpell(SpellConfig{
			ActionID: shatterCurseID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				DefaultCast: Cast{
					GCD: GCDDefault,
				},
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				shatterCurseAura.Activate(sim)
			},

			RelatedSelfBuff: shatterCurseAura,
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: shatterCurse,
			Type:  CooldownTypeSurvival,
		})
	case proto.Race_RaceTauren:
		character.MultiplyStat(stats.Health, 1.05)
		// Endurance: 2% increased critical damage and healing (was 1% hit).
		character.PseudoStats.CritDamageMultiplier *= 1.02
	case proto.Race_RaceTroll:
		applyCreatureTypeSlaying(character, proto.MobType_MobTypeBeast)

		actionID := ActionID{SpellID: 20554}

		berserkingAura := character.RegisterAura(Aura{
			Label:    "Berserking",
			ActionID: actionID,
			Duration: time.Second * 10,
		}).
			AttachMultiplyAttackSpeed(1.1).
			AttachMultiplyCastSpeed(1.1)

		berserking := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				berserkingAura.Activate(sim)
			},

			RelatedSelfBuff: berserkingAura,
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: berserking,
			Type:  CooldownTypeDPS,
		})
	case proto.Race_RaceUndead:
		applyTouchOfTheGrave(character)
	case proto.Race_RaceHighOrderSkyborne:
		applySkyborneSharedRacials(character)

		energizedAura := character.RegisterAura(Aura{
			Label:    "Energized",
			ActionID: ActionID{SpellID: 1270842},
			Duration: time.Second * 15,
			OnGain: func(aura *Aura, sim *Simulation) {
				if aura.Unit.HasManaBar() {
					aura.Unit.MultiplyManaRegenSpeed(sim, 2)
				}
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				if aura.Unit.HasManaBar() {
					aura.Unit.MultiplyManaRegenSpeed(sim, 0.5)
				}
			},
		})

		// Client 1.60.1.70058: Read Ley Line has a 2 sec cast (SpellMisc CastingTimeIndex).
		character.RegisterSpell(SpellConfig{
			ActionID: ActionID{SpellID: 1259705},
			Flags:    SpellFlagAPL | SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				DefaultCast: Cast{
					GCD: GCDDefault,
					CastTime: time.Second * 2,
				},
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				energizedAura.Activate(sim)
			},

			RelatedSelfBuff: energizedAura,
		})
	case proto.Race_RaceWindshaperSkyborne:
		applySkyborneSharedRacials(character)

		elementalBlessingAura := character.RegisterAura(Aura{
			Label:    "Elemental Blessing",
			ActionID: ActionID{SpellID: 1259688},
			Duration: time.Second * 30,
			OnGain: func(aura *Aura, sim *Simulation) {
				aura.Unit.MultiplyMovementSpeed(sim, 1.1)
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				aura.Unit.MultiplyMovementSpeed(sim, 1/1.1)
			},
		})

		// Client 1.60.1.70058: Skysight has a 0.5 sec cast (SpellMisc CastingTimeIndex).
		character.RegisterSpell(SpellConfig{
			ActionID: ActionID{SpellID: 1259686},
			Flags:    SpellFlagAPL | SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				DefaultCast: Cast{
					GCD: GCDDefault,
					CastTime: time.Millisecond * 500,
				},
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				elementalBlessingAura.Activate(sim)
			},

			RelatedSelfBuff: elementalBlessingAura,
		})
	}
}

func applySkyborneSharedRacials(character *Character) {
	applyCreatureTypeSlaying(character, proto.MobType_MobTypeElemental)

	// Wind Blessed: 3% haste (was 1%).
	character.PseudoStats.AttackSpeedMultiplier *= 1.03
	character.PseudoStats.CastSpeedMultiplier *= 1.03
}

// Beast Slaying (troll 20557), Big Game Hunter (dwarf 1259721) and Elemental Insight (Skyborne
// 1259707) are one A_MOD_DAMAGE_DONE_VERSUS 5 each in client 1.60.1.70009, with no crit damage
// effect; the hunter's Humanoid/Monster Slaying carry that as a second aura (547), these don't.
func applyCreatureTypeSlaying(character *Character, mobType proto.MobType) {
	character.Env.RegisterPostFinalizeEffect(func() {
		for _, at := range character.AttackTables {
			if at.Defender.MobType == mobType {
				at.DamageDealtMultiplier *= 1.05
			}
		}
	})
}

func applyExpansiveMind(character *Character) {
	switch character.Class {
	case proto.Class_ClassWarrior:
		character.rageBar.maxRage *= 1.05
	case proto.Class_ClassRogue:
		character.energyBar.maxEnergy *= 1.05
	case proto.Class_ClassPriest, proto.Class_ClassMage, proto.Class_ClassWarlock:
		character.MultiplyStat(stats.Mana, 1.05)
	}
}

// Eureka! names the spells it changes with each class's spell family mask. Build 70170 masks, as
// effect 0 (SPELLMOD_COST), effect 1 (SPELLMOD_DAMAGE) and effect 2 (SPELLMOD_DOT) state them; core
// cannot import spelldata, so they are mirrored here. 70170 made it "non-periodic damaging abilities":
// the DoTs (Rend, Rupture, Garrote, Corruption, Shadow Word: Pain, Devouring Plague...) and the utility
// spells (Sunder Armor, Life Tap...) take no bonus and spend no charge. SPELLMOD_DAMAGE is the direct
// part only: 70094 named Rend, Corruption and Immolate in both effects 1 and 2, which only adds up if
// effect 1 skips their ticks. So Immolate, Fireball, Pyroblast and Holy Fire get +10% on the hit and
// none on the DoT, while Drain Life/Drain Soul and Mind Flay keep theirs through effect 2.
func applyEureka(character *Character) {
	var spellID int32
	var resourceType proto.ResourceType
	var costMask, directMask, dotMask [4]uint32
	// -10% cost and +10% damage on all five class variants.
	const costReduction, damageBonus = 0.1, 0.1

	switch character.Class {
	case proto.Class_ClassRogue:
		spellID, resourceType = 1259812, proto.ResourceType_ResourceTypeEnergy
		costMask = [4]uint32{100794910, 2097408}
		directMask = costMask
	case proto.Class_ClassWarrior:
		spellID, resourceType = 1259813, proto.ResourceType_ResourceTypeRage
		costMask = [4]uint32{1784679630, 5, 1}
		directMask = costMask
	case proto.Class_ClassMage:
		spellID, resourceType = 1259817, proto.ResourceType_ResourceTypeMana
		costMask = [4]uint32{549591799, 4096}
		directMask = [4]uint32{551688947, 4096, 0, 8}
	case proto.Class_ClassWarlock:
		spellID, resourceType = 1259821, proto.ResourceType_ResourceTypeMana
		costMask = [4]uint32{541165, 8650880}
		directMask, dotMask = [4]uint32{524773, 8388736}, [4]uint32{16392, 262144}
	case proto.Class_ClassPriest:
		spellID, resourceType = 1259823, proto.ResourceType_ResourceTypeMana
		costMask = [4]uint32{150224528, 8388614}
		directMask, dotMask = [4]uint32{150224528, 8486918}, [4]uint32{10485760, 98304}
	default:
		return
	}

	family := ClassSpellFamilies[character.Class]
	names := func(spell *Spell, mask [4]uint32) bool {
		return ClassFlags{Family: family, Mask: mask}.Matches(ClientClassFlags(spell.ActionID.SpellID))
	}
	// Direct hits scale with the direct mask, DoT ticks with the DoT mask; ticks read
	// DamageMultiplier too, so a direct-only spell's DoTs are scaled back out.
	scaleDamage := func(spell *Spell, f float64) {
		direct, periodic := names(spell, directMask), names(spell, dotMask)
		if direct {
			spell.DamageMultiplier *= f
		}
		if direct == periodic {
			return
		}
		if direct {
			f = 1 / f
		}
		for _, dot := range append(DotArray{spell.aoeDot}, spell.dots...) {
			if dot != nil {
				dot.PeriodicDamageMultiplier *= f
			}
		}
	}

	actionID := ActionID{SpellID: spellID}
	const anyClassSpell = math.MaxInt64

	costMod := character.AddDynamicMod(SpellModConfig{
		Kind:         SpellMod_Custom,
		ClassMask:    anyClassSpell,
		ResourceType: resourceType,
		ApplyCustom: func(_ *SpellMod, spell *Spell) {
			if names(spell, costMask) {
				spell.Cost.PercentModifier *= 1 - costReduction
			}
		},
		RemoveCustom: func(_ *SpellMod, spell *Spell) {
			if names(spell, costMask) {
				spell.Cost.PercentModifier /= 1 - costReduction
			}
		},
	})
	damageMod := character.AddDynamicMod(SpellModConfig{
		Kind:      SpellMod_Custom,
		ClassMask: anyClassSpell,
		ApplyCustom: func(_ *SpellMod, spell *Spell) {
			scaleDamage(spell, 1+damageBonus)
		},
		RemoveCustom: func(_ *SpellMod, spell *Spell) {
			scaleDamage(spell, 1/(1+damageBonus))
		},
	})

	aura := character.RegisterAura(Aura{
		Label:     "Eureka!",
		ActionID:  actionID,
		Duration:  time.Second * 15,
		MaxStacks: 3,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.SetStacks(sim, 3)
			costMod.Activate()
			damageMod.Activate()
		},
		OnExpire: func(_ *Aura, _ *Simulation) {
			costMod.Deactivate()
			damageMod.Deactivate()
		},
		OnCastComplete: func(aura *Aura, sim *Simulation, spell *Spell) {
			if spell.Matches(anyClassSpell) && (names(spell, costMask) || names(spell, directMask) || names(spell, dotMask)) {
				aura.RemoveStack(sim)
			}
		},
	})

	spell := character.RegisterSpell(SpellConfig{
		ActionID: actionID,
		Flags:    SpellFlagNoOnCastComplete,
		Cast: CastConfig{
			CD: Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Minute * 2,
			},
		},
		ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
			aura.Activate(sim)
		},

		RelatedSelfBuff: aura,
	})

	character.AddMajorCooldown(MajorCooldown{
		Spell: spell,
		Type:  CooldownTypeDPS,
	})
}

func applyTouchOfTheGrave(character *Character) {
	// 16 Oct build: 30% for casters, 20% for paladins, rogues and warriors, 25 sec internal cooldown.
	auraID, procChance := int32(1260201), 0.3
	switch character.Class {
	case proto.Class_ClassWarrior, proto.Class_ClassPaladin, proto.Class_ClassRogue:
		auraID, procChance = 1260189, 0.2
	}

	drainID := ActionID{SpellID: 1260198}
	healthMetrics := character.NewHealthMetrics(drainID)

	drain := character.RegisterSpell(SpellConfig{
		ActionID:    drainID,
		SpellSchool: SpellSchoolShadow,
		DefenseType: DefenseTypeMagic,
		ProcMask:    ProcMaskEmpty,
		Flags:       SpellFlagProc | SpellFlagPassiveSpell | SpellFlagIgnoreAttackerModifiers | SpellFlagNoSpellMods,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			result := spell.CalcAndDealDamage(sim, target, 0.05*spell.Unit.MaxHealth(), spell.OutcomeMagicHit)
			if result.Landed() && spell.Unit.HasHealthBar() {
				spell.Unit.GainHealth(sim, result.Damage, healthMetrics)
			}
		},
	})

	character.MakeProcTriggerAura(ProcTrigger{
		Name:       "Touch of the Grave",
		ActionID:   ActionID{SpellID: auraID},
		Callback:   CallbackOnSpellHitDealt,
		ProcMask:   ProcMaskMelee | ProcMaskRanged | ProcMaskSpellDamage,
		Outcome:    OutcomeLanded,
		ProcChance: procChance,
		ICD:        time.Second * 25,
		Handler: func(sim *Simulation, spell *Spell, result *SpellResult) {
			drain.Cast(sim, result.Target)
		},
	})
}
