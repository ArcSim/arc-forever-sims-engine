import * as PresetUtils from '@app/preset_utils';
import { Debuffs, IndividualBuffs, PartyBuffs, RaidBuffs } from '@generated/proto/buffs';
import { ConsumesSpec, Profession, Race, TristateEffect } from '@generated/proto/common';
import { EnhancementShaman_Options as EnhancementShamanOptions, ShamanImbue, ShamanSyncType } from '@generated/proto/shaman';
import { SavedTalents } from '@generated/proto/ui';

import ForeverApl from './apls/forever.apl.json';
import LaunchGear from './gear_sets/launch.gear.json';
import Phase1Gear from './gear_sets/phase_1.gear.json';
import Phase2Gear from './gear_sets/phase_2.gear.json';

// Our Forever APL (the parity tool's nextApl). The TBC default.apl.json it replaces never cast
// Earth Shock or Lightning Bolt (459 DPS at defaults, master 576).
export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', ForeverApl);

// Defaults below are what master's ui/enhancement_shaman opens with (currentSettings on a fresh
// profile); its raid-wide Battle Shout, Leader of the Pack and totems are party buffs here.
export const DefaultIndividualBuffs = IndividualBuffs.create({});

export const DefaultOptions = EnhancementShamanOptions.create({
	classOptions: {
		shieldProcrate: 0,
		imbueMh: ShamanImbue.WindfuryWeapon,
	},
	imbueOh: ShamanImbue.WindfuryWeapon,
	syncType: ShamanSyncType.Auto,
});

export const OtherDefaults = {
	reactionTime: 200, // master's default
	distanceFromTarget: 5,
	profession1: Profession.Alchemy,
	profession2: Profession.Enchanting,
	race: Race.RaceOrc,
};

// Master's consumables, as the Forever client's items.
export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 13512, // Flask of Supreme Power
	battleElixirId: 13452, // Elixir of the Mongoose
	spellPowerElixirId: 13454, // Greater Arcane Elixir
	schoolElixirId: 6373, // Elixir of Fire Power
	guardianElixirId: 20007, // Mageblood Elixir
	strengthBuffId: 12451, // Juju Power
	attackPowerBuffId: 12460, // Juju Might
	zanzaId: 8410, // R.O.I.D.S.
	dragonbreathChili: true,
	foodId: 13810, // Blessed Sunfruit
	potId: 13444, // Major Mana Potion
	conjuredId: 12662, // Demonic Rune
});

export const DefaultPartyBuffs = PartyBuffs.create({
	battleShout: TristateEffect.TristateEffectImproved,
	leaderOfThePack: true,
	manaSpringTotem: TristateEffect.TristateEffectRegular,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	fireResistanceTotem: true,
	giftOfTheWild: true,
	prayerOfFortitude: true,
	prayerOfSpirit: true,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfRecklessness: true,
	exposeArmor: true,
	faerieFire: true,
	sunderArmor: true,
});

// Talent presets, from master's ui/shaman spec.
export const TalentsLevel60 = PresetUtils.makePresetTalents('Level 60', SavedTalents.create({ talentsString: '5505031-053030031005112251' }));
export const TalentsEnhancement = PresetUtils.makePresetTalents('Enhancement 16/35/0', SavedTalents.create({ talentsString: '2532031-055030030205112251' }));
export const TalentPresets = [TalentsLevel60, TalentsEnhancement];

// Our Forever sim's gear presets (master ui/<spec>/gear_sets).
export const GEAR_LAUNCH = PresetUtils.makePresetGear('Launch', LaunchGear);
export const GEAR_PHASE_1 = PresetUtils.makePresetGear('Phase 1', Phase1Gear);
export const GEAR_PHASE_2 = PresetUtils.makePresetGear('Phase 2', Phase2Gear);
export const DEFAULT_GEAR = GEAR_LAUNCH;
export const GEAR_PRESETS = [GEAR_LAUNCH, GEAR_PHASE_1, GEAR_PHASE_2];
