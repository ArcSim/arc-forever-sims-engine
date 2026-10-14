package core

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
)

func weaponSkillTable(defenderLevel int32, skill [weaponSkillHands]float64) (*AttackTable, *Unit) {
	attacker := &Unit{Type: PlayerUnit, Level: CharacterLevel, weaponSkill: skill}
	defender := &Unit{Type: EnemyUnit, Level: defenderLevel}
	return NewAttackTable(attacker, defender), attacker
}

// With no skill above base every lookup is the table's own field, bit for bit, so sims without
// weapon skill do not move.
func TestWeaponSkillAtBaseIsTheTable(t *testing.T) {
	for level := CharacterLevel - 2; level <= CharacterLevel+3; level++ {
		table, attacker := weaponSkillTable(int32(level), [weaponSkillHands]float64{})
		for _, mask := range []ProcMask{ProcMaskMeleeMHAuto, ProcMaskMeleeOHSpecial, ProcMaskRangedAuto, ProcMaskSpellDamage} {
			spell := &Spell{Unit: attacker, ProcMask: mask}
			multiplier, spread := table.Glance(spell)
			if table.MissChance(spell) != table.BaseMissChance || table.DodgeChance(spell) != table.BaseDodgeChance ||
				table.HitSuppressionFor(spell) != table.HitSuppression || multiplier != table.GlanceMultiplier || spread != table.GlanceSpread {
				t.Errorf("level %d, %v: skill lookups differ from the table at base skill", level, mask)
			}
		}
	}
}

// The Classic formulas reproduce the boss table at base skill (a 15 point gap), so the table and the
// formulas agree where they meet.
func TestWeaponSkillFormulasMeetTheBossTable(t *testing.T) {
	table, _ := weaponSkillTable(CharacterLevel+3, [weaponSkillHands]float64{})
	low, high := max(min(1.3-0.05*15, 0.91), 0.01), max(min(1.2-0.03*15, 0.99), 0.2)
	for name, pair := range map[string][2]float64{
		"miss":            {weaponSkillAvoidance(15, 0.002), table.BaseMissChance},
		"dodge":           {weaponSkillAvoidance(15, 0.001), table.BaseDodgeChance},
		"hit suppression": {(15 - 10) * 0.002, table.HitSuppression},
		"glance":          {(low + high) / 2, table.GlanceMultiplier},
		"glance spread":   {(high - low) / 2, table.GlanceSpread},
	} {
		if math.Abs(pair[0]-pair[1]) > 1e-12 {
			t.Errorf("%s: formula %v, table %v", name, pair[0], pair[1])
		}
	}
}

// Against a level 63 boss (315 defense), per hand.
func TestWeaponSkillAgainstABoss(t *testing.T) {
	near := func(a, b float64) bool { return math.Abs(a-b) < 1e-9 }
	// Main hand 311 (Aged Core Leather Gloves and Distracting Dagger with a dagger), off hand 305,
	// ranged 300.
	table, attacker := weaponSkillTable(CharacterLevel+3, [weaponSkillHands]float64{11, 5, 0})
	for _, tc := range []struct {
		name                       string
		mask                       ProcMask
		miss, dodge, hitSuppressed float64
		glance, spread             float64
	}{
		// 4 point gap: 5.4% miss and dodge, no hit suppression, glances keep 91-99%.
		{"main hand", ProcMaskMeleeMHAuto, 0.054, 0.054, 0, 0.95, 0.04},
		{"main hand ability", ProcMaskMeleeMHSpecial, 0.054, 0.054, 0, 0.95, 0.04},
		// 10 point gap: 6%, still no suppression, glances keep 80-90%.
		{"off hand", ProcMaskMeleeOHAuto, 0.06, 0.06, 0, 0.85, 0.05},
		// Base skill: the table.
		{"ranged", ProcMaskRangedAuto, 0.08, 0.065, 0.01, 0.65, 0.10},
		{"no weapon", ProcMaskSpellDamage, 0.08, 0.065, 0.01, 0.65, 0.10},
	} {
		spell := &Spell{Unit: attacker, ProcMask: tc.mask}
		multiplier, spread := table.Glance(spell)
		if !near(table.MissChance(spell), tc.miss) || !near(table.DodgeChance(spell), tc.dodge) ||
			!near(table.HitSuppressionFor(spell), tc.hitSuppressed) || !near(multiplier, tc.glance) || !near(spread, tc.spread) {
			t.Errorf("%s: miss %v dodge %v suppression %v glance %v +- %v, want %v %v %v %v +- %v", tc.name,
				table.MissChance(spell), table.DodgeChance(spell), table.HitSuppressionFor(spell), multiplier, spread,
				tc.miss, tc.dodge, tc.hitSuppressed, tc.glance, tc.spread)
		}
		// Glancing chance and crit suppression stay on the base skill.
		if table.BaseGlanceChance != 0.40 || table.MeleeCritSuppression != 0.048 {
			t.Errorf("%s: glance chance %v, crit suppression %v", tc.name, table.BaseGlanceChance, table.MeleeCritSuppression)
		}
	}

	// A 12 point gap (303) is past 10: 0.2% a point and 0.4% of hit ignored.
	table, attacker = weaponSkillTable(CharacterLevel+3, [weaponSkillHands]float64{3})
	spell := &Spell{Unit: attacker, ProcMask: ProcMaskMeleeMHAuto}
	if !near(table.MissChance(spell), 0.074) || !near(table.HitSuppressionFor(spell), 0.004) || !near(table.DodgeChance(spell), 0.062) {
		t.Errorf("303 skill: miss %v suppression %v dodge %v", table.MissChance(spell), table.HitSuppressionFor(spell), table.DodgeChance(spell))
	}

	// Skill past the target's defense (320 against 315) takes 0.04% a point off.
	table, attacker = weaponSkillTable(CharacterLevel+3, [weaponSkillHands]float64{20})
	spell = &Spell{Unit: attacker, ProcMask: ProcMaskMeleeMHAuto}
	if !near(table.MissChance(spell), 0.048) || !near(table.DodgeChance(spell), 0.048) {
		t.Errorf("320 skill: miss %v dodge %v", table.MissChance(spell), table.DodgeChance(spell))
	}

	// Encounter.NoHitSuppression stays off whatever the skill.
	table, attacker = weaponSkillTable(CharacterLevel+3, [weaponSkillHands]float64{3})
	table.HitSuppression = 0
	if got := table.HitSuppressionFor(&Spell{Unit: attacker, ProcMask: ProcMaskMeleeMHAuto}); got != 0 {
		t.Errorf("NoHitSuppression: %v", got)
	}
}

