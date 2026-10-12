package shaman

import (
	"fmt"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/buffs"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

const (
	frostbrandEnchantID  int32 = 2
	flametongueEnchantID int32 = 5
	windfuryEnchantID    int32 = 283
	earthlivingEnchantID int32 = 3345
	rockbiterEnchantID   int32 = 3021
)

func (shaman *Shaman) RegisterOnItemSwapWithImbue(effectID int32, procMask *core.ProcMask, aura *core.Aura) {
	shaman.RegisterItemSwapCallback(core.AllWeaponSlots(), func(sim *core.Simulation, slot proto.ItemSlot) {
		mask := core.ProcMaskUnknown
		if shaman.MainHand().TempEnchant == effectID {
			mask |= core.ProcMaskMeleeMH
		}
		if shaman.OffHand().TempEnchant == effectID {
			mask |= core.ProcMaskMeleeOH
		}
		*procMask = mask

		if mask == core.ProcMaskUnknown {
			aura.Deactivate(sim)
		} else {
			aura.Activate(sim)
		}
	})
}

func (shaman *Shaman) setupItemSwapImbue(imbue proto.ShamanImbue, imbueID int32) {
	if shaman.ItemSwap.IsEnabled() {
		if mhSwap := shaman.ItemSwap.GetUnequippedItemBySlot(proto.ItemSlot_ItemSlotMainHand); mhSwap != nil && shaman.SelfBuffs.ImbueMHSwap == imbue {
			mhSwap.TempEnchant = imbueID
			shaman.ItemSwap.AddTempEnchant(imbueID, proto.ItemSlot_ItemSlotMainHand, true)
		}
		if ohSwap := shaman.ItemSwap.GetUnequippedItemBySlot(proto.ItemSlot_ItemSlotOffHand); ohSwap != nil && shaman.SelfBuffs.ImbueOHSwap == imbue {
			ohSwap.TempEnchant = imbueID
			shaman.ItemSwap.AddTempEnchant(imbueID, proto.ItemSlot_ItemSlotOffHand, true)
		}
	}
}

// newWindfuryAttackSpell is the Windfury attack: 439440 (main hand) or 439441 (Attributes[3] 0x1000000,
// off hand), a SPELL_EFFECT_WEAPON_DAMAGE special hit carrying the rank's extra attack power (16361:
// 333 at 60). Beta log 2708 (foreverlogs.gg, Toma, 22 Enhancement, level 30) bears this out: 62 procs,
// each two 439440 hits in the same instant that miss, dodge and parry on their own; the next auto
// still lands one weapon speed after the last (0.2 sec after a Stormstrike proc), so the swing timer
// is untouched; and no 16361 buff is ever applied, so there is no third, lifted swing.
func (shaman *Shaman) newWindfuryAttackSpell(isMH bool) *core.Spell {
	apBonus := shaman.WindfuryAPBonus * (1 + spellData.ElementalWeapons.EffectAt(3).FractionAt(shaman.Talents.ElementalWeapons))
	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: core.TernaryInt32(isMH, 439440, 439441)},
		SpellSchool:      core.SpellSchoolPhysical,
		DefenseType:      core.DefenseTypeMelee,
		ProcMask:         core.Ternary(isMH, core.ProcMaskMeleeMHSpecial, core.ProcMaskMeleeOHSpecial),
		Flags:            core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			weaponDamage := core.Ternary(isMH, spell.Unit.MHWeaponDamage, spell.Unit.OHWeaponDamage)
			baseDamage := weaponDamage(sim, spell.MeleeAttackPower(target)+apBonus)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	})
}

// The enchant's equip aura (439431 in every rank's SpellItemEnchantment) rolls 20% on landed melee
// autos and abilities (ProcTypeMask 20) with a 1.5 sec ProcCategoryRecovery and strikes twice with
// the proc's weapon. It is an enchant on the weapon, so only that weapon's hits roll it; Forever
// shamans cannot dual wield, so in practice that is the main hand. buffs.WindfuryProcsFromSpecials
// set false leaves only the weapon's auto attacks.
func (shaman *Shaman) makeWFProcTriggerAura(dpm *core.DynamicProcManager, procMask *core.ProcMask) *core.Aura {
	mhAttack := shaman.newWindfuryAttackSpell(true)
	ohAttack := shaman.newWindfuryAttackSpell(false)
	var excludeSpecials core.ProcMask
	if !buffs.WindfuryProcsFromSpecials {
		excludeSpecials = core.ProcMaskMeleeSpecial
	}
	aura := shaman.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Windfury Imbue",
		Callback:           core.CallbackOnSpellHitDealt,
		ProcMask:           *procMask,
		ProcMaskExclude:    excludeSpecials,
		IsWeaponProc:       true,
		Outcome:            core.OutcomeLanded,
		ICD:                time.Millisecond * 1500,
		DPM:                dpm,
		TriggerImmediately: true,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			attack := core.Ternary(spell.IsMH(), mhAttack, ohAttack)
			attack.Cast(sim, result.Target)
			attack.Cast(sim, result.Target)
		},
	})
	return aura
}

