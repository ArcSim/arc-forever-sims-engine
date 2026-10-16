package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

func (warrior *Warrior) registerProtectionTalents() {
	// Tier 1
	// Improved Bloodrage: bloodrage.go
	warrior.registerShieldSpecialization()
	warrior.registerIronWill()

	// Tier 2
	warrior.registerAnticipation()
	warrior.registerImprovedRevenge()
	warrior.registerImprovedThunderClap()

	// Tier 3
	warrior.registerLastStand()
	warrior.registerMasterOfDefense()
	warrior.registerImprovedDisarm()
	// Defiance: stances.go

	// Tier 4
	warrior.registerImprovedSunderArmor()
	// Vanguard: charge.go
	warrior.registerImprovedShieldBash()

	// Tier 5
	warrior.registerImprovedShieldWall()
	warrior.registerConcussionBlow()
	warrior.registerFocusedRage()

	// Tier 6
	warrior.registerBastion()

	// Tier 7
	warrior.registerShieldSlam()
}

func (warrior *Warrior) registerAnticipation() {
	if warrior.Talents.Anticipation == 0 {
		return
	}

	warrior.AddStat(stats.DefenseRating, spellData.Anticipation.ValueAt(warrior.Talents.Anticipation)*core.DefenseRatingPerDefenseLevel)
}

func (warrior *Warrior) registerShieldSpecialization() {
	if warrior.Talents.ShieldSpecialization == 0 {
		return
	}

	shieldSpecializationEnergize := spellData.ShieldSpecializationTriggered.Highest()

	// BlockPercent is in percent (5 = 5%) since the core's avoidance rework.
	warrior.AddStat(stats.BlockPercent, spellData.ShieldSpecialization.Effect(dbcenums.A_MOD_BLOCK_PERCENT, 0).ValueAt(warrior.Talents.ShieldSpecialization))

	warrior.registerRageOnAvoid(
		"Shield Specialization",
		shieldSpecializationEnergize,
		spellData.ShieldSpecialization.EffectAt(2).FractionAt(warrior.Talents.ShieldSpecialization),
		core.OutcomeBlock,
		nil,
	)
}

func (warrior *Warrior) registerRageOnAvoid(name string, energize *spelldata.Spell, chance float64, outcome core.HitOutcome, extra core.ProcExtraCondition) {
	rage := energize.EnergizeEffect().Tenths()
	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: energize.ID})
	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               name,
		ProcChance:         chance,
		TriggerImmediately: true,
		Outcome:            outcome,
		Callback:           core.CallbackOnSpellHitTaken,
		ExtraCondition:     extra,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.AddRage(sim, rage, rageMetrics)
		},
	})
}

func (warrior *Warrior) registerIronWill() {
	if warrior.Talents.IronWill == 0 {
		return
	}
	warrior.PseudoStats.FearDurationMultiplier = spellData.IronWill.Effect(dbcenums.A_MECHANIC_DURATION_MOD, 1).MultiplierAt(warrior.Talents.IronWill)
	warrior.PseudoStats.StunDurationMultiplier = spellData.IronWill.Effect(dbcenums.A_MECHANIC_DURATION_MOD, 12).MultiplierAt(warrior.Talents.IronWill)
}

func (warrior *Warrior) registerLastStand() {
	if !warrior.Talents.LastStand {
		return
	}

	lastStandRank := spellData.LastStand.Highest()
	lastStandBuff := spellData.LastStandTriggered.Highest()
	actionID := core.ActionID{SpellID: lastStandRank.ID}
	healthMetrics := warrior.NewHealthMetrics(actionID)

	// Client 1.60.1.70291 moved the buff from A_MOD_MAX_HEALTH to A_MOD_INCREASE_HEALTH, still 30; the
	// talent tooltip reads "30% of your maximum health", so it stays a share of maximum health.
	healthPercent := lastStandBuff.Effect(dbcenums.A_MOD_INCREASE_HEALTH, 0).Percent()
	var bonusHealth float64
	aura := warrior.RegisterAura(core.Aura{
		Label:    "Last Stand",
		ActionID: actionID,
		Duration: lastStandBuff.Duration(),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			bonusHealth = warrior.MaxHealth() * healthPercent
			warrior.UpdateMaxHealth(sim, bonusHealth, healthMetrics)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.UpdateMaxHealth(sim, -bonusHealth, healthMetrics)
		},
	})

	spell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: SpellMaskLastStand,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownOf(lastStandRank),
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			aura.Activate(sim)
		},

		RelatedSelfBuff: aura,
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
		BuffAura: &core.StatBuffAura{
			Aura:            aura,
			BuffedStatTypes: []stats.Stat{stats.Health},
		},
	})
}

func (warrior *Warrior) registerImprovedSunderArmor() {
	if warrior.Talents.ImprovedSunderArmor == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskSunderArmor,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  int32(spellData.ImprovedSunderArmor.TenthsAt(warrior.Talents.ImprovedSunderArmor)),
	})
}

