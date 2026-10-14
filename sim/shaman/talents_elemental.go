package shaman

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/stats"
)

func (shaman *Shaman) registerElementalTalents() {
	// Tier 1
	shaman.applyConvection()
	shaman.applyConcussion()

	// Tier 2
	shaman.applyElementalWarding()
	shaman.applyReverberation()
	shaman.applyCallOfFlame()
	shaman.applyElementalDevastation()

	// Tier 3
	shaman.applyElementalFocus()
	shaman.applyElementalAlacrity()

	// Tier 4
	shaman.applyImprovedFireNova()
	shaman.applyEyeOfTheStorm()
	shaman.applyCallOfThunder()

	// Tier 5
	shaman.applyElementalReach()
	shaman.applyLightningOverload()
	shaman.applyEarthbound()

	// Tier 6
	shaman.applyElementalFury()

	// Tier 7
	shaman.applyLavaBurst()
}

func (shaman *Shaman) applyCallOfFlame() {
	if shaman.Talents.CallOfFlame == 0 {
		return
	}

	shaman.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.CallOfFlame.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DAMAGE)).FractionAt(shaman.Talents.CallOfFlame),
		ClassMask:  SpellMaskFireTotem | SpellMaskFlameShock | SpellMaskFireNova | SpellMaskLavaBurst,
	})
}

func (shaman *Shaman) applyCallOfThunder() {
	if !shaman.Talents.CallOfThunder {
		return
	}

	shaman.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: spellData.CallOfThunder.Effect(dbcenums.A_ADD_FLAT_MODIFIER, int32(dbcenums.SPELLMOD_CRITICAL_CHANCE)).ValueAt(1),
		ClassMask:  SpellMaskLightningBolt | SpellMaskChainLightning | SpellMaskOverload,
	})
}

func (shaman *Shaman) applyConcussion() {
	if shaman.Talents.Concussion == 0 {
		return
	}

	shaman.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.Concussion.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DAMAGE)).FractionAt(shaman.Talents.Concussion),
		// Client 16035's mask is Lightning Bolt, Chain Lightning and Earth Shock: not Flame or Frost Shock.
		ClassMask: SpellMaskLightningBolt | SpellMaskChainLightning | SpellMaskOverload | SpellMaskEarthShock,
	})
}

func (shaman *Shaman) applyConvection() {
	if shaman.Talents.Convection == 0 {
		return
	}

	shaman.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		FloatValue: spellData.Convection.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_COST)).FractionAt(shaman.Talents.Convection),
		ClassMask:  SpellMaskLightningBolt | SpellMaskChainLightning | SpellMaskOverload | SpellMaskShock | SpellMaskLavaBurst,
	})
}

func (shaman *Shaman) applyElementalDevastation() {
	if shaman.Talents.ElementalDevastation == 0 {
		return
	}

	critBuffAura := shaman.RegisterAura(core.Aura{
		Label:    "Elemental Devastation",
		ActionID: core.ActionID{SpellID: 29178},
		Duration: time.Second * 10,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: spellData.ElementalDevastation.Effect(dbcenums.A_DUMMY, 0).ValueAt(shaman.Talents.ElementalDevastation),
		ProcMask:   core.ProcMaskMelee,
	})
	shaman.MakeProcTriggerAura(core.ProcTrigger{
		Name:             "Elemental Devastation Trigger",
		CanProcFromProcs: spellData.ElementalDevastation.Highest().CanProcFromProcs(), // Forever's 30160 lacks the bit (Era's 29179/29180 carry it): overload crits don't count.
		Callback:         core.CallbackOnSpellHitDealt,
		ProcMask:         core.ProcMaskSpellDamage,
		Outcome:          core.OutcomeCrit,
		// Totems and what they cast never trigger it, Fire Nova included (dev AidanZMoon, forever-bugs #334).
		ExtraCondition: func(_ *core.Simulation, spell *core.Spell, _ *core.SpellResult) bool {
			return !spell.Matches(SpellMaskFireNova)
		},
		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			critBuffAura.Activate(sim)
		},
	})
}

