package shaman

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

func (shaman *Shaman) registerShieldsSpells() {
	shaman.registerWaterShieldSpell()
	shaman.registerLightningShieldSpell()
	shaman.registerShieldEffectTriggerSpell()
}

// Only one shield can be up at a time, and either aura may be unregistered (Water Shield is a talent).
func (shaman *Shaman) deactivateShields(sim *core.Simulation) {
	if shaman.LightningShieldAura != nil {
		shaman.LightningShieldAura.Deactivate(sim)
	}
	if shaman.WaterShieldAura != nil {
		shaman.WaterShieldAura.Deactivate(sim)
	}
}

func (shaman *Shaman) registerShieldEffectTriggerSpell() {
	shaman.ShieldSelfProcSpell = shaman.RegisterSpell(core.SpellConfig{
		Flags:          core.SpellFlagNoMetrics | core.SpellFlagNoLogs,
		ClassSpellMask: SpellMaskShieldSelfProc,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, 0, spell.OutcomeAlwaysHit)
		},
	})
}

func (shaman *Shaman) startShieldProcPeriodicAction(sim *core.Simulation) {
	if shaman.SelfBuffs.ShieldProcrate > 0 {
		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period:   60 * time.Second / time.Duration(shaman.SelfBuffs.ShieldProcrate),
			Priority: core.ActionPriorityGCD,
			OnAction: func(sim *core.Simulation) {
				shaman.ShieldSelfProcSpell.Cast(sim, &shaman.Unit)
			},
		})
	}
}

// Water Shield (408510) sits on the Restoration talent line with no rank subtext, so gen_spelldata
// makes no table for it and the ids and values are pinned from the client: three globes of 2% maximum
// mana, one every 3.5 sec at most, no mana cost, 15 sec cooldown.
const waterShieldSpellID = 408510
const waterShieldGlobes = 3
const waterShieldManaFraction = 0.02

// WaterShieldProcsOnSpellCrits (MythicSim patch 145) lets the shaman's own critical spells spend a globe,
// as the 16 October 2026 beta notes say ("Water Shield: Can now trigger from any critical spell cast by
// the Shaman (was only critical heals)"). Client 70425's 408510 states it too: ProcFlags 0x362a8 add
// done magic spells, harmful and helpful, to the hits taken, with ProcHint Crit|Heals. Kerani's
// Elemental sim on the MythicSim Discord kept all three globes for the whole fight. False leaves only
// the hits the shaman takes, as before.
var WaterShieldProcsOnSpellCrits = true

func (shaman *Shaman) registerWaterShieldSpell() {
	if !shaman.Talents.WaterShield {
		return
	}

	actionID := core.ActionID{SpellID: waterShieldSpellID}
	waterShieldManaMetrics := shaman.NewManaMetrics(actionID)
	row := spelldata.MustFind(waterShieldSpellID)

	// One globe every 3.5 sec at most, whichever trigger spends it: the row has one internal cooldown.
	icd := core.Cooldown{Timer: shaman.NewTimer(), Duration: time.Duration(row.ICDMs) * time.Millisecond}
	ready := func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) bool {
		return icd.IsReady(sim)
	}
	spendGlobe := func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
		if !icd.IsReady(sim) || !shaman.WaterShieldAura.IsActive() {
			return
		}
		icd.Use(sim)
		shaman.WaterShieldAura.RemoveStack(sim)
		shaman.AddMana(sim, shaman.MaxMana()*waterShieldManaFraction, waterShieldManaMetrics)
	}

	shaman.WaterShieldAura = shaman.RegisterAura(core.Aura{
		Label:     "Water Shield",
		ActionID:  actionID,
		Duration:  10 * time.Minute,
		MaxStacks: waterShieldGlobes,
	}).AttachProcTrigger(core.ProcTrigger{
		Name:           "Water Shield Trigger",
		Callback:       core.CallbackOnSpellHitTaken,
		ClassSpellMask: SpellMaskShieldSelfProc,
		ExtraCondition: ready,
		Handler:        spendGlobe,
	})

	if WaterShieldProcsOnSpellCrits {
		// "Any critical spell cast by the Shaman": a direct spell hit or heal that crits, the done magic
		// flags of the row. Melee specials (Stormstrike), weapon imbue procs (Flametongue's is a
		// SpellDamageProc), Lightning Shield's orbs and the totems' attacks are not spells the shaman
		// casts, and periodic ticks are not in the row's flags. The row lacks Can Proc From Procs, so
		// Elemental Overload's copies do not spend a globe either. Fire Nova is left out as Elemental
		// Devastation leaves it: a totem's effect, not the shaman's (forever-bugs #334).
		shaman.WaterShieldAura.AttachProcTrigger(core.ProcTrigger{
			Name:             "Water Shield Spell Crit Trigger",
			Callback:         core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
			ProcMask:         core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
			Outcome:          core.OutcomeCrit,
			CanProcFromProcs: row.CanProcFromProcs(),
			ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
				return spell.Unit == &shaman.Unit && !spell.Matches(SpellMaskFireNova) && ready(sim, spell, result)
			},
			Handler: spendGlobe,
		})
	}

	shaman.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMagic,
		Flags:       core.SpellFlagAPL | SpellFlagInstant,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 15,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			shaman.deactivateShields(sim)
			shaman.WaterShieldAura.Activate(sim)
			shaman.WaterShieldAura.SetStacks(sim, waterShieldGlobes)
		},
		RelatedSelfBuff: shaman.WaterShieldAura,
	})
}

