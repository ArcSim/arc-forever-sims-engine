package dps

import (
	"fmt"
	"strings"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// Diamond Flask's Strength comes from finishing its 5 sec channel, so cast at the pull it holds
// for a full minute from the 5 sec mark.
func TestDiamondFlask(t *testing.T) {
	equipment := weaponsOnly(15240, 15238)
	equipment.Items[proto.ItemSlot_ItemSlotTrinket1] = &proto.ItemSpec{Id: 20130}
	rotation := core.GetAplRotation("../../../ui/specs/warrior/dps/apls", "dps_reck").Rotation
	rotation.PriorityList = append([]*proto.APLListItem{{Action: &proto.APLAction{Action: &proto.APLAction_CastSpell{
		CastSpell: &proto.APLActionCastSpell{SpellId: &proto.ActionID{RawId: &proto.ActionID_ItemId{ItemId: 20130}}},
	}}}}, rotation.PriorityList...)

	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceOrc,
		Class:         proto.Class_ClassWarrior,
		Equipment:     equipment,
		Consumables:   &proto.ConsumesSpec{},
		TalentsString: DpsTalents,
		Rotation:      rotation,
	}, DefaultOptions)
	raid := core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{})
	result := core.RunRaidSim(&proto.RaidSimRequest{Raid: raid, Encounter: core.MakeSingleTargetEncounter(0),
		SimOptions: &proto.SimOptions{Iterations: 5, RandomSeed: 101}})
	if result.Error != nil {
		t.Fatal(result.Error.Message)
	}
	for _, aura := range result.RaidMetrics.Parties[0].Players[0].Auras {
		if aura.Id.GetSpellId() == 1318070 {
			if aura.UptimeSecondsAvg != 60 {
				t.Errorf("Diamond Flask's Strength up %v sec, want 60", aura.UptimeSecondsAvg)
			}
			return
		}
	}
	t.Error("Diamond Flask never granted its Strength")
}

// Shadowstrike, The Cruel Hand of Timmy and Skullforge Reaver drain life: a health leech (the first
// two) or a damage aura (Skullforge Brand). Each proc has to land its damage on the target.
func TestLifeDrainWeaponProcs(t *testing.T) {
	for _, weapon := range []struct {
		itemID, spellID int32
	}{{17074, 21170}, {13401, 17505}, {13361, 17484}} {
		player := core.WithSpec(&proto.Player{
			Race:          proto.Race_RaceOrc,
			Class:         proto.Class_ClassWarrior,
			Equipment:     weaponsOnly(weapon.itemID, 0),
			Consumables:   &proto.ConsumesSpec{},
			TalentsString: DpsTalents,
			Rotation:      core.GetAplRotation("../../../ui/specs/warrior/dps/apls", "dps_reck").Rotation,
		}, DefaultOptions)
		raid := core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{})
		result := core.RunRaidSim(&proto.RaidSimRequest{Raid: raid, Encounter: core.MakeSingleTargetEncounter(0),
			SimOptions: &proto.SimOptions{Iterations: 5, RandomSeed: 101}})
		if result.Error != nil {
			t.Fatal(result.Error.Message)
		}
		damage := 0.0
		for _, action := range result.RaidMetrics.Parties[0].Players[0].Actions {
			if action.Id.GetSpellId() == weapon.spellID {
				for _, target := range action.Targets {
					damage += target.Damage
				}
			}
		}
		if damage <= 0 {
			t.Errorf("item %d: proc %d dealt no damage", weapon.itemID, weapon.spellID)
		}
	}
}

