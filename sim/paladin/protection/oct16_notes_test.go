package protection

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/paladin"
)

func newProtPaladin(t *testing.T) (*core.Simulation, *paladin.Paladin) {
	t.Helper()
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: core.SinglePlayerRaidProto(&proto.Player{
			Name: "prot", Class: proto.Class_ClassPaladin, Race: proto.Race_RaceHuman, TalentsString: ProtDeepTalents,
			Equipment: WeaponAndShield,
			Spec:      DefaultOptions,
			Rotation:  &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim, sim.Raid.Parties[0].Players[0].(paladin.PaladinAgent).GetPaladin()
}

// Righteous Fury adds 90% threat to Holy attacks (16 October 2026 beta notes; was 60%).
func TestRighteousFuryHolyThreat(t *testing.T) {
	sim, pal := newProtPaladin(t)
	var holy *core.Spell
	for _, spell := range pal.Spellbook {
		if spell.Matches(paladin.SpellMaskConsecration) {
			holy = spell
			break
		}
	}
	if holy == nil {
		t.Fatal("no Consecration to read Holy threat from")
	}
	pal.RighteousFuryAura.Deactivate(sim)
	without := holy.ThreatMultiplier
	pal.RighteousFuryAura.Activate(sim)
	with := holy.ThreatMultiplier
	if math.Abs(with-without-paladin.RighteousFuryHolyThreat) > 1e-9 || paladin.RighteousFuryHolyThreat != 0.9 {
		t.Fatalf("Righteous Fury takes Consecration's threat from %v to %v, want 0.9 more", without, with)
	}
}

// Shield Specialization's block restores 4% of maximum mana (16 October 2026 beta notes; was 6%).
func TestShieldSpecializationRestoresFourPercent(t *testing.T) {
	sim, pal := newProtPaladin(t)
	if pal.Talents.ShieldSpecialization != 3 {
		t.Fatalf("test build has Shield Specialization %d, want 3 for a sure proc", pal.Talents.ShieldSpecialization)
	}
	trigger := pal.GetAura("Shield Specialization" + pal.Label)
	pal.SpendMana(sim, pal.CurrentMana()/2, pal.NewManaMetrics(core.ActionID{OtherID: proto.OtherAction_OtherActionNone}))
	before := pal.CurrentMana()
	block := &core.SpellResult{Outcome: core.OutcomeBlock, Target: &pal.Unit}
	trigger.OnSpellHitTaken(trigger, sim, pal.AutoAttacks.MHAuto(), block)
	if got, want := pal.CurrentMana()-before, 0.04*pal.MaxMana(); math.Abs(got-want) > 1e-6 {
		t.Fatalf("a block restored %v mana, want 4%% of %v (%v)", got, pal.MaxMana(), want)
	}
}
