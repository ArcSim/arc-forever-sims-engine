package paladin

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

// SealOfCommandBonusTakesImprovedSeals lets Improved Seals raise the part of a Seal of Command proc
// that the target's extra Holy damage taken adds (Judgement of the Crusader, damage against the
// target's creature type), as it raises the rest of the proc (MythicSim patch 169). Lightclub's
// live-server sheet (16 October 2026): Seal of Command takes the "full coeff from JotC", and Improved
// Seals "increases all damage caused (by Seal and Judgements)". False adds that part at the bare 0.29,
// outside the talent, as before.
var SealOfCommandBonusTakesImprovedSeals = true

// The share of the target's extra Holy damage taken a Seal of Command proc adds.
func sealOfCommandTargetBonusCoefficient(coefficient, improvedSeals float64) float64 {
	if SealOfCommandBonusTakesImprovedSeals {
		return coefficient * improvedSeals
	}
	return coefficient
}

// Seal of Command (talent)
// https://www.wowhead.com/forever/spell=20920
//
// Gives the Paladin a chance to deal additional Holy damage equal to 70% of normal weapon damage.
// Only one Seal can be active on the Paladin at any one time. Lasts 30 sec.
//
// Unleashing this Seal's energy will judge an enemy, instantly causing Holy damage, double if the
// target is stunned or incapacitated.
//
// The client states the stunned number and halves it otherwise. The judgement's damage sits on the
// spell its own dummy names, which the family's triggered ladder carries one rank above the proc.
// The proc spell 20424 has no rank subtext and heads that ladder: its 70% of weapon damage and the
// 0.29 coefficient on that same effect are read from it. The client marks the seal as procs per
// minute but states no rate; 7 is the number the Classic sim carries, and the beta's public combat
// logs agree (532 procs off white hits of four level 20 paladins, 6.75 +- 0.23 a minute).
//
// The coefficient sits on the weapon-percent effect, and the percent applies to the paladin's own
// spell power as well as the weapon: at 27 spell power the Forever beta measured 5.5 of each proc
// coming from spell power, 20.4%, against the 29% the coefficient states. Improved Seals raises the
// percent itself (the tooltip reads 80% with the talent, 70 * 1.15). The Holy damage the target
// takes extra, Judgement of the Crusader and the damage-against-mob-type gear, is added after the
// percent at the full coefficient. Improved Seals multiplies that part as well
// (SealOfCommandBonusTakesImprovedSeals).
//
// The 1 second internal cooldown is the engine's own; no live test has confirmed it.
func (paladin *Paladin) registerSealOfCommand(_ int32, rank *spelldata.Spell) {
	judgeRank := spellData.JudgementOfCommand.Rank(rank.RankNumber())
	judgeDamage := spellData.SealOfCommandTriggered.ByID(int32(judgeRank.EffectN(1).BaseValue())).DamageEffect()

	// Melee in SpellCategories with No Active Defense: hit and crit on the melee table, never
	// dodged, parried or blocked.
	judgement := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: judgeRank.ID},
		SpellSchool:    judgeRank.SpellSchool(),
		DefenseType:    judgeRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskJudgementOfCommand,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: judgeDamage.Coeff(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := judgeDamage.Roll(sim, core.CharacterLevel)
			if !target.PseudoStats.Stunned {
				baseDamage /= 2
			}
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})

	procRank := spellData.SealOfCommandTriggered.Rank(1)
	procEffect := procRank.EffectN(1)
	improvedSeals := spellData.ImprovedSeals.MultiplierAt(paladin.Talents.ImprovedSeals)
	weaponPercent := procEffect.Percent() * improvedSeals
	coefficient := procEffect.Coeff()
	bonusCoefficient := sealOfCommandTargetBonusCoefficient(coefficient, improvedSeals)
	procSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: procRank.ID},
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		// 20424 carries Not a Proc: it is an ability hit to every listener, weapon procs included.
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,
		ClassSpellMask: SpellMaskSealOfCommandProc,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		// No BonusCoefficient: the spell power goes through the weapon percent, so the proc adds
		// it by hand below rather than letting CalcDamage add it on top.

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			weaponDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
			baseDamage := (weaponDamage + coefficient*spell.SpellDamage(target)) * weaponPercent

			targetBonus := target.PseudoStats.SchoolBonusSpellDamage[stats.SchoolIndexHoly] +
				spell.Unit.AttackTables[target.UnitIndex].MobTypeBonusStats[target.MobType][stats.SpellDamage]
			baseDamage += bonusCoefficient * targetBonus

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			dealAfterBatch(sim, spell, result)
		},
	})

	// The seal and its Echo roll the same chance: an Echo of Command "empowers your next melee
	// attack with a chance to activate Seal of Command".
	dpm := paladin.NewLegacyPPMManager(7, core.ProcMaskMeleeWhiteHit)
	icd := core.Cooldown{Timer: paladin.NewTimer(), Duration: time.Second}
	tryProc := func(sim *core.Simulation, target *core.Unit) {
		if icd.IsReady(sim) && dpm.Proc(sim, core.ProcMaskMeleeMHAuto, "Seal of Command") {
			icd.Use(sim)
			procSpell.Cast(sim, target)
		}
	}

	aura := paladin.makeSealExclusive(paladin.RegisterAura(core.Aura{
		Label:    sealLabel("Seal of Command", paladin, rank),
		ActionID: core.ActionID{SpellID: rank.ID},
		Duration: sealDuration,
	}).AttachProcTrigger(core.ProcTrigger{
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeWhiteHit,
		Outcome:  core.OutcomeLanded,
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			tryProc(sim, result.Target)
		},
	}))

	paladin.registerSealSpell(&sealConfig{
		rank:      rank,
		classMask: SpellMaskSealOfCommand,
		aura:      aura,
		judgement: judgement,
		echoID:    echoOfCommandID,
		echo:      tryProc,
	})
}
