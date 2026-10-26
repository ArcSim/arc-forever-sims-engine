package core

import (
	"testing"

	"github.com/wowsims/classic/sim/core/proto"
)

// The Weave Play rotation value reads the unit's Weaving setting, so a rotation line can change with it.
func TestAPLValueWeavePlay(t *testing.T) {
	for _, play := range []int32{0, 1, 2} {
		rot := &APLRotation{unit: &Unit{WeavePlay: play}}
		v := rot.newValueWeavePlay(&proto.APLValueWeavePlay{})
		if v == nil {
			t.Fatalf("Weave Play value did not build")
		}
		if got := v.GetInt(nil); got != play {
			t.Errorf("Weave Play reads %d with the setting at %d", got, play)
		}
	}
}
