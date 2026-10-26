package hunter

import (
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/simsignals"
	"github.com/wowsims/classic/sim/core/stats"
	"google.golang.org/protobuf/encoding/protojson"
	googleProto "google.golang.org/protobuf/proto"
)

// The Survival weave talents (Strider Kick, Resourcefulness).
const patchTestTalents = "-00531500101-500230230250222151"

type patchTestSetup struct {
	race        proto.Race
	petType     proto.Hunter_Options_PetType
	buffs       *proto.RaidBuffs
	talents     string
	gear        bool
	distance    float64
	stone       bool
	ibuffs      *proto.IndividualBuffs
	debuffs     *proto.Debuffs
	petBug      bool // the old "Beta pet speed bug" option
	autoSpam    bool
	gearSet     string // a gear set other than 2h_weave (implies gear)
	mhImbue     proto.WeaponImbue
	ohImbue     proto.WeaponImbue
	ohItem      int32 // replaces the gear set's off hand
	runSpeed    float64
	consumes    *proto.Consumes          // starting consumables (the stone and imbue fields still apply on top)
	mhEnchant   int32                    // replaces the gear set's main hand enchant
	rngEnchant  int32                    // replaces the gear set's ranged enchant
	skyborne    bool                     // the Skyborne 15 min racial option
	outside     bool                     // the encounter is outside the Natural flasks' zones
	pieces      map[proto.ItemSlot]int32 // items put in these slots (on top of any gear set)
	engineer    bool                     // Engineering as a profession
	msSkipsAuto bool                     // the Mortal Shots skips Auto Shot option
	meleeResets bool                     // the beta bug: melee swings reset Auto Shot
	shotsDelay  bool                     // the beta bug: a shot in an Auto Shot cycle pushes it back
	quiver      proto.Hunter_Options_QuiverBonus
	mobType     proto.MobType    // the target's creature type
	profession  proto.Profession // a second profession
	manualStart bool             // the rotation starts the melee attack itself (manual_start_attack)
}

func newPatchTestHunter(t *testing.T, setup patchTestSetup) (*core.Simulation, *Hunter) {
	if setup.race == proto.Race_RaceUnknown {
		setup.race = proto.Race_RaceHuman
	}
	if setup.petType == proto.Hunter_Options_PetNone {
		setup.petType = proto.Hunter_Options_Cat
	}
	if setup.buffs == nil {
		setup.buffs = &proto.RaidBuffs{}
	}
	if setup.ibuffs == nil {
		setup.ibuffs = &proto.IndividualBuffs{}
	}
	if setup.debuffs == nil {
		setup.debuffs = &proto.Debuffs{}
	}
	if setup.talents == "" {
		setup.talents = "50230005001-005005-5002302300500221"
	}
	equipment := &proto.EquipmentSpec{}
	if setup.gearSet != "" {
		equipment = core.GetGearSet("../../ui/hunter/gear_sets", setup.gearSet).GearSet
	} else if setup.gear {
		equipment = core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet
	}
	if len(setup.pieces) > 0 && len(equipment.Items) == 0 {
		equipment.Items = make([]*proto.ItemSpec, 17)
		for i := range equipment.Items {
			equipment.Items[i] = &proto.ItemSpec{}
		}
	}
	for slot, id := range setup.pieces {
		equipment.Items[slot] = &proto.ItemSpec{Id: id}
	}
	consumes := &proto.Consumes{}
	if setup.consumes != nil {
		consumes = setup.consumes
	}
	if setup.mhEnchant != 0 {
		equipment.Items[proto.ItemSlot_ItemSlotMainHand].Enchant = setup.mhEnchant
	}
	if setup.rngEnchant != 0 {
		equipment.Items[proto.ItemSlot_ItemSlotRanged].Enchant = setup.rngEnchant
	}
	if setup.stone {
		consumes.MainHandImbue = proto.WeaponImbue_ElementalSharpeningStone
	}
	if setup.mhImbue != proto.WeaponImbue_WeaponImbueUnknown {
		consumes.MainHandImbue = setup.mhImbue
	}
	consumes.OffHandImbue = setup.ohImbue
	if setup.ohItem != 0 {
		equipment.Items[proto.ItemSlot_ItemSlotOffHand] = &proto.ItemSpec{Id: setup.ohItem}
	}
	player := &proto.Player{
		Name:               "Player",
		Race:               setup.race,
		Class:              proto.Class_ClassHunter,
		Equipment:          equipment,
		Consumes:           consumes,
		TalentsString:      setup.talents,
		DistanceFromTarget: setup.distance,
		// an empty rotation: a main-hand swing asks it for the next action (Raptor Strike)
		Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL, ManualStartAttack: setup.manualStart},
		Buffs:    setup.ibuffs,
		Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{
			PetType:                  setup.petType,
			PetUptime:                1,
			PetAttackSpeed:           proto.Hunter_Options_One,
			PetWindfuryTotem:         true,
			PetFlametongueTotem:      true,
			PetSpeedBug:              setup.petBug,
			AutoShotSpamBug:          setup.autoSpam,
			RunSpeed:                 setup.runSpeed,
			SkyborneConvergence:      setup.skyborne,
			MortalShotsSkipsAutoShot: setup.msSkipsAuto,
			MeleeResetsAutoShot:      setup.meleeResets,
			ShotsDelayAutoShot:       setup.shotsDelay,
			QuiverBonus:              setup.quiver,
		}}},
	}
	if setup.engineer {
		player.Profession1 = proto.Profession_Engineering
	}
	player.Profession2 = setup.profession
	encounter := core.MakeSingleTargetEncounter(0)
	encounter.OutsideHyjal = setup.outside
	encounter.Targets[0].MobType = setup.mobType
	sim := core.NewSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, setup.buffs, setup.debuffs),
		Encounter:  encounter,
		SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 1},
	}, simsignals.CreateSignals())
	sim.Reset()
	hunter, ok := sim.Raid.Parties[0].Players[0].(*Hunter)
	if !ok {
		t.Fatal("no hunter")
	}
	if hunter.pet == nil && !hunter.Talents.LoneWolf {
		t.Fatal("no pet")
	}
	return sim, hunter
}

// Windfury Totem and Flametongue Totem do not stack (Windfury wins); Grace of Air stacks with both; pets get
// neither weapon totem, even with the old pet totem switches on.
func TestWeaponTotems(t *testing.T) {
	for _, tc := range []struct {
		name        string
		buffs       *proto.RaidBuffs
		ft, goa, wf bool
	}{
		{"Windfury + Flametongue", &proto.RaidBuffs{WindfuryTotem: true, FlametongueTotem: true}, false, false, true},
		{"Windfury + Grace of Air", &proto.RaidBuffs{WindfuryTotem: true, GraceOfAirTotem: proto.TristateEffect_TristateEffectRegular}, false, true, true},
		{"Flametongue + Grace of Air", &proto.RaidBuffs{FlametongueTotem: true, GraceOfAirTotem: proto.TristateEffect_TristateEffectRegular}, true, true, false},
	} {
		_, hunter := newPatchTestHunter(t, patchTestSetup{buffs: tc.buffs})
		if got := hunter.HasAura("Flametongue Totem"); got != tc.ft {
			t.Errorf("%s: Flametongue Totem %v, want %v", tc.name, got, tc.ft)
		}
		if got := hunter.HasAura("Grace of Air Totem"); got != tc.goa {
			t.Errorf("%s: Grace of Air %v, want %v", tc.name, got, tc.goa)
		}
		if got := hunter.HasAura("Windfury"); got != tc.wf {
			t.Errorf("%s: Windfury %v, want %v", tc.name, got, tc.wf)
		}
		if hunter.pet.HasAura("Flametongue Totem") || hunter.pet.HasAura("Windfury") {
			t.Errorf("%s: the pet got a weapon totem", tc.name)
		}
	}
}

// Badge of the Swarmguard's armor penetration stacks (Insight of the Qiraji, 26481) have no duration of their own in the
// client data: they last only while the badge's 30 s buff (26480) is up, as its tooltip's "for 30 sec" says.
func TestSwarmguardStacksEndWithTheBadge(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, pieces: map[proto.ItemSlot]int32{proto.ItemSlot_ItemSlotTrinket1: 21670}})
	badge, insight := hunter.GetAura("Badge of the Swarmguard"), hunter.GetAura("Insight of the Qiraji")
	if badge == nil || insight == nil {
		t.Fatal("Badge of the Swarmguard auras not registered")
	}
	sim.PrePull()
	// past the pull's own use (30 s), with the badge still on its 3 min cooldown
	for i := 0; i < 1000000 && sim.CurrentTime < 40*time.Second; i++ {
		sim.Step()
	}
	base := hunter.GetStat(stats.ArmorPenetration)
	badge.Activate(sim)
	insight.Activate(sim)
	insight.SetStacks(sim, 6)
	if got := hunter.GetStat(stats.ArmorPenetration) - base; math.Abs(got-1200) > 1e-9 {
		t.Fatalf("6 stacks give %.0f armor penetration, want 1200", got)
	}
	for i := 0; i < 1000000 && sim.CurrentTime < 75*time.Second; i++ {
		sim.Step()
	}
	if insight.IsActive() || hunter.GetStat(stats.ArmorPenetration) != base {
		t.Errorf("35 s after the badge's 30 s buff began: stacks active %v, armor penetration %.0f over base, want gone",
			insight.IsActive(), hunter.GetStat(stats.ArmorPenetration)-base)
	}
}

// Furious Howl gives the hunter 40% less attack power than the client value of 83.6.
func TestFuriousHowlAttackPower(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{petType: proto.Hunter_Options_Wolf})
	howl := hunter.GetAura("Furious Howl")
	if howl == nil {
		t.Fatal("no Furious Howl aura on the hunter")
	}
	before := hunter.GetStat(stats.AttackPower)
	howl.Activate(sim)
	if got := hunter.GetStat(stats.AttackPower) - before; math.Abs(got-83.6*0.6) > 1e-9 {
		t.Errorf("Furious Howl +%.2f attack power, want %.2f", got, 83.6*0.6)
	}
}

// Pet ability cooldowns: Sonic Blast (bat) 1 min, Dismember (crocolisk) 15 s, Mine! (bird of prey) 90 s.
func TestPetAbilityCooldowns(t *testing.T) {
	for _, tc := range []struct {
		pet     proto.Hunter_Options_PetType
		spellID int32
		cd      time.Duration
	}{
		{proto.Hunter_Options_Bat, 1264482, time.Minute},
		{proto.Hunter_Options_Crocolisk, 1264933, time.Second * 15},
		{proto.Hunter_Options_BirdOfPrey, 1265058, time.Second * 90},
	} {
		_, hunter := newPatchTestHunter(t, patchTestSetup{petType: tc.pet})
		spell := hunter.pet.GetSpell(core.ActionID{SpellID: tc.spellID})
		if spell == nil {
			t.Errorf("%v: no spell %d", tc.pet, tc.spellID)
			continue
		}
		if spell.CD.Duration != tc.cd {
			t.Errorf("%v: spell %d cooldown %v, want %v", tc.pet, tc.spellID, spell.CD.Duration, tc.cd)
		}
	}
}

// Racials after Blizzard's tuning pass of 26 October 2026: Ruthless (Human) +1% crit; Orc has no crit (Savage is gone)
// but Command, +5% pet damage; Dwarf has Artisan's Precision, +1% hit, and no critical damage (Earth Fury is gone);
// Tauren's Endurance gives no critical damage or hit; Wind Blessed (Skyborne) 2% haste; Blood Fury on a 3 min cooldown.
func TestRacials(t *testing.T) {
	critPerAgi := core.CritPerAgiAtLevel[proto.Class_ClassHunter] * core.CritRatingPerCritChance
	racialCrit := func(race proto.Race) float64 {
		_, hunter := newPatchTestHunter(t, patchTestSetup{race: race})
		return hunter.GetStat(stats.MeleeCrit) - hunter.GetStat(stats.Agility)*critPerAgi
	}
	nightElf := racialCrit(proto.Race_RaceNightElf)
	for race, want := range map[proto.Race]float64{proto.Race_RaceHuman: 1, proto.Race_RaceOrc: 0, proto.Race_RaceDwarf: 0, proto.Race_RaceTauren: 0} {
		if got := racialCrit(race) - nightElf; math.Abs(got-want) > 1e-6 {
			t.Errorf("%v: +%.3f%% racial crit, want %.0f", race, got, want)
		}
	}
	for race, want := range map[proto.Race]float64{proto.Race_RaceDwarf: 1, proto.Race_RaceTauren: 1, proto.Race_RaceHuman: 1, proto.Race_RaceOrc: 1} {
		_, hunter := newPatchTestHunter(t, patchTestSetup{race: race})
		if got := hunter.PseudoStats.MeleeCritMultiplier; math.Abs(got-want) > 1e-9 {
			t.Errorf("%v: critical damage multiplier %.3f, want %.2f", race, got, want)
		}
	}
	_, tauren := newPatchTestHunter(t, patchTestSetup{race: proto.Race_RaceTauren})
	_, nelf := newPatchTestHunter(t, patchTestSetup{race: proto.Race_RaceNightElf})
	if tauren.GetStat(stats.MeleeHit) != nelf.GetStat(stats.MeleeHit) {
		t.Errorf("Tauren still gets extra hit: %.2f vs %.2f", tauren.GetStat(stats.MeleeHit), nelf.GetStat(stats.MeleeHit))
	}
	_, human := newPatchTestHunter(t, patchTestSetup{race: proto.Race_RaceHuman})
	baseSpirit := core.ClassBaseStats[proto.Class_ClassHunter][stats.Spirit] + core.RaceOffsets[proto.Race_RaceHuman][stats.Spirit]
	if got := human.GetStat(stats.Spirit); math.Abs(got-baseSpirit*1.10) > 1e-6 {
		t.Errorf("Human Spirit %.2f from a base of %.0f, want %.2f (The Human Spirit, +10%%)", got, baseSpirit, baseSpirit*1.10)
	}
	_, dwarf := newPatchTestHunter(t, patchTestSetup{race: proto.Race_RaceDwarf})
	for _, s := range []stats.Stat{stats.MeleeHit, stats.SpellHit} {
		if got := dwarf.GetStat(s) - nelf.GetStat(s); math.Abs(got-1) > 1e-9 {
			t.Errorf("Dwarf %v +%.2f over Night Elf, want +1 (Artisan's Precision)", s, got)
		}
	}
	_, skyborne := newPatchTestHunter(t, patchTestSetup{race: proto.Race_RaceSkyborneAlliance})
	if got := skyborne.PseudoStats.RangedSpeedMultiplier / nelf.PseudoStats.RangedSpeedMultiplier; math.Abs(got-1.02) > 1e-9 {
		t.Errorf("Skyborne ranged speed x%.4f, want 1.02", got)
	}
	_, orc := newPatchTestHunter(t, patchTestSetup{race: proto.Race_RaceOrc})
	if got := orc.pet.PseudoStats.DamageDealtMultiplier / nelf.pet.PseudoStats.DamageDealtMultiplier; math.Abs(got-1.05) > 1e-9 {
		t.Errorf("Orc pet damage x%.4f over Night Elf, want 1.05 (Command)", got)
	}
	if bf := orc.GetSpell(core.ActionID{SpellID: 20572}); bf == nil || bf.CD.Duration != time.Minute*3 {
		t.Errorf("Blood Fury cooldown is not 3 min")
	}
}

// Quickness (Night Elf) stacks with other speed bonuses: Strider Kick's 30% on top of it.
func TestQuicknessStacksWithStriderKick(t *testing.T) {
	sim, nelf := newPatchTestHunter(t, patchTestSetup{race: proto.Race_RaceNightElf, talents: patchTestTalents})
	kick := nelf.GetAura("Strider Kick Speed")
	if kick == nil {
		t.Fatal("no Strider Kick speed aura (talent string)")
	}
	kick.Activate(sim)
	if got := nelf.MovementHandler.MoveSpeed; math.Abs(got-7*1.02*1.3) > 1e-9 {
		t.Errorf("Night Elf under Strider Kick runs %.4f yd/s, want %.4f", got, 7*1.02*1.3)
	}
}

// Auto Shot fires 0.6 s after a move back to range at the soonest: the windup and the moment to start it.
func TestAutoShotStandStill(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 12})
	sim.PrePull()
	for guard := 0; sim.CurrentTime < time.Second*3 && guard < 100000; guard++ {
		sim.Step()
	}
	step := func(to float64) time.Duration {
		hunter.MoveTo(to, sim)
		for guard := 0; hunter.IsMoving() && guard < 100000; guard++ {
			sim.Step()
		}
		return sim.CurrentTime
	}
	step(3)
	stoppedAt := step(12)
	if next := hunter.AutoAttacks.NextRangedAttackAt(); next < stoppedAt+time.Millisecond*600 {
		t.Errorf("Auto Shot due %.3fs after stopping at range, want at least 0.6 s", (next - stoppedAt).Seconds())
	}
}

// Pets take the higher of the hunter's melee and ranged crit (the main-hand stone counts on the melee side)
// and inherit the hunter's expertise.
func TestPetInheritance(t *testing.T) {
	_, plain := newPatchTestHunter(t, patchTestSetup{gear: true})
	_, stoned := newPatchTestHunter(t, patchTestSetup{gear: true, stone: true})
	if got := stoned.pet.GetStat(stats.MeleeCrit) - plain.pet.GetStat(stats.MeleeCrit); math.Abs(got-2) > 1e-9 {
		t.Errorf("main-hand Elemental Sharpening Stone gives the pet +%.2f crit, want 2", got)
	}
	// the stone crit is added once: the hunter gaining and losing stats must not pile it up on the pet
	sim, again := newPatchTestHunter(t, patchTestSetup{gear: true, stone: true})
	before := again.pet.GetStat(stats.MeleeCrit)
	for i := 0; i < 3; i++ {
		again.AddStatsDynamic(sim, stats.Stats{stats.AttackPower: 100})
		again.AddStatsDynamic(sim, stats.Stats{stats.AttackPower: -100})
	}
	if got := again.pet.GetStat(stats.MeleeCrit); math.Abs(got-before) > 1e-9 {
		t.Errorf("pet crit drifted from %.2f to %.2f as the hunter gained and lost attack power", before, got)
	}
	inherited := plain.makeStatInheritance()(stats.Stats{stats.Expertise: 5})
	if inherited[stats.Expertise] != 5 {
		t.Errorf("pet inherits %.1f expertise of 5", inherited[stats.Expertise])
	}
}

// A walk decided and still in its move reaction rules out a spell with a cast time: the walk would cancel it, and
// starting one there let it finish mid-walk as if the hunter stood still. Instant shots stay usable.
func TestNoCastDuringMoveReaction(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 12})
	hunter.MoveDelay = 100 * time.Millisecond
	target := hunter.CurrentTarget
	if !hunter.MultiShot.CanCast(sim, target) {
		t.Fatal("Multi-Shot not castable standing at 12 yd")
	}
	hunter.MoveTo(4, sim)
	if !hunter.MovementHandler.PendingMove {
		t.Fatal("no move reaction pending after deciding to walk in")
	}
	if hunter.MultiShot.CanCast(sim, target) {
		t.Error("Multi-Shot, a cast, can start during the move reaction")
	}
	if !hunter.ArcaneShot.CanCast(sim, target) {
		t.Error("Arcane Shot, an instant, cannot be used during the move reaction")
	}
}

// The pet's attack power is 10% of the hunter's higher attack power at every moment: a buff that ends comes off the
// pet too, whichever of the two it raised, so temporary attack power never piles up on the pet over a fight.
func TestPetAttackPowerFollowsOwner(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true})
	share := func() float64 {
		return max(hunter.GetStat(stats.AttackPower), hunter.GetStat(stats.RangedAttackPower)) * 0.10
	}
	start, startShare := hunter.pet.GetStat(stats.AttackPower), share()
	for _, buff := range []stats.Stats{
		{stats.AttackPower: 100},
		{stats.RangedAttackPower: 100},
		{stats.AttackPower: 3000},
		{stats.AttackPower: 50, stats.RangedAttackPower: 120},
	} {
		for i := 0; i < 3; i++ {
			hunter.AddStatsDynamic(sim, buff)
			if got, want := hunter.pet.GetStat(stats.AttackPower)-start, share()-startShare; math.Abs(got-want) > 1e-6 {
				t.Errorf("%s on the hunter gives the pet %+.1f attack power, want %+.1f", buff.FlatString(), got, want)
			}
			hunter.AddStatsDynamic(sim, buff.Invert())
			if got := hunter.pet.GetStat(stats.AttackPower); math.Abs(got-start) > 1e-6 {
				t.Errorf("pet attack power %.1f after %s ended, want %.1f", got, buff.FlatString(), start)
			}
		}
	}
}

// Resourcefulness procs on melee and ranged ability crits only, not on auto attacks or traps.
func TestResourcefulnessAbilityCritsOnly(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{talents: patchTestTalents, gear: true})
	if hunter.Talents.Resourcefulness == 0 {
		t.Fatal("talent string has no Resourcefulness")
	}
	proc := hunter.GetAura("Resourcefulness Proc")
	buff := hunter.GetAura("Resourcefulness")
	crit := &core.SpellResult{Outcome: core.OutcomeCrit, Target: hunter.CurrentTarget}
	procs := func(spell *core.Spell) bool {
		for i := 0; i < 50; i++ {
			buff.Deactivate(sim)
			proc.OnSpellHitDealt(proc, sim, spell, crit)
			if buff.IsActive() {
				return true
			}
		}
		return false
	}
	if procs(hunter.AutoAttacks.RangedAuto()) || procs(hunter.AutoAttacks.MHAuto()) {
		t.Error("an auto attack crit procs Resourcefulness")
	}
	if !procs(hunter.RaptorStrikeHit) || !procs(hunter.MongooseBite) {
		t.Error("a Raptor Strike or Mongoose Bite crit never procs Resourcefulness")
	}
}

// Resourcefulness and Rapid Recuperation's sting buff let spirit regen continue while casting: the mana ticks inside the
// five second rule grow as soon as one goes up, go back when it fades, and the extra shows as the buff's own mana row.
func TestRegenWhileCastingBuffs(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{talents: "-000000000002-500230230250222151", gear: true})
	drain := hunter.NewManaMetrics(core.ActionID{SpellID: 1})
	tick := func() float64 {
		hunter.SpendMana(sim, hunter.CurrentMana(), drain)
		hunter.PseudoStats.FiveSecondRuleRefreshTime = sim.CurrentTime + time.Second*5
		hunter.ManaTick(sim)
		return hunter.CurrentMana()
	}
	base := tick()
	for _, buff := range []struct {
		name  string
		regen *core.CastingRegenSource
	}{{"Resourcefulness", hunter.resourcefulRegen}, {"Rapid Recuperation Sting", hunter.rapidRecupRegen}} {
		aura := hunter.GetAura(buff.name)
		if aura == nil || buff.regen == nil {
			t.Fatalf("no %s", buff.name)
		}
		aura.Activate(sim)
		gained := buff.regen.Metrics.ActualGain
		with := tick()
		if with < base+5 {
			t.Errorf("%s up: casting tick %.1f, %.1f without it", buff.name, with, base)
		}
		if row := buff.regen.Metrics.ActualGain - gained; math.Abs(row-(with-base)) > 1e-6 {
			t.Errorf("%s row got %.1f of the %.1f extra", buff.name, row, with-base)
		}
		aura.Deactivate(sim)
		if after := tick(); math.Abs(after-base) > 1e-6 {
			t.Errorf("%s faded: casting tick %.1f, want %.1f", buff.name, after, base)
		}
	}
}

