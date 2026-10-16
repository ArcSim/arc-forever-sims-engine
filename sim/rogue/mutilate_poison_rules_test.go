package rogue

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
)

// Assassination checks asked for on the MythicSim Discord (2 October 2026): what the poisons scale with,
// whether they crit, and what Mutilate rolls. The expectations are the beta client's rows (build
// 1.60.1.69893, read with tools/data_watch/spell_client.py) and what the engine code does with them:
//
//   - Instant Poison VI (damage spell 11337): Nature, 88 base with a 27.7% spread (76 to 100), spell power
//     coefficient 0 and attack power coefficient 0, no "cannot crit" attribute.
//   - Deadly Poison V (dot 25349): Nature, 23 a tick every 3 sec for 4 ticks, both coefficients 0, "periodic
//     can crit" set.
//   - Mutilate (parent 1310707 to 1241584) awards 2 Combo Points and triggers two hit spells, one a hand, that
//     are each 75% weapon damage plus a flat rank amount. Seal Fate (14186): "Your critical strikes from
//     abilities that add Combo Points have a chance to add an additional Combo Point", with a 0.5 sec
//     internal cooldown on the proc row.

func newRogueProbe(t *testing.T, consumables *proto.ConsumesSpec) (*core.Simulation, *Rogue, *core.Unit) {
	t.Helper()
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceHuman,
		Class:         proto.Class_ClassRogue,
		Equipment:     daggersOnly(),
		Consumables:   consumables,
		TalentsString: AssassinationTalents,
		Rotation:      &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
	}, DefaultOptions)
	// A fight long enough that the thousands of casts below never reach its end.
	encounter := core.MakeSingleTargetEncounter(0)
	encounter.Duration = 1_000_000
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 7},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  encounter,
	}, simsignals.CreateSignals())
	sim.Reset()
	rogue := sim.Raid.Parties[0].Players[0].(RogueAgent).GetRogue()
	return sim, rogue, rogue.CurrentTarget
}

// advance runs the sim forward by d, whatever else is pending.
func advance(sim *core.Simulation, d time.Duration) {
	until := sim.CurrentTime + d
	sim.AddPendingAction(&core.PendingAction{NextActionAt: until, OnAction: func(*core.Simulation) {}})
	for sim.CurrentTime < until {
		if sim.Step() {
			panic("the fight ended while the probe was still running")
		}
	}
}

