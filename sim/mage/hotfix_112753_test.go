package mage

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

// Client hotfix push 112753 on build 70450 (raidbots DBCache.bin, 20 October 2026), as wowsims/forever
// #594 carries it: Elemental Precision is three ranks of 2/4/6%, Piercing Ice is gone, Ice Barrier's
// ranks gain two 3% Frost damage effects on Piercing Ice's masks, Frigid Power (1326891) gives the same
// for a minute, and Frost Warding's new effect 3 raises the absorb by 10/20% (patch 187).
func TestHotfix112753MageRows(t *testing.T) {
	if got := spellData.ElementalPrecision.Len(); got != 3 {
		t.Errorf("Elemental Precision has %d ranks, want 3", got)
	}
	for rank, want := range []float64{2, 4, 6} {
		if got := spellData.ElementalPrecision.ValueAt(int32(rank + 1)); got != want {
			t.Errorf("Elemental Precision rank %d gives %.0f%% hit, want %.0f%%", rank+1, got, want)
		}
	}
	for rank, want := range []float64{0.10, 0.20} {
		if got := spellData.FrostWarding.EffectAt(3).FractionAt(int32(rank + 1)); math.Abs(got-want) > 1e-12 {
			t.Errorf("Frost Warding rank %d raises the absorb %.2f, want %.2f", rank+1, got, want)
		}
	}
	if TalentTreeSizes != [3]int{19, 17, 18} {
		t.Errorf("tree sizes %v, want Frost at 18 talents", TalentTreeSizes)
	}

	barrier := spellData.IceBarrier.Highest()
	frigid := spellData.IceBarrierTriggered.ByID(frigidPowerSpellID)
	if barrier.ID != 13033 || frigid == nil || frigid.Name != "Frigid Power" || frigid.Duration() != time.Minute {
		t.Fatalf("Ice Barrier %d, Frigid Power %v: want rank 4 (13033) and a 1 min Frigid Power", barrier.ID, frigid)
	}
	if got := barrier.EffectN(1).Average(core.CharacterLevel); got != 819 {
		t.Errorf("Ice Barrier rank 4 absorbs %.0f at level 60, want 811 + 2 x 4 = 819", got)
	}
	for _, row := range []*struct {
		name string
		hit  float64
		dot  float64
	}{
		{"Ice Barrier", barrier.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DAMAGE)).Average(core.CharacterLevel),
			barrier.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DOT)).Average(core.CharacterLevel)},
		{"Frigid Power", frigid.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DAMAGE)).Average(core.CharacterLevel),
			frigid.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DOT)).Average(core.CharacterLevel)},
	} {
		if row.hit != 3 || row.dot != 3 {
			t.Errorf("%s: %.0f%% Frost damage and %.0f%% periodic, want 3 and 3", row.name, row.hit, row.dot)
		}
	}
}

