package enhancement

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/shaman"
)

// Mythicsim patch 156: Windfury Weapon on only the item swap's main hand (the Manual Crowd Pummeler
// opener's case) switches the party's Windfury Totem off once the imbued weapon is in the hand, and not
// before.
func TestWindfuryOnTheSwapWeaponReplacesTheTotemAfterTheSwap(t *testing.T) {
	const none, windfury = proto.ShamanImbue_NoImbue, proto.ShamanImbue_WindfuryWeapon
	type state struct{ totem, imbue bool }
	states := func(mh, mhSwap proto.ShamanImbue) (before, after state) {
		sim := core.NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 1},
			Raid: core.SinglePlayerRaidProto(imbueSwapPlayer(mh, mhSwap, nil, ""),
				&proto.PartyBuffs{WindfuryTotem: true}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter: core.MakeSingleTargetEncounter(0),
		}, simsignals.CreateSignals())
		sim.Reset()
		c := sim.Raid.Parties[0].Players[0].GetCharacter()
		read := func() state {
			totem, imbue := c.GetAura("Windfury Totem Trigger"), c.GetAura("Windfury Imbue")
			return state{totem != nil && totem.IsActive(), imbue != nil && imbue.IsActive()}
		}
		before = read()
		c.ItemSwap.SwapItems(sim, proto.APLActionItemSwap_Swap1, false)
		return before, read()
	}

	for _, c := range []struct {
		name          string
		mh, mhSwap    proto.ShamanImbue
		before, after state
	}{
		{"swap's weapon only", none, windfury, state{totem: true}, state{imbue: true}},
		// The totem comes back at its next drop, every 5 seconds (unchanged by the patch).
		{"pull's weapon only", windfury, none, state{imbue: true}, state{}},
		{"both weapons", windfury, windfury, state{imbue: true}, state{imbue: true}},
	} {
		if before, after := states(c.mh, c.mhSwap); before != c.before || after != c.after {
			t.Errorf("%s: (totem, imbue) active %+v before and %+v after the swap, want %+v and %+v",
				c.name, before, after, c.before, c.after)
		}
	}

	shaman.WindfuryImbueReplacesTotemAfterSwap = false
	defer func() { shaman.WindfuryImbueReplacesTotemAfterSwap = true }()
	if before, after := states(none, windfury); !before.totem || !after.totem || !after.imbue {
		t.Errorf("switched off, swap's weapon only: (totem, imbue) %+v before and %+v after the swap, want the totem on throughout beside the imbue", before, after)
	}
}
