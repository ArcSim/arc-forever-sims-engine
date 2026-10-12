package hunter

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

func petSim(t *testing.T, pet proto.HunterOptions_PetType, talents, apl string, distance float64) *proto.RaidSimResult {
	t.Helper()
	player := &proto.Player{
		Name: "h", Class: proto.Class_ClassHunter, Race: proto.Race_RaceOrc, TalentsString: talents,
		Equipment: WeaponsOnly, DistanceFromTarget: distance,
		Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{ClassOptions: &proto.HunterOptions{
			Ammo: proto.HunterOptions_Doomshot, QuiverBonus: proto.HunterOptions_Speed15, PetType: pet,
			PetAttackSpeed: proto.HunterOptions_OneTwo, PetUptime: 1}}}},
		Rotation: core.GetAplRotation("../../ui/specs/hunter/dps/apls", apl).Rotation,
	}
	raid := &proto.Raid{Parties: []*proto.Party{{Players: []*proto.Player{player}, Buffs: &proto.PartyBuffs{}}}, Buffs: &proto.RaidBuffs{}, Debuffs: &proto.Debuffs{}, NumActiveParties: 1}
	res := core.RunRaidSim(&proto.RaidSimRequest{Raid: raid, Encounter: core.MakeSingleTargetEncounter(0), SimOptions: &proto.SimOptions{Iterations: 20, RandomSeed: 1}})
	if res.Error != nil {
		t.Fatal(res.Error.Message)
	}
	return res
}

func petCasts(res *proto.RaidSimResult, id int32) float64 {
	var casts int32
	for _, pet := range res.RaidMetrics.Parties[0].Players[0].Pets {
		for _, action := range pet.Actions {
			if action.Id.GetSpellId() == id {
				for _, target := range action.Targets {
					casts += target.Casts
				}
			}
		}
	}
	return float64(casts) / 20
}

// Patch 118: with PetPoolsFocusForCooldowns the Cat holds Claw while it would leave Bite (35 focus, 10 sec
// cooldown) short, so Bite goes on cooldown. Without it Claw (25) spends every 25 focus and Bite lands once,
// off the full bar the pet starts with (Kragtor's sim: one Bite in 240 s).
func TestCatBitesOnCooldown(t *testing.T) {
	defer func(old bool) { PetPoolsFocusForCooldowns = old }(PetPoolsFocusForCooldowns)
	fight := core.LongDuration

	PetPoolsFocusForCooldowns = false
	unpooled := petCasts(petSim(t, proto.HunterOptions_Cat, MarksmanshipTalents, "mm", 30), 17261)
	PetPoolsFocusForCooldowns = true
	pooled := petSim(t, proto.HunterOptions_Cat, MarksmanshipTalents, "mm", 30)
	bites, claws := petCasts(pooled, 17261), petCasts(pooled, 3009)

	t.Logf("Bites a fight: %.1f unpooled, %.1f pooled (%d s fight, Claws %.1f)", unpooled, bites, fight, claws)
	if unpooled > 2 {
		t.Errorf("unpooled Cat Bit %.1f times, want the opener only", unpooled)
	}
	if want := 0.9 * float64(fight) / 10; bites < want {
		t.Errorf("pooled Cat Bit %.1f times in %d s, want at least %.0f (every 10 s)", bites, fight, want)
	}
	// Focus still all goes somewhere: 10 a second, 35 a Bite, 25 a Claw.
	if spent := 35*bites + 25*claws - 100; spent < 9.3*float64(fight) {
		t.Errorf("pooled Cat spent %.0f focus in %d s, want about 10 a second", spent, fight)
	}
}

// Patch 119: the Wolf howls (client 24597), keeping the buff up, and the melee Hunter beside it gains its
// attack power; a ranged Hunter at 30 yards is out of its 15 yard reach.
func TestWolfFuriousHowl(t *testing.T) {
	defer func(old bool) { WolfCastsFuriousHowl = old }(WolfCastsFuriousHowl)
	howl := spellData.FuriousHowlTriggered.Highest()
	fight := core.LongDuration

	ownerUptime := func(res *proto.RaidSimResult) float64 {
		for _, aura := range res.RaidMetrics.Parties[0].Players[0].Auras {
			if aura.Id.GetSpellId() == howl.ID {
				return aura.UptimeSecondsAvg
			}
		}
		return 0
	}

	WolfCastsFuriousHowl = true
	melee := petSim(t, proto.HunterOptions_Wolf, SurvivalTalents, "sv_melee", core.MaxMeleeRange)
	ranged := petSim(t, proto.HunterOptions_Wolf, MarksmanshipTalents, "mm", 30)
	WolfCastsFuriousHowl = false
	without := petSim(t, proto.HunterOptions_Wolf, SurvivalTalents, "sv_melee", core.MaxMeleeRange)

	casts, up := petCasts(melee, howl.ID), ownerUptime(melee)
	t.Logf("Furious Howl: %.1f casts, up %.0f s of %d on the melee Hunter, %.0f s on the ranged one; DPS %.1f with, %.1f without",
		casts, up, fight, ownerUptime(ranged), melee.RaidMetrics.Dps.Avg, without.RaidMetrics.Dps.Avg)
	if casts < float64(fight)/60 {
		t.Errorf("Wolf howled %.1f times in %d s, want the 60 s buff kept up", casts, fight)
	}
	if up < 0.9*float64(fight) {
		t.Errorf("Furious Howl up %.0f s of %d on the melee Hunter, want nearly all", up, fight)
	}
	if r := ownerUptime(ranged); r != 0 {
		t.Errorf("Furious Howl up %.0f s on a Hunter 30 yards out, want none", r)
	}
	if melee.RaidMetrics.Dps.Avg <= without.RaidMetrics.Dps.Avg {
		t.Errorf("melee Hunter beside a howling Wolf %.1f DPS, want more than without %.1f", melee.RaidMetrics.Dps.Avg, without.RaidMetrics.Dps.Avg)
	}
}

// Patch 120: Strider Kick (1317257) speeds the Hunter up 30% for 3 sec after each kick. The weapons-only
// Hunter runs short of mana, so it kicks less than every 8 sec: the buff is up a sixth to 3/8 of the fight.
func TestStriderKickSpeedsUp(t *testing.T) {
	res := petSim(t, proto.HunterOptions_PetNone, SurvivalTalents, "sv_melee", core.MaxMeleeRange)
	var up float64
	for _, aura := range res.RaidMetrics.Parties[0].Players[0].Auras {
		if aura.Id.GetSpellId() == 1317257 {
			up = aura.UptimeSecondsAvg
		}
	}
	fight := float64(core.LongDuration)
	t.Logf("Strider Kick speed up %.0f s of %.0f", up, fight)
	if up < 0.15*fight || up > 0.4*fight {
		t.Errorf("Strider Kick speed up %.0f s of %.0f, want a sixth to 3/8", up, fight)
	}
}
