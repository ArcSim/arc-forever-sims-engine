package druid

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

type DruidForm uint8

const (
	Humanoid DruidForm = 1 << iota
	Bear
	Cat
	Moonkin
	Tree
	Any = Humanoid | Bear | Cat | Moonkin | Tree
)

// Converts from 0.009327 to 0.0085
const AnimalSpiritRegenSuppression = 0.911337

// Dire Bear Form: 360% more armor from items, and 50% more threat from Bear Form (Passive2) (21178:
// client 70374 raised it from Classic's 30%, as the 8 Oct beta notes say). Moonkin Form carries the
// same armor multiplier.
const BaseBearArmorMulti = 4.6
const BearFormThreatMultiplier = 1.5
const CatFormThreatMultiplier = 0.71
const MoonkinFormArmorMultiplier = 4.6

func (form DruidForm) Matches(other DruidForm) bool {
	return (form & other) != 0
}

func (druid *Druid) InForm(form DruidForm) bool {
	return druid.form.Matches(form)
}

func (druid *Druid) ClearForm(sim *core.Simulation) {
	if druid.InForm(Cat) {
		druid.CatFormAura.Deactivate(sim)
	} else if druid.InForm(Bear) {
		druid.BearFormAura.Deactivate(sim)
	} else if druid.InForm(Moonkin) {
		druid.MoonkinFormAura.Deactivate(sim)
	}

	druid.form = Humanoid
	druid.SetCurrentPowerBar(core.ManaBar)
}

// Forever: "While Shapeshifted, the Druid's melee auto attack DPS while in Bear Form, Cat Form, or
// Dire Bear Form is now the same as the DPS of the Druid's equipped weapon. But the speed of attack
// is changed to 1.0 or 2.5 seconds, depending on the form. Abilities that deal weapon damage
// likewise use these damage values." (Druid class deep dive, 30 September 2026.) The equipped main
// hand's range, weapon damage enchant and bonus DPS are rescaled to the form's swing, keeping its
// DPS and spread. Nothing equipped is the unarmed fist. Ported from MythicSim patch 27 (sage3648).
func (druid *Druid) formWeapon(swingSpeed float64) core.Weapon {
	weapon := druid.WeaponFromMainHand()
	scale := swingSpeed / weapon.SwingSpeed
	return core.Weapon{
		BaseDamageMin:        weapon.BaseDamageMin * scale,
		BaseDamageMax:        weapon.BaseDamageMax * scale,
		SwingSpeed:           swingSpeed,
		NormalizedSwingSpeed: swingSpeed,
		AttackPowerPerDPS:    core.DefaultAttackPowerPerDPS,
		MaxRange:             core.MaxMeleeRange,
	}
}

func (druid *Druid) GetCatWeapon() core.Weapon {
	return druid.formWeapon(1.0)
}

func (druid *Druid) GetBearWeapon() core.Weapon {
	return druid.formWeapon(2.5)
}

// The stats both animal forms grant: Predatory Strikes' attack power off level and Sharpened Claws'
// critical strike chance, both read from the client's rank ladders.
func (druid *Druid) formShiftStats() stats.Stats {
	return stats.Stats{
		stats.AttackPower:         predatoryStrikesAPPerLevel(druid.Talents.PredatoryStrikes) * float64(core.CharacterLevel),
		stats.PhysicalCritPercent: sharpenedClawsCritPercent(druid.Talents.SharpenedClaws),
		stats.SpellCritPercent:    sharpenedClawsCritPercent(druid.Talents.SharpenedClaws),
	}
}

