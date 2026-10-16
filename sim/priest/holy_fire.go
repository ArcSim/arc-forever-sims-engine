package priest

import (
	"fmt"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// The Smite build's opener and the trigger Power in Light and Searing Light both read.
// The client's rank 5 dot (13 a tick) is larger than rank 6's (10 a tick); the table is taken as it is.
var HolyFireRankMap = spellData.HolyFire

func (priest *Priest) registerHolyFireSpell(rank *spelldata.Spell) {
	tick := rank.PeriodicEffect()
	tickLength := tick.Period()

	priest.HolyFire = append(priest.HolyFire, priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: PriestSpellHolyFire,
		Rank:           rank.RankNumber(),
		MaxRange:       float64(rank.MaxRange),

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
		BonusCoefficient: rank.DamageEffect().Coeff(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("HolyFire-%d", rank.RankNumber()),
			},
			NumberOfTicks:       int32(rank.Duration() / tickLength),
			TickLength:          tickLength,
			AffectedByCastSpeed: false,
			BonusCoefficient:    tick.Coeff(),

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, tick.Average(core.CharacterLevel))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, priest.tickOutcome(rank.PeriodicCanCrit(), dot))
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, rank.DamageEffect().Roll(sim, core.CharacterLevel), spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealDamage(sim, result)
		},
	}))
}

// Whether this priest's Holy Fire or Shadow Word: Pain is on the target, which is what Power in Light asks.
func (priest *Priest) hasPowerInLightDot(target *core.Unit) bool {
	for _, spells := range [][]*core.Spell{priest.HolyFire, priest.ShadowWordPain} {
		for _, spell := range spells {
			if spell.Dot(target).IsActive() {
				return true
			}
		}
	}
	return false
}
