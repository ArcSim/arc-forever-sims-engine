package hunter

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

// Patch 172: Mortal Shots raises Sniper Shot's crit damage (client 70450's mask, the 16 October
// notes), and with the switch off it does not.
func TestMortalShotsIncludesSniperShot(t *testing.T) {
	critAdditive := func(on bool) (sniper, aimed float64) {
		defer func(old bool) { MortalShotsIncludesSniperShot = old }(MortalShotsIncludesSniperShot)
		MortalShotsIncludesSniperShot = on
		sim := core.NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 1},
			Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
				Name: "mm", Class: proto.Class_ClassHunter, Race: proto.Race_RaceOrc, TalentsString: MarksmanshipTalents,
				Equipment: WeaponsOnly, Buffs: &proto.IndividualBuffs{},
				Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{ClassOptions: &proto.HunterOptions{
					Ammo: proto.HunterOptions_Doomshot, PetType: proto.HunterOptions_PetNone}}}},
				Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
			}}}}},
			Encounter: core.MakeSingleTargetEncounter(0),
		}, simsignals.CreateSignals())
		sim.Reset()

		hunter := sim.Raid.Parties[0].Players[0].(HunterAgent).GetHunter()
		if hunter.Talents.MortalShots == 0 || !hunter.Talents.SniperShot || hunter.SniperShot == nil {
			t.Fatal("the Marksmanship talents need Mortal Shots and Sniper Shot")
		}
		return hunter.SniperShot.CritMultiplierAdditive, hunter.AimedShot.CritMultiplierAdditive
	}

	sniperOn, aimed := critAdditive(true)
	if aimed <= 0 {
		t.Fatalf("Aimed Shot takes no Mortal Shots bonus (%v)", aimed)
	}
	if sniperOn != aimed {
		t.Errorf("switch on: Sniper Shot crit bonus %v, want Aimed Shot's %v", sniperOn, aimed)
	}
	if sniperOff, _ := critAdditive(false); sniperOff != sniperOn-aimed {
		t.Errorf("switch off: Sniper Shot crit bonus %v, want %v", sniperOff, sniperOn-aimed)
	}
}
