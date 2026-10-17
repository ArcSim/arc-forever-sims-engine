package forever

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// RevelationEnchantModelled (mythicsim patch 157): Enchant Weapon - Revelation (enchant 8217) triggers
// its class buff. False registers nothing for the enchant, as before (the generator's "states no rate").
var RevelationEnchantModelled = true

// The chance a Revelation trigger rolls on each eligible hit that does not crit. Measured on public
// foreverlogs.gg reports (a Frost Mage, a Warlock and a Shaman, levels 20 to 35, 7.5% to 13% spell
// crit): 111 triggers in 1443 non-crit hits, 95% CI 0.073 to 0.106. The chance per hit is the same for
// Frostbolt, Searing Pain and Hellfire ticks, so it is not a PPM, and it does not fall across that crit
// range. "Diminished as your critical strike chance increases" is already the crit itself: a crit never
// triggers, so the chance per spell hit is this times (1 - crit). A variable so tests can set it.
var RevelationTriggerChance = 0.0885

// The equip aura the enchant applies (1248806, A_PROC_TRIGGER_SPELL, ProcFlags 0x14000, ProcChance 100,
// no PPM row and no trigger spell: the server's script picks the buff and the chance).
const revelationTriggerSpellID = 1248806

// Each class's buff, client 1.60.1.70425: one charge (SpellAuraOptions ProcCharges 1), 15 sec (duration
// index 8), and its only effect is A_ADD_FLAT_MODIFIER SPELLMOD_CRITICAL_CHANCE 100 on this class mask.
// The store does not carry these rows, so the masks are written out here. Warriors and rogues have none.
// The mask names the class's "non-periodic" spells: the mage's leaves out Fireball, Pyroblast and
// Frostfire Bolt, the warlock's Immolate and the dots, the shaman's Flame Shock; channels are on it
// (Arcane Missiles, Blizzard, Hellfire, Mind Flay, Hurricane), and so are the spells their ticks cast.
var revelationBuffs = map[proto.Class]struct {
	spellID int32
	mask    [4]uint32
}{
	proto.Class_ClassMage:    {1248808, [4]uint32{547494642, 0, 0, 8}},
	proto.Class_ClassPriest:  {1323377, [4]uint32{161758864, 8486934, 128, 0}},
	proto.Class_ClassWarlock: {1323392, [4]uint32{541161, 8650816, 0, 0}},
	proto.Class_ClassDruid:   {1323400, [4]uint32{4194341, 2, 0, 0}},
	proto.Class_ClassHunter:  {1323410, [4]uint32{10240, 0, 0, 0}},
	proto.Class_ClassShaman:  {1323418, [4]uint32{0x801001c3, 266240, 0, 0}},
	proto.Class_ClassPaladin: {1323419, [4]uint32{0xc0200400, 2097792, 64, 67108864}},
}

// The client's class flags for a sim spell: its own where spelldata resolved it, else its client row's.
func revelationSpellFlags(spell *core.Spell) core.ClassFlags {
	if !spell.ClassFlags.IsZero() {
		return spell.ClassFlags
	}
	return core.ClientClassFlags(spell.ActionID.SpellID)
}

// Whether a row carries a periodic effect of its own (a dot or hot part).
func hasPeriodicPart(row *spelldata.Spell) bool {
	for i := range row.Effects {
		switch row.Effects[i].Aura {
		case dbcenums.A_PERIODIC_DAMAGE, dbcenums.A_PERIODIC_LEECH, dbcenums.A_PERIODIC_HEAL:
			return true
		}
	}
	return false
}

// Whether a row ticks by casting a spell (an area trigger or a periodic trigger): Hellfire's ticks are
// Hellfire Effect (5857) and Flamestrike's burn is 1279985 in the game, direct hits, which the sim
// deals as the spell's own periodic ticks.
func ticksByCastingASpell(row *spelldata.Spell) bool {
	for i := range row.Effects {
		switch row.Effects[i].Aura {
		case dbcenums.A_PERIODIC_TRIGGER_SPELL, dbcenums.A_PERIODIC_DUMMY:
			return true
		}
	}
	return false
}

