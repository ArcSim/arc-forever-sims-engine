//go:build with_db

package core

import (
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// Forever's caster weapons carry spell power their client rows do not, and master's database states
// "damage and healing" as SpellDamage alone (tools/database/gen_db/planner_spell_power.go). Checks
// the embedded binary used by real simulations.
func TestForeverItemSpellPower(t *testing.T) {
	for id, want := range map[int32][2]float64{ // SpellDamage, HealingPower
		19102: {99, 99},   // Crackling Staff: 25 from the client row plus 74, as in game
		19884: {94, 94},   // Jin'do's Judgement
		20069: {88, 88},   // Ironbark Staff: the caster spell power row (patch 180), as Wowhead Forever shows
		18608: {65, 195},  // Benediction: a healer's weapon, a third as damage
		23451: {106, 106}, // Grand Marshal's Mageblade
		18534: {23, 23},   // Rod of the Ogre Magi, from master's database: healing was missing
		22383: {94, 94},   // Sageblade
		19379: {44, 44},   // Neltharion's Tear, from the client: unchanged
	} {
		item := NewItem(ItemSpec{ID: id})
		if got := [2]float64{item.Stats[stats.SpellDamage], item.Stats[stats.HealingPower]}; got != want {
			t.Errorf("%d %s spell damage, healing = %v, want %v", id, item.Name, got, want)
		}
	}
	if got := NewItem(ItemSpec{ID: 20536}).Stats[stats.ShadowDamage]; got != 73 {
		t.Errorf("Soul Harvester shadow damage = %v, want 73", got)
	}
}

// Ironbark Staff's 2% critical strike is stated as the generic rating, which the database reads as
// melee crit rating (28). Forever pays gear crit to spells too, so a caster wearing it has 2% spell crit.
func TestIronbarkStaffCritCountsForSpells(t *testing.T) {
	var equipment Equipment
	equipment[proto.ItemSlot_ItemSlotMainHand] = NewItem(ItemSpec{ID: 20069})
	got := equipment.Stats(proto.Spec_SpecWarlock)
	if got[stats.SpellCritRating] != 28 || got[stats.MeleeCritRating] != 28 {
		t.Errorf("Ironbark Staff crit rating: spell %v, melee %v, want 28 and 28", got[stats.SpellCritRating], got[stats.MeleeCritRating])
	}
}
