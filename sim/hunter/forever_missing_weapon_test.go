package hunter

import (
	"math"
	"testing"
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

// rangedWeaponSpellIDs: every rank the sim casts of the spells the client ties to a bow, gun or crossbow
// (SpellEquippedItems)
var rangedWeaponSpellIDs = map[int32]string{
	3044: "Arcane Shot", 14281: "Arcane Shot", 14282: "Arcane Shot", 14283: "Arcane Shot", 14284: "Arcane Shot",
	14285: "Arcane Shot", 14286: "Arcane Shot", 14287: "Arcane Shot",
	19434: "Aimed Shot", 20900: "Aimed Shot", 20901: "Aimed Shot", 20902: "Aimed Shot", 20903: "Aimed Shot", 20904: "Aimed Shot",
	2643: "Multi-Shot",
	1978: "Serpent Sting", 13549: "Serpent Sting", 13550: "Serpent Sting", 13551: "Serpent Sting", 13552: "Serpent Sting",
	13553: "Serpent Sting", 13554: "Serpent Sting", 13555: "Serpent Sting", 25295: "Serpent Sting",
	1310687: "Sniper Shot", 1310785: "Sniper Shot", 1310786: "Sniper Shot",
	1510: "Volley", 14294: "Volley", 14295: "Volley",
}

// Each of those spells refuses to cast with no ranged weapon, also the ones no rotation here reaches.
func TestNoRangedWeaponBlocksEveryShot(t *testing.T) {
	// the patch test talents with Sniper Shot, the last Marksmanship talent; at 20 yd the shots' own range checks pass
	sim, hunter := newPatchTestHunter(t, patchTestSetup{talents: "-0053150010100001-500230230250222151", distance: 20})
	if hunter.HasRangedWeapon() {
		t.Fatal("the test hunter has a ranged weapon")
	}
	blocked := map[string]bool{}
	for _, spell := range hunter.Spellbook {
		name, ok := rangedWeaponSpellIDs[spell.ActionID.SpellID]
		if !ok {
			continue
		}
		if spell.ExtraCastCondition != nil && !spell.ExtraCastCondition(sim, hunter.CurrentTarget) {
			blocked[name] = true
		} else {
			t.Errorf("%s (%d) can be cast with no ranged weapon", name, spell.ActionID.SpellID)
		}
	}
	for _, name := range []string{"Arcane Shot", "Aimed Shot", "Multi-Shot", "Serpent Sting", "Sniper Shot", "Volley"} {
		if !blocked[name] {
			t.Errorf("%s is not blocked (or not registered)", name)
		}
	}
}

// A hunter with no ranged weapon has no Auto Shot, casts none of the shots or stings that need one, and fights in melee.
// An empty ranged slot gave Auto Shot a 0 s cycle, so it fired again at the same instant and the fight never ended; the
// run sits in a goroutine so a hang fails the test instead of the suite. The shots then dealt NaN damage (an empty
// weapon's damage over its 0 attack power per DPS).
func TestNoRangedWeaponFightsInMelee(t *testing.T) {
	noRanged := core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet
	noRanged.Items[proto.ItemSlot_ItemSlotRanged] = &proto.ItemSpec{}
	// the weave walks out to 11 yd before the pull; the basic one reads the bow's swing time
	weave := core.GetAplRotation("../../ui/hunter/apls", "arc_weave_2h").Rotation
	basic := core.GetAplRotation("../../ui/hunter/apls", "arc_weave_2h_basic").Rotation
	idle := &proto.APLRotation{}
	if err := protojson.Unmarshal([]byte(`{"type":"TypeAPL"}`), idle); err != nil {
		t.Fatalf("rotation: %v", err)
	}
	for _, c := range []struct {
		name     string
		gear     *proto.EquipmentSpec
		rotation *proto.APLRotation
		distance float64
	}{
		{"no ranged weapon, in melee", noRanged, idle, 3},
		{"no ranged weapon, weave", noRanged, weave, 3},
		{"no ranged weapon, bow swing time", noRanged, basic, 3},
		{"no ranged weapon, at range", noRanged, idle, 12},
		{"no gear, in melee", &proto.EquipmentSpec{}, idle, 3},
		{"no gear, weave", &proto.EquipmentSpec{}, weave, 3},
	} {
		t.Run(c.name, func(t *testing.T) {
			player := core.WithSpec(&proto.Player{
				Race:               proto.Race_RaceTroll,
				Class:              proto.Class_ClassHunter,
				Equipment:          c.gear,
				TalentsString:      patchTestTalents,
				Rotation:           c.rotation,
				DistanceFromTarget: c.distance,
			}, abOptions)
			req := &proto.RaidSimRequest{
				Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
				Encounter:  abEncounter(),
				SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 20261025},
			}
			done := make(chan *proto.RaidSimResult, 1)
			go func() { done <- core.RunRaidSim(req) }()
			var res *proto.RaidSimResult
			select {
			case res = <-done:
			case <-time.After(10 * time.Second):
				t.Fatal("one iteration did not finish in 10 s")
			}
			if res.Error != nil {
				t.Fatalf("sim error: %s", res.Error.Message)
			}
			um := res.RaidMetrics.Parties[0].Players[0]
			var autoShots, swings, shots int32
			for _, a := range um.Actions {
				for _, tg := range a.Targets {
					switch {
					case a.Id.GetOtherId() == proto.OtherAction_OtherActionShoot:
						autoShots += tg.Casts
					case a.Id.GetOtherId() == proto.OtherAction_OtherActionAttack && a.Id.Tag == 1:
						swings += tg.Casts
					case rangedWeaponSpellIDs[a.Id.GetSpellId()] != "":
						shots += tg.Casts
					}
				}
			}
			if autoShots != 0 {
				t.Errorf("%d Auto Shots with no ranged weapon", autoShots)
			}
			if shots != 0 {
				t.Errorf("%d shots and stings with no ranged weapon", shots)
			}
			if math.IsNaN(um.Dps.Avg) {
				t.Errorf("DPS is NaN")
			}
			if c.rotation == idle && c.distance <= core.MaxMeleeAttackDistance && swings == 0 {
				t.Errorf("no main-hand swings in melee")
			}
		})
	}
}

