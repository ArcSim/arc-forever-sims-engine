package warlock

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

type WarlockPet struct {
	core.Pet

	owner *Warlock

	AutoCastAbilities []*core.Spell
	MinMana           float64 // The minimum mana the AI keeps before it casts again

	DemonicBrandAura *core.Aura
	SoulLinkAura     *core.Aura
}

// Level 60 pet stats, from our Forever sim. Forever's demons inherit nothing from the warlock, so
// there is no stat inheritance function to speak of.
var petBaseStats = map[proto.WarlockOptions_Summon]stats.Stats{
	proto.WarlockOptions_Imp: {
		stats.Strength:  122,
		stats.Agility:   35,
		stats.Stamina:   86,
		stats.Intellect: 264,
		stats.Spirit:    260,
		stats.Mana:      576,
	},
	proto.WarlockOptions_Voidwalker: {
		stats.Strength:  129,
		stats.Agility:   85,
		stats.Stamina:   234,
		stats.Intellect: 70,
		stats.Spirit:    150,
		stats.Mana:      1066,
	},
	// The Succubus' Lash of Pain crits ~6% at the beta cap, master's crit or not: 13 of 205 vs
	// level 30 Cave Yetis (forever-bugs #331, staff skill 1-150, 4.58% sheet spell crit) and 11 of
	// 149 rank 6 casts in foreverlogs 2698-2771 (scratch wl_petcrit.py). Dev AidanZMoon there: demons
	// carry their own base crit and don't take the warlock's yet, 'base spell crit around normal levels'.
	proto.WarlockOptions_Succubus: {
		stats.Strength:         129,
		stats.Agility:          85,
		stats.Stamina:          234,
		stats.Intellect:        70,
		stats.Spirit:           150,
		stats.Mana:             1066,
		stats.SpellCritPercent: 6.2,
	},
	proto.WarlockOptions_Felhunter: {
		stats.Strength:  129,
		stats.Agility:   85,
		stats.Stamina:   234,
		stats.Intellect: 70,
		stats.Spirit:    150,
		stats.Mana:      1066,
	},
}

var petAutoAttacks = map[proto.WarlockOptions_Summon][2]float64{
	proto.WarlockOptions_Voidwalker: {31, 46},
	proto.WarlockOptions_Succubus:   {95, 131},
	proto.WarlockOptions_Felhunter:  {70, 97},
}

func AutoAttackConfig(min float64, max float64) *core.AutoAttackOptions {
	return &core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin: min,
			BaseDamageMax: max,
			SwingSpeed:    2.0,
		},
		AutoSwingMelee: true,
	}
}

func (warlock *Warlock) registerPets() {
	warlock.Imp = warlock.registerPet(proto.WarlockOptions_Imp)
	warlock.Voidwalker = warlock.registerPet(proto.WarlockOptions_Voidwalker)
	warlock.Succubus = warlock.registerPet(proto.WarlockOptions_Succubus)
	warlock.Felhunter = warlock.registerPet(proto.WarlockOptions_Felhunter)

	warlock.BasePets = []*WarlockPet{warlock.Imp, warlock.Voidwalker, warlock.Succubus, warlock.Felhunter}
}

func (warlock *Warlock) registerPet(summon proto.WarlockOptions_Summon) *WarlockPet {
	enabledOnStart := warlock.Options.Summon == summon

	pet := &WarlockPet{
		owner: warlock,
		Pet: core.NewPet(core.PetConfig{
			Name:            proto.WarlockOptions_Summon_name[int32(summon)],
			Owner:           &warlock.Character,
			BaseStats:       petBaseStats[summon],
			StatInheritance: func(_ stats.Stats) stats.Stats { return stats.Stats{} },
			EnabledOnStart:  enabledOnStart,
		}),
	}

	pet.EnableManaBar()
	if summon == proto.WarlockOptions_Imp {
		pet.AddStatDependency(stats.Intellect, stats.SpellCritPercent, core.CritPerIntMaxLevel[proto.Class_ClassMage])
	} else {
		pet.AddStatDependency(stats.Strength, stats.AttackPower, 2)
		pet.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[proto.Class_ClassWarrior])
		pet.AddStatDependency(stats.Intellect, stats.SpellCritPercent, core.CritPerIntMaxLevel[proto.Class_ClassWarrior])

		damage := petAutoAttacks[summon]
		pet.EnableAutoAttacks(pet, *AutoAttackConfig(damage[0], damage[1]))
	}

	if enabledOnStart {
		warlock.ActivePet = pet
		warlock.RegisterResetEffect(func(sim *core.Simulation) {
			warlock.ActivePet = pet
		})
	}

	warlock.AddPet(pet)
	return pet
}

