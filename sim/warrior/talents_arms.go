package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
)

func (warrior *Warrior) registerArmsTalents() {
	// Tier 1
	warrior.registerImprovedHeroicStrike()
	warrior.registerDeflection()
	warrior.registerImprovedRend()

	// Tier 2
	// Improved Charge: charge.go
	// Improved Tactical Mastery: stances.go
	warrior.registerImprovedOverpower()

	// Tier 3
	warrior.registerAngerManagement()
	warrior.registerDeepWounds()

	// Tier 4
	warrior.registerSpearingStrike()
	warrior.registerTwoHandedWeaponSpecialization()
	warrior.registerImpale()

	// Tier 5
	warrior.registerBloodthrill()
	warrior.registerSweepingStrikes()
	warrior.registerWeaponmaster()

	// Tier 6
	warrior.registerImprovedSlam()
	warrior.registerImprovedHamstring()

	// Tier 7
	warrior.registerMortalStrike()
}

/*
 * Arms
 */
func (warrior *Warrior) registerImprovedHeroicStrike() {
	if warrior.Talents.ImprovedHeroicStrike == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskHeroicStrike,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  int32(spellData.ImprovedHeroicStrike.TenthsAt(warrior.Talents.ImprovedHeroicStrike)),
	})
}
func (warrior *Warrior) registerDeflection() {
	if warrior.Talents.Deflection == 0 {
		return
	}

	warrior.PseudoStats.BaseParryChance += spellData.Deflection.FractionAt(warrior.Talents.Deflection)
}

func (warrior *Warrior) registerImprovedRend() {
	if warrior.Talents.ImprovedRend == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskRend,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.ImprovedRend.FractionAt(warrior.Talents.ImprovedRend),
	})
}

func (warrior *Warrior) registerImprovedOverpower() {
	if warrior.Talents.ImprovedOverpower == 0 {
		return
	}

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label:    "Improved Overpower",
		ActionID: core.ActionID{SpellID: 12963}.WithTag(warrior.Talents.ImprovedOverpower),
	})).AttachSpellMod(core.SpellModConfig{
		ClassMask:  SpellMaskOverpower,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: spellData.ImprovedOverpower.ValueAt(warrior.Talents.ImprovedOverpower),
	})
}

func (warrior *Warrior) registerAngerManagement() {
	if !warrior.Talents.AngerManagement {
		return
	}

	angerManagementRank := spellData.AngerManagement.Highest()
	angerManagementRage := angerManagementRank.Effects[1].BasePoints
	angerManagementPeriod := time.Duration(angerManagementRank.Effects[2].BasePoints) * time.Second

	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: angerManagementRank.ID})

	warrior.RegisterResetEffect(func(sim *core.Simulation) {
		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period: angerManagementPeriod,
			OnAction: func(sim *core.Simulation) {
				if sim.CurrentTime > 0 {
					warrior.AddRage(sim, angerManagementRage, rageMetrics)
				}
			},
		})
	})
}

// DeepWoundsUsesStrikingWeapon makes an off-hand crit's Deep Wounds a share of the off-hand weapon's
// average damage. Forever's beta development notes: "Fixed a bug where Deep Wounds was scaling with
// attack power for your Offhand weapon", with the dev note "Deep Wounds does not scale with Attack Power
// but is a rolling DoT compared to Classic/Vanilla". So an off-hand crit's bleed is sized off the off
// hand, without attack power, as the main hand's is. False sizes every crit's bleed off the main hand,
// as before. A shield strike (Shield Slam, Shield Bash) is a main-hand attack to the server and stays
// on the main hand.
var DeepWoundsUsesStrikingWeapon = true

