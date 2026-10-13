package hunter

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/buffs"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

func (hunter *Hunter) registerSurvivalTalents() {
	// Tier 1
	hunter.registerImprovedTracking()
	hunter.registerDeflection()

	// Tier 2
	hunter.registerEntrapment()
	hunter.registerSavageStrikes()
	hunter.registerSurvivalist()
	hunter.registerImprovedWingClip()

	// Tier 3
	hunter.registerCleverTraps()
	hunter.registerSurefooted()
	hunter.registerDeterrence()

	// Tier 4
	hunter.registerSurvivalTactics()
	hunter.registerPredatorsEdge()
	hunter.registerCounterattack()

	// Tier 5
	hunter.registerResourcefulness()
	hunter.registerExposePrey()
	hunter.registerSurvivalistsDiscipline()
	// Strider Kick: strider_kick.go

	// Tier 6
	hunter.registerLightningReflexes()

	// Tier 7
	// Lacerating Strikes: lacerating_strikes.go
}

func (hunter *Hunter) registerSavageStrikes() {
	if hunter.Talents.SavageStrikes == 0 {
		return
	}

	// 19159's mask (word 1, 4096) also names the Lacerating Strikes bleed (1310536), whose ticks crit.
	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		ClassMask:  HunterSpellsMelee | HunterSpellLaceratingStrikes,
		FloatValue: spellData.SavageStrikes.ValueAt(hunter.Talents.SavageStrikes),
	})
}

func (hunter *Hunter) registerSurvivalist() {
	if hunter.Talents.Survivalist == 0 {
		return
	}

	hunter.MultiplyStat(stats.Health, spellData.Survivalist.MultiplierAt(hunter.Talents.Survivalist))
}

// 19290's melee hit effect carries the 1/2/3% curve. Its spell hit effect has no curve and sits at
// rank 3's 3 for every rank, so spell hit reads the melee curve too, as the tooltip's one number does.
func (hunter *Hunter) registerSurefooted() {
	if hunter.Talents.Surefooted == 0 {
		return
	}

	hit := spellData.Surefooted.Effect(dbcenums.A_MOD_HIT_CHANCE, 0).ValueAt(hunter.Talents.Surefooted)
	hunter.AddStat(stats.PhysicalHitPercent, hit)
	hunter.AddStat(stats.SpellHitPercent, hit)
}

func (hunter *Hunter) registerResourcefulness() {
	if hunter.Talents.Resourcefulness == 0 {
		return
	}

	// 440529's cost mask leaves out Strider Kick.
	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		ClassMask:  HunterSpellsTraps | HunterSpellsMelee&^HunterSpellStriderKick,
		FloatValue: spellData.Resourcefulness.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_COST)).FractionAt(hunter.Talents.Resourcefulness),
	})

	// The buff (1242688) is 50% mana regen while casting for 30 sec at both ranks; the points buy
	// the proc chance, 50/100%.
	buff := spellData.ResourcefulnessTriggered.Highest()
	regen := buff.Effect(dbcenums.A_MOD_MANA_REGEN_INTERRUPT, 0).Average(core.CharacterLevel) / 100

	procAura := hunter.RegisterAura(core.Aura{
		Label:    "Resourcefulness",
		ActionID: core.ActionID{SpellID: buff.ID},
		Duration: buff.Duration(),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.PseudoStats.SpiritRegenRateCasting += regen
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.PseudoStats.SpiritRegenRateCasting -= regen
		},
	})

	hunter.MakeProcTriggerAura(core.ProcTrigger{
		Name:       "Resourcefulness Trigger",
		Callback:   core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.5 * float64(hunter.Talents.Resourcefulness),
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			procAura.Activate(sim)
		},
	})
}

// The client curve reads 3% a rank, where our sim had 2%.
func (hunter *Hunter) registerLightningReflexes() {
	if hunter.Talents.LightningReflexes == 0 {
		return
	}

	hunter.MultiplyStat(stats.Agility, spellData.LightningReflexes.MultiplierAt(hunter.Talents.LightningReflexes))
}

