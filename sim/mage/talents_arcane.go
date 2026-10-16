package mage

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/stats"
)

func (mage *Mage) registerArcaneTalents() {
	// Tier 1
	mage.registerWandSpecialization()
	mage.registerArcaneFocus()
	mage.registerImprovedChanneling()

	// Tier 2
	mage.registerArcaneSubtlety()
	mage.registerMagicAbsorption()
	mage.registerArcaneConcentration()
	mage.registerArcaneResilience()

	// Tier 3
	mage.registerArcaneGeometry()
	mage.registerArcaneImpact()
	// Arcane Blast: arcane_blast.go and arcane_charge.go

	// Tier 4
	mage.registerArcaneShielding()
	mage.registerImprovedCounterspell()
	mage.registerArcaneMeditation()
	mage.registerMissileBarrage()

	// Tier 5
	// Presence of Mind: presence_of_mind.go
	mage.registerArcaneMind()

	// Tier 6
	mage.registerArcaneInstability()

	// Tier 7
	// Arcane Power: arcane_power.go
}

// registerWandSpecialization implements Wand Specialization, new in Forever.
//
// Not modelled: 6057 raises wand damage 13/25%, and the sim has no wand attack (see the priest's).
func (mage *Mage) registerWandSpecialization() {
	if mage.Talents.WandSpecialization == 0 {
		return
	}
}

func (mage *Mage) registerArcaneFocus() {
	if mage.Talents.ArcaneFocus == 0 {
		return
	}

	mage.PseudoStats.SchoolBonusHitChance[stats.SchoolIndexArcane] += spellData.ArcaneFocus.ValueAt(mage.Talents.ArcaneFocus)
}

// registerImprovedChanneling implements Improved Channeling, new in Forever.
//
// Client 11237 (row values): 20% to keep channeling Arcane Missiles through damage taken, 14% for Arcane
// Blast. The sim's mage takes no damage, so it changes no number.
func (mage *Mage) registerImprovedChanneling() {
	if mage.Talents.ImprovedChanneling == 0 {
		return
	}
}

func (mage *Mage) registerArcaneSubtlety() {
	if mage.Talents.ArcaneSubtlety == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll,
		School:     core.SpellSchoolArcane,
		FloatValue: spellData.ArcaneSubtlety.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_THREAT)).FractionAt(mage.Talents.ArcaneSubtlety),
		Kind:       core.SpellMod_ThreatMultiplier_Pct,
	})

	// The client states the penetration as a negative target resistance.
	mage.AddStat(stats.SpellPiercing, -spellData.ArcaneSubtlety.Effect(dbcenums.A_MOD_TARGET_RESISTANCE, 126).ValueAt(mage.Talents.ArcaneSubtlety))
}

// Only the resistance half; the mana returned on a resisted spell is not modelled.
func (mage *Mage) registerMagicAbsorption() {
	if mage.Talents.MagicAbsorption == 0 {
		return
	}

	resist := spellData.MagicAbsorption.Effect(dbcenums.A_MOD_RESISTANCE, 124).ValueAt(mage.Talents.MagicAbsorption)
	mage.AddStats(stats.Stats{
		stats.ArcaneResistance: resist,
		stats.FireResistance:   resist,
		stats.FrostResistance:  resist,
		stats.NatureResistance: resist,
		stats.ShadowResistance: resist,
	})
}

