package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

type Stance uint8

const (
	StanceNone          = 0
	BattleStance Stance = 1 << iota
	DefensiveStance
	BerserkerStance
)

const stanceEffectCategory = "Stance"

var battleStanceRank = spellData.BattleStance.Highest()
var defensiveStanceRank = spellData.DefensiveStance.Highest()
var berserkerStanceRank = spellData.BerserkerStance.Highest()

func (warrior *Warrior) StanceMatches(other Stance) bool {
	return (warrior.Stance & other) != 0
}

func (warrior *Warrior) makeStanceSpell(stance Stance, mask int64, rank *spelldata.Spell, aura *core.Aura, stanceCD *core.Timer) *core.Spell {
	actionID := aura.ActionID
	rageMetrics := warrior.NewRageMetrics(actionID)
	maxRetainedRage := spellData.TacticalMastery.ValueAt(1) + spellData.ImprovedTacticalMastery.ValueAt(warrior.Talents.ImprovedTacticalMastery)

	return warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		DefenseType:    rank.DefenseTypeCore(),
		ClassSpellMask: mask,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    stanceCD,
				Duration: cooldownOf(rank),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.Stance != stance
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if warrior.WarriorInputs.StanceSnapshot {
				// Delayed, so same-GCD casts are affected by the current aura.
				// Alternatively, those casts could just (artificially) happen before the stance change.
				pa := sim.GetConsumedPendingActionFromPool()
				pa.NextActionAt = sim.CurrentTime + 10*time.Millisecond
				pa.OnAction = aura.Activate
				sim.AddPendingAction(pa)
			} else {
				aura.Activate(sim)
			}

			if warrior.CurrentRage() > maxRetainedRage {
				warrior.SpendRage(sim, warrior.CurrentRage()-maxRetainedRage, rageMetrics)
			}

			warrior.Stance = stance
		},

		RelatedSelfBuff: aura,
	})
}

func (warrior *Warrior) registerBattleStanceAura() *core.Aura {
	actionID := core.ActionID{SpellID: battleStanceRank.ID}

	aura := warrior.RegisterAura(core.Aura{
		Label:      "Battle Stance",
		ActionID:   actionID,
		Duration:   core.NeverExpires,
		BuildPhase: core.Ternary(warrior.DefaultStance == proto.WarriorStance_WarriorStanceBattle, core.CharacterBuildPhaseBuffs, core.CharacterBuildPhaseNone),
	}).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.ThreatMultiplier,
		spellData.BattleStancePassive.Effect(dbcenums.A_MOD_THREAT, 127).MultiplierAt(1),
	)

	aura.NewExclusiveEffect(stanceEffectCategory, true, core.ExclusiveEffect{})

	return aura
}

// FlatThreatOnlyInDefensiveStance (MythicSim patch 151) gives the Warrior's abilities their flat
// bonus threat (Sunder Armor's client E_THREAT, Heroic Strike's, Cleave's, Revenge's, Shield Slam's,
// Hamstring's, Slam's and the shouts' engine values) only in Defensive Stance; their damage threat is
// unchanged. Ashreon on the MythicSim Discord (18 October 2026) says the bonus is Defensive Stance
// only in Forever; the client's only sign is the tank marker aura (668) on Defensive Stance 71, Bear
// and Dire Bear Form and Righteous Fury from client 70170. Unverified, so off by default.
var FlatThreatOnlyInDefensiveStance = false

