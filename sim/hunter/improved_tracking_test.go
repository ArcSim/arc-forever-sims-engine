package hunter

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

// Improved Tracking (24293, 5% at 5/5) counts on the site's default target, whose mob type is
// Unknown (issue #768: the Survival builds' 5/5 did nothing), and not on Mechanical, which a hunter
// can't track.
func TestImprovedTrackingCountsOnTheDefaultTarget(t *testing.T) {
	multiplier := func(mob proto.MobType) float64 {
		encounter := core.MakeSingleTargetEncounter(0)
		encounter.Targets[0].MobType = mob
		sim := core.NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 1},
			Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
				Name: "sv", Class: proto.Class_ClassHunter, Race: proto.Race_RaceOrc, TalentsString: SurvivalTalents,
				Equipment: WeaponsOnly, Buffs: &proto.IndividualBuffs{},
				Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{ClassOptions: &proto.HunterOptions{
					Ammo: proto.HunterOptions_Doomshot, PetType: proto.HunterOptions_PetNone}}}},
				Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
			}}}}},
			Encounter: encounter,
		}, simsignals.CreateSignals())
		sim.Reset()

		hunter := sim.Raid.Parties[0].Players[0].(HunterAgent).GetHunter()
		if hunter.Talents.ImprovedTracking != 5 {
			t.Fatal("the Survival talents take no 5/5 Improved Tracking")
		}
		return hunter.AttackTables[hunter.CurrentTarget.UnitIndex].DamageDealtMultiplier
	}

	if got := multiplier(proto.MobType_MobTypeUnknown); got != 1.05 {
		t.Errorf("Unknown target: damage multiplier %v, want 1.05", got)
	}
	if got := multiplier(proto.MobType_MobTypeMechanical); got != 1 {
		t.Errorf("Mechanical target: damage multiplier %v, want 1", got)
	}
}
