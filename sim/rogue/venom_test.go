package rogue

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

// Venom (1310703) is a finisher: while its aura is up, Instant Poison's damage (SPELLMOD_DAMAGE +30 on
// its family flag) and Deadly Poison's periodic damage (SPELLMOD_DOT +30 on 0x10000) rise by 30%, and
// the poisons' application chance by 10 points. Talenting it does nothing until it is cast: a custom
// rotation without a Venom line never raises the poisons (the MythicSim report on sim 67731164,
// 17 October 2026). Deadly Poison's ticks read the multiplier as they tick, so a stack applied before
// Venom ticks harder once Venom is up.
func TestVenomRaisesPoisonsWhileUp(t *testing.T) {
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceHuman,
		Class:         proto.Class_ClassRogue,
		Equipment:     daggersOnly(),
		Consumables:   &proto.ConsumesSpec{MhImbueId: instantImbueID, OhImbueId: deadlyImbueID},
		TalentsString: "00532310551511051-302303-002", // MythicSim's Assassination reference, Venom 1/1
		Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, DefaultOptions)
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	rogue := sim.Raid.Parties[0].Players[0].(RogueAgent).GetRogue()
	if rogue.Venom == nil {
		t.Fatal("the reference talents register no Venom")
	}
	at := rogue.AttackTables[rogue.CurrentTarget.UnitIndex]

	deadlyAdditive := rogue.DeadlyPoison.DamageMultiplierAdditive
	deadly := rogue.DeadlyPoison.AttackerDamageMultiplier(at, true)
	instant := rogue.InstantPoison.AttackerDamageMultiplier(at, false)
	chance := rogue.poisonProcChance(0.3)

	rogue.VenomAura.Activate(sim)
	if got := rogue.DeadlyPoison.DamageMultiplierAdditive - deadlyAdditive; !core.WithinToleranceFloat64(0.30, got, 1e-9) {
		t.Errorf("Venom adds %.3f to Deadly Poison's damage bucket, want 0.30", got)
	}
	if got := rogue.DeadlyPoison.AttackerDamageMultiplier(at, true) / deadly; got < 1.2 {
		t.Errorf("Deadly Poison's ticks under Venom: x%.4f", got)
	}
	if got := rogue.InstantPoison.AttackerDamageMultiplier(at, false) / instant; got < 1.2 {
		t.Errorf("Instant Poison under Venom: x%.4f", got)
	}
	if got := rogue.poisonProcChance(0.3) - chance; !core.WithinToleranceFloat64(0.10, got, 1e-9) {
		t.Errorf("Venom adds %.3f to the poison chance, want 0.10", got)
	}

	rogue.VenomAura.Deactivate(sim)
	if got := rogue.DeadlyPoison.AttackerDamageMultiplier(at, true); !core.WithinToleranceFloat64(deadly, got, 1e-9) {
		t.Errorf("Deadly Poison's tick multiplier after Venom: %.4f, want %.4f", got, deadly)
	}
}
