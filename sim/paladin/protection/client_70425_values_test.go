package protection

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/paladin"
)

// A two-hander alone, so One-Handed Weapon Specialization is off and the rest of the build is the same.
var twoHanderOnly = &proto.EquipmentSpec{
	Items: []*proto.ItemSpec{
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{Id: 12784}, // Arcanite Reaper
	},
}

func newProtPaladinWith(t *testing.T, equipment *proto.EquipmentSpec) *paladin.Paladin {
	t.Helper()
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: core.SinglePlayerRaidProto(&proto.Player{
			Name: "prot", Class: proto.Class_ClassPaladin, Race: proto.Race_RaceHuman, TalentsString: ProtDeepTalents,
			Equipment: equipment,
			Spec:      DefaultOptions,
			Rotation:  &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim.Raid.Parties[0].Players[0].(paladin.PaladinAgent).GetPaladin()
}

// The highest rank of the spells a class mask names, with the tag the caller wants.
func spellByMask(t *testing.T, pal *paladin.Paladin, mask int64, tag int32) *core.Spell {
	t.Helper()
	var found *core.Spell
	for _, spell := range pal.Spellbook {
		if spell.Matches(mask) && spell.ActionID.Tag == tag && (found == nil || spell.ActionID.SpellID > found.ActionID.SpellID) {
			found = spell
		}
	}
	if found == nil {
		t.Fatalf("no spell for mask %d tag %d", mask, tag)
	}
	return found
}

// Spell power shares measured in game on client 70425 (Ashreon, MythicSim Discord, 16 October 2026),
// with the row switches off (patch 168 turns the Judgement of Righteousness and Seal of Fury ones on;
// TestPaladinSpellPowerFromRows covers both settings).
func TestPaladinSpellPowerCoefficientsFromClient70425(t *testing.T) {
	oldJudgement, oldFury := paladin.JudgementOfRighteousnessSpellPowerFromRow, paladin.SealOfFurySpellPowerFromRows
	paladin.JudgementOfRighteousnessSpellPowerFromRow, paladin.SealOfFurySpellPowerFromRows = false, false
	defer func() {
		paladin.JudgementOfRighteousnessSpellPowerFromRow, paladin.SealOfFurySpellPowerFromRows = oldJudgement, oldFury
	}()
	pal := newProtPaladinWith(t, WeaponAndShield)
	for _, c := range []struct {
		name string
		mask int64
		tag  int32
		want float64
	}{
		{"Holy Strike", paladin.SpellMaskHolyStrike, 0, 0.286},
		{"Judgement of Righteousness", paladin.SpellMaskJudgementOfRighteousness, 0, 0.70},
		{"Seal of Fury, each hit", paladin.SpellMaskSealOfFuryProc, 0, 0.15},
		{"Seal of Fury, judged", paladin.SpellMaskJudgementOfFury, 0, 0.67},
		{"Holy Shield", paladin.SpellMaskHolyShieldProc, 2, 0.07},
	} {
		if got := spellByMask(t, pal, c.mask, c.tag).BonusCoefficient; math.Abs(got-c.want) > 1e-9 {
			t.Errorf("%s spell power coefficient %v, want %v", c.name, got, c.want)
		}
	}
	// Consecration's 0.042 a tick, 0.336 over its 8 ticks, already matched the measured 33%.
}

// One-Handed Weapon Specialization (3/3, 10%) raises every school's damage, Holy included, but not
// Consecration's, Retribution Aura's or any heal.
func TestOneHandedWeaponSpecializationCoversEverySchoolButConsecrationAndRetributionAura(t *testing.T) {
	oneHand := newProtPaladinWith(t, WeaponAndShield)
	twoHand := newProtPaladinWith(t, twoHanderOnly)
	if oneHand.Talents.OneHandedWeaponSpecialization != 3 {
		t.Fatalf("test build has One-Handed Weapon Specialization %d, want 3", oneHand.Talents.OneHandedWeaponSpecialization)
	}
	ratio := func(pick func(*paladin.Paladin) *core.Spell) float64 {
		return pick(oneHand).DamageMultiplier / pick(twoHand).DamageMultiplier
	}
	byMask := func(mask int64, tag int32) func(*paladin.Paladin) *core.Spell {
		return func(pal *paladin.Paladin) *core.Spell { return spellByMask(t, pal, mask, tag) }
	}
	retributionAura := func(pal *paladin.Paladin) *core.Spell {
		for _, spell := range pal.Spellbook {
			if spell.ActionID.SpellID == paladin.RetributionAuraRankMap.Highest().ID && spell.ActionID.Tag != 0 {
				return spell
			}
		}
		t.Fatal("no Retribution Aura damage spell")
		return nil
	}
	for _, c := range []struct {
		name string
		pick func(*paladin.Paladin) *core.Spell
		want float64
	}{
		{"melee auto attack", func(pal *paladin.Paladin) *core.Spell { return pal.AutoAttacks.MHAuto() }, 1.1},
		{"Holy Strike", byMask(paladin.SpellMaskHolyStrike, 0), 1.1},
		{"Judgement of Righteousness", byMask(paladin.SpellMaskJudgementOfRighteousness, 0), 1.1},
		{"Seal of Fury's hit", byMask(paladin.SpellMaskSealOfFuryProc, 0), 1.1},
		{"Holy Shield", byMask(paladin.SpellMaskHolyShieldProc, 2), 1.1},
		{"Exorcism", byMask(paladin.SpellMaskExorcism, 0), 1.1},
		{"Consecration", byMask(paladin.SpellMaskConsecration, 0), 1},
		{"Retribution Aura", retributionAura, 1},
		{"Holy Light", byMask(paladin.SpellMaskHolyLight, 0), 1},
	} {
		if got := ratio(c.pick); math.Abs(got-c.want) > 1e-9 {
			t.Errorf("%s: one-hander over two-hander damage multiplier %v, want %v", c.name, got, c.want)
		}
	}
}

// The Holy-prot reference's talents: Holy Shock with One-Handed Weapon Specialization 3/3.
const holyProtTalents = "5500300300312125-05302133013-"

func newHolyProtPaladinWith(t *testing.T, equipment *proto.EquipmentSpec) *paladin.Paladin {
	t.Helper()
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: core.SinglePlayerRaidProto(&proto.Player{
			Name: "holy-prot", Class: proto.Class_ClassPaladin, Race: proto.Race_RaceHuman, TalentsString: holyProtTalents,
			Equipment: equipment,
			Spec:      DefaultOptions,
			Rotation:  &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim.Raid.Parties[0].Players[0].(paladin.PaladinAgent).GetPaladin()
}

// MythicSim patch 149: Holy Shock, damage and heal, takes no One-Handed Weapon Specialization until it
// is tested (Ashreon, 17 October 2026); Holy Strike beside it still does.
func TestOneHandedWeaponSpecializationSkipsHolyShock(t *testing.T) {
	ratios := func() (shock, heal, strike float64) {
		oneHand := newHolyProtPaladinWith(t, WeaponAndShield)
		twoHand := newHolyProtPaladinWith(t, twoHanderOnly)
		if oneHand.Talents.OneHandedWeaponSpecialization != 3 || !oneHand.Talents.HolyShock {
			t.Fatalf("test build has One-Handed Weapon Specialization %d, Holy Shock %v", oneHand.Talents.OneHandedWeaponSpecialization, oneHand.Talents.HolyShock)
		}
		ratio := func(mask int64) float64 {
			return spellByMask(t, oneHand, mask, 0).DamageMultiplier / spellByMask(t, twoHand, mask, 0).DamageMultiplier
		}
		return ratio(paladin.SpellMaskHolyShock), ratio(paladin.SpellMaskHolyShockHeal), ratio(paladin.SpellMaskHolyStrike)
	}
	shock, heal, strike := ratios()
	if shock != 1 || heal != 1 || math.Abs(strike-1.1) > 1e-9 {
		t.Errorf("one-hander over two-hander: Holy Shock %v, its heal %v, Holy Strike %v; want 1, 1, 1.1", shock, heal, strike)
	}

	paladin.OneHandedWeaponSpecializationSkipsHolyShock = false
	defer func() { paladin.OneHandedWeaponSpecializationSkipsHolyShock = true }()
	shock, heal, _ = ratios()
	if math.Abs(shock-1.1) > 1e-9 || heal != 1 {
		t.Errorf("switched off: Holy Shock %v, its heal %v; want 1.1 (patch 143) and 1", shock, heal)
	}
}
