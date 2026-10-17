package hunter

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/simsignals"
	"github.com/wowsims/classic/sim/core/stats"
	"google.golang.org/protobuf/encoding/protojson"
)

// The Survival weave talents (Strider Kick, Resourcefulness).
const patchTestTalents = "-00531500101-500230230250222151"

type patchTestSetup struct {
	race      proto.Race
	petType   proto.Hunter_Options_PetType
	buffs     *proto.RaidBuffs
	talents   string
	gear      bool
	distance  float64
	stone     bool
	ibuffs    *proto.IndividualBuffs
	debuffs   *proto.Debuffs
	petBug    bool // the old "Beta pet speed bug" option
	autoSpam  bool
	gearSet   string // a gear set other than 2h_weave (implies gear)
	mhImbue   proto.WeaponImbue
	ohImbue   proto.WeaponImbue
	ohItem    int32 // replaces the gear set's off hand
	runSpeed  float64
	consumes  *proto.Consumes          // starting consumables (the stone and imbue fields still apply on top)
	mhEnchant int32                    // replaces the gear set's main hand enchant
	skyborne  bool                     // the Skyborne 15 min racial option
	outside   bool                     // the encounter is outside the Natural flasks' zones
	pieces    map[proto.ItemSlot]int32 // items put in these slots (on top of any gear set)
	engineer  bool                     // Engineering as a profession
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
		Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		Buffs:    setup.ibuffs,
		Spec: &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{
			PetType:             setup.petType,
			PetUptime:           1,
			PetAttackSpeed:      proto.Hunter_Options_One,
			PetWindfuryTotem:    true,
			PetFlametongueTotem: true,
			PetSpeedBug:         setup.petBug,
			AutoShotSpamBug:     setup.autoSpam,
			RunSpeed:            setup.runSpeed,
			SkyborneConvergence: setup.skyborne,
		}}},
	}
	if setup.engineer {
		player.Profession1 = proto.Profession_Engineering
	}
	encounter := core.MakeSingleTargetEncounter(0)
	encounter.OutsideHyjal = setup.outside
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

// Racials: Ruthless (Human) +2% crit and Savage (Orc) +1% crit with any weapon; Earth Fury (Dwarf) and
// Endurance (Tauren) +2% critical damage, Endurance gives no hit; Wind Blessed (Skyborne) 3% haste;
// Blood Fury on a 3 min cooldown.
func TestRacials(t *testing.T) {
	critPerAgi := core.CritPerAgiAtLevel[proto.Class_ClassHunter] * core.CritRatingPerCritChance
	racialCrit := func(race proto.Race) float64 {
		_, hunter := newPatchTestHunter(t, patchTestSetup{race: race})
		return hunter.GetStat(stats.MeleeCrit) - hunter.GetStat(stats.Agility)*critPerAgi
	}
	nightElf := racialCrit(proto.Race_RaceNightElf)
	for race, want := range map[proto.Race]float64{proto.Race_RaceHuman: 2, proto.Race_RaceOrc: 1, proto.Race_RaceDwarf: 0, proto.Race_RaceTauren: 0} {
		if got := racialCrit(race) - nightElf; math.Abs(got-want) > 1e-6 {
			t.Errorf("%v: +%.3f%% racial crit, want %.0f", race, got, want)
		}
	}
	for race, want := range map[proto.Race]float64{proto.Race_RaceDwarf: 1.02, proto.Race_RaceTauren: 1.02, proto.Race_RaceHuman: 1, proto.Race_RaceOrc: 1} {
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
	_, skyborne := newPatchTestHunter(t, patchTestSetup{race: proto.Race_RaceSkyborneAlliance})
	if got := skyborne.PseudoStats.RangedSpeedMultiplier / nelf.PseudoStats.RangedSpeedMultiplier; math.Abs(got-1.03) > 1e-9 {
		t.Errorf("Skyborne ranged speed x%.4f, want 1.03", got)
	}
	_, orc := newPatchTestHunter(t, patchTestSetup{race: proto.Race_RaceOrc})
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

// Auto Shot needs 0.5 s of standing still: after a move back to range it cannot fire sooner.
func TestAutoShotStandStill(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 9})
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
	stoppedAt := step(9)
	if next := hunter.AutoAttacks.NextRangedAttackAt(); next < stoppedAt+time.Millisecond*500 {
		t.Errorf("Auto Shot due %.3fs after stopping at range, want at least 0.5 s", (next - stoppedAt).Seconds())
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

// Immolation Trap is a real cast (placed on a GCD), so the results count its casts and cast time; without them the
// table cannot show its damage per cast or per execute time.
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
			Equipment:          core.GetGearSet("../../ui/hunter/gear_sets", "2h_weave").GearSet,
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
// 52 +-15% (44 to 60) plus 0.56 of spell power.
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
	if lo < 44-1e-6 || hi > 60+1e-6 || lo > 46 || hi < 58 {
		t.Errorf("Shadow Oil bolt hits %.1f to %.1f, want 44 to 60", lo, hi)
	}
}

// Frost Oil's bolt is spell 350025 in Forever (the Classic 1191 is not in its data, so the site showed a nameless row):
// 33 +-12% (29 to 37) plus 0.269 of spell power.
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
	if lo < 29-1e-6 || hi > 37+1e-6 || lo > 30.5 || hi < 35.5 {
		t.Errorf("Frost Oil bolt hits %.1f to %.1f, want 29 to 37", lo, hi)
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
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 9})
	if !hunter.RaptorStrike.CanCast(sim, hunter.CurrentTarget) {
		t.Fatal("Raptor Strike cannot be queued at 9 yd")
	}
	hunter.RaptorStrike.Cast(sim, hunter.CurrentTarget)
	if hunter.curQueueAura == nil || !hunter.curQueueAura.IsActive() {
		t.Fatal("Raptor Strike is not queued after casting it at range")
	}
	if got := hunter.TryRaptorStrike(sim, hunter.AutoAttacks.MHAuto()); got != hunter.RaptorStrikeHit {
		t.Errorf("the next melee swing is %v, want the queued Raptor Strike", got.ActionID)
	}
}

