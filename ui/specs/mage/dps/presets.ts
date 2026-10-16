import * as PresetUtils from '@app/preset_utils';
import { Debuffs, IndividualBuffs, PartyBuffs, RaidBuffs } from '@generated/proto/buffs';
import { ConsumesSpec, Profession, Race, Spec, TristateEffect } from '@generated/proto/common';
import { Mage_Options as MageOptions, Mage_Rotation, MageArmor } from '@generated/proto/mage';
import { SavedTalents } from '@generated/proto/ui';

import ArcaneApl from './apls/arcane.apl.json';
import BlankAPL from './apls/default.apl.json';
import FireApl from './apls/fire.apl.json';
import FireLowRankApl from './apls/fire_lowrank.apl.json';
import FrostApl from './apls/frost.apl.json';
import FrostfireApl from './apls/frostfire.apl.json';
import ArcaneLaunchGear from './gear_sets/arcane_launch.gear.json';
import FireLaunchGear from './gear_sets/fire_launch.gear.json';
import FrostLaunchGear from './gear_sets/frost_launch.gear.json';
import P0BisGear from './gear_sets/p0.bis.gear.json';
import P1BisGear from './gear_sets/p1.bis.gear.json';

// Defaults are master's ui/mage (the Forever site before the switch), on our Forever APLs.

export const BLANK_APL = PresetUtils.makePresetAPLRotation('Blank', BlankAPL);

export const ROTATION_PRESET_FROST = PresetUtils.makePresetAPLRotation('Frost', FrostApl);
export const ROTATION_PRESET_ARCANE = PresetUtils.makePresetAPLRotation('Arcane', ArcaneApl);
export const ROTATION_PRESET_FIRE = PresetUtils.makePresetAPLRotation('Fire', FireApl);
export const ROTATION_PRESET_FIRE_LOWRANK = PresetUtils.makePresetAPLRotation('Fire (low rank)', FireLowRankApl);
export const ROTATION_PRESET_FROSTFIRE = PresetUtils.makePresetAPLRotation('Frostfire hybrid (experimental)', FrostfireApl);
export const ROTATION_PRESETS = [ROTATION_PRESET_FROSTFIRE, ROTATION_PRESET_FROST, ROTATION_PRESET_ARCANE, ROTATION_PRESET_FIRE, ROTATION_PRESET_FIRE_LOWRANK];

export const ArcaneMageSimpleRotation = Mage_Rotation.create({
	conserveStart: 20,
	conserveEnd: 30,
	delayMajorCDs: 10,
});

export const APL_ARCANE_SIMPLE = PresetUtils.makePresetSimpleRotation('Simple', Spec.SpecMage, ArcaneMageSimpleRotation);

export const TalentsP1Frost = PresetUtils.makePresetTalents('Frost DPS', SavedTalents.create({ talentsString: '05002050030003--055500033100030024' }));
export const TalentsP1Arcane = PresetUtils.makePresetTalents('Arcane DPS', SavedTalents.create({ talentsString: '0501215003100301531-2305003202003-' }));
export const TalentsP1Fire = PresetUtils.makePresetTalents('Fire DPS', SavedTalents.create({ talentsString: '05002252000003-23550000130135031-' }));
// The community builds our Forever sim ranks: Arcane 35/0/16, Fire 0/35/16 and Frost 14/0/37.
export const FireTalents = PresetUtils.makePresetTalents('Fire 0/35/16', SavedTalents.create({ talentsString: '-03552020130135131-005500033' }));
export const FrostTalents = PresetUtils.makePresetTalents('Frost 14/0/37', SavedTalents.create({ talentsString: '0500005013--0555300301001301251' }));
export const ArcaneTalents = PresetUtils.makePresetTalents('Arcane 35/0/16', SavedTalents.create({ talentsString: '0551005023100301531--005500033' }));
export const TALENT_PRESETS = [TalentsP1Frost, TalentsP1Arcane, TalentsP1Fire, FireTalents, FrostTalents, ArcaneTalents];
export const DefaultTalents = TalentsP1Frost;

// Master defaults to Molten Armor, which Forever does not have (master applies nothing for it).
// Mage Armor is the Forever armor a caster runs, and what tools/parity sims.
export const DefaultOptions = MageOptions.create({
	classOptions: {
		defaultMageArmor: MageArmor.MageArmorMageArmor,
	},
});

export const OtherDefaults = {
	reactionTime: 200, // master's default
	distanceFromTarget: 20,
	profession1: Profession.Alchemy,
	profession2: Profession.Tailoring,
	race: Race.RaceTroll,
};

// Master's consumables, as the Forever client's items. Master's Greater Firepower is item 21546,
// which the Forever client makes an Elixir of Holy Power; the Frost default takes Frost Power.
export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 13512, // Flask of Supreme Power
	spellPowerElixirId: 13454, // Greater Arcane Elixir
	schoolElixirId: 17708, // Elixir of Frost Power
	guardianElixirId: 20007, // Mageblood Elixir
	zanzaId: 8423, // Cerebral Cortex Compound
	foodId: 18254, // Runn Tum Tuber Surprise
	potId: 13444, // Major Mana Potion
	conjuredId: 12662, // Demonic Rune
	mhImbueId: 25122, // Brilliant Wizard Oil
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	prayerOfSpirit: true,
	giftOfTheWild: true,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	moonkinAura: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	greaterBlessingOfWisdom: true,
});

// Improved Scorch and Winter's Chill only help the mage that applied them in Forever. Master's
// Judgement of Wisdom default does not survive its own page load, so the mage opens with none.
export const DefaultDebuffs = Debuffs.create({});

// Our Forever sim's gear presets (master ui/<spec>/gear_sets).
// One launch set per tree (tools/launch_gear): each is valued for its own school's damage.
export const GEAR_FROST_LAUNCH = PresetUtils.makePresetGear('Frost Launch', FrostLaunchGear);
export const GEAR_FIRE_LAUNCH = PresetUtils.makePresetGear('Fire Launch', FireLaunchGear);
export const GEAR_ARCANE_LAUNCH = PresetUtils.makePresetGear('Arcane Launch', ArcaneLaunchGear);
export const GEAR_P0_BIS = PresetUtils.makePresetGear('Pre-BiS', P0BisGear);
export const GEAR_P1_BIS = PresetUtils.makePresetGear('P1 BiS', P1BisGear);
export const DEFAULT_GEAR = GEAR_P0_BIS;
export const GEAR_PRESETS = [GEAR_FROST_LAUNCH, GEAR_FIRE_LAUNCH, GEAR_ARCANE_LAUNCH, GEAR_P0_BIS, GEAR_P1_BIS];
