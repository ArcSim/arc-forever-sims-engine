package core

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// mythicsim patch 165: every Merchant's Favor recipe's crafted item is in the database, with its
// client 1.60.1.70450 effect and the picker it goes in.
func TestMerchantsFavorConsumablesLoad(t *testing.T) {
	for _, tc := range []struct {
		id    int32
		name  string
		kind  proto.ConsumableType
		stats stats.Stats
	}{
		{250937, "Major Spellblasting Potion", proto.ConsumableType_ConsumableTypePotion, stats.Stats{stats.SpellDamage: 47}},
		{250943, "Major Frenzy Potion", proto.ConsumableType_ConsumableTypePotion, stats.Stats{stats.AttackPower: 80, stats.RangedAttackPower: 80}},
		{250949, "Major Mender's Potion", proto.ConsumableType_ConsumableTypePotion, stats.Stats{stats.HealingPower: 88}},
		{250950, "Potion of Venomous Blood", proto.ConsumableType_ConsumableTypePotion, stats.Stats{}},
		{250951, "Caustic Smog Potion", proto.ConsumableType_ConsumableTypePotion, stats.Stats{}},
		{250953, "Dragonfire Potion", proto.ConsumableType_ConsumableTypePotion, stats.Stats{}},
		{250954, "Potion of Elemental Purging", proto.ConsumableType_ConsumableTypePotion, stats.Stats{}},
		{250350, "Elixir of Ferocity", proto.ConsumableType_ConsumableTypeBattleElixir, stats.Stats{stats.Strength: 18, stats.Agility: 18}},
		{250328, "Elixir of Cunning", proto.ConsumableType_ConsumableTypeBattleElixir, stats.Stats{stats.Agility: 25, stats.Intellect: 25}},
		{250351, "Elixir of the Grizzly", proto.ConsumableType_ConsumableTypeStrengthBuff, stats.Stats{stats.Strength: 25, stats.MeleeCritRating: 2 * PhysicalCritRatingPerCritPercent}},
		{250333, "Greater Cleric's Elixir", proto.ConsumableType_ConsumableTypeSpellPowerElixir, stats.Stats{stats.HealingPower: 40}},
		{250335, "Elixir of Greater Fortitude", proto.ConsumableType_ConsumableTypeGuardianElixir, stats.Stats{stats.Health: 400}},
		{250336, "Elixir of Wicked Regeneration", proto.ConsumableType_ConsumableTypeGuardianElixir, stats.Stats{stats.MP5: 10}},
		{250337, "Elixir of the Owl", proto.ConsumableType_ConsumableTypeGuardianElixir, stats.Stats{stats.Intellect: 25, stats.MeleeCritRating: 2 * PhysicalCritRatingPerCritPercent}},
		{250338, "Elixir of Sages", proto.ConsumableType_ConsumableTypeGuardianElixir, stats.Stats{stats.Spirit: 25, stats.MeleeCritRating: 2 * PhysicalCritRatingPerCritPercent}},
		{250348, "Elixir of the Whale", proto.ConsumableType_ConsumableTypeGuardianElixir, stats.Stats{stats.Spirit: 25, stats.MeleeCritRating: 2 * PhysicalCritRatingPerCritPercent}},
		{250347, "Elixir of Greater Spirit", proto.ConsumableType_ConsumableTypeGuardianElixir, stats.Stats{stats.Spirit: 18}},
		{250329, "Elixir of the Phalanx", proto.ConsumableType_ConsumableTypeDefenseElixir, stats.Stats{stats.Health: 400, stats.Armor: 500}},
		{250343, "Elixir of Nature Power", proto.ConsumableType_ConsumableTypeSchoolElixir, stats.Stats{stats.NatureDamage: 40}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := GetConsumableByID(tc.id)
			if c.Id != tc.id {
				t.Fatalf("%d is not in the database", tc.id)
			}
			if c.Name != tc.name || c.Type != tc.kind {
				t.Fatalf("%d: got %q type %v, want %q type %v", tc.id, c.Name, c.Type, tc.name, tc.kind)
			}
			if c.Stats != tc.stats {
				t.Fatalf("%s: stats %v, want %v", tc.name, c.Stats.FlatString(), tc.stats.FlatString())
			}
		})
	}
}

