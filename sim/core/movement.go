package core

import (
	"container/heap"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/wowsims/classic/sim/core/proto"
)

type MoveModifier struct {
	ActionId *ActionID
	Modifier float64
}

type MoveHeap []MoveModifier

func (h MoveHeap) Len() int           { return len(h) }
func (h MoveHeap) Less(i, j int) bool { return h[i].Modifier > h[j].Modifier }
func (h MoveHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MoveHeap) Push(x any) {
	*h = append(*h, x.(MoveModifier))
}

func (h *MoveHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *MoveHeap) activeModifier() float64 {
	heap := *h
	if len(heap) < 1 {
		panic("Movement Heaps should never be missing their base case!!!")
	}
	n := heap[0]
	return n.Modifier
}

func (h *MoveHeap) Find(actionId *ActionID) int {
	for i, mod := range *h {
		if mod.ActionId == actionId {
			return i
		}
	}
	return -1
}

func (move *MovementHandler) addMoveSpeedModifier(moveHeap *MoveHeap, moveMod MoveModifier) {
	heap.Push(moveHeap, moveMod)
	move.updateMoveSpeed()
}

func (move *MovementHandler) removeMoveSpeedModifier(moveHeap *MoveHeap, actionID *ActionID) {
	index := moveHeap.Find(actionID)
	if index == -1 {
		return
	}
	heap.Remove(moveHeap, index)
	move.updateMoveSpeed()
}

func (move *MovementHandler) updateMoveSpeed() {
	move.MoveSpeed = move.baseSpeed * move.stackingBonus * move.getActveModifier(move.moveSpeedBonuses) * (1 - move.getActveModifier(move.moveSpeedPenalties))
}

func (move *MovementHandler) getActveModifier(moveHeap *MoveHeap) float64 {
	return moveHeap.activeModifier()
}

type MovementHandler struct {
	Moving    bool
	MoveSpeed float64
	// a move decided but not started yet (reaction delay experiment)
	PendingMove bool
	// the destination of the walk under way, and a walk decided during it that starts on arrival
	moveTarget  float64
	queuedMove  bool
	queuedRange float64
	// the move began inside melee range, so the melee attack never turned off
	startedInMelee bool
	// a melee swing landed during this walk
	swungWhileMoving bool
	// a walk into melee turns the melee attack on as it starts, and the off hand's wait is set then
	moveStartedAt time.Duration
	offHandWait   time.Duration

	baseSpeed float64
	// speed bonuses that stack with the rest instead of competing for the highest
	stackingBonus      float64
	moveAura           *Aura
	moveSpell          *Spell
	moveSpeedBonuses   *MoveHeap
	moveSpeedPenalties *MoveHeap
}

func (unit *Unit) initMovement() {
	unit.MovementHandler = &MovementHandler{
		moveSpeedBonuses: &MoveHeap{
			MoveModifier{
				Modifier: 1,
			},
		},
		moveSpeedPenalties: &MoveHeap{
			MoveModifier{
				Modifier: 0,
			},
		},
		baseSpeed:     TernaryFloat64(unit.BaseRunSpeed > 0, unit.BaseRunSpeed, 7.0),
		stackingBonus: 1,
	}
	unit.MovementHandler.updateMoveSpeed()

	unit.MovementHandler.moveAura = unit.GetOrRegisterAura(Aura{
		Label:     "Movement",
		ActionID:  ActionID{OtherID: proto.OtherAction_OtherActionMove},
		Duration:  NeverExpires,
		MaxStacks: 30,

		OnGain: func(aura *Aura, sim *Simulation) {
			if unit.IsChanneling(sim) {
				unit.ChanneledDot.Cancel(sim)
			}
			unit.AutoAttacks.CancelAutoSwing(sim)
			unit.MovementHandler.Moving = true
			unit.MovementHandler.swungWhileMoving = false
			unit.MovementHandler.startedInMelee = unit.DistanceFromTarget <= MaxMeleeAttackDistance
			unit.MovementHandler.moveStartedAt = sim.CurrentTime
			unit.MovementHandler.offHandWait = unit.AutoAttacks.OffHandWait()
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			unit.MovementHandler.Moving = false
			startedMidWalk := unit.AutoAttacks.MeleeStartedMidMove()
			unit.AutoAttacks.EnableAutoSwing(sim)

			// a walk that reached 5 yd already started swinging (its attack delay and off-hand wait are set)
			if !startedMidWalk {
				// Simulate the delay from starting attack (AB_ATTACK_START_MS sets it for engine comparisons)
				unit.AutoAttacks.DelayMeleeBy(sim, attackStartDelay)
				// a walk in from outside melee range turned the melee attack on when it started
				if !unit.MovementHandler.startedInMelee && unit.DistanceFromTarget <= MaxMeleeAttackDistance {
					unit.AutoAttacks.OffHandWaitsForAttackStart(sim, unit.MovementHandler.moveStartedAt, unit.MovementHandler.offHandWait)
				}
			}
			// Forever: Auto Shot fires 0.6 s after you stop, the 0.5 s windup and the moment to start it (measured in game:
			// 0.57-0.73 s after stopping)
			unit.AutoAttacks.DelayRangedUntil(sim, sim.CurrentTime+AutoShotStandStill)
		},
	})

	unit.MovementHandler.moveSpell = unit.GetOrRegisterSpell(SpellConfig{
		ActionID: ActionID{OtherID: proto.OtherAction_OtherActionMove},
		Flags:    SpellFlagMeleeMetrics,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			unit.MovementHandler.moveAura.Activate(sim)
			unit.MovementHandler.moveAura.SetStacks(sim, int32(unit.DistanceFromTarget))
		},
	})
}

