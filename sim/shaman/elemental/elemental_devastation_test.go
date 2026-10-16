package elemental

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// Fire Nova crits don't trigger Elemental Devastation: no totem or totem attack does (dev AidanZMoon,
// forever-bugs #334). Earth Shock crits still do.
func TestElementalDevastationSkipsFireNova(t *testing.T) {
	uptime := func(spellID int32) float64 {
		player := &proto.Player{
			Name: "ele", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc, TalentsString: "003",
			Equipment: &proto.EquipmentSpec{},
			Spec:      &proto.Player_ElementalShaman{ElementalShaman: &proto.ElementalShaman{Options: &proto.ElementalShaman_Options{ClassOptions: &proto.ShamanOptions{}}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL, PriorityList: []*proto.APLListItem{{Action: &proto.APLAction{Action: &proto.APLAction_CastSpell{
				CastSpell: &proto.APLActionCastSpell{SpellId: &proto.ActionID{RawId: &proto.ActionID_SpellId{SpellId: spellID}}},
			}}}}},
		}
		res := core.RunRaidSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
			SimOptions: &proto.SimOptions{Iterations: 20, RandomSeed: 1},
		})
		if res.Error != nil {
			t.Fatal(res.Error.Message)
		}
		for _, a := range res.RaidMetrics.Parties[0].Players[0].Auras {
			if a.Id.GetSpellId() == 29178 {
				return a.UptimeSecondsAvg
			}
		}
		return 0
	}

	if got := uptime(408345); got != 0 {
		t.Errorf("Fire Nova only: Elemental Devastation up %.1f s, want 0", got)
	}
	if got := uptime(10414); got == 0 {
		t.Errorf("Earth Shock only: Elemental Devastation never up")
	}
}