func (shaman *Shaman) getWindfuryFixedProcChance(procMask core.ProcMask) float64 {
	return 0.2
}

func (shaman *Shaman) RegisterWindfuryImbue(procMask core.ProcMask) {
	if procMask == core.ProcMaskUnknown && !shaman.ItemSwap.IsEnabled() {
		return
	}

	mask := core.ProcMaskUnknown

	mH := shaman.MainHand()
	if mH != nil && shaman.SelfBuffs.ImbueMH == proto.ShamanImbue_WindfuryWeapon {
		mH.TempEnchant = windfuryEnchantID
		if shaman.ItemSwap.IsEnabled() {
			shaman.ItemSwap.AddTempEnchant(windfuryEnchantID, proto.ItemSlot_ItemSlotMainHand, false)
		}
		mask |= core.ProcMaskMeleeMH
	}
	oH := shaman.OffHand()
	if oH != nil && shaman.SelfBuffs.ImbueOH == proto.ShamanImbue_WindfuryWeapon {
		oH.TempEnchant = windfuryEnchantID
		if shaman.ItemSwap.IsEnabled() {
			shaman.ItemSwap.AddTempEnchant(windfuryEnchantID, proto.ItemSlot_ItemSlotOffHand, false)
		}
		mask |= core.ProcMaskMeleeOH
	}

	shaman.setupItemSwapImbue(proto.ShamanImbue_WindfuryWeapon, windfuryEnchantID)

	dpm := shaman.NewDynamicLegacyProcForTempEnchant(windfuryEnchantID, 0, shaman.getWindfuryFixedProcChance)

	aura := shaman.makeWFProcTriggerAura(dpm, &mask)

	if mask.Matches(core.ProcMaskMeleeMH) {
		aura.NewExclusiveEffect(buffs.WindfuryTotemCategory, false, core.ExclusiveEffect{
			Priority: shaman.WindfuryAPBonus * 2, // Need to be higher than Windfury Totem priority
		})
	}

	shaman.RegisterOnItemSwapWithImbue(windfuryEnchantID, &mask, aura)
}

var windfuryImbue = spellData.WindfuryWeaponTriggered.Highest()

// Rank 6's proc dummy. Since the 70170 hotfixes the triggered ladder also lists the Flametongue Attack
// damage spells (10444, 29469, 29470), so Highest() would be the attack rather than the dummy.
var flametongueImbue = spellData.FlametongueWeaponTriggered.ByID(16344)
var frostbrandImbue = spellData.FrostbrandWeaponTriggered.Highest()
var rockbiterImbue = spellData.RockbiterWeaponTriggered.Highest()

// A Flametongue Totem hit lands as Flametongue Attack 16368 (beta log 2713, upstream #682), whose client row
// has no spell power coefficient and a class mask (bit 25) that Elemental Fury and Elemental Weapons do not
// name, so a shaman's carries no talent mask. It keeps the shaman spell flag, so Natural Grace's threat cut
// reaches it.
func init() {
	buffs.SetFlametongueAttackTraits(proto.Class_ClassShaman, buffs.FlametongueAttackTraits{
		Flags: SpellFlagShamanSpell,
	})
}

func (shaman *Shaman) newFlametongueImbueSpell(weapon *core.Item) *core.Spell {
	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: flametongueImbue.ID},
		SpellSchool: core.SpellSchoolFire,
		// The damage logs as Flametongue Attack (10444), Magic in SpellCategories; it crits for 1.5x
		// (2.0x with Elemental Fury, see talents_elemental.go).
		DefenseType:      core.DefenseTypeMagic,
		ProcMask:         core.ProcMaskSpellDamageProc,
		ClassSpellMask:   SpellMaskFlametongueWeapon,
		Flags:            core.SpellFlagPassiveSpell | core.SpellFlagProc | SpellFlagShamanSpell,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 0.10000000149,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if weapon.SwingSpeed != 0 {
				// The proc's dummy value is hundredths of damage per second of weapon speed, speed held
				// to 1.3-4.0 (the tooltip's "(X / 77 - 1) to (X / 25)"): 16344's 2810 is 35 to 112 at 60.
				// A beta log bears the scale out at rank 1 (8026, 4.4 a second): 100 hits from a 2.5
				// speed weapon with no spell power averaged 11.2 (foreverlogs.gg report 2671).
				speed := min(max(weapon.SwingSpeed, 1.3), 4)
				baseDamage := speed * flametongueImbue.EffectN(1).Average(core.CharacterLevel) / 100
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}

