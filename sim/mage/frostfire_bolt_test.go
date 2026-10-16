package mage

import (
	"math"
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
)

const frostfireHybridTalents = "-1355000013013403-00550003210013002"

func frostfireSim(talents string) (*core.Simulation, *Mage) {
	player := &proto.Player{
		Name: "Frostfire", Class: proto.Class_ClassMage, Race: proto.Race_RaceHuman,
		TalentsString: talents, Equipment: &proto.EquipmentSpec{},
		Spec:     &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{ClassOptions: &proto.MageOptions{}}}},
		Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 42},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim, sim.Raid.Parties[0].Players[0].(MageAgent).GetMage()
}

func TestFrostfireRanksAndTalents(t *testing.T) {
	_, plain := frostfireSim("")
	sim, hybrid := frostfireSim(frostfireHybridTalents)
	for i, id := range []int32{401502, 1237312, 1237313} {
		spell := plain.GetSpell(core.ActionID{SpellID: id})
		if spell == nil {
			t.Fatalf("baseline rank %d missing", id)
		}
		if spell.Rank != int32(i+1) || spell.SpellSchool != core.SpellSchoolFrostfire || spell.DefaultCast.CastTime != 3*time.Second {
			t.Errorf("bad baseline rank %d: rank %d, school %v, cast %v", id, spell.Rank, spell.SpellSchool, spell.DefaultCast.CastTime)
		}
		dot := spell.Dot(plain.CurrentTarget)
		if dot.BaseTickCount != 3 || dot.BaseTickLength != 3*time.Second || dot.BonusCoefficient != 0 {
			t.Errorf("rank %d: expected three 3s ticks without spell-power scaling", id)
		}
	}
	spell := hybrid.GetSpell(core.ActionID{SpellID: 1237313})
	if spell.DefaultCast.CastTime != 2500*time.Millisecond {
		t.Errorf("Improved Fireball cast = %v", spell.DefaultCast.CastTime)
	}
	if got := spell.Cost.GetCurrentCost(); math.Abs(got-333) > 1e-9 {
		t.Errorf("Frost Channeling cost = %v, want 333", got)
	}
	if got := spell.SpellHitChance(hybrid.CurrentTarget); math.Abs(got-.05) > 1e-9 {
		t.Errorf("hit = %v, want .05", got)
	}
	// Piercing Ice adds 6% to the hit and 2% to the DoT. Fire Power (client 70425) counts the Fire
	// effects on the target, and a fresh target has none.
	table := hybrid.AttackTables[hybrid.CurrentTarget.UnitIndex]
	if got := spell.AttackerDamageMultiplier(table, false); math.Abs(got-1.06) > 1e-9 {
		t.Errorf("Piercing Ice direct multiplier = %v, want 1.06", got)
	}
	if got := spell.AttackerDamageMultiplier(table, true); math.Abs(got-1.02) > 1e-9 {
		t.Errorf("Piercing Ice periodic multiplier = %v, want 1.02", got)
	}
	if got := spell.TargetDamageMultiplier(sim, table, false); math.Abs(got-1) > 1e-9 {
		t.Errorf("Fire Power with no Fire effect on the target = %v, want 1", got)
	}
	if spell.BonusCritPercent != 8 {
		t.Errorf("Critical Mass = %v, want 8", spell.BonusCritPercent)
	}
	// A guaranteed direct crit proves both the Frost crit bonus and Fire proc interactions.
	spell.BonusCritPercent += 100
	spell.Flags |= core.SpellFlagIgnoreResists
	result := spell.CalcDamage(sim, hybrid.CurrentTarget, 100, spell.OutcomeMagicCrit)
	if math.Abs(result.Damage-212) > 1e-6 {
		t.Errorf("Ice Shards crit = %v, want 212", result.Damage)
	}
	spell.CurCast.Cost = 333
	hybrid.SpendMana(sim, 333, hybrid.NewManaMetrics(spell.ActionID))
	manaBefore := hybrid.CurrentMana()
	spell.DealDamage(sim, result)
	if got := hybrid.CurrentMana() - manaBefore; math.Abs(got-111) > 1e-9 {
		t.Errorf("Master of Elements refund = %v, want 111", got)
	}
	if hybrid.HotStreakAura.GetStacks() != 1 || !hybrid.Ignite.Dot(hybrid.CurrentTarget).IsActive() {
		t.Fatal("Frostfire crit did not trigger Hot Streak and Ignite")
	}
	// Periodic crits must not add Hot Streak or another Ignite contribution.
	igniteBefore := hybrid.Ignite.Dot(hybrid.CurrentTarget).SnapshotBaseDamage
	dot := spell.Dot(hybrid.CurrentTarget)
	dot.Apply(sim)
	dot.CalcAndDealPeriodicSnapshotDamage(sim, hybrid.CurrentTarget, spell.OutcomeTickMagicCrit)
	if hybrid.HotStreakAura.GetStacks() != 1 || hybrid.Ignite.Dot(hybrid.CurrentTarget).SnapshotBaseDamage != igniteBefore {
		t.Fatal("periodic Frostfire damage triggered a direct-crit talent")
	}
}

