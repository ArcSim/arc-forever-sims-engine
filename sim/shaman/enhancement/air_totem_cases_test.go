package enhancement

import (
	"fmt"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// Grace of Air Totem is no longer exclusive with Windfury Totem (16 October 2026 beta notes), so the
// party's Grace of Air and Windfury Totem stack, and a totem the shaman casts stacks with the other one
// from the party. The shaman's own air slot still holds one totem: of its own casts, the last stands.
// Windfury Weapon is not an air totem: it only disables the Windfury Totem proc its own wielder would get,
// so it works beside Grace of Air.

const (
	graceOfAirCast    = 25359 // the totem spell the shaman casts
	windfuryTotemCast = 10614
	graceOfAirParty   = 25360  // the party aura, tag -1
	windfuryParty     = 25587  // the party Windfury Totem aura, tag -1
	windfuryWeaponHit = 439440 // Windfury Weapon's main-hand attack, two a proc
	partyTotemExtra   = 25584  // the extra attack tag the party Windfury Totem's proc swings under
	castTotemExtra    = 10610  // the extra attack tag a cast Windfury Totem's proc swings under
)

type airTotemRun struct {
	// Seconds each air totem aura was up on average, by what supplied it.
	partyWindfury, partyGrace, castWindfury, castGrace float64
	// Windfury Totem extra attacks a fight, and Windfury Weapon procs a fight.
	partyExtra, castExtra, weaponProcs float64
}

// casts lists the air totem spells the shaman casts in the pre-pull, in order.
func runAirTotems(t *testing.T, imbue proto.ShamanImbue, party *proto.PartyBuffs, casts ...int32) airTotemRun {
	t.Helper()
	prepull := ""
	for i, id := range casts {
		if i > 0 {
			prepull += ","
		}
		prepull += fmt.Sprintf(`{"action":{"castSpell":{"spellId":{"spellId":%d}}},"doAtValue":{"const":{"val":"-%.1fs"}}}`, id, 4.0-float64(i)*1.5)
	}
	player := &proto.Player{
		Name: "enh", Class: proto.Class_ClassShaman, Race: proto.Race_RaceOrc, TalentsString: DefaultTalents,
		Equipment: &proto.EquipmentSpec{Items: []*proto.ItemSpec{
			{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
			{Id: 12784}, // Arcanite Reaper
		}},
		Buffs: &proto.IndividualBuffs{}, Consumables: &proto.ConsumesSpec{},
		Spec: &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{
			ClassOptions: &proto.ShamanOptions{ImbueMh: imbue},
		}}},
		Rotation: core.APLRotationFromJsonString(`{"type":"TypeAPL","prepullActions":[` + prepull + `],"priorityList":[]}`),
	}
	req := &proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, party, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
		SimOptions: &proto.SimOptions{Iterations: 30, RandomSeed: 1},
	}
	res := core.RunRaidSim(req)
	if res.Error != nil {
		t.Fatal(res.Error.Message)
	}
	m := res.RaidMetrics.Parties[0].Players[0]
	var out airTotemRun
	for _, a := range m.Auras {
		up := a.UptimeSecondsAvg
		if up < 0 {
			t.Errorf("aura %v has a negative uptime of %v s", a.Id, up)
		}
		switch {
		case a.Id.GetSpellId() == windfuryParty && a.Id.Tag == -1:
			out.partyWindfury = up
		case a.Id.GetSpellId() == graceOfAirParty && a.Id.Tag == -1:
			out.partyGrace = up
		case a.Id.GetSpellId() == windfuryTotemCast && a.Id.Tag == 0 && up > out.castWindfury:
			out.castWindfury = up
		case a.Id.GetSpellId() == graceOfAirCast && a.Id.Tag == 0 && up > out.castGrace:
			out.castGrace = up
		}
	}
	for _, a := range m.Actions {
		var n float64
		for _, tgt := range a.Targets {
			n += float64(tgt.Casts)
			if a.Id.GetSpellId() == windfuryWeaponHit {
				out.weaponProcs += float64(tgt.Hits+tgt.Crits+tgt.Misses+tgt.Dodges+tgt.Parries+tgt.Blocks+tgt.Glances) / 2 / float64(res.IterationsDone)
			}
		}
		if a.Id.GetOtherId() != proto.OtherAction_OtherActionAttack {
			continue
		}
		switch a.Id.Tag {
		case partyTotemExtra:
			out.partyExtra = n / float64(res.IterationsDone)
		case castTotemExtra:
			out.castExtra = n / float64(res.IterationsDone)
		}
	}
	return out
}

const (
	rockbiter    = proto.ShamanImbue_RockbiterWeapon
	windfuryWpn  = proto.ShamanImbue_WindfuryWeapon
	totemPlayed  = 100.0 // seconds a totem must be up to count as standing the fight
	totemAbsent  = 0.5   // seconds under which a totem counts as never up
	procExpected = 1.0   // extra attacks or procs a fight that count as "procs"
)

// With the party's Windfury Totem alone the totem stands and procs; that is the baseline each case
// below changes.
func TestPartyWindfuryTotemAloneProcs(t *testing.T) {
	r := runAirTotems(t, rockbiter, &proto.PartyBuffs{WindfuryTotem: true})
	if r.partyWindfury < totemPlayed || r.partyExtra < procExpected {
		t.Fatalf("party Windfury Totem alone: up %.1f s with %.2f extra attacks a fight, want it up and proccing", r.partyWindfury, r.partyExtra)
	}
}

