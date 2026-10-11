package feralbear

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/druid"
)

// In Forever a druid drinks a potion in Bear Form and stays in form, so Mighty
// Rage Potion (13442) lands its rage on the bear.
func TestBearDrinksMightyRageInForm(t *testing.T) {
	player := &proto.Player{Name: "Potion", Class: proto.Class_ClassDruid, Race: proto.Race_RaceNightElf,
		Equipment:   &proto.EquipmentSpec{},
		Consumables: &proto.ConsumesSpec{PotId: 13442, Potions: []int32{13442}},
		Spec:        DefaultSpecOptions, Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL}}
	sim := core.NewSim(&proto.RaidSimRequest{SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid:      core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter: core.MakeSingleTargetEncounter(0)}, simsignals.CreateSignals())
	sim.Reset()
	d := sim.Raid.Parties[0].Players[0].(druid.DruidAgent).GetDruid()
	sim.CurrentTime = time.Second
	d.ClearForm(sim)
	if !d.BearForm.Cast(sim, d.CurrentTarget) {
		t.Fatal("Bear Form failed")
	}
	potion := d.GetSpell(core.ActionID{ItemID: 13442})
	if potion == nil {
		t.Fatal("no Mighty Rage Potion")
	}
	before := d.CurrentRage()
	if !potion.Cast(sim, d.CurrentTarget) {
		t.Fatal("Mighty Rage Potion failed")
	}
	if !d.InForm(druid.Bear) {
		t.Fatal("the potion dropped Bear Form")
	}
	if d.CurrentRage() <= before {
		t.Fatalf("rage %v after the potion, %v before", d.CurrentRage(), before)
	}
}