func (warrior *Warrior) registerDeepWounds() {
	if warrior.Talents.DeepWounds == 0 {
		return
	}

	deepWoundsBleed := spellData.DeepWoundsTriggered.ByID(412609)

	share := spellData.DeepWounds.FractionAt(warrior.Talents.DeepWounds)
	tick := deepWoundsBleed.EffectN(1)

	// Beta logs (foreverlogs reports 33/35, two level 20 warriors with 1 point) fit a tick of share x the main hand's average
	// weapon damage / 4, attack power left out, and a crit that lands on a running bleed adds what it still owed to the new one.
	// Report 2705 confirms the top rank at level 30: 3 points on a 102-154 axe tick 19.25 against 60% x 128 / 4 = 19.2.
	// The tick spell 412613 carries Attributes[6] 0x20000000 (ignore caster damage modifiers), and reports 2705/2708 agree:
	// the same warrior's ticks hold ~19.3 in Defensive Stance while Thunder Clap and Rend lose their 10% there.
	// A refresh keeps the bleed's tick timer (Zirene on forever-bugs #234): the 12 sec start over, the
	// next tick lands when it was due, and the owed and new damage spread over the ticks left (patch 89).
	// An off-hand crit sizes its share off the off-hand weapon (DeepWoundsUsesStrikingWeapon).
	var struck *core.Weapon
	warrior.DeepWounds = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: deepWoundsBleed.ID},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskEmpty,
		ClassSpellMask: SpellMaskDeepWounds,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreResists | core.SpellFlagProc | core.SpellFlagIgnoreAttackerModifiers, // 12162 and 412609 lack Not a Proc.

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "DeepWounds",
			},
			NumberOfTicks: int32(deepWoundsBleed.Duration() / tick.Period()),
			TickLength:    tick.Period(),

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, dot.SnapshotBaseDamage, deepWoundsBleed.TickOutcome(dot))
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHitNoHitCounter)
			dot := spell.Dot(target)
			owed := dot.OutstandingDmg()
			dot.ApplyKeepingTickTimer(sim)
			weapon := warrior.AutoAttacks.MH()
			if struck != nil {
				weapon, struck = struck, nil
			}
			dot.SnapshotBaseDamage = (owed + weapon.AverageDamage()*share) / float64(dot.RemainingTicks())
		},
	})

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Deep Wounds - Trigger",
		TriggerImmediately: true,
		ProcMaskExclude:    core.ProcMaskEmpty,
		Outcome:            core.OutcomeCrit,
		Callback:           core.CallbackOnSpellHitDealt,
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) bool {
			return spell.SpellSchool.Matches(core.SpellSchoolPhysical)
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if DeepWoundsUsesStrikingWeapon && spell.ProcMask.Matches(core.ProcMaskMeleeOH) && warrior.AutoAttacks.IsDualWielding {
				struck = warrior.AutoAttacks.OH()
			}
			warrior.DeepWounds.Cast(sim, result.Target)
			struck = nil
		},
	})

}

func (warrior *Warrior) registerTwoHandedWeaponSpecialization() {
	if warrior.Talents.TwoHandedWeaponSpecialization == 0 {
		return
	}

	// The effect is all physical damage, auto attacks included, so no mask narrows it.
	weaponMod := warrior.AddDynamicMod(core.SpellModConfig{
		School:     core.SpellSchoolPhysical,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: spellData.TwoHandedWeaponSpecialization.Effect(dbcenums.A_MOD_DAMAGE_PERCENT_DONE, 1).FractionAt(warrior.Talents.TwoHandedWeaponSpecialization),
	})

	if warrior.GetMainHandType() == proto.HandType_HandTypeTwoHand {
		weaponMod.Activate()
	}

	warrior.RegisterItemSwapCallback(core.AllMeleeWeaponSlots(), func(sim *core.Simulation, slot proto.ItemSlot) {
		if warrior.GetMainHandType() == proto.HandType_HandTypeTwoHand {
			weaponMod.Activate()
		} else {
			weaponMod.Deactivate()
		}
	})
}

func (warrior *Warrior) registerImpale() {
	if warrior.Talents.Impale == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskDamageSpells,
		Kind:       core.SpellMod_CritMultiplier_Flat,
		FloatValue: spellData.Impale.FractionAt(warrior.Talents.Impale),
	})
}

func (warrior *Warrior) registerMortalStrike() {
	if !warrior.Talents.MortalStrike {
		return
	}

	mortalStrikeRank := spellData.MortalStrike.Highest()
	mortalStrikeBaseDamage := mortalStrikeRank.DamageEffect().Average(core.CharacterLevel)

	warrior.MortalStrike = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: mortalStrikeRank.ID},
		SpellSchool:    mortalStrikeRank.SpellSchool(),
		DefenseType:    mortalStrikeRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskMortalStrike,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   int32(mortalStrikeRank.Cost()),
			Refund: mortalStrikeRank.MissRefund(),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: mortalStrikeRank.GCD(),
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownOf(mortalStrikeRank),
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := mortalStrikeBaseDamage + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}

