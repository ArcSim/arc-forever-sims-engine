package druid

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// The spells client 16870's cost modifier names, as the sim knows them. Wrath, Faerie Fire and the
// forms are not among them.
const clearcastingSpells = DruidSpellClaw | DruidSpellEntanglingRoots | DruidSpellDemoralizingRoar | DruidSpellHurricane |
	DruidSpellFerociousBite | DruidSpellInsectSwarm | DruidSpellLacerate | DruidSpellPrimalBite | DruidSpellMaul |
	DruidSpellMoonfire | DruidSpellRake | DruidSpellRavage | DruidSpellRip | DruidSpellShred | DruidSpellStarfire |
	DruidSpellSwipe | DruidSpellThorns | DruidSpellHealingTouch | DruidSpellRegrowth | DruidSpellLifebloom |
	DruidSpellRejuvenation | DruidSpellTranquility | DruidSpellSwiftmend

// The client states no rate for 16864 (its chance column of 100 is the "no roll here" convention).
// Beta logs fit 2 procs a minute: 87 procs off ~1,500 landed hits of five level 20 druids outside the
// 10 s cooldown, 8% a hit in Bear Form (2.5 s swing) and 4% in Cat Form (1.0 s), the same for
// white hits and specials (foreverlogs reports 32, 2668, 2674).
const omenOfClarityPPM = 2.0

// Omen of Clarity (16864), a baseline Balance passive in Forever: melee hits and spells can grant
// Clearcasting (16870), making the next ability that costs something and is in its mask free.
// Moonkin Form (24858) doubles the chance and halves the 10 s cooldown.
//
// Melee hits, white or special, take the rate off the current swing (the paw's in a form), not the
// equipped weapon: Maul, Swipe and Claw from druids holding 2.4 to 3.6 s weapons gave 36 procs off
// 542 hits, where the paw's speed predicts 36 and the weapon's 59.
// Spells take the rate off their cast time, instants off a 1.5 s GCD; a cast under 1.5 s is NOT floored
// to the GCD. Saries (level 20 resto, reports 33-35/44/45/2667/2676/2677), casts outside the cooldown:
// instants 69 procs / 1,644 (4.2%, 5% predicted), 1.0 s Healing Touch 8 / 409 (2.0%; 3.3% predicted, the old
// GCD floor gave 5%, 20 expected), 2.0 s casts 26 / 388 (6.7%), 2.5-3.5 s 8 / 57 (14%, 8-12% predicted).
func (druid *Druid) applyOmenOfClarity() {
	clearcasting := spellData.OmenOfClarityTriggered.Highest()

	druid.ClearcastingAura = druid.RegisterAura(core.Aura{
		Label:    "Clearcasting",
		ActionID: core.ActionID{SpellID: clearcasting.ID},
		Duration: clearcasting.Duration(),
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// "Not consumed by Wrath or by spells or abilities that cost no resources."
			if spell.Matches(clearcastingSpells) && spell.Cost != nil && spell.Cost.BaseCost > 0 {
				aura.Deactivate(sim)
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct_Add,
		ClassMask:  clearcastingSpells,
		FloatValue: clearcasting.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_COST)).Percent(),
	})

	omen := spellData.OmenOfClarity.Highest()
	moonkin := spellData.MoonkinForm.Highest()
	moonkinChance := 1 + moonkin.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_CHANCE_OF_SUCCESS)).Percent()
	moonkinCooldown := 1 + moonkin.Effect(dbcenums.A_ADD_PCT_MODIFIER, int32(dbcenums.SPELLMOD_PROC_COOLDOWN)).Percent()

	var omenAura *core.Aura
	trigger := spelldata.ProcTrigger(&druid.Character, omen, func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
		druid.ClearcastingAura.Activate(sim)
	}, spelldata.Chance(1))

	trigger.ExtraCondition = func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) bool {
		// "Omen of Clarity can no longer be triggered by Faerie Fire while in Cat, Bear or Dire Bear
		// Form" (16 October 2026 beta notes). The sim's Bear Form is Dire Bear Form too.
		if spell.Matches(DruidSpellFaerieFire|DruidSpellFaerieFireFeral) && druid.InForm(Cat|Bear) {
			return false
		}

		var seconds float64
		switch {
		case spell.ProcMask.Matches(core.ProcMaskSpellDamage | core.ProcMaskSpellHealing):
			seconds = core.GCDDefault.Seconds()
			if spell.DefaultCast.CastTime > 0 {
				seconds = spell.DefaultCast.CastTime.Seconds()
			}
		default:
			seconds = druid.AutoAttacks.MH().SwingSpeed
		}

		chance := omenOfClarityPPM * seconds / 60
		omenAura.Icd.Duration = omen.ICD()
		if druid.MoonkinFormAura.IsActive() {
			chance *= moonkinChance
			omenAura.Icd.Duration = time.Duration(float64(omen.ICD()) * moonkinCooldown)
		}

		return sim.Proc(chance, "Omen of Clarity")
	}

	omenAura = druid.MakeProcTriggerAura(trigger)
}
