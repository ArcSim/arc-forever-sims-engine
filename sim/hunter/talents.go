package hunter

import (
	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
	"os"
	"time"
)

// false = old behavior: keep the stronger bleed and restart the timer.
// true = rollover: add unticked damage, keep the existing tick time.
const laceratingRollover = true

func (hunter *Hunter) ApplyTalents() {
	if hunter.talentsApplied {
		return
	}
	hunter.talentsApplied = true
	if hunter.pet != nil {
		hunter.applyFrenzy()
		hunter.registerBestialWrathCD()
		hunter.registerIntimidationCD()

		// Forever per-rank values: 2% crit per rank for pets and hawks
		ferocity := 2.0 * float64(hunter.Talents.Ferocity)
		hunter.pet.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*ferocity)
		hunter.pet.AddStat(stats.SpellCrit, core.SpellCritRatingPerCritChance*ferocity)

		hunter.pet.PseudoStats.DamageDealtMultiplier *= 1 + 0.03*float64(hunter.Talents.UnleashedFury)

		if hunter.Talents.EnduranceTraining > 0 {
			mult := 1 + 0.03*float64(hunter.Talents.EnduranceTraining)
			hunter.pet.MultiplyStat(stats.Health, mult)
			hunter.pet.MultiplyStat(stats.Armor, mult)
		}

		if hunter.Talents.FocusedFire > 0 {
			ff := 1 + 0.01*float64(hunter.Talents.FocusedFire)
			hunter.PseudoStats.DamageDealtMultiplier *= ff
			hunter.pet.PseudoStats.DamageDealtMultiplier *= ff
		}
	}

	if hunter.Talents.SavageStrikes > 0 {
		bonus := 2 * float64(hunter.Talents.SavageStrikes) * core.CritRatingPerCritChance
		hunter.OnSpellRegistered(func(spell *core.Spell) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeMHSpecial|core.ProcMaskMeleeOHSpecial) && !spell.Flags.Matches(SpellFlagTrap) {
				spell.BonusCritRating += bonus
			}
		})
	}

	if hunter.Talents.BestialDiscipline > 0 {
		core.MakePermanent(hunter.RegisterAura(core.Aura{
			Label: "Bestial Discipline",
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				// raises the pet passive's focus regen percent, so it scales the doubled Forever rate
				if hunter.pet != nil {
					hunter.pet.MultiplyFocusRegen(1 + 0.1*float64(hunter.Talents.BestialDiscipline))
				}
			},
		}))
		hunter.PseudoStats.SpiritRegenRateCasting += 0.25 * float64(hunter.Talents.BestialDiscipline)
	}

	if hunter.Talents.CarefulAim > 0 {
		intelToAP := 0.20 * float64(hunter.Talents.CarefulAim)
		hunter.AddStatDependency(stats.Intellect, stats.AttackPower, intelToAP)
		hunter.AddStatDependency(stats.Intellect, stats.RangedAttackPower, intelToAP)
	}

	hunter.AddStat(stats.MeleeHit, float64(hunter.Talents.Surefooted)*core.MeleeHitRatingPerHitChance)
	hunter.AddStat(stats.SpellHit, float64(hunter.Talents.Surefooted)*core.SpellHitRatingPerHitChance)

	if hunter.Talents.LethalAttacks > 0 {
		bonus := float64(hunter.Talents.LethalAttacks) * core.CritRatingPerCritChance
		hunter.AddStat(stats.MeleeCrit, bonus)
		hunter.AddStat(stats.SpellCrit, bonus)
	}

	if hunter.Talents.RangedWeaponSpecialization > 0 {
		mult := 1 + 0.01*float64(hunter.Talents.RangedWeaponSpecialization)
		hunter.OnSpellRegistered(func(spell *core.Spell) {
			if spell.ProcMask.Matches(core.ProcMaskRanged) && spell.SpellCode != SpellCode_HunterSerpentSting {
				spell.DamageMultiplier *= mult
			}
		})
	}

	if hunter.Talents.Survivalist > 0 {
		hunter.MultiplyStat(stats.Health, 1.0+0.02*float64(hunter.Talents.Survivalist))
	}

	// 2% Agility per rank: the node in the tree (a 3% copy sits off the tree where it cannot be taken)
	if hunter.Talents.LightningReflexes > 0 {
		hunter.MultiplyStat(stats.Agility, 1.0+0.02*float64(hunter.Talents.LightningReflexes))
	}

	if hunter.Talents.PredatorsEdge > 0 {
		critDmg := 0.06 * float64(hunter.Talents.PredatorsEdge)
		ohMult := 1 + 0.10*float64(hunter.Talents.PredatorsEdge)
		apply := func(spell *core.Spell) {
			if spell == nil {
				return
			}
			if spell.ProcMask.Matches(core.ProcMaskMeleeMHSpecial|core.ProcMaskMeleeOHSpecial) && !spell.Flags.Matches(SpellFlagTrap) {
				spell.CritDamageBonus += critDmg
			}
			if spell.ProcMask.Matches(core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeOHSpecial) {
				spell.DamageMultiplier *= ohMult
			}
		}
		hunter.OnSpellRegistered(apply)
		apply(hunter.AutoAttacks.OHAuto())
	}

	if hunter.Talents.ImprovedTracking > 0 {
		hunter.PseudoStats.DamageDealtMultiplier *= 1 + 0.01*float64(hunter.Talents.ImprovedTracking)
	}

	hunter.applyEfficiency()
	hunter.applyCleverTraps()
	hunter.applySurvivalTactics()
	hunter.applyResourcefulness()
	hunter.applySurvivalistDiscipline()
	hunter.applyRapidRecuperation()
	hunter.applyLaceratingStrikes()
	hunter.applyTrueshotAura()
}

