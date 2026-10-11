//go:build with_db

package hunter

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// The melee Survival reference's weapons (MythicSim's survival-hunter): two one-handers.
var DualWield = &proto.EquipmentSpec{
	Items: []*proto.ItemSpec{
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{Id: 13952}, // Iceblade Hacker
		{Id: 13368}, // Bonescraper
		{Id: 18713}, // Rhok'delar, Longbow of the Ancient Keepers
	},
}

func dualWieldSurvival(t *testing.T) *Hunter {
	player := &proto.Player{
		Name: "sv", Class: proto.Class_ClassHunter, Race: proto.Race_RaceOrc, TalentsString: SurvivalMeleeTalents,
		Equipment: DualWield,
		Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{ClassOptions: &proto.HunterOptions{
			Ammo: proto.HunterOptions_Doomshot, PetType: proto.HunterOptions_PetNone}}}},
	}
	raid := &proto.Raid{Parties: []*proto.Party{{Players: []*proto.Player{player}, Buffs: &proto.PartyBuffs{}}}, Buffs: &proto.RaidBuffs{}, Debuffs: &proto.Debuffs{}, NumActiveParties: 1}
	env, _, _ := core.NewEnvironment(raid, core.MakeSingleTargetEncounter(0), false, true)
	hunter := env.Raid.Parties[0].Players[0].(HunterAgent).GetHunter()
	if !hunter.AutoAttacks.IsDualWielding {
		t.Fatal("the test hunter is not dual wielding")
	}
	return hunter
}

// Savage Strikes rank 2 (19159) is +4% crit on every melee ability and on no auto attack. Discord
// read it as +2% from the per-attempt crit column: a yellow hit rolls crit only once it has landed
// (two-roll), so ~6.5% boss dodges dilute its crits per attempt, while a white swing's crit share of
// the one-roll table is the sheet's.
func TestSavageStrikesRank2OnEveryMeleeAbility(t *testing.T) {
	hunter := dualWieldSurvival(t)
	if hunter.Talents.SavageStrikes != 2 {
		t.Fatalf("the melee build takes Savage Strikes %d/2", hunter.Talents.SavageStrikes)
	}
	for _, spell := range []*core.Spell{hunter.RaptorStrike, hunter.MongooseBite, hunter.StriderKick, hunter.WingClip, hunter.LaceratingStrikes} {
		if spell.BonusCritPercent != 4 {
			t.Errorf("%v: Savage Strikes gives %v%% crit, want 4%%", spell.ActionID, spell.BonusCritPercent)
		}
	}
	for _, auto := range []*core.Spell{hunter.AutoAttacks.MHAuto(), hunter.AutoAttacks.OHAuto()} {
		if auto.BonusCritPercent != 0 {
			t.Errorf("%v: auto attack has %v%% bonus crit", auto.ActionID, auto.BonusCritPercent)
		}
	}
}

// Predator's Edge's "+50% Off Hand weapon damage" (1310627, A_MOD_OFFHAND_DAMAGE_PCT) is the off-hand
// weapon's alone. Mongoose Bite strikes with the main hand, but carried the off-hand special proc
// mask too, and took the +50%.
func TestPredatorsEdgeOffHandBonusSkipsMongooseBite(t *testing.T) {
	hunter := dualWieldSurvival(t)
	if got := hunter.AutoAttacks.OHAuto().DamageMultiplier; got != 1.5 {
		t.Errorf("off-hand auto attack multiplier %v, want Predator's Edge's 1.5", got)
	}
	for _, spell := range []*core.Spell{hunter.MongooseBite, hunter.StriderKick, hunter.WingClip, hunter.AutoAttacks.MHAuto()} {
		if spell.DamageMultiplier != 1 {
			t.Errorf("%v: damage multiplier %v, want 1 (the off-hand bonus is the off hand's)", spell.ActionID, spell.DamageMultiplier)
		}
	}
}
