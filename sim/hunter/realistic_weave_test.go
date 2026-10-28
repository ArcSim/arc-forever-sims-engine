package hunter

import (
	"math"
	"testing"
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
)

func weavePlayHunter(t *testing.T, play proto.Hunter_Options_WeavePlay, distance float64) (*core.Simulation, *Hunter) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: distance})
	hunter.Options.WeavePlay = play
	hunter.WeaveReactionMin, hunter.WeaveReactionSpan, hunter.WeaveOvershoot = 0, 0, 0
	switch play {
	case proto.Hunter_Options_WeavePerfectPlay:
		hunter.WeaveReactionMin, hunter.WeaveReactionSpan, hunter.WeaveOvershoot = 100*time.Millisecond, 100*time.Millisecond, 2.0/3
	case proto.Hunter_Options_WeaveNormalPlay:
		hunter.WeaveReactionMin, hunter.WeaveReactionSpan, hunter.WeaveOvershoot = 200*time.Millisecond, 100*time.Millisecond, 2.0/3
	}
	for i := 0; i < 100000 && sim.CurrentTime < time.Second; i++ {
		sim.Step()
	}
	return sim, hunter
}

// walkTo decides a walk now and steps until it ends: how long from deciding to the first step, and from then to arriving.
func walkTo(sim *core.Simulation, hunter *Hunter, to float64) (reaction, run time.Duration) {
	decided := sim.CurrentTime
	hunter.MoveTo(to, sim)
	for i := 0; i < 100000 && !hunter.IsMoving(); i++ {
		sim.Step()
	}
	started := sim.CurrentTime
	for i := 0; i < 100000 && hunter.IsMoving(); i++ {
		sim.Step()
	}
	return started - decided, sim.CurrentTime - started
}

// The weave play options: a reaction from their range between deciding to move and the first step, and again after
// arriving before any button press (Pro Player 100-200 ms, Regular Player 200-300 ms). God Mode walks at once (no move
// delay set) and a press is ready on arrival.
func TestWeavePlayReactions(t *testing.T) {
	for _, c := range []struct {
		play     proto.Hunter_Options_WeavePlay
		min, max time.Duration
	}{
		{proto.Hunter_Options_WeaveGodMode, 0, 0},
		{proto.Hunter_Options_WeavePerfectPlay, 100 * time.Millisecond, 200 * time.Millisecond},
		{proto.Hunter_Options_WeaveNormalPlay, 200 * time.Millisecond, 300 * time.Millisecond},
	} {
		sim, hunter := weavePlayHunter(t, c.play, 11)
		moveWait, _ := walkTo(sim, hunter, 5)
		arrived := sim.CurrentTime
		readyOnArrival := hunter.PressReady(sim)
		for i := 0; i < 100000 && !hunter.PressReady(sim); i++ {
			sim.Step()
		}
		pressWait := sim.CurrentTime - arrived
		if moveWait < c.min || moveWait > c.max || pressWait < c.min || pressWait > c.max || readyOnArrival != (c.max == 0) {
			t.Errorf("%v: walk started %v after deciding, first press %v after arriving (ready on arrival %v), want both %v-%v",
				c.play, moveWait, pressWait, readyOnArrival, c.min, c.max)
		}
	}
}

// A walk can stop between whole yards: the last step is the fraction left. Pro Player's walk out to 11 yd stops 2/3 yd
// past it, so each way across the dead zone is 6.67 yd (0.95 s at 7 yd/s), and the move to 11 then counts as done.
// God Mode stops exactly where asked, 11.6 yd included.
func TestFractionalWalkStops(t *testing.T) {
	speed := 7.0
	yard := float64(time.Second) / speed // the engine's time per yard at 7 yd/s
	near := func(got time.Duration, yards float64) bool { return math.Abs(float64(got)-yards*yard) < float64(2*time.Millisecond) }

	sim, hunter := weavePlayHunter(t, proto.Hunter_Options_WeavePerfectPlay, 5)
	_, out := walkTo(sim, hunter, 11)
	if math.Abs(hunter.DistanceFromTarget-11.67) > 1e-9 || !near(out, 6.67) {
		t.Errorf("Pro Player: walk out to 11 yd ended at %v yd after %v, want 11.67 yd after %v", hunter.DistanceFromTarget, out, time.Duration(6.67*yard))
	}
	if hunter.SpotYards(11) != hunter.DistanceFromTarget {
		t.Errorf("Pro Player: standing at %v yd, a move to 11 yd stops at %v (it would walk again)", hunter.DistanceFromTarget, hunter.SpotYards(11))
	}
	for i := 0; i < 100000 && !hunter.PressReady(sim); i++ {
		sim.Step()
	}
	_, in := walkTo(sim, hunter, 5)
	if hunter.DistanceFromTarget != 5 || !near(in, 6.67) {
		t.Errorf("Pro Player: walk in ended at %v yd after %v, want 5 yd after %v", hunter.DistanceFromTarget, in, time.Duration(6.67*yard))
	}

	sim, hunter = weavePlayHunter(t, proto.Hunter_Options_WeaveGodMode, 5)
	_, out = walkTo(sim, hunter, 11.6)
	if math.Abs(hunter.DistanceFromTarget-11.6) > 1e-9 || !near(out, 6.6) {
		t.Errorf("God Mode: walk to 11.6 yd ended at %v yd after %v, want 11.6 yd after %v", hunter.DistanceFromTarget, out, time.Duration(6.6*yard))
	}
	_, out = walkTo(sim, hunter, 5)
	_, out = walkTo(sim, hunter, 11)
	if hunter.DistanceFromTarget != 11 || !near(out, 6) {
		t.Errorf("God Mode: walk out to 11 yd ended at %v yd after %v, want 11 yd after %v", hunter.DistanceFromTarget, out, time.Duration(6*yard))
	}
}
