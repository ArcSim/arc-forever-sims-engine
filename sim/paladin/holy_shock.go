package paladin

import (
	"time"

	"github.com/wowsims/forever/sim/core"
)

// One rank of Holy Shock: the spell the paladin casts, and the damage and heal it triggers. The
// client ships the damage and heal chains as two ladders under one name, which the generator
// refuses, so the table is by hand from the client rows: the cast's cost, and each chain's number
// with its 0.429 coefficient. The damage is the client's min-max roll at the rank's max level.
type holyShockRank struct {
	rank     int32
	spellID  int32
	cost     int32
	damageID int32
	damage   [2]float64
	healID   int32
	heal     float64
}

var HolyShockRanks = []holyShockRank{
	{rank: 1, spellID: 1311606, cost: 160, damageID: 1311604, damage: [2]float64{128, 140}, healID: 1311605, heal: 114},
	{rank: 2, spellID: 20473, cost: 225, damageID: 25912, damage: [2]float64{175, 189}, healID: 25914, heal: 156},
	{rank: 3, spellID: 20929, cost: 275, damageID: 25911, damage: [2]float64{248, 268}, healID: 25913, heal: 230},
	{rank: 4, spellID: 20930, cost: 325, damageID: 25902, damage: [2]float64{334, 362}, healID: 25903, heal: 320},
}

const (
	holyShockCoefficient = 0.429
	holyShockCooldown    = time.Second * 10
	holyShockRange       = 20
)

// Holy Shock (talent)
// https://www.wowhead.com/forever/spell=1311606
//
// Blasts the target with Holy energy, causing X Holy damage to an enemy, or Y healing to an ally.
//
// The damage and the heal are two spells that share the cast's cooldown: the damage carries the
// cast's id so the rotation names it as the client does, the heal carries its own. Light's Vigil
// on the enemy turns the damage into the vigil's strike, refunds the vigil, and skips the cooldown.
func (paladin *Paladin) registerHolyShock() {
	for _, rank := range HolyShockRanks {
		paladin.registerHolyShockRank(rank)
	}
}

func (paladin *Paladin) registerHolyShockRank(rank holyShockRank) {
	timer := paladin.sharedTimer(&paladin.holyShockTimer)

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.spellID},
		SpellSchool:    core.SpellSchoolHoly,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskHolyShock,
		Rank:           rank.rank,
		MaxRange:       holyShockRange,

		ManaCost: core.ManaCostOptions{FlatCost: rank.cost},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    timer,
				Duration: holyShockCooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: holyShockCoefficient,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if vigil := paladin.activeLightsVigil(target); vigil != nil {
				vigil.consume(sim, target)
				spell.CD.Reset()
				return
			}
			spell.CalcAndDealDamage(sim, target, sim.Roll(rank.damage[0], rank.damage[1]), spell.OutcomeMagicHitAndCrit)
		},
	})

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.healID},
		SpellSchool:    core.SpellSchoolHoly,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellHealing,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: SpellMaskHolyShockHeal,
		Rank:           rank.rank,
		MaxRange:       holyShockRange,

		ManaCost: core.ManaCostOptions{FlatCost: rank.cost},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    timer,
				Duration: holyShockCooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: holyShockCoefficient,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealHealing(sim, target, rank.heal, spell.OutcomeHealingCrit)
		},
	})
}
