package forever

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

// The eight slots of a crafted plate set.
var platesetSlots = []proto.ItemSlot{
	proto.ItemSlot_ItemSlotHead,
	proto.ItemSlot_ItemSlotShoulder,
	proto.ItemSlot_ItemSlotChest,
	proto.ItemSlot_ItemSlotWrist,
	proto.ItemSlot_ItemSlotHands,
	proto.ItemSlot_ItemSlotWaist,
	proto.ItemSlot_ItemSlotLegs,
	proto.ItemSlot_ItemSlotFeet,
}

// Blacksmithing - Plate
// https://www.wowhead.com/forever/item-set=321/imperial-plate
//
// Forever reworked the set: Imperial Plate Gauntlets (250589) joined the seven Classic pieces and
// every bonus is new. The 4-piece bonus keeps movement speed from dropping below 80% and has no
// place in the sim.
var ItemSetImperialPlate = core.NewItemSet(core.ItemSet{
	ID:    321,
	Name:  "Imperial Plate",
	Slots: platesetSlots,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			// Increased Defense +7.
			setBonusAura.
				AttachStatBuff(stats.DefenseRating, 7*core.DefenseRatingPerDefenseLevel).
				ExposeToAPL(13385)
		},
		3: func(_ core.Agent, setBonusAura *core.Aura) {
			// Improves your chance to hit by 1%.
			setBonusAura.
				AttachStatsBuff(stats.Stats{stats.PhysicalHitPercent: 1, stats.SpellHitPercent: 1}).
				ExposeToAPL(1251990)
		},
		5: func(_ core.Agent, setBonusAura *core.Aura) {
			// Increased Strength +20.
			setBonusAura.
				AttachStatBuff(stats.Strength, 20).
				ExposeToAPL(1251984)
		},
		6: func(_ core.Agent, setBonusAura *core.Aura) {
			// Improves your chance to get a critical strike by 1%.
			setBonusAura.
				AttachStatsBuff(stats.Stats{stats.PhysicalCritPercent: 1, stats.SpellCritPercent: 1}).
				ExposeToAPL(1251991)
		},
	},
})

// Blacksmithing - Plate
// https://www.wowhead.com/forever/item-set=1968/blessed-plate
//
// New in Forever. The 4-piece bonus restores 20 health every 5 sec to party members within 10
// yards; the sim has no health regeneration stat, so it is left out.
var ItemSetBlessedPlate = core.NewItemSet(core.ItemSet{
	ID:    1968,
	Name:  "Blessed Plate",
	Slots: platesetSlots,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			// Increased Shadow Resistance +10.
			setBonusAura.
				AttachStatBuff(stats.ShadowResistance, 10).
				ExposeToAPL(14673)
		},
		3: func(_ core.Agent, setBonusAura *core.Aura) {
			// Increased Strength +20.
			setBonusAura.
				AttachStatBuff(stats.Strength, 20).
				ExposeToAPL(1252009)
		},
		5: func(_ core.Agent, setBonusAura *core.Aura) {
			// Improves your chance to hit by 2%.
			setBonusAura.
				AttachStatsBuff(stats.Stats{stats.PhysicalHitPercent: 2, stats.SpellHitPercent: 2}).
				ExposeToAPL(1252013)
		},
		6: func(_ core.Agent, setBonusAura *core.Aura) {
			// Increases damage done by Holy spells and effects by up to 29.
			setBonusAura.
				AttachStatBuff(stats.HolyDamage, 29).
				ExposeToAPL(21518)
		},
	},
})

// Leather
// https://www.wowhead.com/forever/item-set=161/defias-leather
//
// Forever rebuilt the Deadmines set: the 3 and 4 piece bonuses are new and the old +1 dagger skill
// moved to five pieces.
var ItemSetDefiasLeather = core.NewItemSet(core.ItemSet{
	ID:   161,
	Name: "Defias Leather",
	Bonuses: map[int32]core.ApplySetBonus{
		// Increased Arcane Resistance +5 (14707).
		2: setStats(stats.Stats{stats.ArcaneResistance: 5}),
		// +15 Attack Power against Humanoids (1292025).
		3: setHumanoidAttackPower(15),
		// Devious Strike (1292028): white hits from behind have a 5% chance to make the target
		// bleed 15 a second for 5 sec (1292029), ticks that can crit since client 70009.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			character := agent.GetCharacter()
			if character.PseudoStats.InFrontOfTarget {
				return
			}
			bleed := character.RegisterSpell(shared.SpellDataProcDamageSpell(character, spelldata.MustFind(1292029)))
			trigger := spelldata.ProcTrigger(character, spelldata.MustFind(1292028), func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				bleed.Cast(sim, result.Target)
			})
			trigger.TriggerImmediately = true
			setBonusAura.AttachProcTrigger(trigger)
		},
		// Increased Daggers +1 (7534).
		5: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachWeaponSkill(proto.PseudoStat_PseudoStatDaggersSkill, 1)
		},
	},
})
