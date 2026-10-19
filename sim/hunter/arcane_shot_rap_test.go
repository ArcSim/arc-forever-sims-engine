package hunter

import (
	"math"
	"strconv"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// Patch 174: Arcane Shot takes 0.10 of ranged attack power (the Hunter Discord's three tests, 0.10,
// 0.10 and 0.11, with the logs' 0.10-0.12); off, the logs-only fit of 0.11. Attack power does not
// change a roll, so with the same seed the damage +1000 ranged attack power adds is the share times
// the same hit and crit multipliers: the switch-on gain is 10/11 of the switch-off gain.
func TestArcaneShotRAPFromHunterTests(t *testing.T) {
	arcane := spellData.ArcaneShot.Highest().ID

	damage := func(on bool, rap float64) float64 {
		defer func(old bool) { ArcaneShotRAPFromHunterTests = old }(ArcaneShotRAPFromHunterTests)
		ArcaneShotRAPFromHunterTests = on
		player := &proto.Player{
			Name: "hunter", Class: proto.Class_ClassHunter, Race: proto.Race_RaceOrc, TalentsString: MarksmanshipTalents,
			Equipment: WeaponsOnly, DistanceFromTarget: 30,
			BonusStats: &proto.UnitStats{Stats: stats.Stats{stats.RangedAttackPower: rap}.ToProtoArray()},
			Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{ClassOptions: &proto.HunterOptions{
				Ammo: proto.HunterOptions_Doomshot, QuiverBonus: proto.HunterOptions_Speed15, PetType: proto.HunterOptions_PetNone}}}},
			Rotation: core.APLRotationFromJsonString(`{"type":"TypeAPL","priorityList":[
				{"action":{"castSpell":{"spellId":{"spellId":` + strconv.Itoa(int(arcane)) + `}}}}]}`),
		}
		result := core.RunRaidSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
			SimOptions: &proto.SimOptions{Iterations: 20, RandomSeed: 3},
		})
		if result.Error != nil {
			t.Fatal(result.Error.Message)
		}
		total := 0.0
		for _, action := range result.RaidMetrics.Parties[0].Players[0].Actions {
			if action.Id.GetSpellId() != arcane {
				continue
			}
			for _, target := range action.Targets {
				total += target.Damage
			}
		}
		if total == 0 {
			t.Fatal("Arcane Shot dealt no damage")
		}
		return total
	}

	onGain := damage(true, 1000) - damage(true, 0)
	offGain := damage(false, 1000) - damage(false, 0)
	t.Logf("+1000 ranged attack power adds %.1f Arcane Shot damage with the switch on, %.1f off", onGain, offGain)
	if offGain <= 0 {
		t.Fatalf("switch off: +1000 ranged attack power added %.2f Arcane Shot damage, want more", offGain)
	}
	if ratio := onGain / offGain; math.Abs(ratio-10.0/11.0) > 1e-9 {
		t.Errorf("switch on gains %.6f of the switch-off gain, want %.6f (0.10 against 0.11)", ratio, 10.0/11.0)
	}
}
