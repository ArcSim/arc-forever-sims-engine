package priest

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/buffs"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

func (priest *Priest) registerDisciplineTalents() {
	// Tier 1
	priest.applyPowerInLight()
	priest.applyWandSpecialization()
	priest.applyTwinDisciplines()

	// Tier 2
	priest.applySilentResolve()
	priest.applyHolyPrecision()
	priest.applyImprovedPowerWordShield()
	priest.applyMartyrdom()

	// Tier 3
	priest.applyMentalAgility()
	priest.applyInnerFocus()
	priest.applyMeditation()

	// Tier 4
	priest.applyImprovedInnerFire()
	priest.applyMentalStrength()
	priest.applySoulWarding()
	priest.applyImprovedManaBurn()

	// Tier 5
	priest.applyPenance()
	priest.applyRenewedHope()

	// Tier 6
	priest.applyDivineAegis()

	// Tier 7
	priest.applyPowerInfusion()
}

// Power in Light is new in Forever: Smite and Penance hit 2% harder per point while this priest's
// Holy Fire or Shadow Word: Pain is on the target (SW:P added by the 16 Oct beta notes). The client
// states the ladder on a dummy effect (1309969), so the spells it names come from the tooltip.
func (priest *Priest) applyPowerInLight() {
	if priest.Talents.PowerInLight == 0 {
		return
	}

	multiplier := spellData.PowerInLight.MultiplierAt(priest.Talents.PowerInLight)

	for _, target := range priest.Env.Encounter.AllTargetUnits {
		target.AddDynamicDamageTakenModifier(func(_ *core.Simulation, spell *core.Spell, result *core.SpellResult, _ bool) {
			if spell.Unit == &priest.Unit && spell.Matches(PriestSpellSmite|PriestSpellPenance) && priest.hasPowerInLightDot(result.Target) {
				result.Damage *= multiplier
			}
		})
	}
}

// applyWandSpecialization implements Wand Specialization, new in Forever.
//
// TODO: To be implemented. The talent raises wand damage (14524, +13/25%) and the sim has never
// modelled a wand attack, so there is nothing for it to act on yet.
func (priest *Priest) applyWandSpecialization() {
	if priest.Talents.WandSpecialization == 0 {
		return
	}
}

// Twin Disciplines is new in Forever: +1% damage and healing per point on instant spells. The client
// states it as two SPELL_AURA_ADD_PCT_MODIFIER effects, damage and dot, so it joins the additive
// bucket rather than multiplying. 1225132's masks name Holy Nova, Chastise, Divine Grace and
// Contingency Plan (damage) and Shadow Word: Pain and Devouring Plague (dot). The devs called the
// missing Mind Flay and Shadow Word: Death bugs (forever-bugs #322, #328, 2026-10-13), and the 16 Oct
// beta notes add Starshards and Penance too, so all four join the masks.
func (priest *Priest) applyTwinDisciplines() {
	if priest.Talents.TwinDisciplines == 0 {
		return
	}

	priest.AddStaticMod(core.SpellModConfig{
		ClassMask: PriestSpellHolyNova | PriestSpellChastise | PriestSpellDivineGrace | PriestSpellContingencyPlan |
			PriestSpellShadowWordPain | PriestSpellDevouringPlague | PriestSpellMindFlay | PriestSpellShadowWordDeath |
			PriestSpellStarshards | PriestSpellPenance,
		FloatValue: spellData.TwinDisciplines.EffectAt(1).FractionAt(priest.Talents.TwinDisciplines),
		Kind:       core.SpellMod_DamageDone_Flat,
	})
}

func (priest *Priest) applySilentResolve() {
	if priest.Talents.SilentResolve == 0 {
		return
	}

	// Client 14523 still names Holy (misc 2); the 16 Oct beta notes widen it to all spells, so read
	// the first effect by position and drop the school.
	priest.AddStaticMod(core.SpellModConfig{
		ClassMask:  PriestSpellsAll,
		FloatValue: spellData.SilentResolve.EffectAt(1).FractionAt(priest.Talents.SilentResolve),
		Kind:       core.SpellMod_ThreatMultiplier_Pct,
	})
}

