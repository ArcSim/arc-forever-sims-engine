import { Phase } from '../core/constants/other.js';
import * as PresetUtils from '../core/preset_utils.js';
import {
	AgilityElixir,
	Alcohol,
	AttackPowerBuff,
	Conjured,
	Consumes,
	Debuffs,
	Food,
	HealthElixir,
	IndividualBuffs,
	ManaRegenElixir,
	Potions,
	Profession,
	Race,
	RaidBuffs,
	SaygesFortune,
	SpellPowerBuff,
	StrengthBuff,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common.js';
import {
	Hunter_Options as HunterOptions,
	Hunter_Options_Ammo as Ammo,
	Hunter_Options_PetAttackSpeed as PetAttackSpeed,
	Hunter_Options_PetType as PetType,
	Hunter_Options_QuiverBonus,
} from '../core/proto/hunter.js';
import { SavedTalents } from '../core/proto/ui.js';

import SurvDWGear from './gear_sets/surv_dw.gear.json';
import TwoHWeaveGear from './gear_sets/2h_weave.gear.json';
import MMGear from './gear_sets/mm_weave.gear.json';

import TwoHWeaveAPL from './apls/2h_weave.apl.json';
import EasyWeaveAPL from './apls/easy_weave.apl.json';
import SurvDWAPL from './apls/SurvDW.apl.json';
import MMSniperAPL from './apls/mm_sniper.apl.json';

export const APL2HWeave = PresetUtils.makePresetAPLRotation('2H Max Weave', TwoHWeaveAPL);
export const APLSurvDW = PresetUtils.makePresetAPLRotation('Surv DW', SurvDWAPL);
export const APLMMSniper = PresetUtils.makePresetAPLRotation('MM Sniper', MMSniperAPL);
export const APLEasyWeave = PresetUtils.makePresetAPLRotation('2H Easy Weave', EasyWeaveAPL);

export const APLPresets = {
	[Phase.Phase1]: [APL2HWeave, APLSurvDW, APLMMSniper, APLEasyWeave],
};

export const DefaultAPL = APL2HWeave;

export const Gear2HWeave = PresetUtils.makePresetGear('2H Weave', TwoHWeaveGear);
export const GearSurvDW = PresetUtils.makePresetGear('Surv DW', SurvDWGear);
export const GearMMTurret = PresetUtils.makePresetGear('MM Turret', MMGear);

export const GearPresets = {
	[Phase.Phase1]: [GearSurvDW, Gear2HWeave, GearMMTurret],
};

export const DefaultGear = Gear2HWeave;



export const SurvDWHawk = PresetUtils.makePresetTalents(
	'Surv DW Hawk',
	SavedTalents.create({ talentsString: '53200005001-0050052-5002302300500201' }),
);
export const SurvDWLacerate = PresetUtils.makePresetTalents(
	'Surv DW LS',
	SavedTalents.create({ talentsString: '5-005005201-500230230050222151' }),
);
export const TwoHWeave = PresetUtils.makePresetTalents(
	'2H Weave',
	SavedTalents.create({ talentsString: '-00531500114-500230230050220151' }),
);
export const TalentsMMHawk = PresetUtils.makePresetTalents(
	'MM Hawk',
	SavedTalents.create({ talentsString: '53200005001-005055200150205-5' }),
);
export const TalentsMMSniper = PresetUtils.makePresetTalents(
	'MM Sniper',
	SavedTalents.create({ talentsString: '5-0053552511503051-5' }),
);
export const TalentPresets = {
	[Phase.Phase1]: [SurvDWHawk, SurvDWLacerate, TwoHWeave, TalentsMMHawk, TalentsMMSniper],
};

export const DefaultTalents = TwoHWeave;

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	battleShout: TristateEffect.TristateEffectRegular,
	divineSpirit: true,
	fireResistanceAura: true,
	fireResistanceTotem: true,
	giftOfTheWild: TristateEffect.TristateEffectRegular,
	leaderOfThePack: true,
	manaSpringTotem: TristateEffect.TristateEffectRegular,
	strengthOfEarthTotem: TristateEffect.TristateEffectRegular,
	trueshotAura: true,
	windfuryTotem: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfMight: TristateEffect.TristateEffectRegular,
	blessingOfWisdom: TristateEffect.TristateEffectRegular,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfElements: true,
	faerieFire: true,
	giftOfArthas: true,
	huntersMark: TristateEffect.TristateEffectRegular,
	judgementOfWisdom: true,
	sunderArmor: true,
});

export const DefaultConsumes = Consumes.create({
	agilityElixir: AgilityElixir.ElixirOfTheMongoose,
	alcohol: Alcohol.AlcoholRumseyRumBlackLabel,
	attackPowerBuff: AttackPowerBuff.WinterfallFirewater,
	defaultConjured: Conjured.ConjuredDemonicRune,
	defaultPotion: Potions.MajorManaPotion,
	dragonBreathChili: true,
	food: Food.FoodSmokedDesertDumpling,
	healthElixir: HealthElixir.ElixirOfFortitude,
	mainHandImbue: WeaponImbue.DenseSharpeningStone,
	manaRegenElixir: ManaRegenElixir.MagebloodPotion,
	miscConsumes: { raptorPunch: true },
	offHandImbue: WeaponImbue.DenseSharpeningStone,
	spellPowerBuff: SpellPowerBuff.GreaterArcaneElixir,
	strengthBuff: StrengthBuff.ElixirOfGiants,
	zanzaBuff: ZanzaBuff.GroundScorpokAssay,
});

export const DefaultOptions = HunterOptions.create({
	ammo: Ammo.ThoriumHeadedArrow,
	petAttackSpeed: PetAttackSpeed.OneTwo,
	petType: PetType.Cat,
	petUptime: 1,
	quiverBonus: Hunter_Options_QuiverBonus.Speed15,
});

export const OtherDefaults = {
	distanceFromTarget: 3,
	profession1: Profession.Enchanting,
	profession2: Profession.Engineering,
	race: Race.RaceOrc,
};