package protection

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/paladin"
)

// MythicSim patch 161: with JudgementOfRighteousnessSpellPowerFromRow on, every rank of the judgement
// takes its client row's share (0.429 on every rank in client 70425); off, patch 143's measured 0.70.
func TestPaladinSpellPowerFromRows(t *testing.T) {
	if paladin.JudgementOfRighteousnessSpellPowerFromRow {
		t.Fatal("the row switch is on by default; patch 161 leaves it off")
	}
	cases := []struct {
		name     string
		mask     int64
		sw       *bool
		measured float64
	}{
		{"Judgement of Righteousness", paladin.SpellMaskJudgementOfRighteousness, &paladin.JudgementOfRighteousnessSpellPowerFromRow, 0.70},
	}
	for _, c := range cases {
		for _, on := range []bool{false, true} {
			want := c.measured
			if on {
				want = 0.429
			}
			func() {
				old := *c.sw
				*c.sw = on
				defer func() { *c.sw = old }()
				pal := newProtPaladinWith(t, WeaponAndShield)
				ranks := 0
				for _, spell := range pal.Spellbook {
					if !spell.Matches(c.mask) || spell.ActionID.Tag != 0 {
						continue
					}
					ranks++
					if math.Abs(spell.BonusCoefficient-want) > 1e-6 {
						t.Errorf("%s %d with the row switch %v: coefficient %v, want %v", c.name, spell.ActionID.SpellID, on, spell.BonusCoefficient, want)
					}
				}
				if ranks < 2 {
					t.Errorf("%s: %d ranks registered", c.name, ranks)
				}
			}()
		}
	}
}
