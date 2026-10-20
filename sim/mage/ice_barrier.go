package mage

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

// Frigid Power (1326891), the Frost damage buff every Ice Barrier cast applies for a minute. No effect
// of Ice Barrier triggers it; its tooltip names it ("$1326891s1%"), which is what the row's RefIDs carry.
const frigidPowerSpellID = 1326891

// Ice Barrier (MythicSim patch 188, client hotfix push 112753 on build 70450, 20 October 2026):
//
// "Shield yourself, absorbing ${$s1+($SPFR*0.1)*($11189s3/100+1)} damage and increasing your Frost
// damage by $1326891s1% for $d. While shielded, your Frost damage is increased by an additional $s2%
// and you do not lose casting time when you take damage."
//
// The shield's effects 2 and 3 raise Frost damage while it holds, and Frigid Power raises it again for
// the full minute, however soon the shield breaks. Both carry Piercing Ice's two masks, the talent the
// hotfix removed, so both are applied as Piercing Ice was: effect 2 (SPELLMOD_DAMAGE) to every Frost
// spell, Blizzard's ticks included since they are direct hits of 1279949, and effect 3 (SPELLMOD_DOT) to
// Frostfire Bolt, whose hit takes effect 2 and whose DoT takes effect 3. The spell power share of the
// absorb is the tooltip's, since the row states no coefficient, and Frost Warding's effect 3 raises the
// whole absorb. The sim models no pushback, so the casting time clause does nothing. Ported from
// wowsims/forever #594 (12f3f1db0a), which keeps upstream's absorb and its Frigid Power.
//
// The spell is the highest rank, castable from the rotation only (not a major cooldown): a fight in
// which nothing hits the mage never breaks the shield, so a rotation recasts it before Frigid Power
// drops.
func (mage *Mage) registerIceBarrierSpell() {
	rank := spellData.IceBarrier.Highest()
	frigidPowerRank := spellData.IceBarrierTriggered.ByID(frigidPowerSpellID)

	frigidPower := mage.RegisterAura(spelldata.AuraConfig(frigidPowerRank))
	attachFrostDamageBonus(frigidPower,
		frigidPowerRank.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DAMAGE)).Average(core.CharacterLevel)/100,
		frigidPowerRank.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DOT)).Average(core.CharacterLevel)/100)

	warding := 1 + spellData.FrostWarding.EffectAt(3).FractionAt(mage.Talents.FrostWarding)
	base := rank.EffectN(1).Average(core.CharacterLevel)
	var absorb float64

	shield := mage.NewDamageAbsorptionAura(core.AbsorptionAuraConfig{
		Aura: spelldata.AuraConfig(rank),
		ShieldStrengthCalculator: func(_ *core.Unit) float64 {
			return absorb
		},
	})
	attachFrostDamageBonus(shield.Aura,
		rank.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DAMAGE)).Average(core.CharacterLevel)/100,
		rank.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_DOT)).Average(core.CharacterLevel)/100)

	config := spelldata.SpellConfig(&mage.Unit, rank, spelldata.Flags(core.SpellFlagAPL))
	config.ClassSpellMask = MageSpellIceBarrier
	config.RelatedSelfBuff = shield.Aura
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		absorb = warding * (base + 0.1*mage.frostSpellPower())
		shield.Activate(sim)
		frigidPower.Activate(sim)
	}

	mage.RegisterSpell(config)
}

// Piercing Ice's shape (sim/mage/talents_frost.go before patch 187): hit on every Frost spell but
// Frostfire Bolt, dot on all of Frostfire Bolt and the rest of hit on its direct damage alone.
func attachFrostDamageBonus(aura *core.Aura, hit float64, dot float64) {
	aura.AttachSpellMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll &^ MageSpellFrostfireBolt,
		School:     core.SpellSchoolFrost,
		FloatValue: hit,
		Kind:       core.SpellMod_DamageDone_Flat,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask:  MageSpellFrostfireBolt,
		FloatValue: dot,
		Kind:       core.SpellMod_DamageDone_Flat,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask:  MageSpellFrostfireBolt,
		FloatValue: hit - dot,
		Kind:       core.SpellMod_DirectDamageDone_Flat,
	})
}

// The spell power a Frost spell reads ($SPFR in a tooltip): the general bonus plus the Frost one.
func (mage *Mage) frostSpellPower() float64 {
	return mage.GetStat(stats.SpellDamage) + mage.GetStat(stats.FrostDamage)
}
