package paladin

import (
	"fmt"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/buffs"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

func (paladin *Paladin) registerRetributionTalents() {
	// Tier 1
	paladin.applyDeflection()
	paladin.applyBenediction()

	// Tier 2
	paladin.applyImprovedJudgement()
	paladin.applyHolyConduit()
	paladin.applyConviction()

	// Tier 3
	paladin.applyVindication()
	paladin.applySanctifiedJudgement()
	// Seal of Command registered in registerTalentSpells
	paladin.applyPursuitOfJustice()

	// Tier 4
	paladin.applyEyeForAnEye()
	paladin.applySacredArbiter()

	// Tier 5
	paladin.applyTwoHandedWeaponSpecialization()
	paladin.applyVengeance()

	// Tier 6
	paladin.applyChampionOfTheLight()
	paladin.applyInstrumentOfLaw()

	// Tier 7
	paladin.applyTwistOfLight()
}

// Deflection - Increases your Parry chance by 1/2/3/4/5%.
func (paladin *Paladin) applyDeflection() {
	if paladin.Talents.Deflection == 0 {
		return
	}

	paladin.PseudoStats.BaseParryChance += spellData.Deflection.FractionAt(paladin.Talents.Deflection)
}

// Benediction - Reduces the Mana cost of all instant cast spells and abilities by 2/4/6/8/10%.
func (paladin *Paladin) applyBenediction() {
	if paladin.Talents.Benediction == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskBenediction,
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		FloatValue: spellData.Benediction.FractionAt(paladin.Talents.Benediction),
	})
}

// Improved Judgement - Decreases the cooldown of your Judgement ability by 1/2 sec.
func (paladin *Paladin) applyImprovedJudgement() {
	if paladin.Talents.ImprovedJudgement == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskJudgement,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: time.Duration(spellData.ImprovedJudgement.ValueAt(paladin.Talents.ImprovedJudgement)) * time.Millisecond,
	})
}

// Holy Conduit - Reduces the mana cost of your Consecration, Holy Wrath, Exorcism, and Hammer of
// Wrath spells by 20/40%.
func (paladin *Paladin) applyHolyConduit() {
	if paladin.Talents.HolyConduit == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskConsecration | SpellMaskHolyWrath | SpellMaskExorcism | SpellMaskHammerOfWrath,
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		FloatValue: spellData.HolyConduit.FractionAt(paladin.Talents.HolyConduit),
	})
}

// Conviction - Improves your chance to get a critical strike with melee attacks by 1/2/3/4/5%.
func (paladin *Paladin) applyConviction() {
	if paladin.Talents.Conviction == 0 {
		return
	}

	paladin.AddStat(stats.PhysicalCritPercent, spellData.Conviction.ValueAt(paladin.Talents.Conviction))
}

// Vindication - Gives your damaging melee attacks a chance to reduce the target's Attack Power by
// 67/134/201, and increase your Attack Power by 1/2/3% for 30 sec.
//
// The Classic, TBC and Forever clients all state the chance as 100%, and the Forever beta client
// (1.60.1.69893) is what this sim follows: every damaging melee attack that lands procs it.
const vindicationProcChance = 1.0

func (paladin *Paladin) applyVindication() {
	if paladin.Talents.Vindication == 0 {
		return
	}

	rank := spellData.VindicationTriggered.Highest()
	points := spellData.Vindication.EffectAt(1).ValueAt(paladin.Talents.Vindication)
	// The tooltip's ${$m1/-3*$440667m1}: the trigger's -201 scaled by the points over three.
	targetAttackPower := rank.Effect(dbcenums.A_MOD_ATTACK_POWER, 0).Average(core.CharacterLevel) * points / 3

	targetAuras := paladin.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Vindication" + paladin.Label,
			ActionID: core.ActionID{SpellID: rank.ID},
			Duration: rank.Duration(),
		}).AttachStatBuff(stats.AttackPower, targetAttackPower)
	})

	attackPowerDep := paladin.NewDynamicMultiplyStat(stats.AttackPower, 1+points/100)
	selfAura := paladin.RegisterAura(core.Aura{
		Label:    "Vindication" + paladin.Label,
		ActionID: core.ActionID{SpellID: rank.ID}.WithTag(1),
		Duration: rank.Duration(),
	}).AttachStatDependency(attackPowerDep)

	paladin.MakeProcTriggerAura(core.ProcTrigger{
		Name:       "Vindication - Trigger" + paladin.Label,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: vindicationProcChance,
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			targetAuras.Get(result.Target).Activate(sim)
			selfAura.Activate(sim)
		},
	})
}

