package feralcat

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/druid"
)

// Howling Idol (272427) takes 1 second off Shifting Power's 16-second cooldown.
func TestHowlingIdolShortensShiftingPower(t *testing.T) {
	for _, idol := range []int32{0, 272427} {
		items := make([]*proto.ItemSpec, proto.ItemSlot_ItemSlotRanged+1)
		for i := range items {
			items[i] = &proto.ItemSpec{}
		}
		items[proto.ItemSlot_ItemSlotRanged] = &proto.ItemSpec{Id: idol}
		player := &proto.Player{Name: "Cat", Class: proto.Class_ClassDruid, Race: proto.Race_RaceNightElf,
			TalentsString: shiftingPowerOnly, Equipment: &proto.EquipmentSpec{Items: items},
			Spec: DefaultSpecOptions, Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL}}
		sim := core.NewSim(&proto.RaidSimRequest{SimOptions: &proto.SimOptions{RandomSeed: 1},
			Raid:      core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter: core.MakeSingleTargetEncounter(0)}, simsignals.CreateSignals())
		sim.Reset()
		cat := sim.Raid.Parties[0].Players[0].(druid.DruidAgent).GetDruid()
		sim.CurrentTime = time.Second
		// Shifting Power can't be used at full Energy (8 October beta notes, upstream #805).
		cat.SpendEnergy(sim, 50, cat.EnergyRefundMetrics)
		if !cat.ShiftingPower.Cast(sim, cat.CurrentTarget) {
			t.Fatal("cast failed")
		}
		want := 16 * time.Second
		if idol != 0 {
			want = 15 * time.Second
		}
		got := cat.ShiftingPower.TimeToReady(sim)
		t.Logf("idol %d: cooldown %v (want %v)", idol, got, want)
		if got != want {
			t.Errorf("idol %d: cooldown %v, want %v", idol, got, want)
		}
	}
}
