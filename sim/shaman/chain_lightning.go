package shaman

import (
	"math"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

var ChainLightningRankMap = spellData.ChainLightning

// Lightning Overload's Chain Lightning, one client row per rank (client 1.60.1.70205). Not exactly
// half the rank: its damage grows 1.25-1.5x faster per level, so rank 4 at 60 averages 66, not 63.
var chainLightningOverloadRanks = spelldata.Ranked(408479, 408481, 408482, 408484)

// "Chain Lightning mana cost reduced by approximately 35%" on every rank (16 October 2026 beta notes). The
// client data still carries the old costs, so the rank's cost is cut here until it moves.
const chainLightningCostMultiplier = 0.65

func chainLightningCost(config *spelldata.Spell) int32 {
	return int32(math.Round(config.Cost() * chainLightningCostMultiplier))
}

func (shaman *Shaman) registerChainLightningSpell() {
	maxHits := min(3, shaman.Env.TotalTargetCount())
	sharedCDTimer := shaman.NewTimer()
	shaman.ChainLightningOverloads = make(map[int32][]*core.Spell, ChainLightningRankMap.Len())
	ChainLightningRankMap.Each(func(rank int32, config *spelldata.Spell) {
		shaman.newChainLightningSpell(config, rank, false, sharedCDTimer)
		for range maxHits {
			shaman.ChainLightningOverloads[rank] = append(shaman.ChainLightningOverloads[rank], shaman.newChainLightningSpell(config, rank, true, nil))
		}
	})
}

func (shaman *Shaman) newChainLightningSpell(config *spelldata.Spell, rank int32, isElementalOverload bool, sharedCDTimer *core.Timer) *core.Spell {
	damage := config.DamageEffect()
	if isElementalOverload {
		damage = chainLightningOverloadRanks.Rank(rank).DamageEffect()
	}
	shamConfig := ShamSpellConfig{
		ActionID:            core.ActionID{SpellID: config.ID},
		Rank:                rank,
		IsElementalOverload: isElementalOverload,
		BaseFlatCost:        chainLightningCost(config),
		BonusCoefficient:    damage.Coeff(),
		SpellSchool:         core.SpellSchoolNature,
		Overloads:           shaman.ChainLightningOverloads,
		BounceReduction:     0.7,
		ClassSpellMask:      core.TernaryInt64(isElementalOverload, SpellMaskChainLightningOverload, SpellMaskChainLightning),
		BaseCastTime:        config.CastTime(),
	}
	spellConfig := shaman.newElectricSpellConfig(shamConfig)
	if !isElementalOverload {
		spellConfig.MaxRange = float64(config.MaxRange)
		spellConfig.Cast.CD = core.Cooldown{
			Timer:    sharedCDTimer,
			Duration: max(config.Cooldown(), config.CategoryCooldown()),
		}
	}
	maxHits := int32(3)
	maxHits = min(maxHits, shaman.Env.TotalTargetCount())

	spellConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		curTarget := target
		bounce := shamConfig.BounceReduction + shaman.ChainLightningBounceBonus

		// Damage calculation and DealDamage are in separate loops so that e.g. a spell power proc
		// can't proc on the first target and apply to the second
		numHits := min(maxHits, shaman.Env.ActiveTargetCount())
		results := make([]*core.SpellResult, numHits)
		for hitIndex := range numHits {
			baseDamage := damage.Roll(sim, core.CharacterLevel)
			results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMagicHitAndCrit)

			curTarget = sim.Environment.NextActiveTargetUnit(curTarget)
			spell.DamageMultiplier *= bounce
		}

		for hitIndex := range numHits {
			if !isElementalOverload && results[hitIndex].Landed() && sim.Proc(shaman.GetOverloadChance()/3, "Chain Lightning Elemental Overload") {
				shamConfig.Overloads[rank][hitIndex].Cast(sim, results[hitIndex].Target)
			}
			spell.DealDamage(sim, results[hitIndex])
			spell.DamageMultiplier /= bounce
		}
	}

	return shaman.RegisterSpell(spellConfig)
}
