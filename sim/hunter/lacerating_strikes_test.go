package hunter

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// Savage Strikes (19159) and Predator's Edge (1310627) name the Lacerating Strikes bleed (1310536) in
// their class masks, so its ticks get the same crit chance and crit damage as Mongoose Bite.
func TestLaceratingStrikesGetsSavageStrikesAndPredatorsEdge(t *testing.T) {
	player := &proto.Player{
		Name: "sv", Class: proto.Class_ClassHunter, Race: proto.Race_RaceNightElf, TalentsString: SurvivalTalents,
		Equipment: WeaponsOnly,
		Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{ClassOptions: &proto.HunterOptions{
			Ammo: proto.HunterOptions_Doomshot, PetType: proto.HunterOptions_Cat}}}},
	}
	raid := &proto.Raid{Parties: []*proto.Party{{Players: []*proto.Player{player}, Buffs: &proto.PartyBuffs{}}}, Buffs: &proto.RaidBuffs{}, Debuffs: &proto.Debuffs{}, NumActiveParties: 1}
	env, _, _ := core.NewEnvironment(raid, core.MakeSingleTargetEncounter(0), false, true)
	hunter := env.Raid.Parties[0].Players[0].(HunterAgent).GetHunter()

	bite, bleed := hunter.MongooseBite, hunter.LaceratingStrikes
	if bleed.BonusCritPercent == 0 || bleed.BonusCritPercent != bite.BonusCritPercent {
		t.Errorf("bleed bonus crit %v%%, Mongoose Bite %v%%", bleed.BonusCritPercent, bite.BonusCritPercent)
	}
	if bleed.CritMultiplierAdditive == 0 || bleed.CritMultiplierAdditive != bite.CritMultiplierAdditive {
		t.Errorf("bleed crit damage bonus %v, Mongoose Bite %v", bleed.CritMultiplierAdditive, bite.CritMultiplierAdditive)
	}
}
