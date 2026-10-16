package core

import (
	"math"
	"slices"
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
)

func setupRacialSim(race proto.Race) (*Simulation, *FakeRageWarrior) {
	sim := NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{
			RandomSeed: 100,
		},
		Raid: &proto.Raid{
			Parties: []*proto.Party{
				{
					Players: []*proto.Player{
						{
							Name:      "Warrior",
							Class:     proto.Class_ClassWarrior,
							Race:      race,
							Buffs:     &proto.IndividualBuffs{},
							Spec:      &proto.Player_DpsWarrior{},
							Equipment: &proto.EquipmentSpec{},
							Rotation:  &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
						},
					},
					Buffs: &proto.PartyBuffs{},
				},
			},
		},
		Encounter: &proto.Encounter{
			Targets: []*proto.Target{
				{Name: "target", Level: 63, MobType: proto.MobType_MobTypeElemental},
			},
			Duration: 180,
		},
	}, simsignals.CreateSignals())
	sim.Reset()

	return sim, sim.Raid.Parties[0].Players[0].(*FakeRageWarrior)
}

func castRacial(t *testing.T, sim *Simulation, fw *FakeRageWarrior, spellID int32) {
	t.Helper()
	spell := fw.GetSpell(ActionID{SpellID: spellID})
	if spell == nil {
		t.Fatalf("Racial %d is not registered", spellID)
	}
	if !spell.Cast(sim, fw.CurrentTarget) {
		t.Fatalf("Racial %d could not be cast", spellID)
	}
}

func TestPlayableRaces(t *testing.T) {
	for class, races := range ClassRaceCapabilities {
		for _, race := range races {
			if race == proto.Race_RaceBloodElf {
				t.Errorf("%s is not playable, but %s lists it", race, class)
			}
			if _, ok := BaseStats[BaseStatsKey{Race: race, Class: class}]; !ok {
				t.Errorf("No base stats for %s %s", race, class)
			}
		}
	}

	for class, race := range map[proto.Class]proto.Race{
		proto.Class_ClassPaladin: proto.Race_RaceUndead,
		proto.Class_ClassShaman:  proto.Race_RaceDwarf,
		proto.Class_ClassHunter:  proto.Race_RaceHuman,
		proto.Class_ClassMage:    proto.Race_RaceOrc,
		proto.Class_ClassPriest:  proto.Race_RaceGnome,
		proto.Class_ClassWarlock: proto.Race_RaceTroll,
		proto.Class_ClassDruid:   proto.Race_RaceHighOrderSkyborne,
	} {
		if !slices.Contains(ClassRaceCapabilities[class], race) {
			t.Errorf("%s should be playable as %s", class, race)
		}
	}
	if slices.Contains(ClassRaceCapabilities[proto.Class_ClassShaman], proto.Race_RaceHighOrderSkyborne) {
		t.Errorf("Only the Windshaper Skyborne can be shamans")
	}
	if slices.Contains(ClassRaceCapabilities[proto.Class_ClassMage], proto.Race_RaceWindshaperSkyborne) {
		t.Errorf("Only the High Order Skyborne can be mages")
	}
}

func TestBloodFuryIsPercentBased(t *testing.T) {
	sim, fw := setupRacialSim(proto.Race_RaceOrc)

	apBefore := fw.GetStat(stats.AttackPower)
	castRacial(t, sim, fw, 20572)

	if !WithinToleranceFloat64(apBefore*1.1, fw.GetStat(stats.AttackPower), 0.01) {
		t.Fatalf("Blood Fury should grant 10%% Attack Power: before %0.2f, after %0.2f", apBefore, fw.GetStat(stats.AttackPower))
	}
}

func TestBerserkingIsFlatTenPercent(t *testing.T) {
	sim, fw := setupRacialSim(proto.Race_RaceTroll)

	speedBefore := fw.PseudoStats.AttackSpeedMultiplier
	castRacial(t, sim, fw, 20554)

	if !WithinToleranceFloat64(speedBefore*1.1, fw.PseudoStats.AttackSpeedMultiplier, 0.0001) {
		t.Fatalf("Berserking should grant 10%% attack speed, got x%0.4f", fw.PseudoStats.AttackSpeedMultiplier/speedBefore)
	}
	if fw.CurrentRage() != 0 || fw.GetSpell(ActionID{SpellID: 20554}).Cost != nil {
		t.Fatalf("Berserking should have no cost")
	}
}