func TestFrostfireFingersAndMissileBarrage(t *testing.T) {
	sim, mage := frostfireSim(frostfireHybridTalents)
	spell := mage.GetSpell(core.ActionID{SpellID: 1237313})
	// Fingers of Frost rolls when the bolt launches, on the hit rolled there (upstream #776).
	for i := 0; i < 200 && !mage.FingersOfFrostAura.IsActive(); i++ {
		spell.SkipCastAndApplyEffects(sim, mage.CurrentTarget)
	}
	if mage.FingersOfFrostAura.GetStacks() != 2 {
		t.Fatal("Frostfire chill did not grant two Fingers of Frost charges")
	}
	// The launching cast completes without spending a charge.
	mage.FingersOfFrostAura.OnCastComplete(mage.FingersOfFrostAura, sim, spell)
	if spell.BonusCritPercent != 58 {
		t.Errorf("Shatter crit bonus = %v, want 58", spell.BonusCritPercent)
	}
	// Finish one real cast. Turn off the launch roll so it cannot replenish the charge.
	mage.procFingersOfFrost = nil
	if !spell.Cast(sim, mage.CurrentTarget) {
		t.Fatal("Frostfire did not cast")
	}
	for sim.CurrentTime < 3*time.Second {
		if sim.Step() {
			break
		}
	}
	if mage.FingersOfFrostAura.GetStacks() != 1 {
		t.Errorf("Frostfire left %d charges, want 1", mage.FingersOfFrostAura.GetStacks())
	}
	// An Arcane build can use the same baseline spell to proc Missile Barrage, which rolls when the
	// bolt lands (upstream #643's 400588 proc flags).
	sim, mage = frostfireSim(ArcaneTalents)
	spell = mage.GetSpell(core.ActionID{SpellID: 1237313})
	for i := 0; i < 100 && !mage.MissileBarrageAura.IsActive(); i++ {
		spell.CalcAndDealDamage(sim, mage.CurrentTarget, 1, spell.OutcomeAlwaysHit)
	}
	if !mage.MissileBarrageAura.IsActive() {
		t.Fatal("Frostfire did not trigger Missile Barrage")
	}
}

func TestFrostfireDamageAndManaRefund(t *testing.T) {
	sim, mage := frostfireSim(frostfireHybridTalents)
	spell := mage.GetSpell(core.ActionID{SpellID: 1237313})
	mage.AddStatsDynamic(sim, stats.Stats{stats.SpellCritPercent: 100})
	before := mage.CurrentMana()
	if !spell.Cast(sim, mage.CurrentTarget) {
		t.Fatal("Frostfire did not cast")
	}
	for sim.CurrentTime < 12*time.Second {
		if sim.Step() {
			break
		}
	}
	if spell.SpellMetrics[mage.CurrentTarget.UnitIndex].TotalDamage <= 0 {
		t.Fatal("Frostfire dealt no damage")
	}
	if mage.CurrentMana() <= before-333 {
		t.Fatal("no mana returned after the crit")
	}
}
