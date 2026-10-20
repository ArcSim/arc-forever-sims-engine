package hunter

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// Patch 178: a Hunter pet takes its master's hit, the higher of melee and melee plus ranged (wowsims/forever
// a205077834; Hunter Pet Scaling, 415429, lists hit). Melee hit raised mid-fight reaches the pet, and so does
// ranged hit until it comes off again; with the switch off the pet keeps its own hit.
func TestPetsInheritHit(t *testing.T) {
	const arcaniteReaper = 12784
	near := func(a, b float64) bool { return math.Abs(a-b) < 1e-9 }
	run := func(on bool) (start, melee, ranged float64) {
		defer func(old bool) { PetsInheritHit = old }(PetsInheritHit)
		PetsInheritHit = on
		sim, hunter, pet := petRangedCritSim(t, proto.Race_RaceOrc, BeastMasteryTalents, arcaniteReaper)
		start = pet.GetStat(stats.PhysicalHitPercent)
		owner := hunter.GetStats()
		hunter.AddStatDynamic(sim, stats.PhysicalHitPercent, 3)
		pet.AddOwnerStats(sim, hunter.GetStats().Subtract(owner))
		melee = pet.GetStat(stats.PhysicalHitPercent) - start
		owner = hunter.GetStats()
		hunter.AddStatDynamic(sim, stats.RangedHitPercent, 2)
		pet.AddOwnerStats(sim, hunter.GetStats().Subtract(owner))
		ranged = pet.GetStat(stats.PhysicalHitPercent) - start - melee
		// Taking the ranged hit off again takes it off the pet.
		owner = hunter.GetStats()
		hunter.AddStatDynamic(sim, stats.RangedHitPercent, -2)
		pet.AddOwnerStats(sim, hunter.GetStats().Subtract(owner))
		if got := pet.GetStat(stats.PhysicalHitPercent) - start; !near(got, melee) {
			t.Errorf("switch %v: after the ranged hit came off the pet holds %.3f, want %.3f", on, got, melee)
		}
		return start, melee, ranged
	}
	_, melee, ranged := run(true)
	if !near(melee, 3) || !near(ranged, 2) {
		t.Errorf("switch on: +3 melee hit gives the pet %.3f, +2 ranged hit %.3f; want 3 and 2", melee, ranged)
	}
	_, melee, ranged = run(false)
	if !near(melee, 0) || !near(ranged, 0) {
		t.Errorf("switch off: the pet took %.3f melee and %.3f ranged hit, want none", melee, ranged)
	}
}
