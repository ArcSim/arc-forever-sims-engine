package hunter

import (
	"github.com/wowsims/forever/sim/core"
)

func (hunter *Hunter) registerMongooseBiteSpell() {
	rank := spellData.MongooseBite.Highest()
	baseDamage := rank.DamageEffect().Average(core.CharacterLevel)

	// The aura is only a pre-requisite for Mongoose Bite: a dodge opens the window (client
	// SpellAuraRestrictions CasterAuraState 1 on every rank), and Expose Prey opens it off any
	// landed hit on a marked target. The dodge watch is its own permanent aura, since an aura's
	// OnSpellHitTaken only runs while that aura is up.
	defensiveWindow := spellData.ExposePreyTriggered.Rank(1)
	hunter.DefensiveState = hunter.RegisterAura(core.Aura{
		Label:    "Defensive State",
		ActionID: core.ActionID{SpellID: defensiveWindow.ID},
		Duration: defensiveWindow.Duration(),
	})
	hunter.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Defensive State - Trigger",
		TriggerImmediately: true,
		Outcome:            core.OutcomeDodge,
		Callback:           core.CallbackOnSpellHitTaken,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			hunter.DefensiveState.Activate(sim)
		},
	})

	// A main-hand attack: it deals the main-hand weapon's damage. ProcMaskMeleeSpecial also named the
	// off hand, so a dual-wielding hunter's bite took Predator's Edge's +50% off-hand damage and
	// rolled the off-hand weapon's procs.
	hunter.MongooseBite = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ClassSpellMask: HunterSpellMongooseBite,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MaxRange:       float64(rank.MaxRange),

		ManaCost: core.ManaCostOptions{
			FlatCost: int32(rank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD(),
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: max(rank.Cooldown(), rank.CategoryCooldown()),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DefensiveState.IsActive()
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			hunter.DefensiveState.Deactivate(sim)
			// Forever: normalized melee weapon damage plus a smaller flat amount, where Classic
			// dealt the flat amount alone.
			damage := baseDamage + hunter.AutoAttacks.MH().CalculateNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target))
			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if hunter.LaceratingStrikes != nil && result.Landed() {
				hunter.procLaceratingStrikes(sim, result)
			}
		},
	})
}