func TestElunesLight(t *testing.T) {
	sim, fw := setupRacialSim(proto.Race_RaceNightElf)

	critBefore := fw.GetStat(stats.PhysicalCritPercent)
	castRacial(t, sim, fw, 1259799)

	if !WithinToleranceFloat64(critBefore+10, fw.GetStat(stats.PhysicalCritPercent), 0.0001) {
		t.Fatalf("Elune's Light should grant 10%% crit: before %0.2f, after %0.2f", critBefore, fw.GetStat(stats.PhysicalCritPercent))
	}
}

func TestGnomeWarrior(t *testing.T) {
	sim, fw := setupRacialSim(proto.Race_RaceGnome)

	if fw.MaximumRage() != 105 {
		t.Fatalf("Expansive Mind should raise maximum Rage to 105, got %0.1f", fw.MaximumRage())
	}

	// Client 70170 family masks: Heroic Strike is named, Rend (a DoT) no longer is.
	defer func(lookup func(int32) ClassFlags) { ClientClassFlags = lookup }(ClientClassFlags)
	ClientClassFlags = func(id int32) ClassFlags {
		return map[int32]ClassFlags{78: {Family: 4, Mask: [4]uint32{64}}, 99: {Family: 4, Mask: [4]uint32{64}}, 772: {Family: 4, Mask: [4]uint32{32}}}[id]
	}

	// A spell with a hit and a DoT; the DoT is attached by hand, as the test sim takes no new auras.
	hybrid := fw.RegisterSpell(SpellConfig{
		ActionID:         ActionID{SpellID: 99},
		ClassSpellMask:   4,
		ProcMask:         ProcMaskSpellDamage,
		SpellSchool:      SpellSchoolFire,
		DamageMultiplier: 1,
	})
	hybridDot := &Dot{PeriodicDamageMultiplier: 1}
	hybrid.dots = DotArray{hybridDot}

	castRacial(t, sim, fw, 1259813)
	eureka := fw.GetAura("Eureka!")
	if !eureka.IsActive() || eureka.GetStacks() != 3 {
		t.Fatalf("Eureka! should start at 3 charges, got %d", eureka.GetStacks())
	}


	ability := fw.RegisterSpell(SpellConfig{
		ActionID:         ActionID{SpellID: 78},
		ClassSpellMask:   1,
		ProcMask:         ProcMaskMeleeMHSpecial,
		SpellSchool:      SpellSchoolPhysical,
		DamageMultiplier: 1,
	})
	rend := fw.RegisterSpell(SpellConfig{
		ActionID:         ActionID{SpellID: 772},
		ClassSpellMask:   2,
		ProcMask:         ProcMaskMeleeMHSpecial,
		SpellSchool:      SpellSchoolPhysical,
		DamageMultiplier: 1,
	})
	itemSpell := fw.RegisterSpell(SpellConfig{
		ActionID:         ActionID{SpellID: 2},
		Flags:            SpellFlagAPL,
		ProcMask:         ProcMaskMeleeMHSpecial,
		SpellSchool:      SpellSchoolPhysical,
		DamageMultiplier: 1,
	})
	rend.Cast(sim, fw.CurrentTarget)
	if rend.DamageMultiplier != 1 || eureka.GetStacks() != 3 {
		t.Fatalf("Eureka! should leave Rend alone since build 70170 (non-periodic abilities only)")
	}
	itemSpell.Cast(sim, fw.CurrentTarget)
	if itemSpell.DamageMultiplier != 1 || eureka.GetStacks() != 3 {
		t.Fatalf("Eureka! should ignore spells that are not class abilities")
	}
	if !WithinToleranceFloat64(1.1, ability.DamageMultiplier, 0.0001) {
		t.Fatalf("Eureka! should raise ability damage by 10%%, got x%0.2f", ability.DamageMultiplier)
	}

	// SPELLMOD_DAMAGE is the direct part: a named spell's DoT ticks (Immolate, Fireball) take nothing.
	if !WithinToleranceFloat64(1.1, hybrid.DamageMultiplier, 0.0001) ||
		!WithinToleranceFloat64(1, hybrid.DamageMultiplier*hybridDot.PeriodicDamageMultiplier, 0.0001) {
		t.Fatalf("Eureka! should raise the hit by 10%% and leave the DoT alone, got x%0.2f / x%0.2f",
			hybrid.DamageMultiplier, hybrid.DamageMultiplier*hybridDot.PeriodicDamageMultiplier)
	}

	for range 3 {
		ability.Cast(sim, fw.CurrentTarget)
	}
	if eureka.IsActive() {
		t.Fatalf("Eureka! should end after 3 abilities")
	}
	if !WithinToleranceFloat64(1, ability.DamageMultiplier, 0.0001) {
		t.Fatalf("Eureka! should drop its damage bonus when it ends, got x%0.2f", ability.DamageMultiplier)
	}
	if !WithinToleranceFloat64(1, hybridDot.PeriodicDamageMultiplier, 0.0001) {
		t.Fatalf("Eureka! should restore the DoT multiplier when it ends, got x%0.2f", hybridDot.PeriodicDamageMultiplier)
	}
}

