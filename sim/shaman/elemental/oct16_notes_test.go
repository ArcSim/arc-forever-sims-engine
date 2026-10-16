package elemental

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/shaman"
)

func newElementalShaman(t *testing.T, talents string) (*core.Simulation, *shaman.Shaman) {
	t.Helper()
	player := core.WithSpec(&proto.Player{
		Race: proto.Race_RaceTroll, Class: proto.Class_ClassShaman, Equipment: &proto.EquipmentSpec{},
		Consumables: &proto.ConsumesSpec{}, TalentsString: talents,
		Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, &proto.Player_ElementalShaman{ElementalShaman: &proto.ElementalShaman{Options: &proto.ElementalShaman_Options{ClassOptions: &proto.ShamanOptions{}}}})
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim, sim.Raid.Parties[0].Players[0].(shaman.ShamanAgent).GetShaman()
}

// "Chain Lightning mana cost reduced by approximately 35%" on all four ranks (16 October 2026 beta notes).
// Client 70425's rows carry the cut: 145, 200, 250 and 315, about 65% of 225, 305, 390 and 485, and the
// sim must not cut them again.
func TestChainLightningCostCut(t *testing.T) {
	_, sham := newElementalShaman(t, DefaultTalents)
	ranks := 0
	want := []int32{145, 200, 250, 315}
	shaman.ChainLightningRankMap.Each(func(rank int32, config *spelldata.Spell) {
		spell := sham.GetSpell(core.ActionID{SpellID: config.ID})
		if spell == nil {
			t.Fatalf("Chain Lightning rank %d (%d) is not registered", rank, config.ID)
		}
		if spell.Cost.BaseCost != want[rank-1] || int32(math.Round(config.Cost())) != want[rank-1] {
			t.Errorf("Chain Lightning rank %d costs %d before talents (client %v), want %d", rank, spell.Cost.BaseCost, config.Cost(), want[rank-1])
		}
		ranks++
	})
	if ranks != 4 {
		t.Fatalf("checked %d Chain Lightning ranks, want 4", ranks)
	}
}

// Fire Nova's spell power coefficient is 0.143 (16 October 2026 beta notes), the client row's since 70425
// (stored as a float32, so to a millionth).
func TestFireNovaCoefficient(t *testing.T) {
	_, sham := newElementalShaman(t, DefaultTalents)
	var nova *core.Spell
	for _, spell := range sham.Spellbook {
		if spell.ClassSpellMask&shaman.SpellMaskFireNova != 0 && spell.ActionID.Tag == 0 {
			nova = spell
		}
	}
	if nova == nil {
		t.Fatal("no Fire Nova is registered")
	}
	if math.Abs(nova.BonusCoefficient-0.143) > 1e-6 {
		t.Fatalf("Fire Nova's coefficient is %v, want 0.143", nova.BonusCoefficient)
	}
}

// Call of Flame adds 3, 6 or 9 sec to Flame Shock (16 October 2026 beta notes): one 3 sec tick a rank on
// the 12 sec, four tick dot.
func TestCallOfFlameLengthensFlameShock(t *testing.T) {
	for _, row := range []struct {
		talents string
		rank    int32
		ticks   int32
	}{
		{"5505301000001351--553302001", 0, 4},
		{"5505311000001351--553302001", 1, 5},
		{"5505321000001351--553302001", 2, 6},
		{DefaultTalents, 3, 7},
	} {
		_, sham := newElementalShaman(t, row.talents)
		if sham.Talents.CallOfFlame != row.rank {
			t.Fatalf("talent string %q has Call of Flame %d, want %d", row.talents, sham.Talents.CallOfFlame, row.rank)
		}
		dot := sham.FlameShock.RelatedDotSpell.Dot(sham.CurrentTarget)
		if dot.BaseTickCount != row.ticks {
			t.Errorf("Call of Flame %d: Flame Shock ticks %d times, want %d", row.rank, dot.BaseTickCount, row.ticks)
		}
	}
}
