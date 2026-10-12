package core

import (
	"math"
	"testing"
	"time"
)

// Patch 121: a speed change in the middle of a move keeps its destination. The warrior walks from 5 to
// 12 yards at 7 a second; half a second in it speeds up 30%, so it covers the last 3.5 yards at 9.1 and
// stops at 12. Upstream's re-plan sent it to 7, the length of the move.
func TestSpeedChangeMidMoveKeepsDestination(t *testing.T) {
	for _, keep := range []bool{true, false} {
		MidMoveSpeedChangeKeepsDestination = keep
		sim, fw := setupFakeFearSim()
		fw.DistanceFromTarget = 5
		fw.MoveTo(12, sim)
		pa := &PendingAction{NextActionAt: 500 * time.Millisecond, OnAction: func(sim *Simulation) {
			fw.MultiplyMovementSpeed(sim, 1.3)
		}}
		sim.AddPendingAction(pa)
		for i := 0; i < 10000 && (fw.Moving || sim.CurrentTime < pa.NextActionAt); i++ {
			sim.Step()
		}
		want := 12.0
		if !keep {
			want = 7
		}
		if math.Abs(fw.DistanceFromTarget-want) > 0.01 {
			t.Errorf("keep=%v: stopped at %.2f yards, want %.0f", keep, fw.DistanceFromTarget, want)
		}
	}
	MidMoveSpeedChangeKeepsDestination = true
}
