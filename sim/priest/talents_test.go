package priest

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

func penanceCost(t *testing.T, talents string) float64 {
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceUndead,
		Class:         proto.Class_ClassPriest,
		Equipment:     &proto.EquipmentSpec{},
		Consumables:   &proto.ConsumesSpec{},
		TalentsString: talents,
		Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, &proto.Player_DpsPriest{DpsPriest: &proto.DpsPriest{Options: &proto.DpsPriest_Options{ClassOptions: &proto.PriestOptions{}}}})
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 100},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	priest := sim.Raid.Parties[0].Players[0].(PriestAgent).GetPriest()
	penance := priest.GetSpell(core.ActionID{SpellID: spellData.Penance.Highest().ID})
	if penance == nil {
		t.Fatalf("%s: no Penance", talents)
	}
	return penance.Cost.GetCurrentCost()
}

// Improved Healing 3/3 takes 15% off Penance.
func TestImprovedHealingDiscountsPenance(t *testing.T) {
	base := penanceCost(t, "504020031305001-13505100202-50002")
	discounted := penanceCost(t, "504020031305001-13505100032-50002")
	if !core.WithinToleranceFloat64(base*0.85, discounted, 1e-6) {
		t.Errorf("Penance costs %v with Improved Healing 3, %v without; want 15%% off", discounted, base)
	}
}

// The Smite rotations cast the Penance rank whose bolt is largest: rank 4 since client 1.60.1.70291 (92
// against rank 3's 72; in 1.60.1.70205 it was rank 3, 180 against 131). A client update that reorders the bolts
// fails here, and the rotations' Penance id wants changing with it.
func TestSmiteRotationsCastTheLargestPenanceRank(t *testing.T) {
	var best *spelldata.Spell
	spellData.Penance.Each(func(_ int32, rank *spelldata.Spell) {
		if best == nil || rank.Refs()[0].DamageEffect().Average(core.CharacterLevel) > best.Refs()[0].DamageEffect().Average(core.CharacterLevel) {
			best = rank
		}
	})
	want := fmt.Sprintf(`"spellId": %d,`, best.ID)
	for _, name := range []string{"smite", "smite_lowrank", "smite_lowrank_mindblast"} {
		apl, err := os.ReadFile("../../ui/specs/priest/dps/apls/" + name + ".apl.json")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(apl), want) {
			t.Errorf("%s does not cast Penance %d (rank %d), the largest bolt", name, best.ID, best.RankNumber())
		}
	}
}

// Chastise hits humanoids only, and its five ranks share one 2 minute cooldown.
func TestChastiseHumanoidsAndSharedCooldown(t *testing.T) {
	player := core.WithSpec(&proto.Player{
		Race:        proto.Race_RaceUndead,
		Class:       proto.Class_ClassPriest,
		Equipment:   &proto.EquipmentSpec{},
		Consumables: &proto.ConsumesSpec{},
		Rotation:    &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, &proto.Player_DpsPriest{DpsPriest: &proto.DpsPriest{Options: &proto.DpsPriest_Options{ClassOptions: &proto.PriestOptions{}}}})
	encounter := core.MakeSingleTargetEncounter(0)
	encounter.Targets = []*proto.Target{
		{Name: "beast", Level: 63, MobType: proto.MobType_MobTypeBeast},
		{Name: "humanoid", Level: 63, MobType: proto.MobType_MobTypeHumanoid},
	}
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 100},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  encounter,
	}, simsignals.CreateSignals())
	sim.Reset()

	priest := sim.Raid.Parties[0].Players[0].(PriestAgent).GetPriest()
	other, humanoid := sim.Encounter.AllTargetUnits[0], sim.Encounter.AllTargetUnits[1]
	top := priest.GetSpell(core.ActionID{SpellID: 1277335})
	low := priest.GetSpell(core.ActionID{SpellID: 1277331})
	if top.CanCast(sim, other) {
		t.Fatal("Chastise castable on a non-humanoid")
	}
	if !top.Cast(sim, humanoid) {
		t.Fatal("Chastise not castable on a humanoid")
	}
	if low.IsReady(sim) {
		t.Error("rank 1 ready after rank 5 was cast: the ranks should share the cooldown")
	}
	if dealt := top.SpellMetrics[humanoid.UnitIndex].TotalDamage; dealt < 271 {
		t.Errorf("Chastise dealt %.1f, want at least rank 5's 272", dealt)
	}
}