func TestMutilateRollsEachHandAndSealFateAddsOnePoint(t *testing.T) {
	sim, rogue, target := newRogueProbe(t, &proto.ConsumesSpec{})
	if rogue.Mutilate == nil || rogue.Talents.SealFate != 5 {
		t.Fatalf("the Assassination build has Mutilate %v and Seal Fate %d/5", rogue.Mutilate != nil, rogue.Talents.SealFate)
	}
	// Plenty of crits, so the double crit casts are common enough to count.
	rogue.AddStatDynamic(sim, stats.PhysicalCritPercent, 25)

	id := strconv.Itoa(int(MutilateSpellID))
	hit := func(tag string) *regexp.Regexp {
		return regexp.MustCompile(`\{SpellID: ` + id + `, Tag: ` + tag + `\} (Crit|Hit|Miss|Dodge|Parry|Block|Glance|Crush)`)
	}
	mainHand, offHand := hit("1"), hit("2")
	sealFate := regexp.MustCompile(`Gained 1 combo points from \{SpellID: 14195\}`)
	parentPoints := regexp.MustCompile(`Gained 2 combo points from \{SpellID: ` + id + `\}`)

	var lines []string
	sim.Log = func(format string, args ...interface{}) { lines = append(lines, fmt.Sprintf(format, args...)) }
	energy := rogue.NewEnergyMetrics(core.ActionID{OtherID: proto.OtherAction_OtherActionEnergyRegen})
	points := rogue.NewComboPointMetrics(core.ActionID{OtherID: proto.OtherAction_OtherActionRefund})

	var casts, mhCrits, ohCrits, bothCrits, anyCrit, extraOnAnyCrit, extraOnBoth, extraWithoutCrit, mostExtra int
	for i := 0; i < 4000; i++ {
		advance(sim, 1400*time.Millisecond) // the global cooldown and the Seal Fate cooldown are both over
		if missing := rogue.MaximumEnergy() - rogue.CurrentEnergy(); missing > 0 {
			rogue.AddEnergy(sim, missing, energy)
		}
		rogue.SpendPartialComboPoints(sim, rogue.ComboPoints(), points)
		lines = lines[:0]
		if !rogue.Mutilate.Cast(sim, target) {
			t.Fatal("Mutilate did not cast")
		}
		// A proc handler runs one spell batch window after the hit that triggered it.
		advance(sim, 100*time.Millisecond)
		var mh, oh string
		extra, parent := 0, false
		for _, line := range lines {
			if m := mainHand.FindStringSubmatch(line); m != nil {
				mh = m[1]
			}
			if m := offHand.FindStringSubmatch(line); m != nil {
				oh = m[1]
			}
			if sealFate.MatchString(line) {
				extra++
			}
			if parentPoints.MatchString(line) {
				parent = true
			}
		}
		if !parent {
			continue // the parent strike missed: no hits, no points
		}
		if mh == "" || oh == "" {
			t.Fatalf("Mutilate landed but a hand rolled nothing (main hand %q, off hand %q)", mh, oh)
		}
		casts++
		mhCrit, ohCrit := mh == "Crit", oh == "Crit"
		if mhCrit {
			mhCrits++
		}
		if ohCrit {
			ohCrits++
		}
		mostExtra = max(mostExtra, extra)
		switch {
		case mhCrit && ohCrit:
			bothCrits++
			anyCrit++
			extraOnBoth += extra
			extraOnAnyCrit += extra
		case mhCrit || ohCrit:
			anyCrit++
			extraOnAnyCrit += extra
		default:
			extraWithoutCrit += extra
		}
	}

	pMH, pOH, pBoth := float64(mhCrits)/float64(casts), float64(ohCrits)/float64(casts), float64(bothCrits)/float64(casts)
	t.Logf("%d landed Mutilates: main hand crit %.3f, off hand crit %.3f, both %.3f (independent rolls predict %.3f)", casts, pMH, pOH, pBoth, pMH*pOH)
	t.Logf("Seal Fate points: %d on %d casts with a crit, %d on the %d double crits, %d on casts with no crit; most on one cast: %d",
		extraOnAnyCrit, anyCrit, extraOnBoth, bothCrits, extraWithoutCrit, mostExtra)

	if bothCrits < 100 {
		t.Fatalf("only %d double crits; the sample is too thin", bothCrits)
	}
	// Two rolls: the double crit rate is the product of the two hands' rates (about 6% either way here).
	if math.Abs(pBoth-pMH*pOH) > 0.02 {
		t.Errorf("double crit rate %.3f against %.3f if the hands rolled independently", pBoth, pMH*pOH)
	}
	if math.Abs(pMH-pOH) > 0.04 {
		t.Errorf("main hand crit %.3f and off hand crit %.3f differ more than the same table should", pMH, pOH)
	}
	if extraWithoutCrit != 0 {
		t.Errorf("%d Seal Fate points on casts that did not crit", extraWithoutCrit)
	}
	// 5/5 Seal Fate is a certain extra point on a crit, and the 0.5 sec cooldown keeps a double crit
	// at one extra point, since both hands land at the same instant.
	if mostExtra != 1 {
		t.Errorf("the most Seal Fate points on one Mutilate was %d, want 1", mostExtra)
	}
	if extraOnAnyCrit != anyCrit {
		t.Errorf("Seal Fate added %d points over %d casts with a crit, want one each", extraOnAnyCrit, anyCrit)
	}
}

