package mage

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
)

// Patch 175: Missile Barrage follows the 19 October 2026 in-game tooltip. Frost Blast procs it at half
// of Arcane Blast's chance and Arcane Missiles' own missiles no longer do; off, client 70425's list
// (missiles at a quarter, no Frost Blast). Each trial lands one hit with the buff down and records
// whether it came up.
func TestMissileBarrageHotfixedProcs(t *testing.T) {
	const trials = 4000

	var chance float64 // Arcane Blast's, at the build's rank
	rates := func(on bool) map[string]float64 {
		defer func(old bool) { MissileBarrageHotfixedProcs = old }(MissileBarrageHotfixedProcs)
		MissileBarrageHotfixedProcs = on
		sim, mage := frostfireSim(ArcaneTalents)
		if mage.Talents.MissileBarrage == 0 {
			t.Fatal("ArcaneTalents take no Missile Barrage; pick a build that does")
		}
		chance = spellData.MissileBarrage.Effect(dbcenums.A_PROC_TRIGGER_SPELL, 0).ValueAt(mage.Talents.MissileBarrage) / 100

		var missile *core.Spell
		for _, spell := range mage.Spellbook {
			if spell.Matches(MageSpellArcaneMissilesTick) && (missile == nil || spell.Rank > missile.Rank) {
				missile = spell
			}
		}
		spells := map[string]*core.Spell{
			"Arcane Blast":    mage.GetSpell(core.ActionID{SpellID: spellData.ArcaneBlast.Highest().ID}),
			"Frost Blast":     mage.GetSpell(core.ActionID{SpellID: spellData.FrostBlast.Highest().ID}),
			"Fireball":        mage.GetSpell(core.ActionID{SpellID: spellData.Fireball.Highest().ID}),
			"Arcane Missiles": missile,
		}

		out := make(map[string]float64, len(spells))
		for name, spell := range spells {
			if spell == nil {
				t.Fatalf("%s not registered", name)
			}
			procs := 0
			for i := 0; i < trials; i++ {
				mage.MissileBarrageAura.Deactivate(sim)
				spell.CalcAndDealDamage(sim, mage.CurrentTarget, 1, spell.OutcomeAlwaysHit)
				if mage.MissileBarrageAura.IsActive() {
					procs++
				}
			}
			out[name] = float64(procs) / trials
		}
		return out
	}

	near := func(got, want float64) bool { return got > want-0.03 && got < want+0.03 }

	on := rates(true)
	t.Logf("switch on: %v", on)
	if !near(on["Arcane Blast"], chance) || !near(on["Frost Blast"], chance/2) || !near(on["Fireball"], chance/2) {
		t.Errorf("switch on: Arcane Blast %.3f, Frost Blast %.3f, Fireball %.3f; want %.2f, %.2f, %.2f",
			on["Arcane Blast"], on["Frost Blast"], on["Fireball"], chance, chance/2, chance/2)
	}
	if on["Arcane Missiles"] != 0 {
		t.Errorf("switch on: an Arcane Missiles missile procs Missile Barrage %.3f of the time, want never", on["Arcane Missiles"])
	}

	off := rates(false)
	t.Logf("switch off: %v", off)
	if !near(off["Arcane Missiles"], chance/4) || off["Arcane Missiles"] == 0 || off["Frost Blast"] != 0 {
		t.Errorf("switch off: Arcane Missiles %.3f, Frost Blast %.3f; want %.3f and never", off["Arcane Missiles"], off["Frost Blast"], chance/4)
	}
}
