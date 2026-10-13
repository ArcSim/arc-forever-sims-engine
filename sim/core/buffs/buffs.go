// Package buffs holds the raid buffs and debuffs: the constructors tools/database/gen_spelldata
// generates from tools/database/buffmanifest, the drivers that add what the client does not state,
// and the hand-written auras that share categories with them. They sit above sim/core so that they
// can read the spell store, which imports sim/core, so core cannot call them: init registers them
// as core's buff hooks, and sim/common imports the package so that every sim links it.
package buffs

import (
	"github.com/wowsims/forever/sim/core"
)

func init() {
	core.RegisterBuffHooks(core.BuffHooks{
		ApplyBuffs:    applyGeneratedBuffs,
		ApplyPetBuffs: applyPetPartyBuffs,
		ApplyDebuffs:  applyDebuffs,
		GiftOfArthasAura: func(target *core.Unit) *core.Aura {
			return GiftOfArthasAura(target, true, 0)
		},
	})
}