// MultiplyStackingMoveSpeed scales the run speed on top of the strongest regular speed bonus.
func (unit *Unit) MultiplyStackingMoveSpeed(mult float64) {
	unit.MovementHandler.stackingBonus *= mult
	unit.MovementHandler.updateMoveSpeed()
}

var noShotRecheck = os.Getenv("AB_NO_SHOT_RECHECK") != ""

var shotRecheckDelay = func() time.Duration {
	if ms, err := strconv.Atoi(os.Getenv("AB_SHOT_RECHECK_MS")); err == nil {
		return time.Duration(ms) * time.Millisecond
	}
	return 0
}()

var attackStartDelay = func() time.Duration {
	if ms, err := strconv.Atoi(os.Getenv("AB_ATTACK_START_MS")); err == nil {
		return time.Duration(ms) * time.Millisecond
	}
	return time.Millisecond * 50
}()

// AutoShotStandStill: AB_NO_STAND_STILL=1 turns it off for comparisons, AB_STAND_STILL_MS sets it.
var AutoShotStandStill = func() time.Duration {
	if os.Getenv("AB_NO_STAND_STILL") != "" {
		return 0
	}
	if ms, err := strconv.Atoi(os.Getenv("AB_STAND_STILL_MS")); err == nil {
		return time.Duration(ms) * time.Millisecond
	}
	return time.Millisecond * 600
}()

func (unit *Unit) IsMoving() bool {
	return unit.MovementHandler.Moving
}

// Footwork experiment (local knob, off by default): AB_MOVE_OUT_SPEED sets the speed in yd/s for
// moves away from the target (4.5 = backpedal).
func envFloat(name string) float64 {
	v, err := strconv.ParseFloat(os.Getenv(name), 64)
	if err != nil {
		return 0
	}
	return v
}

// WholeYards: distances are whole yards (a walk steps a yard at a time, so a 4.5 yd spot was never reached and froze the
// rotation for the rest of the fight)
func WholeYards(distance float64) float64 {
	return math.Round(distance)
}

func (unit *Unit) MoveTo(moveRange float64, sim *Simulation) {
	moveRange = WholeYards(moveRange)
	mh := unit.MovementHandler
	// decided during a walk: the reaction runs while the walk goes on, and the new walk starts on arrival (or once the
	// reaction is over, if that is later), the way a player turns around at the spot they were heading for
	if mh.Moving {
		if moveRange == mh.moveTarget || mh.PendingMove || mh.queuedMove {
			return
		}
		queue := func(sim *Simulation) {
			if mh.Moving {
				mh.queuedMove, mh.queuedRange = true, moveRange
			} else {
				unit.moveToNow(moveRange, sim)
			}
		}
		if unit.MoveDelay <= 0 || sim.CurrentTime < 0 {
			queue(sim)
			return
		}
		mh.PendingMove = true
		waitID := ActionID{OtherID: proto.OtherAction_OtherActionMove, Tag: 1}
		if sim.Log != nil {
			unit.Log(sim, "Casting %s (Cost = 0.000, Cast Time = %s, Effective Time = 0s)", waitID, unit.MoveDelay)
		}
		StartDelayedAction(sim, DelayedActionOptions{
			DoAt: sim.CurrentTime + unit.MoveDelay,
			OnAction: func(sim *Simulation) {
				if sim.Log != nil {
					unit.Log(sim, "Completed cast %s", waitID)
				}
				mh.PendingMove = false
				queue(sim)
			},
		})
		return
	}
	if moveRange == unit.DistanceFromTarget {
		return
	}
	// the unit's MoveDelay passes between the decision to move and the first step
	if unit.MoveDelay > 0 && sim.CurrentTime >= 0 && !unit.MovementHandler.PendingMove {
		unit.MovementHandler.PendingMove = true
		// log-only markers so the timeline can draw the wait in front of the move (no cast, no effects)
		waitID := ActionID{OtherID: proto.OtherAction_OtherActionMove, Tag: 1}
		if sim.Log != nil {
			unit.Log(sim, "Casting %s (Cost = 0.000, Cast Time = %s, Effective Time = 0s)", waitID, unit.MoveDelay)
		}
		StartDelayedAction(sim, DelayedActionOptions{
			DoAt: sim.CurrentTime + unit.MoveDelay,
			OnAction: func(sim *Simulation) {
				if sim.Log != nil {
					unit.Log(sim, "Completed cast %s", waitID)
				}
				unit.MovementHandler.PendingMove = false
				unit.moveToNow(moveRange, sim)
			},
		})
		return
	}
	unit.moveToNow(moveRange, sim)
}