func (shaman *Shaman) applyElementalFocus() {
	if !shaman.Talents.ElementalFocus {
		return
	}

	var triggeringSpell *core.Spell
	var triggerTime time.Duration

	// 16246's class mask: Lightning Bolt, Chain Lightning, Lava Burst, the shocks and Fire Nova (408345).
	canConsumeSpells := SpellMaskLightningBolt | SpellMaskChainLightning | SpellMaskLavaBurst | SpellMaskFireNova | (SpellMaskShock & ^SpellMaskFlameShockDot)

	clearcasting := spellData.ElementalFocusTriggered.Highest()
	maxStacks := int32(clearcasting.ProcCharges)

	clearcastingAura := shaman.RegisterAura(core.Aura{
		Label:     "Clearcasting",
		ActionID:  core.ActionID{SpellID: 16246},
		Duration:  clearcasting.Duration(),
		MaxStacks: maxStacks,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Matches(canConsumeSpells) {
				return
			}
			if spell == triggeringSpell && sim.CurrentTime == triggerTime {
				return
			}
			aura.RemoveStack(sim)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		ClassMask:  canConsumeSpells,
		FloatValue: clearcasting.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_COST)).Percent(),
	})

	// Client 16164: "a chance to enter a Clearcasting state after casting any Fire, Frost, or Nature
	// damage spell", so it rolls when the cast completes, hit or miss, and a proc off a Lightning Bolt
	// is there for the next cast to see instead of arriving mid-cast with the missile.
	shaman.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Elemental Focus",
		ProcChance:         float64(spellData.ElementalFocus.Rank(1).ProcChance) / 100,
		Callback:           core.CallbackOnCastComplete,
		ProcMask:           core.ProcMaskSpellDamage,
		CanProcFromProcs:   true, // 16164 carries the bit.
		TriggerImmediately: true,

		Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			// Searing and Magma Totem attacks are the totem's, not the shaman's: they never proc the shaman's talents.
			// Flame Shock's periodic half is cast alongside its hit, one cast of the spell.
			if !spell.SpellSchool.Matches(core.SpellSchoolElemental) || spell.Matches(SpellMaskFireTotem|SpellMaskFlameShockDot) {
				return
			}
			triggeringSpell = spell
			triggerTime = sim.CurrentTime
			clearcastingAura.Activate(sim)
			clearcastingAura.SetStacks(sim, maxStacks)
		},
	})
}

func (shaman *Shaman) applyElementalFury() {
	if shaman.Talents.ElementalFury == 0 {
		return
	}

	// The talent's class mask (16089) also covers Flametongue Attack (bit 21) and Frostbrand
	// Attack (bit 24): shamans' Flametongue Weapon hits crit for 2.0x in logs while the same
	// attack granted by Flametongue Totem crits for 1.5x on other players. Bit 10 is Lightning
	// Shield, the orbs (26363..26370) included.
	shaman.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_CritMultiplier_Flat,
		FloatValue: spellData.ElementalFury.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_CRIT_DAMAGE_BONUS)).FractionAt(shaman.Talents.ElementalFury),
		ClassMask:  SpellMaskFireTotem | SpellMaskFire | SpellMaskNature | SpellMaskFrost | SpellMaskFlametongueWeapon | SpellMaskFrostbrandWeapon | SpellMaskLightningShield,
	})
}
func (shaman *Shaman) applyLightningOverload() {
	if shaman.Talents.LightningOverload == 0 {
		return
	}
	// In shaman.go -> GetOverloadChance()
}

func (shaman *Shaman) applyReverberation() {
	if shaman.Talents.Reverberation == 0 {
		return
	}

	shaman.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: time.Millisecond * time.Duration(spellData.Reverberation.Effect(dbcenums.A_ADD_FLAT_MODIFIER, int32(dbcenums.SPELLMOD_COOLDOWN)).ValueAt(shaman.Talents.Reverberation)),
		ClassMask: SpellMaskShock,
	})
}

