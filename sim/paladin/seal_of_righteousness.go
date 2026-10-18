package paladin

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// The damage spell each rank fires on a hit. They carry no rank subtext, so no ladder holds them.
//
// Rank 1 is 21084, the row that names Judgement of Righteousness rank 1 (20187), and fires 25741.
// 25742 belongs to 20154, the older Rank 1 row with no judgement, which no ladder registers: the
// client's tooltips label 25742 "Rank 1" and 25741 "Rank 1.5" (20154 and 21084 in that order), and
// the logs show a rank 1 seal hitting as 25741 (foreverlogs 2768, Amicus). MythicSim patch 158.
var sealOfRighteousnessProcIDs = map[int32]int32{1: 25741, 2: 25740, 3: 25739, 4: 25738, 5: 25737, 6: 25736, 7: 25735, 8: 25713}

// SealOfRighteousnessFollowsSwappedWeapon reads the main hand at each hit for the spell power share,
// as the base damage already does, so a swap between a two-hander and a one-hander (a Manual Crowd
// Pummeler opener) changes the coefficient with the weapon. False keeps the coefficient of the
// weapon worn at the pull for the whole fight, as it was registered before (MythicSim patch 129).
var SealOfRighteousnessFollowsSwappedWeapon = true

// SealOfRighteousnessTriggersWeaponProcs lets the seal's hit roll "Chance on hit" weapon effects and
// combat enchants (Crusader, Wolfsbane, Windfury Totem) as a swing does (MythicSim patch 160). The proc
// rows (25713 and down, unchanged in client 70425) carry no Suppress Weapon Procs and no Suppress
// Caster Procs, unlike their TBC rows, and the pre-16 October logs show it: Wolfsbane's Blazewind Blast
// lands on the seal hit's millisecond, 20 to 220 ms after the swing that fired the seal, 21 times in
// three reports (2756 George, 2670 Jack, 2755 Conde), at about the rate it lands on swings. Equip
// effects, talents and racials are aura procs and never hear the seal (it is a proc without Not a
// Proc): Touch of the Grave, a 20 to 45 ms aura proc, followed swings, Holy Strikes and judgements 89
// times in five reports (963 seal hits) and a seal hit alone never. False adds Suppress Weapon Procs to the hit, the reading of
// Lightclub's 16 October test ("does not proc equip effects") if it covered chance on hit effects.
var SealOfRighteousnessTriggersWeaponProcs = true

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
// under the debuff, so a one-hander keeps the bare 0.185 (Ms's 1.7 one-hander fits it). Those logs
// are all from before Blizzard's 16 October fix, so they measure the bugged share; patch 158
// (SealOfRighteousnessSpellPowerFixed) halves it.

// JudgementOfRighteousnessSpellPowerCoefficient is the judgement's spell power share as measured in
// game on client 70425 (MythicSim patch 143); the rows carry 0.429. The seal's per-hit share is the
// dummy's, halved (SealOfRighteousnessSpellPowerFixed, MythicSim patch 158).
const JudgementOfRighteousnessSpellPowerCoefficient = 0.70

// JudgementOfRighteousnessSpellPowerFromRow gives the judgement the share its client row carries, 0.429 on
// every rank in client 70425 (0.50 before it), instead of patch 143's measured 0.70 (MythicSim patch 161).
// Lightclub prints "42.9% of spellpower (49.3% improved)", the row's 0.429 and that times Improved Seals'
// 15%. The pre-16 October logs show the server dealing the row's share of the day: Judgement of the
// Crusader steps of 0.51 to 0.53 a point on rank 2 and 3 judgements of Retribution paladins with no damage
// talent for them (reports 2693, 2744, 2756), against that client's 0.50. No log from after 16 October
// exists. Lightclub's 16 October test on the live server gives 0.429 (49.3% with Improved Seals, whose
// damage mod scales the whole hit, spell power included), so the switch is on (MythicSim patch 168).
// False gives patch 143's 0.70.
var JudgementOfRighteousnessSpellPowerFromRow = true

// The judgement's spell power share under the switch, from its damage effect.
func judgementOfRighteousnessCoefficient(damage *spelldata.Effect) float64 {
	if JudgementOfRighteousnessSpellPowerFromRow {
		return damage.Coeff()
	}
	return JudgementOfRighteousnessSpellPowerCoefficient
}

// SealOfRighteousnessBasesAtSpellLevel reads the judgement's damage and the seal's per-hit value at
// the rank's own spell level, without the rows' per-level gain (MythicSim patch 148). Ashreon's bugs
// forum post ("Oh the prot paladin nerfs keeps giving", 17 October 2026) gives rank 8 as "21-71 /
// 162-178": the seal's tooltip range ($20286s3/87 to $20286s3/25) and the judgement's. Both are the
// rows' base points exactly, 1786/87 = 20.5 and 1786/25 = 71.4, and 170 +/- 4.7% = 162.0 to 178.0;
// the engine read them at level 60, two levels past rank 8's spell level 58: 1880 a hundred speed
// (21.6 to 75.2) and 178 (170 to 186). False keeps the per-level gain.
var SealOfRighteousnessBasesAtSpellLevel = true