func (druid *Druid) RegisterCatFormAura() {
	actionID := core.ActionID{SpellID: 768}

	statBonus := druid.formShiftStats().Add(stats.Stats{
		stats.AttackPower: 2 * float64(core.CharacterLevel),
	})

	// In Cat Form each point of Agility gives 1 AP, and Feral Attack Power converts 1:1.
	agiApDep := druid.NewDynamicStatDependency(stats.Agility, stats.AttackPower, 1)
	feralApDep := druid.NewDynamicStatDependency(stats.FeralAttackPower, stats.AttackPower, 1)

	// Heart of the Wild: +2% Strength a rank while in Cat Form.
	var hotwDep *stats.StatDependency
	if druid.Talents.HeartOfTheWild > 0 {
		hotwDep = druid.NewDynamicMultiplyStat(stats.Strength, heartOfTheWildFormMultiplier(druid.Talents.HeartOfTheWild))
	}

	druid.CatFormAura = druid.RegisterAura(core.Aura{
		Label:      "Cat Form",
		ActionID:   actionID,
		Duration:   core.NeverExpires,
		BuildPhase: core.Ternary(druid.StartingForm.Matches(Cat), core.CharacterBuildPhaseBase, core.CharacterBuildPhaseNone),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if !druid.Env.MeasuringStats && druid.form != Humanoid {
				druid.ClearForm(sim)
			}
			druid.form = Cat
			druid.SetCurrentPowerBar(core.EnergyBar)

			druid.PseudoStats.ThreatMultiplier *= CatFormThreatMultiplier
			druid.PseudoStats.SpiritRegenMultiplier *= AnimalSpiritRegenSuppression

			druid.AddStatsDynamic(sim, statBonus)
			druid.EnableBuildPhaseStatDep(sim, agiApDep)
			druid.EnableBuildPhaseStatDep(sim, feralApDep)
			if hotwDep != nil {
				druid.EnableBuildPhaseStatDep(sim, hotwDep)
			}

			if !druid.Env.MeasuringStats {
				druid.AutoAttacks.SetMH(druid.GetCatWeapon())
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.UpdateManaRegenRates()
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Humanoid

			druid.PseudoStats.ThreatMultiplier /= CatFormThreatMultiplier
			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression

			druid.AddStatsDynamic(sim, statBonus.Invert())
			druid.DisableBuildPhaseStatDep(sim, agiApDep)
			druid.DisableBuildPhaseStatDep(sim, feralApDep)
			if hotwDep != nil {
				druid.DisableBuildPhaseStatDep(sim, hotwDep)
			}

			if druid.TigersFuryAura != nil {
				druid.TigersFuryAura.Deactivate(sim)
			}

			if !druid.Env.MeasuringStats {
				druid.lastCatFormEnergy = druid.CurrentEnergy()
				druid.lastCatFormExitAt = sim.CurrentTime

				druid.AutoAttacks.SetMH(druid.WeaponFromMainHand())
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.UpdateManaRegenRates()
			}
		},
	})

	druid.CatFormAura.NewPassiveMovementSpeedEffect(0.25)

	// Cat Form (Passive) 3025 in client 70009: Faerie Fire costs nothing and its global cooldown is
	// 0.5 sec shorter in the form.
	druid.attachFormFaerieFireMods(druid.CatFormAura)
	druid.CatFormAura.AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_GlobalCooldown_Flat,
		ClassMask: DruidSpellFaerieFire,
		TimeValue: -500 * time.Millisecond,
	})
}

// Cat and Bear Form (Passive) (3025, 1178, 9635) make Faerie Fire free. They also give it a 6 sec
// cooldown, left out: the debuff lasts 40 sec, so the cooldown never binds.
func (druid *Druid) attachFormFaerieFireMods(formAura *core.Aura) {
	formAura.AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		ClassMask:  DruidSpellFaerieFire,
		FloatValue: -1,
	})
}

func (druid *Druid) registerCatFormSpell() {
	actionID := core.ActionID{SpellID: 768}
	energyMetrics := druid.NewEnergyMetrics(actionID)

	druid.CatForm = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: DruidSpellCatForm,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		ManaCost: spellData.CatForm.Highest().ManaCost(),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if druid.CatFormAura.IsActive() {
				druid.CatFormAura.Deactivate(sim)
			}

			if sim.CurrentTime > 0 {
				target := druid.furorShiftEnergy(sim) + druid.WolfsheadEnergyBonus
				if delta := target - druid.CurrentEnergy(); delta > 0 {
					druid.AddEnergy(sim, delta, energyMetrics)
				} else if delta < 0 {
					druid.SpendEnergy(sim, -delta, energyMetrics)
				}
			}

			druid.CatFormAura.Activate(sim)
		},
	})
}

