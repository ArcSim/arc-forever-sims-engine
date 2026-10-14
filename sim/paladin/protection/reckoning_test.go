package protection

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/paladin"
)

// A Reckoning proc swings at once and leaves the swing timer alone (beta logs 2747-2749), at most
// once per 1.5 s (client 20177 ProcCategoryRecovery).
func TestReckoningSwingsWithoutResettingTheTimer(t *testing.T) {
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
	next := pal.AutoAttacks.NextAttackAt()

	reck := pal.GetAura("Reckoning - Crit" + pal.Label)
	mh := pal.AutoAttacks.MHAuto()
	casts := func() int32 { return mh.SpellMetrics[pal.CurrentTarget.UnitIndex].Casts }
	crit := &core.SpellResult{Outcome: core.OutcomeCrit, Target: &pal.Unit}

	// 50 crits taken at once (enough to proc), then read the result once the procs' batch window
	// has passed, well before the next normal swing.
	critsAt := func(at time.Duration) (swings int32, next time.Duration) {
		done := false
		sim.AddPendingAction(core.NewDelayedAction(core.DelayedActionOptions{DoAt: at, OnAction: func(sim *core.Simulation) {
			for i := 0; i < 50; i++ {
				reck.OnSpellHitTaken(reck, sim, mh, crit)
			}
		}}))
		sim.AddPendingAction(core.NewDelayedAction(core.DelayedActionOptions{DoAt: at + core.SpellBatchWindow + time.Millisecond, OnAction: func(sim *core.Simulation) {
			swings, next, done = casts(), pal.AutoAttacks.NextAttackAt(), true
		}}))
		for !done {
			sim.Step()
		}
		return
	}

	if swings, got := critsAt(sim.CurrentTime); swings != 1 || got != next {
		t.Errorf("Reckoning: %d swings, next normal swing %s; want 1 swing and the swing at %s untouched", swings, got, next)
	}
	if swings, _ := critsAt(sim.CurrentTime + 1400*time.Millisecond); swings != 1 {
		t.Errorf("Reckoning swings inside the 1.5 s lockout: %d in total, want 1", swings)
	}
	if swings, _ := critsAt(sim.CurrentTime + 100*time.Millisecond); swings != 2 {
		t.Errorf("Reckoning swings after the 1.5 s lockout: %d in total, want 2", swings)
	}
}
