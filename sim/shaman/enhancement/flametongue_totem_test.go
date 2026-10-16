package enhancement

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
)

// Flametongue Totem (patch 70): rank 4 (16387) is a 5 min fire totem whose party aura (15036) adds the fire
// hit 16389 states, 13.63 damage a second of weapon speed, to each landed main-hand auto attack, and a
// main-hand Flametongue Weapon turns the benefit off.

const (
	flametongueTotemCast = 16387
	flametongueTotemHit  = 16389
	flametongueImbueHit  = 16344
	searingTotemCast     = 10438
	searingTotemHit      = 10436
	windfuryTotemCastID  = 10614
	crestedScepter       = 3414  // main-hand mace, 2.6 speed
	stonevaultShiv       = 9384  // dagger, 1.5 speed
	arcaniteReaper       = 12784 // two-hand axe, 3.8 speed
)

type ftCase struct {
	mh, oh  int32
	imbue   proto.ShamanImbue
	imbueOh proto.ShamanImbue
	party   *proto.PartyBuffs
	casts   []int32 // pre-pull casts, in order
	// No talents unless the case names some, so a hit is its base damage and nothing else.
	talents string
	elixir  int32
}

func (c ftCase) player() *proto.Player {
	if c.mh == 0 {
		c.mh = arcaniteReaper
	}
	prepull := ""
	for i, id := range c.casts {
		if i > 0 {
			prepull += ","
		}
		prepull += fmt.Sprintf(`{"action":{"castSpell":{"spellId":{"spellId":%d}}},"doAtValue":{"const":{"val":"-%.1fs"}}}`, id, 4.0-float64(i)*1.5)
	}
	items := make([]*proto.ItemSpec, 16)
	for i := range items {
		items[i] = &proto.ItemSpec{}
	}
	items[14] = &proto.ItemSpec{Id: c.mh}
	if c.oh != 0 {
		items[15] = &proto.ItemSpec{Id: c.oh}
	}
	return &proto.Player{
		Name: "enh", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc, TalentsString: c.talents,
		Equipment: &proto.EquipmentSpec{Items: items},
		Buffs:     &proto.IndividualBuffs{}, Consumables: &proto.ConsumesSpec{BattleElixirId: c.elixir},
		Spec: &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{
			ClassOptions: &proto.ShamanOptions{ImbueMh: c.imbue},
			ImbueOh:      c.imbueOh,
		}}},
		Rotation: core.APLRotationFromJsonString(`{"type":"TypeAPL","prepullActions":[` + prepull + `],"priorityList":[]}`),
	}
}

func (c ftCase) raid() *proto.Raid {
	if c.party == nil {
		c.party = &proto.PartyBuffs{}
	}
	return core.SinglePlayerRaidProto(c.player(), c.party, &proto.RaidBuffs{}, &proto.Debuffs{})
}

// The shaman's spell damage at the pull.
func (c ftCase) spellDamage() float64 {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid:       c.raid(),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim.Raid.Parties[0].Players[0].GetCharacter().GetStat(stats.SpellDamage)
}

type ftRun struct {
	m *proto.UnitMetrics
}

func runFlametongueTotem(t *testing.T, c ftCase) ftRun {
	t.Helper()
	res := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       c.raid(),
		Encounter:  core.MakeSingleTargetEncounter(0),
		SimOptions: &proto.SimOptions{Iterations: 40, RandomSeed: 7},
	})
	if res.Error != nil {
		t.Fatal(res.Error.Message)
	}
	return ftRun{m: res.RaidMetrics.Parties[0].Players[0]}
}

// A row of the action table, summed over targets and iterations.
type ftAction struct {
	casts, hits, crits, misses, dodges, parries, damage, critDamage float64
	resistedHits, resistedCrits, resistedDamage, resistedCritDamage float64
}

// Every hit that dealt damage or missed: what a trigger on landed hits fires once for.
func (a ftAction) rolled() float64 { return a.hits + a.crits + a.misses }

