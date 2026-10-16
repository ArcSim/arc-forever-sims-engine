package feralcat

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
	"github.com/wowsims/forever/sim/druid"
)

// Rake adds 0.05 of attack power to its hit and to every tick, and a running tick reads the attack
// power the druid has then (Hameru's beta logs, 7 October 2026). Measured through the live spell: the
// non-crit hit and a tick, each divided by the multipliers that sit on top of its base damage.
func TestRakeScalesWithAttackPower(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Cat", Class: proto.Class_ClassDruid, Race: proto.Race_RaceNightElf, TalentsString: DefaultTalents,
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{}, Spec: DefaultSpecOptions,
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	cat := sim.Raid.Parties[0].Players[0].(druid.DruidAgent).GetDruid()
	target := cat.CurrentTarget
	rake := cat.Rake.Spell
	attackTable := cat.AttackTables[target.UnitIndex]
	cat.AddStatDynamic(sim, stats.PhysicalCritPercent, -200)

	// The base of a non-crit hit at the druid's current attack power. The last cast leaves the dot up.
	hitBase := func() float64 {
		before := rake.SpellMetrics[target.UnitIndex]
		for range 200 {
			// Each hit lands on a target that is not bleeding, out of Rend and Tear's reach.
			cat.Rake.Dot(target).Deactivate(sim)
			rake.SkipCastAndApplyEffects(sim, target)
		}
		after := rake.SpellMetrics[target.UnitIndex]
		if after.Crits != before.Crits || after.Hits == before.Hits {
			t.Fatalf("Rake: %d crits, %d hits over 200 casts, want only hits", after.Crits-before.Crits, after.Hits-before.Hits)
		}
		damage := (after.TotalDamage - before.TotalDamage) / float64(after.Hits-before.Hits)
		armor, _ := rake.ResistanceMultiplier(sim, false, attackTable)
		return damage / (rake.AttackerDamageMultiplier(attackTable, false) * armor * rake.TargetDamageMultiplier(sim, attackTable, false))
	}

	const share, extra = 0.05, 99.0
	ap1 := rake.MeleeAttackPower(target)
	hit1 := hitBase()
	dot := cat.Rake.Dot(target)
	snapshot1 := dot.SnapshotBaseDamage
	tickMultiplier := dot.CalcSnapshotDamage(sim, target, dot.OutcomeTick).Damage / snapshot1

	cat.AddStatDynamic(sim, stats.AttackPower, extra)
	ap2 := rake.MeleeAttackPower(target)
	if ap2-ap1 != extra {
		t.Fatalf("attack power went from %v to %v, want +%v", ap1, ap2, extra)
	}
	want := share * extra

	// The running dot, snapshotted at ap1, ticks on ap2.
	if got := dot.CalcSnapshotDamage(sim, target, dot.OutcomeTick).Damage/tickMultiplier - snapshot1; math.Abs(got-want) > 1e-6 {
		t.Errorf("running Rake tick gained %.4f from %v attack power, want %.4f", got, extra, want)
	}
	hit2 := hitBase()
	if got := hit2 - hit1; math.Abs(got-want) > 1e-6 {
		t.Errorf("Rake hit gained %.4f from %v attack power, want %.4f (%.2f at %v, %.2f at %v)", got, extra, want, hit1, ap1, hit2, ap2)
	}
	// A fresh dot snapshots the new attack power once, with no share counted twice.
	if got := dot.SnapshotBaseDamage - snapshot1; math.Abs(got-want) > 1e-6 {
		t.Errorf("Rake tick snapshot gained %.4f from %v attack power, want %.4f", got, extra, want)
	}
	if got := dot.CalcSnapshotDamage(sim, target, dot.OutcomeTick).Damage/tickMultiplier - snapshot1; math.Abs(got-want) > 1e-6 {
		t.Errorf("fresh Rake tick gained %.4f from %v attack power, want %.4f", got, extra, want)
	}
	t.Logf("hit base %.2f at %v AP, tick base %.2f", hit1, ap1, snapshot1)
}
