package shaman

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/buffs"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

// A totem's buff is a spell of its own; the value it gives lives on that spell, not on the totem.
var windfuryTotemRank = spellData.WindfuryTotem.Highest()

// Build 70009 renamed the totem's party aura (10612, which triggers 10610) to "Windfury Totem", so the
// Triggered ladder now carries both; the attack power buff is the one the aura triggers.
var windfuryTotemBuff = spellData.WindfuryTotemTriggered.ByID(10610)
var windfuryTotemPartyAura = spellData.WindfuryTotemTriggered.ByID(10612)
var strengthOfEarthTotemRank = spellData.StrengthOfEarthTotem.Highest()
var strengthOfEarthTotemBuff = spellData.StrengthOfEarthTotemTriggered.Highest()
var graceOfAirTotemRank = spellData.GraceOfAirTotem.Highest()
var graceOfAirTotemBuff = spellData.GraceOfAirTotemTriggered.Highest()
var manaSpringTotemRank = spellData.ManaSpringTotem.Highest()
var manaSpringTotemBuff = spellData.ManaSpringTotemTriggered.Highest()

func (shaman *Shaman) newTotemSpellConfig(flatCost int32, spellID int32, spellMask int64, gcd time.Duration) core.SpellConfig {
	return core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellID},
		DefenseType:    core.DefenseTypeMagic,
		Flags:          core.SpellFlagAPL | SpellFlagInstant,
		ClassSpellMask: spellMask,

		ManaCost: core.ManaCostOptions{
			FlatCost: flatCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: gcd,
			},
		},
	}
}

