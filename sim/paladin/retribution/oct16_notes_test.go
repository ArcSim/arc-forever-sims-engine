package retribution

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
	"github.com/wowsims/forever/sim/paladin"
)

// Champion of the Light gives 10/20/30% of Intellect as spell damage (16 October 2026 beta notes; was
// 20/40/60%).
func TestChampionOfTheLightShare(t *testing.T) {
	player := core.WithSpec(&proto.Player{
		Race: proto.Race_RaceHuman, Class: proto.Class_ClassPaladin, Equipment: WeaponOnly,
		TalentsString: RetDeepTalents,
		Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, DefaultOptions)
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	pal := sim.Raid.Parties[0].Players[0].(paladin.PaladinAgent).GetPaladin()
	rank := pal.Talents.ChampionOfTheLight
	if rank == 0 {
		t.Fatal("test build has no Champion of the Light")
	}
	sheet := pal.GetStats()
	if got, want := sheet[stats.SpellDamage], math.Floor(0.1*float64(rank)*sheet[stats.Intellect]); math.Abs(got-want) > 1 {
		t.Fatalf("Champion of the Light %d gives %v spell damage on %v Intellect, want %v", rank, got, sheet[stats.Intellect], want)
	}
}
