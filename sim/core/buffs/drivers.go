package buffs

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

// The buff drivers the generator calls by name. A generated apply block for a
// cooldown, proc, uptime, item-count or manual row, or for a row the manifest
// marks as driven, calls drive<Go>(unit, scope); declaring that function here
// is what makes the row compile. A buff-scope row
// hands the driver the *Character the buffs are applied to and the message its
// field lives on, a debuff-scope row the *Unit the debuffs land on plus the
// debuffs and the raid. The driver reads its own field out of the message, and
// any other field it needs: Grace of Air lasts 9 seconds while the party is
// twisting totems.

// Spell 29166 states 100% mana regen while casting (aura 134) and +400% of it
// (aura 110); neither is a stat or a pseudo-stat, so the regen is the driver's.
const innervateSpiritRegenMultiplier = 5.0

// Three pieces of Battlegear of Wrath are worth 30 more attack power on Battle
// Shout: item set 218's ItemSetSpell at three pieces is 23563, which adds a flat
// 30 to every effect of the Battle Shout family. The resolver reads no
// ItemSetSpell, so the amount is stated here.
const BattleShoutT2Bonus = 30.0

// The party's Battle Shout is the external caster's copy, which chains behind
// the player's own shout rather than being up from the start. Its improved state
// says that warrior shouted with the set on.
func driveBattleShout(char *core.Character, party *proto.PartyBuffs) {
	aura := BattleShoutAura(&char.Unit, false, 0)
	if party.BattleShout == proto.TristateEffect_TristateEffectImproved {
		core.AddGeneratedFlatBonus(aura, stats.AttackPower, BattleShoutValue(0), BattleShoutT2Bonus)
	}
	core.ApplyFixedShoutAura(char, aura, BattleShoutCategory)
}

// A druid innervates a character who is nearly out of mana, so that every other
// mana cooldown is spent first. The aura forces full spirit regen while it is
// up. Its mana is regen (29166 has no energize effect), so it lands in the regen
// metrics and makes no threat.
func driveInnervates(char *core.Character, individual *proto.IndividualBuffs) {
	aura := InnervatesAura(&char.Unit, false, 0)

	threshold := 0.0
	char.Env.RegisterPostFinalizeEffect(func() {
		threshold = innervateManaThreshold(char)
	})

	aura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
		char.PseudoStats.ForceFullSpiritRegen = true
		char.PseudoStats.SpiritRegenMultiplier *= innervateSpiritRegenMultiplier
		char.UpdateManaRegenRates()
	}).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
		char.PseudoStats.ForceFullSpiritRegen = false
		char.PseudoStats.SpiritRegenMultiplier /= innervateSpiritRegenMultiplier
		char.UpdateManaRegenRates()
	})

	core.NewGeneratedExternalCD(char, aura, core.GeneratedExternalCD{
		NumSources: individual.Innervates,
		Cooldown:   InnervatesCooldown(),
		Type:       core.CooldownTypeMana,
		ShouldActivate: func(_ *core.Simulation, char *core.Character) bool {
			return char.CurrentMana() <= threshold
		},
	})
}

// A mage burns mana fast enough that waiting for a flat thousand left would
// waste most of the innervate.
func innervateManaThreshold(char *core.Character) float64 {
	if char.Class == proto.Class_ClassMage {
		return char.MaxMana() * 0.4
	}
	return 1000
}

// The priest has nothing to hold Power Infusion for, so it goes out on cooldown.
func drivePowerInfusions(char *core.Character, individual *proto.IndividualBuffs) {
	core.NewGeneratedExternalCD(char, PowerInfusionsAura(&char.Unit, false, 0), core.GeneratedExternalCD{
		NumSources: individual.PowerInfusions,
		Cooldown:   PowerInfusionsCooldown(),
		Type:       core.CooldownTypeDPS,
	})
}

// A restoration shaman drops Mana Tide once the party has mana to refill, which
// is 40 seconds in, or halfway through a fight shorter than that.
func driveManaTideTotems(char *core.Character, party *proto.PartyBuffs) {
	initialDelay := time.Duration(0)
	char.Env.RegisterPostFinalizeEffect(func() {
		initialDelay = min(char.Env.BaseDuration/2, time.Second*40)
	})

	core.NewGeneratedExternalCD(char, ManaTideTotemsAura(&char.Unit, false, 0), core.GeneratedExternalCD{
		NumSources: party.ManaTideTotems,
		Cooldown:   ManaTideTotemsCooldown(),
		Type:       core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, _ *core.Character) bool {
			return sim.CurrentTime >= initialDelay
		},
	})
}

