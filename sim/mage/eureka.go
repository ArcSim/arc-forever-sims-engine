package mage

import "github.com/wowsims/forever/sim/core"

// A Gnome mage's Eureka! (1259817) in client 1.60.1.70170. Effect 0 cuts the cost and effect 1 raises the
// damage of the same fifteen spells; effect 2, the periodic bonus Fireball and Frostfire Bolt had, now
// names none, so the ticks of Fireball, Frostfire Bolt, Pyroblast and Flamestrike keep their damage and no
// spell is in Tick. Arcane Missiles and Blizzard deal their damage as hits of a triggered spell, which the
// damage list names (the Arcane Missile tick, 1308937, is new on it). Pyroblast and Frost Nova are new on
// both lists; Combustion, Ignite and the mana gems were never on either. Frost Blast (client 70425) carries
// Frostbolt's bit, which both of 1259817's masks name (MythicSim patch 153).
func (mage *Mage) EurekaSpells() core.EurekaSpells {
	cost := MageSpellArcaneBlast | MageSpellArcaneExplosion | MageSpellArcaneMissilesCast | MageSpellBlastWave |
		MageSpellBlizzard | MageSpellConeOfCold | MageSpellFireBlast | MageSpellFireball | MageSpellFlamestrike |
		MageSpellFrostNova | MageSpellFrostbolt | MageSpellFrostfireBolt | MageSpellIceLance | MageSpellPyroblast |
		MageSpellScorch | MageSpellFrostBlast
	return core.EurekaSpells{
		Cost:   cost,
		Damage: cost | MageSpellArcaneMissilesTick,
	}
}