// A melee hit restarts the Auto Shot timer from that moment: the next Auto Shot is one full cycle after it (confirmed in
// game with the swing timer events). Raptor Strike goes through the normal spell path; white swings restart it in the
// core swing loop (attack.go), which a test cannot drive directly.
func TestMeleeSwingRestartsAutoShot(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 3})
	want := sim.CurrentTime + hunter.AutoAttacks.RangedSwingSpeed()
	// a miss or dodge does not count: strike until one lands
	for i := 0; i < 50 && hunter.AutoAttacks.NextRangedAttackAt() != want; i++ {
		hunter.RaptorStrikeHit.Cast(sim, hunter.CurrentTarget)
	}
	if got := hunter.AutoAttacks.NextRangedAttackAt(); got != want {
		t.Errorf("after a Raptor Strike hit the next Auto Shot is at %v, want %v (one cycle after the hit)", got, want)
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

// Stepping into melee turns the melee attack on as the walk starts: the main hand swings as soon as it is ready, the
// off hand only half its swing time after the walk began (or once in range when the walk took longer).
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
	attackOn := sim.CurrentTime + 50*time.Millisecond
	if got, want := hunter.AutoAttacks.OffhandSwingAt(), max(attackOn, walkStart+hunter.AutoAttacks.OffhandSwingSpeed()/2); got != want {
		t.Errorf("off hand swings at %v, want %v (walk start + half its swing time)", got, want)
	}
	if got := hunter.AutoAttacks.MainhandSwingAt(); got > attackOn {
		t.Errorf("main hand swings at %v, want by %v (ready, no wait)", got, attackOn)
	}

	// out and straight back in without a shot: the attack stayed on, so a ready off hand swings once in range
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
	hunter.AutoAttacks.SetOffhandSwingAt(sim.CurrentTime)
	hunter.MoveTo(3, sim)
	walkUntilMeleeRange(sim, hunter)
	if got, want := hunter.AutoAttacks.OffhandSwingAt(), sim.CurrentTime+50*time.Millisecond; got != want {
		t.Errorf("no shot between: off hand swings at %v, want %v (once in range, no wait)", got, want)
	}
}

// A walk into melee swings the moment it reaches 5 yd, before it ends at 3 yd.
func TestMeleeSwingsStartAtFiveYards(t *testing.T) {
	sim, hunter := newPatchTestHunter(t, patchTestSetup{gearSet: "surv_dw", distance: 9})
	sim.PrePull()
	for i := 0; i < 100000 && sim.CurrentTime < 5*time.Second; i++ {
		sim.Step()
	}
	hunter.MoveTo(3, sim)
	walkUntilMeleeRange(sim, hunter)
	if !hunter.MovementHandler.Moving {
		t.Fatal("the walk ended at 5 yd")
	}
	swing := hunter.AutoAttacks.MainhandSwingAt()
	for i := 0; i < 100000 && hunter.MovementHandler.Moving; i++ {
		sim.Step()
	}
	if arrival := sim.CurrentTime; swing >= arrival {
		t.Errorf("main hand swing at %v, want before the walk ends at %v", swing, arrival)
	}
	if got, want := hunter.AutoAttacks.MainhandSwingAt(), swing+hunter.AutoAttacks.MainhandSwingSpeed(); got != want {
		t.Errorf("next main hand swing at %v, want %v (one swing mid-walk at %v)", got, want, swing)
	}
}