// Holy Precision is new in Forever: +6% hit per point. 1309957 is a miss chance mod on a class mask,
// not school hit: Smite, Holy Fire, Holy Nova and the Penance bolts are in, Chastise (1277335) is not
// (client 1.60.1.70205). The mask also catches Vampiric Embrace and Shadowfiend, which never roll hit here.
func (priest *Priest) applyHolyPrecision() {
	if priest.Talents.HolyPrecision == 0 {
		return
	}

	priest.AddStaticMod(core.SpellModConfig{
		ClassMask:  PriestHolySpells,
		Kind:       core.SpellMod_BonusHit_Percent,
		FloatValue: spellData.HolyPrecision.Effect(dbcenums.A_ADD_FLAT_MODIFIER, int32(dbcenums.SPELLMOD_RESIST_MISS_CHANCE)).ValueAt(priest.Talents.HolyPrecision),
	})
}

// applyImprovedPowerWordShield implements Improved Power Word: Shield, new in Forever.
//
// TODO: To be implemented. Power Word: Shield is not modelled, so the talent has nothing to raise.
func (priest *Priest) applyImprovedPowerWordShield() {
	if priest.Talents.ImprovedPowerWordShield == 0 {
		return
	}
}

// applyMartyrdom implements Martyrdom, new in Forever.
//
// TODO: To be implemented. It triggers Focused Casting (27828) off damage taken, and nothing in a
// DPS sim takes damage yet.
func (priest *Priest) applyMartyrdom() {
	if priest.Talents.Martyrdom == 0 {
		return
	}
}

// Instant spells, and the two cast-time spells the Forever tooltip adds: Smite and Holy Fire. 14520's
// mask leaves out the channels and cooldowns: Mind Flay, Penance, Shadow Word: Death, Starshards
// and Shadowfiend pay full price.
func (priest *Priest) applyMentalAgility() {
	if priest.Talents.MentalAgility == 0 {
		return
	}

	priest.AddStaticMod(core.SpellModConfig{
		ClassMask: PriestSpellSmite | PriestSpellHolyFire | PriestSpellHolyNova | PriestSpellShadowWordPain |
			PriestSpellDevouringPlague | PriestSpellVampiricEmbrace | PriestSpellPowerInfusion | PriestSpellShadowform |
			PriestSpellFade | PriestSpellChastise | PriestSpellConfoundingFlash | PriestSpellContingencyPlan,
		FloatValue: spellData.MentalAgility.FractionAt(priest.Talents.MentalAgility),
		Kind:       core.SpellMod_PowerCost_Pct_Add,
	})
}

func (priest *Priest) applyInnerFocus() {
	if !priest.Talents.InnerFocus {
		return
	}

	// The cost cut (14751 e0) covers every priest spell; the crit (e1) leaves out Mind Flay,
	// Shadow Word: Death and Starshards.
	rank := spellData.InnerFocus.Highest()
	critPercent := rank.Effect(dbcenums.A_ADD_FLAT_MODIFIER, int32(dbcenums.SPELLMOD_CRITICAL_CHANCE)).Average(core.CharacterLevel)
	critMod := priest.AddDynamicMod(core.SpellModConfig{
		ClassMask:  PriestSpellsAll &^ (PriestSpellMindFlay | PriestSpellShadowWordDeath | PriestSpellStarshards),
		FloatValue: critPercent,
		Kind:       core.SpellMod_BonusCrit_Percent,
	})
	priest.innerFocusCrit, priest.innerFocusCritChance = critMod, critPercent/100

	var innerFocusSpell *core.Spell
	costPercent := int32(rank.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_COST)).Average(core.CharacterLevel))

	priest.InnerFocusAura = priest.RegisterAura(core.Aura{
		Label:    "Inner Focus",
		ActionID: core.ActionID{SpellID: rank.ID},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, _ *core.Simulation) {
			aura.Unit.PseudoStats.SpellCostPercentModifier += costPercent
			critMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SpellCostPercentModifier -= costPercent
			critMod.Deactivate()
			innerFocusSpell.CD.Use(sim)
		},
		// The buff is spent by the next priest spell, whatever it is.
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(PriestSpellsAll) {
				aura.Deactivate(sim)
			}
		},
	})

	innerFocusSpell = priest.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: rank.ID},
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: max(rank.Cooldown(), rank.CategoryCooldown()),
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			priest.InnerFocusAura.Activate(sim)
		},

		RelatedSelfBuff: priest.InnerFocusAura,
	})

	priest.AddMajorCooldown(core.MajorCooldown{
		Spell: innerFocusSpell,
		Type:  core.CooldownTypeMana,
	})
}