func (hunter *Hunter) registerImprovedTracking() {
	if hunter.Talents.ImprovedTracking == 0 {
		return
	}

	// Everything a raid encounter can be is trackable apart from Mechanical. Damage only: 24293 has
	// no crit damage effect. Unknown counts: it's the site's default "Level 60" target, a boss of no
	// set type, and the hunter tracks whatever the boss is (pick Mechanical to turn it off).
	multiplier := spellData.ImprovedTracking.MultiplierAt(hunter.Talents.ImprovedTracking)
	hunter.Env.RegisterPostFinalizeEffect(func() {
		for _, t := range hunter.Env.Encounter.AllTargets {
			switch t.MobType {
			case proto.MobType_MobTypeUnknown, proto.MobType_MobTypeBeast, proto.MobType_MobTypeDemon,
				proto.MobType_MobTypeDragonkin, proto.MobType_MobTypeElemental, proto.MobType_MobTypeGiant,
				proto.MobType_MobTypeHumanoid, proto.MobType_MobTypeUndead:
				at := hunter.AttackTables[t.UnitIndex]
				at.DamageDealtMultiplier *= multiplier
			}
		}
	})
}

func (hunter *Hunter) registerDeflection() {
	if hunter.Talents.Deflection == 0 {
		return
	}

	hunter.PseudoStats.BaseParryChance += spellData.Deflection.FractionAt(hunter.Talents.Deflection)
}

func (hunter *Hunter) registerCleverTraps() {
	if hunter.Talents.CleverTraps == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  HunterSpellsTraps,
		FloatValue: spellData.CleverTraps.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_ALL_EFFECTS)).FractionAt(hunter.Talents.CleverTraps),
	})
}

func (hunter *Hunter) registerSurvivalTactics() {
	if hunter.Talents.SurvivalTactics == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusHit_Percent,
		ClassMask:  HunterSpellsTraps,
		FloatValue: spellData.SurvivalTactics.ValueAt(hunter.Talents.SurvivalTactics),
	})
}

func (hunter *Hunter) registerPredatorsEdge() {
	if hunter.Talents.PredatorsEdge == 0 {
		return
	}

	// 1310627's crit damage mask is the melee abilities and the Lacerating Strikes bleed (word 1, 4096):
	// no auto attacks, no hawks.
	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_CritMultiplier_Flat,
		ClassMask:  HunterSpellsMelee | HunterSpellLaceratingStrikes,
		FloatValue: spellData.PredatorsEdge.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_CRIT_DAMAGE_BONUS)).FractionAt(hunter.Talents.PredatorsEdge),
	})

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ProcMask:   core.ProcMaskMeleeOH,
		FloatValue: spellData.PredatorsEdge.Effect(dbcenums.A_MOD_OFFHAND_DAMAGE_PCT, 0).FractionAt(hunter.Talents.PredatorsEdge),
	})
}

func (hunter *Hunter) registerSurvivalistsDiscipline() {
	if hunter.Talents.SurvivalistsDiscipline == 0 {
		return
	}

	hunter.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_Cooldown_Multiplier,
		ClassMask:  HunterSpellsTraps,
		FloatValue: spellData.SurvivalistsDiscipline.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_COOLDOWN)).MultiplierAt(hunter.Talents.SurvivalistsDiscipline),
	})
}

// Expose Prey opens the Mongoose Bite window off a landed melee or ranged attack on a target with
// Hunter's Mark (client 1310532 ProcTypeMask 340 = melee/ranged autos and specials; Wowhead Forever:
// "targets with Hunter's Mark"), where only a dodge opens it otherwise.
func (hunter *Hunter) registerExposePrey() {
	if hunter.Talents.ExposePrey == 0 {
		return
	}

	hunter.MakeProcTriggerAura(core.ProcTrigger{
		Name:       "Expose Prey",
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		ProcChance: spellData.ExposePrey.FractionAt(hunter.Talents.ExposePrey),
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && result.Target.HasActiveAuraWithTag(buffs.HuntersMarkCategory) {
				hunter.DefensiveState.Activate(sim)
			}
		},
	})
}

// registerEntrapment implements Entrapment, new in Forever.
//
// TODO: a root on trap targets; bosses are immune.
func (hunter *Hunter) registerEntrapment() {
	if hunter.Talents.Entrapment == 0 {
		return
	}
}

// registerImprovedWingClip implements Improved Wing Clip, new in Forever.
//
// TODO: a root chance on Wing Clip; bosses are immune.
func (hunter *Hunter) registerImprovedWingClip() {
	if hunter.Talents.ImprovedWingClip == 0 {
		return
	}
}

// registerDeterrence implements Deterrence, new in Forever.
//
// TODO: a defensive cooldown, 25% parry for 10 sec; no effect on damage.
func (hunter *Hunter) registerDeterrence() {
	if !hunter.Talents.Deterrence {
		return
	}
}

// registerCounterattack implements Counterattack, new in Forever.
//
// TODO: a parry-window strike the sim's rotations never reach, as the boss is the one parrying.
func (hunter *Hunter) registerCounterattack() {
	if !hunter.Talents.Counterattack {
		return
	}
}
