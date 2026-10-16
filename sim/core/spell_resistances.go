package core

import (
	"github.com/wowsims/forever/sim/core/stats"
)

func (result *SpellResult) applyResistances(sim *Simulation, spell *Spell, isPeriodic bool, attackTable *AttackTable) {
	resistanceMultiplier, outcome := spell.ResistanceMultiplier(sim, isPeriodic, attackTable)

	result.Damage *= resistanceMultiplier
	result.Outcome |= outcome

	result.ArmorAndResistanceMultiplier = resistanceMultiplier
	result.PostArmorAndResistanceMultiplier = result.Damage
}

// Modifies damage based on Armor or Magic resistances, depending on the damage type.
func (spell *Spell) ResistanceMultiplier(sim *Simulation, isPeriodic bool, attackTable *AttackTable) (float64, HitOutcome) {
	if spell.Flags.Matches(SpellFlagIgnoreResists) {
		return 1, OutcomeEmpty
	}

	if spell.SpellSchool.Matches(SpellSchoolPhysical) {
		// All physical dots (Bleeds) ignore armor.
		if isPeriodic && !spell.Flags.Matches(SpellFlagApplyArmorReduction) {
			return 1, OutcomeEmpty
		}

		// Physical resistance (armor).
		return attackTable.GetArmorDamageModifier(spell), OutcomeEmpty
	}

	// Spell Vulnerability: the attacker's penetration past the target's resistance adds damage, on
	// binary spells too, on top of the partial resist roll below.
	vulnerability := attackTable.Defender.spellVulnerabilityMultiplier(spell, attackTable.Attacker)

	// Magical resistance.
	if spell.Flags.Matches(SpellFlagBinary) {
		return vulnerability, OutcomeEmpty
	}

	resistanceRoll := sim.RandomFloat("Partial Resist")

	threshold00, threshold25, threshold50 := attackTable.GetPartialResistThresholds(spell)
	//if sim.Log != nil {
	//	sim.Log("Resist thresholds: %0.04f, %0.04f, %0.04f", threshold00, threshold25, threshold50)
	//}

	if resistanceRoll > threshold00 {
		// No partial resist.
		return vulnerability, OutcomeEmpty
	} else if resistanceRoll > threshold25 {
		return 0.75 * vulnerability, OutcomePartial1_4
	} else if resistanceRoll > threshold50 {
		return 0.5 * vulnerability, OutcomePartial2_4
	} else {
		return 0.25 * vulnerability, OutcomePartial3_4
	}
}

// SpellVulnerabilityBelowZeroResistance (MythicSim patch 146): Spell Piercing past the target's
// resistance raises the spell's damage by 1 / (5 x target level) per point. Client tooltip (Spell
// Piercing 7, beta 70425): "Spell Piercing may bring a creature below 0 Resistances, becoming Spell
// Vulnerability, increasing their spell damage taken". Targetdummy's combat log (15 October 2026): a
// Frostbolt from a mage with 7 against a level 28 target at 0 resistance hit for 460 with resisted -21,
// 5.0% of 439 (7 / 140), truncated. Only the attacker's penetration counts, measured from the target's
// resistance floored at 0: a debuff that pushes the resistance itself below 0 (Curse of the Elements)
// adds nothing, as there is no evidence for it. False is the old clamp at 0.
var SpellVulnerabilityBelowZeroResistance = true

// The target's resistance to the spell's school, before the attacker's penetration.
func (unit *Unit) rawMagicResistance(spell *Spell) float64 {
	// Frostfire Bolt explicitly checks the lower resistance, even when both are nonzero.
	if spell.SpellSchool == SpellSchoolFrostfire {
		return min(unit.GetStat(stats.FireResistance), unit.GetStat(stats.FrostResistance))
	}
	return unit.GetStat(spell.SpellSchool.ResistanceStat())
}

// The damage multiplier of Spell Vulnerability, 1 without it. Schools with no resistance (Holy, and
// combinations other than Frostfire) take none.
func (unit *Unit) spellVulnerabilityMultiplier(spell *Spell, attacker *Unit) float64 {
	if !SpellVulnerabilityBelowZeroResistance || spell.SchoolIndex <= stats.SchoolIndexPhysical || unit.Level <= 0 {
		return 1
	}
	if spell.SpellSchool != SpellSchoolFrostfire && spell.SpellSchool.ResistanceStat() == 0 {
		return 1
	}
	beyond := attacker.stats[stats.SpellPiercing] - max(0, unit.rawMagicResistance(spell))
	if beyond <= 0 {
		return 1
	}
	return 1 + beyond/(5*float64(unit.Level))
}

