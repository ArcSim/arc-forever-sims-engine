package forever

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// EnableFrostOil registers Frost Oil as a weapon imbue (mythicsim patch 126). Off, the imbue IDs
// below do nothing, as before the patch.
var EnableFrostOil = true

// Frost Oil (item 3829, use spell 3595) puts temporary enchantment 26 on a melee weapon. The client's
// SpellItemEnchantment row for 26 casts 350025, a Frostbolt that rolls 33 + 0.7 a level (levels 8 to
// 12) with 12% variance and 0.269 of spell power, at the 10% chance the client states for it
// (ProcChance 10 on 350025; the tooltip's "10% chance of casting Frostbolt at the opponent when it
// hits"). It is a flat chance on each landed hit of the oiled weapon, not a rate per minute.
const (
	frostOilUseSpell = 3595
	frostOilItem     = 3829
	frostOilBolt     = 350025
)

func init() {
	for _, id := range []int32{frostOilUseSpell, frostOilItem} {
		core.StaticImbueProcs[id] = registerFrostOil
	}
}

func registerFrostOil(agent core.Agent, hand core.ProcMask) {
	if !EnableFrostOil {
		return
	}
	character := agent.GetCharacter()
	row := spelldata.MustFind(frostOilBolt)
	bolt := character.GetOrRegisterSpell(shared.SpellDataProcDamageSpell(character, row))
	bolt.Flags |= core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagProc
	chance := float64(row.ProcChance) / 100

	name := "Frost Oil (main hand)"
	if hand == core.ProcMaskMeleeOH {
		name = "Frost Oil (off hand)"
	}
	character.MakeProcTriggerAura(core.ProcTrigger{
		Name:               name,
		Callback:           core.CallbackOnSpellHitDealt,
		ProcMask:           hand,
		IsWeaponProc:       true,
		Outcome:            core.OutcomeLanded,
		ProcChance:         chance,
		TriggerImmediately: true,
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			bolt.Cast(sim, result.Target)
		},
	})
}