func (warrior *Warrior) registerImprovedShieldWall() {
	if warrior.Talents.ImprovedShieldWall == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskShieldWall,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: time.Duration(spellData.ImprovedShieldWall.ValueAt(warrior.Talents.ImprovedShieldWall)) * time.Millisecond,
	})
}

// TODO: In-game testing if this generates threat
func (warrior *Warrior) registerConcussionBlow() {
	if !warrior.Talents.ConcussionBlow {
		return
	}

	concussionBlowRank := spellData.ConcussionBlow.Highest()

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: concussionBlowRank.ID},
		ClassSpellMask: SpellMaskConcussionBlow,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   int32(concussionBlowRank.Cost()),
			Refund: concussionBlowRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownOf(concussionBlowRank),
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}

func (warrior *Warrior) registerShieldSlam() {
	if !warrior.Talents.ShieldSlam {
		return
	}

	shieldSlamRank := spellData.ShieldSlam.Highest()

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: shieldSlamRank.ID},
		ClassSpellMask: SpellMaskShieldSlam,
		SpellSchool:    shieldSlamRank.SpellSchool(),
		DefenseType:    shieldSlamRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskMeleeOHSpecial, // struck with the shield (upstream b6ce30c972)
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   int32(shieldSlamRank.Cost()),
			Refund: shieldSlamRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: shieldSlamRank.GCD(),
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownOf(shieldSlamRank),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.PseudoStats.CanBlock
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		// Not in the client table; our Classic value until measured in game.
		FlatThreatBonus: 254 * 2,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Rank 4 rolls 655 +-2.3% (640-670) off the row.
			baseDamage := shieldSlamRank.DamageEffect().Roll(sim, core.CharacterLevel) + warrior.BlockDamageReduction()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}

func (warrior *Warrior) registerFocusedRage() {
	if warrior.Talents.FocusedRage == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskFocusedRage,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  int32(spellData.FocusedRage.TenthsAt(warrior.Talents.FocusedRage)),
	})
}

func (warrior *Warrior) registerMasterOfDefense() {
	if warrior.Talents.MasterOfDefense == 0 {
		return
	}

	masterOfDefenseEnergize := spellData.MasterOfDefenseTriggered.Highest()

	warrior.registerRageOnAvoid(
		"Master of Defense",
		masterOfDefenseEnergize,
		spellData.MasterOfDefense.FractionAt(warrior.Talents.MasterOfDefense),
		core.OutcomeDodge|core.OutcomeParry,
		func(_ *core.Simulation, _ *core.Spell, _ *core.SpellResult) bool {
			return warrior.PseudoStats.CanBlock
		},
	)
}

func (warrior *Warrior) registerImprovedRevenge() {
	if warrior.Talents.ImprovedRevenge == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskRevenge,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.ImprovedRevenge.FractionAt(warrior.Talents.ImprovedRevenge),
	})
}

func (warrior *Warrior) registerImprovedDisarm() {
	if warrior.Talents.ImprovedDisarm == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskDisarm,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: time.Duration(spellData.ImprovedDisarm.ValueAt(warrior.Talents.ImprovedDisarm)) * time.Millisecond,
	})
}

func (warrior *Warrior) registerImprovedShieldBash() {
	if warrior.Talents.ImprovedShieldBash == 0 {
		return
	}

	improvedShieldBashSilence := spellData.ImprovedShieldBashTriggered.Highest()

	// TODO: nothing in the sim reads a silence on an enemy, so the aura only shows up in metrics.
	silenceAuras := warrior.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Shield Bash - Silence",
			ActionID: core.ActionID{SpellID: improvedShieldBashSilence.ID},
			Duration: improvedShieldBashSilence.Duration(),
		})
	})

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Improved Shield Bash",
		ProcChance:         spellData.ImprovedShieldBash.FractionAt(warrior.Talents.ImprovedShieldBash),
		TriggerImmediately: true,
		ClassSpellMask:     SpellMaskShieldBash,
		Outcome:            core.OutcomeLanded,
		Callback:           core.CallbackOnSpellHitDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			silenceAuras.Get(result.Target).Activate(sim)
		},
	})
}

func (warrior *Warrior) registerBastion() {
	if warrior.Talents.Bastion == 0 {
		return
	}

	damageMod := warrior.AddDynamicMod(core.SpellModConfig{
		School:     core.SpellSchoolPhysical,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: spellData.Bastion.FractionAt(warrior.Talents.Bastion),
	})

	if warrior.PseudoStats.CanBlock {
		damageMod.Activate()
	}

	warrior.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotOffHand}, func(sim *core.Simulation, slot proto.ItemSlot) {
		if warrior.PseudoStats.CanBlock {
			damageMod.Activate()
		} else {
			damageMod.Deactivate()
		}
	})
}

func (warrior *Warrior) registerImprovedThunderClap() {
	if warrior.Talents.ImprovedThunderClap == 0 {
		return
	}

	// Slowing effect implemented in thunder_clap.go
	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskThunderClap,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  int32(spellData.ImprovedThunderClap.TenthsAt(warrior.Talents.ImprovedThunderClap)),
	})
}
