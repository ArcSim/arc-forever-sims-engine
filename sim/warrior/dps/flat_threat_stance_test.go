package dps

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/warrior"
)

func flatThreatTestWarrior(t *testing.T, only bool) (*core.Simulation, *DpsWarrior) {
	t.Helper()
	old := warrior.FlatThreatOnlyInDefensiveStance
	warrior.FlatThreatOnlyInDefensiveStance = only
	t.Cleanup(func() { warrior.FlatThreatOnlyInDefensiveStance = old })
	sim := core.NewSim(&proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(&proto.Player{
			Race: proto.Race_RaceOrc, Class: proto.Class_ClassWarrior,
			Equipment: &proto.EquipmentSpec{}, TalentsString: ArmsTalents, Spec: DefaultOptions,
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
		SimOptions: &proto.SimOptions{RandomSeed: 1},
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim, sim.Raid.Parties[0].Players[0].(*DpsWarrior)
}

// Patch 151: with FlatThreatOnlyInDefensiveStance a Warrior's abilities add their flat bonus threat
// only in Defensive Stance (Sunder Armor's client 206, Heroic Strike's 173); damage threat is kept.
func TestFlatThreatOnlyInDefensiveStance(t *testing.T) {
	sim, war := flatThreatTestWarrior(t, true)
	at := war.AttackTables[war.CurrentTarget.UnitIndex]
	sunder := war.GetSpell(core.ActionID{SpellID: 11597})
	heroicStrike := war.GetSpell(core.ActionID{SpellID: 25286})
	if sunder == nil || heroicStrike == nil {
		t.Fatal("the warrior has no Sunder Armor or Heroic Strike")
	}
	if sunder.FlatThreatBonus != 206 {
		t.Fatalf("Sunder Armor's flat threat is %.1f, want the client's 206", sunder.FlatThreatBonus)
	}
	if war.DefensiveStance.RelatedSelfBuff.IsActive() {
		t.Fatal("the DPS warrior starts in Defensive Stance")
	}

	outside := war.PseudoStats.ThreatMultiplier
	if got := sunder.ThreatFromDamage(sim, core.OutcomeHit, 0, at); got != 0 {
		t.Errorf("Sunder Armor outside Defensive Stance: %.2f threat, want 0", got)
	}
	if got, want := heroicStrike.ThreatFromDamage(sim, core.OutcomeHit, 100, at), 100*heroicStrike.ThreatMultiplier*outside; !core.WithinToleranceFloat64(want, got, 1e-9) {
		t.Errorf("Heroic Strike for 100 outside Defensive Stance: %.2f threat, want %.2f (damage only)", got, want)
	}

	war.DefensiveStance.RelatedSelfBuff.Activate(sim)
	inside := war.PseudoStats.ThreatMultiplier
	if got, want := sunder.ThreatFromDamage(sim, core.OutcomeHit, 0, at), 206*inside; !core.WithinToleranceFloat64(want, got, 1e-9) {
		t.Errorf("Sunder Armor in Defensive Stance: %.2f threat, want %.2f", got, want)
	}
	if got, want := heroicStrike.ThreatFromDamage(sim, core.OutcomeHit, 100, at), (100*heroicStrike.ThreatMultiplier+173)*inside; !core.WithinToleranceFloat64(want, got, 1e-9) {
		t.Errorf("Heroic Strike for 100 in Defensive Stance: %.2f threat, want %.2f", got, want)
	}

	war.BattleStance.RelatedSelfBuff.Activate(sim)
	if got := sunder.ThreatFromDamage(sim, core.OutcomeHit, 0, at); got != 0 {
		t.Errorf("Sunder Armor back in Battle Stance: %.2f threat, want 0", got)
	}
}

// Off (the default), every stance keeps the flat bonus.
func TestFlatThreatInEveryStanceByDefault(t *testing.T) {
	sim, war := flatThreatTestWarrior(t, false)
	at := war.AttackTables[war.CurrentTarget.UnitIndex]
	sunder := war.GetSpell(core.ActionID{SpellID: 11597})
	if got, want := sunder.ThreatFromDamage(sim, core.OutcomeHit, 0, at), 206*war.PseudoStats.ThreatMultiplier; !core.WithinToleranceFloat64(want, got, 1e-9) || got == 0 {
		t.Errorf("Sunder Armor outside Defensive Stance with the switch off: %.2f threat, want %.2f", got, want)
	}
}