// What a melee swing landed on the target: not a miss, a dodge or a parry.
func (a ftAction) landed() float64 { return a.casts - a.misses - a.dodges - a.parries }

// The damage of a hit that took no partial resist, which is the base damage times the multipliers.
func (a ftAction) cleanHit() float64 {
	return (a.damage - a.critDamage - (a.resistedDamage - a.resistedCritDamage)) / (a.hits - a.resistedHits)
}

// The same for a crit: the base damage times the crit multiplier.
func (a ftAction) cleanCrit() float64 {
	return (a.critDamage - a.resistedCritDamage) / (a.crits - a.resistedCrits)
}

func (r ftRun) rows(match func(*proto.ActionMetrics) bool) ftAction {
	var out ftAction
	for _, a := range r.m.Actions {
		if !match(a) {
			continue
		}
		for _, tgt := range a.Targets {
			out.casts += float64(tgt.Casts)
			out.hits += float64(tgt.Hits)
			out.crits += float64(tgt.Crits)
			out.misses += float64(tgt.Misses)
			out.dodges += float64(tgt.Dodges)
			out.parries += float64(tgt.Parries)
			out.damage += tgt.Damage
			out.critDamage += tgt.CritDamage
			out.resistedHits += float64(tgt.ResistedHits)
			out.resistedCrits += float64(tgt.ResistedCrits)
			out.resistedDamage += tgt.ResistedDamage
			out.resistedCritDamage += tgt.ResistedCritDamage
		}
	}
	return out
}

func (r ftRun) spell(id int32) ftAction {
	return r.rows(func(a *proto.ActionMetrics) bool { return a.Id.GetSpellId() == id && a.Id.Tag == 0 })
}

// The swings of one hand's timer: tag 1 is the main hand, tag 2 the off hand.
func (r ftRun) swings(tag int32) ftAction {
	return r.rows(func(a *proto.ActionMetrics) bool {
		return a.Id.GetOtherId() == proto.OtherAction_OtherActionAttack && a.Id.Tag == tag
	})
}

func (r ftRun) uptime(spellID int32) float64 {
	var up float64
	for _, a := range r.m.Auras {
		if a.Id.GetSpellId() == spellID && a.Id.Tag == 0 && a.UptimeSecondsAvg > up {
			up = a.UptimeSecondsAvg
		}
	}
	return up
}

func near(got, want float64) bool { return math.Abs(got-want) < 1e-6 }

// The totem hit is 13.63 a second of the main hand's speed, 1363 hundredths: a 2.6 speed mace hits for
// 35.438, a 3.8 speed axe for 51.794, a 1.5 speed dagger for 20.445, and a crit takes the school's 1.5. With
// no talents and no spell power these are the whole damage.
func TestFlametongueTotemHitDamage(t *testing.T) {
	for _, row := range []struct {
		name  string
		mh    int32
		speed float64
	}{
		{"a 2.6 speed mace", crestedScepter, 2.6},
		{"a 3.8 speed axe", arcaniteReaper, 3.8},
		{"a 1.5 speed dagger", stonevaultShiv, 1.5},
	} {
		t.Run(row.name, func(t *testing.T) {
			hit := runFlametongueTotem(t, ftCase{mh: row.mh, casts: []int32{flametongueTotemCast}}).spell(flametongueTotemHit)
			if hit.hits-hit.resistedHits < 20 || hit.crits-hit.resistedCrits < 5 {
				t.Fatalf("only %v clean hits and %v clean crits, too few to read the damage off", hit.hits-hit.resistedHits, hit.crits-hit.resistedCrits)
			}
			want := row.speed * 13.63
			if got := hit.cleanHit(); !near(got, want) {
				t.Errorf("a hit deals %v, want %v (%v speed at 13.63)", got, want, row.speed)
			}
			if got := hit.cleanCrit(); !near(got, want*1.5) {
				t.Errorf("a crit deals %v, want %v (1.5 times)", got, want*1.5)
			}
		})
	}
}

