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
	curQueuedAutoSpell *core.Spell
	raptorFromSwing    bool
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

	DefensiveState      *core.Aura
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
	hunter.registerWingClipSpell()
	hunter.registerVolleySpell()

	traps := hunter.NewTimer()
	hunter.registerExplosiveTrapSpell(traps)
	hunter.registerImmolationTrapSpell(traps)
	hunter.registerFreezingTrapSpell(traps)
	hunter.registerRapidFire()
	hunter.registerRaptorStrikeSpell()
	hunter.applyMeleeResetsRangedAuto()
}

func (hunter *Hunter) Reset(sim *core.Simulation) {
	hunter.activeHawks = 0
	hunter.raptorFromSwing = false
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
	hunter.AutoAttacks.RangedConfig().Flags |= core.SpellFlagCastTimeNoGCD
	hunter.AutoAttacks.RangedConfig().Flags |= core.SpellFlagMeleeMetrics
	// Mortal Shots does not touch Auto Shot in game (tested); the client data mask says it would.
	hunter.AutoAttacks.RangedConfig().CritDamageBonus = 0
	hunter.AutoAttacks.RangedConfig().BonusCoefficient = 1
	hunter.AutoAttacks.RangedConfig().ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		// the Auto Shot takes you out of melee mode too
		hunter.unqueueRaptorStrike(sim)
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

func (hunter *Hunter) applyMeleeResetsRangedAuto() {
	hunter.RegisterAura(core.Aura{
		Label:    "Melee Resets Ranged Auto",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			id := spell.ActionID.SpellID
			if id == 1317257 || id == 1495 || id == 14269 || id == 14270 || id == 14271 {
				return
			}
			if spell.SpellCode == SpellCode_HunterMongooseBite {
				return
			}

			isWhite := spell.ProcMask.Matches(core.ProcMaskMeleeMHAuto)
			isRaptor := spell.SpellCode == SpellCode_HunterRaptorStrikeHit

			if !isWhite && !isRaptor {
				return
			}

			hunter.AutoAttacks.RestartRangedSwing(sim)
		},
	})
}

// AutoShotSpamCut: spamming Auto Shot in game gave 2.07 s between shots at best on a 2.52 s cycle (median 2.16, the
// rest is press timing); fast spam lands near the best. A melee restart is not shortened: no log shows it.
const AutoShotSpamCut = 450 * time.Millisecond

type HunterAgent interface {
	GetHunter() *Hunter
}
