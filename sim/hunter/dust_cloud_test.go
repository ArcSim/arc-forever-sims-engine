package hunter

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// A Tallstrider keeps Dust Cloud (client 1265904, armor -505) on the target, so the hunter's own
// physical shots land harder than beside a Turtle, the same pet without it.
func TestTallstriderDustCloud(t *testing.T) {
	run := func(pet proto.HunterOptions_PetType) *proto.RaidSimResult {
		player := &proto.Player{
			Name: "bm", Class: proto.Class_ClassHunter, Race: proto.Race_RaceOrc, TalentsString: BeastMasteryTalents,
			Equipment: WeaponsOnly, DistanceFromTarget: 30,
			Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{ClassOptions: &proto.HunterOptions{
				Ammo: proto.HunterOptions_Doomshot, QuiverBonus: proto.HunterOptions_Speed15, PetType: pet,
				PetAttackSpeed: proto.HunterOptions_OneTwo, PetUptime: 1}}}},
			Rotation: core.GetAplRotation("../../ui/specs/hunter/dps/apls", "bm").Rotation,
		}
		raid := &proto.Raid{Parties: []*proto.Party{{Players: []*proto.Player{player}, Buffs: &proto.PartyBuffs{}}}, Buffs: &proto.RaidBuffs{}, Debuffs: &proto.Debuffs{}, NumActiveParties: 1}
		res := core.RunRaidSim(&proto.RaidSimRequest{Raid: raid, Encounter: core.MakeSingleTargetEncounter(0), SimOptions: &proto.SimOptions{Iterations: 20, RandomSeed: 1}})
		if res.Error != nil {
			t.Fatal(res.Error.Message)
		}
		return res
	}

	strider, turtle := run(proto.HunterOptions_Tallstrider), run(proto.HunterOptions_Turtle)

	var uptime float64
	for _, aura := range strider.EncounterMetrics.Targets[0].Auras {
		if aura.Id.GetSpellId() == spellData.DustCloudTriggered.Highest().ID {
			uptime = aura.UptimeSecondsAvg
		}
	}
	if fight := core.LongDuration; uptime < 0.9*float64(fight) {
		t.Errorf("Dust Cloud up %.0f s of %d s, want nearly all of it", uptime, fight)
	}

	s, u := strider.RaidMetrics.Parties[0].Players[0].Dps.Avg, turtle.RaidMetrics.Parties[0].Players[0].Dps.Avg
	t.Logf("hunter DPS beside a Tallstrider %.1f, a Turtle %.1f; Dust Cloud up %.0f s", s, u, uptime)
	if s < u*1.02 {
		t.Errorf("hunter DPS beside a Tallstrider %.1f, want at least 2%% over a Turtle's %.1f", s, u)
	}
}

// Sunder Armor, Expose Armor and Dust Cloud don't stack in Forever (DustCloudYieldsToSunderAndExpose):
// beside the raid's Sunder Armor or Expose Armor the Tallstrider leaves Dust Cloud alone. Faerie Fire
// and Curse of Recklessness are magic and don't stop it. With the switch off Dust Cloud goes back on
// top of Sunder, and the hunter gains its 505 armor.
func TestDustCloudYieldsToSunderAndExpose(t *testing.T) {
	run := func(pet proto.HunterOptions_PetType, debuffs *proto.Debuffs) (float64, float64) {
		player := &proto.Player{
			Name: "bm", Class: proto.Class_ClassHunter, Race: proto.Race_RaceOrc, TalentsString: BeastMasteryTalents,
			Equipment: WeaponsOnly, DistanceFromTarget: 30,
			Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{ClassOptions: &proto.HunterOptions{
				Ammo: proto.HunterOptions_Doomshot, QuiverBonus: proto.HunterOptions_Speed15, PetType: pet,
				PetAttackSpeed: proto.HunterOptions_OneTwo, PetUptime: 1}}}},
			Rotation: core.GetAplRotation("../../ui/specs/hunter/dps/apls", "bm").Rotation,
		}
		raid := &proto.Raid{Parties: []*proto.Party{{Players: []*proto.Player{player}, Buffs: &proto.PartyBuffs{}}}, Buffs: &proto.RaidBuffs{}, Debuffs: debuffs, NumActiveParties: 1}
		res := core.RunRaidSim(&proto.RaidSimRequest{Raid: raid, Encounter: core.MakeSingleTargetEncounter(0), SimOptions: &proto.SimOptions{Iterations: 20, RandomSeed: 1}})
		if res.Error != nil {
			t.Fatal(res.Error.Message)
		}
		var uptime float64
		for _, aura := range res.EncounterMetrics.Targets[0].Auras {
			if aura.Id.GetSpellId() == spellData.DustCloudTriggered.Highest().ID {
				uptime = aura.UptimeSecondsAvg
			}
		}
		return uptime, res.RaidMetrics.Parties[0].Players[0].Dps.Avg
	}
	fight := float64(core.LongDuration)

	for name, debuffs := range map[string]*proto.Debuffs{
		"Sunder Armor": {SunderArmor: true},
		"Expose Armor": {ExposeArmor: true},
	} {
		if uptime, _ := run(proto.HunterOptions_Tallstrider, debuffs); uptime > 0 {
			t.Errorf("beside %s Dust Cloud is up %.0f s, want 0", name, uptime)
		}
	}

	if uptime, _ := run(proto.HunterOptions_Tallstrider, &proto.Debuffs{FaerieFire: true, CurseOfRecklessness: true}); uptime < 0.9*fight {
		t.Errorf("beside Faerie Fire and Curse of Recklessness Dust Cloud is up %.0f s of %.0f, want nearly all of it", uptime, fight)
	}

	_, yielding := run(proto.HunterOptions_Tallstrider, &proto.Debuffs{SunderArmor: true})
	DustCloudYieldsToSunderAndExpose = false
	defer func() { DustCloudYieldsToSunderAndExpose = true }()
	uptime, stacking := run(proto.HunterOptions_Tallstrider, &proto.Debuffs{SunderArmor: true})
	if uptime < 0.9*fight {
		t.Errorf("switched off, beside Sunder Armor Dust Cloud is up %.0f s of %.0f, want nearly all of it", uptime, fight)
	}
	t.Logf("Tallstrider hunter beside Sunder Armor: %.1f with Dust Cloud held, %.1f with it stacking", yielding, stacking)
	if stacking <= yielding {
		t.Errorf("Dust Cloud stacking on Sunder gave %.1f, not above %.1f without it", stacking, yielding)
	}
}
