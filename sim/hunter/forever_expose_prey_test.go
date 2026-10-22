package hunter

import (
	"testing"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	googleProto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

// Expose Prey only opens Mongoose Bite on a target carrying a Hunter's Mark (any hunter's): with no mark up and none
// cast, a melee hunter with the talent never bites; with the mark it does.
func TestExposePreyNeedsHuntersMark(t *testing.T) {
	rotation := &proto.APLRotation{}
	if err := protojson.Unmarshal([]byte(`{"type":"TypeAPL","priorityList":[`+
		`{"action":{"castSpell":{"spellId":{"spellId":14271,"rank":4}}}},`+
		`{"action":{"castSpell":{"spellId":{"spellId":14266,"tag":1,"rank":8}}}}]}`), rotation); err != nil {
		t.Fatal(err)
	}
	bites := func(mark bool) float64 {
		debuffs := googleProto.Clone(abDebuffs).(*proto.Debuffs)
		debuffs.HuntersMark = proto.TristateEffect_TristateEffectMissing
		if mark {
			debuffs.HuntersMark = proto.TristateEffect_TristateEffectRegular
		}
		player := core.WithSpec(&proto.Player{
			Race:               proto.Race_RaceHuman,
			Class:              proto.Class_ClassHunter,
			Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "surv_dw").GearSet,
			TalentsString:      "-00531500113-500230230050122151",
			Consumes:           abConsumes,
			Buffs:              abIndividualBuffs,
			Rotation:           rotation,
			DistanceFromTarget: 3,
		}, abOptions)
		const iterations = 100
		res := core.RunRaidSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, debuffs),
			Encounter:  abEncounter(),
			SimOptions: &proto.SimOptions{Iterations: iterations, RandomSeed: 20261022},
		})
		if res.Error != nil {
			t.Fatalf("sim error: %s", res.Error.Message)
		}
		var n int32
		for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
			if a.Id.GetSpellId() == 14271 {
				for _, tm := range a.Targets {
					n += tm.Casts
				}
			}
		}
		return float64(n) / iterations
	}
	without, with := bites(false), bites(true)
	t.Logf("Mongoose Bites a fight: %.2f with no Hunter's Mark, %.2f with one", without, with)
	if without != 0 {
		t.Errorf("%.2f Mongoose Bites a fight with no Hunter's Mark on the target, want 0", without)
	}
	if with < 3 {
		t.Errorf("%.2f Mongoose Bites a fight with Hunter's Mark up, want Expose Prey to open some", with)
	}
}
