package core

import (
	"time"

	"github.com/wowsims/classic/sim/core/proto"
)

type APLValueAutoTimeToNext struct {
	DefaultAPLValueImpl
	unit     *Unit
	autoType proto.APLValueAutoTimeToNext_AttackType
}

func (rot *APLRotation) newValueAutoTimeToNext(config *proto.APLValueAutoTimeToNext) APLValue {
	return &APLValueAutoTimeToNext{
		unit:     rot.unit,
		autoType: config.AutoType,
	}
}
func (value *APLValueAutoTimeToNext) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueAutoTimeToNext) GetDuration(sim *Simulation) time.Duration {
	switch value.autoType {
	case proto.APLValueAutoTimeToNext_Melee:
		return max(0, value.unit.AutoAttacks.NextAttackAt()-sim.CurrentTime)
	case proto.APLValueAutoTimeToNext_MainHand:
		return max(0, value.unit.AutoAttacks.MainhandSwingAt()-sim.CurrentTime)
	case proto.APLValueAutoTimeToNext_OffHand:
		return max(0, value.unit.AutoAttacks.OffhandSwingAt()-sim.CurrentTime)
	case proto.APLValueAutoTimeToNext_Ranged:
		return max(0, value.unit.AutoAttacks.NextRangedAttackAt()-sim.CurrentTime)
	}
	// defaults to Any
	return max(0, value.unit.AutoAttacks.NextAnyAttackAt()-sim.CurrentTime)
}
func (value *APLValueAutoTimeToNext) String() string {
	return "Auto Time To Next"
}

type APLValueAutoSwingTime struct {
	DefaultAPLValueImpl
	unit     *Unit
	autoType proto.APLValueAutoSwingTime_SwingType
}

func (rot *APLRotation) newValueAutoSwingTime(config *proto.APLValueAutoSwingTime) APLValue {
	return &APLValueAutoSwingTime{
		unit:     rot.unit,
		autoType: config.AutoType,
	}
}
func (value *APLValueAutoSwingTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueAutoSwingTime) GetDuration(sim *Simulation) time.Duration {
	switch value.autoType {
	case proto.APLValueAutoSwingTime_MainHand:
		return max(0, value.unit.AutoAttacks.MainhandSwingSpeed())
	case proto.APLValueAutoSwingTime_OffHand:
		return max(0, value.unit.AutoAttacks.OffhandSwingSpeed())
	case proto.APLValueAutoSwingTime_Ranged:
		return max(0, value.unit.AutoAttacks.RangedSwingSpeed())
	}
	// defaults to 0
	return 0
}
func (value *APLValueAutoSwingTime) String() string {
	return "Auto Swing Time"
}

// APLValueAutoTimeSinceLast: time since the weapon last swung (or the Auto Shot fired). Swing time minus time to the
// next swing reads short while a swing retry or the windup holds the next one, which made a trip look swung before it.
type APLValueAutoTimeSinceLast struct {
	DefaultAPLValueImpl
	unit     *Unit
	autoType proto.APLValueAutoSwingTime_SwingType
}

func (rot *APLRotation) newValueAutoTimeSinceLast(config *proto.APLValueAutoTimeSinceLast) APLValue {
	return &APLValueAutoTimeSinceLast{
		unit:     rot.unit,
		autoType: config.AutoType,
	}
}
func (value *APLValueAutoTimeSinceLast) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueAutoTimeSinceLast) GetDuration(sim *Simulation) time.Duration {
	return value.unit.AutoAttacks.TimeSinceLastAttack(sim, value.autoType)
}

// TimeSinceLastAttack: time since that weapon last swung, or the Auto Shot last fired.
func (aa *AutoAttacks) TimeSinceLastAttack(sim *Simulation, which proto.APLValueAutoSwingTime_SwingType) time.Duration {
	switch which {
	case proto.APLValueAutoSwingTime_MainHand:
		return max(0, sim.CurrentTime-aa.mh.lastAttackAt)
	case proto.APLValueAutoSwingTime_OffHand:
		return max(0, sim.CurrentTime-aa.oh.lastAttackAt)
	case proto.APLValueAutoSwingTime_Ranged:
		return max(0, sim.CurrentTime-aa.ranged.lastAttackAt)
	}
	return 0
}
func (value *APLValueAutoTimeSinceLast) String() string {
	return "Auto Time Since Last"
}
