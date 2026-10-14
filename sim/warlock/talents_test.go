package warlock

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

// The Demonic Pact 2/31/18 preset: Demonic Sacrifice, Soul Link and Demonic Pact.
const pactTalents = "113-0005003221220311351-0550005"

func sacrificeAura(t *testing.T, talents string, pact proto.WarlockOptions_Summon) *core.Aura {
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceOrc,
		Class:         proto.Class_ClassWarlock,
		Equipment:     &proto.EquipmentSpec{},
		Consumables:   &proto.ConsumesSpec{},
		TalentsString: talents,
		Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, &proto.Player_Warlock{Warlock: &proto.Warlock{Options: &proto.Warlock_Options{ClassOptions: &proto.WarlockOptions{
		Summon:        proto.WarlockOptions_Succubus,
		PactSacrifice: pact,
	}}}})
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 100},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim.Raid.Parties[0].Players[0].(WarlockAgent).GetWarlock().GetAura("Demonic Sacrifice")
}

// Demonic Pact keeps a sacrificed Imp's Shadow buff with the Succubus out, but not a sacrificed
// Succubus (resummoning the sacrificed demon cancels it), and nothing without the talent.
func TestDemonicPactKeepsTheSacrifice(t *testing.T) {
	if aura := sacrificeAura(t, pactTalents, proto.WarlockOptions_Imp); aura == nil || aura.ActionID.SpellID != 18789 {
		t.Errorf("Pact + sacrificed Imp: got %v, want Touch of Shadow (18789)", aura)
	}
	if aura := sacrificeAura(t, pactTalents, proto.WarlockOptions_Succubus); aura != nil {
		t.Errorf("Pact + sacrificed Succubus with the Succubus out: got %v, want none", aura)
	}
	if aura := sacrificeAura(t, "113-0005003221220311350-0550005", proto.WarlockOptions_Imp); aura != nil {
		t.Errorf("no Demonic Pact: got %v, want none", aura)
	}
}

// The Imp has no melee, so its Firebolt spends the brand, and the branded hit is Fire (1293698);
// the other demons' is Shadow (1293697). Both rows are Always Hit without Cannot Crit, so the hit
// never misses and can crit.
func TestDemonicBrandImpSpendsWithFirebolt(t *testing.T) {
	for _, tc := range []struct {
		summon proto.WarlockOptions_Summon
		school core.SpellSchool
	}{{proto.WarlockOptions_Imp, core.SpellSchoolFire}, {proto.WarlockOptions_Succubus, core.SpellSchoolShadow}} {
		player := core.WithSpec(&proto.Player{
			Race:          proto.Race_RaceOrc,
			Class:         proto.Class_ClassWarlock,
			Equipment:     &proto.EquipmentSpec{},
			Consumables:   &proto.ConsumesSpec{},
			TalentsString: "-00000000000003",
			Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.Player_Warlock{Warlock: &proto.Warlock{Options: &proto.Warlock_Options{ClassOptions: &proto.WarlockOptions{
			Summon: tc.summon,
		}}}})
		sim := core.NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 100},
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
		}, simsignals.CreateSignals())
		sim.Reset()

		pet := sim.Raid.Parties[0].Players[0].(WarlockAgent).GetWarlock().ActivePet
		brand := pet.DemonicBrandAura
		brand.Activate(sim)
		brand.SetStacks(sim, brand.MaxStacks)
		for _, spell := range pet.AutoCastAbilities {
			spell.SkipCastAndApplyEffects(sim, pet.CurrentTarget)
		}

		hit := pet.GetSpell(brand.ActionID)
		if hit.SpellSchool != tc.school {
			t.Errorf("%v: branded hit is %v, want %v", tc.summon, hit.SpellSchool, tc.school)
		}
		if hit.SpellMetrics[0].Casts == 0 || brand.GetStacks() != brand.MaxStacks-hit.SpellMetrics[0].Casts {
			t.Errorf("%v: pet spell spent %d of %d charges", tc.summon, brand.MaxStacks-brand.GetStacks(), brand.MaxStacks)
		}

		for i := 0; i < 2000; i++ {
			hit.SkipCastAndApplyEffects(sim, pet.CurrentTarget)
		}
		if m := hit.SpellMetrics[0]; m.Misses != 0 || m.Crits == 0 {
			t.Errorf("%v: branded hit missed %d and crit %d times in 2000, want 0 misses and some crits", tc.summon, m.Misses, m.Crits)
		}
	}
}

