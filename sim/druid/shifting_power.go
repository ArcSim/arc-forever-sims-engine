package druid

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
)

// Shifting Power (1322605), a Feral talent of client 1.60.1.70170 that replaces Tiger's Fury: "Instantly
// convert $c Mana into $s2 Energy. Shifting Power's cost is reduced by effects that reduce the cost of
// Shapeshifting." The row states 40 Energy (effect 1, energize), 55% of base mana (SpellPower, the same
// share Cat Form costs), a 16 second cooldown and a 1 second global cooldown. It needs Cat Form (caster
// aura 768 and shapeshift mask 1) and is not a form change: the druid stays in Cat Form.
//
// Natural Shapeshifter's cost reduction reaches it: the talent's class mask names the family bit Cat, Bear,
// Moonkin and Shifting Power share (word 0, 0x20000000). Since the 8 Oct beta build it can't be used at
// full Energy (dev notes, upstream #805).
var shiftingPowerRank = spellData.ShiftingPower.Highest()

func (druid *Druid) registerShiftingPowerSpell() {
	if !druid.Talents.ShiftingPower {
		return
	}

	actionID := core.ActionID{SpellID: shiftingPowerRank.ID}
	energy := shiftingPowerRank.EnergizeEffect().Average(core.CharacterLevel)
	energyMetrics := druid.NewEnergyMetrics(actionID)

	druid.ShiftingPower = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: DruidSpellShiftingPower,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		ManaCost: shiftingPowerRank.ManaCost(),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: shiftingPowerRank.GCD(),
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: max(shiftingPowerRank.Cooldown(), shiftingPowerRank.CategoryCooldown()),
			},
		},

		ExtraCastCondition: func(_ *core.Simulation, _ *core.Unit) bool {
			return druid.CurrentEnergy() < druid.MaximumEnergy()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			// Wolfshead Helm (8345): "an additional 20 Energy from activating Shifting Power".
			druid.AddEnergy(sim, energy+druid.WolfsheadShiftingPowerEnergy, energyMetrics)
		},
	})
}

// Improved Shifting Power (1322670): "Reduces the cooldown of your Shifting Power spell by ${$m1/-1000}
// sec", 4 and 8 seconds at the two ranks (the talent tree's curve states -4000 and -8000 milliseconds).
func (druid *Druid) applyImprovedShiftingPower() {
	if druid.Talents.ImprovedShiftingPower == 0 {
		return
	}

	druid.AddStaticMod(core.SpellModConfig{
		ClassMask: DruidSpellShiftingPower,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: time.Millisecond * time.Duration(spellData.ImprovedShiftingPower.Effect(dbcenums.A_ADD_FLAT_MODIFIER, int32(dbcenums.SPELLMOD_COOLDOWN)).ValueAt(druid.Talents.ImprovedShiftingPower)),
	})
}