// The hit takes none of the shaman's spell damage: Hameru tested the totem's hit on the beta (MythicSim
// Discord, 6 October 2026) and spell power did not move it (patch 91). A Greater Arcane Elixir changes
// neither a hit nor a crit. The imbue's own hit keeps Flametongue Attack's 0.1 coefficient.
func TestFlametongueTotemHitIgnoresSpellDamage(t *testing.T) {
	plain := ftCase{mh: crestedScepter, casts: []int32{flametongueTotemCast}}
	boosted := plain
	boosted.elixir = 13454
	extra := boosted.spellDamage() - plain.spellDamage()
	if extra <= 0 {
		t.Fatalf("the elixir adds %v spell damage, want some", extra)
	}
	base, withElixir := runFlametongueTotem(t, plain).spell(flametongueTotemHit), runFlametongueTotem(t, boosted).spell(flametongueTotemHit)
	if got := withElixir.cleanHit() - base.cleanHit(); !near(got, 0) {
		t.Errorf("%v spell damage adds %v to a hit, want nothing", extra, got)
	}
	if got := withElixir.cleanCrit() - base.cleanCrit(); !near(got, 0) {
		t.Errorf("%v spell damage adds %v to a crit, want nothing", extra, got)
	}
}

// Neither Elemental Fury nor Elemental Weapons reaches the totem's hit: it lands as Flametongue Attack 16368
// (beta log 2713, upstream #682), whose class mask (bit 25) neither talent names.
func TestFlametongueTotemHitTakesNoElementalFuryOrElementalWeapons(t *testing.T) {
	const base = 2.6 * 13.63
	fury := runFlametongueTotem(t, ftCase{mh: crestedScepter, talents: "000000000000005", casts: []int32{flametongueTotemCast}}).spell(flametongueTotemHit)
	if got := fury.cleanHit(); !near(got, base) {
		t.Errorf("Elemental Fury 5/5 changes a hit to %v, want it at %v", got, base)
	}
	if got := fury.cleanCrit(); !near(got, base*1.5) {
		t.Errorf("Elemental Fury 5/5 crits for %v, want %v (1.5 times)", got, base*1.5)
	}
	weapons := runFlametongueTotem(t, ftCase{mh: crestedScepter, talents: "-00000003", casts: []int32{flametongueTotemCast}}).spell(flametongueTotemHit)
	if got := weapons.cleanHit(); !near(got, base) {
		t.Errorf("Elemental Weapons 3/3 hits for %v, want %v (unchanged)", got, base)
	}
	if got := weapons.cleanCrit(); !near(got, base*1.5) {
		t.Errorf("Elemental Weapons 3/3 crits for %v, want %v", got, base*1.5)
	}
}

// Only a landed main-hand auto attack adds the hit. The off hand's swings, which are more of them, add none.
func TestFlametongueTotemHitsOnlyMainHandAutoAttacks(t *testing.T) {
	r := runFlametongueTotem(t, ftCase{mh: crestedScepter, oh: stonevaultShiv, casts: []int32{flametongueTotemCast}})
	mh, oh, hit := r.swings(1), r.swings(2), r.spell(flametongueTotemHit)
	if oh.landed() <= mh.landed() {
		t.Fatalf("the off hand landed %v swings against the main hand's %v; the test needs it to swing more", oh.landed(), mh.landed())
	}
	if hit.rolled() != mh.landed() {
		t.Fatalf("%v totem hits for %v landed main-hand swings (%v off-hand), want one each and none for the off hand", hit.rolled(), mh.landed(), oh.landed())
	}
}

