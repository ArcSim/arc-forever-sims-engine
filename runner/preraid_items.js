// The one pre-raid item rule every gear and weapon script shares: something a hunter can farm before raids.
// No Molten Core / Onyxia / Blackwing Lair / Zul'Gurub drops, world bosses, raid quests or crafts from raid
// materials; phase 1-2; hunter-usable; the character's faction only; no gear bought with PvP rank.
// Battleground reputation rewards count up to Revered (Exalted is a long farm): env GO_PVP=all adds Exalted,
// GO_NO_PVP=1 leaves them all out. GO_ALL_FOREVER=1 lets unscreened new Forever items in, for screening.
const RAID_ZONES = new Set([1977, 2159, 2677, 2717]);
// 40-man world bosses: Azuregos, Lord Kazzak and the four Emerald Dragons
const WORLD_BOSSES = new Set([6109, 12397, 14887, 14888, 14889, 14890]);
const HORDE_RACES = new Set(['RaceOrc', 'RaceTroll', 'RaceTauren', 'RaceUndead', 'RaceSkyborneHorde']);
// raid-gated items the data gives no raid drop for: quest rewards from raid content (Onyxia's head, Zul'Gurub coins
// and bijous, the hunter epic quest's Molten Core leaf, Thunderfury, Thunderstrike (Shadowstrike from Molten Core, transformed), the Hydraxian Waterlords' Molten Core rewards)
const RAID_QUESTS =
	/Onyxia|Dragonslayer|Black Dragon Sinew|Quel'Serrar|Zandalar|Hakkari|Primal|Rhok'delar|Lok'delar|Thunderfury|^Thunderstrike$|Ragnaros|^Tidal Loop$|^Ocean's Breeze$/i;
// bought with PvP rank: a rank title in the name and no reputation source
const PVP_RANK =
	/^(Private|Corporal|Sergeant|Master Sergeant|Sergeant Major|Knight|Knight-Lieutenant|Knight-Captain|Knight-Champion|Lieutenant Commander|Commander|Marshal|Field Marshal|Grand Marshal|Scout|Grunt|Senior Sergeant|First Sergeant|Stone Guard|Blood Guard|Legionnaire|Centurion|Champion|Lieutenant General|General|Warlord|High Warlord)'s /;
// later content than pre-raid (Silithus and Ahn'Qiraj era quest rewards)
const LATER = /^(Maelstrom's Tendril|Malfurion's Signet Ring|Band of Icy Depths|Staff of Rampant Growth)$/;
// crafts that need Molten Core, Blackwing Lair or Zul'Gurub materials
const RAID_MATS =
	/\b(Core|Flarecore|Molten|Lava|Fiery|Corehound|Volcanic|Blazefury|Bloodvine|Bloodsoul|Primal Bat|Dragonstalker|Black Dragonscale|Chromatic|Ironvine|Onyxia)\b/i;
// Forever crafts whose names match RAID_MATS but whose recipes use no raid materials (Wowhead Forever reagents)
const FOREVER_CRAFTS = /^Black Dragonscale Helm$/;
// Thorium Brotherhood weapons made with Fiery / Lava Cores or Elementium
const RAID_CRAFTS = /^(Dark Iron Reaver|Dark Iron Destroyer|Black Amnesty|Blackfury|Blackguard|Nightfall)$/;
const PVP_FACTIONS = [729, 730, 889, 890, 509, 510]; // Frostwolf, Stormpike, Warsong, Silverwing, Arathor, Defilers
const MELEE_TYPES = new Set([1, 2, 3, 6, 8, 9]); // axe dagger fist polearm staff sword

const isHorde = race => HORDE_RACES.has(race);
// new Forever items the gear lists show but no build may use: everything scripts/forever_item_sources.json does not add
// (gen_forever_new_items.js writes the list)
const FOREVER_UNAPPROVED = new Set(JSON.parse(require('fs').readFileSync(require('path').join(__dirname, 'forever_unapproved_ids.json'), 'utf8')));

function preRaid(item, race) {
	if ((item.phase || 1) > 2 || item.quality < 2 || item.ilvl < 45) return false;
	// GO_ALL_FOREVER=1 (screening only): new Forever items count before their source is checked
	if (!process.env.GO_ALL_FOREVER && FOREVER_UNAPPROVED.has(item.id)) return false;
	// Class enum: hunter is 2
	if (item.classAllowlist && item.classAllowlist.length && !item.classAllowlist.includes('ClassHunter') && !item.classAllowlist.includes(2)) return false;
	const sources = item.sources || [];
	if (sources.some(s => s.drop && (RAID_ZONES.has(s.drop.zoneId) || WORLD_BOSSES.has(s.drop.npcId)))) return false;
	// faction-locked rewards (1 Alliance, 2 Horde) follow the character's race
	if (item.factionRestriction && item.factionRestriction !== (isHorde(race) ? 2 : 1)) return false;
	if (sources.some(s => s.crafted) && !FOREVER_CRAFTS.test(item.name) && (RAID_MATS.test(item.name) || RAID_CRAFTS.test(item.name))) return false;
	const pvpRep = sources.filter(s => s.rep && PVP_FACTIONS.includes(s.rep.repFactionId));
	if (pvpRep.length && (process.env.GO_NO_PVP || (process.env.GO_PVP !== 'all' && pvpRep.every(s => s.rep.repLevel > 7)))) return false;
	if (PVP_RANK.test(item.name) && !sources.some(s => s.rep)) return false;
	if (LATER.test(item.name)) return false;
	if (RAID_QUESTS.test(item.name)) return false;
	if (item.armorType && item.armorType > 3) return false; // cloth 1, leather 2, mail 3
	// held-in-off-hand frills and shields are their own weapon types; off-hand-only weapons stay in
	if (item.type === 13 && !MELEE_TYPES.has(item.weaponType)) return false;
	if (item.type === 14 && ![1, 2, 3].includes(item.rangedWeaponType)) return false;
	return true;
}

module.exports = { preRaid, isHorde, MELEE_TYPES };
