package mage

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/stats"
)

func (mage *Mage) registerFrostTalents() {
	// Tier 1
	mage.registerFrostWarding()
	mage.registerImprovedFrostbolt()
	mage.registerElementalPrecision()

	// Tier 2
	mage.registerIceShards()
	mage.registerPermafrost()
	mage.registerImprovedFrostNova()
	mage.registerFrostbite()

	// Tier 3
	mage.registerPiercingIce()
	mage.registerFrostChanneling()
	// Ice Lance: ice_lance.go
	// Improved Blizzard: blizzard.go

	// Tier 4
	mage.registerArcticReach()
	mage.registerIceBlock()
	// Shatter: with Fingers of Frost below

	// Tier 5
	mage.registerImprovedConeOfCold()
	// Cold Snap: cold_snap.go
	mage.registerFingersOfFrost()

	// Tier 6
	mage.registerWinterChill()

	// Tier 7
	mage.registerIceBarrier()
}

// registerFrostWarding implements Frost Warding, new in Forever.
//
// Not modelled: 11189 raises Frost/Ice Armor's armor and resistance and lets Frost Ward reflect
// Frost spells; neither changes damage dealt.
func (mage *Mage) registerFrostWarding() {
	if mage.Talents.FrostWarding == 0 {
		return
	}
}

func (mage *Mage) registerImprovedFrostbolt() {
	if mage.Talents.ImprovedFrostbolt == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask: MageSpellFrostbolt,
		TimeValue: time.Millisecond * time.Duration(spellData.ImprovedFrostbolt.Effect(dbcenums.A_ADD_FLAT_MODIFIER, int32(dbcenums.SPELLMOD_CASTING_TIME)).ValueAt(mage.Talents.ImprovedFrostbolt)),
		Kind:      core.SpellMod_CastTime_Flat,
	})
}

func (mage *Mage) registerElementalPrecision() {
	if mage.Talents.ElementalPrecision == 0 {
		return
	}

	hit := spellData.ElementalPrecision.ValueAt(mage.Talents.ElementalPrecision)
	mage.PseudoStats.SchoolBonusHitChance[stats.SchoolIndexFire] += hit
	mage.PseudoStats.SchoolBonusHitChance[stats.SchoolIndexFrost] += hit
}

func (mage *Mage) registerIceShards() {
	if mage.Talents.IceShards == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll,
		School:     core.SpellSchoolFrost,
		FloatValue: spellData.IceShards.FractionAt(mage.Talents.IceShards),
		Kind:       core.SpellMod_CritMultiplier_Flat,
	})
}

// registerPermafrost implements Permafrost, new in Forever.
//
// Not modelled: 11175 lengthens Chill effects 11/22/33% and slows 3/7/10% more; neither changes
// damage in the sim.
func (mage *Mage) registerPermafrost() {
	if mage.Talents.Permafrost == 0 {
		return
	}
}

// Improved Frost Nova (11165): Frost Nova's cooldown -2/-4 s.
func (mage *Mage) registerImprovedFrostNova() {
	if mage.Talents.ImprovedFrostNova == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask: MageSpellFrostNova,
		TimeValue: time.Millisecond * time.Duration(spellData.ImprovedFrostNova.Effect(dbcenums.A_ADD_FLAT_MODIFIER, int32(dbcenums.SPELLMOD_COOLDOWN)).ValueAt(mage.Talents.ImprovedFrostNova)),
		Kind:      core.SpellMod_Cooldown_Flat,
	})
}

// registerFrostbite implements Frostbite, new in Forever.
//
// Not modelled: 11071 gives Chill effects a 5/10/15% chance to freeze (12494); bosses are immune.
func (mage *Mage) registerFrostbite() {
	if mage.Talents.Frostbite == 0 {
		return
	}
}

// registerPiercingIce: effect 0 (SPELLMOD_DAMAGE, 2/4/6%) raises Frost spell hits, Blizzard's too (its
// ticks are direct hits of 1279949). Effect 1 (SPELLMOD_DOT) names Blizzard and Frostfire Bolt but has
// no rank curve in client 1.60.1.70170, so Frostfire Bolt's DoT gets its base 2% at every rank.
func (mage *Mage) registerPiercingIce() {
	if mage.Talents.PiercingIce == 0 {
		return
	}

	hit := spellData.PiercingIce.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DAMAGE)).FractionAt(mage.Talents.PiercingIce)
	dot := spellData.PiercingIce.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DOT)).FractionAt(mage.Talents.PiercingIce)
	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll &^ MageSpellFrostfireBolt,
		School:     core.SpellSchoolFrost,
		FloatValue: hit,
		Kind:       core.SpellMod_DamageDone_Flat,
	})
	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellFrostfireBolt,
		FloatValue: dot,
		Kind:       core.SpellMod_DamageDone_Flat,
	})
	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellFrostfireBolt,
		FloatValue: hit - dot,
		Kind:       core.SpellMod_DirectDamageDone_Flat,
	})
}

