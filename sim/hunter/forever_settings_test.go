package hunter

import (
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

// TestForeverSettings runs a UI settings export (IndividualSimSettings JSON, e.g. the browser's
// localStorage currentSettings) the way the browser would, then prints the per-action table.
// Knobs: AB_SETTINGS (file path, required), AB_SETTINGS_APL (rotation file base name in
// ui/hunter/apls to swap in), AB_TALENTS (talents string to swap in), AB_ITERS (default 5000).
func TestForeverSettings(t *testing.T) {
	path := os.Getenv("AB_SETTINGS")
	if path == "" {
		t.Skip("AB_SETTINGS not set")
	}
	// comparison switches: the walk-in swing at 5 yd off, melee restarting Auto Shot off, spell haste on the GCD off
	if os.Getenv("AB_NO_GCD_HASTE") != "" {
		core.ForeverSpellHasteGCD = false
	}
	if os.Getenv("AB_NO_MIDWALK") != "" {
		core.ForeverMeleeMidWalk = false
	}
	if os.Getenv("AB_NO_RESTART") != "" {
		core.ForeverMeleeRestartsAutoShot = false
	}
	if os.Getenv("AB_NO_MOVE_POLL") != "" {
		core.ForeverMovePoll = false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	settings := &proto.IndividualSimSettings{}
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(data, settings); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	player := settings.Player
	if player == nil {
		t.Fatal("no player in settings")
	}
	if v := os.Getenv("AB_TALENTS"); v != "" {
		player.TalentsString = v
	}
	label := "settings " + path
	if os.Getenv("AB_WORLDBUFFS") != "" {
		abWorldBuffs(player)
		label += " + world buffs"
	}
	if v := os.Getenv("AB_SETTINGS_APL"); v != "" {
		player.Rotation = core.GetAplRotation("../../ui/hunter/apls", v).Rotation
		label += " + rotation " + v
	}
	iters := int32(5000)
	if v := os.Getenv("AB_ITERS"); v != "" {
		var n int32
		fmt.Sscanf(v, "%d", &n)
		if n > 0 {
			iters = n
		}
	}
	if os.Getenv("AB_EP") != "" {
		abPrintEP(settings, iters)
		return
	}
	rsr := &proto.RaidSimRequest{
		Raid:       core.SinglePlayerRaidProto(player, settings.PartyBuffs, settings.RaidBuffs, settings.Debuffs),
		Encounter:  settings.Encounter,
		SimOptions: &proto.SimOptions{Iterations: iters, RandomSeed: 20261012},
	}
	if os.Getenv("AB_DEBUG") != "" {
		// one logged iteration: the event log goes to stdout for grepping
		rsr.SimOptions.Debug = true
		rsr.SimOptions.Iterations = 1
		iters = 1
	}
	if v := os.Getenv("AB_SEED"); v != "" {
		fmt.Sscanf(v, "%d", &rsr.SimOptions.RandomSeed)
	}
	res := core.RunRaidSim(rsr)
	if res.Error != nil {
		t.Fatalf("sim error: %s", res.Error.Message)
	}
	if res.Logs != "" {
		fmt.Println(res.Logs)
	}
	abPrint(label, res.RaidMetrics.Parties[0].Players[0], iters)
	if os.Getenv("AB_IDS") != "" {
		abPrintIDs("player", res.RaidMetrics.Parties[0].Players[0])
	}
	// AB_RES=1: every resource row (mana ticks, potions, talents) as gained per fight
	if os.Getenv("AB_RES") != "" {
		for _, r := range res.RaidMetrics.Parties[0].Players[0].Resources {
			fmt.Printf("RES %v %s events %.1f gain %.1f actual %.1f\n", r.Type, r.Id.String(), float64(r.Events)/float64(iters), r.Gain/float64(iters), r.ActualGain/float64(iters))
		}
	}
	fmt.Printf("  raid DPS incl. pets %.1f\n", res.RaidMetrics.Dps.Avg)
}

// abPrintIDs lists every action, aura and resource source id in the results (the unit and its pets), one per line,
// for scripts/check_ids.js: each must be a spell or item the site can look up.
func abPrintIDs(who string, um *proto.UnitMetrics) {
	print := func(kind string, id *proto.ActionID) {
		if id == nil {
			return
		}
		fmt.Printf("ID %s %s spell=%d item=%d other=%d tag=%d\n", who, kind, id.GetSpellId(), id.GetItemId(), id.GetOtherId(), id.Tag)
	}
	for _, a := range um.Actions {
		print("action", a.Id)
	}
	for _, a := range um.Auras {
		print("aura", a.Id)
	}
	for _, r := range um.Resources {
		print("resource", r.Id)
	}
	for _, pet := range um.Pets {
		abPrintIDs("pet:"+pet.Name, pet)
	}
}

// abPrint prints DPS and the per-action table (damage share, casts, outcomes, ticks per fight).
func abPrint(label string, pm *proto.UnitMetrics, iters int32) {
	fmt.Printf("\n### %s\nDPS %.1f (stdev %.1f)  oom at %.1fs\n", label, pm.Dps.Avg, pm.Dps.Stdev, pm.SecondsOomAvg)
	type row struct {
		name                                                        string
		dmg                                                         float64
		casts, hits, crits, misses, dodges, glances, parries, ticks int32
		critTicks                                                   int32
	}
	var rows []row
	var total float64
	for _, a := range pm.Actions {
		r := row{}
		id := a.Id
		if id != nil {
			if sp := id.GetSpellId(); sp != 0 {
				if n, ok := abSpellNames[sp]; ok {
					r.name = n
				} else {
					r.name = fmt.Sprintf("spell %d", sp)
				}
			} else if it := id.GetItemId(); it != 0 {
				r.name = fmt.Sprintf("item %d", it)
			} else {
				r.name = fmt.Sprintf("other %d", id.GetOtherId())
			}
			if id.Tag != 0 {
				r.name += fmt.Sprintf(" t%d", id.Tag)
			}
		}
		for _, tg := range a.Targets {
			r.dmg += tg.Damage
			r.casts += tg.Casts
			r.hits += tg.Hits
			r.crits += tg.Crits
			r.misses += tg.Misses
			r.dodges += tg.Dodges
			r.glances += tg.Glances
			r.parries += tg.Parries
			// Ticks and CritTicks split the ticks; the Resisted counts are subsets of them
			r.ticks += tg.Ticks + tg.CritTicks
			r.critTicks += tg.CritTicks
		}
		total += r.dmg
		if r.dmg > 0 || (os.Getenv("AB_ALLROWS") != "" && (r.casts > 0 || r.hits > 0)) {
			rows = append(rows, r)
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].dmg > rows[j].dmg })
	fi := float64(iters)
	for i, r := range rows {
		if i >= 16 {
			break
		}
		fmt.Printf("  %-18s %5.1f%%  casts %5.1f hits %5.1f crit %5.1f miss %4.1f dodge %4.1f glance %4.1f ticks %5.1f (crit %4.1f) dmg/fight %7.0f\n",
			r.name, 100*r.dmg/total, float64(r.casts)/fi, float64(r.hits)/fi, float64(r.crits)/fi,
			float64(r.misses)/fi, float64(r.dodges)/fi, float64(r.glances)/fi, float64(r.ticks)/fi, float64(r.critTicks)/fi, r.dmg/fi)
	}
	for _, pet := range pm.Pets {
		fmt.Printf("  pet %s: %.1f dps\n", pet.Name, pet.Dps.Avg)
	}
	if os.Getenv("AB_MANA") != "" {
		// mana gained (+) and spent (-) per source, per fight
		type res struct {
			name         string
			gain, actual float64
			events       int32
		}
		var rs []res
		for _, r := range pm.Resources {
			if r.Type != proto.ResourceType_ResourceTypeMana {
				continue
			}
			name := ""
			if sp := r.Id.GetSpellId(); sp != 0 {
				if n, ok := abSpellNames[sp]; ok {
					name = n
				} else {
					name = fmt.Sprintf("spell %d", sp)
				}
			} else if it := r.Id.GetItemId(); it != 0 {
				name = fmt.Sprintf("item %d", it)
			} else {
				name = fmt.Sprintf("other %d t%d", r.Id.GetOtherId(), r.Id.Tag)
			}
			rs = append(rs, res{name, r.Gain / fi, r.ActualGain / fi, r.Events / iters})
		}
		sort.Slice(rs, func(i, j int) bool { return rs[i].actual > rs[j].actual })
		var plus, minus float64
		for _, r := range rs {
			if r.actual > 0 {
				plus += r.actual
			} else {
				minus += r.actual
			}
			fmt.Printf("  mana %-20s %8.0f\n", r.name, r.actual)
		}
		fmt.Printf("  mana gained %.0f, spent %.0f per fight\n", plus, -minus)
	}
}

// AB_EP=1: the site's stat weights (EP against Ranged Attack Power) for the stats the hunter page weighs, as the
// Calculate button would give them, so the default EP list can come from a real build.
func abPrintEP(settings *proto.IndividualSimSettings, iters int32) {
	statsToWeigh := []proto.Stat{proto.Stat_StatStrength, proto.Stat_StatAgility, proto.Stat_StatIntellect, proto.Stat_StatAttackPower,
		proto.Stat_StatRangedAttackPower, proto.Stat_StatMeleeHit, proto.Stat_StatMeleeCrit, proto.Stat_StatExpertise, proto.Stat_StatSpellPower,
		proto.Stat_StatSpellDamage, proto.Stat_StatNaturePower, proto.Stat_StatArcanePower, proto.Stat_StatSpellCrit, proto.Stat_StatMP5}
	pseudoToWeigh := []proto.PseudoStat{proto.PseudoStat_PseudoStatMainHandDps, proto.PseudoStat_PseudoStatOffHandDps, proto.PseudoStat_PseudoStatRangedDps,
		proto.PseudoStat_PseudoStatMeleeSpeedMultiplier, proto.PseudoStat_PseudoStatRangedSpeedMultiplier}
	res := core.StatWeights(&proto.StatWeightsRequest{
		Player:             settings.Player,
		RaidBuffs:          settings.RaidBuffs,
		PartyBuffs:         settings.PartyBuffs,
		Debuffs:            settings.Debuffs,
		Encounter:          settings.Encounter,
		SimOptions:         &proto.SimOptions{Iterations: iters, RandomSeed: 20261012},
		StatsToWeigh:       statsToWeigh,
		PseudoStatsToWeigh: pseudoToWeigh,
		EpReferenceStat:    proto.Stat_StatRangedAttackPower,
	})
	if res.Error != nil {
		fmt.Println("stat weights error:", res.Error.Message)
		return
	}
	ep := res.Dps.EpValues
	for _, s := range statsToWeigh {
		fmt.Printf("EP %s %.2f\n", s, ep.Stats[s])
	}
	for _, s := range pseudoToWeigh {
		fmt.Printf("EP %s %.2f\n", s, ep.PseudoStats[s])
	}
}