// The talent teaches Trueshot Aura, cast before the pull (it lasts 30 minutes), so the hunter has it even when the raid
// buffs leave it out.
func (hunter *Hunter) applyTrueshotAura() {
	if hunter.Talents.TrueshotAura {
		core.TrueshotAura(&hunter.Unit)
	}
}

func (hunter *Hunter) applyExposePrey() {
	if hunter.Talents.ExposePrey == 0 || hunter.DefensiveState == nil {
		return
	}
	procChance := 0.05 * float64(hunter.Talents.ExposePrey)

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: "Expose Prey",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}
			if sim.RandomFloat("Expose Prey") < procChance {
				hunter.DefensiveState.Activate(sim)
			}
		},
	}))
}

func (hunter *Hunter) applyLaceratingStrikes() {
	if !hunter.Talents.LaceratingStrikes {
		return
	}

	bank := make([]float64, len(hunter.Env.AllUnits))

	// AB_LAC_MODE=refresh|strongest|rolling overrides the setting (harness)
	mode := hunter.Options.LaceratingMode
	switch os.Getenv("AB_LAC_MODE") {
	case "refresh":
		mode = proto.Hunter_Options_LaceratingRefresh
	case "strongest":
		mode = proto.Hunter_Options_LaceratingStrongest
	case "rolling":
		mode = proto.Hunter_Options_LaceratingRolling
	}

	hunter.LaceratingBleed = hunter.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 1310533},
		SpellSchool:      core.SpellSchoolPhysical,
		DefenseType:      core.DefenseTypeMelee,
		ProcMask:         core.ProcMaskEmpty,
		Flags:            core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagIgnoreTargetModifiers,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Lacerating Strikes",
				ActionID: core.ActionID{SpellID: 1310533},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					bank[aura.Unit.UnitIndex] = 0
				},
			},
			NumberOfTicks: 7,
			TickLength:    time.Second * 3,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				at := dot.Spell.Unit.AttackTables[target.UnitIndex][proto.CastType_CastTypeMainHand]
				dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(at)
				// AB_LAC_NOCRIT=1: bleed ticks cannot crit (as Classic bleeds)
				if os.Getenv("AB_LAC_NOCRIT") != "" {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				} else {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickPhysicalCrit)
				}
				bank[target.UnitIndex] -= dot.SnapshotBaseDamage
				if bank[target.UnitIndex] < 0 {
					bank[target.UnitIndex] = 0
				}
			},
		},
	})

	if hunter.Talents.SavageStrikes > 0 {
		hunter.LaceratingBleed.BonusCritRating += 2 * float64(hunter.Talents.SavageStrikes) * core.CritRatingPerCritChance
	}
	if hunter.Talents.PredatorsEdge > 0 {
		hunter.LaceratingBleed.CritDamageBonus += 0.06 * float64(hunter.Talents.PredatorsEdge)
	}

	hunter.RegisterAura(core.Aura{
		Label:    "Lacerating Strikes Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			for i := range bank {
				bank[i] = 0
			}
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell != hunter.MongooseBite || !result.Landed() || result.Damage <= 0 {
				return
			}
			dot := hunter.LaceratingBleed.Dot(result.Target)
			idx := result.Target.UnitIndex
			at := hunter.AttackTables[result.Target.UnitIndex][proto.CastType_CastTypeMainHand]
			fresh := result.Damage * 0.40

			if mode != proto.Hunter_Options_LaceratingRolling {
				// a weaker bleed cannot overwrite a stronger one (Strongest mode)
				if mode == proto.Hunter_Options_LaceratingStrongest && dot.IsActive() && fresh/7 < dot.SnapshotBaseDamage {
					return
				}
				// the new bleed replaces the old one: unticked damage is lost, the tick timer restarts
				dot.SnapshotBaseDamage = fresh / 7
				dot.SnapshotAttackerMultiplier = 1
				dot.SnapshotCritChance = hunter.LaceratingBleed.PhysicalCritChance(at)
				dot.NumberOfTicks = 7
				dot.ApplyOrReset(sim)
				bank[idx] = fresh
				return
			}

			total := fresh + bank[idx]
			dot.SnapshotBaseDamage = total / 7
			dot.SnapshotAttackerMultiplier = 1
			dot.SnapshotCritChance = hunter.LaceratingBleed.PhysicalCritChance(at)
			bank[idx] = total

			if !dot.IsActive() {
				dot.Apply(sim)
				return
			}

			nextTick := dot.NextTickAt()
			dot.TickCount = 0
			dot.NumberOfTicks = 7
			dot.RecomputeAuraDuration()
			dot.Aura.Refresh(sim)
			if pa := dot.TickAction(); pa != nil {
				pa.NextActionAt = nextTick
			}
		},
	})
}

