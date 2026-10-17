package core

import (
	"math"
	"testing"

	"github.com/wowsims/classic/sim/core/proto"
)

// One worker can register an aura mid-fight that another worker never did (an effect that registers on its first
// proc), so the workers' aura lists can differ in length and order and must be merged by id, not by position.
func TestCombineConcurrentResultsAurasByID(t *testing.T) {
	dist := func(n int32) *proto.DistributionMetrics {
		return &proto.DistributionMetrics{AggregatorData: &proto.AggregatorData{N: n}}
	}
	aura := func(spellID int32, uptime float64, n int32) *proto.AuraMetrics {
		return &proto.AuraMetrics{
			Id:               &proto.ActionID{RawId: &proto.ActionID_SpellId{SpellId: spellID}},
			UptimeSecondsAvg: uptime,
			AggregatorData:   &proto.AggregatorData{N: n, SumSq: uptime * uptime * float64(n)},
		}
	}
	result := func(n int32, auras ...*proto.AuraMetrics) *proto.RaidSimResult {
		player := &proto.UnitMetrics{
			Dps: dist(n), Dpasp: dist(n), Threat: dist(n), Dtps: dist(n), Tmi: dist(n), Hps: dist(n), Tto: dist(n),
			Auras: auras,
		}
		return &proto.RaidSimResult{
			RaidMetrics: &proto.RaidMetrics{
				Dps: dist(n), Hps: dist(n),
				Parties: []*proto.PartyMetrics{{Dps: dist(n), Hps: dist(n), Players: []*proto.UnitMetrics{player}}},
			},
			EncounterMetrics: &proto.EncounterMetrics{},
			IterationsDone:   n,
		}
	}

	// the first worker has aura 1 only; the second has 2 before 1 and one more than the first
	combined := CombineConcurrentSimResults([]*proto.RaidSimResult{
		result(100, aura(1, 10, 100)),
		result(100, aura(2, 5, 100), aura(1, 20, 100)),
	}, false)

	got := map[int32]*proto.AuraMetrics{}
	for _, a := range combined.RaidMetrics.Parties[0].Players[0].Auras {
		got[a.Id.GetSpellId()] = a
	}
	if len(got) != 2 {
		t.Fatalf("combined %d auras, want 2", len(got))
	}
	if a := got[1]; math.Abs(a.UptimeSecondsAvg-15) > 1e-9 || a.AggregatorData.N != 200 {
		t.Errorf("aura 1: uptime %.2f over %d iterations, want 15.00 over 200", a.UptimeSecondsAvg, a.AggregatorData.N)
	}
	// aura 2 had zero uptime in the first worker's 100 iterations
	if a := got[2]; math.Abs(a.UptimeSecondsAvg-2.5) > 1e-9 || a.AggregatorData.N != 200 {
		t.Errorf("aura 2: uptime %.2f over %d iterations, want 2.50 over 200", a.UptimeSecondsAvg, a.AggregatorData.N)
	}
	if a := got[2]; math.Abs(a.UptimeSecondsStdev-2.5) > 1e-9 {
		t.Errorf("aura 2: uptime spread %.3f, want 2.500 (half the fights 5 s, half 0)", a.UptimeSecondsStdev)
	}
}