// Forever reworks powershifting: instead of a chance at a flat 40 Energy, shifting into Cat Form
// carries over a share of the Energy you left the form with, plus a small amount for every second
// spent out of form. The client's curve is 20-100 and the text builds all three from it: that share
// of the Energy, a tenth of it a second, and the whole of it as the cap.
func (druid *Druid) furorShiftEnergy(sim *core.Simulation) float64 {
	if druid.Talents.Furor == 0 {
		return 0
	}

	// Client 17056 eff 1 (EffectAt(2), 20-100; eff 0 is the bear rage chance): "up to a maximum of m2 Energy" caps the whole refund, so 3/5 tops out at 60.
	m2 := spellData.Furor.EffectAt(2).ValueAt(druid.Talents.Furor)
	carryOver := druid.lastCatFormEnergy * m2 / 100
	outOfForm := 0.0
	if druid.lastCatFormExitAt > 0 {
		outOfForm = m2 / 10 * (sim.CurrentTime - druid.lastCatFormExitAt).Seconds()
	}

	return min(m2, carryOver+outOfForm)
}

func (druid *Druid) RegisterBearFormAura() {
	actionID := core.ActionID{SpellID: 9634} // Dire Bear Form
	healthMetrics := druid.NewHealthMetrics(actionID)

	// Dire Bear Form: 180 attack power at level 60 and 1240 health, both flat.
	statBonus := druid.formShiftStats().Add(stats.Stats{
		stats.AttackPower: 3 * float64(core.CharacterLevel),
		stats.Health:      1240,
	})

	feralApDep := druid.NewDynamicStatDependency(stats.FeralAttackPower, stats.AttackPower, 1)

	// Heart of the Wild: +4% Stamina a rank while in Bear Form.
	var hotwDep *stats.StatDependency
	if druid.Talents.HeartOfTheWild > 0 {
		hotwDep = druid.NewDynamicMultiplyStat(stats.Stamina, heartOfTheWildBearStaminaMultiplier(druid.Talents.HeartOfTheWild))
	}

	druid.BearFormAura = druid.RegisterAura(core.Aura{
		Label:      "Bear Form",
		ActionID:   actionID,
		Duration:   core.NeverExpires,
		BuildPhase: core.Ternary(druid.StartingForm.Matches(Bear), core.CharacterBuildPhaseBase, core.CharacterBuildPhaseNone),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if !druid.Env.MeasuringStats && druid.form != Humanoid {
				druid.ClearForm(sim)
			}
			druid.form = Bear
			druid.SetCurrentPowerBar(core.RageBar)

			druid.PseudoStats.ThreatMultiplier *= BearFormThreatMultiplier
			druid.PseudoStats.SpiritRegenMultiplier *= AnimalSpiritRegenSuppression

			druid.AddStatsDynamic(sim, statBonus)
			druid.ApplyDynamicEquipScaling(sim, stats.Armor, BaseBearArmorMulti)
			druid.ApplyDynamicEquipScaling(sim, stats.BonusArmor, BaseBearArmorMulti)
			druid.EnableBuildPhaseStatDep(sim, feralApDep)

			// Preserve the fraction of max health when shifting.
			healthFrac := druid.CurrentHealth() / druid.MaxHealth()
			if hotwDep != nil {
				druid.EnableBuildPhaseStatDep(sim, hotwDep)
			}

			if !druid.Env.MeasuringStats {
				if sim.CurrentTime > 0 {
					druid.restoreHealthFraction(sim, healthFrac, healthMetrics)
				}
				druid.AutoAttacks.SetMH(druid.GetBearWeapon())
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.UpdateManaRegenRates()
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Humanoid

			druid.PseudoStats.ThreatMultiplier /= BearFormThreatMultiplier
			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression

			druid.AddStatsDynamic(sim, statBonus.Invert())
			druid.RemoveDynamicEquipScaling(sim, stats.Armor, BaseBearArmorMulti)
			druid.RemoveDynamicEquipScaling(sim, stats.BonusArmor, BaseBearArmorMulti)
			druid.DisableBuildPhaseStatDep(sim, feralApDep)

			healthFrac := druid.CurrentHealth() / druid.MaxHealth()
			if hotwDep != nil {
				druid.DisableBuildPhaseStatDep(sim, hotwDep)
			}

			if !druid.Env.MeasuringStats {
				if sim.CurrentTime > 0 {
					druid.restoreHealthFraction(sim, healthFrac, healthMetrics)
				}
				if druid.EnrageAura != nil {
					druid.EnrageAura.Deactivate(sim)
				}
				if druid.FrenziedRegenerationAura != nil {
					druid.FrenziedRegenerationAura.Deactivate(sim)
				}
				if druid.maulQueueAura != nil {
					druid.maulQueueAura.Deactivate(sim)
				}
				druid.AutoAttacks.SetMH(druid.WeaponFromMainHand())
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.UpdateManaRegenRates()
			}
		},
	})

	druid.attachFormFaerieFireMods(druid.BearFormAura)
}