// Dark Sacrifice pays the row's 320 plus a fifth of Spirit on each of its five ticks, as a beta log
// shows at rank 1 (80 + 72 / 5 = 94.4 a tick).
func TestDarkSacrificeTicksBasePlusSpirit(t *testing.T) {
	player := core.WithSpec(&proto.Player{
		Race:        proto.Race_RaceUndead,
		Class:       proto.Class_ClassPriest,
		Equipment:   &proto.EquipmentSpec{},
		Consumables: &proto.ConsumesSpec{},
		Rotation:    &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, &proto.Player_DpsPriest{DpsPriest: &proto.DpsPriest{Options: &proto.DpsPriest_Options{ClassOptions: &proto.PriestOptions{}}}})
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 100},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	priest := sim.Raid.Parties[0].Players[0].(PriestAgent).GetPriest()
	spell := priest.GetSpell(core.ActionID{SpellID: 1277328})
	priest.SpendMana(sim, priest.CurrentMana(), priest.NewManaMetrics(core.ActionID{OtherID: proto.OtherAction_OtherActionNone}))
	spell.Cast(sim, nil)
	hot := spell.SelfHot()
	if hot.BaseTickCount != 5 {
		t.Fatalf("Dark Sacrifice ticks %d times, want 5", hot.BaseTickCount)
	}
	for range 5 {
		hot.TickOnce(sim)
	}
	want := 5 * (320 + priest.GetStat(stats.Spirit)/5)
	if got := priest.CurrentMana(); !core.WithinToleranceFloat64(want, got, 1e-6) {
		t.Errorf("Dark Sacrifice restored %.1f mana, want %.1f", got, want)
	}
	// Every rank is flagged No Threat (Attributes[1] 0x400), so the mana adds none.
	sim.Cleanup()
	if got := priest.GetSpell(core.ActionID{OtherID: proto.OtherAction_OtherActionManaGain}).SpellMetrics[0].TotalThreat; got != 0 {
		t.Errorf("Dark Sacrifice's mana added %.1f threat, want 0", got)
	}
}

// Holy Precision (1309957) and Holy Specialization (14889) are class-mask mods: Smite gets the hit and
// crit, Chastise (1277335) is outside both masks in client 1.60.1.70205 and gets neither.
func TestHolyTalentsSkipChastise(t *testing.T) {
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceUndead,
		Class:         proto.Class_ClassPriest,
		Equipment:     &proto.EquipmentSpec{},
		Consumables:   &proto.ConsumesSpec{},
		TalentsString: SmiteTalents,
		Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, &proto.Player_DpsPriest{DpsPriest: &proto.DpsPriest{Options: &proto.DpsPriest_Options{ClassOptions: &proto.PriestOptions{}}}})
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 100},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	priest := sim.Raid.Parties[0].Players[0].(PriestAgent).GetPriest()
	if priest.Talents.HolyPrecision == 0 || priest.Talents.HolySpecialization == 0 {
		t.Fatalf("SmiteTalents take Holy Precision %d, Holy Specialization %d; want both", priest.Talents.HolyPrecision, priest.Talents.HolySpecialization)
	}
	smite := priest.GetSpell(core.ActionID{SpellID: spellData.Smite.Highest().ID})
	chastise := priest.GetSpell(core.ActionID{SpellID: 1277335})
	if smite.BonusHitPercent == 0 || smite.BonusCritPercent == 0 {
		t.Errorf("Smite hit %v crit %v, want both raised", smite.BonusHitPercent, smite.BonusCritPercent)
	}
	if chastise.BonusHitPercent != 0 || chastise.BonusCritPercent != 0 || priest.PseudoStats.SchoolBonusHitChance[stats.SchoolIndexHoly] != 0 {
		t.Errorf("Chastise hit %v crit %v, Holy school hit %v; want none", chastise.BonusHitPercent, chastise.BonusCritPercent, priest.PseudoStats.SchoolBonusHitChance[stats.SchoolIndexHoly])
	}
}

