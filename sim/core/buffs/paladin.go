package buffs

import (
	"fmt"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

// One rank of a paladin aura: the spell the caster used and the number its row states. Retribution
// Aura's damage shield deals Value; every other aura reads its amount off the spell. The paladin
// registers an aura per rank; each joins the same categories as the generated party-buff copy, which
// is the top rank.
type PaladinAuraRank struct {
	SpellID int32
	Rank    int32
	Value   float64
}

var RetributionAuraMaxRank = PaladinAuraRank{SpellID: 10301, Value: RetributionAuraValue(0)}

// A name with the rank the paladin cast; the max rank the raid config applies carries none.
func paladinRankName(name string, rank int32) string {
	if rank > 0 {
		return fmt.Sprintf("%s Rank %d", name, rank)
	}
	return name
}

// A rank the paladin casts is the party-buff row on the rank's own spell, labelled with the rank.
func paladinAuraMeta(base *Meta, rank PaladinAuraRank) *Meta {
	m := *base
	m.Label = paladinRankName(base.Label, rank.Rank)
	m.Spell = spelldata.MustFind(rank.SpellID)
	return &m
}

func DevotionAuraBuff(char *core.Character, isPlayer bool, rank PaladinAuraRank) *core.Aura {
	return newBuff(&char.Unit, paladinAuraMeta(devotionAuraMeta, rank), isPlayer, 0)
}

func ConcentrationAura(char *core.Character, isPlayer bool, rank PaladinAuraRank) *core.Aura {
	return newBuff(&char.Unit, paladinAuraMeta(concentrationAuraMeta, rank), isPlayer, 0)
}

func FireResistanceAura(char *core.Character, isPlayer bool, rank PaladinAuraRank) *core.Aura {
	return newBuff(&char.Unit, paladinAuraMeta(fireResistanceAuraMeta, rank), isPlayer, 0)
}

func FrostResistanceAura(char *core.Character, isPlayer bool, rank PaladinAuraRank) *core.Aura {
	return newBuff(&char.Unit, paladinAuraMeta(frostResistanceAuraMeta, rank), isPlayer, 0)
}

func ShadowResistanceAura(char *core.Character, isPlayer bool, rank PaladinAuraRank) *core.Aura {
	return newBuff(&char.Unit, paladinAuraMeta(shadowResistanceAuraMeta, rank), isPlayer, 0)
}

// Retribution Aura scales with the paladin's Holy spell power at 6%, though its client row carries
// no coefficient (EffectBonusCoefficient 0 on every rank). Blizzard's beta notes of 8 Oct 2026 set
// it: "the spell power ratio has been reduced to 6% (was 13.3%)", the base back to Vanilla's (client
// 70291 has rank 5 at 20), and a fix so it reads "the casting Paladin's spell power" instead of the
// struck member's. That retires the level-20 reading of 0.1357 and the holder split in beta logs
// 2729-2736, both from the build before.
const RetributionAuraSpellPowerCoefficient = 0.06

// RetributionAuraBuff is the aura on the unit the shield protects. The paladin's own copy carries
// the paladin's spell power. The party option's copy deals the base: it is on a party member, whose
// spell power is not the paladin's, and the party config does not say the paladin's.
func RetributionAuraBuff(char *core.Character, isPlayer bool, rank PaladinAuraRank) *core.Aura {
	m := paladinAuraMeta(retributionAuraMeta, rank)
	label := m.label(isPlayer)
	if char.HasAura(label) {
		return char.GetAura(label)
	}

	aura := core.NewDamageShield(&char.Unit, label, m.actionID(isPlayer), core.NeverExpires, m.Category, m.SingleAura,
		core.SpellSchoolHoly, rank.Value, core.TernaryFloat64(isPlayer, RetributionAuraSpellPowerCoefficient, 0))
	core.JoinSharedCategory(aura, m.SharedCategory, isPlayer)
	return aura
}

// One rank of a judgement debuff: the spell the target shows and the number its row states. The
// paladin registers a rank per row; the debuff panel applies the max rank, carried by the *MaxRank
// values with Rank 0.
type JudgementRank struct {
	SpellID int32
	Rank    int32
	Value   float64
}

var (
	JudgementOfTheCrusaderMaxRank = JudgementRank{SpellID: judgementOfTheCrusaderSpell.ID,
		Value: amount(judgementOfTheCrusaderSpell.Effect(dbcenums.A_MOD_DAMAGE_TAKEN, int32(core.SpellSchoolHoly)))}
	JudgementOfLightMaxRank = JudgementRank{SpellID: judgementOfLightSpell.ID,
		Value: amount(spelldata.MustFind(judgementOfLightHealID).HealEffect())}
	JudgementOfWisdomMaxRank = JudgementRank{SpellID: judgementOfWisdomSpell.ID,
		Value: amount(spelldata.MustFind(judgementOfWisdomManaID).EnergizeEffect())}
)

// The raid's Judgement of the Crusader, which no generated row reads: the client states its bonus
// for the holy school alone, which the manifest has no pseudo-stat for.
var judgementOfTheCrusaderSpell = spelldata.MustFind(20303)

// The top ranks' heal and mana. The judgement states neither: its trigger is a dummy, and the
// spells the paladin's own ranks pair with it by hand carry the amounts.
const (
	judgementOfLightHealID  = 20343
	judgementOfWisdomManaID = 20353
)

// Every judgement debuff a paladin puts up carries the tag, so an effect that refreshes "all
// Judgement effects on the target" can find them.
const JudgementAuraTag = "JudgementAura"

// The client says the judgements proc on a chance without stating it; the sim uses 50% until
// in-game testing says otherwise.
const judgementProcChance = 0.5

// Judgement of the Crusader raises the Holy damage the target takes by a flat amount. Every rank
// and every paladin share one exclusive category, so the strongest active one is the one that
// counts.
func JudgementOfTheCrusaderAura(target *core.Unit, rank JudgementRank) *core.Aura {
	bonus := rank.Value
	label := paladinRankName("Judgement of the Crusader", rank.Rank)
	if target.HasAura(label) {
		return target.GetAura(label)
	}

	aura := target.GetOrRegisterAura(core.Aura{
		Label:    label,
		ActionID: core.ActionID{SpellID: rank.SpellID},
		Tag:      JudgementAuraTag,
		Duration: auraDuration(judgementOfTheCrusaderSpell),
	})

	aura.NewExclusiveEffect("Judgement of the Crusader", true, core.ExclusiveEffect{
		Priority: bonus,
		OnGain: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			target.PseudoStats.SchoolBonusSpellDamage[stats.SchoolIndexHoly] += bonus
		},
		OnExpire: func(ee *core.ExclusiveEffect, sim *core.Simulation) {
			target.PseudoStats.SchoolBonusSpellDamage[stats.SchoolIndexHoly] -= bonus
		},
	})

	return aura
}

