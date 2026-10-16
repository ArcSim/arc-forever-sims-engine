package enhancement

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

// Stormstrike costs 8% of base mana, not a flat 125 (16 October 2026 beta notes).
func TestStormstrikeCostsEightPercentOfBaseMana(t *testing.T) {
	player := &proto.Player{
		Name: "enh", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc, TalentsString: DefaultTalents,
		Equipment: &proto.EquipmentSpec{},
		Spec: &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{
			ClassOptions: &proto.ShamanOptions{},
		}}},
		Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	char := sim.Raid.Parties[0].Players[0].GetCharacter()
	spell := char.GetSpell(core.ActionID{SpellID: 17364})
	if spell == nil {
		t.Fatal("no Stormstrike on the default Enhancement build")
	}
	want := int32(8*char.BaseMana) / 100
	if want < 50 || spell.Cost.BaseCost != want {
		t.Fatalf("Stormstrike costs %d, want 8%% of %v base mana (%d)", spell.Cost.BaseCost, char.BaseMana, want)
	}
}