// 16 Oct build: Endurance (tauren) and Earth Fury (dwarf) are 2% critical damage and healing; Endurance lost its 1% hit.
func TestCritDamageRacials(t *testing.T) {
	_, human := setupRacialSim(proto.Race_RaceHuman)
	for _, race := range []proto.Race{proto.Race_RaceTauren, proto.Race_RaceDwarf} {
		_, fw := setupRacialSim(race)
		if got := fw.PseudoStats.CritDamageMultiplier; got != 1.02 {
			t.Errorf("%s: crit damage x%.4f, want x1.02", race, got)
		}
		if got, want := fw.GetStat(stats.PhysicalHitPercent), human.GetStat(stats.PhysicalHitPercent); got != want {
			t.Errorf("%s: hit %.2f, want %.2f (no hit racial)", race, got, want)
		}
	}
}

func TestTouchOfTheGraveUsesTheMeleeVariant(t *testing.T) {
	_, fw := setupRacialSim(proto.Race_RaceUndead)

	if aura := fw.GetAura("Touch of the Grave"); aura == nil || aura.ActionIDForProc.SpellID != 1260189 {
		t.Fatalf("An undead warrior should carry the 5%% Touch of the Grave")
	} else if aura.OnSpellHitDealt == nil {
		t.Fatalf("Touch of the Grave listens for no hit, so it never procs")
	}
	if fw.GetSpell(ActionID{SpellID: 1260198}) == nil {
		t.Fatalf("Touch of the Grave's drain is not registered")
	}
}

func TestSkyborneRacials(t *testing.T) {
	_, human := setupRacialSim(proto.Race_RaceHuman)

	for _, race := range []proto.Race{proto.Race_RaceHighOrderSkyborne, proto.Race_RaceWindshaperSkyborne} {
		_, fw := setupRacialSim(race)

		swingRatio := human.AutoAttacks.MainhandSwingSpeed().Seconds() / fw.AutoAttacks.MainhandSwingSpeed().Seconds()
		if !WithinToleranceFloat64(1.03, swingRatio, 0.001) {
			t.Errorf("%s: Wind Blessed should swing 3%% faster, got x%0.4f", race, swingRatio)
		}
		if at := fw.AttackTables[fw.CurrentTarget.UnitIndex]; !WithinToleranceFloat64(1.05, at.DamageDealtMultiplier, 0.0001) {
			t.Errorf("%s: Elemental Insight should grant 5%% damage against elementals, got x%0.4f", race, at.DamageDealtMultiplier)
		}
		if at := fw.AttackTables[fw.CurrentTarget.UnitIndex]; at.CritMultiplier != 1 {
			t.Errorf("%s: Elemental Insight has no crit damage effect, got crit x%0.4f", race, at.CritMultiplier)
		}

		hasReadLeyLine := fw.GetSpell(ActionID{SpellID: 1259705}) != nil
		hasSkysight := fw.GetSpell(ActionID{SpellID: 1259686}) != nil
		if hasReadLeyLine != (race == proto.Race_RaceHighOrderSkyborne) || hasSkysight != (race == proto.Race_RaceWindshaperSkyborne) {
			t.Errorf("%s: Read Ley Line is the High Order's and Skysight the Windshapers', got %t / %t", race, hasReadLeyLine, hasSkysight)
		}
		// Client 1.60.1.70058 SpellMisc: Read Ley Line casts in 2 sec, Skysight in 0.5 sec.
		for id, want := range map[int32]time.Duration{1259705: time.Second * 2, 1259686: time.Millisecond * 500} {
			if spell := fw.GetSpell(ActionID{SpellID: id}); spell != nil && spell.DefaultCast.CastTime != want {
				t.Errorf("%s: spell %d should cast in %s, got %s", race, id, want, spell.DefaultCast.CastTime)
			}
		}
	}
}

