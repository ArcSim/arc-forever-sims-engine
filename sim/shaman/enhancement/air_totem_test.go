package enhancement

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

// Grace of Air Totem is no longer exclusive with Windfury Totem (16 October 2026 beta notes): the party's
// two air totems both stand, and a Grace of Air the shaman casts leaves the party's Windfury Totem alone.
func TestPartyAirTotemsStackOnTheUnit(t *testing.T) {
	player := &proto.Player{
		Name: "enh", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc, TalentsString: DefaultTalents,
		Equipment: &proto.EquipmentSpec{Items: []*proto.ItemSpec{
			{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
			{Id: 17182}, // Sulfuras, Hand of Ragnaros
		}},
		Spec: &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{
			ClassOptions: &proto.ShamanOptions{},
		}}},
		Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{WindfuryTotem: true, GraceOfAirTotem: true}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	char := sim.Raid.Parties[0].Players[0].GetCharacter()
	up := func(label string) bool { return char.GetAura(label).IsActive() }
	check := func(when string, windfury, grace, ownGrace bool) {
		t.Helper()
		if up("Windfury Totem Trigger") != windfury || up("Grace of Air Totem (External)") != grace || up("Grace Of Air Totem (Self)") != ownGrace {
			t.Errorf("%s: party Windfury procs %v, party Grace of Air %v, own Grace of Air %v; want %v, %v, %v", when,
				up("Windfury Totem Trigger"), up("Grace of Air Totem (External)"), up("Grace Of Air Totem (Self)"), windfury, grace, ownGrace)
		}
	}
	check("both party totems", true, true, false)

	char.GetSpell(core.ActionID{SpellID: 25359}).Cast(sim, char.CurrentTarget)
	check("after casting Grace of Air", true, true, true)
}
