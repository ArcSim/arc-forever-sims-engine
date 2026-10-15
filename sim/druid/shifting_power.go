package druid

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// Shifting Power, the Feral talent build 70170 put where King of the Jungle was: a Cat Form spell
// that turns 55% of base mana into 40 Energy on a 16 sec cooldown (client 1322605). Natural
// Shapeshifter's cost cut reaches it, and Improved Shifting Power (1322670) takes 4 sec a rank off
// the cooldown. Since the 8 Oct beta build it can't be used at full Energy (dev notes).
func (druid *Druid) registerShiftingPowerSpell() {
	if !druid.Talents.ShiftingPower {
		return
	}

	row := spelldata.Ranked(1322605).Highest()
	actionID := core.ActionID{SpellID: row.ID}
	energyMetrics := druid.NewEnergyMetrics(actionID)
	energy := row.EnergizeEffect().BaseValue()

	druid.ShiftingPower = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: DruidSpellShiftingPower,
		Flags:          core.SpellFlagAPL,

		ManaCost: row.ManaCost(),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: row.GCD(),
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: row.Cooldown(),
			},
		},

		ExtraCastCondition: func(_ *core.Simulation, _ *core.Unit) bool {
			return druid.CurrentEnergy() < druid.MaximumEnergy()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.AddEnergy(sim, energy, energyMetrics)
		},
	})

	if druid.Talents.ImprovedShiftingPower > 0 {
		druid.AddStaticMod(core.SpellModConfig{
			ClassMask: DruidSpellShiftingPower,
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: time.Duration(spellData.ImprovedShiftingPower.Effect(dbcenums.A_ADD_FLAT_MODIFIER, int32(dbcenums.SPELLMOD_COOLDOWN)).ValueAt(druid.Talents.ImprovedShiftingPower)) * time.Millisecond,
		})
	}
}