// Blood Fury's multipliers report the stats they add on gain and the stats they remove on expire.
func TestBloodFuryReportsItsTemporaryStats(t *testing.T) {
	sim := NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{
			Players: []*proto.Player{{
				Name:      "Orc",
				Class:     proto.Class_ClassWarrior,
				Race:      proto.Race_RaceOrc,
				Buffs:     &proto.IndividualBuffs{},
				Spec:      &proto.Player_DpsWarrior{},
				Equipment: &proto.EquipmentSpec{},
				Rotation:  &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
			}},
			Buffs: &proto.PartyBuffs{},
		}}},
		Encounter: &proto.Encounter{
			Targets:  []*proto.Target{{Name: "target", Level: 63, MobType: proto.MobType_MobTypeHumanoid}},
			Duration: 180,
		},
	}, simsignals.CreateSignals())
	sim.Reset()

	character := sim.Raid.Parties[0].Players[0].GetCharacter()
	bloodFury := character.GetAura("Blood Fury")
	if bloodFury == nil {
		t.Fatal("no Blood Fury aura on an Orc")
	}

	var heard []stats.Stats
	character.OnTemporaryStatsChanges = append(character.OnTemporaryStatsChanges, func(_ *Simulation, aura *Aura, change stats.Stats) {
		if aura != bloodFury {
			t.Errorf("callback names aura %s, want Blood Fury", aura.Label)
		}
		heard = append(heard, change)
	})

	before := character.GetStats()
	bloodFury.Activate(sim)
	during := character.GetStats()
	bloodFury.Deactivate(sim)
	after := character.GetStats()

	if len(heard) != 2 {
		t.Fatalf("callback heard %d changes, want one on gain and one on expire", len(heard))
	}
	if heard[0] != during.Subtract(before) {
		t.Errorf("gain reported %s, want the stats the multipliers added: %s",
			heard[0].FlatString(), during.Subtract(before).FlatString())
	}
	if heard[1] != after.Subtract(during) {
		t.Errorf("expire reported %s, want the stats the multipliers removed: %s",
			heard[1].FlatString(), after.Subtract(during).FlatString())
	}
	if after != before {
		t.Errorf("stats after expire %s, want the ones before gain %s", after.FlatString(), before.FlatString())
	}

	attackPower := before[stats.AttackPower]
	if attackPower <= 0 {
		t.Fatalf("attack power is %v before Blood Fury, nothing for 10%% to show on", attackPower)
	}
	if got := heard[0][stats.AttackPower]; math.Abs(got-attackPower*0.1) > 1 {
		t.Errorf("gain reported %v attack power, want 10%% of %v", got, attackPower)
	}
	if got := heard[1][stats.AttackPower]; got != -heard[0][stats.AttackPower] {
		t.Errorf("expire reported %v attack power, want the reverse of the gain's %v", got, heard[0][stats.AttackPower])
	}
}

// 16 Oct build: Ruthless (human, 2%) and Savage (orc, 1%) replaced the weapon specializations, so the
// crit no longer needs a sword or an axe and reaches melee, ranged and spells alike.
func TestCritChanceRacialsNeedNoWeapon(t *testing.T) {
	_, undead := setupRacialSim(proto.Race_RaceUndead)
	for race, bonus := range map[proto.Race]float64{proto.Race_RaceHuman: 2, proto.Race_RaceOrc: 1} {
		_, fw := setupRacialSim(race)
		for _, stat := range []stats.Stat{stats.PhysicalCritPercent, stats.SpellCritPercent} {
			if got, want := fw.GetStat(stat), undead.GetStat(stat)+bonus; !WithinToleranceFloat64(want, got, 1e-9) {
				t.Errorf("%s: stat %d %.2f unarmed, want %.2f", race, stat, got, want)
			}
		}
	}
}
