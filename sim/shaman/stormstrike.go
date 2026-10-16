package shaman

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
)

var stormstrikeRank = spellData.Stormstrike.Highest()
var StormstrikeActionID = core.ActionID{SpellID: stormstrikeRank.ID}

func (shaman *Shaman) StormstrikeDebuffAura(target *core.Unit) *core.Aura {
	aura := target.GetOrRegisterAura(core.Aura{
		Label:     "Stormstrike-" + shaman.Label,
		ActionID:  StormstrikeActionID,
		Duration:  stormstrikeRank.Duration(),
		MaxStacks: int32(stormstrikeRank.ProcCharges),
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Unit != &shaman.Unit || !spell.Matches(stormstrikeSpells) {
				return
			}
			if !result.Landed() || result.Damage == 0 {
				return
			}
			aura.RemoveStack(sim)
		},
	})
	// Client 17364 (aura 271): only this shaman's Lightning Bolt, Chain Lightning and Earth Shock take the bonus.
	multiplier := 1 + stormstrikeRank.Effect(dbcenums.A_MOD_SPELL_DAMAGE_FROM_CASTER, 0).Percent()
	return aura.AttachDDBC(0, 1, &shaman.AttackTables, func(_ *core.Simulation, spell *core.Spell, _ *core.AttackTable) float64 {
		return core.Ternary(spell.Matches(stormstrikeSpells), multiplier, 1)
	})
}

const stormstrikeSpells = SpellMaskLightningBolt | SpellMaskChainLightning | SpellMaskEarthShock | SpellMaskOverload

func (shaman *Shaman) newStormstrikeHitSpellConfig(spellID int32, isMH bool) core.SpellConfig {
	var procMask core.ProcMask
	var actionTag int32

	procMask = core.Ternary(isMH, core.ProcMaskMeleeMHSpecial, core.ProcMaskMeleeOHSpecial)
	actionTag = core.TernaryInt32(isMH, 1, 2)

	return core.SpellConfig{
		ActionID:         core.ActionID{SpellID: spellID}.WithTag(actionTag),
		SpellSchool:      core.SpellSchoolPhysical,
		DefenseType:      core.DefenseTypeMelee,
		ProcMask:         procMask,
		Flags:            core.SpellFlagMeleeMetrics,
		ClassSpellMask:   SpellMaskStormstrikeDamage,
		ThreatMultiplier: 1,
		DamageMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			weaponDamage := core.Ternary(isMH, spell.Unit.MHWeaponDamage, spell.Unit.OHWeaponDamage)
			baseDamage := weaponDamage(sim, spell.MeleeAttackPower(target))
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialBlockAndCrit)
		},
	}
}

func (shaman *Shaman) newStormstrikeHitSpell(isMH bool) *core.Spell {
	return shaman.RegisterSpell(shaman.newStormstrikeHitSpellConfig(stormstrikeRank.ID, isMH))
}

func (shaman *Shaman) newStormstrikeSpellConfig(spellID int32, ssDebuffAuras *core.AuraArray, mhHit *core.Spell, ohHit *core.Spell) core.SpellConfig {
	stormstrikeSpellConfig := core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellID},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskStormstrikeCast,
		// 8% of base mana, not the client's flat 125 (16 October 2026 beta notes).
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: max(stormstrikeRank.Cooldown(), stormstrikeRank.CategoryCooldown()),
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			shaman.StormstrikeCastResult = spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHitNoHitCounter)
			if shaman.StormstrikeCastResult.Landed() {
				ssDebuffAura := ssDebuffAuras.Get(target)
				ssDebuffAura.Activate(sim)
				ssDebuffAura.SetStacks(sim, ssDebuffAura.MaxStacks)

				if shaman.HasMHWeapon() {
					mhHit.Cast(sim, target)
				}

			}
			spell.DisposeResult(shaman.StormstrikeCastResult)
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return (shaman.HasMHWeapon() || shaman.HasOHWeapon())
		},
	}
	return stormstrikeSpellConfig
}

func (shaman *Shaman) registerStormstrikeSpell() {
	mhHit := shaman.newStormstrikeHitSpell(true)
	ohHit := shaman.newStormstrikeHitSpell(false)

	shaman.StormStrikeDebuffAuras = shaman.NewEnemyAuraArray(shaman.StormstrikeDebuffAura)

	shaman.Stormstrike = shaman.RegisterSpell(shaman.newStormstrikeSpellConfig(stormstrikeRank.ID, &shaman.StormStrikeDebuffAuras, mhHit, ohHit))
}
