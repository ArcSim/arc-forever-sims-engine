package enhancement

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// Patch 141: Windfury Weapon's main-hand strikes (439440) proc Flametongue Totem as main-hand autos do
// (beta testing relayed by Redfall, MythicSim Discord, 16 October 2026: "Windfury extra attacks appear to be
// proccing FTT"), and they still never glance ("WFW, despite being changed to a white attack, does not
// appear to glance").

type wfFtCounts struct {
	mhAutos              int32 // landed main-hand swings
	autoGlances          int32
	strikes              int32 // every Windfury Weapon strike, either hand, whatever its outcome
	mhStrikesLanded      int32
	strikeGlances        int32
	flametongueTotemHits int32 // every roll of the totem's hit
}

func runWindfuryWithFlametongueTotem(t *testing.T, mh, oh int32, imbueMh, imbueOh proto.ShamanImbue, party *proto.PartyBuffs) wfFtCounts {
	t.Helper()
	items := make([]*proto.ItemSpec, 16)
	for i := range items {
		items[i] = &proto.ItemSpec{}
	}
	items[14] = &proto.ItemSpec{Id: mh}
	if oh != 0 {
		items[15] = &proto.ItemSpec{Id: oh}
	}
	player := &proto.Player{
		Name: "enh", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc, TalentsString: DefaultTalents,
		Equipment: &proto.EquipmentSpec{Items: items},
		Spec: &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{
			ClassOptions: &proto.ShamanOptions{ImbueMh: imbueMh},
			ImbueOh:      imbueOh,
		}}},
		Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}
	raid := &proto.Raid{Parties: []*proto.Party{{Players: []*proto.Player{player}, Buffs: party}}, Buffs: &proto.RaidBuffs{}, Debuffs: &proto.Debuffs{}, NumActiveParties: 1}
	res := core.RunRaidSim(&proto.RaidSimRequest{Raid: raid, Encounter: core.MakeSingleTargetEncounter(0), SimOptions: &proto.SimOptions{Iterations: 20, RandomSeed: 1}})
	if res.Error != nil {
		t.Fatal(res.Error.Message)
	}
	var c wfFtCounts
	for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
		for _, tgt := range a.Targets {
			switch {
			case a.Id.GetOtherId() == proto.OtherAction_OtherActionAttack && a.Id.Tag == 1:
				c.mhAutos += tgt.Hits + tgt.Crits + tgt.Glances + tgt.Blocks
				c.autoGlances += tgt.Glances
			case a.Id.GetSpellId() == 439440:
				c.strikes += tgt.Hits + tgt.Crits + tgt.Misses + tgt.Dodges + tgt.Parries + tgt.Blocks + tgt.Glances
				c.mhStrikesLanded += tgt.Hits + tgt.Crits + tgt.Blocks + tgt.Glances
				c.strikeGlances += tgt.Glances
			case a.Id.GetSpellId() == 439441:
				c.strikes += tgt.Hits + tgt.Crits + tgt.Misses + tgt.Dodges + tgt.Parries + tgt.Blocks + tgt.Glances
				c.strikeGlances += tgt.Glances
			case a.Id.GetSpellId() == flametongueTotemHit && a.Id.Tag == 0:
				c.flametongueTotemHits += tgt.Hits + tgt.Crits + tgt.Misses
			}
		}
	}
	return c
}

// One totem hit for every landed main-hand auto and every landed main-hand Windfury Weapon strike.
func TestWindfuryWeaponStrikesProcFlametongueTotem(t *testing.T) {
	c := runWindfuryWithFlametongueTotem(t, 17182, 0, proto.ShamanImbue_WindfuryWeapon, proto.ShamanImbue_NoImbue, &proto.PartyBuffs{FlametongueTotem: true})
	if c.mhStrikesLanded == 0 || c.mhAutos == 0 {
		t.Fatalf("%d landed autos and %d landed Windfury strikes; the test needs both", c.mhAutos, c.mhStrikesLanded)
	}
	if want := c.mhAutos + c.mhStrikesLanded; c.flametongueTotemHits != want {
		t.Errorf("%d totem hits for %d landed autos and %d landed Windfury strikes, want %d", c.flametongueTotemHits, c.mhAutos, c.mhStrikesLanded, want)
	}
}

// The totem is main hand only: an off-hand Windfury Weapon's strikes (439441) add nothing.
func TestOffHandWindfuryWeaponDoesNotProcFlametongueTotem(t *testing.T) {
	c := runWindfuryWithFlametongueTotem(t, crestedScepter, stonevaultShiv, proto.ShamanImbue_NoImbue, proto.ShamanImbue_WindfuryWeapon, &proto.PartyBuffs{FlametongueTotem: true})
	if c.strikes == 0 {
		t.Fatal("the off hand's Windfury Weapon never struck")
	}
	if c.flametongueTotemHits != c.mhAutos {
		t.Errorf("%d totem hits for %d landed main-hand autos, want one each and none for the off hand's strikes", c.flametongueTotemHits, c.mhAutos)
	}
}

// Windfury Weapon's strikes roll the special attack table: the boss glances the autos and never a strike.
func TestWindfuryWeaponStrikesDoNotGlance(t *testing.T) {
	c := runWindfuryWithFlametongueTotem(t, 17182, 0, proto.ShamanImbue_WindfuryWeapon, proto.ShamanImbue_NoImbue, &proto.PartyBuffs{FlametongueTotem: true})
	if c.autoGlances == 0 {
		t.Fatal("no auto glanced; the target does not roll glancing blows")
	}
	if c.strikes == 0 || c.strikeGlances != 0 {
		t.Errorf("%d of %d Windfury strikes glanced, want none", c.strikeGlances, c.strikes)
	}
}
