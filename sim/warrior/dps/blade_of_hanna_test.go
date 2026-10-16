package dps

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
)

// Blade of Hanna (client 70425) adds a stack of Improving Form per landed main-hand swing and perfects the
// form at 30: one 490 Physical strike and +20 to all stats for 20 sec, with the stacks spent. The item
// generator had read it as the strike on every swing.
func TestBladeOfHannaPerfectsTheFormAtThirtyStacks(t *testing.T) {
	player := &proto.Player{
		Name: "Hanna", Race: proto.Race_RaceHuman, Class: proto.Class_ClassWarrior,
		Equipment: &proto.EquipmentSpec{Items: []*proto.ItemSpec{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {Id: 2801}}},
		// No talents: Sword Specialization's extra swings would add stacks of their own.
		Spec:     DefaultOptions,
		Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}
	sim := core.NewSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
		SimOptions: &proto.SimOptions{RandomSeed: 1},
	}, simsignals.CreateSignals())
	sim.Reset()
	war := sim.Raid.Parties[0].Players[0].(*DpsWarrior)
	target := sim.Encounter.AllTargetUnits[0]

	strike := war.GetSpell(core.ActionID{SpellID: 1325508})
	stacks := war.GetAura("Improving Form")
	buff := war.GetAura("Perfected Form")
	if strike == nil || stacks == nil || buff == nil {
		t.Fatal("Blade of Hanna's effect is not registered")
	}
	swing := war.AutoAttacks.MHAuto()
	strength := war.GetStat(stats.Strength)
	for i := 1; i <= 29; i++ {
		swing.CalcAndDealDamage(sim, target, 100, swing.OutcomeAlwaysHit)
	}
	if got := strike.SpellMetrics[target.UnitIndex].Casts; got != 0 || stacks.GetStacks() != 29 || buff.IsActive() {
		t.Fatalf("after 29 swings: %d strikes, %d stacks, buff %v", got, stacks.GetStacks(), buff.IsActive())
	}
	swing.CalcAndDealDamage(sim, target, 100, swing.OutcomeAlwaysHit)
	if got := strike.SpellMetrics[target.UnitIndex].Casts; got != 1 || stacks.IsActive() || !buff.IsActive() {
		t.Fatalf("after 30 swings: %d strikes, stacks active %v, buff %v", got, stacks.IsActive(), buff.IsActive())
	}
	if got := war.GetStat(stats.Strength) - strength; got < 20 {
		t.Fatalf("Perfected Form added %v Strength, want at least 20", got)
	}
}
