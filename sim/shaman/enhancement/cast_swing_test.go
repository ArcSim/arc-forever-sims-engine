package enhancement

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/shaman"
)

// A hard-cast Lightning Bolt restarts the swing timer when it completes, even when the next swing
// was not due until after it: beta logs 2701 and 2712 (foreverlogs.gg) show the next swing one
// weapon speed after every completed hard cast. An instant bolt (5 Maelstrom stacks) leaves it alone.
// Lava Burst is a hard cast too.
func TestHardCastRestartsTheSwingTimer(t *testing.T) {
	// The next swing, counted from the end of the cast.
	nextSwingAfterCast := func(spellID int32, maelstrom int32) (next, speed time.Duration) {
		sim := core.NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 1},
			Raid: core.SinglePlayerRaidProto(&proto.Player{
				Name: "enh", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc, TalentsString: "5505031000001-053030031005112251", // DefaultTalents + Lava Burst
				Equipment: &proto.EquipmentSpec{Items: []*proto.ItemSpec{
					{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
					{Id: 17182}, // Sulfuras, Hand of Ragnaros
				}},
				Spec:     &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{ClassOptions: &proto.ShamanOptions{}}}},
				Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
			}, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter: core.MakeSingleTargetEncounter(0),
		}, simsignals.CreateSignals())
		sim.Reset()
		sham := sim.Raid.Parties[0].Players[0].(shaman.ShamanAgent).GetShaman()
		sham.AutoAttacks.EnableAutoSwing(sim)
		sham.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime) // a swing just landed
		if maelstrom > 0 {
			mw := sham.GetAura("Maelstrom Weapon")
			mw.Activate(sim)
			mw.SetStacks(sim, maelstrom)
		}
		spell := sham.GetSpell(core.ActionID{SpellID: spellID})
		if !spell.Cast(sim, sham.CurrentTarget) {
			t.Fatalf("%d did not cast", spellID)
		}
		castEnd := sim.CurrentTime + spell.CurCast.CastTime
		return sham.AutoAttacks.NextAttackAt() - castEnd, sham.AutoAttacks.MainhandSwingSpeed()
	}

	if got, speed := nextSwingAfterCast(403, 0); got != speed { // Lightning Bolt rank 1
		t.Errorf("hard cast: next swing %s after the cast ends, want %s (a full swing)", got, speed)
	}
	if got, speed := nextSwingAfterCast(1326047, 0); got != speed { // Lava Burst rank 3
		t.Errorf("Lava Burst: next swing %s after the cast ends, want %s (a full swing)", got, speed)
	}
	if got, speed := nextSwingAfterCast(403, 5); got != speed {
		t.Errorf("instant bolt: next swing %s from the cast, want %s (untouched)", got, speed)
	}
}
