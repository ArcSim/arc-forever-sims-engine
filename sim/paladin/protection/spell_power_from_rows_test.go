package protection

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/paladin"
)

// MythicSim patches 161, 162, 168 and 170: with JudgementOfRighteousnessSpellPowerFromRow,
// HolyStrikeSpellPowerFromRow, SealOfFurySpellPowerFromRows and HolyShieldSpellPowerFromRow on, every rank
// of those spells takes its client row's share (0.429 for the judgements and Holy Strike, 0.10 for the
// Seal of Fury hit, 0.08 for Holy Shield's per-block damage in client 70450); off, patch 143's measured
// 0.70, 0.286, 0.67, 0.15 and 0.07. Patch 168 turns the Judgement of Righteousness and Seal of Fury
// switches on (Lightclub's 16 October test), patch 170 Holy Strike's and Holy Shield's (official rows).
func TestPaladinSpellPowerFromRows(t *testing.T) {
	if !paladin.JudgementOfRighteousnessSpellPowerFromRow || !paladin.SealOfFurySpellPowerFromRows || !paladin.HolyStrikeSpellPowerFromRow || !paladin.HolyShieldSpellPowerFromRow {
		t.Fatal("row switch defaults moved; patch 170 has every row switch on")
	}
	cases := []struct {
		name     string
		mask     int64
		tag      int32
		sw       *bool
		measured float64
		row      float64
	}{
		{"Judgement of Righteousness", paladin.SpellMaskJudgementOfRighteousness, 0, &paladin.JudgementOfRighteousnessSpellPowerFromRow, 0.70, 0.429},
		{"Holy Strike", paladin.SpellMaskHolyStrike, 0, &paladin.HolyStrikeSpellPowerFromRow, 0.286, 0.429},
		{"Judgement of Fury", paladin.SpellMaskJudgementOfFury, 0, &paladin.SealOfFurySpellPowerFromRows, 0.67, 0.429},
		{"Seal of Fury hit", paladin.SpellMaskSealOfFuryProc, 0, &paladin.SealOfFurySpellPowerFromRows, 0.15, 0.10},
		{"Holy Shield per block", paladin.SpellMaskHolyShieldProc, 2, &paladin.HolyShieldSpellPowerFromRow, 0.07, 0.08},
	}
	for _, c := range cases {
		for _, on := range []bool{false, true} {
			want := c.measured
			if on {
				want = c.row
			}
			func() {
				old := *c.sw
				*c.sw = on
				defer func() { *c.sw = old }()
				pal := newProtPaladinWith(t, WeaponAndShield)
				ranks := 0
				for _, spell := range pal.Spellbook {
					if !spell.Matches(c.mask) || spell.ActionID.Tag != c.tag {
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
