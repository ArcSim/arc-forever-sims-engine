package enhancement

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
	"github.com/wowsims/forever/sim/shaman"
)

// Mythicsim patch 133: the imbues, stones and oils follow the weapon in the hand at an item swap.

const (
	elementalSharpeningStone = 22756
	frostOil                 = 3829
	denseSharpeningStone     = 16138
	frostOilAura             = "Frost Oil (main hand)"
	flametongueMainHandAura  = "Flametongue Imbue ItemSlotMainHand"
)

func weaponOnly(mainHand int32) []*proto.ItemSpec {
	items := make([]*proto.ItemSpec, 15)
	for i := range items {
		items[i] = &proto.ItemSpec{}
	}
	items[14] = &proto.ItemSpec{Id: mainHand}
	return items
}

func imbueSwapPlayer(mh, mhSwap proto.ShamanImbue, consumes *proto.ConsumesSpec, rotation string) *proto.Player {
	if rotation == "" {
		rotation = `{"type":"TypeAPL"}`
	}
	return &proto.Player{
		Name: "enh", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc,
		Equipment:      &proto.EquipmentSpec{Items: weaponOnly(crestedScepter)},
		EnableItemSwap: true,
		ItemSwap:       &proto.ItemSwap{Items: weaponOnly(arcaniteReaper)},
		Consumables:    consumes,
		Spec: &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{
			ClassOptions: &proto.ShamanOptions{ImbueMh: mh, ImbueMhSwap: mhSwap},
		}}},
		Rotation: core.APLRotationFromJsonString(rotation),
	}
}

// imbueSwapStates returns the character's state, read by measure, with the weapon worn at the pull
// (Crested Scepter) in the hand and after a swap to Arcanite Reaper.
func imbueSwapStates[T any](mh, mhSwap proto.ShamanImbue, consumes *proto.ConsumesSpec, measure func(*core.Character) T) (before, after T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid:       core.SinglePlayerRaidProto(imbueSwapPlayer(mh, mhSwap, consumes, ""), &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	character := sim.Raid.Parties[0].Players[0].GetCharacter()
	before = measure(character)
	character.ItemSwap.SwapItems(sim, proto.APLActionItemSwap_Swap1, false)
	after = measure(character)
	return before, after
}

func attackPower(c *core.Character) float64 { return c.GetStat(stats.AttackPower) }

func TestRockbiterFollowsTheWeaponInTheHand(t *testing.T) {
	const none, rockbiter = proto.ShamanImbue_NoImbue, proto.ShamanImbue_RockbiterWeapon
	baseBefore, baseAfter := imbueSwapStates(none, none, nil, attackPower)
	for _, c := range []struct {
		name          string
		mh, mhSwap    proto.ShamanImbue
		before, after float64
	}{
		{"swap's weapon only", none, rockbiter, 0, 653},
		{"pull's weapon only", rockbiter, none, 653, 0},
		{"both weapons", rockbiter, rockbiter, 653, 653},
	} {
		before, after := imbueSwapStates(c.mh, c.mhSwap, nil, attackPower)
		if got := before - baseBefore; !core.WithinToleranceFloat64(c.before, got, 0.01) {
			t.Errorf("%s: %.2f attack power before the swap, want %v", c.name, got, c.before)
		}
		if got := after - baseAfter; !core.WithinToleranceFloat64(c.after, got, 0.01) {
			t.Errorf("%s: %.2f attack power after the swap, want %v", c.name, got, c.after)
		}
	}

	shaman.PerWeaponImbues = false
	defer func() { shaman.PerWeaponImbues = true }()
	for _, c := range []struct {
		name          string
		mh, mhSwap    proto.ShamanImbue
		before, after float64
	}{
		{"switched off, swap's weapon only", none, rockbiter, 0, 0},
		{"switched off, pull's weapon only", rockbiter, none, 653, 653},
	} {
		before, after := imbueSwapStates(c.mh, c.mhSwap, nil, attackPower)
		if got := before - baseBefore; !core.WithinToleranceFloat64(c.before, got, 0.01) {
			t.Errorf("%s: %.2f attack power before the swap, want %v", c.name, got, c.before)
		}
		if got := after - baseAfter; !core.WithinToleranceFloat64(c.after, got, 0.01) {
			t.Errorf("%s: %.2f attack power after the swap, want %v", c.name, got, c.after)
		}
	}
}

// Spirit Weapons' threat bonus under Rockbiter comes and goes with it.
func TestRockbiterThreatFollowsTheWeaponInTheHand(t *testing.T) {
	threat := func(mh, mhSwap proto.ShamanImbue) (float64, float64) {
		sim := core.NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 1},
			Raid: core.SinglePlayerRaidProto(func() *proto.Player {
				p := imbueSwapPlayer(mh, mhSwap, nil, "")
				p.TalentsString = DefaultTalents
				return p
			}(), &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter: core.MakeSingleTargetEncounter(0),
		}, simsignals.CreateSignals())
		sim.Reset()
		character := sim.Raid.Parties[0].Players[0].GetCharacter()
		before := character.PseudoStats.ThreatMultiplier
		character.ItemSwap.SwapItems(sim, proto.APLActionItemSwap_Swap1, false)
		return before, character.PseudoStats.ThreatMultiplier
	}
	plain, _ := threat(proto.ShamanImbue_NoImbue, proto.ShamanImbue_NoImbue)
	before, after := threat(proto.ShamanImbue_NoImbue, proto.ShamanImbue_RockbiterWeapon)
	if !core.WithinToleranceFloat64(plain, before, 1e-9) || !core.WithinToleranceFloat64(plain*1.86, after, 1e-3) {
		t.Errorf("Rockbiter on the swap's weapon: threat x%.3f before and x%.3f after the swap, want x%.3f and x%.3f", before, after, plain, plain*1.86)
	}
}