func (druid *Druid) registerBearFormSpell() {
	actionID := core.ActionID{SpellID: 9634} // Dire Bear Form
	rageMetrics := druid.NewRageMetrics(actionID)

	druid.BearForm = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: DruidSpellBearForm,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		ManaCost: spellData.DireBearForm.Highest().ManaCost(),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if cur := druid.CurrentRage(); cur > 0 {
				// Shifting into Bear Form resets Rage.
				druid.SpendRage(sim, cur, rageMetrics)
			}

			// Wolfshead Helm gives 5 Rage; the Bear half of Furor is still a chance at 10.
			rageGain := druid.WolfsheadRageBonus
			if sim.Proc(druid.FurorProcChance, "Furor") {
				rageGain += 10
			}
			if rageGain > 0 {
				druid.AddRage(sim, rageGain, rageMetrics)
			}

			druid.BearFormAura.Activate(sim)
		},
	})
}

func (druid *Druid) RegisterMoonkinFormAura() {
	if !druid.Talents.MoonkinForm {
		return
	}

	druid.MoonkinFormAura = druid.RegisterAura(core.Aura{
		Label:      "Moonkin Form",
		ActionID:   core.ActionID{SpellID: 24858},
		Duration:   core.NeverExpires,
		BuildPhase: core.Ternary(druid.StartingForm.Matches(Moonkin), core.CharacterBuildPhaseBase, core.CharacterBuildPhaseNone),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if !druid.Env.MeasuringStats && druid.form != Moonkin {
				druid.ClearForm(sim)
			}

			druid.ApplyDynamicEquipScaling(sim, stats.Armor, MoonkinFormArmorMultiplier)

			druid.form = Moonkin
			druid.SetCurrentPowerBar(core.ManaBar)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.RemoveDynamicEquipScaling(sim, stats.Armor, MoonkinFormArmorMultiplier)
			druid.form = Humanoid
		},
	})
}

func (druid *Druid) RegisterMoonkinFormSpell() {
	if !druid.Talents.MoonkinForm {
		return
	}

	druid.MoonkinForm = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 24858},
		ClassSpellMask: DruidSpellMoonkinForm,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		ManaCost: spellData.MoonkinForm.Highest().ManaCost(),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.MoonkinFormAura.Activate(sim)
		},
	})
}

// Puts health back at the same fraction of max health it was at before a shift changed max health.
// Gains or removes depending on which way max health moved rather than assuming: a second stamina
// multiplier can make leaving Bear Form raise max health, and RemoveHealth panics on a negative.
func (druid *Druid) restoreHealthFraction(sim *core.Simulation, fraction float64, metrics *core.ResourceMetrics) {
	// Max health is rebuilt from stat changes on every shift, so its last bits depend on the
	// iteration's history: a zero delta can come out as +-1e-12 and log a phantom health event,
	// which made single- and multi-threaded runs disagree on the event count.
	delta := fraction*druid.MaxHealth() - druid.CurrentHealth()
	if delta > 1e-6 {
		druid.GainHealth(sim, delta, metrics)
	} else if delta < -1e-6 {
		druid.RemoveHealth(sim, -delta)
	}
}
