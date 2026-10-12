package classic

import (
	"time"

	"github.com/wowsims/forever/sim/common/itemhelpers"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

// IronfoeMainHandOnly and IronfoePPM are the switches for Ironfoe's Fury of Forgewright. The client
// row is the default: 1301046 is an equip aura (ItemEffect TriggerType 1 on Ironfoe, 11684) with
// ProcChance 6, ProcsPerMinuteID 0, a 100 ms ProcCategoryRecovery and ProcTypeMask 0x14 (melee autos
// and specials), and no SpellEquippedItems row ties it to the weapon, so off-hand hits and specials
// roll it at 6%. Its tooltip reads "Attacks have a chance to grant 2 extra attacks on your next
// swing". A player report (MythicSim Discord, 12 October 2026) wants Era's rate instead: 0.8 PPM off
// the main hand (3.2% a swing on its 2.4 sec speed). Set IronfoeMainHandOnly to drop the off-hand
// weapon's rolls and IronfoePPM above 0 to replace the 6% with that rate.
var IronfoeMainHandOnly = false
var IronfoePPM = 0.0

func init() {
	// Thunderfury, Blessed Blade of the Windseeker
	itemhelpers.CreateWeaponProcTrigger(itemhelpers.WeaponProcTrigger{
		ItemID: 19019,
		Name:   "Thunderfury",
		PPM:    6,
		Handler: func(character *core.Character) core.ProcHandler {
			procActionID := core.ActionID{SpellID: 21992}

			attackSpeedDebuffAura := character.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
				aura := target.GetOrRegisterAura(core.Aura{
					Label:    "Cyclone",
					ActionID: core.ActionID{SpellID: 27648},
					Duration: time.Second * 12,
				})

				core.AtkSpeedReductionEffect(aura, core.SlowedTimeMultiplier(-20))

				return aura
			})

			singleTargetSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    procActionID.WithTag(1),
				SpellSchool: core.SpellSchoolNature,
				DefenseType: core.DefenseTypeMagic,
				ProcMask:    core.ProcMaskSpellDamage,
				Flags:       core.SpellFlagProc,

				DamageMultiplier: 1,
				ThreatMultiplier: 0.5,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					result := spell.CalcAndDealDamage(sim, target, 300, spell.OutcomeMagicHitAndCrit)
					if result.Landed() {
						attackSpeedDebuffAura.Get(result.Target).Activate(sim)
					}
				},
			})

			resistanceDebuffAura := character.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
				return target.GetOrRegisterAura(core.Aura{
					Label:    "Thunderfury",
					ActionID: procActionID,
					Duration: time.Second * 12,
					OnGain: func(aura *core.Aura, sim *core.Simulation) {
						target.AddStatDynamic(sim, stats.NatureResistance, -25)
					},
					OnExpire: func(aura *core.Aura, sim *core.Simulation) {
						target.AddStatDynamic(sim, stats.NatureResistance, 25)
					},
				})
			})

			bounceSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    procActionID.WithTag(2),
				SpellSchool: core.SpellSchoolNature,
				DefenseType: core.DefenseTypeMagic,
				ProcMask:    core.ProcMaskEmpty,

				ThreatMultiplier: 1,
				FlatThreatBonus:  63,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					results := spell.CalcCleaveDamage(sim, target, 5, 0, spell.OutcomeMagicHit)
					for _, result := range results {
						if result.Landed() {
							resistanceDebuffAura.Get(result.Target).Activate(sim)
						}
					}
					spell.DealBatchedAoeDamage(sim)
				},
			})

			return func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				singleTargetSpell.Cast(sim, result.Target)
				bounceSpell.Cast(sim, result.Target)
			}
		},
	})

	// Dragon's Call: chance on hit, Dragon's Call (13049) summons an Emerald Dragon Whelp for 15 sec.
	// The client gives no proc rate; 1 PPM is master's (Armaments Discord). 13049's category cooldown
	// gates the proc: Forever retuned it from Classic Era's 60 sec to 45 sec, and on Classic Era the
	// whelp could not proc again for a minute after a summon. Without it the 1 PPM roll kept
	// refreshing the whelp: 50% uptime on the melee Survival Hunter reference, 25% with it.
	summonSpell := spelldata.MustFind(13049)
	itemhelpers.CreateWeaponProcTrigger(itemhelpers.WeaponProcTrigger{
		ItemID: DragonsCall,
		Name:   "Emerald Dragon Whelp",
		PPM:    1,
		ICD:    time.Millisecond * time.Duration(summonSpell.CategoryCooldownMs),
		Handler: func(character *core.Character) core.ProcHandler {
			for _, petAgent := range character.PetAgents {
				if whelp, ok := petAgent.(*EmeraldDragonWhelp); ok {
					return func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
						whelp.summon(sim, time.Millisecond*time.Duration(summonSpell.DurationMs))
					}
				}
			}
			return nil
		},
	})

	// Ironfoe: Forever made it an equip, Fury of Forgewright (1301046): 6% chance on melee hit with
	// a 100 ms proc cooldown to grant 2 extra attacks (15494). Master had 0.8 PPM off Ironfoe's own
	// hits. The "twice as likely against Orcs" part is not modelled.
	core.NewItemEffect(11684, func(agent core.Agent) {
		character := agent.GetCharacter()
		trigger := core.ProcTrigger{
			Name:               "Fury of Forgewright",
			ActionID:           core.ActionID{SpellID: 15494},
			Callback:           core.CallbackOnSpellHitDealt,
			ProcMask:           core.ProcMaskMelee,
			Outcome:            core.OutcomeLanded,
			ProcChance:         0.06,
			ICD:                time.Millisecond * 100,
			TriggerImmediately: true,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				character.AutoAttacks.ExtraMHAttacks(sim, 2)
			},
		}
		if IronfoeMainHandOnly {
			// A shield strike (an off-hand special with no off-hand weapon) is a main-hand attack to
			// the server, as for Windfury Totem (buffs.NarrowWindfuryProc).
			trigger.ProcMask = core.ProcMaskMeleeMH | core.ProcMaskMeleeOHSpecial
			trigger.ExtraCondition = func(_ *core.Simulation, spell *core.Spell, _ *core.SpellResult) bool {
				return !spell.ProcMask.Matches(core.ProcMaskMeleeOHSpecial) || !character.HasOHWeapon()
			}
		}
		if IronfoePPM > 0 {
			trigger.ProcChance = 0
			trigger.DPM = character.NewLegacyPPMManager(IronfoePPM, trigger.ProcMask)
		}
		aura := character.MakeProcTriggerAura(trigger)
		character.ItemSwap.RegisterProc(11684, aura)
	})

	// Runeblade of Baron Rivendare: Unholy Aura (17625), 60 health every 5 sec and 8% movement speed.
	core.NewItemEffect(13505, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 17625}
		healthMetrics := character.NewHealthMetrics(actionID)

		var regen *core.PendingAction
		aura := character.NewPassiveMovementSpeedAura("Unholy Aura", actionID, 0.08).
			ApplyOnGain(func(_ *core.Aura, sim *core.Simulation) {
				regen = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
					Period:   time.Second * 5,
					Priority: core.ActionPriorityAuto,
					OnAction: func(sim *core.Simulation) {
						character.GainHealth(sim, 60, healthMetrics)
					},
				})
			}).
			ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
				regen.Cancel(sim)
			})
		character.ItemSwap.RegisterProc(13505, aura)
	})

	// Headmaster's Charge: use, 20 Intellect for 15 min on a 10 min cooldown (18264). It is a party
	// aura; only the wearer is simulated. Master had it commented out, at 25.
	core.NewItemEffect(13937, func(agent core.Agent) {
		character := agent.GetCharacter()
		core.RegisterTemporaryStatsOnUseCD(character, "Headmaster's Charge", stats.Stats{stats.Intellect: 20}, time.Minute*15, core.SpellConfig{
			ActionID: core.ActionID{ItemID: 13937},
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 10,
				},
			},
		})
	})

	// Ebon Hilt of Marduk: chance on hit, Corruption (18656), 28 Shadow every 3 sec for 9 sec. Forever
	// adds an equip, 1% less threat (1298501). 1 PPM is master's guess.
	itemhelpers.CreateWeaponProcSpell(itemhelpers.WeaponProcSpell{
		ItemID: 14576,
		Name:   "Ebon Hilt of Marduk",
		PPM:    1,
		Spell: func(character *core.Character) *core.Spell {
			threatAura := core.MakePermanent(character.RegisterAura(core.Aura{
				Label:    "Decrease Threat All 01",
				ActionID: core.ActionID{SpellID: 1298501},
			})).AttachMultiplicativePseudoStatBuff(&character.PseudoStats.ThreatMultiplier, 0.99)
			character.ItemSwap.RegisterProc(14576, threatAura)

			return character.GetOrRegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 18656},
				SpellSchool: core.SpellSchoolShadow,
				DefenseType: core.DefenseTypeMagic,
				ProcMask:    core.ProcMaskEmpty,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				Dot: core.DotConfig{
					Aura: core.Aura{
						Label: "Corruption (Ebon Hilt of Marduk)",
					},
					TickLength:    time.Second * 3,
					NumberOfTicks: 3,
					OnSnapshot: func(_ *core.Simulation, target *core.Unit, dot *core.Dot) {
						dot.Snapshot(target, 28)
					},
					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
					},
				},

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
					if result.Landed() {
						spell.Dot(target).Apply(sim)
					}
					spell.DealOutcome(sim, result)
				},
			})
		},
	})

	// Sulfuras, Hand of Ragnaros: chance on hit, Fireball (21162), 303 +-19.8% (273 to 333) Fire plus
	// 15 every 2 sec for 10 sec; 1 PPM is master's (Armaments Discord). Equip, Immolation (21142): 5
	// Fire to every melee attacker. 21142 has no defense type, so unlike master it cannot miss.
	itemhelpers.CreateWeaponProcSpell(itemhelpers.WeaponProcSpell{
		ItemID: 17182,
		Name:   "Sulfuras, Hand of Ragnaros",
		PPM:    1,
		Spell: func(character *core.Character) *core.Spell {
			immolationSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 21142},
				SpellSchool: core.SpellSchoolFire,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealDamage(sim, target, 5, spell.OutcomeAlwaysHit)
				},
			})
			immolationAura := character.MakeProcTriggerAura(core.ProcTrigger{
				Name:     "Immolation (Hand of Ragnaros)",
				Callback: core.CallbackOnSpellHitTaken,
				ProcMask: core.ProcMaskMelee,
				Outcome:  core.OutcomeLanded,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					immolationSpell.Cast(sim, spell.Unit)
				},
			})
			character.ItemSwap.RegisterProc(17182, immolationAura)

			return character.GetOrRegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 21162},
				SpellSchool: core.SpellSchoolFire,
				DefenseType: core.DefenseTypeMagic,
				ProcMask:    core.ProcMaskEmpty,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				Dot: core.DotConfig{
					Aura: core.Aura{
						Label: "Fireball (Hand of Ragnaros)",
					},
					TickLength:    time.Second * 2,
					NumberOfTicks: 5,
					OnSnapshot: func(_ *core.Simulation, target *core.Unit, dot *core.Dot) {
						dot.Snapshot(target, 15)
					},
					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
					},
				},

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					result := spell.CalcAndDealDamage(sim, target, sim.Roll(273, 333), spell.OutcomeMagicHitAndCrit)
					if result.Landed() {
						spell.Dot(target).Apply(sim)
					}
				},
			})
		},
	})
}
