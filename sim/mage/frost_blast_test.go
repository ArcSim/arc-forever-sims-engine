package mage

import (
	"math"
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
)

// Frost Blast (patch 153): every rank is a castable Frost spell with the client's cast, cost and
// coefficient, and Improved Frostbolt takes the same time off it as off Frostbolt.
func TestFrostBlastIsCastable(t *testing.T) {
	_, mage := newThreeTargetMage(FrostTalents)
	if mage.Talents.ImprovedFrostbolt == 0 {
		t.Fatal("FrostTalents no longer take Improved Frostbolt; pick a build that does")
	}
	ids := []int32{1326053, 1326054, 1326055}
	for _, id := range ids {
		if mage.GetSpell(core.ActionID{SpellID: id}) == nil {
			t.Fatalf("Frost Blast %d not registered", id)
		}
	}
	blast := mage.GetSpell(core.ActionID{SpellID: spellData.FrostBlast.Highest().ID})
	if blast.ActionID.SpellID != 1326055 || !blast.Flags.Matches(core.SpellFlagAPL) || !blast.Matches(MageSpellFrost|MageSpellChill) ||
		blast.SpellSchool != core.SpellSchoolFrost {
		t.Fatalf("Frost Blast rank 3: id %d, flags %v, school %v", blast.ActionID.SpellID, blast.Flags, blast.SpellSchool)
	}
	if got := blast.DefaultCast.Cost; math.Abs(got-90) > 1e-9 {
		t.Errorf("Frost Blast costs %v, want 90", got)
	}
	if got := blast.BonusCoefficient; math.Abs(got-0.271) > 1e-6 {
		t.Errorf("Frost Blast coefficient %v, want 0.271", got)
	}
	frostbolt := mage.GetSpell(core.ActionID{SpellID: spellData.Frostbolt.Highest().ID})
	cut := spellData.Frostbolt.Highest().CastTime() - frostbolt.CastTime()
	if cut <= 0 {
		t.Fatalf("Improved Frostbolt took nothing off Frostbolt (%v)", frostbolt.CastTime())
	}
	if got, want := blast.CastTime(), 1500*time.Millisecond-cut; got != want {
		t.Errorf("Frost Blast casts in %v with Improved Frostbolt, want %v", got, want)
	}
}

func TestFrostBlastSwitch(t *testing.T) {
	defer func(old bool) { FrostBlastCastable = old }(FrostBlastCastable)
	FrostBlastCastable = false
	_, mage := newThreeTargetMage(FrostTalents)
	if mage.GetSpell(core.ActionID{SpellID: 1326055}) != nil {
		t.Fatal("Frost Blast registered with FrostBlastCastable off")
	}
}

// Frost Blast's slow is a chill: a landed Frost Blast rolls Fingers of Frost at its 15%.
func TestFrostBlastRollsFingersOfFrost(t *testing.T) {
	sim, mage := newThreeTargetMage(FrostTalents)
	if mage.FingersOfFrostAura == nil {
		t.Fatal("FrostTalents no longer take Fingers of Frost; pick a build that does")
	}
	blast := mage.GetSpell(core.ActionID{SpellID: 1326055})
	const trials = 4000
	procs := 0
	for i := 0; i < trials; i++ {
		mage.FingersOfFrostAura.Deactivate(sim)
		result := blast.CalcDamage(sim, mage.CurrentTarget, 1, blast.OutcomeAlwaysHit)
		mage.procFingersOfFrost(sim, blast, result)
		if mage.FingersOfFrostAura.IsActive() {
			procs++
		}
	}
	want := spellData.FingersOfFrost.EffectAt(2).FractionAt(mage.Talents.FingersOfFrost)
	if got := float64(procs) / trials; math.Abs(got-want) > 0.02 {
		t.Errorf("Frost Blast procced Fingers of Frost %.3f, want %.2f", got, want)
	}
}
