package hunter

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
)

// From 16 October 2026 a pet's crit is its master's ranged crit (dev, forever-bugs #331), plus Ferocity on
// the pet; its own Agility gives none.
func TestPetCritIsTheOwnersRangedCrit(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "bm", Class: proto.Class_ClassHunter, Race: proto.Race_RaceOrc, TalentsString: BeastMasteryTalents,
			Equipment: WeaponsOnly, Buffs: &proto.IndividualBuffs{},
			Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{ClassOptions: &proto.HunterOptions{
				Ammo: proto.HunterOptions_Doomshot, PetType: proto.HunterOptions_Cat, PetUptime: 1}}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	hunter := sim.Raid.Parties[0].Players[0].(HunterAgent).GetHunter()
	pet := hunter.Pet
	if hunter.Talents.Ferocity == 0 {
		t.Fatal("the Beast Mastery build no longer takes Ferocity")
	}

	autoShot := hunter.GetStat(stats.PhysicalCritPercent) + hunter.GetStat(stats.RangedCritPercent)
	want := autoShot + spellData.Ferocity.ValueAt(hunter.Talents.Ferocity)
	if got := pet.GetStat(stats.PhysicalCritPercent); !core.WithinToleranceFloat64(want, got, 1e-9) {
		t.Errorf("pet crit %.4f%%, want the Auto Shot's %.4f%% plus Ferocity = %.4f%%", got, autoShot, want)
	}

	pet.AddStatDynamic(sim, stats.Agility, 200)
	if got := pet.GetStat(stats.PhysicalCritPercent); !core.WithinToleranceFloat64(want, got, 1e-9) {
		t.Errorf("200 pet Agility moved its crit to %.4f%%, want %.4f%%", got, want)
	}
}