func (shaman *Shaman) makeFTProcTriggerAura(itemSlot proto.ItemSlot, triggerProcMask core.ProcMask, flameTongueSpell *core.Spell) *core.Aura {
	aura := shaman.MakeProcTriggerAura(core.ProcTrigger{
		Name:               fmt.Sprintf("Flametongue Imbue %s", itemSlot),
		ProcMask:           triggerProcMask,
		IsWeaponProc:       true,
		Outcome:            core.OutcomeLanded,
		Callback:           core.CallbackOnSpellHitDealt,
		TriggerImmediately: true,

		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			flameTongueSpell.Cast(sim, result.Target)
		},
	})

	shaman.RegisterItemSwapCallback([]proto.ItemSlot{itemSlot}, func(sim *core.Simulation, is proto.ItemSlot) {
		if is == proto.ItemSlot_ItemSlotMainHand {
			mh := shaman.MainHand()
			mhSwap := shaman.ItemSwap.GetUnequippedItemBySlot(is)
			if mh.TempEnchant != flametongueEnchantID {
				// The new main hand does not have flametongue on, so deactivate
				aura.Deactivate(sim)
				return
			}
			if mhSwap.TempEnchant != flametongueEnchantID {
				// The new main hand has flametongue on and the swapped one does not, so need to activate
				aura.Activate(sim)
				return
			}
		}
		if is == proto.ItemSlot_ItemSlotOffHand {
			oh := shaman.OffHand()
			ohSwap := shaman.ItemSwap.GetUnequippedItemBySlot(is)
			if oh.TempEnchant != flametongueEnchantID {
				// The new offhand does not have flametongue on, so deactivate
				aura.Deactivate(sim)
				return
			}
			if ohSwap.TempEnchant != flametongueEnchantID {
				// The new offhand has flametongue on and the swapped one does not, so need to activate
				aura.Activate(sim)
				return
			}

		}
	})

	return aura
}

func (shaman *Shaman) RegisterFlametongueImbue(procMask core.ProcMask) {
	if procMask == core.ProcMaskUnknown && !shaman.ItemSwap.IsEnabled() {
		return
	}

	for _, itemSlot := range core.AllWeaponSlots() {
		var weapon *core.Item
		var triggerProcMask core.ProcMask
		switch {
		case shaman.SelfBuffs.ImbueMH == proto.ShamanImbue_FlametongueWeapon && itemSlot == proto.ItemSlot_ItemSlotMainHand:
			weapon = shaman.MainHand()
			triggerProcMask = core.ProcMaskMeleeMH
		case shaman.SelfBuffs.ImbueOH == proto.ShamanImbue_FlametongueWeapon && itemSlot == proto.ItemSlot_ItemSlotOffHand:
			weapon = shaman.OffHand()
			triggerProcMask = core.ProcMaskMeleeOH
		}

		if weapon == nil {
			continue
		}

		weapon.TempEnchant = flametongueEnchantID

		if shaman.ItemSwap.IsEnabled() {
			shaman.ItemSwap.AddTempEnchant(flametongueEnchantID, itemSlot, false)
		}

		flameTongueSpell := shaman.newFlametongueImbueSpell(weapon)
		// Classic's Windfury Totem enchanted the weapon, so a main-hand imbue displaced it. Forever's is a
		// party aura, and only Windfury Weapon names Windfury Totem (upstream #674).
		aura := shaman.makeFTProcTriggerAura(itemSlot, triggerProcMask, flameTongueSpell)
		// It is in Flametongue Totem's: "When applied to main hand, disables any benefit you personally
		// receive from Flametongue Totem" (patch 70). An off-hand Flametongue Weapon leaves the totem on.
		if itemSlot == proto.ItemSlot_ItemSlotMainHand {
			buffs.DisableFlametongueTotem(aura)
		}
	}

	shaman.setupItemSwapImbue(proto.ShamanImbue_FlametongueWeapon, flametongueEnchantID)
}

