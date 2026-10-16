package feralcat

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/druid"
)

// Client 1.60.1.70170: "Bear Form and Dire Bear Form now generate 75% increased Rage when landing a
// Critical Strike"; the Druid development notes of 8 October 2026 raise it to 100%. The Cat registers the Bear Form aura without the spell, so its Rage bar is only ever
// fed in that form; the bonus is on the same bar every Warrior and Bear uses (core.CritRageBonus) and
// applies exactly once.
func TestCatsBearFormCritAutoAttackGivesTwiceTheRage(t *testing.T) {
	player := &proto.Player{Name: "Cat", Class: proto.Class_ClassDruid, Race: proto.Race_RaceNightElf,
		TalentsString: DefaultTalents, Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{},
		Spec: DefaultSpecOptions, Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL}}
	sim := core.NewSim(&proto.RaidSimRequest{SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid:      core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter: core.MakeSingleTargetEncounter(0)}, simsignals.CreateSignals())
	sim.Reset()
	cat := sim.Raid.Parties[0].Players[0].(druid.DruidAgent).GetDruid()
	sim.CurrentTime = time.Second
	cat.ClearForm(sim)
	cat.BearFormAura.Activate(sim)
	if !cat.InForm(druid.Bear) {
		t.Fatal("not in Bear Form")
	}

	swing := func(outcome core.HitOutcome, damage float64) float64 {
		cat.ResetRageBar(sim, 0)
		cat.Unit.OnSpellHitDealt(sim, cat.AutoAttacks.MHAuto(), &core.SpellResult{
			Target: sim.Encounter.ActiveTargetUnits[0], Outcome: outcome, Damage: damage, PostArmorAndResistanceMultiplier: damage,
		})
		return cat.CurrentRage()
	}
	hit := swing(core.OutcomeHit, 1000)
	crit := swing(core.OutcomeCrit, 2000)
	glance := swing(core.OutcomeGlance, 700)

	if want := core.BaseRageHitFactor * cat.AutoAttacks.MH().SwingSpeed; !core.WithinToleranceFloat64(want, hit, 0.000001) {
		t.Errorf("hit gave %0.5f Rage, want %0.5f", hit, want)
	}
	if !core.WithinToleranceFloat64(2, crit/hit, 0.0000001) {
		t.Errorf("crit gave %0.6f times a hit in Bear Form, want 2", crit/hit)
	}
	if glance != hit {
		t.Errorf("glance gave %0.5f Rage, want the hit's %0.5f", glance, hit)
	}
}
