package warlock

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

func (warlock *Warlock) registerDemonologyTalents() {
	// Tier 1
	warlock.applyImprovedHealthFunnel()
	warlock.applyImprovedImp()
	warlock.applyDemonicEmbrace()
	warlock.applyUnholyPower()

	// Tier 2
	// Demonic Aegis: armors.go
	warlock.applyImprovedVoidwalker()
	warlock.applyFelVitality()
	warlock.applyDemonicEnergies()

	// Tier 3
	warlock.applyImprovedSayaad()
	warlock.applyDemonicSacrifice()
	warlock.applyMasterSummoner()

	// Tier 4
	warlock.applyDecimation()
	warlock.applyFelDomination()
	warlock.applyDemonicBrand()

	// Tier 5
	warlock.applyImprovedFelhunter()
	warlock.applySoulLink()
	warlock.applyDemonicKnowledge()

	// Tier 6
	warlock.applyMasterDemonologist()

	// Tier 7
	warlock.applyDemonicPact()
}

// The Firebolt half of 18694; its first effect carries the same ladder for Blood Pact, which the
// raid buff handles.
func (warlock *Warlock) applyImprovedImp() {
	if warlock.Talents.ImprovedImp == 0 || warlock.Options.SacrificeSummon {
		return
	}

	warlock.Imp.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.ImprovedImp.EffectAt(2).FractionAt(warlock.Talents.ImprovedImp),
		ClassMask:  WarlockSpellImpFireBolt,
	})
}

// 18697 raises stamina 3% a rank. The tooltip no longer names Classic's spirit penalty, but the client
// keeps it as a second effect with no rank curve: -1% spirit at every rank.
func (warlock *Warlock) applyDemonicEmbrace() {
	if warlock.Talents.DemonicEmbrace == 0 {
		return
	}

	warlock.MultiplyStat(stats.Stamina, spellData.DemonicEmbrace.EffectAt(1).MultiplierAt(warlock.Talents.DemonicEmbrace))
	warlock.MultiplyStat(stats.Spirit, spellData.DemonicEmbrace.EffectAt(2).MultiplierAt(warlock.Talents.DemonicEmbrace))
}

// 2% more pet damage a point (18769).
func (warlock *Warlock) applyUnholyPower() {
	if warlock.Talents.UnholyPower == 0 || warlock.Options.SacrificeSummon {
		return
	}

	multiplier := spellData.UnholyPower.MultiplierAt(warlock.Talents.UnholyPower)
	for _, pet := range warlock.BasePets {
		pet.PseudoStats.DamageDealtMultiplier *= multiplier
	}
}

// 5% more mana for the warlock and 5% more health and mana for the demon, a point (18731).
func (warlock *Warlock) applyFelVitality() {
	if warlock.Talents.FelVitality == 0 {
		return
	}

	multiplier := spellData.FelVitality.EffectAt(1).MultiplierAt(warlock.Talents.FelVitality)
	warlock.MultiplyStat(stats.Mana, multiplier)
	for _, pet := range warlock.BasePets {
		pet.MultiplyStat(stats.Health, multiplier)
		pet.MultiplyStat(stats.Mana, multiplier)
	}
}

// 10% more Lash of Pain damage a point: 18754 effect 0 (the tooltip's $s1, SPELLMOD_ALL_EFFECTS);
// effect 1 (SPELLMOD_DURATION) is Seduction's duration.
func (warlock *Warlock) applyImprovedSayaad() {
	if warlock.Talents.ImprovedSayaad == 0 || warlock.Options.SacrificeSummon {
		return
	}

	warlock.Succubus.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.ImprovedSayaad.EffectAt(1).FractionAt(warlock.Talents.ImprovedSayaad),
		ClassMask:  WarlockSpellSuccubusLashOfPain,
	})
}

