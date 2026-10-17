package hunter

import (
	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/stats"
	"strconv"
	"time"
)

func (hunter *Hunter) registerAspectOfTheBeastSpell() {
	spellId := [5]int32{0, 13161, 1299445, 1299446, 1299447}
	manaCost := [5]float64{0, 50, 70, 90, 110}
	apBonus := [5]float64{0, 50, 70, 90, 110}
	level := [5]int{0, 30, 40, 50, 60}

	deadlyChance := 0.02 * float64(hunter.Talents.DeadlyAspects)
	var meleeHasteAura *core.Aura
	if hunter.Talents.DeadlyAspects > 0 {
		meleeHasteAura = hunter.GetOrRegisterAura(core.Aura{
			Label: "Quick Strikes Beast",
			// the Forever Quick Strikes buff
			ActionID: core.ActionID{SpellID: 1299448},
			Duration: time.Second * 12,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.MultiplyMeleeSpeed(sim, 1.3)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.MultiplyMeleeSpeed(sim, 1/1.3)
			},
		})
	}

	maxRank := 1
	for r := 4; r >= 1; r-- {
		if level[r] <= int(hunter.Level) {
			maxRank = r
			break
		}
	}

	meleeAP := apBonus[maxRank]
	actionID := core.ActionID{SpellID: spellId[maxRank]}

	beastAura := hunter.GetOrRegisterAura(core.Aura{
		Label:    "Aspect of the Beast",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.AttackPower, meleeAP)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.AttackPower, -meleeAP)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if meleeHasteAura == nil || !result.Landed() {
				return
			}
			// melee auto attacks, Raptor Strike included (it replaces a swing)
			if !spell.ProcMask.Matches(core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto) {
				return
			}
			if sim.Proc(deadlyChance, "Deadly Aspects Beast") {
				meleeHasteAura.Activate(sim)
			}
		},
	})
	beastAura.NewExclusiveEffect("Aspect", true, core.ExclusiveEffect{})

	hunter.GetOrRegisterSpell(core.SpellConfig{
		ActionID:      actionID,
		Flags:         core.SpellFlagAPL,
		Rank:          maxRank,
		RequiredLevel: level[maxRank],
		ManaCost:      core.ManaCostOptions{FlatCost: manaCost[maxRank]},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{GCD: core.GCDDefault},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !beastAura.IsActive()
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			beastAura.Activate(sim)
		},
	})
}

// Aspect of the Cheetah (Forever client data, spell 5118): +30% run speed for 40 mana, Pathfinding adds 3% a point
// (Forever's 19559 adds 6 at 2/2). The daze when struck is left out: a boss you are not tanking does not hit you.
// Like other run speed bonuses it does not stack with Strider Kick; the stronger one counts.
func (hunter *Hunter) registerAspectOfTheCheetahSpell() {
	actionID := core.ActionID{SpellID: 5118}
	speed := 1.30 + 0.03*float64(hunter.Talents.Pathfinding)

	cheetahAura := hunter.GetOrRegisterAura(core.Aura{
		Label:    "Aspect of the Cheetah",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddMoveSpeedModifier(&aura.ActionID, speed)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.RemoveMoveSpeedModifier(&aura.ActionID)
		},
	})
	cheetahAura.NewExclusiveEffect("Aspect", true, core.ExclusiveEffect{})

	hunter.GetOrRegisterSpell(core.SpellConfig{
		ActionID:      actionID,
		Flags:         core.SpellFlagAPL,
		RequiredLevel: 20,
		ManaCost:      core.ManaCostOptions{FlatCost: 40},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{GCD: core.GCDDefault},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !cheetahAura.IsActive()
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			cheetahAura.Activate(sim)
		},
	})
}

// Utility function to create an Improved Hawk Aura
func (hunter *Hunter) createImprovedHawkAura(auraLabel string, actionID core.ActionID) *core.Aura {
	bonusMultiplier := 1.3
	return hunter.GetOrRegisterAura(core.Aura{
		Label:    auraLabel,
		ActionID: actionID,
		Duration: time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyRangedSpeed(sim, bonusMultiplier)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyRangedSpeed(sim, 1/bonusMultiplier)
		},
	})
}

// Function to get the maximum attack power for Aspect of the Hawk based on rank
func (hunter *Hunter) getMaxAspectOfTheHawkAttackPower(rank int) float64 {
	attackPower := [8]float64{0, 20, 35, 50, 70, 90, 110, 120}

	if rank < 1 || rank > 7 {
		return 0.0
	}

	return attackPower[rank]
}

func (hunter *Hunter) getMaxHawkRank() int {
	maxRank := core.TernaryInt(core.IncludeAQ, 7, 6)

	for i := maxRank; i > 0; i-- {
		config := hunter.getAspectOfTheHawkSpellConfig(i)
		if config.RequiredLevel <= int(hunter.Level) {
			return i
		}
	}
	return 1
}

func (hunter *Hunter) getAspectOfTheHawkSpellConfig(rank int) core.SpellConfig {
	var impHawkAura *core.Aura
	improvedHawkProcChance := 0.02 * float64(hunter.Talents.DeadlyAspects)

	spellIds := [8]int32{0, 13165, 14318, 14319, 14320, 14321, 14322, 25296}
	levels := [8]int{0, 10, 18, 28, 38, 48, 58, 60}

	spellId := spellIds[rank]
	level := levels[rank]

	if hunter.Talents.DeadlyAspects > 0 {
		impHawkAura = hunter.QuickShotsAura
	}

	rap := hunter.getMaxAspectOfTheHawkAttackPower(rank)

	actionID := core.ActionID{SpellID: spellId}
	aspectOfTheHawkAura := hunter.GetOrRegisterAura(core.Aura{
		Label:    "Aspect of the Hawk" + strconv.Itoa(rank),
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.RangedAttackPower, rap*hunter.AspectOfTheHawkAPMultiplier)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.RangedAttackPower, -rap*hunter.AspectOfTheHawkAPMultiplier)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}
			if !spell.ProcMask.Matches(core.ProcMaskRangedAuto) {
				return
			}
			if impHawkAura != nil && sim.Proc(improvedHawkProcChance, "Deadly Aspects") {
				impHawkAura.Activate(sim)
			}
		},
	})

	aspectOfTheHawkAura.NewExclusiveEffect("Aspect", true, core.ExclusiveEffect{})

	return core.SpellConfig{
		ActionID:      actionID,
		Flags:         core.SpellFlagAPL,
		Rank:          rank,
		RequiredLevel: level,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !aspectOfTheHawkAura.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			aspectOfTheHawkAura.Activate(sim)
		},
	}
}

func (hunter *Hunter) registerAspectOfTheHawkSpell() {
	hunter.AspectOfTheHawkAPMultiplier = 1.0
	if hunter.Talents.DeadlyAspects > 0 {
		hunter.QuickShotsAura = hunter.createImprovedHawkAura(
			"Quick Shots",
			core.ActionID{SpellID: 6150},
		)
	}
	maxRank := hunter.getMaxHawkRank()
	config := hunter.getAspectOfTheHawkSpellConfig(maxRank)
	hunter.GetOrRegisterSpell(config)
}
