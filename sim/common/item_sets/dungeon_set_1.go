package item_sets

import (
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/stats"
)

//We can place these in the class items_sets_pve.go if wanted.

// Forever's Dungeon Set 1 has bonuses at 2-6 pieces. Left out: the 4-piece when-struck effects and the caster 5-piece
// procs, which do nothing for a hunter's damage.

var ItemSetWildheartRaiment = core.NewItemSet(core.ItemSet{
	Name: "Wildheart Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		// +8 All Resistances.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddResistances(8)
		},
		// +30 Attack Power, and damage and healing done by magical spells and effects up to 18.
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       30,
				stats.RangedAttackPower: 30,
				stats.SpellDamage:       18,
				stats.HealingPower:      18,
			})
		},
		// Restores 8 mana per 5 sec.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MP5, 8)
		},
	},
})

// Forever's Dungeon Set 1 bonuses are at 2-6 pieces (the 3-piece attack power is reported not to
// apply in game; the sim follows the client data).
var ItemSetBeaststalkerArmor = core.NewItemSet(core.ItemSet{
	Name: "Beaststalker Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// +8 All Resistances.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddResistances(8)
		},
		// +30 Attack Power.
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       30,
				stats.RangedAttackPower: 30,
			})
		},
		// Your melee and ranged autoattacks have a 5% chance to energize you for 200 mana.
		5: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450577}
			manaMetrics := c.NewManaMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "Hunter Armor Energize",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskWhiteHit,
				ProcChance: 0.05,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if c.HasManaBar() {
						c.AddMana(sim, 200, manaMetrics)
					}
				},
			})
		},
		// Restores 8 mana per 5 sec.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MP5, 8)
		},
	},
})

var ItemSetMagistersRegalia = core.NewItemSet(core.ItemSet{
	Name: "Magister's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		// +8 All Resistances.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddResistances(8)
		},
		// Increases damage and healing done by magical spells and effects by up to 18.
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.SpellDamage:  18,
				stats.HealingPower: 18,
			})
		},
		// Restores 8 mana per 5 sec.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MP5, 8)
		},
	},
})

var ItemSetLightforgeArmor = core.NewItemSet(core.ItemSet{
	Name: "Lightforge Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// +200 Armor.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
		// +40 Attack Power
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
		// Chance on melee attack to increase your damage and healing done by magical spells and effects by up to 95 for 10 sec.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 27498}

			procAura := c.NewTemporaryStatsAura("Crusader's Wrath", actionID, stats.Stats{stats.SpellPower: 95}, time.Second*10)
			handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			}

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "Item - Crusader's Wrath Proc - Lightforge Armor",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskMeleeWhiteHit,
				PPM:        3, //6% is the sod proc rate, best data I could find online says 3 ppm
				Handler:    handler,
			})
		},
		// +8 All Resistances.
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddResistances(8)
		},
	},
})

var ItemSetVestmentsOfTheDevout = core.NewItemSet(core.ItemSet{
	Name: "Vestments of the Devout",
	Bonuses: map[int32]core.ApplyEffect{
		// +8 All Resistances.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddResistances(8)
		},
		// Increases damage and healing done by magical spells and effects by up to 18.
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.SpellDamage:  18,
				stats.HealingPower: 18,
			})
		},
		// Restores 8 mana per 5 sec.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MP5, 8)
		},
	},
})

var ItemSetShadowcraftArmor = core.NewItemSet(core.ItemSet{
	Name: "Shadowcraft Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// +8 All Resistances.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddResistances(8)
		},
		// +30 Attack Power.
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       30,
				stats.RangedAttackPower: 30,
			})
		},
		// Chance on melee attack to restore 35 energy.
		5: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 27787}
			energyMetrics := c.NewEnergyMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID: actionID,
				Name:     "Rogue Armor Energize",
				Callback: core.CallbackOnSpellHitDealt,
				Outcome:  core.OutcomeLanded,
				ProcMask: core.ProcMaskMeleeWhiteHit,
				PPM:      1,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if c.HasEnergyBar() {
						c.AddEnergy(sim, 35, energyMetrics)
					}
				},
			})
		},
	},
})

var ItemSetTheElements = core.NewItemSet(core.ItemSet{
	Name: "The Elements",
	Bonuses: map[int32]core.ApplyEffect{
		// +8 All Resistances.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddResistances(8)
		},
		// Increases damage and healing done by magical spells and effects by up to 18.
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.SpellDamage:  18,
				stats.HealingPower: 18,
			})
		},
		// Restores 8 mana per 5 sec.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MP5, 8)
		},
	},
})

var ItemSetDreadmistRaiment = core.NewItemSet(core.ItemSet{
	Name: "Dreadmist Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		// +8 All Resistances.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddResistances(8)
		},
		// Increases damage and healing done by magical spells and effects by up to 18.
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.SpellDamage:  18,
				stats.HealingPower: 18,
			})
		},
		// Restores 8 mana per 5 sec.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MP5, 8)
		},
	},
})

var ItemSetBattlegearOfValor = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Valor",
	Bonuses: map[int32]core.ApplyEffect{
		// +200 Armor.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Armor, 200)
		},
		// +40 Attack Power.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
		// Chance on melee attack to heal you for 88 to 132
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 27419}
			healthMetrics := c.NewHealthMetrics(core.ActionID{SpellID: 27419})

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID: actionID,
				Name:     "Warrior's Resolve",
				Callback: core.CallbackOnSpellHitDealt,
				Outcome:  core.OutcomeLanded,
				ProcMask: core.ProcMaskMelee,
				PPM:      1,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					c.GainHealth(sim, sim.Roll(88, 133), healthMetrics)
				},
			})
		},
		// +8 All Resistances.
		8: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddResistances(8)
		},
	},
})
