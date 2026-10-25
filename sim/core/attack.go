package core

import (
	"os"
	"slices"
	"time"

	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
)

// ReplaceMHSwing is called right before a main hand auto attack fires.
// It must never return nil, but either a replacement spell or the passed in regular mhSwingSpell.
// This allows for abilities that convert a white attack into a yellow attack.
type ReplaceMHSwing func(sim *Simulation, mhSwingSpell *Spell) *Spell

// Represents a generic weapon. Pets / unarmed / various other cases don't use
// actual weapon items so this is an abstraction of a Weapon.
type Weapon struct {
	BaseDamageMin        float64
	BaseDamageMax        float64
	AttackPowerPerDPS    float64
	SwingSpeed           float64
	NormalizedSwingSpeed float64
	SpellSchool          SpellSchool
}

func (weapon *Weapon) DPS() float64 {
	if weapon.SwingSpeed == 0 {
		return 0
	}
	return (weapon.BaseDamageMin + weapon.BaseDamageMax) / 2.0 / weapon.SwingSpeed
}

func newWeaponFromUnarmed() Weapon {
	// These numbers are probably wrong but nobody cares.
	return Weapon{
		BaseDamageMin:        0,
		BaseDamageMax:        0,
		SwingSpeed:           1,
		NormalizedSwingSpeed: 1,
		AttackPowerPerDPS:    DefaultAttackPowerPerDPS,
	}
}

func newWeaponFromItem(item *Item, bonusDps float64) Weapon {
	normalizedWeaponSpeed := 2.4
	if item.WeaponType == proto.WeaponType_WeaponTypeDagger {
		normalizedWeaponSpeed = 1.7
	} else if item.HandType == proto.HandType_HandTypeTwoHand {
		normalizedWeaponSpeed = 3.3
	} else if item.RangedWeaponType != proto.RangedWeaponType_RangedWeaponTypeUnknown {
		normalizedWeaponSpeed = 2.8
	}

	return Weapon{
		BaseDamageMin:        item.WeaponDamageMin + bonusDps*item.SwingSpeed,
		BaseDamageMax:        item.WeaponDamageMax + bonusDps*item.SwingSpeed,
		SwingSpeed:           item.SwingSpeed,
		NormalizedSwingSpeed: normalizedWeaponSpeed,
		AttackPowerPerDPS:    DefaultAttackPowerPerDPS,
	}
}

// Returns weapon stats using the main hand equipped weapon.
func (character *Character) WeaponFromMainHand() Weapon {
	if weapon := character.GetMHWeapon(); weapon != nil {
		return newWeaponFromItem(weapon, character.PseudoStats.BonusMHDps)
	} else {
		return newWeaponFromUnarmed()
	}
}

// Returns weapon stats using the off-hand equipped weapon.
func (character *Character) WeaponFromOffHand() Weapon {
	if weapon := character.GetOHWeapon(); weapon != nil {
		return newWeaponFromItem(weapon, character.PseudoStats.BonusOHDps)
	} else {
		return Weapon{}
	}
}

// Returns weapon stats using the ranged equipped weapon.
func (character *Character) WeaponFromRanged() Weapon {
	if weapon := character.GetRangedWeapon(); weapon != nil {
		return newWeaponFromItem(weapon, character.PseudoStats.BonusRangedDps)
	} else {
		return Weapon{}
	}
}

func (weapon *Weapon) GetSpellSchool() SpellSchool {
	if weapon.SpellSchool == SpellSchoolNone {
		return SpellSchoolPhysical
	} else {
		return weapon.SpellSchool
	}
}

func (weapon *Weapon) EnemyWeaponDamage(sim *Simulation, attackPower float64, damageSpread float64) float64 {
	// Maximum damage range is 133% of minimum damage; AP contribution is % of minimum damage roll.
	// Patchwerk follows special damage range rules.
	// TODO: Scrape more logs to determine these values more accurately. AP defined in constants.go

	rand := 1 + damageSpread*sim.RandomFloat("Enemy Weapon Damage")

	return weapon.BaseDamageMin * (rand + attackPower*EnemyAutoAttackAPCoefficient)
}

func (weapon *Weapon) BaseDamage(sim *Simulation) float64 {
	return weapon.BaseDamageMin + (weapon.BaseDamageMax-weapon.BaseDamageMin)*sim.RandomFloat("Weapon Base Damage")
}

func (weapon *Weapon) AverageDamage() float64 {
	return (weapon.BaseDamageMin + weapon.BaseDamageMax) / 2
}

func (weapon *Weapon) CalculateWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return weapon.BaseDamage(sim) + (weapon.SwingSpeed*attackPower)/weapon.AttackPowerPerDPS
}

func (weapon *Weapon) CalculateAverageWeaponDamage(attackPower float64) float64 {
	return weapon.AverageDamage() + (weapon.SwingSpeed*attackPower)/weapon.AttackPowerPerDPS
}

// Normalized instants (Strider Kick, Mongoose Bite) take the plain normalized speed, slow weapons too: a 3.8
// two-hander's Strider Kicks land at its 3.3 (16 kicks at 192.9, 192.2 predicted at 3.3, 204.7 at 3.8), and a 2.8 and
// a 1.7 one-hander in the same session hit 24.7 apart, their base damage gap (24.5), not the 33 a kept 2.8 would add.
// Knob AB_NORM_MAX=1 gives a slower weapon its real speed instead, for comparisons.
func (weapon *Weapon) CalculateNormalizedWeaponDamage(sim *Simulation, attackPower float64) float64 {
	speed := weapon.NormalizedSwingSpeed
	if os.Getenv("AB_NORM_MAX") == "1" {
		speed = max(speed, weapon.SwingSpeed)
	}
	return weapon.BaseDamage(sim) + (speed*attackPower)/weapon.AttackPowerPerDPS
}

func (unit *Unit) MHWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return unit.AutoAttacks.mh.CalculateWeaponDamage(sim, attackPower)
}
func (unit *Unit) MHNormalizedWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return unit.AutoAttacks.mh.CalculateNormalizedWeaponDamage(sim, attackPower)
}

func (unit *Unit) OHWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return 0.5 * unit.AutoAttacks.oh.CalculateWeaponDamage(sim, attackPower)
}
func (unit *Unit) OHNormalizedWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return 0.5 * unit.AutoAttacks.oh.CalculateNormalizedWeaponDamage(sim, attackPower)
}

func (unit *Unit) RangedWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return unit.AutoAttacks.ranged.CalculateWeaponDamage(sim, attackPower)
}

// Returns whether this hit effect is associated with the main-hand weapon.
func (spell *Spell) IsMH() bool {
	return spell.ProcMask.Matches(ProcMaskMeleeMH)
}

// Returns whether this hit effect is associated with the off-hand weapon.
func (spell *Spell) IsOH() bool {
	return spell.ProcMask.Matches(ProcMaskMeleeOH)
}

// Returns whether this hit effect is associated with either melee weapon.
func (spell *Spell) IsMelee() bool {
	return spell.ProcMask.Matches(ProcMaskMelee)
}

func (aa *AutoAttacks) MH() *Weapon {
	return aa.mh.getWeapon()
}

func (aa *AutoAttacks) SetMH(weapon Weapon) {
	aa.mh.setWeapon(weapon)

	if aa.mh.extraAttacksAura == nil {
		aa.mh.extraAttacksAura = aa.mh.unit.GetAuraByID(ActionID{SpellID: 21919})
	}
}

func (aa *AutoAttacks) OH() *Weapon {
	return aa.oh.getWeapon()
}

func (aa *AutoAttacks) SetOH(weapon Weapon) {
	aa.oh.setWeapon(weapon)
}

func (aa *AutoAttacks) Ranged() *Weapon {
	return aa.ranged.getWeapon()
}

func (aa *AutoAttacks) SetRanged(weapon Weapon) {
	aa.ranged.setWeapon(weapon)
}

func (aa *AutoAttacks) MHAuto() *Spell {
	return aa.mh.spell
}

