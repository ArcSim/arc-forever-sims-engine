package core

import (
	"fmt"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// EnablePerWeaponMainHandImbue gives the item swap's main-hand weapon its own stone or oil
// (ConsumesSpec.mh_imbue_swap_id, mythicsim patch 133). Off, mh_imbue_swap_id is ignored and the
// main-hand imbue counts as before the patch.
var EnablePerWeaponMainHandImbue = true

// staticImbueStats is the character stat a stone or oil gives, for the ones registerStaticImbue
// adds as plain stats.
func staticImbueStats(imbueId int32) stats.Stats {
	var s stats.Stats
	switch imbueId {
	case 25123: // Mana Oil
		s[stats.HealingPower] = 30
		s[stats.MP5] = 15
	case 25122, 20749: // Brilliant Wizard Oil
		s[stats.SpellDamage] = 36
		s[stats.HealingPower] = 36
		s[stats.SpellCritPercent] = 1
	case 25121: // Wizard Oil
		s[stats.SpellDamage] = 24
	case 22756, 18262: // Elemental Sharpening Stone (melee crit only)
		s[stats.PhysicalCritPercent] = 2
		s[stats.RangedCritPercent] = -2
	}
	return s
}

// A main-hand stone or oil that counts only while one of the two weapons is in the main hand.
type weaponImbueEffect struct {
	auras      []*Aura
	flat       float64
	pull, swap bool // carried by the weapon worn at the pull / by the swap's weapon
}

func (effect *weaponImbueEffect) inHand(pull, swap bool) bool {
	return (pull && effect.pull) || (swap && effect.swap)
}

// registerPerWeaponMainHandImbue registers consumables.MhImbueId on the weapon worn at the pull and
// consumables.MhImbueSwapId on the item swap's main-hand weapon, each counting only while its weapon
// is in the main hand: the stat stones and oils (Elemental Sharpening Stone, the wizard oils, Mana
// Oil) as an aura, the proc oils (Frost Oil, StaticImbueProcs) as their proc aura switched on and off,
// and a Dense stone's flat damage on its weapon's swing (a Feral Druid's paw keeps
// MHImbueFlatWeaponDamage, as before). A stone with any other effect (the undead-only Blessed Wizard
// Oil and Consecrated Sharpening Stone) counts all fight. It returns false, leaving the main hand to
// registerStaticImbue as before, when mh_imbue_swap_id is unset or the swap has no main-hand weapon
// of its own.
func registerPerWeaponMainHandImbue(agent Agent, consumables *proto.ConsumesSpec) bool {
	character := agent.GetCharacter()
	if !EnablePerWeaponMainHandImbue || consumables.MhImbueSwapId == 0 || !character.ItemSwap.IsEnabled() {
		return false
	}
	const mainHand = proto.ItemSlot_ItemSlotMainHand
	pullID := character.Equipment[mainHand].ID
	swapID := character.ItemSwap.GetUnequippedItemBySlot(mainHand).ID
	if swapID == 0 || swapID == pullID {
		return false
	}

	var effects []*weaponImbueEffect
	byID := map[int32]*weaponImbueEffect{}
	for _, imbue := range []struct {
		id         int32
		pull, swap bool
	}{
		{consumables.MhImbueId, true, false},
		{max(consumables.MhImbueSwapId, 0), false, true},
	} {
		if imbue.id == 0 {
			continue
		}
		if effect, ok := byID[imbue.id]; ok {
			// The same stone on both weapons: one effect, in hand all fight.
			effect.pull = effect.pull || imbue.pull
			effect.swap = effect.swap || imbue.swap
			continue
		}
		effect := &weaponImbueEffect{pull: imbue.pull, swap: imbue.swap, flat: imbueFlatWeaponDamage(imbue.id)}
		known := effect.flat != 0
		if s := staticImbueStats(imbue.id); s != (stats.Stats{}) {
			known = true
			effect.auras = append(effect.auras, character.NewTemporaryStatsAura(
				fmt.Sprintf("Main Hand Imbue %d", imbue.id), ActionID{ItemID: imbue.id}.WithTag(133), s, NeverExpires).Aura)
		}
		if proc, ok := StaticImbueProcs[imbue.id]; ok {
			known = true
			before := len(character.auras)
			proc(agent, ProcMaskMeleeMH)
			effect.auras = append(effect.auras, character.auras[before:]...)
		}
		if !known {
			registerStaticImbue(agent, imbue.id, character.AutoAttacks.MH(), ProcMaskMeleeMH)
			continue
		}
		byID[imbue.id] = effect
		effects = append(effects, effect)
	}

	character.RegisterItemSwapCallback([]proto.ItemSlot{mainHand}, func(sim *Simulation, _ proto.ItemSlot) {
		inHand := character.MainHand().ID
		pull, swap := inHand == pullID, inHand == swapID
		// The core rebuilds the swing's weapon at a swap in combat and at the reset, not before the pull.
		rebuilt := character.AutoAttacks.AutoSwingMelee && !character.ItemSwap.isFeralDruid && sim.CurrentTime >= 0
		for _, effect := range effects {
			on := effect.inHand(pull, swap)
			for _, aura := range effect.auras {
				if !on {
					aura.Deactivate(sim)
				} else if !aura.IsActive() {
					aura.Activate(sim)
				}
			}
			if on && effect.flat != 0 && rebuilt {
				weapon := character.AutoAttacks.MH()
				weapon.BaseDamageMin += effect.flat
				weapon.BaseDamageMax += effect.flat
			}
		}
	})
	return true
}
