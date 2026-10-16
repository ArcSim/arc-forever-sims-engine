package priest

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// 16 October 2026 beta notes: Devouring Plague 24 sec (was 60), Shadow Word: Death 18 sec (was 15),
// Vampiric Embrace 30 sec (was 60). Every rank carries the new cooldown.
func TestOct16PriestCooldowns(t *testing.T) {
	// Vampiric Embrace is shadow 47, the 12th shadow talent.
	_, priest := newEurekaPriest(t, proto.Race_RaceUndead, "--00000000001")
	if priest.VampiricEmbrace == nil {
		t.Fatal("talents did not register Vampiric Embrace")
	}
	check := func(name string, ladder spelldata.Ladder, want time.Duration) {
		t.Helper()
		ladder.Each(func(_ int32, rank *spelldata.Spell) {
			spell := priest.GetSpell(core.ActionID{SpellID: rank.ID})
			if spell == nil {
				t.Fatalf("%s rank %d not registered", name, rank.RankNumber())
			}
			if spell.CD.Duration != want {
				t.Errorf("%s rank %d cooldown %v, want %v", name, rank.RankNumber(), spell.CD.Duration, want)
			}
		})
	}
	check("Devouring Plague", spellData.DevouringPlague, 24*time.Second)
	check("Shadow Word: Death", spellData.ShadowWordDeath, 18*time.Second)
	if got := priest.VampiricEmbrace.CD.Duration; got != 30*time.Second {
		t.Errorf("Vampiric Embrace cooldown %v, want 30s", got)
	}
}

// 16 October 2026 beta notes: Vampiric Embrace can no longer be resisted. Against a +3 target a spell
// hit roll misses often enough that 200 casts would show it.
func TestVampiricEmbraceAlwaysLands(t *testing.T) {
	sim, priest := newEurekaPriest(t, proto.Race_RaceUndead, "--00000000001")
	target := priest.CurrentTarget
	for i := 0; i < 200; i++ {
		priest.VampiricEmbrace.SkipCastAndApplyEffects(sim, target)
		veAura := target.GetAura("Vampiric Embrace - " + target.Label)
		if veAura == nil || !veAura.IsActive() {
			t.Fatalf("Vampiric Embrace missed on cast %d", i+1)
		}
		veAura.Deactivate(sim)
	}
}

// 16 October 2026 beta notes: Power in Light also counts the priest's Shadow Word: Pain on the target,
// not only Holy Fire.
func TestPowerInLightCountsShadowWordPain(t *testing.T) {
	// Power in Light 1 and Penance.
	sim, priest := newEurekaPriest(t, proto.Race_RaceUndead, "100000000000001")
	target := priest.CurrentTarget
	smite := priest.GetSpell(core.ActionID{SpellID: spellData.Smite.Highest().ID})
	swp := priest.ShadowWordPain[len(priest.ShadowWordPain)-1]

	// Level based partial resists only take damage away, so the largest of many hits is the unresisted one.
	hit := func() float64 {
		best := 0.0
		for i := 0; i < 200; i++ {
			best = max(best, smite.CalcDamage(sim, target, 100, smite.OutcomeAlwaysHit).Damage)
		}
		return best
	}
	without := hit()
	swp.Dot(target).Apply(sim)
	with := hit()
	want := without * spellData.PowerInLight.MultiplierAt(1)
	if !core.WithinToleranceFloat64(want, with, 1e-9) {
		t.Errorf("Smite hits for %v under Shadow Word: Pain, %v without; want %v", with, without, want)
	}
}

// 16 October 2026 beta notes: Twin Disciplines also covers Shadow Word: Death, Mind Flay, Starshards
// and Penance.
func TestTwinDisciplinesCoversTheNotesSpells(t *testing.T) {
	// Penance (disc 15) and Mind Flay (shadow 9), with and without Twin Disciplines 1 (disc 3). Night Elf
	// for Starshards.
	_, plain := newEurekaPriest(t, proto.Race_RaceNightElf, "000000000000001--00000001")
	_, twin := newEurekaPriest(t, proto.Race_RaceNightElf, "001000000000001--00000001")
	bonus := spellData.TwinDisciplines.EffectAt(1).FractionAt(1)
	if bonus <= 0 {
		t.Fatal("Twin Disciplines 1 states no bonus")
	}
	for _, c := range []struct {
		name   string
		ladder spelldata.Ladder
	}{
		{"Shadow Word: Death", spellData.ShadowWordDeath},
		{"Mind Flay", spellData.MindFlay},
		{"Starshards", spellData.Starshards},
		{"Penance", spellData.Penance},
	} {
		id := core.ActionID{SpellID: c.ladder.Highest().ID}
		a, b := plain.GetSpell(id), twin.GetSpell(id)
		if a == nil || b == nil {
			t.Fatalf("%s not registered", c.name)
		}
		if got := b.DamageMultiplierAdditive - a.DamageMultiplierAdditive; !core.WithinToleranceFloat64(bonus, got, 1e-9) {
			t.Errorf("Twin Disciplines adds %v to %s, want %v", got, c.name, bonus)
		}
	}
}