// Whether a hit of this spell can trigger Revelation. The trigger fires on non-periodic magic spell
// hits, damage or heal (ProcFlags 0x14000), including triggered spells (Attributes[3] 0x4000000, Can
// Proc From Procs: the logs show Frostbrand Attack and Lightning Shield triggering it):
//   - direct hits: magic damage or a heal from a spell whose row has no dot or hot part. The logs show
//     Immolate's and Flame Shock's direct hits triggering 1 time in 45 against 4 expected, so a spell
//     with a periodic part is left out (Fireball, Pyroblast, Frostfire Bolt, Moonfire, Holy Fire too);
//     this is the one part the logs do not settle.
//   - periodic ticks: only the area channels and burns whose ticks are cast spells in the game
//     (Hellfire, Flamestrike's burn, Consecration). Dot ticks (Corruption, Curse of Agony, Immolate,
//     Flame Shock) and periodic channels (Drain Life, Mind Flay) never trigger it: 3 in 525 in the logs.
func revelationCanTriggerFrom(spell *core.Spell, periodic bool) bool {
	row := spelldata.Find(spell.ActionID.SpellID)
	if periodic {
		return spell.AOEDot() != nil && ticksByCastingASpell(row)
	}
	if spell.ProcMask.Matches(core.ProcMaskSpellHealing) {
		return !hasPeriodicPart(row)
	}
	return spell.DefenseType == core.DefenseTypeMagic && !hasPeriodicPart(row)
}

func applyRevelation(agent core.Agent) {
	if !RevelationEnchantModelled {
		return
	}
	character := agent.GetCharacter()
	buff, ok := revelationBuffs[character.Class]
	if !ok {
		return
	}
	mask := core.ClassFlags{Family: core.ClassSpellFamilies[character.Class], Mask: buff.mask}

	var affected []*core.Spell
	isAffected := make(map[*core.Spell]bool)

	buffAura := character.RegisterAura(core.Aura{
		Label:    "Revelation",
		ActionID: core.ActionID{SpellID: buff.spellID},
		Duration: time.Second * 15,
		OnGain: func(_ *core.Aura, _ *core.Simulation) {
			for _, spell := range affected {
				spell.BonusCritPercent += 100
			}
		},
		OnExpire: func(_ *core.Aura, _ *core.Simulation) {
			for _, spell := range affected {
				spell.BonusCritPercent -= 100
			}
		},
		// The charge goes on the next cast of a spell the buff names, after that cast's hit is
		// computed (OnCastComplete follows ApplyEffects). A channel holds it to the channel's end
		// instead (below): Hellfire with Revelation up crits on 121 of 121 hits and the buff falls
		// when the channel ends; Blizzard 36 of 36. Triggered spells (a tick, Lightning Overload)
		// are not casts.
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !isAffected[spell] || spell.Flags.Matches(core.SpellFlagChanneled|core.SpellFlagProc) {
				return
			}
			aura.Deactivate(sim)
		},
	})

	character.OnSpellRegistered(func(spell *core.Spell) {
		if spell.ActionID.SpellID == 0 || !mask.Matches(revelationSpellFlags(spell)) {
			return
		}
		affected = append(affected, spell)
		isAffected[spell] = true

		if !spell.Flags.Matches(core.SpellFlagChanneled) {
			return
		}
		endChannel := func(_ *core.Aura, sim *core.Simulation) {
			buffAura.Deactivate(sim)
		}
		if dot := spell.AOEDot(); dot != nil {
			dot.ApplyOnExpire(endChannel)
		}
		for _, target := range character.Env.AllUnits {
			if dot := spell.Dot(target); dot != nil {
				dot.ApplyOnExpire(endChannel)
			}
		}
	})

	// Each target of an area spell rolls on its own: Hellfire casts with 2, 3 and 4 or more non-crit
	// targets triggered 5/33, 9/35 and 29/65 times against 5.1, 7.7 and 29 expected. The buff comes one
	// batch window after the hit (the trigger's default), so the rest of that cast's targets do not take it.
	makeTrigger := func(name string, callback core.AuraCallback, periodic bool) *core.Aura {
		eligible := make(map[*core.Spell]bool)
		return character.MakeProcTriggerAura(core.ProcTrigger{
			Name:             name,
			ActionID:         core.ActionID{SpellID: revelationTriggerSpellID},
			Callback:         callback,
			ProcMask:         core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc | core.ProcMaskSpellHealing,
			Outcome:          core.OutcomeLanded,
			CanProcFromProcs: true,
			ProcChance:       RevelationTriggerChance,
			ExtraCondition: func(_ *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
				if result.DidCrit() || result.Damage <= 0 {
					return false
				}
				ok, seen := eligible[spell]
				if !seen {
					ok = revelationCanTriggerFrom(spell, periodic)
					eligible[spell] = ok
				}
				return ok
			},
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				buffAura.Activate(sim)
			},
		})
	}
	direct := makeTrigger("Enchant Weapon - Revelation", core.CallbackOnSpellHitDealt|core.CallbackOnHealDealt, false)
	ticks := makeTrigger("Enchant Weapon - Revelation (ticks)", core.CallbackOnPeriodicDamageDealt, true)
	character.ItemSwap.RegisterEnchantProc(8217, direct)
	character.ItemSwap.RegisterEnchantProc(8217, ticks)
}

func init() {
	core.NewEnchantEffect(8217, applyRevelation)
}