// The totem is a cast: with none down, a shaman adds nothing.
func TestFlametongueTotemNeedsTheCast(t *testing.T) {
	if hit := runFlametongueTotem(t, ftCase{}).spell(flametongueTotemHit); hit.rolled() != 0 {
		t.Fatalf("%v totem hits with no totem down", hit.rolled())
	}
	r := runFlametongueTotem(t, ftCase{casts: []int32{flametongueTotemCast}})
	if hit := r.spell(flametongueTotemHit); hit.rolled() == 0 {
		t.Fatal("no totem hits with the totem down")
	}
	if up := r.uptime(flametongueTotemCast); up < 100 {
		t.Fatalf("the totem is up %.1f s of the fight, want it up for the whole of it", up)
	}
}

// A shaman holds one fire totem: Flametongue Totem takes the slot from Searing Totem and gives it back.
func TestFlametongueTotemReplacesSearingTotem(t *testing.T) {
	searingOnly := runFlametongueTotem(t, ftCase{casts: []int32{searingTotemCast}})
	if searingOnly.spell(searingTotemHit).rolled() < 20 || searingOnly.spell(flametongueTotemHit).rolled() != 0 {
		t.Fatalf("Searing Totem alone: %v Searing hits and %v Flametongue hits, want many and none",
			searingOnly.spell(searingTotemHit).rolled(), searingOnly.spell(flametongueTotemHit).rolled())
	}

	flametongueAfter := runFlametongueTotem(t, ftCase{casts: []int32{searingTotemCast, flametongueTotemCast}})
	if got := flametongueAfter.spell(searingTotemHit).rolled(); got != 0 {
		t.Errorf("Searing Totem hit %v times after Flametongue Totem replaced it, want none", got)
	}
	if got := flametongueAfter.spell(flametongueTotemHit).rolled(); got == 0 {
		t.Error("Flametongue Totem added no hits after replacing Searing Totem")
	}

	searingAfter := runFlametongueTotem(t, ftCase{casts: []int32{flametongueTotemCast, searingTotemCast}})
	if got := searingAfter.spell(flametongueTotemHit).rolled(); got != 0 {
		t.Errorf("Flametongue Totem added %v hits after Searing Totem replaced it, want none", got)
	}
	if got := searingAfter.spell(searingTotemHit).rolled(); got < 20 {
		t.Errorf("Searing Totem hit %v times after replacing Flametongue Totem, want many", got)
	}
}

// "When applied to main hand, disables any benefit you personally receive from Flametongue Totem": a
// main-hand Flametongue Weapon leaves the totem up and adds no totem hits, from the shaman's own totem or
// the party's. Windfury Weapon, Frostbrand Weapon and Rockbiter Weapon do not, and an off-hand Flametongue
// Weapon does not either.
func TestOnlyAMainHandFlametongueWeaponDisablesFlametongueTotem(t *testing.T) {
	for _, row := range []struct {
		name      string
		c         ftCase
		totemHits bool
		imbueHits bool
	}{
		{"no imbue", ftCase{}, true, false},
		{"Rockbiter Weapon", ftCase{imbue: proto.ShamanImbue_RockbiterWeapon}, true, false},
		{"Windfury Weapon", ftCase{imbue: proto.ShamanImbue_WindfuryWeapon}, true, false},
		{"Frostbrand Weapon", ftCase{imbue: proto.ShamanImbue_FrostbrandWeapon}, true, false},
		{"Flametongue Weapon in the main hand", ftCase{imbue: proto.ShamanImbue_FlametongueWeapon}, false, true},
		{"Flametongue Weapon in the off hand", ftCase{mh: crestedScepter, oh: stonevaultShiv, imbue: proto.ShamanImbue_RockbiterWeapon, imbueOh: proto.ShamanImbue_FlametongueWeapon}, true, true},
	} {
		for _, source := range []string{"cast", "party"} {
			t.Run(row.name+", the "+source+" totem", func(t *testing.T) {
				c := row.c
				if source == "cast" {
					c.casts = []int32{flametongueTotemCast}
				} else {
					c.party = &proto.PartyBuffs{FlametongueTotem: true}
				}
				r := runFlametongueTotem(t, c)
				if got := r.spell(flametongueTotemHit).rolled() > 0; got != row.totemHits {
					t.Errorf("totem hits %v (%v of them), want %v", got, r.spell(flametongueTotemHit).rolled(), row.totemHits)
				}
				if got := r.spell(flametongueImbueHit).rolled() > 0; got != row.imbueHits {
					t.Errorf("Flametongue Weapon hits %v, want %v", got, row.imbueHits)
				}
				if source == "cast" && r.uptime(flametongueTotemCast) < 100 {
					t.Errorf("the cast totem is up %.1f s, want it standing whatever the imbue", r.uptime(flametongueTotemCast))
				}
			})
		}
	}
}

