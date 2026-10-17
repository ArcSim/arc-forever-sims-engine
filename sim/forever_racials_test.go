package sim

import (
	"math"
	"testing"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/simsignals"
	"github.com/wowsims/classic/sim/core/stats"
)

// Forever racials and raid debuffs checked on several classes, so they hold for every class and not only the
// hunter they were first built for.

func newForeverCharacter(t *testing.T, race proto.Race, class proto.Class, debuffs *proto.Debuffs) (*core.Simulation, *core.Character) {
	t.Helper()
	player := &proto.Player{
		Name:      "Player",
		Race:      race,
		Class:     class,
		Equipment: &proto.EquipmentSpec{},
		Consumes:  &proto.Consumes{},
		Buffs:     &proto.IndividualBuffs{},
	}
	switch class {
	case proto.Class_ClassMage:
		player.Spec = &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{}}}
	case proto.Class_ClassWarrior:
		player.Spec = &proto.Player_Warrior{Warrior: &proto.Warrior{Options: &proto.Warrior_Options{}}}
	case proto.Class_ClassWarlock:
		player.Spec = &proto.Player_Warlock{Warlock: &proto.Warlock{Options: &proto.WarlockOptions{}}}
	default:
		t.Fatalf("no spec for %v", class)
	}
	if debuffs == nil {
		debuffs = &proto.Debuffs{}
	}
	sim := core.NewSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, debuffs),
		Encounter:  core.MakeSingleTargetEncounter(0),
		SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 1},
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim, sim.Raid.Parties[0].Players[0].GetCharacter()
}

// a damaging spell the class pays for (and whether it is periodic)
func classAbility(char *core.Character, periodic bool) *core.Spell {
	for _, spell := range char.Spellbook {
		if !spell.Flags.Matches(core.SpellFlagAPL) || spell.Cost == nil || !spell.ProcMask.Matches(core.ProcMaskMeleeOrRangedSpecial|core.ProcMaskSpellDamage) {
			continue
		}
		if (len(spell.Dots()) > 0 || spell.AOEDot() != nil) == periodic {
			return spell
		}
	}
	return nil
}

// Earth Fury (Dwarf) and Endurance (Tauren): +2% critical damage and healing on attacks and spells.
func TestForeverCritDamageRacials(t *testing.T) {
	for _, tc := range []struct {
		race proto.Race
		want float64
	}{
		{proto.Race_RaceDwarf, 1.02},
		{proto.Race_RaceTauren, 1.02},
		{proto.Race_RaceHuman, 1},
	} {
		_, mage := newForeverCharacter(t, tc.race, proto.Class_ClassMage, nil)
		spell := classAbility(mage, false)
		if spell == nil {
			t.Fatal("no mage damage spell")
		}
		at := &core.AttackTable{Attacker: &mage.Unit, CritMultiplier: 1}
		if got, want := spell.CritMultiplier(at), 1+(1.5*tc.want-1)*spell.CritDamageBonus; math.Abs(got-want) > 1e-9 {
			t.Errorf("%v spell crit multiplier %.4f, want %.4f", tc.race, got, want)
		}
		if got := mage.PseudoStats.MeleeCritMultiplier; math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("%v melee crit multiplier %.4f, want %.4f", tc.race, got, tc.want)
		}
	}
}

// Blood Fury (Forever 20572): +10% attack power and spell power (damage and healing).
func TestBloodFurySpellPower(t *testing.T) {
	sim, mage := newForeverCharacter(t, proto.Race_RaceOrc, proto.Class_ClassMage, nil)
	mage.AddStatsDynamic(sim, stats.Stats{stats.SpellPower: 100, stats.SpellDamage: 50, stats.HealingPower: 40, stats.AttackPower: 200})
	before := mage.GetStats()
	aura := mage.GetAura("Blood Fury")
	if aura == nil {
		t.Fatal("no Blood Fury")
	}
	aura.Activate(sim)
	for _, s := range []stats.Stat{stats.SpellPower, stats.SpellDamage, stats.HealingPower, stats.AttackPower} {
		if got, want := mage.GetStat(s), before[s]*1.1; math.Abs(got-want) > 1e-6 {
			t.Errorf("Blood Fury %v %.1f, want %.1f", s, got, want)
		}
	}
}

