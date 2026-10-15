package hunter

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
)

func petRangedCritSim(t *testing.T, race proto.Race, talents string, mainHand int32) (*core.Simulation, *Hunter, *HunterPet) {
	t.Helper()
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1, Iterations: 1},
		Raid: core.SinglePlayerRaidProto(&proto.Player{
			Name: "pet crit", Class: proto.Class_ClassHunter, Race: race, TalentsString: talents,
			Equipment: &proto.EquipmentSpec{Items: []*proto.ItemSpec{
				{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
				{Id: mainHand},
				{},
				{Id: 18713}, // Rhok'delar, Longbow of the Ancient Keepers
			}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
			Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{ClassOptions: &proto.HunterOptions{
				Ammo: proto.HunterOptions_Doomshot, PetType: proto.HunterOptions_Cat, PetUptime: 1}}}},
		}, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	hunter := sim.Raid.Parties[0].Players[0].(HunterAgent).GetHunter()
	if hunter.Pet == nil || !hunter.Pet.IsEnabled() {
		t.Fatal("the hunter has no pet out")
	}
	return sim, hunter, hunter.Pet
}

// autoShotCritPercent is the crit the Hunter's Auto Shot rolls before the target's side of the table.
func autoShotCritPercent(hunter *Hunter) float64 {
	return hunter.GetStat(stats.PhysicalCritPercent) + hunter.GetStat(stats.RangedCritPercent) +
		hunter.AutoAttacks.RangedAuto().BonusCritPercent
}

// Patch 136: a Hunter pet's crit is its master's ranged crit (Forever, 16 October 2026). Its own 5% from
// 100 base Agility is gone, the Auto Shot's crit is what it takes (an Orc's axe racial stays out, a scope
// goes in, TestPetCritFollowsTheOwner), the target's crit suppression applies to both the same, and
// Ferocity still adds on top.
func TestPetCritIsTheOwnersRangedCrit(t *testing.T) {
	const arcaniteReaper = 12784 // two-hand axe
	for _, tc := range []struct {
		name     string
		race     proto.Race
		talents  string
		mainHand int32
	}{
		{"Night Elf Marksmanship", proto.Race_RaceNightElf, MarksmanshipTalents, arcaniteReaper},
		{"Orc with an axe (racial skips Auto Shot)", proto.Race_RaceOrc, MarksmanshipTalents, arcaniteReaper},
		{"Orc Beast Mastery (Ferocity)", proto.Race_RaceOrc, BeastMasteryTalents, arcaniteReaper},
		{"Troll Survival", proto.Race_RaceTroll, SurvivalTalents, arcaniteReaper},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, hunter, pet := petRangedCritSim(t, tc.race, tc.talents, tc.mainHand)
			if tc.talents == BeastMasteryTalents && hunter.Talents.Ferocity == 0 {
				t.Fatal("the Beast Mastery build no longer takes Ferocity")
			}
			ferocity := 0.0
			if hunter.Talents.Ferocity > 0 {
				ferocity = spellData.Ferocity.ValueAt(hunter.Talents.Ferocity)
			}
			want := autoShotCritPercent(hunter) + ferocity
			if got := pet.GetStat(stats.PhysicalCritPercent); !core.WithinToleranceFloat64(want, got, 1e-9) {
				t.Errorf("pet crit %.4f%%, want the Auto Shot's %.4f%% plus Ferocity %.0f%%", got, want-ferocity, ferocity)
			}

			hunterTable := hunter.AttackTables[hunter.CurrentTarget.UnitIndex]
			petTable := pet.AttackTables[pet.CurrentTarget.UnitIndex]
			autoShot := hunter.AutoAttacks.RangedAuto().PhysicalCritChance(hunterTable)
			petSwing := pet.AutoAttacks.MHAuto().PhysicalCritChance(petTable)
			if !core.WithinToleranceFloat64(autoShot+ferocity/100, petSwing, 1e-9) {
				t.Errorf("against the boss the pet's swing crits %.4f, the Auto Shot %.4f (+%.2f Ferocity)", petSwing, autoShot, ferocity/100)
			}
			if tc.race == proto.Race_RaceOrc && hunter.AutoAttacks.RangedAuto().BonusCritPercent != -1 {
				t.Errorf("the Orc's Auto Shot carries %.2f%% bonus crit, want the axe racial's -1%%", hunter.AutoAttacks.RangedAuto().BonusCritPercent)
			}
		})
	}
}

// The pet's Agility no longer gives it crit, and crit the Hunter gains in the fight (Agility, crit, ranged-
// only crit such as a scope's) reaches the pet through the dynamic inheritance (the pet heartbeat calls
// AddOwnerStats with the owner's stat change); a melee-only crit (a sharpening stone's +2% physical, -2%
// ranged) does not.
func TestPetCritFollowsTheOwner(t *testing.T) {
	sim, hunter, pet := petRangedCritSim(t, proto.Race_RaceNightElf, MarksmanshipTalents, 12784)

	before := pet.GetStat(stats.PhysicalCritPercent)
	pet.AddStatDynamic(sim, stats.Agility, 200)
	if got := pet.GetStat(stats.PhysicalCritPercent); got != before {
		t.Errorf("200 pet Agility moved its crit from %.4f%% to %.4f%%, want no change", before, got)
	}

	ownerBefore := hunter.GetStats()
	hunter.AddStatDynamic(sim, stats.Agility, 106)
	hunter.AddStatDynamic(sim, stats.PhysicalCritPercent, 3)
	hunter.AddStatDynamic(sim, stats.RangedCritPercent, 2)
	hunter.AddStatsDynamic(sim, stats.Stats{stats.PhysicalCritPercent: 2, stats.RangedCritPercent: -2})
	pet.AddOwnerStats(sim, hunter.GetStats().Subtract(ownerBefore))
	want := autoShotCritPercent(hunter)
	if got := pet.GetStat(stats.PhysicalCritPercent); !core.WithinToleranceFloat64(want, got, 1e-9) {
		t.Errorf("after the Hunter gained 106 Agility and 3%% crit the pet crits %.4f%%, want %.4f%%", got, want)
	}
	wantGain := 106*core.CritPerAgiMaxLevel[proto.Class_ClassHunter] + 3 + 2
	if gain := pet.GetStat(stats.PhysicalCritPercent) - before; !core.WithinToleranceFloat64(wantGain, gain, 1e-9) {
		t.Errorf("pet crit rose %.4f%%, want %.4f%%", gain, wantGain)
	}
}
