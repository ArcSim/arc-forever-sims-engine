package hunter

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

// Savage Strikes (19159) and Predator's Edge (1310627) name the Lacerating Strikes bleed (1310536) in
// their class masks, so its ticks get the same crit chance and crit damage as Mongoose Bite.
func TestLaceratingStrikesGetsSavageStrikesAndPredatorsEdge(t *testing.T) {
	player := &proto.Player{
		Name: "sv", Class: proto.Class_ClassHunter, Race: proto.Race_RaceOrc, TalentsString: SurvivalTalents,
		Equipment: WeaponsOnly,
		Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{ClassOptions: &proto.HunterOptions{
			Ammo: proto.HunterOptions_Doomshot, PetType: proto.HunterOptions_Cat}}}},
	}
	raid := &proto.Raid{Parties: []*proto.Party{{Players: []*proto.Player{player}, Buffs: &proto.PartyBuffs{}}}, Buffs: &proto.RaidBuffs{}, Debuffs: &proto.Debuffs{}, NumActiveParties: 1}
	env, _, _ := core.NewEnvironment(raid, core.MakeSingleTargetEncounter(0), false, true)
	hunter := env.Raid.Parties[0].Players[0].(HunterAgent).GetHunter()

	bite, bleed := hunter.MongooseBite, hunter.LaceratingStrikes
	if bleed.BonusCritPercent == 0 || bleed.BonusCritPercent != bite.BonusCritPercent {
		t.Errorf("bleed bonus crit %v%%, Mongoose Bite %v%%", bleed.BonusCritPercent, bite.BonusCritPercent)
	}
	if bleed.CritMultiplierAdditive == 0 || bleed.CritMultiplierAdditive != bite.CritMultiplierAdditive {
		t.Errorf("bleed crit damage bonus %v, Mongoose Bite %v", bleed.CritMultiplierAdditive, bite.CritMultiplierAdditive)
	}
}

// A Mongoose Bite on a running bleed rolls what the bleed still owes into the new one and keeps its
// tick timer, as Forever's Deep Wounds does (patch 132): the 21 sec start over from the bite, the next
// tick lands when it was due, and the owed damage joins the new bite's 40%, spread over the ticks that
// fit before the bleed runs out. Switched off, the new bite's 40% replaces the bleed (patch 87).
func TestLaceratingStrikesRefreshRollsOver(t *testing.T) {
	for _, rollsOver := range []bool{true, false} {
		t.Run(fmt.Sprint("rolls over ", rollsOver), func(t *testing.T) {
			defer func(old bool) { LaceratingStrikesRollsOver = old }(LaceratingStrikesRollsOver)
			LaceratingStrikesRollsOver = rollsOver

			sim := core.NewSim(&proto.RaidSimRequest{
				SimOptions: &proto.SimOptions{RandomSeed: 1, Iterations: 1},
				Raid: core.SinglePlayerRaidProto(&proto.Player{
					Name: "sv", Class: proto.Class_ClassHunter, Race: proto.Race_RaceOrc, TalentsString: SurvivalMeleeTalents,
					Equipment: WeaponsOnly,
					Rotation:  &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
					Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{ClassOptions: &proto.HunterOptions{
						Ammo: proto.HunterOptions_Doomshot, PetType: proto.HunterOptions_PetNone}}}},
				}, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
				Encounter: core.MakeSingleTargetEncounter(0),
			}, simsignals.CreateSignals())
			sim.Reset()
			hunter := sim.Raid.Parties[0].Players[0].(HunterAgent).GetHunter()
			if hunter.LaceratingStrikes == nil {
				t.Fatal("the melee Survival build no longer takes Lacerating Strikes")
			}
			hunter.AutoAttacks.CancelAutoSwing(sim)
			target := hunter.CurrentTarget
			dot := hunter.LaceratingStrikes.Dot(target)

			hunter.procLaceratingStrikes(sim, &core.SpellResult{Target: target, Damage: 1000})
			firstTick := dot.SnapshotBaseDamage
			if want := 400.0 / 7; !core.WithinToleranceFloat64(want, firstTick, 1e-9) || dot.RemainingTicks() != 7 || dot.SnapshotAttackerMultiplier != 1 {
				t.Fatalf("fresh bleed: %d ticks of %.3f at x%v, want 7 of %.3f at x1", dot.RemainingTicks(), firstTick, dot.SnapshotAttackerMultiplier, want)
			}

			refreshAt := 4500 * time.Millisecond
			sim.AddPendingAction(&core.PendingAction{
				NextActionAt: refreshAt,
				OnAction: func(sim *core.Simulation) {
					if dot.NextTickAt() != 6*time.Second || dot.RemainingTicks() != 6 {
						t.Fatalf("before the refresh: next tick %s with %d left, want 6s with 6", dot.NextTickAt(), dot.RemainingTicks())
					}
					hunter.procLaceratingStrikes(sim, &core.SpellResult{Target: target, Damage: 1400})
					if got, want := dot.ExpiresAt(), refreshAt+21*time.Second; got != want {
						t.Errorf("refreshed bleed runs out at %s, want %s", got, want)
					}
					want := 560.0 / 7
					if rollsOver {
						// Ticks at 6, 9, ..., 24 sec fit before 25.5.
						want = (6*firstTick + 560) / 7
						if dot.NextTickAt() != 6*time.Second || dot.RemainingTicks() != 7 {
							t.Errorf("refreshed bleed: next tick %s with %d left, want 6s with 7", dot.NextTickAt(), dot.RemainingTicks())
						}
					}
					if !core.WithinToleranceFloat64(want, dot.SnapshotBaseDamage, 1e-9) || dot.SnapshotAttackerMultiplier != 1 {
						t.Errorf("refreshed bleed ticks %.3f at x%v, want %.3f at x1", dot.SnapshotBaseDamage, dot.SnapshotAttackerMultiplier, want)
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
			want := []time.Duration{3 * time.Second}
			if rollsOver {
				for at := 6 * time.Second; at <= 24*time.Second; at += 3 * time.Second {
					want = append(want, at)
				}
			} else {
				for at := 7500 * time.Millisecond; at <= 25500*time.Millisecond; at += 3 * time.Second {
					want = append(want, at)
				}
			}
			if !slices.Equal(ticks, want) {
				t.Errorf("ticks at %v, want %v", ticks, want)
			}
		})
	}
}