// The party's Flametongue Totem and the shaman's own are one effect: a hit for each main-hand swing, not
// two.
func TestPartyAndCastFlametongueTotemAddOneHit(t *testing.T) {
	for _, row := range []struct {
		name string
		c    ftCase
	}{
		{"the party's alone", ftCase{mh: crestedScepter, party: &proto.PartyBuffs{FlametongueTotem: true}}},
		{"the cast one alone", ftCase{mh: crestedScepter, casts: []int32{flametongueTotemCast}}},
		{"both", ftCase{mh: crestedScepter, party: &proto.PartyBuffs{FlametongueTotem: true}, casts: []int32{flametongueTotemCast}}},
	} {
		t.Run(row.name, func(t *testing.T) {
			r := runFlametongueTotem(t, row.c)
			if hit, mh := r.spell(flametongueTotemHit), r.swings(1); hit.rolled() == 0 || hit.rolled() != mh.landed() {
				t.Fatalf("%v totem hits for %v landed main-hand swings, want one each", hit.rolled(), mh.landed())
			}
		})
	}
}

// "Flametongue Totem no longer stacks with Windfury Totem" (Forever beta development notes, upstream #677),
// and Windfury holds (patch 98): whichever way each totem reaches the shaman, the party's or its own cast,
// in either cast order, Windfury Totem keeps granting extra attacks, the Flametongue Totem stays down and
// adds no hits beside it.
func TestWindfuryTotemSwitchesFlametongueTotemOff(t *testing.T) {
	for _, row := range []struct {
		name string
		c    ftCase
	}{
		{"both cast, Windfury first", ftCase{casts: []int32{windfuryTotemCastID, flametongueTotemCast}}},
		{"both cast, Flametongue first", ftCase{casts: []int32{flametongueTotemCast, windfuryTotemCastID}}},
		{"the party's Windfury Totem and the cast Flametongue Totem", ftCase{party: &proto.PartyBuffs{WindfuryTotem: true}, casts: []int32{flametongueTotemCast}}},
		{"the cast Windfury Totem and the party's Flametongue Totem", ftCase{party: &proto.PartyBuffs{FlametongueTotem: true}, casts: []int32{windfuryTotemCastID}}},
		{"both the party's", ftCase{party: &proto.PartyBuffs{WindfuryTotem: true, FlametongueTotem: true}}},
		{"the party's Windfury Totem and Grace of Air", ftCase{party: &proto.PartyBuffs{WindfuryTotem: true, GraceOfAirTotem: true, FlametongueTotem: true}}},
		{"the party's Windfury Totem beside a cast Grace of Air", ftCase{party: &proto.PartyBuffs{WindfuryTotem: true, FlametongueTotem: true}, casts: []int32{graceOfAirCast}}},
	} {
		t.Run(row.name, func(t *testing.T) {
			r := runFlametongueTotem(t, row.c)
			extra := r.rows(func(a *proto.ActionMetrics) bool {
				return a.Id.GetOtherId() == proto.OtherAction_OtherActionAttack && (a.Id.Tag == partyTotemExtra || a.Id.Tag == castTotemExtra)
			})
			if extra.casts == 0 {
				t.Error("Windfury Totem granted no extra attacks beside Flametongue Totem")
			}
			if hit := r.spell(flametongueTotemHit); hit.rolled() != 0 {
				t.Errorf("Flametongue Totem added %v hits beside Windfury Totem, want none", hit.rolled())
			}
			if len(row.c.casts) == 2 && r.uptime(flametongueTotemCast) < 100 {
				t.Errorf("the cast Flametongue Totem is up %.1f s, want it standing with no benefit", r.uptime(flametongueTotemCast))
			}
		})
	}
}

