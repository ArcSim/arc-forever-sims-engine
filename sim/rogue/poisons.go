package rogue

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
)

// Poisons are consumable imbues, not class spells, so gen_spelldata writes no table for them.
// The numbers below are the ones our Forever sim reads off beta client 1.60.1.69893, where each
// poison hits for about a third less than Classic's:
//
//	Instant Poison rank 6 (11340): 76-100 damage, 20% chance
//	Deadly Poison rank 5 (25347): 23 damage a tick, 5 stacks over 12 sec, 30% chance
//	Wound Poison rank 4 (13227): no damage, -135 healing, 5 stacks, 30% chance
//
// Attack power scaling is server-side, so the client rows carry no coefficient. The devs give 0.5%
// of attack power for Instant and 0.45% for Deadly (forever-bugs #100, 2026-09-23). Deadly's share
// is per tick per stack: beta log 2701 (Salty, AP ~363, rank 1 at 6 a tick, Vile Poisons ~1.2 from
// Instant Poison II) ticks ~8.8 a stack, against 9.2 per tick and 7.7 spread over the 4 ticks.
const (
	instantImbueID = 26891
	woundImbueID   = 27188
	deadlyImbueID  = 27186

	instantPoisonSpellID int32 = 11340
	deadlyPoisonSpellID  int32 = 25347
	woundPoisonSpellID   int32 = 13227

	instantPoisonBaseDamage = 76.0
	instantPoisonVariance   = 24.0
	deadlyPoisonTickDamage  = 23.0
	instantPoisonAPShare    = 0.005
	deadlyPoisonAPShare     = 0.0045
)

func (rogue *Rogue) applyPoisons() {
	rogue.applyDeadlyPoison()
	rogue.applyWoundPoison()
	rogue.applyInstantPoison()
}

// The chance a weapon hit applies a poison: the poison's own base, plus Improved Poisons, plus
// Venom while it is up. Read at proc time so Venom's share can come and go.
func (rogue *Rogue) poisonProcChance(base float64) float64 {
	return base + spellData.ImprovedPoisons.Effect(dbcenums.A_ADD_FLAT_MODIFIER, int32(dbcenums.SPELLMOD_CHANCE_OF_SUCCESS)).FractionAt(rogue.Talents.ImprovedPoisons) +
		rogue.additivePoisonBonusChance
}

func (rogue *Rogue) registerDeadlyPoisonSpell() {
	if rogue.getPoisonProcMask(deadlyImbueID) == core.ProcMaskUnknown {
		return
	}

	rogue.deadlyPoisonTick = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: deadlyPoisonSpellID, Tag: 100},
		SpellSchool:    core.SpellSchoolNature,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamageProc,
		ClassSpellMask: RogueSpellDeadlyPoison,
		Flags:          core.SpellFlagPoison | core.SpellFlagPassiveSpell | core.SpellFlagProc,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Deadly Poison",
				MaxStacks: 5,
				Duration:  time.Second * 12,
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				share := deadlyPoisonAPShare * float64(dot.GetStacks())
				dot.Snapshot(target, deadlyPoisonTickDamage*float64(dot.GetStacks())+share*dot.Spell.MeleeAttackPower(target))
				dot.SnapshotAttackPowerShare(target, share, false)
			},
			// The dot is 25349, which carries Periodic Can Crit in the client (25347 is the imbue).
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeTickMagicCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if !result.Landed() {
				return
			}

			dot := spell.Dot(target)
			if dot.IsActive() {
				dot.Refresh(sim)
				if dot.GetStacks() < dot.MaxStacks {
					dot.AddStack(sim)
				}
			} else {
				dot.Apply(sim)
				dot.SetStacks(sim, 1)
			}
			dot.TakeSnapshot(sim)
		},
	})
	rogue.DeadlyPoison = rogue.deadlyPoisonTick
}