func (aa *AutoAttacks) OHAuto() *Spell {
	return aa.oh.spell
}

func (aa *AutoAttacks) RangedAuto() *Spell {
	return aa.ranged.spell
}

func (aa *AutoAttacks) MainhandSwingAt() time.Duration {
	return aa.mh.swingAt
}

func (aa *AutoAttacks) OffhandSwingAt() time.Duration {
	return aa.oh.swingAt
}

func (aa *AutoAttacks) SetOffhandSwingAt(offhandSwingAt time.Duration) {
	aa.oh.swingAt = offhandSwingAt
}

func (aa *AutoAttacks) SetReplaceMHSwing(replaceSwing ReplaceMHSwing) {
	aa.mh.replaceSwing = replaceSwing
}

func (aa *AutoAttacks) MHConfig() *SpellConfig {
	return &aa.mh.config
}

// RegisterExtraMHSwing hands register a copy of the main-hand swing under actionID once the auto attacks are finalized,
// so it carries every change made to the main-hand config: an extra auto attack that rolls on the white table and
// leaves the swing timer alone. Nothing is registered for a unit that does not swing in melee.
func (aa *AutoAttacks) RegisterExtraMHSwing(actionID ActionID, register func(*Spell)) {
	aa.extraMHSwings = append(aa.extraMHSwings, extraMHSwing{actionID: actionID, register: register})
}

// MaybeReplaceMHSwing gives a queued next-melee ability (Raptor Strike) the swing, as a main-hand auto attack would.
func (aa *AutoAttacks) MaybeReplaceMHSwing(sim *Simulation, mhSwingSpell *Spell) *Spell {
	if aa.mh.replaceSwing == nil {
		return mhSwingSpell
	}
	return aa.mh.replaceSwing(sim, mhSwingSpell)
}

func (aa *AutoAttacks) OHConfig() *SpellConfig {
	return &aa.oh.config
}

func (aa *AutoAttacks) RangedConfig() *SpellConfig {
	return &aa.ranged.config
}

const (
	tagMainhand    = 1
	tagOffhand     = 2
	tagExtraAttack = 3
)

type WeaponAttack struct {
	Weapon

	agent Agent
	unit  *Unit

	config SpellConfig
	spell  *Spell

	replaceSwing ReplaceMHSwing

	swingAt             time.Duration
	lastSwingAt         time.Duration
	// when the weapon last really swung (or the Auto Shot fired); a restarted timer leaves it (autoTimeSinceLast)
	lastAttackAt time.Duration
	// when the running timer is ready with nothing in the way (walking, the stand-still wait, the off-hand wait): the
	// timeline draws the timer up to here and any time after it that the weapon sat ready
	readyAt time.Duration
	// a swing retry wait since retryFrom (steps from retryOrigin): the timeline gets it once the swing lands, or as
	// left if a walk takes the unit out first
	retryPending bool
	retryFrom    time.Duration
	retryOrigin  time.Duration
	extraAttacks        int32 // extra attacks that happen right away
	extraAttacksStored  int32 // extra attack that happen on next auto (e.g reckoning)
	extraAttacksPending int32 // extraAttacks prior to previous ones resolving for spell metrics
	extraAttacksAura    *Aura

	curSwingSpeed    float64
	curSwingDuration time.Duration
}

func (wa *WeaponAttack) getWeapon() *Weapon {
	return &wa.Weapon
}

func (wa *WeaponAttack) setWeapon(weapon Weapon) {
	wa.Weapon = weapon
	wa.updateSwingDuration(wa.curSwingSpeed)
}

// startTimer starts a timer cycle now, ready at readyAt.
func (wa *WeaponAttack) startTimer(sim *Simulation, readyAt time.Duration) {
	wa.readyAt = readyAt
	if sim.Log != nil && wa.spell != nil {
		wa.unit.Log(sim, "Swing timer %s started, ready at %.3fs", wa.timerID(), readyAt.Seconds())
	}
}

// moveReadyAt moves the running timer's ready time (an attack speed change).
func (wa *WeaponAttack) moveReadyAt(sim *Simulation, readyAt time.Duration) {
	wa.readyAt = readyAt
	if sim.Log != nil && wa.spell != nil {
		wa.unit.Log(sim, "Swing timer %s ready at %.3fs", wa.timerID(), readyAt.Seconds())
	}
}

// timerID names the weapon in the log whatever metrics tag its spell carries at the moment (an extra attack from Hand of
// Justice swings with the extra-attack tag but restarts the main hand's own timer).
func (wa *WeaponAttack) timerID() ActionID {
	switch aa := &wa.unit.AutoAttacks; wa {
	case &aa.mh:
		return ActionID{OtherID: proto.OtherAction_OtherActionAttack, Tag: tagMainhand}
	case &aa.oh:
		return ActionID{OtherID: proto.OtherAction_OtherActionAttack, Tag: tagOffhand}
	}
	return ActionID{OtherID: proto.OtherAction_OtherActionShoot}
}

// inlineable stub for swing
func (wa *WeaponAttack) trySwing(sim *Simulation) time.Duration {
	if sim.CurrentTime < wa.swingAt {
		return wa.swingAt
	}
	return wa.swing(sim)
}

func (wa *WeaponAttack) castExtraAttacksStored(sim *Simulation) {
	if wa.extraAttacksAura != nil {
		wa.extraAttacksStored = wa.extraAttacksAura.GetStacks()
		wa.extraAttacksAura.SetStacks(sim, 0)
	}

	wa.castExtraAttacks(sim, wa.extraAttacksStored, 0)
	wa.extraAttacksStored = 0

}

func (wa *WeaponAttack) castExtraAttacksTriggered(sim *Simulation, moreAttacks bool) {
	if moreAttacks {
		wa.castExtraAttacks(sim, wa.extraAttacks, 0)
	} else {
		wa.castExtraAttacks(sim, wa.extraAttacks, 1)
	}
	wa.extraAttacks = 0
}

func (wa *WeaponAttack) castExtraAttacks(sim *Simulation, numExtraAttacks int32, startIndex int32) bool {
	if numExtraAttacks > 0 {
		// if startIndex==1, Ignore the first extra attack, that was used to speed up next attack
		wa.spell.SetMetricsSplit(1)

		for i := int32(startIndex); i < numExtraAttacks; i++ {
			// use original attacks for subsequent extra Attacks
			wa.spell.Cast(sim, wa.unit.CurrentTarget)
		}

		return true
	}
	return false

}