func (unit *Unit) moveToNow(moveRange float64, sim *Simulation) {
	if moveRange == unit.DistanceFromTarget {
		return
	}

	unit.MovementHandler.moveTarget = moveRange
	moveDistance := moveRange - unit.DistanceFromTarget
	moveTicks := math.Abs(moveDistance)
	moveInterval := moveDistance / float64(moveTicks)

	speed := unit.MovementHandler.MoveSpeed
	if out := envFloat("AB_MOVE_OUT_SPEED"); out > 0 && moveDistance > 0 {
		speed = out * unit.MovementHandler.MoveSpeed / unit.MovementHandler.baseSpeed
	}

	// a walk into melee turns the melee attack on now: the timeline draws the off hand's wait from here
	if sim.Log != nil && unit.DistanceFromTarget > MaxMeleeAttackDistance && moveRange <= MaxMeleeAttackDistance {
		unit.AutoAttacks.logOffHandWait(sim)
	}
	unit.MovementHandler.moveSpell.Cast(sim, unit.CurrentTarget)

	sim.AddPendingAction(NewPeriodicAction(sim, PeriodicActionOptions{
		Period:          time.Millisecond * time.Duration(1000/speed),
		NumTicks:        int(moveTicks),
		TickImmediately: false,

		OnAction: func(sim *Simulation) {
			wasInMelee := unit.DistanceFromTarget <= 5
			unit.DistanceFromTarget += moveInterval
			unit.MovementHandler.moveAura.SetStacks(sim, int32(unit.DistanceFromTarget))

			if !wasInMelee && unit.DistanceFromTarget <= 5 {
				unit.AutoAttacks.ForceMeleeSwing(sim)
				// swinging starts here, mid-walk, if the walk goes on into melee
				if unit.DistanceFromTarget != moveRange {
					unit.AutoAttacks.StartMeleeMidMove(sim, unit.MovementHandler.moveStartedAt, unit.MovementHandler.offHandWait)
				}
			}

			if unit.DistanceFromTarget == moveRange {
				swung := unit.MovementHandler.swungWhileMoving
				unit.MovementHandler.moveAura.Deactivate(sim)
				if mh := unit.MovementHandler; mh.queuedMove {
					mh.queuedMove = false
					unit.moveToNow(mh.queuedRange, sim)
					return
				}
				// a swing that landed mid-walk (melee starts at 5 yd) could not start the walk back then: its decision runs
				// on arrival instead of waiting out a busy GCD. Only then: a walk that ends before its swing keeps waiting
				// for the swing as before. Here and not in the aura's expiry, which also runs when a fight ends mid-walk.
				if swung && !sim.Options.Interactive && unit.Rotation != nil {
					unit.Rotation.DoNextAction(sim)
				}
			}
		},
	}))
}

// A move speed increase of 30% should be represented as 1.30 and a move speed slow of 70% should be respresented as 0.70
func (unit *Unit) AddMoveSpeedModifier(actionId *ActionID, modifier float64) {
	moveSpeedMod := MoveModifier{
		ActionId: actionId,
		Modifier: modifier,
	}
	if moveSpeedMod.Modifier < 1 {
		unit.MovementHandler.addMoveSpeedModifier(unit.MovementHandler.moveSpeedPenalties, moveSpeedMod)
	} else {
		unit.MovementHandler.addMoveSpeedModifier(unit.MovementHandler.moveSpeedBonuses, moveSpeedMod)
	}

}

func (unit *Unit) RemoveMoveSpeedModifier(actionID *ActionID) {
	unit.MovementHandler.removeMoveSpeedModifier(unit.MovementHandler.moveSpeedPenalties, actionID)
	unit.MovementHandler.removeMoveSpeedModifier(unit.MovementHandler.moveSpeedBonuses, actionID)
}
