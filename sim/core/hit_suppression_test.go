package core

import (
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
)

// MythicSim patch 125: a level 63 target suppresses the first 1% of melee and ranged hit unless the
// request's Encounter.no_hit_suppression drops it. With 8% hit a single-handed attack misses the
// boss 1% of the time with the suppression and never without it; spells are not affected.
func TestNoHitSuppressionSwitch(t *testing.T) {
	for _, tc := range []struct {
		off          bool
		suppression  float64
		missAt8, at9 float64
	}{
		{off: false, suppression: 0.01, missAt8: 0.01, at9: 0},
		{off: true, suppression: 0, missAt8: 0, at9: 0},
	} {
		sim := NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 100},
			Raid: &proto.Raid{Parties: []*proto.Party{{
				Players: []*proto.Player{{
					Race:      proto.Race_RaceHuman,
					Class:     proto.Class_ClassWarrior,
					Buffs:     &proto.IndividualBuffs{},
					Spec:      &proto.Player_DpsWarrior{},
					Equipment: &proto.EquipmentSpec{},
					Rotation:  &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
				}},
				Buffs: &proto.PartyBuffs{},
			}}},
			Encounter: &proto.Encounter{Targets: []*proto.Target{{Level: 63}}, Duration: 180, NoHitSuppression: tc.off},
		}, simsignals.CreateSignals())
		sim.Reset()
		warrior := sim.Raid.Parties[0].Players[0].GetCharacter()
		boss := sim.Encounter.AllTargetUnits[0]
		table := warrior.AttackTables[boss.UnitIndex]
		if table.HitSuppression != tc.suppression {
			t.Errorf("no_hit_suppression=%v: suppression %.3f, want %.3f", tc.off, table.HitSuppression, tc.suppression)
		}
		if table.BaseSpellMissChance != 0.17 {
			t.Errorf("no_hit_suppression=%v: base spell miss %.3f, want 0.17", tc.off, table.BaseSpellMissChance)
		}
		// The naked test warrior counts as dual wielding; leave the 19% penalty out.
		warrior.PseudoStats.DisableDWMissPenalty = true
		spell := &Spell{Unit: &warrior.Unit, ProcMask: ProcMaskMeleeMHSpecial}
		for hit, want := range map[float64]float64{8: tc.missAt8, 9: tc.at9} {
			warrior.stats[stats.PhysicalHitPercent] = hit
			if got := spell.GetPhysicalMissChance(table); !WithinToleranceFloat64(want, got, 1e-9) {
				t.Errorf("no_hit_suppression=%v, %.0f%% hit: miss %.4f, want %.4f", tc.off, hit, got, want)
			}
		}
	}
}
