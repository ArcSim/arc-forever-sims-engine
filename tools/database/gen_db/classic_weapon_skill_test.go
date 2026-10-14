package main

import (
	"encoding/json"
	"maps"
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
)

func TestPlannerWeaponSkills(t *testing.T) {
	for _, tc := range []struct {
		name  string
		stats map[string]json.RawMessage
		want  map[int]float64
	}{
		{"Aged Core Leather Gloves", map[string]json.RawMessage{"skillBuff": json.RawMessage(`{"173":5}`)},
			map[int]float64{int(proto.PseudoStat_PseudoStatDaggersSkill): 5}},
		{"Edgemaster's Handguards", map[string]json.RawMessage{"skillBuff": json.RawMessage(`{"44":7,"173":7,"43":7}`)},
			map[int]float64{int(proto.PseudoStat_PseudoStatAxesSkill): 7, int(proto.PseudoStat_PseudoStatDaggersSkill): 7, int(proto.PseudoStat_PseudoStatSwordsSkill): 7}},
		{"Devilsaur Claws: unarmed", map[string]json.RawMessage{"skillBuff": json.RawMessage(`{"162":4}`)},
			map[int]float64{int(proto.PseudoStat_PseudoStatUnarmedSkill): 4}},
		{"Dark Leather Gloves: Lockpicking is not a weapon", map[string]json.RawMessage{"skillBuff": json.RawMessage(`{"633":6}`)},
			map[int]float64{}},
		{"no skillBuff", map[string]json.RawMessage{"str": json.RawMessage(`15`)}, nil},
	} {
		if got := plannerWeaponSkills("1", tc.stats); !maps.Equal(got, tc.want) {
			t.Errorf("%s: %v, want %v", tc.name, got, tc.want)
		}
	}
}