// meleeWeaponSpellIDs: every rank the sim casts of the strikes the client ties to a melee weapon in the main hand
// (SpellEquippedItems; Mongoose Bite's rank 4 has no row there, its ranks 1-3 do)
var meleeWeaponSpellIDs = map[int32]string{
	2973: "Raptor Strike", 14260: "Raptor Strike", 14261: "Raptor Strike", 14262: "Raptor Strike", 14263: "Raptor Strike",
	14264: "Raptor Strike", 14265: "Raptor Strike", 14266: "Raptor Strike",
	1495: "Mongoose Bite", 14269: "Mongoose Bite", 14270: "Mongoose Bite", 14271: "Mongoose Bite",
	2974: "Wing Clip", 14267: "Wing Clip", 14268: "Wing Clip",
	1317257: "Strider Kick",
}

// With no main-hand weapon the strikes refuse to cast (the bow still shoots, the fists still swing): each one the
// rotation can press, and in a fight Raptor Strike never queues.
func TestNoMeleeWeaponBlocksEveryStrike(t *testing.T) {
	noMainHand := map[proto.ItemSlot]int32{proto.ItemSlot_ItemSlotMainHand: 0}
	// the patch test talents carry Strider Kick
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, pieces: noMainHand, talents: patchTestTalents, distance: 3})
	if hunter.HasMHWeapon() || !hunter.HasRangedWeapon() {
		t.Fatal("the test hunter needs a bow and no main hand")
	}
	// Mongoose Bite's own condition wants Defensive State (or Expose Prey)
	hunter.DefensiveState.Activate(sim)
	blocked := map[string]bool{}
	for _, spell := range hunter.Spellbook {
		name, ok := meleeWeaponSpellIDs[spell.ActionID.SpellID]
		// what a rotation presses: Raptor Strike's queue, not the hit that takes the swing
		if !ok || !spell.Flags.Matches(core.SpellFlagAPL) {
			continue
		}
		if spell.ExtraCastCondition != nil && !spell.ExtraCastCondition(sim, hunter.CurrentTarget) {
			blocked[name] = true
		} else {
			t.Errorf("%s (%s) can be cast with no main-hand weapon", name, spell.ActionID)
		}
	}
	for _, name := range []string{"Raptor Strike", "Mongoose Bite", "Wing Clip", "Strider Kick"} {
		if !blocked[name] {
			t.Errorf("%s is not blocked (or not registered)", name)
		}
	}

	gear := core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet
	gear.Items[proto.ItemSlot_ItemSlotMainHand] = &proto.ItemSpec{}
	raptorOnly := &proto.APLRotation{}
	if err := protojson.Unmarshal([]byte(`{"type":"TypeAPL","priorityList":[{"action":{"castSpell":{"spellId":{"spellId":14266,"tag":1}}}}]}`), raptorOnly); err != nil {
		t.Fatalf("rotation: %v", err)
	}
	for _, c := range []struct {
		name     string
		rotation *proto.APLRotation
	}{
		{"Raptor Strike only", raptorOnly},
		{"weave", core.GetAplRotation("../../ui/hunter/apls", "arc_weave_2h").Rotation},
	} {
		player := core.WithSpec(&proto.Player{
			Race:               proto.Race_RaceTroll,
			Class:              proto.Class_ClassHunter,
			Equipment:          gear,
			TalentsString:      patchTestTalents,
			Rotation:           c.rotation,
			DistanceFromTarget: 3,
		}, abOptions)
		res := core.RunRaidSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  abEncounter(),
			SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 20261026},
		})
		if res.Error != nil {
			t.Fatalf("%s: sim error: %s", c.name, res.Error.Message)
		}
		var strikes, swings, shots int32
		for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
			for _, tg := range a.Targets {
				switch {
				case meleeWeaponSpellIDs[a.Id.GetSpellId()] != "":
					strikes += tg.Casts
				case a.Id.GetOtherId() == proto.OtherAction_OtherActionAttack && a.Id.Tag == 1:
					swings += tg.Casts
				case rangedWeaponSpellIDs[a.Id.GetSpellId()] != "":
					shots += tg.Casts
				}
			}
		}
		if strikes != 0 {
			t.Errorf("%s: %d strikes with no main-hand weapon", c.name, strikes)
		}
		if c.rotation == raptorOnly && swings == 0 {
			t.Errorf("%s: the fists never swung", c.name)
		}
		if c.rotation != raptorOnly && shots == 0 {
			t.Errorf("%s: no shots with the bow on", c.name)
		}
	}
}
