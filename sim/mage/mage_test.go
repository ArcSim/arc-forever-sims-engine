package mage

import (
	"math"
	"testing"
	"time"

	"github.com/wowsims/forever/sim/arenalib"
	"github.com/wowsims/forever/sim/common"
	_ "github.com/wowsims/forever/sim/common"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
)

func init() {
	RegisterMage()
	common.RegisterAllEffects()
}

func TestArcane(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{mageSuite("arcane", ArcaneTalents)}))
}

func TestFire(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{mageSuite("fire", FireTalents)}))
}

func TestFrost(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{mageSuite("frost", FrostTalents)}))
}

// The community builds our Forever sim ranks: Arcane 35/0/16, Fire 0/35/16 and Frost 14/0/37.
var ArcaneTalents = "055005023100311531--005500033"
var FireTalents = "-03552020130133151-005500033"
var FrostTalents = "050005013--0555003301001301251"

func mageSuite(apl string, talents string) core.CharacterSuiteConfig {
	return core.CharacterSuiteConfig{
		Class:      proto.Class_ClassMage,
		Race:       proto.Race_RaceGnome,
		OtherRaces: []proto.Race{proto.Race_RaceTroll},
		SpecOptions: core.SpecOptionsCombo{Label: "MageArmor", SpecOptions: &proto.Player_Mage{
			Mage: &proto.Mage{
				Options: &proto.Mage_Options{
					ClassOptions: &proto.MageOptions{
						DefaultMageArmor: proto.MageArmor_MageArmorMageArmor,
					},
				},
			},
		}},
		// Naked: the generated item database does not carry most of the pre-raid set our Forever sim
		// tests with yet, and gives the rest TBC-shaped stats.
		GearSet:  core.GearSetCombo{Label: "Naked", GearSet: &proto.EquipmentSpec{}},
		Talents:  talents,
		Rotation: core.GetAplRotation("../../ui/specs/mage/dps/apls", apl),
		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeStaff,
			},
			ArmorType: proto.ArmorType_ArmorTypeCloth,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeWand,
			},
			EnchantBlacklist: []int32{2673, 3225, 3273},
		},
	}
}

// Hot Streak (400625) has one charge: the Pyroblast its stacks speed up spends all of them.
func TestHotStreakSpentByPyroblast(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Mage", Class: proto.Class_ClassMage, Race: proto.Race_RaceGnome, TalentsString: FireTalents,
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{},
			Spec:     &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{ClassOptions: &proto.MageOptions{}}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	mage := sim.Raid.Parties[0].Players[0].(MageAgent).GetMage()
	pyroblastRank := spellData.Pyroblast.Highest()
	pyroblast := mage.GetSpell(core.ActionID{SpellID: pyroblastRank.ID})

	mage.HotStreakAura.Activate(sim)
	mage.HotStreakAura.SetStacks(sim, 3)
	if !pyroblast.Cast(sim, mage.CurrentTarget) {
		t.Fatal("Pyroblast did not cast")
	}
	if want := pyroblastRank.CastTime() / 4; mage.Hardcast.Expires != want {
		t.Errorf("Pyroblast cast ends at %v, want %v with 3 stacks", mage.Hardcast.Expires, want)
	}

	for sim.CurrentTime < 5*time.Second && mage.HotStreakAura.IsActive() {
		sim.Step()
	}
	if mage.HotStreakAura.IsActive() {
		t.Errorf("Hot Streak still up at %v after Pyroblast finished", sim.CurrentTime)
	}
}