// Grace of Air is no longer exclusive with Flametongue Totem (16 October 2026 beta notes): beside it,
// Flametongue Totem adds its hit on every landed main-hand swing.
func TestGraceOfAirLeavesFlametongueTotemAlone(t *testing.T) {
	for _, row := range []struct {
		name string
		c    ftCase
	}{
		{"both cast", ftCase{mh: crestedScepter, casts: []int32{graceOfAirCast, flametongueTotemCast}}},
		{"the party's Grace of Air and the cast Flametongue Totem", ftCase{mh: crestedScepter, party: &proto.PartyBuffs{GraceOfAirTotem: true}, casts: []int32{flametongueTotemCast}}},
		{"both the party's", ftCase{mh: crestedScepter, party: &proto.PartyBuffs{GraceOfAirTotem: true, FlametongueTotem: true}}},
	} {
		t.Run(row.name, func(t *testing.T) {
			r := runFlametongueTotem(t, row.c)
			extra := r.rows(func(a *proto.ActionMetrics) bool {
				return a.Id.GetOtherId() == proto.OtherAction_OtherActionAttack && (a.Id.Tag == partyTotemExtra || a.Id.Tag == castTotemExtra)
			})
			if extra.casts != 0 {
				t.Fatalf("%v Windfury Totem extra attacks, want none with Grace of Air in the air slot", extra.casts)
			}
			if hit, mh := r.spell(flametongueTotemHit), r.swings(1); hit.rolled() == 0 || hit.rolled() != mh.landed() {
				t.Errorf("%v totem hits for %v landed main-hand swings, want one each", hit.rolled(), mh.landed())
			}
		})
	}
}

// Flametongue Totem (16387) adds 16389's fire hit to each landed main-hand auto. It no longer stacks
// with Windfury Totem or Flametongue Weapon (Forever beta development notes), and a main-hand
// Flametongue Weapon disables it (tooltips 8024/16342).
func TestFlametongueTotem(t *testing.T) {
	run := func(mh proto.ShamanImbue, party *proto.PartyBuffs) (autos, hits int32) {
		player := &proto.Player{
			Name: "enh", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc, TalentsString: DefaultTalents,
			Equipment: &proto.EquipmentSpec{Items: []*proto.ItemSpec{
				{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
				{Id: 17182}, // Sulfuras, Hand of Ragnaros
			}},
			Spec: &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{
				ClassOptions: &proto.ShamanOptions{ImbueMh: mh},
			}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL, PriorityList: []*proto.APLListItem{{Action: &proto.APLAction{Action: &proto.APLAction_CastSpell{
				CastSpell: &proto.APLActionCastSpell{SpellId: &proto.ActionID{RawId: &proto.ActionID_SpellId{SpellId: 16387}}},
			}}}}},
		}
		raid := &proto.Raid{Parties: []*proto.Party{{Players: []*proto.Player{player}, Buffs: party}}, Buffs: &proto.RaidBuffs{}, Debuffs: &proto.Debuffs{}, NumActiveParties: 1}
		res := core.RunRaidSim(&proto.RaidSimRequest{Raid: raid, Encounter: core.MakeSingleTargetEncounter(0), SimOptions: &proto.SimOptions{Iterations: 5, RandomSeed: 1}})
		if res.Error != nil {
			t.Fatal(res.Error.Message)
		}
		for _, a := range res.RaidMetrics.Parties[0].Players[0].Actions {
			for _, tgt := range a.Targets {
				if a.Id.GetOtherId() == proto.OtherAction_OtherActionAttack && a.Id.Tag == 1 {
					autos += tgt.Hits + tgt.Crits + tgt.Glances + tgt.Blocks
				}
				if a.Id.GetSpellId() == 16389 {
					hits += tgt.Hits + tgt.Crits + tgt.Misses
				}
			}
		}
		return
	}

	if autos, hits := run(proto.ShamanImbue_NoImbue, &proto.PartyBuffs{}); hits == 0 || hits != autos {
		t.Errorf("alone: %d totem hits for %d landed main-hand autos, want one each", hits, autos)
	}
	if autos, hits := run(proto.ShamanImbue_FrostbrandWeapon, &proto.PartyBuffs{}); hits == 0 || hits != autos {
		t.Errorf("Frostbrand main hand: %d totem hits for %d landed main-hand autos, want one each", hits, autos)
	}
	if _, hits := run(proto.ShamanImbue_NoImbue, &proto.PartyBuffs{WindfuryTotem: true}); hits != 0 {
		t.Errorf("with the party's Windfury Totem: %d totem hits, want 0", hits)
	}
	if _, hits := run(proto.ShamanImbue_FlametongueWeapon, &proto.PartyBuffs{}); hits != 0 {
		t.Errorf("Flametongue main hand: %d totem hits, want 0", hits)
	}
}

