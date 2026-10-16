package protection

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/paladin"
)

// Reckoning works with charges (16 October 2026 beta notes): each trigger adds one, up to 4, and the
// next main-hand swing spends them all for one extra attack each, leaving the swing timer alone.
func TestReckoningChargesSpendOnTheNextSwing(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: core.SinglePlayerRaidProto(&proto.Player{
			Name: "prot", Class: proto.Class_ClassPaladin, Race: proto.Race_RaceHuman, TalentsString: ProtDeepTalents,
			Equipment: WeaponAndShield,
			Spec:      DefaultOptions,
			Rotation:  &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	pal := sim.Raid.Parties[0].Players[0].(paladin.PaladinAgent).GetPaladin()
	if pal.Talents.Reckoning == 0 {
		t.Fatal("test build has no Reckoning")
	}
	pal.AutoAttacks.EnableAutoSwing(sim)
	pal.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime) // a swing just landed
	swingAt := pal.AutoAttacks.NextAttackAt()
	speed := pal.AutoAttacks.MainhandSwingSpeed()

	trigger := pal.GetAura("Reckoning - Crit" + pal.Label)
	charges := pal.GetAura("Reckoning" + pal.Label)
	mh := pal.AutoAttacks.MHAuto()
	extra := pal.GetSpell(core.ActionID{OtherID: proto.OtherAction_OtherActionAttack}.WithTag(20178))
	if extra == nil {
		t.Fatal("no Reckoning extra attack spell is registered")
	}
	casts := func(s *core.Spell) int32 { return s.SpellMetrics[pal.CurrentTarget.UnitIndex].Casts }
	crit := &core.SpellResult{Outcome: core.OutcomeCrit, Target: &pal.Unit}

	// Crits taken add charges, capped at 4, and swing nothing on their own. The triggers land after the
	// spell batch window, so the charges are read once it has passed.
	done := false
	sim.AddPendingAction(core.NewDelayedAction(core.DelayedActionOptions{DoAt: sim.CurrentTime, OnAction: func(sim *core.Simulation) {
		for range 50 {
			trigger.OnSpellHitTaken(trigger, sim, mh, crit)
		}
	}}))
	sim.AddPendingAction(core.NewDelayedAction(core.DelayedActionOptions{DoAt: sim.CurrentTime + core.SpellBatchWindow + time.Millisecond, OnAction: func(sim *core.Simulation) {
		done = true
	}}))
	for !done {
		sim.Step()
	}
	if !charges.IsActive() || charges.GetStacks() != paladin.ReckoningMaxCharges {
		t.Fatalf("Reckoning holds %d charges after many crits taken, want %d", charges.GetStacks(), paladin.ReckoningMaxCharges)
	}
	if casts(mh) != 0 || casts(extra) != 0 {
		t.Fatalf("a crit taken swung %d normal and %d extra attacks, want none before the next swing", casts(mh), casts(extra))
	}
	if got := pal.AutoAttacks.NextAttackAt(); got != swingAt {
		t.Fatalf("the charges moved the next swing to %s, want it left at %s", got, swingAt)
	}

	for casts(mh) == 0 {
		sim.Step()
	}
	if casts(extra) != paladin.ReckoningMaxCharges {
		t.Errorf("the next swing spent the charges on %d extra attacks, want %d", casts(extra), paladin.ReckoningMaxCharges)
	}
	if charges.IsActive() {
		t.Errorf("Reckoning still holds %d charges after the swing", charges.GetStacks())
	}
	if got, want := pal.AutoAttacks.NextAttackAt(), swingAt+speed; got != want {
		t.Errorf("the swing after the extra attacks is at %s, want one weapon speed after the last normal swing, %s", got, want)
	}
	if casts(mh) != 1 {
		t.Errorf("%d normal swings by %s, want 1", casts(mh), sim.CurrentTime)
	}
}
