package spelldata

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// Every number a tooltip on the evidence page states (ui/sim/spells/*.json) has to be one the client
// store can produce for that row: an effect's amount at level 60, its low or high roll, a tick times the
// tick count, a talent's value at some rank, a cost, a duration. The daily Update DB rebuilds the store,
// so a hotfix that moves a number the page still quotes fails here rather than going quiet.
//
// A number the store cannot produce is either wrong or a rule the client does not carry (a threat
// multiplier, a proc chance from a beta tooltip, a measured coefficient). Those are listed in
// tooltipAllowed with the numbers alone, so a changed number on the same row still fails.
var tooltipAllowed = map[int32][]float64{
	// Client 1.60.1.70291 cut the low ranks' gain per level to fractions (#758). The tooltip rounds the gain,
	// the store truncates it, so these rows read one point apart; beta logs should settle which the server does.
	585:     {14},                                     // Smite rank 1
	591:     {26},                                     // Smite rank 2
	686:     {16},                                     // Shadow Bolt rank 1
	5143:    {25},                                     // Arcane Missiles rank 2
	5145:    {49},                                     // Arcane Missiles rank 4
	10301:   {6},                                      // Retribution Aura: spell power share, beta notes 8 Oct 2026
	14893:   {25},                                     // Inspiration: Classic's 25%, the client buff says 8% (assumed row)
	17794:   {20},                                     // Improved Shadow Bolt: 4% a point, 20% at five
	20346:   {50, 61},                                 // Judgement of Light: 50% chance not in the client (assumed row)
	20355:   {50, 59},                                 // Judgement of Wisdom: same
	22959:   {15, 33},                                 // Improved Scorch: 3% x5 stacks, 33% a point
	408255:  {2},                                      // Eclipse: "next 2 Starfire spells"
	417053:  {20},                                     // Natural Reaction: 20% a point
	5229:    {27},                                     // Enrage: the armor cut lives on the form auras
	9896:    {3, 4, 243, 396, 549, 702, 855},          // Rip: per combo point totals over 6 ticks
	1242512: {2, 25},                                  // Rapid Recuperation: rank 2 wording
	10187:   {1168},                                   // Blizzard: ticks come from an area trigger
	10216:   {332},                                    // Flamestrike: the dot comes from an area trigger
	12484:   {15, 25, 40},                             // Improved Blizzard: chill per rank
	400669:  {17, 33, 50},                             // Fingers of Frost: Shatter's crit per rank
	412538:  {5, 8, 40},                               // Ignite: 8% a point
	1310925: {3, 10, 33, 66},                          // Shield Specialization: per point wording
	1311033: {5},                                      // Iron Creed: per point
	1311074: {20, 33},                                 // Sanctified Judgement: per point
	20116:   {56},                                     // Consecration: area trigger ticks
	20922:   {48, 88},                                 // Consecration
	20923:   {160},                                    // Consecration
	20924:   {96, 216},                                // Consecration
	20925:   {20},                                     // Holy Shield: the 20% extra threat is scripted; the 70170 block chance (30) no longer matches it
	20928:   {20},                                     // Holy Shield: the 20% extra threat is scripted; the 70170 block chance (30) no longer matches it
	20128:   {2, 6},                                   // Redoubt: per point
	20178:   {8, 20},                                  // Reckoning: per point
	24239:   {20},                                     // Hammer of Wrath: the 20% health condition
	24274:   {20},                                     // Hammer of Wrath
	24275:   {20},                                     // Hammer of Wrath
	9799:    {5, 10, 50},                              // Eye for an Eye: per rank, and the 50% cap
	1284536: {5},                                      // Holy Purpose: per rank
	11275:   {8, 12, 14, 16, 159, 222, 295, 377, 469}, // Rupture: per combo point durations and totals
	1241584: {50},                                     // Mutilate: 67 per hand before the 75%
	1310710: {3, 5, 10, 15},                           // Puncturing Wounds: per rank
	13976:   {3, 33, 67},                              // Initiative: per rank
	16257:   {5, 20, 25},                              // Flurry (shaman): per rank
	16344:   {35, 112},                                // Flametongue Weapon: by weapon speed
	16387:   {54},                                     // Flametongue Totem: by weapon speed (1363 / 25)
	408345:  {443},                                    // Fire Nova: damage is on 408428, via hand triggers
	408505:  {4},                                      // Maelstrom Weapon: per point
	1293696: {7.8, 65, 68},                            // Demonic Brand: assumed row
	17926:   {460},                                    // Death Coil: health leech effect
	23785:   {2},                                      // Master Demonologist: per point
	412732:  {33, 67},                                 // Demonic Knowledge: per rank
	426311:  {2},                                      // Shadow and Flame: per point
	440873:  {3, 20, 35, 45},                          // Decimation: per point
	12319:   {5},                                      // Flurry (warrior): assumed row
	12966:   {5},                                      // Flurry buff: same
	12880:   {2, 30},                                  // Enrage: per point and proc chance
	1290261: {3},                                      // Weaponmaster: per point
	12964:   {2, 12},                                  // Unbridled Wrath: per point, two-hander rage
	1310222: {80},                                     // Spearing Strike: bonus against creature types
	1310318: {20},                                     // Shield Specialization: per point
	1318070: {5, 6},                                   // Diamond Flask: channel and cooldown
	23925:   {50},                                     // Shield Slam: dispel chance
	871:     {5.5},
	12579:   {20},                         // Winter's Chill: 20% a point to proc
	24858:   {35},                         // Moonkin Form: the cost is a percentage of base mana the store rounds away
	9634:    {50, 1240, 55},               // Dire Bear Form: threat is not in the client; health is on 9635
	19577:   {8},                          // Intimidation: base mana share
	1310703: {3, 4, 5, 9, 12, 15, 18, 21}, // Venom: per combo point durations
	16488:   {20},                         // Blood Craze: the 20% health condition
	2687:    {25},                         // Bloodrage: Improved Bloodrage per point                                    // Shield Wall: Improved Shield Wall per point
}