func (shaman *Shaman) registerWindfuryTotemSpell() {
	duration := windfuryTotemRank.Duration()
	// Forever drops Improved Weapon Totems, so the buff's own attack power is the whole value.
	value := windfuryTotemBuff.EffectN(1).Average(core.CharacterLevel)

	wfProcAura := shaman.NewTemporaryStatsAura("Windfury Totem Proc (Self)", core.ActionID{SpellID: windfuryTotemBuff.ID}, stats.Stats{stats.AttackPower: value}, windfuryTotemBuff.Duration())
	wfProcAura.MaxStacks = int32(windfuryTotemBuff.ProcCharges)
	// The buff row's proc flags say what spends a charge, as for the party's (buffs.driveWindfuryTotem):
	// a melee auto that lands, so a missed, dodged or parried swing keeps it.
	spender := spelldata.ProcTrigger(&shaman.Character, windfuryTotemBuff, func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
		wfProcAura.RemoveStack(sim)
	}, spelldata.Chance(1))
	spender.Name = "Windfury Attack (Self)"
	spender.TriggerImmediately = true
	wfProcAura.AttachProcTrigger(spender)

	config := shaman.newTotemSpellConfig(int32(windfuryTotemRank.Cost()), windfuryTotemRank.ID, SpellMaskBasicTotem, windfuryTotemRank.GCD())

	// The party aura's own row (see buffs.driveWindfuryTotem): 20% on a melee auto or special, a
	// 100 ms internal cooldown, and the extra attack is always a main-hand one. The tooltip's main hand
	// and the yellow-hit question are buffs' switches.
	var windfurySpell *core.Spell
	trigger := spelldata.ProcTrigger(&shaman.Character, windfuryTotemPartyAura, func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		wfProcAura.Activate(sim)
		// An auto that procs it has spent the first charge itself; a special leaves both.
		charges := wfProcAura.MaxStacks
		if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
			charges--
		}
		wfProcAura.SetStacks(sim, charges)
		shaman.AutoAttacks.MaybeReplaceMHSwing(sim, windfurySpell).Cast(sim, result.Target)
	})
	buffs.NarrowWindfuryProc(&shaman.Character, &trigger, buffs.WindfuryTotemMainHandOnly)
	trigger.Name = "Windfury Totem Trigger (Self)"
	trigger.ActionID = core.ActionID{}
	trigger.MetricsActionID = core.ActionID{SpellID: windfuryTotemRank.ID}
	trigger.Duration = core.NeverExpires
	trigger.TriggerImmediately = true
	wfProcTrigger := shaman.MakeProcTriggerAura(trigger)

	wfIntermediateAuraForExclusitivity := shaman.RegisterAura(core.Aura{
		Label:    "Windfury Dummy Aura (self)",
		Duration: time.Second * 10,
	})

	wfPartyWeaponBuffTrackingAura := shaman.RegisterAura(core.Aura{
		Label:    "Windfury Party Weapon Buff Tracking Aura",
		Duration: time.Second * 10,
		ActionID: core.ActionID{SpellID: windfuryTotemRank.ID, Tag: 1},
	})

	wfAura := shaman.RegisterAura(core.Aura{
		Label:    "Windfury Totem (Self)",
		ActionID: config.ActionID,
		Duration: duration,
	})
	wfAura.NewExclusiveEffect(buffs.AirTotemCategory, true, core.ExclusiveEffect{Priority: buffs.AirTotemCastWindfury})
	buffs.WindfuryTotemDisablesFlametongueTotem(wfAura)
	wfAura.ApplyOnInit(func(aura *core.Aura, sim *core.Simulation) {
		mhConfig := *shaman.AutoAttacks.MHConfig()
		mhConfig.ActionID = mhConfig.ActionID.WithTag(windfuryTotemBuff.ID)
		windfurySpell = shaman.GetOrRegisterSpell(mhConfig)
	}).ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
		// Since build 70009 the totem's effect is a party aura that no other air totem stacks with, so it
		// ends with the totem instead of lingering on the weapon for a twist.
		wfIntermediateAuraForExclusitivity.Deactivate(sim)
		wfPartyWeaponBuffTrackingAura.Deactivate(sim)
	}).AttachPeriodicAction(core.PeriodicActionOptions{
		Period:          time.Second * 5,
		TickImmediately: true,
		Priority:        core.ActionPriorityAuto,
		OnAction: func(sim *core.Simulation) {
			wfPartyWeaponBuffTrackingAura.Activate(sim)
			wfIntermediateAuraForExclusitivity.Activate(sim)
		},
	})

	wfIntermediateAuraForExclusitivity.NewExclusiveEffect(buffs.WindfuryTotemCategory, false, core.ExclusiveEffect{
		Priority: value,
		OnGain: func(_ *core.ExclusiveEffect, sim *core.Simulation) {
			wfProcTrigger.Activate(sim)
		},
		OnExpire: func(_ *core.ExclusiveEffect, sim *core.Simulation) {
			wfProcTrigger.Deactivate(sim)
			wfIntermediateAuraForExclusitivity.Deactivate(sim)
		},
	})

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		if shaman.AirTotemAura != nil {
			shaman.AirTotemAura.Deactivate(sim)
		}
		shaman.TotemExpirations[AirTotem] = sim.CurrentTime + duration
		shaman.AirTotemAura = wfAura
		wfAura.Activate(sim)
	}

	shaman.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand}, func(sim *core.Simulation, slot proto.ItemSlot) {
		wfIntermediateAuraForExclusitivity.Deactivate(sim)
	})

	shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerStrengthOfEarthTotemSpell() {
	duration := strengthOfEarthTotemRank.Duration()
	// Forever drops Enhancing Totems, so the buff's own value is the whole value.
	value := strengthOfEarthTotemBuff.EffectN(1).Average(core.CharacterLevel)
	config := shaman.newTotemSpellConfig(int32(strengthOfEarthTotemRank.Cost()), strengthOfEarthTotemRank.ID, SpellMaskBasicTotem, strengthOfEarthTotemRank.GCD())
	buffAura := shaman.RegisterAura(core.Aura{
		Label:    "Strength Of Earth Totem (Self)",
		ActionID: config.ActionID,
		Duration: duration,
	})
	buffAura.NewExclusiveEffect(buffs.StrengthOfEarthTotemCategory+stats.Strength.StatName()+"Add", false, core.ExclusiveEffect{
		Priority: value,
		OnGain: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			ee.Aura.Unit.AddStatDynamic(sim, stats.Strength, value)
		},
		OnExpire: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			ee.Aura.Unit.AddStatDynamic(sim, stats.Strength, -value)
		},
	})
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		if shaman.EarthTotemAura != nil {
			shaman.EarthTotemAura.Deactivate(sim)
		}
		shaman.TotemExpirations[EarthTotem] = sim.CurrentTime + duration
		shaman.EarthTotemAura = buffAura
		buffAura.Activate(sim)
	}
	shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerGraceOfAirTotemSpell() {
	duration := graceOfAirTotemRank.Duration()
	value := graceOfAirTotemBuff.EffectN(1).Average(core.CharacterLevel)
	config := shaman.newTotemSpellConfig(int32(graceOfAirTotemRank.Cost()), graceOfAirTotemRank.ID, SpellMaskBasicTotem, graceOfAirTotemRank.GCD())
	buffAura := shaman.RegisterAura(core.Aura{
		Label:    "Grace Of Air Totem (Self)",
		ActionID: config.ActionID,
		Duration: duration,
	})
	buffAura.NewExclusiveEffect(buffs.AirTotemCategory, true, core.ExclusiveEffect{Priority: buffs.AirTotemCastGraceOfAir})
	buffAura.NewExclusiveEffect(buffs.GraceOfAirTotemCategory+stats.Agility.StatName()+"Add", false, core.ExclusiveEffect{
		Priority: value,
		OnGain: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			ee.Aura.Unit.AddStatDynamic(sim, stats.Agility, value)
		},
		OnExpire: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			ee.Aura.Unit.AddStatDynamic(sim, stats.Agility, -value)
		},
	})
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		if shaman.AirTotemAura != nil {
			shaman.AirTotemAura.Deactivate(sim)
		}
		shaman.TotemExpirations[AirTotem] = sim.CurrentTime + duration
		shaman.AirTotemAura = buffAura
		buffAura.Activate(sim)
	}
	shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerManaSpringTotemSpell() {
	duration := manaSpringTotemRank.Duration()
	// The buff ticks its value every 2 sec, and MP5 is the form the sim takes; Restorative Totems
	// raises it by the ladder the client states.
	tick := manaSpringTotemBuff.Effect(dbcenums.A_PERIODIC_ENERGIZE, 0)
	value := tick.Average(core.CharacterLevel) * (5 / tick.Period().Seconds()) *
		spellData.RestorativeTotems.EffectAt(1).MultiplierAt(shaman.Talents.RestorativeTotems)
	config := shaman.newTotemSpellConfig(int32(manaSpringTotemRank.Cost()), manaSpringTotemRank.ID, SpellMaskBasicTotem, manaSpringTotemRank.GCD())
	buffAura := shaman.RegisterAura(core.Aura{
		Label:    "Mana Spring Totem (Self)",
		ActionID: config.ActionID,
		Duration: duration,
	})
	buffAura.NewExclusiveEffect(buffs.ManaSpringTotemCategory+stats.MP5.StatName()+"Add", false, core.ExclusiveEffect{
		Priority: value,
		OnGain: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			ee.Aura.Unit.AddStatDynamic(sim, stats.MP5, value)
		},
		OnExpire: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			ee.Aura.Unit.AddStatDynamic(sim, stats.MP5, -value)
		},
	})
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		if shaman.WaterTotemAura != nil {
			shaman.WaterTotemAura.Deactivate(sim)
		}
		shaman.TotemExpirations[WaterTotem] = sim.CurrentTime + duration
		shaman.WaterTotemAura = buffAura
		buffAura.Activate(sim)
	}
	shaman.RegisterSpell(config)
}

