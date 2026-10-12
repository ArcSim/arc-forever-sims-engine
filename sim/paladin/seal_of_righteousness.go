package paladin

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// The damage spell each rank fires on a hit. They carry no rank subtext, so no ladder holds them.
var sealOfRighteousnessProcIDs = map[int32]int32{1: 25742, 2: 25740, 3: 25739, 4: 25738, 5: 25737, 6: 25736, 7: 25735, 8: 25713}

// SealOfRighteousnessFollowsSwappedWeapon reads the main hand at each hit for the spell power share,
// as the base damage already does, so a swap between a two-hander and a one-hander (a Manual Crowd
// Pummeler opener) changes the coefficient with the weapon. False keeps the coefficient of the
// weapon worn at the pull for the whole fight, as it was registered before (MythicSim patch 129).
var SealOfRighteousnessFollowsSwappedWeapon = true

// Seal of Righteousness
// https://www.wowhead.com/forever/spell=20293
//
// Fills the Paladin with holy spirit for 30 sec, granting each melee attack an additional X to Y
// Holy damage. Slower weapons cause more Holy damage per swing. Only one Seal can be active on
// the Paladin at any one time.
//
// Unleashing this Seal's energy will cause Holy damage to an enemy.
//
// The seal's second effect names its judgement (the row has no effect at the client's index 1).
// The judgement's dummy, its second effect, is what the seal's tooltip renders each hit from: the
// number and the coefficient the per-hit damage scales on. The hit itself is a function of weapon
// speed: the value per hundred, times the swing speed, times 0.85 for a one-hander or 1.2 for a
// two-hander, the formula the Classic sim settled on from testing. The beta's public combat logs
// agree at level 20 with one-handers of 1.7 and 2.4 speed and a ~3.4 two-hander, within half a
// point a hit, and they want the dummy's coefficient rather than the proc row's 0.1.
//
// A two-hander's 1.2 lands on the spell power share as well: 0.185 x 1.2 = 0.222. Exordium (3.3 s
// two-hander, no Holy points; foreverlogs 2693) hit for 17, and for 24-25 while their rank 2
// Judgement of the Crusader (20188, +35 Holy taken) was up: +7.4 to +7.8 a hit, where 0.185 gives
// 6.5. Their Judgement of Righteousness rose 17.9 on the same debuff, the 0.5 it should. Jack's ~3.4
// two-hander (2670) also hit 25s on top of 16-18, the same step. No one-hander has been logged
// under the debuff, so a one-hander keeps the bare 0.185 (Ms's 1.7 one-hander fits it).

func (paladin *Paladin) registerSealOfRighteousness(_ int32, rank *spelldata.Spell) {
	judgeRank := spellData.JudgementOfRighteousness.ByID(int32(rank.EffectN(2).BaseValue()))
	judgeDamage := judgeRank.DamageEffect()
	perHit := judgeRank.EffectN(2)

	// The judgement is Melee in SpellCategories and carries No Active Defense: it rolls hit and
	// crit on the melee table and cannot be dodged, parried or blocked.
	judgement := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: judgeRank.ID},
		SpellSchool:    judgeRank.SpellSchool(),
		DefenseType:    judgeRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagBinary,
		ClassSpellMask: SpellMaskJudgementOfRighteousness,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: judgeDamage.Coeff(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, judgeDamage.Roll(sim, core.CharacterLevel), spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})

	value := perHit.Average(core.CharacterLevel)
	coefficient := perHit.Coeff()
	twoHandCoefficient := func() float64 {
		if paladin.GetMainHandType() == proto.HandType_HandTypeTwoHand {
			return coefficient * 1.2
		}
		return coefficient
	}
	procSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: sealOfRighteousnessProcIDs[rank.RankNumber()]},
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		// The damage spells are procs, so auras without Can Proc From Procs never hear them. They
		// carry no Suppress Weapon Procs, unlike their TBC rows: a "Chance on hit" weapon effect
		// rolls on the seal's hit as well as on the swing, and the game shows both procs landing.
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell | core.SpellFlagProc,
		ClassSpellMask: SpellMaskSealOfRighteousnessProc,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: twoHandCoefficient(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if SealOfRighteousnessFollowsSwappedWeapon {
				spell.BonusCoefficient = twoHandCoefficient()
			}
			mh := paladin.MainHand()
			handMultiplier := core.TernaryFloat64(mh.HandType == proto.HandType_HandTypeTwoHand, 1.2, 0.85)
			baseDamage := value / 100 * handMultiplier * mh.SwingSpeed
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
			dealAfterBatch(sim, spell, result)
		},
	})

	aura := paladin.makeSealExclusive(paladin.MakeProcTriggerAura(core.ProcTrigger{
		Name:            sealLabel("Seal of Righteousness", paladin, rank),
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
		classMask: SpellMaskSealOfRighteousness,
		aura:      aura,
		judgement: judgement,
		echoID:    echoOfRighteousnessID,
		echo: func(sim *core.Simulation, target *core.Unit) {
			procSpell.Cast(sim, target)
		},
	})
}