// Since the beta build of 2026-10-16 (Blizzard's beta notes, forums t/2360696 post 6) every Fire Blast
// adds a Heating Up stack, crit or not, and a Pyroblast made instant by Presence of Mind keeps them.
func TestHeatingUpFromFireBlastAndKeptThroughPresenceOfMind(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Mage", Class: proto.Class_ClassMage, Race: proto.Race_RaceGnome, TalentsString: "000000000000001-03552020130133151",
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{},
			Spec:     &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{ClassOptions: &proto.MageOptions{}}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	mage := sim.Raid.Parties[0].Players[0].(MageAgent).GetMage()
	fireBlast := mage.GetSpell(core.ActionID{SpellID: spellData.FireBlast.Highest().ID})
	pom := mage.GetSpell(core.ActionID{SpellID: spellData.PresenceOfMind.Highest().ID})
	pyroblast := mage.GetSpell(core.ActionID{SpellID: spellData.Pyroblast.Highest().ID})

	if !fireBlast.Cast(sim, mage.CurrentTarget) {
		t.Fatal("Fire Blast did not cast")
	}
	if !mage.HotStreakAura.IsActive() || mage.HotStreakAura.GetStacks() < 1 {
		t.Fatal("Fire Blast gave no Heating Up stack")
	}

	for !mage.GCD.IsReady(sim) {
		sim.Step()
	}
	mage.HotStreakAura.SetStacks(sim, 3)
	if !pom.Cast(sim, mage.CurrentTarget) || !pyroblast.Cast(sim, mage.CurrentTarget) {
		t.Fatal("Presence of Mind + Pyroblast did not cast")
	}
	if mage.Hardcast.Expires > sim.CurrentTime {
		t.Fatal("Presence of Mind left Pyroblast a cast time")
	}
	if got := mage.HotStreakAura.GetStacks(); !mage.HotStreakAura.IsActive() || got != 3 {
		t.Errorf("Heating Up at %d stacks after an instant Pyroblast, want 3 kept", got)
	}
}

// Frostfire Bolt (1237313) casts, lands its bolt and dot, and takes Improved Fireball's cast time cut
// exactly as Fireball does: the client's mask on 11069 names both.
func TestFrostfireBolt(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Mage", Class: proto.Class_ClassMage, Race: proto.Race_RaceGnome, TalentsString: FireTalents,
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{},
			Spec:     &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{ClassOptions: &proto.MageOptions{}}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	mage := sim.Raid.Parties[0].Players[0].(MageAgent).GetMage()
	if mage.Talents.ImprovedFireball == 0 {
		t.Fatal("FireTalents no longer take Improved Fireball; pick a build that does")
	}
	fireball := mage.GetSpell(core.ActionID{SpellID: spellData.Fireball.Highest().ID})
	ffbRank := spellData.FrostfireBolt.Highest()
	ffb := mage.GetSpell(core.ActionID{SpellID: ffbRank.ID})
	if ffb == nil {
		t.Fatal("Frostfire Bolt is not registered")
	}
	if want := ffbRank.CastTime() - (spellData.Fireball.Highest().CastTime() - fireball.DefaultCast.CastTime); ffb.DefaultCast.CastTime != want {
		t.Errorf("Frostfire Bolt casts in %v, want %v (Improved Fireball's cut)", ffb.DefaultCast.CastTime, want)
	}

	if !ffb.Cast(sim, mage.CurrentTarget) {
		t.Fatal("Frostfire Bolt did not cast")
	}
	for sim.CurrentTime < 8*time.Second {
		sim.Step()
	}
	if ffb.SpellMetrics[0].TotalDamage <= 0 || !ffb.Dot(mage.CurrentTarget).IsActive() {
		t.Errorf("Frostfire Bolt dealt %v and its dot is up: %v", ffb.SpellMetrics[0].TotalDamage, ffb.Dot(mage.CurrentTarget).IsActive())
	}
}

// The arena entry for this spec. Without ARENA_OUT set it only checks every build's damage against the spell manifest; see sim/arenalib.
func TestArena(t *testing.T) {
	arenalib.Run(t, arenalib.Spec{
		Dir:   "mage",
		UI:    "mage/dps",
		Class: proto.Class_ClassMage,
		Race:  proto.Race_RaceGnome,
		SpecOptions: &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{
			ClassOptions: &proto.MageOptions{DefaultMageArmor: proto.MageArmor_MageArmorMageArmor},
		}}},
		Role:               arenalib.Caster,
		DistanceFromTarget: 30,
	})
}