// Eureka! (Gnome): 8 s of 10% more damage and 10% lower cost on damaging abilities; Rage users keep their
// periodic abilities out of it.
func TestEureka(t *testing.T) {
	for _, tc := range []struct {
		class         proto.Class
		periodicGains bool
	}{
		{proto.Class_ClassMage, true},
		{proto.Class_ClassWarrior, false},
	} {
		sim, char := newForeverCharacter(t, proto.Race_RaceGnome, tc.class, nil)
		var eureka *core.Aura
		for _, label := range []string{"Eureka!"} {
			eureka = char.GetAura(label)
		}
		if eureka == nil {
			t.Fatalf("%v: no Eureka!", tc.class)
		}
		if eureka.Duration.Seconds() != 8 {
			t.Errorf("%v: Eureka! lasts %v, want 8s", tc.class, eureka.Duration)
		}
		direct := classAbility(char, false)
		if direct == nil {
			t.Fatalf("%v: no direct damage ability", tc.class)
		}
		dmg, cost := direct.DamageMultiplierAdditive, direct.Cost.Multiplier
		periodic := classAbility(char, true)
		var pDmg float64
		if periodic != nil {
			pDmg = periodic.DamageMultiplierAdditive
		}
		eureka.Activate(sim)
		if math.Abs(direct.DamageMultiplierAdditive-dmg-0.1) > 1e-9 || direct.Cost.Multiplier != cost-10 {
			t.Errorf("%v: Eureka! on %v: damage %+.2f cost %+d, want +0.10 and -10", tc.class, direct.ActionID, direct.DamageMultiplierAdditive-dmg, direct.Cost.Multiplier-cost)
		}
		if periodic != nil {
			if gain := periodic.DamageMultiplierAdditive - pDmg; (gain > 0) != tc.periodicGains {
				t.Errorf("%v: Eureka! on periodic %v: %+.2f, want gain %v", tc.class, periodic.ActionID, gain, tc.periodicGains)
			}
		}
		eureka.Deactivate(sim)
		if math.Abs(direct.DamageMultiplierAdditive-dmg) > 1e-9 || direct.Cost.Multiplier != cost {
			t.Errorf("%v: Eureka! did not undo itself", tc.class)
		}
	}
}

// Touch of the Grave (Undead): the drain takes 5% of the player's maximum health and ignores damage bonuses.
func TestTouchOfTheGraveDrain(t *testing.T) {
	sim, lock := newForeverCharacter(t, proto.Race_RaceUndead, proto.Class_ClassWarlock, nil)
	if lock.GetAura("Touch of the Grave") == nil {
		t.Fatal("no Touch of the Grave")
	}
	drain := lock.GetSpell(core.ActionID{SpellID: 1260189})
	if drain == nil {
		t.Fatal("no drain spell")
	}
	target := lock.CurrentTarget
	// shadow bonuses on the target must not raise it
	target.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= 1.5
	drain.Cast(sim, target)
	if got, want := drain.SpellMetrics[target.UnitIndex].TotalDamage, 0.05*lock.MaxHealth(); math.Abs(got-want) > 1e-6 {
		t.Errorf("drain %.1f, want %.1f (5%% of %.0f health)", got, want, lock.MaxHealth())
	}
}

// Improved Shadow Bolt, Shadow Weaving, Improved Scorch, Winter's Chill and Stormstrike only raise damage taken from
// the player who applied them: switched on as debuffs from someone else, they put nothing on the target.
func TestCasterOnlyDebuffsFromOthers(t *testing.T) {
	_, warrior := newForeverCharacter(t, proto.Race_RaceHuman, proto.Class_ClassWarrior, &proto.Debuffs{
		ImprovedShadowBolt: true, ShadowWeaving: true, ImprovedScorch: true, WintersChill: true, Stormstrike: true,
	})
	for _, label := range []string{"Improved Shadow Bolt", "Shadow Weaving", "Improved Scorch", "Winter's Chill", "Stormstrike"} {
		if warrior.CurrentTarget.GetAura(label) != nil {
			t.Errorf("%s from another player is on the target", label)
		}
	}
}
