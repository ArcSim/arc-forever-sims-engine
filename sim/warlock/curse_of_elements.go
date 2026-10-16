package warlock

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/buffs"
)

// Forever folded Curse of Shadow into Curse of the Elements: its top rank (1311680) drops every
// magic resistance and raises all magic damage taken by 10%, which is what core's aura already
// does. Malediction is no longer a bonus on the curse (the 16 October 2026 beta notes replace the
// talent with Grim Reach, range only), so no ranks are passed in.
func (warlock *Warlock) registerCurseOfElements() {
	rank := spellData.CurseOfTheElements.Highest()

	// The player copy (untagged 1311680, one label per target): the APLs ask for aura 1311680 on the
	// target, and every warlock of a raid shares it, one Curse of the Elements per target being the
	// game's rule. Tagged per caster, every warlock but the raid's first never saw its own curse and
	// recast it every GCD.
	warlock.CurseOfElementsAuras = warlock.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return buffs.CurseOfElementsAura(target, true, 0)
	})

	warlock.CurseOfElements = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellCurseOfElements,
		MaxRange:       float64(rank.MaxRange),

		ManaCost: core.ManaCostOptions{FlatCost: int32(rank.Cost())},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD(),
			},
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				aura := warlock.CurseOfElementsAuras.Get(target)
				warlock.takeCurseSlot(sim, target, aura)
				aura.Activate(sim)
			}

			spell.DealOutcome(sim, result)
		},

		RelatedAuraArrays: warlock.CurseOfElementsAuras.ToMap(),
	})
}
