package paladin

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

var HolyStrikeRankMap = spellData.HolyStrike

// Holy Strike
// https://www.wowhead.com/forever/spell=10333
//
// An instant strike that causes 50% weapon damage plus an additional 93 as Holy damage. 10 sec
// cooldown.
//
// The row states the flat part as its normalized-weapon-damage effect (121) and the percentage as a
// weapon-percent-damage effect (31). The client adds the flat amount to the normalized swing and
// then takes the percentage of the sum, the way it does for Backstab. Build 70009's ranks are
// 25/29/32/36/39/43/46/50%, so rank 8 is 50% of (weapon + 81 to 105).
// HolyStrikeSpellPowerCoefficient is Holy Strike's spell power share as measured in game on client
// 70425 (MythicSim patch 143); the rows carry 0.429 on every rank. The base, 50% of the weapon plus the
// flat Holy amount, already matched.
const HolyStrikeSpellPowerCoefficient = 0.286

// HolyStrikeSpellPowerFromRow gives Holy Strike the share its client rows carry, 0.429 on every rank in
// client 70450 and every client before it, instead of patch 143's measured 0.286 (MythicSim patch 162).
// The pre-16 October logs show the server dealing the row's share: Judgement of the Crusader steps of 0.37
// to 0.45 a point (mean 0.41) for five paladins with no Holy Strike talent (reports 2744, 2693, 32, 2755,
// 2782: 0/0/5, 0/0/10, 0/0/11, 0/0/13 and 0/24/0), where 0.286 would need a 1.43 multiplier none of them
// can have. On by default since MythicSim patch 170 (Sage, 18 October 2026: the beta stops at level 30,
// so official data settles it). The 16 October coefficient changes reached the rows (Judgement of
// Righteousness 0.5 to 0.429, Consecration 0.042) and left Holy Strike at 0.429; no note changes it.
// Ashreon's 0.286 is exactly 2/3 of the row, the same ratio as her Judgement of Righteousness, Judgement
// of Fury and Seal of Fury numbers that Lightclub's live sheet put back on the rows (patch 168).
// False gives 0.286.
var HolyStrikeSpellPowerFromRow = true

// HolyStrikeFlatBeforeWeaponPercentIs40To53 reads Ashreon's "Holy Strike 40-53 (was 94-121)" (bugs
// forum, 17 October 2026) as the flat amount the client adds before taking 50% of the sum, so the
// hit's Holy part is half of that, 20 to 26.5 (MythicSim patch 148). Every rank's flat amount is
// halved, as rank 8's 40-53 is half of its row's 81-105. Off by default: the engine's Holy part is
// already 50% of 81-105, 40.7 to 52.3, which is the 40-53 her post gives, so the other reading
// (her number is the hit's Holy part) needs no change. Her "~50% lower in game" fits this switch;
// which reading is right needs a hit with a known weapon and spell power.
var HolyStrikeFlatBeforeWeaponPercentIs40To53 = false

func (paladin *Paladin) registerHolyStrike(_ int32, rank *spelldata.Spell) {
	flat := rank.DamageEffect()
	weaponPercent := rank.EffectN(2).Percent()
	flatMultiplier := core.TernaryFloat64(HolyStrikeFlatBeforeWeaponPercentIs40To53, 0.5, 1)

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskHolyStrike,
		Rank:           rank.RankNumber(),
		MaxRange:       core.MaxMeleeRange,

		ManaCost: rank.ManaCost(),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD(),
			},
			// Client cooldown category 2404 holds Holy Strike and Hammer of the Righteous.
			CD: core.Cooldown{
				Timer:    paladin.sharedTimer(&paladin.holyStrikeTimer),
				Duration: cooldown(rank),
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: core.TernaryFloat64(HolyStrikeSpellPowerFromRow, flat.Coeff(), HolyStrikeSpellPowerCoefficient),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := weaponPercent * (spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower(target)) + flatMultiplier*flat.Roll(sim, core.CharacterLevel))
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}
