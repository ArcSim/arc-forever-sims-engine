package forever

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

// Items whose client 70425 rows the item generator cannot read, hand-written so it skips them.
func init() {
	// Blade of Hanna (2801): "Each melee hit grants you a stack of Improving Form. At 30 stacks perfect your
	// form, inflicting 490 Physical damage to your current target and increasing all your stats by 20 for
	// 20s." 1325535 hears main-hand white swings and adds a stack of Improving Form (1325509, up to 30),
	// whose A_TRIGGER_SPELL_ON_STACK_AMOUNT casts Perfected Form (1325508) at 30; Perfected Form deals the
	// damage, grants the stats and removes the stacks. The generator read 1325535 as a proc of Perfected
	// Form on every swing, which dealt the 490 each time.
	core.NewItemEffect(2801, func(agent core.Agent) {
		character := agent.GetCharacter()
		perfected := spelldata.MustFind(1325508)
		improving := spelldata.MustFind(1325509)
		damage := perfected.EffectN(1).BasePoints
		allStats := perfected.EffectN(2).BasePoints

		buff := character.RegisterAura(core.Aura{
			Label:    "Perfected Form",
			ActionID: core.ActionID{SpellID: perfected.ID},
			Duration: perfected.Duration(),
		}).AttachStatsBuff(stats.Stats{
			stats.Strength:  allStats,
			stats.Agility:   allStats,
			stats.Stamina:   allStats,
			stats.Intellect: allStats,
			stats.Spirit:    allStats,
		})

		strike := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: perfected.ID},
			SpellSchool: core.SpellSchoolPhysical,
			DefenseType: core.DefenseTypeMelee,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeSpecialHitAndCrit)
			},
		})

		stacks := character.RegisterAura(core.Aura{
			Label:     "Improving Form",
			ActionID:  core.ActionID{SpellID: improving.ID},
			Duration:  core.NeverExpires,
			MaxStacks: int32(improving.MaxStack),
		})

		character.MakeProcTriggerAura(core.ProcTrigger{
			Name:     "Blade of Hanna",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskMeleeMHAuto,
			Outcome:  core.OutcomeLanded,
			// A stack per swing; nothing about it needs the spell batch window.
			TriggerImmediately: true,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				stacks.Activate(sim)
				stacks.AddStack(sim)
				if stacks.GetStacks() < stacks.MaxStacks {
					return
				}
				stacks.Deactivate(sim)
				strike.Cast(sim, result.Target)
				buff.Activate(sim)
			},
		})
	})

	// Lord Sakrasis' Scepter (5028): Holy Smite (1326019) casts 1312330, 56 Holy damage, off harmful spells
	// and abilities, and its row states no chance. A proc on every cast is not credible for an item level
	// 47 off-hand and nothing yet states the rate, so it is not modelled. Registering it keeps the item
	// generator from reading it as a proc on every cast.
	core.NewItemEffect(5028, func(core.Agent) {})
}
