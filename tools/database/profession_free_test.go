package database

import (
	"os"
	"testing"

	"github.com/wowsims/classic/sim/core/proto"
)

// Forever lets anyone wear the Zul'Gurub crafted pieces Classic kept for their crafters (ItemSparse RequiredSkill 0 in
// client build 1.60.1.70450): the shipped item database asks no profession for them, while an item Forever still locks
// keeps its profession.
func TestForeverProfessionFreeItems(t *testing.T) {
	data, err := os.ReadFile("../../assets/database/db.json")
	if err != nil {
		t.Fatal(err)
	}
	db := ReadDatabaseFromJson(string(data))
	unlocked := map[int32]string{
		19682: "Bloodvine Vest", 19683: "Bloodvine Leggings", 19684: "Bloodvine Boots",
		19685: "Primal Batskin Jerkin", 19686: "Primal Batskin Gloves", 19687: "Primal Batskin Bracers",
		19688: "Blood Tiger Breastplate", 19689: "Blood Tiger Shoulders",
		19693: "Darksoul Breastplate", 19694: "Darksoul Leggings", 19695: "Darksoul Shoulders",
	}
	for id, name := range unlocked {
		item, ok := db.Items[id]
		if !ok {
			t.Errorf("%s (%d) is not in the item database", name, id)
			continue
		}
		if item.RequiredProfession != proto.Profession_ProfessionUnknown {
			t.Errorf("%s (%d) still asks for profession %v", name, id, item.RequiredProfession)
		}
	}
	if item, ok := db.Items[16022]; !ok || item.RequiredProfession != proto.Profession_Engineering {
		t.Error("Arcanite Dragonling (16022) lost its Engineering lock, which Forever keeps")
	}
}
