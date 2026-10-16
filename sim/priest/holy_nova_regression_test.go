package priest

import (
	"fmt"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

func TestHolyNovaHealingCritAllRanks(t *testing.T) {
	for _, ids := range [][2]int32{{15237, 23455}, {15430, 23458}, {15431, 23459}, {27799, 27803}, {27800, 27804}, {27801, 27805}} {
		t.Run(fmt.Sprint(ids[0]), func(t *testing.T) {
			player := &proto.Player{
				Race:          proto.Race_RaceUndead,
				Class:         proto.Class_ClassPriest,
				Equipment:     &proto.EquipmentSpec{},
				Consumables:   &proto.ConsumesSpec{},
				TalentsString: "504025031030001-13505100202-50002",
				Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
				Spec:          &proto.Player_DpsPriest{DpsPriest: &proto.DpsPriest{Options: &proto.DpsPriest_Options{ClassOptions: &proto.PriestOptions{}}}},
			}
			sim := core.NewSim(&proto.RaidSimRequest{
				SimOptions: &proto.SimOptions{RandomSeed: 42},
				Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
				Encounter:  core.MakeSingleTargetEncounter(0),
			}, simsignals.CreateSignals())
			sim.Reset()
			character := sim.Raid.Parties[0].Players[0].GetCharacter()
			damage := character.GetSpell(core.ActionID{SpellID: ids[0]})
			heal := character.GetSpell(core.ActionID{SpellID: ids[1]})
			if damage == nil || heal == nil {
				t.Fatal("Holy Nova damage or healing spell is missing")
			}
			heal.BonusCritPercent = 100
			if !damage.Cast(sim, character.CurrentTarget) {
				t.Fatal("Holy Nova did not cast")
			}
			metrics := heal.SpellMetrics[character.UnitIndex]
			if metrics.Crits != 1 || metrics.TotalHealing <= 0 {
				t.Fatalf("expected one positive healing crit, got %+v", metrics)
			}
			if got := heal.CritDamageMultiplier(nil); got != 1.5 {
				t.Fatalf("healing crit multiplier = %v, want 1.5", got)
			}
		})
	}
}
