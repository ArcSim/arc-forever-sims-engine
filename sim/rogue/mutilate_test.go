package rogue

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// Mutilate's cast rolls miss/dodge/parry once; its two hand strikes carry the client's No Attack
// Miss/Dodge/Parry attributes and never roll them again. Beta logs agree: in foreverlogs 2701/2702
// the level 30 cast (1310707) missed, was dodged or parried 16 times in 162, the hand strikes
// (1310705/1310706) 0 times in 280.
func TestMutilateHitsDoNotRollAvoidance(t *testing.T) {
	result := core.RunRaidSim(&proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(&proto.Player{
			Race: proto.Race_RaceHuman, Class: proto.Class_ClassRogue,
			Equipment: DefaultGear.GearSet, TalentsString: AssassinationTalents, Spec: DefaultOptions,
			Rotation: core.GetAplRotation("../../ui/specs/rogue/dps/apls", "forever_mutilate").Rotation,
		}, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
		SimOptions: &proto.SimOptions{Iterations: 5, RandomSeed: 1},
	})
	if result.Error != nil {
		t.Fatal(result.Error.Message)
	}
	casts, castAvoided, hits := int32(0), int32(0), int32(0)
	for _, action := range result.RaidMetrics.Parties[0].Players[0].Actions {
		if action.Id.GetSpellId() != MutilateSpellID {
			continue
		}
		for _, target := range action.Targets {
			avoided := target.Misses + target.Dodges + target.Parries
			if action.Id.GetTag() == 0 {
				casts += target.Casts
				castAvoided += avoided
			} else {
				hits += target.Hits + target.Crits
				if avoided > 0 {
					t.Errorf("Mutilate hand strike (tag %d) was avoided %d times", action.Id.GetTag(), avoided)
				}
			}
		}
	}
	if casts == 0 || castAvoided == 0 || hits == 0 {
		t.Errorf("want Mutilate casts, some avoided and hand strikes: casts %d, avoided %d, strikes %d", casts, castAvoided, hits)
	}
}

// Before the 16 October 2026 beta notes (ColdBloodSpentByMutilate false): Cold Blood crits both Mutilate
// hands and isn't spent by them; the next Eviscerate spends it.
func TestColdBloodLastsThroughMutilateBeforeTheFix(t *testing.T) {
	ColdBloodSpentByMutilate = false
	defer func() { ColdBloodSpentByMutilate = true }()

	sim, rogue := kidneyShotSim(AssassinationTalents)
	target := rogue.CurrentTarget
	coldBlood := rogue.GetAura("Cold Blood")

	rogue.ColdBlood.SkipCastAndApplyEffects(sim, target)
	for i := 0; i < 10; i++ {
		rogue.Mutilate.SkipCastAndApplyEffects(sim, target)
	}
	if !coldBlood.IsActive() {
		t.Fatal("Cold Blood was spent by Mutilate")
	}
	for _, hand := range []*core.Spell{rogue.MutilateMH, rogue.MutilateOH} {
		m := hand.SpellMetrics[target.UnitIndex]
		if m.Crits == 0 || m.Hits != 0 {
			t.Errorf("%v under Cold Blood: %d crits, %d hits; want only crits", hand.ActionID, m.Crits, m.Hits)
		}
	}

	metrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 31016, Tag: 9})
	for i := 0; i < 20 && coldBlood.IsActive(); i++ {
		rogue.AddComboPoints(sim, 5, metrics)
		rogue.Eviscerate.SkipCastAndApplyEffects(sim, target)
	}
	if coldBlood.IsActive() {
		t.Error("Cold Blood survived Eviscerate")
	}
}

// 16 October 2026 beta notes: "Cold Blood is now consumed when Mutilate crits". One Mutilate crits with
// both hands and spends Cold Blood, once for the cast: the main hand's crit leaves it up for the off
// hand.
func TestColdBloodSpentByMutilate(t *testing.T) {
	if !ColdBloodSpentByMutilate {
		t.Fatal("ColdBloodSpentByMutilate is off by default")
	}
	sim, rogue := kidneyShotSim(AssassinationTalents)
	target := rogue.CurrentTarget
	coldBlood := rogue.GetAura("Cold Blood")

	rogue.ColdBlood.SkipCastAndApplyEffects(sim, target)
	casts := 0
	for ; casts < 20 && coldBlood.IsActive(); casts++ {
		before := rogue.MutilateMH.SpellMetrics[target.UnitIndex].Crits
		rogue.Mutilate.SkipCastAndApplyEffects(sim, target)
		if rogue.MutilateMH.SpellMetrics[target.UnitIndex].Crits == before {
			continue // the cast missed, was dodged or parried: no hand strikes, Cold Blood stays
		}
		break
	}
	if coldBlood.IsActive() {
		t.Fatal("Cold Blood survived a Mutilate crit")
	}
	for _, hand := range []*core.Spell{rogue.MutilateMH, rogue.MutilateOH} {
		m := hand.SpellMetrics[target.UnitIndex]
		if m.Crits != 1 || m.Hits != 0 {
			t.Errorf("%v under Cold Blood: %d crits, %d hits; want 1 crit", hand.ActionID, m.Crits, m.Hits)
		}
	}

	// The next Mutilate has no Cold Blood: it is not spent again, and a recast is spent by the next crit.
	rogue.ColdBlood.CD.Reset()
	rogue.ColdBlood.SkipCastAndApplyEffects(sim, target)
	for i := 0; i < 20 && coldBlood.IsActive(); i++ {
		rogue.Mutilate.SkipCastAndApplyEffects(sim, target)
	}
	if coldBlood.IsActive() {
		t.Error("a second Cold Blood survived Mutilate")
	}
}
