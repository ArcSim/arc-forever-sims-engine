package dps

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

func newRageWarrior(t *testing.T, gear *proto.EquipmentSpec, talents string) (*core.Simulation, *DpsWarrior) {
	t.Helper()
	player := &proto.Player{
		Name:          "Rage",
		Race:          proto.Race_RaceHuman,
		Class:         proto.Class_ClassWarrior,
		Equipment:     gear,
		TalentsString: talents,
		Spec:          DefaultOptions,
		Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim, sim.Raid.Parties[0].Players[0].(*DpsWarrior)
}

// The Rage a landed swing of this outcome gives, from an empty bar.
func rageFromSwing(sim *core.Simulation, war *DpsWarrior, spell *core.Spell, outcome core.HitOutcome, damage float64) float64 {
	war.ResetRageBar(sim, 0)
	result := &core.SpellResult{
		Target:                           sim.Encounter.ActiveTargetUnits[0],
		Outcome:                          outcome,
		Damage:                           damage,
		PostArmorAndResistanceMultiplier: damage,
	}
	war.Unit.OnSpellHitDealt(sim, spell, result)
	return war.CurrentRage()
}

// Client 1.60.1.70170 brought the crit bonus back at 75%; the Warrior development notes of 1 October
// 2026 raise it: "Additional Rage generated from landing Critical Strikes increased to 100% increased
// Rage (was 75%)". A Warrior's critical auto attack, in either hand, gives twice the Rage of the same
// swing as a plain hit, and a crit that deals twice the damage still gives twice.
func TestCritAutoAttackGivesAWarriorTwiceTheRage(t *testing.T) {
	sim, war := newRageWarrior(t, DualWieldGear.GearSet, FuryTalents)
	mhSpeed := war.AutoAttacks.MH().SwingSpeed
	ohSpeed := war.AutoAttacks.OH().SwingSpeed

	for _, tc := range []struct {
		name    string
		spell   *core.Spell
		wantHit float64
	}{
		{"main hand", war.AutoAttacks.MHAuto(), core.BaseRageHitFactor * mhSpeed},
		// Dual Wield Specialization 5/5 pays the off hand 10% a rank more Rage since the 70170 hotfixes.
		{"off hand", war.AutoAttacks.OHAuto(), core.BaseRageHitFactor * ohSpeed / 2 * 1.5},
	} {
		hit := rageFromSwing(sim, war, tc.spell, core.OutcomeHit, 1000)
		crit := rageFromSwing(sim, war, tc.spell, core.OutcomeCrit, 2000)

		if !core.WithinToleranceFloat64(tc.wantHit, hit, 0.000001) {
			t.Errorf("%s hit gave %0.5f Rage, want %0.5f", tc.name, hit, tc.wantHit)
		}
		if !core.WithinToleranceFloat64(tc.wantHit*2, crit, 0.000001) {
			t.Errorf("%s crit gave %0.5f Rage, want %0.5f", tc.name, crit, tc.wantHit*2)
		}
		if !core.WithinToleranceFloat64(2, crit/hit, 0.0000001) {
			t.Errorf("%s crit gave %0.6f times a hit, want 2", tc.name, crit/hit)
		}
	}
}

func TestCritAutoAttackGivesATwoHanderTwiceTheRage(t *testing.T) {
	sim, war := newRageWarrior(t, TwoHandGear.GearSet, ArmsTalents)
	speed := war.AutoAttacks.MH().SwingSpeed
	want := core.TwoHandRageHitFactor * speed

	hit := rageFromSwing(sim, war, war.AutoAttacks.MHAuto(), core.OutcomeHit, 1000)
	crit := rageFromSwing(sim, war, war.AutoAttacks.MHAuto(), core.OutcomeCrit, 2000)

	if !core.WithinToleranceFloat64(want, hit, 0.000001) {
		t.Errorf("2H hit gave %0.5f Rage, want %0.5f", hit, want)
	}
	if !core.WithinToleranceFloat64(want*2, crit, 0.000001) {
		t.Errorf("2H crit gave %0.5f Rage, want %0.5f", crit, want*2)
	}
}

// Only a swing that crits is paid more: a glancing blow and a miss are not, and an ability that crits
// (Heroic Strike, Mortal Strike) gives no Rage on hit at all, as before.
func TestCritRageBonusIsForAutoAttackCritsOnly(t *testing.T) {
	sim, war := newRageWarrior(t, DualWieldGear.GearSet, FuryTalents)
	hit := rageFromSwing(sim, war, war.AutoAttacks.MHAuto(), core.OutcomeHit, 1000)

	if got := rageFromSwing(sim, war, war.AutoAttacks.MHAuto(), core.OutcomeGlance, 1000); got != hit {
		t.Errorf("glance gave %0.5f Rage, want the hit's %0.5f", got, hit)
	}
	if got := rageFromSwing(sim, war, war.AutoAttacks.MHAuto(), core.OutcomeMiss, 0); got != 0 {
		t.Errorf("miss gave %0.5f Rage, want 0", got)
	}
	heroicStrike := war.GetSpell(core.ActionID{SpellID: 25286})
	if heroicStrike == nil {
		t.Fatal("Heroic Strike is not registered")
	}
	if got := rageFromSwing(sim, war, heroicStrike, core.OutcomeCrit, 2000); got != 0 {
		t.Errorf("Heroic Strike crit gave %0.5f Rage, want 0", got)
	}
}
