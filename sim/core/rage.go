package core

import (
	"fmt"

	"github.com/wowsims/forever/sim/core/proto"
)

const ThreatPerRageGained = 5

// TODO: Ingame test needed at higher levels
const DamageTakenRageFactor = 10
// A critical auto attack adds its swing's rage again, before the two-hander's factor. Blizzard's
// beta development notes (1 and 8 October, warrior and druid): "Additional Rage generated from
// landing Critical Strikes increased to 100% increased Rage (was 75%)". The client hook, Rule of
// Rage (DND) 1322574, is an A_DUMMY of 100 on melee autos. Magey's beta logs
// (github.com/magey/forever-warrior/issues/3) fit crit = swing + 3.46 x speed: a one-hander doubles,
// a two-hander's 3.3 sec swing gets 11.42 on top of 14.85 (measured 11.3-11.5).
const CritRageBonus = 1.0

const (
	BaseRageHitFactor     = 3.46
	TwoHandRageHitFactor  = 4.5
	TwoHandRageMultiplier = TwoHandRageHitFactor / BaseRageHitFactor
)

type rageBar struct {
	unit *Unit

	maxRage      float64
	startingRage float64
	currentRage  float64

	offHandRageMultiplier     float64
	damageTakenRageMultiplier float64

	RageRefundMetrics     *ResourceMetrics
	EncounterStartMetrics *ResourceMetrics
}

type RageBarOptions struct {
	MaxRage            float64
	StartingRage       float64
	BaseRageMultiplier float64
}

func (unit *Unit) EnableRageBar(options RageBarOptions) {
	rageFromDamageTakenMetrics := unit.NewRageMetrics(ActionID{OtherID: proto.OtherAction_OtherActionDamageTaken})

	unit.SetCurrentPowerBar(RageBar)
	unit.RegisterAura(Aura{
		Label:    "RageBar",
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if unit.GetCurrentPowerBar() != RageBar {
				return
			}
			if !result.Landed() {
				return
			}

			hitFactor := BaseRageHitFactor
			handMultiplier := 1.0
			var weapon *Weapon

			switch spell.ProcMask {
			case ProcMaskMeleeMHAuto:
				weapon = unit.AutoAttacks.MH()
			case ProcMaskMeleeOHAuto:
				// OH hits generate 50% of the rage they would if they were MH hits
				hitFactor /= 2
				handMultiplier = unit.rageBar.offHandRageMultiplier
				weapon = unit.AutoAttacks.OH()
			default:
				return
			}

			// rage is normalized so it only depends on weapon swing speed and some multipliers
			baseRage := hitFactor * weapon.SwingSpeed * options.BaseRageMultiplier * handMultiplier
			generatedRage := baseRage
			if weapon.NormalizedSwingSpeed == TwoHandNormalizedSwingSpeed {
				generatedRage *= TwoHandRageMultiplier
			}
			if result.DidCrit() {
				generatedRage += baseRage * CritRageBonus
			}

			var metrics *ResourceMetrics
			if spell.Cost != nil {
				metrics = spell.Cost.ResourceCostImpl.(*RageCost).ResourceMetrics
			} else {
				if spell.ResourceMetrics == nil {
					spell.ResourceMetrics = spell.Unit.NewRageMetrics(spell.ActionID)
				}
				metrics = spell.ResourceMetrics
			}
			unit.AddRage(sim, generatedRage, metrics)
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if unit.GetCurrentPowerBar() != RageBar {
				return
			}
			// "A hit taken" (f9f9f21883): a missed, dodged or parried swing pays nothing, though its
			// PostArmorAndResistanceMultiplier is still filled in.
			if !result.Landed() {
				return
			}

			preArmorDamage := result.PostArmorAndResistanceMultiplier
			if result.ArmorAndResistanceMultiplier > 0 {
				preArmorDamage /= result.ArmorAndResistanceMultiplier
			}
			generatedRage := preArmorDamage * DamageTakenRageFactor / unit.MaxHealth() * unit.rageBar.damageTakenRageMultiplier

			unit.AddRage(sim, generatedRage, rageFromDamageTakenMetrics)
		},
	})

	// Not a real spell, just holds metrics from rage gain threat.
	unit.RegisterSpell(SpellConfig{
		ActionID: ActionID{OtherID: proto.OtherAction_OtherActionRageGain},
	})

	maxRage := max(100.0, options.MaxRage)

	unit.rageBar = rageBar{
		unit:                      unit,
		maxRage:                   maxRage,
		startingRage:              max(0, min(options.StartingRage, maxRage)),
		offHandRageMultiplier:     1,
		damageTakenRageMultiplier: 1,
		RageRefundMetrics:         unit.NewRageMetrics(ActionID{OtherID: proto.OtherAction_OtherActionRefund}),
		EncounterStartMetrics:     unit.NewRageMetrics(ActionID{OtherID: proto.OtherAction_OtherActionEncounterStart}),
	}
}

func (unit *Unit) HasRageBar() bool {
	return unit.rageBar.unit != nil
}

func (rb *rageBar) CurrentRage() float64 {
	return rb.currentRage
}

func (rb *rageBar) MaximumRage() float64 {
	return rb.maxRage
}

func (rb *rageBar) SetOffHandRageMultiplier(multiplier float64) {
	rb.offHandRageMultiplier = multiplier
}

func (rb *rageBar) MultiplyDamageTakenRageGen(multiplier float64) {
	rb.damageTakenRageMultiplier *= multiplier
}

