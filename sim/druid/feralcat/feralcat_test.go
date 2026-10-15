package feralcat

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/arenalib"
	"github.com/wowsims/forever/sim/common"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/druid"
)

func init() {
	RegisterFeralCatDruid()
	common.RegisterAllEffects()
}

func TestFeralCat(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassDruid,
			Race:       proto.Race_RaceNightElf,
			OtherRaces: []proto.Race{proto.Race_RaceTauren},

			// Naked: the generated item database does not carry the Forever gear our sim tests
			// with yet, and gives the rest TBC-shaped stats.
			GearSet: core.GearSetCombo{Label: "Naked", GearSet: &proto.EquipmentSpec{}},

			Talents: DefaultTalents,
			OtherTalentSets: []core.TalentsCombo{
				{Label: "FeralCat", Talents: FeralCatTalents},
			},

			SpecOptions: core.SpecOptionsCombo{Label: "Standard", SpecOptions: DefaultSpecOptions},

			Rotation: core.GetAplRotation("../../../ui/specs/druid/feralcat/apls", "default"),

			Consumables: DefaultConsumables,

			Profession1: proto.Profession_Engineering,
			Profession2: proto.Profession_Enchanting,

			ItemFilter: core.ItemFilter{
				ArmorType: proto.ArmorType_ArmorTypeLeather,
				WeaponTypes: []proto.WeaponType{
					proto.WeaponType_WeaponTypeDagger,
					proto.WeaponType_WeaponTypeFist,
					proto.WeaponType_WeaponTypeMace,
					proto.WeaponType_WeaponTypeStaff,
				},
				RangedWeaponTypes: []proto.RangedWeaponType{
					proto.RangedWeaponType_RangedWeaponTypeIdol,
				},
			},

			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh: []proto.Stat{
				proto.Stat_StatAgility,
				proto.Stat_StatStrength,
				proto.Stat_StatAttackPower,
				proto.Stat_StatFeralAttackPower,
				proto.Stat_StatMeleeHitRating,
				proto.Stat_StatExpertiseRating,
				proto.Stat_StatMeleeCritRating,
				proto.Stat_StatMeleeHasteRating,
				proto.Stat_StatArmorPenetration,
			},
		},
	}))
}

// Our Forever sim's feral builds.
const DefaultTalents = "-55210032021132212051-05503"
const FeralCatTalents = "050022-55000032121032212051-052"

var DefaultSpecOptions = &proto.Player_FeralCatDruid{
	FeralCatDruid: &proto.FeralCatDruid{
		Rotation: &proto.FeralCatDruid_Rotation{
			FinishingMove:      proto.FeralCatDruid_Rotation_Rip,
			Biteweave:          true,
			RipMinComboPoints:  5,
			BiteMinComboPoints: 5,
			MangleTrick:        true,
			MaintainFaerieFire: false,
		},
		Options: &proto.FeralCatDruid_Options{},
	},
}

var DefaultConsumables = &proto.ConsumesSpec{
	ConjuredId:   12662, // Demonic Rune
	GoblinSapper: true,
}

// Clearcasting (16870, one charge) makes the next ability in its mask free and is spent by it; one
// outside the mask leaves it up.
func TestClearcastingSpentByNextCostedAbility(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Cat", Class: proto.Class_ClassDruid, Race: proto.Race_RaceNightElf, TalentsString: DefaultTalents,
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{}, Spec: DefaultSpecOptions,
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	cat := sim.Raid.Parties[0].Players[0].(druid.DruidAgent).GetDruid()
	cat.ClearcastingAura.Activate(sim)

	if !cat.TigersFury.Cast(sim, cat.CurrentTarget) || !cat.ClearcastingAura.IsActive() {
		t.Fatal("Tiger's Fury, outside the mask, did not cast or spent Clearcasting")
	}

	for !cat.GCD.IsReady(sim) && sim.CurrentTime < 5*time.Second {
		sim.Step()
	}
	energy := cat.CurrentEnergy()
	if !cat.Shred.Cast(sim, cat.CurrentTarget) {
		t.Fatal("Shred did not cast")
	}
	if cat.CurrentEnergy() != energy {
		t.Errorf("Shred cost %v Energy under Clearcasting, want 0", energy-cat.CurrentEnergy())
	}
	if cat.ClearcastingAura.IsActive() {
		t.Error("Shred did not spend Clearcasting")
	}
}

