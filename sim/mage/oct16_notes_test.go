package mage

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

// Presence of Mind (arcane 15), Pyroblast (fire 9) and Heating Up (fire 12), nothing else.
const pomPyroTalents = "000000000000001-000000001001"

func newMageWithTalents(t *testing.T, talents string) (*core.Simulation, *Mage) {
	t.Helper()
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Mage", Class: proto.Class_ClassMage, Race: proto.Race_RaceGnome, TalentsString: talents,
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{},
			Spec:     &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{ClassOptions: &proto.MageOptions{}}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim, sim.Raid.Parties[0].Players[0].(MageAgent).GetMage()
}

// 16 October 2026 beta notes: Ice Lance deals +200% to frozen targets (was +300%).
func TestIceLanceFrozenBonusIsTwoHundredPercent(t *testing.T) {
	if IceLanceFrozenMultiplier != 3 {
		t.Errorf("Ice Lance frozen multiplier %v, want 3 (+200%%)", IceLanceFrozenMultiplier)
	}
}

// 16 October 2026 beta notes: a Fire Blast crit grants 2 Heating Up charges, other crits 1, capped at 3.
func TestFireBlastCritGrantsTwoHeatingUpCharges(t *testing.T) {
	sim, mage := newGnomeMage(t)
	if mage.HotStreakAura == nil {
		t.Fatal("FireTalents no longer take Heating Up")
	}
	target := mage.CurrentTarget
	fireBlast := mage.GetSpell(core.ActionID{SpellID: spellData.FireBlast.Highest().ID})
	scorch := mage.GetSpell(core.ActionID{SpellID: spellData.Scorch.Highest().ID})
	fireBlast.BonusCritPercent = 200
	scorch.BonusCritPercent = 200

	castUntilCrit := func(spell *core.Spell) {
		t.Helper()
		for i := 0; i < 50; i++ {
			before := spell.SpellMetrics[target.UnitIndex].Crits
			spell.SkipCastAndApplyEffects(sim, target)
			if spell.SpellMetrics[target.UnitIndex].Crits > before {
				return
			}
		}
		t.Fatalf("%v never crit", spell.ActionID)
	}

	castUntilCrit(scorch)
	if got := mage.HotStreakAura.GetStacks(); got != 1 {
		t.Fatalf("Heating Up at %d stacks after a Scorch crit, want 1", got)
	}
	mage.HotStreakAura.Deactivate(sim)

	castUntilCrit(fireBlast)
	if got := mage.HotStreakAura.GetStacks(); got != 2 {
		t.Fatalf("Heating Up at %d stacks after a Fire Blast crit, want 2", got)
	}
	castUntilCrit(fireBlast)
	if got := mage.HotStreakAura.GetStacks(); got != 3 {
		t.Errorf("Heating Up at %d stacks after two Fire Blast crits, want the cap of 3", got)
	}
}

// 16 October 2026 beta notes: a Pyroblast cast under Presence of Mind does not spend Heating Up, while
// Presence of Mind itself is still spent; the next Pyroblast without it spends the stacks as before.
func TestPresenceOfMindPyroblastKeepsHeatingUp(t *testing.T) {
	sim, mage := newMageWithTalents(t, pomPyroTalents)
	if mage.HotStreakAura == nil || mage.PresenceOfMindAura == nil {
		t.Fatal("talents did not register Heating Up and Presence of Mind")
	}
	pyroblast := mage.GetSpell(core.ActionID{SpellID: spellData.Pyroblast.Highest().ID})
	pom := mage.GetSpell(core.ActionID{SpellID: spellData.PresenceOfMind.Highest().ID})

	mage.HotStreakAura.Activate(sim)
	mage.HotStreakAura.SetStacks(sim, 2)
	if !pom.Cast(sim, mage.CurrentTarget) {
		t.Fatal("Presence of Mind did not cast")
	}
	if !pyroblast.Cast(sim, mage.CurrentTarget) {
		t.Fatal("Pyroblast did not cast under Presence of Mind")
	}
	if mage.Hardcast.Expires > sim.CurrentTime {
		t.Fatalf("Pyroblast under Presence of Mind is a hardcast ending at %v", mage.Hardcast.Expires)
	}
	if mage.PresenceOfMindAura.IsActive() {
		t.Error("Presence of Mind still up after the Pyroblast it made instant")
	}
	if got := mage.HotStreakAura.GetStacks(); !mage.HotStreakAura.IsActive() || got != 2 {
		t.Fatalf("Heating Up at %d stacks after a Presence of Mind Pyroblast, want 2 kept", got)
	}

	// The next, hardcast Pyroblast spends the stacks.
	for i := 0; i < 100 && !mage.GCD.IsReady(sim); i++ {
		sim.Step()
	}
	if !pyroblast.Cast(sim, mage.CurrentTarget) {
		t.Fatal("second Pyroblast did not cast")
	}
	end := mage.Hardcast.Expires
	if end <= sim.CurrentTime {
		t.Fatal("second Pyroblast was instant without Presence of Mind")
	}
	for i := 0; i < 100 && sim.CurrentTime <= end && mage.HotStreakAura.IsActive(); i++ {
		sim.Step()
	}
	if mage.HotStreakAura.IsActive() {
		t.Errorf("Heating Up still up at %v after a Pyroblast without Presence of Mind", sim.CurrentTime)
	}
}