func (wa *WeaponAttack) swing(sim *Simulation) time.Duration {
	if wa.spell != nil && wa.spell.ProcMask.Matches(ProcMaskRangedAuto) {
		if target := wa.unit.CurrentTarget; target != nil {
			if !wa.spell.Cast(sim, target) {
				wa.spell.ApplyEffects(sim, target, wa.spell)
			}
		}
		wa.updateSwingDuration(wa.unit.RangedSwingSpeed())
		wa.swingAt = sim.CurrentTime + max(wa.curSwingDuration-wa.unit.AutoAttacks.RangedCycleCut, 0)
		wa.lastSwingAt = sim.CurrentTime
		wa.lastAttackAt = sim.CurrentTime
		wa.startTimer(sim, wa.swingAt)
		wa.unit.AutoAttacks.shotSinceMelee = true
		wa.unit.AutoAttacks.startAttackOn = false
		// the shot is a decision point like a melee swing: off-GCD actions (a move into melee) can start right after it
		// instead of waiting out the GCD (AB_NO_SHOT_RECHECK=1 turns it off, AB_SHOT_RECHECK_MS delays it, for engine
		// comparisons)
		if !sim.Options.Interactive && wa.unit.Rotation != nil && !noShotRecheck {
			if shotRecheckDelay > 0 {
				StartDelayedAction(sim, DelayedActionOptions{
					DoAt:     sim.CurrentTime + shotRecheckDelay,
					OnAction: func(sim *Simulation) { wa.unit.Rotation.DoNextAction(sim) },
				})
			} else {
				wa.unit.Rotation.DoNextAction(sim)
			}
		}
		return wa.swingAt
	}

	isExtraAttack := wa.extraAttacksPending > 0

	if isExtraAttack {
		wa.extraAttacks = wa.extraAttacksPending
		wa.extraAttacksPending = 0
	}

	wa.castExtraAttacksStored(sim)

	attackSpell := wa.spell

	if wa.replaceSwing != nil {
		wa.unit.Rotation.DoNextAction(sim)
		attackSpell = wa.replaceSwing(sim, attackSpell)
	}

	if attackSpell.CanCast(sim, wa.unit.CurrentTarget) {
		wa.endRetryWait(sim, "retry")
		wa.swingAt = sim.CurrentTime + wa.curSwingDuration
		wa.lastSwingAt = sim.CurrentTime
		wa.lastAttackAt = sim.CurrentTime
		wa.startTimer(sim, wa.swingAt)
		wa.unit.AutoAttacks.shotSinceMelee = false
		// a swing landing mid-walk (melee starts at 5 yd) cannot start a move now; the walk's end replays the decision
		if wa.unit.MovementHandler != nil && wa.unit.MovementHandler.Moving {
			wa.unit.MovementHandler.swungWhileMoving = true
		}

		var originalCastTime time.Duration
		if isExtraAttack {
			originalCastTime = wa.spell.DefaultCast.CastTime
			wa.spell.DefaultCast.CastTime = 0
			wa.spell.SetMetricsSplit(1)
		} else {
			wa.spell.SetMetricsSplit(0)
		}

		attackSpell.Cast(sim, wa.unit.CurrentTarget)

		// Forever: melee restarts Auto Shot from this instant (does not add a full bow cycle).

		if ForeverMeleeRestartsAutoShot && wa.unit.AutoAttacks.MeleeResetsAutoShot && wa.unit.AutoAttacks.AutoSwingRanged &&
			!attackSpell.ProcMask.Matches(ProcMaskRangedAuto) {

			// a white swing, or a queued ability that took the swing (Raptor Strike): one restart, from the swing
			isWhite := attackSpell.ProcMask.Matches(ProcMaskMeleeMHAuto | ProcMaskMeleeOHAuto)
			isQueued := attackSpell != wa.spell

			if isWhite || isQueued {
				wa.unit.AutoAttacks.RestartRangedSwing(sim)
			}
		}

		moreAttacks := !isExtraAttack && wa.extraAttacksPending > 0
		wa.castExtraAttacksTriggered(sim, moreAttacks)

		if isExtraAttack {
			if originalCastTime > 0 {
				wa.spell.WaitTravelTime(sim, func(sim *Simulation) {
					wa.spell.DefaultCast.CastTime = originalCastTime
					wa.spell.SetMetricsSplit(0)
				})
			} else {
				wa.spell.SetMetricsSplit(0)
			}
		}

		if wa.extraAttacksPending > 0 {
			wa.spell.SetMetricsSplit(1)
			wa.swingAt = sim.CurrentTime + SpellBatchWindow
			wa.lastSwingAt = sim.CurrentTime
			sim.rescheduleWeaponAttack(wa.swingAt)
		}

		if !sim.Options.Interactive && wa.unit.Rotation != nil {
			wa.unit.Rotation.DoNextAction(sim)
		}
	} else {
		wa.swingAt = max(wa.unit.Hardcast.Expires, sim.CurrentTime+time.Millisecond*100)
	}

	return wa.swingAt
}

func (wa *WeaponAttack) updateSwingDuration(curSwingSpeed float64) {
	wa.curSwingSpeed = curSwingSpeed
	wa.curSwingDuration = DurationFromSeconds(wa.SwingSpeed / wa.curSwingSpeed)
}

func (wa *WeaponAttack) addWeaponAttack(sim *Simulation, swingSpeed float64) {
	wa.updateSwingDuration(swingSpeed)
	sim.addWeaponAttack(wa)
	sim.rescheduleWeaponAttack(wa.swingAt)
}

type AutoAttacks struct {
	AutoSwingMelee  bool
	AutoSwingRanged bool

	IsDualWielding bool

	// RangedCycleCut comes off every Auto Shot timer that a shot starts (Forever's Auto Shot spam bug).
	RangedCycleCut time.Duration
	// MeleeResetsAutoShot: melee swings restart the Auto Shot timer, a live beta bug Blizzard confirmed (forever-bugs
	// #14); off, they leave it alone, as once it is fixed.
	MeleeResetsAutoShot bool
	// an Auto Shot since the last melee swing: Auto Shot turning on is what turns the melee attack off
	shotSinceMelee bool
	// the rotation's Start Attack turned the melee attack on at startAttackAt, after the last Auto Shot
	startAttackOn bool
	startAttackAt time.Duration
	// melee swings started on a walk in at 5 yd, before the walk ended
	meleeWhileMoving bool

	mh     WeaponAttack
	oh     WeaponAttack
	ranged WeaponAttack

	// copies of the main-hand swing for extra attacks, registered at finalize
	extraMHSwings []extraMHSwing

	enabled bool
}

type extraMHSwing struct {
	actionID ActionID
	register func(*Spell)
}

// Options for initializing auto attacks.
type AutoAttackOptions struct {
	MainHand        Weapon
	OffHand         Weapon
	Ranged          Weapon
	AutoSwingMelee  bool // If true, core engine will handle calling SwingMelee() for you.
	AutoSwingRanged bool // If true, core engine will handle calling SwingRanged() for you.
	ReplaceMHSwing  ReplaceMHSwing
}

