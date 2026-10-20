package warrior

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/buffs"
	"github.com/wowsims/forever/sim/core/dbcenums"
)

// SunderArmorThreatFromAttackPower adds 0.05 of the warrior's attack power to Sunder Armor's threat
// (mythicsim patch 185, ported from wowsims/forever's 2026-10-07 overrides): upstream measured Rank 1
// at 132, 158, 170 and 180 attack power, in Battle and Defensive Stance, as 34 + 0.05 x attack power
// times the stance, the "small increase from attack power" of the 2026-09-24 notes. The other ranks
// carry Rank 1's share until measured. Off, the client's flat value alone.
var SunderArmorThreatFromAttackPower = true

const sunderThreatPerAttackPower = 0.05

func (warrior *Warrior) registerSunderArmor() {
	// The client supplies Sunder Armor's flat threat per rank: 34/75/117/158/206 (build 70009).
	// The 2026-09-24 notes add "a small increase from attack power", which the E_THREAT effect does
	// not carry: patch 185 adds upstream's measured share.
	sunderArmorRank := spellData.SunderArmor.Highest()
	sunderThreat := sunderArmorRank.FindEffect(dbcenums.E_THREAT, 0, 0).Average(core.CharacterLevel)

	actionId := core.ActionID{SpellID: sunderArmorRank.ID}

	// The generated aura is the client's rank 5 (11597, 450 armor a stack), the rank cast here.
	warrior.SunderArmorAuras = warrior.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return buffs.SunderArmorAura(target, true, 0)
	})

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskSunderArmor,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   int32(sunderArmorRank.Cost()),
			Refund: sunderArmorRank.MissRefund(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: sunderArmorRank.GCD(),
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.CanApplySunderAura(target)
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		FlatThreatBonus:  sunderThreat,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if SunderArmorThreatFromAttackPower {
				spell.FlatThreatBonus = sunderThreat + sunderThreatPerAttackPower*spell.MeleeAttackPower(target)
			}
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)

			if result.Landed() {
				aura := warrior.SunderArmorAuras.Get(target)
				aura.Activate(sim)
				aura.AddStack(sim)
			} else {
				spell.IssueRefund(sim)
			}

			spell.DealOutcome(sim, result)
		},

		RelatedAuraArrays: warrior.SunderArmorAuras.ToMap(),
	})
}

func (warrior *Warrior) CanApplySunderAura(target *core.Unit) bool {
	return warrior.SunderArmorAuras.Get(target).IsActive() || !warrior.SunderArmorAuras.Get(target).ExclusiveEffects[0].Category.AnyActive()
}