// Bane of Havoc copies 15% of the warlock's damage to other targets onto the baned target (client
// 1225228), and nothing of the damage the baned target takes itself.
func TestBaneOfHavocCopiesDamage(t *testing.T) {
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceOrc,
		Class:         proto.Class_ClassWarlock,
		Equipment:     &proto.EquipmentSpec{},
		Consumables:   &proto.ConsumesSpec{},
		TalentsString: "--0000000000001",
		Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, &proto.Player_Warlock{Warlock: &proto.Warlock{Options: &proto.Warlock_Options{ClassOptions: &proto.WarlockOptions{
		Summon: proto.WarlockOptions_NoSummon,
	}}}})
	encounter := core.MakeSingleTargetEncounter(0)
	encounter.Targets = append(encounter.Targets, core.NewDefaultTarget())
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 100},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  encounter,
	}, simsignals.CreateSignals())
	sim.Reset()

	warlock := sim.Raid.Parties[0].Players[0].(WarlockAgent).GetWarlock()
	main, havoc := sim.Encounter.AllTargetUnits[0], sim.Encounter.AllTargetUnits[1]
	bane := warlock.GetSpell(core.ActionID{SpellID: 1225228})
	copied := warlock.GetSpell(core.ActionID{SpellID: 1225228, Tag: 1})
	for i := 0; i < 20 && !havoc.HasActiveAura("Bane of Havoc-"+warlock.Label); i++ {
		bane.SkipCastAndApplyEffects(sim, havoc)
	}

	for i := 0; i < 20 && warlock.SearingPain.SpellMetrics[main.UnitIndex].TotalDamage == 0; i++ {
		warlock.SearingPain.SkipCastAndApplyEffects(sim, main)
	}
	dealt := warlock.SearingPain.SpellMetrics[main.UnitIndex].TotalDamage
	got := copied.SpellMetrics[havoc.UnitIndex].TotalDamage
	if dealt == 0 || math.Abs(got-dealt*0.15) > 1e-6 {
		t.Fatalf("Searing Pain dealt %.2f to the main target, bane copied %.2f, want %.2f", dealt, got, dealt*0.15)
	}

	warlock.SearingPain.SkipCastAndApplyEffects(sim, havoc)
	if again := copied.SpellMetrics[havoc.UnitIndex].TotalDamage; again != got {
		t.Errorf("damage to the baned target itself copied: %.2f -> %.2f", got, again)
	}
}

// The warlock's own Imp gives its party Blood Pact (11767); a sacrificed Imp or another demon doesn't.
func TestImpGivesBloodPact(t *testing.T) {
	for _, tc := range []struct {
		summon    proto.WarlockOptions_Summon
		sacrifice bool
		want      bool
	}{{proto.WarlockOptions_Imp, false, true}, {proto.WarlockOptions_Imp, true, false}, {proto.WarlockOptions_Succubus, false, false}} {
		player := core.WithSpec(&proto.Player{
			Race:        proto.Race_RaceOrc,
			Class:       proto.Class_ClassWarlock,
			Equipment:   &proto.EquipmentSpec{},
			Consumables: &proto.ConsumesSpec{},
			Rotation:    &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.Player_Warlock{Warlock: &proto.Warlock{Options: &proto.Warlock_Options{ClassOptions: &proto.WarlockOptions{
			Summon:          tc.summon,
			SacrificeSummon: tc.sacrifice,
		}}}})
		sim := core.NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 100},
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
		}, simsignals.CreateSignals())
		sim.Reset()
		aura := sim.Raid.Parties[0].Players[0].(WarlockAgent).GetWarlock().GetAuraByID(core.ActionID{SpellID: 11767, Tag: -1})
		if got := aura != nil && aura.IsActive(); got != tc.want {
			t.Errorf("summon %v, sacrificed %v: Blood Pact up = %v, want %v", tc.summon, tc.sacrifice, got, tc.want)
		}
	}
}

// The Imp's Firebolt rolls client 11763's row: 45 +-11.4% at 60, and the level scaling stops at 63.
func TestImpFireboltRollsItsRow(t *testing.T) {
	if got := impFireboltEffect.Average(60); got != 45 {
		t.Errorf("average at 60 = %v, want 45", got)
	}
	if lo, hi := impFireboltEffect.Min(60), impFireboltEffect.Max(60); math.Abs(lo-42.44) > 0.01 || math.Abs(hi-47.56) > 0.01 {
		t.Errorf("roll at 60 = %.2f to %.2f, want 42.44 to 47.56", lo, hi)
	}
	if got := impFireboltEffect.Average(70); got != 47 {
		t.Errorf("average at 70 = %v, want 47 (capped at 63)", got)
	}
}