func (at *AttackTable) GetPartialResistThresholds(spell *Spell) (float64, float64, float64) {
	return at.Defender.partialResistRollThresholds(spell, at.Attacker)
}

func (at *AttackTable) GetBinaryHitChance(spell *Spell) float64 {
	return at.Defender.binaryHitChance(spell, at.Attacker)
}

// All of the following calculations are based on this guide:
// https://royalgiraffe.github.io/resist-guide

func (unit *Unit) resistCoeff(spell *Spell, attacker *Unit, binary bool) float64 {
	if spell.SchoolIndex <= stats.SchoolIndexPhysical {
		return 0
	}

	resistance := max(0, unit.rawMagicResistance(spell)-attacker.stats[stats.SpellPiercing])
	if resistance <= 0 {
		// Level-based resistance does not apply to binary spells.
		if binary {
			return 0
		}
		// Scaled like the branch below (coef runs 0 to 1, mitigation 0 to 0.75), so a level 63
		// boss averages 6% as the guide says, not 4.5%.
		return min(1, unit.levelBasedResist(attacker)/0.75)
	}

	resistanceCap := float64(attacker.Level * 5)
	resistanceCoef := resistance / resistanceCap

	if !binary && unit.Type == EnemyUnit && unit.Level > attacker.Level {
		avgMitigationAdded := unit.levelBasedResist(attacker)
		// coef is scaled 0 to 1, not 0 to 0.75
		resistanceCoef += avgMitigationAdded * 1 / 0.75
	}

	return min(1, resistanceCoef)
}

func (unit *Unit) levelBasedResist(attacker *Unit) float64 {
	if unit.Type == EnemyUnit && unit.Level > attacker.Level {
		// 2% average mitigation per level difference
		return 0.02 * float64(unit.Level-attacker.Level)
	}
	return 0
}

func (unit *Unit) binaryHitChance(spell *Spell, attacker *Unit) float64 {
	resistCoeff := unit.resistCoeff(spell, attacker, true)
	return 1 - 0.75*resistCoeff
}

// Roll threshold for each type of partial resist.
func (unit *Unit) partialResistRollThresholds(spell *Spell, attacker *Unit) (float64, float64, float64) {
	resistCoeff := unit.resistCoeff(spell, attacker, false)

	// Based on the piecewise linear regression estimates at https://royalgiraffe.github.io/partial-resist-table.
	if val := resistCoeff * 3; val <= 1 {
		return 0.76 * val, 0.21 * val, 0.03 * val
	} else if val <= 2 {
		val -= 1
		return 0.76 + 0.24*val, 0.21 + 0.57*val, 0.03 + 0.19*val
	} else {
		val -= 2
		return 1, 0.78 + 0.18*val, 0.22 + 0.58*val
	}
}

// https://web.archive.org/web/20130208043756/http://elitistjerks.com/f15/t29453-combat_ratings_level_85_cataclysm/
// https://web.archive.org/web/20110309163709/http://elitistjerks.com/f78/t105429-cataclysm_mechanics_testing/
func (at *AttackTable) GetArmorDamageModifier(spell *Spell) float64 {
	if at.IgnoreArmor {
		return 1.0
	}

	ignoreArmorFactor := Clamp(at.ArmorIgnoreFactor, 0.0, 1.0)

	// Level 60 armor: 400 + 85 x attacker level (5500 for a level 60 attacker). The
	// 467.5 x level - 22167.5 this engine had is TBC's formula for attackers of level 60+.
	armorConstant := 400 + 85*float64(at.Attacker.Level)
	defenderArmor := at.Defender.Armor() - (at.Defender.Armor() * ignoreArmorFactor)
	// TODO: Forever possibly has different rules as there is Armor Piercing Flat and % values
	defenderArmor = max(defenderArmor-at.Attacker.stats[stats.ArmorPenetration], 0)
	// Damage reduction from armor is capped at 75%.
	return max(1-defenderArmor/(defenderArmor+armorConstant), 0.25)
}