// Immolation Trap is a magic spell: its ticks crit for 1.5x like the live logs (65 -> 98), not 2x.
func TestImmolationTrapCritMultiplier(t *testing.T) {
	_, hunter := newPatchTestHunter(t, patchTestSetup{talents: patchTestTalents})
	if hunter.ImmolationTrap == nil {
		t.Fatal("no Immolation Trap")
	}
	at := hunter.AttackTables[hunter.CurrentTarget.UnitIndex][proto.CastType_CastTypeMainHand]
	if got := hunter.ImmolationTrap.CritMultiplier(at); math.Abs(got-1.5) > 1e-9 {
		t.Errorf("Immolation Trap crits for %.2fx, want 1.5x", got)
	}
}

// Immolation Trap burns for its base plus 0.132 of the higher of melee and ranged AP (live logs: Hawk vs Beast
// traps 350 vs 342; naked with Aspect of the Beast, melee above ranged, 324 vs 317 without).
func TestImmolationTrapAttackPower(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{talents: patchTestTalents, gear: true})
	trap := hunter.ImmolationTrap
	target := hunter.CurrentTarget
	// extra melee AP only (like Aspect of the Beast), enough that melee is the higher one with any shipped gear
	hunter.AddStatsDynamic(sim, stats.Stats{stats.AttackPower: 600})
	dot := trap.Dot(target)
	dot.Apply(sim)
	base := []float64{0, 105, 215, 340, 510, 690}[trap.Rank]
	want := (base + 0.132*max(trap.RangedAttackPower(target, false), trap.MeleeAttackPower(target))) / 5
	if got := dot.SnapshotBaseDamage; math.Abs(got-want) > 1e-6 {
		t.Errorf("Immolation Trap tick %.3f, want %.3f (base %.0f + 0.132 x the higher AP)", got, want, base)
	}
	if trap.MeleeAttackPower(target) <= trap.RangedAttackPower(target, false) {
		t.Errorf("melee AP %.0f is not above ranged AP %.0f, so this cannot tell the two rules apart", trap.MeleeAttackPower(target), trap.RangedAttackPower(target, false))
	}
}

// Arcane Shot (rank 8) hits for its base 217 plus 0.10 of ranged attack power, the target's Hunter's Mark (+110)
// included, read from each hit's logged base damage (the logged attack power leaves the mark out).
func TestArcaneShotScaling(t *testing.T) {
	const aplJSON = `{"type":"TypeAPL","priorityList":[{"action":{"castSpell":{"spellId":{"spellId":14287,"rank":8}}}}]}`
	rot := &proto.APLRotation{}
	if err := protojson.Unmarshal([]byte(aplJSON), rot); err != nil {
		t.Fatalf("rotation: %v", err)
	}
	player := core.WithSpec(&proto.Player{
		Race:               proto.Race_RaceHuman,
		Class:              proto.Class_ClassHunter,
		Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
		TalentsString:      patchTestTalents,
		Consumes:           abConsumes,
		Buffs:              abIndividualBuffs,
		Rotation:           rot,
		DistanceFromTarget: 12,
	}, abOptions)
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, abDebuffs),
		Encounter:  abEncounter(),
		SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 20261020, Debug: true},
	})
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
	}
	hit := regexp.MustCompile(`\{SpellID: 14287\} \[DEBUG\] .*RAP\+Bonus: ([0-9.]+), .*BaseDamage:([0-9.]+),`)
	n := 0
	for _, line := range strings.Split(res.Logs, "\n") {
		m := hit.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		rap, _ := strconv.ParseFloat(m[1], 64)
		base, _ := strconv.ParseFloat(m[2], 64)
		if want := 217 + 0.10*(rap+110); math.Abs(base-want) > 0.11 {
			t.Fatalf("Arcane Shot base damage %.1f at %.1f ranged attack power + 110 from the mark, want %.1f (217 + 0.10 x RAP)", base, rap, want)
		}
		n++
	}
	if n == 0 {
		t.Fatalf("no Arcane Shot hits logged")
	}
}

// Strider Kick and Mongoose Bite's weapon part takes the plain normalized speed whatever the weapon's own: a 3.8
// two-hander 3.3, a 2.8 and a 1.7 one-hander 2.4 (Forever logs, a Corpsemaker test and a two-one-hander session).
func TestNormalizedStrikeSpeed(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{talents: patchTestTalents, gear: true})
	mh := hunter.AutoAttacks.MH()
	saved := *mh
	defer func() { *mh = saved }()
	// a fixed base of 100 and 1,400 attack power: 100 + speed x 100
	mh.BaseDamageMin, mh.BaseDamageMax, mh.AttackPowerPerDPS = 100, 100, core.DefaultAttackPowerPerDPS
	for _, tc := range []struct {
		name             string
		real, normalized float64
		want             float64
	}{
		{"3.8 two-hander", 3.8, 3.3, 430},
		{"2.8 one-hander", 2.8, 2.4, 340},
		{"1.7 one-hander", 1.7, 2.4, 340},
	} {
		mh.SwingSpeed, mh.NormalizedSwingSpeed = tc.real, tc.normalized
		if got := mh.CalculateNormalizedWeaponDamage(sim, 1400); math.Abs(got-tc.want) > 1e-6 {
			t.Errorf("%s: strike weapon damage %.1f, want %.1f", tc.name, got, tc.want)
		}
	}
}

// Serpent Sting ticks grow by 0.030 per point of ranged attack power and not with spell power (in game: rank 4 ticked
// 27 at 174 ranged attack power, 32 at 361 before rounding, unchanged with 20 spell power); the rank 4 base tick is the
// client's 22.
func TestSerpentStingScaling(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{talents: patchTestTalents, gear: true})
	target := hunter.CurrentTarget
	sting := hunter.SerpentSting
	tick := func() float64 {
		dot := sting.Dot(target)
		dot.Apply(sim)
		return dot.SnapshotBaseDamage
	}
	rap := sting.RangedAttackPower(target, false)
	if got, want := tick(), 111+0.030*rap; math.Abs(got-want) > 1e-6 {
		t.Errorf("rank 9 tick %.3f at %.0f ranged attack power, want %.3f (111 + 0.030 x RAP)", got, rap, want)
	}
	hunter.AddStatsDynamic(sim, stats.Stats{stats.SpellPower: 200, stats.SpellDamage: 200, stats.NaturePower: 200})
	if got, want := tick(), 111+0.030*rap; math.Abs(got-want) > 1e-6 {
		t.Errorf("rank 9 tick %.3f with 200 spell power, want %.3f (no spell power share)", got, want)
	}
	if c := sting.Dot(target).BonusCoefficient; c != 0 {
		t.Errorf("Serpent Sting spell power share %.3f, want 0", c)
	}
	if low, high := serpentStingTick(4, 174), serpentStingTick(4, 361); math.Floor(low) != 27 || math.Floor(high) != 32 {
		t.Errorf("rank 4 ticks %.2f at 174 and %.2f at 361 ranged attack power, the game had 27 and 32 before rounding", low, high)
	}
}

// A hunter with the Trueshot Aura talent has the aura (+100 ranged attack power) without the raid buff, and the raid buff
// on top does not add a second one.
func TestTrueshotAuraFromTalent(t *testing.T) {
	rap := func(talents string, raid bool) float64 {
		_, hunter := newPatchTestHunter(t, patchTestSetup{talents: talents, gear: true, buffs: &proto.RaidBuffs{TrueshotAura: raid}})
		return hunter.GetStat(stats.RangedAttackPower)
	}
	with, without := "-0053150011-500230230250222151", "-0053150010-500230230250222151"
	if d := rap(with, false) - rap(without, false); math.Abs(d-100) > 1e-6 {
		t.Errorf("talented Trueshot Aura without the raid buff adds %.1f ranged attack power, want 100", d)
	}
	if d := rap(with, true) - rap(without, false); math.Abs(d-100) > 1e-6 {
		t.Errorf("talented Trueshot Aura with the raid buff adds %.1f ranged attack power, want 100 (one aura)", d)
	}
	if d := rap(without, true) - rap(without, false); math.Abs(d-100) > 1e-6 {
		t.Errorf("the raid buff alone adds %.1f ranged attack power, want 100", d)
	}
}

// Immolation Trap is a real cast (placed on a GCD), so the results count its casts and cast time; without them the
// table cannot show its damage per cast or per execute time.
// The burn lands 1.2 s after the trap is dropped and first ticks 3 s later: a trap dropped at the pull from 3 yd has not
// ticked by 4 s (a throw alone would tick at 3.1 s) and has by 4.5 s.
func TestImmolationTrapArming(t *testing.T) {
	const aplJSON = `{"type":"TypeAPL","priorityList":[{"action":{"castSpell":{"spellId":{"spellId":14305}}}}]}`
	trapDamage := func(seconds float64) float64 {
		rot := &proto.APLRotation{}
		if err := protojson.Unmarshal([]byte(aplJSON), rot); err != nil {
			t.Fatalf("rotation: %v", err)
		}
		player := core.WithSpec(&proto.Player{
			Race:               proto.Race_RaceHuman,
			Class:              proto.Class_ClassHunter,
			Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
			TalentsString:      patchTestTalents,
			Consumes:           abConsumes,
			Buffs:              abIndividualBuffs,
			Rotation:           rot,
			DistanceFromTarget: 3,
		}, abOptions)
		enc := abEncounter()
		enc.Duration, enc.DurationVariation = seconds, 0
		res := core.RunRaidSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, abDebuffs),
			Encounter:  enc,
			SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 20261020},
		})
		if res.Error != nil {
			t.Fatalf("sim error: %s", res.Error.Message)
		}
		var damage float64
		for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
			if a.Id.GetSpellId() == 14305 {
				for _, tm := range a.Targets {
					damage += tm.Damage
				}
			}
		}
		return damage
	}
	if d := trapDamage(4); d != 0 {
		t.Errorf("4 s fight: trap did %.0f damage; want none (burn lands at 1.2 s, first tick at 4.2 s)", d)
	}
	if d := trapDamage(4.5); d == 0 {
		t.Error("4.5 s fight: trap did no damage; want its first tick at 4.2 s")
	}
}

func TestImmolationTrapCastMetrics(t *testing.T) {
	player := core.WithSpec(&proto.Player{
		Race:               proto.Race_RaceHuman,
		Class:              proto.Class_ClassHunter,
		Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
		TalentsString:      patchTestTalents,
		Consumes:           abConsumes,
		Buffs:              abIndividualBuffs,
		Rotation:           core.GetAplRotation("../../ui/hunter/apls", "arc_weave_2h").Rotation,
		DistanceFromTarget: 3,
	}, abOptions)
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, abDebuffs),
		Encounter:  abEncounter(),
		SimOptions: &proto.SimOptions{Iterations: 20, RandomSeed: 20261012},
	})
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
	}
	for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
		if a.Id.GetSpellId() != 14305 {
			continue
		}
		var casts int32
		var castTime float64
		for _, tm := range a.Targets {
			casts += tm.Casts
			castTime += tm.CastTimeMs
		}
		if a.IsPassive || casts == 0 || castTime == 0 {
			t.Errorf("Immolation Trap metrics: passive %v, casts %d, cast time %.0f ms; want a counted cast", a.IsPassive, casts, castTime)
		}
		return
	}
	t.Fatal("no Immolation Trap in the results")
}

// Explosive Trap's first hit gets the same T.N.T. share of the higher of melee and ranged AP as Immolation Trap's burn
// (428717 names the first hit). Its burn gets it too only with the "burn scales" setting, spread over its 10 ticks.
// Same fight, same seed, +500 melee and ranged AP (so the higher one rises by 500 whichever it is): each first hit rises
// by about 0.132 x 500; each tick by a tenth of it or nothing.
func TestExplosiveTrapAttackPower(t *testing.T) {
	// first-hit damage per trap and damage per burn tick
	run := func(bonusAP float64, burnAP bool) (float64, float64) {
		rotation := &proto.APLRotation{}
		if err := protojson.Unmarshal([]byte(`{"type":"TypeAPL","priorityList":[{"action":{"castSpell":{"spellId":{"spellId":14317}}}}]}`), rotation); err != nil {
			t.Fatal(err)
		}
		bonus := make([]float64, 44)
		bonus[17] = bonusAP // melee attack power
		bonus[27] = bonusAP // ranged attack power
		player := core.WithSpec(&proto.Player{
			Race:               proto.Race_RaceHuman,
			Class:              proto.Class_ClassHunter,
			Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_light_weave").GearSet,
			TalentsString:      patchTestTalents,
			BonusStats:         &proto.UnitStats{Stats: bonus},
			Rotation:           rotation,
			DistanceFromTarget: 3,
		}, abOptions)
		player.GetHunter().Options.ExplosiveTrapBurnAp = burnAP
		res := core.RunRaidSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  abEncounter(),
			SimOptions: &proto.SimOptions{Iterations: 50, RandomSeed: 20261012},
		})
		if res.Error != nil {
			t.Fatalf("sim error: %s", res.Error.Message)
		}
		for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
			if a.Id.GetSpellId() != 14317 {
				continue
			}
			var casts, ticks int32
			var damage, tickDamage float64
			for _, tm := range a.Targets {
				casts += tm.Casts
				ticks += tm.Ticks + tm.CritTicks
				damage += tm.Damage
				tickDamage += tm.TickDamage // crit ticks included
			}
			if casts == 0 || ticks == 0 {
				t.Fatalf("Explosive Trap: %d casts, %d ticks", casts, ticks)
			}
			return (damage - tickDamage) / float64(casts), tickDamage / float64(ticks)
		}
		t.Fatal("no Explosive Trap in the results")
		return 0, 0
	}
	// the default: the first hit scales, the burn does not
	hit0, tick0 := run(0, false)
	hit1, tick1 := run(500, false)
	// generous bounds: talents raise it, misses and partial resists lower the first hit
	if gain, want := hit1-hit0, 0.132*500; gain < want*0.6 || gain > want*1.6 {
		t.Errorf("+500 AP adds %.1f to Explosive Trap's first hit, want about %.0f", gain, want)
	}
	if gain := tick1 - tick0; math.Abs(gain) > 1e-6 {
		t.Errorf("by default +500 AP adds %.2f to each Explosive Trap burn tick, want 0", gain)
	}

	// the "burn scales" setting: the hit and each tick scale
	hitA, tickA := run(0, true)
	hitB, tickB := run(500, true)
	if gain, want := hitB-hitA, 0.132*500; gain < want*0.6 || gain > want*1.6 {
		t.Errorf("burn scales: +500 AP adds %.1f to the first hit, want about %.0f", gain, want)
	}
	// a tenth of the first hit's gain: both get the same trap talents, so this holds whatever they add
	if gain, want := tickB-tickA, (hitB-hitA)/10; gain < want*0.6 || gain > want*1.6 {
		t.Errorf("burn scales: +500 AP adds %.2f to each burn tick, want about %.1f", gain, want)
	}
}

// Berserking is spell 20554 in Forever (26297 is not in its client data, so the site cannot name it or show its
// icon); a rotation naming 20554 must find it to cast.
func TestBerserkingForeverSpell(t *testing.T) {
	_, troll := newPatchTestHunter(t, patchTestSetup{race: proto.Race_RaceTroll})
	if troll.GetSpell(core.ActionID{SpellID: 20554}) == nil {
		t.Error("a Troll has no Berserking spell 20554")
	}
	if troll.GetSpell(core.ActionID{SpellID: 26297}) != nil {
		t.Error("Berserking still uses the Classic spell id 26297")
	}
}

// Explosive Trap: 30 s base cooldown in the Forever data (category 411, shared with Immolation Trap), 18 s with
// Survivalist's Discipline (in-game tooltip); its burn can crit (Periodic Can Crit).
func TestExplosiveTrapForeverData(t *testing.T) {
	_, hunter := newPatchTestHunter(t, patchTestSetup{talents: patchTestTalents})
	if hunter.Talents.SurvivalistDiscipline != 2 {
		t.Fatalf("test talents have Survivalist's Discipline %d/2", hunter.Talents.SurvivalistDiscipline)
	}
	if got := hunter.ExplosiveTrap.CD.Duration; got != 18*time.Second {
		t.Errorf("Explosive Trap cooldown %v, want 18s (30 s base, -40%%)", got)
	}
	if hunter.ExplosiveTrap.CD.Timer != hunter.ImmolationTrap.CD.Timer {
		t.Error("Explosive Trap and Immolation Trap do not share a cooldown")
	}

	rotation := &proto.APLRotation{}
	if err := protojson.Unmarshal([]byte(`{"type":"TypeAPL","priorityList":[{"action":{"castSpell":{"spellId":{"spellId":14317}}}}]}`), rotation); err != nil {
		t.Fatal(err)
	}
	player := core.WithSpec(&proto.Player{
		Race:               proto.Race_RaceHuman,
		Class:              proto.Class_ClassHunter,
		Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
		TalentsString:      patchTestTalents,
		Rotation:           rotation,
		DistanceFromTarget: 3,
	}, abOptions)
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  abEncounter(),
		SimOptions: &proto.SimOptions{Iterations: 50, RandomSeed: 20261012},
	})
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
	}
	for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
		if a.Id.GetSpellId() != 14317 {
			continue
		}
		var ticks, critTicks int32
		for _, tm := range a.Targets {
			ticks += tm.Ticks
			critTicks += tm.CritTicks
		}
		if critTicks == 0 {
			t.Errorf("Explosive Trap burn: %d ticks, none crit", ticks+critTicks)
		}
		return
	}
	t.Fatal("no Explosive Trap in the results")
}

// A saved "Beta pet speed bug" setting changes nothing: pets hit for their own speed, as in game.
func TestPetSpeedBugFixed(t *testing.T) {
	_, clean := newPatchTestHunter(t, patchTestSetup{})
	_, old := newPatchTestHunter(t, patchTestSetup{petBug: true})
	a, b := clean.pet.AutoAttacks.MH(), old.pet.AutoAttacks.MH()
	if a.BaseDamageMin != b.BaseDamageMin || a.BaseDamageMax != b.BaseDamageMax || a.AttackPowerPerDPS != b.AttackPowerPerDPS {
		t.Errorf("pet weapon with the old setting %.2f-%.2f (%.2f AP/DPS), without %.2f-%.2f (%.2f); want the same",
			b.BaseDamageMin, b.BaseDamageMax, b.AttackPowerPerDPS, a.BaseDamageMin, a.BaseDamageMax, a.AttackPowerPerDPS)
	}
}

// Dust Cloud (tallstrider, -505 armor) does not stack with Faerie Fire (tested in game: one -505 debuff); on its own
// it still takes the armor off.
func TestDustCloudSharesFaerieFireGroup(t *testing.T) {
	for _, withFF := range []bool{true, false} {
		sim, hunter := newPatchTestHunter(t, patchTestSetup{petType: proto.Hunter_Options_Tallstrider, debuffs: &proto.Debuffs{FaerieFire: withFF}})
		target := hunter.CurrentTarget
		if ff := target.GetAura("Faerie Fire"); withFF && !ff.IsActive() {
			ff.Activate(sim)
		}
		before := target.GetStat(stats.Armor)
		target.GetAura("Dust Cloud").Activate(sim)
		want := 505.0
		if withFF {
			want = 0
		}
		if got := before - target.GetStat(stats.Armor); math.Abs(got-want) > 1e-6 {
			t.Errorf("Faerie Fire up %v: Dust Cloud removes %.0f armor, want %.0f", withFF, got, want)
		}
	}
}

// Curse of Recklessness, Faerie Fire and Dust Cloud leave one -505 armor debuff between them in Forever (tested in
// game); Curse of Recklessness alone still takes 505 off.
func TestCurseOfRecklessnessSharesFaerieFireGroup(t *testing.T) {
	armorOff := func(debuffs *proto.Debuffs) float64 {
		_, bare := newPatchTestHunter(t, patchTestSetup{})
		sim, hunter := newPatchTestHunter(t, patchTestSetup{debuffs: debuffs})
		target := hunter.CurrentTarget
		for _, label := range []string{"Faerie Fire", "Curse of Recklessness"} {
			if aura := target.GetAura(label); aura != nil && !aura.IsActive() {
				aura.Activate(sim)
			}
		}
		return bare.CurrentTarget.GetStat(stats.Armor) - target.GetStat(stats.Armor)
	}
	for name, tc := range map[string]struct {
		debuffs *proto.Debuffs
		want    float64
	}{
		"Curse of Recklessness":                 {&proto.Debuffs{CurseOfRecklessness: true}, 505},
		"Faerie Fire":                           {&proto.Debuffs{FaerieFire: true}, 505},
		"Curse of Recklessness and Faerie Fire": {&proto.Debuffs{CurseOfRecklessness: true, FaerieFire: true}, 505},
	} {
		if got := armorOff(tc.debuffs); math.Abs(got-tc.want) > 1e-6 {
			t.Errorf("%s: target armor -%.0f, want -%.0f", name, got, tc.want)
		}
	}
}

// Shadow Oil's bolt is spell 350026 in Forever (the Classic 1382 is not in its data, so the site cannot name it):
// 52 +-7.7% (48 to 56; the client's variance field is the whole spread) plus 0.56 of spell power.
func TestShadowOilForeverData(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, mhImbue: proto.WeaponImbue_ShadowOil})
	if hunter.GetSpell(core.ActionID{SpellID: 1382}) != nil {
		t.Error("Shadow Oil still uses the Classic spell id 1382")
	}
	bolt := hunter.GetSpell(core.ActionID{SpellID: 350026})
	if bolt == nil {
		t.Fatal("no Shadow Oil bolt 350026")
	}
	if bolt.BonusCoefficient != 0.56 {
		t.Errorf("Shadow Oil bolt spell power share %.2f, want 0.56", bolt.BonusCoefficient)
	}
	// the 2H gear has no spell power: an unresisted, non-crit bolt is the base roll times the talents
	target := sim.Encounter.TargetUnits[0]
	m := &bolt.SpellMetrics[target.UnitIndex]
	lo, hi := math.Inf(1), 0.0
	for i := 0; i < 2000; i++ {
		hits, resisted, crits, damage := m.Hits, m.ResistedHits, m.Crits, m.TotalDamage
		bolt.Cast(sim, target)
		if m.Hits == hits+1 && m.ResistedHits == resisted && m.Crits == crits {
			d := (m.TotalDamage - damage) / hunter.PseudoStats.DamageDealtMultiplier // talents raise all damage
			lo, hi = math.Min(lo, d), math.Max(hi, d)
		}
	}
	if lo < 48-1e-6 || hi > 56+1e-6 || lo > 49 || hi < 55 {
		t.Errorf("Shadow Oil bolt hits %.1f to %.1f, want 48 to 56", lo, hi)
	}
}