// The totem's hit lands as Flametongue Attack 16368 (beta log 2713), whose client row has no spell
// power coefficient and a class mask Elemental Weapons and Elemental Fury don't name; the imbue's
// hit (29469) has both.
func TestFlametongueTotemHitIsNotTheImbuesHit(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: core.SinglePlayerRaidProto(&proto.Player{
			Name: "enh", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc, TalentsString: DefaultTalents,
			Equipment: &proto.EquipmentSpec{Items: []*proto.ItemSpec{
				{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
				{Id: 17182}, // Sulfuras, Hand of Ragnaros
			}},
			Spec: &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{
				ClassOptions: &proto.ShamanOptions{ImbueMh: proto.ShamanImbue_FlametongueWeapon},
			}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	var totem, imbue *core.Spell
	for _, spell := range sim.Raid.Parties[0].Players[0].GetCharacter().Spellbook {
		switch spell.ActionID.SpellID {
		case 16389:
			totem = spell
		case 16344:
			imbue = spell
		}
	}
	if totem == nil || imbue == nil {
		t.Fatalf("totem hit %v, imbue hit %v: want both registered", totem, imbue)
	}
	if totem.BonusCoefficient != 0 || totem.DamageMultiplierAdditive != 1 {
		t.Errorf("totem hit: coefficient %v, additive multiplier %v, want 0 and 1", totem.BonusCoefficient, totem.DamageMultiplierAdditive)
	}
	if imbue.BonusCoefficient == 0 || imbue.DamageMultiplierAdditive == 1 {
		t.Errorf("imbue hit: coefficient %v, additive multiplier %v, want 0.1 and Elemental Weapons", imbue.BonusCoefficient, imbue.DamageMultiplierAdditive)
	}
}

// Searing Totem attacks every 2.435 sec on beta logs (2.2 sec cast plus ~0.23 sec between casts, 1,267
// attacks), so rank 6's 55 sec totem lands 22 attacks, not the 25 a bare 2.2 sec cast gives.
func TestSearingTotemAttackInterval(t *testing.T) {
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid: core.SinglePlayerRaidProto(&proto.Player{
			Name: "enh", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc, TalentsString: DefaultTalents,
			Equipment: &proto.EquipmentSpec{},
			Spec:      &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{ClassOptions: &proto.ShamanOptions{}}}},
			Rotation:  &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter: core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()

	enh := sim.Raid.Parties[0].Players[0].(*EnhancementShaman)
	dot := enh.SearingTotem.Dot(enh.CurrentTarget)
	if dot.BaseTickLength != 2430*time.Millisecond || dot.BaseTickCount != 22 {
		t.Fatalf("Searing Totem attacks every %v, %d times, want 2.43s and 22", dot.BaseTickLength, dot.BaseTickCount)
	}
}
