package hunter

import (
	"github.com/wowsims/classic/sim/core"
)

// Hunter's Mark (rank 4), castable for a rotation's pull: it puts up the same mark as the raid debuff setting
// (core.HuntersMarkAura, 2 min), so with that setting on it only refreshes it.
func (hunter *Hunter) registerHuntersMarkSpell() {
	marks := hunter.NewEnemyAuraArray(core.HuntersMarkAura)

	hunter.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 14325},
		Flags:    core.SpellFlagAPL,
		ManaCost: core.ManaCostOptions{FlatCost: 60},
		Cast:     core.CastConfig{DefaultCast: core.Cast{GCD: core.GCDDefault}},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			marks.Get(target).Activate(sim)
		},
	})
}
