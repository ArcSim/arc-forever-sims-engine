package hunter

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
)

// 16 October 2026 beta notes: "Resourcefulness now reduces the mana cost of Strider Kick" (440529's cost
// mask left it out), and "Resourcefulness now only triggers on critical strikes from melee and ranged
// abilities, not traps or auto attacks".
func TestResourcefulness(t *testing.T) {
	sim, hunter, _ := petRangedCritSim(t, proto.Race_RaceNightElf, SurvivalMeleeTalents, 12784)
	if hunter.Talents.Resourcefulness != 2 || hunter.StriderKick == nil {
		t.Fatalf("test setup: the build has Resourcefulness %d, Strider Kick %v", hunter.Talents.Resourcefulness, hunter.StriderKick != nil)
	}

	cut := spellData.Resourcefulness.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_COST)).FractionAt(2)
	if got := hunter.StriderKick.Cost.AdditivePercentModifier - 1; !core.WithinToleranceFloat64(cut, got, 1e-9) {
		t.Errorf("Strider Kick cost modifier %.4f, want Resourcefulness's %.4f", got, cut)
	}

	trigger := hunter.GetAura("Resourcefulness Trigger")
	buff := hunter.GetAura("Resourcefulness")
	if trigger == nil || buff == nil {
		t.Fatal("no Resourcefulness auras")
	}
	for _, tc := range []struct {
		name  string
		spell *core.Spell
		procs bool
	}{
		{"Auto Shot", hunter.AutoAttacks.RangedAuto(), false},
		{"melee swing", hunter.AutoAttacks.MHAuto(), false},
		{"Immolation Trap", hunter.ImmolationTrap, false},
		{"Explosive Trap", hunter.ExplosiveTrap, false},
		{"Mongoose Bite", hunter.MongooseBite, true},
		{"Strider Kick", hunter.StriderKick, true},
		{"Arcane Shot", hunter.ArcaneShot, true},
		{"Multi-Shot", hunter.MultiShot, true},
	} {
		if tc.spell == nil {
			t.Errorf("%s: no spell", tc.name)
			continue
		}
		buff.Deactivate(sim)
		result := &core.SpellResult{Target: hunter.CurrentTarget, Outcome: core.OutcomeCrit}
		trigger.OnSpellHitDealt(trigger, sim, tc.spell, result)
		// The handler runs a spell batch later.
		for at := sim.CurrentTime + core.SpellBatchWindow; sim.CurrentTime < at; {
			sim.Step()
		}
		if got := buff.IsActive(); got != tc.procs {
			t.Errorf("%s crit: Resourcefulness up %v, want %v", tc.name, got, tc.procs)
		}
	}
}
