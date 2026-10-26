package hunter

import (
	"strconv"
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
)

func (hunter *Hunter) improvedStingsMultiplier() float64 {
	return []float64{0, 0.06, 0.13, 0.20}[hunter.Talents.ImprovedStings]
}

// The client data gives Serpent Sting no attack power or spell power share, but the ticks grow with ranged attack
// power: tested in game (rank 4, no Improved Stings), 27 at 174 and 32 at 361 before rounding, so 0.030 a tick; 20
// spell power left them unchanged.
const serpentStingAPPerTick = 0.030

// tick damage per rank from the Forever client data
var serpentStingBaseTick = [10]float64{0, 2, 6, 12, 22, 34, 48, 64, 83, 111}

func serpentStingTick(rank int, rangedAttackPower float64) float64 {
	return serpentStingBaseTick[rank] + rangedAttackPower*serpentStingAPPerTick
}

func (hunter *Hunter) getSerpentStingConfig(rank int) core.SpellConfig {
	spellId := [10]int32{0, 1978, 13549, 13550, 13551, 13552, 13553, 13554, 13555, 25295}[rank]
	manaCost := [10]float64{0, 15, 30, 50, 80, 115, 150, 190, 230, 250}[rank]
	level := [10]int{0, 4, 10, 18, 26, 34, 42, 50, 58, 60}[rank]

	return core.SpellConfig{
		SpellCode:     SpellCode_HunterSerpentSting,
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolNature,
		DefenseType:   core.DefenseTypeRanged,
		ProcMask:      core.ProcMaskRangedSpecial,
		Flags:         core.SpellFlagAPL | core.SpellFlagPureDot | core.SpellFlagPoison | SpellFlagSting,
		CastType:      proto.CastType_CastTypeRanged,
		Rank:          rank,
		RequiredLevel: level,
		MissileSpeed:  40,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceNow(sim) >= core.MinRangedAttackDistance
		},

		DamageMultiplier: 1 + hunter.improvedStingsMultiplier(),
		// Mortal Shots' crit damage modifier covers the sting in the client data, so its crit ticks get it
		CritDamageBonus: hunter.mortalShots(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "SerpentSting" + hunter.Label + strconv.Itoa(rank),
				Tag:   "SerpentSting",
			},
			NumberOfTicks:    5,
			TickLength:       time.Second * 3,
			BonusCoefficient: 0,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, serpentStingTick(rank, dot.Spell.RangedAttackPower(target, false)), isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickPhysicalCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeRangedHitNoHitCounter)

			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.DealOutcome(sim, result)

				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			})
		},
	}
}

func (hunter *Hunter) registerSerpentStingSpell() {
	maxRank := core.TernaryInt(core.IncludeAQ, 9, 8)
	for rank := maxRank; rank >= 0; rank-- {
		config := hunter.getSerpentStingConfig(rank)

		if config.RequiredLevel <= int(hunter.Level) {
			hunter.SerpentSting = hunter.GetOrRegisterSpell(config)
			break
		}
	}
}