func (warrior *Warrior) registerDefensiveStanceAura() *core.Aura {
	actionID := core.ActionID{SpellID: defensiveStanceRank.ID}

	aura := warrior.RegisterAura(core.Aura{
		Label:      "Defensive Stance",
		ActionID:   actionID,
		Duration:   core.NeverExpires,
		BuildPhase: core.Ternary(warrior.DefaultStance == proto.WarriorStance_WarriorStanceDefensive, core.CharacterBuildPhaseBuffs, core.CharacterBuildPhaseNone),
	}).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.ThreatMultiplier,
		spellData.DefensiveStancePassive.Effect(dbcenums.A_MOD_THREAT, 127).MultiplierAt(1),
	).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.DamageTakenMultiplier,
		spellData.DefensiveStancePassive.Effect(dbcenums.A_MOD_DAMAGE_PERCENT_TAKEN, 127).MultiplierAt(1),
	).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.DamageDealtMultiplier,
		spellData.DefensiveStancePassive.Effect(dbcenums.A_MOD_DAMAGE_PERCENT_DONE, 127).MultiplierAt(1),
	)
	if FlatThreatOnlyInDefensiveStance {
		set := func(inStance bool) {
			warrior.PseudoStats.ClassFlatThreatMultiplier = core.TernaryFloat64(inStance, 1, 0)
		}
		set(false)
		aura.ApplyOnGain(func(_ *core.Aura, _ *core.Simulation) { set(true) })
		aura.ApplyOnExpire(func(_ *core.Aura, _ *core.Simulation) { set(false) })
	}
	if warrior.Talents.Defiance > 0 {
		defiance := spellData.Defiance.Effect(dbcenums.A_MOD_THREAT, 127).MultiplierAt(warrior.Talents.Defiance)
		applied := false
		refresh := func(inStance bool) {
			want := inStance && warrior.PseudoStats.CanBlock
			if want && !applied {
				warrior.PseudoStats.ThreatMultiplier *= defiance
			} else if !want && applied {
				warrior.PseudoStats.ThreatMultiplier /= defiance
			}
			applied = want
		}
		aura.ApplyOnGain(func(_ *core.Aura, _ *core.Simulation) { refresh(true) })
		aura.ApplyOnExpire(func(_ *core.Aura, _ *core.Simulation) { refresh(false) })
		warrior.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotOffHand}, func(_ *core.Simulation, _ proto.ItemSlot) {
			refresh(aura.IsActive())
		})
	}

	aura.NewExclusiveEffect(stanceEffectCategory, true, core.ExclusiveEffect{})

	return aura
}

func (warrior *Warrior) registerBerserkerStanceAura() *core.Aura {
	actionId := core.ActionID{SpellID: berserkerStanceRank.ID}

	aura := warrior.RegisterAura(core.Aura{
		Label:      "Berserker Stance",
		ActionID:   actionId,
		Duration:   core.NeverExpires,
		BuildPhase: core.Ternary(warrior.DefaultStance == proto.WarriorStance_WarriorStanceBerserker, core.CharacterBuildPhaseBuffs, core.CharacterBuildPhaseNone),
	}).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.ThreatMultiplier,
		spellData.BerserkerStancePassive.Effect(dbcenums.A_MOD_THREAT, 127).MultiplierAt(1),
	).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.DamageTakenMultiplier,
		spellData.BerserkerStancePassive.Effect(dbcenums.A_MOD_DAMAGE_PERCENT_TAKEN, 127).MultiplierAt(1),
	).AttachStatBuff(stats.PhysicalCritPercent, spellData.BerserkerStancePassive.Effect(dbcenums.A_MOD_CRIT_PCT, 0).ValueAt(1))

	aura.NewExclusiveEffect(stanceEffectCategory, true, core.ExclusiveEffect{})

	return aura
}

func (warrior *Warrior) registerStances() {
	stanceCD := warrior.NewTimer()
	battleStanceAura := warrior.registerBattleStanceAura()
	defensiveStanceAura := warrior.registerDefensiveStanceAura()
	berserkerStanceAura := warrior.registerBerserkerStanceAura()
	warrior.BattleStance = warrior.makeStanceSpell(BattleStance, SpellMaskBattleStance, battleStanceRank, battleStanceAura, stanceCD)
	warrior.DefensiveStance = warrior.makeStanceSpell(DefensiveStance, SpellMaskDefensiveStance, defensiveStanceRank, defensiveStanceAura, stanceCD)
	warrior.BerserkerStance = warrior.makeStanceSpell(BerserkerStance, SpellMaskBerserkerStance, berserkerStanceRank, berserkerStanceAura, stanceCD)

	switch warrior.DefaultStance {
	case proto.WarriorStance_WarriorStanceBattle:
		core.MakePermanent(warrior.BattleStance.RelatedSelfBuff)
	case proto.WarriorStance_WarriorStanceDefensive:
		core.MakePermanent(warrior.DefensiveStance.RelatedSelfBuff)
	case proto.WarriorStance_WarriorStanceBerserker:
		core.MakePermanent(warrior.BerserkerStance.RelatedSelfBuff)
	}
}