// Each demon leaves behind the opposing aspect, and Forever's pairing is the reverse of Classic's:
// the Imp leaves Shadow damage (18789, school mask 32), the Succubus Fire (18791, mask 4), the
// Voidwalker mana (18792) and the Felhunter health (18790). The demon is sacrificed before the pull,
// so the buff is simply permanent and no pet is ever summoned.
func (warlock *Warlock) applyDemonicSacrifice() {
	if !warlock.Talents.DemonicSacrifice {
		return
	}

	demon := warlock.Options.Summon
	if !warlock.Options.SacrificeSummon {
		// Demonic Pact (425464): a demon sacrificed before the pull keeps its buff while a different
		// one is out; summoning the sacrificed demon again cancels it.
		if !warlock.Talents.DemonicPact || warlock.Options.PactSacrifice == warlock.Options.Summon {
			return
		}
		demon = warlock.Options.PactSacrifice
	}

	var spellID int32
	var school stats.SchoolIndex
	switch demon {
	case proto.WarlockOptions_Imp:
		spellID, school = 18789, stats.SchoolIndexShadow
	case proto.WarlockOptions_Succubus:
		spellID, school = 18791, stats.SchoolIndexFire
	case proto.WarlockOptions_Voidwalker:
		warlock.applyFelEnergy()
		return
	default:
		// The Felhunter's health is survival only; it is left out until the sim needs it.
		return
	}

	row := spellData.DemonicSacrificeTriggered.ByID(spellID)
	multiplier := 1 + row.EffectN(1).Percent()

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label:    "Demonic Sacrifice",
		ActionID: core.ActionID{SpellID: spellID},
		Duration: row.Duration(),
	}).AttachMultiplicativePseudoStatBuff(&warlock.PseudoStats.SchoolDamageDealtMultiplier[school], multiplier))
}

// The Voidwalker's sacrifice, Fel Energy (18792): 2% of total mana every 4 s.
func (warlock *Warlock) applyFelEnergy() {
	row := spellData.DemonicSacrificeTriggered.ByID(18792)
	manaFraction := row.EffectN(1).Percent()
	period := row.EffectN(1).Period()
	manaMetrics := warlock.NewManaMetrics(core.ActionID{SpellID: row.ID})

	core.MakePermanent(warlock.RegisterAura(core.Aura{
		Label:    "Demonic Sacrifice",
		ActionID: core.ActionID{SpellID: row.ID},
		Duration: row.Duration(),
		OnGain: func(_ *core.Aura, sim *core.Simulation) {
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   period,
				Priority: core.ActionPriorityRegen,
				OnAction: func(sim *core.Simulation) {
					warlock.AddMana(sim, warlock.MaxMana()*manaFraction, manaMetrics)
				},
			})
		},
	}))
}

// Shadow Bolt and Searing Pain hit 3% harder a point below 35% health, and Soul Fire casts 20% a
// point faster and comes off cooldown 45% a point sooner (440870 / 440873).
func (warlock *Warlock) applyDecimation() {
	if warlock.Talents.Decimation == 0 {
		return
	}

	points := warlock.Talents.Decimation
	triggered := spellData.DecimationTriggered.Highest()

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_Cooldown_Multiplier,
		FloatValue: 1 + spellData.Decimation.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_COOLDOWN)).FractionAt(points),
		ClassMask:  WarlockSpellSoulFire,
	})

	warlock.DecimationAura = warlock.RegisterAura(core.Aura{
		Label:    "Decimation",
		ActionID: core.ActionID{SpellID: triggered.ID},
		Duration: triggered.Duration(),
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.Decimation.EffectAt(4).FractionAt(points),
		ClassMask:  WarlockSpellShadowBolt | WarlockSpellSearingPain,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: spellData.Decimation.EffectAt(1).FractionAt(points),
		ClassMask:  WarlockSpellSoulFire,
	})

	warlock.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Decimation Trigger",
		Callback:           core.CallbackOnSpellHitDealt,
		ClassSpellMask:     WarlockSpellShadowBolt | WarlockSpellSearingPain,
		Outcome:            core.OutcomeLanded,
		TriggerImmediately: true,
		ExtraCondition: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) bool {
			return sim.IsExecutePhase35()
		},
		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			warlock.DecimationAura.Activate(sim)
		},
	})
}

