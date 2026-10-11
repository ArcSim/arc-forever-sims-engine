package druid

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

var TalentTreeSizes = [3]int{16, 20, 16}

type Druid struct {
	core.Character
	SelfBuffs

	Talents *proto.DruidTalents

	StartingForm DruidForm

	CannotShredTarget bool

	WolfsheadShiftingPowerEnergy float64 // Wolfshead Helm (8345): +20 energy from Shifting Power
	WolfsheadEnrageRage          float64 // Wolfshead Helm (8345): +5 rage from Enrage

	MHAutoSpell *core.Spell

	Barkskin             *DruidSpell
	Claw                 *DruidSpell
	Dash                 *DruidSpell
	DemoralizingRoar     *DruidSpell
	FaerieFire           *DruidSpell
	FaerieFireFeral      *DruidSpell
	FerociousBite        *DruidSpell
	Enrage               *DruidSpell
	EnrageAura           *core.Aura
	FrenziedRegeneration *DruidSpell
	Hurricane            *DruidSpell
	Innervate            *DruidSpell
	InsectSwarm          *DruidSpell
	Lacerate             *DruidSpell
	PrimalBite           *DruidSpell
	Maul                 *DruidSpell
	Moonfire             *DruidSpell
	NaturesSwiftness     *DruidSpell
	Prowl                *DruidSpell
	Rake                 *DruidSpell
	Ravage               *DruidSpell
	Rejuvenation         *DruidSpell
	Rip                  *DruidSpell
	Shred                *DruidSpell
	Starfire             []*DruidSpell
	ShiftingPower        *DruidSpell
	Swipe                *DruidSpell
	Wrath                *DruidSpell

	CatForm     *DruidSpell
	BearForm    *DruidSpell
	MoonkinForm *DruidSpell

	BearFormAura             *core.Aura
	CatFormAura              *core.Aura
	ClearcastingAura         *core.Aura
	DashAura                 *core.Aura
	FrenziedRegenerationAura *core.Aura
	DemoralizingRoarAuras    core.AuraArray
	FaerieFireAuras          core.AuraArray
	MangleAuras              core.AuraArray // always nil: Forever's Mangle has no debuff; the legacy cat rotation reads it as inactive
	BerserkAura              *core.Aura
	MoonkinFormAura          *core.Aura
	EclipseAura              *core.Aura
	ProwlAura                *core.Aura

	form DruidForm

	IntensityEnrageRageBonus float64

	// Furor: chance to gain energy/rage when shifting into Cat/Bear Form.
	FurorProcChance float64

	// Maul queue (fires on next auto-attack swing, like warrior Heroic Strike)
	maulQueueAura  *core.Aura
	maulQueueSpell *core.Spell
	maulRealismICD *core.Cooldown

	// Forever's Furor carries Energy across a powershift, so the shift needs to know how much
	// was left behind and when.
	lastCatFormEnergy float64
	lastCatFormExitAt time.Duration
}

const (
	DruidSpellFlagNone        int64 = 0
	DruidSpellEntanglingRoots int64 = 1 << iota
	DruidSpellClaw
	DruidSpellDemoralizingRoar
	DruidSpellFaerieFire
	DruidSpellFaerieFireFeral
	DruidSpellHurricane
	DruidSpellFerociousBite
	DruidSpellFrenziedRegeneration
	DruidSpellInnervate
	DruidSpellInsectSwarm
	DruidSpellLacerate
	DruidSpellPrimalBite
	DruidSpellMaul
	DruidSpellMoonfireInitial
	DruidSpellMoonfireDoT
	DruidSpellRake
	DruidSpellRavage
	DruidSpellRip
	DruidSpellShred
	DruidSpellStarfire
	DruidSpellSwipe
	DruidSpellThorns
	DruidSpellWrath
	DruidSpellEnrage
	DruidSpellShiftingPower
	DruidSpellCatForm
	DruidSpellBearForm
	DruidSpellMoonkinForm

	DruidSpellHealingTouch
	DruidSpellRegrowth
	DruidSpellLifebloom
	DruidSpellRejuvenation
	DruidSpellTranquility
	DruidSpellMarkOfTheWild
	DruidSpellSwiftmend
	DruidSpellCenarionWard

	// TODO: Forever abilities the sim does not model yet; see the stub file named for each.
	DruidSpellRevive

	DruidSpellLast
	DruidSpellsAll = DruidSpellLast<<1 - 1

	DruidSpellMoonfire           = DruidSpellMoonfireInitial | DruidSpellMoonfireDoT
	DruidSpellDoT                = DruidSpellMoonfireDoT | DruidSpellInsectSwarm
	DruidSpellHoT                = DruidSpellRejuvenation | DruidSpellLifebloom | DruidSpellRegrowth
	DruidSpellInstant            = DruidSpellMoonfire | DruidSpellFaerieFire
	DruidSpellBuilder            = DruidSpellClaw | DruidSpellPrimalBite | DruidSpellShred | DruidSpellRake | DruidSpellRavage
	DruidSpellFinisher           = DruidSpellFerociousBite | DruidSpellRip
	DruidArcaneSpells            = DruidSpellMoonfire | DruidSpellMoonfireDoT | DruidSpellStarfire
	DruidNatureSpells            = DruidSpellWrath | DruidSpellHurricane | DruidSpellInsectSwarm
	DruidHealingNonInstantSpells = DruidSpellHealingTouch | DruidSpellRegrowth
	DruidHealingSpells           = DruidHealingNonInstantSpells | DruidSpellRejuvenation | DruidSpellLifebloom | DruidSpellSwiftmend
	DruidDamagingSpells          = DruidArcaneSpells | DruidNatureSpells
)