// applyElementalAlacrity implements Elemental Alacrity, new in Forever: a flat cast-time cut.
func (shaman *Shaman) applyElementalAlacrity() {
	if shaman.Talents.ElementalAlacrity == 0 {
		return
	}

	shaman.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Millisecond * time.Duration(spellData.ElementalAlacrity.Effect(dbcenums.A_ADD_FLAT_MODIFIER, int32(dbcenums.SPELLMOD_CASTING_TIME)).ValueAt(shaman.Talents.ElementalAlacrity)),
		ClassMask: SpellMaskLightningBolt | SpellMaskChainLightning | SpellMaskLavaBurst,
	})
}

// applyElementalReach implements Elemental Reach, new in Forever: 28999 adds 3/6 yards to Lightning
// Bolt, Chain Lightning and Lava Burst range and 8/15 to Flame Shock's (its two class masks). The first
// mask also reaches Fire Nova, which centres on the shaman and has no range here.
func (shaman *Shaman) applyElementalReach() {
	if shaman.Talents.ElementalReach == 0 {
		return
	}

	shaman.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskLightningBolt | SpellMaskChainLightning | SpellMaskLavaBurst,
		FloatValue: spellData.ElementalReach.EffectAt(1).ValueAt(shaman.Talents.ElementalReach),
		Kind:       core.SpellMod_Range_Flat,
	})
	shaman.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskFlameShock,
		FloatValue: spellData.ElementalReach.EffectAt(2).ValueAt(shaman.Talents.ElementalReach),
		Kind:       core.SpellMod_Range_Flat,
	})
}

// applyElementalWarding implements Elemental Warding, new in Forever: less fire, frost and nature
// damage taken (the client states one modifier over the three-school mask 28).
func (shaman *Shaman) applyElementalWarding() {
	if shaman.Talents.ElementalWarding == 0 {
		return
	}

	multiplier := spellData.ElementalWarding.Effect(dbcenums.A_MOD_DAMAGE_PERCENT_TAKEN, 28).MultiplierAt(shaman.Talents.ElementalWarding)
	for _, school := range []stats.SchoolIndex{stats.SchoolIndexFire, stats.SchoolIndexFrost, stats.SchoolIndexNature} {
		shaman.PseudoStats.SchoolDamageTakenMultiplier[school] *= multiplier
	}
}

// applyEyeOfTheStorm is Forever's pushback-resistance talent (29062, 23/47/70%). The sim models no
// spell pushback, so it changes no number here.
func (shaman *Shaman) applyEyeOfTheStorm() {
	if shaman.Talents.EyeOfTheStorm == 0 {
		return
	}
}

// applyImprovedFireNova implements Improved Fire Nova, new in Forever: more Fire Nova damage and a
// shorter cooldown on it.
func (shaman *Shaman) applyImprovedFireNova() {
	if shaman.Talents.ImprovedFireNova == 0 {
		return
	}

	shaman.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: spellData.ImprovedFireNova.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DAMAGE)).FractionAt(shaman.Talents.ImprovedFireNova),
		ClassMask:  SpellMaskFireNova,
	})
	shaman.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: time.Millisecond * time.Duration(spellData.ImprovedFireNova.Effect(dbcenums.A_ADD_FLAT_MODIFIER, int32(dbcenums.SPELLMOD_COOLDOWN)).ValueAt(shaman.Talents.ImprovedFireNova)),
		ClassMask: SpellMaskFireNova,
	})
}

// applyEarthbound implements Earthbound, new in Forever.
//
// Client 1222988: Earthbind Totem also roots nearby targets (1238289) when cast. The sim casts no
// Earthbind Totem and a boss has nothing to root, so it changes no number.
func (shaman *Shaman) applyEarthbound() {
	if !shaman.Talents.Earthbound {
		return
	}
}

// applyLavaBurst grants Lava Burst, new in Forever. See lava_burst.go.
func (shaman *Shaman) applyLavaBurst() {
	shaman.registerLavaBurstSpell()
}
