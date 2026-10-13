package mage

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

var TalentTreeSizes = [3]int{18, 17, 19}

type Mage struct {
	core.Character

	Talents *proto.MageTalents
	Options *proto.MageOptions

	ArcaneBlastAura    *core.Aura
	ArcanePowerAura    *core.Aura
	ClearcastingAura   *core.Aura
	FingersOfFrostAura *core.Aura
	HotStreakAura      *core.Aura
	ImprovedScorchAura *core.Aura
	MissileBarrageAura *core.Aura
	PresenceOfMindAura *core.Aura
	WintersChillAura   *core.Aura

	// procFingersOfFrost rolls Fingers of Frost for a chilling bolt at launch; nil without the talent.
	procFingersOfFrost func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult)

	ArcaneBlast *core.Spell
	Ignite      *core.Spell
	Flamestrike []*core.Spell
}

func (mage *Mage) GetCharacter() *core.Character {
	return &mage.Character
}

func (mage *Mage) GetMage() *Mage {
	return mage
}

func RegisterMage() {
	core.RegisterAgentFactory(
		proto.Player_Mage{},
		proto.Spec_SpecMage,
		func(character *core.Character, options *proto.Player, _ *proto.Raid) core.Agent {
			return NewMage(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Mage)
			if !ok {
				panic("Invalid spec value for Mage!")
			}
			player.Spec = playerSpec
		},
	)
}

func (mage *Mage) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.ArcaneBrilliance = true
}

func (mage *Mage) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (mage *Mage) Initialize() {
	mage.registerPassives()
	mage.registerSpells()
}

func (mage *Mage) registerPassives() {
	mage.registerArcaneCharges()
}

func (mage *Mage) registerSpells() {
	mage.registerArcaneBlastSpell()
	mage.registerArcaneExplosionSpell()
	mage.registerArcaneMissilesSpell()
	mage.registerArmorSpells()
	mage.registerBlizzardSpell()
	mage.registerConeOfColdSpell()
	mage.registerFrostboltSpell()
	mage.registerEvocation()
	mage.registerFireballSpell()
	mage.registerFrostfireBoltSpell()
	mage.registerFireBlastSpell()
	mage.registerFrostNovaSpell()
	mage.registerIceLanceSpell()
	mage.registerManaGems()
	mage.registerScorchSpell()

	FlameStrikeRankMap.Each(func(_ int32, rankConfig *spelldata.Spell) { mage.registerFlamestrike(rankConfig) })

	//TalentSpells
	mage.registerPresenceOfMindSpell()
	mage.registerArcanePowerSpell()

	mage.registerBlastWaveSpell()
	mage.registerPyroblastSpell()
	mage.registerCombustionSpell()

	mage.registerColdSnapSpell()
}

func (mage *Mage) Reset(sim *core.Simulation) {
}

func (mage *Mage) OnEncounterStart(sim *core.Simulation) {
}

func NewMage(character *core.Character, options *proto.Player) *Mage {
	mageOptions := options.GetMage().Options.ClassOptions
	mage := &Mage{
		Character: *character,
		Talents:   &proto.MageTalents{},
		Options:   mageOptions,
	}

	core.FillTalentsProto(mage.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)

	mage.EnableManaBar()
	mage.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[character.Class])

	// Forever has no Water Elemental: client 70009 carries no Summon Water Elemental (31687) at all,
	// so the pet is never created.

	return mage
}

// Agent is a generic way to access underlying mage on any of the agents.
type MageAgent interface {
	GetMage() *Mage
}

const (
	FireSpellMaxTimeUntilResult       = 750 * time.Millisecond
	MageSpellFlagNone           int64 = 0
	MageSpellArcaneBlast        int64 = 1 << iota
	MageSpellArcaneExplosion
	MageSpellArcanePower
	MageSpellArcaneMissilesCast
	MageSpellArcaneMissilesTick
	MageSpellBlastWave
	MageSpellBlizzard
	MageSpellColdSnap
	MageSpellConeOfCold
	MageSpellEvocation
	MageSpellFireBlast
	MageSpellFireball
	MageSpellFlamestrike
	MageSpellFlamestrikeDot
	MageSpellFrostArmor
	MageSpellFrostbolt
	MageSpellFrostNova
	MageSpellIceBarrier
	MageSpellIceBlock
	MageSpellIceLance
	MageSpellIgnite
	MageSpellMageArmor
	MageSpellManaGems
	MageSpellMoltenArmor
	MageSpellPresenceOfMind
	MageSpellPyroblast
	MageSpellPyroblastDot
	MageSpellScorch
	MageSpellManaGem
	MageSpellCombustion
	MageSpellImprovedBlizzard
	MageSpellFrostfireBolt

	MageSpellLast
	MageSpellsAll  = MageSpellLast<<1 - 1
	MageSpellFrost = MageSpellFrostbolt | MageSpellBlizzard | MageSpellFrostNova | MageSpellConeOfCold | MageSpellIceLance
	MageSpellFire  = MageSpellFireball | MageSpellCombustion |
		MageSpellFireBlast | MageSpellFlamestrike | MageSpellIgnite | MageSpellPyroblast | MageSpellScorch
	MageSpellsAllDamaging = MageSpellArcaneBlast | MageSpellArcaneExplosion | MageSpellArcaneMissilesTick | MageSpellBlizzard |
		MageSpellFireBlast | MageSpellFireball | MageSpellFlamestrike | MageSpellFrostbolt |
		MageSpellIceLance | MageSpellPyroblast | MageSpellPyroblastDot | MageSpellScorch |
		MageSpellBlastWave | MageSpellConeOfCold | MageSpellFrostNova | MageSpellFrostfireBolt
	MageSpellInstantCast = MageSpellArcaneMissilesCast | MageSpellArcaneMissilesTick | MageSpellFireBlast | MageSpellArcaneExplosion | MageSpellPyroblastDot |
		MageSpellCombustion | MageSpellConeOfCold | MageSpellIceLance | MageSpellManaGems | MageSpellPresenceOfMind
	MageSpellExtraResult    = MageSpellArcaneMissilesTick | MageSpellBlizzard
	FireSpellIgnitable      = MageSpellFireball | MageSpellScorch | MageSpellPyroblast
	MageSpellArcaneMissiles = MageSpellArcaneMissilesCast | MageSpellArcaneMissilesTick

	// The chill effects Fingers of Frost rolls on: Frostbolt's slow, Cone of Cold's, Frostfire Bolt's
	// (all three carry the client's chill bit, 0x100000) and Improved Blizzard's.
	MageSpellChill = MageSpellFrostbolt | MageSpellConeOfCold | MageSpellFrostfireBolt | MageSpellImprovedBlizzard
)

// Arcane Explosion, Blast Wave, Cone of Cold and Frost Nova reach the targets within 10 yards of the
// mage and state no cast range, so the radius is their range. Client 1.60.1.70245: Cone of Cold,
// Blast Wave and Frost Nova state 10 yd on their area effects; Arcane Explosion's row carries no
// radius, its Era one 10 yd.
const pointBlankRadius = 10