func (mage *Mage) registerFrostChanneling() {
	if mage.Talents.FrostChanneling == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll,
		School:     core.SpellSchoolFrost,
		FloatValue: spellData.FrostChanneling.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_COST)).FractionAt(mage.Talents.FrostChanneling),
		Kind:       core.SpellMod_PowerCost_Pct_Add,
	})
	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll,
		School:     core.SpellSchoolFrost,
		FloatValue: -spellData.FrostChanneling.Effect(dbcenums.A_MOD_THREAT, 16).FractionAt(mage.Talents.FrostChanneling),
		Kind:       core.SpellMod_ThreatMultiplier_Pct,
	})
}

// registerArcticReach implements Arctic Reach, new in Forever: 16757 adds 10/20% to Frostbolt and
// Blizzard range and to Frost Nova and Cone of Cold radius (its two class masks). The sim has no
// percent range mod, so the share is taken of each spell's own range here.
func (mage *Mage) registerArcticReach() {
	if mage.Talents.ArcticReach == 0 {
		return
	}

	rangeShare := spellData.ArcticReach.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_RANGE)).FractionAt(mage.Talents.ArcticReach)
	radiusShare := spellData.ArcticReach.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_RADIUS)).FractionAt(mage.Talents.ArcticReach)
	for mask, yards := range map[int64]float64{
		MageSpellFrostbolt:                       rangeShare * float64(spellData.Frostbolt.Highest().MaxRange),
		MageSpellBlizzard:                        rangeShare * float64(spellData.Blizzard.Highest().MaxRange),
		MageSpellFrostNova | MageSpellConeOfCold: radiusShare * pointBlankRadius,
	} {
		mage.AddStaticMod(core.SpellModConfig{ClassMask: mask, FloatValue: yards, Kind: core.SpellMod_Range_Flat})
	}
}

// registerIceBlock implements Ice Block, new in Forever.
//
// Not modelled: 11958 is a defensive immunity during which the mage cannot act.
func (mage *Mage) registerIceBlock() {
	if !mage.Talents.IceBlock {
		return
	}
}

// Improved Cone of Cold (11190): Cone of Cold damage +12/23/35%.
func (mage *Mage) registerImprovedConeOfCold() {
	if mage.Talents.ImprovedConeOfCold == 0 {
		return
	}

	mage.AddStaticMod(core.SpellModConfig{
		ClassMask:  MageSpellConeOfCold,
		FloatValue: spellData.ImprovedConeOfCold.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DAMAGE)).FractionAt(mage.Talents.ImprovedConeOfCold),
		Kind:       core.SpellMod_DamageDone_Flat,
	})
}

