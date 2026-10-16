package elemental

import (
	"testing"

	"github.com/wowsims/forever/sim/arenalib"
	"github.com/wowsims/forever/sim/common"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

func init() {
	RegisterElementalShaman()
	common.RegisterAllEffects()
}

func TestElemental(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassShaman,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceOrc, proto.Race_RaceDwarf},
			SpecOptions: core.SpecOptionsCombo{Label: "Standard", SpecOptions: &proto.Player_ElementalShaman{
				ElementalShaman: &proto.ElementalShaman{
					Options: &proto.ElementalShaman_Options{
						ClassOptions: &proto.ShamanOptions{
							ShieldProcrate: 0.0,
						},
					},
				},
			}},
			// Naked: the generated item database does not carry the Forever gear our sim tests with
			// yet, and gives the rest TBC-shaped stats. Same call as the merged Mage port.
			GearSet:  core.GearSetCombo{Label: "Naked", GearSet: &proto.EquipmentSpec{}},
			Talents:  DefaultTalents,
			Rotation: core.GetAplRotation("../../../ui/specs/shaman/elemental/apls", "forever"),
			ItemFilter: core.ItemFilter{
				WeaponTypes:       DefaultWeaponTypes,
				ArmorType:         DefaultArmorType,
				RangedWeaponTypes: DefaultRangedWeaponTypes,
			},
		},
	}))
}

// The community build our Forever sim ranks Elemental with.
const DefaultTalents = "5505331000001351--553302001"

const DefaultArmorType = proto.ArmorType_ArmorTypeMail

var DefaultWeaponTypes = []proto.WeaponType{
	proto.WeaponType_WeaponTypeAxe,
	proto.WeaponType_WeaponTypeDagger,
	proto.WeaponType_WeaponTypeFist,
	proto.WeaponType_WeaponTypeMace,
	proto.WeaponType_WeaponTypeStaff,
	proto.WeaponType_WeaponTypeShield,
}

var DefaultRangedWeaponTypes = []proto.RangedWeaponType{
	proto.RangedWeaponType_RangedWeaponTypeTotem,
}

// The arena entry for this spec. Without ARENA_OUT set it only checks every build's damage against the spell manifest; see sim/arenalib.
func TestArena(t *testing.T) {
	arenalib.Run(t, arenalib.Spec{
		Dir:   "elemental_shaman",
		UI:    "shaman/elemental",
		Class: proto.Class_ClassShaman,
		Race:  proto.Race_RaceOrc,
		SpecOptions: &proto.Player_ElementalShaman{ElementalShaman: &proto.ElementalShaman{
			Options: &proto.ElementalShaman_Options{ClassOptions: &proto.ShamanOptions{}},
		}},
		Role: arenalib.Caster,
		// Flame Shock is 20 yd in the client (SpellRange 3).
		DistanceFromTarget: 20,
	})
}