// The paladin's own Judgement of Light at one rank: the bare 40-second debuff, plus the heal it
// grants whoever strikes the target.
func JudgementOfLightRankAura(target *core.Unit, rank JudgementRank) *core.Aura {
	label := paladinRankName("Judgement of Light", rank.Rank)
	if target.HasAura(label) {
		return target.GetAura(label)
	}

	return AttachJudgementOfLightHeal(target.GetOrRegisterAura(core.Aura{
		Label:    label,
		ActionID: core.ActionID{SpellID: rank.SpellID},
		Tag:      JudgementAuraTag,
		Duration: JudgementOfLightDuration(0),
	}), rank)
}

// Whoever lands a melee hit on the judged target has a chance to be healed for the rank's amount.
func AttachJudgementOfLightHeal(aura *core.Aura, rank JudgementRank) *core.Aura {
	healthMetrics := aura.Unit.NewHealthMetrics(core.ActionID{SpellID: rank.SpellID})
	heal := rank.Value

	return aura.AttachProcTrigger(core.ProcTrigger{
		Name:     aura.Label + " - Heal",
		Callback: core.CallbackOnSpellHitTaken,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.Proc(judgementProcChance, "Judgement of Light - Heal") {
				spell.Unit.GainHealth(sim, heal, healthMetrics)
			}
		},
	})
}

// The paladin's own Judgement of Wisdom at one rank.
func JudgementOfWisdomRankAura(target *core.Unit, rank JudgementRank) *core.Aura {
	label := paladinRankName("Judgement of Wisdom", rank.Rank)
	if target.HasAura(label) {
		return target.GetAura(label)
	}

	return AttachJudgementOfWisdomMana(target.GetOrRegisterAura(core.Aura{
		Label:    label,
		ActionID: core.ActionID{SpellID: rank.SpellID},
		Tag:      JudgementAuraTag,
		Duration: JudgementOfWisdomDuration(0),
	}), rank)
}

// Whoever lands an attack or spell on the judged target has a chance to regain the rank's mana.
// Melee claim it returns mana on a miss as well.
func AttachJudgementOfWisdomMana(aura *core.Aura, rank JudgementRank) *core.Aura {
	actionID := core.ActionID{SpellID: rank.SpellID}
	mana := rank.Value

	return aura.AttachProcTrigger(core.ProcTrigger{
		Name:            aura.Label,
		ActionID:        actionID,
		MetricsActionID: actionID,
		ProcChance:      judgementProcChance,
		ProcMask:        core.ProcMaskDirect,
		Callback:        core.CallbackOnSpellHitTaken,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) && !result.Landed() {
				return
			}

			unit := spell.Unit
			if !unit.HasManaBar() {
				return
			}
			if unit.JowManaMetrics == nil {
				unit.JowManaMetrics = unit.NewManaMetrics(actionID)
			}
			unit.AddMana(sim, mana, unit.JowManaMetrics)
		},
	})
}