// Barrage per rank: 3/7/10% damage to Multi-Shot, Aimed Shot and Volley
func (hunter *Hunter) barrageBonus() float64 {
	return []float64{0, 0.03, 0.07, 0.10}[hunter.Talents.Barrage]
}

func (hunter *Hunter) improvedStingsSerpentBonus() float64 {
	return []float64{0, 0.06, 0.13, 0.20}[hunter.Talents.ImprovedStings]
}

func (hunter *Hunter) applyFrenzy() {
	if hunter.Talents.Frenzy == 0 || hunter.pet == nil {
		return
	}

	procChance := 0.2 * float64(hunter.Talents.Frenzy)

	procAura := hunter.pet.RegisterAura(core.Aura{
		Label: "Frenzy Proc",
		// the Forever Frenzy buff (19625 is the Classic talent rank)
		ActionID: core.ActionID{SpellID: 19615},
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, 1.3)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, 1/1.3)
		},
	})

	hunter.pet.RegisterAura(core.Aura{
		Label:    "Frenzy",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellResult *core.SpellResult) {
			if !spellResult.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			if procChance == 1 || sim.RandomFloat("Frenzy") < procChance {
				procAura.Activate(sim)
			}
		},
	})
}

func (hunter *Hunter) registerBestialWrathCD() {
	if !hunter.Talents.BestialWrath || hunter.pet == nil {
		return
	}

	actionID := core.ActionID{SpellID: 19574}

	hunter.BestialWrathPetAura = hunter.pet.RegisterAura(core.Aura{
		Label:    "Bestial Wrath Pet",
		ActionID: actionID,
		Duration: time.Second * 18,
	}).AttachMultiplicativePseudoStatBuff(&hunter.pet.PseudoStats.DamageDealtMultiplier, 1.5)

	bwSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,
		ManaCost: core.ManaCostOptions{BaseCost: 0.12},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Minute * 2,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.BestialWrathPetAura.Activate(sim)
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: bwSpell,
		Type:  core.CooldownTypeDPS,
	})
}

// Intimidation (19577): the pet's next landed attack gains 100% critical strike chance (its stun does not
// matter on a boss). 8% base mana, 1 min cooldown, on the global cooldown, waits up to 15 s for the attack.
func (hunter *Hunter) registerIntimidationCD() {
	if !hunter.Talents.Intimidation || hunter.pet == nil {
		return
	}

	actionID := core.ActionID{SpellID: 19577}
	critBonus := float64(100 * core.CritRatingPerCritChance)

	// a white swing rolled just before the cast lands a moment later and does not use the buff
	var gainedAt time.Duration
	hunter.IntimidationPetAura = hunter.pet.RegisterAura(core.Aura{
		Label:    "Intimidation Pet",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			gainedAt = sim.CurrentTime
			aura.Unit.AddStatDynamic(sim, stats.MeleeCrit, critBonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.MeleeCrit, -critBonus)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) || !result.Landed() {
				return
			}
			if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) && sim.CurrentTime-gainedAt < 50*time.Millisecond {
				return
			}
			aura.Deactivate(sim)
		},
	})

	intimidation := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,
		ManaCost: core.ManaCostOptions{BaseCost: 0.08},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{GCD: core.GCDDefault},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Minute,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.IntimidationPetAura.Activate(sim)
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: intimidation,
		Type:  core.CooldownTypeDPS,
	})
}

