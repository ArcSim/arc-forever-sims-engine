package mage

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/stats"
)

// Forever's Arcane Blast buff (400573) as client 70425 (16 October 2026 beta) redesigned it: each stack
// gives the mage 15% of its Intellect as spell damage, for every spell including Arcane Blast, and
// raises Arcane Blast's own cost by 175%. Up to four stacks, lasting 15 sec "from time of initial
// application": a new stack does not refresh it, and no other spell spends it any more. The spell
// damage follows the mage's Intellect while the buff is up.
func (mage *Mage) registerArcaneCharges() {
	if !mage.Talents.ArcaneBlast {
		return
	}

	buffRank := spellData.ArcaneBlastTriggered.Highest()
	intellectPerStack := buffRank.EffectN(1).Percent()
	costPerStack := buffRank.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_COST)).Average(core.CharacterLevel) / 100

	spellDamage := mage.NewDynamicStatDependency(stats.Intellect, stats.SpellDamage, 0)
	costMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask: MageSpellArcaneBlast,
		Kind:      core.SpellMod_PowerCost_Pct_Add,
	})

	mage.ArcaneBlastAura = mage.RegisterAura(core.Aura{
		Label:     "Arcane Blast",
		ActionID:  core.ActionID{SpellID: buffRank.ID},
		Duration:  buffRank.Duration(),
		MaxStacks: int32(buffRank.MaxStack),
		OnGain: func(_ *core.Aura, sim *core.Simulation) {
			mage.EnableDynamicStatDep(sim, spellDamage)
			costMod.Activate()
		},
		OnExpire: func(_ *core.Aura, sim *core.Simulation) {
			mage.DisableDynamicStatDep(sim, spellDamage)
			mage.UpdateDynamicStatDep(sim, spellDamage, 0)
			costMod.Deactivate()
		},
		OnStacksChange: func(_ *core.Aura, sim *core.Simulation, _ int32, newStacks int32) {
			mage.UpdateDynamicStatDep(sim, spellDamage, intellectPerStack*float64(newStacks))
			costMod.UpdateFloatValue(costPerStack * float64(newStacks))
		},
	})
}
