package hunter

import (
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
)

func (hunter *Hunter) getMongooseBiteConfig(rank int) core.SpellConfig {
	spellId := [5]int32{0, 1495, 14269, 14270, 14271}[rank]
	bonus := [5]float64{0, 15, 22, 37, 57}[rank]
	manaCost := [5]float64{0, 30, 40, 50, 65}[rank]
	level := [5]int{0, 16, 30, 44, 58}[rank]

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolPhysical,
		DefenseType:   core.DefenseTypeMelee,
		ProcMask:      core.ProcMaskMeleeMHSpecial,
		Flags:         core.SpellFlagMeleeMetrics | core.SpellFlagAPL | SpellFlagStrike,
		CastType:      proto.CastType_CastTypeMainHand,
		Rank:          rank,
		RequiredLevel: level,

		ManaCost: core.ManaCostOptions{FlatCost: manaCost},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{GCD: core.GCDDefault},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 5,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget <= 5 && (hunter.DefensiveState.IsActive() || hunter.ExposePrey.IsActive())
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			base := bonus + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			result := spell.CalcAndDealDamage(sim, target, base, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			// one bite spends the ready state, whichever buff gave it
			if result.Landed() {
				if hunter.ExposePrey.IsActive() {
					hunter.exposePreyUsed = true
					hunter.ExposePrey.Deactivate(sim)
				}
				hunter.DefensiveState.Deactivate(sim)
			}
		},
	}
}

func (hunter *Hunter) registerMongooseBiteSpell() {
	// dodges and parries (nothing attacks the hunter in a raid sim, so this one stays down there)
	hunter.DefensiveState = hunter.RegisterAura(core.Aura{
		Label:    "Defensive State",
		ActionID: core.ActionID{SpellID: 5302},
		Duration: time.Second * 10,
	})
	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: "Mongoose Bite Trigger",
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeDodge | core.OutcomeParry) {
				hunter.DefensiveState.Activate(sim)
			}
		},
	}))

	// Expose Prey's own buff (1310726, 10 s), with short markers so the buff table counts (procs per fight) and the
	// timeline shows what became of each one: used by a Mongoose Bite, run out unused, or refreshed by a new proc
	// while it was still up
	marker := func(tag int32, label string) *core.Aura {
		return hunter.RegisterAura(core.Aura{Label: label, ActionID: core.ActionID{SpellID: 1310726, Tag: tag}, Duration: time.Millisecond * 100})
	}
	used := marker(1, "Expose Prey Used")
	expired := marker(2, "Expose Prey Expired")
	refreshed := marker(3, "Expose Prey Refreshed")
	hunter.ExposePrey = hunter.RegisterAura(core.Aura{
		Label:    "Expose Prey",
		ActionID: core.ActionID{SpellID: 1310726},
		Duration: time.Second * 10,
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			// the fight's end also takes it down, and that wastes nothing
			if hunter.exposePreyUsed {
				used.Activate(sim)
			} else if sim.CurrentTime < sim.Duration {
				expired.Activate(sim)
			}
			hunter.exposePreyUsed = false
		},
	})

	if hunter.Talents.ExposePrey > 0 {
		procChance := 0.05 * float64(hunter.Talents.ExposePrey)
		core.MakePermanent(hunter.RegisterAura(core.Aura{
			Label: "Expose Prey Talent",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				// the hunter's own melee and ranged attacks only (not the hawk's swings or traps)
				if !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) || spell.ActionID == (core.ActionID{SpellID: 1293241}.WithTag(1)) {
					return
				}
				// only on a target carrying a Hunter's Mark, any hunter's (the raid debuff setting counts as another hunter's)
				if !result.Landed() || !result.Target.HasActiveAuraWithTag(core.HuntersMarkAuraTag) {
					return
				}
				if sim.RandomFloat("Expose Prey") >= procChance {
					return
				}
				if hunter.ExposePrey.IsActive() {
					refreshed.Activate(sim)
					hunter.ExposePrey.Refresh(sim)
				} else {
					hunter.ExposePrey.Activate(sim)
				}
			},
		}))
	}

	best := 0
	for rank := 1; rank <= 4; rank++ {
		if [5]int{0, 16, 30, 44, 58}[rank] <= int(hunter.Level) {
			best = rank
		}
	}
	if best == 0 {
		return
	}
	hunter.MongooseBite = hunter.GetOrRegisterSpell(hunter.getMongooseBiteConfig(best))
}
