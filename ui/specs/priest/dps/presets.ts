import * as PresetUtils from '@app/preset_utils';
import { Debuffs, IndividualBuffs, PartyBuffs, RaidBuffs } from '@generated/proto/buffs';
import { ConsumesSpec, Profession, Race, TristateEffect } from '@generated/proto/common';
import { DpsPriest_Options as Options } from '@generated/proto/priest';
import { SavedTalents } from '@generated/proto/ui';

import ShadowApl from './apls/shadow.apl.json';
import SmiteApl from './apls/smite.apl.json';
import LaunchGear from './gear_sets/launch.gear.json';
import P0BisGear from './gear_sets/p0.bis.gear.json';
import P1BisGear from './gear_sets/p1.bis.gear.json';
import SmiteLaunchGear from './gear_sets/smite_launch.gear.json';

// Defaults are master's ui/shadow_priest and ui/smite_priest (the Forever site before the switch),
// on our Forever APLs. Both pages sim the one dps priest spec; ../smite/spec.ts is the Smite page.

export const ROTATION_PRESET_SHADOW = PresetUtils.makePresetAPLRotation('Shadow', ShadowApl);
export const ROTATION_PRESET_SMITE = PresetUtils.makePresetAPLRotation('Smite', SmiteApl);

export const TalentsP1Shadow = PresetUtils.makePresetTalents('Shadow', SavedTalents.create({ talentsString: '00530323103--50552001211300051' }));
// The community builds the rankings page runs on our Forever sim: Shadow 15/0/36 and Smite 31/17/3.
export const ShadowTalents = PresetUtils.makePresetTalents('Shadow 15/0/36', SavedTalents.create({ talentsString: '02530003101--55052201201302251' }));
// Thirty-one points of Discipline for Power Infusion, twenty of Holy for the pieces that scale Smite.
export const TalentsLaunchSmite = PresetUtils.makePresetTalents('Smite', SavedTalents.create({ talentsString: '515335031030001001-30505113002' }));
export const SmiteTalents = PresetUtils.makePresetTalents('Smite 31/17/3', SavedTalents.create({ talentsString: '515035031030001031-00505023002-003' }));

// Master's priest options are empty: no Inner Fire, and the Shadow rotation casts Shadowform itself.
export const DefaultOptions = Options.create({
	classOptions: {},
});

export const ShadowConsumables = ConsumesSpec.create({
	flaskId: 13512, // Flask of Supreme Power
	spellPowerElixirId: 13454, // Greater Arcane Elixir
	schoolElixirId: 9264, // Elixir of Shadow Power
	guardianElixirId: 20007, // Mageblood Elixir
	zanzaId: 8423, // Cerebral Cortex Compound
	foodId: 18254, // Runn Tum Tuber Surprise
	potId: 13444, // Major Mana Potion
	conjuredId: 12662, // Demonic Rune
	mhImbueId: 20749, // Brilliant Wizard Oil
});
export const SmiteConsumables = ShadowConsumables;

// What master's page opens with once it has applied its presets (it drops the Fire Resistance
// Aura, Blessing of Wisdom and Judgement of Wisdom its presets name).
export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	giftOfTheWild: true,
	prayerOfSpirit: true,
	fireResistanceTotem: true, // a raid buff in the new buffs proto
});

export const DefaultPartyBuffs = PartyBuffs.create({
	manaSpringTotem: TristateEffect.TristateEffectImproved,
	moonkinAura: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({});

export const DefaultDebuffs = Debuffs.create({});

export const ShadowOtherDefaults = {
	reactionTime: 200, // master's default
	channelClipDelay: 100,
	distanceFromTarget: 20, // Mind Flay is 20 yd in the client
	profession1: Profession.Alchemy,
	profession2: Profession.Enchanting,
	race: Race.RaceTroll,
};
export const SmiteOtherDefaults = { ...ShadowOtherDefaults, distanceFromTarget: 30 };

// Our Forever sim's gear presets (master ui/<spec>/gear_sets).
export const GEAR_LAUNCH = PresetUtils.makePresetGear('Launch', LaunchGear);
export const GEAR_P0_BIS = PresetUtils.makePresetGear('Pre-BiS', P0BisGear);
export const GEAR_P1_BIS = PresetUtils.makePresetGear('P1 BiS', P1BisGear);
export const GEAR_SMITE_LAUNCH = PresetUtils.makePresetGear('Launch', SmiteLaunchGear);
export const DEFAULT_GEAR = GEAR_P0_BIS;
export const GEAR_PRESETS = [GEAR_LAUNCH, GEAR_P0_BIS, GEAR_P1_BIS];
