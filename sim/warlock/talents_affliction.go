package warlock

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/stats"
)

func (warlock *Warlock) registerAfflictionTalents() {
	// Tier 1
	// Improved Life Tap: lifetap.go
	warlock.applySuppression()
	warlock.applyImprovedCorruption()

	// Tier 2
	warlock.applyMalediction()
	warlock.applySoulHarvesting()
	warlock.applyImprovedDrains()

	// Tier 3
	warlock.applyImprovedBaneOfAgony()
	warlock.applyFelConcentration()
	warlock.registerAmplifyCurse()
	warlock.applyPandemic()

	// Tier 4
	warlock.applyMalevolence()
	warlock.applyNightfall()
	warlock.applyCurseOfExhaustion()

	// Tier 5
	// Siphon Life: siphon_life.go
	// Soul Siphon: drain_life.go

	// Tier 6
	warlock.applyShadowMastery()

	// Tier 7
	// Wrack: wrack.go
}

// 1% hit a point on every school the warlock casts (A_MOD_SPELL_HIT_CHANCE with no family
// restriction) plus 4% less threat a point.
func (warlock *Warlock) applySuppression() {
	if warlock.Talents.Suppression == 0 {
		return
	}

	points := warlock.Talents.Suppression
	warlock.AddStat(stats.SpellHitPercent, spellData.Suppression.Effect(dbcenums.A_MOD_SPELL_HIT_CHANCE, 0).ValueAt(points))
	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_ThreatMultiplier_Pct,
		FloatValue: spellData.Suppression.Effect(dbcenums.A_MOD_THREAT, 127).FractionAt(points),
		ClassMask:  WarlockSpellAll,
	})
}

// A cast time cut and, in Forever, 2% more dot damage a point.
func (warlock *Warlock) applyImprovedCorruption() {
	if warlock.Talents.ImprovedCorruption == 0 {
		return
	}

	points := warlock.Talents.ImprovedCorruption
	warlock.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Millisecond * time.Duration(spellData.ImprovedCorruption.Effect(dbcenums.A_ADD_FLAT_MODIFIER, int32(dbcenums.SPELLMOD_CASTING_TIME)).ValueAt(points)),
		ClassMask: WarlockSpellCorruption,
	})
	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DotDamageDone_Pct,
		FloatValue: spellData.ImprovedCorruption.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DOT)).FractionAt(points),
		ClassMask:  WarlockSpellCorruption,
	})
}

// Forever moved Malediction off the curse: the beta client's 1225177 is a flat 1% a point on the
// warlock's own damage and dot damage, not a bonus on Curse of the Elements.
func (warlock *Warlock) applyMalediction() {
	if warlock.Talents.Malediction == 0 {
		return
	}

	// Periodic damage only ("Increases all periodic damage done"). The row also states the same
	// value on op 0, but a tick already takes DamageDone, so taking both would count it twice.
	// Hellfire's area hits are Hellfire Effect (11682), a direct School Damage effect, so in the
	// client the dot modifier only reached the self-burn. Dev AidanZMoon called that a bug, 'Fixed
	// in next build' (forever-bugs #311, 2026-10-13), so the area hits take it too.
	points := warlock.Talents.Malediction
	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DotDamageDone_Pct,
		FloatValue: spellData.Malediction.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DOT)).FractionAt(points),
		ClassMask:  WarlockSpellAll,
	})
}

// Improved Drains: 7/13/20% more drain damage (403511, a dot modifier).
func (warlock *Warlock) applyImprovedDrains() {
	if warlock.Talents.ImprovedDrains == 0 {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DotDamageDone_Pct,
		FloatValue: spellData.ImprovedDrains.FractionAt(warlock.Talents.ImprovedDrains),
		ClassMask:  WarlockDrainSpells,
	})
}

// Improved Bane of Agony: 5/10% more periodic damage (18827).
func (warlock *Warlock) applyImprovedBaneOfAgony() {
	if warlock.Talents.ImprovedBaneOfAgony == 0 {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DotDamageDone_Pct,
		FloatValue: spellData.ImprovedBaneOfAgony.FractionAt(warlock.Talents.ImprovedBaneOfAgony),
		ClassMask:  WarlockSpellCurseOfAgony,
	})
}

// Pandemic lets the Affliction dots crit for 33/67/100% more, the beta client's SPELLMOD_CRIT_DAMAGE
// ladder (427712). Forever's core already lets a dot crit.
func (warlock *Warlock) applyPandemic() {
	if warlock.Talents.Pandemic == 0 {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_CritMultiplier_Flat,
		FloatValue: spellData.Pandemic.FractionAt(warlock.Talents.Pandemic),
		ClassMask:  WarlockPeriodicShadowDamage,
	})
}