// Frost Oil's bolt is spell 350025 in Forever (the Classic 1191 is not in its data, so the site showed a nameless row):
// 33 +-6% (31 to 35; the client's variance field is the whole spread) plus 0.269 of spell power.
func TestFrostOilForeverData(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, mhImbue: proto.WeaponImbue_FrostOil})
	if hunter.GetSpell(core.ActionID{SpellID: 1191}) != nil {
		t.Error("Frost Oil still uses the Classic spell id 1191")
	}
	bolt := hunter.GetSpell(core.ActionID{SpellID: 350025})
	if bolt == nil {
		t.Fatal("no Frost Oil bolt 350025")
	}
	if bolt.BonusCoefficient != 0.269 {
		t.Errorf("Frost Oil bolt spell power share %.3f, want 0.269", bolt.BonusCoefficient)
	}
	target := sim.Encounter.TargetUnits[0]
	m := &bolt.SpellMetrics[target.UnitIndex]
	lo, hi := math.Inf(1), 0.0
	for i := 0; i < 2000; i++ {
		hits, resisted, crits, damage := m.Hits, m.ResistedHits, m.Crits, m.TotalDamage
		bolt.Cast(sim, target)
		if m.Hits == hits+1 && m.ResistedHits == resisted && m.Crits == crits {
			d := (m.TotalDamage - damage) / hunter.PseudoStats.DamageDealtMultiplier
			lo, hi = math.Min(lo, d), math.Max(hi, d)
		}
	}
	if lo < 31-1e-6 || hi > 35+1e-6 || lo > 31.5 || hi < 34.5 {
		t.Errorf("Frost Oil bolt hits %.1f to %.1f, want 31 to 35", lo, hi)
	}
}

// Aspect of the Cheetah: +30% run speed (+3% a Pathfinding point), it shares the aspect slot with Aspect of the Beast,
// and it does not stack with Strider Kick (the stronger bonus counts).
func TestAspectOfTheCheetah(t *testing.T) {
	speedWith := func(talents string, kick bool) float64 {
		sim, hunter := newPatchTestHunter(t, patchTestSetup{talents: talents})
		cheetah := hunter.GetSpell(core.ActionID{SpellID: 5118})
		if cheetah == nil {
			t.Fatal("no Aspect of the Cheetah spell 5118")
		}
		cheetah.Cast(sim, hunter.CurrentTarget)
		if kick {
			hunter.GetAura("Strider Kick Speed").Activate(sim)
		}
		return hunter.MovementHandler.MoveSpeed
	}
	if got := speedWith("50230005001-005005-5002302300500221", false); math.Abs(got-7*1.3) > 1e-9 {
		t.Errorf("Cheetah runs %.3f yd/s, want %.3f", got, 7*1.3)
	}
	if got := speedWith("50232005001-005005-5002302300500221", false); math.Abs(got-7*1.36) > 1e-9 {
		t.Errorf("Cheetah with Pathfinding 2 runs %.3f yd/s, want %.3f", got, 7*1.36)
	}
	if got := speedWith(patchTestTalents, true); math.Abs(got-7*1.3) > 1e-9 {
		t.Errorf("Cheetah + Strider Kick runs %.3f yd/s, want %.3f (no stacking)", got, 7*1.3)
	}

	sim, hunter := newPatchTestHunter(t, patchTestSetup{})
	hunter.GetSpell(core.ActionID{SpellID: 5118}).Cast(sim, hunter.CurrentTarget)
	hunter.GetAura("Aspect of the Beast").Activate(sim)
	if hunter.GetAura("Aspect of the Cheetah").IsActive() || math.Abs(hunter.MovementHandler.MoveSpeed-7) > 1e-9 {
		t.Errorf("Aspect of the Beast did not replace Cheetah (speed %.3f)", hunter.MovementHandler.MoveSpeed)
	}
}

// Raptor Strike's 6 s cooldown starts when the queued hit lands; rotations read it through the castable (queue) spell,
// which must report the running cooldown, not 0.
func TestRaptorStrikeCooldownReadable(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 3})
	if got := hunter.RaptorStrike.TimeToReady(sim); got != 0 {
		t.Fatalf("Raptor Strike not ready at the start (%v)", got)
	}
	hunter.RaptorStrikeHit.Cast(sim, hunter.CurrentTarget)
	if got := hunter.RaptorStrike.TimeToReady(sim); got != 6*time.Second {
		t.Errorf("after a Raptor Strike hit the rotation sees %v until it is ready, want 6s", got)
	}
}

// Raptor Strike queues from range in Forever and stays queued until the next melee swing, which it replaces.
func TestRaptorStrikeQueuesFromRange(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 12})
	if !hunter.RaptorStrike.CanCast(sim, hunter.CurrentTarget) {
		t.Fatal("Raptor Strike cannot be queued at 12 yd")
	}
	hunter.RaptorStrike.Cast(sim, hunter.CurrentTarget)
	if hunter.curQueueAura == nil || !hunter.curQueueAura.IsActive() {
		t.Fatal("Raptor Strike is not queued after casting it at range")
	}
	if got := hunter.TryRaptorStrike(sim, hunter.AutoAttacks.MHAuto()); got != hunter.RaptorStrikeHit {
		t.Errorf("the next melee swing is %v, want the queued Raptor Strike", got.ActionID)
	}
}

// A melee swing restarts the Auto Shot timer from that moment when the beta bug is on: the next Auto Shot is one full
// cycle after it (confirmed in game with the swing timer events), a queued Raptor Strike taking the swing as well. One
// restart, from the swing (at 3 yd no Auto Shot fires, so the timer only moves by the restart).
func TestMeleeSwingRestartsAutoShot(t *testing.T) {
	raptorSwing := func(meleeResets bool) (before, hitAt, after time.Duration) {
		sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 3, meleeResets: meleeResets})
		sim.PrePull()
		for guard := 0; sim.CurrentTime < time.Second && guard < 100000; guard++ {
			sim.Step()
		}
		before = hunter.AutoAttacks.NextRangedAttackAt()
		if !hunter.RaptorStrike.Cast(sim, hunter.CurrentTarget) {
			t.Fatal("Raptor Strike did not queue")
		}
		casts := func() int32 { return hunter.RaptorStrikeHit.SpellMetrics[hunter.CurrentTarget.UnitIndex].Casts }
		for guard := 0; casts() == 0 && guard < 100000; guard++ {
			sim.Step()
		}
		return before, sim.CurrentTime, hunter.AutoAttacks.NextRangedAttackAt()
	}
	_, hitAt, after := raptorSwing(true)
	_, h := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 3, meleeResets: true})
	if want := hitAt + h.AutoAttacks.RangedSwingSpeed(); after != want {
		t.Errorf("bug on: after the Raptor Strike swing at %v the next Auto Shot is at %v, want %v (one cycle after it)", hitAt, after, want)
	}
	// a white swing restarts it once too, from the swing (it used to restart again as the hit landed)
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 3, meleeResets: true})
	sim.PrePull()
	white := func() int32 { return hunter.AutoAttacks.MHAuto().SpellMetrics[hunter.CurrentTarget.UnitIndex].Casts }
	for guard := 0; white() == 0 && guard < 100000; guard++ {
		sim.Step()
	}
	swingAt := sim.CurrentTime
	for guard := 0; sim.CurrentTime < swingAt+50*time.Millisecond && guard < 100000; guard++ {
		sim.Step()
	}
	if got, want := hunter.AutoAttacks.NextRangedAttackAt(), swingAt+hunter.AutoAttacks.RangedSwingSpeed(); got != want {
		t.Errorf("bug on: after the white swing at %v the next Auto Shot is at %v, want %v (one cycle after the swing)", swingAt, got, want)
	}

	// by default (the bug fixed, forever-bugs #14) the swing leaves the Auto Shot timer where it was
	before, _, after := raptorSwing(false)
	if after != before {
		t.Errorf("bug off: the Raptor Strike swing moved the next Auto Shot from %v to %v", before, after)
	}
}

// The Auto Shot spam bug switch takes the cut off every Auto Shot timer a shot starts; off, the full cycle stays.
func TestAutoShotSpamBugCutsTheCycle(t *testing.T) {
	for _, spam := range []bool{false, true} {
		sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 30, autoSpam: spam})
		sim.PrePull()
		shots := 0
		for i := 0; i < 100000 && shots < 4; i++ {
			prev := hunter.AutoAttacks.NextRangedAttackAt()
			sim.Step()
			next := hunter.AutoAttacks.NextRangedAttackAt()
			if next == prev || sim.CurrentTime != prev {
				continue
			}
			shots++
			want := hunter.AutoAttacks.RangedSwingSpeed()
			if spam {
				want -= AutoShotSpamCut
			}
			if got := next - sim.CurrentTime; got != want {
				t.Errorf("spam bug %v: shot %d starts a %v timer, want %v", spam, shots, got, want)
			}
		}
		if shots < 4 {
			t.Fatalf("spam bug %v: only %d Auto Shots fired", spam, shots)
		}
	}
}

// swingRetry: the first half-second swing retry from origin that the server sees you in range for (seen).
func swingRetry(origin, seen time.Duration) time.Duration {
	const beat = 500 * time.Millisecond
	if seen <= origin {
		return origin
	}
	return origin + (seen-origin+beat-1)/beat*beat
}

// Stepping into melee turns the melee attack on as the walk starts: the main hand, ready long ago, swings on the first
// retry after the server sees you at 5 yd, the off hand half its swing time after the walk began (or on the first
// retry after that once in range).
func TestOffHandWaitsHalfItsSwingOnStepIn(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gearSet: "surv_dw", distance: 9})
	if !hunter.AutoAttacks.IsDualWielding {
		t.Fatal("not dual wielding")
	}
	sim.PrePull()
	for i := 0; i < 100000 && sim.CurrentTime < 5*time.Second; i++ {
		sim.Step()
	}
	walkStart := sim.CurrentTime
	hunter.MoveTo(3, sim)
	walkUntilMeleeRange(sim, hunter)
	seen := sim.CurrentTime + 50*time.Millisecond
	if got, want := hunter.AutoAttacks.OffhandSwingAt(), swingRetry(walkStart+hunter.AutoAttacks.OffhandSwingSpeed()/2, seen); got != want {
		t.Errorf("off hand swings at %v, want %v (walk start + half its swing time, or the retry after)", got, want)
	}
	if got, want := hunter.AutoAttacks.MainhandSwingAt(), swingRetry(walkStart, seen); got != want {
		t.Errorf("main hand swings at %v, want %v (the retry after reaching 5 yd)", got, want)
	}

	// out and straight back in without a shot: the attack stayed on, so a ready off hand retries from its timer's end
	for i := 0; i < 100000 && hunter.MovementHandler.Moving; i++ {
		sim.Step()
	}
	ohSwing := hunter.AutoAttacks.OffhandSwingAt()
	for i := 0; i < 100000 && sim.CurrentTime <= ohSwing; i++ {
		sim.Step()
	}
	hunter.MoveTo(9, sim)
	for i := 0; i < 100000 && hunter.MovementHandler.Moving; i++ {
		sim.Step()
	}
	ohReady := sim.CurrentTime
	hunter.AutoAttacks.SetOffhandSwingAt(ohReady)
	hunter.MoveTo(3, sim)
	walkUntilMeleeRange(sim, hunter)
	if got, want := hunter.AutoAttacks.OffhandSwingAt(), swingRetry(ohReady, sim.CurrentTime+50*time.Millisecond); got != want {
		t.Errorf("no shot between: off hand swings at %v, want %v (the retry from its timer's end)", got, want)
	}
}

// The Start Attack action and the walk-in switch (manualStartAttack). Off, the walk in turns the melee attack on as it
// starts (or at an earlier Start Attack); on, only Start Attack does, else reaching melee turns it on and the main hand
// swings as soon as the server sees it there, the off hand half its swing time later.
func TestStartAttack(t *testing.T) {
	type run struct{ press, walkStart, seen, mh, oh, ohSpeed time.Duration }
	step := func(t *testing.T, manual, press bool) run {
		sim, hunter := newPatchTestHunter(t, patchTestSetup{gearSet: "surv_dw", distance: 9, manualStart: manual})
		if !hunter.AutoAttacks.IsDualWielding {
			t.Fatal("not dual wielding")
		}
		sim.PrePull()
		for i := 0; i < 100000 && sim.CurrentTime < 5*time.Second; i++ {
			sim.Step()
		}
		// no swing since the pull: the attack is off, as after an Auto Shot
		if !hunter.AutoAttacks.MeleeAttackOff() {
			t.Fatal("test setup: the melee attack is on")
		}
		var r run
		if press {
			r.press = sim.CurrentTime
			hunter.AutoAttacks.StartAttack(sim)
			if hunter.AutoAttacks.MeleeAttackOff() {
				t.Fatal("Start Attack left the melee attack off")
			}
		}
		for i := 0; i < 100000 && sim.CurrentTime < r.press+400*time.Millisecond; i++ {
			sim.Step()
		}
		r.walkStart = sim.CurrentTime
		hunter.MoveTo(3, sim)
		walkUntilMeleeRange(sim, hunter)
		r.seen = sim.CurrentTime + 50*time.Millisecond
		r.mh, r.oh, r.ohSpeed = hunter.AutoAttacks.MainhandSwingAt(), hunter.AutoAttacks.OffhandSwingAt(), hunter.AutoAttacks.OffhandSwingSpeed()
		return r
	}
	for _, c := range []struct {
		name          string
		manual, press bool
		mh, oh        func(r run) time.Duration
	}{
		{"walk-in switch off, no Start Attack: on from the walk's start", false, false,
			func(r run) time.Duration { return swingRetry(r.walkStart, r.seen) },
			func(r run) time.Duration { return swingRetry(r.walkStart+r.ohSpeed/2, r.seen) }},
		{"walk-in switch off, Start Attack before the walk: on from the press", false, true,
			func(r run) time.Duration { return swingRetry(r.press, r.seen) },
			func(r run) time.Duration { return swingRetry(r.press+r.ohSpeed/2, r.seen) }},
		{"walk-in switch on, no Start Attack: on at arrival, no retry", true, false,
			func(r run) time.Duration { return r.seen },
			func(r run) time.Duration { return r.seen - 50*time.Millisecond + r.ohSpeed/2 }},
		{"walk-in switch on, Start Attack before the walk: on from the press", true, true,
			func(r run) time.Duration { return swingRetry(r.press, r.seen) },
			func(r run) time.Duration { return swingRetry(r.press+r.ohSpeed/2, r.seen) }},
	} {
		t.Run(c.name, func(t *testing.T) {
			r := step(t, c.manual, c.press)
			if want := c.mh(r); r.mh != want {
				t.Errorf("main hand swings at %v, want %v (walk %v, seen in range %v)", r.mh, want, r.walkStart, r.seen)
			}
			if want := c.oh(r); r.oh != want {
				t.Errorf("off hand swings at %v, want %v (walk %v, seen in range %v)", r.oh, want, r.walkStart, r.seen)
			}
		})
	}
}

// The Start Attack action in a rotation: pressed while the melee attack is off, never again until an Auto Shot turns
// the attack off.
func TestStartAttackActionOncePerAutoShot(t *testing.T) {
	const aplJSON = `{"type":"TypeAPL","recheckAfterAutoShot":true,"manualStartAttack":true,"priorityList":[{"action":{"startAttack":{}}}]}`
	rot := &proto.APLRotation{}
	if err := protojson.Unmarshal([]byte(aplJSON), rot); err != nil {
		t.Fatalf("rotation: %v", err)
	}
	player := core.WithSpec(&proto.Player{
		Race:               proto.Race_RaceHuman,
		Class:              proto.Class_ClassHunter,
		Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
		TalentsString:      patchTestTalents,
		Consumes:           abConsumes,
		Buffs:              abIndividualBuffs,
		Rotation:           rot,
		DistanceFromTarget: 12,
	}, abOptions)
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, abDebuffs),
		Encounter:  abEncounter(),
		SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 20261024, DebugFirstIteration: true},
	})
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
	}
	presses, shots, pressedSinceShot := 0, 0, false
	for _, line := range strings.Split(res.Logs, "\n") {
		switch {
		case strings.Contains(line, "Melee attack on (Start Attack)"):
			if pressedSinceShot {
				t.Fatalf("Start Attack pressed twice without an Auto Shot between: %s", line)
			}
			presses, pressedSinceShot = presses+1, true
		case strings.Contains(line, "Completed cast {OtherID: 8}"):
			shots, pressedSinceShot = shots+1, false
		}
	}
	if shots < 10 || presses < shots {
		t.Errorf("%d Start Attack presses for %d Auto Shots, want one after each", presses, shots)
	}
}

// A walk into melee swings from 5 yd on, mid-walk when the retry comes before the walk ends: the main hand, ready long
// ago, swings on the first retry after the server sees it at 5 yd, counted from the walk's start.
func TestMeleeSwingsStartAtFiveYards(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gearSet: "surv_dw", distance: 9})
	sim.PrePull()
	for i := 0; i < 100000 && sim.CurrentTime < 5*time.Second; i++ {
		sim.Step()
	}
	walkStart := sim.CurrentTime
	hunter.MoveTo(1, sim)
	walkUntilMeleeRange(sim, hunter)
	if !hunter.MovementHandler.Moving {
		t.Fatal("the walk ended at 5 yd")
	}
	swing := hunter.AutoAttacks.MainhandSwingAt()
	if want := swingRetry(walkStart, sim.CurrentTime+50*time.Millisecond); swing != want {
		t.Errorf("main hand swings at %v, want the retry at %v", swing, want)
	}
	for i := 0; i < 100000 && hunter.MovementHandler.Moving; i++ {
		sim.Step()
	}
	if arrival := sim.CurrentTime; swing >= arrival {
		t.Fatalf("test setup: the retry at %v comes after the walk ends at %v", swing, arrival)
	}
	if got, want := hunter.AutoAttacks.MainhandSwingAt(), swing+hunter.AutoAttacks.MainhandSwingSpeed(); got != want {
		t.Errorf("next main hand swing at %v, want %v (one swing mid-walk at %v)", got, want, swing)
	}
}

// A walk out of melee keeps swinging while still within 5 yd: a main-hand swing due in the first yard lands.
func TestMeleeSwingsWhileWalkingOut(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 4})
	hunter.MoveDelay = 0
	sim.PrePull()
	for i := 0; i < 100000; i++ {
		sim.Step()
		left := hunter.AutoAttacks.MainhandSwingAt() - sim.CurrentTime
		if sim.CurrentTime > 3*time.Second && left > 0 && left < 100*time.Millisecond {
			break
		}
	}
	due := hunter.AutoAttacks.MainhandSwingAt()
	hunter.MoveTo(11, sim)
	for i := 0; i < 100000 && hunter.DistanceFromTarget <= core.MaxMeleeAttackDistance; i++ {
		sim.Step()
	}
	if got := hunter.AutoAttacks.MainhandSwingAt(); got <= due {
		t.Errorf("the main-hand swing due at %v did not land on the walk out (next swing at %v, out of melee at %v)", due, got, sim.CurrentTime)
	}
}

// walkUntilMeleeRange steps a walk in until it reaches 5 yd.
func walkUntilMeleeRange(sim *core.Simulation, hunter *Hunter) {
	for i := 0; i < 100000 && hunter.DistanceFromTarget > core.MaxMeleeAttackDistance; i++ {
		sim.Step()
	}
}

// A walk moves you a yard at a time, but range checks see where you are mid-yard: walking off the 5 yd edge you are out
// of melee as you start (no Wing Clip or Mongoose Bite until the first yard ends), and walking in from 11 yd you are
// out of shot range as you start (no Arcane Shot).
func TestRangeChecksMidWalk(t *testing.T) {
	midFirstYard := func(sim *core.Simulation, hunter *Hunter, to float64) {
		t.Helper()
		start := sim.CurrentTime
		hunter.MoveTo(to, sim)
		reached := false
		core.StartDelayedAction(sim, core.DelayedActionOptions{DoAt: start + 50*time.Millisecond, OnAction: func(*core.Simulation) { reached = true }})
		for i := 0; i < 100000 && !reached; i++ {
			sim.Step()
		}
		if !hunter.IsMoving() || sim.CurrentTime != start+50*time.Millisecond {
			t.Fatalf("walk to %v yd: not 50 ms into its first yard (moving %v, at %v)", to, hunter.IsMoving(), sim.CurrentTime-start)
		}
	}
	settle := func(sim *core.Simulation, hunter *Hunter) {
		for i := 0; i < 100000 && (hunter.IsMoving() || sim.CurrentTime < 5*time.Second || !hunter.GCD.IsReady(sim)); i++ {
			sim.Step()
		}
	}

	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 5})
	sim.PrePull()
	settle(sim, hunter)
	target := hunter.CurrentTarget
	if !hunter.WingClip.CanCast(sim, target) {
		t.Fatal("Wing Clip not castable standing at 5 yd")
	}
	midFirstYard(sim, hunter, 11)
	if d := hunter.DistanceNow(sim); d <= 5 || d >= 6 {
		t.Errorf("50 ms into the walk out from 5 yd, at %.2f yd, want between 5 and 6", d)
	}
	if hunter.WingClip.CanCast(sim, target) {
		t.Error("Wing Clip castable 50 ms into a walk out from 5 yd")
	}
	settle(sim, hunter)

	if !hunter.ArcaneShot.CanCast(sim, target) {
		t.Fatalf("Arcane Shot not castable standing at 11 yd (at %v yd, moving %v, GCD ready %v, cooldown ready %v, mana %.0f, at %v)", hunter.DistanceNow(sim),
			hunter.IsMoving(), hunter.GCD.IsReady(sim), core.BothTimersReady(hunter.ArcaneShot.CD.Timer, hunter.ArcaneShot.SharedCD.Timer, sim), hunter.CurrentMana(), sim.CurrentTime)
	}
	midFirstYard(sim, hunter, 5)
	if d := hunter.DistanceNow(sim); d >= 11 || d <= 10 {
		t.Errorf("50 ms into the walk in from 11 yd, at %.2f yd, want between 10 and 11", d)
	}
	if hunter.ArcaneShot.CanCast(sim, target) {
		t.Error("Arcane Shot castable 50 ms into a walk in from 11 yd")
	}
}

