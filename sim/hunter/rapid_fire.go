package hunter

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (hunter *Hunter) registerRapidFire() {
	if hunter.Level < 26 {
		return
	}

	actionID := core.ActionID{SpellID: 3045}
	// Forever: Rapid Killing (415405) takes 2 min off at 2/2 (the data stores the max-rank value).
	cooldown := time.Minute*5 - time.Minute*time.Duration(hunter.Talents.RapidKilling)

	hunter.RapidFireAura = hunter.RegisterAura(core.Aura{
		Label:    "Rapid Fire",
		ActionID: actionID,
		Duration: time.Second * 15,

		// Forever Rapid Fire: +40% ranged and +40% melee attack speed (effects 1 and 4).
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyRangedSpeed(sim, 1.4)
			aura.Unit.MultiplyMeleeSpeed(sim, 1.4)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyRangedSpeed(sim, 1/1.4)
			aura.Unit.MultiplyMeleeSpeed(sim, 1/1.4)
		},
	})

	hunter.RapidFire = hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ManaCost: core.ManaCostOptions{
			FlatCost: 100,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: cooldown,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.RapidFireAura.Activate(sim)
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: hunter.RapidFire,
		Type:  core.CooldownTypeDPS,
	})
}