func (unit *Unit) EnableAutoAttacks(agent Agent, options AutoAttackOptions) {
	if options.MainHand.AttackPowerPerDPS == 0 {
		options.MainHand.AttackPowerPerDPS = DefaultAttackPowerPerDPS
	}
	if options.OffHand.AttackPowerPerDPS == 0 {
		options.OffHand.AttackPowerPerDPS = DefaultAttackPowerPerDPS
	}

	unit.AutoAttacks = AutoAttacks{
		AutoSwingMelee:  options.AutoSwingMelee,
		AutoSwingRanged: options.AutoSwingRanged,

		IsDualWielding: options.OffHand.SwingSpeed != 0,

		mh: WeaponAttack{
			agent:        agent,
			unit:         unit,
			Weapon:       options.MainHand,
			replaceSwing: options.ReplaceMHSwing,
		},
		oh: WeaponAttack{
			agent:  agent,
			unit:   unit,
			Weapon: options.OffHand,
		},
		ranged: WeaponAttack{
			agent:  agent,
			unit:   unit,
			Weapon: options.Ranged,
		},
	}

	unit.AutoAttacks.mh.config = SpellConfig{
		ActionID:    ActionID{OtherID: proto.OtherAction_OtherActionAttack, Tag: tagMainhand},
		SpellSchool: options.MainHand.GetSpellSchool(),
		DefenseType: DefenseTypeMelee,
		ProcMask:    ProcMaskMeleeMHAuto,
		Flags:       SpellFlagMeleeMetrics | SpellFlagNoOnCastComplete,
		CastType:    proto.CastType_CastTypeMainHand,

		MetricSplits: 2,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: TernaryFloat64(options.MainHand.GetSpellSchool() == SpellSchoolPhysical, 1, 0),

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {

			autoInProgress := *spell

			baseDamage := autoInProgress.Unit.MHWeaponDamage(sim, autoInProgress.MeleeAttackPower(target))
			result := autoInProgress.CalcDamage(sim, target, baseDamage, autoInProgress.OutcomeMeleeWhite)

			StartDelayedAction(sim, DelayedActionOptions{
				DoAt: sim.CurrentTime + SpellBatchWindow,
				OnAction: func(s *Simulation) {
					autoInProgress.DealDamage(sim, result)
				},
			})
		},
	}

	unit.AutoAttacks.oh.config = SpellConfig{
		ActionID:    ActionID{OtherID: proto.OtherAction_OtherActionAttack, Tag: tagOffhand},
		SpellSchool: options.OffHand.GetSpellSchool(),
		DefenseType: DefenseTypeMelee,
		ProcMask:    ProcMaskMeleeOHAuto,
		Flags:       SpellFlagMeleeMetrics | SpellFlagNoOnCastComplete,
		CastType:    proto.CastType_CastTypeOffHand,

		MetricSplits: 2,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: TernaryFloat64(options.OffHand.GetSpellSchool() == SpellSchoolPhysical, 1, 0),

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {

			autoInProgress := *spell
			baseDamage := autoInProgress.Unit.OHWeaponDamage(sim, autoInProgress.MeleeAttackPower(target))
			result := autoInProgress.CalcDamage(sim, target, baseDamage, autoInProgress.OutcomeMeleeWhite)

			StartDelayedAction(sim, DelayedActionOptions{
				DoAt: sim.CurrentTime + SpellBatchWindow,
				OnAction: func(s *Simulation) {
					autoInProgress.DealDamage(sim, result)
				},
			})
		},
	}

	unit.AutoAttacks.ranged.config = SpellConfig{
		ActionID:     ActionID{OtherID: proto.OtherAction_OtherActionShoot},
		SpellSchool:  options.Ranged.GetSpellSchool(),
		DefenseType:  DefenseTypeRanged,
		ProcMask:     ProcMaskRangedAuto,
		Flags:        SpellFlagMeleeMetrics | SpellFlagNoOnCastComplete | SpellFlagPassiveSpell,
		CastType:     proto.CastType_CastTypeRanged,
		MissileSpeed: 24,

		MetricSplits: 2,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: TernaryFloat64(options.Ranged.GetSpellSchool() == SpellSchoolPhysical, 1, 0),

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			baseDamage := spell.Unit.RangedWeaponDamage(sim, spell.RangedAttackPower(target, false))
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	}

	if unit.Type == EnemyUnit {
		unit.AutoAttacks.mh.config.ApplyEffects = func(sim *Simulation, target *Unit, spell *Spell) {
			ap := max(0, spell.Unit.stats[stats.AttackPower])
			baseDamage := spell.Unit.AutoAttacks.mh.EnemyWeaponDamage(sim, ap, spell.Unit.PseudoStats.DamageSpread)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeEnemyMeleeWhite)
		}
		unit.AutoAttacks.oh.config.ApplyEffects = func(sim *Simulation, target *Unit, spell *Spell) {
			ap := max(0, spell.Unit.stats[stats.AttackPower])
			baseDamage := spell.Unit.AutoAttacks.mh.EnemyWeaponDamage(sim, ap, spell.Unit.PseudoStats.DamageSpread) * 0.5

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeEnemyMeleeWhite)
		}
	}
}

func (aa *AutoAttacks) finalize() {
	if aa.AutoSwingMelee {
		aa.mh.spell = aa.mh.unit.GetOrRegisterSpell(aa.mh.config)
		aa.mh.spell.TagSplitMetric(1, tagExtraAttack)
		for _, extra := range aa.extraMHSwings {
			config := aa.mh.config
			config.ActionID = extra.actionID
			config.MetricSplits = 0
			extra.register(aa.mh.unit.GetOrRegisterSpell(config))
		}
		// Will keep the OH spell registered for Item swapping
		aa.oh.spell = aa.oh.unit.GetOrRegisterSpell(aa.oh.config)
		aa.oh.spell.TagSplitMetric(1, tagExtraAttack)
	}
	if aa.AutoSwingRanged {
		aa.ranged.spell = aa.ranged.unit.GetOrRegisterSpell(aa.ranged.config)
		aa.ranged.spell.TagSplitMetric(1, tagExtraAttack)
	}
}

func (aa *AutoAttacks) reset(sim *Simulation) {
	if !aa.AutoSwingMelee && !aa.AutoSwingRanged {
		return
	}

	aa.enabled = false

	aa.mh.swingAt = NeverExpires
	aa.oh.swingAt = NeverExpires
	aa.mh.retryPending, aa.oh.retryPending = false, false
	// no swing yet this fight: long ago
	aa.mh.lastAttackAt, aa.oh.lastAttackAt, aa.ranged.lastAttackAt = -time.Hour, -time.Hour, -time.Hour

	// Make sure extra attacks are reset
	aa.mh.extraAttacks = 0
	aa.mh.extraAttacksPending = 0
	aa.mh.spell.SetMetricsSplit(0)

	if aa.ranged.spell != nil {
		aa.ranged.extraAttacks = 0
		aa.ranged.extraAttacksPending = 0
		aa.ranged.spell.SetMetricsSplit(0)
	}

	if aa.AutoSwingMelee {
		aa.mh.updateSwingDuration(aa.mh.unit.SwingSpeed())
		aa.mh.swingAt = 0

		if aa.IsDualWielding {
			aa.oh.extraAttacks = 0
			aa.oh.extraAttacksPending = 0
			aa.oh.spell.SetMetricsSplit(0)
			aa.oh.updateSwingDuration(aa.mh.curSwingSpeed)
			aa.oh.swingAt = 0

			// Apply random delay of 0 - 50% swing time, to one of the weapons if dual wielding
			if ForeverOffHandWaitsOnAttackStart && aa.oh.unit.Type != EnemyUnit {
				// the attack turns on at the pull (startPull) or with the first walk in
			} else if aa.oh.unit.Type == EnemyUnit {
				aa.oh.swingAt = DurationFromSeconds(aa.mh.SwingSpeed / 2)
			} else {
				if sim.RandomFloat("SwingResetWeapon") < 0.5 {
					aa.mh.swingAt = DurationFromSeconds(sim.RandomFloat("SwingResetDelay") * aa.mh.SwingSpeed / 2)
				} else {
					aa.oh.swingAt = DurationFromSeconds(sim.RandomFloat("SwingResetDelay") * aa.mh.SwingSpeed / 2)
				}
			}
		}

	}

	// the melee attack is not on before the pull
	aa.shotSinceMelee = true
	aa.startAttackOn = false
	aa.meleeWhileMoving = false
	aa.ranged.swingAt = NeverExpires

	if aa.AutoSwingRanged {
		aa.ranged.updateSwingDuration(aa.ranged.unit.RangedSwingSpeed())
		aa.ranged.swingAt = 0
	}
}

func (aa *AutoAttacks) startPull(sim *Simulation) {

	if !aa.AutoSwingMelee && !aa.AutoSwingRanged {
		return
	}

	if aa.enabled {
		return
	}

	// a walk already under way at the pull (it can start in the same instant, from the first Auto Shot) turns the attacks
	// back on where it ends; turning them on here would use the distance it started from
	if aa.mh.unit != nil && aa.mh.unit.IsMoving() {
		return
	}

	aa.enabled = true

	if aa.AutoSwingMelee && aa.mh.unit.DistanceFromTarget <= MaxMeleeAttackDistance {
		aa.mh.addWeaponAttack(sim, aa.mh.unit.SwingSpeed())
		if aa.IsDualWielding {
			if ForeverOffHandWaitsOnAttackStart && aa.oh.unit.Type != EnemyUnit {
				// the melee attack turns on at the pull
				aa.oh.swingAt = max(aa.oh.swingAt, sim.CurrentTime+aa.OffHandWait())
			}
			aa.oh.addWeaponAttack(sim, aa.mh.curSwingSpeed)
		}
	}

	if aa.AutoSwingRanged && aa.mh.unit.DistanceFromTarget >= MinRangedAttackDistance {
		aa.ranged.addWeaponAttack(sim, aa.ranged.unit.RangedSwingSpeed())
	}
}

// Stops the auto swing action until EnableAutoSwing (pets that are dismissed, walks)
// after being disabled.
func (aa *AutoAttacks) CancelAutoSwing(sim *Simulation) {
	if !aa.AutoSwingMelee && !aa.AutoSwingRanged {
		return
	}

	if aa.meleeWhileMoving {
		aa.stopMeleeMidMove(sim)
	}

	if !aa.enabled {
		return
	}

	aa.enabled = false

	if aa.AutoSwingMelee {
		sim.removeWeaponAttack(&aa.mh)
		if aa.IsDualWielding {
			sim.removeWeaponAttack(&aa.oh)
		}
	}

	if aa.AutoSwingRanged {
		sim.removeWeaponAttack(&aa.ranged)
	}
}

