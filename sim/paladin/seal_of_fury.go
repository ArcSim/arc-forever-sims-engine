package paladin

import (
	"fmt"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// Seal of Fury
// https://www.wowhead.com/forever/spell=20423
//
// Fills the Paladin with divine fury for 30 sec, causing melee attacks to deal an additional X
// Holy damage. While a shield is equipped, each attack also grants an absorb shield equal to 50%
// of the Holy damage dealt. Only one Seal can be active on the Paladin at any one time.
//
// Unleashing this Seal's energy causes Holy damage to an enemy and taunts the target to attack
// you for 4 sec.
//
// The per-hit damage is the spell the tooltip names first, a flat number with its row's 10% coefficient
// (SealOfFurySpellPowerFromRows); the seal's third effect names its judgement. The seal also carries a
// weapon-speed dummy in Seal of Righteousness's shape that the tooltip never references, and it
// is left alone. The taunt has no place in the sim.
// Seal of Fury's spell power shares as measured in game on client 70425 (MythicSim patch 143): 15% on
// each hit and 67% on the judgement, where the rows carry 0.10 and 0.429. The bases already matched.
const (
	SealOfFurySpellPowerCoefficient      = 0.15
	JudgementOfFurySpellPowerCoefficient = 0.67
)

// SealOfFurySpellPowerFromRows gives the seal's hit and its judgement the shares their client rows
// carry, 0.10 and 0.429 in client 70425, instead of patch 143's 0.15 and 0.67 (MythicSim patch 168).
// Lightclub's 16 October test on the live server prints "10% spellpower (11.5% improved)" for the seal
// and "42.9% of spellpower (49.3% improved)" for the judgement: the rows, and the rows times Improved
// Seals. Patch 143's numbers are both exactly 1.5 times the pre-16 October rows (0.10, 0.45), and the
// pre-16 October logs read the rows (patch 163). False keeps 0.15 and 0.67.
var SealOfFurySpellPowerFromRows = true

func sealOfFuryCoefficients(proc, judgement *spelldata.Effect) (float64, float64) {
	if SealOfFurySpellPowerFromRows {
		return proc.Coeff(), judgement.Coeff()
	}
	return SealOfFurySpellPowerCoefficient, JudgementOfFurySpellPowerCoefficient
}

func (paladin *Paladin) registerSealOfFury(_ int32, rank *spelldata.Spell) {
	judgementRank := spellData.SealOfFuryTriggered.ByID(int32(rank.EffectN(3).BaseValue()))
	judgementDamage := judgementRank.DamageEffect()
	procRank := rank.Refs()[0]
	procDamage := procRank.DamageEffect()
	procCoefficient, judgementCoefficient := sealOfFuryCoefficients(procDamage, judgementDamage)

	// Melee in SpellCategories with No Active Defense: hit and crit on the melee table, never
	// dodged, parried or blocked.
	judgement := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: judgementRank.ID},
		SpellSchool:    judgementRank.SpellSchool(),
		DefenseType:    judgementRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagBinary,
		ClassSpellMask: SpellMaskJudgementOfFury,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: judgementCoefficient,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, judgementDamage.Roll(sim, core.CharacterLevel), spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})

	shieldShare := rank.EffectN(2).Percent()
	var pendingShield float64
	shield := paladin.NewDamageAbsorptionAura(core.AbsorptionAuraConfig{
		Aura: core.Aura{
			Label:    fmt.Sprintf("Seal of Fury Shield%s Rank %d", paladin.Label, rank.RankNumber()),
			ActionID: core.ActionID{SpellID: rank.ID}.WithTag(1),
			Duration: sealDuration,
		},
		ShieldStrengthCalculator: func(_ *core.Unit) float64 {
			return pendingShield
		},
	})
	paladin.applyImprovedSealOfFury(shield)

	damage := procDamage.Average(core.CharacterLevel)
	procSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: procRank.ID},
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		// The damage spells carry the same flags as Seal of Righteousness's: procs themselves,
		// and no Suppress Weapon Procs, so a weapon's "Chance on hit" rolls on the seal's hit too.
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell | core.SpellFlagProc,
		ClassSpellMask: SpellMaskSealOfFuryProc,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: procCoefficient,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, damage, spell.OutcomeMeleeSpecialCritOnly)
			sim.AddPendingAction(core.NewDelayedAction(core.DelayedActionOptions{
				DoAt:     sim.CurrentTime + core.SpellBatchWindow,
				Priority: core.ActionPriorityLow,
				OnAction: func(sim *core.Simulation) {
					spell.DealDamage(sim, result)
					if paladin.PseudoStats.CanBlock && result.Damage > 0 {
						pendingShield = result.Damage * shieldShare
						shield.Activate(sim)
					}
				},
			}))
		},
	})

	aura := paladin.makeSealExclusive(paladin.MakeProcTriggerAura(core.ProcTrigger{
		Name:            sealLabel("Seal of Fury", paladin, rank),
		ActionID:        core.ActionID{SpellID: rank.ID},
		MetricsActionID: core.ActionID{SpellID: rank.ID},
		Duration:        sealDuration,
		Callback:        core.CallbackOnSpellHitDealt,
		ProcMask:        core.ProcMaskMeleeWhiteHit,
		Outcome:         core.OutcomeLanded,
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			procSpell.Cast(sim, result.Target)
		},
	}))

	paladin.registerSealSpell(&sealConfig{
		rank:      rank,
		classMask: SpellMaskSealOfFury,
		aura:      aura,
		judgement: judgement,
		echoID:    echoOfFuryID,
		echo: func(sim *core.Simulation, target *core.Unit) {
			procSpell.Cast(sim, target)
		},
	})
}