// Sword of Zeal and Argent Avenger put their client row's buff on the wearer when they proc.
func TestWeaponProcBuffsFromTheirRows(t *testing.T) {
	for _, weapon := range []struct {
		itemID, spellID int32
	}{{6622, 8191}, {13246, 17352}} {
		player := core.WithSpec(&proto.Player{
			Race:          proto.Race_RaceOrc,
			Class:         proto.Class_ClassWarrior,
			Equipment:     weaponsOnly(weapon.itemID, 0),
			Consumables:   &proto.ConsumesSpec{},
			TalentsString: DpsTalents,
			Rotation:      core.GetAplRotation("../../../ui/specs/warrior/dps/apls", "dps_reck").Rotation,
		}, DefaultOptions)
		raid := core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{})
		result := core.RunRaidSim(&proto.RaidSimRequest{Raid: raid, Encounter: core.MakeSingleTargetEncounter(0),
			SimOptions: &proto.SimOptions{Iterations: 5, RandomSeed: 101}})
		if result.Error != nil {
			t.Fatal(result.Error.Message)
		}
		uptime := 0.0
		for _, aura := range result.RaidMetrics.Parties[0].Players[0].Auras {
			if aura.Id.GetSpellId() == weapon.spellID {
				uptime = aura.UptimeSecondsAvg
			}
		}
		if uptime <= 0 {
			t.Errorf("item %d never applied %d", weapon.itemID, weapon.spellID)
		}
	}
}

// Fiery Weapon and Lifestealing land their client row's damage on the target when they proc.
func TestWeaponEnchantDamageProcs(t *testing.T) {
	for _, enchant := range []struct {
		effectID, spellID int32
	}{{803, 13897}, {1898, 20004}} {
		equipment := weaponsOnly(15240, 0)
		equipment.Items[proto.ItemSlot_ItemSlotMainHand].Enchant = enchant.effectID
		player := core.WithSpec(&proto.Player{
			Race:          proto.Race_RaceOrc,
			Class:         proto.Class_ClassWarrior,
			Equipment:     equipment,
			Consumables:   &proto.ConsumesSpec{},
			TalentsString: DpsTalents,
			Rotation:      core.GetAplRotation("../../../ui/specs/warrior/dps/apls", "dps_reck").Rotation,
		}, DefaultOptions)
		raid := core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{})
		result := core.RunRaidSim(&proto.RaidSimRequest{Raid: raid, Encounter: core.MakeSingleTargetEncounter(0),
			SimOptions: &proto.SimOptions{Iterations: 5, RandomSeed: 101}})
		if result.Error != nil {
			t.Fatal(result.Error.Message)
		}
		damage := 0.0
		for _, action := range result.RaidMetrics.Parties[0].Players[0].Actions {
			if action.Id.GetSpellId() == enchant.spellID {
				for _, target := range action.Targets {
					damage += target.Damage
				}
			}
		}
		if damage <= 0 {
			t.Errorf("enchant %d: proc %d dealt no damage", enchant.effectID, enchant.spellID)
		}
	}
}

// Flurry Axe's proc is an extra swing (18797): passive procs keep no cast count, so read the log.
// Electrified Dagger's is a 45 Nature bolt (23592) that has to deal damage.
func TestFlurryAxeAndElectrifiedDagger(t *testing.T) {
	for _, weapon := range []struct {
		itemID, spellID int32
	}{{871, 18797}, {19100, 23592}} {
		player := core.WithSpec(&proto.Player{
			Race:          proto.Race_RaceOrc,
			Class:         proto.Class_ClassWarrior,
			Equipment:     weaponsOnly(weapon.itemID, 0),
			Consumables:   &proto.ConsumesSpec{},
			TalentsString: DpsTalents,
			Rotation:      core.GetAplRotation("../../../ui/specs/warrior/dps/apls", "dps_reck").Rotation,
		}, DefaultOptions)
		raid := core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{})
		result := core.RunRaidSim(&proto.RaidSimRequest{Raid: raid, Encounter: core.MakeSingleTargetEncounter(0),
			SimOptions: &proto.SimOptions{Iterations: 5, RandomSeed: 101, Debug: true}})
		if result.Error != nil {
			t.Fatal(result.Error.Message)
		}
		damage := 0.0
		for _, action := range result.RaidMetrics.Parties[0].Players[0].Actions {
			if action.Id.GetSpellId() == weapon.spellID {
				for _, target := range action.Targets {
					damage += target.Damage
				}
			}
		}
		casts := strings.Count(result.Logs, fmt.Sprintf("Casting {SpellID: %d}", weapon.spellID))
		if casts == 0 || (weapon.spellID == 23592 && damage <= 0) {
			t.Errorf("item %d: proc %d cast %d times for %.0f damage", weapon.itemID, weapon.spellID, casts, damage)
		}
	}
}
