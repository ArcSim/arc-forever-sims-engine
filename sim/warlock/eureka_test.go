package warlock

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// Eureka! 1259821's three effects as client 1.60.1.70170 states them (EffectSpellClassMask_0..3, family 5).
// The engine's spell store does not carry the racial rows, so the masks are copied here; the 70124 masks
// were {542127, 8388994, 0, 0} on the first two (cost, damage) and {17422, 258, 0, 0} on the third.
var (
	eurekaCost     = core.ClassFlags{Family: 5, Mask: [4]uint32{541165, 8650880, 0, 0}}
	eurekaDamage   = core.ClassFlags{Family: 5, Mask: [4]uint32{524773, 8388736, 0, 0}}
	eurekaPeriodic = core.ClassFlags{Family: 5, Mask: [4]uint32{16392, 262144, 0, 0}}
)

// Hellfire is a channel here whose ticks are the hits of the client's Hellfire Effect, so it sits on the
// periodic list where the client puts the hit on the damage effect.
var eurekaDealtAsTicks = map[string]bool{"Hellfire": true}

// The lists Eureka! applies are the client's: every warlock spell the sim registers is on them exactly when
// the row's class masks name it. Client 1.60.1.70170 took Corruption, Bane of Agony, Bane of Doom and
// Unstable Affliction off every effect and Immolate and the drains' periodic bonus off effect 2.
func TestEurekaListsMatchTheClientRow(t *testing.T) {
	lists := (&Warlock{}).EurekaSpells()
	checked := 0
	for _, talents := range []string{AfflictionTalents, DestructionTalents} {
		_, warlock := newClient70170Warlock(t, proto.Race_RaceGnome, talents)
		for _, spell := range warlock.Spellbook {
			client := spelldata.Find(spell.ActionID.SpellID)
			if client == nil || client.ClassFlags.IsZero() || spell.ClassSpellMask == 0 {
				continue
			}
			// Immolate's dot is registered as its own spell on the hit's id (tag 1): the dot half answers to
			// the periodic effect, which no longer names Immolate.
			if spell.ClassSpellMask&WarlockSpellImmolateDot != 0 {
				if spell.ClassSpellMask&(lists.Cost|lists.Damage|lists.Tick) != 0 {
					t.Errorf("Immolate's dot is on an Eureka! list")
				}
				continue
			}
			checked++
			name := client.Name
			onCost, onDamage, onTick := spell.ClassSpellMask&lists.Cost != 0, spell.ClassSpellMask&lists.Damage != 0, spell.ClassSpellMask&lists.Tick != 0
			clientCost := eurekaCost.Matches(client.ClassFlags)
			clientDamage := eurekaDamage.Matches(client.ClassFlags)
			clientPeriodic := eurekaPeriodic.Matches(client.ClassFlags)

			if onCost != clientCost {
				t.Errorf("%s %d: on the sim's cost list = %v, client row says %v", name, client.ID, onCost, clientCost)
			}
			if (onDamage || onTick) != (clientDamage || clientPeriodic) {
				t.Errorf("%s %d: bonus damage in the sim = %v, client row says %v", name, client.ID, onDamage || onTick, clientDamage || clientPeriodic)
			}
			if clientPeriodic && !onTick {
				t.Errorf("%s %d: the client's periodic effect names it and the sim's tick list does not", name, client.ID)
			}
			if clientDamage && !clientPeriodic && !onDamage && !eurekaDealtAsTicks[name] {
				t.Errorf("%s %d: the client's damage effect names it and the sim's damage list does not", name, client.ID)
			}
			if eurekaDealtAsTicks[name] && !onTick {
				t.Errorf("%s %d: dealt as ticks here and missing from the tick list", name, client.ID)
			}
		}
	}
	if checked < 20 {
		t.Fatalf("only %d warlock spells were comparable with the client rows", checked)
	}
}

// Corruption casts don't cost less any more, and its ticks keep their damage.
func TestEurekaPassesOverTheDots(t *testing.T) {
	sim, warlock := newClient70170Warlock(t, proto.Race_RaceGnome, AfflictionTalents)
	corruption := warlock.Corruption
	shadowBolt := warlock.GetSpell(core.ActionID{SpellID: spellData.ShadowBolt.Highest().ID})
	drain := warlock.GetSpell(core.ActionID{SpellID: spellData.DrainLife.Highest().ID})

	corruptionCost, boltCost := corruption.Cost.GetCurrentCost(), shadowBolt.Cost.GetCurrentCost()
	corruptionMult := corruption.DamageMultiplier
	drainMult := drain.DamageMultiplier

	aura := warlock.GetAura("Eureka!")
	aura.Activate(sim)
	if got := corruption.Cost.GetCurrentCost(); got != corruptionCost {
		t.Errorf("Corruption costs %v under Eureka!, want %v", got, corruptionCost)
	}
	if got := corruption.DamageMultiplier; got != corruptionMult {
		t.Errorf("Corruption's damage multiplier is %v under Eureka!, want %v", got, corruptionMult)
	}
	if got, want := shadowBolt.Cost.GetCurrentCost(), boltCost*0.9; got < want-1 || got > want+1 {
		t.Errorf("Shadow Bolt costs %v under Eureka!, want %v", got, want)
	}
	// Drain Life's ticks are the channel's, which stay on the periodic effect.
	if got := drain.DamageMultiplier / drainMult; got < 1.1-1e-9 || got > 1.1+1e-9 {
		t.Errorf("Drain Life's damage multiplier rose by %v under Eureka!, want 1.1", got)
	}

}