type SelfBuffs struct {
	InnervateTarget *proto.UnitReference
}

func (druid *Druid) GetCharacter() *core.Character {
	return &druid.Character
}

func (druid *Druid) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
	if druid.InForm(Cat|Bear) && druid.Talents.LeaderOfThePack {
		partyBuffs.LeaderOfThePack = true
	} else if druid.InForm(Moonkin) && druid.Talents.MoonkinForm {
		partyBuffs.MoonkinAura = true
	}
}

func (druid *Druid) RegisterSpell(formMask DruidForm, config core.SpellConfig) *DruidSpell {
	prev := config.ExtraCastCondition
	prevModify := config.Cast.ModifyCast

	ds := &DruidSpell{FormMask: formMask}
	config.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		// Check if we're in allowed form to cast
		// Allow 'humanoid' auto unshift casts
		if (ds.FormMask != Any && !druid.InForm(ds.FormMask)) && !ds.FormMask.Matches(Humanoid) {
			if sim.Log != nil {
				sim.Log("Failed cast to spell %s, wrong form", ds.ActionID)
			}
			return false
		}
		return prev == nil || prev(sim, target)
	}
	config.Cast.ModifyCast = func(sim *core.Simulation, s *core.Spell, c *core.Cast) {
		if !druid.InForm(ds.FormMask) && ds.FormMask.Matches(Humanoid) {
			druid.ClearForm(sim)
		}
		if prevModify != nil {
			prevModify(sim, s, c)
		}
		// A hard cast restarts the swing timer as it completes (core.AutoAttacks.HoldMeleeForCast), as
		// for a caster-form druid swinging a staff. A druid that doesn't auto attack is untouched.
		if castTime := druid.ApplyCastSpeedForSpell(c.CastTime, s); castTime > 0 {
			druid.AutoAttacks.HoldMeleeForCast(sim, sim.CurrentTime+castTime)
		}
	}

	ds.Spell = druid.Unit.RegisterSpell(config)

	return ds
}

func (druid *Druid) Initialize() {
	druid.form = druid.StartingForm

	druid.Env.RegisterPostFinalizeEffect(func() {
		druid.MHAutoSpell = druid.AutoAttacks.MHAuto()
	})

	druid.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand}, func(sim *core.Simulation, slot proto.ItemSlot) {
		switch {
		case druid.InForm(Cat):
			druid.AutoAttacks.SetMH(druid.GetCatWeapon())
		case druid.InForm(Bear):
			druid.AutoAttacks.SetMH(druid.GetBearWeapon())
		}
	})

	druid.RegisterBaselineSpells()
}

func (druid *Druid) RegisterBaselineSpells() {
	druid.registerInnervateCD()
	druid.registerThornsSpell()
	druid.registerFormBreakingConsumes()
}

// registerFormBreakingConsumes patches ApplyEffects on potions, conjured items,
// and engineering explosives to drop Bear/Cat form when used. These spells all
// carry SpellFlagNoOnCastComplete, so OnCastComplete aura hooks never fire for
// them - we must wrap ApplyEffects directly instead.
func (druid *Druid) registerFormBreakingConsumes() {
	druid.Env.RegisterPostFinalizeEffect(func() {
		breakFlags := core.SpellFlagPotion | core.SpellFlagConjured | core.SpellFlagExplosive
		for _, spell := range druid.Spellbook {
			if !spell.Flags.Matches(breakFlags) {
				continue
			}
			prev := spell.ApplyEffects
			spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, sp *core.Spell) {
				prev(sim, target, sp)
				if druid.InForm(Bear) || druid.InForm(Cat) {
					druid.ClearForm(sim)
				}
			}
			// Nobody pots out of form on a timer: left to fire on cooldown, Mighty Rage Potion
			// dropped the bear out of form every two minutes and cut it from 410 DPS to 153.
			// The automatic use waits for caster form; a rotation can still drink one in form on
			// purpose, as the cat's powershifting does.
			druid.AddActivationCondition(spell.ActionID, func(_ *core.Simulation, _ *core.Character) bool {
				return !druid.InForm(Bear) && !druid.InForm(Cat)
			})
		}
	})
}

func (druid *Druid) RegisterBalanceSpells() {
	StarfireRankMap.Each(func(_ int32, r *spelldata.Spell) { druid.registerStarfireSpell(r) })
	druid.registerMoonfireSpell()
	spellData.Wrath.Each(func(_ int32, r *spelldata.Spell) { druid.registerWrathSpell(r) })
	druid.registerHurricaneSpell()
	druid.registerFaerieFireSpell()
}