func (shaman *Shaman) newFrostbrandImbueSpell() *core.Spell {
	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: frostbrandImbue.ID},
		SpellSchool:    frostbrandImbue.SpellSchool(),
		DefenseType:    core.DefenseTypeMagic, // Frostbrand Attack (25501 / 38617) is Magic in SpellCategories
		ClassSpellMask: SpellMaskFrostbrandWeapon,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagPassiveSpell | core.SpellFlagProc | SpellFlagShamanSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: frostbrandImbue.DamageEffect().Coeff(),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, frostbrandImbue.DamageEffect().Average(core.CharacterLevel), spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (shaman *Shaman) RegisterFrostbrandImbue(procMask core.ProcMask) {
	if procMask == core.ProcMaskUnknown && !shaman.ItemSwap.IsEnabled() {
		return
	}

	mH := shaman.MainHand()
	if mH != nil && shaman.SelfBuffs.ImbueMH == proto.ShamanImbue_FrostbrandWeapon {
		mH.TempEnchant = frostbrandEnchantID
		if shaman.ItemSwap.IsEnabled() {
			shaman.ItemSwap.AddTempEnchant(frostbrandEnchantID, proto.ItemSlot_ItemSlotMainHand, false)
		}
	}
	oH := shaman.OffHand()
	if oH != nil && shaman.SelfBuffs.ImbueOH == proto.ShamanImbue_FrostbrandWeapon {
		oH.TempEnchant = frostbrandEnchantID
		if shaman.ItemSwap.IsEnabled() {
			shaman.ItemSwap.AddTempEnchant(frostbrandEnchantID, proto.ItemSlot_ItemSlotOffHand, false)
		}
	}

	shaman.setupItemSwapImbue(proto.ShamanImbue_FrostbrandWeapon, frostbrandEnchantID)

	// 8 procs a minute, not Classic's 9 (the client stores no rate): a beta log gives 636 procs from 1,601
	// landed swings of a 3.0 speed axe, 0.397 a hit = 7.95 +- 0.24 PPM (foreverlogs.gg report 2681, Cihan).
	dpm := shaman.NewDynamicLegacyProcForTempEnchant(frostbrandEnchantID, 8.0, func(pm core.ProcMask) float64 { return 0 })

	fbSpell := shaman.newFrostbrandImbueSpell()

	aura := shaman.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Frostbrand Imbue",
		Callback:           core.CallbackOnSpellHitDealt,
		IsWeaponProc:       true,
		Outcome:            core.OutcomeLanded,
		DPM:                dpm,
		TriggerImmediately: true,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			fbSpell.Cast(sim, result.Target)
		},
	})

	shaman.RegisterOnItemSwapWithImbue(frostbrandEnchantID, &procMask, aura)
}

// Rockbiter Weapon's passive (16313, rank 7) is a standing melee attack power aura, 554 + 16.5 a level
// = 653 at 60 (Wowhead env 16 prints 686, its MaxLevel 62), not Classic's per-hit damage. Elemental
// Weapons' first effect adds 7/13/20% to it (SPELLMOD_EFFECT1 on the passive's mask). Spirit Weapons'
// second effect (+86% to the passive's threat aura) turns the talent's -30% threat into +30% (0.7 x
// 1.86) while Rockbiter is up, as its tooltip says. One passive whatever the hands, so a second
// Rockbiter weapon adds nothing. Ported from MythicSim patch 19 (sage3648).
func (shaman *Shaman) RegisterRockbiterImbue(procMask core.ProcMask) {
	if procMask == core.ProcMaskUnknown && !shaman.ItemSwap.IsEnabled() {
		return
	}

	imbued := false
	mH := shaman.MainHand()
	if mH != nil && shaman.SelfBuffs.ImbueMH == proto.ShamanImbue_RockbiterWeapon {
		mH.TempEnchant = rockbiterEnchantID
		if shaman.ItemSwap.IsEnabled() {
			shaman.ItemSwap.AddTempEnchant(rockbiterEnchantID, proto.ItemSlot_ItemSlotMainHand, false)
		}
		imbued = true
	}
	oH := shaman.OffHand()
	if oH != nil && shaman.SelfBuffs.ImbueOH == proto.ShamanImbue_RockbiterWeapon {
		oH.TempEnchant = rockbiterEnchantID
		if shaman.ItemSwap.IsEnabled() {
			shaman.ItemSwap.AddTempEnchant(rockbiterEnchantID, proto.ItemSlot_ItemSlotOffHand, false)
		}
		imbued = true
	}

	shaman.setupItemSwapImbue(proto.ShamanImbue_RockbiterWeapon, rockbiterEnchantID)

	if !imbued {
		return
	}

	attackPower := rockbiterImbue.Effect(dbcenums.A_MOD_ATTACK_POWER, 0).Average(core.CharacterLevel) *
		(1 + spellData.ElementalWeapons.EffectAt(1).FractionAt(shaman.Talents.ElementalWeapons))
	core.MakePermanent(shaman.NewTemporaryStatsAura("Rockbiter Weapon", core.ActionID{SpellID: rockbiterImbue.ID},
		stats.Stats{stats.AttackPower: attackPower}, core.NeverExpires).Aura)

	if shaman.Talents.SpiritWeapons {
		shaman.PseudoStats.ThreatMultiplier *= spellData.SpiritWeapons.EffectAt(2).MultiplierAt(1)
	}
}
