package priest

import (
	"fmt"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// Penance is the Discipline talent the Smite build goes deep for: three Holy bolts over the channel,
// on a 12 second cooldown the ranks share (Category 2414). Every rank is registered.
//
// The client's rank 3 bolt (1240730, 180 for 270 mana) is larger than rank 4's (1316993, 131 for
// 355); the table is taken as it is, so the Smite rotations cast rank 3. Wowhead Forever prints the
// same. The bolts land "instantly and every 1 sec for 2 sec" (1316995's tooltip): the channel 1316994 lasts
// 2000 ms, fires 1316993 every 1000 ms and ticks on application (attribute 5, 0x200). So one bolt on
// cast and two channel ticks.
const PenanceTicks = 3

func (priest *Priest) registerPenanceSpell(rank *spelldata.Spell) {
	// The cast's tooltip names its damage bolt first, then the heal bolt.
	bolt := rank.Refs()[0]

	priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    core.SpellSchoolHoly,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagChanneled,
		ClassSpellMask: PriestSpellPenance,
		Rank:           rank.RankNumber(),
		MaxRange:       float64(rank.MaxRange),

		ManaCost: core.ManaCostOptions{
			FlatCost: int32(rank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD(),
			},
			CD: core.Cooldown{
				Timer:    priest.CategoryTimer(int32(rank.Category)),
				Duration: max(rank.Cooldown(), rank.CategoryCooldown()),
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("Penance-%d", rank.RankNumber()),
			},
			NumberOfTicks:       PenanceTicks - 1,
			TickLength:          time.Second,
			AffectedByCastSpeed: false,
			BonusCoefficient:    bolt.DamageEffect().Coeff(),

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, bolt.DamageEffect().Average(core.CharacterLevel))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// Each bolt is its own direct Holy hit (1316993, School Damage), so every one can crit.
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, priest.tickOutcome(true, dot))
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				dot := spell.Dot(target)
				dot.Apply(sim)
				dot.TickOnce(sim)
			}
			spell.DealOutcome(sim, result)
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				return spell.Dot(target).CalcSnapshotDamage(sim, target, spell.OutcomeExpectedMagicHit)
			}
			return spell.CalcPeriodicDamage(sim, target, bolt.DamageEffect().Average(core.CharacterLevel), spell.OutcomeExpectedMagicHit)
		},
	})
}