// Client 1.60.1.70205 SpellShapeshift: Holy Nova and Chastise exclude form 28, which Shadowform puts
// the priest in. Smite, Holy Fire and Penance do not.
func TestShadowformRefusesHolyNovaAndChastise(t *testing.T) {
	newPriest := func(talents string) (*core.Simulation, *Priest) {
		player := core.WithSpec(&proto.Player{
			Race:          proto.Race_RaceUndead,
			Class:         proto.Class_ClassPriest,
			Equipment:     &proto.EquipmentSpec{},
			Consumables:   &proto.ConsumesSpec{},
			TalentsString: talents,
			Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.Player_DpsPriest{DpsPriest: &proto.DpsPriest{Options: &proto.DpsPriest_Options{ClassOptions: &proto.PriestOptions{}}}})
		encounter := core.MakeSingleTargetEncounter(0)
		encounter.Targets[0].MobType = proto.MobType_MobTypeHumanoid
		sim := core.NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 100},
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  encounter,
		}, simsignals.CreateSignals())
		sim.Reset()
		return sim, sim.Raid.Parties[0].Players[0].(PriestAgent).GetPriest()
	}

	sim, priest := newPriest(ShadowTalents)
	target := sim.Encounter.AllTargetUnits[0]
	chastise := priest.GetSpell(core.ActionID{SpellID: spellData.Chastise.Highest().ID})
	smite := priest.GetSpell(core.ActionID{SpellID: spellData.Smite.Highest().ID})
	priest.ShadowformAura.Activate(sim)
	if chastise.CanCast(sim, target) {
		t.Error("Chastise castable in Shadowform")
	}
	if !smite.CanCast(sim, target) {
		t.Error("Smite refused in Shadowform; the client does not exclude it")
	}
	priest.ShadowformAura.Deactivate(sim)
	if !chastise.CanCast(sim, target) {
		t.Error("Chastise still refused after Shadowform dropped")
	}

	sim, priest = newPriest("-000001") // Holy Nova only
	priest.ShapeshiftForm = spellData.Shadowform.Highest().ShapeshiftForm()
	holyNova := priest.GetSpell(core.ActionID{SpellID: spellData.HolyNova.Highest().ID})
	if holyNova.CanCast(sim, sim.Encounter.AllTargetUnits[0]) {
		t.Error("Holy Nova castable in Shadowform")
	}
}

// Holy Nova's party heal crits like any spell, and a crit reads the spell's DefenseType for its
// multiplier: without one the first heal crit panicked the sim.
func TestHolyNovaHealCanCrit(t *testing.T) {
	holyNova := &proto.ActionID{RawId: &proto.ActionID_SpellId{SpellId: spellData.HolyNova.Highest().ID}}
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceUndead,
		Class:         proto.Class_ClassPriest,
		Equipment:     &proto.EquipmentSpec{},
		Consumables:   &proto.ConsumesSpec{},
		TalentsString: "-000001", // Holy Nova only
		Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL, PriorityList: []*proto.APLListItem{
			{Action: &proto.APLAction{Action: &proto.APLAction_CastSpell{CastSpell: &proto.APLActionCastSpell{SpellId: holyNova}}}},
		}},
	}, &proto.Player_DpsPriest{DpsPriest: &proto.DpsPriest{Options: &proto.DpsPriest_Options{ClassOptions: &proto.PriestOptions{}}}})
	result := core.RunRaidSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{Iterations: 20, RandomSeed: 100},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	})
	if result.Error != nil {
		t.Fatal(result.Error.Message)
	}
}

// Client 1.60.1.70205: Shadowform (15473) has no proc or cancel rule and Power Infusion (10060) no
// form exclusion, so casting it keeps the form.
func TestPowerInfusionKeepsShadowform(t *testing.T) {
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceUndead,
		Class:         proto.Class_ClassPriest,
		Equipment:     &proto.EquipmentSpec{},
		Consumables:   &proto.ConsumesSpec{},
		TalentsString: "000000000000000001--000000000000000001", // Power Infusion + Shadowform
		Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, &proto.Player_DpsPriest{DpsPriest: &proto.DpsPriest{Options: &proto.DpsPriest_Options{ClassOptions: &proto.PriestOptions{}}}})
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 100},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	priest := sim.Raid.Parties[0].Players[0].(PriestAgent).GetPriest()

	priest.ShadowformAura.Activate(sim)
	priest.GetSpell(core.ActionID{SpellID: spellData.PowerInfusion.Highest().ID, Tag: priest.Index}).Cast(sim, &priest.Unit)
	if !priest.ShadowformAura.IsActive() {
		t.Error("Power Infusion dropped Shadowform")
	}
}