// Omen of Clarity rolls its 2 procs a minute off the paw's 1.0 s swing in Cat Form, for specials too, not
// off the equipped weapon: beta logs gave 36 procs off 542 Maul, Swipe and Claw hits where the weapon's
// speed predicts 59 (sim/druid/omen_of_clarity.go).
func TestOmenOfClarityIgnoresWeaponSpeed(t *testing.T) {
	items := make([]*proto.ItemSpec, proto.ItemSlot_ItemSlotRanged+1)
	for i := range items {
		items[i] = &proto.ItemSpec{}
	}
	items[proto.ItemSlot_ItemSlotMainHand] = &proto.ItemSpec{Id: 7230} // Smite's Mighty Hammer, 3.5 s

	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Cat", Class: proto.Class_ClassDruid, Race: proto.Race_RaceNightElf, TalentsString: DefaultTalents,
			Equipment: &proto.EquipmentSpec{Items: items}, Buffs: &proto.IndividualBuffs{}, Spec: DefaultSpecOptions,
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	cat := sim.Raid.Parties[0].Players[0].(druid.DruidAgent).GetDruid()
	if weapon := cat.GetMHWeapon(); weapon == nil || weapon.SwingSpeed != 3.5 {
		t.Fatalf("want the 3.5 s Smite's Mighty Hammer in the main hand, got %v", weapon)
	}
	omen := cat.GetAura("Omen of Clarity")
	if omen == nil {
		t.Fatal("Omen of Clarity is not registered")
	}

	const casts = 4000
	procs := 0
	for range casts {
		cat.ClearcastingAura.Deactivate(sim)
		omen.Icd.Reset()
		cat.GCD.Reset()
		cat.AddEnergy(sim, 100, cat.EnergyRefundMetrics)
		if !cat.Shred.Cast(sim, cat.CurrentTarget) {
			t.Fatal("Shred did not cast")
		}
		// Procs land one spell batch window after the hit.
		settled := false
		sim.AddPendingAction(core.NewDelayedAction(core.DelayedActionOptions{
			DoAt:     sim.CurrentTime + core.SpellBatchWindow,
			Priority: core.ActionPriorityLow,
			OnAction: func(*core.Simulation) { settled = true },
		}))
		for !settled {
			sim.Step()
		}
		if cat.ClearcastingAura.IsActive() {
			procs++
		}
	}

	// 2 a minute is 3.3% a landed hit on the paw's 1.0 s and 11.7% on the hammer's 3.5 s.
	if rate := float64(procs) / casts; rate < 0.02 || rate > 0.05 {
		t.Errorf("Shred procced Clearcasting %.1f%% of casts, want about 3%%", rate*100)
	}
}

// The arena entry for this spec. Without ARENA_OUT set it only checks every build's damage against the spell manifest; see sim/arenalib.
func TestArena(t *testing.T) {
	arenalib.Run(t, arenalib.Spec{
		Dir:         "feral_druid",
		UI:          "druid/feralcat",
		Class:       proto.Class_ClassDruid,
		Race:        proto.Race_RaceTauren,
		SpecOptions: DefaultSpecOptions,
		Role:        arenalib.Melee,
	})
}