// Frostbolt rolls the spread its client row states, as the beta's do (Gromnie, foreverlogs 2679:
// rank 4 non-crits land anywhere from 65 to 74), where it used to deal the row's average every cast.
func TestFrostboltRollsItsRow(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Mage", Class: proto.Class_ClassMage, Race: proto.Race_RaceGnome, TalentsString: FrostTalents,
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{},
			Spec:     &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{ClassOptions: &proto.MageOptions{}}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	mage := sim.Raid.Parties[0].Players[0].(MageAgent).GetMage()
	frostbolt := mage.GetSpell(core.ActionID{SpellID: spellData.Frostbolt.Highest().ID})
	m := &frostbolt.SpellMetrics[0]
	seen := map[float64]bool{}
	for i := 0; i < 30; i++ {
		damage, hits, crits, resisted := m.TotalDamage, m.Hits, m.Crits, m.ResistedHits
		frostbolt.ApplyEffects(sim, mage.CurrentTarget, frostbolt)
		for end := sim.CurrentTime + 5*time.Second; sim.CurrentTime < end; {
			sim.Step()
		}
		if m.Hits > hits && m.Crits == crits && m.ResistedHits == resisted {
			seen[math.Round(m.TotalDamage-damage)] = true
		}
	}
	if len(seen) < 3 {
		t.Errorf("Frostbolt non-crits dealt only %v; it should roll its row's spread", seen)
	}
}

// Master of Elements (29074) refunds once per cast: its 9 ms ProcCategoryRecovery stops a Cone of
// Cold that crits three targets from refunding three times.
func TestMasterOfElementsRefundsOncePerCast(t *testing.T) {
	encounter := core.MakeSingleTargetEncounter(0)
	encounter.Targets = append(encounter.Targets, core.NewDefaultTarget(), core.NewDefaultTarget())
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Mage", Class: proto.Class_ClassMage, Race: proto.Race_RaceGnome, TalentsString: FireTalents,
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{},
			Spec:     &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{ClassOptions: &proto.MageOptions{}}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: encounter,
	}, simsignals.CreateSignals())
	sim.Reset()

	mage := sim.Raid.Parties[0].Players[0].(MageAgent).GetMage()
	if mage.Talents.MasterOfElements == 0 {
		t.Fatal("FireTalents no longer take Master of Elements; pick a build that does")
	}
	coneOfCold := mage.GetSpell(core.ActionID{SpellID: spellData.ConeOfCold.Highest().ID})
	mage.AddStatDynamic(sim, stats.SpellCritPercent, 100)
	mage.SpendMana(sim, mage.CurrentMana()/2, mage.NewManaMetrics(coneOfCold.ActionID))

	before, crits := mage.CurrentMana(), coneOfCold.SpellMetrics
	if !coneOfCold.Cast(sim, mage.CurrentTarget) {
		t.Fatal("Cone of Cold did not cast")
	}
	totalCrits := int32(0)
	for i := range crits {
		totalCrits += crits[i].Crits
	}
	if totalCrits < 2 {
		t.Fatalf("Cone of Cold crit %d targets; the test needs at least 2", totalCrits)
	}
	want := float64(coneOfCold.Cost.BaseCost) * spellData.MasterOfElements.FractionAt(mage.Talents.MasterOfElements)
	if got := mage.CurrentMana() - before + coneOfCold.CurCast.Cost; math.Abs(got-want) > 0.01 {
		t.Errorf("Master of Elements refunded %.1f mana off %d crits, want one refund of %.1f", got, totalCrits, want)
	}
}

