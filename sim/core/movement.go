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
	if move.repace != nil {
		move.repace()
	}
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
	// how far from the target the walk began (the timeline names a ready Auto Shot's wait from it)
	moveStartDistance float64
	// the walk's steps: the action, its time per yard, when the last step ended, where the walk began, its stops and how
	// many are left, so a run speed change during the walk re-paces the rest of it (repace, set for the walk under way)
	walkSteps  *PendingAction
	walkPeriod time.Duration
	walkStepAt time.Duration
	walkStart  float64
	walkPoints []float64
	walkLeft   int
	// +1 walking away from the target, -1 walking towards it
	walkDir float64
	repace  func()

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
			walkingOut := unit.MovementHandler.moveTarget > MaxMeleeAttackDistance && unit.DistanceFromTarget <= MaxMeleeAttackDistance
			// walking out of melee before a swing retry landed: the timeline shows that wait as left
			if unit.MovementHandler.moveTarget > MaxMeleeAttackDistance && !(walkingOut && ForeverMeleeWhileWalkingOut) {
				unit.AutoAttacks.mh.endRetryWait(sim, "left")
				unit.AutoAttacks.oh.endRetryWait(sim, "left")
			}
			unit.AutoAttacks.CancelAutoSwing(sim)
			if walkingOut {
				unit.AutoAttacks.KeepMeleeWhileWalkingOut(sim)
			}
			unit.MovementHandler.Moving = true
			unit.MovementHandler.swungWhileMoving = false
			unit.MovementHandler.startedInMelee = unit.DistanceFromTarget <= MaxMeleeAttackDistance
			unit.MovementHandler.moveStartedAt = sim.CurrentTime
			unit.MovementHandler.moveStartDistance = unit.DistanceFromTarget
			unit.MovementHandler.offHandWait = unit.AutoAttacks.OffHandWait()
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			unit.MovementHandler.Moving = false
			unit.MovementHandler.repace, unit.MovementHandler.walkSteps = nil, nil
			startedMidWalk := unit.AutoAttacks.MeleeStartedMidMove()
			mhReady, ohReady := unit.AutoAttacks.MeleeTimers()
			unit.AutoAttacks.EnableAutoSwing(sim)

			// a walk that reached 5 yd already started swinging (its attack delay and off-hand wait are set)
			if !startedMidWalk {
				// a walk in from outside melee range turned the melee attack on when it started
				walkedIn := !unit.MovementHandler.startedInMelee && unit.DistanceFromTarget <= MaxMeleeAttackDistance
				if ForeverSwingRetryBeat {
					if walkedIn {
						unit.AutoAttacks.SwingOnArrival(sim, mhReady, ohReady, unit.MovementHandler.moveStartedAt, unit.MovementHandler.offHandWait)
					} else if unit.DistanceFromTarget <= MaxMeleeAttackDistance {
						// a step inside melee range; a walk out leaves the timers running untouched
						unit.AutoAttacks.DelayMeleeBy(sim, attackStartDelay)
					}
				} else {
					// Simulate the delay from starting attack (AB_ATTACK_START_MS sets it for engine comparisons)
					unit.AutoAttacks.DelayMeleeBy(sim, attackStartDelay)
					if walkedIn {
						unit.AutoAttacks.OffHandWaitsForAttackStart(sim, unit.MovementHandler.moveStartedAt, unit.MovementHandler.offHandWait)
					}
				}
			}
			// Forever: Auto Shot fires 0.6 s after you stop, the 0.5 s windup and the moment to start it (measured in game:
			// 0.57-0.73 s after stopping)
			unit.AutoAttacks.logShotWaits(sim, unit.MovementHandler.moveStartedAt, unit.MovementHandler.moveStartDistance)
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

// attackStartDelay: from getting in melee range to the server seeing it, so a swing retry this close after arrival is
// missed (AB_ATTACK_START_MS sets it for comparisons)
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

// DistanceNow is where the unit is right now: DistanceFromTarget only changes as each yard of a walk ends, so mid-yard
// it counts the share already walked. Range checks read this, so a walk off the 5 yd edge leaves melee as it starts
// (not a yard later) and a walk in from 11 yd is no longer in shot range once it starts.
func (unit *Unit) DistanceNow(sim *Simulation) float64 {
	mh := unit.MovementHandler
	if mh == nil || !mh.Moving || mh.walkSteps == nil || mh.walkPeriod <= 0 || mh.walkLeft <= 0 {
		return unit.DistanceFromTarget
	}
	step := mh.stepYards()
	left := float64(mh.walkSteps.NextActionAt-sim.CurrentTime) / float64(mh.walkPeriod)
	return unit.DistanceFromTarget + mh.walkDir*(step-max(0, min(step, left)))
}

