package mage

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/spelldata"
)

func newGnomeMage(t *testing.T) (*core.Simulation, *Mage) {
	t.Helper()
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Mage", Class: proto.Class_ClassMage, Race: proto.Race_RaceGnome, TalentsString: FireTalents,
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{},
			Spec:     &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{ClassOptions: &proto.MageOptions{}}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim, sim.Raid.Parties[0].Players[0].(MageAgent).GetMage()
}

// Eureka! 1259817's three effects as client 1.60.1.70170 states them (EffectSpellClassMask_0..3, family 3).
// The engine's spell store does not carry the racial rows, so the masks are copied here; the 70124 masks
// were {545397431, 4096, 0, 0} on both of the first two and {1, 4096, 0, 0} on the third.
var (
	eurekaCost     = core.ClassFlags{Family: 3, Mask: [4]uint32{549591799, 4096, 0, 0}}
	eurekaDamage   = core.ClassFlags{Family: 3, Mask: [4]uint32{551688947, 4096, 0, 8}}
	eurekaPeriodic = core.ClassFlags{Family: 3}
)

// The lists Eureka! applies are the client's: every mage spell the sim registers is on the cost, damage
// and (none) periodic lists exactly when the row's three class masks name it. Client 1.60.1.70170 added
// Pyroblast, Frost Nova and the Arcane Missile tick, and emptied the periodic mask.
func TestEurekaListsMatchTheClientRow(t *testing.T) {
	_, mage := newGnomeMage(t)
	cost, damage, periodic := eurekaCost, eurekaDamage, eurekaPeriodic
	lists := mage.EurekaSpells()

	checked := 0
	for _, spell := range mage.Spellbook {
		client := spelldata.Find(spell.ActionID.SpellID)
		if client == nil || client.ClassFlags.IsZero() || spell.ClassSpellMask == 0 {
			continue
		}
		checked++
		name := client.Name
		if got, want := spell.ClassSpellMask&lists.Cost != 0, cost.Matches(client.ClassFlags); got != want {
			t.Errorf("%s %d: on the sim's cost list = %v, client row says %v", name, client.ID, got, want)
		}
		if got, want := spell.ClassSpellMask&lists.Damage != 0, damage.Matches(client.ClassFlags); got != want {
			t.Errorf("%s %d: on the sim's damage list = %v, client row says %v", name, client.ID, got, want)
		}
		if got, want := spell.ClassSpellMask&lists.Tick != 0, periodic.Matches(client.ClassFlags); got != want {
			t.Errorf("%s %d: on the sim's periodic list = %v, client row says %v", name, client.ID, got, want)
		}
	}
	if checked < 15 {
		t.Fatalf("only %d mage spells were comparable with the client rows", checked)
	}
}

// Eureka! gives three casts of its listed spells -10% cost and +10% damage, and nothing to the dots: a
// Fireball's hit gains the 10%, its dot does not, and Pyroblast's separate dot spell is off the lists.
func TestEurekaRaisesHitsButNotDots(t *testing.T) {
	sim, mage := newGnomeMage(t)
	fireball := mage.GetSpell(core.ActionID{SpellID: spellData.Fireball.Highest().ID})
	dot := fireball.Dot(mage.CurrentTarget)
	hitBefore := fireball.DamageMultiplier
	tickBefore := fireball.DamageMultiplier * dot.PeriodicDamageMultiplier
	costBefore := fireball.Cost.GetCurrentCost()

	aura := mage.GetAura("Eureka!")
	if aura == nil {
		t.Fatal("no Eureka! aura on a Gnome mage")
	}
	aura.Activate(sim)
	if got, want := fireball.DamageMultiplier/hitBefore, 1.1; got < want-1e-9 || got > want+1e-9 {
		t.Errorf("Fireball's hit multiplier rose by %v under Eureka!, want 1.1", got)
	}
	// A tick is dealt on the spell's multiplier and the dot's own, which takes the bonus back.
	if got := fireball.DamageMultiplier * dot.PeriodicDamageMultiplier; got < tickBefore-1e-9 || got > tickBefore+1e-9 {
		t.Errorf("Fireball's tick multiplier %v under Eureka!, want it unchanged at %v", got, tickBefore)
	}
	if got, want := fireball.Cost.GetCurrentCost(), costBefore*0.9; got < want-1 || got > want+1 {
		t.Errorf("Fireball costs %v under Eureka!, want %v", got, want)
	}
	pyroblast := mage.GetSpell(core.ActionID{SpellID: spellData.Pyroblast.Highest().ID})
	if got := pyroblast.Cost.GetCurrentCost(); got >= float64(spellData.Pyroblast.Highest().Cost()) {
		t.Errorf("Pyroblast costs %v under Eureka!, want 10%% less than %v", got, spellData.Pyroblast.Highest().Cost())
	}

	aura.Deactivate(sim)
	if got := fireball.DamageMultiplier; got < hitBefore-1e-9 || got > hitBefore+1e-9 {
		t.Errorf("Eureka! left Fireball's hit multiplier at %v, want %v", got, hitBefore)
	}
	if got := fireball.DamageMultiplier * dot.PeriodicDamageMultiplier; got < tickBefore-1e-9 || got > tickBefore+1e-9 {
		t.Errorf("Eureka! left Fireball's tick multiplier at %v, want %v", got, tickBefore)
	}
}
