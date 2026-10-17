package warlock

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// A Destruction build with Incinerate (412758) talented.
const incinerateTalents = "2502-0050203001-1550305103101051"

// incinerateMeanHit is Incinerate's mean landed hit, never a crit or a partial resist.
func incinerateMeanHit(t *testing.T, immolated bool, spellPower float64) float64 {
	t.Helper()
	sim, warlock := newClient70170Warlock(t, proto.Race_RaceHuman, incinerateTalents)
	if warlock.Incinerate == nil {
		t.Fatal("incinerateTalents no longer take Incinerate")
	}
	warlock.AddStatDynamic(sim, stats.SpellDamage, spellPower)
	spell := warlock.Incinerate
	spell.BonusCritPercent -= 1000
	spell.BonusHitPercent += 100
	spell.Flags |= core.SpellFlagIgnoreResists
	target := warlock.CurrentTarget
	if immolated {
		warlock.Immolate.Dot(target).Apply(sim)
	}
	for i := 0; i < 2000; i++ {
		spell.SkipCastAndApplyEffects(sim, target)
		for j := 0; j < 20 && !sim.Step(); j++ {
		}
		if immolated && !warlock.Immolate.Dot(target).IsActive() {
			warlock.Immolate.Dot(target).Apply(sim)
		}
	}
	m := spell.SpellMetrics[target.UnitIndex]
	return m.TotalDamage / float64(m.Hits)
}

// Patch 155: Immolate on the target raises Incinerate's whole hit by 25%, spell power included.
func TestIncinerateImmolateBonusOnWholeHit(t *testing.T) {
	defer func(old bool) { IncinerateImmolateBonusOnWholeHit = old }(IncinerateImmolateBonusOnWholeHit)
	for _, whole := range []bool{true, false} {
		IncinerateImmolateBonusOnWholeHit = whole
		plain := incinerateMeanHit(t, false, 1000)
		burning := incinerateMeanHit(t, true, 1000)
		plainBase := incinerateMeanHit(t, false, 0)
		t.Logf("whole=%v: plain %.1f, immolated %.1f, base %.1f", whole, plain, burning, plainBase)
		if whole {
			if got := burning / plain; math.Abs(got-1.25) > 0.01 {
				t.Errorf("whole hit: Immolate raises Incinerate x%.4f at 1000 spell power, want x1.25", got)
			}
			continue
		}
		// Base roll only: the bonus is a quarter of the base hit, whatever the spell power.
		if got, want := burning-plain, 0.25*plainBase; math.Abs(got-want) > 0.03*want {
			t.Errorf("base only: Immolate adds %.1f at 1000 spell power, want about %.1f", got, want)
		}
	}
}
