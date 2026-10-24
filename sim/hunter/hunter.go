package hunter

import (
	"os"
	"strconv"
	"time"

	"github.com/wowsims/classic/sim/common/guardians"
	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
)

var TalentTreeSizes = [3]int{16, 16, 18}

const (
	SpellFlagShot   = core.SpellFlagAgentReserved1
	SpellFlagStrike = core.SpellFlagAgentReserved2
	SpellFlagSting  = core.SpellFlagAgentReserved3
	SpellFlagTrap   = core.SpellFlagAgentReserved4
)

const (
	SpellCode_HunterNone int32 = iota

	SpellCode_HunterAimedShot
	SpellCode_HunterArcaneShot
	SpellCode_HunterMultiShot

	SpellCode_HunterRaptorStrike
	SpellCode_HunterRaptorStrikeHit

	SpellCode_HunterSerpentSting

	SpellCode_HunterExplosiveTrap
	SpellCode_HunterFreezingTrap
	SpellCode_HunterImmolationTrap

	SpellCode_HunterMongooseBite
	SpellCode_HunterWingClip
	SpellCode_HunterVolley

	SpellCode_HunterPetClaw
	SpellCode_HunterPetBite
	SpellCode_HunterPetLightningBreath
	SpellCode_HunterPetScorpidPoison
)

func RegisterHunter() {
	core.RegisterAgentFactory(
		proto.Player_Hunter{},
		proto.Spec_SpecHunter,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewHunter(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Hunter)
			if !ok {
				panic("Invalid spec value for Hunter!")
			}
			player.Spec = playerSpec
		},
	)
}

type Hunter struct {
	core.Character

	Talents *proto.HunterTalents
	Options *proto.Hunter_Options

	pet *HunterPet

	AmmoDPS                   float64
	AmmoDamageBonus           float64
	NormalizedAmmoDamageBonus float64
	quiverBonus               float64

	AspectOfTheHawkAPMultiplier float64

	curQueueAura       *core.Aura
	raptorQueueAura    *core.Aura

	AimedShot       *core.Spell
	ArcaneShot      *core.Spell
	ExplosiveTrap   *core.Spell
	ImmolationTrap  *core.Spell
	FreezingTrap    *core.Spell
	KillCommand     *core.Spell
	MultiShot       *core.Spell
	RapidFire       *core.Spell
	RaptorStrike    *core.Spell
	RaptorStrikeHit *core.Spell
	MongooseBite    *core.Spell
	ScorpidSting    *core.Spell
	SerpentSting    *core.Spell
	SilencingShot   *core.Spell
	Volley          *core.Spell
	WingClip        *core.Spell
	SniperShot      *core.Spell
	StriderKick     *core.Spell
	SummonHawk      *core.Spell
	activeHawks     int32
	LaceratingBleed *core.Spell

	Shots       []*core.Spell
	Strikes     []*core.Spell
	MeleeSpells []*core.Spell
	LastShot    *core.Spell
	// a shot already pushed this Auto Shot cycle back (the Shots delay Auto Shot option)
	autoShotDelayed bool

	DefensiveState *core.Aura
	ExposePrey     *core.Aura
	// a landed Mongoose Bite is what ended this Expose Prey (for the used / expired markers)
	exposePreyUsed      bool
	RapidFireAura       *core.Aura
	BestialWrathPetAura *core.Aura
	IntimidationPetAura *core.Aura
	QuickShotsAura      *core.Aura
	rapidRecupAura      *core.Aura
	rapidRecupRegen     *core.CastingRegenSource
	resourcefulRegen    *core.CastingRegenSource
	ExposePreyAura      *core.Aura
	talentsApplied      bool
}

func (hunter *Hunter) GetCharacter() *core.Character {
	return &hunter.Character
}

func (hunter *Hunter) GetHunter() *Hunter {
	return hunter
}

func (hunter *Hunter) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
}