func (warrior *Warrior) registerSpearingStrike() {
	if !warrior.Talents.SpearingStrike {
		return
	}

	spearingStrikeRank := spellData.SpearingStrike.Highest()
	// The tooltip reads "deals $s2% weapon damage" and "an additional ${$s2*$s3}%" against Giants and
	// Dragonkin, and the effects share an aura and misc value, so both are taken by effect index.
	spearingStrikeWeaponShare := spearingStrikeRank.Effects[1].Percent()
	spearingStrikeMobtypeMultiplier := 1 + spearingStrikeRank.Effects[2].BasePoints

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spearingStrikeRank.ID},
		SpellSchool:    spearingStrikeRank.SpellSchool(),
		DefenseType:    spearingStrikeRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskSpearingStrike,
		MaxRange:       float64(spearingStrikeRank.MaxRange),

		RageCost: core.RageCostOptions{
			Cost:   int32(spearingStrikeRank.Cost()),
			Refund: spearingStrikeRank.MissRefund(),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: spearingStrikeRank.GCD(),
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownOf(spearingStrikeRank),
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		// Client 1.60.1.70170 (StanceMask 0x10000 on 1310222): "Spearing Strike no longer requires a 2handed
		// weapon. Spearing Strike requires Battle Stance." The weapon half needs nothing here, since the sim
		// never checked the weapon, and the weapon's own damage is used whatever it is.
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BattleStance)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spearingStrikeWeaponShare * spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			if target.MobType == proto.MobType_MobTypeGiant || target.MobType == proto.MobType_MobTypeDragonkin {
				baseDamage *= spearingStrikeMobtypeMultiplier
			}

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}

func (warrior *Warrior) registerBloodthrill() {
	if warrior.Talents.Bloodthrill == 0 {
		return
	}

	// Client 70425 dropped the tooltip's "Lasts $1289681d" (a 6s aura): the proc opens the 5s Overpower
	// window (1282733) the talent's effect triggers.
	bloodthrillProc := spellData.BloodthrillTriggered.ByID(1282733)

	// The proc makes Overpower usable for the buff's duration; the cast consumes it like a dodge
	// would.
	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:     "Bloodthrill - Trigger",
		ActionID: core.ActionID{SpellID: 1289682},
		Callback: core.CallbackOnSpellHitDealt,
		// 1289682's proc flags are 0x14 (melee auto attacks and melee abilities) with attr3 0x400,
		// main hand only: white MH swings plus MH abilities, Heroic Strike and Cleave included.
		ProcMask:   core.ProcMaskMeleeMH,
		Outcome:    core.OutcomeLanded,
		ProcChance: spellData.Bloodthrill.FractionAt(warrior.Talents.Bloodthrill),
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return warrior.Rend.Dot(result.Target).IsActive()
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.OverpowerAura.Activate(sim)
			warrior.OverpowerAura.UpdateExpires(sim.CurrentTime + bloodthrillProc.Duration())
		},
	})
}

