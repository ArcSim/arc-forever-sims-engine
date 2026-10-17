package hunter

import (
	"math"
	"testing"
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

// The Hawk Weave talents (Summon Hawk, no Ferocity).
const hawkTestTalents = "50230005001-005005001-500230230050021100"

// A hawk attacks every 2.57 s divided by the hunter's ranged attack speed (level 30 logs: 2.27 s with a 13%
// quiver, 2.23 s with a 15% one; 40% faster under Rapid Fire, 30% under Quick Shots, unchanged under the melee
// proc Quick Strikes).
func TestHawkSwingInterval(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{talents: hawkTestTalents})
	want := func(rangedSpeed float64) time.Duration {
		base := 2570 * time.Millisecond
		return time.Duration(float64(base) / rangedSpeed)
	}
	check := func(label string, rangedSpeed float64) {
		if got := hunter.hawkSwingInterval(); got != want(rangedSpeed) {
			t.Errorf("%s: hawk attacks every %v, want %v", label, got, want(rangedSpeed))
		}
	}
	check("no quiver", 1)
	hunter.PseudoStats.RangedSpeedMultiplier *= 1.15
	check("15% quiver", 1.15)
	hunter.RapidFireAura.Activate(sim)
	check("15% quiver, Rapid Fire", 1.15*1.4)
	hunter.RapidFireAura.Deactivate(sim)
	hunter.QuickShotsAura.Activate(sim)
	check("15% quiver, Quick Shots", 1.15*1.3)
	hunter.QuickShotsAura.Deactivate(sim)
	strikes := hunter.GetAura("Quick Strikes Beast")
	if strikes == nil {
		t.Fatal("no Quick Strikes aura")
	}
	if strikes.ActionID != (core.ActionID{SpellID: 1299448}) {
		t.Errorf("Quick Strikes uses %v, want the Forever spell 1299448", strikes.ActionID)
	}
	strikes.Activate(sim)
	check("15% quiver, Quick Strikes (melee only)", 1.15)
}

// Each hawk attack: about 36 at level 60 (the level 30 hawk's ~14.7 base x2.45, ExpectedStat creature melee damage
// from level 30 to 60) plus 0.5% of ranged attack power (level 30 logs: +0.8 per attack for +109 ranged AP), 91-109%
// spread. The dive: rank damage + 5% of ranged AP, Unleashed Fury a flat +3 per point (Forever data; logs: 66 raw at
// 373 ranged AP, 72 at 482).
func TestHawkDamage(t *testing.T) {
	for _, tc := range []struct {
		roll, rap, want float64
	}{{0, 0, 36 * 0.91}, {1, 0, 36 * 1.09}, {0.5, 1000, 36 + 5}} {
		if got := hawkSwingBase(tc.roll, tc.rap); math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("hawk attack at roll %.1f, ranged AP %.0f: %.3f, want %.3f", tc.roll, tc.rap, got, tc.want)
		}
	}
	if got, want := hawkDiveBase(32, 5, 373), 32.0+15+0.05*373; math.Abs(got-want) > 1e-9 {
		t.Errorf("rank 1 dive with Unleashed Fury 5/5 at 373 ranged AP: %.2f, want %.2f", got, want)
	}
}

// The hawk's attacks resolve like a pet's white melee, by the target's level: against the level 63 dummy about 40%
// glance, some are dodged or missed, and only about 1% crit (logs: 1 crit in 130 attacks), not the hunter's own
// crit chance.
func TestHawkSwingOutcomes(t *testing.T) {
	rotation := &proto.APLRotation{}
	if err := protojson.Unmarshal([]byte(`{"type":"TypeAPL","priorityList":[{"action":{"castSpell":{"spellId":{"spellId":1293241}}}}]}`), rotation); err != nil {
		t.Fatal(err)
	}
	player := core.WithSpec(&proto.Player{
		Race:               proto.Race_RaceHuman,
		Class:              proto.Class_ClassHunter,
		Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
		TalentsString:      hawkTestTalents,
		Rotation:           rotation,
		DistanceFromTarget: 3,
	}, abOptions)
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  abEncounter(),
		SimOptions: &proto.SimOptions{Iterations: 50, RandomSeed: 20261012},
	})
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
	}
	for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
		if a.Id.GetSpellId() != 1293241 || a.Id.Tag != 1 {
			continue
		}
		var hits, crits, glances, dodges, misses int32
		for _, tm := range a.Targets {
			hits += tm.Hits
			crits += tm.Crits
			glances += tm.Glances
			dodges += tm.Dodges
			misses += tm.Misses
		}
		total := float64(hits + crits + glances + dodges + misses)
		if total < 500 {
			t.Fatalf("only %.0f hawk attacks", total)
		}
		if rate := float64(crits) / total; rate > 0.03 {
			t.Errorf("hawk attacks crit %.1f%% of the time, want about 1%%", rate*100)
		}
		if rate := float64(glances) / total; rate < 0.3 || rate > 0.5 {
			t.Errorf("hawk attacks glance %.1f%% of the time against the level 63 dummy, want about 40%%", rate*100)
		}
		if dodges == 0 || misses == 0 {
			t.Errorf("hawk attacks: %d dodged, %d missed; want some of each", dodges, misses)
		}
		return
	}
	t.Fatal("no hawk attacks in the results")
}
