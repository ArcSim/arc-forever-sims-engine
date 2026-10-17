package warlock

import (
	"github.com/wowsims/forever/sim/core"
)

// Forever's Incinerate is a Destruction talent (412758 and up), 2.5 sec, and hits 25% harder on a
// target carrying Immolate - the second effect on every rank of the client's row.
//
// IncinerateImmolateBonusOnWholeHit (MythicSim patch 155) puts that 25% on the whole hit, spell power
// included, rather than on the base roll alone (upstream's reading). Client 70425's rows (412758,
// 1293812, 1293813) hold the 25 on a SPELL_EFFECT_DUMMY (effect 3) beside the SCHOOL_DAMAGE effect, with
// the tooltip "an additional $m2% damage if the target is afflicted by Immolate": no second damage
// effect or base points carry it, so the server script scales the damage the spell computed, as
// TrinityCore's Cataclysm Incinerate script does for its Immolate bonus (SetHitDamage(GetHitDamage() + ...)). The
// question came from Ambiwar (MythicSim bugs forum, 15 October 2026) and ShirodawaDaisuke and ProvenTea
// (#warlock). False multiplies the base roll only, as before.
var IncinerateImmolateBonusOnWholeHit = true

func (warlock *Warlock) registerIncinerate() {
	if !warlock.Talents.Incinerate {
		return
	}

	rank := spellData.Incinerate.Highest()
	immolateBonus := 1 + rank.EffectN(2).Percent()

	warlock.Incinerate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellIncinerate,
		MaxRange:       float64(rank.MaxRange),
		MissileSpeed:   float64(rank.Speed),

		ManaCost: core.ManaCostOptions{FlatCost: int32(rank.Cost())},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      rank.GCD(),
				CastTime: rank.CastTime(),
			},
		},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,
		BonusCoefficient:         rank.DamageEffect().Coeff(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := rank.DamageEffect().Roll(sim, core.CharacterLevel)
			immolated := warlock.Immolate.Dot(target).IsActive()
			if immolated && !IncinerateImmolateBonusOnWholeHit {
				baseDamage *= immolateBonus
			}

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if immolated && IncinerateImmolateBonusOnWholeHit {
				result.Damage *= immolateBonus
			}
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
