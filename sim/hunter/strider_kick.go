package hunter

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
)

// StriderKickSpeedsUp (MythicSim patch 120) gives Strider Kick its movement speed: the client row's third
// effect (1317257 effect 2, A_MOD_INCREASE_SPEED 30 on the caster for the spell's 3 sec), and the Sep 24
// beta notes ("Strider Kick now increases movement speed by 30% for 3 seconds."). Without it the Hunter
// walks out of melee at the base 7 yards a second after a kick; it only matters to a rotation that moves.
// Set it false for the kick without the speed.
var StriderKickSpeedsUp = true

// Strider Kick from the beta client (1317257): 100% normalized melee weapon damage, 8 sec cooldown,
// 5.81% of base mana.
func (hunter *Hunter) registerStriderKickSpell() {
	if !hunter.Talents.StriderKick {
		return
	}

	rank := spellData.StriderKick.Highest()

	var speed *core.Aura
	if e := rank.Effect(dbcenums.A_MOD_INCREASE_SPEED, 0); StriderKickSpeedsUp && e != nil {
		speed = hunter.RegisterAura(core.Aura{
			Label:    "Strider Kick",
			ActionID: core.ActionID{SpellID: rank.ID},
			Duration: rank.Duration(),
		})
		speed.NewActiveMovementSpeedEffect(e.Percent())
	}

	hunter.StriderKick = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ClassSpellMask: HunterSpellStriderKick,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MaxRange:       float64(rank.MaxRange),

		ManaCost: rank.ManaCost(),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD(),
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: max(rank.Cooldown(), rank.CategoryCooldown()),
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := hunter.AutoAttacks.MH().CalculateNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			if speed != nil {
				speed.Activate(sim)
			}
		},
	})
}