// Re-enables the auto swing action for the iteration
func (aa *AutoAttacks) EnableAutoSwing(sim *Simulation) {
	if !aa.AutoSwingMelee && !aa.AutoSwingRanged {
		return
	}

	if aa.enabled {
		return
	}

	aa.enabled = true

	if aa.meleeWhileMoving {
		// already swinging since the walk reached 5 yd
		aa.meleeWhileMoving = false
	} else if aa.AutoSwingMelee && aa.mh.unit.DistanceFromTarget <= MaxMeleeAttackDistance {
		aa.mh.swingAt = max(aa.mh.swingAt, sim.CurrentTime, 0)
		aa.mh.addWeaponAttack(sim, aa.mh.unit.SwingSpeed())
		if aa.IsDualWielding {
			aa.oh.swingAt = max(aa.oh.swingAt, sim.CurrentTime, 0)
			aa.oh.addWeaponAttack(sim, aa.mh.unit.SwingSpeed())
		}
	}

	if aa.AutoSwingRanged && aa.mh.unit.DistanceFromTarget >= MinRangedAttackDistance {
		aa.ranged.swingAt = max(aa.ranged.swingAt, sim.CurrentTime, 0)
		aa.ranged.addWeaponAttack(sim, aa.ranged.unit.RangedSwingSpeed())
	}
}

// The amount of time between two MH swings.
func (aa *AutoAttacks) MainhandSwingSpeed() time.Duration {
	return aa.mh.curSwingDuration
}

// The amount of time between two OH swings.
func (aa *AutoAttacks) OffhandSwingSpeed() time.Duration {
	return aa.oh.curSwingDuration
}

// The amount of time between two Ranged swings.
func (aa *AutoAttacks) RangedSwingSpeed() time.Duration {
	return aa.ranged.curSwingDuration - aa.RangedAuto().CastTime()
}

func (aa *AutoAttacks) UpdateSwingTimers(sim *Simulation) {
	if !aa.enabled && !aa.meleeWhileMoving {
		return
	}

	if aa.AutoSwingRanged {
		aa.ranged.updateSwingDuration(aa.ranged.unit.RangedSwingSpeed())
		// ranged attack speed changes aren't applied mid-"swing"
	}

	if aa.AutoSwingMelee {
		oldSwingSpeed := aa.mh.curSwingSpeed
		aa.mh.updateSwingDuration(aa.mh.unit.SwingSpeed())
		f := oldSwingSpeed / aa.mh.curSwingSpeed

		if remainingSwingTime := aa.mh.swingAt - sim.CurrentTime; remainingSwingTime > 0 {
			aa.mh.swingAt = sim.CurrentTime + time.Duration(float64(remainingSwingTime)*f)
		}
		if remaining := aa.mh.readyAt - sim.CurrentTime; remaining > 0 && f != 1 {
			aa.mh.moveReadyAt(sim, sim.CurrentTime+time.Duration(float64(remaining)*f))
		}

		sim.rescheduleWeaponAttack(aa.mh.swingAt)

		if aa.IsDualWielding {
			aa.oh.updateSwingDuration(aa.mh.curSwingSpeed)

			if remainingSwingTime := aa.oh.swingAt - sim.CurrentTime; remainingSwingTime > 0 {
				aa.oh.swingAt = sim.CurrentTime + time.Duration(float64(remainingSwingTime)*f)
			}
			if remaining := aa.oh.readyAt - sim.CurrentTime; remaining > 0 && f != 1 {
				aa.oh.moveReadyAt(sim, sim.CurrentTime+time.Duration(float64(remaining)*f))
			}

			sim.rescheduleWeaponAttack(aa.oh.swingAt)
		}
	}
}

func (aa *AutoAttacks) ExtraMHAttackProc(sim *Simulation, attacks int32, actionID ActionID, spell *Spell) {
	if spell.Flags.Matches(SpellFlagBatchStopAttackMacro) {
		aa.StoreExtraMHAttack(sim, attacks, actionID, spell.ActionID)
	} else {
		aa.ExtraMHAttack(sim, attacks, actionID, spell.ActionID)
	}
}

// ExtraMHAttack should be used for all "extra attack" procs in Classic Era versions, including Wild Strikes and Hand of Justice. In vanilla, these procs don't actually grant a full extra attack, but instead just advance the MH swing timer.
func (aa *AutoAttacks) ExtraMHAttack(sim *Simulation, attacks int32, actionID ActionID, triggerAction ActionID) {
	if attacks == 0 {
		return
	}
	if sim.Log != nil {
		attacksText := Ternary(attacks == 1, "attack", "attacks")
		aa.mh.unit.Log(sim, "gained %d extra main-hand %s from %s triggered by %s", attacks, attacksText, actionID, triggerAction)
	}

	aa.mh.swingAt = sim.CurrentTime
	aa.mh.spell.SetMetricsSplit(1)
	sim.rescheduleWeaponAttack(aa.mh.swingAt)
	aa.mh.extraAttacksPending += attacks
}

func (aa *AutoAttacks) StoreExtraMHAttack(sim *Simulation, attacks int32, actionID ActionID, triggerAction ActionID) {
	if attacks == 0 {
		return
	}

	if aa.mh.extraAttacksAura == nil {
		aa.mh.extraAttacksAura = aa.mh.unit.GetAuraByID(ActionID{SpellID: 21919})
	}

	if aa.mh.extraAttacksAura == nil {
		aa.mh.extraAttacksAura = MakePermanent(aa.mh.unit.GetOrRegisterAura(Aura{
			Label:     "Extra Attacks (Main Hand)", // Tracks Stored Extra Attacks from all sources
			ActionID:  ActionID{SpellID: 21919},    // Thrash ID
			Duration:  NeverExpires,
			MaxStacks: 4, // Max is 4 extra attacks stored - more can proc after
		}))
	}

	if !aa.mh.extraAttacksAura.IsActive() {
		aa.mh.extraAttacksAura.Activate(sim)
	}

	aa.mh.extraAttacksAura.AddStacks(sim, attacks)

	if sim.Log != nil {
		aa.mh.unit.Log(sim, "stored %d extra main-hand attacks from %s triggered by %s, total is %d", attacks, actionID, triggerAction, aa.mh.extraAttacksStored)
	}
}

// ExtraMHAttack should be used for all "extra attack" procs in Classic Era versions, including Wild Strikes and Hand of Justice. In vanilla, these procs don't actually grant a full extra attack, but instead just advance the MH swing timer.
func (aa *AutoAttacks) ExtraOHAttack(sim *Simulation, attacks int32, actionID ActionID, triggerAction ActionID) {
	if attacks == 0 {
		return
	}
	if sim.Log != nil {
		attacksText := Ternary(attacks == 1, "attack", "attacks")
		aa.oh.unit.Log(sim, "gained %d extra off-hand %s from %s triggered by %s", attacks, attacksText, actionID, triggerAction)
	}
	aa.oh.swingAt = sim.CurrentTime + SpellBatchWindow
	aa.oh.spell.SetMetricsSplit(1)
	sim.rescheduleWeaponAttack(aa.oh.swingAt)
	aa.oh.extraAttacksPending += attacks
}

// ExtraRangedAttack should be used for all "extra ranged attack" procs in Classic Era versions, including Hand of Injustice. In vanilla, these procs don't actually grant a full extra attack, but instead just advance the Ranged swing timer.
// Note that Hand of Injustice doesn't seem to reset the swing timer however.
func (aa *AutoAttacks) ExtraRangedAttack(sim *Simulation, attacks int32, actionID ActionID, triggerAction ActionID) {
	if sim.Log != nil {
		attacksText := Ternary(attacks == 1, "attack", "attacks")
		aa.mh.unit.Log(sim, "gained %d extra ranged %s from %s triggered by %s", attacks, attacksText, actionID, triggerAction)
	}
	aa.ranged.swingAt = sim.CurrentTime + SpellBatchWindow
	aa.ranged.spell.SetMetricsSplit(1)
	sim.rescheduleWeaponAttack(aa.ranged.swingAt)
	aa.ranged.extraAttacks += attacks
}