func newHotfixMage(t *testing.T, talents string) (*core.Simulation, *Mage) {
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

// Ice Barrier (patch 188): the cast shields the mage for 1 min, absorbing 819 + 0.1 x Frost spell
// power, raised 20% by Frost Warding 2/2; while the shield holds Frost damage is +3% and Frostfire
// Bolt's hit and DoT +3% each, and Frigid Power adds the same again for the full minute, shield or no
// shield. The cast spends a GCD and starts the 30 sec category cooldown.
func TestIceBarrierShieldsAndRaisesFrostDamage(t *testing.T) {
	// FrostTalents with Frost Warding 2/2.
	sim, mage := newHotfixMage(t, "0500005013--253530000101301251")
	if !mage.Talents.IceBarrier || mage.Talents.FrostWarding != 2 {
		t.Fatal("the build does not take Ice Barrier and Frost Warding 2/2")
	}
	barrier := mage.GetSpell(core.ActionID{SpellID: spellData.IceBarrier.Highest().ID})
	if barrier == nil {
		t.Fatal("Ice Barrier is not castable")
	}
	shield := mage.GetAura("Ice Barrier")
	frigid := mage.GetAuraByID(core.ActionID{SpellID: frigidPowerSpellID})
	if shield == nil || frigid == nil {
		t.Fatal("Ice Barrier or Frigid Power aura not registered")
	}

	table := mage.AttackTables[mage.CurrentTarget.UnitIndex]
	frostbolt := mage.GetSpell(core.ActionID{SpellID: spellData.Frostbolt.Highest().ID})
	ffb := mage.GetSpell(core.ActionID{SpellID: spellData.FrostfireBolt.Highest().ID})
	fireball := mage.GetSpell(core.ActionID{SpellID: spellData.Fireball.Highest().ID})
	multipliers := func() [4]float64 {
		return [4]float64{
			frostbolt.AttackerDamageMultiplier(table, false),
			ffb.AttackerDamageMultiplier(table, false),
			ffb.AttackerDamageMultiplier(table, true),
			fireball.AttackerDamageMultiplier(table, false),
		}
	}
	before := multipliers()
	expect := func(when string, bonus float64) {
		t.Helper()
		got := multipliers()
		for i, name := range []string{"Frostbolt", "Frostfire Bolt's hit", "Frostfire Bolt's DoT"} {
			if math.Abs(got[i]-before[i]-bonus) > 1e-9 {
				t.Errorf("%s: %s multiplier +%.4f, want +%.2f", when, name, got[i]-before[i], bonus)
			}
		}
		if got[3] != before[3] {
			t.Errorf("%s: Fireball multiplier moved %.4f; Ice Barrier names only Frost spells", when, got[3]-before[3])
		}
	}

	if !barrier.Cast(sim, &mage.Unit) {
		t.Fatal("Ice Barrier did not cast")
	}
	if !shield.IsActive() || !frigid.IsActive() || shield.RemainingDuration(sim) != time.Minute || frigid.RemainingDuration(sim) != time.Minute {
		t.Fatal("the cast did not apply both a 1 min shield and a 1 min Frigid Power")
	}
	if want := int32(1.2 * (819 + 0.1*mage.frostSpellPower())); shield.GetStacks() != want {
		t.Errorf("shield holds %d, want %d", shield.GetStacks(), want)
	}
	if got := barrier.CD.TimeToReady(sim); got != 30*time.Second {
		t.Errorf("Ice Barrier ready again in %v, want 30s", got)
	}
	if got := mage.GCD.TimeToReady(sim); got != 1500*time.Millisecond {
		t.Errorf("Ice Barrier spent a %v GCD, want 1.5s", got)
	}
	expect("shield and Frigid Power", 0.06)

	// A broken shield leaves Frigid Power up.
	shield.Deactivate(sim)
	expect("Frigid Power alone", 0.03)
	frigid.Deactivate(sim)
	expect("neither", 0)
}

// A Frost string from before the hotfix gives 19 digits; field 55 no longer exists, and the first 18
// would be read in the new order, so the sim refuses it with a reason (patch 189), as it refuses a
// pre-hotfix Elemental Precision 4/5 or 5/5. A current 18 digit string runs.
func TestPreHotfixFrostTalentStringsAreRejected(t *testing.T) {
	run := func(talents string) *proto.RaidSimResult {
		return core.RunRaidSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 1, Iterations: 1},
			Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
				Name: "Mage", Class: proto.Class_ClassMage, Race: proto.Race_RaceGnome, TalentsString: talents,
				Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{},
				Spec:     &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{ClassOptions: &proto.MageOptions{}}}},
				Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
			}}}}},
			Encounter: core.MakeSingleTargetEncounter(0),
		})
	}
	for _, c := range []struct{ talents, reason string }{
		// The live frost-mage reference before the hotfix.
		{"0502105003-03-0555000301001301251", "tree 3 has 18 talents, the string gives 19"},
		{"0502105003-03-055500030100130125", "Elemental Precision has 3 ranks, the string gives 5"},
		{"05022050001001234-235500001301350310-003", "tree 2 has 17 talents, the string gives 18"},
		{"0-0-0-0", "4 trees, a mage has 3"},
	} {
		result := run(c.talents)
		if result.Error == nil || !strings.Contains(result.Error.Message, c.reason) {
			t.Errorf("%q: error %v, want one saying %q", c.talents, result.Error, c.reason)
		}
	}
	if result := run("0502105003-03-053500000101301251"); result.Error != nil {
		t.Errorf("the re-encoded string fails: %s", result.Error.Message)
	}
}
