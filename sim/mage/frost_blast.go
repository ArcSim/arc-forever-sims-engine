package mage

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// FrostBlastCastable (MythicSim patch 153) registers Frost Blast (1326053, ranks to 1326055), a trained
// Frost spell new in client 70425 (SkillLineAbility 59660-59662, the Frost skill line, every mage): a
// 1.5 sec cast for 90 mana at rank 3 that deals 77 Frost damage (+0.6 a level from 48, 10.5% variance)
// with a 0.271 spell power coefficient and slows the target 40% for 9 sec. The client gave it Frostbolt's
// class mask (0x40180020), the chill bit 0x100000 included, so the mage's Frostbolt modifiers reach it:
// Improved Frostbolt's cast time, Arctic Reach's range, Eureka!'s cost and damage, and every Frost
// school modifier the engine already applies by school (Ice Shards, Piercing Ice, Frost Channeling,
// Elemental Precision). Its slow is a chill, so it rolls Fingers of Frost at launch like Frostbolt.
// Reported missing from the rotation options by Targetdummy (MythicSim Discord, #mage, 17 October 2026).
// False leaves it unregistered, as before.
var FrostBlastCastable = true

// Every rank is registered, as Frostbolt's are: a rotation can pick a cheaper one.
func (mage *Mage) registerFrostBlastSpell() {
	if !FrostBlastCastable {
		return
	}
	spellData.FrostBlast.Each(func(_ int32, rank *spelldata.Spell) { mage.registerFrostBlastRank(rank) })
}

func (mage *Mage) registerFrostBlastRank(rank *spelldata.Spell) {
	damage := rank.DamageEffect()
	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagBinary,
		ClassSpellMask: MageSpellFrostBlast,
		MaxRange:       float64(rank.MaxRange),
		Rank:           rank.RankNumber(),
		MissileSpeed:   float64(rank.Speed),

		ManaCost: core.ManaCostOptions{
			FlatCost: int32(rank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      rank.GCD(),
				CastTime: rank.CastTime(),
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: damage.Coeff(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, damage.Roll(sim, core.CharacterLevel), spell.OutcomeMagicHitAndCrit)
			if mage.procFingersOfFrost != nil {
				mage.procFingersOfFrost(sim, spell, result)
			}
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
