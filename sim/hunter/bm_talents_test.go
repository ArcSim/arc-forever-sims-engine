package hunter

import (
	"math"
	"testing"
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/simsignals"
	"github.com/wowsims/classic/sim/core/stats"
)

func newTestHunter(t *testing.T, talents string) (*core.Simulation, *Hunter) {
	player := &proto.Player{
		Name:          "Player",
		Race:          proto.Race_RaceHuman,
		Class:         proto.Class_ClassHunter,
		Equipment:     &proto.EquipmentSpec{},
		TalentsString: talents,
		Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{
			PetType:        proto.Hunter_Options_Cat,
			PetUptime:      1,
			PetAttackSpeed: proto.Hunter_Options_One,
		}}},
	}
	sim := core.NewSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
		SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 1},
	}, simsignals.CreateSignals())
	sim.Reset()
	hunter, ok := sim.Raid.Parties[0].Players[0].(*Hunter)
	if !ok || hunter.pet == nil {
		t.Fatal("no hunter with a pet")
	}
	return sim, hunter
}

func bestialDisciplineHunter(t *testing.T, talents string) *Hunter {
	_, hunter := newTestHunter(t, talents)
	return hunter
}

// Bestial Discipline (19590) raises the pet's focus regen percent (Tamed Pet Passive 19589) by 10 per rank,
// which scales the doubled Forever rate (10 focus/s), and keeps 25% of mana regen per rank while casting.
func TestBestialDiscipline(t *testing.T) {
	for _, tc := range []struct {
		name    string
		talents string
		focus   float64
		casting float64
	}{
		{"no points", "50230005001-005005-5002302300500221", 10, 0},
		{"2/2", "50230005001002-005005-5002302300500221", 12, 0.5},
	} {
		hunter := bestialDisciplineHunter(t, tc.talents)
		if got := hunter.pet.CurrentFocusPerSecond(); math.Abs(got-tc.focus) > 1e-9 {
			t.Errorf("%s: pet focus %.3f/s, want %.3f", tc.name, got, tc.focus)
		}
		ratio := hunter.ManaRegenPerSecondWhileCasting() / hunter.ManaRegenPerSecondWhileNotCasting()
		if math.Abs(ratio-tc.casting) > 1e-9 {
			t.Errorf("%s: casting keeps %.3f of mana regen, want %.3f", tc.name, ratio, tc.casting)
		}
	}
}

// Regen while casting stacks across sources (Bestial Discipline, Resourcefulness, Rapid Recuperation) but
// can never exceed full regen.
func TestCastingRegenCap(t *testing.T) {
	hunter := bestialDisciplineHunter(t, "50230005001002-005005-5002302300500221")
	hunter.PseudoStats.SpiritRegenRateCasting += 1.0
	ratio := hunter.ManaRegenPerSecondWhileCasting() / hunter.ManaRegenPerSecondWhileNotCasting()
	if math.Abs(ratio-1) > 1e-9 {
		t.Errorf("casting keeps %.3f of mana regen with 150%% stacked, want 1", ratio)
	}
}

// Intimidation (19577): the pet's next landed attack gains 100% crit. A white swing landing within 50 ms of
// the cast was rolled before it and leaves the buff up; a miss or dodge leaves it up; a landed attack uses it.
func TestIntimidation(t *testing.T) {
	sim, hunter := newTestHunter(t, "5023000505111251-00500523005")
	aura := hunter.IntimidationPetAura
	if aura == nil {
		t.Fatal("Intimidation not registered for a build with the talent")
	}
	white := &core.Spell{ProcMask: core.ProcMaskMeleeMHAuto}
	special := &core.Spell{ProcMask: core.ProcMaskMeleeMHSpecial}
	landed := &core.SpellResult{Outcome: core.OutcomeHit}
	dodged := &core.SpellResult{Outcome: core.OutcomeDodge}
	baseCrit := hunter.pet.GetStat(stats.MeleeCrit)

	aura.Activate(sim)
	if got := hunter.pet.GetStat(stats.MeleeCrit) - baseCrit; math.Abs(got-100*core.CritRatingPerCritChance) > 1e-9 {
		t.Fatalf("pet crit +%.1f, want +100%%", got/core.CritRatingPerCritChance)
	}
	sim.CurrentTime += 10 * time.Millisecond
	aura.OnSpellHitDealt(aura, sim, white, landed)
	if !aura.IsActive() {
		t.Fatal("a swing rolled before the cast used the buff")
	}
	sim.CurrentTime += time.Second
	aura.OnSpellHitDealt(aura, sim, special, dodged)
	if !aura.IsActive() {
		t.Fatal("a dodged attack used the buff")
	}
	aura.OnSpellHitDealt(aura, sim, white, landed)
	if aura.IsActive() || math.Abs(hunter.pet.GetStat(stats.MeleeCrit)-baseCrit) > 1e-9 {
		t.Fatal("the next landed attack did not use the buff")
	}
}
