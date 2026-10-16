package healer

import (
	"testing"

	"github.com/wowsims/forever/sim/common"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

func init() {
	RegisterHealerPriest()
	common.RegisterAllEffects()
}

// Stats-only suite: this spec is a gear planner, it has no healing rotation.
// Pins the final character stats for each gear preset so the passives stay covered. The empty APL
// rotation and the fake prepull (no SkipRotation) make it exercise a full environment reset, the
// path the UI's stats request takes.
func TestHealerPriest(t *testing.T) {
	var generators []core.TestGenerator
	// Naked and one talent build per preset: the generated item database does not carry the Forever
	// gear our sim plans with, and gives the rest TBC-shaped stats.
	for _, build := range []struct{ name, talents string }{{"holy-19-32-0", HolyTalents}, {"discipline-35-16-0", DisciplineTalents}} {
		player := core.WithSpec(
			&proto.Player{
				Class:         proto.Class_ClassPriest,
				Race:          proto.Race_RaceDwarf,
				Equipment:     &proto.EquipmentSpec{},
				Consumables:   FullConsumes,
				Buffs:         core.FullIndividualBuffs,
				TalentsString: build.talents,
				Profession1:   proto.Profession_Tailoring,
				Profession2:   proto.Profession_Enchanting,
				Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
			},
			PlayerOptions,
		)
		generators = append(generators, &core.SingleCharacterStatsTestGenerator{
			Name: build.name,
			Request: &proto.ComputeStatsRequest{
				Raid: core.SinglePlayerRaidProto(player, core.FullPartyBuffs, core.FullRaidBuffs, core.FullDebuffs),
			},
		})
	}
	core.RunTestSuite(t, t.Name(), generators)
}

// ui/specs/priest/healer/presets.ts.
var HolyTalents = "00520303123-2350510323000053"
var DisciplineTalents = "005205031332011531-03505003"

// The UI sets no consumables for this spec (the TBC flask, food and potion are not in the Forever client).
var FullConsumes = &proto.ConsumesSpec{}

var PlayerOptions = &proto.Player_HealerPriest{
	HealerPriest: &proto.HealerPriest{
		Options: &proto.HealerPriest_Options{
			ClassOptions: &proto.PriestOptions{},
		},
	},
}
