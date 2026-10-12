package dps

import (
	"slices"
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	warrior_pkg "github.com/wowsims/forever/sim/warrior"
)

// Deep Wound's tick (412613) ignores caster damage modifiers (client Attributes[6] 0x20000000):
// beta logs 2705/2708 show its ticks unchanged in Defensive Stance while Thunder Clap and Rend lose 10%.
func TestDeepWoundsIgnoresStanceDamage(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(&proto.Player{
			Race: proto.Race_RaceOrc, Class: proto.Class_ClassWarrior,
			Equipment: TwoHandGear.GearSet, TalentsString: ArmsTalents, Spec: DefaultOptions,
		}, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
		SimOptions: &proto.SimOptions{RandomSeed: 1},
	}, simsignals.CreateSignals())
	sim.Reset()
	warrior := sim.Raid.Parties[0].Players[0].(*DpsWarrior)
	deepWounds := warrior.GetSpell(core.ActionID{SpellID: 412609})
	if deepWounds == nil {
		t.Fatal("Deep Wounds not registered")
	}
	warrior.PseudoStats.DamageDealtMultiplier *= 0.9 // Defensive Stance
	attackTable := warrior.AttackTables[warrior.CurrentTarget.UnitIndex]
	if got := deepWounds.AttackerDamageMultiplier(attackTable, true); got != 1 {
		t.Errorf("Deep Wounds tick multiplier %v in Defensive Stance, want 1", got)
	}
}

// A crit on a running Deep Wounds keeps its tick timer (Zirene on forever-bugs #234, patch 89): the 12
// sec start over from the crit, the next tick lands when it was due, and what the bleed still owed
// joins the new crit's share, spread over the ticks that fit before it runs out.
func TestDeepWoundsRefreshKeepsItsTickTimer(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(&proto.Player{
			Race: proto.Race_RaceOrc, Class: proto.Class_ClassWarrior,
			Equipment: TwoHandGear.GearSet, TalentsString: ArmsTalents, Spec: DefaultOptions,
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
		SimOptions: &proto.SimOptions{RandomSeed: 1},
	}, simsignals.CreateSignals())
	sim.Reset()
	warrior := sim.Raid.Parties[0].Players[0].(*DpsWarrior)
	warrior.AutoAttacks.CancelAutoSwing(sim) // no crits of its own to refresh the bleed
	target := warrior.CurrentTarget
	dot := warrior.DeepWounds.Dot(target)
	share := warrior.AutoAttacks.MH().AverageDamage() * 0.6 // 3 points

	warrior.DeepWounds.Cast(sim, target)
	firstTick := dot.SnapshotBaseDamage
	if want := share / 4; !core.WithinToleranceFloat64(want, firstTick, 1e-9) || dot.RemainingTicks() != 4 {
		t.Fatalf("fresh bleed: %d ticks of %.3f, want 4 of %.3f", dot.RemainingTicks(), firstTick, want)
	}

	refreshAt := 4500 * time.Millisecond
	sim.AddPendingAction(&core.PendingAction{
		NextActionAt: refreshAt,
		OnAction: func(sim *core.Simulation) {
			nextTick := dot.NextTickAt()
			if nextTick != 6*time.Second || dot.RemainingTicks() != 3 {
				t.Fatalf("before the refresh: next tick %s with %d left, want 6s with 3", nextTick, dot.RemainingTicks())
			}
			warrior.DeepWounds.Cast(sim, target)
			if got := dot.NextTickAt(); got != nextTick {
				t.Errorf("the refresh moved the next tick from %s to %s", nextTick, got)
			}
			if got, want := dot.ExpiresAt(), refreshAt+12*time.Second; got != want {
				t.Errorf("refreshed bleed runs out at %s, want %s", got, want)
			}
			// Ticks at 6, 9, 12 and 15 sec fit before 16.5.
			if dot.RemainingTicks() != 4 {
				t.Errorf("refreshed bleed has %d ticks left, want 4", dot.RemainingTicks())
			}
			if want := (3*firstTick + share) / 4; !core.WithinToleranceFloat64(want, dot.SnapshotBaseDamage, 1e-9) {
				t.Errorf("refreshed bleed ticks %.3f, want %.3f", dot.SnapshotBaseDamage, want)
			}
		},
	})

	var ticks []time.Duration
	left := dot.RemainingTicks()
	for steps := 0; dot.IsActive(); steps++ {
		if steps > 1_000_000 || sim.Step() {
			t.Fatal("the bleed never ran out")
		}
		if dot.RemainingTicks() < left {
			ticks = append(ticks, sim.CurrentTime)
		}
		left = dot.RemainingTicks()
	}
	want := []time.Duration{3 * time.Second, 6 * time.Second, 9 * time.Second, 12 * time.Second, 15 * time.Second}
	if !slices.Equal(ticks, want) {
		t.Errorf("ticks at %v, want %v", ticks, want)
	}
}

