package core

import (
	"testing"

	"github.com/wowsims/forever/sim/core/stats"
)

// Beta development notes, 1 October (Warrior) and 8 October 2026 (Druid): "Additional Rage generated
// from landing Critical Strikes increased to 100% increased Rage (was 75%)". A critical strike with an
// auto attack pays the swing's rage again on top, and the damage it dealt does not matter: a crit that
// deals three times a hit's damage still pays the same. The extra is the swing's rage before a
// two-hander's factor (Magey, client 70170, CritRageBonusSkipsTwoHand), so a two-hander gets about 1.77
// times its swing, not 2.
func TestCritAutoAttackRageIsTwiceAHit(t *testing.T) {
	const hitDamage = 500.0

	for _, tc := range []struct {
		name     string
		prepare  func(fw *FakeRageWarrior)
		spell    func(fw *FakeRageWarrior) *Spell
		wantHit  float64
		wantCrit float64
	}{
		{
			name:    "one-hand main hand",
			prepare: func(fw *FakeRageWarrior) {},
			spell:   func(fw *FakeRageWarrior) *Spell { return fw.AutoAttacks.MHAuto() },
			// 2.6 * 3.46, then times 2
			wantHit: 8.996, wantCrit: 17.992,
		},
		{
			name: "two-hand main hand",
			prepare: func(fw *FakeRageWarrior) {
				fw.AutoAttacks.MH().NormalizedSwingSpeed = TwoHandNormalizedSwingSpeed
			},
			spell: func(fw *FakeRageWarrior) *Spell { return fw.AutoAttacks.MHAuto() },
			// 2.6 * 4.5, plus 2.6 * 3.46 on a crit
			wantHit: 11.7, wantCrit: 20.696,
		},
		{
			name:    "off hand",
			prepare: func(fw *FakeRageWarrior) {},
			spell:   func(fw *FakeRageWarrior) *Spell { return fw.AutoAttacks.OHAuto() },
			// 1.8 * 3.46 * 0.5
			wantHit: 3.114, wantCrit: 6.228,
		},
		{
			name:    "multiplied off hand",
			prepare: func(fw *FakeRageWarrior) { fw.SetOffHandRageMultiplier(1.5) },
			spell:   func(fw *FakeRageWarrior) *Spell { return fw.AutoAttacks.OHAuto() },
			// 1.8 * 3.46 * 0.5 * 1.5
			wantHit: 4.671, wantCrit: 9.342,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sim := SetupFakeRageSim()
			fw := sim.Raid.Parties[0].Players[0].(*FakeRageWarrior)
			tc.prepare(fw)
			spell := tc.spell(fw)

			hit := rageFromAutoAttack(sim, fw, spell, OutcomeHit, hitDamage)
			fw.ResetRageBar(sim, 0)
			crit := rageFromAutoAttack(sim, fw, spell, OutcomeCrit, 3*hitDamage)

			if !WithinToleranceFloat64(tc.wantHit, hit, 0.0005) {
				t.Fatalf("hit paid %0.4f Rage, want %0.4f", hit, tc.wantHit)
			}
			if !WithinToleranceFloat64(tc.wantCrit, crit, 0.0005) {
				t.Fatalf("crit paid %0.4f Rage, want %0.4f", crit, tc.wantCrit)
			}
		})
	}
}

// CritRageBonusSkipsTwoHand off multiplies a two-hander's whole swing again: twice a hit.
func TestCritRageBonusOnTheWholeTwoHandSwingWhenSwitchedOff(t *testing.T) {
	defer func(old bool) { CritRageBonusSkipsTwoHand = old }(CritRageBonusSkipsTwoHand)
	CritRageBonusSkipsTwoHand = false

	sim := SetupFakeRageSim()
	fw := sim.Raid.Parties[0].Players[0].(*FakeRageWarrior)
	fw.AutoAttacks.MH().NormalizedSwingSpeed = TwoHandNormalizedSwingSpeed
	spell := fw.AutoAttacks.MHAuto()
	hit := rageFromAutoAttack(sim, fw, spell, OutcomeHit, 500)
	fw.ResetRageBar(sim, 0)
	crit := rageFromAutoAttack(sim, fw, spell, OutcomeCrit, 1000)
	if !WithinToleranceFloat64(2, crit/hit, 0.0000001) {
		t.Fatalf("crit paid %0.6f times a hit, want 2", crit/hit)
	}
}

// CritRageMultiplier set back to client 70170's 1.75 pays the 75% bonus again.
func TestCritRageMultiplierSwitchesBackToSeventyFivePercent(t *testing.T) {
	defer func(old float64) { CritRageMultiplier = old }(CritRageMultiplier)
	CritRageMultiplier = 1.75

	sim := SetupFakeRageSim()
	fw := sim.Raid.Parties[0].Players[0].(*FakeRageWarrior)
	spell := fw.AutoAttacks.MHAuto()
	hit := rageFromAutoAttack(sim, fw, spell, OutcomeHit, 500)
	fw.ResetRageBar(sim, 0)
	crit := rageFromAutoAttack(sim, fw, spell, OutcomeCrit, 1000)
	if !WithinToleranceFloat64(1.75, crit/hit, 0.0000001) {
		t.Fatalf("crit paid %0.6f times a hit, want 1.75", crit/hit)
	}
}

// Rage from a hit taken is a separate rule, and a critical hit taken pays by the damage it dealt like
// any other: 1000 damage on 10000 maximum health is 1000 * 10 / 10000 = 1 Rage.
func TestCritBonusLeavesDamageTakenRageAlone(t *testing.T) {
	sim := SetupFakeRageSim()
	fw := sim.Raid.Parties[0].Players[0].(*FakeRageWarrior)
	fw.Unit.stats[stats.Health] = 10000
	spell := sim.Encounter.ActiveTargetUnits[0].AutoAttacks.MHAuto()

	taken := func(outcome HitOutcome) float64 {
		fw.ResetRageBar(sim, 0)
		result := &SpellResult{
			Target:                           &fw.Unit,
			Outcome:                          outcome,
			PostArmorAndResistanceMultiplier: 1000,
			ArmorAndResistanceMultiplier:     1,
			Damage:                           1000,
		}
		before := fw.CurrentRage()
		fw.Unit.OnSpellHitTaken(sim, spell, result)
		return fw.CurrentRage() - before
	}

	hit := taken(OutcomeHit)
	crit := taken(OutcomeCrit)
	if !WithinToleranceFloat64(1.0, hit, 0.0000001) {
		t.Fatalf("a hit taken paid %0.6f Rage, want 1", hit)
	}
	if hit != crit {
		t.Fatalf("a crit taken paid %0.6f Rage and a hit %0.6f, want the same for the same damage", crit, hit)
	}
}
