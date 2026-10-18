package retribution

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/paladin"
)

// MythicSim patch 160: the seal's hit rolls a weapon proc (Fiery Weapon on Arcanite Reaper; the logs
// show Wolfsbane's Blazewind Blast, which the test item database lacks) as a swing does, unless SealOfRighteousnessTriggersWeaponProcs is off. With only Seal of
// Righteousness in the rotation every landed swing fires one seal hit, so the switch on gives about
// twice the procs it gives off.
func TestSealOfRighteousnessTriggersWeaponProcs(t *testing.T) {
	if !paladin.SealOfRighteousnessTriggersWeaponProcs {
		t.Fatal("SealOfRighteousnessTriggersWeaponProcs is off by default; patch 160 leaves it on")
	}
	const (
		arcaniteReaper = 12784
		fieryWeapon    = 803
		fireball       = 13897
	)
	run := func() (procs, swings float64) {
		items := make([]*proto.ItemSpec, 15)
		for i := range items {
			items[i] = &proto.ItemSpec{}
		}
		items[14] = &proto.ItemSpec{Id: arcaniteReaper, Enchant: fieryWeapon}
		player := core.WithSpec(&proto.Player{
			Race:          proto.Race_RaceHuman,
			Class:         proto.Class_ClassPaladin,
			Equipment:     &proto.EquipmentSpec{Items: items},
			TalentsString: RetTalents,
			Rotation: core.APLRotationFromJsonString(`{"type":"TypeAPL","priorityList":[
				{"action":{"castSpell":{"spellId":{"spellId":20293}},"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":20293}}}}}}}
			]}`),
		}, DefaultOptions)
		result := core.RunRaidSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 160, Iterations: 200},
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
		})
		if result.Error != nil {
			t.Fatal(result.Error.Message)
		}
		for _, action := range result.RaidMetrics.Parties[0].Players[0].Actions {
			for _, target := range action.Targets {
				switch {
				case action.Id.GetSpellId() == fireball:
					procs += float64(target.Hits + target.Crits)
				case action.Id.GetOtherId() == proto.OtherAction_OtherActionAttack:
					swings += float64(target.Hits + target.Crits)
				}
			}
		}
		return procs, swings
	}

	on, swingsOn := run()
	var off, swingsOff float64
	withSwitch(&paladin.SealOfRighteousnessTriggersWeaponProcs, false, func() {
		off, swingsOff = run()
	})
	if on == 0 || off == 0 || swingsOn == 0 || swingsOff == 0 {
		t.Fatalf("no procs or swings: on %v (%v swings), off %v (%v swings)", on, swingsOn, off, swingsOff)
	}
	ratio := (on / swingsOn) / (off / swingsOff)
	t.Logf("procs a landed swing: %.4f with the seal hits rolling, %.4f without (x%.2f)", on/swingsOn, off/swingsOff, ratio)
	if ratio < 1.7 || ratio > 2.3 {
		t.Errorf("procs a landed swing with the seal hits rolling %.4f, without %.4f: ratio %.2f, want about 2", on/swingsOn, off/swingsOff, ratio)
	}
}