func (druid *Druid) RegisterFeralCatSpells() {
	druid.registerCatFormSpell()

	// Forever has no Cat-form Mangle.
	druid.registerRakeSpell()
	druid.registerRipSpell()
	druid.registerFerociousBiteSpell()
	// Forever drops Faerie Fire (Feral); the Balance version is the only one.
	druid.registerFaerieFireSpell()
	druid.registerShredSpell()
	druid.registerClawSpell()
	druid.registerProwlSpell()
	druid.registerRavageSpell()
	druid.registerShiftingPowerSpell()
	// A cat that leaves its form to refresh Moonfire (the Furor powershift with a Moonfire cast as
	// the form drop). The spell is castable in caster form only, so casting it from Cat Form
	// unshifts first, as every other caster spell does here.
	druid.registerMoonfireSpell()
}

func (druid *Druid) RegisterFeralTankSpells() {
	druid.registerBearFormSpell()
	druid.registerBarkskin()
	druid.registerDemoralizingRoarSpell()
	// Forever drops Faerie Fire (Feral); the Balance version is the only one.
	druid.registerFaerieFireSpell()
	druid.registerEnrageSpell()
	druid.registerFrenziedRegenerationSpell()
	druid.registerLacerateSpell()
	druid.registerPrimalBiteSpell()
	druid.registerMaulSpell()
	druid.registerSwipeBearSpell()
}

func (druid *Druid) Reset(_ *core.Simulation) {
	druid.form = druid.StartingForm
}

func (druid *Druid) OnEncounterStart(sim *core.Simulation) {
}

func New(char *core.Character, form DruidForm, selfBuffs SelfBuffs, talents string) *Druid {
	druid := &Druid{
		Character:    *char,
		SelfBuffs:    selfBuffs,
		Talents:      &proto.DruidTalents{},
		StartingForm: form,
		form:         form,
	}

	core.FillTalentsProto(druid.Talents.ProtoReflect(), talents, TalentTreeSizes)
	druid.EnableManaBar()

	// Two attack power a point of Strength in every form, as on master (a Classic druid; the 1 was
	// carried over from the MoP-era engine).
	druid.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	druid.AddStatDependency(stats.BonusArmor, stats.Armor, 1)
	druid.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[char.Class])
	// Dodge is Classic's at level 60, as on master: 0.9% base and 20 Agility a percent (the same rate
	// as crit). The 14.7 Agility and -1.87% were TBC's level 70 fit.
	druid.AddStatDependency(stats.Agility, stats.DodgeRating, core.CritPerAgiMaxLevel[char.Class]*core.DodgeRatingPerDodgePercent)
	druid.PseudoStats.BaseDodgeChance += 0.009

	return druid
}

type DruidSpell struct {
	*core.Spell
	FormMask DruidForm

	// Optional fields used in snapshotting calculations
	CurrentSnapshotPower float64
	NewSnapshotPower     float64
	ShortName            string
}

func (ds *DruidSpell) IsReady(sim *core.Simulation) bool {
	if ds == nil {
		return false
	}
	return ds.Spell.IsReady(sim)
}

func (ds *DruidSpell) CanCast(sim *core.Simulation, target *core.Unit) bool {
	if ds == nil {
		return false
	}
	return ds.Spell.CanCast(sim, target)
}

func (ds *DruidSpell) IsEqual(s *core.Spell) bool {
	if ds == nil || s == nil {
		return false
	}
	return ds.Spell == s
}

func (druid *Druid) UpdateBleedPower(bleedSpell *DruidSpell, sim *core.Simulation, target *core.Unit, updateCurrent bool, updateNew bool) {
	snapshotPower := bleedSpell.ExpectedTickDamage(sim, target)

	if updateCurrent {
		bleedSpell.CurrentSnapshotPower = snapshotPower

		if sim.Log != nil {
			druid.Log(sim, "%s Snapshot Power: %.1f", bleedSpell.ShortName, snapshotPower)
		}
	}

	if updateNew {
		bleedSpell.NewSnapshotPower = snapshotPower

		if (sim.Log != nil) && !updateCurrent {
			druid.Log(sim, "%s Projected Power: %.1f", bleedSpell.ShortName, snapshotPower)
		}
	}
}

// Agent is a generic way to access underlying druid on any of the agents (for example balance druid.)
type DruidAgent interface {
	GetDruid() *Druid
}

// The tick outcome the family tables' shared.PeriodicTickOutcome picked: a crit roll where the row
// marks Periodic Can Crit, a plain tick otherwise, and never a hit roll. The store's TickOutcome rolls
// hit on magic ticks, which would move every dot, so the old choice is kept here.
func periodicTickOutcome(rank *spelldata.Spell, dot *core.Dot) core.OutcomeApplier {
	magic := rank.DefenseTypeCore() == core.DefenseTypeMagic
	switch {
	case rank.PeriodicCanCrit() && magic:
		return dot.Spell.OutcomeTickMagicCrit
	case rank.PeriodicCanCrit():
		return dot.Spell.OutcomeTickPhysicalCrit
	default:
		return dot.OutcomeTick
	}
}
