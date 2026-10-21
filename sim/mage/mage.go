package mage

import (
	"fmt"
	"strings"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

// Arcane, Fire and Frost. Frost has 18 talents since client hotfix push 112753 (20 October 2026), which
// removed Piercing Ice (patch 187).
var TalentTreeSizes = [3]int{19, 17, 18}

// checkTalentsString (MythicSim patch 189) rejects a talent string that gives a tree more digits than
// the tree has talents, or more than three trees. core.FillTalentsProto would carry the extra digits
// into the next tree's fields, or panic on a field the proto does not have: a Frost string from before
// hotfix push 112753 has 19 digits, and its 19th asks for field 55, which is gone. Its first 18 would
// also be read in the new order (Frost Channeling where Piercing Ice was), so it is refused as a whole
// rather than trimmed. A pre-hotfix string with 18 or fewer Frost digits cannot be told apart from a
// current one and is read in the current order.
func checkTalentsString(talents string) {
	trees := strings.Split(talents, "-")
	if len(trees) > len(TalentTreeSizes) {
		panic(fmt.Sprintf("mage talents %q: %d trees, a mage has %d", talents, len(trees), len(TalentTreeSizes)))
	}
	for i, tree := range trees {
		if len(tree) > TalentTreeSizes[i] {
			panic(fmt.Sprintf("mage talents %q: tree %d has %d talents, the string gives %d (a Frost string from before"+
				" the 20 October 2026 hotfix has 19)", talents, i+1, TalentTreeSizes[i], len(tree)))
		}
	}
}

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
	mage.registerFrostfireBolt()
	mage.registerEvocation()
	mage.registerFireballSpell()
	mage.registerFireBlastSpell()
	mage.registerFrostNovaSpell()
	mage.registerIceLanceSpell()
	mage.registerManaGems()
	mage.registerScorchSpell()

	FlameStrikeRankMap.Each(func(_ int32, rankConfig *spelldata.Spell) { mage.registerFlamestrike(rankConfig) })

	//TalentSpells
	mage.registerPresenceOfMindSpell()
	mage.registerArcanePowerSpell()
	mage.registerArcaneWrathSpell()
	mage.registerFrostBlastSpell()

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

	checkTalentsString(options.TalentsString)
	core.FillTalentsProto(mage.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)
	// Elemental Precision went from five ranks to three in the same hotfix; a pre-hotfix 4 or 5 would
	// otherwise panic in its ladder with no word of the talent.
	if maxRanks := spellData.ElementalPrecision.Len(); mage.Talents.ElementalPrecision > maxRanks {
		panic(fmt.Sprintf("mage talents %q: Elemental Precision has %d ranks, the string gives %d (five before"+
			" the 20 October 2026 hotfix)", options.TalentsString, maxRanks, mage.Talents.ElementalPrecision))
	}

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

	// Forever baseline spell; keep its mask stable for existing spells.
	MageSpellFrostfireBolt
	// Arcane talent of client 70425 (16 October 2026 beta).
	MageSpellArcaneWrath
	// Frost Blast, a trained Frost spell of client 70425 (MythicSim patch 153).
	MageSpellFrostBlast

	MageSpellLast
	MageSpellsAll  = MageSpellLast<<1 - 1
	MageSpellFrost = MageSpellFrostfireBolt | MageSpellFrostbolt | MageSpellBlizzard | MageSpellFrostNova | MageSpellConeOfCold | MageSpellIceLance |
		MageSpellFrostBlast
	MageSpellFire = MageSpellFrostfireBolt | MageSpellFireball | MageSpellCombustion |
		MageSpellFireBlast | MageSpellFlamestrike | MageSpellIgnite | MageSpellPyroblast | MageSpellScorch
	MageSpellsAllDamaging = MageSpellArcaneBlast | MageSpellArcaneExplosion | MageSpellArcaneMissilesTick | MageSpellBlizzard |
		MageSpellFireBlast | MageSpellFireball | MageSpellFlamestrike | MageSpellFrostbolt |
		MageSpellIceLance | MageSpellPyroblast | MageSpellPyroblastDot | MageSpellScorch |
		MageSpellBlastWave | MageSpellConeOfCold | MageSpellFrostNova | MageSpellFrostfireBolt | MageSpellArcaneWrath |
		MageSpellFrostBlast
	MageSpellInstantCast = MageSpellArcaneMissilesCast | MageSpellArcaneMissilesTick | MageSpellFireBlast | MageSpellArcaneExplosion | MageSpellPyroblastDot |
		MageSpellCombustion | MageSpellConeOfCold | MageSpellIceLance | MageSpellManaGems | MageSpellPresenceOfMind
	MageSpellExtraResult    = MageSpellArcaneMissilesTick | MageSpellBlizzard
	FireSpellIgnitable      = MageSpellFrostfireBolt | MageSpellFireball | MageSpellScorch | MageSpellPyroblast
	MageSpellArcaneMissiles = MageSpellArcaneMissilesCast | MageSpellArcaneMissilesTick

	// The chill effects Fingers of Frost rolls on: Frostbolt's slow, Cone of Cold's, Frostfire Bolt's,
	// Frost Blast's (all four carry the client's chill bit, 0x100000) and Improved Blizzard's.
	MageSpellChill = MageSpellFrostbolt | MageSpellConeOfCold | MageSpellFrostfireBolt | MageSpellImprovedBlizzard |
		MageSpellFrostBlast
)

// Arcane Explosion, Blast Wave, Cone of Cold and Frost Nova reach the targets within 10 yards of the
// mage and state no cast range, so the radius is their range. Client 1.60.1.70245: Cone of Cold,
// Blast Wave and Frost Nova state 10 yd on their area effects; Arcane Explosion's row carries no
// radius, its Era one 10 yd.
const pointBlankRadius = 10