// StopMeleeUntil should be used whenever a non-melee spell is cast. It stops melee, then restarts it
// at end of cast, but with a reset swing timer (as if swings had just landed).
func (aa *AutoAttacks) StopMeleeUntil(sim *Simulation, readyAt time.Duration, desyncOH bool) {
	if !aa.AutoSwingMelee { // if not auto swinging, don't auto restart.
		return
	}

	aa.mh.swingAt = readyAt + aa.mh.curSwingDuration
	sim.rescheduleWeaponAttack(aa.mh.swingAt)

	if aa.IsDualWielding {
		aa.oh.swingAt = readyAt + aa.oh.curSwingDuration
		if desyncOH {
			// Used by warrior to desync offhand after unglyphed Shattering Throw.
			aa.oh.swingAt += aa.oh.curSwingDuration / 2
		}
		sim.rescheduleWeaponAttack(aa.oh.swingAt)
	}
}

func (aa *AutoAttacks) StopRangedUntil(sim *Simulation, readyAt time.Duration) {
	if !aa.AutoSwingRanged { // if not auto swinging, don't auto restart.
		return
	}
	aa.ranged.swingAt = readyAt + aa.ranged.curSwingDuration
	sim.rescheduleWeaponAttack(aa.ranged.swingAt)
}

// Delays all swing timers for the specified amount (Slam, and the start of melee after a step inside melee range).
func (aa *AutoAttacks) DelayMeleeBy(sim *Simulation, delay time.Duration) {
	if delay <= 0 {
		return
	}

	aa.mh.swingAt += delay
	sim.rescheduleWeaponAttack(aa.mh.swingAt)

	if aa.IsDualWielding {
		aa.oh.swingAt += delay
		sim.rescheduleWeaponAttack(aa.oh.swingAt)
	}
}

func (aa *AutoAttacks) DelayRangedUntil(sim *Simulation, readyAt time.Duration) {
	if readyAt <= aa.ranged.swingAt {
		return
	}

	aa.ranged.swingAt = readyAt
	sim.rescheduleWeaponAttack(aa.ranged.swingAt)
}
func (aa *AutoAttacks) RestartRangedSwing(sim *Simulation) {
	if !ForeverMeleeRestartsAutoShot || !aa.MeleeResetsAutoShot || !aa.AutoSwingRanged || aa.ranged.spell == nil {
		return
	}
	aa.ranged.updateSwingDuration(aa.ranged.unit.RangedSwingSpeed())
	aa.ranged.swingAt = sim.CurrentTime + aa.ranged.curSwingDuration
	aa.ranged.lastSwingAt = sim.CurrentTime
	aa.ranged.startTimer(sim, aa.ranged.swingAt)
	sim.rescheduleWeaponAttack(aa.ranged.swingAt)
}

// Returns the time at which the next melee attack will occur.
func (aa *AutoAttacks) NextAttackAt() time.Duration {
	return min(aa.mh.swingAt, aa.oh.swingAt)
}

// Returns the time at which the next attack will occur.
func (aa *AutoAttacks) NextAnyAttackAt() time.Duration {
	return min(min(aa.mh.swingAt, aa.oh.swingAt), aa.ranged.swingAt)
}

// Returns the time at which the next ranged attack will occur.
func (aa *AutoAttacks) NextRangedAttackAt() time.Duration {
	return aa.ranged.swingAt
}

type PPMManager struct {
	ppm         float64
	procMasks   []ProcMask
	procChances []float64

	// For feral druids, certain PPM effects use their equipped weapon speed
	// instead of their paw attack speed.
	mhSpecialProcChance float64
	ohSpecialProcChance float64
}

// Returns whether the effect procced.
func (ppmm *PPMManager) Proc(sim *Simulation, procMask ProcMask, label string) bool {
	for i, m := range ppmm.procMasks {
		if m.Matches(procMask) {
			return sim.RandomFloat(label) < ppmm.procChances[i]
		}
	}
	return false
}

// Returns whether the effect procced.
// This is different from Proc() in that yellow melee hits use a proc chance based on the equipped
// weapon speed rather than the base attack speed. This distinction matters for feral druids.
func (ppmm *PPMManager) ProcWithWeaponSpecials(sim *Simulation, procMask ProcMask, label string) bool {
	if procMask.Matches(ProcMaskMeleeMHSpecial) {
		return sim.RandomFloat(label) < ppmm.mhSpecialProcChance
	} else if procMask.Matches(ProcMaskMeleeOHSpecial) {
		return sim.RandomFloat(label) < ppmm.ohSpecialProcChance
	} else {
		return ppmm.Proc(sim, procMask, label)
	}
}

func (ppmm *PPMManager) Chance(procMask ProcMask) float64 {
	for i, m := range ppmm.procMasks {
		if m.Matches(procMask) {
			return ppmm.procChances[i]
		}
	}
	return 0
}

func (ppmm *PPMManager) GetPPM() float64 {
	return ppmm.ppm
}

func (aa *AutoAttacks) NewPPMManager(ppm float64, procMask ProcMask) PPMManager {
	if !aa.AutoSwingMelee && !aa.AutoSwingRanged {
		return PPMManager{}
	}

	ppmm := PPMManager{ppm: ppm, procMasks: make([]ProcMask, 0, 2), procChances: make([]float64, 0, 2)}

	mergeOrAppend := func(speed float64, mask ProcMask) {
		if speed == 0 || mask == 0 {
			return
		}

		if i := slices.Index(ppmm.procChances, speed); i != -1 {
			ppmm.procMasks[i] |= mask
			return
		}

		ppmm.procMasks = append(ppmm.procMasks, mask)
		ppmm.procChances = append(ppmm.procChances, speed)
	}

	mergeOrAppend(aa.mh.SwingSpeed, procMask&^ProcMaskRanged&^ProcMaskMeleeOH) // "everything else", even if not explicitly flagged MH
	mergeOrAppend(aa.oh.SwingSpeed, procMask&ProcMaskMeleeOH)
	mergeOrAppend(aa.ranged.SwingSpeed, procMask&ProcMaskRanged)

	for i := range ppmm.procChances {
		ppmm.procChances[i] *= ppm / 60
	}

	character := aa.mh.agent.GetCharacter()
	if procMask.Matches(ProcMaskMeleeMH) {
		if character != nil {
			if mhWeapon := character.GetMHWeapon(); mhWeapon != nil {
				ppmm.mhSpecialProcChance = mhWeapon.SwingSpeed * ppm / 60
			}
		}
	}
	if procMask.Matches(ProcMaskMeleeOH) {
		if character != nil {
			if ohWeapon := character.GetOHWeapon(); ohWeapon != nil {
				ppmm.ohSpecialProcChance = ohWeapon.SwingSpeed * ppm / 60
			}
		}
	}
	return ppmm
}

// Returns whether a PPM-based effect procced.
// Using NewPPMManager() is preferred; this function should only be used when
// the attacker is not known at initialization time.
func (aa *AutoAttacks) PPMProc(sim *Simulation, ppm float64, procMask ProcMask, label string, spell *Spell) bool {
	if !aa.AutoSwingMelee && !aa.AutoSwingRanged {
		return false
	}

	switch {
	case spell.ProcMask.Matches(procMask &^ ProcMaskMeleeOH &^ ProcMaskRanged):
		return sim.RandomFloat(label) < ppm*aa.mh.SwingSpeed/60.0
	case spell.ProcMask.Matches(procMask & ProcMaskMeleeOH):
		return sim.RandomFloat(label) < ppm*aa.oh.SwingSpeed/60.0
	case spell.ProcMask.Matches(procMask & ProcMaskRanged):
		return sim.RandomFloat(label) < ppm*aa.ranged.SwingSpeed/60.0
	}
	return false
}