func TestFlametongueOnTheSwapWeaponProcsOnlyAfterTheSwap(t *testing.T) {
	const none, flametongue = proto.ShamanImbue_NoImbue, proto.ShamanImbue_FlametongueWeapon
	active := func(c *core.Character) bool {
		aura := c.GetAura(flametongueMainHandAura)
		return aura != nil && aura.IsActive()
	}
	for _, c := range []struct {
		name          string
		mh, mhSwap    proto.ShamanImbue
		before, after bool
	}{
		{"swap's weapon only", none, flametongue, false, true},
		{"pull's weapon only", flametongue, none, true, false},
		{"both weapons", flametongue, flametongue, true, true},
	} {
		if before, after := imbueSwapStates(c.mh, c.mhSwap, nil, active); before != c.before || after != c.after {
			t.Errorf("%s: proc aura active %v before and %v after the swap, want %v and %v", c.name, before, after, c.before, c.after)
		}
	}

	// A fight: with the swap's weapon imbued and no swap, no Flametongue hit; swapping as the fight
	// starts, Flametongue hits.
	hits := func(rotation string) int32 {
		result := core.RunRaidSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 7, Iterations: 5},
			Raid: core.SinglePlayerRaidProto(imbueSwapPlayer(none, flametongue, nil, rotation),
				&proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter: core.MakeSingleTargetEncounter(0),
		})
		if result.Error != nil {
			t.Fatal(result.Error.Message)
		}
		var n int32
		for _, action := range result.RaidMetrics.Parties[0].Players[0].Actions {
			if action.Id.GetSpellId() == flametongueImbueHit {
				for _, target := range action.Targets {
					n += target.Hits
				}
			}
		}
		return n
	}
	if n := hits(""); n != 0 {
		t.Errorf("Flametongue on the swap's weapon only, never swapped: %d hits, want none", n)
	}
	if n := hits(`{"type":"TypeAPL","priorityList":[{"action":{"itemSwap":{"swapSet":"Swap1"}}}]}`); n == 0 {
		t.Errorf("Flametongue on the swap's weapon only, swapped at the pull: no hits")
	}
}

func TestStoneAndOilFollowTheWeaponInTheHand(t *testing.T) {
	const none = proto.ShamanImbue_NoImbue
	crit := func(c *core.Character) float64 { return c.GetStat(stats.PhysicalCritPercent) }
	baseBefore, baseAfter := imbueSwapStates(none, none, nil, crit)

	// Frost Oil on the pull's weapon, an Elemental Sharpening Stone on the swap's weapon.
	split := &proto.ConsumesSpec{MhImbueId: frostOil, MhImbueSwapId: elementalSharpeningStone}
	before, after := imbueSwapStates(none, none, split, crit)
	if !core.WithinToleranceFloat64(baseBefore, before, 1e-6) || !core.WithinToleranceFloat64(baseAfter+2, after, 1e-6) {
		t.Errorf("stone on the swap's weapon: crit %+.2f before and %+.2f after the swap, want +0 and +2", before-baseBefore, after-baseAfter)
	}
	oil := func(c *core.Character) bool {
		aura := c.GetAura(frostOilAura)
		return aura != nil && aura.IsActive()
	}
	if before, after := imbueSwapStates(none, none, split, oil); !before || after {
		t.Errorf("Frost Oil on the pull's weapon: active %v before and %v after the swap, want true and false", before, after)
	}
	// -1: nothing on the swap's weapon.
	if before, after := imbueSwapStates(none, none, &proto.ConsumesSpec{MhImbueId: frostOil, MhImbueSwapId: -1}, oil); !before || after {
		t.Errorf("Frost Oil on the pull's weapon, nothing on the swap's: active %v before and %v after, want true and false", before, after)
	}
	// The same oil on both weapons stays on.
	if before, after := imbueSwapStates(none, none, &proto.ConsumesSpec{MhImbueId: frostOil, MhImbueSwapId: frostOil}, oil); !before || !after {
		t.Errorf("Frost Oil on both weapons: active %v before and %v after, want true and true", before, after)
	}
	// Unset, the main-hand stone counts all fight, as before.
	if before, after := imbueSwapStates(none, none, &proto.ConsumesSpec{MhImbueId: elementalSharpeningStone}, crit); !core.WithinToleranceFloat64(baseBefore+2, before, 1e-6) || !core.WithinToleranceFloat64(baseAfter+2, after, 1e-6) {
		t.Errorf("stone without mhImbueSwapId: crit %+.2f before and %+.2f after the swap, want +2 and +2", before-baseBefore, after-baseAfter)
	}

	// A Dense stone on the swap's weapon is on its swing only.
	damage := func(c *core.Character) float64 { return c.AutoAttacks.MH().BaseDamageMin }
	plainBefore, plainAfter := imbueSwapStates(none, none, nil, damage)
	before, after = imbueSwapStates(none, none, &proto.ConsumesSpec{MhImbueSwapId: denseSharpeningStone}, damage)
	if before != plainBefore || after != plainAfter+8 {
		t.Errorf("Dense stone on the swap's weapon: %+.1f damage before and %+.1f after the swap, want +0 and +8", before-plainBefore, after-plainAfter)
	}
	before, after = imbueSwapStates(none, none, &proto.ConsumesSpec{MhImbueId: denseSharpeningStone, MhImbueSwapId: -1}, damage)
	if before != plainBefore+8 || after != plainAfter {
		t.Errorf("Dense stone on the pull's weapon: %+.1f damage before and %+.1f after the swap, want +8 and +0", before-plainBefore, after-plainAfter)
	}
}
