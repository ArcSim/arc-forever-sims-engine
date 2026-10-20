package protection

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
	"github.com/wowsims/forever/sim/warrior"
	googleProto "google.golang.org/protobuf/proto"
)

func threatRequest(targets int) *proto.RaidSimRequest {
	encounter := core.MakeSingleTargetEncounter(0)
	for len(encounter.Targets) < targets {
		encounter.Targets = append(encounter.Targets, core.NewDefaultTarget())
	}
	player := &proto.Player{
		Name: "Threat", Race: proto.Race_RaceOrc, Class: proto.Class_ClassWarrior,
		Equipment:     googleProto.Clone(core.GetGearSet("../../../ui/specs/warrior/protection/gear_sets", "preraid").GearSet).(*proto.EquipmentSpec),
		TalentsString: DefaultProtectionTalents, Spec: DefaultOptions,
		Rotation: core.GetAplRotation("../../../ui/specs/warrior/protection/apls", "protection").Rotation,
	}
	return &proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  encounter,
		SimOptions: &proto.SimOptions{RandomSeed: 1, Iterations: 20},
	}
}

func threatSpells(t *testing.T) (*core.Simulation, *core.Character, map[string]*core.Spell) {
	t.Helper()
	sim := core.NewSim(threatRequest(1), simsignals.CreateSignals())
	sim.Reset()
	war := sim.Raid.Parties[0].Players[0].GetCharacter()
	spells := map[string]*core.Spell{}
	for name, mask := range map[string]int64{
		"Revenge": warrior.SpellMaskRevenge, "Heroic Strike": warrior.SpellMaskHeroicStrike,
		"Cleave": warrior.SpellMaskCleave, "Sunder Armor": warrior.SpellMaskSunderArmor,
	} {
		for _, spell := range war.Spellbook {
			if spell.Matches(mask) && spell.ActionID.Tag == 0 {
				spells[name] = spell
			}
		}
		if spells[name] == nil {
			t.Fatalf("the warrior has no %s", name)
		}
	}
	return sim, war, spells
}

// Patches 183 to 185 (wowsims/forever's in-game threat captures): Revenge is 2.25 times its damage,
// Heroic Strike adds 2.54 and Cleave 1.53 a spell level, and Sunder Armor adds 0.05 of the attack power
// to the client's 206. Off, the old values.
func TestWarriorThreatCaptures(t *testing.T) {
	type want struct{ revengeMult, revengeFlat, heroicStrike, cleave float64 }
	for _, c := range []struct {
		on   bool
		want want
	}{
		{true, want{2.25, 0, 2.54 * 60, 1.53 * 60}},
		{false, want{1, 120, 173, 100}},
	} {
		func() {
			defer func(a, b, s bool) {
				warrior.RevengeThreatFromDamage, warrior.HeroicStrikeCleaveThreatPerLevel, warrior.SunderArmorThreatFromAttackPower = a, b, s
			}(warrior.RevengeThreatFromDamage, warrior.HeroicStrikeCleaveThreatPerLevel, warrior.SunderArmorThreatFromAttackPower)
			warrior.RevengeThreatFromDamage, warrior.HeroicStrikeCleaveThreatPerLevel, warrior.SunderArmorThreatFromAttackPower = c.on, c.on, c.on

			sim, war, spells := threatSpells(t)
			near := func(name string, got, want float64) {
				if math.Abs(got-want) > 0.01 {
					t.Errorf("switches %v: %s is %.3f, want %.3f", c.on, name, got, want)
				}
			}
			near("Revenge's threat multiplier", spells["Revenge"].ThreatMultiplier, c.want.revengeMult)
			near("Revenge's flat threat", spells["Revenge"].FlatThreatBonus, c.want.revengeFlat)
			near("Heroic Strike's flat threat", spells["Heroic Strike"].FlatThreatBonus, c.want.heroicStrike)
			near("Cleave's flat threat", spells["Cleave"].FlatThreatBonus, c.want.cleave)

			target := sim.Encounter.AllTargetUnits[0]
			sunder := spells["Sunder Armor"]
			sunder.SkipCastAndApplyEffects(sim, target)
			wantSunder := 206.0
			if c.on {
				wantSunder += 0.05 * sunder.MeleeAttackPower(target)
			}
			near("Sunder Armor's flat threat", sunder.FlatThreatBonus, wantSunder)
			if c.on && war.GetStat(stats.AttackPower) < 100 {
				t.Errorf("the preraid warrior has %.0f attack power; the Sunder check needs some", war.GetStat(stats.AttackPower))
			}
		}()
	}
}

// Patch 182 (wowsims/forever 913a5bcf69): rage gained is 5 threat times the warrior's threat multiplier
// (Defensive Stance x1.3, Defiance), split between the targets in combat. Off, 5 per rage on each target.
func TestRageGainThreat(t *testing.T) {
	rageThreat := func(on bool, targets int) []float64 {
		defer func(old bool) { core.RageGainThreatMultiplied = old }(core.RageGainThreatMultiplied)
		core.RageGainThreatMultiplied = on
		result := core.RunRaidSim(threatRequest(targets))
		if result.Error != nil {
			t.Fatalf("sim failed: %s", result.Error.Message)
		}
		var out []float64
		for _, action := range result.RaidMetrics.Parties[0].Players[0].Actions {
			if action.Id.GetOtherId() == proto.OtherAction_OtherActionRageGain {
				for _, target := range action.Targets {
					out = append(out, target.Threat)
				}
			}
		}
		// The metrics carry a slot past the enemies, which takes no threat.
		if len(out) < targets {
			t.Fatalf("rage gain threat on %d targets, want %d", len(out), targets)
		}
		return out[:targets]
	}

	on, off := rageThreat(true, 1)[0], rageThreat(false, 1)[0]
	if ratio := on / off; ratio < 1.29 || ratio > 1.5 {
		t.Errorf("one target: rage gain threat x%.3f of the unmultiplied, want Defensive Stance's 1.3 to 1.495 with Defiance", ratio)
	}
	two := rageThreat(true, 2)
	twoOff := rageThreat(false, 2)
	if math.Abs(two[0]-two[1]) > 1e-6 || math.Abs(twoOff[0]-twoOff[1]) > 1e-6 {
		t.Errorf("two targets: rage gain threat %v (off %v), want the same on each", two, twoOff)
	}
	if ratio := two[0] / twoOff[0]; ratio < 0.645 || ratio > 0.75 {
		t.Errorf("two targets: rage gain threat x%.3f of the unsplit, want half of 1.3 to 1.495", ratio)
	}
}