// Arcane Instability (15058) and Arcane Power (12042) name Frostfire Bolt in their SPELLMOD_DAMAGE
// mask but not their SPELLMOD_DOT one (client 1.60.1.70170): the bolt's hit takes the bonus, its DoT
// does not.
func TestArcaneBonusesSkipFrostfireBoltDot(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Mage", Class: proto.Class_ClassMage, Race: proto.Race_RaceGnome, TalentsString: ArcaneTalents,
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{},
			Spec:     &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{ClassOptions: &proto.MageOptions{}}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	mage := sim.Raid.Parties[0].Players[0].(MageAgent).GetMage()
	if mage.Talents.ArcaneInstability == 0 || !mage.Talents.ArcanePower {
		t.Fatal("ArcaneTalents no longer take Arcane Instability and Arcane Power; pick a build that does")
	}
	ffb := mage.GetSpell(core.ActionID{SpellID: spellData.FrostfireBolt.Highest().ID})
	table := mage.AttackTables[mage.CurrentTarget.UnitIndex]
	gap := func() float64 {
		return ffb.AttackerDamageMultiplier(table, false) - ffb.AttackerDamageMultiplier(table, true)
	}

	instability := spellData.ArcaneInstability.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DAMAGE)).FractionAt(mage.Talents.ArcaneInstability)
	// Piercing Ice's hit/DoT split, if the build takes it (TestPiercingIceDotBonusStaysAtBasePoints).
	if points := mage.Talents.PiercingIce; points > 0 {
		instability += spellData.PiercingIce.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DAMAGE)).FractionAt(points) -
			spellData.PiercingIce.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DOT)).FractionAt(points)
	}
	if got := gap(); math.Abs(got-instability) > 1e-9 {
		t.Errorf("Frostfire Bolt hit - DoT multiplier = %.4f, want %.4f on the hit alone", got, instability)
	}
	mage.ArcanePowerAura.Activate(sim)
	power := spellData.ArcanePower.Highest().Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DAMAGE)).Average(core.CharacterLevel) / 100
	if got := gap(); math.Abs(got-instability-power) > 1e-9 {
		t.Errorf("with Arcane Power, Frostfire Bolt hit - DoT multiplier = %.4f, want %.4f", got, instability+power)
	}
}

// Piercing Ice (11151): effect 0 (SPELLMOD_DAMAGE) scales 2/4/6%, but effect 1 (SPELLMOD_DOT, Blizzard
// and Frostfire Bolt) has no rank curve in client 1.60.1.70170 and stays at 2%. Frostfire Bolt's hit
// takes the full bonus, its DoT 2%.
func TestPiercingIceDotBonusStaysAtBasePoints(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Mage", Class: proto.Class_ClassMage, Race: proto.Race_RaceGnome, TalentsString: FrostTalents,
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{},
			Spec:     &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{ClassOptions: &proto.MageOptions{}}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	mage := sim.Raid.Parties[0].Players[0].(MageAgent).GetMage()
	if mage.Talents.PiercingIce < 2 {
		t.Fatal("FrostTalents no longer take 2+ points of Piercing Ice; pick a build that does")
	}
	hit := spellData.PiercingIce.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DAMAGE)).FractionAt(mage.Talents.PiercingIce)
	dot := spellData.PiercingIce.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DOT)).FractionAt(mage.Talents.PiercingIce)
	if dot != 0.02 || hit <= dot {
		t.Fatalf("Piercing Ice rank %d: hit %.3f, DoT %.3f; client 70170 has 2/4/6%% and a flat 2%%", mage.Talents.PiercingIce, hit, dot)
	}

	ffb := mage.GetSpell(core.ActionID{SpellID: spellData.FrostfireBolt.Highest().ID})
	table := mage.AttackTables[mage.CurrentTarget.UnitIndex]
	if got := ffb.AttackerDamageMultiplier(table, false) - ffb.AttackerDamageMultiplier(table, true); math.Abs(got-(hit-dot)) > 1e-9 {
		t.Errorf("Frostfire Bolt hit - DoT multiplier = %.4f, want %.4f (Piercing Ice %.2f on the hit, %.2f on the DoT)", got, hit-dot, hit, dot)
	}
}