var lightningShieldRank = spellData.LightningShield.Highest()

// The shield spell itself carries no damage (its Direct row is the placeholder 1 with a 0 coefficient);
// the orb that fires is a separate spell, and rank 7's is 26363.
var lightningShieldOrb = spellData.LightningShieldTriggered.ByID(26363)

func (shaman *Shaman) registerLightningShieldSpell() {
	actionID := core.ActionID{SpellID: lightningShieldRank.ID}

	lsDamage := shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: lightningShieldOrb.ID},
		SpellSchool:      lightningShieldOrb.SpellSchool(),
		DefenseType:      lightningShieldOrb.DefenseTypeCore(),
		ProcMask:         core.ProcMaskEmpty,
		Flags:            SpellFlagShamanSpell | core.SpellFlagPassiveSpell,
		ClassSpellMask:   SpellMaskLightningShield,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: lightningShieldOrb.DamageEffect().Coeff(),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := lightningShieldOrb.DamageEffect().Average(core.CharacterLevel)
			// Every orb rank carries the client's cannot-crit attribute; beta logs agree (0 crits in
			// ~200 orb hits from 12 shamans, foreverlogs 2650-2680).
			outcome := spell.OutcomeMagicHitAndCrit
			if lightningShieldOrb.CannotCrit() {
				outcome = spell.OutcomeMagicHit
			}
			spell.CalcAndDealDamage(sim, target, baseDamage, outcome)
		},
	})

	shaman.LightningShieldAura = shaman.RegisterAura(core.Aura{
		Label:     "Lightning Shield",
		ActionID:  actionID,
		Duration:  lightningShieldRank.Duration(),
		MaxStacks: int32(lightningShieldRank.ProcCharges),
	}).AttachProcTrigger(core.ProcTrigger{
		Name:           "Lightning Shield Trigger",
		Callback:       core.CallbackOnSpellHitTaken,
		ICD:            3500 * time.Millisecond,
		ClassSpellMask: SpellMaskShieldSelfProc,
		Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			shaman.LightningShieldAura.RemoveStack(sim)
			lsDamage.Cast(sim, shaman.CurrentTarget)
		},
	})

	shaman.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMagic,
		Flags:       core.SpellFlagAPL | SpellFlagInstant,
		ManaCost: core.ManaCostOptions{
			FlatCost: int32(lightningShieldRank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: lightningShieldRank.GCD(),
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			shaman.deactivateShields(sim)
			shaman.LightningShieldAura.Activate(sim)
			shaman.LightningShieldAura.SetStacks(sim, int32(lightningShieldRank.ProcCharges))
		},
		RelatedSelfBuff: shaman.LightningShieldAura,
	})
}
