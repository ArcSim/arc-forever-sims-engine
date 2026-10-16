package priest

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// The priest rotations open with a strict sequence "Inner Focus, then X". Without the talent the
// sequence used to shrink to an unconditional X: smite's Smite outranked Shadow Word: Pain and
// everything below it.
func TestNoInnerFocusDropsTheSequence(t *testing.T) {
	casts := func(apl string, talents string, spellID int32) int32 {
		player := core.WithSpec(&proto.Player{
			Race:          proto.Race_RaceUndead,
			Class:         proto.Class_ClassPriest,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: talents,
			Rotation:      core.GetAplRotation("../../ui/specs/priest/dps/apls", apl).Rotation,
		}, arenaPriestOptions)
		result := core.RunRaidSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 101, Iterations: 1},
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
		})
		if result.Error != nil {
			t.Fatal(result.Error.Message)
		}
		n := int32(0)
		for _, action := range result.RaidMetrics.Parties[0].Players[0].Actions {
			if action.Id.GetSpellId() == spellID {
				for _, target := range action.Targets {
					n += target.Casts
				}
			}
		}
		return n
	}

	if n := casts("smite_lowrank", "-33555133102-5550200003", 10894); n == 0 {
		t.Errorf("smite without Inner Focus never cast Shadow Word: Pain")
	}
	// Shadow keeps Devouring Plague on its own row for builds without Inner Focus.
	if n := casts("shadow", "02530103002--55042201201302251", 19280); n == 0 {
		t.Errorf("shadow without Inner Focus never cast Devouring Plague")
	}
	// Forever's Devouring Plague is a 1 min cooldown (19280 category 691), Inner Focus 3 min: the
	// rank 6 plague must not wait for Inner Focus.
	if n := casts("shadow", ShadowTalents, 19280); n < 2 {
		t.Errorf("shadow with Inner Focus cast Devouring Plague %d times in %d s, want it again on its 1 min cooldown", n, core.LongDuration)
	}
}