func (warrior *Warrior) registerWeaponmaster() {
	if warrior.Talents.Weaponmaster == 0 {
		return
	}

	rank := warrior.Talents.Weaponmaster
	actionID := core.ActionID{SpellID: 1290261}

	mainHandIs := func(weaponTypes ...proto.WeaponType) bool {
		return warrior.GetProcMaskForTypes(weaponTypes...).Matches(core.ProcMaskMeleeMH)
	}
	var critMask core.ProcMask
	var armorIgnoreOn bool
	var swordMask core.ProcMask
	readWeapons := func() {
		critMask = warrior.GetProcMaskForTypes(proto.WeaponType_WeaponTypeAxe, proto.WeaponType_WeaponTypePolearm)
		armorIgnoreOn = mainHandIs(proto.WeaponType_WeaponTypeMace, proto.WeaponType_WeaponTypeStaff)
		swordMask = warrior.GetProcMaskForTypes(proto.WeaponType_WeaponTypeSword)
	}
	readWeapons()

	// The crit goes with the weapon, as rogue Hack and Slash does (#654): an axe in one hand crits
	// more with that hand only. (The client effects are dummies; the tooltip and Classic's Axe
	// Specialization both scope it to the weapon.)
	crit := spellData.Weaponmaster.EffectAt(1).ValueAt(rank)
	critMods := map[core.ProcMask]*core.SpellMod{
		core.ProcMaskMeleeMH: warrior.AddDynamicMod(core.SpellModConfig{Kind: core.SpellMod_BonusCrit_Percent, ProcMask: core.ProcMaskMeleeMH, FloatValue: crit}),
		core.ProcMaskMeleeOH: warrior.AddDynamicMod(core.SpellModConfig{Kind: core.SpellMod_BonusCrit_Percent, ProcMask: core.ProcMaskMeleeOH, FloatValue: crit}),
	}
	setCritMods := func() {
		for hand, mod := range critMods {
			if critMask.Matches(hand) {
				mod.Activate()
			} else {
				mod.Deactivate()
			}
		}
	}
	setCritMods()

	armorIgnore := spellData.Weaponmaster.EffectAt(2).FractionAt(rank)
	addArmorIgnore := func(delta float64) {
		for _, attackTable := range warrior.AttackTables {
			attackTable.ArmorIgnoreFactor += delta
		}
	}
	armorIgnoreAura := warrior.RegisterAura(core.Aura{
		Label:    "Weaponmaster (Mace/Staff)",
		ActionID: actionID.WithTag(2),
		Duration: core.NeverExpires,
		OnGain: func(_ *core.Aura, _ *core.Simulation) {
			addArmorIgnore(armorIgnore)
		},
		OnExpire: func(_ *core.Aura, _ *core.Simulation) {
			addArmorIgnore(-armorIgnore)
		},
	})
	if armorIgnoreOn {
		core.MakePermanent(armorIgnoreAura)
	}

	var extraAttack *core.Spell
	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Weaponmaster (Sword)",
		ActionID:           actionID.WithTag(3),
		MetricsActionID:    actionID.WithTag(3),
		Callback:           core.CallbackOnSpellHitDealt,
		ProcMask:           core.ProcMaskMelee,
		Outcome:            core.OutcomeLanded,
		ProcChance:         spellData.Weaponmaster.EffectAt(3).FractionAt(rank),
		TriggerImmediately: true,
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) bool {
			// 16 October 2026 beta notes: Sword Specialization no longer triggers from Shield Bash
			// and Shield Slam.
			return spell.ProcMask.Matches(swordMask) && spell != extraAttack &&
				!spell.Matches(SpellMaskShieldBash|SpellMaskShieldSlam)
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.AutoAttacks.MaybeReplaceMHSwing(sim, extraAttack).Cast(sim, result.Target)
		},
	}).ApplyOnInit(func(aura *core.Aura, sim *core.Simulation) {
		config := *warrior.AutoAttacks.MHConfig()
		config.ActionID = config.ActionID.WithTag(actionID.SpellID)
		extraAttack = warrior.GetOrRegisterSpell(config)
	})

	setActive := func(sim *core.Simulation, aura *core.Aura, on bool) {
		if on {
			aura.Activate(sim)
		} else {
			aura.Deactivate(sim)
		}
	}
	warrior.RegisterItemSwapCallback(core.AllMeleeWeaponSlots(), func(sim *core.Simulation, slot proto.ItemSlot) {
		readWeapons()
		setCritMods()
		setActive(sim, armorIgnoreAura, armorIgnoreOn)
	})
}

func (warrior *Warrior) registerImprovedHamstring() {
	if warrior.Talents.ImprovedHamstring == 0 {
		return
	}

	improvedHamstringRoot := spellData.ImprovedHamstringTriggered.Highest()

	immobilizeAuras := warrior.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Improved Hamstring-" + warrior.Label,
			ActionID: core.ActionID{SpellID: improvedHamstringRoot.ID},
			Duration: improvedHamstringRoot.Duration(),
		})
	})

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:           "Improved Hamstring - Trigger",
		ActionID:       core.ActionID{SpellID: 12289},
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: SpellMaskHamstring,
		Outcome:        core.OutcomeLanded,
		ProcChance:     spellData.ImprovedHamstring.FractionAt(warrior.Talents.ImprovedHamstring),
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			immobilizeAuras.Get(result.Target).Activate(sim)
		},
	})
}

