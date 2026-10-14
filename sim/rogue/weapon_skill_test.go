package rogue

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// Weapon skill moves a level 63 boss's table for the hand that carries it (sim/core/weapon_skill.go):
// with Aged Core Leather Gloves (+5 daggers) a dagger's white swings are dodged 6% of the time
// instead of 6.5%, and its glancing blows keep 80-90% of the damage instead of 55-75%. The Defias
// Leather 5 piece's +1 dagger takes it to 6.4% and 60-78%.
func TestWeaponSkillMovesTheBossTable(t *testing.T) {
	type swings struct{ landed, dodges, glances, glanceDamage, hitDamage, hits float64 }
	run := func(extra map[proto.ItemSlot]int32) (mainHand swings) {
		equipment := daggersOnly()
		for slot, id := range extra {
			equipment.Items[slot] = &proto.ItemSpec{Id: id}
		}
		player := core.WithSpec(&proto.Player{
			Race:          proto.Race_RaceHuman,
			Class:         proto.Class_ClassRogue,
			Equipment:     equipment,
			TalentsString: SubtletyTalents,
			Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, DefaultOptions)
		res := core.RunRaidSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{Iterations: 300, RandomSeed: 7},
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
		})
		if res.Error != nil {
			t.Fatal(res.Error.Message)
		}
		for _, action := range res.RaidMetrics.Parties[0].Players[0].Actions {
			if action.Id.GetOtherId() != proto.OtherAction_OtherActionAttack || action.Id.Tag != 1 {
				continue
			}
			for _, target := range action.Targets {
				mainHand.dodges += float64(target.Dodges)
				mainHand.glances += float64(target.Glances)
				mainHand.hits += float64(target.Hits)
				mainHand.glanceDamage += target.GlanceDamage
				mainHand.hitDamage += target.Damage - target.CritDamage - target.GlanceDamage
				mainHand.landed += float64(target.Hits + target.Crits + target.Glances + target.Misses + target.Dodges)
			}
		}
		return mainHand
	}

	defias := map[proto.ItemSlot]int32{
		proto.ItemSlot_ItemSlotChest: 10399, proto.ItemSlot_ItemSlotLegs: 10400, proto.ItemSlot_ItemSlotHands: 10401,
		proto.ItemSlot_ItemSlotFeet: 10402, proto.ItemSlot_ItemSlotWaist: 10403,
	}
	for _, tc := range []struct {
		name   string
		extra  map[proto.ItemSlot]int32
		dodge  float64
		glance float64 // a glance's average share of a normal hit
	}{
		{"base skill", nil, 0.065, 0.65},
		{"Aged Core Leather Gloves", map[proto.ItemSlot]int32{proto.ItemSlot_ItemSlotHands: 18823}, 0.06, 0.85},
		{"Defias Leather 5 piece", defias, 0.064, 0.69},
	} {
		got := run(tc.extra)
		dodge := got.dodges / got.landed
		glance := (got.glanceDamage / got.glances) / (got.hitDamage / got.hits)
		// Four standard errors of the dodge rate over a few thousand swings; glance damage is a ratio
		// of averages over thousands of rolls each.
		if tolerance := 4 * math.Sqrt(tc.dodge*(1-tc.dodge)/got.landed); math.Abs(dodge-tc.dodge) > tolerance {
			t.Errorf("%s: main hand dodged %.4f of %v swings, want %.3f +/- %.4f", tc.name, dodge, got.landed, tc.dodge, tolerance)
		}
		if math.Abs(glance-tc.glance) > 0.01 {
			t.Errorf("%s: glances kept %.3f of a hit, want %.2f", tc.name, glance, tc.glance)
		}
	}
}