// 17/33/50, not the 17/34/51 that multiplying rank 1 would give: the client states each rank.
func (priest *Priest) applyMeditation() {
	if priest.Talents.Meditation == 0 {
		return
	}

	priest.PseudoStats.SpiritRegenRateCasting += spellData.Meditation.FractionAt(priest.Talents.Meditation)
	priest.UpdateManaRegenRates()
}

// applyImprovedInnerFire implements Improved Inner Fire, new in Forever.
//
// TODO: To be implemented. It raises Inner Fire's armour and charges (14747); the sim gives Inner
// Fire a flat armour value and never spends its charges, so neither half changes a DPS run.
func (priest *Priest) applyImprovedInnerFire() {
	if priest.Talents.ImprovedInnerFire == 0 {
		return
	}
}

// +3% Intellect per point in the Forever client, where TBC's raised total mana.
func (priest *Priest) applyMentalStrength() {
	if priest.Talents.MentalStrength == 0 {
		return
	}

	priest.MultiplyStat(stats.Intellect, spellData.MentalStrength.MultiplierAt(priest.Talents.MentalStrength))
}

// applySoulWarding implements Soul Warding, new in Forever.
//
// TODO: To be implemented. Power Word: Shield is not modelled, so the talent has nothing to act on.
func (priest *Priest) applySoulWarding() {
	if !priest.Talents.SoulWarding {
		return
	}
}

// applyImprovedManaBurn implements Improved Mana Burn, new in Forever.
//
// TODO: To be implemented. Mana Burn is not modelled; it drains a target's mana, which no boss has.
func (priest *Priest) applyImprovedManaBurn() {
	if priest.Talents.ImprovedManaBurn == 0 {
		return
	}
}

func (priest *Priest) applyPenance() {
	if !priest.Talents.Penance {
		return
	}

	spellData.Penance.Each(func(_ int32, rank *spelldata.Spell) { priest.registerPenanceSpell(rank) })
}

// applyRenewedHope implements Renewed Hope, new in Forever.
//
// TODO: To be implemented. It buffs Power Word: Shield and Flash Heal, neither of which is modelled.
func (priest *Priest) applyRenewedHope() {
	if priest.Talents.RenewedHope == 0 {
		return
	}
}

// applyDivineAegis implements Divine Aegis, new in Forever.
//
// TODO: To be implemented. The shield it leaves on a healing crit needs a healing spell to crit.
func (priest *Priest) applyDivineAegis() {
	if priest.Talents.DivineAegis == 0 {
		return
	}
}

// The priest casts Power Infusion on itself, which is the reason the Smite build goes thirty-one
// points deep. The proto still carries a target option for the healing specs, which the sim has no
// way to act on yet.
// TODO: let the option pick a raid member once buffing another player is modelled.
func (priest *Priest) applyPowerInfusion() {
	if !priest.Talents.PowerInfusion {
		return
	}

	rank := spellData.PowerInfusion.Highest()
	piAura := buffs.PowerInfusionsAura(&priest.Unit, true, 0)

	piSpell := priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID, Tag: priest.Index},
		SpellSchool:    core.SpellSchoolHoly,
		Flags:          core.SpellFlagHelpful | core.SpellFlagAPL,
		ClassSpellMask: PriestSpellPowerInfusion,

		ManaCost: rank.ManaCost(),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: max(rank.Cooldown(), rank.CategoryCooldown()),
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			piAura.Activate(sim)
		},

		RelatedSelfBuff: piAura,
	})

	priest.AddMajorCooldown(core.MajorCooldown{
		Spell:    piSpell,
		Priority: core.CooldownPriorityDefault,
		Type:     core.CooldownTypeDPS,
	})
}
