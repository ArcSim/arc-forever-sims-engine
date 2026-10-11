package rogue

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

// Iceblade Hacker's "Melee attacks with this weapon deal 41 Frost damage" is the main hand's alone:
// an off-hand swing beside it must not deal the Frost damage.
func TestIcebladeHackerProcsOnlyFromItsHand(t *testing.T) {
	gear := daggersOnly()
	gear.Items[proto.ItemSlot_ItemSlotMainHand] = &proto.ItemSpec{Id: 13952} // Iceblade Hacker
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceHuman,
		Class:         proto.Class_ClassRogue,
		Equipment:     gear,
		Consumables:   &proto.ConsumesSpec{},
		TalentsString: SubtletyTalents,
		Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, DefaultOptions)
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	rogue := sim.Raid.Parties[0].Players[0].(RogueAgent).GetRogue()

	trigger := rogue.GetAura("Iceblade Hacker")
	if trigger == nil || trigger.Dpm == nil {
		t.Fatal("Iceblade Hacker registered no proc manager")
	}
	for _, mask := range []core.ProcMask{core.ProcMaskMeleeMHAuto, core.ProcMaskMeleeMHSpecial} {
		if got := trigger.Dpm.Chance(mask, sim); got != 1 {
			t.Errorf("a main-hand hit (%v) procs Iceblade Hacker at %v, want every time", mask, got)
		}
	}
	for _, mask := range []core.ProcMask{core.ProcMaskMeleeOHAuto, core.ProcMaskMeleeOHSpecial} {
		if got := trigger.Dpm.Chance(mask, sim); got != 0 {
			t.Errorf("an off-hand hit (%v) procs the main hand's Iceblade Hacker at %v, want never", mask, got)
		}
	}
}

// Iceblade Hacker's Frost (1298414) carries the client's always-hit attribute: it lands on every swing
// that landed, and only rolls to crit. Forever Logs report 2701 shows the server honouring the same
// attribute on a physical attack: Trapz's Summon Hawk dive bomb landed 20 of 20 while Auto Shot
// missed 7 of 75.
func TestIcebladeHackerFrostAlwaysHits(t *testing.T) {
	gear := daggersOnly()
	gear.Items[proto.ItemSlot_ItemSlotMainHand] = &proto.ItemSpec{Id: 13952} // Iceblade Hacker
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceHuman,
		Class:         proto.Class_ClassRogue,
		Equipment:     gear,
		Consumables:   &proto.ConsumesSpec{},
		TalentsString: SubtletyTalents,
		Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, DefaultOptions)
	res := core.RunRaidSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{Iterations: 20, RandomSeed: 1},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	})
	if res.Error != nil {
		t.Fatal(res.Error.Message)
	}

	var landed, avoided, crits int32
	for _, action := range res.RaidMetrics.Parties[0].Players[0].Actions {
		if action.Id.GetSpellId() != 1298414 {
			continue
		}
		for _, target := range action.Targets {
			landed += target.Hits + target.Crits
			avoided += target.Misses + target.Dodges + target.Parries + target.Blocks + target.Glances
			crits += target.Crits
		}
	}
	if landed == 0 {
		t.Fatal("Iceblade Hacker's Frost never landed")
	}
	if avoided != 0 {
		t.Errorf("Iceblade Hacker's Frost was avoided %d times in %d, want never: the client says it always hits", avoided, landed+avoided)
	}
	if crits == 0 {
		t.Errorf("Iceblade Hacker's Frost never crit in %d hits, want it to roll for crits", landed)
	}
}