// Sanctified Judgement - Gives your Judgement ability a 33/66/100% chance to return 20/40/60% of
// the Mana cost of the judged seal. The chance is the first effect, the refund the second.
//
// The game returns 10/9 of the stated share of what the seal cost: 22/44/67%. Beta logs (reports
// 2699/2701/2712, 1311077) show it on 13 paladins at all three ranks and five seal costs: Seal
// of Command 65 -> 14-15 / 29 / 43, Seal of the Crusader 40 -> 9 at rank 1, Seal of Fury 60 -> 40
// and 90 -> 60, Seal of Command rank 2 110 -> 73 at rank 3. The chance reads true (Ligma 73 of
// 223 Judgements at rank 1, 3-point paladins every one). Two paladins land on exactly 40% (Inco,
// rank 2) and 60% (Tristram, rank 3): 10/9 of a seal 10% cheaper, five points of Benediction, so
// the share applies to the cost paid.
const sanctifiedJudgementLogScale = 10.0 / 9

func (paladin *Paladin) applySanctifiedJudgement() {
	if paladin.Talents.SanctifiedJudgement == 0 {
		return
	}

	manaMetrics := paladin.NewManaMetrics(core.ActionID{SpellID: spellData.SanctifiedJudgement.Highest().ID})
	refund := spellData.SanctifiedJudgement.EffectAt(2).FractionAt(paladin.Talents.SanctifiedJudgement) * sanctifiedJudgementLogScale

	paladin.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Sanctified Judgement" + paladin.Label,
		Callback:           core.CallbackOnCastComplete,
		ClassSpellMask:     SpellMaskJudgement,
		ProcChance:         spellData.SanctifiedJudgement.EffectAt(1).FractionAt(paladin.Talents.SanctifiedJudgement),
		TriggerImmediately: true,
		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			if seal := paladin.activeSeal(); seal != nil {
				paladin.AddMana(sim, seal.spell.CurCast.Cost*refund, manaMetrics)
			}
		},
	})
}

// Eye for an Eye - All critical strikes against you cause 5/10% of the damage taken to the attacker
// as well. The damage caused by Eye for an Eye will not exceed 50% of the Paladin's total health.
func (paladin *Paladin) applyEyeForAnEye() {
	if paladin.Talents.EyeForAnEye == 0 {
		return
	}

	rank := spellData.EyeForAnEye.Highest()
	share := spellData.EyeForAnEye.FractionAt(paladin.Talents.EyeForAnEye)

	var reflected float64
	reflect := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: rank.ID},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagBinary | core.SpellFlagPassiveSpell | core.SpellFlagIgnoreModifiers,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, reflected, spell.OutcomeAlwaysHit)
		},
	})

	paladin.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Eye for an Eye" + paladin.Label,
		Callback:           core.CallbackOnSpellHitTaken,
		Outcome:            core.OutcomeCrit,
		RequireDamageDealt: true,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			reflected = min(result.Damage*share, paladin.MaxHealth()/2)
			reflect.Cast(sim, spell.Unit)
		},
	})
}

// Pursuit of Justice - Increases movement speed and mounted movement speed by 8/15%. This does not
// stack with other movement speed increasing effects.
func (paladin *Paladin) applyPursuitOfJustice() {
	if paladin.Talents.PursuitOfJustice == 0 {
		return
	}

	rank := spellData.PursuitOfJustice.Highest()
	paladin.NewPassiveMovementSpeedAura(
		"Pursuit of Justice",
		core.ActionID{SpellID: rank.ID},
		spellData.PursuitOfJustice.Effect(dbcenums.A_MOD_INCREASE_SPEED, 0).FractionAt(paladin.Talents.PursuitOfJustice),
	)
}

