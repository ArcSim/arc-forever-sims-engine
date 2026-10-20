package feralcat

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/druid"
)

// Howling Idol (272427) takes its row's 0.5 seconds off Shifting Power's 16-second cooldown (patch 179); with
// HowlingIdolFromRow off, patch 103's 1 second.
func TestHowlingIdolShortensShiftingPower(t *testing.T) {
	for _, c := range []struct {
		idol    int32
		fromRow bool
		want    time.Duration
	}{{0, true, 16 * time.Second}, {272427, true, 15500 * time.Millisecond}, {272427, false, 15 * time.Second}} {
		idol := c.idol
		druid.HowlingIdolFromRow = c.fromRow
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
		want := c.want
		got := cat.ShiftingPower.TimeToReady(sim)
		t.Logf("idol %d: cooldown %v (want %v)", idol, got, want)
		if got != want {
			t.Errorf("idol %d: cooldown %v, want %v", idol, got, want)
		}
	}
	druid.HowlingIdolFromRow = true
}
