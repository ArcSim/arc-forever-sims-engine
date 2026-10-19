package hunter

import (
	"strconv"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// Patch 173: no Hunter ability scales with spell power (client 70425/70450 give Arcane Shot, Serpent
// Sting and Volley's tick spells 0; Bae's in-game tests). With the switch on, +1000 spell damage
// leaves every one of them where it was, to the last digit; with it off, Classic's coefficients
// raise all three.
func TestHunterAbilitiesIgnoreSpellPower(t *testing.T) {
	arcane := spellData.ArcaneShot.Highest().ID
	sting := spellData.SerpentSting.Highest().ID
	volley := spellData.Volley.Highest().ID
	id := func(spell int32) string { return strconv.Itoa(int(spell)) }

	damage := func(on bool, spellDamage float64) map[int32]float64 {
		defer func(old bool) { HunterAbilitiesIgnoreSpellPower = old }(HunterAbilitiesIgnoreSpellPower)
		HunterAbilitiesIgnoreSpellPower = on
		player := &proto.Player{
			Name: "hunter", Class: proto.Class_ClassHunter, Race: proto.Race_RaceOrc, TalentsString: MarksmanshipTalents,
			Equipment: WeaponsOnly, DistanceFromTarget: 30,
			BonusStats: &proto.UnitStats{Stats: stats.Stats{stats.SpellDamage: spellDamage}.ToProtoArray()},
			Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{ClassOptions: &proto.HunterOptions{
				Ammo: proto.HunterOptions_Doomshot, QuiverBonus: proto.HunterOptions_Speed15, PetType: proto.HunterOptions_PetNone}}}},
			Rotation: core.APLRotationFromJsonString(`{"type":"TypeAPL","priorityList":[
				{"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":` + id(sting) + `}}}}},"castSpell":{"spellId":{"spellId":` + id(sting) + `}}}},
				{"action":{"castSpell":{"spellId":{"spellId":` + id(arcane) + `}}}},
				{"action":{"castSpell":{"spellId":{"spellId":` + id(volley) + `}}}}]}`),
		}
		result := core.RunRaidSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
			SimOptions: &proto.SimOptions{Iterations: 20, RandomSeed: 3},
		})
		if result.Error != nil {
			t.Fatal(result.Error.Message)
		}
		byID := map[int32]float64{}
		for _, action := range result.RaidMetrics.Parties[0].Players[0].Actions {
			for _, target := range action.Targets {
				byID[action.Id.GetSpellId()] += target.Damage
			}
		}
		return byID
	}

	names := map[int32]string{arcane: "Arcane Shot", sting: "Serpent Sting", volley: "Volley"}
	base, withSpellPower := damage(true, 0), damage(true, 1000)
	for spell, name := range names {
		if base[spell] == 0 {
			t.Fatalf("%s dealt no damage", name)
		}
		if withSpellPower[spell] != base[spell] {
			t.Errorf("switch on: %s dealt %.2f with +1000 spell damage, want %.2f (none)", name, withSpellPower[spell], base[spell])
		}
	}

	classicBase, classicSpellPower := damage(false, 0), damage(false, 1000)
	for spell, name := range names {
		t.Logf("%s: %.1f, +1000 spell damage %.1f; switch off %.1f, +1000 %.1f", name, base[spell], withSpellPower[spell], classicBase[spell], classicSpellPower[spell])
		if classicSpellPower[spell] <= classicBase[spell] {
			t.Errorf("switch off: %s dealt %.2f with +1000 spell damage, want more than %.2f (Classic's coefficient)", name, classicSpellPower[spell], classicBase[spell])
		}
	}
}
