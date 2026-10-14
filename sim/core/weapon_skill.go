package core

import (
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// Weapon skill, ported from wowsims/classic (sim/core/target.go NewAttackTable, which follows
// https://github.com/magey/classic-warrior/wiki/Attack-table).
//
// Forever keeps Classic's weapon skill: items carry it as item stats 90-103 (Death's Sting "Increased
// Daggers +3"), set bonuses and suffixes as A_MOD_SKILL spells (Defias Leather 5 piece). Only the
// racials (now +1% crit) and Weapon Expertise (now expertise) stopped granting it.
//
// Every character has 5 skill a level in every weapon; what this tracks is the skill above that,
// per skill line, kept as the item PseudoStats proto.PseudoStat_PseudoStat*Skill. An attack uses the
// skill of the weapon it is made with: the main hand's for ProcMaskMeleeMH, the off hand's for
// ProcMaskMeleeOH and the ranged weapon's for ProcMaskRanged. Attacks with no weapon, pets and
// druids (whose forms fight with paws) attack at the base skill.
//
// Against a target at or above the attacker's level, defense minus weapon skill sets the miss and
// dodge chances, the hit suppression past a 10 point gap and how much damage a glancing blow
// keeps. Glancing chance, parry and crit suppression stay on the base skill, as in wowsims/classic.
// With no skill above base every value is the attack table's own, so a sim without weapon skill is
// unchanged.

// The three weapons an attack can be made with, in Unit.weaponSkill.
const (
	weaponSkillMainHand = iota
	weaponSkillOffHand
	weaponSkillRanged
	weaponSkillHands
)

// Every skill line, in proto order from PseudoStatAxesSkill.
const numWeaponSkills = 14

var WeaponSkillPseudoStats = [numWeaponSkills]proto.PseudoStat{
	proto.PseudoStat_PseudoStatAxesSkill,
	proto.PseudoStat_PseudoStatSwordsSkill,
	proto.PseudoStat_PseudoStatMacesSkill,
	proto.PseudoStat_PseudoStatDaggersSkill,
	proto.PseudoStat_PseudoStatUnarmedSkill,
	proto.PseudoStat_PseudoStatTwoHandedAxesSkill,
	proto.PseudoStat_PseudoStatTwoHandedSwordsSkill,
	proto.PseudoStat_PseudoStatTwoHandedMacesSkill,
	proto.PseudoStat_PseudoStatPolearmsSkill,
	proto.PseudoStat_PseudoStatStavesSkill,
	proto.PseudoStat_PseudoStatBowsSkill,
	proto.PseudoStat_PseudoStatGunsSkill,
	proto.PseudoStat_PseudoStatCrossbowsSkill,
	proto.PseudoStat_PseudoStatThrownSkill,
}

// The skill line a weapon is used with, as its PseudoStat; false for no weapon, a shield, an
// off-hand item, a wand or a relic.
func WeaponSkillOf(item *Item) (proto.PseudoStat, bool) {
	if item == nil || item.ID == 0 {
		return 0, false
	}
	switch item.RangedWeaponType {
	case proto.RangedWeaponType_RangedWeaponTypeUnknown:
	case proto.RangedWeaponType_RangedWeaponTypeBow:
		return proto.PseudoStat_PseudoStatBowsSkill, true
	case proto.RangedWeaponType_RangedWeaponTypeGun:
		return proto.PseudoStat_PseudoStatGunsSkill, true
	case proto.RangedWeaponType_RangedWeaponTypeCrossbow:
		return proto.PseudoStat_PseudoStatCrossbowsSkill, true
	case proto.RangedWeaponType_RangedWeaponTypeThrown:
		return proto.PseudoStat_PseudoStatThrownSkill, true
	default:
		return 0, false
	}

	if item.HandType == proto.HandType_HandTypeTwoHand {
		switch item.WeaponType {
		case proto.WeaponType_WeaponTypeAxe:
			return proto.PseudoStat_PseudoStatTwoHandedAxesSkill, true
		case proto.WeaponType_WeaponTypeSword:
			return proto.PseudoStat_PseudoStatTwoHandedSwordsSkill, true
		case proto.WeaponType_WeaponTypeMace:
			return proto.PseudoStat_PseudoStatTwoHandedMacesSkill, true
		case proto.WeaponType_WeaponTypePolearm:
			return proto.PseudoStat_PseudoStatPolearmsSkill, true
		case proto.WeaponType_WeaponTypeStaff:
			return proto.PseudoStat_PseudoStatStavesSkill, true
		}
		return 0, false
	}

	switch item.WeaponType {
	case proto.WeaponType_WeaponTypeAxe:
		return proto.PseudoStat_PseudoStatAxesSkill, true
	case proto.WeaponType_WeaponTypeSword:
		return proto.PseudoStat_PseudoStatSwordsSkill, true
	case proto.WeaponType_WeaponTypeMace:
		return proto.PseudoStat_PseudoStatMacesSkill, true
	case proto.WeaponType_WeaponTypeDagger:
		return proto.PseudoStat_PseudoStatDaggersSkill, true
	case proto.WeaponType_WeaponTypeFist:
		return proto.PseudoStat_PseudoStatUnarmedSkill, true
	}
	return 0, false
}

// The character's skill above base in one skill line: its equipment's plus the bonuses of active
// auras (AddWeaponSkill).
func (character *Character) WeaponSkill(skill proto.PseudoStat) float64 {
	total := character.weaponSkillBonus[skill-proto.PseudoStat_PseudoStatAxesSkill]
	for i := range character.Equipment {
		total += stats.PseudoStatValue(character.Equipment[i].PseudoStats, skill)
	}
	return total
}

// Recomputes the skill each hand attacks with, after the equipment or a skill bonus changed.
func (character *Character) updateWeaponSkills() {
	if character.Class == proto.Class_ClassDruid {
		return
	}
	for hand, item := range [weaponSkillHands]*Item{character.MainHand(), character.OffHand(), character.Ranged()} {
		character.weaponSkill[hand] = 0
		if skill, ok := WeaponSkillOf(item); ok {
			character.weaponSkill[hand] = character.WeaponSkill(skill)
		}
	}
}

// Raises (or, negative, lowers) one skill line, for set bonuses and other auras.
func (unit *Unit) AddWeaponSkill(skill proto.PseudoStat, amount float64) {
	unit.weaponSkillBonus[skill-proto.PseudoStat_PseudoStatAxesSkill] += amount
	if unit.onWeaponSkillChange != nil {
		unit.onWeaponSkillChange()
	}
}

// Attaches a weapon skill bonus that lasts as long as the aura is active.
// Returns parent aura for chaining
func (parentAura *Aura) AttachWeaponSkill(skill proto.PseudoStat, amount float64) *Aura {
	parentAura.ApplyOnGain(func(aura *Aura, _ *Simulation) {
		aura.Unit.AddWeaponSkill(skill, amount)
	})

	parentAura.ApplyOnExpire(func(aura *Aura, _ *Simulation) {
		aura.Unit.AddWeaponSkill(skill, -amount)
	})

	if parentAura.IsActive() {
		parentAura.Unit.AddWeaponSkill(skill, amount)
	}

	return parentAura
}

// The skill above base the attack is made with.
func (at *AttackTable) weaponSkillBonus(spell *Spell) float64 {
	if at.Attacker != spell.Unit {
		return 0
	}
	switch {
	case spell.ProcMask.Matches(ProcMaskMeleeMH):
		return spell.Unit.weaponSkill[weaponSkillMainHand]
	case spell.ProcMask.Matches(ProcMaskMeleeOH):
		return spell.Unit.weaponSkill[weaponSkillOffHand]
	case spell.ProcMask.Matches(ProcMaskRanged):
		return spell.Unit.weaponSkill[weaponSkillRanged]
	}
	return 0
}

// Target defense minus the attack's weapon skill, and whether weapon skill moves the table at all:
// only against an enemy at or above the attacker's level, and only with skill above base.
func (at *AttackTable) weaponSkillGap(spell *Spell) (float64, bool) {
	bonus := at.weaponSkillBonus(spell)
	if bonus == 0 || at.Defender.Type != EnemyUnit || at.Defender.Level < at.Attacker.Level {
		return 0, false
	}
	return float64(at.Defender.Level-at.Attacker.Level)*5 - bonus, true
}

// Below the defense the gap costs 0.1% a point, 0.2% once it passes 10; above it, skill takes
// 0.04% a point off.
func weaponSkillAvoidance(gap float64, perPointFar float64) float64 {
	switch {
	case gap > 10:
		return 0.05 + gap*perPointFar
	case gap >= 0:
		return 0.05 + gap*0.001
	default:
		return 0.05 + gap*0.0004
	}
}

// The chance to miss before hit, for an attack with the given spell.
func (at *AttackTable) MissChance(spell *Spell) float64 {
	gap, ok := at.weaponSkillGap(spell)
	if !ok {
		return at.BaseMissChance
	}
	return weaponSkillAvoidance(gap, 0.002)
}

// How much of the attacker's hit the target ignores: 0.2% a point past a 10 point gap, which is
// the 1% a boss takes off an attack at base skill.
func (at *AttackTable) HitSuppressionFor(spell *Spell) float64 {
	gap, ok := at.weaponSkillGap(spell)
	if !ok || at.HitSuppression == 0 {
		return at.HitSuppression
	}
	return max(gap-10, 0) * 0.002
}

// The target's chance to dodge, before expertise.
func (at *AttackTable) DodgeChance(spell *Spell) float64 {
	gap, ok := at.weaponSkillGap(spell)
	if !ok {
		return at.BaseDodgeChance
	}
	return weaponSkillAvoidance(gap, 0.001)
}

// The average share of damage a glancing blow keeps, and the spread a roll can take it either side.
func (at *AttackTable) Glance(spell *Spell) (multiplier float64, spread float64) {
	gap, ok := at.weaponSkillGap(spell)
	if !ok {
		return at.GlanceMultiplier, at.GlanceSpread
	}
	low := max(min(1.3-0.05*gap, 0.91), 0.01)
	high := max(min(1.2-0.03*gap, 0.99), 0.2)
	return (low + high) / 2, (high - low) / 2
}