// mythicsim patch 166: a consumable's crit counts for spells too.
func TestConsumableCritIsAllCrit(t *testing.T) {
	critDeltas := func(set func(*proto.ConsumesSpec)) (physical, spell float64) {
		_, without := setupConsumesSim(nil)
		_, with := setupConsumesSim(func(request *proto.RaidSimRequest) { set(consumesOf(request)) })
		return with.GetStat(stats.PhysicalCritPercent) - without.GetStat(stats.PhysicalCritPercent),
			with.GetStat(stats.SpellCritPercent) - without.GetStat(stats.SpellCritPercent)
	}
	for _, tc := range []struct {
		name string
		set  func(*proto.ConsumesSpec)
	}{
		{"Elixir of the Mongoose", func(c *proto.ConsumesSpec) { c.BattleElixirId = 13452 }},
		{"Elixir of the Owl", func(c *proto.ConsumesSpec) { c.GuardianElixirId = 250337 }},
		{"Elixir of Sages", func(c *proto.ConsumesSpec) { c.GuardianElixirId = 250338 }},
		{"Elixir of the Grizzly", func(c *proto.ConsumesSpec) { c.StrengthBuffId = 250351 }},
		{"Grilled Squid", func(c *proto.ConsumesSpec) { c.FoodId = 13928 }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			want := 2.0
			if tc.name == "Grilled Squid" {
				want = 1
			}
			physical, spell := critDeltas(tc.set)
			if !WithinToleranceFloat64(want, physical, 0.001) || !WithinToleranceFloat64(want, spell, 0.001) {
				t.Fatalf("%s: physical %+.3f%%, spell %+.3f%%, want %+.0f%% each", tc.name, physical, spell, want)
			}

			ConsumableCritIsAllCrit = false
			defer func() { ConsumableCritIsAllCrit = true }()
			physical, spell = critDeltas(tc.set)
			if !WithinToleranceFloat64(want, physical, 0.001) || spell != 0 {
				t.Fatalf("%s with the switch off: physical %+.3f%%, spell %+.3f%%, want %+.0f%% and 0", tc.name, physical, spell, want)
			}
		})
	}
}

func drinkPotion(t *testing.T, potion int32, tweak func(*proto.RaidSimRequest)) (*Simulation, *FakeRageWarrior) {
	t.Helper()
	sim, fw := setupConsumesSim(func(request *proto.RaidSimRequest) {
		consumesOf(request).PotId = potion
		consumesOf(request).Potions = []int32{potion}
		if tweak != nil {
			tweak(request)
		}
	})
	spell := fw.GetSpell(ActionID{ItemID: potion})
	if spell == nil {
		t.Fatalf("potion %d is not registered", potion)
	}
	if !spell.Cast(sim, fw.CurrentTarget) {
		t.Fatalf("potion %d could not be drunk", potion)
	}
	return sim, fw
}

// A level 60 target has no level-based resistance, so a hit lands whole.
func levelSixtyTarget(request *proto.RaidSimRequest) {
	request.Encounter.Targets[0].Level = 60
}

func stepUntil(sim *Simulation, at time.Duration) {
	for sim.CurrentTime < at {
		if sim.Step() { // true once the encounter is over
			return
		}
	}
}

// mythicsim patch 167: Dragonfire Potion's 280 Fire at once and four ticks of 35.
func TestDragonfirePotion(t *testing.T) {
	sim, fw := drinkPotion(t, DragonfirePotionItemID, levelSixtyTarget)
	stepUntil(sim, 9*time.Second)
	damage := fw.GetSpell(ActionID{SpellID: 1252554})
	metrics := &damage.SpellMetrics[fw.CurrentTarget.UnitIndex]
	if metrics.Hits+metrics.Crits != 1 {
		t.Fatalf("the direct hit should land once (seed 100 lands it), got %d hits and %d crits", metrics.Hits, metrics.Crits)
	}
	scale := 1.0
	if metrics.Crits == 1 {
		scale = damage.CritDamageMultiplier(fw.AttackTables[fw.CurrentTarget.UnitIndex])
	}
	direct := metrics.TotalDamage - metrics.TotalTickDamage
	if !WithinToleranceFloat64(280*scale, direct, 0.01) {
		t.Fatalf("direct damage %.2f, want %.2f", direct, 280*scale)
	}
	if metrics.Ticks != 4 || !WithinToleranceFloat64(140, metrics.TotalTickDamage, 0.01) {
		t.Fatalf("got %d ticks for %.2f, want 4 for 140", metrics.Ticks, metrics.TotalTickDamage)
	}
}

// mythicsim patch 167: Caustic Smog Potion's twelve ticks of 33 Nature (+-10%).
func TestCausticSmogPotion(t *testing.T) {
	sim, fw := drinkPotion(t, CausticSmogPotionItemID, levelSixtyTarget)
	stepUntil(sim, 13*time.Second)
	smog := fw.GetSpell(ActionID{SpellID: 1252543})
	metrics := &smog.SpellMetrics[fw.CurrentTarget.UnitIndex]
	if metrics.Ticks != 12 {
		t.Fatalf("got %d ticks, want 12", metrics.Ticks)
	}
	if metrics.TotalTickDamage < 12*29.7-0.01 || metrics.TotalTickDamage > 12*36.3+0.01 {
		t.Fatalf("tick damage %.2f outside 12 x 29.7 to 36.3", metrics.TotalTickDamage)
	}
}

