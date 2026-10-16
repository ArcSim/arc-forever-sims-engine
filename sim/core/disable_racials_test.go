package core

import (
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
)

// racialWarrior builds a finalized (fake) warrior of the given race, so a test can read what
// its racials registered. mutate, if set, edits the player proto first.
func racialWarrior(race proto.Race, mutate func(*proto.Player)) *Character {
	player := &proto.Player{
		Name:        "Warrior",
		Race:        race,
		Class:       proto.Class_ClassWarrior,
		Buffs:       &proto.IndividualBuffs{},
		Consumables: &proto.ConsumesSpec{},
		Spec:        &proto.Player_DpsWarrior{},
		Equipment:   &proto.EquipmentSpec{},
		Rotation:    &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}
	if mutate != nil {
		mutate(player)
	}
	sim := NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 100},
		Raid: &proto.Raid{Parties: []*proto.Party{{
			Players: []*proto.Player{player},
			Buffs:   &proto.PartyBuffs{},
		}}},
		Encounter: &proto.Encounter{Targets: []*proto.Target{{Level: 63}}, Duration: 180},
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim.Raid.Parties[0].Players[0].GetCharacter()
}

func disableRacials(player *proto.Player) { player.DisableRacials = true }

// DisableRacials measures what a race's racials are worth: the same character, base stats
// included, with none of its racial effects.
func TestDisableRacialsDropsRacialsButKeepsBaseStats(t *testing.T) {
	withRacials := racialWarrior(proto.Race_RaceOrc, nil)
	orc := racialWarrior(proto.Race_RaceOrc, disableRacials)

	if withRacials.GetAura("Blood Fury") == nil {
		t.Fatal("an orc warrior has no Blood Fury, so this test checks nothing")
	}
	if orc.GetAura("Blood Fury") != nil {
		t.Error("an orc with racials disabled still has Blood Fury")
	}
	if got, want := orc.GetBaseStats(), withRacials.GetBaseStats(); got != want {
		t.Errorf("disabling racials changed the orc's base stats: %v, want %v", got, want)
	}

	// The human's spirit racial is a multiplier on top of the base stats, so it goes too.
	human := racialWarrior(proto.Race_RaceHuman, nil)
	humanWithout := racialWarrior(proto.Race_RaceHuman, disableRacials)
	if human.GetStat(stats.Spirit) <= humanWithout.GetStat(stats.Spirit) {
		t.Errorf("human spirit %.1f with racials, %.1f without; the racial should raise it",
			human.GetStat(stats.Spirit), humanWithout.GetStat(stats.Spirit))
	}
	if got, want := humanWithout.GetStat(stats.Spirit), humanWithout.GetBaseStats()[stats.Spirit]; got != want {
		t.Errorf("human spirit with racials disabled is %.1f, want the base %.1f", got, want)
	}
}

// A warrior wielding a mace (Ironfoe) and a sword (Cho'Rush's Blade), the Fury reference's weapons.
func weaponWarrior(race proto.Race, mutate func(*proto.Player)) *Character {
	return racialWarrior(race, func(player *proto.Player) {
		player.Equipment = &proto.EquipmentSpec{Items: []*proto.ItemSpec{
			{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
			{Id: 11684}, {Id: 18484},
		}}
		if mutate != nil {
			mutate(player)
		}
	})
}

// The override still retypes the weapons, but since the 16 October 2026 beta notes no racial depends on
// the weapon type, so neither it nor DisableWeaponSpecialization moves a race's crit.
func TestWeaponTypeOverrideNoLongerMovesRacials(t *testing.T) {
	crit := func(c *Character) float64 { return c.GetStat(stats.PhysicalCritPercent) }
	orc := weaponWarrior(proto.Race_RaceOrc, nil)
	orcAxes := weaponWarrior(proto.Race_RaceOrc, func(p *proto.Player) { p.WeaponTypeOverride = proto.WeaponType_WeaponTypeAxe })
	if got := orc.MainHand().WeaponType; got != proto.WeaponType_WeaponTypeMace {
		t.Fatalf("test setup: main hand is %v, want a mace", got)
	}
	if got := orcAxes.MainHand().WeaponType; got != proto.WeaponType_WeaponTypeAxe {
		t.Errorf("override left the main hand a %v", got)
	}
	if got := orcAxes.OffHand().WeaponType; got != proto.WeaponType_WeaponTypeAxe {
		t.Errorf("override left the off hand a %v", got)
	}
	if d := crit(orcAxes) - crit(orc); d != 0 {
		t.Errorf("an orc with axes has %.2f more crit than with a mace and a sword, want none", d)
	}
	if orcAxes.MainHand().ID != orc.MainHand().ID {
		t.Error("override changed the weapon itself")
	}

	offHand := weaponWarrior(proto.Race_RaceOrc, func(p *proto.Player) {
		p.WeaponTypeOverride = proto.WeaponType_WeaponTypeAxe
		p.WeaponTypeOverrideOffHandOnly = true
	})
	if offHand.MainHand().WeaponType != proto.WeaponType_WeaponTypeMace || offHand.OffHand().WeaponType != proto.WeaponType_WeaponTypeAxe {
		t.Error("off-hand-only override changed the main hand or left the off hand")
	}

	human := weaponWarrior(proto.Race_RaceHuman, nil)
	humanWithheld := weaponWarrior(proto.Race_RaceHuman, func(p *proto.Player) { p.DisableWeaponSpecialization = true })
	if d := crit(human) - crit(humanWithheld); d != 0 {
		t.Errorf("DisableWeaponSpecialization changed a human's crit by %.2f, want none", d)
	}
}