func (hunter *Hunter) mortalShots() float64 {
	return 0.06 * float64(hunter.Talents.MortalShots)
}

func (hunter *Hunter) applyCleverTraps() {
	if hunter.Talents.CleverTraps == 0 {
		return
	}
	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagTrap) {
			spell.DamageMultiplier *= 1 + 0.15*float64(hunter.Talents.CleverTraps)
		}
	})
}

// Survival Tactics: +5% chance to hit with traps per rank. The trap outcome strips the unit's spell
// hit on purpose, so the bonus lives on the trap spells themselves.
func (hunter *Hunter) applySurvivalTactics() {
	if hunter.Talents.SurvivalTactics == 0 {
		return
	}
	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagTrap) {
			spell.BonusHitRating += 5 * float64(hunter.Talents.SurvivalTactics) * core.SpellHitRatingPerHitChance
		}
	})
}

func (hunter *Hunter) applyEfficiency() {
	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Cost == nil || spell.ActionID.SpellID == StriderKickSpellID || isSniperShot(spell) {
			return
		}
		if spell.Flags.Matches(SpellFlagSting|SpellFlagShot|SpellFlagStrike) || spell.SpellCode == SpellCode_HunterVolley || spell.SpellCode == SpellCode_HunterRaptorStrikeHit || spell.SpellCode == SpellCode_HunterWingClip {
			spell.Cost.Multiplier -= 3 * hunter.Talents.Efficiency
		}
	})
}

func (hunter *Hunter) applyResourcefulness() {
	if hunter.Talents.Resourcefulness == 0 {
		return
	}

	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Cost == nil {
			return
		}
		// traps and melee abilities, Strider Kick included
		if spell.Flags.Matches(SpellFlagTrap|SpellFlagStrike) || spell.SpellCode == SpellCode_HunterRaptorStrikeHit || spell.SpellCode == SpellCode_HunterWingClip {
			spell.Cost.Multiplier -= 30 * hunter.Talents.Resourcefulness
		}
	})

	procChance := 0.50 * float64(hunter.Talents.Resourcefulness)
	hunter.resourcefulRegen = hunter.NewCastingRegenSource(core.ActionID{SpellID: 1242688}, 0.5)
	aura := hunter.RegisterAura(core.Aura{
		Label:    "Resourcefulness",
		ActionID: core.ActionID{SpellID: 1242688},
		Duration: time.Second * 30,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.resourcefulRegen.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.resourcefulRegen.Deactivate()
		},
	})

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: "Resourcefulness Proc",
		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// melee and ranged ability crits only: not auto attacks, not traps
			if !result.DidCrit() || !spell.ProcMask.Matches(core.ProcMaskMeleeSpecial|core.ProcMaskRangedSpecial) || spell.Flags.Matches(SpellFlagTrap) {
				return
			}
			if sim.RandomFloat("Resourcefulness") < procChance {
				aura.Activate(sim)
			}
		},
	}))
}

func (hunter *Hunter) applySurvivalistDiscipline() {
	if hunter.Talents.SurvivalistDiscipline == 0 {
		return
	}
	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagTrap) && spell.CD.Timer != nil {
			spell.CD.Duration = time.Duration(float64(spell.CD.Duration) * (1 - 0.20*float64(hunter.Talents.SurvivalistDiscipline)))
		}
	})
}

func (hunter *Hunter) applyRapidRecuperation() {
	if hunter.Talents.RapidRecuperation == 0 {
		return
	}
	if hunter.rapidRecupAura != nil {
		return
	}

	hunter.rapidRecupRegen = hunter.NewCastingRegenSource(core.ActionID{SpellID: 1223987}, 0.25*float64(hunter.Talents.RapidRecuperation))
	hunter.rapidRecupAura = hunter.RegisterAura(core.Aura{
		Label:    "Rapid Recuperation Sting",
		ActionID: core.ActionID{SpellID: 1223987},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.rapidRecupRegen.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.rapidRecupRegen.Deactivate()
		},
	})

	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellCode != SpellCode_HunterSerpentSting {
			return
		}
		old := spell.ApplyEffects
		spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, sp *core.Spell) {
			old(sim, target, sp)
			hunter.rapidRecupAura.Activate(sim)
		}
	})
}

// Strider Kick and Sniper Shot are outside the Efficiency spell list.
const StriderKickSpellID = 1317257

func isSniperShot(spell *core.Spell) bool {
	id := spell.ActionID.SpellID
	return id == 1310687 || id == 1310785 || id == 1310786
}
