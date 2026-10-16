import * as PresetUtils from '@app/preset_utils';
import { Debuffs, IndividualBuffs, PartyBuffs, RaidBuffs } from '@generated/proto/buffs';
import { ConsumesSpec, Profession, Race, TristateEffect } from '@generated/proto/common';
import { ElementalShaman_Options as ElementalShamanOptions } from '@generated/proto/shaman';
import { SavedTalents } from '@generated/proto/ui';

import ForeverApl from './apls/forever.apl.json';
import LaunchGear from './gear_sets/launch.gear.json';
import Phase1Gear from './gear_sets/phase_1.gear.json';
import Phase2Gear from './gear_sets/phase_2.gear.json';

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', ForeverApl);

export const DefaultOptions = ElementalShamanOptions.create({
	classOptions: {
		shieldProcrate: 0,
	},
});

// Defaults below are master's ui/elemental_shaman (the Forever site before the switch).
export const OtherDefaults = {
	reactionTime: 200, // master's default
	distanceFromTarget: 15,
	profession1: Profession.Enchanting,
	profession2: Profession.Alchemy,
	race: Race.RaceTroll,
};

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	giftOfTheWild: true,
	prayerOfFortitude: true,
	prayerOfSpirit: true,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	moonkinAura: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({});

// Master also sets Stormstrike (+20% Nature damage taken), which this engine has no debuff for.
export const DefaultDebuffs = Debuffs.create({
	curseOfElements: true,
});

// Master's consumables, as the Forever client's items; one school elixir (Fire Power).
export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 13512, // Flask of Supreme Power
	spellPowerElixirId: 13454, // Greater Arcane Elixir
	schoolElixirId: 6373, // Elixir of Fire Power
	guardianElixirId: 20007, // Mageblood Elixir
	strengthBuffId: 12451, // Juju Power
	attackPowerBuffId: 12460, // Juju Might
	zanzaId: 8423, // Cerebral Cortex Compound
	foodId: 18254, // Runn Tum Tuber Surprise
	potId: 13444, // Major Mana Potion
	conjuredId: 12662, // Demonic Rune
});

export const TalentsLevel60 = PresetUtils.makePresetTalents('Level 60', SavedTalents.create({ talentsString: '5505331000001351--553302001' }));
export const TalentsElemental = PresetUtils.makePresetTalents(
	'Elemental 31/6/14',
	SavedTalents.create({ talentsString: '3505331000011351-0500001-053050001' }),
);
export const TalentsStormcaller = PresetUtils.makePresetTalents(
	'Stormcaller 28/23/0',
	SavedTalents.create({ talentsString: '253533100000033-055030030004102' }),
);
export const TalentPresets = [TalentsLevel60, TalentsElemental, TalentsStormcaller];

// Our Forever sim's gear presets (master ui/<spec>/gear_sets).
export const GEAR_LAUNCH = PresetUtils.makePresetGear('Launch', LaunchGear);
export const GEAR_PHASE_1 = PresetUtils.makePresetGear('Phase 1', Phase1Gear);
export const GEAR_PHASE_2 = PresetUtils.makePresetGear('Phase 2', Phase2Gear);
export const DEFAULT_GEAR = GEAR_LAUNCH;
export const GEAR_PRESETS = [GEAR_PHASE_2, GEAR_LAUNCH, GEAR_PHASE_1];