// Shadow Reach (17322) and Holy Reach (27789) add 10/20% to the range of the spells their client
// class masks reach; both were empty stubs, so the priest cast from as far as without them.
func TestReachTalentsExtendRange(t *testing.T) {
	for _, c := range []struct {
		talents string
		spell   int32
		want    float64
	}{
		{"", spellData.MindBlast.Highest().ID, 30},
		{"--000002", spellData.MindBlast.Highest().ID, 36}, // Shadow Reach 2/2
		{"--000002", spellData.Smite.Highest().ID, 30},
		{"-000000002", spellData.Smite.Highest().ID, 36}, // Holy Reach 2/2
		{"-000000002", spellData.MindBlast.Highest().ID, 30},
	} {
		player := core.WithSpec(&proto.Player{
			Race: proto.Race_RaceUndead, Class: proto.Class_ClassPriest, Equipment: &proto.EquipmentSpec{},
			Consumables: &proto.ConsumesSpec{}, TalentsString: c.talents, Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.Player_DpsPriest{DpsPriest: &proto.DpsPriest{Options: &proto.DpsPriest_Options{ClassOptions: &proto.PriestOptions{}}}})
		sim := core.NewSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{RandomSeed: 100},
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
		}, simsignals.CreateSignals())
		sim.Reset()
		priest := sim.Raid.Parties[0].Players[0].(PriestAgent).GetPriest()
		if got := priest.GetSpell(core.ActionID{SpellID: c.spell}).MaxRange; got != c.want {
			t.Errorf("%q: spell %d reaches %.1f yd, want %.1f", c.talents, c.spell, got, c.want)
		}
	}
}

// Shadow Weaving's stacks (15258) sit on the target and raise only the priest's Shadow damage there:
// five on one mob leave a second mob at +0%. Beta log 2726 shows each dotted mob with its own stacks.
func TestShadowWeavingStacksPerTarget(t *testing.T) {
	player := core.WithSpec(&proto.Player{
		Race: proto.Race_RaceUndead, Class: proto.Class_ClassPriest, Equipment: &proto.EquipmentSpec{},
		Consumables: &proto.ConsumesSpec{}, TalentsString: ShadowTalents, Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, &proto.Player_DpsPriest{DpsPriest: &proto.DpsPriest{Options: &proto.DpsPriest_Options{ClassOptions: &proto.PriestOptions{}}}})
	encounter := core.MakeSingleTargetEncounter(0)
	encounter.Targets = append(encounter.Targets, core.NewDefaultTarget())
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 100},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  encounter,
	}, simsignals.CreateSignals())
	sim.Reset()
	priest := sim.Raid.Parties[0].Players[0].(PriestAgent).GetPriest()
	first, second := sim.Encounter.AllTargetUnits[0], sim.Encounter.AllTargetUnits[1]

	aura := priest.ShadowWeavingAuras.Get(first)
	aura.Activate(sim)
	aura.SetStacks(sim, 5)

	mindBlast := priest.GetSpell(core.ActionID{SpellID: spellData.MindBlast.Highest().ID})
	for _, c := range []struct {
		target *core.Unit
		want   float64
	}{{first, 110}, {second, 100}} {
		result := &core.SpellResult{Target: c.target, Damage: 100}
		mindBlast.ApplyPostOutcomeDamageModifiers(sim, result, false)
		if result.Damage < c.want-1e-9 || result.Damage > c.want+1e-9 {
			t.Errorf("Mind Blast on %s: %.2f, want %.2f", c.target.Label, result.Damage, c.want)
		}
	}
}

// Inner Focus's +25% crit stays on the next spell but no longer reaches periodic ticks (dev
// AidanZMoon, forever-bugs #137): a Shadow Word: Pain running while it is up ticks at its own crit.
func TestInnerFocusSkipsPeriodicCrit(t *testing.T) {
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceUndead,
		Class:         proto.Class_ClassPriest,
		Equipment:     &proto.EquipmentSpec{},
		Consumables:   &proto.ConsumesSpec{},
		TalentsString: "000000001", // Inner Focus
		Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, &proto.Player_DpsPriest{DpsPriest: &proto.DpsPriest{Options: &proto.DpsPriest_Options{ClassOptions: &proto.PriestOptions{}}}})
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 100},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	priest := sim.Raid.Parties[0].Players[0].(PriestAgent).GetPriest()
	target := priest.CurrentTarget
	swp := priest.GetSpell(core.ActionID{SpellID: spellData.ShadowWordPain.Highest().ID})

	priest.InnerFocusAura.Activate(sim)
	withIF := swp.SpellCritChance(target)
	tick := priest.tickOutcome(true, swp.Dot(target))
	crits := 0
	for range 4000 {
		result := &core.SpellResult{Target: target, Damage: 1}
		tick(sim, result, priest.AttackTables[target.UnitIndex])
		if result.Outcome == core.OutcomeCrit {
			crits++
		}
	}
	if rate := float64(crits) / 4000; rate > withIF-0.2 {
		t.Errorf("ticks crit %.3f with Inner Focus up (spell %.3f), want its 25%% left out", rate, withIF)
	}
}
