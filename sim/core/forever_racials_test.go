package core

import (
	"testing"

	"github.com/wowsims/classic/sim/core/proto"
)

// Touch of the Grave: 30% for casters, 20% for Paladins, Rogues and Warriors (Forever patch notes).
func TestTouchOfTheGraveChance(t *testing.T) {
	for class, want := range map[proto.Class]float64{
		proto.Class_ClassMage: 0.3, proto.Class_ClassPriest: 0.3, proto.Class_ClassWarlock: 0.3,
		proto.Class_ClassPaladin: 0.2, proto.Class_ClassRogue: 0.2, proto.Class_ClassWarrior: 0.2,
	} {
		if got := touchOfTheGraveChance(class); got != want {
			t.Errorf("%v: %.2f, want %.2f", class, got, want)
		}
	}
}
