package priest

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// Shadow talents with Inner Focus, Shadow Weaving 3/3 and Early Demise 2/2 (the benchmark build).
const client70170Talents = "025300011303--500222501201302251"

func newClient70170Sim(t *testing.T, race proto.Race, executeFromStart bool) (*core.Simulation, *core.Character) {
	t.Helper()
	player := &proto.Player{
		Race:          race,
		Class:         proto.Class_ClassPriest,
		Equipment:     &proto.EquipmentSpec{},
		Consumables:   &proto.ConsumesSpec{},
		TalentsString: client70170Talents,
		Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		Spec:          &proto.Player_DpsPriest{DpsPriest: &proto.DpsPriest{Options: &proto.DpsPriest_Options{ClassOptions: &proto.PriestOptions{}}}},
	}
	encounter := core.MakeSingleTargetEncounter(0)
	if executeFromStart {
		// Every execute proportion at 100% opens every execute phase at the first advance of the clock.
		encounter.ExecuteProportion_90 = 1
		encounter.ExecuteProportion_45 = 1
		encounter.ExecuteProportion_35 = 1
		encounter.ExecuteProportion_25 = 1
		encounter.ExecuteProportion_20 = 1
	}
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 7},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  encounter,
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim, sim.Raid.Parties[0].Players[0].GetCharacter()
}

// Shadow Word: Death's crit under Early Demise, measured as the share of casts that crit with the
// target out of and in the execute range.
func shadowWordDeathCritShare(t *testing.T, executePhase bool) float64 {
	t.Helper()
	sim, character := newClient70170Sim(t, proto.Race_RaceUndead, executePhase)
	target := character.CurrentTarget
	if executePhase {
		sim.Step() // the first advance of the clock opens the execute phases
		if !sim.IsExecutePhase20() {
			t.Fatal("the sim is not in the 20% execute phase after the first step")
		}
	} else if sim.IsExecutePhase20() {
		t.Fatal("the sim starts in the 20% execute phase")
	}

	spell := character.GetSpell(core.ActionID{SpellID: ShadowWordDeathRankMap.Highest().ID})
	if spell == nil {
		t.Fatal("Shadow Word: Death is not registered")
	}
	const casts = 6000
	for range casts {
		spell.ApplyEffects(sim, target, spell)
	}
	m := spell.SpellMetrics[target.UnitIndex]
	landed := m.Hits + m.Crits
	if landed < casts*8/10 {
		t.Fatalf("only %d of %d casts landed", landed, casts)
	}
	return float64(m.Crits) / float64(landed)
}

// Client 1.60.1.70170 fixed Shadow Word: Death taking Early Demise's crit at any target health. The
// talent states 30 at rank 2 for targets at or below 20%, and the sim reads the execute phase.
func TestEarlyDemiseOnlyBelowItsHealthThreshold(t *testing.T) {
	early := spellData.EarlyDemise
	if got := early.EffectAt(1).ValueAt(2); got != 30 {
		t.Fatalf("Early Demise rank 2 crit = %v, want 30", got)
	}
	if got := early.EffectAt(2).ValueAt(2); got != 20 {
		t.Fatalf("Early Demise health threshold = %v, want 20 (the 20%% execute phase the sim reads)", got)
	}

	outside := shadowWordDeathCritShare(t, false)
	inside := shadowWordDeathCritShare(t, true)
	t.Logf("Shadow Word: Death crit share %.3f above 20%% health, %.3f at or below", outside, inside)
	if inside-outside < 0.22 || inside-outside > 0.38 {
		t.Errorf("crit share %.3f above 20%% health and %.3f at or below it: want a rise of about 0.30", outside, inside)
	}
}

