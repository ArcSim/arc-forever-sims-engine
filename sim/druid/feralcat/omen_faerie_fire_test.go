package feralcat

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/druid"
)

// "Omen of Clarity can no longer be triggered by Faerie Fire while in Cat, Bear or Dire Bear Form"
// (16 October 2026 beta notes). Many landed Faerie Fires in Cat Form grant no Clearcasting, where as
// many landed Claws do.
func TestFaerieFireDoesNotTriggerOmenInCatForm(t *testing.T) {
	procs := func(spellOf func(*druid.Druid) *core.Spell) bool {
		sim := core.NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 1},
			Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
				Name: "Cat", Class: proto.Class_ClassDruid, Race: proto.Race_RaceNightElf, TalentsString: DefaultTalents,
				Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{}, Spec: DefaultSpecOptions,
				Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
			}}}}},
			Encounter: core.MakeSingleTargetEncounter(0),
		}, simsignals.CreateSignals())
		sim.Reset()
		cat := sim.Raid.Parties[0].Players[0].(druid.DruidAgent).GetDruid()
		if !cat.InForm(druid.Cat) {
			t.Fatal("the Cat druid does not start in Cat Form")
		}
		omen := cat.GetAura("Omen of Clarity")
		if omen == nil || omen.OnSpellHitDealt == nil {
			t.Fatal("no Omen of Clarity trigger on the Cat druid")
		}
		spell := spellOf(cat)
		landed := &core.SpellResult{Outcome: core.OutcomeHit, Target: cat.CurrentTarget}

		done := false
		sim.AddPendingAction(core.NewDelayedAction(core.DelayedActionOptions{DoAt: sim.CurrentTime, OnAction: func(sim *core.Simulation) {
			for range 2000 {
				omen.OnSpellHitDealt(omen, sim, spell, landed)
			}
		}}))
		sim.AddPendingAction(core.NewDelayedAction(core.DelayedActionOptions{DoAt: sim.CurrentTime + core.SpellBatchWindow + time.Millisecond, OnAction: func(sim *core.Simulation) {
			done = true
		}}))
		for !done {
			sim.Step()
		}
		return cat.ClearcastingAura.IsActive()
	}

	if procs(func(d *druid.Druid) *core.Spell { return d.FaerieFire.Spell }) {
		t.Error("Faerie Fire in Cat Form granted Clearcasting")
	}
	if !procs(func(d *druid.Druid) *core.Spell { return d.AutoAttacks.MHAuto() }) {
		t.Error("2000 landed auto attacks in Cat Form granted no Clearcasting; the check above proves nothing")
	}
}