// walkUntilMeleeRange steps a walk in until it reaches 5 yd.
func walkUntilMeleeRange(sim *core.Simulation, hunter *Hunter) {
	for i := 0; i < 100000 && hunter.DistanceFromTarget > core.MaxMeleeAttackDistance; i++ {
		sim.Step()
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
// Agility, the wolf foods give Agility, and the new foods and elixirs give what their buffs say.
func TestForeverFoodAndElixirValues(t *testing.T) {
	_, base := newPatchTestHunter(t, patchTestSetup{})
	type want struct{ agi, str, intel, crit float64 }
	cases := map[string]struct {
		consumes *proto.Consumes
		want     want
	}{
		"Grilled Squid":         {&proto.Consumes{Food: proto.Food_FoodGrilledSquid}, want{crit: 1}},
		"Hot Wolf Ribs":         {&proto.Consumes{Food: proto.Food_FoodHotWolfRibs}, want{agi: 10}},
		"Tender Wolf Steak":     {&proto.Consumes{Food: proto.Food_FoodTenderWolfSteak}, want{agi: 15}},
		"Flank au Poivre":       {&proto.Consumes{Food: proto.Food_FoodFlankAuPoivre}, want{agi: 20}},
		"Elixir of Cunning":     {&proto.Consumes{AgilityElixir: proto.AgilityElixir_ElixirOfCunning}, want{agi: 25, intel: 25}},
		"Elixir of Ferocity":    {&proto.Consumes{AgilityElixir: proto.AgilityElixir_ElixirOfFerocity}, want{agi: 18, str: 18}},
		"Elixir of the Grizzly": {&proto.Consumes{StrengthBuff: proto.StrengthBuff_ElixirOfTheGrizzly}, want{str: 25, crit: 2}},
	}
	critPerAgi := core.CritPerAgiAtLevel[proto.Class_ClassHunter] * core.CritRatingPerCritChance
	for name, tc := range cases {
		_, hunter := newPatchTestHunter(t, patchTestSetup{consumes: tc.consumes})
		agi := hunter.GetStat(stats.Agility) - base.GetStat(stats.Agility)
		got := want{
			agi:   agi,
			str:   hunter.GetStat(stats.Strength) - base.GetStat(stats.Strength),
			intel: hunter.GetStat(stats.Intellect) - base.GetStat(stats.Intellect),
			crit:  (hunter.GetStat(stats.MeleeCrit) - base.GetStat(stats.MeleeCrit) - agi*critPerAgi) / core.CritRatingPerCritChance,
		}
		if math.Abs(got.agi-tc.want.agi) > 1e-6 || math.Abs(got.str-tc.want.str) > 1e-6 || math.Abs(got.intel-tc.want.intel) > 1e-6 || math.Abs(got.crit-tc.want.crit) > 1e-6 {
			t.Errorf("%s: %+v, want %+v", name, got, tc.want)
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

// Skysight (Horde Skyborne): +10% run speed, 30 s on a 2 min cooldown on the GCD. It is a regular speed bonus, so under
// Strider Kick's 30% it adds nothing (AB_SKYSIGHT_STACKS=1 stacks it, for comparing). With the 15 min option it is up
// from the start and is not cast again.
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
	if got := hunter.MovementHandler.MoveSpeed; math.Abs(got-7*1.1) > 1e-9 {
		t.Errorf("Skysight: %.4f yd/s, want %.4f", got, 7*1.1)
	}
	if aura.Duration != 30*time.Second {
		t.Errorf("Skysight lasts %v, want 30s", aura.Duration)
	}
	hunter.GetAura("Strider Kick Speed").Activate(sim)
	want := 7 * 1.3
	if core.SkysightStacks {
		want *= 1.1
	}
	if got := hunter.MovementHandler.MoveSpeed; math.Abs(got-want) > 1e-9 {
		t.Errorf("Skysight under Strider Kick: %.4f yd/s, want %.4f", got, want)
	}

	sim, long := newPatchTestHunter(t, patchTestSetup{race: proto.Race_RaceSkyborneHorde, skyborne: true})
	if a := long.GetAura("Skysight"); !a.IsActive() || a.Duration != core.NeverExpires {
		t.Error("with the 15 min option Skysight is not up the whole fight")
	}
	if long.GetSpell(core.ActionID{SpellID: 1259686}).CanCast(sim, long.CurrentTarget) {
		t.Error("with the 15 min option Skysight can still be cast")
	}
	if got := long.MovementHandler.MoveSpeed; math.Abs(got-7*1.1) > 1e-9 {
		t.Errorf("15 min Skysight: %.4f yd/s, want %.4f", got, 7*1.1)
	}
}

// Read Leyline (Alliance Skyborne): +100% mana regeneration (double the spirit regen), up the whole fight with the
// 15 min option; Horde has Skysight instead.
func TestReadLeyline(t *testing.T) {
	_, plain := newPatchTestHunter(t, patchTestSetup{race: proto.Race_RaceSkyborneAlliance})
	_, long := newPatchTestHunter(t, patchTestSetup{race: proto.Race_RaceSkyborneAlliance, skyborne: true})
	if plain.GetAura("Read Leyline") == nil || plain.GetAura("Skysight") != nil {
		t.Fatal("Alliance Skyborne should have Read Leyline, not Skysight")
	}
	if got := long.PseudoStats.SpiritRegenMultiplier - plain.PseudoStats.SpiritRegenMultiplier; math.Abs(got-1) > 1e-9 {
		t.Errorf("Read Leyline adds %.2f to the spirit regen multiplier, want 1 (+100%%)", got)
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
}

// The log carries every weapon's swing timer for the timeline: each main-hand swing starts a timer ready one swing time
// later, a melee swing restarts the Auto Shot timer on the spot, ready one bow cycle later, and an extra attack (Hand of
// Justice) restarts the main hand's own timer, never one of its own. Only the player's lines count (the pet swings too).
func TestSwingTimerLogs(t *testing.T) {
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
		SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 6, DebugFirstIteration: true},
	})
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
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

// The Walk Time rotation value reads the real run speed: from 9 yd, walking to 5 yd takes 4 / 7 = 0.571 s at normal
// speed and 4 / 7.14 = 0.560 s with Night Elf's 2%, so a cast gated on "under 0.565 s" goes out for a Night Elf only.
func TestWalkTimeValue(t *testing.T) {
	const aplJSON = `{"type":"TypeAPL","recheckAfterAutoShot":true,"priorityList":[{"action":{"condition":{"cmp":{"op":"OpLt","lhs":{"walkTime":{"toDistance":{"const":{"val":"5"}}}},"rhs":{"const":{"val":"0.565s"}}}},"castSpell":{"spellId":{"spellId":14287,"rank":8}}}}]}`
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
			DistanceFromTarget: 9,
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
			DistanceFromTarget: 9,
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

// A walk in that starts at the pull, from inside the first Auto Shot (after the pre-pull walk out to 8 yd and Aimed Shot),
// is under way when the pull turns the attacks on: they wait for the walk to end. Turned on mid-walk, Auto Shot was put
// back from 8 yd and kept firing from inside melee range, and the melee swings never started.
func TestMoveAtPull(t *testing.T) {
	const aplJSON = `{"type":"TypeAPL","prepullActions":[{"action":{"move":{"rangeFromTarget":{"const":{"val":"8"}}}},"doAtValue":{"const":{"val":"-5s"}}},{"action":{"castSpell":{"spellId":{"spellId":20904}}},"doAtValue":{"const":{"val":"-2s"}}}],"priorityList":[{"action":{"move":{"rangeFromTarget":{"const":{"val":"3"}}}}}]}`
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

	sim, hunter := newPatchTestHunter(t, patchTestSetup{gear: true, distance: 9})
	hunter.RaptorStrike.Cast(sim, hunter.CurrentTarget)
	if !queued(hunter) {
		t.Fatal("Raptor Strike did not queue at 9 yd")
	}
	arcane := hunter.GetSpell(core.ActionID{SpellID: 14287})
	if arcane == nil {
		t.Fatal("no Arcane Shot")
	}
	arcane.Cast(sim, hunter.CurrentTarget)
	if queued(hunter) {
		t.Error("Raptor Strike is still queued after Arcane Shot")
	}

	sim, hunter = newPatchTestHunter(t, patchTestSetup{gear: true, distance: 9})
	hunter.RaptorStrike.Cast(sim, hunter.CurrentTarget)
	hunter.AutoAttacks.RangedAuto().Cast(sim, hunter.CurrentTarget)
	if queued(hunter) {
		t.Error("Raptor Strike is still queued after an Auto Shot")
	}
}