func (mage *Mage) registerArcaneConcentration() {
	if mage.Talents.ArcaneConcentration == 0 {
		return
	}

	clearcastingRank := spellData.ArcaneConcentrationTriggered.Highest()

	mage.ClearcastingAura = mage.RegisterAura(core.Aura{
		Label:    "Clearcasting",
		ActionID: core.ActionID{SpellID: clearcastingRank.ID},
		Duration: clearcastingRank.Duration(),
		OnGain: func(aura *core.Aura, _ *core.Simulation) {
			aura.Unit.PseudoStats.SpellCostPercentModifier -= 100
		},
		OnExpire: func(aura *core.Aura, _ *core.Simulation) {
			aura.Unit.PseudoStats.SpellCostPercentModifier += 100
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// OnCastComplete runs after the hit that may have just procced this, so a fresh aura
			// belongs to the next cast.
			if aura.RemainingDuration(sim) == aura.Duration {
				return
			}
			if !spell.Matches(MageSpellsAll) || spell.Cost == nil || spell.Cost.BaseCost == 0 {
				return
			}
			aura.Deactivate(sim)
		},
	})

	// Forever states a flat SpellAuraOptions.ProcChance of 100 on the talent spell and puts the
	// real per-rank chance on the effect, so ProcChanceAt would read 100% at every rank. 11213's
	// ProcCategoryRecovery holds it to one proc a second.
	//
	// Spells another spell triggers do not proc it: Arcane Missiles' missiles and Blizzard's ticks.
	// Blizzard rolls once per enemy when cast instead. Beta log 2706 (Jamal, 18 Arcane, level 26):
	// 70 procs on 246 Arcane Explosions (722 targets hit, 11% a target), 12 on 37 Blizzard casts,
	// 0 on 642 Blizzard tick hits (62 if they rolled) and 0 on 32 missiles.
	mage.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Arcane Concentration",
		Callback:           core.CallbackOnSpellHitDealt,
		ClassSpellMask:     MageSpellsAllDamaging &^ MageSpellArcaneMissilesTick,
		Outcome:            core.OutcomeLanded,
		ProcChance:         spellData.ArcaneConcentration.EffectAt(1).FractionAt(mage.Talents.ArcaneConcentration),
		ICD:                spellData.ArcaneConcentration.Highest().ICD(),
		TriggerImmediately: true,
		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			mage.ClearcastingAura.Activate(sim)
		},
	})
}

func (mage *Mage) registerArcaneResilience() {
	if mage.Talents.ArcaneResilience == 0 {
		return
	}

	mage.AddStatDependency(stats.Intellect, stats.Armor, spellData.ArcaneResilience.FractionAt(mage.Talents.ArcaneResilience))
}

// registerArcaneGeometry implements Arcane Geometry, new in Forever: 11247 adds 3/6 yards to the
// range of the spells its class mask reaches, of the mage's Arcane Blast and Arcane Missiles.
func (mage *Mage) registerArcaneGeometry() {
	if mage.Talents.ArcaneGeometry == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellArcaneBlast | MageSpellArcaneMissilesCast,
		FloatValue: spellData.ArcaneGeometry.Effect(dbcenums.A_ADD_FLAT_MODIFIER, int32(dbcenums.SPELLMOD_RANGE)).ValueAt(mage.Talents.ArcaneGeometry),
		Kind:       core.SpellMod_Range_Flat,
	})
}

// Every arcane spell of the mage's, not TBC's Arcane Blast and Arcane Explosion only.
func (mage *Mage) registerArcaneImpact() {
	if mage.Talents.ArcaneImpact == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll,
		School:     core.SpellSchoolArcane,
		FloatValue: spellData.ArcaneImpact.ValueAt(mage.Talents.ArcaneImpact),
		Kind:       core.SpellMod_BonusCrit_Percent,
	})
}

// registerArcaneShielding implements Arcane Shielding, new in Forever.
//
// Not modelled: 11252 cuts Mana Shield's mana per damage 17/33% and raises Mage Armor's
// resistances 25/50%; the sim casts neither.
func (mage *Mage) registerArcaneShielding() {
	if mage.Talents.ArcaneShielding == 0 {
		return
	}
}

// registerImprovedCounterspell implements Improved Counterspell, new in Forever.
//
// Not modelled: 11255 adds a 2/4 sec silence to Counterspell; the sim casts no interrupts.
func (mage *Mage) registerImprovedCounterspell() {
	if mage.Talents.ImprovedCounterspell == 0 {
		return
	}
}

func (mage *Mage) registerArcaneMeditation() {
	if mage.Talents.ArcaneMeditation == 0 {
		return
	}

	mage.PseudoStats.SpiritRegenRateCasting += spellData.ArcaneMeditation.FractionAt(mage.Talents.ArcaneMeditation)
}