func (warrior *Warrior) registerImprovedSlam() {
	if warrior.Talents.ImprovedSlam == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskSlam,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Millisecond * time.Duration(spellData.ImprovedSlam.Effect(dbcenums.A_ADD_FLAT_MODIFIER, int32(dbcenums.SPELLMOD_CASTING_TIME)).ValueAt(warrior.Talents.ImprovedSlam)),
	})

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskSlam,
		Kind:      core.SpellMod_GlobalCooldown_Flat,
		TimeValue: time.Millisecond * time.Duration(spellData.ImprovedSlam.Effect(dbcenums.A_ADD_FLAT_MODIFIER, int32(dbcenums.SPELLMOD_GLOBAL_COOLDOWN)).ValueAt(warrior.Talents.ImprovedSlam)),
	})

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskSlam,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: time.Millisecond * time.Duration(spellData.ImprovedSlam.Effect(dbcenums.A_ADD_FLAT_MODIFIER, int32(dbcenums.SPELLMOD_COOLDOWN)).ValueAt(warrior.Talents.ImprovedSlam)),
	})
}

func (warrior *Warrior) registerSweepingStrikes() {
	if !warrior.Talents.SweepingStrikes {
		return
	}

	sweepingStrikesRank := spellData.SweepingStrikes.Highest()

	actionID := core.ActionID{SpellID: 12723}

	var copyDamage float64
	hitSpell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: SpellMaskSweepingStrikesHit,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeSpecial,
		Flags:          core.SpellFlagIgnoreModifiers | core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, copyDamage, spell.OutcomeAlwaysHit)
		},
	})

	warrior.SweepingStrikesNormalizedAttack = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(1), // Real SpellID: 26654
		ClassSpellMask: SpellMaskSweepingStrikesNormalizedHit,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
		},
	})

	warrior.SweepingStrikesAura = warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Sweeping Strikes",
		ActionID:           actionID,
		MetricsActionID:    actionID,
		Duration:           sweepingStrikesRank.Duration(),
		Callback:           core.CallbackOnSpellHitDealt,
		ProcMask:           core.ProcMaskMelee,
		Outcome:            core.OutcomeLanded,
		TriggerImmediately: true,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if warrior.Env.ActiveTargetCount() < 2 || warrior.SweepingStrikesAura.GetStacks() == 0 || result.PostOutcomeDamage <= 0 {
				return
			}

			if spell.Matches(SpellMaskSweepingStrikesHit | SpellMaskSweepingStrikesNormalizedHit | SpellMaskThunderClap | SpellMaskWhirlwind | SpellMaskWhirlwindOh) {
				return
			}

			nextTarget := warrior.Env.NextActiveTargetUnit(result.Target)
			if spell.Matches(SpellMaskExecute) && sim.IsExecutePhase20() {
				warrior.SweepingStrikesNormalizedAttack.Cast(sim, nextTarget)
			} else {
				copyDamage = result.Damage / result.ArmorAndResistanceMultiplier
				hitSpell.Cast(sim, nextTarget)
			}

			warrior.SweepingStrikesAura.RemoveStack(sim)
		},
	})
	warrior.SweepingStrikesAura.MaxStacks = int32(sweepingStrikesRank.ProcCharges)

	ssCD := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: SpellMaskSweepingStrikes,
		SpellSchool:    core.SpellSchoolPhysical,

		RageCost: core.RageCostOptions{
			Cost: int32(sweepingStrikesRank.Cost()),
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownOf(sweepingStrikesRank),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BattleStance)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
			warrior.SweepingStrikesAura.SetStacks(sim, int32(sweepingStrikesRank.ProcCharges))
		},

		RelatedSelfBuff: warrior.SweepingStrikesAura,
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: ssCD,
		Type:  core.CooldownTypeDPS,
		// It only strikes "an additional nearby opponent": 30 rage for nothing on one target.
		ShouldActivate: func(sim *core.Simulation, _ *core.Character) bool {
			return sim.Environment.ActiveTargetCount() >= 2
		},
	})
}