// SealOfRighteousnessSpellPowerFixed halves the seal hit's spell power share, Blizzard's 16 October
// fix: "Fixed an issue causing Seal of Righteousness to gain almost twice as much benefit from Spell
// Damage than intended" (MythicSim patch 158). Every public log is from before the fix and gives the
// judgement dummy's share, 0.18 to 0.25 a hit on one- and two-handers (the 18 October log review), so
// the dummy's number is the bugged one. Half of it is 0.10 a hit from rank 4 with a one-hander (0.0925
// at rank 3) and 0.12 with a two-hander (the dummy x1.2, halved): the two-hander's tooltip
// ($<2h> = 0.12) exactly, and the one-hander's logs halved (0.09 to 0.11), where the one-hander's
// tooltip ($<1h> = 0.0765) is 15-30% under them. False keeps the dummy's share, the pre-fix server.
var SealOfRighteousnessSpellPowerFixed = true

// SealOfRighteousnessPerHitTooltipCoefficient gives each seal hit the spell power share client
// 70425's tooltip prints, $<1h> = 0.0765 with a one-hander and $<2h> = 0.12 with a two-hander
// (SpellDescriptionVariables 1023), at every rank (MythicSim patch 148). When on it replaces both the
// dummy's share and patch 158's half of it, so SealOfRighteousnessSpellPowerFixed does nothing; the two
// agree on a two-hander from rank 4 and differ on a one-hander (0.0765 against 0.10). Off by default:
// it is the other reading of the fixed one-hander until a log from after 16 October settles it.
var SealOfRighteousnessPerHitTooltipCoefficient = false

const (
	sealOfRighteousnessTooltipOneHandCoefficient = 0.0765
	sealOfRighteousnessTooltipTwoHandCoefficient = 0.12
)

// SealOfRighteousnessHitCoefficient is the seal hit's spell power share under the switches as they
// stand, from the judgement dummy's coefficient (0.058, 0.125, 0.185 and from rank 4 0.20) and the
// main hand. The tooltip switch comes first; then the dummy's share, x1.2 on a two-hander, is halved
// when SealOfRighteousnessSpellPowerFixed is on.
func SealOfRighteousnessHitCoefficient(dummy float64, twoHand bool) float64 {
	return sealOfRighteousnessHitCoefficient(dummy, twoHand, SealOfRighteousnessPerHitTooltipCoefficient, SealOfRighteousnessSpellPowerFixed)
}

func sealOfRighteousnessHitCoefficient(dummy float64, twoHand, tooltip, fixed bool) float64 {
	if tooltip {
		return core.TernaryFloat64(twoHand, sealOfRighteousnessTooltipTwoHandCoefficient, sealOfRighteousnessTooltipOneHandCoefficient)
	}
	coefficient := dummy
	if twoHand {
		coefficient *= 1.2
	}
	if fixed {
		coefficient /= 2
	}
	return coefficient
}

// The level a Seal of Righteousness row is read at.
func sealOfRighteousnessLevel(e *spelldata.Effect) int32 {
	if SealOfRighteousnessBasesAtSpellLevel {
		return int32(e.SpellLevel)
	}
	return core.CharacterLevel
}

func (paladin *Paladin) registerSealOfRighteousness(_ int32, rank *spelldata.Spell) {
	judgeRank := spellData.JudgementOfRighteousness.ByID(int32(rank.EffectN(2).BaseValue()))
	judgeDamage := judgeRank.DamageEffect()
	perHit := judgeRank.EffectN(2)
	judgeLevel := sealOfRighteousnessLevel(judgeDamage)

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
		BonusCoefficient: judgementOfRighteousnessCoefficient(judgeDamage),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, judgeDamage.Roll(sim, judgeLevel), spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})

	value := perHit.Average(sealOfRighteousnessLevel(perHit))
	coefficient := perHit.Coeff()
	tooltipCoefficient := SealOfRighteousnessPerHitTooltipCoefficient
	spellPowerFixed := SealOfRighteousnessSpellPowerFixed
	twoHandCoefficient := func() float64 {
		twoHand := paladin.GetMainHandType() == proto.HandType_HandTypeTwoHand
		return sealOfRighteousnessHitCoefficient(coefficient, twoHand, tooltipCoefficient, spellPowerFixed)
	}
	procSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: sealOfRighteousnessProcIDs[rank.RankNumber()]},
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		// The damage spells are procs, so auras without Can Proc From Procs never hear them. They
		// carry no Suppress Weapon Procs, unlike their TBC rows: a "Chance on hit" weapon effect
		// rolls on the seal's hit as well as on the swing, and the logs show both procs landing
		// (SealOfRighteousnessTriggersWeaponProcs).
		Flags: core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell | core.SpellFlagProc |
			core.Ternary(SealOfRighteousnessTriggersWeaponProcs, core.SpellFlagNone, core.SpellFlagSuppressWeaponProcs),
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
