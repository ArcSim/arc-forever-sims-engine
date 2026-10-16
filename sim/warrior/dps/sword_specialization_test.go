package dps

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/warrior"
)

// 16 October 2026 beta notes: Weaponmaster's Sword Specialization no longer triggers from Shield Bash
// and Shield Slam. The engine strikes both with the shield (the off-hand proc mask), so they already
// missed a main-hand sword's mask; this drives the trigger with main-hand strikes carrying their class
// masks, so the exclusion holds whatever hand they are booked to. Heroic Strike is the control.
func TestSwordSpecializationSkipsShieldBashAndShieldSlam(t *testing.T) {
	const (
		swordMH = 19019 // Thunderfury, Blessed Blade of the Windseeker
		shield  = 18168 // Force Reactive Disk
	)
	sim := core.NewSim(&proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(&proto.Player{
			Race: proto.Race_RaceOrc, Class: proto.Class_ClassWarrior,
			Equipment: weaponsOnly(swordMH, shield), TalentsString: ArmsTalents, Spec: DefaultOptions,
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
		SimOptions: &proto.SimOptions{RandomSeed: 1},
	}, simsignals.CreateSignals())
	sim.Reset()
	war := sim.Raid.Parties[0].Players[0].(*DpsWarrior)
	if war.Talents.Weaponmaster == 0 || war.MainHand().WeaponType != proto.WeaponType_WeaponTypeSword ||
		war.OffHand().WeaponType != proto.WeaponType_WeaponTypeShield {
		t.Fatalf("test setup: Weaponmaster %d, main hand %v, off hand %v", war.Talents.Weaponmaster,
			war.MainHand().WeaponType, war.OffHand().WeaponType)
	}
	trigger := war.GetAura("Weaponmaster (Sword)")
	extra := war.GetSpell(war.AutoAttacks.MHConfig().ActionID.WithTag(1290261))
	if trigger == nil || extra == nil {
		t.Fatal("no Sword Specialization trigger or extra attack")
	}
	target := war.CurrentTarget

	for _, tc := range []struct {
		name  string
		mask  int64
		procs bool
	}{
		{"Shield Bash", warrior.SpellMaskShieldBash, false},
		{"Shield Slam", warrior.SpellMaskShieldSlam, false},
		{"Heroic Strike", warrior.SpellMaskHeroicStrike, true},
	} {
		strike := &core.Spell{ProcMask: core.ProcMaskMeleeMHSpecial, ClassSpellMask: tc.mask}
		before := extra.SpellMetrics[target.UnitIndex].Casts
		for i := 0; i < 2000; i++ {
			trigger.OnSpellHitDealt(trigger, sim, strike, &core.SpellResult{Target: target, Outcome: core.OutcomeHit})
		}
		if procced := extra.SpellMetrics[target.UnitIndex].Casts > before; procced != tc.procs {
			t.Errorf("%s: Sword Specialization procced %v, want %v", tc.name, procced, tc.procs)
		}
	}
}
