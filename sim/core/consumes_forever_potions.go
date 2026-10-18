package core

import (
	"time"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// Forever's Merchant's Favor potions whose effect is not a stat buff (mythicsim patch 173), all
// from client 1.60.1.70450. Each is a combat potion (SpellCategory 4, the 2 min potion cooldown);
// the stat potions (Major Spellblasting, Frenzy and Mender's) take the generic path.
const (
	DragonfirePotionItemID        = 250953
	CausticSmogPotionItemID       = 250951
	ElementalPurgingPotionItemID  = 250954
	VenomousBloodPotionItemID     = 250950
	elementalPurgingSpellDamage   = 105
	dragonfireDirectDamage        = 280
	dragonfireTickDamage          = 35
	causticSmogTickDamage         = 33
	causticSmogTickDamageVariance = 0.1
)

// foreverPotionEffect returns what one of these potions does when drunk, given the target the
// potion is used on, and the aura to show as its buff (or nil). nil for every other potion.
func (character *Character) foreverPotionEffect(itemID int32) (func(*Simulation, *Unit), *Aura) {
	switch itemID {
	case DragonfirePotionItemID:
		spell := character.newDragonfirePotionSpell()
		return func(sim *Simulation, target *Unit) { spell.Cast(sim, target) }, nil
	case CausticSmogPotionItemID:
		spell := character.newCausticSmogPotionSpell()
		return func(sim *Simulation, target *Unit) { spell.Cast(sim, target) }, nil
	case ElementalPurgingPotionItemID:
		aura := character.newElementalPurgingAura()
		return func(sim *Simulation, _ *Unit) { aura.Activate(sim) }, aura
	}
	return nil, nil
}

// Dragonfire Potion (1252554): 280 Fire damage at once (SCHOOL_DAMAGE, no variance and no spell
// power coefficient) and 35 Fire every 2 sec for 8 sec (A_PERIODIC_DAMAGE), Magic defense, in a
// narrow cone in front of the drinker (TARGET_UNIT_CONE_CASTER_TO_DEST_ENEMY). Every target of
// the encounter stands in the cone, as Dragonbreath Chili's every target does.
func (character *Character) newDragonfirePotionSpell() *Spell {
	return character.GetOrRegisterSpell(SpellConfig{
		ActionID:    ActionID{SpellID: 1252554},
		SpellSchool: SpellSchoolFire,
		DefenseType: DefenseTypeMagic,
		ProcMask:    ProcMaskEmpty,
		Flags:       SpellFlagNoOnCastComplete | SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: DotConfig{
			Aura: Aura{
				Label: "Dragonfire Potion",
			},
			TickLength:    time.Second * 2,
			NumberOfTicks: 4,
			OnTick: func(sim *Simulation, target *Unit, dot *Dot) {
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, dragonfireTickDamage, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *Simulation, _ *Unit, spell *Spell) {
			for _, aoeTarget := range sim.Encounter.ActiveTargetUnits {
				result := spell.CalcDamage(sim, aoeTarget, dragonfireDirectDamage, spell.OutcomeMagicHitAndCrit)
				if result.Landed() {
					spell.Dot(aoeTarget).Apply(sim)
				}
				spell.DealDamage(sim, result)
			}
		},
	})
}

// Caustic Smog Potion (1252525): thrown at the target's feet, a cloud (1252526, 12 sec) whose
// enemies take 33 Nature (+-10%) every sec (1252543, A_PERIODIC_DAMAGE, Magic defense, no spell
// power coefficient). Every target of the encounter stands in the cloud for all of it: twelve ticks
// after one hit roll on entering.
func (character *Character) newCausticSmogPotionSpell() *Spell {
	return character.GetOrRegisterSpell(SpellConfig{
		ActionID:    ActionID{SpellID: 1252543},
		SpellSchool: SpellSchoolNature,
		DefenseType: DefenseTypeMagic,
		ProcMask:    ProcMaskEmpty,
		Flags:       SpellFlagNoOnCastComplete | SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: DotConfig{
			Aura: Aura{
				Label: "Caustic Smog Potion",
			},
			TickLength:    time.Second,
			NumberOfTicks: 12,
			OnTick: func(sim *Simulation, target *Unit, dot *Dot) {
				damage := causticSmogTickDamage * (1 + causticSmogTickDamageVariance*(2*sim.RandomFloat("Caustic Smog Potion")-1))
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, damage, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *Simulation, _ *Unit, spell *Spell) {
			for _, aoeTarget := range sim.Encounter.ActiveTargetUnits {
				result := spell.CalcOutcome(sim, aoeTarget, spell.OutcomeMagicHit)
				if result.Landed() {
					spell.Dot(aoeTarget).Apply(sim)
				}
				spell.DealOutcome(sim, result)
			}
		},
	})
}

// Potion of Elemental Purging (1252558): 105 spell damage against Elementals for 5 min
// (A_MOD_FLAT_SPELL_DAMAGE_VERSUS, creature type mask 8), paid the way the Elixir of Demonslaying
// pays its attack power against Demons.
func (character *Character) newElementalPurgingAura() *Aura {
	bonus := stats.Stats{stats.SpellDamage: elementalPurgingSpellDamage}
	return character.GetOrRegisterAura(Aura{
		Label:    "Potion of Elemental Purging",
		ActionID: ActionID{ItemID: ElementalPurgingPotionItemID},
		Duration: time.Minute * 5,
		OnGain: func(_ *Aura, _ *Simulation) {
			for _, at := range character.AttackTables {
				at.MobTypeBonusStats[proto.MobType_MobTypeElemental] = at.MobTypeBonusStats[proto.MobType_MobTypeElemental].Add(bonus)
			}
		},
		OnExpire: func(_ *Aura, _ *Simulation) {
			for _, at := range character.AttackTables {
				at.MobTypeBonusStats[proto.MobType_MobTypeElemental] = at.MobTypeBonusStats[proto.MobType_MobTypeElemental].Subtract(bonus)
			}
		},
	})
}