// Raid bosses cannot be chilled or frozen, so Fingers of Frost is the only thing that gets Shatter
// and the Ice Lance bonus going on one; Shatter is folded in here because the two only ever fire
// together. A chill effect has a 15% chance (the talent row's second effect, flat at both ranks) to
// treat the next spells, one per point, as if the target were frozen.
func (mage *Mage) registerFingersOfFrost() {
	if mage.Talents.FingersOfFrost == 0 {
		return
	}

	procChance := spellData.FingersOfFrost.EffectAt(2).FractionAt(mage.Talents.FingersOfFrost)
	fofRank := spellData.FingersOfFrostTriggered.Highest()
	shatterCrit := spellData.Shatter.ValueAt(mage.Talents.Shatter)

	shatterMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll,
		FloatValue: shatterCrit,
		Kind:       core.SpellMod_BonusCrit_Percent,
	})

	// The beta's logs (foreverlogs 2748, 5 procs off Frostbolt, 4 mages) gain the aura at Frostbolt's cast
	// end, before it lands (0.6-0.9 sec later): the next cast takes the bonus and a charge, the bolt's
	// own impact neither. So Frostbolt and Frostfire Bolt roll it at launch (procFingersOfFrost, on
	// the hit already rolled there); Cone of Cold and Improved Blizzard's ticks land at once.
	var launcher *core.Spell
	mage.FingersOfFrostAura = mage.RegisterAura(core.Aura{
		Label:     "Fingers of Frost",
		ActionID:  core.ActionID{SpellID: fofRank.ID},
		Duration:  fofRank.Duration(),
		MaxStacks: int32(spellData.FingersOfFrost.EffectAt(1).ValueAt(mage.Talents.FingersOfFrost)),
		OnGain: func(_ *core.Aura, _ *core.Simulation) {
			shatterMod.Activate()
		},
		OnExpire: func(_ *core.Aura, _ *core.Simulation) {
			shatterMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Matches(MageSpellsAllDamaging) {
				return
			}

			if spell == launcher {
				launcher = nil
				return
			}

			// OnCastComplete runs after the damage is rolled, so the consuming cast keeps the bonus.
			aura.RemoveStack(sim)
		},
	})

	grant := func(sim *core.Simulation) {
		mage.FingersOfFrostAura.Activate(sim)
		mage.FingersOfFrostAura.SetStacks(sim, mage.FingersOfFrostAura.MaxStacks)
	}

	mage.procFingersOfFrost = func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if result.Landed() && sim.Proc(procChance, "Fingers of Frost") {
			launcher = spell
			grant(sim)
		}
	}

	mage.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Fingers of Frost Trigger",
		Callback:           core.CallbackOnSpellHitDealt,
		ClassSpellMask:     MageSpellConeOfCold | MageSpellImprovedBlizzard,
		Outcome:            core.OutcomeLanded,
		ProcChance:         procChance,
		TriggerImmediately: true,
		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			grant(sim)
		},
	})
}

// IsTargetFrozen reports whether the mage's next spell is treated as hitting a frozen target.
func (mage *Mage) IsTargetFrozen() bool {
	return mage.FingersOfFrostAura != nil && mage.FingersOfFrostAura.IsActive()
}

// In Forever Winter's Chill only helps the mage's own Frostbolt and Ice Lance, where Classic and TBC
// made it a raid debuff: 2% crit a stack, one stack per talent point.
func (mage *Mage) registerWinterChill() {
	if mage.Talents.WintersChill == 0 {
		return
	}

	chillRank := spellData.WintersChillTriggered.Highest()
	critPerStack := chillRank.EffectN(1).Average(core.CharacterLevel)

	critMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask: MageSpellFrostbolt | MageSpellIceLance,
		Kind:      core.SpellMod_BonusCrit_Percent,
	})

	mage.WintersChillAura = mage.RegisterAura(core.Aura{
		Label:     "Winter's Chill",
		ActionID:  core.ActionID{SpellID: chillRank.ID},
		Duration:  chillRank.Duration(),
		MaxStacks: int32(spellData.WintersChill.EffectAt(1).ValueAt(mage.Talents.WintersChill)),
		OnGain: func(_ *core.Aura, _ *core.Simulation) {
			critMod.Activate()
		},
		OnExpire: func(_ *core.Aura, _ *core.Simulation) {
			critMod.Deactivate()
		},
		OnStacksChange: func(_ *core.Aura, _ *core.Simulation, _ int32, newStacks int32) {
			critMod.UpdateFloatValue(critPerStack * float64(newStacks))
		},
	})

	// Forever states a flat SpellAuraOptions.ProcChance of 100 on the talent spell and puts the
	// real per-rank chance on effect 2 (by position); effect 1 is the stack count.
	mage.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Winters Chill Talent",
		Callback:           core.CallbackOnSpellHitDealt,
		ProcMask:           core.ProcMaskSpellDamage,
		Outcome:            core.OutcomeLanded,
		ProcChance:         spellData.WintersChill.EffectAt(2).FractionAt(mage.Talents.WintersChill),
		TriggerImmediately: true,
		ExtraCondition: func(_ *core.Simulation, spell *core.Spell, _ *core.SpellResult) bool {
			return spell.SpellSchool.Matches(core.SpellSchoolFrost)
		},
		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			mage.WintersChillAura.Activate(sim)
			mage.WintersChillAura.AddStack(sim)
		},
	})
}

// registerIceBarrier implements Ice Barrier, new in Forever.
//
// Not modelled: 11426 is an absorb shield (rank 1 431); the sim models no damage taken for a mage.
func (mage *Mage) registerIceBarrier() {
	if !mage.Talents.IceBarrier {
		return
	}
}