// A spell with a cast time cannot start while walking, nor with a walk decided and in its move reaction: the walk
// would cancel it.
func (unit *Unit) castBlockedByMove() bool {
	return unit.MovementHandler.Moving || unit.MovementHandler.PendingMove
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

// WholeYards: a fight's starting distance in whole yards.
func WholeYards(distance float64) float64 {
	return math.Round(distance)
}

// SpotYards: where a move stops, to the hundredth of a yard. With realistic weaving a walk out to shot range ends
// WeaveOvershoot past the minimum range: nobody stops exactly on its edge.
func (unit *Unit) SpotYards(distance float64) float64 {
	spot := math.Round(distance*100) / 100
	if edge := MinRangedAttackDistance + unit.WeaveOvershoot; unit.WeaveOvershoot > 0 && spot >= MinRangedAttackDistance && spot < edge {
		spot = math.Round(edge*100) / 100
	}
	return spot
}

func (unit *Unit) realisticWeave() bool {
	return unit.WeaveReactionMin > 0
}

// weaveReaction: a human reaction on a weave step, drawn evenly from its range (measured weaves come out 2.5-2.6 s at
// best and 2.8-3.0 s for most players, the reactions spread over the four steps)
func (unit *Unit) weaveReaction(sim *Simulation) time.Duration {
	return unit.WeaveReactionMin + time.Duration(sim.RandomFloat("Weave Reaction")*float64(unit.WeaveReactionSpan))
}

// moveReaction: the wait between deciding to move and the first step.
func (unit *Unit) moveReaction(sim *Simulation) time.Duration {
	if unit.realisticWeave() {
		return unit.weaveReaction(sim)
	}
	return unit.MoveDelay
}

// PressReady: with realistic weaving, a button press waits out the reaction to arriving.
func (unit *Unit) PressReady(sim *Simulation) bool {
	return unit.pressReactionUntil == 0 || sim.CurrentTime >= unit.pressReactionUntil
}

// reactToArrival: realistic weaving holds button presses for a reaction after a walk ends, then looks at the rotation
// again right as it is over.
func (unit *Unit) reactToArrival(sim *Simulation) {
	if !unit.realisticWeave() || sim.CurrentTime < 0 {
		return
	}
	unit.pressReactionUntil = sim.CurrentTime + unit.weaveReaction(sim)
	if unit.Rotation == nil || sim.Options.Interactive {
		return
	}
	StartDelayedAction(sim, DelayedActionOptions{
		DoAt:     unit.pressReactionUntil,
		Priority: ActionPriorityLow,
		OnAction: unit.Rotation.DoNextAction,
	})
}

func (unit *Unit) MoveTo(moveRange float64, sim *Simulation) {
	moveRange = unit.SpotYards(moveRange)
	mh := unit.MovementHandler
	// decided during a walk: the reaction runs while the walk goes on, and the new walk starts on arrival (or once the
	// reaction is over, if that is later), the way a player turns around at the spot they were heading for
	if mh.Moving {
		if moveRange == mh.moveTarget || mh.PendingMove || mh.queuedMove {
			return
		}
		unit.countPress(sim)
		queue := func(sim *Simulation) {
			if mh.Moving {
				mh.queuedMove, mh.queuedRange = true, moveRange
			} else {
				unit.moveToNow(moveRange, sim)
			}
		}
		delay := unit.moveReaction(sim)
		if delay <= 0 || sim.CurrentTime < 0 {
			queue(sim)
			return
		}
		mh.PendingMove = true
		waitID := ActionID{OtherID: proto.OtherAction_OtherActionMove, Tag: 1}
		if sim.Log != nil {
			unit.Log(sim, "Casting %s (Cost = 0.000, Cast Time = %s, Effective Time = 0s)", waitID, delay)
		}
		StartDelayedAction(sim, DelayedActionOptions{
			DoAt: sim.CurrentTime + delay,
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
	if !unit.MovementHandler.PendingMove {
		unit.countPress(sim)
	}
	// the move reaction passes between the decision to move and the first step
	if delay := unit.moveReaction(sim); delay > 0 && sim.CurrentTime >= 0 && !unit.MovementHandler.PendingMove {
		unit.MovementHandler.PendingMove = true
		// log-only markers so the timeline can draw the wait in front of the move (no cast, no effects)
		waitID := ActionID{OtherID: proto.OtherAction_OtherActionMove, Tag: 1}
		if sim.Log != nil {
			unit.Log(sim, "Casting %s (Cost = 0.000, Cast Time = %s, Effective Time = 0s)", waitID, delay)
		}
		StartDelayedAction(sim, DelayedActionOptions{
			DoAt: sim.CurrentTime + delay,
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
	moveInterval := math.Copysign(1, moveDistance)
	// the stops on the way: every whole yard between here and the spot, then the spot, so a fractional start or spot
	// makes a step of less than a yard at that end
	points := []float64{}
	next := math.Floor(unit.DistanceFromTarget+1e-9) + 1
	if moveInterval < 0 {
		next = math.Ceil(unit.DistanceFromTarget-1e-9) - 1
	}
	for ; moveInterval*(moveRange-next) > 1e-9; next += moveInterval {
		points = append(points, next)
	}
	points = append(points, moveRange)

	mh := unit.MovementHandler
	// the time per yard at the run speed of the moment
	yardTime := func() time.Duration {
		speed := mh.MoveSpeed
		if out := envFloat("AB_MOVE_OUT_SPEED"); out > 0 && moveDistance > 0 {
			speed = out * mh.MoveSpeed / mh.baseSpeed
		}
		return time.Millisecond * time.Duration(1000/speed)
	}

	// a walk into melee turns the melee attack on now: the timeline draws the off hand's wait from here
	if sim.Log != nil && unit.DistanceFromTarget > MaxMeleeAttackDistance && moveRange <= MaxMeleeAttackDistance {
		unit.AutoAttacks.logOffHandWait(sim)
	}
	mh.moveSpell.Cast(sim, unit.CurrentTarget)

	var step func(sim *Simulation)
	// the next stop is reached at first, each one after it its share of a yard time later (re-added after the step, as a
	// periodic action is)
	startSteps := func(first time.Duration) {
		mh.walkPeriod = yardTime()
		pa := &PendingAction{NextActionAt: first}
		pa.OnAction = func(sim *Simulation) {
			step(sim)
			if mh.walkSteps == pa && mh.walkLeft > 0 {
				pa.NextActionAt = sim.CurrentTime + time.Duration(mh.stepYards()*float64(mh.walkPeriod))
				sim.AddPendingAction(pa)
			} else {
				pa.Cancel(sim)
			}
		}
		mh.walkSteps = pa
		sim.AddPendingAction(pa)
	}
	// a run speed change during the walk (a speed buff gained or running out): the step under way keeps the share
	// already walked, the rest go at the new speed
	mh.repace = func() {
		if !mh.Moving || mh.walkSteps == nil || mh.walkLeft <= 0 {
			return
		}
		period := yardTime()
		if period == mh.walkPeriod {
			return
		}
		left := mh.stepYards() - float64(sim.CurrentTime-mh.walkStepAt)/float64(mh.walkPeriod)
		mh.walkSteps.Cancel(sim)
		startSteps(sim.CurrentTime + time.Duration(left*float64(period)))
	}
	mh.walkStepAt, mh.walkStart, mh.walkPoints, mh.walkLeft, mh.walkDir = sim.CurrentTime, unit.DistanceFromTarget, points, len(points), moveInterval
	step = func(sim *Simulation) {
		wasInMelee := unit.DistanceFromTarget <= 5
		unit.DistanceFromTarget = mh.walkPoints[len(mh.walkPoints)-mh.walkLeft]
		mh.walkStepAt, mh.walkLeft = sim.CurrentTime, mh.walkLeft-1
		unit.MovementHandler.moveAura.SetStacks(sim, int32(unit.DistanceFromTarget))

		// a walk out that kept swinging stops at the step onto 5 yd: walking on, it is past melee range the moment after
		if moveInterval > 0 && unit.DistanceFromTarget >= 5 && unit.AutoAttacks.MeleeStartedMidMove() {
			unit.AutoAttacks.mh.endRetryWait(sim, "left")
			unit.AutoAttacks.oh.endRetryWait(sim, "left")
			unit.AutoAttacks.stopMeleeMidMove(sim)
		}
		if !wasInMelee && unit.DistanceFromTarget <= 5 {
			unit.AutoAttacks.ForceMeleeSwing(sim)
			// swinging starts here, mid-walk, if the walk goes on into melee
			if mh.walkLeft > 0 {
				unit.AutoAttacks.StartMeleeMidMove(sim, unit.MovementHandler.moveStartedAt, unit.MovementHandler.offHandWait)
			}
		}

		if mh.walkLeft == 0 {
			swung := unit.MovementHandler.swungWhileMoving
			unit.MovementHandler.moveAura.Deactivate(sim)
			if mh := unit.MovementHandler; mh.queuedMove {
				mh.queuedMove = false
				unit.moveToNow(mh.queuedRange, sim)
				return
			}
			unit.reactToArrival(sim)
			// a swing that landed mid-walk (melee starts at 5 yd) could not start the walk back then: its decision runs
			// on arrival instead of waiting out a busy GCD. Only then: a walk that ends before its swing keeps waiting
			// for the swing as before. Here and not in the aura's expiry, which also runs when a fight ends mid-walk.
			if swung && !sim.Options.Interactive && unit.Rotation != nil {
				unit.Rotation.DoNextAction(sim)
			}
		}
	}
	startSteps(sim.CurrentTime + time.Duration(mh.stepYards()*float64(yardTime())))
}

// stepYards: the length of the walk's next step, a yard or the fraction left to a whole yard or to the spot.
func (mh *MovementHandler) stepYards() float64 {
	if mh.walkLeft <= 0 {
		return 0
	}
	i := len(mh.walkPoints) - mh.walkLeft
	from := mh.walkStart
	if i > 0 {
		from = mh.walkPoints[i-1]
	}
	return math.Abs(mh.walkPoints[i] - from)
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