// Arcane Concentration (11213) does not proc off spells another spell triggers: Blizzard rolls once
// an enemy on its cast, never on its ticks, and Arcane Missiles' missiles never roll (beta log 2706).
func TestArcaneConcentrationSkipsTriggeredSpells(t *testing.T) {
	sim, mage := newThreeTargetMage(ArcaneTalents)
	if mage.Talents.ArcaneConcentration == 0 {
		t.Fatal("ArcaneTalents no longer take Arcane Concentration; pick a build that does")
	}
	var castAt time.Duration
	onCast, later := 0, 0
	mage.ClearcastingAura.ApplyOnGain(func(_ *core.Aura, sim *core.Simulation) {
		if sim.CurrentTime == castAt {
			onCast++
		} else {
			later++
		}
	})
	for _, spell := range []*core.Spell{
		mage.GetSpell(core.ActionID{SpellID: spellData.Blizzard.Highest().ID}),
		mage.GetSpell(core.ActionID{SpellID: spellData.ArcaneMissiles.Highest().ID}),
	} {
		onCast, later = 0, 0
		for i := 0; i < 30; i++ {
			mage.ClearcastingAura.Deactivate(sim)
			castAt = sim.CurrentTime
			spell.SkipCastAndApplyEffects(sim, mage.CurrentTarget)
			for sim.CurrentTime < castAt+9*time.Second {
				sim.Step()
			}
		}
		if later != 0 {
			t.Errorf("%v: Clearcasting procced %d times after the cast (ticks or missiles)", spell.ActionID, later)
		}
		if blizzard := spell.Matches(MageSpellBlizzard); blizzard != (onCast > 0) {
			t.Errorf("%v: Clearcasting procced %d times on the cast", spell.ActionID, onCast)
		}
	}
}

// A naked gnome mage facing three targets, reset and ready for hand-cast spells.
func newThreeTargetMage(talents string) (*core.Simulation, *Mage) {
	encounter := core.MakeSingleTargetEncounter(0)
	encounter.Duration = 600
	encounter.Targets = append(encounter.Targets, core.NewDefaultTarget(), core.NewDefaultTarget())
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Mage", Class: proto.Class_ClassMage, Race: proto.Race_RaceGnome, TalentsString: talents,
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{},
			Spec:     &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{ClassOptions: &proto.MageOptions{}}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: encounter,
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim, sim.Raid.Parties[0].Players[0].(MageAgent).GetMage()
}

// Blizzard's tick rows (1279976...) lack Not a Proc and Winter's Chill (11180) cannot proc from procs,
// so the debuff only rolls on the cast's own hit on each enemy, never on a tick.
func TestWintersChillSkipsBlizzardTicks(t *testing.T) {
	sim, mage := newThreeTargetMage(FrostTalents)
	if mage.Talents.WintersChill == 0 {
		t.Fatal("FrostTalents no longer take Winter's Chill; pick a build that does")
	}
	blizzard := mage.GetSpell(core.ActionID{SpellID: spellData.Blizzard.Highest().ID})
	var castAt time.Duration
	onCast, later := 0, 0
	mage.WintersChillAura.ApplyOnStacksChange(func(_ *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
		if newStacks <= oldStacks {
			return
		}
		if sim.CurrentTime == castAt {
			onCast++
		} else {
			later++
		}
	})
	for i := 0; i < 30; i++ {
		mage.WintersChillAura.Deactivate(sim)
		castAt = sim.CurrentTime
		blizzard.SkipCastAndApplyEffects(sim, mage.CurrentTarget)
		for sim.CurrentTime < castAt+9*time.Second {
			sim.Step()
		}
	}
	if later != 0 || onCast == 0 {
		t.Errorf("Winter's Chill stacked %d times on Blizzard casts and %d times on its ticks; want some and 0", onCast, later)
	}
}

// Presence of Mind and Combustion share client category 1151's 3 min cooldown, so casting one locks
// the other.
func TestPresenceOfMindAndCombustionShareACooldown(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Mage", Class: proto.Class_ClassMage, Race: proto.Race_RaceGnome, TalentsString: "000000000000001-00000000000000001",
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{},
			Spec:     &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{ClassOptions: &proto.MageOptions{}}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	mage := sim.Raid.Parties[0].Players[0].(MageAgent).GetMage()
	pom := mage.GetSpell(core.ActionID{SpellID: spellData.PresenceOfMind.Highest().ID})
	combustion := mage.GetSpell(core.ActionID{SpellID: spellData.Combustion.Highest().ID})
	if pom == nil || combustion == nil {
		t.Fatal("talents did not register Presence of Mind and Combustion")
	}
	if !combustion.CD.IsReady(sim) {
		t.Fatal("Combustion on cooldown before anything was cast")
	}
	if !pom.Cast(sim, mage.CurrentTarget) {
		t.Fatal("Presence of Mind did not cast")
	}
	if got := combustion.CD.TimeToReady(sim); got != 3*time.Minute {
		t.Errorf("Combustion ready in %v after Presence of Mind, want 3m0s", got)
	}
}

