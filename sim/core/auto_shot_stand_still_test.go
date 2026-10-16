package core

import (
	"testing"
	"time"
)

// 16 October 2026 beta notes: Auto Shot needs 0.5 seconds of standing still (sim/hunter
// TestAutoShotNeedsHalfASecondStandingStill). Melee swings do not: a main-hand swing due 0.2 seconds
// after a move ends lands on time, and the move's end is recorded.
func TestStandStillOnlyHoldsTheRangedAuto(t *testing.T) {
	sim, fw := setupFakeFearSim()
	fw.DistanceFromTarget = MaxMeleeRange
	sim.PrePull()
	aa := &fw.AutoAttacks
	for i := 0; i < 100000 && aa.mh.previousSwing <= 0; i++ {
		sim.Step()
	}
	due := aa.mh.swingAt
	if due-sim.CurrentTime < time.Second {
		t.Fatalf("test setup: the next swing is due %v after now, want at least 1s", due-sim.CurrentTime)
	}
	moveEnd := due - 200*time.Millisecond
	sim.AddPendingAction(&PendingAction{NextActionAt: due - time.Second, OnAction: func(sim *Simulation) {
		fw.MoveDuration(moveEnd-sim.CurrentTime, sim)
	}})
	for i := 0; i < 100000 && aa.mh.previousSwing < due; i++ {
		sim.Step()
	}
	if fw.movementEndedAt != moveEnd {
		t.Errorf("move end recorded at %v, want %v", fw.movementEndedAt, moveEnd)
	}
	if aa.mh.previousSwing != due {
		t.Errorf("main-hand swing landed at %v, want %v", aa.mh.previousSwing, due)
	}
}