// An off-hand crit's Deep Wounds is a share of the off-hand weapon's average damage, without attack
// power (Forever: "Fixed a bug where Deep Wounds was scaling with attack power for your Offhand weapon"),
// and a main-hand crit's of the main hand's. DeepWoundsUsesStrikingWeapon false sizes both off the main
// hand.
func TestDeepWoundsFromAnOffHandCritUsesTheOffHand(t *testing.T) {
	bleedAfterCrit := func(t *testing.T, procMask core.ProcMask) (tick, mhAvg, ohAvg float64) {
		t.Helper()
		sim := core.NewSim(&proto.RaidSimRequest{
			Raid: core.SinglePlayerRaidProto(&proto.Player{
				Race: proto.Race_RaceHuman, Class: proto.Class_ClassWarrior,
				Equipment: DualWieldGear.GearSet, TalentsString: FuryTalents, Spec: DefaultOptions,
				Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
			}, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
			SimOptions: &proto.SimOptions{RandomSeed: 1},
		}, simsignals.CreateSignals())
		sim.Reset()
		warrior := sim.Raid.Parties[0].Players[0].(*DpsWarrior)
		trigger := warrior.GetAura("Deep Wounds - Trigger")
		if trigger == nil {
			t.Fatal("Deep Wounds trigger not registered")
		}
		target := warrior.CurrentTarget
		spell := &core.Spell{ProcMask: procMask, SpellSchool: core.SpellSchoolPhysical}
		trigger.OnSpellHitDealt(trigger, sim, spell, &core.SpellResult{Target: target, Outcome: core.OutcomeCrit, Damage: 100})
		dot := warrior.DeepWounds.Dot(target)
		if !dot.IsActive() || dot.RemainingTicks() != 4 {
			t.Fatalf("the crit left Deep Wounds active %v with %d ticks, want 4", dot.IsActive(), dot.RemainingTicks())
		}
		return dot.SnapshotBaseDamage, warrior.AutoAttacks.MH().AverageDamage(), warrior.AutoAttacks.OH().AverageDamage()
	}
	const share = 0.6 // 3 points

	tick, mhAvg, ohAvg := bleedAfterCrit(t, core.ProcMaskMeleeOHAuto)
	if mhAvg == ohAvg {
		t.Fatalf("the test gear's hands average the same, %.1f", mhAvg)
	}
	if want := ohAvg * share / 4; !core.WithinToleranceFloat64(want, tick, 1e-9) {
		t.Errorf("an off-hand auto crit ticks %.3f, want the off hand's %.3f", tick, want)
	}
	if tick, _, _ = bleedAfterCrit(t, core.ProcMaskMeleeOHSpecial); !core.WithinToleranceFloat64(ohAvg*share/4, tick, 1e-9) {
		t.Errorf("an off-hand special crit ticks %.3f, want the off hand's %.3f", tick, ohAvg*share/4)
	}
	if tick, _, _ = bleedAfterCrit(t, core.ProcMaskMeleeMHAuto); !core.WithinToleranceFloat64(mhAvg*share/4, tick, 1e-9) {
		t.Errorf("a main-hand auto crit ticks %.3f, want the main hand's %.3f", tick, mhAvg*share/4)
	}

	defer func(old bool) { warrior_pkg.DeepWoundsUsesStrikingWeapon = old }(warrior_pkg.DeepWoundsUsesStrikingWeapon)
	warrior_pkg.DeepWoundsUsesStrikingWeapon = false
	if tick, _, _ = bleedAfterCrit(t, core.ProcMaskMeleeOHAuto); !core.WithinToleranceFloat64(mhAvg*share/4, tick, 1e-9) {
		t.Errorf("switched off, an off-hand auto crit ticks %.3f, want the main hand's %.3f", tick, mhAvg*share/4)
	}
}