func (hunter *Hunter) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (hunter *Hunter) Initialize() {
	// Local experiment hook: AB_HASTE_PCT=1 adds +1% melee and ranged attack speed (e.g. a racial).
	if v, err := strconv.ParseFloat(os.Getenv("AB_HASTE_PCT"), 64); err == nil && v != 0 {
		hunter.PseudoStats.MeleeSpeedMultiplier *= 1 + v/100
		hunter.PseudoStats.RangedSpeedMultiplier *= 1 + v/100
	}
	hunter.registerAspectOfTheHawkSpell()
	hunter.registerAspectOfTheBeastSpell()
	hunter.registerAspectOfTheCheetahSpell()

	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagShot) {
			hunter.Shots = append(hunter.Shots, spell)
		}
	})
	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagStrike) {
			hunter.Strikes = append(hunter.Strikes, spell)
		}
	})
	hunter.OnSpellRegistered(func(spell *core.Spell) {
		if spell.ProcMask.Matches(core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial) {
			hunter.MeleeSpells = append(hunter.MeleeSpells, spell)
		}
	})

	multiShotTimer := hunter.NewTimer()
	arcaneShotTimer := hunter.NewTimer()

	hunter.registerSerpentStingSpell()
	hunter.registerArcaneShotSpell(arcaneShotTimer)
	hunter.registerAimedShotSpell(multiShotTimer)
	hunter.registerMultiShotSpell(multiShotTimer)
	hunter.registerSniperShotSpell()
	hunter.registerSummonHawkSpell(arcaneShotTimer)
	hunter.registerStriderKickSpell()
	hunter.registerMongooseBiteSpell()
	hunter.registerHuntersMarkSpell()
	hunter.registerWingClipSpell()
	hunter.registerVolleySpell()

	traps := hunter.NewTimer()
	hunter.registerExplosiveTrapSpell(traps)
	hunter.registerImmolationTrapSpell(traps)
	hunter.registerFreezingTrapSpell(traps)
	hunter.registerRapidFire()
	hunter.registerRaptorStrikeSpell()
	hunter.applyShotsDelayAutoShot()
}

// rangedHastedCast: a shot's cast time under the hunter's ranged haste (quiver, Rapid Fire and the rest).
func (hunter *Hunter) rangedHastedCast(base time.Duration) time.Duration {
	mult := hunter.PseudoStats.RangedSpeedMultiplier
	if mult < 0.01 {
		mult = 1
	}
	return time.Duration(float64(base) / mult)
}

// shotAutoShotDelay: how far a shot or sting cast in an Auto Shot cycle pushes that cycle's Auto Shot back (beta logs,
// standing still: +0.19 s middle with one cast in the cycle, +0.00 s with none)
const shotAutoShotDelay = 200 * time.Millisecond

// applyShotsDelayAutoShot: the "Shots delay Auto Shot" option (forever-bugs #43), once per Auto Shot cycle.
func (hunter *Hunter) applyShotsDelayAutoShot() {
	if !hunter.Options.ShotsDelayAutoShot {
		return
	}
	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: "Shots Delay Auto Shot",
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// the Auto Shot itself never reaches here (no cast-complete callbacks): its effects start the new cycle
			if hunter.autoShotDelayed || spell.DefenseType != core.DefenseTypeRanged || !spell.ProcMask.Matches(core.ProcMaskRangedSpecial) {
				return
			}
			hunter.autoShotDelayed = true
			due := hunter.AutoAttacks.NextRangedAttackAt()
			hunter.AutoAttacks.DelayRangedUntil(sim, due+shotAutoShotDelay)
			hunter.AutoAttacks.LogShotDelay(sim, due, due+shotAutoShotDelay)
		},
	}))
}

func (hunter *Hunter) Reset(sim *core.Simulation) {
	hunter.autoShotDelayed = false
	hunter.activeHawks = 0
}

