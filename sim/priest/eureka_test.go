package priest

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// Eureka! 1259823's three effects as client 1.60.1.70170 states them (EffectSpellClassMask_0..3, family 6).
// The engine's spell store does not carry the racial rows, so the masks are copied here; the 70124 masks
// were {183811792, 8486916, 64, 0} on the first two and {43024448, 8454144, 64, 0} on the third.
var (
	eurekaCost     = core.ClassFlags{Family: 6, Mask: [4]uint32{150224528, 8388614, 0, 0}}
	eurekaDamage   = core.ClassFlags{Family: 6, Mask: [4]uint32{150224528, 8486918, 0, 0}}
	eurekaPeriodic = core.ClassFlags{Family: 6, Mask: [4]uint32{10485760, 98304, 0, 0}}
)

func newEurekaPriest(t *testing.T, race proto.Race, talents string) (*core.Simulation, *Priest) {
	t.Helper()
	player := &proto.Player{
		Race: race, Class: proto.Class_ClassPriest, Equipment: &proto.EquipmentSpec{}, Consumables: &proto.ConsumesSpec{},
		TalentsString: talents, Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		Spec: &proto.Player_DpsPriest{DpsPriest: &proto.DpsPriest{Options: &proto.DpsPriest_Options{ClassOptions: &proto.PriestOptions{}}}},
	}
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 3},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim, sim.Raid.Parties[0].Players[0].(PriestAgent).GetPriest()
}

// The lists Eureka! applies are the client's: every priest spell the sim registers is on them exactly when
// the row's class masks name it. Client 1.60.1.70170 took Devouring Plague, Shadow Word: Pain and Renew off
// every effect, Holy Fire's dot off the periodic one, and put Shadow Word: Death on the first two. The
// channels Mind Flay, Penance and Starshards keep the periodic bonus.
func TestEurekaListsMatchTheClientRow(t *testing.T) {
	lists := (&Priest{}).EurekaSpells()
	checked := 0
	seen := map[string]bool{}
	for _, c := range []struct {
		race     proto.Race
		talents  string
		wantSeen []string
	}{
		{proto.Race_RaceNightElf, ShadowTalents, []string{"Starshards", "Mind Flay", "Shadow Word: Death", "Mind Blast", "Shadow Word: Pain", "Devouring Plague"}},
		{proto.Race_RaceGnome, SmiteTalents, []string{"Penance", "Smite", "Holy Fire"}},
	} {
		_, priest := newEurekaPriest(t, c.race, c.talents)
		for _, spell := range priest.Spellbook {
			client := spelldata.Find(spell.ActionID.SpellID)
			if client == nil || client.ClassFlags.IsZero() || spell.ClassSpellMask == 0 {
				continue
			}
			checked++
			seen[client.Name] = true
			name := client.Name
			onCost, onDamage, onTick := spell.ClassSpellMask&lists.Cost != 0, spell.ClassSpellMask&lists.Damage != 0, spell.ClassSpellMask&lists.Tick != 0
			clientCost := eurekaCost.Matches(client.ClassFlags)
			clientDamage := eurekaDamage.Matches(client.ClassFlags)
			clientPeriodic := eurekaPeriodic.Matches(client.ClassFlags)

			if onCost != clientCost {
				t.Errorf("%s %d: on the sim's cost list = %v, client row says %v", name, client.ID, onCost, clientCost)
			}
			if (onDamage || onTick) != (clientDamage || clientPeriodic) {
				t.Errorf("%s %d: bonus damage in the sim = %v, client row says %v", name, client.ID, onDamage || onTick, clientDamage || clientPeriodic)
			}
			if clientPeriodic && !onTick {
				t.Errorf("%s %d: the client's periodic effect names it and the sim's tick list does not", name, client.ID)
			}
			// The client deals Penance's bolts from the channel 1316994, which the periodic effect names; the sim
			// registers the learned rank 1316995 as the cast and its dot.
			if onTick && !clientPeriodic && name != "Penance" {
				t.Errorf("%s %d: the sim bonuses its ticks and the client's periodic effect does not name it", name, client.ID)
			}
		}
		for _, name := range c.wantSeen {
			if !seen[name] {
				t.Errorf("%s was not among the spells compared", name)
			}
		}
	}
	if checked < 8 {
		t.Fatalf("only %d priest spells were comparable with the client rows", checked)
	}
}

// Shadow Word: Pain's and Devouring Plague's casts don't cost less, Mind Flay's ticks
// keep the bonus, and Holy Fire's dot gives it back.
func TestEurekaPassesOverTheDots(t *testing.T) {
	sim, priest := newEurekaPriest(t, proto.Race_RaceGnome, ShadowTalents)
	get := func(id int32) *core.Spell {
		spell := priest.GetSpell(core.ActionID{SpellID: id})
		if spell == nil {
			t.Fatalf("spell %d is not registered", id)
		}
		return spell
	}
	pain := get(ShadowWordPainRankMap.Highest().ID)
	plague := get(DevouringPlagueRankMap.Highest().ID)
	blast := get(MindBlastRankMap.Highest().ID)
	flay := get(MindFlayRankMap.Highest().ID)

	painCost, plagueCost, blastCost := pain.Cost.GetCurrentCost(), plague.Cost.GetCurrentCost(), blast.Cost.GetCurrentCost()
	painMult, flayMult := pain.DamageMultiplier, flay.DamageMultiplier

	aura := priest.GetAura("Eureka!")
	if aura == nil {
		t.Fatal("no Eureka! aura on a Gnome priest")
	}
	aura.Activate(sim)
	if got := pain.Cost.GetCurrentCost(); got != painCost {
		t.Errorf("Shadow Word: Pain costs %v under Eureka!, want %v", got, painCost)
	}
	if got := pain.DamageMultiplier; got != painMult {
		t.Errorf("Shadow Word: Pain's damage multiplier is %v under Eureka!, want %v", got, painMult)
	}
	if got, want := blast.Cost.GetCurrentCost(), blastCost*0.9; got < want-1 || got > want+1 {
		t.Errorf("Mind Blast costs %v under Eureka!, want %v", got, want)
	}
	if got := flay.DamageMultiplier / flayMult; got < 1.1-1e-9 || got > 1.1+1e-9 {
		t.Errorf("Mind Flay's damage multiplier rose by %v under Eureka!, want 1.1", got)
	}

	if got := plague.Cost.GetCurrentCost(); got != plagueCost {
		t.Errorf("Devouring Plague costs %v under Eureka!, want %v", got, plagueCost)
	}
}
