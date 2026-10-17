package paladin

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core/spelldata"
)

// MythicSim patch 158: each rank's seal hit takes half the judgement dummy's spell power share (x1.2 on
// a two-hander first), and the tooltip switch, when on, replaces it at every rank.
func TestSealOfRighteousnessSpellPowerFixedShares(t *testing.T) {
	if !SealOfRighteousnessSpellPowerFixed || SealOfRighteousnessPerHitTooltipCoefficient {
		t.Fatal("patch 158 has SealOfRighteousnessSpellPowerFixed on and SealOfRighteousnessPerHitTooltipCoefficient off by default")
	}
	// rank: the dummy's share, then the seal hit's with a one-hander and a two-hander.
	want := map[int32][3]float64{
		1: {0.058, 0.029, 0.0348},
		2: {0.125, 0.0625, 0.075},
		3: {0.185, 0.0925, 0.111},
		4: {0.20, 0.10, 0.12},
		5: {0.20, 0.10, 0.12},
		6: {0.20, 0.10, 0.12},
		7: {0.20, 0.10, 0.12},
		8: {0.20, 0.10, 0.12},
	}
	if got := spellData.SealOfRighteousness.Len(); got != int32(len(want)) {
		t.Fatalf("Seal of Righteousness has %d ranks, want %d", got, len(want))
	}
	near := func(a, b float64) bool { return math.Abs(a-b) < 1e-6 }
	spellData.SealOfRighteousness.Each(func(rank int32, seal *spelldata.Spell) {
		judge := spellData.JudgementOfRighteousness.ByID(int32(seal.EffectN(2).BaseValue()))
		dummy := judge.EffectN(2).Coeff()
		w := want[rank]
		if !near(dummy, w[0]) {
			t.Errorf("rank %d (%d): the judgement dummy's share is %.4f, want %.4f", rank, seal.ID, dummy, w[0])
		}
		oneHand, twoHand := SealOfRighteousnessHitCoefficient(dummy, false), SealOfRighteousnessHitCoefficient(dummy, true)
		if !near(oneHand, w[1]) || !near(twoHand, w[2]) {
			t.Errorf("rank %d: the seal hit takes %.4f (one-hander) and %.4f (two-hander), want %.4f and %.4f", rank, oneHand, twoHand, w[1], w[2])
		}
		// Switched off: the dummy's share, the pre-fix server.
		if got := sealOfRighteousnessHitCoefficient(dummy, false, false, false); !near(got, w[0]) {
			t.Errorf("rank %d, fix off: one-hander %.4f, want %.4f", rank, got, w[0])
		}
		if got := sealOfRighteousnessHitCoefficient(dummy, true, false, false); !near(got, w[0]*1.2) {
			t.Errorf("rank %d, fix off: two-hander %.4f, want %.4f", rank, got, w[0]*1.2)
		}
		// The tooltip switch comes first, with the fix on or off.
		for _, fixed := range []bool{true, false} {
			if got := sealOfRighteousnessHitCoefficient(dummy, false, true, fixed); !near(got, 0.0765) {
				t.Errorf("rank %d, tooltip on (fix %v): one-hander %.4f, want 0.0765", rank, fixed, got)
			}
			if got := sealOfRighteousnessHitCoefficient(dummy, true, true, fixed); !near(got, 0.12) {
				t.Errorf("rank %d, tooltip on (fix %v): two-hander %.4f, want 0.12", rank, fixed, got)
			}
		}
	})
}

// Rank 1 is 21084, the row with a judgement (20187), and its hit is 25741 as the logs show
// (foreverlogs 2768); 25742 is the older 20154's.
func TestSealOfRighteousnessRankOneHitID(t *testing.T) {
	rank1 := spellData.SealOfRighteousness.Rank(1)
	if rank1.ID != 21084 {
		t.Fatalf("Seal of Righteousness rank 1 is %d, want 21084", rank1.ID)
	}
	if got := int32(rank1.EffectN(2).BaseValue()); got != 20187 {
		t.Errorf("rank 1 judges as %d, want 20187", got)
	}
	if got := sealOfRighteousnessProcIDs[1]; got != 25741 {
		t.Errorf("rank 1's hit is %d, want 25741", got)
	}
	want := map[int32]int32{1: 25741, 2: 25740, 3: 25739, 4: 25738, 5: 25737, 6: 25736, 7: 25735, 8: 25713}
	for rank, id := range want {
		if got := sealOfRighteousnessProcIDs[rank]; got != id {
			t.Errorf("rank %d's hit is %d, want %d", rank, got, id)
		}
	}
}
