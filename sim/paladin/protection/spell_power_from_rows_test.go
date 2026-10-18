package protection

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/paladin"
)

// MythicSim patches 161, 162 and 168: with JudgementOfRighteousnessSpellPowerFromRow,
// HolyStrikeSpellPowerFromRow and SealOfFurySpellPowerFromRows on, every rank of those spells takes its
// client row's share (0.429 for the judgements and Holy Strike, 0.10 for the Seal of Fury hit in client
// 70425); off, patch 143's measured 0.70, 0.286, 0.67 and 0.15. Patch 168 turns the Judgement of
// Righteousness and Seal of Fury switches on (Lightclub's 16 October test); Holy Strike's stays off.
func TestPaladinSpellPowerFromRows(t *testing.T) {
	if !paladin.JudgementOfRighteousnessSpellPowerFromRow || !paladin.SealOfFurySpellPowerFromRows || paladin.HolyStrikeSpellPowerFromRow {
		t.Fatal("row switch defaults moved; patch 168 has Judgement of Righteousness and Seal of Fury on, Holy Strike off")
	}
	cases := []struct {
		name     string
		mask     int64
		sw       *bool
		measured float64
		row      float64
	}{
		{"Judgement of Righteousness", paladin.SpellMaskJudgementOfRighteousness, &paladin.JudgementOfRighteousnessSpellPowerFromRow, 0.70, 0.429},
		{"Holy Strike", paladin.SpellMaskHolyStrike, &paladin.HolyStrikeSpellPowerFromRow, 0.286, 0.429},
		{"Judgement of Fury", paladin.SpellMaskJudgementOfFury, &paladin.SealOfFurySpellPowerFromRows, 0.67, 0.429},
		{"Seal of Fury hit", paladin.SpellMaskSealOfFuryProc, &paladin.SealOfFurySpellPowerFromRows, 0.15, 0.10},
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
