package enhancement

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// The shipped rotation casts Lightning Bolt at 5 Maelstrom Weapon stacks. A build without the talent
// has no such aura; the condition used to drop out and leave Lightning Bolt cast on every GCD.
func TestNoMaelstromWeaponNoLightningBolt(t *testing.T) {
	stacks := &proto.APLValue{Value: &proto.APLValue_Cmp{Cmp: &proto.APLValueCompare{
		Op:  proto.APLValueCompare_OpGe,
		Lhs: &proto.APLValue{Value: &proto.APLValue_AuraNumStacks{AuraNumStacks: &proto.APLValueAuraNumStacks{AuraId: core.ActionID{SpellID: 408505}.ToProto()}}},
		Rhs: &proto.APLValue{Value: &proto.APLValue_Const{Const: &proto.APLValueConst{Val: "5"}}},
	}}}
	player := &proto.Player{
		Name: "enh", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc,
		TalentsString: "05323303001-0530300310051-05005", // 21 Enhancement on the UI presets: no Maelstrom Weapon
		Equipment:     &proto.EquipmentSpec{Items: []*proto.ItemSpec{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {Id: 17182}}},
		Spec:          &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{ClassOptions: &proto.ShamanOptions{}}}},
		Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL, PriorityList: []*proto.APLListItem{{Action: &proto.APLAction{
			Condition: stacks,
			Action:    &proto.APLAction_CastSpell{CastSpell: &proto.APLActionCastSpell{SpellId: core.ActionID{SpellID: 15208}.ToProto()}},
		}}}},
	}
	raid := &proto.Raid{Parties: []*proto.Party{{Players: []*proto.Player{player}, Buffs: &proto.PartyBuffs{}}}, Buffs: &proto.RaidBuffs{}, Debuffs: &proto.Debuffs{}, NumActiveParties: 1}
	res := core.RunRaidSim(&proto.RaidSimRequest{Raid: raid, Encounter: core.MakeSingleTargetEncounter(0), SimOptions: &proto.SimOptions{Iterations: 5, RandomSeed: 1}})
	if res.Error != nil {
		t.Fatal(res.Error.Message)
	}
	for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
		for _, tgt := range a.Targets {
			if a.Id.GetSpellId() == 15208 && tgt.Casts > 0 {
				t.Fatalf("Lightning Bolt cast %d times without Maelstrom Weapon", tgt.Casts)
			}
		}
	}
}

// Maelstrom Weapon procs at 10 PPM whatever the talent rank: Steele (foreverlogs 2748, 1 point)
// stacked it on 0.54 of landed hits with a 3.6 sec two-hander. 2 PPM a point gave 1 point a fifth.
func TestMaelstromWeaponRateDoesNotScaleWithPoints(t *testing.T) {
	procs := func(talents string) float64 {
		player := &proto.Player{
			Name: "enh", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc,
			TalentsString: talents,
			Equipment:     &proto.EquipmentSpec{Items: []*proto.ItemSpec{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {Id: 17182}}},
			Spec:          &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{ClassOptions: &proto.ShamanOptions{}}}},
			Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}
		raid := &proto.Raid{Parties: []*proto.Party{{Players: []*proto.Player{player}, Buffs: &proto.PartyBuffs{}}}, Buffs: &proto.RaidBuffs{}, Debuffs: &proto.Debuffs{}, NumActiveParties: 1}
		res := core.RunRaidSim(&proto.RaidSimRequest{Raid: raid, Encounter: core.MakeSingleTargetEncounter(0), SimOptions: &proto.SimOptions{Iterations: 200, RandomSeed: 1}})
		if res.Error != nil {
			t.Fatal(res.Error.Message)
		}
		for _, a := range res.RaidMetrics.Parties[0].Players[0].Auras {
			if a.Id.GetSpellId() == 408505 {
				return a.ProcsAvg
			}
		}
		t.Fatal("no Maelstrom Weapon aura")
		return 0
	}
	one, five := procs("-00000000000000001"), procs("-00000000000000005")
	if one < 0.9*five || one > 1.1*five || five == 0 {
		t.Fatalf("Maelstrom Weapon procs a fight: 1 point %.1f, 5 points %.1f; want the same rate", one, five)
	}
}
