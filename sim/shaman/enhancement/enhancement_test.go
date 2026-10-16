package enhancement

import (
	"testing"

	"github.com/wowsims/forever/sim/arenalib"
	"github.com/wowsims/forever/sim/common"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
	"github.com/wowsims/forever/sim/shaman"
)

func init() {
	RegisterEnhancementShaman()
	common.RegisterAllEffects()
}

func TestEnhancement(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassShaman,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceOrc, proto.Race_RaceDwarf},
			SpecOptions: core.SpecOptionsCombo{Label: "Standard", SpecOptions: &proto.Player_EnhancementShaman{
				EnhancementShaman: &proto.EnhancementShaman{
					Options: &proto.EnhancementShaman_Options{
						SyncType:    proto.ShamanSyncType_Auto,
						ImbueOh:     proto.ShamanImbue_WindfuryWeapon,
						ImbueOhSwap: proto.ShamanImbue_WindfuryWeapon,
						ClassOptions: &proto.ShamanOptions{
							ImbueMh:        proto.ShamanImbue_WindfuryWeapon,
							ImbueMhSwap:    proto.ShamanImbue_WindfuryWeapon,
							ShieldProcrate: 0.0,
						},
					},
				},
			}},
			// Arcanite Reaper alone, as Retribution: the generated item database does not carry the
			// Forever gear our sim tests with, and a naked shaman never casts Stormstrike or procs
			// Windfury, Flurry or Elemental Weapons.
			GearSet:  core.GearSetCombo{Label: "Weapon", GearSet: WeaponOnly},
			Talents:  DefaultTalents,
			Rotation: core.GetAplRotation("../../../ui/specs/shaman/enhancement/apls", "forever"),
			ItemFilter: core.ItemFilter{
				WeaponTypes: []proto.WeaponType{
					proto.WeaponType_WeaponTypeAxe,
					proto.WeaponType_WeaponTypeDagger,
					proto.WeaponType_WeaponTypeFist,
					proto.WeaponType_WeaponTypeMace,
					proto.WeaponType_WeaponTypeOffHand,
				},
				ArmorType: proto.ArmorType_ArmorTypeMail,
				RangedWeaponTypes: []proto.RangedWeaponType{
					proto.RangedWeaponType_RangedWeaponTypeTotem,
				},
			},
		},
	}))
}

var WeaponOnly = &proto.EquipmentSpec{
	Items: []*proto.ItemSpec{
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{Id: 12784}, // Arcanite Reaper
	},
}

// The community build our Forever sim ranks Enhancement with.
const DefaultTalents = "5505301-053030031005112251"

// Stormstrike (17364, aura 271) raises only its caster's Lightning Bolt, Chain Lightning and Earth
// Shock. Only the bolts spend the charge: Earth Shock takes the bonus and leaves it (patch 140, beta
// testing relayed by Redfall, 16 October 2026).
func TestStormstrikeOnlyBoostsCastersBoltsAndEarthShock(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Shaman", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc, TalentsString: DefaultTalents,
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{},
			Spec: &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{
				Options: &proto.EnhancementShaman_Options{ClassOptions: &proto.ShamanOptions{}},
			}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	sham := sim.Raid.Parties[0].Players[0].(shaman.ShamanAgent).GetShaman()
	target := sham.CurrentTarget
	table := sham.AttackTables[target.UnitIndex]
	debuff := sham.StormStrikeDebuffAuras.Get(target)
	earthShockBase := sham.EarthShock.TargetDamageMultiplier(sim, table, false)

	debuff.Activate(sim)
	debuff.SetStacks(sim, debuff.MaxStacks)
	if got := sham.EarthShock.TargetDamageMultiplier(sim, table, false) / earthShockBase; got != 1.2 {
		t.Errorf("Earth Shock multiplier under Stormstrike = %v, want 1.2", got)
	}
	if got := target.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature]; got != 1 {
		t.Errorf("target takes %v Nature damage from everyone, want 1", got)
	}

	sham.FrostShock.CalcAndDealDamage(sim, target, 100, sham.FrostShock.OutcomeAlwaysHit)
	if !debuff.IsActive() {
		t.Fatal("Frost Shock spent the Stormstrike charge")
	}
	sham.EarthShock.CalcAndDealDamage(sim, target, 100, sham.EarthShock.OutcomeAlwaysHit)
	if !debuff.IsActive() {
		t.Fatal("Earth Shock spent the Stormstrike charge")
	}
	if got := sham.EarthShock.TargetDamageMultiplier(sim, table, false) / earthShockBase; got != 1.2 {
		t.Errorf("second Earth Shock multiplier under Stormstrike = %v, want 1.2", got)
	}
	lightningBolt := sham.GetSpell(core.ActionID{SpellID: 15208})
	if lightningBolt == nil {
		t.Fatal("no Lightning Bolt rank 10 (15208)")
	}
	lightningBolt.CalcAndDealDamage(sim, target, 100, lightningBolt.OutcomeAlwaysHit)
	if debuff.IsActive() {
		t.Error("Lightning Bolt did not spend the Stormstrike charge")
	}
}

// The arena entry for this spec. Without ARENA_OUT set it only checks every build's damage against the spell manifest; see sim/arenalib.
func TestArena(t *testing.T) {
	arenalib.Run(t, arenalib.Spec{
		Dir:   "enhancement_shaman",
		UI:    "shaman/enhancement",
		Class: proto.Class_ClassShaman,
		Race:  proto.Race_RaceDwarf,
		SpecOptions: &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{
			Options: &proto.EnhancementShaman_Options{
				SyncType:     proto.ShamanSyncType_Auto,
				ImbueOh:      proto.ShamanImbue_WindfuryWeapon,
				ClassOptions: &proto.ShamanOptions{ImbueMh: proto.ShamanImbue_WindfuryWeapon},
			},
		}},
		Role: arenalib.Melee,
		// Windfury Weapon is the shaman casting on their own weapons, not a totem somebody drops.
		ClassImbues: arenalib.ClassImbues{Windfury: true},
	})
}