func (aa *AutoAttacks) ForceMeleeSwing(sim *Simulation) {
	if aa.mh.spell == nil || !aa.AutoSwingMelee {
		return
	}
	// the arrival sets the swing from the timer that ran out at range (ForeverSwingRetryBeat)
	if ForeverSwingRetryBeat {
		return
	}
	// Entering melee must not reset a 3.8s timer that still has time left.
	if aa.mh.swingAt > sim.CurrentTime {
		return
	}
	aa.mh.swingAt = sim.CurrentTime
	sim.rescheduleWeaponAttack(aa.mh.swingAt)
}

func (aa *AutoAttacks) ForceRangedSwing(sim *Simulation) {
	if aa.ranged.spell == nil || !aa.AutoSwingRanged {
		return
	}
	if aa.ranged.swingAt <= sim.CurrentTime {
		return
	}
	aa.ranged.swingAt = sim.CurrentTime
	sim.rescheduleWeaponAttack(aa.ranged.swingAt)
}

func (unit *Unit) applyParryHaste() {
	if !unit.PseudoStats.ParryHaste || !unit.AutoAttacks.AutoSwingMelee {
		return
	}

	unit.RegisterAura(Aura{
		Label:    "Parry Haste",
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if !result.DidParry() {
				return
			}

			currentSwingTime := aura.Unit.AutoAttacks.mh.swingAt - aura.Unit.AutoAttacks.mh.lastSwingAt
			swingSpeed := aura.Unit.AutoAttacks.mh.curSwingDuration
			minRemainingTime := time.Duration(float64(swingSpeed) * 0.2) // 20% of Swing Speed
			defaultReduction := minRemainingTime * 2                     // 40% of Swing Speed

			if currentSwingTime <= minRemainingTime {
				return
			}

			newReadyAt := max(aura.Unit.AutoAttacks.mh.swingAt-defaultReduction, aura.Unit.AutoAttacks.mh.lastSwingAt+minRemainingTime)
			parryHasteReduction := newReadyAt - aura.Unit.AutoAttacks.mh.swingAt
			if sim.Log != nil {
				aura.Unit.Log(sim, "MH Swing reduced by %s due to parry haste, will now occur at %s", parryHasteReduction, newReadyAt)
			}

			aura.Unit.AutoAttacks.mh.swingAt = newReadyAt
			sim.rescheduleWeaponAttack(newReadyAt)
		},
	})
}

// ForeverMeleeRestartsAutoShot: on the beta every melee swing and Raptor Strike restarts the Auto Shot timer (the swing
// timer events and logged Auto Shot casts both put the next shot one full cycle after the swing), a bug Blizzard
// confirmed (forever-bugs #14). A hunter models it only with AutoAttacks.MeleeResetsAutoShot; false turns it off for all.
var ForeverMeleeRestartsAutoShot = true

// ForeverOffHandWaitsOnAttackStart: when the melee attack turns on, the main hand swings as soon as it is ready but the
// off hand's first swing comes half its swing time later, however long it has been ready. The wait runs from the moment
// the attack turns on, even while still walking in (measured with the swing timer events for a 1.55 and a 2.33 off hand:
// attack on + half the off-hand swing time, or arrival when the walk took longer). Shooting turns the attack off; with it
// left on, walking out and back in swings both hands on arrival. False only for comparisons.
var ForeverOffHandWaitsOnAttackStart = true

// ForeverMeleeMidWalk: a walk into melee starts swinging the moment it reaches 5 yd, not when it ends (the swing timer
// events show swings while moving). False only for comparisons.
var ForeverMeleeMidWalk = true

// ForeverSwingRetryBeat: a weapon whose timer has run out while the target is out of reach retries the swing every half
// second, counted from the later of the timer running out and the melee attack turning on, and swings on the first
// retry after the server sees you in range (attackStartDelay after you get there). In range before the timer runs out,
// it swings as it runs out. Measured on Forever: probe walk-ins land within 0.02 s of a retry, a Strider Kick that
// passed its range check still waited for one, and 13 combat logs show the steps. False only for comparisons.
var ForeverSwingRetryBeat = true

const swingRetryPeriod = 500 * time.Millisecond

// retrySwingAt: when a weapon whose timer runs out at ready swings for a unit in range from now, its attack on since
// attackOn.
func retrySwingAt(now, ready, attackOn time.Duration) time.Duration {
	origin := max(ready, attackOn)
	seen := now + attackStartDelay
	if seen <= origin {
		return origin
	}
	steps := (seen - origin + swingRetryPeriod - 1) / swingRetryPeriod
	return origin + steps*swingRetryPeriod
}

// swingOnArrival sets each hand's next swing for a unit that just got in melee range after a walk that turned the melee
// attack on at attackOn (the off hand also waits ohWait from then), from the timers each hand had at range.
func (aa *AutoAttacks) swingOnArrival(sim *Simulation, mhReady, ohReady, walkStart, ohWait time.Duration) {
	attackOn, startsNow := aa.attackOnForWalkIn(sim, walkStart)
	// no Auto Shot since the last swing: the attack stayed on and only the timers count
	if !aa.shotSinceMelee {
		attackOn, startsNow = 0, false
	}
	// an attack turned on in melee range swings as soon as the server sees it there, with no retry
	swingAt := func(ready time.Duration) time.Duration {
		if startsNow {
			return max(ready, sim.CurrentTime+attackStartDelay)
		}
		return retrySwingAt(sim.CurrentTime, ready, attackOn)
	}
	if startsNow {
		aa.startAttackAt, aa.startAttackOn = attackOn, true
	}
	aa.mh.swingAt = swingAt(mhReady)
	sim.rescheduleWeaponAttack(aa.mh.swingAt)
	aa.mh.logSwingWaits(sim, mhReady, max(mhReady, attackOn))
	if aa.IsDualWielding {
		timerEnd := ohReady
		if ForeverOffHandWaitsOnAttackStart && aa.shotSinceMelee {
			ohReady = max(ohReady, attackOn+ohWait)
		}
		aa.oh.swingAt = swingAt(ohReady)
		sim.rescheduleWeaponAttack(aa.oh.swingAt)
		aa.oh.logSwingWaits(sim, timerEnd, max(ohReady, attackOn))
	}
	if startsNow {
		aa.mh.retryPending, aa.oh.retryPending = false, false
	}
}

// attackOnForWalkIn: when the melee attack turned on for a walk into melee that began at walkStart: at the rotation's
// Start Attack since the last Auto Shot, else as the walk began. A rotation that starts the attack itself
// (manualStartAttack) and has not done so turns it on now, on reaching melee (startsNow).
func (aa *AutoAttacks) attackOnForWalkIn(sim *Simulation, walkStart time.Duration) (time.Duration, bool) {
	manual := aa.mh.unit.Rotation != nil && aa.mh.unit.Rotation.manualStartAttack
	switch {
	case aa.startAttackOn && (manual || aa.startAttackAt < walkStart):
		return aa.startAttackAt, false
	case manual:
		return sim.CurrentTime, true
	}
	return walkStart, false
}

// StartAttack turns the melee attack on now (the rotation's Start Attack action), until the next Auto Shot.
func (aa *AutoAttacks) StartAttack(sim *Simulation) {
	aa.startAttackAt, aa.startAttackOn = sim.CurrentTime, true
	if sim.Log != nil && aa.mh.spell != nil {
		aa.mh.unit.Log(sim, "Melee attack on (Start Attack)")
	}
}

// MeleeAttackOff: no melee swing and no Start Attack since the last Auto Shot, or since the pull.
func (aa *AutoAttacks) MeleeAttackOff() bool {
	return aa.shotSinceMelee && !aa.startAttackOn
}

