package hunter

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

func petSim(t *testing.T, pet proto.HunterOptions_PetType, talents, apl string, distance float64) *proto.RaidSimResult {
	t.Helper()
	player := &proto.Player{
		Name: "h", Class: proto.Class_ClassHunter, Race: proto.Race_RaceOrc, TalentsString: talents,
		Equipment: WeaponsOnly, DistanceFromTarget: distance,
		Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{ClassOptions: &proto.HunterOptions{
			Ammo: proto.HunterOptions_Doomshot, QuiverBonus: proto.HunterOptions_Speed15, PetType: pet,
			PetAttackSpeed: proto.HunterOptions_OneTwo, PetUptime: 1}}}},
		Rotation: core.GetAplRotation("../../ui/specs/hunter/dps/apls", apl).Rotation,
	}
	raid := &proto.Raid{Parties: []*proto.Party{{Players: []*proto.Player{player}, Buffs: &proto.PartyBuffs{}}}, Buffs: &proto.RaidBuffs{}, Debuffs: &proto.Debuffs{}, NumActiveParties: 1}
	res := core.RunRaidSim(&proto.RaidSimRequest{Raid: raid, Encounter: core.MakeSingleTargetEncounter(0), SimOptions: &proto.SimOptions{Iterations: 20, RandomSeed: 1}})
	if res.Error != nil {
		t.Fatal(res.Error.Message)
	}
	return res
}

func petCasts(res *proto.RaidSimResult, id int32) float64 {
	var casts int32
	for _, pet := range res.RaidMetrics.Parties[0].Players[0].Pets {
		for _, action := range pet.Actions {
			if action.Id.GetSpellId() == id {
				for _, target := range action.Targets {
					casts += target.Casts
				}
			}
		}
	}
	return float64(casts) / 20
}

// Patch 118: with PetPoolsFocusForCooldowns the Cat holds Claw while it would leave Bite (35 focus, 10 sec
// cooldown) short, so Bite goes on cooldown. Without it Claw (25) spends every 25 focus and Bite lands once,
// off the full bar the pet starts with (Kragtor's sim: one Bite in 240 s).
func TestCatBitesOnCooldown(t *testing.T) {
	defer func(old bool) { PetPoolsFocusForCooldowns = old }(PetPoolsFocusForCooldowns)
	fight := core.LongDuration

	PetPoolsFocusForCooldowns = false
	unpooled := petCasts(petSim(t, proto.HunterOptions_Cat, MarksmanshipTalents, "mm", 30), 17261)
	PetPoolsFocusForCooldowns = true
	pooled := petSim(t, proto.HunterOptions_Cat, MarksmanshipTalents, "mm", 30)
	bites, claws := petCasts(pooled, 17261), petCasts(pooled, 3009)

	t.Logf("Bites a fight: %.1f unpooled, %.1f pooled (%d s fight, Claws %.1f)", unpooled, bites, fight, claws)
	if unpooled > 2 {
		t.Errorf("unpooled Cat Bit %.1f times, want the opener only", unpooled)
	}
	if want := 0.9 * float64(fight) / 10; bites < want {
		t.Errorf("pooled Cat Bit %.1f times in %d s, want at least %.0f (every 10 s)", bites, fight, want)
	}
	// Focus still all goes somewhere: 10 a second, 35 a Bite, 25 a Claw.
	if spent := 35*bites + 25*claws - 100; spent < 9.3*float64(fight) {
		t.Errorf("pooled Cat spent %.0f focus in %d s, want about 10 a second", spent, fight)
	}
}
