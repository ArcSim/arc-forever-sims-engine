package hunter

import (
	"fmt"
	"os"
	"strings"
	"testing"

	_ "github.com/wowsims/classic/sim/common"
	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
)

// Local A/B harness (not a regression test): DW vs 2H under matched talents and rotations.
// Run: go test ./sim/hunter -run TestForeverAB -v --tags=with_db -count=1

var abConsumes = &proto.Consumes{
	AgilityElixir:     proto.AgilityElixir_ElixirOfTheMongoose,
	Alcohol:           proto.Alcohol_AlcoholRumseyRumBlackLabel,
	AttackPowerBuff:   proto.AttackPowerBuff_WinterfallFirewater,
	DefaultConjured:   proto.Conjured_ConjuredDemonicRune,
	DefaultPotion:     proto.Potions_MajorManaPotion,
	DragonBreathChili: true,
	Food:              proto.Food_FoodSmokedDesertDumpling,
	HealthElixir:      proto.HealthElixir_ElixirOfFortitude,
	MainHandImbue:     proto.WeaponImbue_DenseSharpeningStone,
	OffHandImbue:      proto.WeaponImbue_DenseSharpeningStone,
	ManaRegenElixir:   proto.ManaRegenElixir_MagebloodPotion,
	MiscConsumes:      &proto.MiscConsumes{RaptorPunch: true},
	SpellPowerBuff:    proto.SpellPowerBuff_GreaterArcaneElixir,
	StrengthBuff:      proto.StrengthBuff_ElixirOfGiants,
	ZanzaBuff:         proto.ZanzaBuff_GroundScorpokAssay,
}

var abRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance:     true,
	BattleShout:          proto.TristateEffect_TristateEffectRegular,
	DivineSpirit:         true,
	FireResistanceAura:   true,
	FireResistanceTotem:  true,
	GiftOfTheWild:        proto.TristateEffect_TristateEffectRegular,
	LeaderOfThePack:      true,
	ManaSpringTotem:      proto.TristateEffect_TristateEffectRegular,
	StrengthOfEarthTotem: proto.TristateEffect_TristateEffectRegular,
	TrueshotAura:         true,
	WindfuryTotem:        true,
}

var abIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:  true,
	BlessingOfMight:  proto.TristateEffect_TristateEffectRegular,
	BlessingOfWisdom: proto.TristateEffect_TristateEffectRegular,
}

var abDebuffs = &proto.Debuffs{
	CurseOfElements:   true,
	FaerieFire:        true,
	GiftOfArthas:      true,
	HuntersMark:       proto.TristateEffect_TristateEffectRegular,
	JudgementOfWisdom: true,
	SunderArmor:       true,
}

var abOptions = &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{
	Ammo:               proto.Hunter_Options_ThoriumHeadedArrow,
	QuiverBonus:        proto.Hunter_Options_Speed15,
	PetType:            proto.Hunter_Options_Cat,
	PetUptime:          1,
	PetAttackSpeed:     proto.Hunter_Options_OneTwo,
	LaceratingRollover: true,
}}}

type abCase struct {
	label, gear, talents, apl string
}

var abSpellNames = map[int32]string{
	14266: "Raptor Strike", 14271: "Mongoose Bite", 1317257: "Strider Kick", 14305: "Immolation Trap", 14317: "Explosive Trap",
	1293241: "Summon Hawk", 14287: "Arcane Shot", 2643: "Multi-Shot", 25295: "Serpent Sting",
	14268: "Wing Clip", 1310536: "Lacerating", 8515: "Windfury Totem", 3045: "Rapid Fire", 20572: "Blood Fury",
}