func (rogue *Rogue) registerWoundPoisonSpell() {
	if rogue.getPoisonProcMask(woundImbueID) == core.ProcMaskUnknown {
		return
	}

	woundPoisonDebuffAura := core.Aura{
		Label:     "Wound Poison",
		ActionID:  core.ActionID{SpellID: woundPoisonSpellID},
		Duration:  time.Second * 15,
		MaxStacks: 5,
		// The healing debuff has no effect on a DPS sim.
	}

	rogue.WoundPoisonDebuffAuras = rogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.RegisterAura(woundPoisonDebuffAura)
	})

	rogue.WoundPoison = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: woundPoisonSpellID},
		SpellSchool:    core.SpellSchoolNature,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamageProc,
		ClassSpellMask: RogueSpellWoundPoison,
		Flags:          core.SpellFlagPoison | core.SpellFlagPassiveSpell | core.SpellFlagProc,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Forever's Wound Poison deals no damage; it only stacks the healing debuff.
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if !result.Landed() {
				return
			}

			aura := rogue.WoundPoisonDebuffAuras.Get(target)
			if !aura.IsActive() {
				aura.Activate(sim)
				aura.SetStacks(sim, 1)
				return
			}
			aura.Refresh(sim)
			if aura.GetStacks() < aura.MaxStacks {
				aura.AddStack(sim)
			}
		},

		RelatedAuraArrays: rogue.WoundPoisonDebuffAuras.ToMap(),
	})
}

func (rogue *Rogue) registerInstantPoisonSpell() {
	if rogue.getPoisonProcMask(instantImbueID) == core.ProcMaskUnknown {
		return
	}

	rogue.InstantPoison = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: instantPoisonSpellID},
		SpellSchool:    core.SpellSchoolNature,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamageProc,
		ClassSpellMask: RogueSpellInstantPoison,
		Flags:          core.SpellFlagPoison | core.SpellFlagPassiveSpell | core.SpellFlagProc,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := sim.Roll(instantPoisonBaseDamage, instantPoisonBaseDamage+instantPoisonVariance) + instantPoisonAPShare*spell.MeleeAttackPower(target)
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (rogue *Rogue) getPoisonProcMask(poisonId int32) core.ProcMask {
	var mask core.ProcMask
	if rogue.Consumables.MhImbueId == poisonId {
		mask |= core.ProcMaskMeleeMH
	}
	if rogue.Consumables.OhImbueId == poisonId {
		mask |= core.ProcMaskMeleeOH
	}
	return mask
}

// The three imbues share a shape: a weapon proc on the hand it is on, rolled against a chance
// Improved Poisons and Venom both add to, which is why the roll is inside the handler.
func (rogue *Rogue) applyPoisonProc(name string, imbueID int32, baseChance float64, cast func(*core.Simulation, *core.Unit)) {
	procMask := rogue.getPoisonProcMask(imbueID)
	if procMask == core.ProcMaskUnknown {
		return
	}

	rogue.MakeProcTriggerAura(core.ProcTrigger{
		Name:               name,
		Outcome:            core.OutcomeLanded,
		Callback:           core.CallbackOnSpellHitDealt,
		TriggerImmediately: true,
		ProcMask:           procMask,
		IsWeaponProc:       true,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.Proc(rogue.poisonProcChance(baseChance), name) {
				cast(sim, result.Target)
			}
		},
	})
}

func (rogue *Rogue) applyDeadlyPoison() {
	rogue.applyPoisonProc("Deadly Poison", deadlyImbueID, 0.3, func(sim *core.Simulation, target *core.Unit) {
		rogue.DeadlyPoison.Cast(sim, target)
	})
}

func (rogue *Rogue) applyWoundPoison() {
	rogue.applyPoisonProc("Wound Poison", woundImbueID, 0.3, func(sim *core.Simulation, target *core.Unit) {
		rogue.WoundPoison.Cast(sim, target)
	})
}

func (rogue *Rogue) applyInstantPoison() {
	rogue.applyPoisonProc("Instant Poison", instantImbueID, 0.2, func(sim *core.Simulation, target *core.Unit) {
		rogue.InstantPoison.Cast(sim, target)
	})
}

// Instant Poison leaves nothing behind, so only the two lingering poisons count.
func (rogue *Rogue) isPoisoned(target *core.Unit) bool {
	if rogue.deadlyPoisonTick != nil && rogue.deadlyPoisonTick.Dot(target).IsActive() {
		return true
	}
	return rogue.WoundPoisonDebuffAuras != nil && rogue.WoundPoisonDebuffAuras.Get(target).IsActive()
}
