package warlock

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
)

// Hellfire's area hits crit only where the Hellfire Effect row lets them: client 70124 marks it
// Cannot Crit, 70170 does not. The test flips the flag both ways, so it holds on either build.
func TestHellfireCritsWhereItsEffectRowAllows(t *testing.T) {
	burn := spellData.Hellfire.Highest().Effect(dbcenums.A_PERIODIC_TRIGGER_SPELL, 0).Trigger()
	saved := burn.Attr[dbcenums.ATTR_INDEX_EX_2]
	defer func() { burn.Attr[dbcenums.ATTR_INDEX_EX_2] = saved }()

	for _, cannotCrit := range []bool{true, false} {
		burn.Attr[dbcenums.ATTR_INDEX_EX_2] = saved &^ dbcenums.ATTR_EX_2_CANT_CRIT
		if cannotCrit {
			burn.Attr[dbcenums.ATTR_INDEX_EX_2] |= dbcenums.ATTR_EX_2_CANT_CRIT
		}

		player := core.WithSpec(&proto.Player{
			Race:        proto.Race_RaceOrc,
			Class:       proto.Class_ClassWarlock,
			Equipment:   &proto.EquipmentSpec{},
			Consumables: &proto.ConsumesSpec{},
			Rotation:    &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.Player_Warlock{Warlock: &proto.Warlock{Options: &proto.Warlock_Options{ClassOptions: &proto.WarlockOptions{
			Summon: proto.WarlockOptions_NoSummon,
		}}}})
		sim := core.NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 1},
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
		}, simsignals.CreateSignals())
		sim.Reset()

		warlock := sim.Raid.Parties[0].Players[0].(WarlockAgent).GetWarlock()
		warlock.AddStatDynamic(sim, stats.SpellCritPercent, 100)
		dot := warlock.Hellfire.AOEDot()
		dot.Apply(sim)
		for range 5 {
			dot.TickOnce(sim)
		}

		m := warlock.Hellfire.SpellMetrics[0]
		if gotCrits := m.CritTicks > 0; gotCrits == cannotCrit {
			t.Errorf("Cannot Crit %v: %d crit ticks of %d", cannotCrit, m.CritTicks, m.Ticks+m.CritTicks)
		}
	}
}

// Malediction raises Hellfire's area hits like Corruption's ticks: the client left them out (they
// are Hellfire Effect 11682's direct School Damage), dev fixed it in the next build (forever-bugs #311).
func TestMaledictionRaisesHellfire(t *testing.T) {
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceOrc,
		Class:         proto.Class_ClassWarlock,
		Equipment:     &proto.EquipmentSpec{},
		Consumables:   &proto.ConsumesSpec{},
		TalentsString: "0005",
		Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, &proto.Player_Warlock{Warlock: &proto.Warlock{Options: &proto.Warlock_Options{ClassOptions: &proto.WarlockOptions{
		Summon: proto.WarlockOptions_NoSummon,
	}}}})
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	warlock := sim.Raid.Parties[0].Players[0].(WarlockAgent).GetWarlock()
	if got := warlock.Hellfire.AOEDot().PeriodicDamageMultiplier; got <= 1 {
		t.Errorf("Hellfire periodic multiplier with 5 Malediction = %v, want > 1", got)
	}
	if got := warlock.Corruption.Dot(warlock.CurrentTarget).PeriodicDamageMultiplier; got <= 1 {
		t.Errorf("Corruption periodic multiplier with 5 Malediction = %v, want > 1", got)
	}
}

// The sim has no healer, so Hellfire's self-burn must not cut the channel short: it used to stop once
// a tick outburned the warlock's health, which turned every channel past ~18 s into one-tick recasts.
func TestHellfireChannelsAtLowHealth(t *testing.T) {
	player := core.WithSpec(&proto.Player{
		Race:        proto.Race_RaceOrc,
		Class:       proto.Class_ClassWarlock,
		Equipment:   &proto.EquipmentSpec{},
		Consumables: &proto.ConsumesSpec{},
		Rotation:    &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, &proto.Player_Warlock{Warlock: &proto.Warlock{Options: &proto.Warlock_Options{ClassOptions: &proto.WarlockOptions{
		Summon: proto.WarlockOptions_NoSummon,
	}}}})
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	warlock := sim.Raid.Parties[0].Players[0].(WarlockAgent).GetWarlock()
	warlock.RemoveHealth(sim, warlock.CurrentHealth()-1)
	dot := warlock.Hellfire.AOEDot()
	dot.Apply(sim)
	for range 3 {
		dot.TickOnce(sim)
	}
	if !dot.IsActive() {
		t.Error("Hellfire stopped channelling at low health")
	}
}