// Poisons gain a share of attack power the client rows do not carry: 0.5% for Instant Poison and 0.45% a
// stack a tick for Deadly Poison (upstream #781, from the devs, forever-bugs #100). Spell power adds nothing,
// and both crit.
func TestPoisonsScaleWithAttackPowerOnlyAndCrit(t *testing.T) {
	sim, rogue, target := newRogueProbe(t, &proto.ConsumesSpec{MhImbueId: 26891, OhImbueId: 27186}) // Instant Poison, Deadly Poison
	if rogue.InstantPoison == nil || rogue.DeadlyPoison == nil {
		t.Fatal("the rogue has no Instant Poison or Deadly Poison")
	}

	instant := regexp.MustCompile(`\{SpellID: 11340\} \[DEBUG\] MAP: (-?[0-9.]+), RAP: (-?[0-9.]+), SP: (-?[0-9.]+), BaseDamage:(-?[0-9.]+), AfterAttackerMods:(-?[0-9.]+)`)
	instantOutcome := regexp.MustCompile(`\{SpellID: 11340\} (Crit|Hit)`)
	var lines []string
	sim.Log = func(format string, args ...interface{}) { lines = append(lines, fmt.Sprintf(format, args...)) }

	// Instant Poison: 76 to 100 plus 0.5% of attack power, whatever the spell power, and it crits.
	type sample struct{ meanBase, minBase, maxBase, critRate, meanAP float64 }
	instantSample := func() sample {
		s := sample{minBase: math.MaxFloat64}
		n, crits := 0, 0
		for i := 0; i < 3000; i++ {
			advance(sim, 100*time.Millisecond)
			lines = lines[:0]
			rogue.InstantPoison.Cast(sim, target)
			for _, line := range lines {
				if m := instant.FindStringSubmatch(line); m != nil {
					ap, _ := strconv.ParseFloat(m[1], 64)
					base, _ := strconv.ParseFloat(m[4], 64)
					s.meanBase += base
					s.meanAP += ap
					s.minBase, s.maxBase = min(s.minBase, base), max(s.maxBase, base)
					n++
				}
				if m := instantOutcome.FindStringSubmatch(line); m != nil && m[1] == "Crit" {
					crits++
				}
			}
		}
		s.meanBase /= float64(n)
		s.meanAP /= float64(n)
		s.critRate = float64(crits) / float64(n)
		return s
	}
	before := instantSample()
	rogue.AddStatDynamic(sim, stats.AttackPower, 2000)
	rogue.AddStatDynamic(sim, stats.SpellDamage, 2000)
	after := instantSample()
	t.Logf("Instant Poison: base %.1f to %.1f (mean %.2f) at %.0f attack power; %.1f to %.1f (mean %.2f) at %.0f; crit rate %.3f then %.3f; spell crit chance %.3f",
		before.minBase, before.maxBase, before.meanBase, before.meanAP, after.minBase, after.maxBase, after.meanBase, after.meanAP, before.critRate, after.critRate,
		rogue.InstantPoison.SpellCritChance(target))
	if after.meanAP < before.meanAP+1900 {
		t.Fatalf("the attack power never rose (%.0f then %.0f)", before.meanAP, after.meanAP)
	}
	const instantShare, deadlyShare = 0.005, 0.0045
	for _, s := range []sample{before, after} {
		if lo, hi := 76+instantShare*s.meanAP, 100+instantShare*s.meanAP; s.minBase < lo-0.1 || s.maxBase > hi+0.1 {
			t.Errorf("Instant Poison's base damage %.1f to %.1f at %.0f attack power, want %.1f to %.1f", s.minBase, s.maxBase, s.meanAP, lo, hi)
		}
	}
	if want := instantShare * (after.meanAP - before.meanAP); math.Abs(after.meanBase-before.meanBase-want) > 1.5 {
		t.Errorf("Instant Poison's mean base damage moved %.2f -> %.2f with 2000 attack power and spell power, want +%.2f", before.meanBase, after.meanBase, want)
	}
	if before.critRate < 0.02 || math.Abs(before.critRate-rogue.InstantPoison.SpellCritChance(target)) > 0.04 {
		t.Errorf("Instant Poison crit rate %.3f against a spell crit chance of %.3f", before.critRate, rogue.InstantPoison.SpellCritChance(target))
	}

	// Deadly Poison: 23 plus 0.45% of attack power a tick a stack, and its ticks crit too.
	tick := regexp.MustCompile(`\{SpellID: 25347, Tag: 100\} \[DEBUG\] MAP: (-?[0-9.]+), RAP: (-?[0-9.]+), SP: (-?[0-9.]+), BaseDamage:(-?[0-9.]+), AfterAttackerMods:`)
	tickOutcome := regexp.MustCompile(`\{SpellID: 25347, Tag: 100\} tick (Crit|Hit)`)
	var bases, perStack []float64
	var tickCrits, ticks int
	lines = lines[:0]
	rogue.DeadlyPoison.Cast(sim, target)
	for i := 0; i < 400000 && len(bases) < 400; i++ {
		sim.Step()
		for _, line := range lines {
			if m := tick.FindStringSubmatch(line); m != nil {
				ap, _ := strconv.ParseFloat(m[1], 64)
				base, _ := strconv.ParseFloat(m[4], 64)
				bases = append(bases, base)
				perStack = append(perStack, 23+deadlyShare*ap)
			}
			if m := tickOutcome.FindStringSubmatch(line); m != nil {
				ticks++
				if m[1] == "Crit" {
					tickCrits++
				}
			}
		}
		lines = lines[:0]
		if !rogue.DeadlyPoison.Dot(target).IsActive() {
			rogue.DeadlyPoison.Cast(sim, target)
		}
	}
	if len(bases) < 100 {
		t.Fatalf("only %d Deadly Poison ticks logged", len(bases))
	}
	for i, base := range bases {
		if stacks := base / perStack[i]; math.Abs(stacks-math.Round(stacks)) > 0.01 {
			t.Errorf("a Deadly Poison tick of %.2f is not a whole number of %.2f-damage stacks, so something other than stacks and attack power moved it", base, perStack[i])
			break
		}
	}
	t.Logf("Deadly Poison: %d ticks, base damage always a whole number of stacks with %.0f extra attack power and %.0f spell power up; %d of %d ticks crit (%.3f)",
		len(bases), 2000.0, 2000.0, tickCrits, ticks, float64(tickCrits)/math.Max(1, float64(ticks)))
	if tickCrits == 0 {
		t.Error("no Deadly Poison tick crit")
	}
}
