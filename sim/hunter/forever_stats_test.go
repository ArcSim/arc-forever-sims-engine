package hunter

import (
	"fmt"
	"os"
	"testing"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
	"google.golang.org/protobuf/encoding/protojson"
)

// abWorldBuffs switches on every world buff the engine knows (AB_WORLDBUFFS=1).
func abWorldBuffs(player *proto.Player) {
	if player.Buffs == nil {
		player.Buffs = &proto.IndividualBuffs{}
	}
	b := player.Buffs
	b.RallyingCryOfTheDragonslayer = true
	b.SpiritOfZandalar = true
	b.SongflowerSerenade = true
	b.WarchiefsBlessing = true
	b.FengusFerocity = true
	b.MoldarsMoxie = true
	b.SlipkiksSavvy = true
	b.SaygesFortune = proto.SaygesFortune_SaygesDamage
}

// TestForeverStats prints the final character stats of a UI settings export (AB_SETTINGS), with
// AB_WORLDBUFFS=1 adding every world buff, so buff claims can be checked against the engine.
func TestForeverStats(t *testing.T) {
	path := os.Getenv("AB_SETTINGS")
	if path == "" {
		t.Skip("AB_SETTINGS not set")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	settings := &proto.IndividualSimSettings{}
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(data, settings); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	player := settings.Player
	label := "no world buffs"
	if os.Getenv("AB_WORLDBUFFS") != "" {
		abWorldBuffs(player)
		label = "all world buffs"
	}
	res := core.ComputeStats(&proto.ComputeStatsRequest{
		Raid:      core.SinglePlayerRaidProto(player, settings.PartyBuffs, settings.RaidBuffs, settings.Debuffs),
		Encounter: settings.Encounter,
	})
	if res.ErrorResult != "" {
		t.Fatalf("compute stats: %s", res.ErrorResult)
	}
	ps := res.RaidStats.Parties[0].Players[0]
	show := func(name string, us *proto.UnitStats) {
		if us == nil {
			return
		}
		s := us.Stats
		fmt.Printf("  %-6s str %6.0f  agi %6.0f  AP %6.0f  RAP %6.0f  melee crit %5.2f%%  melee hit %5.2f%%\n",
			name, s[stats.Strength], s[stats.Agility], s[stats.AttackPower], s[stats.RangedAttackPower], s[stats.MeleeCrit], s[stats.MeleeHit])
	}
	fmt.Printf("\n### stats, %s\n", label)
	show("gear", ps.GearStats)
	show("buffs", ps.BuffsStats)
	show("final", ps.FinalStats)
}
