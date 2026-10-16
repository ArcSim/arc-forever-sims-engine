package priest

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// Shadow Word: Death is new in Forever: trained on the Shadow Magic line, four ranks, an instant
// shadow nuke on an 18 second cooldown (15 before the 16 October 2026 beta notes). The generated table already carries the level-60 damage of
// every rank, which matches the beta client 1.60.1.69893 our sim reads.
//
// Not modelled, deliberately:
//   - The backlash (effect 2, 10% of max health when the target survives). It costs health, not
//     damage, and nothing in a DPS sim reads the priest's health.
//   - Effect 1, a script effect worth 150 that no tooltip names. It could be an execute bonus;
//     until a combat log says what it does, guessing would put an invented number in the DPS.
var ShadowWordDeathRankMap = spellData.ShadowWordDeath

// ShadowWordDeathCooldown is the 16 October 2026 beta notes' cooldown (was 15 sec).
const ShadowWordDeathCooldown = 18 * time.Second

func (priest *Priest) registerShadowWordDeathSpell(rank *spelldata.Spell, cdTimer *core.Timer) {
	// Early Demise (1310076): +15% crit per point against a target at or below 20% health.
	earlyDemise := spellData.EarlyDemise.EffectAt(1).ValueAt(priest.Talents.EarlyDemise)

	priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: PriestSpellShadowWordDeath,
		Rank:           rank.RankNumber(),
		MaxRange:       float64(rank.MaxRange),

		ManaCost: core.ManaCostOptions{
			FlatCost: int32(rank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD(),
			},
			// The 16 October 2026 beta notes raise the cooldown to 18 sec; the client row still states 15.
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: ShadowWordDeathCooldown,
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: rank.DamageEffect().Coeff(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			bonus := 0.0
			if earlyDemise > 0 && sim.IsExecutePhase20() {
				bonus = earlyDemise
			}
			spell.BonusCritPercent += bonus
			spell.CalcAndDealDamage(sim, target, rank.DamageEffect().Roll(sim, core.CharacterLevel), spell.OutcomeMagicHitAndCrit)
			spell.BonusCritPercent -= bonus
		},

		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			return spell.CalcDamage(sim, target, rank.DamageEffect().Average(core.CharacterLevel), spell.OutcomeExpectedMagicHitAndCrit)
		},
	})
}
