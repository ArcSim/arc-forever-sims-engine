package hunter

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
)

func newTrapTestHunter(t *testing.T, bonus stats.Stats) (*core.Simulation, *Hunter) {
	t.Helper()
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 7},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "sv", Class: proto.Class_ClassHunter, Race: proto.Race_RaceOrc, TalentsString: SurvivalTalents,
			Equipment: WeaponsOnly, Buffs: &proto.IndividualBuffs{},
			Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{ClassOptions: &proto.HunterOptions{
				Ammo: proto.HunterOptions_Doomshot, PetType: proto.HunterOptions_PetNone}}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	hunter := sim.Raid.Parties[0].Players[0].(HunterAgent).GetHunter()
	hunter.AddStatsDynamic(sim, bonus)
	return sim, hunter
}

// Hmm, 16 October: Immolation Trap takes 13% of the higher of ranged and melee attack power over the
// whole burn, so 1000 more attack power adds 130 across its 5 ticks, 26 a tick.
func TestImmolationTrapAttackPowerShare(t *testing.T) {
	sim, hunter := newTrapTestHunter(t, stats.Stats{})
	target := hunter.CurrentTarget
	dot := hunter.ImmolationTrap.Dot(target)
	if dot.BaseTickCount != 5 {
		t.Fatalf("Immolation Trap ticks %d, want 5", dot.BaseTickCount)
	}
	dot.Apply(sim)
	base := dot.SnapshotBaseDamage
	rap, ap := hunter.ImmolationTrap.RangedAttackPower(target), hunter.ImmolationTrap.MeleeAttackPower(target)
	dot.Deactivate(sim)

	for _, tc := range []struct {
		name  string
		bonus stats.Stats
		want  float64
	}{
		// Ranged is the higher one: only its gain counts.
		{"ranged", stats.Stats{stats.RangedAttackPower: 1000}, 0.026 * (math.Max(rap+1000, ap) - math.Max(rap, ap))},
		// Melee overtakes ranged.
		{"melee", stats.Stats{stats.AttackPower: 3000}, 0.026 * (math.Max(rap, ap+3000) - math.Max(rap, ap))},
	} {
		sim, hunter := newTrapTestHunter(t, tc.bonus)
		dot := hunter.ImmolationTrap.Dot(hunter.CurrentTarget)
		dot.Apply(sim)
		if got := dot.SnapshotBaseDamage - base; math.Abs(got-tc.want) > 1e-6 {
			t.Errorf("%s: tick base +%.4f, want +%.4f", tc.name, got, tc.want)
		}
	}
}

// Explosive Trap's explosion takes the same 13%; the burning ground it leaves takes none.
func TestExplosiveTrapAttackPowerShare(t *testing.T) {
	lowSim, low := newTrapTestHunter(t, stats.Stats{})
	highSim, high := newTrapTestHunter(t, stats.Stats{stats.RangedAttackPower: 1000, stats.AttackPower: 1000})
	damageRange := explosiveTrapRange[spellData.ExplosiveTrap.Highest().RankNumber()]

	// Same seed, so both draw the same roll.
	lowHit := explosiveTrapHitDamage(lowSim, low.ExplosiveTrap, low.CurrentTarget, damageRange)
	highHit := explosiveTrapHitDamage(highSim, high.ExplosiveTrap, high.CurrentTarget, damageRange)
	if got := highHit - lowHit; math.Abs(got-130) > 1e-6 {
		t.Errorf("explosion +%.4f for 1000 attack power, want +130", got)
	}

	lowDot, highDot := low.ExplosiveTrap.AOEDot(), high.ExplosiveTrap.AOEDot()
	lowDot.Apply(lowSim)
	highDot.Apply(highSim)
	if lowDot.SnapshotBaseDamage != highDot.SnapshotBaseDamage {
		t.Errorf("burning ground tick %.4f with 1000 more attack power, want %.4f", highDot.SnapshotBaseDamage, lowDot.SnapshotBaseDamage)
	}
}