// Arcane Blast landing (20/40% at 1/2 points), an Arcane Missiles missile hitting (a quarter of it), or
// Arcane Wrath, Fireball, Frostbolt or Frostfire Bolt (half of it) can make the next Arcane Missiles
// free and fire its missiles every 0.5 sec instead of every second: the same count in half the channel.
// Client 70425 (16 October 2026 beta) gave the talent a second rank and added the missile and Arcane
// Wrath procs. 400588 procs on a landed harmful spell (ProcFlags 0x10000), so a miss can't, and a bolt
// procs when it arrives. Its trigger effect holds Arcane Blast's chance; the tooltip's quarter and half
// are fractions of it the row doesn't state. A missile of the channel that spends the buff can proc it
// again for the next one. The buff (400589) holds the duration, cost and tick changes. 44404 stays the
// aura id because saved rotations key on it.
func (mage *Mage) registerMissileBarrage() {
	if mage.Talents.MissileBarrage == 0 {
		return
	}

	buff := spellData.MissileBarrageTriggered.Highest()
	arcaneBlastChance := spellData.MissileBarrage.Effect(dbcenums.A_PROC_TRIGGER_SPELL, 0).ValueAt(mage.Talents.MissileBarrage) / 100

	mage.MissileBarrageAura = mage.RegisterAura(core.Aura{
		Label:    "Missile Barrage",
		ActionID: core.ActionID{SpellID: 44404},
		Duration: buff.Duration(),
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(MageSpellArcaneMissilesCast) {
				aura.Deactivate(sim)
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask:  MageSpellArcaneMissilesCast,
		FloatValue: buff.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_COST)).Average(core.CharacterLevel) / 100,
		Kind:       core.SpellMod_PowerCost_Pct_Add,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: MageSpellArcaneMissilesCast,
		TimeValue: time.Duration(buff.Effect(dbcenums.A_ADD_FLAT_MODIFIER, int32(dbcenums.SPELLMOD_ACTIVATION_TIME)).Average(core.CharacterLevel)) * time.Millisecond,
		Kind:      core.SpellMod_DotTickLength_Flat,
	})

	mage.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Missile Barrage Trigger",
		Callback:           core.CallbackOnSpellHitDealt,
		ClassSpellMask:     MageSpellArcaneBlast | MageSpellArcaneMissilesTick | MageSpellArcaneWrath | MageSpellFireball | MageSpellFrostbolt | MageSpellFrostfireBolt,
		Outcome:            core.OutcomeLanded,
		TriggerImmediately: true,
		Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			procChance := arcaneBlastChance
			switch {
			case spell.Matches(MageSpellArcaneBlast):
			case spell.Matches(MageSpellArcaneMissilesTick):
				procChance /= 4
			default:
				procChance /= 2
			}
			if sim.Proc(procChance, "Missile Barrage") {
				mage.MissileBarrageAura.Activate(sim)
			}
		},
	})
}

func (mage *Mage) registerArcaneMind() {
	if mage.Talents.ArcaneMind == 0 {
		return
	}

	mage.MultiplyStat(stats.Intellect, spellData.ArcaneMind.Effect(dbcenums.A_MOD_TOTAL_STAT_PERCENTAGE, 0).MultiplierAt(mage.Talents.ArcaneMind))
	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll,
		School:     core.SpellSchoolArcane,
		FloatValue: spellData.ArcaneMind.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_CRIT_DAMAGE_BONUS)).FractionAt(mage.Talents.ArcaneMind),
		Kind:       core.SpellMod_CritMultiplier_Flat,
	})
}

func (mage *Mage) registerArcaneInstability() {
	if mage.Talents.ArcaneInstability == 0 {
		return
	}

	damage := spellData.ArcaneInstability.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DAMAGE)).FractionAt(mage.Talents.ArcaneInstability)
	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll &^ MageSpellFrostfireBolt,
		FloatValue: damage,
		Kind:       core.SpellMod_DamageDone_Flat,
	})
	// The SPELLMOD_DOT mask leaves out Frostfire Bolt: its hit takes the bonus, its DoT does not.
	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellFrostfireBolt,
		FloatValue: damage,
		Kind:       core.SpellMod_DirectDamageDone_Flat,
	})
	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll,
		FloatValue: spellData.ArcaneInstability.Effect(dbcenums.A_MOD_CRIT_PCT, 0).ValueAt(mage.Talents.ArcaneInstability),
		Kind:       core.SpellMod_BonusCrit_Percent,
	})
}
