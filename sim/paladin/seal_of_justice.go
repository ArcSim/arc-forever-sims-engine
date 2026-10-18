package paladin

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/buffs"
)

// Seal of Justice
// https://www.wowhead.com/forever/spell=20164
//
// Fills the Paladin with the spirit of justice for 30 sec, giving each melee attack a chance to
// stun for 2 sec. Only one Seal can be active on the Paladin at any one time.
//
// Unleashing this Seal's energy will judge an enemy for 10 sec, preventing them from fleeing.
// Your melee strikes will refresh the spell's duration. Only one Judgement per Paladin can be
// active at any one time.
//
// The seal's second effect names its judgement (the row has no effect at the client's index 1).
// Raid bosses are immune to the stun and never flee, so the seal is the aura and the judgement
// is the debuff, nothing more. Twist of Light still names it, so replacing it leaves an Echo of
// Justice with nothing to replay.
func (paladin *Paladin) registerSealOfJustice() {
	rank := spellData.SealOfJustice.Highest()
	judgementRank := spellData.SealOfJusticeTriggered.ByID(int32(rank.EffectN(2).BaseValue()))
	judgementAuras := paladin.newJudgementAuras(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Judgement of Justice",
			ActionID: core.ActionID{SpellID: judgementRank.ID},
			Tag:      buffs.JudgementAuraTag,
			Duration: judgementRank.Duration(),
		})
	})

	// Melee in SpellCategories and Always Hit: the debuff lands without a roll.
	judgement := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: judgementRank.ID},
		SpellSchool:    judgementRank.SpellSchool(),
		DefenseType:    judgementRank.DefenseTypeCore(),
		ProcMask:       debuffJudgementProcMask(),
		Flags:          debuffJudgementFlags(),
		ClassSpellMask: SpellMaskJudgementOfJustice,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
			judgementAuras.Get(target).Activate(sim)
		},

		RelatedAuraArrays: judgementAuras.ToMap(),
	})

	aura := paladin.makeSealExclusive(paladin.RegisterAura(core.Aura{
		Label:    sealLabel("Seal of Justice", paladin, rank),
		ActionID: core.ActionID{SpellID: rank.ID},
		Duration: sealDuration,
	}))

	paladin.registerSealSpell(&sealConfig{
		rank:      rank,
		classMask: SpellMaskSealOfJustice,
		aura:      aura,
		judgement: judgement,
		echoID:    echoOfJusticeID,
	})
}
