package feralcat

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/druid"
)

func RegisterFeralCatDruid() {
	core.RegisterAgentFactory(
		proto.Player_FeralCatDruid{},
		proto.Spec_SpecFeralCatDruid,
		func(character *core.Character, options *proto.Player, _ *proto.Raid) core.Agent {
			return NewFeralCatDruid(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_FeralCatDruid)
			if !ok {
				panic("Invalid spec value for Feral Druid!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewFeralCatDruid(character *core.Character, options *proto.Player) *FeralDruid {
	feralOptions := options.GetFeralCatDruid()
	selfBuffs := druid.SelfBuffs{}

	cat := &FeralDruid{
		Druid: druid.New(character, druid.Cat, selfBuffs, options.TalentsString),
	}

	cat.CannotShredTarget = feralOptions.Options.CannotShredTarget

	cat.EnableEnergyBar(core.EnergyBarOptions{
		MaxComboPoints: 5,
		MaxEnergy:      100.0,
		UnitClass:      proto.Class_ClassDruid,
	})
	// The Cat never swings in Bear Form unless a custom rotation shifts. Its bar pays the 100% crit Rage
	// bonus like the Bear's (core.CritRageMultiplier applies to every rage bar).
	cat.EnableRageBar(core.RageBarOptions{BaseRageMultiplier: 1})

	cat.EnableAutoAttacks(cat, core.AutoAttackOptions{
		// Base paw weapon.
		MainHand:       cat.GetCatWeapon(),
		AutoSwingMelee: true,
	})

	cat.RegisterCatFormAura()
	cat.RegisterBearFormAura()

	return cat
}

type FeralDruid struct {
	*druid.Druid

	Rotation FeralDruidRotation

	readyToShift   bool
	waitingForTick bool
}

func (cat *FeralDruid) GetDruid() *druid.Druid {
	return cat.Druid
}

func (cat *FeralDruid) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
}

// AddPartyBuffs auto-applies Leader of the Pack from the druid's own talent.
// Windfury Totem stays: Forever's 8515 is a party proc aura, not a weapon
// enchant, so a cat procs it (client 1.60.1.70205).
func (cat *FeralDruid) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
	if cat.Talents.LeaderOfThePack {
		partyBuffs.LeaderOfThePack = true
	}
}

func (cat *FeralDruid) Initialize() {
	cat.Druid.Initialize()
	cat.RegisterFeralCatSpells()
}

func (cat *FeralDruid) ApplyTalents() {
	cat.Druid.ApplyTalents()
}

func (cat *FeralDruid) Reset(sim *core.Simulation) {
	cat.Druid.Reset(sim)
	cat.Druid.ClearForm(sim)
	cat.CatFormAura.Activate(sim)
	cat.readyToShift = false
	cat.waitingForTick = false
}