// The site lets a sharpening stone onto a sharp weapon only and a weightstone onto a blunt one, and clears any other
// pairing; the sim gives a misfit stone nothing, so its numbers match what the site loads.
func TestImbueMustFitWeapon(t *testing.T) {
	const baghNakh = 250618 // fist weapon
	ohMin := func(imbue proto.WeaponImbue) float64 {
		_, hunter := newPatchTestHunter(t, patchTestSetup{gearSet: "surv_dw", ohItem: baghNakh, ohImbue: imbue})
		return hunter.AutoAttacks.OH().BaseDamageMin
	}
	bare := ohMin(proto.WeaponImbue_WeaponImbueUnknown)
	if got := ohMin(proto.WeaponImbue_DenseSharpeningStone); got != bare {
		t.Errorf("sharpening stone on a fist weapon: off-hand damage %v, want %v (no bonus)", got, bare)
	}
	if got := ohMin(proto.WeaponImbue_DenseWeightstone); got != bare+8 {
		t.Errorf("weightstone on a fist weapon: off-hand damage %v, want %v", got, bare+8)
	}
}

// The Run speed setting sets the base speed walks use; 0 keeps the game's 7 yd/s.
func TestRunSpeedSetting(t *testing.T) {
	_, standard := newPatchTestHunter(t, patchTestSetup{})
	if got := standard.MovementHandler.MoveSpeed; got != 7 {
		t.Errorf("default run speed %v, want 7", got)
	}
	_, slow := newPatchTestHunter(t, patchTestSetup{runSpeed: 3.5})
	if got := slow.MovementHandler.MoveSpeed; got != 3.5 {
		t.Errorf("run speed setting 3.5 gives %v", got)
	}
}

// Forever's own food and elixir values, from its client data (build 1.60.1.70450): Grilled Squid gives crit instead of
// Agility, the wolf foods give Agility, and the new foods and elixirs give what their buffs say. Spirit is the elixir's
// 25 times the Human 10%.
func TestForeverFoodAndElixirValues(t *testing.T) {
	_, base := newPatchTestHunter(t, patchTestSetup{})
	type want struct{ agi, str, intel, spi, mp5, crit float64 }
	const spirit = 25 * 1.10
	cases := map[string]struct {
		consumes *proto.Consumes
		want     want
	}{
		"Grilled Squid":         {&proto.Consumes{Food: proto.Food_FoodGrilledSquid}, want{crit: 1}},
		"Hot Wolf Ribs":         {&proto.Consumes{Food: proto.Food_FoodHotWolfRibs}, want{agi: 10}},
		"Tender Wolf Steak":     {&proto.Consumes{Food: proto.Food_FoodTenderWolfSteak}, want{agi: 15}},
		"Flank au Poivre":       {&proto.Consumes{Food: proto.Food_FoodFlankAuPoivre}, want{agi: 20}},
		"Elixir of Cunning":     {&proto.Consumes{AgilityElixir: proto.AgilityElixir_ElixirOfCunning}, want{agi: 25, intel: 25}},
		"Elixir of Ferocity":    {&proto.Consumes{StrengthBuff: proto.StrengthBuff_FerocityElixir}, want{agi: 18, str: 18}},
		"Elixir of the Grizzly": {&proto.Consumes{StrengthBuff: proto.StrengthBuff_ElixirOfTheGrizzly}, want{str: 25, crit: 2}},
		"Elixir of the Owl":     {&proto.Consumes{IntellectElixir: proto.ManaRegenElixir_ElixirOfTheOwl}, want{intel: 25, crit: 2}},
		"Soaring Pamplona":      {&proto.Consumes{Food: proto.Food_FoodSoaringPamplona}, want{intel: 20}},
		"Greater Mageblood":     {&proto.Consumes{ManaRegenElixir: proto.ManaRegenElixir_GreaterMagebloodElixir}, want{mp5: 20}},
		"Elixir of Sages":       {&proto.Consumes{SpiritElixir: proto.ManaRegenElixir_ElixirOfSages}, want{spi: spirit, crit: 2}},
		"Elixir of the Whale":   {&proto.Consumes{SpiritElixir: proto.ManaRegenElixir_ElixirOfTheWhale}, want{spi: spirit, crit: 2}},
		// the three elixirs' +2% crit is one effect and never stacks
		"Mongoose + Owl":     {&proto.Consumes{AgilityElixir: proto.AgilityElixir_ElixirOfTheMongoose, IntellectElixir: proto.ManaRegenElixir_ElixirOfTheOwl}, want{agi: 25, intel: 25, crit: 2}},
		"Mongoose + Grizzly": {&proto.Consumes{AgilityElixir: proto.AgilityElixir_ElixirOfTheMongoose, StrengthBuff: proto.StrengthBuff_ElixirOfTheGrizzly}, want{agi: 25, str: 25, crit: 2}},
		"Grizzly + Owl":      {&proto.Consumes{StrengthBuff: proto.StrengthBuff_ElixirOfTheGrizzly, IntellectElixir: proto.ManaRegenElixir_ElixirOfTheOwl}, want{str: 25, intel: 25, crit: 2}},
		"Cunning + Owl":      {&proto.Consumes{AgilityElixir: proto.AgilityElixir_ElixirOfCunning, IntellectElixir: proto.ManaRegenElixir_ElixirOfTheOwl}, want{agi: 25, intel: 50, crit: 2}},
		"Mongoose + Sages":   {&proto.Consumes{AgilityElixir: proto.AgilityElixir_ElixirOfTheMongoose, SpiritElixir: proto.ManaRegenElixir_ElixirOfSages}, want{agi: 25, spi: spirit, crit: 2}},
		"Grizzly + Whale":    {&proto.Consumes{StrengthBuff: proto.StrengthBuff_ElixirOfTheGrizzly, SpiritElixir: proto.ManaRegenElixir_ElixirOfTheWhale}, want{str: 25, spi: spirit, crit: 2}},
		// one elixir per stat: the Owl, a Spirit elixir and a Mageblood all apply, the crit still once
		"Owl + Whale + Greater Mageblood": {&proto.Consumes{IntellectElixir: proto.ManaRegenElixir_ElixirOfTheOwl, SpiritElixir: proto.ManaRegenElixir_ElixirOfTheWhale, ManaRegenElixir: proto.ManaRegenElixir_GreaterMagebloodElixir}, want{intel: 25, spi: spirit, mp5: 20, crit: 2}},
		"Mongoose + Ferocity":             {&proto.Consumes{AgilityElixir: proto.AgilityElixir_ElixirOfTheMongoose, StrengthBuff: proto.StrengthBuff_FerocityElixir}, want{agi: 43, str: 18, crit: 2}},
		// settings saved before the slot per stat: each elixir moves to its stat's slot, or drops out if that is taken
		"old: Owl in the mana slot":     {&proto.Consumes{ManaRegenElixir: proto.ManaRegenElixir_ElixirOfTheOwl}, want{intel: 25, crit: 2}},
		"old: Whale in the mana slot":   {&proto.Consumes{ManaRegenElixir: proto.ManaRegenElixir_ElixirOfTheWhale, IntellectElixir: proto.ManaRegenElixir_ElixirOfTheOwl}, want{intel: 25, spi: spirit, crit: 2}},
		"old: Ferocity in Agility":      {&proto.Consumes{AgilityElixir: proto.AgilityElixir_ElixirOfFerocity}, want{agi: 18, str: 18}},
		"old: Ferocity with Juju Power": {&proto.Consumes{AgilityElixir: proto.AgilityElixir_ElixirOfFerocity, StrengthBuff: proto.StrengthBuff_JujuPower}, want{str: 30}},
	}
	critPerAgi := core.CritPerAgiAtLevel[proto.Class_ClassHunter] * core.CritRatingPerCritChance
	for name, tc := range cases {
		_, hunter := newPatchTestHunter(t, patchTestSetup{consumes: tc.consumes})
		agi := hunter.GetStat(stats.Agility) - base.GetStat(stats.Agility)
		got := want{
			agi:   agi,
			str:   hunter.GetStat(stats.Strength) - base.GetStat(stats.Strength),
			intel: hunter.GetStat(stats.Intellect) - base.GetStat(stats.Intellect),
			spi:   hunter.GetStat(stats.Spirit) - base.GetStat(stats.Spirit),
			mp5:   hunter.GetStat(stats.MP5) - base.GetStat(stats.MP5),
			crit:  (hunter.GetStat(stats.MeleeCrit) - base.GetStat(stats.MeleeCrit) - agi*critPerAgi) / core.CritRatingPerCritChance,
		}
		if math.Abs(got.agi-tc.want.agi) > 1e-6 || math.Abs(got.str-tc.want.str) > 1e-6 || math.Abs(got.intel-tc.want.intel) > 1e-6 ||
			math.Abs(got.spi-tc.want.spi) > 1e-6 || math.Abs(got.mp5-tc.want.mp5) > 1e-6 || math.Abs(got.crit-tc.want.crit) > 1e-6 {
			t.Errorf("%s: %+v, want %+v", name, got, tc.want)
		}
	}
	// the elixir crit is crit of every kind (aura 290), so spell crit too, and once
	for name, c := range map[string]*proto.Consumes{
		"Elixir of the Owl": {IntellectElixir: proto.ManaRegenElixir_ElixirOfTheOwl},
		"Mongoose + Owl":    {AgilityElixir: proto.AgilityElixir_ElixirOfTheMongoose, IntellectElixir: proto.ManaRegenElixir_ElixirOfTheOwl},
		"Grizzly + Owl":     {StrengthBuff: proto.StrengthBuff_ElixirOfTheGrizzly, IntellectElixir: proto.ManaRegenElixir_ElixirOfTheOwl},
		"Owl + Sages":       {IntellectElixir: proto.ManaRegenElixir_ElixirOfTheOwl, SpiritElixir: proto.ManaRegenElixir_ElixirOfSages},
	} {
		_, hunter := newPatchTestHunter(t, patchTestSetup{consumes: c})
		intel := hunter.GetStat(stats.Intellect) - base.GetStat(stats.Intellect)
		spellCrit := hunter.GetStat(stats.SpellCrit) - base.GetStat(stats.SpellCrit) - intel*core.CritPerIntAtLevel[proto.Class_ClassHunter]*core.SpellCritRatingPerCritChance
		if math.Abs(spellCrit-2*core.SpellCritRatingPerCritChance) > 1e-6 {
			t.Errorf("%s: spell crit +%.2f, want +2 (once)", name, spellCrit)
		}
	}
}

// Forever's 2H Impact is its own enchant at +6 weapon damage; Grand Crusader (+120 Strength on proc) is Forever's too.
func TestForeverWeaponEnchants(t *testing.T) {
	mhMin := func(enchant int32) float64 {
		_, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, mhEnchant: enchant})
		return hunter.AutoAttacks.MH().BaseDamageMin
	}
	if got := mhMin(8205) - mhMin(1897); got != 1 {
		t.Errorf("Impact (8205) adds %v more than the +5 enchant, want 1 (+6)", got)
	}
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, mhEnchant: 7940})
	proc := hunter.GetAura("Grand Crusader MH")
	if proc == nil {
		t.Fatal("Grand Crusader registers no proc aura")
	}
	before := hunter.GetStat(stats.Strength)
	proc.Activate(sim)
	if got := hunter.GetStat(stats.Strength) - before; got != 120 {
		t.Errorf("Grand Crusader proc gives %v Strength, want 120", got)
	}
	if proc.Duration != 20*time.Second {
		t.Errorf("Grand Crusader proc lasts %v, want 20s", proc.Duration)
	}
}

// Forever's Natural flasks: +60 Stamina each, plus their Hyjal effect (every encounter counts as Hyjal): 5% hit, 4%
// crit, 5% less dodge and parry, 5% attack speed.
func TestForeverNaturalFlasks(t *testing.T) {
	_, base := newPatchTestHunter(t, patchTestSetup{})
	flask := func(f proto.Flask) *Hunter {
		_, hunter := newPatchTestHunter(t, patchTestSetup{consumes: &proto.Consumes{Flask: f}})
		if got := hunter.GetStat(stats.Stamina) - base.GetStat(stats.Stamina); math.Abs(got-60) > 1e-6 {
			t.Errorf("%v: +%.1f Stamina, want 60", f, got)
		}
		return hunter
	}
	if h := flask(proto.Flask_FlaskOfNaturalAccuracy); math.Abs(h.GetStat(stats.MeleeHit)-base.GetStat(stats.MeleeHit)-5*core.MeleeHitRatingPerHitChance) > 1e-6 {
		t.Errorf("Natural Accuracy: hit +%.2f, want 5%%", (h.GetStat(stats.MeleeHit)-base.GetStat(stats.MeleeHit))/core.MeleeHitRatingPerHitChance)
	}
	if h := flask(proto.Flask_FlaskOfNaturalAggression); math.Abs(h.GetStat(stats.MeleeCrit)-base.GetStat(stats.MeleeCrit)-4*core.CritRatingPerCritChance) > 1e-6 {
		t.Errorf("Natural Aggression: crit +%.2f, want 4%%", (h.GetStat(stats.MeleeCrit)-base.GetStat(stats.MeleeCrit))/core.CritRatingPerCritChance)
	}
	if h := flask(proto.Flask_FlaskOfNaturalPrecision); math.Abs(h.GetStat(stats.Expertise)-base.GetStat(stats.Expertise)-5*core.ExpertiseRatingPerExpertiseChance) > 1e-6 {
		t.Errorf("Natural Precision: expertise +%.2f, want 5", h.GetStat(stats.Expertise)-base.GetStat(stats.Expertise))
	}
	h := flask(proto.Flask_FlaskOfNaturalSwiftness)
	for name, ratio := range map[string]float64{
		"melee":  h.PseudoStats.MeleeSpeedMultiplier / base.PseudoStats.MeleeSpeedMultiplier,
		"ranged": h.PseudoStats.RangedSpeedMultiplier / base.PseudoStats.RangedSpeedMultiplier,
	} {
		if math.Abs(ratio-1.05) > 1e-9 {
			t.Errorf("Natural Swiftness: %s speed x%.3f, want x1.05", name, ratio)
		}
	}

	// outside Hyjal (the encounter option) only the Stamina is left
	_, outBase := newPatchTestHunter(t, patchTestSetup{outside: true})
	_, out := newPatchTestHunter(t, patchTestSetup{outside: true, consumes: &proto.Consumes{Flask: proto.Flask_FlaskOfNaturalAggression}})
	if got := out.GetStat(stats.Stamina) - outBase.GetStat(stats.Stamina); math.Abs(got-60) > 1e-6 {
		t.Errorf("outside Hyjal: Natural Aggression +%.1f Stamina, want 60", got)
	}
	if got := out.GetStat(stats.MeleeCrit) - outBase.GetStat(stats.MeleeCrit); math.Abs(got) > 1e-6 {
		t.Errorf("outside Hyjal: Natural Aggression still adds %.2f crit", got/core.CritRatingPerCritChance)
	}
}

// Mortal Shots covers Auto Shot (client mask; Arc's log 2026-10-21: 5/5 Auto Shot crits at 2.30-2.32x, as Arcane Shot's),
// so Auto Shot crits get its bonus (4/5: +24% on the base +100%) like Arcane Shot; the option for the earlier "it skips
// Auto Shot" report takes it off Auto Shot only.
func TestMortalShotsOnAutoShot(t *testing.T) {
	const talents = "-00531500114-530200230050022151" // Mortal Shots 4/5
	for _, skip := range []bool{false, true} {
		_, hunter := newPatchTestHunter(t, patchTestSetup{talents: talents, gear: true, msSkipsAuto: skip})
		want := 1.24
		if skip {
			want = 1
		}
		if got := hunter.AutoAttacks.RangedAuto().CritDamageBonus; math.Abs(got-want) > 1e-9 {
			t.Errorf("skip %v: Auto Shot crit damage bonus %.2f, want %.2f", skip, got, want)
		}
		if got := hunter.ArcaneShot.CritDamageBonus; math.Abs(got-1.24) > 1e-9 {
			t.Errorf("skip %v: Arcane Shot crit damage bonus %.2f, want 1.24", skip, got)
		}
	}
}

// Mortal Shots' crit damage modifier covers Serpent Sting in the client data (its class mask holds the sting's), so a
// sting tick that crits gets its bonus too: 4/5 makes the double damage 2.24x.
func TestMortalShotsOnSerpentStingTicks(t *testing.T) {
	const talents = "-00531500114-530200230050022151" // Mortal Shots 4/5
	_, hunter := newPatchTestHunter(t, patchTestSetup{talents: talents, gear: true})
	if got := hunter.SerpentSting.CritDamageBonus; math.Abs(got-1.24) > 1e-9 {
		t.Errorf("Serpent Sting crit damage bonus %.2f, want 1.24 (Mortal Shots 4/5)", got)
	}
	at := hunter.AttackTables[hunter.CurrentTarget.UnitIndex][proto.CastType_CastTypeRanged]
	if got := hunter.SerpentSting.CritMultiplier(at); math.Abs(got-2.24) > 1e-9 {
		t.Errorf("Serpent Sting tick crit multiplier %.3f, want 2.24", got)
	}
}

// Spell haste shortens the GCD (forever-bugs #389): Skyborne's 3% casting speed takes Arcane Shot's 1.5 s GCD to 1.456 s,
// a big haste stops at 1 s, and a cast speed slow (Curse of Tongues) leaves it at 1.5 s (#284). With the switch off a
// Skyborne GCD stays 1.5 s.
func TestSpellHasteShortensGCD(t *testing.T) {
	gcdAfterArcaneShot := func(race proto.Race, castSpeed float64) time.Duration {
		sim, hunter := newPatchTestHunter(t, patchTestSetup{race: race, talents: patchTestTalents, gear: true})
		if castSpeed != 1 {
			hunter.MultiplyCastSpeed(castSpeed)
		}
		hunter.DistanceFromTarget = 20
		if !hunter.ArcaneShot.Cast(sim, hunter.CurrentTarget) {
			t.Fatal("Arcane Shot did not cast")
		}
		return hunter.GCD.TimeToReady(sim)
	}
	for _, tc := range []struct {
		name      string
		race      proto.Race
		castSpeed float64
		want      time.Duration
	}{
		{"Human", proto.Race_RaceHuman, 1, 1500 * time.Millisecond},
		{"Skyborne", proto.Race_RaceSkyborneAlliance, 1, time.Duration(math.Round(1.5e9 / 1.02))},
		{"double haste", proto.Race_RaceHuman, 2, time.Second},
		{"cast speed slow", proto.Race_RaceHuman, 1 / 1.6, 1500 * time.Millisecond},
	} {
		if got := gcdAfterArcaneShot(tc.race, tc.castSpeed); (got - tc.want).Abs() > time.Millisecond {
			t.Errorf("%s: GCD %v, want %v", tc.name, got, tc.want)
		}
	}
	core.ForeverSpellHasteGCD = false
	defer func() { core.ForeverSpellHasteGCD = true }()
	if got := gcdAfterArcaneShot(proto.Race_RaceSkyborneAlliance, 1); got != 1500*time.Millisecond {
		t.Errorf("switch off: Skyborne GCD %v, want 1.5s", got)
	}
}

// Skysight (Horde Skyborne): +15% run speed, 15 s on a 2 min cooldown on the GCD. It is a regular speed bonus, so under
// Strider Kick's 30% it adds nothing (AB_SKYSIGHT_STACKS=1 stacks it, for comparing). Charged (cast at an elemental
// convergence before the fight) makes each cast +30%, still for 15 s: nothing is up until it is cast.
func TestSkysight(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{race: proto.Race_RaceSkyborneHorde, talents: patchTestTalents})
	aura, spell := hunter.GetAura("Skysight"), hunter.GetSpell(core.ActionID{SpellID: 1259686})
	if aura == nil || spell == nil {
		t.Fatal("no Skysight")
	}
	if aura.IsActive() {
		t.Error("Skysight is up before it is cast")
	}
	if spell.CD.Duration != 2*time.Minute || spell.DefaultCast.GCD == 0 {
		t.Errorf("Skysight cooldown %v, GCD %v; want 2 min on the GCD", spell.CD.Duration, spell.DefaultCast.GCD)
	}
	spell.Cast(sim, hunter.CurrentTarget)
	if got := hunter.MovementHandler.MoveSpeed; math.Abs(got-7*1.15) > 1e-9 {
		t.Errorf("Skysight: %.4f yd/s, want %.4f", got, 7*1.15)
	}
	if aura.Duration != 15*time.Second {
		t.Errorf("Skysight lasts %v, want 15s", aura.Duration)
	}
	hunter.GetAura("Strider Kick Speed").Activate(sim)
	want := 7 * 1.3
	if core.SkysightStacks {
		want *= 1.15
	}
	if got := hunter.MovementHandler.MoveSpeed; math.Abs(got-want) > 1e-9 {
		t.Errorf("Skysight under Strider Kick: %.4f yd/s, want %.4f", got, want)
	}

	sim, charged := newPatchTestHunter(t, patchTestSetup{race: proto.Race_RaceSkyborneHorde, skyborne: true})
	if a := charged.GetAura("Skysight"); a.IsActive() {
		t.Error("Charged: Skysight is up before it is cast")
	}
	if got := charged.MovementHandler.MoveSpeed; math.Abs(got-7) > 1e-9 {
		t.Errorf("Charged, not cast: %.4f yd/s, want 7", got)
	}
	charged.GetSpell(core.ActionID{SpellID: 1259686}).Cast(sim, charged.CurrentTarget)
	if got := charged.MovementHandler.MoveSpeed; math.Abs(got-7*1.3) > 1e-9 {
		t.Errorf("Charged Skysight: %.4f yd/s, want %.4f", got, 7*1.3)
	}
	if d := charged.GetAura("Skysight").Duration; d != 15*time.Second {
		t.Errorf("Charged Skysight lasts %v, want 15s", d)
	}
}

// Read Leyline (Alliance Skyborne): +50% mana regeneration for 15 s when cast; Attuned (cast at a leyline before the
// fight) makes each cast +100% (double the spirit regen), still for 15 s; Horde has Skysight instead.
func TestReadLeyline(t *testing.T) {
	sim, plain := newPatchTestHunter(t, patchTestSetup{race: proto.Race_RaceSkyborneAlliance})
	asim, attuned := newPatchTestHunter(t, patchTestSetup{race: proto.Race_RaceSkyborneAlliance, skyborne: true})
	if plain.GetAura("Read Leyline") == nil || plain.GetAura("Skysight") != nil {
		t.Fatal("Alliance Skyborne should have Read Leyline, not Skysight")
	}
	if attuned.PseudoStats.SpiritRegenMultiplier != plain.PseudoStats.SpiritRegenMultiplier {
		t.Error("Attuned: Read Leyline is up before it is cast")
	}
	abefore := attuned.PseudoStats.SpiritRegenMultiplier
	attuned.GetSpell(core.ActionID{SpellID: 1259705}).Cast(asim, attuned.CurrentTarget)
	if got := attuned.PseudoStats.SpiritRegenMultiplier - abefore; math.Abs(got-1) > 1e-9 {
		t.Errorf("Attuned Read Leyline adds %.2f to the spirit regen multiplier, want 1 (+100%%)", got)
	}
	before := plain.PseudoStats.SpiritRegenMultiplier
	plain.GetSpell(core.ActionID{SpellID: 1259705}).Cast(sim, plain.CurrentTarget)
	if got := plain.PseudoStats.SpiritRegenMultiplier - before; math.Abs(got-0.5) > 1e-9 {
		t.Errorf("cast Read Leyline adds %.2f to the spirit regen multiplier, want 0.5 (+50%%)", got)
	}
	if d := plain.GetAura("Read Leyline").Duration; d != 15*time.Second {
		t.Errorf("Read Leyline lasts %v, want 15s", d)
	}
}

