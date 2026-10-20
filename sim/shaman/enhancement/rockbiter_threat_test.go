package enhancement

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/shaman"
)

// Patch 181 (wowsims/forever f128c3d268): each landed melee hit of a Rockbiter weapon casts Rockbiter
// Weapon Proc (20871), 72 threat for every second of the main hand's swing. Off, there is no such spell.
func TestRockbiterThreatProc(t *testing.T) {
	run := func(on bool, mh proto.ShamanImbue, outcome func(*core.Spell) core.OutcomeApplier) (threat float64, registered bool) {
		defer func(old bool) { shaman.RockbiterThreatProc = old }(shaman.RockbiterThreatProc)
		shaman.RockbiterThreatProc = on
		player := &proto.Player{
			Name: "enh", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc,
			Equipment: &proto.EquipmentSpec{Items: []*proto.ItemSpec{
				{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
				{Id: 17182}, // Sulfuras, Hand of Ragnaros, 3.7 speed
			}},
			Spec: &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{
				ClassOptions: &proto.ShamanOptions{ImbueMh: mh},
			}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}
		sim := core.NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 1},
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
		}, simsignals.CreateSignals())
		sim.Reset()
		character := sim.Raid.Parties[0].Players[0].GetCharacter()
		target := sim.Encounter.AllTargetUnits[0]
		proc := character.GetSpell(core.ActionID{SpellID: 20871})
		if proc == nil {
			return 0, false
		}
		swing := character.AutoAttacks.MHAuto()
		before := proc.SpellMetrics[target.UnitIndex].TotalThreat
		swing.CalcAndDealDamage(sim, target, 100, outcome(swing))
		return (proc.SpellMetrics[target.UnitIndex].TotalThreat - before) / character.PseudoStats.ThreatMultiplier, true
	}
	hit := func(s *core.Spell) core.OutcomeApplier { return s.OutcomeAlwaysHit }
	miss := func(s *core.Spell) core.OutcomeApplier { return s.OutcomeAlwaysMiss }

	if got, ok := run(true, proto.ShamanImbue_RockbiterWeapon, hit); !ok || math.Abs(got-72*3.7) > 0.01 {
		t.Errorf("a landed Rockbiter swing adds %.2f threat (registered %v), want 72 x 3.7 = 266.4", got, ok)
	}
	if got, _ := run(true, proto.ShamanImbue_RockbiterWeapon, miss); got != 0 {
		t.Errorf("a missed Rockbiter swing adds %.2f threat, want 0", got)
	}
	if got, _ := run(true, proto.ShamanImbue_WindfuryWeapon, hit); got != 0 {
		t.Errorf("a Windfury Weapon swing adds %.2f Rockbiter threat, want 0", got)
	}
	if _, ok := run(false, proto.ShamanImbue_RockbiterWeapon, hit); ok {
		t.Errorf("switch off: Rockbiter Weapon Proc is registered, want none")
	}
}
