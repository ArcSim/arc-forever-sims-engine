package hunter

import (
	"os"
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (hunter *Hunter) registerStriderKickSpell() {
	if !hunter.Talents.StriderKick {
		return
	}

	// Forever data: Strider Kick also grants +30% movement speed for 3 s (effect 3, duration index 27).
	// AB_NO_SK_SPEED=1 turns it off for comparisons.
	noStriderSpeed := os.Getenv("AB_NO_SK_SPEED") != ""
	speedID := core.ActionID{SpellID: 1317257}.WithTag(1)
	speedAura := hunter.RegisterAura(core.Aura{
		Label:    "Strider Kick Speed",
		ActionID: speedID,
		Duration: time.Second * 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddMoveSpeedModifier(&speedID, 1.3)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.RemoveMoveSpeedModifier(&speedID)
		},
	})

	hunter.StriderKick = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1317257}, // Classic icon; change when you have the Forever ID
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | SpellFlagStrike,

		RageCost: core.RageCostOptions{}, // ignore if hunter uses mana
		ManaCost: core.ManaCostOptions{
			BaseCost: 0.0581,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{GCD: core.GCDDefault},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget <= 5
		},
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			base := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			spell.CalcAndDealDamage(sim, target, base, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			if !noStriderSpeed {
				speedAura.Activate(sim)
			}
		},
	})
}
