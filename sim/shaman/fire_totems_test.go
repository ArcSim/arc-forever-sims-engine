package shaman

import (
	"math"
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/buffs"
)

// Flametongue Totem rank 4 (client 1.60.1.70170): 275 mana, a 1 sec global, 5 min. Its proc 16389 is one dummy of
// 1363 that does not scale with level, read as hundredths of damage per second of weapon speed held to 1.3 to 4.0
// (the tooltip's "(1363 / 77 - 1) to (1363 / 25)"), with no spell power coefficient.
func TestFlametongueTotemRank4(t *testing.T) {
	if flametongueTotemRank.ID != 16387 || flametongueTotemRank.Cost() != 275 || flametongueTotemRank.GCD() != time.Second ||
		flametongueTotemRank.Duration() != 5*time.Minute {
		t.Fatalf("Flametongue Totem is %d at %v mana, %v global, %v; want 16387 at 275, 1s, 5m0s",
			flametongueTotemRank.ID, flametongueTotemRank.Cost(), flametongueTotemRank.GCD(), flametongueTotemRank.Duration())
	}
	proc := spellData.FlametongueTotemTriggered.ByID(16389)
	if got := proc.EffectN(1).Average(core.CharacterLevel); got != 1363 {
		t.Fatalf("the totem's proc is %v at level 60, want 1363", got)
	}
	if proc.EffectN(1).Coeff() != 0 {
		t.Fatalf("the totem's proc has a spell power coefficient of %v, want none", proc.EffectN(1).Coeff())
	}
}

func TestFlametongueTotemBaseDamage(t *testing.T) {
	for _, row := range []struct {
		speed, want float64
	}{
		{1.3, 17.719}, // the tooltip's low end, 1363 / 77 to the nearest tenth of a point
		{1.0, 17.719}, // held to the floor
		{1.5, 20.445}, // a dagger
		{2.6, 35.438}, // a one-hand mace
		{3.8, 51.794}, // Arcanite Reaper
		{4.0, 54.52},  // the tooltip's high end, 1363 / 25
		{4.5, 54.52},  // held to the cap
		{0.0, 17.719}, // no weapon speed is not a weapon, which the hit checks first
		{3.0, 40.89},  // a staff
		{2.0, 27.26},  // a fast one-hander
	} {
		if got := buffs.FlametongueTotemBaseDamage(row.speed); math.Abs(got-row.want) > 1e-9 {
			t.Errorf("a %v speed weapon hits for %v, want %v", row.speed, got, row.want)
		}
	}
}

// Rank 5 hits through 408428, not the Era row 11307 its tooltip cites. Client 70425 cut it to 239 + 1.3 a
// level (245 at 60) at a 0.143 coefficient, the 16 October 2026 beta notes' reduction.
func TestFireNovaDamage(t *testing.T) {
	e := fireNovaDamage.DamageEffect()
	if avg := e.Average(core.CharacterLevel); math.Abs(avg-245) > 0.01 || math.Abs(e.Coeff()-0.143) > 1e-6 || fireNovaCoefficient != e.Coeff() {
		t.Fatalf("Fire Nova rank 5 averages %v at %v (sim %v), want 245 at 0.143", avg, e.Coeff(), fireNovaCoefficient)
	}
}