// Sacred Arbiter - Increases the damage of your Holy Strike ability by 20% and causes it to refresh
// all Judgement effects on the target. The paladin's own melee strikes already refresh its own
// judgements; Holy Strike with the talent refreshes every judgement on the target, whoever put it
// there, the way Crusader Strike did in TBC.
func (paladin *Paladin) applySacredArbiter() {
	if !paladin.Talents.SacredArbiter {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskHolyStrike,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.SacredArbiter.FractionAt(1),
	})

	paladin.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Sacred Arbiter" + paladin.Label,
		Callback:           core.CallbackOnSpellHitDealt,
		ClassSpellMask:     SpellMaskHolyStrike,
		Outcome:            core.OutcomeLanded,
		TriggerImmediately: true,
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			for _, aura := range result.Target.GetAurasWithTag(buffs.JudgementAuraTag) {
				if aura.IsActive() {
					aura.Refresh(sim)
				}
			}
		},
	})
}

// Two-Handed Weapon Specialization - Increases the damage you deal with two-handed melee weapons
// by 2/4/6%. The client puts it on the Physical school alone.
func (paladin *Paladin) applyTwoHandedWeaponSpecialization() {
	if paladin.Talents.TwoHandedWeaponSpecialization == 0 {
		return
	}

	paladin.applyWeaponSpecialization(
		spellData.TwoHandedWeaponSpecialization.FractionAt(paladin.Talents.TwoHandedWeaponSpecialization),
		proto.HandType_HandTypeTwoHand,
	)
}

// Vengeance - Increases your Physical and Holy damage dealt by 1/2/3% for 30 sec after landing a
// non-periodic critical strike. Stacks up to 3 times.
//
// 20049's proc flags (69972) carry no periodic flag, and the hit callback never sees a tick. The
// stack cap is 20050's CumulativeAura.
func (paladin *Paladin) applyVengeance() {
	if paladin.Talents.Vengeance == 0 {
		return
	}

	rank := spellData.VengeanceTriggered.Highest()
	perStack := spellData.Vengeance.FractionAt(paladin.Talents.Vengeance)

	// 20050 is damage done (A79), which never raises a heal.
	damageMod := paladin.AddDynamicMod(core.SpellModConfig{
		School:     core.SpellSchoolHoly | core.SpellSchoolPhysical,
		ProcMask:   ^core.ProcMaskSpellHealing,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: perStack,
	})

	vengeance := paladin.RegisterAura(core.Aura{
		Label:     "Vengeance" + paladin.Label,
		ActionID:  core.ActionID{SpellID: rank.ID},
		Duration:  rank.Duration(),
		MaxStacks: int32(rank.MaxStack),
		OnGain: func(_ *core.Aura, _ *core.Simulation) {
			damageMod.Activate()
		},
		OnExpire: func(_ *core.Aura, _ *core.Simulation) {
			damageMod.Deactivate()
		},
		OnStacksChange: func(_ *core.Aura, _ *core.Simulation, _, newStacks int32) {
			damageMod.UpdateFloatValue(perStack * float64(newStacks))
		},
	})

	paladin.MakeProcTriggerAura(core.ProcTrigger{
		Name:             "Vengeance - Trigger" + paladin.Label,
		Callback:         core.CallbackOnSpellHitDealt,
		Outcome:          core.OutcomeCrit,
		CanProcFromProcs: true, // 20049 carries the bit: seal crits count.
		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			vengeance.Activate(sim)
			vengeance.AddStack(sim)
		},
	})
}

// Champion of the Light - Increases your spell damage by up to 10/20/30% of your Intellect. Client
// 70170 cut it from 100% and dropped the healing effect, as the WoW: Forever Podcast (episode 2)
// announced: "60% instead of 100, and it doesn't increase healing power anymore". The 16 October 2026
// beta notes halve it again to 10/20/30%; the client row still says 20/40/60, so the share is set here.
const championOfTheLightIntellectShare = 0.10