// logSwingWaits tells the timeline why a weapon that got in range now waits for its swing: out of melee range since its
// timer ran out (timerEnd), then the swing retry (counted from origin) until the swing.
func (wa *WeaponAttack) logSwingWaits(sim *Simulation, timerEnd, origin time.Duration) {
	if sim.Log == nil || wa.spell == nil {
		return
	}
	now := sim.CurrentTime
	if timerEnd < now-time.Millisecond {
		wa.unit.Log(sim, "Swing wait %s range from %.3fs to %.3fs", wa.timerID(), max(timerEnd, 0).Seconds(), now.Seconds())
	}
	if from := max(now, timerEnd); wa.swingAt > from+time.Millisecond {
		wa.retryPending, wa.retryFrom, wa.retryOrigin = true, from, origin
	}
}

// endRetryWait logs a pending swing retry wait up to now: as retry when the swing lands, as left when a walk takes the
// unit out before it.
func (wa *WeaponAttack) endRetryWait(sim *Simulation, how string) {
	if !wa.retryPending {
		return
	}
	wa.retryPending = false
	if sim.Log != nil && wa.spell != nil && sim.CurrentTime > wa.retryFrom+time.Millisecond {
		wa.unit.Log(sim, "Swing wait %s %s from %.3fs to %.3fs, steps from %.3fs", wa.timerID(), how, wa.retryFrom.Seconds(), sim.CurrentTime.Seconds(), wa.retryOrigin.Seconds())
	}
}

// logShotWaits tells the timeline why a ready Auto Shot waits when a walk ends at range: too close to shoot before the
// walk began, the walk itself, then the windup after stopping.
func (aa *AutoAttacks) logShotWaits(sim *Simulation, walkStart time.Duration, startDistance float64) {
	wa := &aa.ranged
	if sim.Log == nil || wa.spell == nil || !aa.AutoSwingRanged || wa.unit.DistanceFromTarget < MinRangedAttackDistance {
		return
	}
	ready, stop := max(wa.swingAt, 0), sim.CurrentTime
	if startDistance < MinRangedAttackDistance && ready < walkStart-time.Millisecond {
		wa.unit.Log(sim, "Swing wait %s close from %.3fs to %.3fs", wa.timerID(), ready.Seconds(), walkStart.Seconds())
	}
	if ready < stop-time.Millisecond {
		wa.unit.Log(sim, "Swing wait %s moving from %.3fs to %.3fs", wa.timerID(), max(ready, walkStart).Seconds(), stop.Seconds())
	}
	// the whole windup runs from the stop; it holds the shot only past the timer's ready (from)
	fire := stop + AutoShotStandStill
	from := min(max(ready, stop), fire)
	wa.unit.Log(sim, "Swing wait %s windup from %.3fs to %.3fs, stopped at %.3fs", wa.timerID(), from.Seconds(), fire.Seconds(), stop.Seconds())
}

// LogShotDelay tells the timeline a shot pushed the Auto Shot back (the Shots delay Auto Shot option).
func (aa *AutoAttacks) LogShotDelay(sim *Simulation, from, to time.Duration) {
	if sim.Log != nil && aa.ranged.spell != nil {
		aa.ranged.unit.Log(sim, "Swing wait %s shotdelay from %.3fs to %.3fs", aa.ranged.timerID(), from.Seconds(), to.Seconds())
	}
}

// SwingOnArrival: a walk from outside melee range that ended in it (see swingOnArrival); mhReady and ohReady are the
// timers from before the walk's end turned the swings back on.
func (aa *AutoAttacks) SwingOnArrival(sim *Simulation, mhReady, ohReady, attackOn, ohWait time.Duration) {
	if !aa.AutoSwingMelee || !aa.enabled || aa.mh.spell == nil {
		return
	}
	aa.swingOnArrival(sim, mhReady, ohReady, attackOn, ohWait)
}

// MeleeTimers: when each hand's timer runs out, for SwingOnArrival.
func (aa *AutoAttacks) MeleeTimers() (time.Duration, time.Duration) {
	return aa.mh.swingAt, aa.oh.swingAt
}

// StartMeleeMidMove starts the melee swings when a walk in reaches melee range; the walk turned the attack on as it
// started (attackOn), so the off hand keeps that wait.
func (aa *AutoAttacks) StartMeleeMidMove(sim *Simulation, attackOn time.Duration, wait time.Duration) {
	if !ForeverMeleeMidWalk || !aa.AutoSwingMelee || aa.enabled || aa.meleeWhileMoving || aa.mh.spell == nil {
		return
	}
	aa.meleeWhileMoving = true
	if ForeverSwingRetryBeat {
		aa.swingOnArrival(sim, aa.mh.swingAt, aa.oh.swingAt, attackOn, wait)
		aa.mh.addWeaponAttack(sim, aa.mh.unit.SwingSpeed())
		if aa.IsDualWielding {
			aa.oh.addWeaponAttack(sim, aa.mh.unit.SwingSpeed())
		}
		return
	}
	// the same 50 ms start-of-attack delay as at the end of a walk
	aa.mh.swingAt = max(aa.mh.swingAt, sim.CurrentTime) + 50*time.Millisecond
	aa.mh.addWeaponAttack(sim, aa.mh.unit.SwingSpeed())
	if aa.IsDualWielding {
		aa.oh.swingAt = max(aa.oh.swingAt, sim.CurrentTime) + 50*time.Millisecond
		aa.oh.addWeaponAttack(sim, aa.mh.unit.SwingSpeed())
		aa.OffHandWaitsForAttackStart(sim, attackOn, wait)
	}
}

// MeleeStartedMidMove: the walk that just ended already started the melee swings at 5 yd.
func (aa *AutoAttacks) MeleeStartedMidMove() bool {
	return aa.meleeWhileMoving
}

func (aa *AutoAttacks) stopMeleeMidMove(sim *Simulation) {
	aa.meleeWhileMoving = false
	sim.removeWeaponAttack(&aa.mh)
	if aa.IsDualWielding {
		sim.removeWeaponAttack(&aa.oh)
	}
}

// logOffHandWait writes log-only markers (no cast, no effects) so the timeline can draw the off hand's wait as a bar
// from the moment the melee attack turns on.
func (aa *AutoAttacks) logOffHandWait(sim *Simulation) {
	if !ForeverOffHandWaitsOnAttackStart || !aa.AutoSwingMelee || !aa.IsDualWielding || !aa.shotSinceMelee {
		return
	}
	waitID := ActionID{OtherID: proto.OtherAction_OtherActionAttack, Tag: 4}
	wait := aa.OffHandWait()
	aa.oh.unit.Log(sim, "Casting %s (Cost = 0.000, Cast Time = %s, Effective Time = 0s)", waitID, wait)
	StartDelayedAction(sim, DelayedActionOptions{
		DoAt: sim.CurrentTime + wait,
		OnAction: func(sim *Simulation) {
			aa.oh.unit.Log(sim, "Completed cast %s", waitID)
		},
	})
}

// OffHandWait is half the off hand's swing time at the current attack speed: its wait when the melee attack turns on now.
func (aa *AutoAttacks) OffHandWait() time.Duration {
	if !aa.IsDualWielding {
		return 0
	}
	return DurationFromSeconds(aa.oh.SwingSpeed/aa.oh.unit.SwingSpeed()) / 2
}

// OffHandWaitsForAttackStart holds the off hand's next swing until the wait that started when the melee attack turned on.
// The attack only turned off if Auto Shot fired: walking out and back in with the attack left on swings both hands on
// arrival.
func (aa *AutoAttacks) OffHandWaitsForAttackStart(sim *Simulation, attackOn time.Duration, wait time.Duration) {
	if !ForeverOffHandWaitsOnAttackStart || !(aa.enabled || aa.meleeWhileMoving) || !aa.AutoSwingMelee || !aa.IsDualWielding || !aa.shotSinceMelee {
		return
	}
	aa.oh.swingAt = max(aa.oh.swingAt, attackOn+wait)
	sim.rescheduleWeaponAttack(aa.oh.swingAt)
}
