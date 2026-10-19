package hunter

import (
	"github.com/wowsims/forever/sim/core"
)

// Classic's spell power coefficients. The beta client (70425, 70450) gives Arcane Shot none (0) and
// Bae's in-game tests agree, so the shot takes none (HunterAbilitiesIgnoreSpellPower, patch 173);
// these apply only with that switch off.
var arcaneShotCoefficients = [9]float64{0, .204, .3, .429, .429, .429, .429, .429, .429}

// Nor any attack power share, but the beta's combat logs show one: four level 20 hunters with known
// agility and gear (foreverlogs.gg, 2026-09-26) hit for base + 0.10-0.12 of ranged attack power.
// An aspect swap settles that it scales: one hunter (report 2673) lost 0.10-0.12 a point of the 35
// ranged attack power Aspect of the Hawk rank 2 gives while Cheetah was up instead.
// Rechecked at levels 24-30 (2026-10-06): Hunter's Mark rank 2 (+59) moves Miss (report 2712,
// 413 hits) 0.107 a point, Lokronn (2695) and Consumer (2717) 0.12, and Miss's own combatant info
// (214 ranged attack power by our formula, +35 Hawk) puts her at 0.112 absolute. That fit gave 0.11.
// Three players in the Hunter Discord then ran dedicated tests (relayed by Bae, 19 October 2026) and
// landed at 0.10, 0.10 and 0.11; with the logs' 0.10-0.12, 0.10 fits all seven best, so the shot
// takes 0.10 (ArcaneShotRAPFromHunterTests, patch 174).
// ponytail: fitted at levels 20-30; refit when a level 60 log shows otherwise.
const (
	arcaneShotRAPCoefficient        = 0.10
	arcaneShotRAPCoefficientLogsFit = 0.11
)

// ArcaneShotRAPFromHunterTests gives Arcane Shot 0.10 of ranged attack power (mythicsim patch 174),
// the Hunter Discord's three dedicated tests (0.10, 0.10, 0.11; via Bae, 19 October 2026) taken with
// the earlier fit from four beta logs (0.10-0.12). False restores the logs-only fit, 0.11.
var ArcaneShotRAPFromHunterTests = true

// arcaneShotRAPShare is the share of ranged attack power Arcane Shot adds to its base damage.
func arcaneShotRAPShare() float64 {
	if ArcaneShotRAPFromHunterTests {
		return arcaneShotRAPCoefficient
	}
	return arcaneShotRAPCoefficientLogsFit
}

func (hunter *Hunter) registerArcaneShotSpell(timer *core.Timer) {
	rank := spellData.ArcaneShot.Highest()
	baseDamage := rank.DamageEffect().Average(core.CharacterLevel)
	rapShare := arcaneShotRAPShare()

	hunter.ArcaneShot = hunter.RegisterRangedSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ClassSpellMask: HunterSpellArcaneShot,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MissileSpeed:   float64(rank.Speed),

		ManaCost: core.ManaCostOptions{
			FlatCost: int32(rank.Cost()),
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    timer,
				Duration: max(rank.Cooldown(), rank.CategoryCooldown()),
			},
		},

		BonusCoefficient: hunterSpellPowerCoefficient(arcaneShotCoefficients[rank.RankNumber()]),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, baseDamage+rapShare*spell.RangedAttackPower(target), spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
