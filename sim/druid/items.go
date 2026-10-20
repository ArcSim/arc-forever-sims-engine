package druid

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// HowlingIdolFromRow (mythicsim patch 179, as wowsims/forever 7bf2911cbe): Howling Idol takes the cooldown
// 1291059's row states off Shifting Power, 500 ms at client 70425. False restores patch 103's 1 second.
var HowlingIdolFromRow = true

func init() {
	// Idol of the Moon
	core.NewItemEffect(23197, func(agent core.Agent) {
		character := agent.GetCharacter()
		aura := core.MakePermanent(character.RegisterAura(core.Aura{
			Label: "Improved Moonfire",
		}).AttachSpellMod(core.SpellModConfig{
			ClassMask:  DruidSpellMoonfire,
			Kind:       core.SpellMod_BaseDamage_Flat,
			FloatValue: 33.0,
		}))

		// TODO: this registers the swap proc against 32330, not 23197. That item is not in
		// the Forever database either, so the registration is dead whichever was intended.
		character.ItemSwap.RegisterProc(32330, aura)
	})

	// Idol of Ferocity
	// https://www.wowhead.com/forever/item=22397/idol-of-ferocity
	//
	// Reduces the energy cost of Claw and Rake by 2 (27851); Classic's took 3.
	core.NewItemEffect(22397, func(agent core.Agent) {
		agent.GetCharacter().AddStaticMod(core.SpellModConfig{
			ClassMask: DruidSpellClaw | DruidSpellRake,
			Kind:      core.SpellMod_PowerCost_Flat,
			IntValue:  -2,
		})
	})

	// Idol of Brutality
	// https://www.wowhead.com/forever/item=23198/idol-of-brutality
	//
	// Reduces the rage cost of Maul and Swipe by 2 (28855: -20 in the client's tenths of rage, the
	// scale Ferocity's -10 a rank is on); Classic's took 3. The client's mask also names Primal Bite, as
	// Ferocity's does.
	core.NewItemEffect(23198, func(agent core.Agent) {
		agent.GetCharacter().AddStaticMod(core.SpellModConfig{
			ClassMask: DruidSpellMaul | DruidSpellSwipe | DruidSpellPrimalBite,
			Kind:      core.SpellMod_PowerCost_Flat,
			IntValue:  -2,
		})
	})

	// Howling Idol
	// https://www.wowhead.com/forever/item=272427/howling-idol
	//
	// "Reduces the cooldown of your Shifting Power ability by 0.5 sec." (1291059: A_ADD_FLAT_MODIFIER,
	// SPELLMOD_COOLDOWN on Shifting Power's class mask). Stacks with Improved Shifting Power. Patch 103 wrote
	// the -1000 ms of an earlier client; client 70425's row says -500 ms, which HowlingIdolFromRow reads.
	core.NewItemEffect(272427, func(agent core.Agent) {
		reduction := -time.Second
		if HowlingIdolFromRow {
			reduction = time.Duration(spelldata.MustFind(1291059).EffectN(1).BaseValue()) * time.Millisecond
		}
		agent.GetCharacter().AddStaticMod(core.SpellModConfig{
			ClassMask: DruidSpellShiftingPower,
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: reduction,
		})
	})

	// Forever Wolfshead Helm (8345, effect 17768): the bonus belongs to
	// Shifting Power (Tiger's Fury before client 70170) and Enrage, not to entering a form.
	core.NewItemEffect(8345, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label:    "Wolfshead Helm",
			ActionID: core.ActionID{SpellID: 17768},
			OnGain: func(_ *core.Aura, _ *core.Simulation) {
				druid.WolfsheadShiftingPowerEnergy += 20
				druid.WolfsheadEnrageRage += 5
			},
			OnExpire: func(_ *core.Aura, _ *core.Simulation) {
				druid.WolfsheadShiftingPowerEnergy -= 20
				druid.WolfsheadEnrageRage -= 5
			},
		}))
	})
}
