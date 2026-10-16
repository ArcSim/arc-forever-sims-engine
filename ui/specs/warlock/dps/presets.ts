import * as PresetUtils from '@app/preset_utils';
import { Debuffs, IndividualBuffs, PartyBuffs, RaidBuffs } from '@generated/proto/buffs';
import { ConsumesSpec, Profession, Race } from '@generated/proto/common';
import { SavedTalents } from '@generated/proto/ui';
import { Warlock_Options as WarlockOptions, WarlockOptions_Armor, WarlockOptions_CurseOptions, WarlockOptions_Summon } from '@generated/proto/warlock';

import AfflictionRot from './apls/affliction.apl.json';
import BlankAPL from './apls/default.apl.json';
import DemonicPactRot from './apls/demonic_pact.apl.json';
import ShadowAndFlameRot from './apls/destruction_conflag.apl.json';
import DSRuinRot from './apls/ds_ruin.apl.json';
import LaunchGear from './gear_sets/launch.gear.json';
import McGear from './gear_sets/mc.gear.json';
import PrebisGear from './gear_sets/prebis.gear.json';

// Defaults are master's ui/warlock (the Forever site before the switch), on our Forever APLs.

export const BLANK_APL = PresetUtils.makePresetAPLRotation('Blank', BlankAPL);

// Rotations, in master's order. Summoning and sacrificing the demon is the sacrificeSummon
// option here, not a prepull cast.
export const RotationDemonicPact = PresetUtils.makePresetAPLRotation('Demonic Pact', DemonicPactRot);
export const RotationAffliction = PresetUtils.makePresetAPLRotation('Affliction', AfflictionRot);
export const RotationDSRuin = PresetUtils.makePresetAPLRotation('DS/Ruin', DSRuinRot);
// Shadow and Flame 5/5 keeps Immolate up through Conflagrate, so its rotation casts it.
export const RotationShadowAndFlame = PresetUtils.makePresetAPLRotation('Shadow and Flame', ShadowAndFlameRot);
export const APLPresets = [RotationDemonicPact, RotationAffliction, RotationDSRuin, RotationShadowAndFlame];

// The Imp, sacrificed by the DS/Ruin default build (master's DS/Ruin rotation sacrifices it).
export const DefaultOptions = WarlockOptions.create({
	classOptions: {
		armor: WarlockOptions_Armor.DemonArmor,
		curseOptions: WarlockOptions_CurseOptions.Elements,
		sacrificeSummon: true,
		summon: WarlockOptions_Summon.Imp,
	},
});

// Without pet talents the Succubus out-damages the Imp, so the Affliction builds run one.
export const AfflictionOptions = WarlockOptions.create({
	classOptions: {
		armor: WarlockOptions_Armor.DemonArmor,
		curseOptions: WarlockOptions_CurseOptions.Elements,
		summon: WarlockOptions_Summon.Succubus,
	},
});

// Demonic Pact keeps the Succubus out for Master Demonologist and Soul Link, and keeps the buff of
// an Imp sacrificed before the pull (+15% Shadow), which beats Master's Voidwalker (mana) by ~7%.
export const DemonicPactOptions = WarlockOptions.create({
	classOptions: {
		...AfflictionOptions.classOptions!,
		pactSacrifice: WarlockOptions_Summon.Imp,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 13512, // Flask of Supreme Power
	spellPowerElixirId: 13454, // Greater Arcane Elixir
	schoolElixirId: 9264, // Elixir of Shadow Power
	guardianElixirId: 20007, // Mageblood Elixir
	zanzaId: 8423, // Cerebral Cortex Compound
	alcoholId: 21151, // Rumsey Rum Black Label
	foodId: 18254, // Runn Tum Tuber Surprise
	potId: 13444, // Major Mana Potion
	conjuredId: 12662, // Demonic Rune
});

export const OtherDefaults = {
	reactionTime: 200, // master's default
	distanceFromTarget: 25,
	profession1: Profession.Enchanting,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
	race: Race.RaceGnome,
};

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	giftOfTheWild: true,
	prayerOfFortitude: true,
	prayerOfSpirit: true,
	fireResistanceAura: true, // a raid buff in the new buffs proto
});

