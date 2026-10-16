package priest

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// A Shadow priest with Shadowform, on the current benchmark's talents.
func newShadowCritsSim(t *testing.T) (*core.Simulation, *core.Character) {
	t.Helper()
	player := &proto.Player{
		Race:          proto.Race_RaceUndead,
		Class:         proto.Class_ClassPriest,
		Equipment:     &proto.EquipmentSpec{},
		Consumables:   &proto.ConsumesSpec{},
		TalentsString: "02530301103--50052201221302251",
		Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		Spec:          &proto.Player_DpsPriest{DpsPriest: &proto.DpsPriest{Options: &proto.DpsPriest_Options{ClassOptions: &proto.PriestOptions{}}}},
	}
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 42},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim, sim.Raid.Parties[0].Players[0].GetCharacter()
}

// Forever's Devouring Plague ticks crit, although beta client 1.60.1.70009 leaves the rank's
// Periodic Can Crit attribute off.
func TestDevouringPlagueTicksCrit(t *testing.T) {
	sim, character := newShadowCritsSim(t)
	target := character.CurrentTarget

	DevouringPlagueRankMap.Each(func(_ int32, rank *spelldata.Spell) {
		id := rank.ID
		spell := character.GetSpell(core.ActionID{SpellID: id})
		if spell == nil {
			t.Fatalf("Devouring Plague %d is not registered", id)
		}
		spell.BonusCritPercent = 100
		dot := spell.Dot(target)
		dot.Apply(sim)
		dot.TickOnce(sim)

		metrics := spell.SpellMetrics[target.UnitIndex]
		// A crit tick counts in CritTicks only, not in Ticks.
		if metrics.CritTicks != 1 || metrics.Ticks != 0 {
			t.Fatalf("Devouring Plague %d: expected one crit tick at 100%% crit, got %d crit and %d normal ticks", id, metrics.CritTicks, metrics.Ticks)
		}
	})
}

// Forever extends Shadowform's +100% critical strike damage bonus to Shadow Word: Death, which the
// beta client's mask does not name yet. Mind Blast is the reference the mask already covers.
func TestShadowformCritBonusCoversShadowWordDeath(t *testing.T) {
	sim, character := newShadowCritsSim(t)
	attackTable := character.AttackTables[character.CurrentTarget.UnitIndex]
	shadowform := character.GetAuraByID(core.ActionID{SpellID: 15473})
	if shadowform == nil {
		t.Fatal("Shadowform aura is not registered")
	}

	check := func(name string, id int32) {
		spell := character.GetSpell(core.ActionID{SpellID: id})
		if spell == nil {
			t.Fatalf("%s %d is not registered", name, id)
		}
		shadowform.Deactivate(sim)
		if got := spell.CritDamageMultiplier(attackTable); got != 1.5 {
			t.Fatalf("%s without Shadowform: crit multiplier %v, want 1.5", name, got)
		}
		shadowform.Activate(sim)
		if got := spell.CritDamageMultiplier(attackTable); got != 2.0 {
			t.Fatalf("%s in Shadowform: crit multiplier %v, want 2.0", name, got)
		}
	}
	ShadowWordDeathRankMap.Each(func(_ int32, rank *spelldata.Spell) {
		check("Shadow Word: Death", rank.ID)
	})
	check("Mind Blast", MindBlastRankMap.Highest().ID)
}