// An enemy's table against a player, and a target below the attacker, ignore the attacker's skill.
func TestWeaponSkillOnlyAgainstEnemiesAtOrAboveTheAttacker(t *testing.T) {
	table, attacker := weaponSkillTable(CharacterLevel-2, [weaponSkillHands]float64{11})
	spell := &Spell{Unit: attacker, ProcMask: ProcMaskMeleeMHAuto}
	if table.MissChance(spell) != table.BaseMissChance {
		t.Errorf("a lower level target moved: %v", table.MissChance(spell))
	}

	player := &Unit{Type: PlayerUnit, Level: CharacterLevel, weaponSkill: [weaponSkillHands]float64{11}}
	enemy := &Unit{Type: EnemyUnit, Level: CharacterLevel + 3}
	enemyTable := NewAttackTable(enemy, player)
	if got := enemyTable.DodgeChance(&Spell{Unit: enemy, ProcMask: ProcMaskMeleeMHAuto}); got != enemyTable.BaseDodgeChance {
		t.Errorf("the enemy's table moved: %v", got)
	}
}

func TestWeaponSkillOf(t *testing.T) {
	for _, tc := range []struct {
		item Item
		want proto.PseudoStat
		ok   bool
	}{
		{Item{ID: 1, WeaponType: proto.WeaponType_WeaponTypeDagger, HandType: proto.HandType_HandTypeOneHand}, proto.PseudoStat_PseudoStatDaggersSkill, true},
		{Item{ID: 1, WeaponType: proto.WeaponType_WeaponTypeSword, HandType: proto.HandType_HandTypeMainHand}, proto.PseudoStat_PseudoStatSwordsSkill, true},
		{Item{ID: 1, WeaponType: proto.WeaponType_WeaponTypeSword, HandType: proto.HandType_HandTypeTwoHand}, proto.PseudoStat_PseudoStatTwoHandedSwordsSkill, true},
		{Item{ID: 1, WeaponType: proto.WeaponType_WeaponTypeAxe, HandType: proto.HandType_HandTypeOffHand}, proto.PseudoStat_PseudoStatAxesSkill, true},
		{Item{ID: 1, WeaponType: proto.WeaponType_WeaponTypeFist, HandType: proto.HandType_HandTypeOneHand}, proto.PseudoStat_PseudoStatUnarmedSkill, true},
		{Item{ID: 1, WeaponType: proto.WeaponType_WeaponTypePolearm, HandType: proto.HandType_HandTypeTwoHand}, proto.PseudoStat_PseudoStatPolearmsSkill, true},
		{Item{ID: 1, WeaponType: proto.WeaponType_WeaponTypeStaff, HandType: proto.HandType_HandTypeTwoHand}, proto.PseudoStat_PseudoStatStavesSkill, true},
		{Item{ID: 1, RangedWeaponType: proto.RangedWeaponType_RangedWeaponTypeGun}, proto.PseudoStat_PseudoStatGunsSkill, true},
		{Item{ID: 1, RangedWeaponType: proto.RangedWeaponType_RangedWeaponTypeWand}, 0, false},
		{Item{ID: 1, WeaponType: proto.WeaponType_WeaponTypeShield, HandType: proto.HandType_HandTypeOffHand}, 0, false},
		{Item{ID: 1, WeaponType: proto.WeaponType_WeaponTypeOffHand, HandType: proto.HandType_HandTypeOffHand}, 0, false},
		{Item{}, 0, false},
	} {
		if got, ok := WeaponSkillOf(&tc.item); got != tc.want || ok != tc.ok {
			t.Errorf("%+v: %v %v, want %v %v", tc.item, got, ok, tc.want, tc.ok)
		}
	}
}
