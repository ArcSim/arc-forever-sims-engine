package warlock

import (
	"regexp"
	"strings"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/buffs"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// The client keeps the raid-wide buffs, the Greater Blessings and Devotion Aura off pets
// (ATTR_EX_3_ONLY_TARGET_PLAYERS); the party auras without the flag reach them.
func TestPartyAuraPetFlags(t *testing.T) {
	for id, playersOnly := range map[int32]bool{
		25918: true,  // Greater Blessing of Wisdom
		23028: true,  // Arcane Brilliance
		27681: true,  // Prayer of Spirit
		21850: true,  // Gift of the Wild
		10293: true,  // Devotion Aura
		25289: true,  // Battle Shout
		10494: false, // Mana Spring Totem's aura
		10301: false, // Retribution Aura
		20905: false, // Trueshot Aura
	} {
		if got := spelldata.MustFind(id).OnlyTargetsPlayers(); got != playersOnly {
			t.Errorf("spell %d: OnlyTargetsPlayers %v, want %v", id, got, playersOnly)
		}
	}
}

// A summoned Imp in a raid with every mana buff gets only Mana Spring Totem: 10 mana every 2 sec,
// which it keeps while it casts. Before, a pet got no party aura at all, so the Imp had no
// regeneration while casting.
func TestImpGetsManaSpringTotemOnly(t *testing.T) {
	run := func() string {
		player := core.WithSpec(&proto.Player{
			Race: proto.Race_RaceOrc, Class: proto.Class_ClassWarlock,
			Equipment: &proto.EquipmentSpec{}, Consumables: &proto.ConsumesSpec{},
			Buffs:    &proto.IndividualBuffs{GreaterBlessingOfWisdom: true, GreaterBlessingOfKings: true},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.Player_Warlock{Warlock: &proto.Warlock{Options: &proto.Warlock_Options{ClassOptions: &proto.WarlockOptions{
			Summon: proto.WarlockOptions_Imp,
		}}}})
		party := &proto.PartyBuffs{ManaSpringTotem: proto.TristateEffect_TristateEffectRegular}
		raid := &proto.RaidBuffs{ArcaneBrilliance: true, PrayerOfSpirit: true, GiftOfTheWild: true}
		res := core.RunRaidSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, party, raid, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
			SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 1, Debug: true},
		})
		if res.Error != nil {
			t.Fatal(res.Error.Message)
		}
		return res.Logs
	}
	petStats := regexp.MustCompile(`- Imp\] Pet stats: (\{[^}]*\})`)
	castingRegen := regexp.MustCompile(`- Imp\] Gained 10\.000 mana from \{OtherID: 2, Tag: 1\}`)

	logs := run()
	stats := petStats.FindStringSubmatch(logs)
	if stats == nil {
		t.Fatal("no Imp stats logged")
	}
	for _, want := range []string{`"Intellect": 264.000`, `"Spirit": 260.000`, `"Stamina": 86.000`, `"MP5": 25.000`} {
		if !strings.Contains(stats[1], want) {
			t.Errorf("Imp stats %s, want %s (Mana Spring Totem only)", stats[1], want)
		}
	}
	if n := len(castingRegen.FindAllString(logs, -1)); n < 10 {
		t.Errorf("%d Mana Spring ticks while casting, want the Imp to regenerate while it casts", n)
	}

	buffs.PartyAurasReachPets = false
	defer func() { buffs.PartyAurasReachPets = true }()
	logs = run()
	if stats := petStats.FindStringSubmatch(logs); stats == nil || strings.Contains(stats[1], `"MP5"`) {
		t.Errorf("switched off, Imp stats %v carry MP5", stats)
	}
	if n := len(castingRegen.FindAllString(logs, -1)); n != 0 {
		t.Errorf("switched off, %d Mana Spring ticks while casting", n)
	}
}
