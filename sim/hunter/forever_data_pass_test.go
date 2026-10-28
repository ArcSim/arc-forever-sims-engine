package hunter

import (
	"math"
	"strings"
	"testing"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
	"google.golang.org/protobuf/encoding/protojson"
)

// Aspect of the Hawk rank 7 and Serpent Sting rank 9 come only from Ruins of Ahn'Qiraj books: a level 60 hunter knows
// the trainer's rank 6 (+55 ranged attack power) and rank 8 (83 a tick).
func TestBookRanksAreNotKnown(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true})
	if hunter.SerpentSting == nil || hunter.SerpentSting.ActionID.SpellID != 13555 {
		t.Errorf("Serpent Sting is %v, want rank 8 (13555)", hunter.SerpentSting.ActionID)
	}
	for _, id := range []int32{25295, 25296} {
		if hunter.GetSpell(core.ActionID{SpellID: id}) != nil {
			t.Errorf("the hunter knows %d, a Ruins of Ahn'Qiraj book rank", id)
		}
	}
	if hunter.GetSpell(core.ActionID{SpellID: 14322}) == nil {
		t.Fatal("Aspect of the Hawk rank 6 (14322) is not known")
	}
	hawk := hunter.GetAura("Aspect of the Hawk6")
	if hawk == nil {
		t.Fatal("no Aspect of the Hawk rank 6 aura")
	}
	before := hunter.GetStat(stats.RangedAttackPower)
	hawk.Activate(sim)
	if d := hunter.GetStat(stats.RangedAttackPower) - before; math.Abs(d-55) > 1e-6 {
		t.Errorf("Aspect of the Hawk rank 6 adds %.1f ranged attack power, want 55", d)
	}
}

// Rotations saved with the book ranks, or with Trueshot Aura's rank 4, still cast: they get the ranks the hunter knows.
func TestRotationBookRanksMigrate(t *testing.T) {
	const aplJSON = `{"type":"TypeAPL","prepullActions":[{"action":{"castSpell":{"spellId":{"spellId":20905}}},"doAtValue":{"const":{"val":"-3s"}}}],` +
		`"priorityList":[{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":25295,"rank":9}}},"rhs":{"const":{"val":"1s"}}}},` +
		`"castSpell":{"spellId":{"spellId":25295,"rank":9}}}},{"action":{"castSpell":{"spellId":{"spellId":25296,"rank":7}}}}]}`
	rot := &proto.APLRotation{}
	if err := protojson.Unmarshal([]byte(aplJSON), rot); err != nil {
		t.Fatalf("rotation: %v", err)
	}
	migrateBookRanks(rot)
	got, _ := protojson.Marshal(rot)
	for _, old := range []string{"25295", "25296", "20905"} {
		if strings.Contains(string(got), old) {
			t.Errorf("the rotation still names %s after the migration: %s", old, got)
		}
	}

	if err := protojson.Unmarshal([]byte(aplJSON), rot); err != nil {
		t.Fatalf("rotation: %v", err)
	}
	player := core.WithSpec(&proto.Player{
		Race:               proto.Race_RaceHuman,
		Class:              proto.Class_ClassHunter,
		Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
		TalentsString:      "-0053150011-500230230250222151",
		Rotation:           rot,
		DistanceFromTarget: 12,
	}, abOptions)
	raid := core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{})
	stats := core.ComputeStats(&proto.ComputeStatsRequest{Raid: raid, Encounter: abEncounter()}).RaidStats.Parties[0].Players[0].RotationStats
	var warnings []string
	for _, a := range append(stats.PriorityList, stats.PrepullActions...) {
		warnings = append(warnings, a.Warnings...)
	}
	if len(warnings) > 0 {
		t.Errorf("a rotation saved with the book ranks warns: %v", warnings)
	}
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       raid,
		Encounter:  abEncounter(),
		SimOptions: &proto.SimOptions{Iterations: 5, RandomSeed: 20261028},
	})
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
	}
	casts := map[int32]int32{}
	for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
		for _, tm := range a.Targets {
			casts[a.Id.GetSpellId()] += tm.Casts
		}
	}
	if casts[13555] == 0 || casts[20906] != 5 {
		t.Errorf("Serpent Sting rank 8 cast %d times, Trueshot Aura rank 5 %d; want some, and once a fight", casts[13555], casts[20906])
	}
}

// Efficiency lowers Strider Kick's cost too (launch realms' Efficiency 19416 covers Strider Kick's spell family flag).
func TestEfficiencyCoversStriderKick(t *testing.T) {
	cost := func(talents string) int32 {
		_, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, talents: talents})
		if hunter.StriderKick == nil {
			t.Fatal("no Strider Kick")
		}
		return hunter.StriderKick.Cost.Multiplier
	}
	// patchTestTalents with Efficiency 0 / 5 (Resourcefulness 2/2 in both)
	none, five := cost("-00530500101-500230230250222151"), cost("-00535500101-500230230250222151")
	if none-five != 15 {
		t.Errorf("Efficiency 5/5 lowers Strider Kick's cost by %d%%, want 15%%", none-five)
	}
}

// Volley: Forever's 70 / 91 / 112 a wave, and no cooldown.
func TestVolleyForeverValues(t *testing.T) {
	if volleyWaveDamage != [4]float64{0, 70, 91, 112} {
		t.Errorf("Volley waves %v, want 70 / 91 / 112", volleyWaveDamage[1:])
	}
	_, hunter := newPatchTestHunter(t, patchTestSetup{gear: true})
	if hunter.Volley == nil || hunter.Volley.ActionID.SpellID != 14295 {
		t.Fatal("Volley rank 3 not known")
	}
	if hunter.Volley.CD.Timer != nil {
		t.Errorf("Volley has a %v cooldown, want none", hunter.Volley.CD.Duration)
	}
}
