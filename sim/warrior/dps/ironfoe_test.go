package dps

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/common/classic"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

// Ironfoe's Fury of Forgewright (1301046) is an equip aura: 6% on a landed melee auto or special
// (ProcTypeMask 0x14), either hand, 100 ms internal cooldown. The switches give Era's 0.8 PPM off the
// main hand only.
func TestIronfoeProcRateAndHands(t *testing.T) {
	rate := func(t *testing.T, procMask core.ProcMask) float64 {
		t.Helper()
		player := &proto.Player{
			Name:          "Ironfoe",
			Race:          proto.Race_RaceHuman,
			Class:         proto.Class_ClassWarrior,
			Equipment:     weaponsOnly(11684, 15238), // Ironfoe, Warlord's Axe
			TalentsString: FuryTalents,
			Spec:          DefaultOptions,
			Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}
		sim := core.NewSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
			SimOptions: &proto.SimOptions{RandomSeed: 7},
		}, simsignals.CreateSignals())
		sim.Reset()
		war := sim.Raid.Parties[0].Players[0].(*DpsWarrior)
		trigger := war.GetAura("Fury of Forgewright")
		if trigger == nil || !trigger.IsActive() {
			t.Fatal("Fury of Forgewright is not up")
		}
		target := sim.Encounter.ActiveTargetUnits[0]
		spell := &core.Spell{ProcMask: procMask}
		const n = 20000
		procs := 0
		for i := 0; i < n; i++ {
			sim.CurrentTime += time.Second
			trigger.OnSpellHitDealt(trigger, sim, spell, &core.SpellResult{Target: target, Outcome: core.OutcomeHit, Damage: 100})
			if !trigger.Icd.IsReady(sim) {
				procs++
			}
		}
		return float64(procs) / n
	}
	near := func(t *testing.T, what string, got, want float64) {
		t.Helper()
		if got < want*0.85 || got > want*1.15 {
			t.Errorf("%s procced on %.2f%% of hits, want about %.2f%%", what, got*100, want*100)
		}
	}

	near(t, "a main-hand auto", rate(t, core.ProcMaskMeleeMHAuto), 0.06)
	near(t, "an off-hand auto", rate(t, core.ProcMaskMeleeOHAuto), 0.06)
	near(t, "a main-hand special", rate(t, core.ProcMaskMeleeMHSpecial), 0.06)

	defer func(hand bool, ppm float64) {
		classic.IronfoeMainHandOnly, classic.IronfoePPM = hand, ppm
	}(classic.IronfoeMainHandOnly, classic.IronfoePPM)
	classic.IronfoeMainHandOnly, classic.IronfoePPM = true, 0.8

	// 0.8 PPM on a 2.4 sec weapon: 0.8 * 2.4 / 60 = 3.2% a hit.
	near(t, "switched, a main-hand auto", rate(t, core.ProcMaskMeleeMHAuto), 0.032)
	if got := rate(t, core.ProcMaskMeleeOHAuto); got != 0 {
		t.Errorf("switched to the main hand, an off-hand auto procced on %.2f%% of hits", got*100)
	}
}