// Forever's set bonuses (client ItemSet data): Beaststalker +30 AP at 3 pieces, The Gladiator +20 AP at 3 and
// +1% crit at 5, Cadaverous Garb +2% hit at 5. Each bonus is measured as (all pieces) - (one piece fewer) - (that piece
// alone) + (nothing), which leaves the bonus without the piece's own stats.
func TestForeverSetBonuses(t *testing.T) {
	stat := func(pieces map[proto.ItemSlot]int32, s stats.Stat) float64 {
		_, hunter := newPatchTestHunter(t, patchTestSetup{pieces: pieces})
		return hunter.GetStat(s)
	}
	bonus := func(set map[proto.ItemSlot]int32, last proto.ItemSlot, s stats.Stat) float64 {
		fewer := map[proto.ItemSlot]int32{}
		for slot, id := range set {
			if slot != last {
				fewer[slot] = id
			}
		}
		alone := map[proto.ItemSlot]int32{last: set[last]}
		return stat(set, s) - stat(fewer, s) - stat(alone, s) + stat(map[proto.ItemSlot]int32{}, s)
	}
	head, chest, hands, legs, feet := proto.ItemSlot_ItemSlotHead, proto.ItemSlot_ItemSlotChest, proto.ItemSlot_ItemSlotHands, proto.ItemSlot_ItemSlotLegs, proto.ItemSlot_ItemSlotFeet
	waist := proto.ItemSlot_ItemSlotWaist

	beast3 := map[proto.ItemSlot]int32{head: 16677, chest: 16674, legs: 16678}
	if got := bonus(beast3, legs, stats.RangedAttackPower); math.Abs(got-30) > 1e-6 {
		t.Errorf("Beaststalker 3 pieces: +%.1f ranged AP, want 30", got)
	}
	glad3 := map[proto.ItemSlot]int32{head: 11729, chest: 11726, legs: 11728}
	if got := bonus(glad3, legs, stats.AttackPower); math.Abs(got-20) > 1e-6 {
		t.Errorf("The Gladiator 3 pieces: +%.1f AP, want 20", got)
	}
	glad5 := map[proto.ItemSlot]int32{head: 11729, chest: 11726, legs: 11728, hands: 11730, feet: 11731}
	if got := bonus(glad5, feet, stats.MeleeCrit) / core.CritRatingPerCritChance; math.Abs(got-1) > 1e-6 {
		t.Errorf("The Gladiator 5 pieces: +%.2f%% crit, want 1", got)
	}
	cadaver5 := map[proto.ItemSlot]int32{waist: 14636, chest: 14637, legs: 14638, hands: 14640, feet: 14641}
	if got := bonus(cadaver5, feet, stats.MeleeHit) / core.MeleeHitRatingPerHitChance; math.Abs(got-2) > 1e-6 {
		t.Errorf("Cadaverous Garb 5 pieces: +%.2f%% hit, want 2", got)
	}
	// Devilsaur Armor 2 pieces (spell 460230): +2% hit and +2% spell hit
	devil2 := map[proto.ItemSlot]int32{legs: 15062, hands: 15063}
	if got := bonus(devil2, hands, stats.MeleeHit) / core.MeleeHitRatingPerHitChance; math.Abs(got-2) > 1e-6 {
		t.Errorf("Devilsaur Armor 2 pieces: +%.2f%% hit, want 2", got)
	}
	if got := bonus(devil2, hands, stats.SpellHit) / core.SpellHitRatingPerHitChance; math.Abs(got-2) > 1e-6 {
		t.Errorf("Devilsaur Armor 2 pieces: +%.2f%% spell hit, want 2", got)
	}
}

// Forever explosives: Crystal Charge (instant) and the EZ-Thro line need no profession, Dark Iron Bomb needs Engineering;
// all but the sapper and Crystal Charge take 1 s to throw.
func TestForeverExplosives(t *testing.T) {
	spell := func(filler proto.Explosive, engineer bool, item int32) *core.Spell {
		_, hunter := newPatchTestHunter(t, patchTestSetup{engineer: engineer, consumes: &proto.Consumes{FillerExplosive: filler}})
		return hunter.GetSpell(core.ActionID{ItemID: item})
	}
	if s := spell(proto.Explosive_ExplosiveCrystalCharge, false, 11566); s == nil || s.DefaultCast.CastTime != 0 {
		t.Error("Crystal Charge: wanted an instant throw without Engineering")
	}
	if s := spell(proto.Explosive_ExplosiveEzThroDarkBomb, false, 260817); s == nil || s.DefaultCast.CastTime != time.Second {
		t.Error("EZ-Thro Dark Bomb: wanted a 1 s throw without Engineering")
	}
	if s := spell(proto.Explosive_ExplosiveDarkIronBomb, false, 16005); s != nil {
		t.Error("Dark Iron Bomb works without Engineering")
	}
	if s := spell(proto.Explosive_ExplosiveDarkIronBomb, true, 16005); s == nil {
		t.Error("Dark Iron Bomb missing for an engineer")
	}
	if s := spell(proto.Explosive_ExplosiveHolyHandGrenade, false, 284868); s == nil || s.SpellSchool != core.SpellSchoolHoly {
		t.Error("Holy Hand Grenade: wanted Holy damage without Engineering")
	}
	if s := spell(proto.Explosive_ExplosiveCrystalCharge, false, 11566); s == nil || s.SpellSchool != core.SpellSchoolFire {
		t.Error("Crystal Charge: wanted Fire damage")
	}
	// a Mage item: a Hunter who picks it throws nothing (the Mage side is in sim/forever_racials_test.go)
	if s := spell(proto.Explosive_ExplosiveScrollOfGreaterCryoblast, true, 277499); s != nil {
		t.Error("a Hunter has the Scroll of Greater Cryoblast, a Mage item")
	}
}

// Forever's damage potions from the client data, on the shared 2 min potion cooldown. Dragonfire (12 yd cone): 280 Fire
// at once, then 35 every 2 sec for 8 sec, always hits, never crits. Caustic Smog (thrown at 10 yd/s from up to 20 yd):
// a cloud that deals 33 Nature +-5% every sec for 12 sec from where it lands. Both always hit and add no spell power.
func TestForeverDamagePotions(t *testing.T) {
	type tick struct {
		at     time.Duration
		damage float64
	}
	// drinks the potion at 5 s from the given distance and records each periodic tick until the fight is 20 s in
	drink := func(potion proto.Potions, item int32, distance float64) (bool, time.Duration, []tick, *core.Simulation, *Hunter) {
		sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: distance, consumes: &proto.Consumes{DefaultPotion: potion}})
		mcd := hunter.GetMajorCooldown(core.ActionID{ItemID: item})
		if mcd == nil {
			t.Fatalf("no potion cooldown for item %d", item)
		}
		sim.PrePull()
		for i := 0; i < 100000 && sim.CurrentTime < 5*time.Second; i++ {
			sim.Step()
		}
		drunkAt := sim.CurrentTime
		if !mcd.TryActivate(sim, hunter.GetCharacter()) {
			return false, 0, nil, sim, hunter
		}
		cd := mcd.TimeToReady(sim)
		m := &mcd.Spell.SpellMetrics[hunter.CurrentTarget.UnitIndex]
		var ticks []tick
		for i := 0; i < 100000 && sim.CurrentTime < 20*time.Second; i++ {
			n, d := m.Ticks, m.TotalTickDamage
			sim.Step()
			if m.Ticks > n {
				ticks = append(ticks, tick{sim.CurrentTime - drunkAt, (m.TotalTickDamage - d) / hunter.PseudoStats.DamageDealtMultiplier})
			}
		}
		return true, cd, ticks, sim, hunter
	}

	if used, _, _, _, _ := drink(proto.Potions_DragonfirePotion, 250953, 12.5); used {
		t.Error("Dragonfire Potion used from 12.5 yd, past its 12 yd cone")
	}
	used, cd, ticks, sim, hunter := drink(proto.Potions_DragonfirePotion, 250953, 11)
	if !used {
		t.Fatal("Dragonfire Potion not used from 11 yd")
	}
	if cd != 2*time.Minute {
		t.Errorf("Dragonfire Potion cooldown %v, want 2m", cd)
	}
	if len(ticks) != 4 {
		t.Errorf("Dragonfire Potion ticked %d times, want 4", len(ticks))
	}
	for i, tk := range ticks {
		if want := time.Duration(i+1) * 2 * time.Second; tk.at != want || math.Abs(tk.damage-35) > 1e-6 {
			t.Errorf("Dragonfire Potion tick %d: %.2f at %v, want 35 at %v", i+1, tk.damage, tk.at, want)
		}
	}
	// the cone at once: every one lands for exactly 280, even at 100% spell crit
	potion := hunter.GetSpell(core.ActionID{ItemID: 250953})
	m := &potion.SpellMetrics[hunter.CurrentTarget.UnitIndex]
	hunter.AddStatDynamic(sim, stats.SpellCrit, 100*core.SpellCritRatingPerCritChance)
	for i := 0; i < 500; i++ {
		hits, damage := m.Hits, m.TotalDamage
		potion.CD.Reset()
		potion.Cast(sim, hunter.CurrentTarget)
		if d := (m.TotalDamage - damage) / hunter.PseudoStats.DamageDealtMultiplier; m.Hits != hits+1 || math.Abs(d-280) > 1e-6 {
			t.Fatalf("Dragonfire Potion cast %d: %.2f damage, %d hits, %d misses, %d crits; want 280 every time", i+1, d, m.Hits-hits, m.Misses, m.Crits)
		}
	}
	if m.Crits != 0 || m.Misses != 0 {
		t.Errorf("Dragonfire Potion: %d crits and %d misses, want none", m.Crits, m.Misses)
	}

	if used, _, _, _, _ := drink(proto.Potions_CausticSmogPotion, 250951, 21); used {
		t.Error("Caustic Smog Potion thrown from 21 yd, past its 20 yd range")
	}
	used, cd, ticks, _, _ = drink(proto.Potions_CausticSmogPotion, 250951, 11)
	if !used {
		t.Fatal("Caustic Smog Potion not thrown from 11 yd")
	}
	if cd != 2*time.Minute {
		t.Errorf("Caustic Smog Potion cooldown %v, want 2m", cd)
	}
	if len(ticks) != 12 {
		t.Errorf("Caustic Smog Potion ticked %d times, want 12", len(ticks))
	}
	lands := 1100 * time.Millisecond // 11 yd at 10 yd/s
	allSame := true
	for i, tk := range ticks {
		if want := lands + time.Duration(i+1)*time.Second; tk.at != want || tk.damage < 33*0.95-1e-6 || tk.damage > 33*1.05+1e-6 {
			t.Errorf("Caustic Smog Potion tick %d: %.2f at %v, want 31.35 to 34.65 at %v", i+1, tk.damage, tk.at, want)
		}
		allSame = allSame && math.Abs(tk.damage-ticks[0].damage) < 1e-6
	}
	if allSame {
		t.Error("Caustic Smog Potion ticks all the same, want +-5% each")
	}
}

// Potion of Beast Slaying (spell 1252559): +180 attack power against Beasts for 5 min, on the shared 2 min potion cooldown.
// Its client aura has no ranged twin, so only melee attack power rises; another creature type gains nothing, and a fight
// with no Beast never drinks it. No second drink while the first lasts.
func TestPotionOfBeastSlaying(t *testing.T) {
	setup := func(mob proto.MobType) (*core.Simulation, *Hunter, *core.MajorCooldown, func() (float64, float64)) {
		sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, mobType: mob, consumes: &proto.Consumes{DefaultPotion: proto.Potions_PotionOfBeastSlaying}})
		mcd := hunter.GetMajorCooldown(core.ActionID{ItemID: 250955})
		if mcd == nil {
			t.Fatal("no Potion of Beast Slaying cooldown")
		}
		sim.PrePull()
		for i := 0; i < 100000 && sim.CurrentTime < time.Second; i++ {
			sim.Step()
		}
		target := hunter.CurrentTarget
		attackPower := func() (float64, float64) {
			return hunter.AutoAttacks.MHAuto().MeleeAttackPower(target), hunter.AutoAttacks.RangedAuto().RangedAttackPower(target, false)
		}
		return sim, hunter, mcd, attackPower
	}

	sim, hunter, mcd, attackPower := setup(proto.MobType_MobTypeBeast)
	melee, ranged := attackPower()
	if !mcd.TryActivate(sim, hunter.GetCharacter()) {
		t.Fatal("Potion of Beast Slaying not drunk against a Beast")
	}
	if m, r := attackPower(); m-melee != 180 || r != ranged {
		t.Errorf("against a Beast: melee attack power +%.0f, ranged +%.0f; want +180 and +0", m-melee, r-ranged)
	}
	aura := hunter.GetAura("Potion of Beast Slaying")
	if aura == nil || aura.Duration != 5*time.Minute {
		t.Fatal("no 5 min Potion of Beast Slaying aura")
	}
	if cd := mcd.TimeToReady(sim); cd != 2*time.Minute {
		t.Errorf("Potion of Beast Slaying cooldown %v, want 2m", cd)
	}
	mcd.Spell.CD.Reset()
	if mcd.TryActivate(sim, hunter.GetCharacter()) {
		t.Error("Potion of Beast Slaying drunk again while it lasts")
	}
	aura.Deactivate(sim)
	if m, _ := attackPower(); m != melee {
		t.Errorf("melee attack power %.0f after the potion ends, want %.0f", m, melee)
	}

	sim, hunter, mcd, attackPower = setup(proto.MobType_MobTypeHumanoid)
	if mcd.TryActivate(sim, hunter.GetCharacter()) {
		t.Error("Potion of Beast Slaying drunk with no Beast to fight")
	}
	melee, _ = attackPower()
	mcd.Spell.Cast(sim, hunter.CurrentTarget)
	if m, _ := attackPower(); m != melee {
		t.Errorf("against a Humanoid: melee attack power +%.0f, want +0", m-melee)
	}
}

// Forever lets anyone wear the Zul'Gurub crafted pieces, but their set bonuses still need the crafting profession at 300
// (ItemSet RequiredSkill): Primal Batskin's 3-piece (8% run speed, not stacking with other speed bonuses) and Blood Tiger's
// 2-piece (1% crit) need Leatherworking, The Darksoul's 2-piece (20 Defense) Blacksmithing.
func TestCraftedSetBonusesNeedTheirProfession(t *testing.T) {
	batskin := map[proto.ItemSlot]int32{proto.ItemSlot_ItemSlotChest: 19685, proto.ItemSlot_ItemSlotWrist: 19687, proto.ItemSlot_ItemSlotHands: 19686}
	for _, prof := range []proto.Profession{proto.Profession_ProfessionUnknown, proto.Profession_Leatherworking} {
		_, hunter := newPatchTestHunter(t, patchTestSetup{pieces: batskin, profession: prof})
		want := 7.0
		if prof == proto.Profession_Leatherworking {
			want = 7 * 1.08
		}
		if got := hunter.MovementHandler.MoveSpeed; math.Abs(got-want) > 1e-9 {
			t.Errorf("Primal Batskin, profession %v: run speed %.3f, want %.3f", prof, got, want)
		}
	}

	bonus := func(pieces map[proto.ItemSlot]int32, prof proto.Profession, stat stats.Stat) float64 {
		_, without := newPatchTestHunter(t, patchTestSetup{pieces: pieces})
		_, with := newPatchTestHunter(t, patchTestSetup{pieces: pieces, profession: prof})
		return with.GetStat(stat) - without.GetStat(stat)
	}
	tiger := map[proto.ItemSlot]int32{proto.ItemSlot_ItemSlotChest: 19688, proto.ItemSlot_ItemSlotShoulder: 19689}
	if got := bonus(tiger, proto.Profession_Leatherworking, stats.MeleeCrit); math.Abs(got-core.CritRatingPerCritChance) > 1e-9 {
		t.Errorf("Blood Tiger Harness: Leatherworking adds %.2f crit, want 1%% (only a leatherworker gets the bonus)", got)
	}
	darksoul := map[proto.ItemSlot]int32{proto.ItemSlot_ItemSlotChest: 19693, proto.ItemSlot_ItemSlotLegs: 19694}
	if got := bonus(darksoul, proto.Profession_Blacksmithing, stats.Defense); got != 20 {
		t.Errorf("The Darksoul: Blacksmithing adds %.0f Defense, want 20 (only a blacksmith gets the bonus)", got)
	}
}

// The log carries every weapon's swing timer for the timeline: each main-hand swing starts a timer ready one swing time
// later, a melee swing restarts the Auto Shot timer on the spot, ready one bow cycle later, and an extra attack (Hand of
// Justice) restarts the main hand's own timer, never one of its own. Only the player's lines count (the pet swings too).
// resetOptions: the A/B hunter options with the beta's melee reset of Auto Shot on.
func resetOptions() *proto.Player_Hunter {
	opts := &proto.Player_Hunter{Hunter: googleProto.Clone(abOptions.Hunter).(*proto.Hunter)}
	opts.Hunter.Options.MeleeResetsAutoShot = true
	return opts
}

func TestSwingTimerLogs(t *testing.T) {
	player := core.WithSpec(&proto.Player{
		Race:               proto.Race_RaceHuman,
		Class:              proto.Class_ClassHunter,
		Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_light_weave").GearSet,
		TalentsString:      patchTestTalents,
		Consumes:           abConsumes,
		Buffs:              abIndividualBuffs,
		Rotation:           core.GetAplRotation("../../ui/hunter/apls", "arc_weave_2h").Rotation,
		DistanceFromTarget: 3,
	}, resetOptions())
	// the first fight with a Hand of Justice proc, so a rotation change does not leave the check without one
	var res *proto.RaidSimResult
	for seed := int64(1); seed <= 40; seed++ {
		res = core.RunRaidSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, abDebuffs),
			Encounter:  abEncounter(),
			SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: seed, DebugFirstIteration: true},
		})
		if res.Error != nil {
			t.Fatalf("sim error: %s", res.Error.Message)
		}
		if strings.Contains(res.Logs, "extra main-hand attack") {
			break
		}
	}
	mh := core.ActionID{OtherID: proto.OtherAction_OtherActionAttack, Tag: 1}.String()
	shot := core.ActionID{OtherID: proto.OtherAction_OtherActionShoot}.String()
	// the player's lines only: "[name (#1)]", not the pet's "[name (#1) - pet]" or the target's
	lineRe := regexp.MustCompile(`^\[(-?[0-9.]+)\] \[[^\]-]*\(#\d+\)\] (.*)$`)
	startRe := regexp.MustCompile(`^Swing timer (\{[^}]*\}) started, ready at (-?[0-9.]+)s$`)

	mhSwings := map[string]bool{}
	type start struct{ at, ready float64 }
	starts := map[string][]start{}
	extraAttacks := 0
	for _, line := range strings.Split(res.Logs, "\n") {
		m := lineRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		if strings.Contains(m[2], "extra main-hand attack") {
			extraAttacks++
		}
		if strings.HasPrefix(m[2], "Casting "+mh) {
			mhSwings[m[1]] = true
		}
		if s := startRe.FindStringSubmatch(m[2]); s != nil {
			at, _ := strconv.ParseFloat(m[1], 64)
			ready, _ := strconv.ParseFloat(s[2], 64)
			starts[s[1]] = append(starts[s[1]], start{at, ready})
		}
	}
	if extraAttacks == 0 {
		t.Fatal("the fight had no extra main-hand attack to check; pick another seed")
	}
	oh := core.ActionID{OtherID: proto.OtherAction_OtherActionAttack, Tag: 2}.String()
	for id := range starts {
		if id != mh && id != oh && id != shot {
			t.Errorf("swing timer logged under %s; want the weapon's own id (main hand, off hand or Auto Shot)", id)
		}
	}
	if len(starts[mh]) < 10 || len(starts[shot]) < 10 {
		t.Fatalf("swing timer starts: main hand %d, Auto Shot %d; want both logged", len(starts[mh]), len(starts[shot]))
	}
	// Arcanite Reaper swings every 3.8 s, less with melee haste
	for _, s := range starts[mh] {
		if d := s.ready - s.at; d < 2.5 || d > 3.81 {
			t.Errorf("main-hand timer at %.2f s ready %.2f s later; want one 3.8 s swing time or less", s.at, d)
		}
	}
	restarts := 0
	for _, s := range starts[shot] {
		if d := s.ready - s.at; d < 1 || d > 3.5 {
			t.Errorf("Auto Shot timer at %.2f s ready %.2f s later; want one bow cycle", s.at, d)
		}
		if mhSwings[strconv.FormatFloat(s.at, 'f', 2, 64)] {
			restarts++
		}
	}
	if restarts == 0 {
		t.Error("no Auto Shot timer started on a main-hand swing; melee restarts it")
	}
}

// The Walk Time rotation value reads the real run speed: from 12 yd, walking to 8 yd takes 4 / 7 = 0.571 s at normal
// speed and 4 / 7.14 = 0.560 s with Night Elf's 2%, so a cast gated on "under 0.565 s" goes out for a Night Elf only.
func TestWalkTimeValue(t *testing.T) {
	const aplJSON = `{"type":"TypeAPL","recheckAfterAutoShot":true,"priorityList":[{"action":{"condition":{"cmp":{"op":"OpLt","lhs":{"walkTime":{"toDistance":{"const":{"val":"8"}}}},"rhs":{"const":{"val":"0.565s"}}}},"castSpell":{"spellId":{"spellId":14287,"rank":8}}}}]}`
	arcaneShots := func(race proto.Race) int32 {
		rot := &proto.APLRotation{}
		if err := protojson.Unmarshal([]byte(aplJSON), rot); err != nil {
			t.Fatalf("rotation: %v", err)
		}
		player := core.WithSpec(&proto.Player{
			Race:               race,
			Class:              proto.Class_ClassHunter,
			Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
			TalentsString:      patchTestTalents,
			Consumes:           abConsumes,
			Buffs:              abIndividualBuffs,
			Rotation:           rot,
			DistanceFromTarget: 12,
		}, abOptions)
		res := core.RunRaidSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, abDebuffs),
			Encounter:  abEncounter(),
			SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 20261017},
		})
		if res.Error != nil {
			t.Fatalf("sim error: %s", res.Error.Message)
		}
		var casts int32
		for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
			if a.Id.GetSpellId() == 14287 {
				for _, tm := range a.Targets {
					casts += tm.Casts
				}
			}
		}
		return casts
	}
	if n := arcaneShots(proto.Race_RaceNightElf); n == 0 {
		t.Error("Night Elf: no Arcane Shot; walking 4 yd at 7.14 yd/s (0.560 s) should pass the 0.565 s gate")
	}
	if n := arcaneShots(proto.Race_RaceHuman); n != 0 {
		t.Errorf("Human: %d Arcane Shots; walking 4 yd at 7 yd/s (0.571 s) should fail the 0.565 s gate", n)
	}
}