func (warlock *Warlock) registerPetAbilities() {
	if warlock.Imp != nil {
		warlock.Imp.registerFireboltSpell()
	}
	if warlock.Succubus != nil {
		warlock.Succubus.registerLashOfPainSpell()
	}
	// The Voidwalker's Torment is a threat ability and is not modelled.
}

func (pet *WarlockPet) GetPet() *core.Pet {
	return &pet.Pet
}

func (pet *WarlockPet) Reset(_ *core.Simulation) {
}

func (pet *WarlockPet) OnEncounterStart(_ *core.Simulation) {
}

func (pet *WarlockPet) ExecuteCustomRotation(sim *core.Simulation) {
	waitUntil := time.Duration(0)

	for _, spell := range pet.AutoCastAbilities {
		if spell.CanCast(sim, pet.CurrentTarget) && pet.CurrentMana() > pet.MinMana {
			spell.Cast(sim, pet.CurrentTarget)
			return
		}

		cost := max(pet.MinMana, spell.Cost.GetCurrentCost())
		regen := pet.ManaRegenPerSecondWhileCasting()
		if regen > 0 {
			timeTillMana := max(0, (cost-pet.CurrentMana())/regen)
			waitUntil = min(waitUntil, time.Duration(float64(time.Second)*timeTillMana))
		}
	}

	// for now average the delay out to 100 ms so we don't need to roll random every time
	pet.WaitUntil(sim, sim.CurrentTime+waitUntil+time.Millisecond*100)
}

// The generator makes no table for pet spells, so Firebolt and Lash of Pain carry our
// client-verified beta 1.60.1 values (Firebolt rank 7, Lash of Pain rank 6). Firebolt's cooldown is the
// Imp's pause between casts: beta logs time 281 Firebolts 2.435 sec apart (median, 2.0 sec cast, next one
// starting 0.43 sec after the last lands; foreverlogs.gg 2687, 2695). The rotation polls every 100 ms, so
// 400 ms lands the next cast 2.4 sec after the last; 430 would round up to 2.5. Improved Imp and Improved Sayaad ride on their
// talents as SpellMods.
// Firebolt rank 7's damage effect, client 1.60.1.70094 11763: 44 + 0.6 a level from 58 (capped at 63),
// +-11.4%, so 45 (42.4 to 47.6) at 60. Beta logs agree on the row: level 20 imps' rank 2 (7799, 7 + 0.2 a
// level to 13) hit for 8 (126 of 138 non-crits; 6 each of 7 and 9), foreverlogs 2674.
var impFireboltEffect = spelldata.Effect{BasePoints: 44, PPL: 0.6000000238418579, Variance: 0.11363636702, SpellLevel: 58, MaxLevel: 63}

func (pet *WarlockPet) registerFireboltSpell() {
	pet.MinMana = 115

	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 11763},
		SpellSchool:    core.SpellSchoolFire,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellImpFireBolt,

		ManaCost: core.ManaCostOptions{FlatCost: 115},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDMin,
				CastTime: time.Second * 2,
			},
			CD: core.Cooldown{
				Timer:    pet.NewTimer(),
				Duration: time.Millisecond * 400,
			},
		},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,
		BonusCoefficient:         0.571,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, impFireboltEffect.Roll(sim, core.CharacterLevel), spell.OutcomeMagicHitAndCrit)
		},
	}))
}

func (pet *WarlockPet) registerLashOfPainSpell() {
	pet.MinMana = 160

	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 11780},
		SpellSchool:    core.SpellSchoolShadow,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellSuccubusLashOfPain,

		ManaCost: core.ManaCostOptions{FlatCost: 160},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    pet.NewTimer(),
				Duration: time.Second * 12,
			},
		},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,
		BonusCoefficient:         0.429,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, 50, spell.OutcomeMagicHitAndCrit)
		},
	}))
}