// Searing Pain sheds 17/33/50% of its threat and arms the demon with 2/4/6 branded attacks
// (1293695 / 1293696).
//
// The branded hit is the client's SpellDescriptionVariables formula for 1293696/1293697 (beta
// 1.60.1.69977, the same text Wowhead's Forever tooltip resolves): ((level-26)*1.5)+14 to +17, plus
// 7.8% of the warlock's spell power, so 65 to 68 at level 60. The client splits the hit by pet
// (1.60.1.70009): the Imp's is Fire (1293698, fire spell power), the others Shadow (1293697, shadow
// spell power). The brand 1293696 procs on the target taking melee or spell hits (ProcTypeMask
// 139944), so the Imp's Firebolt and the Succubus' Lash of Pain use charges too.
func (warlock *Warlock) applyDemonicBrand() {
	if warlock.Talents.DemonicBrand == 0 {
		return
	}

	points := warlock.Talents.DemonicBrand
	triggered := spellData.DemonicBrandTriggered.Highest()
	actionID := core.ActionID{SpellID: triggered.ID}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_ThreatMultiplier_Pct,
		FloatValue: spellData.DemonicBrand.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_THREAT)).FractionAt(points),
		ClassMask:  WarlockSpellSearingPain,
	})

	if warlock.Options.SacrificeSummon {
		return
	}

	charges := int32(spellData.DemonicBrand.Effect(dbcenums.A_ADD_FLAT_MODIFIER, int32(dbcenums.SPELLMOD_CHARGES)).ValueAt(points))

	for _, pet := range warlock.BasePets {
		// Only the Shadow hit carries the high threat (beta notes 8 Oct).
		school, schoolPower, threat := core.SpellSchoolShadow, stats.ShadowDamage, 3.0
		if pet == warlock.Imp {
			school, schoolPower, threat = core.SpellSchoolFire, stats.FireDamage, 1
		}

		brandSpell := pet.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: school,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete,

			DamageMultiplier: 1,
			ThreatMultiplier: threat,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				levelBonus := float64(core.CharacterLevel-26) * 1.5
				spellPower := warlock.GetStat(stats.SpellDamage) + warlock.GetStat(schoolPower)
				damage := sim.Roll(levelBonus+14, levelBonus+17) + 0.078*spellPower
				// 1293697/1293698 (client 1.60.1.70205) carry Always Hit (Attributes[3] 0x40000) and
				// no Cannot Crit, so the hit never misses and crits on the pet's spell crit.
				spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicCrit)
			},
		})

		pet.DemonicBrandAura = pet.RegisterAura(core.Aura{
			Label:     "Demonic Brand",
			ActionID:  actionID,
			Duration:  triggered.Duration(),
			MaxStacks: charges,
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMelee|core.ProcMaskSpellDamage) {
					brandSpell.Cast(sim, result.Target)
					aura.RemoveStack(sim)
				}
			},
		})
	}

	warlock.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Demonic Brand Trigger",
		Callback:           core.CallbackOnSpellHitDealt,
		ClassSpellMask:     WarlockSpellSearingPain,
		Outcome:            core.OutcomeLanded,
		TriggerImmediately: true,
		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			if warlock.ActivePet == nil {
				return
			}
			brandAura := warlock.ActivePet.DemonicBrandAura
			brandAura.Activate(sim)
			brandAura.SetStacks(sim, brandAura.MaxStacks)
		},
	})
}

// 3% more damage dealt and 30% of the damage taken split with the demon (25228).
func (warlock *Warlock) applySoulLink() {
	if !warlock.Talents.SoulLink || warlock.Options.SacrificeSummon {
		return
	}

	row := spellData.SoulLinkTriggered.ByID(25228)
	damageDealt := 1 + row.EffectN(1).Percent()
	damageTaken := 1 - row.EffectN(2).Percent()

	config := func(unit *core.Unit) core.Aura {
		return core.Aura{
			Label:    "Soul Link",
			ActionID: core.ActionID{SpellID: 19028},
			Duration: core.NeverExpires,
			OnGain: func(aura *core.Aura, _ *core.Simulation) {
				aura.Unit.PseudoStats.DamageDealtMultiplier *= damageDealt
				aura.Unit.PseudoStats.DamageTakenMultiplier *= damageTaken
			},
			OnExpire: func(aura *core.Aura, _ *core.Simulation) {
				aura.Unit.PseudoStats.DamageDealtMultiplier /= damageDealt
				aura.Unit.PseudoStats.DamageTakenMultiplier /= damageTaken
			},
		}
	}

	warlock.SoulLinkAura = core.MakePermanent(warlock.RegisterAura(config(&warlock.Unit)))
	for _, pet := range warlock.BasePets {
		pet.SoulLinkAura = core.MakePermanent(pet.RegisterAura(config(&pet.Unit)))
	}
}