// The ranged minimum is 11 yd from the target's center (the client's 8 yd is hitbox to hitbox; in game shots stop inside
// about 10.5-11 yd of a mob's center): standing at 10 yd a hunter fires no Auto Shot and no Arcane Shot, at 11 yd both.
func TestRangedMinimumElevenYards(t *testing.T) {
	const aplJSON = `{"type":"TypeAPL","priorityList":[{"action":{"castSpell":{"spellId":{"spellId":14287,"rank":8}}}}]}`
	shots := func(distance float64) (autos, arcanes int32) {
		rot := &proto.APLRotation{}
		if err := protojson.Unmarshal([]byte(aplJSON), rot); err != nil {
			t.Fatalf("rotation: %v", err)
		}
		player := core.WithSpec(&proto.Player{
			Race:               proto.Race_RaceHuman,
			Class:              proto.Class_ClassHunter,
			Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
			TalentsString:      patchTestTalents,
			Consumes:           abConsumes,
			Buffs:              abIndividualBuffs,
			Rotation:           rot,
			DistanceFromTarget: distance,
		}, abOptions)
		res := core.RunRaidSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, abDebuffs),
			Encounter:  abEncounter(),
			SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 20261018},
		})
		if res.Error != nil {
			t.Fatalf("sim error: %s", res.Error.Message)
		}
		for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
			for _, tm := range a.Targets {
				switch {
				case a.Id.GetOtherId() == proto.OtherAction_OtherActionShoot:
					autos += tm.Casts
				case a.Id.GetSpellId() == 14287:
					arcanes += tm.Casts
				}
			}
		}
		return
	}
	if autos, arcanes := shots(10); autos != 0 || arcanes != 0 {
		t.Errorf("from 10 yd: %d Auto Shots, %d Arcane Shots; want none (inside the 11 yd minimum)", autos, arcanes)
	}
	if autos, arcanes := shots(11); autos == 0 || arcanes == 0 {
		t.Errorf("from 11 yd: %d Auto Shots, %d Arcane Shots; want both firing", autos, arcanes)
	}
}

// Trueshot Aura and Hunter's Mark are castable for a rotation's pull: each goes out once before the pull, and the mark
// is on the target with the raid's Hunter's Mark debuff off. With that debuff on, the mark is up the whole fight
// (the setting used to register it without ever putting it up) and the cast adds nothing on top. A hunter without the
// talent skips Trueshot Aura, and a rotation shared with talented builds raises no warning for it.
func TestPrepullTrueshotAndMark(t *testing.T) {
	const talented, untalented = "-0053150011-500230230250222151", "-0053150010-500230230250222151"
	run := func(talents string, cast bool, debuffs *proto.Debuffs) (casts map[int32]int32, markUptime, fight, dps float64, warnings []string) {
		prepull := ""
		if cast {
			prepull = `{"action":{"castSpell":{"spellId":{"spellId":20905}}},"doAtValue":{"const":{"val":"-13s"}}},` +
				`{"action":{"castSpell":{"spellId":{"spellId":14325}}},"doAtValue":{"const":{"val":"-4s"}}}`
		}
		rot := &proto.APLRotation{}
		aplJSON := `{"type":"TypeAPL","prepullActions":[` + prepull + `],` +
			`"priorityList":[{"action":{"castSpell":{"spellId":{"spellId":14287,"rank":8}}}}]}`
		if err := protojson.Unmarshal([]byte(aplJSON), rot); err != nil {
			t.Fatalf("rotation: %v", err)
		}
		player := core.WithSpec(&proto.Player{
			Race:               proto.Race_RaceHuman,
			Class:              proto.Class_ClassHunter,
			Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
			TalentsString:      talents,
			Consumes:           abConsumes,
			Buffs:              abIndividualBuffs,
			Rotation:           rot,
			DistanceFromTarget: 12,
		}, abOptions)
		raid := core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, debuffs)
		for _, a := range core.ComputeStats(&proto.ComputeStatsRequest{Raid: raid, Encounter: abEncounter()}).RaidStats.Parties[0].Players[0].RotationStats.PrepullActions {
			warnings = append(warnings, a.Warnings...)
		}
		res := core.RunRaidSim(&proto.RaidSimRequest{
			Raid:       raid,
			Encounter:  abEncounter(),
			SimOptions: &proto.SimOptions{Iterations: 50, RandomSeed: 20261019},
		})
		if res.Error != nil {
			t.Fatalf("sim error: %s", res.Error.Message)
		}
		casts = map[int32]int32{}
		for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
			for _, tm := range a.Targets {
				casts[a.Id.GetSpellId()] += tm.Casts
			}
		}
		for _, a := range res.EncounterMetrics.Targets[0].Auras {
			if a.Id.GetSpellId() == 14325 {
				markUptime = a.UptimeSecondsAvg
			}
		}
		return casts, markUptime, res.AvgIterationDuration, res.RaidMetrics.Dps.Avg, warnings
	}

	casts, markUptime, _, _, warnings := run(talented, true, &proto.Debuffs{})
	if len(warnings) > 0 {
		t.Errorf("pull casts warn: %v", warnings)
	}
	if casts[20905] != 50 || casts[14325] != 50 {
		t.Errorf("Trueshot Aura cast %d times in 50 fights, Hunter's Mark %d; want each once a fight, before the pull", casts[20905], casts[14325])
	}
	if markUptime < 60 {
		t.Errorf("Hunter's Mark up %.1f s on the target; want it from the pull on", markUptime)
	}

	casts, _, _, _, warnings = run(untalented, true, &proto.Debuffs{})
	if casts[20905] != 0 || casts[14325] != 50 || len(warnings) > 0 {
		t.Errorf("without the talent: Trueshot Aura cast %d times, Hunter's Mark %d, warnings %v; want 0, 50, none", casts[20905], casts[14325], warnings)
	}

	_, raidUptime, fight, raidDps, _ := run(talented, false, abDebuffs)
	if raidUptime < fight-0.5 {
		t.Errorf("raid Hunter's Mark up %.1f s of a %.1f s fight; want all of it", raidUptime, fight)
	}
	_, _, _, bothDps, _ := run(talented, true, abDebuffs)
	if math.Abs(bothDps-raidDps) > raidDps*0.005 {
		t.Errorf("casting the mark over the raid's: %.1f DPS vs %.1f; want the same (one mark)", bothDps, raidDps)
	}
}

// Expose Prey is its own buff, and every one ends one way: a Mongoose Bite uses it (one that misses too), or it runs out unused (one
// still up when the fight ends counts as neither), and a proc while it is up refreshes it. The markers count each: used
// matches the Mongoose Bites, and used + expired is every proc but at most one a fight.
func TestExposePreyMarkers(t *testing.T) {
	const iterations = 50
	player := core.WithSpec(&proto.Player{
		Race:               proto.Race_RaceHuman,
		Class:              proto.Class_ClassHunter,
		Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "surv_dw").GearSet,
		TalentsString:      "-00531500111-500230230250122151",
		Consumes:           abConsumes,
		Buffs:              abIndividualBuffs,
		Rotation:           core.GetAplRotation("../../ui/hunter/apls", "arc_weave_2h").Rotation,
		DistanceFromTarget: 3,
	}, abOptions)
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, abDebuffs),
		Encounter:  abEncounter(),
		SimOptions: &proto.SimOptions{Iterations: iterations, RandomSeed: 20261019},
	})
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
	}
	um := res.RaidMetrics.Parties[0].Players[0]
	procs := map[int32]float64{}
	for _, a := range um.Auras {
		if a.Id.GetSpellId() == 1310726 {
			procs[a.Id.Tag] = a.ProcsAvg
		}
	}
	gained, used, expired, refreshed := procs[0], procs[1], procs[2], procs[3]
	var landed int32
	for _, a := range um.Actions {
		if a.Id.GetSpellId() == 14271 {
			for _, tm := range a.Targets {
				landed += tm.Casts // a bite that misses spends it too
			}
		}
	}
	if gained == 0 || landed == 0 {
		t.Fatalf("no Expose Prey or no Mongoose Bite (procs %.2f, landed %d); the check needs both", gained, landed)
	}
	t.Logf("Expose Prey a fight: %.2f procs, %.2f used, %.2f expired unused, %.2f refreshed", gained, used, expired, refreshed)
	if diff := used*iterations - float64(landed); diff > 0.01 || diff < -0.01 {
		t.Errorf("used %.2f a fight, Mongoose Bites %.2f a fight; want equal", used, float64(landed)/iterations)
	}
	if left := gained - used - expired; left < -0.001 || left > 1.001 {
		t.Errorf("procs %.2f, used %.2f, expired %.2f, refreshed %.2f a fight: %.2f unaccounted, want 0 to 1 (still up at the end)", gained, used, expired, refreshed, left)
	}
}

// The shipped weaves are tuned for melee swings leaving the Auto Shot timer alone (the beta's reset is off by default):
// the hard weave stays for a swing landing up to 0.6 s after it would have to head back (it beats leaving a quarter
// second early, the reset-era setting), and Arc Auto Weave reads the MH swing timer and the weapon's speed. Each on the
// race its builds ship: the Auto Weave on Night Elf, the hard weave on Human (the Mythic builds).
func TestWeavesTunedWithoutTheMeleeReset(t *testing.T) {
	dps := func(apl []byte, gear string, moveMs int32, race proto.Race) float64 {
		rot := &proto.APLRotation{}
		if err := protojson.Unmarshal(apl, rot); err != nil {
			t.Fatalf("rotation: %v", err)
		}
		opts := &proto.Player_Hunter{Hunter: googleProto.Clone(abOptions.Hunter).(*proto.Hunter)}
		opts.Hunter.Options.MoveDelayMs = moveMs
		opts.Hunter.Options.MeleeResetsAutoShot = false
		player := core.WithSpec(&proto.Player{
			Race:               race,
			Class:              proto.Class_ClassHunter,
			Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", gear).GearSet,
			TalentsString:      "-00531500113-500230230050122151",
			Consumes:           abConsumes,
			Buffs:              abIndividualBuffs,
			Rotation:           rot,
			DistanceFromTarget: 11,
		}, opts)
		res := core.RunRaidSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, abDebuffs),
			Encounter:  abEncounter(),
			SimOptions: &proto.SimOptions{Iterations: 1000, RandomSeed: 20261021},
		})
		if res.Error != nil {
			t.Fatalf("sim error: %s", res.Error.Message)
		}
		return res.RaidMetrics.Dps.Avg
	}
	read := func(file string) []byte {
		raw, err := os.ReadFile("../../ui/hunter/apls/" + file)
		if err != nil {
			t.Fatalf("rotation file: %v", err)
		}
		return raw
	}

	hard := read("arc_weave_2h.apl.json")
	late := regexp.MustCompile(`"val": "-0.6s"`)
	if len(late.FindAll(hard, -1)) == 0 {
		t.Fatalf("no late swing hold in the shipped hard weave")
	}
	// the 2H gain is small since Auto Shot waits 0.6 s after a move, and dual wield's is about 1 since melee abilities
	// stop at the 5 yd edge mid-walk, so both only must not lose
	for gear, gain := range map[string]float64{"2h_weave": 0, "surv_dw": 0} {
		// the site's move reaction: the 5 yd edge wins with it, not with a 10 ms one
		with, early := dps(hard, gear, 100, proto.Race_RaceHuman), dps(late.ReplaceAll(hard, []byte(`"val": "0.25s"`)), gear, 100, proto.Race_RaceHuman)
		t.Logf("hard weave, %s: %.1f DPS staying for a late swing, %.1f leaving a quarter second early", gear, with, early)
		if with < early+gain {
			t.Errorf("hard weave, %s: %.1f DPS with the late swing hold, %.1f without; want at least %.0f more", gear, with, early, gain)
		}
	}

	light := read("arc_light_weave_2h_upd3.apl.json")
	one := func(what string) {
		if n := strings.Count(string(light), what); n != 1 {
			t.Fatalf("%s %d times in the shipped Auto Weave, want once", what, n)
		}
	}
	// the exit sits on Move out and on Wing Clip's copy of its condition; both of its branches read the speed test
	times := func(what string, want int) {
		if n := strings.Count(string(light), what); n != want {
			t.Fatalf("%s %d times in the shipped Auto Weave, want %d", what, n, want)
		}
	}
	const entry, slowAt, freeAt = `"val": "1.5s"`, `"val": "500ms"`, `"val": "0.7s"`
	one(entry)
	times(slowAt, 4)
	times(freeAt, 2)
	spot := regexp.MustCompile(`("rangeFromTarget":\s*\{\s*"const":\s*\{\s*"val":\s*")5(")`)
	if n := len(spot.FindAll(light, -1)); n != 1 {
		t.Fatalf("a move to 5 yd %d times in the shipped Auto Weave, want once", n)
	}
	// the off-hand hold, in the walk out and in the Wing Clip that shares its condition
	ohHold := regexp.MustCompile(`("autoType":\s*"OffHand"\s*\}\s*\},\s*"rhs":\s*\{\s*"const":\s*\{\s*"val":\s*")0\.8s(")`)
	if n := len(ohHold.FindAll(light, -1)); n != 2 {
		t.Fatalf("the off-hand hold %d times in the shipped Auto Weave, want twice", n)
	}
	// dual wield skips the walk-out Wing Clip while Arcane Shot or Multi-Shot is ready; never ready brings the clip back
	shotReady := regexp.MustCompile(`("spellTimeToReady":\s*\{\s*"spellId":\s*\{\s*"spellId":\s*(?:14287|2643),\s*"rank":\s*\d+\s*\}\s*\}\s*\},\s*"rhs":\s*\{\s*"const":\s*\{\s*"val":\s*")0s(")`)
	if n := len(shotReady.FindAll(light, -1)); n != 2 {
		t.Fatalf("the shot-ready test %d times in the shipped Auto Weave, want twice (Arcane Shot, Multi-Shot)", n)
	}
	clipAlways := shotReady.ReplaceAll(light, []byte("${1}-9s${2}"))
	swap := func(from, to string) []byte { return []byte(strings.ReplaceAll(string(light), from, to)) }
	spot4 := spot.ReplaceAll(light, []byte("${1}4${2}"))
	oldExit := []byte(strings.NewReplacer(slowAt, `"val": "9s"`, freeAt, `"val": "9s"`).Replace(string(light)))
	// with the swing retry and the rotation looked at through the GCD (shipped: 2H in at 1.5 s vs every shot +9, vs 1 s
	// +94; DW off-hand hold +71). The 5 yd spot keeps a dual-wield trip inside one swing for a slower player: at a 300 ms
	// move reaction DW +28 over 4 yd, 2H -1. The exit by weapon speed: a main hand half a second slower than the Auto Shot
	// waits for one strike, a faster one uses a ready strike while the next shot is not due
	for _, c := range []struct {
		gear, what string
		other      []byte
		gain       float64
		moveMs     int32
	}{
		{"2h_weave", "in after every Auto Shot", swap(entry, `"val": "9s"`), 0, 100},
		{"2h_weave", "in at 1 s left on the MH timer", swap(entry, `"val": "1s"`), 30, 100},
		// worth less since a walk out still swings inside 5 yd
		{"surv_dw", "no off-hand hold", ohHold.ReplaceAll(light, []byte("${1}0s${2}")), 5, 100},
		{"surv_dw", "a 4 yd melee spot, 300 ms move reaction", spot4, 15, 300},
		{"2h_weave", "a 4 yd melee spot, 300 ms move reaction", spot4, -5, 300},
		{"2h_weave", "the strike wait on every weapon", swap(slowAt, `"val": "-9s"`), 5, 100},
		{"2h_weave", "out as soon as the swing lands", oldExit, 3, 100},
		{"2h_weave_reaper", "out as soon as the swing lands", oldExit, 5, 100},
		{"2h_weave_reaper", "no strike wait", swap(slowAt, `"val": "9s"`), 0, 100},
		{"surv_dw", "Wing Clip while a shot is ready", clipAlways, 3, 100},
		{"surv_dw_no_teebu", "Wing Clip while a shot is ready", clipAlways, 3, 100},
	} {
		shipped, other := dps(light, c.gear, c.moveMs, proto.Race_RaceNightElf), dps(c.other, c.gear, c.moveMs, proto.Race_RaceNightElf)
		t.Logf("Auto Weave, %s: %.1f DPS shipped, %.1f with %s", c.gear, shipped, other, c.what)
		if shipped < other+c.gain {
			t.Errorf("Auto Weave, %s: %.1f DPS shipped, %.1f with %s; want %+.0f or more", c.gear, shipped, other, c.what, c.gain)
		}
	}
	// dual wield never waits for an ability on the way out: it plays the exit as it was before the weapon speed rules
	for _, gear := range []string{"surv_dw", "surv_dw_no_teebu"} {
		if shipped, old := dps(light, gear, 100, proto.Race_RaceNightElf), dps(oldExit, gear, 100, proto.Race_RaceNightElf); shipped != old {
			t.Errorf("Auto Weave, %s: %.2f DPS shipped, %.2f out as soon as the swing lands; dual wield should play the same", gear, shipped, old)
		}
	}
	// a two-hander keeps its Wing Clip
	for _, gear := range []string{"2h_weave", "2h_weave_reaper"} {
		if shipped, clip := dps(light, gear, 100, proto.Race_RaceNightElf), dps(clipAlways, gear, 100, proto.Race_RaceNightElf); shipped != clip {
			t.Errorf("Auto Weave, %s: %.2f DPS shipped, %.2f with Wing Clip while a shot is ready; a two-hander should play the same", gear, shipped, clip)
		}
	}
}

// With a move reaction, the walk back out can be decided on the way in, once inside 5 yd: the reaction runs during the
// rest of the walk, and the walk out starts as the hunter reaches the 3 yd spot (a player turns around there). Deciding
// only on arrival added the whole reaction after it.
func TestWalkOutDecidedOnTheWayIn(t *testing.T) {
	const aplJSON = `{"type":"TypeAPL","priorityList":[` +
		`{"action":{"condition":{"cmp":{"op":"OpGe","lhs":{"distanceToTarget":{}},"rhs":{"const":{"val":"11"}}}},"move":{"rangeFromTarget":{"const":{"val":"3"}}}}},` +
		`{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"distanceToTarget":{}},"rhs":{"const":{"val":"5"}}}},"move":{"rangeFromTarget":{"const":{"val":"12"}}}}}]}`
	rot := &proto.APLRotation{}
	if err := protojson.Unmarshal([]byte(aplJSON), rot); err != nil {
		t.Fatalf("rotation: %v", err)
	}
	opts := &proto.Player_Hunter{Hunter: googleProto.Clone(abOptions.Hunter).(*proto.Hunter)}
	opts.Hunter.Options.MoveDelayMs = 200
	player := core.WithSpec(&proto.Player{
		Race:               proto.Race_RaceHuman,
		Class:              proto.Class_ClassHunter,
		Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
		TalentsString:      patchTestTalents,
		Consumes:           abConsumes,
		Buffs:              abIndividualBuffs,
		Rotation:           rot,
		DistanceFromTarget: 12,
	}, opts)
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, abDebuffs),
		Encounter:  abEncounter(),
		SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 20261020, Debug: true},
	})
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
	}
	stamp := regexp.MustCompile(`^\[([0-9.]+)\]`)
	var arrived, walkOut float64 = -1, -1
	starts := 0
	for _, line := range strings.Split(res.Logs, "\n") {
		m := stamp.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		at, _ := strconv.ParseFloat(m[1], 64)
		if arrived < 0 && strings.Contains(line, "{OtherID: 14} stacks: 4 --> 3") {
			arrived = at
		}
		if strings.Contains(line, "Aura gained: {OtherID: 14}") {
			if starts++; starts == 2 {
				walkOut = at
				break
			}
		}
	}
	if arrived < 0 || walkOut < 0 {
		t.Fatalf("no walk in to 3 yd and back out logged (arrived %.2f, walk out %.2f)", arrived, walkOut)
	}
	if gap := walkOut - arrived; gap > 0.011 {
		t.Errorf("walk out started %.2f s after reaching 3 yd; want on arrival (the reaction ran during the walk in)", gap)
	}
}

// The end of a walk is a decision point: walking in from 11 to 3 yd during a busy GCD, the swing retry lands while the
// walk goes on (no new move then), so the walk back has to start on arrival. Without the recheck it waited for the GCD.
func TestMoveAfterArrivalDuringGcd(t *testing.T) {
	const aplJSON = `{"type":"TypeAPL","priorityList":[` +
		`{"action":{"condition":{"cmp":{"op":"OpGe","lhs":{"distanceToTarget":{}},"rhs":{"const":{"val":"11"}}}},"castSpell":{"spellId":{"spellId":14287,"rank":8}}}},` +
		`{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGe","lhs":{"distanceToTarget":{}},"rhs":{"const":{"val":"11"}}}},{"not":{"val":{"gcdIsReady":{}}}}]}},"move":{"rangeFromTarget":{"const":{"val":"3"}}}}},` +
		`{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"distanceToTarget":{}},"rhs":{"const":{"val":"3"}}}},"move":{"rangeFromTarget":{"const":{"val":"12"}}}}}]}`
	rot := &proto.APLRotation{}
	if err := protojson.Unmarshal([]byte(aplJSON), rot); err != nil {
		t.Fatalf("rotation: %v", err)
	}
	player := core.WithSpec(&proto.Player{
		Race:               proto.Race_RaceHuman,
		Class:              proto.Class_ClassHunter,
		Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
		TalentsString:      patchTestTalents,
		Consumes:           abConsumes,
		Buffs:              abIndividualBuffs,
		Rotation:           rot,
		DistanceFromTarget: 11,
	}, abOptions)
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, abDebuffs),
		Encounter:  abEncounter(),
		SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 20261018, DebugFirstIteration: true},
	})
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
	}
	at := func(line string) float64 {
		v, _ := strconv.ParseFloat(strings.TrimPrefix(strings.SplitN(line, "]", 2)[0], "["), 64)
		return v
	}
	shot, arrived, out := -1.0, -1.0, -1.0
	for _, line := range strings.Split(res.Logs, "\n") {
		switch {
		case shot < 0 && at(line) >= 0 && strings.Contains(line, "Completed cast {SpellID: 14287"):
			shot = at(line)
		case shot >= 0 && arrived < 0 && strings.Contains(line, "Aura faded: {OtherID: 14}"):
			arrived = at(line)
		case arrived >= 0 && strings.Contains(line, "Moving to %!s(float64=12)"):
			out = at(line)
		}
		if out >= 0 {
			break
		}
	}
	if shot < 0 || arrived < 0 || out < 0 {
		t.Fatalf("no Arcane Shot, walk in and walk back in the log (shot %.2f, arrived %.2f, out %.2f)", shot, arrived, out)
	}
	if arrived >= shot+1.45 {
		t.Fatalf("arrived %.2f s after the shot, not inside its GCD; the walk is too long for this check", arrived-shot)
	}
	if out > arrived+0.01 {
		t.Errorf("walked back at %.2f s, %.2f s after arriving at %.2f s; want it on arrival (the GCD ends at %.2f s)", out, out-arrived, arrived, shot+1.5)
	}
}

