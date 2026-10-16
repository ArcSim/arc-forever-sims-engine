import * as PresetUtils from '@app/preset_utils';
import { Debuffs, IndividualBuffs, PartyBuffs, RaidBuffs } from '@generated/proto/buffs';
import { ConsumesSpec, Profession, Race } from '@generated/proto/common';
import { HealerPriest_Options as HealerPriestOptions, PriestOptions_Armor } from '@generated/proto/priest';
import { SavedTalents } from '@generated/proto/ui';

// Defaults follow master's ui/healing_priest.
export const DefaultOptions = HealerPriestOptions.create({
	classOptions: {
		armor: PriestOptions_Armor.InnerFire,
	},
});

// Master sets no consumables for this spec (the old TBC flask, food and potion don't exist in the Forever client).
export const DefaultConsumables = ConsumesSpec.create({});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	giftOfTheWild: true,
	prayerOfFortitude: true,
	prayerOfSpirit: true,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	moonkinAura: true,
	strengthOfEarthTotem: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	greaterBlessingOfKings: true,
	greaterBlessingOfWisdom: true,
});

export const DefaultDebuffs = Debuffs.create({});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Tailoring,
	profession2: Profession.Enchanting,
	race: Race.RaceDwarf,
};

// Talent presets, from master's ui/priest spec.
export const TalentsHolyHealer = PresetUtils.makePresetTalents('Holy 19/32/0', SavedTalents.create({ talentsString: '00520303123-2350510323000053' }));
export const TalentsDisciplineHealer = PresetUtils.makePresetTalents(
	'Discipline 35/16/0',
	SavedTalents.create({ talentsString: '005205031332011531-03505003' }),
);
export const TalentPresets = [TalentsHolyHealer, TalentsDisciplineHealer];