func (rb *rageBar) AddRage(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to add negative rage!")
	}

	newRage := min(rb.currentRage+amount, rb.maxRage)
	metrics.AddEvent(amount, newRage-rb.currentRage)

	if sim.Log != nil {
		rb.unit.Log(sim, "Gained %0.3f rage from %s (%0.3f --> %0.3f) of %0.0f total.", amount, metrics.ActionID, rb.currentRage, newRage, 100.0)
	}

	rb.currentRage = newRage
	if !sim.Options.Interactive {
		rb.unit.ReactToEvent(sim, false, true)
	}
}

func (rb *rageBar) SpendRage(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to spend negative rage!")
	}

	newRage := rb.currentRage - amount
	metrics.AddEvent(-amount, -amount)

	if sim.Log != nil {
		rb.unit.Log(sim, "Spent %0.3f rage from %s (%0.3f --> %0.3f) of %0.0f total.", amount, metrics.ActionID, rb.currentRage, newRage, 100.0)
	}

	rb.currentRage = newRage
}

func (rb *rageBar) ResetRageBar(sim *Simulation, rageToKeep float64) {
	if rb.currentRage > rageToKeep {
		rb.SpendRage(sim, rb.currentRage-rageToKeep, rb.EncounterStartMetrics)
	} else if rageToKeep > rb.currentRage {
		rb.AddRage(sim, rageToKeep-rb.currentRage, rb.EncounterStartMetrics)
	}
}

func (rb *rageBar) reset(_ *Simulation) {
	if rb.unit == nil {
		return
	}

	rb.currentRage = rb.startingRage
	rb.damageTakenRageMultiplier = 1
}

func (rb *rageBar) doneIteration() {
	if rb.unit == nil {
		return
	}

	rageGainSpell := rb.unit.GetSpell(ActionID{OtherID: proto.OtherAction_OtherActionRageGain})

	for _, resourceMetrics := range rb.unit.Metrics.resources {
		if resourceMetrics.Type != proto.ResourceType_ResourceTypeRage {
			continue
		}
		if resourceMetrics.ActionID.SameActionIgnoreTag(ActionID{OtherID: proto.OtherAction_OtherActionDamageTaken}) {
			continue
		}
		if resourceMetrics.ActionID.SameActionIgnoreTag(ActionID{OtherID: proto.OtherAction_OtherActionRefund}) {
			continue
		}
		if resourceMetrics.NoThreat || resourceMetrics.ActualGainForCurrentIteration() <= 0 {
			continue
		}

		// Need to exclude rage gained from white hits. Rather than have a manual list of all IDs that would
		// apply here (autos, WF attack, sword spec procs, etc), just check if the effect caused any damage.
		sourceSpell := rb.unit.GetSpell(resourceMetrics.ActionID)
		if sourceSpell != nil && sourceSpell.SpellMetrics[0].TotalDamage > 0 {
			continue
		}

		rageGainSpell.SpellMetrics[0].Casts += resourceMetrics.EventsForCurrentIteration()
		rageGainSpell.ApplyAOEThreatIgnoreMultipliers(resourceMetrics.ActualGainForCurrentIteration() * ThreatPerRageGained)
	}
}

type RageCostOptions struct {
	Cost int32

	Refund        float64
	RefundMetrics *ResourceMetrics // Optional, will default to unit.RageRefundMetrics if not supplied.
}
type RageCost struct {
	Refund          float64
	RefundMetrics   *ResourceMetrics
	ResourceMetrics *ResourceMetrics
}

func newRageCost(spell *Spell, options RageCostOptions) *SpellCost {
	if options.Refund > 0 && options.RefundMetrics == nil {
		options.RefundMetrics = spell.Unit.RageRefundMetrics
	}

	return &SpellCost{
		spell:                   spell,
		BaseCost:                options.Cost,
		PercentModifier:         1,
		AdditivePercentModifier: 1,
		ResourceCostImpl: &RageCost{
			Refund:          options.Refund,
			RefundMetrics:   options.RefundMetrics,
			ResourceMetrics: spell.Unit.NewRageMetrics(spell.ActionID),
		},
	}
}

func (rc *RageCost) MeetsRequirement(_ *Simulation, spell *Spell) bool {
	spell.CurCast.Cost = spell.Cost.GetCurrentCost()
	return spell.Unit.CurrentRage() >= spell.CurCast.Cost
}
func (rc *RageCost) CostFailureReason(sim *Simulation, spell *Spell) string {
	return fmt.Sprintf("not enough rage (Current Rage = %0.03f, Rage Cost = %0.03f)", spell.Unit.CurrentRage(), spell.CurCast.Cost)
}
func (rc *RageCost) SpendCost(sim *Simulation, spell *Spell) {
	if spell.CurCast.Cost > 0 {
		spell.Unit.SpendRage(sim, spell.CurCast.Cost, rc.ResourceMetrics)
	}
}
func (rc *RageCost) IssueRefund(sim *Simulation, spell *Spell) {
	if rc.Refund > 0 && spell.CurCast.Cost > 0 {
		spell.Unit.AddRage(sim, rc.Refund*spell.CurCast.Cost, rc.RefundMetrics)
	}
}

func (spell *Spell) RageMetrics() *ResourceMetrics {
	return spell.Cost.ResourceCostImpl.(*RageCost).ResourceMetrics
}