func (paladin *Paladin) applyChampionOfTheLight() {
	if paladin.Talents.ChampionOfTheLight == 0 {
		return
	}

	share := championOfTheLightIntellectShare * float64(paladin.Talents.ChampionOfTheLight)
	paladin.AddStatDependency(stats.Intellect, stats.SpellDamage, share)
}

// Instrument of Law - Reduces the cast time of your Hammer of Wrath by 0.5/1.0 sec, and reduces all
// threat you generate by 10/20% while Righteous Fury is not active.
func (paladin *Paladin) applyInstrumentOfLaw() {
	if paladin.Talents.InstrumentOfLaw == 0 {
		return
	}

	rank := spellData.InstrumentOfLaw.Highest()

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskHammerOfWrath,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Duration(spellData.InstrumentOfLaw.EffectAt(1).ValueAt(paladin.Talents.InstrumentOfLaw)) * time.Millisecond,
	})

	// The client states the threat reduction as a positive number Righteous Fury zeroes.
	threat := core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    "Instrument of Law" + paladin.Label,
		ActionID: core.ActionID{SpellID: rank.ID},
	}).AttachMultiplicativePseudoStatBuff(
		&paladin.PseudoStats.ThreatMultiplier,
		1-spellData.InstrumentOfLaw.Effect(dbcenums.A_MOD_THREAT, 127).FractionAt(paladin.Talents.InstrumentOfLaw),
	))

	paladin.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Matches(SpellMaskRighteousFury) {
			spell.RelatedSelfBuff.ApplyOnGain(func(_ *core.Aura, sim *core.Simulation) {
				threat.Deactivate(sim)
			}).ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
				threat.Activate(sim)
			})
		}
	})
}

// Twist of Light - Reduces the mana cost of your Seal spells by 20%. When you replace your Seal of
// Command, Seal of Righteousness, Seal of Fury, or Seal of Justice with a different Seal, gain an
// Echo. Your next melee attack applies the replaced Seal's effects, consuming the Echo.
//
// Each seal leaves its own Echo (Echo of Command, of Fury, of Righteousness, of Justice): one
// charge, no duration, consumed by the next auto attack that lands.
func (paladin *Paladin) applyTwistOfLight() {
	if !paladin.Talents.TwistOfLight {
		return
	}

	// The discount is an A_ADD_PCT_MODIFIER on the cost, so it joins the additive bucket the way
	// Swift Judgement's does.
	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskAllSeals,
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		FloatValue: spellData.TwistOfLight.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_COST)).FractionAt(1),
	})

	paladin.echoes = map[int32]*sealEcho{}
	var echoes []*sealEcho
	for _, id := range []int32{echoOfCommandID, echoOfFuryID, echoOfRighteousnessID, echoOfJusticeID} {
		echo := &sealEcho{}
		echo.aura = paladin.RegisterAura(core.Aura{
			Label:    fmt.Sprintf("Echo (%d)%s", id, paladin.Label),
			ActionID: core.ActionID{SpellID: id},
			Duration: core.NeverExpires,
		})
		paladin.echoes[id] = echo
		echoes = append(echoes, echo)
	}

	// One permanent trigger consumes every Echo that is up, in a fixed order. The Echo auras carry
	// no callbacks of their own: an aura that deactivates itself from inside the hit callbacks
	// reshuffles core's callback list mid-loop, which skipped or repeated other handlers depending
	// on the order earlier iterations had left the list in.
	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label: "Twist of Light" + paladin.Label,
	}).AttachProcTrigger(core.ProcTrigger{
		Callback:           core.CallbackOnSpellHitDealt,
		ProcMask:           core.ProcMaskMeleeWhiteHit,
		Outcome:            core.OutcomeLanded,
		TriggerImmediately: true,
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			for _, echo := range echoes {
				if !echo.aura.IsActive() {
					continue
				}
				if echo.seal != nil && echo.seal.echo != nil {
					echo.seal.echo(sim, result.Target)
				}
				echo.seal = nil
				echo.aura.Deactivate(sim)
			}
		},
	}))
}
