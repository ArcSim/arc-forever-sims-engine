package hunter

import (
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
)

// Level 30 hawk logs (150 hawk attacks): a hawk follows the hunter's ranged attack speed (2.27 s
// with a 13% quiver, 2.23 s with a 15% one; 40% faster under Rapid Fire, 30% under Quick Shots, unchanged
// under the melee proc Quick Strikes); its attacks resolve like a pet's white melee (glancing, dodged, missed) and
// crit about 1% of the time for double damage; each attack is 15-18 raw at level 30 and gains about 0.5% of ranged AP.
const (
	// before ranged haste (quiver, Rapid Fire, Quick Shots)
	hawkSwingSpeed = 2570 * time.Millisecond
	hawkCritChance = 0.01
	// level 60: the level 30 hawk's ~14.7 base scaled by ExpectedStat CreatureAutoAttackDps (24.69 at 30, 60.45 at
	// 60). Forever's creatures follow that table: a level 30 tiger in the logs had its level 30 armor (1212), and the
	// level 63 boss its 3731. Assumes the hawk is the hunter's level.
	hawkSwingBaseDamage = 36.0
	hawkSwingAPShare    = 0.005
)

func (hunter *Hunter) hawkSwingInterval() time.Duration {
	speed := hawkSwingSpeed
	return time.Duration(float64(speed) / hunter.PseudoStats.RangedSpeedMultiplier)
}

// one hawk attack before armor and outcome; roll in [0, 1) spreads it 91-109%
func hawkSwingBase(roll, rangedAP float64) float64 {
	return hawkSwingBaseDamage*(0.91+0.18*roll) + hawkSwingAPShare*rangedAP
}

// the dive: rank damage + 5% of ranged AP; Unleashed Fury adds a flat 3 per point to it (Forever data, a flat
// modifier on Summon Hawk's damage effect; logs: 66 raw at 373 ranged AP, 72 at 482)
func hawkDiveBase(rankDamage float64, unleashedFury int32, rangedAP float64) float64 {
	return rankDamage + 3*float64(unleashedFury) + 0.05*rangedAP
}

