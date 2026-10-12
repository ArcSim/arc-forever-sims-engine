package feralbear

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/druid"
)

func newBearInForm(t *testing.T) (*core.Simulation, *druid.Druid) {
	t.Helper()
	player := &proto.Player{Name: "Bear", Class: proto.Class_ClassDruid, Race: proto.Race_RaceNightElf,
		Equipment: &proto.EquipmentSpec{}, Spec: DefaultSpecOptions,
		Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL}}
	sim := core.NewSim(&proto.RaidSimRequest{SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid:      core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter: core.MakeSingleTargetEncounter(0)}, simsignals.CreateSignals())
	sim.Reset()
	d := sim.Raid.Parties[0].Players[0].(druid.DruidAgent).GetDruid()
	sim.CurrentTime = time.Second
	d.ClearForm(sim)
	if !d.BearForm.Cast(sim, d.CurrentTarget) {
		t.Fatal("Bear Form failed")
	}
	if !d.InForm(druid.Bear) {
		t.Fatal("not in Bear Form")
	}
	return sim, d
}

func bearRageFromSwing(sim *core.Simulation, d *druid.Druid, outcome core.HitOutcome, damage float64) float64 {
	d.ResetRageBar(sim, 0)
	result := &core.SpellResult{
		Target:                           sim.Encounter.ActiveTargetUnits[0],
		Outcome:                          outcome,
		Damage:                           damage,
		PostArmorAndResistanceMultiplier: damage,
	}
	d.Unit.OnSpellHitDealt(sim, d.AutoAttacks.MHAuto(), result)
	return d.CurrentRage()
}

// Client 1.60.1.70170: "Bear Form and Dire Bear Form now generate 75% increased Rage when landing a
// Critical Strike". The client has no row for the number (only a melee-auto proc flag on the form
// passives 1178 and 9635), so it was the patch note's 75%. The Druid development notes of 8 October 2026
// raise it: "Additional Rage generated from landing Critical Strikes increased to 100% increased Rage
// (was 75%)". A Bear's critical auto attack gives twice the Rage of the same swing as a plain hit,
// whatever damage it dealt.
func TestBearCritAutoAttackGivesTwiceTheRage(t *testing.T) {
	sim, d := newBearInForm(t)
	want := core.BaseRageHitFactor * d.AutoAttacks.MH().SwingSpeed

	hit := bearRageFromSwing(sim, d, core.OutcomeHit, 1000)
	crit := bearRageFromSwing(sim, d, core.OutcomeCrit, 2000)
	glance := bearRageFromSwing(sim, d, core.OutcomeGlance, 700)

	if !core.WithinToleranceFloat64(want, hit, 0.000001) {
		t.Errorf("hit gave %0.5f Rage, want %0.5f", hit, want)
	}
	if !core.WithinToleranceFloat64(want*2, crit, 0.000001) {
		t.Errorf("crit gave %0.5f Rage, want %0.5f", crit, want*2)
	}
	if !core.WithinToleranceFloat64(2, crit/hit, 0.0000001) {
		t.Errorf("crit gave %0.6f times a hit, want 2", crit/hit)
	}
	if glance != hit {
		t.Errorf("glance gave %0.5f Rage, want the hit's %0.5f", glance, hit)
	}
}
