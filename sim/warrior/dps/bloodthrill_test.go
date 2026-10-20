package dps

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/warrior"
)

// Patch 176: Bloodthrill opens Overpower from main-hand white swings, Heroic Strike and Cleave only (wowsims/forever
// 788f62ce0f, from combat logs); Mortal Strike, a main-hand special, never does. Off, every landed main-hand hit
// can. Each trial lands one hit on a target with Rend and the window closed, then waits out the spell batch window
// the proc lands in.
func TestBloodthrillWhiteSwingsOnly(t *testing.T) {
	rates := func(on bool) map[string]float64 {
		defer func(old bool) { warrior.BloodthrillWhiteSwingsOnly = old }(warrior.BloodthrillWhiteSwingsOnly)
		warrior.BloodthrillWhiteSwingsOnly = on
		player := &proto.Player{
			Name: "Arms", Race: proto.Race_RaceHuman, Class: proto.Class_ClassWarrior,
			TalentsString: ArmsTalents, Equipment: TwoHandGear.GearSet, Spec: DefaultOptions,
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}
		sim := core.NewSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
			SimOptions: &proto.SimOptions{RandomSeed: 7},
		}, simsignals.CreateSignals())
		sim.Reset()
		war := sim.Raid.Parties[0].Players[0].(*DpsWarrior).GetWarrior()
		if war.Talents.Bloodthrill == 0 {
			t.Fatal("ArmsTalents take no Bloodthrill; pick a build that does")
		}
		target := sim.Encounter.AllTargetUnits[0]
		war.AutoAttacks.CancelAutoSwing(sim)
		war.Rend.Dot(target).Activate(sim)

		var heroicStrike *core.Spell
		for _, spell := range war.Spellbook {
			if spell.Matches(warrior.SpellMaskHeroicStrike) && (heroicStrike == nil || spell.Rank > heroicStrike.Rank) {
				heroicStrike = spell
			}
		}
		out := map[string]float64{}
		for name, spell := range map[string]*core.Spell{
			"swing": war.AutoAttacks.MHAuto(), "Heroic Strike": heroicStrike, "Mortal Strike": war.MortalStrike,
		} {
			if spell == nil {
				t.Fatalf("%s not registered", name)
			}
			procs := 0
			for i := 0; i < 2000; i++ {
				war.OverpowerAura.Deactivate(sim)
				spell.CalcAndDealDamage(sim, target, 1, spell.OutcomeAlwaysHit)
				// The proc lands one spell batch window after the hit.
				for until := sim.CurrentTime + core.SpellBatchWindow; sim.CurrentTime <= until; {
					if sim.Step() {
						break
					}
				}
				war.Rend.Dot(target).Activate(sim)
				if war.OverpowerAura.IsActive() {
					procs++
				}
			}
			out[name] = float64(procs) / 2000
		}
		return out
	}

	on := rates(true)
	t.Logf("switch on: %v", on)
	if on["swing"] == 0 || on["Heroic Strike"] == 0 || on["Mortal Strike"] != 0 {
		t.Errorf("switch on: swing %.3f, Heroic Strike %.3f, Mortal Strike %.3f; want both above 0 and Mortal Strike 0",
			on["swing"], on["Heroic Strike"], on["Mortal Strike"])
	}
	off := rates(false)
	t.Logf("switch off: %v", off)
	if off["Mortal Strike"] == 0 {
		t.Errorf("switch off: Mortal Strike never procs Bloodthrill, want the row's chance")
	}
}
