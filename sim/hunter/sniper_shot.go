package hunter

import (
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
)

func (hunter *Hunter) registerSniperShotSpell() {
	if !hunter.Talents.SniperShot {
		return
	}

	spellId := [4]int32{0, 1310687, 0, 1310786}
	base := [4]float64{0, 160, 225, 295}
	mana := [4]float64{0, 365, 365, 365}
	level := [4]int{0, 40, 48, 58}

	best := 0
	for rank := 1; rank <= 3; rank++ {
		if spellId[rank] != 0 && level[rank] <= int(hunter.Level) {
			best = rank
		}
	}
	if best == 0 {
		return
	}

	bonus := base[best]
	hunter.SniperShot = hunter.GetOrRegisterSpell(core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId[best]},
		SpellSchool:   core.SpellSchoolPhysical,
		DefenseType:   core.DefenseTypeRanged,
		ProcMask:      core.ProcMaskRangedSpecial,
		Flags:         core.SpellFlagMeleeMetrics | core.SpellFlagAPL | SpellFlagShot,
		CastType:      proto.CastType_CastTypeRanged,
		Rank:          best,
		RequiredLevel: level[best],
		MissileSpeed:  24,
		ManaCost:      core.ManaCostOptions{FlatCost: mana[best]},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 4000,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 15,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				mult := hunter.PseudoStats.RangedSpeedMultiplier
				if hunter.quiverBonus > 1 {
					mult /= hunter.quiverBonus
				}
				if mult < 0.01 {
					mult = 1
				}
				cast.CastTime = time.Duration(float64(time.Millisecond*4000) / mult)
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.DistanceFromTarget >= core.MinRangedAttackDistance
		},
		CritDamageBonus:  hunter.mortalShots(),
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dmg := hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower(target, false)) +
				hunter.AmmoDamageBonus + bonus
			result := spell.CalcDamage(sim, target, dmg, spell.OutcomeRangedHitAndCrit)
			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.DealDamage(s, result)
			})
		},
	})
}