// Blizzard's ticks crit: the tick row (1279949) has no Cannot Crit bit, and beta logs show crits
// (foreverlogs 2668: 27 of 301; 2706: 24 of 649).
func TestBlizzardTicksCrit(t *testing.T) {
	result := core.RunRaidSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1, Iterations: 50},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Mage", Class: proto.Class_ClassMage, Race: proto.Race_RaceGnome, TalentsString: FrostTalents,
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{},
			Spec: &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{ClassOptions: &proto.MageOptions{}}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL, PriorityList: []*proto.APLListItem{{Action: &proto.APLAction{
				Action: &proto.APLAction_CastSpell{CastSpell: &proto.APLActionCastSpell{SpellId: &proto.ActionID{RawId: &proto.ActionID_SpellId{SpellId: spellData.Blizzard.Highest().ID}}}},
			}}}},
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	})
	if result.Error != nil {
		t.Fatal(result.Error.Message)
	}
	crits := int32(0)
	for _, action := range result.RaidMetrics.Parties[0].Players[0].Actions {
		for _, target := range action.Targets {
			crits += target.Crits
		}
	}
	if crits == 0 {
		t.Error("no Blizzard tick crit in 50 iterations")
	}
}

// A Goblin Sapper Charge also hits the Mage, and that hit is a fire spell crit
// Ignite hears. Ignite burns enemies only; the self hit must not reach for a
// dot the Mage doesn't have (issue #699).
func TestIgniteIgnoresTheSapperHitOnTheMage(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: &proto.Raid{Parties: []*proto.Party{{Buffs: &proto.PartyBuffs{}, Players: []*proto.Player{{
			Name: "Mage", Class: proto.Class_ClassMage, Race: proto.Race_RaceGnome, TalentsString: FireTalents,
			Profession1: proto.Profession_Engineering, Consumables: &proto.ConsumesSpec{GoblinSapper: true},
			Equipment: &proto.EquipmentSpec{}, Buffs: &proto.IndividualBuffs{},
			Spec:     &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{ClassOptions: &proto.MageOptions{}}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}}}}},
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	mage := sim.Raid.Parties[0].Players[0].(MageAgent).GetMage()
	if mage.Talents.Ignite == 0 {
		t.Fatal("FireTalents no longer take Ignite; pick a build that does")
	}
	mage.AddStatDynamic(sim, stats.SpellCritPercent, 100)

	sapper := mage.GetSpell(core.GoblinSapperActionID)
	if sapper == nil {
		t.Fatal("Goblin Sapper Charge is not registered")
	}
	if !sapper.Cast(sim, mage.CurrentTarget) {
		t.Fatal("Goblin Sapper Charge did not cast")
	}
	for sim.CurrentTime < 2*time.Second && !sim.Step() {
	}
	self := mage.GetSpell(core.GoblinSapperActionID.WithTag(1))
	if self == nil || self.SpellMetrics[mage.UnitIndex].Crits == 0 {
		t.Fatal("the sapper's hit on the Mage did not crit; the test proves nothing")
	}
}

// Arcane Blast is a talent (MageTalents.arcane_blast): without it the arcane APL falls back to Frostbolt
// instead of standing idle.
func TestArcaneWithoutArcaneBlastStillCasts(t *testing.T) {
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceGnome,
		Class:         proto.Class_ClassMage,
		Equipment:     &proto.EquipmentSpec{},
		TalentsString: "055005023000311531--005500033",
		Rotation:      core.GetAplRotation("../../ui/specs/mage/dps/apls", "arcane").Rotation,
	}, &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{ClassOptions: &proto.MageOptions{}}}})
	result := core.RunRaidSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 101, Iterations: 1},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	})
	if result.Error != nil {
		t.Fatal(result.Error.Message)
	}
	if dps := result.RaidMetrics.Dps.Avg; dps <= 0 {
		t.Errorf("arcane mage without Arcane Blast did %.1f DPS", dps)
	}
}

