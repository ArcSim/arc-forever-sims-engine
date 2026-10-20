package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

// RevengeThreatFromDamage makes Revenge's threat 2.25 times its damage with no flat bonus (mythicsim
// patch 183, ported from wowsims/forever 913a5bcf69): upstream's owner measured it on Rank 2 in game,
// where the tooltip states only "a high amount of threat". It replaces our damage x1 plus 2 a level,
// fitted to one Rank 1 reading; that reading (46 damage, ~96 threat) is about 2.1 times the damage
// read without the stance, which fits the new value too. Off, damage x1 plus 120.
var RevengeThreatFromDamage = true

const revengeThreatMultiplier = 2.25

func (warrior *Warrior) registerRevenge() {
	// TODO: Manual review needed -- spell 25288 states "a high amount of threat" with no number;
	// none is modelled until measured in game.
	revengeRank := spellData.Revenge.Highest()

	actionID := core.ActionID{SpellID: revengeRank.ID}

	// TODO: In-game test needed
	aura := warrior.RegisterAura(core.Aura{
		Label:    "Revenge",
		Duration: 5 * time.Second,
		ActionID: actionID,
	})

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Revenge - Trigger",
		TriggerImmediately: true,
		Outcome:            core.OutcomeBlock | core.OutcomeDodge | core.OutcomeParry,
		Callback:           core.CallbackOnSpellHitTaken,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			aura.Activate(sim)
		},
	})

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskRevenge,
		MaxRange:       core.MaxMeleeRange,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: revengeRank.GCD(),
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownOf(revengeRank),
			},
		},

		RageCost: core.RageCostOptions{
			Cost:   int32(revengeRank.Cost()),
			Refund: revengeRank.MissRefund(),
		},

		DamageMultiplier: 1,
		// Measured in the level 20 beta (2026-09-28): damage x1 plus a flat bonus of ~40 at rank 1
		// (46 damage, 96 threat in Defensive Stance), before the stance modifier. Was 2.25x damage
		// plus 270. The flat keeps the 2-per-level shape the old value had, which the rank 1
		// reading (2 x 20 = 40) fits; rank 6's own bonus has not been measured.
		ThreatMultiplier: core.TernaryFloat64(RevengeThreatFromDamage, revengeThreatMultiplier, 1),
		FlatThreatBonus:  core.TernaryFloat64(RevengeThreatFromDamage, 0, 2*60),

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(DefensiveStance) && aura.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Rank 6 rolls 153 +-10% (138-168) off the row. Forever adds 25% of attack power, which the
			// client row doesn't carry: 4 level 20 warriors in beta logs (rank 1, 20-24) land 0.23-0.27
			// AP over the base once armor is taken out, on 1H and 2H alike, with Battle Shout up or not.
			baseDamage := revengeRank.DamageEffect().Roll(sim, core.CharacterLevel) + 0.25*spell.MeleeAttackPower(target)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			aura.Deactivate(sim)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},

		RelatedSelfBuff: aura,
	})
}
