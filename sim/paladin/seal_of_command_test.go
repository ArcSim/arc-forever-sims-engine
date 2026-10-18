package paladin

import (
	"math"
	"testing"
)

// MythicSim patch 169: Improved Seals raises the part of a Seal of Command proc that Judgement of the
// Crusader and damage against the creature type add, as it raises the rest (Lightclub's sheet). Off,
// that part stays at the bare 0.29.
func TestSealOfCommandTargetBonusCoefficient(t *testing.T) {
	if !SealOfCommandBonusTakesImprovedSeals {
		t.Fatal("SealOfCommandBonusTakesImprovedSeals is off by default; patch 169 turns it on")
	}
	coefficient := spellData.SealOfCommandTriggered.Rank(1).EffectN(1).Coeff()
	if math.Abs(coefficient-0.29) > 1e-6 {
		t.Fatalf("the proc's coefficient is %v, want 0.29", coefficient)
	}
	for points, multiplier := range map[int32]float64{0: 1, 1: 1.05, 2: 1.10, 3: 1.15} {
		improvedSeals := spellData.ImprovedSeals.MultiplierAt(points)
		if math.Abs(improvedSeals-multiplier) > 1e-6 {
			t.Errorf("Improved Seals %d: multiplier %v, want %v", points, improvedSeals, multiplier)
		}
		if got := sealOfCommandTargetBonusCoefficient(coefficient, improvedSeals); math.Abs(got-0.29*multiplier) > 1e-6 {
			t.Errorf("Improved Seals %d: the target bonus takes %v, want %v", points, got, 0.29*multiplier)
		}
		func() {
			SealOfCommandBonusTakesImprovedSeals = false
			defer func() { SealOfCommandBonusTakesImprovedSeals = true }()
			if got := sealOfCommandTargetBonusCoefficient(coefficient, improvedSeals); math.Abs(got-0.29) > 1e-6 {
				t.Errorf("Improved Seals %d, switch off: the target bonus takes %v, want 0.29", points, got)
			}
		}()
	}
}
