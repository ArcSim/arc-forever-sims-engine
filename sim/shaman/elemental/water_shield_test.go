package elemental

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
	"github.com/wowsims/forever/sim/shaman"
)

const waterShieldID = 408510

func castSpell(id int32) *proto.APLAction {
	return &proto.APLAction{Action: &proto.APLAction_CastSpell{
		CastSpell: &proto.APLActionCastSpell{SpellId: &proto.ActionID{RawId: &proto.ActionID_SpellId{SpellId: id}}},
	}}
}

// waterShieldRun sims a shaman with Water Shield talented (and nothing else) who casts Water Shield before
// the pull and whenever it is down, then the given spell. It reports the mana Water Shield restored and how
// often it was cast, averaged over the iterations.
func waterShieldRun(t *testing.T, spellID int32) (mana float64, casts float64) {
	t.Helper()
	notActive := &proto.APLValue{Value: &proto.APLValue_Not{Not: &proto.APLValueNot{Val: &proto.APLValue{
		Value: &proto.APLValue_AuraIsActive{AuraIsActive: &proto.APLValueAuraIsActive{AuraId: &proto.ActionID{RawId: &proto.ActionID_SpellId{SpellId: waterShieldID}}}},
	}}}}
	// Enough crit and mana that every fight spends all three globes more than once.
	var bonus stats.Stats
	bonus[stats.SpellCritRating] = 50 * core.SpellCritRatingPerCritPercent
	bonus[stats.Intellect] = 300
	player := &proto.Player{
		BonusStats: &proto.UnitStats{Stats: bonus.ToProtoArray()},
		Name:       "ele", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc, TalentsString: "--000000001",
		Equipment: &proto.EquipmentSpec{},
		Spec:      &proto.Player_ElementalShaman{ElementalShaman: &proto.ElementalShaman{Options: &proto.ElementalShaman_Options{ClassOptions: &proto.ShamanOptions{}}}},
		Rotation: &proto.APLRotation{
			Type: proto.APLRotation_TypeAPL,
			PrepullActions: []*proto.APLPrepullAction{{Action: castSpell(waterShieldID),
				DoAtValue: &proto.APLValue{Value: &proto.APLValue_Const{Const: &proto.APLValueConst{Val: "-2s"}}}}},
			PriorityList: []*proto.APLListItem{
				{Action: &proto.APLAction{Condition: notActive, Action: castSpell(waterShieldID).Action}},
				{Action: castSpell(spellID)},
			},
		},
	}
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
		SimOptions: &proto.SimOptions{Iterations: 20, RandomSeed: 1},
	})
	if res.Error != nil {
		t.Fatal(res.Error.Message)
	}
	p := res.RaidMetrics.Parties[0].Players[0]
	for _, r := range p.Resources {
		if r.Id.GetSpellId() == waterShieldID && r.Type == proto.ResourceType_ResourceTypeMana {
			mana += r.Gain
		}
	}
	for _, a := range p.Actions {
		if a.Id.GetSpellId() == waterShieldID {
			for _, target := range a.Targets {
				casts += float64(target.Casts)
			}
		}
	}
	return mana / 20, casts / 20
}

// The 16 October 2026 beta notes: "Water Shield: Can now trigger from any critical spell cast by the
// Shaman (was only critical heals)". Lightning Bolt crits spend globes for 2% maximum mana each, at most
// one every 3.5 sec, and the shield is recast once the last globe goes (Kerani's Elemental kept all three
// for the whole fight). Fire Nova, a totem's effect, spends none.
func TestWaterShieldSpendsGlobesOnSpellCrits(t *testing.T) {
	mana, casts := waterShieldRun(t, 15208) // Lightning Bolt rank 10
	if mana <= 0 || casts < 3 {
		t.Errorf("Lightning Bolt: Water Shield restored %.0f mana and was cast %.1f times a fight, want mana and recasts", mana, casts)
	}
	if mana, casts := waterShieldRun(t, 408345); mana != 0 || casts != 1 {
		t.Errorf("Fire Nova: Water Shield restored %.0f mana and was cast %.1f times, want 0 and 1", mana, casts)
	}

	defer func(old bool) { shaman.WaterShieldProcsOnSpellCrits = old }(shaman.WaterShieldProcsOnSpellCrits)
	shaman.WaterShieldProcsOnSpellCrits = false
	if mana, casts := waterShieldRun(t, 15208); mana != 0 || casts != 1 {
		t.Errorf("switch off: Water Shield restored %.0f mana and was cast %.1f times, want 0 and 1", mana, casts)
	}
}