// The totem's aura is the attack power a windfury proc grants. Since build
// 70009 the totem hands the proc out as a party aura (10612, an area
// A_PROC_TRIGGER_SPELL) rather than Era's main-hand enchant 564, which nothing
// in the client applies any more; its row states the 20%, the 100 ms internal
// cooldown and what hears it: every melee auto and special, off-hand included.
// The extra attack 10610 grants is always a main-hand one.
var windfuryTotemPartyAura = spelldata.MustFind(10612)

// Since build 70009 Tranquil Air, Windfury and Grace of Air totems no longer
// stack, even from different shamans in the group (Forever beta development
// notes). One air totem stands: a totem the shaman casts replaces the one the
// party buffs assume, and of the two party buffs Windfury holds.
var AirTotemCategory = "AirTotem"

const (
	AirTotemPartyGraceOfAir = iota + 1
	AirTotemPartyWindfury
	AirTotemCastGraceOfAir
	AirTotemCastWindfury
)

// Flametongue Totem no longer stacks with Flametongue Weapon or Windfury Totem (the same notes), and a
// main-hand Flametongue Weapon "disables any benefit you personally receive from Flametongue Totem"
// (tooltips 8024/16342). Neither text says which of the two totems holds; Windfury, the stronger, does.
var FlametongueTotemCategory = "FlametongueTotem"

const (
	FlametongueTotemCast = iota + 1
	FlametongueTotemWindfuryTotem
	FlametongueTotemMainHandImbue
)

func driveWindfuryTotem(char *core.Character, _ *proto.PartyBuffs) {
	procAura := WindfuryTotemAura(&char.Unit, false, 0)
	// The attack power is only there for a moment after a proc, so it is not
	// part of the stats the character sheet is measured with.
	procAura.BuildPhase = core.CharacterBuildPhaseNone

	// The row's own proc flags say what spends a charge: every auto attack that
	// lands. The attack power stays until the charges are gone or the row's
	// duration runs out.
	procAura.MaxStacks = int32(windfuryTotemSpell.ProcCharges)
	spender := spelldata.ProcTrigger(char, windfuryTotemSpell, func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
		procAura.RemoveStack(sim)
	}, spelldata.Chance(1))
	spender.Name = "Windfury Totem Charges"
	spender.TriggerImmediately = true
	procAura.AttachProcTriggerCallback(&char.Unit, spender)

	// An auto that procs it has spent the first charge itself, so the extra
	// attack is the only one buffed; a special hands both charges to the extra
	// attack and the auto after it.
	var windfurySpell *core.Spell
	trigger := spelldata.ProcTrigger(char, windfuryTotemPartyAura, func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		procAura.Activate(sim)
		charges := procAura.MaxStacks
		if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
			charges--
		}
		procAura.SetStacks(sim, charges)
		char.AutoAttacks.MaybeReplaceMHSwing(sim, windfurySpell).Cast(sim, result.Target)
	})
	trigger.Name = "Windfury Totem Trigger"
	trigger.ActionID = core.ActionID{}
	trigger.MetricsActionID = core.ActionID{SpellID: 25580, Tag: -1}
	trigger.Duration = core.NeverExpires
	trigger.TriggerImmediately = true
	procTrigger := char.MakeProcTriggerAura(trigger)

	// The totem stands for 10 seconds and the shaman drops a new one every 5,
	// so the aura that holds the category is simply refreshed.
	totemAura := char.GetOrRegisterAura(core.Aura{
		Label:    "Windfury Totem",
		ActionID: core.ActionID{SpellID: 25587, Tag: -1},
		Duration: time.Second * 10,
	}).ApplyOnInit(func(aura *core.Aura, sim *core.Simulation) {
		config := *char.AutoAttacks.MHConfig()
		config.ActionID = config.ActionID.WithTag(25584)
		windfurySpell = char.GetOrRegisterSpell(config)
	}).ApplyOnReset(func(aura *core.Aura, sim *core.Simulation) {
		aura.Activate(sim)
		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period:   time.Second * 5,
			Priority: core.ActionPriorityAuto,
			OnAction: func(sim *core.Simulation) {
				aura.Activate(sim)
			},
		})
	})

	// Joined first, so a cast air totem refuses the aura before it arms the proc.
	totemAura.NewExclusiveEffect(AirTotemCategory, true, core.ExclusiveEffect{Priority: AirTotemPartyWindfury})
	totemAura.NewExclusiveEffect(FlametongueTotemCategory, false, core.ExclusiveEffect{Priority: FlametongueTotemWindfuryTotem})
	totemAura.NewExclusiveEffect(WindfuryTotemCategory, false, core.ExclusiveEffect{
		Priority: WindfuryTotemValue(0),
		OnGain: func(_ *core.ExclusiveEffect, sim *core.Simulation) {
			procTrigger.Activate(sim)
		},
		OnExpire: func(_ *core.ExclusiveEffect, sim *core.Simulation) {
			procTrigger.Deactivate(sim)
			totemAura.Deactivate(sim)
		},
	})
}