func TestForeverAB(t *testing.T) {
	if os.Getenv("AB_NO_RESTART") != "" {
		core.ForeverMeleeRestartsAutoShot = false
	}
	iters := int32(5000)
	if v := os.Getenv("AB_ITERS"); v != "" {
		var n int32
		fmt.Sscanf(v, "%d", &n)
		if n > 0 {
			iters = n
		}
	}
	twoH, dw := "2h_weave", "surv_dw"
	if v := os.Getenv("AB_GEAR_DW"); v != "" {
		dw = v
	}
	if v := os.Getenv("AB_GEAR_2H"); v != "" {
		twoH = v
	}
	dwAPL := ""
	if v := os.Getenv("AB_APL_DW"); v != "" {
		dwAPL = v
	}
	weaveAPL := "2HWeave"
	if v := os.Getenv("AB_APL"); v != "" {
		weaveAPL = v
	}
	tal2H := "-00531500114-500230230050220151"
	talDWHawk := "53200005001-0050052-5002302300500201"
	talDWLS := "5-005005201-500230230050222151"
	cases := []abCase{
		{"A  2H gear + 2H talents + weave APL (site default)", twoH, tal2H, weaveAPL},
		{"B  DW gear + 2H talents + weave APL (pure weapon swap)", dw, tal2H, weaveAPL},
	}
	cases = append(cases, []abCase{
		{"C  DW gear + DW LS talents + weave APL", dw, talDWLS, weaveAPL},
		{"D  DW gear + DW Hawk talents + SurvDW melee-only APL (shipped DW preset)", dw, talDWHawk, "SurvDW"},
		{"E  2H gear + 2H talents + SurvDW melee-only APL", twoH, tal2H, "SurvDW"},
	}...)
	if dwAPL != "" {
		for i := range cases {
			if cases[i].label[0] == 'B' || cases[i].label[0] == 'C' {
				cases[i].apl = dwAPL
			}
		}
	}
	if only := os.Getenv("AB_CASES"); only != "" {
		var kept []abCase
		for _, c := range cases {
			if strings.Contains(only, c.label[:1]) {
				kept = append(kept, c)
			}
		}
		cases = kept
	}
	if v := os.Getenv("AB_TALENTS"); v != "" {
		for i := range cases {
			cases[i].talents = v
		}
	}
	for _, c := range cases {
		player := core.WithSpec(&proto.Player{
			Race:               abRace(),
			BonusStats:         abBonusStats(),
			Class:              proto.Class_ClassHunter,
			Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", c.gear).GearSet,
			TalentsString:      c.talents,
			Consumes:           abConsumes,
			Buffs:              abIndividualBuffs,
			Profession1:        proto.Profession_Enchanting,
			Profession2:        proto.Profession_Engineering,
			Rotation:           core.GetAplRotation("../../ui/hunter/apls", c.apl).Rotation,
			DistanceFromTarget: 3,
			ReactionTimeMs:     abReaction(),
			ChannelClipDelayMs: 50,
		}, abOptions)
		rsr := &proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, abDebuffs),
			Encounter:  abEncounter(),
			SimOptions: &proto.SimOptions{Iterations: iters, RandomSeed: 20261012},
		}
		res := core.RunRaidSim(rsr)
		if res.Error != nil {
			t.Fatalf("%s: sim error: %s", c.label, res.Error.Message)
		}
		abPrint(c.label, res.RaidMetrics.Parties[0].Players[0], iters)
		fmt.Printf("  raid DPS incl. pets %.1f\n", res.RaidMetrics.Dps.Avg)
	}
}

func abEncounter() *proto.Encounter {
	enc := core.MakeSingleTargetEncounter(0.15)
	enc.Duration = 120
	enc.DurationVariation = 15
	t := enc.Targets[0]
	t.Stats[26] = 3731   // armor, UI "Level 60" dummy
	t.Stats[17] = 805    // attack power
	t.Stats[34] = 127393 // health
	t.MinBaseDamage = 3000
	if os.Getenv("AB_MOBTYPE") == "none" {
		enc.Targets[0].MobType = proto.MobType_MobTypeUnknown
	}
	return enc
}

func abReaction() int32 {
	if v := os.Getenv("AB_REACT"); v != "" {
		var n int32
		fmt.Sscanf(v, "%d", &n)
		return n
	}
	return 200 // the UI default
}

func abRace() proto.Race {
	switch os.Getenv("AB_RACE") {
	case "troll":
		return proto.Race_RaceTroll
	case "dwarf":
		return proto.Race_RaceDwarf
	case "nightelf":
		return proto.Race_RaceNightElf
	case "tauren":
		return proto.Race_RaceTauren
	}
	return proto.Race_RaceOrc
}

func abBonusStats() *proto.UnitStats {
	st := make([]float64, 44)
	var hit, ap float64
	fmt.Sscanf(os.Getenv("AB_BONUS_HIT"), "%f", &hit)
	fmt.Sscanf(os.Getenv("AB_BONUS_AP"), "%f", &ap)
	st[18] = hit * core.MeleeHitRatingPerHitChance // melee hit
	st[13] = hit * core.SpellHitRatingPerHitChance // spell hit (ranged uses melee hit in this engine)
	st[17] = ap
	st[27] = ap // ranged attack power
	return &proto.UnitStats{Stats: st}
}