// Master opens as Alliance, so the Horde totems its presets name are not applied.
export const DefaultPartyBuffs = PartyBuffs.create({
	moonkinAura: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	greaterBlessingOfKings: true,
	greaterBlessingOfWisdom: true,
});

export const DefaultDebuffs = Debuffs.create({
	exposeArmor: true,
	faerieFire: true,
	judgementOfWisdom: true,
	sunderArmor: true,
});

// Talent presets, from master's ui/warlock spec.
export const TalentsDemonicPact = PresetUtils.makePresetTalents('Demonic Pact', SavedTalents.create({ talentsString: '203-0055003221201001351-0550005' }));
export const TalentsAffliction = PresetUtils.makePresetTalents('Affliction', SavedTalents.create({ talentsString: '2552002013520135--0500055' }));
export const TalentsDSRuin = PresetUtils.makePresetTalents('DS/Ruin', SavedTalents.create({ talentsString: '255100201332-0340003001-0550105' }));
export const TalentsPactOptimised = PresetUtils.makePresetTalents(
	'Demonic Pact 2/31/18',
	SavedTalents.create({ talentsString: '113-0005003221220311351-0550005' }),
);
export const TalentsDeepAffliction = PresetUtils.makePresetTalents(
	'Deep Affliction 35/0/16',
	SavedTalents.create({ talentsString: '2552012013521105--05000551' }),
);
export const TalentsDSRuinPandemic = PresetUtils.makePresetTalents(
	'DS/Ruin Pandemic 24/11/16',
	SavedTalents.create({ talentsString: '25220010135201-0025003001-05500051' }),
);
export const TalentsShadowAndFlame = PresetUtils.makePresetTalents(
	'Shadow and Flame 13/11/27',
	SavedTalents.create({ talentsString: '25501-0025003001-055035510010002' }),
);
export const TalentPresets = [
	TalentsDemonicPact,
	TalentsAffliction,
	TalentsDSRuin,
	TalentsPactOptimised,
	TalentsDeepAffliction,
	TalentsDSRuinPandemic,
	TalentsShadowAndFlame,
];
export const DefaultTalents = TalentsDSRuin;

// The community builds with the pet setup and rotation each one is measured with.
const buildOptions = (name: string, specOptions: WarlockOptions) => ({ settings: { name, specOptions } });
export const BuildDemonicPact = PresetUtils.makePresetBuild('Demonic Pact 2/31/18', {
	talents: TalentsPactOptimised,
	rotation: RotationDemonicPact,
	...buildOptions('Demonic Pact 2/31/18', DemonicPactOptions),
});
export const BuildDeepAffliction = PresetUtils.makePresetBuild('Deep Affliction 35/0/16', {
	talents: TalentsDeepAffliction,
	rotation: RotationAffliction,
	...buildOptions('Deep Affliction 35/0/16', AfflictionOptions),
});
export const BuildDSRuinPandemic = PresetUtils.makePresetBuild('DS/Ruin Pandemic 24/11/16', {
	talents: TalentsDSRuinPandemic,
	rotation: RotationDSRuin,
	...buildOptions('DS/Ruin Pandemic 24/11/16', DefaultOptions),
});
export const BuildShadowAndFlame = PresetUtils.makePresetBuild('Shadow and Flame 13/11/27', {
	talents: TalentsShadowAndFlame,
	rotation: RotationShadowAndFlame,
	...buildOptions('Shadow and Flame 13/11/27', DefaultOptions),
});
export const BuildPresets = [BuildDemonicPact, BuildDeepAffliction, BuildDSRuinPandemic, BuildShadowAndFlame];

// Our Forever sim's gear presets (master ui/<spec>/gear_sets).
export const GEAR_BLANK = PresetUtils.makePresetGear('Blank', { items: [] });
export const GEAR_LAUNCH = PresetUtils.makePresetGear('Launch', LaunchGear);
export const GEAR_PREBIS = PresetUtils.makePresetGear('Pre-BIS', PrebisGear);
export const GEAR_MC = PresetUtils.makePresetGear('MC', McGear);
export const DEFAULT_GEAR = GEAR_PREBIS;
export const GEAR_PRESETS = [GEAR_BLANK, GEAR_LAUNCH, GEAR_PREBIS, GEAR_MC];
