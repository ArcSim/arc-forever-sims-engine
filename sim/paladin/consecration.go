package paladin

import (
	"fmt"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
)

var ConsecrationRankMap = spellData.Consecration

// Consecration
// https://www.wowhead.com/forever/spell=20924
//
// Consecrates the land beneath the Paladin, doing X Holy damage over 8 sec to enemies who enter the
// area. The first 4 enemies who enter the area will take an additional Y damage over 8 sec.
//
// The ticks sit on the spell the tooltip names, which the family's triggered ladder carries rank
// for rank: the tick everyone takes on its first effect, and on its second the extra damage the
// first few targets take, the only part of the spell the client gives a spell power coefficient.
// The rank's own periodic dummy carries no damage: it times the ticks and states how many targets
// take the bonus.
// Consecration's per-tick spell power coefficient: 0.042, the 16 October 2026 beta notes' number (was
// 0.095). Client 70425's rows carry 0.042 on the bonus effect as well, so the constant now only pins it.
// The pre-16 October logs measured the old 0.095 (Judgement of the Crusader steps of 0.089 to 0.108 a
// tick on Retribution paladins, MythicSim's 19 October review).
const consecrationTickCoefficient = 0.042

func (paladin *Paladin) registerConsecration(n int32, rank *spelldata.Spell) {
	tickSpell := spellData.ConsecrationTriggered.Rank(n)
	tick := tickSpell.EffectN(1)
	bonus := tickSpell.EffectN(2)

	// The 16 October 2026 beta notes cut the per-tick spell power coefficient to 0.042; client 70425's
	// rows carry the same 0.042 on the bonus effect, the only one with a coefficient.
	bonusCoeff := consecrationTickCoefficient

	dummy := rank.Effect(dbcenums.A_PERIODIC_DUMMY, 0)
	bonusTargets := int(dummy.Average(core.CharacterLevel))
	tickLength := dummy.Period()
	numberOfTicks := int32(rank.Duration() / tickLength)

	// Each tick is its own direct School Damage spell in the client (1280345-1280349, no Can't Crit),
	// so it rolls a spell crit, as on master.
	// The bonus scales on its own coefficient, so it is added to the base damage here rather than
	// through the dot's, which is the base tick's. Consecrated Ground marks the same targets.
	dealTick := func(sim *core.Simulation, dot *core.Dot) {
		for i, target := range sim.Encounter.ActiveTargetUnits {
			damage := tick.Average(core.CharacterLevel)
			if i < bonusTargets {
				damage += bonus.Average(core.CharacterLevel) + bonusCoeff*dot.Spell.BonusDamage(dot.Spell.Unit.AttackTables[target.UnitIndex])
				if paladin.consecratedGroundAuras != nil {
					paladin.consecratedGroundAuras.Get(target).Activate(sim)
				}
			}
			dot.Spell.CalcAndDealPeriodicDamage(sim, target, damage, dot.Spell.OutcomeTickMagicHitAndCrit)
		}
	}

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskConsecration,
		Rank:           rank.RankNumber(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		MaxRange: 8,

		ManaCost: rank.ManaCost(),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: rank.GCD(),
			},
			CD: core.Cooldown{
				Timer:    paladin.sharedTimer(&paladin.consecrationTimer),
				Duration: cooldown(rank),
			},
		},

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				ActionID: core.ActionID{SpellID: rank.ID},
				Label:    fmt.Sprintf("Consecration%s Rank %d", paladin.Label, rank.RankNumber()),
			},
			// The client ticks on a 1 sec period with no tick-on-apply attribute (20924: aura 226,
			// 1000 ms, SpellMisc Attributes[5] 0), so the first tick lands 1 sec in, as on master.
			NumberOfTicks:    numberOfTicks,
			TickLength:       tickLength,
			BonusCoefficient: tick.Coeff(),
			OnTick: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot) {
				dealTick(sim, dot)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Consecration does one hit check on cast but the ground effect will still be applied
			// meaning it's only needed to proc things like Eye of Magtheridon (procs on resist)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)

			spell.AOEDot().Apply(sim)
		},
	})
}
