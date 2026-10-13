package core

import (
	"slices"
	"time"

	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
)

func applyRaceEffects(agent Agent) {
	character := agent.GetCharacter()

	switch character.Race {
	case proto.Race_RaceDwarf:
		mhItem := character.Equipment.MainHand()
		ohItem := character.Equipment.OffHand()
		if mhItem.WeaponType == proto.WeaponType_WeaponTypeMace || ohItem.WeaponType == proto.WeaponType_WeaponTypeMace {
			bonus := 1.0 * CritRatingPerCritChance
			character.OnSpellRegistered(func(spell *Spell) {
				if spell.ProcMask.Matches(ProcMaskRangedAuto) {
					return
				}
				spell.BonusCritRating += bonus
			})
		}

		character.Env.RegisterPostFinalizeEffect(func() {
			for _, t := range character.Env.Encounter.Targets {
				if t.MobType == proto.MobType_MobTypeBeast {
					for _, at := range character.AttackTables[t.UnitIndex] {
						at.DamageDealtMultiplier *= 1.05
					}
				}
			}
		})

		actionID := ActionID{SpellID: 20594}
		stoneFormAura := character.RegisterAura(Aura{
			Label:    "Stoneform",
			ActionID: actionID,
			Duration: time.Second * 8,
			OnGain: func(aura *Aura, sim *Simulation) {
				character.PseudoStats.DamageTakenMultiplier *= 0.90
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				character.PseudoStats.DamageTakenMultiplier /= 0.90
			},
		})

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete | SpellFlagAPL,
			Cast: CastConfig{
				DefaultCast: Cast{GCD: 0},
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
		if character.HasRageBar() {
			character.AddMaxRage(5)
		} else if character.HasEnergyBar() {
			character.AddStat(stats.Energy, 5)
		} else {
			character.MultiplyStat(stats.Mana, 1.05)
		}

	case proto.Race_RaceHuman:
		character.MultiplyStat(stats.Spirit, 1.05)

		mhItem := character.Equipment.MainHand()
		ohItem := character.Equipment.OffHand()
		if mhItem.WeaponType == proto.WeaponType_WeaponTypeSword || ohItem.WeaponType == proto.WeaponType_WeaponTypeSword {
			bonus := 2.0 * CritRatingPerCritChance
			character.OnSpellRegistered(func(spell *Spell) {
				if spell.ProcMask.Matches(ProcMaskRangedAuto) {
					return
				}
				spell.BonusCritRating += bonus
			})
		}

	case proto.Race_RaceNightElf:
		character.AddStat(stats.Dodge, 1)
		// Classic PseudoStats has no MovementSpeedMultiplier. Skip the 2% run speed.

		actionID := ActionID{SpellID: 1259799}
		bonus := 10.0 * CritRatingPerCritChance
		spellBonus := 10.0 * SpellCritRatingPerCritChance
		elunesLightAura := character.RegisterAura(Aura{
			Label:    "Elune's Light",
			ActionID: actionID,
			Duration: time.Second * 15,
			OnGain: func(aura *Aura, sim *Simulation) {
				character.AddStatDynamic(sim, stats.MeleeCrit, bonus)
				character.AddStatDynamic(sim, stats.SpellCrit, spellBonus)
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				character.AddStatDynamic(sim, stats.MeleeCrit, -bonus)
				character.AddStatDynamic(sim, stats.SpellCrit, -spellBonus)
			},
		})

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete | SpellFlagAPL,
			Cast: CastConfig{
				DefaultCast: Cast{GCD: 0},
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				elunesLightAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeDPS,
		})

	case proto.Race_RaceOrc:
		mhItem := character.Equipment.MainHand()
		ohItem := character.Equipment.OffHand()
		if mhItem.WeaponType == proto.WeaponType_WeaponTypeAxe || ohItem.WeaponType == proto.WeaponType_WeaponTypeAxe {
			bonus := 1.0 * CritRatingPerCritChance
			character.OnSpellRegistered(func(spell *Spell) {
				if spell.ProcMask.Matches(ProcMaskRangedAuto) {
					return
				}
				spell.BonusCritRating += bonus
			})
		}

		actionID := ActionID{SpellID: 20572}
		bloodFuryAura := character.RegisterAura(Aura{
			Label:    "Blood Fury",
			ActionID: actionID,
			Duration: time.Second * 15,
			OnGain: func(aura *Aura, sim *Simulation) {
				character.PseudoStats.DamageDealtMultiplier *= 1.10
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				character.PseudoStats.DamageDealtMultiplier /= 1.10
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
		character.MultiplyStat(stats.Health, 1.05)
		character.AddStat(stats.MeleeHit, 1)
		character.AddStat(stats.SpellHit, 1)

	case proto.Race_RaceTroll:
		character.Env.RegisterPostFinalizeEffect(func() {
			for _, t := range character.Env.Encounter.Targets {
				if t.MobType == proto.MobType_MobTypeBeast {
					for _, at := range character.AttackTables[t.UnitIndex] {
						at.DamageDealtMultiplier *= 1.05
					}
				}
			}
		})

		actionID := ActionID{SpellID: 26297}
		berserkingAura := character.RegisterAura(Aura{
			Label:    "Berserking",
			ActionID: actionID,
			Duration: time.Second * 10,
			OnGain: func(aura *Aura, sim *Simulation) {
				character.MultiplyCastSpeed(1.10)
				character.MultiplyAttackSpeed(sim, 1.10)
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				character.MultiplyCastSpeed(1 / 1.10)
				character.MultiplyAttackSpeed(sim, 1/1.10)
			},
		})

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete | SpellFlagAPL,
			Cast: CastConfig{
				DefaultCast: Cast{GCD: 0},
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
			Spell: spell,
			Type:  CooldownTypeDPS,
		})

	case proto.Race_RaceUndead:
	}
}

func (character *Character) GetFaction() proto.Faction {
	if slices.Contains([]proto.Race{proto.Race_RaceHuman, proto.Race_RaceDwarf, proto.Race_RaceGnome, proto.Race_RaceNightElf}, character.Race) {
		return proto.Faction_Alliance
	} else if slices.Contains([]proto.Race{proto.Race_RaceOrc, proto.Race_RaceTroll, proto.Race_RaceTauren, proto.Race_RaceUndead}, character.Race) {
		return proto.Faction_Horde
	}
	return proto.Faction_Unknown
}