// Searing Pain rolls the spread its client row states, as every warlock nuke now does, where it
// used to deal the row's average on every cast.
func TestSearingPainRollsItsRow(t *testing.T) {
	player := core.WithSpec(&proto.Player{
		Race:        proto.Race_RaceOrc,
		Class:       proto.Class_ClassWarlock,
		Equipment:   &proto.EquipmentSpec{},
		Consumables: &proto.ConsumesSpec{},
		Rotation:    &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, &proto.Player_Warlock{Warlock: &proto.Warlock{Options: &proto.Warlock_Options{ClassOptions: &proto.WarlockOptions{
		Summon: proto.WarlockOptions_NoSummon,
	}}}})
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	warlock := sim.Raid.Parties[0].Players[0].(WarlockAgent).GetWarlock()
	m := &warlock.SearingPain.SpellMetrics[0]
	seen := map[float64]bool{}
	for i := 0; i < 30; i++ {
		damage, hits, crits, resisted := m.TotalDamage, m.Hits, m.Crits, m.ResistedHits
		warlock.SearingPain.SkipCastAndApplyEffects(sim, warlock.CurrentTarget)
		if m.Hits > hits && m.Crits == crits && m.ResistedHits == resisted {
			seen[math.Round(m.TotalDamage-damage)] = true
		}
	}
	if len(seen) < 3 {
		t.Errorf("Searing Pain non-crits dealt only %v; it should roll its row's spread", seen)
	}
}

// Life Tap's mana adds no threat: every rank is flagged No Threat in the client (Attributes[1] 0x400),
// while an ordinary mana gain still adds half its amount.
func TestLifeTapAddsNoThreat(t *testing.T) {
	player := core.WithSpec(&proto.Player{
		Race:        proto.Race_RaceOrc,
		Class:       proto.Class_ClassWarlock,
		Equipment:   &proto.EquipmentSpec{},
		Consumables: &proto.ConsumesSpec{},
		Rotation:    &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, &proto.Player_Warlock{Warlock: &proto.Warlock{Options: &proto.Warlock_Options{ClassOptions: &proto.WarlockOptions{
		Summon: proto.WarlockOptions_NoSummon,
	}}}})
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	warlock := sim.Raid.Parties[0].Players[0].(WarlockAgent).GetWarlock()
	threat := func() float64 {
		return warlock.GetSpell(core.ActionID{OtherID: proto.OtherAction_OtherActionManaGain}).SpellMetrics[0].TotalThreat
	}
	warlock.SpendMana(sim, warlock.CurrentMana(), warlock.NewManaMetrics(core.ActionID{OtherID: proto.OtherAction_OtherActionNone}))
	warlock.LifeTap.SkipCastAndApplyEffects(sim, warlock.CurrentTarget)
	if warlock.CurrentMana() <= 0 {
		t.Fatal("Life Tap restored no mana")
	}
	sim.Cleanup()
	if got := threat(); got != 0 {
		t.Errorf("Life Tap's mana added %.1f threat, want 0", got)
	}

	warlock.AddMana(sim, 100, warlock.NewManaMetrics(core.ActionID{SpellID: 20268}))
	sim.Cleanup()
	if got := threat(); got != 50 {
		t.Errorf("an ordinary 100 mana gain added %.1f threat, want 50", got)
	}
}

// Destructive Reach (17917) adds 10/20% to the range of the warlock spells its client class mask
// reaches, Affliction ones included but not Curse of Recklessness; it was an empty stub and the
// spells carried no range, so a warlock cast from any distance.
func TestDestructiveReachExtendsRange(t *testing.T) {
	for _, c := range []struct {
		talents string
		spell   int32
		want    float64
	}{
		{"", spellData.ShadowBolt.Highest().ID, 30},
		{"", spellData.DrainLife.Highest().ID, 20},
		{"--2", spellData.ShadowBolt.Highest().ID, 36},
		{"--2", spellData.Corruption.Highest().ID, 36},
		{"--2", spellData.DrainLife.Highest().ID, 24},
		{"--1", spellData.Immolate.Highest().ID, 33},
		{"--2", spellData.CurseOfRecklessness.Highest().ID, 30},
	} {
		player := core.WithSpec(&proto.Player{
			Race: proto.Race_RaceOrc, Class: proto.Class_ClassWarlock, Equipment: &proto.EquipmentSpec{},
			Consumables: &proto.ConsumesSpec{}, TalentsString: c.talents, Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.Player_Warlock{Warlock: &proto.Warlock{Options: &proto.Warlock_Options{ClassOptions: &proto.WarlockOptions{
			Summon: proto.WarlockOptions_NoSummon,
		}}}})
		sim := core.NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 100},
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
		}, simsignals.CreateSignals())
		sim.Reset()
		warlock := sim.Raid.Parties[0].Players[0].(WarlockAgent).GetWarlock()
		if got := warlock.GetSpell(core.ActionID{SpellID: c.spell}).MaxRange; math.Abs(got-c.want) > 1e-9 {
			t.Errorf("%q: spell %d reaches %.1f yd, want %.1f", c.talents, c.spell, got, c.want)
		}
	}
}