// The default rotation Prowls before the pull and opens with Ravage (9867), once a fight: Prowl
// only casts before combat and the opener breaks it.
func TestRavageOpensFromProwl(t *testing.T) {
	result := core.RunRaidSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1, Iterations: 100},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Cat", Class: proto.Class_ClassDruid, Race: proto.Race_RaceNightElf, TalentsString: FeralCatTalents,
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{}, Spec: DefaultSpecOptions,
			Rotation: core.GetAplRotation("../../../ui/specs/druid/feralcat/apls", "default").Rotation,
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	})
	if result.Error != nil {
		t.Fatal(result.Error.Message)
	}

	casts := int32(0)
	for _, action := range result.RaidMetrics.Parties[0].Players[0].Actions {
		if action.Id.GetSpellId() == 9867 {
			for _, target := range action.Targets {
				casts += target.Casts
			}
		}
	}
	if casts != 100 {
		t.Errorf("Ravage cast %d times over 100 fights, want 100", casts)
	}
}

// Claw (9850) is the builder that works from the front, where Shred cannot: 45 Energy less Ferocity's
// one a rank, one combo point when it lands.
func TestClawFromTheFront(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Cat", Class: proto.Class_ClassDruid, Race: proto.Race_RaceNightElf, TalentsString: DefaultTalents,
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{}, Spec: DefaultSpecOptions,
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL}, InFrontOfTarget: true,
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	cat := sim.Raid.Parties[0].Players[0].(druid.DruidAgent).GetDruid()
	if cat.Shred.Cast(sim, cat.CurrentTarget) {
		t.Fatal("Shred cast from the front")
	}

	landed := 0
	for range 20 {
		cat.GCD.Reset()
		cat.AddEnergy(sim, 100, cat.EnergyRefundMetrics)
		energy, cp := cat.CurrentEnergy(), cat.ComboPoints()
		if !cat.Claw.Cast(sim, cat.CurrentTarget) {
			t.Fatal("Claw did not cast from the front")
		}
		if cat.ComboPoints() > cp {
			landed++
			if spent, want := energy-cat.CurrentEnergy(), float64(45-cat.Talents.Ferocity); spent != want && !cat.ClearcastingAura.IsActive() {
				t.Errorf("Claw cost %v Energy, want %v", spent, want)
			}
		}
		cat.SpendComboPoints(sim, cat.Rip.ComboPointMetrics())
	}
	if landed < 10 {
		t.Errorf("Claw awarded a combo point %d times in 20 casts", landed)
	}
}

// Ravage (9867) carries Attributes[0] 0x200000 on every rank: it can miss but is never dodged, parried
// or blocked. Behind the target already rules out parry and block, so dodges are what this checks.
func TestRavageCannotBeDodged(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Cat", Class: proto.Class_ClassDruid, Race: proto.Race_RaceNightElf, TalentsString: DefaultTalents,
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{}, Spec: DefaultSpecOptions,
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	cat := sim.Raid.Parties[0].Players[0].(druid.DruidAgent).GetDruid()
	target := cat.CurrentTarget
	for range 2000 {
		cat.Ravage.SkipCastAndApplyEffects(sim, target)
	}
	if m := cat.Ravage.SpellMetrics[target.UnitIndex]; m.Dodges != 0 || m.Misses == 0 {
		t.Errorf("Ravage over 2000 attempts: %d dodges, %d misses, want 0 dodges and some misses", m.Dodges, m.Misses)
	}
}

// Shifting Power (1322605) can't be used at full Energy since the 8 Oct beta build (dev notes:
// "can no longer be used while already at Full Energy").
func TestShiftingPowerRefusedAtFullEnergy(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Cat", Class: proto.Class_ClassDruid, Race: proto.Race_RaceNightElf, TalentsString: DefaultTalents,
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{}, Spec: DefaultSpecOptions,
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	cat := sim.Raid.Parties[0].Players[0].(druid.DruidAgent).GetDruid()
	if cat.ShiftingPower == nil {
		t.Fatal("the default build has no Shifting Power")
	}
	cat.AddEnergy(sim, 100, cat.EnergyRefundMetrics)
	if cat.ShiftingPower.Cast(sim, cat.CurrentTarget) {
		t.Fatal("Shifting Power cast at full Energy")
	}
	cat.SpendEnergy(sim, 50, cat.EnergyRefundMetrics)
	if !cat.ShiftingPower.Cast(sim, cat.CurrentTarget) {
		t.Fatal("Shifting Power did not cast at half Energy")
	}
}
