package priest

import (
	"testing"

	"github.com/wowsims/forever/sim/arenalib"
	"github.com/wowsims/forever/sim/common"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

func init() {
	RegisterPriest()
	common.RegisterAllEffects()
}

// The community builds the rankings page runs on our Forever sim: Shadow 15/0/36 and Smite 31/17/3.
var ShadowTalents = "02530003101--55052201201302251"
var SmiteTalents = "515035031030001031-00505023002-003"

func TestShadowPriest(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		priestSuite("shadow", ShadowTalents, true),
	}))
}

func TestSmitePriest(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		priestSuite("smite", SmiteTalents, false),
	}))
}

func priestSuite(apl string, talents string, preShadowform bool) core.CharacterSuiteConfig {
	return core.CharacterSuiteConfig{
		Class:      proto.Class_ClassPriest,
		Race:       proto.Race_RaceTroll,
		OtherRaces: []proto.Race{proto.Race_RaceUndead, proto.Race_RaceNightElf},

		SpecOptions: core.SpecOptionsCombo{
			Label: "Default",
			SpecOptions: &proto.Player_DpsPriest{
				DpsPriest: &proto.DpsPriest{
					Options: &proto.DpsPriest_Options{
						ClassOptions: &proto.PriestOptions{
							Armor: proto.PriestOptions_InnerFire,
							// Off by default in the UI; on here so the pet stays covered.
							UseShadowfiend: true,
							// Begin the sim already in Shadowform so the opener does not spend a
							// GCD casting it.
							PreShadowform: preShadowform,
						},
					},
				},
			},
		},

		// Naked: the generated item database does not carry the pre-raid set our Forever sim tests
		// with yet, and gives the rest TBC-shaped stats.
		GearSet:  core.GearSetCombo{Label: "Naked", GearSet: &proto.EquipmentSpec{}},
		Talents:  talents,
		Rotation: core.GetAplRotation("../../ui/specs/priest/dps/apls", apl),

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeStaff,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeOffHand,
			},
			ArmorType: proto.ArmorType_ArmorTypeCloth,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeWand,
			},
			// Blacklist melee enchants that appear on cloth-relevant slots but are never used by
			// casters.
			EnchantBlacklist: []int32{2673, 3225, 3273},
		},
	}
}

// Both priests share ui/specs/priest/dps, so each takes its own talents, gear and rotation out
// of it. The site's options, which leave Inner Fire and Shadowfiend off.
func TestArenaShadow(t *testing.T) {
	arenalib.Run(t, arenalib.Spec{
		Dir:         "shadow_priest",
		UI:          "priest/dps",
		Class:       proto.Class_ClassPriest,
		Race:        proto.Race_RaceUndead,
		SpecOptions: arenaPriestOptions,
		Role:        arenalib.Caster,
		// Mind Flay is 20 yd in the client (SpellRange 3).
		DistanceFromTarget: 20,
		Talents:            "Shadow",
		GearSets:           []string{"launch", "p0.bis", "p1.bis"},
		Rotations:          []string{"shadow"},
	})
}

func TestArenaSmite(t *testing.T) {
	arenalib.Run(t, arenalib.Spec{
		Dir:                "smite_priest",
		UI:                 "priest/dps",
		Class:              proto.Class_ClassPriest,
		Race:               proto.Race_RaceUndead,
		SpecOptions:        arenaPriestOptions,
		Role:               arenalib.Caster,
		DistanceFromTarget: 30,
		Talents:            "Smite",
		GearSets:           []string{"smite_launch"},
		Rotations:          []string{"smite", "smite_lowrank", "smite_lowrank_mindblast"},
	})
}

var arenaPriestOptions = &proto.Player_DpsPriest{DpsPriest: &proto.DpsPriest{Options: &proto.DpsPriest_Options{
	ClassOptions: &proto.PriestOptions{},
}}}
