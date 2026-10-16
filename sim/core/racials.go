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
		applyEureka(agent)
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

		// 16 October beta notes: "Blood Fury: Cooldown increased to 3 minutes (was 2 minutes)."
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
		// 16 October beta notes: "Magic damage taken reduced by 5% (was 15%)."
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
					GCD:      GCDDefault,
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
					GCD:      GCDDefault,
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

// EurekaSpells names, as the class's own spell masks, the abilities client 1.60.1.70170 lists on the
// three effects of its Eureka! rows (1259812 Rogue, 1259813 Warrior, 1259817 Mage, 1259821 Warlock,
// 1259823 Priest). Core cannot see the client's family masks, so each class that can be a Gnome states
// them, and its tests hold them to the rows.
type EurekaSpells struct {
	// Effect 0, SPELLMOD_COST -10%.
	Cost int64
	// Effect 1, SPELLMOD_DAMAGE +10%: the spells whose direct hits gain the bonus. A spell the sim deals as
	// a hit of a triggered spell, as the client does for Arcane Missiles, Blizzard and Rain of Fire, is named
	// here by the mask its tick spell carries.
	Damage int64
	// Effect 2, SPELLMOD_DOT +10%, and the channels the client deals as hits of a triggered spell (Hellfire)
	// that the sim deals as ticks: the spells whose dot ticks gain the bonus. Client 70170 empties this
	// effect for every class but the Priest and Warlock, whose channels stay on it ("channeled spells do
	// not count as periodics").
	Tick int64
}

// EurekaAgent is implemented by the agents of the classes that can be Gnomes.
type EurekaAgent interface {
	EurekaSpells() EurekaSpells
}

// Eureka! raises the damage of the abilities the client's rows name by 10% and cuts their cost. The 16
// October 2026 beta notes make it "an 8 second buff with no charges. Same effect on damage and ability
// costs." Client 1.60.1.70170 ("no longer benefits periodic effects at all. Channeled spells do not count as
// periodics") takes the bleeds and damage-over-time spells off every class's list, so a dot's ticks keep
// their damage unless the spell is in Tick.
func applyEureka(agent Agent) {
	character := agent.GetCharacter()
	var spellID int32
	var resourceType proto.ResourceType
	procMask := ProcMaskSpecial
	// Build 70009 states -10% cost (effect 0, SPELLMOD_COST) on all five class variants and 70170 keeps
	// it; core cannot import spelldata, so the value is mirrored here.
	const costReduction = 0.1

	switch character.Class {
	case proto.Class_ClassRogue:
		spellID, resourceType = 1259812, proto.ResourceType_ResourceTypeEnergy
	case proto.Class_ClassWarrior:
		spellID, resourceType = 1259813, proto.ResourceType_ResourceTypeRage
	case proto.Class_ClassMage:
		spellID, resourceType = 1259817, proto.ResourceType_ResourceTypeMana
	case proto.Class_ClassWarlock:
		spellID, resourceType = 1259821, proto.ResourceType_ResourceTypeMana
	case proto.Class_ClassPriest:
		spellID, resourceType = 1259823, proto.ResourceType_ResourceTypeMana
		procMask |= ProcMaskSpellHealing
	default:
		return
	}

	actionID := ActionID{SpellID: spellID}
	const anyClassSpell = math.MaxInt64
	spells := EurekaSpells{Cost: anyClassSpell, Damage: anyClassSpell}
	if eurekaAgent, ok := agent.(EurekaAgent); ok {
		spells = eurekaAgent.EurekaSpells()
	}

	costMod := character.AddDynamicMod(SpellModConfig{
		Kind:         SpellMod_PowerCost_Pct,
		ClassMask:    spells.Cost,
		ProcMask:     procMask,
		ResourceType: resourceType,
		FloatValue:   -costReduction,
	})
	// A spell's damage multiplier covers its hits and its ticks alike, so the +10% goes on every spell
	// of either list, and the ticks of a spell that is only on the direct list take it back.
	damageMod := character.AddDynamicMod(SpellModConfig{
		Kind:       SpellMod_DamageDone_Pct,
		ClassMask:  spells.Damage | spells.Tick,
		ProcMask:   procMask,
		FloatValue: 0.1,
	})
	var tickCancelMod *SpellMod
	if directOnly := spells.Damage &^ spells.Tick; directOnly != 0 {
		tickCancelMod = character.AddDynamicMod(SpellModConfig{
			Kind:       SpellMod_DotDamageDone_Pct,
			ClassMask:  directOnly,
			ProcMask:   procMask,
			FloatValue: 1/1.1 - 1,
		})
	}

	aura := character.RegisterAura(Aura{
		// 16 Oct build: an 8 s buff with no charges (was 15 s / 3 charges).
		Label:    "Eureka!",
		ActionID: actionID,
		Duration: time.Second * 8,
		OnGain: func(_ *Aura, _ *Simulation) {
			costMod.Activate()
			damageMod.Activate()
			if tickCancelMod != nil {
				tickCancelMod.Activate()
			}
		},
		OnExpire: func(_ *Aura, _ *Simulation) {
			costMod.Deactivate()
			damageMod.Deactivate()
			if tickCancelMod != nil {
				tickCancelMod.Deactivate()
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

// Touch of the Grave drains 5% of the caster's maximum health. The 16 October 2026 beta notes: "a 30%
// chance to activate for casters and a 20% chance to activate for Paladins/Rogues/Warriors, as well as a
// 25 second internal cooldown" (was 10% and 5% with a 1 sec cooldown), and it "no longer scales with
// Shadow increases on the target or the caster".
func applyTouchOfTheGrave(character *Character) {
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
		Flags:       SpellFlagProc | SpellFlagPassiveSpell | SpellFlagIgnoreModifiers | SpellFlagNoSpellMods,

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