// Distances are whole yards: move spots of 4.6 and 11.4 yd and a start at 11.5 yd round to 5, 11 and 12, so a weave on
// them keeps shooting and swinging all fight. Unrounded, the first walk never reached 4.6 and froze the rotation.
func TestWholeYardMoveSpots(t *testing.T) {
	const aplJSON = `{"type":"TypeAPL","priorityList":[` +
		`{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGe","lhs":{"distanceToTarget":{}},"rhs":{"const":{"val":"11"}}}},{"cmp":{"op":"OpGt","lhs":{"autoTimeToNext":{"autoType":"Ranged"}},"rhs":{"const":{"val":"1.5s"}}}}]}},"move":{"rangeFromTarget":{"const":{"val":"4.6"}}}}},` +
		`{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"distanceToTarget":{}},"rhs":{"const":{"val":"5"}}}},{"cmp":{"op":"OpGt","lhs":{"autoTimeToNext":{"autoType":"MainHand"}},"rhs":{"const":{"val":"0.5s"}}}}]}},"move":{"rangeFromTarget":{"const":{"val":"11.4"}}}}}]}`
	rot := &proto.APLRotation{}
	if err := protojson.Unmarshal([]byte(aplJSON), rot); err != nil {
		t.Fatalf("rotation: %v", err)
	}
	player := core.WithSpec(&proto.Player{
		Race:               proto.Race_RaceHuman,
		Class:              proto.Class_ClassHunter,
		Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
		TalentsString:      patchTestTalents,
		Consumes:           abConsumes,
		Buffs:              abIndividualBuffs,
		Rotation:           rot,
		DistanceFromTarget: 11.5,
	}, abOptions)
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, abDebuffs),
		Encounter:  abEncounter(),
		SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 20261018},
	})
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
	}
	var swings, shots int32
	for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
		for _, tm := range a.Targets {
			switch a.Id.GetOtherId() {
			case proto.OtherAction_OtherActionAttack:
				swings += tm.Casts
			case proto.OtherAction_OtherActionShoot:
				shots += tm.Casts
			}
		}
	}
	if swings < 10 || shots < 10 {
		t.Errorf("%d melee swings and %d Auto Shots weaving between 4.6 and 11.4 yd; want a full fight of both (10+ each)", swings, shots)
	}
}

// A walk to 6 yd (no swing on the way) that ends during a busy GCD walks back within one 50 ms look of arriving: moves
// are looked at through the GCD (ForeverMovePoll). Rotations read "swung this trip" from the swing itself
// (TestTimeSinceLastSwing), so a look on arrival no longer sends them out before their swing.
func TestNoArrivalRecheckWithoutSwing(t *testing.T) {
	const aplJSON = `{"type":"TypeAPL","priorityList":[` +
		`{"action":{"condition":{"cmp":{"op":"OpGe","lhs":{"distanceToTarget":{}},"rhs":{"const":{"val":"11"}}}},"castSpell":{"spellId":{"spellId":14287,"rank":8}}}},` +
		`{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGe","lhs":{"distanceToTarget":{}},"rhs":{"const":{"val":"11"}}}},{"not":{"val":{"gcdIsReady":{}}}}]}},"move":{"rangeFromTarget":{"const":{"val":"6"}}}}},` +
		`{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"distanceToTarget":{}},"rhs":{"const":{"val":"6"}}}},"move":{"rangeFromTarget":{"const":{"val":"12"}}}}}]}`
	rot := &proto.APLRotation{}
	if err := protojson.Unmarshal([]byte(aplJSON), rot); err != nil {
		t.Fatalf("rotation: %v", err)
	}
	player := core.WithSpec(&proto.Player{
		Race:               proto.Race_RaceHuman,
		Class:              proto.Class_ClassHunter,
		Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
		TalentsString:      patchTestTalents,
		Consumes:           abConsumes,
		Buffs:              abIndividualBuffs,
		Rotation:           rot,
		DistanceFromTarget: 12,
	}, abOptions)
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, abDebuffs),
		Encounter:  abEncounter(),
		SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 20261018, DebugFirstIteration: true},
	})
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
	}
	at := func(line string) float64 {
		v, _ := strconv.ParseFloat(strings.TrimPrefix(strings.SplitN(line, "]", 2)[0], "["), 64)
		return v
	}
	shot, arrived, out := -1.0, -1.0, -1.0
	for _, line := range strings.Split(res.Logs, "\n") {
		switch {
		case shot < 0 && at(line) >= 0 && strings.Contains(line, "Completed cast {SpellID: 14287"):
			shot = at(line)
		case shot >= 0 && arrived < 0 && strings.Contains(line, "Aura faded: {OtherID: 14}"):
			arrived = at(line)
		case arrived >= 0 && strings.Contains(line, "Moving to %!s(float64=12)"):
			out = at(line)
		}
		if out >= 0 {
			break
		}
	}
	if shot < 0 || arrived < 0 || out < 0 {
		t.Fatalf("no Arcane Shot, walk in and walk back in the log (shot %.2f, arrived %.2f, out %.2f)", shot, arrived, out)
	}
	if arrived >= shot+1.4 {
		t.Fatalf("arrived %.2f s after the shot, not inside its GCD; the walk is too long for this check", arrived-shot)
	}
	if out > arrived+0.051 {
		t.Errorf("walked back at %.2f s, %.2f s after arriving at %.2f s; want within one 50 ms look (the GCD ends at %.2f s)", out, out-arrived, arrived, shot+1.5)
	}
}

// An Auto Shot is a decision point like a melee swing: a move into melee that only goes while the GCD is busy starts in
// the instant after a shot, on a GCD already running (not one just started in the same pass, and past the pull). Without
// the recheck the rotation is not looked at then, so the hunter never walks in. A rotation saved with the old
// recheckAfterAutoShot field still loads and behaves the same either way.
func TestMoveAfterShotDuringGcd(t *testing.T) {
	const aplJSON = `{"type":"TypeAPL","priorityList":[{"action":{"castSpell":{"spellId":{"spellId":14287,"rank":8}}}},{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGe","lhs":{"currentTime":{}},"rhs":{"const":{"val":"1s"}}}},{"cmp":{"op":"OpLt","lhs":{"gcdTimeToReady":{}},"rhs":{"const":{"val":"1.4s"}}}},{"not":{"val":{"gcdIsReady":{}}}},{"cmp":{"op":"OpLe","lhs":{"math":{"op":"OpSub","lhs":{"autoSwingTime":{"autoType":"Ranged"}},"rhs":{"autoTimeToNext":{"autoType":"Ranged"}}}},"rhs":{"const":{"val":"0.01s"}}}}]}},"move":{"rangeFromTarget":{"const":{"val":"2"}}}}}]}`
	meleeSwings := func(savedSwitch bool) int32 {
		rot := &proto.APLRotation{}
		if err := protojson.Unmarshal([]byte(aplJSON), rot); err != nil {
			t.Fatalf("rotation: %v", err)
		}
		rot.RecheckAfterAutoShot = savedSwitch
		player := core.WithSpec(&proto.Player{
			Race:               proto.Race_RaceNightElf,
			Class:              proto.Class_ClassHunter,
			Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
			TalentsString:      patchTestTalents,
			Consumes:           abConsumes,
			Buffs:              abIndividualBuffs,
			Rotation:           rot,
			DistanceFromTarget: 12,
		}, abOptions)
		res := core.RunRaidSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, abDebuffs),
			Encounter:  abEncounter(),
			SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 20261017},
		})
		if res.Error != nil {
			t.Fatalf("sim error: %s", res.Error.Message)
		}
		var swings int32
		for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
			if a.Id.GetOtherId() == proto.OtherAction_OtherActionAttack {
				for _, tm := range a.Targets {
					swings += tm.Casts
				}
			}
		}
		return swings
	}
	for _, saved := range []bool{false, true} {
		if n := meleeSwings(saved); n == 0 {
			t.Errorf("no melee swings (saved switch %v): the move into melee should start in the instant after an Auto Shot while the GCD is busy", saved)
		}
	}
}

// A trap lands about a yard in front of you, so it catches a target from 6 yd (in game: hits from 6, not from 8): from
// 7 yd it is not dropped and does nothing.
func TestTrapReach(t *testing.T) {
	const aplJSON = `{"type":"TypeAPL","priorityList":[{"action":{"castSpell":{"spellId":{"spellId":14305}}}}]}`
	trapDamage := func(distance float64) (int32, float64) {
		rot := &proto.APLRotation{}
		if err := protojson.Unmarshal([]byte(aplJSON), rot); err != nil {
			t.Fatalf("rotation: %v", err)
		}
		player := core.WithSpec(&proto.Player{
			Race:               proto.Race_RaceHuman,
			Class:              proto.Class_ClassHunter,
			Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
			TalentsString:      patchTestTalents,
			Consumes:           abConsumes,
			Buffs:              abIndividualBuffs,
			Rotation:           rot,
			DistanceFromTarget: distance,
		}, abOptions)
		res := core.RunRaidSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, abDebuffs),
			Encounter:  abEncounter(),
			SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 20261017},
		})
		if res.Error != nil {
			t.Fatalf("sim error: %s", res.Error.Message)
		}
		var casts int32
		var damage float64
		for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
			if a.Id.GetSpellId() == 14305 {
				for _, tm := range a.Targets {
					casts += tm.Casts
					damage += tm.Damage
				}
			}
		}
		return casts, damage
	}
	for _, d := range []float64{4, 6} {
		if casts, damage := trapDamage(d); casts == 0 || damage == 0 {
			t.Errorf("from %.0f yd: %d traps, %.0f damage; want them to land", d, casts, damage)
		}
	}
	if casts, damage := trapDamage(7); casts != 0 || damage != 0 {
		t.Errorf("from 7 yd: %d traps, %.0f damage; want none (6 yd at most)", casts, damage)
	}
}

// A walk in that starts at the pull, from inside the first Auto Shot (after the pre-pull walk out to 11 yd and Aimed Shot),
// is under way when the pull turns the attacks on: they wait for the walk to end. Turned on mid-walk, Auto Shot was put
// back from 11 yd and kept firing from inside melee range, and the melee swings never started.
func TestMoveAtPull(t *testing.T) {
	const aplJSON = `{"type":"TypeAPL","prepullActions":[{"action":{"move":{"rangeFromTarget":{"const":{"val":"11"}}}},"doAtValue":{"const":{"val":"-5s"}}},{"action":{"castSpell":{"spellId":{"spellId":20904}}},"doAtValue":{"const":{"val":"-2s"}}}],"priorityList":[{"action":{"move":{"rangeFromTarget":{"const":{"val":"3"}}}}}]}`
	rot := &proto.APLRotation{}
	if err := protojson.Unmarshal([]byte(aplJSON), rot); err != nil {
		t.Fatalf("rotation: %v", err)
	}
	player := core.WithSpec(&proto.Player{
		Race:               proto.Race_RaceHuman,
		Class:              proto.Class_ClassHunter,
		Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
		TalentsString:      patchTestTalents,
		Consumes:           abConsumes,
		Buffs:              abIndividualBuffs,
		Rotation:           rot,
		DistanceFromTarget: 3,
	}, abOptions)
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, abDebuffs),
		Encounter:  abEncounter(),
		SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 20261017},
	})
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
	}
	var swings, shots int32
	for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
		for _, tm := range a.Targets {
			switch a.Id.GetOtherId() {
			case proto.OtherAction_OtherActionAttack:
				swings += tm.Casts
			case proto.OtherAction_OtherActionShoot:
				shots += tm.Casts
			}
		}
	}
	if swings < 20 {
		t.Errorf("%d melee swings standing at 3 yd for the whole fight; want them all (about 30)", swings)
	}
	if shots > 1 {
		t.Errorf("%d Auto Shots from inside melee range; want only the one at the pull", shots)
	}
}

// Any ranged attack takes you out of melee mode, so it unqueues Raptor Strike: a shot (Arcane Shot here) and the Auto
// Shot each cancel a queued Raptor Strike.
func TestRangedAttackUnqueuesRaptorStrike(t *testing.T) {
	queued := func(hunter *Hunter) bool { return hunter.curQueueAura != nil && hunter.curQueueAura.IsActive() }

	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 12})
	hunter.RaptorStrike.Cast(sim, hunter.CurrentTarget)
	if !queued(hunter) {
		t.Fatal("Raptor Strike did not queue at 12 yd")
	}
	arcane := hunter.GetSpell(core.ActionID{SpellID: 14287})
	if arcane == nil {
		t.Fatal("no Arcane Shot")
	}
	arcane.Cast(sim, hunter.CurrentTarget)
	if queued(hunter) {
		t.Error("Raptor Strike is still queued after Arcane Shot")
	}

	sim, hunter = newPatchTestHunter(t, patchTestSetup{gear: true, distance: 12})
	hunter.RaptorStrike.Cast(sim, hunter.CurrentTarget)
	hunter.AutoAttacks.RangedAuto().Cast(sim, hunter.CurrentTarget)
	if queued(hunter) {
		t.Error("Raptor Strike is still queued after an Auto Shot")
	}
}

// A walk into melee with the main hand ready swings on the first half-second retry after the server sees you there,
// counted from the later of the timer running out and the walk turning the attack on (an Auto Shot turned it off);
// in range before the timer runs out, it swings as it runs out, and a walk out leaves the timer alone (probe walk-ins
// and the beta combat logs).
func TestSwingRetryOnArrival(t *testing.T) {
	const beat = 500 * time.Millisecond
	walk := func(sim *core.Simulation, hunter *Hunter, to float64) (started, arrived time.Duration) {
		hunter.MoveTo(to, sim)
		for guard := 0; !hunter.IsMoving() && guard < 100000; guard++ {
			sim.Step()
		}
		started = sim.CurrentTime
		for guard := 0; hunter.IsMoving() && guard < 100000; guard++ {
			sim.Step()
		}
		return started, sim.CurrentTime
	}

	// timer long run out, attack turned on by the walk: the retries count from the walk's start
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 12})
	sim.PrePull()
	for guard := 0; sim.CurrentTime < time.Second*4 && guard < 100000; guard++ {
		sim.Step()
	}
	started, arrived := walk(sim, hunter, 5)
	seen := arrived + 50*time.Millisecond
	want := started + (seen-started+beat-1)/beat*beat
	if got := hunter.AutoAttacks.MainhandSwingAt(); got != want {
		t.Errorf("walk in from 12 yd (started %.3f, arrived %.3f): main hand swings at %.3f, want the retry at %.3f",
			started.Seconds(), arrived.Seconds(), got.Seconds(), want.Seconds())
	}

	// in range again before the timer runs out (no Auto Shot, the attack stayed on): it swings as the timer runs out
	sim, hunter = newPatchTestHunter(t, patchTestSetup{gear: true, distance: 5})
	sim.PrePull()
	for guard := 0; sim.CurrentTime < time.Second*2 && guard < 100000; guard++ {
		sim.Step()
	}
	end := hunter.AutoAttacks.MainhandSwingAt()
	walk(sim, hunter, 9)
	_, back := walk(sim, hunter, 5)
	if back >= end {
		t.Fatalf("test setup: back in range at %.3f, after the timer ran out at %.3f", back.Seconds(), end.Seconds())
	}
	if got := hunter.AutoAttacks.MainhandSwingAt(); got != end {
		t.Errorf("back in range at %.3f before the timer ran out at %.3f: main hand swings at %.3f, want %.3f",
			back.Seconds(), end.Seconds(), got.Seconds(), end.Seconds())
	}

	// in range just after the timer ran out, attack still on: the next retry from the timer's end
	sim, hunter = newPatchTestHunter(t, patchTestSetup{gear: true, distance: 5})
	sim.PrePull()
	for guard := 0; sim.CurrentTime < time.Second*2 && guard < 100000; guard++ {
		sim.Step()
	}
	end = hunter.AutoAttacks.MainhandSwingAt()
	walk(sim, hunter, 9)
	for guard := 0; sim.CurrentTime < end+300*time.Millisecond && guard < 100000; guard++ {
		sim.Step()
	}
	_, back = walk(sim, hunter, 5)
	seen = back + 50*time.Millisecond
	want = end + (seen-end+beat-1)/beat*beat
	if got := hunter.AutoAttacks.MainhandSwingAt(); got != want {
		t.Errorf("back in range at %.3f after the timer ran out at %.3f: main hand swings at %.3f, want the retry at %.3f",
			back.Seconds(), end.Seconds(), got.Seconds(), want.Seconds())
	}
}

// The "Shots delay Auto Shot" option (forever-bugs #43): a shot cast during an Auto Shot cycle pushes that cycle's
// Auto Shot back 0.2 s, once a cycle, every cycle; off (the default) it does not.
func TestShotsDelayAutoShot(t *testing.T) {
	type dues struct{ base, first, sameCycle, nextBase, next time.Duration }
	run := func(on bool) (d dues) {
		sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 12, shotsDelay: on})
		sim.PrePull()
		for guard := 0; sim.CurrentTime < time.Second*2 && guard < 100000; guard++ {
			sim.Step()
		}
		d.base = hunter.AutoAttacks.NextRangedAttackAt()
		hunter.ArcaneShot.Cast(sim, hunter.CurrentTarget)
		d.first = hunter.AutoAttacks.NextRangedAttackAt()
		// a second shot in the same cycle (its GCD cleared for the check) adds nothing
		hunter.SetGCDTimer(sim, sim.CurrentTime)
		hunter.SerpentSting.Cast(sim, hunter.CurrentTarget)
		d.sameCycle = hunter.AutoAttacks.NextRangedAttackAt()
		// that Auto Shot fires: a shot in the next cycle pushes it back again
		for guard := 0; sim.CurrentTime <= d.sameCycle && guard < 100000; guard++ {
			sim.Step()
		}
		d.nextBase = hunter.AutoAttacks.NextRangedAttackAt()
		hunter.SetGCDTimer(sim, sim.CurrentTime)
		hunter.SerpentSting.Cast(sim, hunter.CurrentTarget)
		d.next = hunter.AutoAttacks.NextRangedAttackAt()
		return
	}
	on := run(true)
	if on.first-on.base != shotAutoShotDelay || on.sameCycle != on.first || on.next-on.nextBase != shotAutoShotDelay {
		t.Errorf("option on: Auto Shot due %.3f, %.3f after an Arcane Shot, %.3f after a second shot that cycle; next cycle %.3f, %.3f after a shot; want +0.2 s once each cycle",
			on.base.Seconds(), on.first.Seconds(), on.sameCycle.Seconds(), on.nextBase.Seconds(), on.next.Seconds())
	}
	off := run(false)
	if off.first != off.base || off.sameCycle != off.base || off.next != off.nextBase {
		t.Errorf("option off: Auto Shot due %.3f moved to %.3f / %.3f, next cycle %.3f moved to %.3f", off.base.Seconds(), off.first.Seconds(), off.sameCycle.Seconds(), off.nextBase.Seconds(), off.next.Seconds())
	}
}

// Aimed Shot and Multi-Shot casts take all ranged haste, the quiver's included (beta logs: 1.81 s and 0.45 s with an
// 11% quiver; 0.32 s Multi-Shots with Rapid Fire on top).
func TestShotCastsTakeRangedHaste(t *testing.T) {
	for _, c := range []struct {
		name  string
		spell func(*Hunter) *core.Spell
		base  time.Duration
	}{{"Multi-Shot", func(h *Hunter) *core.Spell { return h.MultiShot }, 500 * time.Millisecond},
		{"Aimed Shot", func(h *Hunter) *core.Spell { return h.AimedShot }, 2000 * time.Millisecond}} {
		// a hunter each: the two shots share a cooldown
		sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 12, quiver: proto.Hunter_Options_Speed15})
		if hunter.quiverBonus <= 1 {
			t.Fatalf("test setup: no quiver bonus (%.2f)", hunter.quiverBonus)
		}
		sim.PrePull()
		for guard := 0; sim.CurrentTime < time.Second && guard < 100000; guard++ {
			sim.Step()
		}
		haste := hunter.PseudoStats.RangedSpeedMultiplier
		hunter.SetGCDTimer(sim, sim.CurrentTime)
		if !c.spell(hunter).Cast(sim, hunter.CurrentTarget) {
			t.Fatalf("%s did not cast", c.name)
		}
		want := time.Duration(float64(c.base) / haste)
		if got := hunter.Hardcast.Expires - sim.CurrentTime; got != want {
			t.Errorf("%s casts in %.3f s, want %.3f s (%.3f s at %.2fx ranged haste)", c.name, got.Seconds(), want.Seconds(), c.base.Seconds(), haste)
		}
	}
}

// A pre-pull Aimed Shot that ends as the pull's Auto Shot fires completes then: the shot's decision point must not
// start a Serpent Sting first, whose GCD pushed the completion to 1.5 s.
func TestPrepullAimedShotCompletesAtPull(t *testing.T) {
	const aplJSON = `{"type":"TypeAPL","prepullActions":[{"action":{"castSpell":{"spellId":{"spellId":20904}}},"doAtValue":{"const":{"val":"-2s"}}}],` +
		`"priorityList":[{"action":{"castSpell":{"spellId":{"spellId":25295,"rank":9}}}}]}`
	rot := &proto.APLRotation{}
	if err := protojson.Unmarshal([]byte(aplJSON), rot); err != nil {
		t.Fatalf("rotation: %v", err)
	}
	// no quiver: the 2 s cast from -2 s ends exactly as the pull's Auto Shot fires
	noQuiver := &proto.Player_Hunter{Hunter: googleProto.Clone(abOptions.Hunter).(*proto.Hunter)}
	noQuiver.Hunter.Options.QuiverBonus = proto.Hunter_Options_QuiverNone
	player := core.WithSpec(&proto.Player{
		Race:               proto.Race_RaceHuman,
		Class:              proto.Class_ClassHunter,
		Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
		TalentsString:      patchTestTalents,
		Consumes:           abConsumes,
		Buffs:              abIndividualBuffs,
		Rotation:           rot,
		DistanceFromTarget: 12,
	}, noQuiver)
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, abDebuffs),
		Encounter:  abEncounter(),
		SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 20261023, DebugFirstIteration: true},
	})
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
	}
	for _, line := range strings.Split(res.Logs, "\n") {
		if strings.Contains(line, "Completed cast {SpellID: 20904") {
			at, _ := strconv.ParseFloat(strings.TrimPrefix(strings.SplitN(line, "]", 2)[0], "["), 64)
			if at > 0.001 {
				t.Errorf("the pre-pull Aimed Shot completed at %.2f s, want at the pull (0.00)", at)
			}
			return
		}
	}
	t.Fatal("no completed Aimed Shot in the log")
}