// Shadow Weaving 15258 gained the Always Hit attribute: the stack no longer rolls to land. The sim
// applies it to the target (upstream #777) on every landed Shadow hit the talent's proc chance picks, which is every hit
// at rank 3 (33/67/100%).
func TestShadowWeavingAppliesOnEveryHitAtRankThree(t *testing.T) {
	if got := spellData.ShadowWeaving.FractionAt(3); got != 1 {
		t.Fatalf("Shadow Weaving rank 3 chance = %v, want 1", got)
	}
	sim, character := newClient70170Sim(t, proto.Race_RaceUndead, false)
	target := character.CurrentTarget
	priest := sim.Raid.Parties[0].Players[0].(PriestAgent).GetPriest()
	aura := priest.ShadowWeavingAuras.Get(target)
	if aura == nil {
		t.Fatal("no Shadow Weaving aura")
	}
	spell := character.GetSpell(core.ActionID{SpellID: MindBlastRankMap.Highest().ID})
	if spell == nil {
		t.Fatal("Mind Blast is not registered")
	}

	stacks := int32(0)
	for range 40 {
		before := spell.SpellMetrics[target.UnitIndex]
		spell.ApplyEffects(sim, target, spell)
		after := spell.SpellMetrics[target.UnitIndex]
		if after.Hits+after.Crits > before.Hits+before.Crits {
			stacks = min(stacks+1, 5)
		}
		if got := aura.GetStacks(); got != stacks {
			t.Fatalf("Shadow Weaving at %d stacks after a landed Mind Blast, want %d", got, stacks)
		}
	}
	if stacks != 5 {
		t.Fatalf("Shadow Weaving only reached %d stacks in 40 Mind Blasts", stacks)
	}
}

// Client 1.60.1.70170 takes the periodic spells off Inner Focus's crit and puts the channels on:
// "non-periodic" in its tooltip, and Mind Flay and Starshards do not count as periodic.
func TestInnerFocusCritIsNonPeriodic(t *testing.T) {
	sim, character := newClient70170Sim(t, proto.Race_RaceNightElf, false)
	priest := sim.Raid.Parties[0].Players[0].(PriestAgent).GetPriest()
	if priest.InnerFocusAura == nil {
		t.Fatal("no Inner Focus aura")
	}

	type entry struct {
		name string
		id   int32
		want bool
	}
	entries := []entry{
		{"Mind Blast", MindBlastRankMap.Highest().ID, true},
		{"Mind Flay", MindFlayRankMap.Highest().ID, true},
		{"Starshards", StarshardsRankMap.Highest().ID, true},
		{"Smite", SmiteRankMap.Highest().ID, true},
		{"Holy Fire", HolyFireRankMap.Highest().ID, true},
		{"Shadow Word: Pain", ShadowWordPainRankMap.Highest().ID, false},
		{"Devouring Plague", DevouringPlagueRankMap.Highest().ID, false},
		{"Shadow Word: Death", ShadowWordDeathRankMap.Highest().ID, false},
	}

	before := make([]float64, len(entries))
	spells := make([]*core.Spell, len(entries))
	for i, e := range entries {
		spells[i] = character.GetSpell(core.ActionID{SpellID: e.id})
		if spells[i] == nil {
			t.Fatalf("%s %d is not registered", e.name, e.id)
		}
		before[i] = spells[i].BonusCritPercent
	}

	priest.InnerFocusAura.Activate(sim)
	for i, e := range entries {
		got := spells[i].BonusCritPercent - before[i]
		if e.want && got != 25 {
			t.Errorf("%s: Inner Focus adds %v crit, want 25", e.name, got)
		}
		if !e.want && got != 0 {
			t.Errorf("%s: Inner Focus adds %v crit, want none", e.name, got)
		}
	}
}

// The row the sim reads for a Devouring Plague tick's crit: every rank carries Periodic Can Crit
// since client 1.60.1.70170, so priestTickOutcome(rank.PeriodicCanCrit()) is patch 14's unconditional roll.
func TestEveryDevouringPlagueRankCarriesPeriodicCanCrit(t *testing.T) {
	DevouringPlagueRankMap.Each(func(rank int32, s *spelldata.Spell) {
		if !s.PeriodicCanCrit() {
			t.Errorf("Devouring Plague rank %d (%d) lacks Periodic Can Crit; the tick would stop rolling a crit", rank, s.ID)
		}
	})
}
