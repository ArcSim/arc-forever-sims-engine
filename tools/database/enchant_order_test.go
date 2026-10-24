package database

import (
	"slices"
	"testing"

	"github.com/wowsims/classic/sim/core/proto"
)

// Enchants sharing an effect id and slot (Forever gave some two-hand enchants a one-hand enchant's effect id) come out of
// the generator in the same order every run, so db.json does not reshuffle and the site shows the same names.
func TestEnchantOrderIsStable(t *testing.T) {
	db := NewWowDatabase()
	for _, e := range []*proto.UIEnchant{
		{EffectId: 805, SpellId: 13943, Name: "Enchant Weapon - Greater Striking", Type: proto.ItemType_ItemTypeWeapon},
		{EffectId: 805, SpellId: 7745, Name: "Enchant 2H Weapon - Minor Impact", Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
		{EffectId: 1897, SpellId: 20031, Name: "Enchant Weapon - Superior Striking", Type: proto.ItemType_ItemTypeWeapon},
		{EffectId: 1897, SpellId: 13529, Name: "Enchant 2H Weapon - Lesser Impact", Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
		{EffectId: 2563, SpellId: 23799, Name: "Enchant Weapon - Strength", Type: proto.ItemType_ItemTypeWeapon},
		{EffectId: 2563, SpellId: 1248511, Name: "Enchant 2H Weapon - Lesser Strength", Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
		{EffectId: 2563, SpellId: 1248640, Name: "Enchant Gloves - Superior Strength", Type: proto.ItemType_ItemTypeHands},
		{EffectId: 2563, ItemId: 1, SpellId: 1, Name: "same slot and kind, told apart by item", Type: proto.ItemType_ItemTypeWeapon},
	} {
		db.Enchants[EnchantToDBKey(e)] = e
	}
	names := func() []string {
		var out []string
		for _, e := range db.ToUIProto().Enchants {
			out = append(out, e.Name)
		}
		return out
	}
	want := names()
	for i := 0; i < 200; i++ {
		if got := names(); !slices.Equal(got, want) {
			t.Fatalf("enchant order changed between runs:\n%q\n%q", want, got)
		}
	}
	// within an effect id and slot, the plain enchant comes before the two-hand one
	if want[0] != "Enchant Weapon - Greater Striking" || want[1] != "Enchant 2H Weapon - Minor Impact" {
		t.Errorf("effect 805 order %q, want the weapon enchant first", want[:2])
	}
}