// Curse of the Elements lowers each of its five schools' resistance by 75 once (it used to take 75 off all five schools
// for each school it covers: 375 each).
func TestCurseOfElementsResistance(t *testing.T) {
	_, plain := newPatchTestHunter(t, patchTestSetup{gear: true, debuffs: &proto.Debuffs{}})
	sim, cursed := newPatchTestHunter(t, patchTestSetup{gear: true, debuffs: &proto.Debuffs{CurseOfElements: true}})
	sim.PrePull()
	for guard := 0; sim.CurrentTime < time.Second && guard < 100000; guard++ {
		sim.Step()
	}
	for _, s := range []stats.Stat{stats.ArcaneResistance, stats.FireResistance, stats.FrostResistance, stats.NatureResistance, stats.ShadowResistance} {
		if got := cursed.CurrentTarget.GetStat(s) - plain.CurrentTarget.GetStat(s); math.Abs(got+75) > 1e-9 {
			t.Errorf("Curse of the Elements changes %v by %.0f, want -75", s, got)
		}
	}
}

// Projectile speeds from the client's SpellMisc: Auto Shot, Arcane Shot, Serpent Sting and Aimed Shot 40 yd/s,
// Multi-Shot 30.
func TestShotMissileSpeeds(t *testing.T) {
	_, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 12})
	for _, c := range []struct {
		spell *core.Spell
		speed float64
	}{{hunter.AutoAttacks.RangedAuto(), 40}, {hunter.ArcaneShot, 40}, {hunter.SerpentSting, 40}, {hunter.AimedShot, 40}, {hunter.MultiShot, 30}} {
		if c.spell.MissileSpeed != c.speed {
			t.Errorf("%s flies at %.0f yd/s, want %.0f", c.spell.ActionID, c.spell.MissileSpeed, c.speed)
		}
	}
}

// The timeline's wait lines tell the truth: in a weave fight there are swing retry and Auto Shot windup waits, each
// retry wait ends exactly when the main hand swings (one cut short by a walk out is logged as left), and each windup runs
// the full 0.6 s from its stop.
func TestSwingWaitLogs(t *testing.T) {
	data, err := os.ReadFile("../../../builds/variants/ship_2h_auto_weave.json")
	if err != nil {
		t.Fatal(err)
	}
	settings := &proto.IndividualSimSettings{}
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(data, settings); err != nil {
		t.Fatal(err)
	}
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(settings.Player, settings.PartyBuffs, settings.RaidBuffs, settings.Debuffs),
		Encounter:  settings.Encounter,
		SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 7, DebugFirstIteration: true},
	})
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
	}
	waitRe := regexp.MustCompile(`Swing wait \{OtherID: (\d+)(?:, Tag: (\d+))?\} (\w+) from (-?[0-9.]+)s to (-?[0-9.]+)s(?:, steps from (-?[0-9.]+)s)?(?:, stopped at (-?[0-9.]+)s)?`)
	swingRe := regexp.MustCompile(`^\[(-?[0-9.]+)\] \[Player \(#1\)\] Swing timer \{OtherID: 7, Tag: 1\} started`)
	swings := map[string]bool{}
	var retries, windups []string
	for _, line := range strings.Split(res.Logs, "\n") {
		if m := swingRe.FindStringSubmatch(line); m != nil {
			v, _ := strconv.ParseFloat(m[1], 64)
			swings[strconv.FormatFloat(v, 'f', 2, 64)] = true
		}
		if m := waitRe.FindStringSubmatch(line); m != nil {
			switch m[3] {
			case "retry":
				retries = append(retries, m[5])
			case "windup":
				windups = append(windups, m[5])
				to, _ := strconv.ParseFloat(m[5], 64)
				stop, _ := strconv.ParseFloat(m[7], 64)
				if math.Abs(to-stop-core.AutoShotStandStill.Seconds()) > 0.0015 {
					t.Errorf("a windup runs from the stop at %.3fs to %.3fs, want %.1f s", stop, to, core.AutoShotStandStill.Seconds())
				}
			}
		}
	}
	if len(retries) == 0 || len(windups) == 0 {
		t.Fatalf("%d swing retry and %d windup waits in a weave fight, want some of each", len(retries), len(windups))
	}
	// the log stamps lines to 0.01 s: a swing within 6 ms of the wait's end is the one it waited for
	for _, to := range retries {
		v, _ := strconv.ParseFloat(to, 64)
		found := false
		for k := range swings {
			if s, _ := strconv.ParseFloat(k, 64); math.Abs(s-v) <= 0.006 {
				found = true
			}
		}
		if !found {
			t.Errorf("a swing retry wait ends at %ss, but the main hand did not swing then", to)
		}
	}
}

// "Time since last swing" counts from the swing itself. Swing time minus time to the next swing (how the rotations
// read it before) comes out short while a swing retry holds the next swing, so a trip that has not swung looked swung.
func TestTimeSinceLastSwing(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 5})
	sim.PrePull()
	for guard := 0; sim.CurrentTime < time.Second*2 && guard < 100000; guard++ {
		sim.Step()
	}
	// a swing has landed: walk out, wait out the timer at range, walk back in for a retry
	end := hunter.AutoAttacks.MainhandSwingAt()
	lastSwing := end - hunter.AutoAttacks.MainhandSwingSpeed()
	hunter.MoveTo(9, sim)
	for guard := 0; sim.CurrentTime < end+300*time.Millisecond && guard < 100000; guard++ {
		sim.Step()
	}
	hunter.MoveTo(5, sim)
	for guard := 0; (hunter.IsMoving() || sim.CurrentTime == 0) && guard < 100000; guard++ {
		sim.Step()
	}
	if hunter.AutoAttacks.MainhandSwingAt() <= sim.CurrentTime {
		t.Fatal("test setup: no swing retry pending on arrival")
	}
	since := hunter.AutoAttacks.TimeSinceLastAttack(sim, proto.APLValueAutoSwingTime_MainHand)
	if want := sim.CurrentTime - lastSwing; since != want {
		t.Errorf("on arrival, time since the last swing reads %v, want %v (the swing at %v)", since, want, lastSwing)
	}
	old := hunter.AutoAttacks.MainhandSwingSpeed() - (hunter.AutoAttacks.MainhandSwingAt() - sim.CurrentTime)
	if old >= since {
		t.Fatalf("test setup: the old reading %v is not short of %v here", old, since)
	}
}

// A move is looked at every 50 ms while the GCD rolls, as a player steps in when a timer says so: a walk due 0.7 s into
// an Arcane Shot's GCD starts within 50 ms of it. The engine used to look again only as the GCD ended.
func TestMovePollDuringGcd(t *testing.T) {
	const aplJSON = `{"type":"TypeAPL","priorityList":[` +
		`{"action":{"condition":{"cmp":{"op":"OpGe","lhs":{"distanceToTarget":{}},"rhs":{"const":{"val":"11"}}}},"castSpell":{"spellId":{"spellId":14287,"rank":8}}}},` +
		`{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGe","lhs":{"distanceToTarget":{}},"rhs":{"const":{"val":"11"}}}},{"cmp":{"op":"OpGe","lhs":{"currentTime":{}},"rhs":{"const":{"val":"0.7s"}}}}]}},"move":{"rangeFromTarget":{"const":{"val":"5"}}}}}]}`
	rot := &proto.APLRotation{}
	if err := protojson.Unmarshal([]byte(aplJSON), rot); err != nil {
		t.Fatalf("rotation: %v", err)
	}
	player := core.WithSpec(&proto.Player{
		Race:               proto.Race_RaceHuman,
		Class:              proto.Class_ClassHunter,
		Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
		TalentsString:      patchTestTalents,
		Consumes:           abConsumes,
		Buffs:              abIndividualBuffs,
		Rotation:           rot,
		DistanceFromTarget: 11,
	}, abOptions)
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, abRaidBuffs, abDebuffs),
		Encounter:  abEncounter(),
		SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 20261023, DebugFirstIteration: true},
	})
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
	}
	at := func(line string) float64 {
		v, _ := strconv.ParseFloat(strings.TrimPrefix(strings.SplitN(line, "]", 2)[0], "["), 64)
		return v
	}
	shot, move := -1.0, -1.0
	for _, line := range strings.Split(res.Logs, "\n") {
		switch {
		case shot < 0 && at(line) >= 0 && strings.Contains(line, "Completed cast {SpellID: 14287"):
			shot = at(line)
		case shot >= 0 && strings.Contains(line, "Moving to %!s(float64=5)"):
			move = at(line)
		}
		if move >= 0 {
			break
		}
	}
	if shot < 0 || move < 0 {
		t.Fatalf("no Arcane Shot and walk in logged (shot %.2f, move %.2f)", shot, move)
	}
	if shot > 0.2 {
		t.Fatalf("Arcane Shot at %.2f s: its GCD does not cover 0.7 s", shot)
	}
	if move > 0.75+0.001 {
		t.Errorf("walk in started at %.2f s; want by 0.75 s (due at 0.7 s, the GCD ends at %.2f s)", move, shot+1.5)
	}
}

// A pet swings for the level 60 creature dps (60.45), 20% either side, over its swing time; its crit is its owner's
// with nothing from its own Agility, and its spell hit is 17/8 of the hit it inherits (forever-bugs #394).
func TestPetSwingFromCreatureDps(t *testing.T) {
	_, hunter := newPatchTestHunter(t, patchTestSetup{gear: true})
	mh := hunter.pet.AutoAttacks.MH()
	speed := mh.SwingSpeed
	if want := 60.45 * 0.8 * speed; math.Abs(mh.BaseDamageMin-want) > 1e-6 {
		t.Errorf("pet swing minimum %.2f at %.1f s, want %.2f", mh.BaseDamageMin, speed, want)
	}
	if want := 60.45 * 1.2 * speed; math.Abs(mh.BaseDamageMax-want) > 1e-6 {
		t.Errorf("pet swing maximum %.2f at %.1f s, want %.2f", mh.BaseDamageMax, speed, want)
	}
	if got, want := hunter.pet.GetStat(stats.MeleeCrit), hunter.GetStat(stats.MeleeCrit); math.Abs(got-want) > 1e-9 {
		t.Errorf("pet crit %.2f, want its owner's %.2f", got/core.CritRatingPerCritChance, want/core.CritRatingPerCritChance)
	}
	if got, want := hunter.pet.GetStat(stats.SpellHit), hunter.GetStat(stats.MeleeHit)*17/8; math.Abs(got-want) > 1e-9 {
		t.Errorf("pet spell hit %.2f, want 17/8 of the owner's %.2f hit", got, hunter.GetStat(stats.MeleeHit))
	}
}

// The APM counts what the player does: each Arcane Shot and each walk the rotation starts (one walk in at 10 s, one
// out at 20 s) over a 60 s fight, and nothing for the Auto Shots and the trinket-free gear's procs.
func TestAPMCountsPresses(t *testing.T) {
	// a walk to yd, once, in the second from s while the distance compares (op) to yd
	walk := func(s int, op string, yd int) string {
		at := func(v string) string { return `{"const":{"val":"` + v + `"}}` }
		cmp := func(op, lhs, rhs string) string { return `{"cmp":{"op":"` + op + `","lhs":` + lhs + `,"rhs":` + rhs + `}}` }
		return `{"action":{"condition":{"and":{"vals":[` +
			cmp("OpGe", `{"currentTime":{}}`, at(strconv.Itoa(s)+"s")) + `,` +
			cmp("OpLt", `{"currentTime":{}}`, at(strconv.Itoa(s+1)+"s")) + `,` +
			cmp(op, `{"distanceToTarget":{}}`, at(strconv.Itoa(yd))) + `]}},` +
			`"move":{"rangeFromTarget":` + at(strconv.Itoa(yd)) + `}}}`
	}
	aplJSON := `{"type":"TypeAPL","priorityList":[` + walk(10, "OpGt", 5) + `,` + walk(20, "OpLt", 12) + `,` +
		`{"action":{"castSpell":{"spellId":{"spellId":14287,"rank":8}}}}]}`
	rot := &proto.APLRotation{}
	if err := protojson.Unmarshal([]byte(aplJSON), rot); err != nil {
		t.Fatalf("rotation: %v", err)
	}
	gear := core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet
	gear.Items[proto.ItemSlot_ItemSlotTrinket1] = &proto.ItemSpec{}
	gear.Items[proto.ItemSlot_ItemSlotTrinket2] = &proto.ItemSpec{}
	player := core.WithSpec(&proto.Player{
		Race:               proto.Race_RaceHuman,
		Class:              proto.Class_ClassHunter,
		Equipment:          gear,
		TalentsString:      patchTestTalents,
		Rotation:           rot,
		DistanceFromTarget: 12,
	}, abOptions)
	enc := abEncounter()
	enc.Duration, enc.DurationVariation = 60, 0
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  enc,
		SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 20261024},
	})
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
	}
	um := res.RaidMetrics.Parties[0].Players[0]
	shots := int32(0)
	autos := int32(0)
	for _, a := range um.Actions {
		for _, tg := range a.Targets {
			if a.Id.GetSpellId() == 14287 {
				shots += tg.Casts
			}
			if a.Id.GetOtherId() == proto.OtherAction_OtherActionShoot {
				autos += tg.Casts
			}
		}
	}
	if shots == 0 || autos == 0 {
		t.Fatalf("%d Arcane Shots and %d Auto Shots: the fight did not play", shots, autos)
	}
	if want := float64(shots + 2); math.Abs(um.Apm-want) > 1e-9 {
		t.Errorf("APM %.2f, want %.0f (%d Arcane Shots + 2 walks in a 60 s fight; %d Auto Shots not counted)", um.Apm, want, shots, autos)
	}
}

// A pet's hit is its owner's ranged hit (forever-bugs #394), so the Biznicks 247x128 Accurascope's ranged-only +3%
// reaches the pet, 17/8 of it as spell hit, while the hunter's melee stays without it.
func TestPetInheritsRangedOnlyHit(t *testing.T) {
	_, plain := newPatchTestHunter(t, patchTestSetup{gear: true})
	_, scoped := newPatchTestHunter(t, patchTestSetup{gear: true, rngEnchant: 2523})
	if got := scoped.GetStat(stats.MeleeHit) - plain.GetStat(stats.MeleeHit); math.Abs(got) > 1e-9 {
		t.Fatalf("the scope added %.2f to the hunter's melee and ranged hit; it is ranged only", got)
	}
	scope := 3 * float64(core.MeleeHitRatingPerHitChance)
	if got := scoped.pet.GetStat(stats.MeleeHit) - plain.pet.GetStat(stats.MeleeHit); math.Abs(got-scope) > 1e-9 {
		t.Errorf("the scope gave the pet %.2f hit, want %.2f", got, scope)
	}
	if got := scoped.pet.GetStat(stats.SpellHit) - plain.pet.GetStat(stats.SpellHit); math.Abs(got-scope*17/8) > 1e-9 {
		t.Errorf("the scope gave the pet %.2f spell hit, want %.2f", got, scope*17/8)
	}
}

// Each family's damage, armor and health scalars are its Tamed Pet Passive in the client data (the Fox carries the
// Wolf's), and the Bat keeps Sonic Blast (Forever's new ranks sit on the Bat skill line, 653).
func TestPetFamilyScalars(t *testing.T) {
	// passive: damage %, armor %, health %
	client := map[proto.Hunter_Options_PetType][3]float64{
		proto.Hunter_Options_Bat: {7, 0, 0}, proto.Hunter_Options_Bear: {-9, 5, 8}, proto.Hunter_Options_BirdOfPrey: {7, 0, 0},
		proto.Hunter_Options_Boar: {-10, 9, 4}, proto.Hunter_Options_CarrionBird: {0, 5, 0}, proto.Hunter_Options_Cat: {10, 0, -2},
		proto.Hunter_Options_CoreHound: {1, 0, 5}, proto.Hunter_Options_Crab: {-5, 13, -4}, proto.Hunter_Options_Crocolisk: {0, 10, -5},
		proto.Hunter_Options_Fox: {0, 5, 0}, proto.Hunter_Options_Gorilla: {2, 0, 4}, proto.Hunter_Options_Hyena: {0, 5, 0},
		proto.Hunter_Options_Raptor: {10, 3, -5}, proto.Hunter_Options_Scorpid: {-6, 10, 0}, proto.Hunter_Options_Spider: {7, 0, 0},
		proto.Hunter_Options_Tallstrider: {0, 0, 5}, proto.Hunter_Options_Turtle: {-10, 13, 0}, proto.Hunter_Options_WindSerpent: {7, 0, 0},
		proto.Hunter_Options_Wolf: {0, 5, 0},
	}
	for family, want := range client {
		c := PetConfigs[family]
		if math.Abs(c.Damage-(1+want[0]/100)) > 1e-9 || math.Abs(c.Armor-(1+want[1]/100)) > 1e-9 || math.Abs(c.Health-(1+want[2]/100)) > 1e-9 {
			t.Errorf("%s: damage %.2f armor %.2f health %.2f, want %+.0f%% %+.0f%% %+.0f%%", c.Name, c.Damage, c.Armor, c.Health, want[0], want[1], want[2])
		}
	}
	if !slices.Contains(PetConfigs[proto.Hunter_Options_Bat].Abilities, SonicBlast) {
		t.Error("the Bat has no Sonic Blast; Forever's ranks are on the Bat skill line")
	}
}

// Mana per 5 shows in a row per source: Blessing of Wisdom (40, improved x1.2) and Mana Spring Totem (25, improved
// x1.05) split every tick's mana per 5 by their amounts, casting or not; mana per 5 nobody names (here bonus stats)
// lands in the "Mana per 5 (unnamed source)" row rather than hiding in another.
func TestManaPerFiveRowsBySource(t *testing.T) {
	run := func(bonusMP5 float64) map[string]float64 {
		rot := &proto.APLRotation{}
		if err := protojson.Unmarshal([]byte(`{"type":"TypeAPL","priorityList":[{"action":{"castSpell":{"spellId":{"spellId":14287,"rank":8}}}}]}`), rot); err != nil {
			t.Fatalf("rotation: %v", err)
		}
		bonus := make([]float64, stats.Len)
		bonus[stats.MP5] = bonusMP5
		player := core.WithSpec(&proto.Player{
			Race:               proto.Race_RaceHuman,
			Class:              proto.Class_ClassHunter,
			Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
			TalentsString:      patchTestTalents,
			Rotation:           rot,
			DistanceFromTarget: 11,
			Buffs:              &proto.IndividualBuffs{BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved},
			BonusStats:         &proto.UnitStats{Stats: bonus},
		}, abOptions)
		res := core.RunRaidSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{ManaSpringTotem: proto.TristateEffect_TristateEffectImproved}, &proto.Debuffs{}),
			Encounter:  abEncounter(),
			SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 20261023},
		})
		if res.Error != nil {
			t.Fatalf("sim error: %s", res.Error.Message)
		}
		rows := map[string]float64{}
		for _, r := range res.RaidMetrics.Parties[0].Players[0].Resources {
			switch {
			case r.Id.GetSpellId() == 25290:
				rows["wisdom"] += r.Gain
			case r.Id.GetSpellId() == 10497:
				rows["spring"] += r.Gain
			case r.Id.GetOtherId() == proto.OtherAction_OtherActionManaRegen && r.Id.Tag == 3:
				rows["unnamed"] += r.Gain
			}
		}
		return rows
	}
	named := run(0)
	if named["wisdom"] <= 0 || named["spring"] <= 0 {
		t.Fatalf("no Blessing of Wisdom or Mana Spring Totem row: %v", named)
	}
	if got, want := named["wisdom"]/named["spring"], 40/31.25; math.Abs(got-want) > 1e-6 {
		t.Errorf("Blessing of Wisdom / Mana Spring Totem mana %.4f, want %.4f (their mana per 5)", got, want)
	}
	if named["unnamed"] > 0 {
		t.Errorf("%.1f mana per 5 unnamed with only named sources", named["unnamed"])
	}
	extra := run(10)
	if got, want := extra["unnamed"]/extra["wisdom"], 10.0/40; math.Abs(got-want) > 1e-6 {
		t.Errorf("unnamed / Blessing of Wisdom mana %.4f, want %.4f for 10 unnamed mana per 5", got, want)
	}
}

// A run speed change during a walk changes the rest of it: the yard under way keeps the share already walked and the
// yards after it go at the new speed (Strider Kick's +30% landing just after a walk out began; the same buff running
// out during a walk). A walk with no change steps every 142 ms as before. 6 yd from 11 yd at 7 yd/s.
func TestRunSpeedChangeMidWalk(t *testing.T) {
	speedID := core.ActionID{SpellID: 1317257}.WithTag(1)
	yard := func(speed float64) time.Duration { return time.Millisecond * time.Duration(1000/speed) }
	walk := func(t *testing.T, startFast bool, changeAt time.Duration) time.Duration {
		sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 11})
		hunter.MoveDelay = 0
		sim.PrePull()
		for i := 0; i < 100000 && sim.CurrentTime < 2*time.Second; i++ {
			sim.Step()
		}
		if startFast {
			hunter.AddMoveSpeedModifier(&speedID, 1.3)
		}
		start := sim.CurrentTime
		hunter.MoveTo(5, sim)
		if changeAt > 0 {
			sim.AddPendingAction(&core.PendingAction{NextActionAt: start + changeAt, OnAction: func(*core.Simulation) {
				if startFast {
					hunter.RemoveMoveSpeedModifier(&speedID)
				} else {
					hunter.AddMoveSpeedModifier(&speedID, 1.3)
				}
			}})
		}
		for i := 0; i < 100000 && hunter.MovementHandler.Moving; i++ {
			sim.Step()
		}
		return sim.CurrentTime - start
	}
	slow, fast := yard(7), yard(7*1.3)
	changeAt := 110 * time.Millisecond
	for _, c := range []struct {
		name      string
		startFast bool
		changeAt  time.Duration
		want      time.Duration
	}{
		{"no change", false, 0, 6 * slow},
		{"+30% from 110 ms in", false, changeAt, changeAt + time.Duration((1-float64(changeAt)/float64(slow))*float64(fast)) + 5*fast},
		// the first yard ended at 109 ms, so the second is under way
		{"+30% running out 110 ms in", true, changeAt, changeAt + time.Duration((1-float64(changeAt-fast)/float64(fast))*float64(slow)) + 4*slow},
	} {
		t.Run(c.name, func(t *testing.T) {
			if got := walk(t, c.startFast, c.changeAt); got < c.want-time.Millisecond || got > c.want+time.Millisecond {
				t.Errorf("the walk took %v, want %v", got, c.want)
			}
		})
	}
}