func NewHunter(character *core.Character, options *proto.Player) *Hunter {
	hunterOptions := options.GetHunter()

	hunter := &Hunter{
		Character:   *character,
		Talents:     &proto.HunterTalents{},
		Options:     hunterOptions.Options,
		quiverBonus: 1.0,
	}
	core.FillTalentsProto(hunter.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)
	hunter.EnableManaBar()
	hunter.PseudoStats.CanParry = true
	hunter.MoveDelay = time.Duration(hunter.Options.MoveDelayMs) * time.Millisecond
	hunter.BaseRunSpeed = hunter.Options.RunSpeed
	hunter.SkyborneConvergence = hunter.Options.SkyborneConvergence

	rangedWeapon := hunter.WeaponFromRanged()

	if hunter.HasRangedWeapon() {
		switch hunter.Options.Ammo {
		case proto.Hunter_Options_RazorArrow:
			hunter.AmmoDPS = 7.5
		case proto.Hunter_Options_SolidShot:
			hunter.AmmoDPS = 7.5
		case proto.Hunter_Options_JaggedArrow:
			hunter.AmmoDPS = 13
		case proto.Hunter_Options_AccurateSlugs:
			hunter.AmmoDPS = 13
		case proto.Hunter_Options_MithrilGyroShot:
			hunter.AmmoDPS = 15
		case proto.Hunter_Options_IceThreadedArrow:
			hunter.AmmoDPS = 16.5
		case proto.Hunter_Options_IceThreadedBullet:
			hunter.AmmoDPS = 16.5
		case proto.Hunter_Options_ThoriumHeadedArrow:
			hunter.AmmoDPS = 17.5
		case proto.Hunter_Options_ThoriumShells:
			hunter.AmmoDPS = 17.5
		case proto.Hunter_Options_RockshardPellets:
			hunter.AmmoDPS = 18
		case proto.Hunter_Options_Doomshot:
			hunter.AmmoDPS = 20
		case proto.Hunter_Options_MiniatureCannonBalls:
			hunter.AmmoDPS = 20.5
		// Forever Engineering ammo
		case proto.Hunter_Options_SwiftfeatherArrow:
			hunter.AmmoDPS = 24.62
		case proto.Hunter_Options_SwiftstrikeShot:
			hunter.AmmoDPS = 24.62
		}
		hunter.AmmoDamageBonus = hunter.AmmoDPS * rangedWeapon.SwingSpeed
		hunter.NormalizedAmmoDamageBonus = hunter.AmmoDPS * 2.8

		switch hunter.Options.QuiverBonus {
		case proto.Hunter_Options_Speed10:
			hunter.quiverBonus = 1.10
		case proto.Hunter_Options_Speed11:
			hunter.quiverBonus = 1.11
		case proto.Hunter_Options_Speed12:
			hunter.quiverBonus = 1.12
		case proto.Hunter_Options_Speed13:
			hunter.quiverBonus = 1.13
		case proto.Hunter_Options_Speed14:
			hunter.quiverBonus = 1.14
		case proto.Hunter_Options_Speed15:
			hunter.quiverBonus = 1.15
		}
		if hunter.quiverBonus > 1 {
			hunter.PseudoStats.RangedSpeedMultiplier *= hunter.quiverBonus
		}
	}

	// Local experiment knob: AB_OH_SPEED=<s> gives the off-hand that speed at the same damage per second.
	offHand := hunter.WeaponFromOffHand()
	if v, err := strconv.ParseFloat(os.Getenv("AB_OH_SPEED"), 64); err == nil && v > 0 && offHand.SwingSpeed > 0 {
		k := v / offHand.SwingSpeed
		offHand.BaseDamageMin *= k
		offHand.BaseDamageMax *= k
		offHand.SwingSpeed = v
	}

	// AB_MH_SPEED=<s>: same for the main hand (specials keep their normalization rules).
	mainHand := hunter.WeaponFromMainHand()
	if v, err := strconv.ParseFloat(os.Getenv("AB_MH_SPEED"), 64); err == nil && v > 0 && mainHand.SwingSpeed > 0 {
		k := v / mainHand.SwingSpeed
		mainHand.BaseDamageMin *= k
		mainHand.BaseDamageMax *= k
		mainHand.SwingSpeed = v
	}

	hunter.EnableAutoAttacks(hunter, core.AutoAttackOptions{
		MainHand:        mainHand,
		OffHand:         offHand,
		Ranged:          rangedWeapon,
		ReplaceMHSwing:  hunter.TryRaptorStrike,
		AutoSwingRanged: true,
		AutoSwingMelee:  true,
	})

	if hunter.Options.AutoShotSpamBug {
		hunter.AutoAttacks.RangedCycleCut = AutoShotSpamCut
	}
	hunter.AutoAttacks.MeleeResetsAutoShot = hunter.Options.MeleeResetsAutoShot
	hunter.AutoAttacks.RangedConfig().Flags |= core.SpellFlagCastTimeNoGCD
	hunter.AutoAttacks.RangedConfig().Flags |= core.SpellFlagMeleeMetrics
	// Mortal Shots covers Auto Shot: the client data mask, and Arc's log (2026-10-21, 5/5: Auto Shot crits 2.30-2.32x a
	// normal hit, as Arcane Shot's). The option keeps the earlier report that it does not, for comparisons.
	hunter.AutoAttacks.RangedConfig().CritDamageBonus = hunter.mortalShots()
	if hunter.Options.MortalShotsSkipsAutoShot {
		hunter.AutoAttacks.RangedConfig().CritDamageBonus = 0
	}
	hunter.AutoAttacks.RangedConfig().BonusCoefficient = 1
	// Auto Shot's projectile flies at 40 yd/s (client SpellMisc, spell 75), as Arcane Shot's, Serpent Sting's and Aimed
	// Shot's; Multi-Shot's at 30
	hunter.AutoAttacks.RangedConfig().MissileSpeed = 40
	hunter.AutoAttacks.RangedConfig().ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		// the Auto Shot takes you out of melee mode too
		hunter.unqueueRaptorStrike(sim)
		// a new Auto Shot cycle: a shot in it may push the next Auto Shot back again (Shots delay Auto Shot)
		hunter.autoShotDelayed = false
		baseDamage := hunter.RangedWeaponDamage(sim, spell.RangedAttackPower(target, false)) +
			hunter.AmmoDamageBonus
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
		spell.WaitTravelTime(sim, func(sim *core.Simulation) {
			spell.DealDamage(sim, result)
		})
	}

	if hunter.Talents.LoneWolf {
		hunter.pet = nil
		hunter.PseudoStats.DamageDealtMultiplier *= 1.20
	} else {
		hunter.pet = hunter.NewHunterPet()
	}

	hunter.AddStatDependency(stats.Strength, stats.AttackPower, core.APPerStrength[character.Class])
	hunter.AddStatDependency(stats.Agility, stats.AttackPower, 1)
	hunter.AddStatDependency(stats.Agility, stats.RangedAttackPower, 2)
	hunter.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiAtLevel[character.Class]*core.CritRatingPerCritChance)
	hunter.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[character.Class]*core.SpellCritRatingPerCritChance)

	guardians.ConstructGuardians(&hunter.Character)
	return hunter
}

func (hunter *Hunter) OnGCDReady(_ *core.Simulation) {
}

// AutoShotSpamCut: spamming Auto Shot in game gave 2.07 s between shots at best on a 2.52 s cycle (median 2.16, the
// rest is press timing); fast spam lands near the best. A melee restart is not shortened: no log shows it.
const AutoShotSpamCut = 450 * time.Millisecond

type HunterAgent interface {
	GetHunter() *Hunter
}
