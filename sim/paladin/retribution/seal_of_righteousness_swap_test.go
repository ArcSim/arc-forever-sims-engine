package retribution

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
	"github.com/wowsims/forever/sim/paladin"
)

// Seal of Righteousness rank 8 and the damage spell its hits fire.
const (
	sealOfRighteousnessRank8 = 20293
	sealOfRighteousnessHit8  = 25713
	massOfMcGowan            = 13006 // a one-hand mace with no effects, 2.8 speed
	arcaniteReaper           = 12784 // a two-handed axe with no effects, 3.8 speed
)

// sealOfRighteousnessHit returns the average non-critical hit of Seal of Righteousness for a
// paladin wielding firstWeapon who, when swapTo is set, swaps to it as the fight starts. Plenty
// of spell power makes the coefficient the bigger part of the hit.
func sealOfRighteousnessHit(t *testing.T, firstWeapon, swapTo int32) float64 {
	t.Helper()
	weapons := func(mainHand int32) *proto.EquipmentSpec {
		items := make([]*proto.ItemSpec, 15)
		for i := range items {
			items[i] = &proto.ItemSpec{}
		}
		items[14] = &proto.ItemSpec{Id: mainHand}
		return &proto.EquipmentSpec{Items: items}
	}
	rotation := `{"type":"TypeAPL","priorityList":[
		{"action":{"castSpell":{"spellId":{"spellId":20293}},"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":20293}}}}}}}
	]}`
	if swapTo != 0 {
		rotation = `{"type":"TypeAPL","priorityList":[
			{"action":{"itemSwap":{"swapSet":"Swap1"}}},
			{"action":{"castSpell":{"spellId":{"spellId":20293}},"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":20293}}}}}}}
		]}`
	}
	var bonus stats.Stats
	bonus[stats.SpellDamage] = 1000
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceHuman,
		Class:         proto.Class_ClassPaladin,
		Equipment:     weapons(firstWeapon),
		Consumables:   DefaultConsumables,
		BonusStats:    &proto.UnitStats{Stats: bonus.ToProtoArray()},
		TalentsString: RetTalents,
		Rotation:      core.APLRotationFromJsonString(rotation),
	}, DefaultOptions)
	if swapTo != 0 {
		player.EnableItemSwap = true
		player.ItemSwap = &proto.ItemSwap{Items: weapons(swapTo).Items}
	}
	result := core.RunRaidSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 101, Iterations: 20},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	})
	if result.Error != nil {
		t.Fatal(result.Error.Message)
	}
	var damage float64
	var hits int32
	for _, action := range result.RaidMetrics.Parties[0].Players[0].Actions {
		if action.Id.GetSpellId() != sealOfRighteousnessHit8 {
			continue
		}
		for _, target := range action.Targets {
			damage += target.Damage - target.CritDamage
			hits += target.Hits - target.Crits
		}
	}
	if hits == 0 {
		t.Fatalf("no Seal of Righteousness hits with %d swapping to %d", firstWeapon, swapTo)
	}
	return damage / float64(hits)
}

// A two-hander swapped for a one-hander as the fight starts hits for what the one-hander hits for:
// the spell power share follows the weapon in the hand, as the weapon-speed part does.
func TestSealOfRighteousnessFollowsASwappedWeapon(t *testing.T) {
	oneHand := sealOfRighteousnessHit(t, massOfMcGowan, 0)
	swapped := sealOfRighteousnessHit(t, arcaniteReaper, massOfMcGowan)
	t.Logf("one-hander %.1f a hit, two-hander swapped to it %.1f", oneHand, swapped)
	if math.Abs(swapped/oneHand-1) > 0.02 {
		t.Errorf("Seal of Righteousness after a two-hander to one-hander swap hits for %.1f, the one-hander alone %.1f", swapped, oneHand)
	}

	paladin.SealOfRighteousnessFollowsSwappedWeapon = false
	defer func() { paladin.SealOfRighteousnessFollowsSwappedWeapon = true }()
	stale := sealOfRighteousnessHit(t, arcaniteReaper, massOfMcGowan)
	t.Logf("switched off: %.1f", stale)
	// 1000 spell power at 0.185 x 1.2 against 0.185 is about 37 more a hit.
	if stale-oneHand < 25 {
		t.Errorf("switched off, the swapped one-hander should keep the two-hander's coefficient: %.1f against %.1f", stale, oneHand)
	}
}