// mythicsim patch 167: Potion of Elemental Purging's 105 spell damage against Elementals only, for 5 min.
func TestElementalPurgingPotion(t *testing.T) {
	elemental := func(request *proto.RaidSimRequest) {
		request.Encounter.Targets[0].MobType = proto.MobType_MobTypeElemental
	}
	sim, fw := drinkPotion(t, ElementalPurgingPotionItemID, elemental)
	table := fw.AttackTables[fw.CurrentTarget.UnitIndex]
	if got := table.MobTypeBonusStats[proto.MobType_MobTypeElemental][stats.SpellDamage]; got != 105 {
		t.Fatalf("spell damage against Elementals %.0f, want 105", got)
	}
	if got := table.MobTypeBonusStats[proto.MobType_MobTypeUndead][stats.SpellDamage]; got != 0 {
		t.Fatalf("spell damage against Undead %.0f, want 0", got)
	}
	if fw.GetStat(stats.SpellDamage) != 0 {
		t.Fatalf("the potion should add no spell damage against everything, got %.0f", fw.GetStat(stats.SpellDamage))
	}
	aura := fw.GetAura("Potion of Elemental Purging")
	if aura == nil || aura.Duration != 5*time.Minute || !aura.IsActive() {
		t.Fatal("the potion should show a 5 min buff")
	}
	aura.Deactivate(sim)
	if got := table.MobTypeBonusStats[proto.MobType_MobTypeElemental][stats.SpellDamage]; got != 0 {
		t.Fatalf("after the buff, spell damage against Elementals %.0f, want 0", got)
	}
}

// Potion of Venomous Blood only splashes when its drinker takes Bleed damage, which no
// encounter deals: drinking it takes the potion cooldown and does nothing else.
func TestVenomousBloodPotionDoesNothingWithoutBleeds(t *testing.T) {
	_, before := setupConsumesSim(nil)
	sim, fw := drinkPotion(t, VenomousBloodPotionItemID, nil)
	if fw.GetStats() != before.GetStats() {
		t.Fatalf("it should change no stat: %s against %s", fw.GetStats().FlatString(), before.GetStats().FlatString())
	}
	stepUntil(sim, 25*time.Second)
	if fw.GetPotionCD().IsReady(sim) {
		t.Fatal("drinking it should start the potion cooldown")
	}
}

// mythicsim patch 171: the elixirs' identical +2% crit rider applies once; each elixir's own stat
// still applies, and the Grilled Squid's well fed crit adds beside it.
func TestConsumableCritRiderAppliesOnce(t *testing.T) {
	deltas := func(set func(*proto.ConsumesSpec)) (physical, spell, strength, agility, spirit float64) {
		_, without := setupConsumesSim(nil)
		_, with := setupConsumesSim(func(request *proto.RaidSimRequest) { set(consumesOf(request)) })
		d := func(s stats.Stat) float64 { return with.GetStat(s) - without.GetStat(s) }
		return d(stats.PhysicalCritPercent), d(stats.SpellCritPercent), d(stats.Strength), d(stats.Agility), d(stats.Spirit)
	}
	all := func(c *proto.ConsumesSpec) {
		c.BattleElixirId = 13452    // Mongoose
		c.GuardianElixirId = 250338 // Sages
		c.StrengthBuffId = 250351   // Grizzly
	}
	physical, spell, strength, agility, spirit := deltas(all)
	if !WithinToleranceFloat64(2, physical, 0.001) || !WithinToleranceFloat64(2, spell, 0.001) {
		t.Fatalf("three rider elixirs: physical %+.3f%%, spell %+.3f%%, want +2%% each", physical, spell)
	}
	if strength <= 0 || agility <= 0 || spirit <= 0 {
		t.Fatalf("every elixir's own stat should apply: strength %+.1f, agility %+.1f, spirit %+.1f", strength, agility, spirit)
	}
	_, _, aloneStr, _, _ := deltas(func(c *proto.ConsumesSpec) { c.StrengthBuffId = 250351 })
	if !WithinToleranceFloat64(aloneStr, strength, 0.001) {
		t.Fatalf("the Grizzly's Strength with the others %+.2f, alone %+.2f", strength, aloneStr)
	}

	physical, spell, _, _, _ = deltas(func(c *proto.ConsumesSpec) { all(c); c.FoodId = 13928 })
	if !WithinToleranceFloat64(3, physical, 0.001) || !WithinToleranceFloat64(3, spell, 0.001) {
		t.Fatalf("riders and Grilled Squid: physical %+.3f%%, spell %+.3f%%, want +3%% each", physical, spell)
	}

	ConsumableCritRiderStacks = true
	defer func() { ConsumableCritRiderStacks = false }()
	physical, spell, _, _, _ = deltas(all)
	if !WithinToleranceFloat64(6, physical, 0.001) || !WithinToleranceFloat64(6, spell, 0.001) {
		t.Fatalf("with the switch on: physical %+.3f%%, spell %+.3f%%, want +6%% each", physical, spell)
	}
}