var tooltipNumber = regexp.MustCompile(`\d+(?:\.\d+)?`)

func TestManifestTooltipsMatchTheClient(t *testing.T) {
	defer generatedStore()()

	files, err := filepath.Glob(filepath.Join("..", "..", "..", "ui", "sim", "spells", "*.json"))
	if err != nil || len(files) == 0 {
		t.Fatalf("no evidence manifest: %v", err)
	}
	checked := 0
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		var rows map[string]struct {
			Tooltip string `json:"tooltip"`
		}
		if err := json.Unmarshal(data, &rows); err != nil {
			t.Fatalf("%s: %v", file, err)
		}
		for key, row := range rows {
			id, err := strconv.Atoi(key)
			if err != nil || row.Tooltip == "" || Find(int32(id)) == Nil {
				continue
			}
			checked++
			known := producible(Find(int32(id)))
			for _, v := range tooltipAllowed[int32(id)] {
				known[v] = true
			}
			var bad []string
			for _, m := range tooltipNumber.FindAllString(strings.ReplaceAll(row.Tooltip, ",", ""), -1) {
				v, _ := strconv.ParseFloat(m, 64)
				if v > 1 && !known[round2(v)] {
					bad = append(bad, m)
				}
			}
			if len(bad) > 0 {
				sort.Strings(bad)
				t.Errorf("%s %d: the tooltip states %s, which the client store does not produce for this row. Fix the tooltip, or add the number to tooltipAllowed with the reason.", filepath.Base(file), id, strings.Join(bad, ", "))
			}
		}
	}
	if checked < 300 {
		t.Fatalf("checked only %d tooltips; is the manifest where this test expects it?", checked)
	}
}

// Every number a tooltip could fairly quote for the row, rounded the ways tooltips round.
func producible(root *Spell) map[float64]bool {
	var raw []float64
	seen := map[int32]bool{}
	var walk func(s *Spell, depth int)
	walk = func(s *Spell, depth int) {
		if s == Nil || seen[s.ID] || depth > 2 {
			return
		}
		seen[s.ID] = true
		raw = append(raw, s.Cost(), ms(s.CastTimeMs), ms(s.DurationMs), ms(s.CooldownMs), ms(s.CategoryCooldownMs), ms(s.ICDMs), float64(s.ProcChance), float64(s.MaxStack), float64(s.MinRange), float64(s.MaxRange), ms(s.GCDMs), float64(s.ProcCharges))
		for _, p := range s.Powers {
			raw = append(raw, float64(p.CostPct), float64(p.Cost)/10)
		}
		ticks := 1.0
		if s.DurationMs > 0 {
			for _, e := range s.Effects {
				if e.PeriodMs > 0 {
					ticks = math.Max(ticks, float64(s.DurationMs/e.PeriodMs))
				}
			}
		}
		for i := range s.Effects {
			e := &s.Effects[i]
			avg := e.Average(60)
			for _, v := range []float64{e.BasePoints, avg, e.Min(60), e.Max(60), float64(e.PointsPerResource), ms(e.PeriodMs), float64(e.ChainTargets), float64(e.ChainAmp), float64(e.RadiusMax)} {
				raw = append(raw, v, v*ticks)
			}
			for cp := 1.0; cp <= 5; cp++ {
				raw = append(raw, (avg+cp*float64(e.PointsPerResource))*ticks)
			}
			if e.Trigger() != nil {
				walk(e.Trigger(), depth+1)
			}
		}
		for _, curve := range curves[s.ID] {
			raw = append(raw, curve...)
		}
		for _, r := range s.Refs() {
			walk(r, depth+1)
		}
		for _, r := range s.Triggered() {
			walk(r, depth+1)
		}
	}
	walk(root, 0)

	// A bonus stated as a share of weapon damage scales the flat part too: 80 at 155% reads as 124.
	var pcts []float64
	for _, v := range raw {
		if v > 100 && v <= 400 {
			pcts = append(pcts, v/100)
		}
	}
	for _, v := range append([]float64(nil), raw...) {
		for _, p := range pcts {
			raw = append(raw, v*p)
		}
	}

	out := map[float64]bool{}
	for _, v := range raw {
		v = math.Abs(v)
		for _, w := range []float64{v, v / 1000, v * 100, v / 100, v / 60, v / 10} {
			for _, r := range []float64{math.Floor(w + 1e-9), math.Ceil(w - 1e-9), math.Round(w), round2(w), math.Round(w*10) / 10} {
				out[round2(r)] = true
			}
		}
	}
	return out
}

func ms(v int32) float64 { return float64(v) / 1000 }

func round2(v float64) float64 { return math.Round(v*100) / 100 }
