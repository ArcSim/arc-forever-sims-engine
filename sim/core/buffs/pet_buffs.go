package buffs

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// PartyAurasReachPets gives a player's summoned pet the party auras that reach it in game. Forever's
// party auras are area auras on the party (SpellEffect type 35), and the server applies them to the
// party's pets unless the spell carries ATTR_EX_3_ONLY_TARGET_PLAYERS. Beta logs show the rule: in
// foreverlogs.gg report 2695 (2 October) a Warlock's Imp and a Hunter's pet carry Retribution Aura
// (7294, no flag) 90% and 96% of the run while the players beside them carry Devotion Aura (643,
// flagged) 83% to 90% and the pets never do. In report 2748 (8 October) two pets carry Retribution
// Aura and none of the 35 pets carries Devotion Aura, which players in the same raid do.
// The raid-wide buffs (Arcane Brilliance, Prayer of Spirit, Gift of the Wild, Prayer of Fortitude)
// and the Greater Blessings carry the flag, so a pet gets none of them.
//
// Left out: Mana Tide Totem, which the sim casts as the player's external cooldown; Retribution Aura,
// since the sim's pets take no hits; Blood Pact, which adds only health and which the Imp casting it
// does not carry itself (report 2695: the Hunter's pet has it 96% of the run, the Imp never); and
// Windfury and Flametongue Totem, whose procs hang off a weapon the pets do not have.
var PartyAurasReachPets = true

type petPartyBuff struct {
	spell *spelldata.Spell
	apply func(pet *core.Character, party *proto.PartyBuffs)
}

var petPartyBuffs = []petPartyBuff{
	{manaSpringTotemMeta.Spell, func(pet *core.Character, party *proto.PartyBuffs) {
		if party.ManaSpringTotem != proto.TristateEffect_TristateEffectMissing {
			core.MakePermanent(ManaSpringTotemAura(&pet.Unit, false, core.GetTristateValueInt32(party.ManaSpringTotem, 0, 5)))
		}
	}},
	{trueshotAuraMeta.Spell, func(pet *core.Character, party *proto.PartyBuffs) {
		if party.TrueshotAura {
			core.MakePermanent(TrueshotAuraAura(&pet.Unit, false, 0))
		}
	}},
	{atieshMageMeta.Spell, func(pet *core.Character, party *proto.PartyBuffs) {
		if party.AtieshMage > 0 {
			driveAtieshMage(pet, party)
		}
	}},
	{atieshWarlockMeta.Spell, func(pet *core.Character, party *proto.PartyBuffs) {
		if party.AtieshWarlock > 0 {
			driveAtieshWarlock(pet, party)
		}
	}},
	{atieshDruidMeta.Spell, func(pet *core.Character, party *proto.PartyBuffs) {
		if party.AtieshDruid > 0 {
			driveAtieshDruid(pet, party)
		}
	}},
	{atieshPriestMeta.Spell, func(pet *core.Character, party *proto.PartyBuffs) {
		if party.AtieshPriest > 0 {
			driveAtieshPriest(pet, party)
		}
	}},
}

// The party auras a summoned pet carries: those of petPartyBuffs that the party has and the client
// lets reach a pet.
func applyPetPartyBuffs(pet *core.Character, party *proto.PartyBuffs) {
	if !PartyAurasReachPets || party == nil {
		return
	}
	for _, b := range petPartyBuffs {
		if !b.spell.OnlyTargetsPlayers() {
			b.apply(pet, party)
		}
	}
}