func (hunter *Hunter) registerSummonHawkSpell(sharedTimer *core.Timer) {
	if !hunter.Talents.SummonHawk {
		return
	}
	if sharedTimer == nil {
		sharedTimer = hunter.NewTimer()
	}

	// client ranks: levels 25 / 36 / 48 / 60
	baseDamage, manaCost := 32.0, 80.0
	if hunter.Level >= 60 {
		baseDamage, manaCost = 108, 190
	} else if hunter.Level >= 48 {
		baseDamage, manaCost = 85, 135
	} else if hunter.Level >= 36 {
		baseDamage, manaCost = 47, 105
	}

	// Ferocity: +2% crit per point for pets and hawks
	hawkCrit := hawkCritChance + 0.02*float64(hunter.Talents.Ferocity)

	var swingTarget *core.Unit
	expireAt := []time.Duration{0, 0}
	alive := []bool{false, false}
	generation := []int{0, 0}

	// the hawk's white melee by the target's level (the attack table formulas, from the hunter's level with no
	// weapon skill): miss, dodge, glancing; then the hawk's own crit chance for double damage
	var hawkSwing *core.Spell
	hawkOutcome := func(sim *core.Simulation, result *core.SpellResult, _ *core.AttackTable) {
		metrics := &hawkSwing.SpellMetrics[result.Target.UnitIndex]
		diff := float64(result.Target.Level*5 - hunter.Level*5)
		miss := 0.05 + diff*0.001
		if diff > 10 {
			miss = 0.05 + diff*0.002
		}
		dodge := 0.05 + diff*0.001
		glance := 0.1 + diff*0.02
		roll := sim.RandomFloat("Hawk Hit Table")
		switch {
		case roll < miss:
			result.Outcome = core.OutcomeMiss
			result.Damage = 0
			metrics.Misses++
		case roll < miss+dodge:
			result.Outcome = core.OutcomeDodge
			result.Damage = 0
			metrics.Dodges++
		case roll < miss+dodge+glance:
			result.Outcome = core.OutcomeGlance
			low := max(min(1.3-0.05*diff, 0.91), 0.01)
			high := max(min(1.2-0.03*diff, 0.99), 0.2)
			result.Damage *= low + sim.RandomFloat("Hawk Glancing Penalty")*(high-low)
			metrics.Glances++
		case roll < miss+dodge+glance+hawkCrit:
			result.Outcome = core.OutcomeCrit
			result.Damage *= 2
			metrics.Crits++
		default:
			result.Outcome = core.OutcomeHit
			metrics.Hits++
		}
	}

	hawkSwing = hunter.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1293241}.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskRangedSpecial,
		// a guardian's attack: the hunter's own damage talents (Lone Wolf, Improved Tracking) do not apply
		Flags:    core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell | core.SpellFlagIgnoreAttackerModifiers,
		CastType: proto.CastType_CastTypeMainHand,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if swingTarget == nil {
				return
			}
			base := hawkSwingBase(sim.RandomFloat("Hawk Damage"), spell.RangedAttackPower(swingTarget, false))
			spell.CalcAndDealDamage(sim, swingTarget, base, hawkOutcome)
		},
	})

	hawkAura := hunter.RegisterAura(core.Aura{
		Label:     "Summon Hawk",
		ActionID:  core.ActionID{SpellID: 1293241}.WithTag(2),
		Duration:  core.NeverExpires,
		MaxStacks: 2,
		// every iteration starts with no hawks out (old expiry times would block both slots all fight)
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			for i := range expireAt {
				expireAt[i] = 0
				alive[i] = false
				generation[i]++
			}
		},
	})

	hunter.SummonHawk = hunter.GetOrRegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 1293241},
		SpellSchool:  core.SpellSchoolPhysical,
		DefenseType:  core.DefenseTypeRanged,
		ProcMask:     core.ProcMaskRangedSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagAPL | SpellFlagShot,
		CastType:     proto.CastType_CastTypeRanged,
		MissileSpeed: 35,

		ManaCost: core.ManaCostOptions{FlatCost: manaCost},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{GCD: core.GCDDefault},
			IgnoreHaste: true,
			// shares its cooldown with Arcane Shot
			CD: core.Cooldown{
				Timer:    sharedTimer,
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// the dive always hits
			initial := hawkDiveBase(baseDamage, hunter.Talents.UnleashedFury, spell.RangedAttackPower(target, false))
			spell.CalcAndDealDamage(sim, target, initial, spell.OutcomeRangedCritOnly)

			swingTarget = target

			for i := 0; i < 2; i++ {
				if expireAt[i] > sim.CurrentTime {
					continue
				}
				slot := i
				expireAt[slot] = sim.CurrentTime + 18*time.Second
				// a re-summon can land before the old hawk's swing timer noticed it expired: the new hawk
				// takes over the slot (and its stack), the old timer just stops
				generation[slot]++
				myGeneration := generation[slot]
				if !alive[slot] {
					alive[slot] = true
					if !hawkAura.IsActive() {
						hawkAura.Activate(sim)
					}
					hawkAura.AddStack(sim)
				}

				var pa *core.PendingAction
				pa = &core.PendingAction{
					NextActionAt: sim.CurrentTime + hunter.hawkSwingInterval(),
					Priority:     core.ActionPriorityDOT,
					OnAction: func(sim *core.Simulation) {
						if generation[slot] != myGeneration {
							return
						}
						if sim.CurrentTime > expireAt[slot] || swingTarget == nil {
							expireAt[slot] = 0
							alive[slot] = false
							if hawkAura.IsActive() {
								hawkAura.RemoveStack(sim)
								if hawkAura.GetStacks() == 0 {
									hawkAura.Deactivate(sim)
								}
							}
							return
						}
						hawkSwing.Cast(sim, swingTarget)
						pa.NextActionAt = sim.CurrentTime + hunter.hawkSwingInterval()
						sim.AddPendingAction(pa)
					},
				}
				sim.AddPendingAction(pa)
				break
			}
		},
	})
}