// 1% shadow crit a point (1310949).
func (warlock *Warlock) applyMalevolence() {
	if warlock.Talents.Malevolence == 0 {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: spellData.Malevolence.ValueAt(warlock.Talents.Malevolence),
		ClassMask:  WarlockShadowDamage,
	})
}

func (warlock *Warlock) applyNightfall() {
	if warlock.Talents.Nightfall == 0 {
		return
	}

	warlock.NightfallProcAura = warlock.MakeProcTriggerAura(core.ProcTrigger{
		Name:            "Shadow Trance",
		MetricsActionID: core.ActionID{SpellID: spellData.NightfallTriggered.Highest().ID},
		Duration:        spellData.NightfallTriggered.Highest().Duration(),
		ClassSpellMask:  WarlockSpellShadowBolt,
		Callback:        core.CallbackOnCastComplete,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.CurCast.CastTime != 0 {
				return
			}
			warlock.NightfallProcAura.Deactivate(sim)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -1.0,
		ClassMask:  WarlockSpellShadowBolt,
	})

	warlock.MakeProcTriggerAura(core.ProcTrigger{
		Name:           "Nightfall",
		ClassSpellMask: WarlockNightfallSpells,
		// The per-rank chance is on the effect; ProcChanceAt reads the row's flat 100%.
		ProcChance: spellData.Nightfall.FractionAt(warlock.Talents.Nightfall),
		Callback:   core.CallbackOnPeriodicDamageDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warlock.NightfallProcAura.Activate(sim)
		},
	})
}

// Shadow Mastery is 1% a point on both damage and dot damage in Forever: 18271 kept its op 0 and
// op 22 modifiers and dropped Classic's op 8, so every shadow spell takes it the same way. Neither
// mask names Life Tap, so its mana is left alone.
func (warlock *Warlock) applyShadowMastery() {
	if warlock.Talents.ShadowMastery == 0 {
		return
	}

	points := warlock.Talents.ShadowMastery
	warlock.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.ShadowMastery.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DAMAGE)).FractionAt(points),
		ClassMask:  WarlockShadowDamage,
	})
	// The row's op 22 (dot) modifier carries the same value as op 0: one bonus stated for both
	// halves of a spell, not a second one for ticks, which DamageDone already reaches.
}

func (warlock *Warlock) registerAmplifyCurse() {
	if !warlock.Talents.AmplifyCurse {
		return
	}

	rank := spellData.AmplifyCurse.Highest()
	actionID := core.ActionID{SpellID: rank.ID}

	// Spent by Bane of Agony, the only spell in 18288's mask the sim casts; see agony.go.
	warlock.AmplifyCurseAura = warlock.GetOrRegisterAura(core.Aura{
		Label:    "Amplify Curse",
		ActionID: actionID,
		Duration: rank.Duration(),
	})

	warlock.AmplifyCurse = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: rank.SpellSchool(),
		Flags:       core.SpellFlagAPL,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: max(rank.Cooldown(), rank.CategoryCooldown()),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warlock.AmplifyCurseAura.Activate(sim)
		},
	})

	warlock.AddMajorCooldown(core.MajorCooldown{
		Spell: warlock.AmplifyCurse,
		Type:  core.CooldownTypeDPS,
	})
}

// applySoulHarvesting implements Soul Harvest (Soul Harvesting before build 70170), new in Forever.
//
// Client 437032: a kill under Drain Soul grants Soul Harvest (1242853) for 10 sec, +50/100% mana
// regeneration and 50/100% of it kept while casting. Nothing dies under Drain Soul in the sim's
// encounters, so there is nothing to model.
func (warlock *Warlock) applySoulHarvesting() {
	if warlock.Talents.SoulHarvest == 0 {
		return
	}
}

// applyFelConcentration implements Fel Concentration, new in Forever.
//
// Pushback protection (SPELLMOD_NOT_LOSE_CASTING_TIME, 23/47/70%) has nothing to act on in the sim.
func (warlock *Warlock) applyFelConcentration() {
	if warlock.Talents.FelConcentration == 0 {
		return
	}
}

// applyCurseOfExhaustion implements Curse of Exhaustion, new in Forever.
//
// A 30% movement slow (18223); nothing to model.
func (warlock *Warlock) applyCurseOfExhaustion() {
	if !warlock.Talents.CurseOfExhaustion {
		return
	}
}