// Commented out upstream of this fork by "only 1 totem per type" (53970b1d84), not by
// the Forever stubbing pass, and never called since. It is kept rather than reduced to
// the usual empty-body no-op because the body is the only record of how the totem was
// modelled; shaman.HealingStreamTotem stays nil until someone revisits totem slots.
/* func (shaman *Shaman) registerHealingStreamTotemSpell() {
	config := shaman.newTotemSpellConfig(3, 5394, SpellMaskBasicTotem, time.Second)
	hsHeal := shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 5394},
		SpellSchool:      core.SpellSchoolNature,
		ProcMask:         core.ProcMaskEmpty,
		Flags:            core.SpellFlagHelpful | core.SpellFlagNoOnCastComplete | SpellFlagInstant,
		DamageMultiplier: 1,
		CritMultiplier:   1,
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			healing := 28 + spell.HealingPower(target)*0.08272
			spell.CalcAndDealHealing(sim, target, healing, spell.OutcomeHealing)
		},
	})
	config.Hot = core.DotConfig{
		Aura: core.Aura{
			Label: "HealingStreamHot",
		},
		NumberOfTicks: 150,
		TickLength:    time.Second * 2,
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			hsHeal.Cast(sim, target)
		},
	}
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
		shaman.TotemExpirations[WaterTotem] = sim.CurrentTime + time.Second*300
		for _, agent := range shaman.Party.Players {
			spell.Hot(&agent.GetCharacter().Unit).Activate(sim)
		}
	}
	shaman.HealingStreamTotem = shaman.RegisterSpell(config)
} */
