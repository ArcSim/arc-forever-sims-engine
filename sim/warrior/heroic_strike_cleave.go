package warrior

import (
	"github.com/wowsims/forever/sim/core"
)

// HeroicStrikeCleaveThreatPerLevel takes Heroic Strike's and Cleave's flat threat from upstream's
// in-game captures (mythicsim patch 184, ported from wowsims/forever's 2026-10-07 overrides): Heroic
// Strike adds 2.54 threat a spell level beyond its damage (Rank 3: 40 to 40.6 at level 16, Rank 4: 61
// at 24), none with attack power, and Cleave 1.53 a level (Rank 2: ~46 at level 30 on one target,
// ~23 on each of two), split between the targets it lands on. Both are multiplied by the stance with
// the damage. Off, Classic's 173 and 100, Cleave's on every target.
var HeroicStrikeCleaveThreatPerLevel = true

const (
	heroicStrikeThreatPerLevel = 2.54
	cleaveThreatPerLevel       = 1.53
)

// Heroic Strike and Cleave carry the client's On Next Swing attribute (Attributes 0x4 on 25286 and
// 20569, which Bloodthirst and Mortal Strike lack), so both queue onto the next main hand swing
// instead of firing as instant specials.
func (warrior *Warrior) registerHeroicStrike() {
	heroicStrikeRank := spellData.HeroicStrike.Highest()
	heroicStrikeBaseDamage := heroicStrikeRank.DamageEffect().Average(core.CharacterLevel)

	spell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: heroicStrikeRank.ID},
		SpellSchool:    heroicStrikeRank.SpellSchool(),
		DefenseType:    heroicStrikeRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskMeleeMH,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,
		ClassSpellMask: SpellMaskHeroicStrike,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   int32(heroicStrikeRank.Cost()),
			Refund: heroicStrikeRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		// Not in the client table (no E_THREAT effect). Patch 184: 2.54 a spell level (152 at Rank 9).
		FlatThreatBonus: core.TernaryFloat64(HeroicStrikeCleaveThreatPerLevel,
			heroicStrikeThreatPerLevel*float64(heroicStrikeRank.SpellLevel), 173),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := heroicStrikeBaseDamage + warrior.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
			result := warrior.calcQueuedSwing(sim, spell, target, baseDamage)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}

			spell.DealDamage(sim, result)
			if warrior.curQueueAura != nil {
				warrior.curQueueAura.Deactivate(sim)
			}
		},
	})
	warrior.makeQueueSpellsAndAura(spell)
}

func (warrior *Warrior) registerCleave() {
	cleaveRank := spellData.Cleave.Highest()
	cleaveBaseDamage := cleaveRank.DamageEffect().Average(core.CharacterLevel)

	const maxTargets int32 = 2
	results := make(core.SpellResultSlice, 0, maxTargets)

	spell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: cleaveRank.ID},
		SpellSchool:    cleaveRank.SpellSchool(),
		DefenseType:    cleaveRank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskMeleeMH,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,
		ClassSpellMask: SpellMaskCleave,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   int32(cleaveRank.Cost()),
			Refund: cleaveRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		// Not in the client table (no E_THREAT effect). Patch 184: 1.53 a spell level (92 at Rank 5),
		// split between the targets the cleave lands on.
		FlatThreatBonus: core.TernaryFloat64(HeroicStrikeCleaveThreatPerLevel,
			cleaveThreatPerLevel*float64(cleaveRank.SpellLevel), 100),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results = results[:0]
			numTargets := min(maxTargets, sim.Environment.ActiveTargetCount())
			for range numTargets {
				baseDamage := cleaveBaseDamage + warrior.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
				results = append(results, warrior.calcQueuedSwing(sim, spell, target, baseDamage))
				target = sim.Environment.NextActiveTargetUnit(target)
			}

			if HeroicStrikeCleaveThreatPerLevel {
				landed := 0
				for _, result := range results {
					if result.Landed() {
						landed++
					}
				}
				if landed > 1 {
					flatThreat := spell.FlatThreatBonus
					spell.FlatThreatBonus /= float64(landed)
					for _, result := range results {
						result.Threat = spell.ThreatFromDamage(sim, result.Outcome, result.Damage, warrior.AttackTables[result.Target.UnitIndex])
					}
					spell.FlatThreatBonus = flatThreat
				}
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
			}

			if warrior.curQueueAura != nil {
				warrior.curQueueAura.Deactivate(sim)
			}
		},
	})
	warrior.makeQueueSpellsAndAura(spell)
}

// Heroic Strike and Cleave replace the main hand swing but roll on the special attack table, so they
// skip the dual wield miss penalty. The penalty flag is character wide, so it is only lifted for the
// swing itself: an off-hand auto that lands while the queue is up still pays it.
func (warrior *Warrior) calcQueuedSwing(sim *core.Simulation, spell *core.Spell, target *core.Unit, baseDamage float64) *core.SpellResult {
	warrior.PseudoStats.DisableDWMissPenalty = true
	result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
	warrior.PseudoStats.DisableDWMissPenalty = false
	return result
}

func (warrior *Warrior) makeQueueSpellsAndAura(srcSpell *core.Spell) *core.Spell {
	isQueueQueued := false

	queueAura := warrior.RegisterAura(core.Aura{
		Label:    "HS/Cleave Queue Aura-" + srcSpell.ActionID.String(),
		ActionID: srcSpell.ActionID.WithTag(1),
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			isQueueQueued = false
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if warrior.curQueueAura != nil {
				warrior.curQueueAura.Deactivate(sim)
			}
			warrior.curQueueAura = aura
			warrior.curQueuedAutoSpell = srcSpell
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.curQueueAura = nil
			warrior.curQueuedAutoSpell = nil
		},
	})

	return warrior.RegisterSpell(core.SpellConfig{
		ActionID:    srcSpell.ActionID.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: srcSpell.DefenseType,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagNoMetrics,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.curQueueAura == nil &&
				!isQueueQueued &&
				warrior.CurrentRage() >= srcSpell.Cost.GetCurrentCost() &&
				warrior.Hardcast.Expires <= sim.CurrentTime &&
				warrior.queuedRealismICD.IsReady(sim)
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if warrior.queuedRealismICD.IsReady(sim) {
				isQueueQueued = true
				warrior.queuedRealismICD.Use(sim)
				sim.AddPendingAction(&core.PendingAction{
					NextActionAt: sim.CurrentTime + warrior.queuedRealismICD.Duration,
					OnAction: func(sim *core.Simulation) {
						queueAura.Activate(sim)
						isQueueQueued = false
					},
				})
			}
		},
	})
}

// Swaps the main hand swing for the queued Heroic Strike or Cleave, or drops the queue when it can
// no longer be paid for.
func (warrior *Warrior) TryHSOrCleave(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
	if !warrior.curQueueAura.IsActive() {
		return mhSwingSpell
	}

	if !warrior.curQueuedAutoSpell.CanCast(sim, warrior.CurrentTarget) {
		warrior.curQueueAura.Deactivate(sim)
		return mhSwingSpell
	}

	return warrior.curQueuedAutoSpell
}
