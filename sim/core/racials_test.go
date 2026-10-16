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

// 16 October 2026 beta notes: Eureka! is "an 8 second buff with no charges. Same effect on damage and
// ability costs."
func TestGnomeWarrior(t *testing.T) {
	sim, fw := setupRacialSim(proto.Race_RaceGnome)

	if fw.MaximumRage() != 105 {
		t.Fatalf("Expansive Mind should raise maximum Rage to 105, got %0.1f", fw.MaximumRage())
	}

	castRacial(t, sim, fw, 1259813)
	eureka := fw.GetAura("Eureka!")
	if !eureka.IsActive() || eureka.Duration != time.Second*8 || eureka.MaxStacks != 0 {
		t.Fatalf("Eureka! should be an 8 sec buff with no charges, got %s and %d stacks", eureka.Duration, eureka.MaxStacks)
	}

	ability := fw.RegisterSpell(SpellConfig{
		ActionID:         ActionID{SpellID: 1},
		ClassSpellMask:   1,
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
	itemSpell.Cast(sim, fw.CurrentTarget)
	if itemSpell.DamageMultiplier != 1 {
		t.Fatalf("Eureka! should ignore spells that are not class abilities")
	}
	if !WithinToleranceFloat64(1.1, ability.DamageMultiplier, 0.0001) {
		t.Fatalf("Eureka! should raise ability damage by 10%%, got x%0.2f", ability.DamageMultiplier)
	}

	for range 5 {
		ability.Cast(sim, fw.CurrentTarget)
	}
	if !eureka.IsActive() || !WithinToleranceFloat64(1.1, ability.DamageMultiplier, 0.0001) {
		t.Fatalf("Eureka! has no charges, so casts should not end it")
	}

	eureka.Deactivate(sim)
	if !WithinToleranceFloat64(1, ability.DamageMultiplier, 0.0001) {
		t.Fatalf("Eureka! should drop its damage bonus when it ends, got x%0.2f", ability.DamageMultiplier)
	}
}

// 16 October 2026 beta notes: "Endurance: Now grants 2% increased critical damage and healing (was 1% hit
// chance)", and Dwarf "Mace Specialization: Replaced with Earth Fury: 2% increased critical damage and
// healing."
func TestCritDamageRacials(t *testing.T) {
	_, human := setupRacialSim(proto.Race_RaceHuman)

	for _, race := range []proto.Race{proto.Race_RaceTauren, proto.Race_RaceDwarf} {
		_, fw := setupRacialSim(race)
		if got := fw.GetStat(stats.PhysicalHitPercent); got != human.GetStat(stats.PhysicalHitPercent) {
			t.Errorf("%s: hit %0.2f%%, want no racial hit", race, got)
		}
		if !WithinToleranceFloat64(1.02, fw.PseudoStats.CritDamageMultiplier, 1e-9) {
			t.Errorf("%s: crit damage x%0.4f, want x1.02", race, fw.PseudoStats.CritDamageMultiplier)
		}
	}
	if human.PseudoStats.CritDamageMultiplier != 1 {
		t.Errorf("a Human has no crit damage racial, got x%0.4f", human.PseudoStats.CritDamageMultiplier)
	}
}

// 16 October 2026 beta notes: Human Ruthless (+2% crit) and Orc Savage (+1% crit) replace the weapon
// specializations, with no weapon needed.
func TestCritRacials(t *testing.T) {
	_, nightElf := setupRacialSim(proto.Race_RaceNightElf)
	for race, want := range map[proto.Race]float64{proto.Race_RaceHuman: 2, proto.Race_RaceOrc: 1} {
		without := racialWarrior(race, disableRacials)
		with := racialWarrior(race, nil)
		if d := with.GetStat(stats.PhysicalCritPercent) - without.GetStat(stats.PhysicalCritPercent); d != want {
			t.Errorf("%s: racials add %0.2f%% physical crit with no weapon, want %0.0f%%", race, d, want)
		}
		if d := with.GetStat(stats.SpellCritPercent) - without.GetStat(stats.SpellCritPercent); d != want {
			t.Errorf("%s: racials add %0.2f%% spell crit, want %0.0f%%", race, d, want)
		}
	}
	if nightElf.PseudoStats.CritDamageMultiplier != 1 {
		t.Errorf("a Night Elf has no crit damage racial")
	}
}

// 16 October 2026 beta notes: Blood Fury's cooldown is 3 minutes (was 2).
func TestBloodFuryCooldown(t *testing.T) {
	_, fw := setupRacialSim(proto.Race_RaceOrc)
	if got := fw.GetSpell(ActionID{SpellID: 20572}).CD.Duration; got != time.Minute*3 {
		t.Fatalf("Blood Fury's cooldown is %s, want 3m", got)
	}
}

func TestTouchOfTheGraveUsesTheMeleeVariant(t *testing.T) {
	_, fw := setupRacialSim(proto.Race_RaceUndead)

	if aura := fw.GetAura("Touch of the Grave"); aura == nil || aura.ActionIDForProc.SpellID != 1260189 {
		t.Fatalf("An undead warrior should carry the melee Touch of the Grave")
	} else if aura.OnSpellHitDealt == nil {
		t.Fatalf("Touch of the Grave listens for no hit, so it never procs")
	} else if aura.Icd == nil || aura.Icd.Duration != time.Second*25 {
		t.Fatalf("Touch of the Grave should have a 25 sec internal cooldown (16 October 2026 beta notes)")
	}
	drain := fw.GetSpell(ActionID{SpellID: 1260198})
	if drain == nil {
		t.Fatalf("Touch of the Grave's drain is not registered")
	}
	if !drain.Flags.Matches(SpellFlagIgnoreModifiers) {
		t.Fatalf("Touch of the Grave no longer scales with Shadow increases on the target or the caster")
	}
}

func TestSkyborneRacials(t *testing.T) {
	_, human := setupRacialSim(proto.Race_RaceHuman)

	for _, race := range []proto.Race{proto.Race_RaceHighOrderSkyborne, proto.Race_RaceWindshaperSkyborne} {
		_, fw := setupRacialSim(race)

		swingRatio := human.AutoAttacks.MainhandSwingSpeed().Seconds() / fw.AutoAttacks.MainhandSwingSpeed().Seconds()
		if !WithinToleranceFloat64(1.03, swingRatio, 0.001) {
			t.Errorf("%s: Wind Blessed should swing 3%% faster (16 October 2026 beta notes), got x%0.4f", race, swingRatio)
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