// 33/67/100% of the warlock's level in spell power for the warlock and the demon while it is out
// (412732).
func (warlock *Warlock) applyDemonicKnowledge() {
	if warlock.Talents.DemonicKnowledge == 0 || warlock.Options.SacrificeSummon {
		return
	}

	bonus := spellData.DemonicKnowledge.FractionAt(warlock.Talents.DemonicKnowledge) * float64(core.CharacterLevel)

	config := core.Aura{
		Label:    "Demonic Knowledge",
		ActionID: core.ActionID{SpellID: 412732},
		Duration: core.NeverExpires,
	}
	core.MakePermanent(warlock.RegisterAura(config).AttachStatBuff(stats.SpellDamage, bonus))
	for _, pet := range warlock.BasePets {
		core.MakePermanent(pet.RegisterAura(config).AttachStatBuff(stats.SpellDamage, bonus))
	}
}

// 2% a point, on the school the demon out matches (23785): Fire for the Imp, Shadow for the
// Succubus, damage taken for the Voidwalker and the Felhunter. Both the warlock and the demon get it.
func (warlock *Warlock) applyMasterDemonologist() {
	if warlock.Talents.MasterDemonologist == 0 || warlock.Options.SacrificeSummon {
		return
	}

	fraction := spellData.MasterDemonologist.EffectAt(1).FractionAt(warlock.Talents.MasterDemonologist)

	var label string
	var tag int32
	var school stats.SchoolIndex
	switch warlock.Options.Summon {
	case proto.WarlockOptions_Imp:
		label, tag, school = "Master Demonologist (Imp)", 1, stats.SchoolIndexFire
	case proto.WarlockOptions_Succubus:
		label, tag, school = "Master Demonologist (Succubus)", 3, stats.SchoolIndexShadow
	default:
		// The Voidwalker's and the Felhunter's halves only cut damage taken.
		return
	}

	buff := func(unit *core.Unit) *core.Aura {
		return core.MakePermanent(unit.RegisterAura(core.Aura{
			Label:    label,
			ActionID: core.ActionID{SpellID: 23785, Tag: tag},
			Duration: core.NeverExpires,
		}).AttachMultiplicativePseudoStatBuff(&unit.PseudoStats.SchoolDamageDealtMultiplier[school], 1+fraction))
	}

	warlock.MasterDemonologistAura = buff(&warlock.Unit)
	buff(&warlock.ActivePet.Unit)
}

// applyImprovedHealthFunnel implements Improved Health Funnel, new in Forever.
//
// Health Funnel is not modelled, so its cost, threat and healing bonuses have nothing to act on.
func (warlock *Warlock) applyImprovedHealthFunnel() {
	if warlock.Talents.ImprovedHealthFunnel == 0 {
		return
	}
}

// applyImprovedVoidwalker implements Improved Voidwalker, new in Forever.
//
// 10% a point on the Voidwalker's Torment and Sacrifice, neither of which is modelled.
func (warlock *Warlock) applyImprovedVoidwalker() {
	if warlock.Talents.ImprovedVoidwalker == 0 {
		return
	}
}

// applyDemonicEnergies implements Demonic Energies, new in Forever.
//
// The pet's share of Life Tap is in lifetap.go. The first effect (client 1225214, 8/15%) heals the
// pet for that share of the warlock's spell damage; the sim's demon takes no damage, so it is left out.
func (warlock *Warlock) applyDemonicEnergies() {
	if warlock.Talents.DemonicEnergies == 0 {
		return
	}
}

// applyMasterSummoner implements Master Summoner, new in Forever.
//
// The summon spells are not modelled - the demon is out from the start - so the cast time and cost
// cuts have nothing to act on.
func (warlock *Warlock) applyMasterSummoner() {
	if warlock.Talents.MasterSummoner == 0 {
		return
	}
}

// applyFelDomination implements Fel Domination, new in Forever.
//
// A summon cooldown; nothing to model while the demon never has to be resummoned.
func (warlock *Warlock) applyFelDomination() {
	if !warlock.Talents.FelDomination {
		return
	}
}

// applyImprovedFelhunter implements Improved Felhunter, new in Forever.
//
// 10% a point on the Felhunter's abilities, none of which are modelled.
func (warlock *Warlock) applyImprovedFelhunter() {
	if warlock.Talents.ImprovedFelhunter == 0 {
		return
	}
}

// applyDemonicPact implements Demonic Pact, new in Forever. The sacrifice it keeps (PactSacrifice)
// is applied in applyDemonicSacrifice.
func (warlock *Warlock) applyDemonicPact() {
	if !warlock.Talents.DemonicPact {
		return
	}
}
