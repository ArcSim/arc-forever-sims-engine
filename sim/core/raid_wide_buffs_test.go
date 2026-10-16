package core

import (
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
)

// Leader of the Pack and Moonkin Aura are raid-wide since the 16 October 2026 beta notes: either set
// on one party reaches every party, and nothing else spreads.
func TestDruidCritAurasReachEveryParty(t *testing.T) {
	all := []*proto.PartyBuffs{
		{LeaderOfThePack: true},
		{MoonkinAura: true, GraceOfAirTotem: true},
		{},
	}
	spreadRaidWidePartyBuffs(all)
	for i, buffs := range all {
		if !buffs.LeaderOfThePack || !buffs.MoonkinAura {
			t.Errorf("party %d: Leader of the Pack %v, Moonkin Aura %v, want both", i, buffs.LeaderOfThePack, buffs.MoonkinAura)
		}
		if buffs.GraceOfAirTotem != (i == 1) {
			t.Errorf("party %d: Grace of Air %v, want it on party 1 alone", i, buffs.GraceOfAirTotem)
		}
	}
}