// Redfall and Kerani's report had a Windfury Totem from the party buffs and a Grace of Air the rotation
// casts both up. Since the 16 October 2026 notes that is right: both stand and Windfury procs.
func TestCastGraceOfAirStacksWithPartyWindfuryTotem(t *testing.T) {
	r := runAirTotems(t, rockbiter, &proto.PartyBuffs{WindfuryTotem: true}, graceOfAirCast)
	if r.castGrace < totemPlayed {
		t.Fatalf("cast Grace of Air up %.1f s, want it standing", r.castGrace)
	}
	if r.partyWindfury < totemPlayed || r.partyExtra < procExpected {
		t.Fatalf("party Windfury Totem up %.1f s with %.2f extra attacks beside a cast Grace of Air, want it up and proccing", r.partyWindfury, r.partyExtra)
	}
}

func TestCastWindfuryTotemStacksWithPartyGraceOfAir(t *testing.T) {
	r := runAirTotems(t, rockbiter, &proto.PartyBuffs{GraceOfAirTotem: true}, windfuryTotemCast)
	if r.castWindfury < totemPlayed || r.castExtra < procExpected {
		t.Fatalf("cast Windfury Totem up %.1f s with %.2f extra attacks a fight, want it up and proccing", r.castWindfury, r.castExtra)
	}
	if r.partyGrace < totemPlayed {
		t.Fatalf("party Grace of Air up %.1f s beside a cast Windfury Totem, want it standing", r.partyGrace)
	}
}

// A shaman who casts one air totem and then the other holds the second only: one air slot.
func TestCastAirTotemsReplaceEachOther(t *testing.T) {
	r := runAirTotems(t, rockbiter, &proto.PartyBuffs{}, graceOfAirCast, windfuryTotemCast)
	if r.castWindfury < totemPlayed || r.castGrace > totemAbsent {
		t.Fatalf("Grace of Air then Windfury Totem: Windfury up %.1f s, Grace of Air up %.1f s", r.castWindfury, r.castGrace)
	}
	r = runAirTotems(t, rockbiter, &proto.PartyBuffs{}, windfuryTotemCast, graceOfAirCast)
	if r.castGrace < totemPlayed || r.castWindfury > totemAbsent {
		t.Fatalf("Windfury Totem then Grace of Air: Grace up %.1f s, Windfury up %.1f s", r.castGrace, r.castWindfury)
	}
}

// A request that sets both party air totems holds both, and the twisting toggle takes nothing off Grace
// of Air when Windfury Totem is the party's as well.
func TestPartyAirTotemsStack(t *testing.T) {
	for name, party := range map[string]*proto.PartyBuffs{
		"plain":    {WindfuryTotem: true, GraceOfAirTotem: true},
		"twisting": {WindfuryTotem: true, GraceOfAirTotem: true, TotemTwisting: true},
	} {
		r := runAirTotems(t, rockbiter, party)
		if r.partyWindfury < totemPlayed || r.partyExtra < procExpected {
			t.Fatalf("%s, both party air totems: Windfury up %.1f s with %.2f extra attacks, want it up and proccing", name, r.partyWindfury, r.partyExtra)
		}
		if r.partyGrace < 0.99*r.partyWindfury {
			t.Fatalf("%s, both party air totems: Grace of Air up %.1f s against Windfury's %.1f s, want it standing the fight", name, r.partyGrace, r.partyWindfury)
		}
	}
}

// Redfall: Windfury Weapon "should be useable with Grace of Air". Both the party's and a cast Grace of
// Air stand beside it, and its procs match a Windfury Weapon with no totem at all.
func TestWindfuryWeaponWorksWithGraceOfAir(t *testing.T) {
	alone := runAirTotems(t, windfuryWpn, &proto.PartyBuffs{})
	if alone.weaponProcs < procExpected {
		t.Fatalf("Windfury Weapon alone procs %.2f times a fight", alone.weaponProcs)
	}
	for name, r := range map[string]airTotemRun{
		"party Grace of Air": runAirTotems(t, windfuryWpn, &proto.PartyBuffs{GraceOfAirTotem: true}),
		"cast Grace of Air":  runAirTotems(t, windfuryWpn, &proto.PartyBuffs{}, graceOfAirCast),
		"both":               runAirTotems(t, windfuryWpn, &proto.PartyBuffs{GraceOfAirTotem: true}, graceOfAirCast),
	} {
		if r.partyGrace < totemPlayed && r.castGrace < totemPlayed {
			t.Errorf("%s: Grace of Air is not up (%.1f s party, %.1f s cast)", name, r.partyGrace, r.castGrace)
		}
		if r.weaponProcs < 0.7*alone.weaponProcs || r.weaponProcs > 1.5*alone.weaponProcs {
			t.Errorf("%s: Windfury Weapon procs %.2f times a fight, %.2f without the totem", name, r.weaponProcs, alone.weaponProcs)
		}
	}
}

// Windfury Weapon in the main hand disables the Windfury Totem its wielder would get (a party's or a
// cast one), and only that.
func TestWindfuryWeaponStillDisablesWindfuryTotem(t *testing.T) {
	for name, r := range map[string]airTotemRun{
		"party": runAirTotems(t, windfuryWpn, &proto.PartyBuffs{WindfuryTotem: true}),
		"cast":  runAirTotems(t, windfuryWpn, &proto.PartyBuffs{}, windfuryTotemCast),
	} {
		if r.partyExtra != 0 || r.castExtra != 0 {
			t.Errorf("%s Windfury Totem beside Windfury Weapon swung %.2f / %.2f extra attacks, want none", name, r.partyExtra, r.castExtra)
		}
		if r.weaponProcs < procExpected {
			t.Errorf("%s Windfury Totem beside Windfury Weapon: the weapon procs %.2f times a fight", name, r.weaponProcs)
		}
	}
}
