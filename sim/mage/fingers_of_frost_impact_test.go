package mage

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

// lastChargeFrostboltThenIceLance casts Frostbolt on the last Fingers of Frost charge from 20 yards and
// Ice Lance as soon as the Frostbolt leaves the hand. It reports whether the charge was still up when
// Ice Lance went out, and whether it is gone once both have landed.
func lastChargeFrostboltThenIceLance(t *testing.T) (frozenForIceLance bool, activeAfter bool) {
	t.Helper()
	player := &proto.Player{
		Name: "Mage", Class: proto.Class_ClassMage, Race: proto.Race_RaceGnome,
		TalentsString: FrostTalents, Equipment: &proto.EquipmentSpec{}, DistanceFromTarget: 20,
		Spec:     &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{ClassOptions: &proto.MageOptions{}}}},
		Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 42},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	mage := sim.Raid.Parties[0].Players[0].(MageAgent).GetMage()
	if mage.FingersOfFrostAura == nil {
		t.Fatal("FrostTalents no longer take Fingers of Frost; pick a build that does")
	}
	// No fresh charges: none from a chill landing, none from the Frostbolt's launch roll (upstream #776).
	mage.GetAura("Fingers of Frost Trigger").Deactivate(sim)
	mage.procFingersOfFrost = nil
	mage.FingersOfFrostAura.Activate(sim)
	mage.FingersOfFrostAura.SetStacks(sim, 1)

	frostbolt := mage.GetSpell(core.ActionID{SpellID: 25304})
	iceLance := mage.GetSpell(core.ActionID{SpellID: 1240047})
	if frostbolt.TravelTime() <= iceLance.TravelTime() || iceLance.TravelTime() <= 0 {
		t.Fatalf("travel times: Frostbolt %v, Ice Lance %v", frostbolt.TravelTime(), iceLance.TravelTime())
	}
	if !frostbolt.Cast(sim, mage.CurrentTarget) {
		t.Fatal("Frostbolt did not cast")
	}
	for mage.Hardcast.Expires > sim.CurrentTime || !iceLance.CanCast(sim, mage.CurrentTarget) {
		if sim.Step() {
			t.Fatal("simulation ended before the Frostbolt left the hand")
		}
	}
	if frostbolt.SpellMetrics[0].Casts != 1 || frostbolt.Unit.SpellsInFlight[frostbolt] != 1 {
		t.Fatal("Frostbolt is not in flight")
	}
	frozenForIceLance = mage.IsTargetFrozen()
	if !iceLance.Cast(sim, mage.CurrentTarget) {
		t.Fatal("Ice Lance did not cast")
	}
	end := sim.CurrentTime + time.Second
	for sim.CurrentTime < end {
		if sim.Step() {
			break
		}
	}
	return frozenForIceLance, mage.FingersOfFrostAura.IsActive()
}

func TestFingersOfFrostLastChargeSpentOnImpact(t *testing.T) {
	defer func(old bool) { FingersOfFrostSpentOnImpact = old }(FingersOfFrostSpentOnImpact)

	FingersOfFrostSpentOnImpact = false
	if frozen, active := lastChargeFrostboltThenIceLance(t); frozen || active {
		t.Errorf("spent on cast: Ice Lance frozen = %v, charge left = %v; want false, false", frozen, active)
	}

	FingersOfFrostSpentOnImpact = true
	if frozen, active := lastChargeFrostboltThenIceLance(t); !frozen || active {
		t.Errorf("spent on impact: Ice Lance frozen = %v, charge left = %v; want true, false", frozen, active)
	}
}

// The beta (mixarxrt, 16 October 2026: "fof last charge is that i can cast fb+il with 1 fof proc"):
// by default a Frostbolt on the last charge and an Ice Lance thrown while it travels are both frozen.
func TestFingersOfFrostLastChargeCoversFrostboltAndIceLance(t *testing.T) {
	if frozen, active := lastChargeFrostboltThenIceLance(t); !frozen || active {
		t.Errorf("default: Ice Lance frozen = %v, charge left = %v; want true, false", frozen, active)
	}
}

// Upstream #776 is kept: a Frostbolt cast on the last charge that procs Fingers of Frost at launch refills
// the charges and spends none of them, at cast end or when it lands, so two are up after it hits.
func TestFingersOfFrostLaunchProcSpendsNothingOnLanding(t *testing.T) {
	for seed := int64(1); seed <= 400; seed++ {
		player := &proto.Player{
			Name: "Mage", Class: proto.Class_ClassMage, Race: proto.Race_RaceGnome,
			TalentsString: FrostTalents, Equipment: &proto.EquipmentSpec{}, DistanceFromTarget: 20,
			Spec:     &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{ClassOptions: &proto.MageOptions{}}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}
		sim := core.NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: seed},
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
		}, simsignals.CreateSignals())
		sim.Reset()
		mage := sim.Raid.Parties[0].Players[0].(MageAgent).GetMage()
		mage.FingersOfFrostAura.Activate(sim)
		mage.FingersOfFrostAura.SetStacks(sim, 1)
		frostbolt := mage.GetSpell(core.ActionID{SpellID: 25304})
		if !frostbolt.Cast(sim, mage.CurrentTarget) {
			t.Fatal("Frostbolt did not cast")
		}
		for mage.Hardcast.Expires > sim.CurrentTime {
			if sim.Step() {
				t.Fatal("simulation ended during the cast")
			}
		}
		if mage.FingersOfFrostAura.GetStacks() != 2 {
			continue // no launch proc on this seed
		}
		end := sim.CurrentTime + frostbolt.TravelTime() + 100*time.Millisecond
		for sim.CurrentTime < end {
			if sim.Step() {
				break
			}
		}
		if got := mage.FingersOfFrostAura.GetStacks(); !mage.FingersOfFrostAura.IsActive() || got != 2 {
			t.Errorf("seed %d: %d charges after a launch-procced Frostbolt landed, want 2", seed, got)
		}
		return
	}
	t.Fatal("no seed procced Fingers of Frost at Frostbolt's launch")
}
