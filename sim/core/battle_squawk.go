package core

import (
	"math"
	"time"
)

// EnableBattleSquawk applies Party.battle_squawks (mythicsim patch 127). Off, the field is ignored.
var EnableBattleSquawk = true

// MaxBattleSquawks is the most Battle Squawk stacks a party can assume: one per party member.
const MaxBattleSquawks = 5

// Battle Squawk (23060), the buff an engineer's Gnomish Battle Chicken (item 10725, summon 13166) gives
// its party: 5% melee attack speed for 4 minutes. Forever's client states the aura as melee haste
// (SpellEffect 699978: aura 319, 5 points; Era's 1.15.9 row has aura 138, 4 + 1 die), and its buff
// tooltip reads "Attack speed increased by 5%." Each chicken's squawk is its own aura, so a party
// with several engineers stacks them, multiplying as separate haste auras do. The squawks are taken
// as cast at the pull and last their 4 minutes; a longer fight loses them, as nobody can summon a
// chicken again within its 30 minute cooldown.
func registerBattleSquawk(character *Character, stacks int32) {
	if !EnableBattleSquawk || stacks <= 0 {
		return
	}
	stacks = min(stacks, MaxBattleSquawks)
	const perStack = 1.05

	character.GetOrRegisterAura(Aura{
		Label:     "Battle Squawk",
		ActionID:  ActionID{SpellID: 23060},
		Duration:  time.Minute * 4,
		MaxStacks: MaxBattleSquawks,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
			aura.SetStacks(sim, stacks)
		},
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.MultiplyMeleeSpeed(sim, math.Pow(perStack, float64(newStacks-oldStacks)))
		},
	})
}
