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
func TestChainLightningCostCut(t *testing.T) {
	_, sham := newElementalShaman(t, DefaultTalents)
	ranks := 0
	shaman.ChainLightningRankMap.Each(func(rank int32, config *spelldata.Spell) {
		spell := sham.GetSpell(core.ActionID{SpellID: config.ID})
		if spell == nil {
			t.Fatalf("Chain Lightning rank %d (%d) is not registered", rank, config.ID)
		}
		want := int32(math.Round(config.Cost() * 0.65))
		if spell.Cost.BaseCost != want {
			t.Errorf("Chain Lightning rank %d costs %d before talents, want %d (65%% of the client's %v)", rank, spell.Cost.BaseCost, want, config.Cost())
		}
		ranks++
	})
	if ranks != 4 {
		t.Fatalf("checked %d Chain Lightning ranks, want 4", ranks)
	}
}

// Fire Nova's spell power coefficient is 0.143 (16 October 2026 beta notes), not the client row's 0.214.
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
	if math.Abs(nova.BonusCoefficient-0.143) > 1e-9 {
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
		{"5505001300103051--503352001", 0, 4},
		{"5505101300103051--503352001", 1, 5},
		{"5505201300103051--503352001", 2, 6},
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