// A shaman twisting totems keeps Grace of Air up for 9 seconds out of every 10,
// because the air slot is holding another totem the rest of the time; a shaman
// who is not twisting leaves it standing.
func driveGraceOfAirTotem(char *core.Character, party *proto.PartyBuffs) {
	aura := GraceOfAirTotemAura(&char.Unit, false, 0)
	aura.NewExclusiveEffect(AirTotemCategory, true, core.ExclusiveEffect{Priority: AirTotemPartyGraceOfAir})

	if !party.TotemTwisting {
		core.MakePermanent(aura)
		return
	}

	// The first cast lands a totem cycle into the fight, because the shaman
	// spends the opening one on the totem being twisted with.
	aura.Duration = time.Second * 9
	aura.ApplyOnReset(func(aura *core.Aura, sim *core.Simulation) {
		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period:   time.Second * 10,
			Priority: core.ActionPriorityAuto,
			OnAction: func(sim *core.Simulation) {
				aura.Activate(sim)
			},
		})
	})
}

// The party's Retribution Aura is the top rank at its base damage: it scales with
// the paladin's spell power, which the party config does not carry.
func driveRetributionAura(char *core.Character, _ *proto.PartyBuffs) {
	core.MakePermanent(RetributionAuraBuff(char, false, RetributionAuraMaxRank))
}

// The blessing does nothing on its own: the paladin's Holy Light and Flash of
// Light read their bonus off their own rows when the target carries it.
func driveGreaterBlessingOfLight(char *core.Character, _ *proto.IndividualBuffs) {
	core.MakePermanent(GreaterBlessingOfLightAura(&char.Unit, false, 0))
}

// The staff's aura is worth its amounts once per Atiesh in the party.
func driveAtieshDruid(char *core.Character, party *proto.PartyBuffs) {
	core.MakePermanent(AtieshDruidAura(&char.Unit, false, 0, float64(party.AtieshDruid)))
}

func driveAtieshMage(char *core.Character, party *proto.PartyBuffs) {
	core.MakePermanent(AtieshMageAura(&char.Unit, false, 0, float64(party.AtieshMage)))
}

func driveAtieshPriest(char *core.Character, party *proto.PartyBuffs) {
	core.MakePermanent(AtieshPriestAura(&char.Unit, false, 0, float64(party.AtieshPriest)))
}

func driveAtieshWarlock(char *core.Character, party *proto.PartyBuffs) {
	core.MakePermanent(AtieshWarlockAura(&char.Unit, false, 0, float64(party.AtieshWarlock)))
}

// The judgement the paladin leaves on the target heals whoever strikes it; the
// client's trigger spell 5373 is a dummy, so how much and how often is the
// driver's. The raid's copy is the top rank the paladin's own judgement states.
func driveJudgementOfLight(target *core.Unit, _ *proto.Debuffs, _ *proto.Raid) {
	AttachJudgementOfLightHeal(core.MakePermanent(JudgementOfLightAura(target, false, 0)), JudgementOfLightMaxRank)
}

// Judgement of Wisdom returns mana to whoever strikes the target, on the same
// terms: 1826 is a dummy, so the amount and the chance stay with the paladin's
// top rank.
func driveJudgementOfWisdom(target *core.Unit, _ *proto.Debuffs, _ *proto.Raid) {
	AttachJudgementOfWisdomMana(core.MakePermanent(JudgementOfWisdomAura(target, false, 0)), JudgementOfWisdomMaxRank)
}

// A stack of Sunder Armor is worth nothing until it is on the target, so the
// raid's copy is ramped to five over the first five global cooldowns, which is
// how long a warrior takes to stack it.
func driveSunderArmor(target *core.Unit, _ *proto.Debuffs, _ *proto.Raid) {
	aura := core.MakePermanent(SunderArmorAura(target, false, 0))

	core.ScheduledAura(aura, core.PeriodicActionOptions{
		Period:          core.GCDDefault,
		NumTicks:        5,
		TickImmediately: true,
		Priority:        core.ActionPriorityDOT,
		OnAction: func(sim *core.Simulation) {
			aura.Activate(sim)
			if aura.IsActive() {
				aura.AddStack(sim)
			}
		},
	})
}