// Arcane Missiles spends the Arcane Blast stacks as the channel starts, but its missiles keep 15% a
// stack (beta log 2721: 1 stack 99-101 a missile, 2 stacks 112-114).
func TestArcaneMissilesKeepTheArcaneBlastStacks(t *testing.T) {
	sim, mage := newThreeTargetMage(ArcaneTalents)
	if mage.ArcaneBlastAura == nil {
		t.Fatal("ArcaneTalents no longer take Arcane Blast; pick a build that does")
	}
	missiles := mage.GetSpell(core.ActionID{SpellID: spellData.ArcaneMissiles.Highest().ID})
	missile := mage.GetSpell(core.ActionID{SpellID: spellData.ArcaneMissilesTriggered.Highest().ID})
	metrics := &missile.SpellMetrics[mage.CurrentTarget.UnitIndex]

	// A missile has no damage range, so every plain hit (no crit, no partial resist) deals the same.
	plain := func() (float64, int32) {
		return metrics.TotalDamage - metrics.TotalCritDamage - (metrics.TotalResistedDamage - metrics.TotalResistedCritDamage),
			metrics.Hits - metrics.ResistedHits
	}
	perMissile := make([]float64, 3)
	for stacks := range perMissile {
		damage, hits := plain()
		for i := 0; i < 10; i++ {
			mage.ArcaneBlastAura.Deactivate(sim)
			if stacks > 0 {
				mage.ArcaneBlastAura.Activate(sim)
				mage.ArcaneBlastAura.SetStacks(sim, int32(stacks))
			}
			start := sim.CurrentTime
			missiles.SkipCastAndApplyEffects(sim, mage.CurrentTarget)
			if mage.ArcaneBlastAura.IsActive() {
				t.Fatal("Arcane Missiles left the Arcane Blast stacks up")
			}
			for sim.CurrentTime < start+6*time.Second {
				sim.Step()
			}
		}
		newDamage, newHits := plain()
		perMissile[stacks] = (newDamage - damage) / float64(newHits-hits)
	}
	// The stacks join the missile's other additive bonuses (Arcane Instability and the like).
	base := missile.DamageMultiplierAdditive
	for stacks := range perMissile {
		want := (base + 0.15*float64(stacks)) / base
		if got := perMissile[stacks] / perMissile[0]; math.Abs(got-want) > 1e-9 {
			t.Errorf("missiles after %d Arcane Blast stacks = %.4f x those after none, want %.4f", stacks, got, want)
		}
	}
}

// Mage spells reach as far as their client rows state, and the point-blank ones only the targets
// within their 10 yd radius; Arcane Geometry, Flame Throwing and Arctic Reach extend them. They used
// to cast at any distance, so the three talents did nothing.
func TestSpellRangesFollowTheClient(t *testing.T) {
	for _, c := range []struct {
		talents string
		spell   int32
		want    float64
	}{
		{"", spellData.Frostbolt.Highest().ID, 30},
		{"", spellData.FireBlast.Highest().ID, 20},
		{"", spellData.ConeOfCold.Highest().ID, 10},
		{ArcaneTalents, spellData.ArcaneBlast.Highest().ID, 36},  // Arcane Geometry 2/2: +6 yd
		{FireTalents, spellData.FireBlast.Highest().ID, 26},      // Flame Throwing 2/2: +6 yd
		{"--000000000002", spellData.Blizzard.Highest().ID, 36},  // Arctic Reach 2/2: +20% range
		{"--000000000002", spellData.FrostNova.Highest().ID, 12}, // and radius
	} {
		sim, mage := newThreeTargetMage(c.talents)
		spell := mage.GetSpell(core.ActionID{SpellID: c.spell})
		if spell.MaxRange != c.want {
			t.Errorf("%q: %s reaches %.1f yd, want %.1f", c.talents, spell.ActionID, spell.MaxRange, c.want)
		}
		mage.DistanceFromTarget = c.want + 1
		if spell.CanCast(sim, mage.CurrentTarget) {
			t.Errorf("%q: %s casts from %.1f yd", c.talents, spell.ActionID, mage.DistanceFromTarget)
		}
	}
}
