import * as InputHelpers from '../core/components/input_helpers.js';
import { Player } from '../core/player.js';
import { Spec } from '../core/proto/common.js';
import {
	Hunter_Options_Ammo as Ammo,
	Hunter_Options_PetAttackSpeed as PetAttackSpeed,
	Hunter_Options_PetType,
	Hunter_Options_QuiverBonus as QuiverBonus,
	Hunter_Rotation_RotationType as RotationType,
	Hunter_Rotation_StingType as StingType,
} from '../core/proto/hunter.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { makePetTypeInputConfig } from '../core/talents/hunter_pet.js';
import { TypedEvent } from '../core/typed_event.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const WeaponAmmo = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecHunter, Ammo>({
	fieldName: 'ammo',
	values: [
		{ value: Ammo.AmmoNone, tooltip: 'No Ammo' },
		{ actionId: () => ActionId.fromItemId(12654), value: Ammo.Doomshot },
		{ actionId: () => ActionId.fromItemId(19317), value: Ammo.IceThreadedBullet },
		{ actionId: () => ActionId.fromItemId(15997), value: Ammo.ThoriumShells },
		{ actionId: () => ActionId.fromItemId(13377), value: Ammo.MiniatureCannonBalls },
		{ actionId: () => ActionId.fromItemId(274387), value: Ammo.SwiftfeatherArrow },
		{ actionId: () => ActionId.fromItemId(274388), value: Ammo.SwiftstrikeShot },
	],
});

export const QuiverInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecHunter, QuiverBonus>({
	extraCssClasses: ['quiver-picker'],
	fieldName: 'quiverBonus',
	numColumns: 2,
	values: [
		{ color: '82e89d', value: QuiverBonus.QuiverNone },
		{ actionId: () => ActionId.fromItemId(18714), value: QuiverBonus.Speed15 },
		{ actionId: () => ActionId.fromItemId(2662), value: QuiverBonus.Speed14 },
		{ actionId: () => ActionId.fromItemId(8217), value: QuiverBonus.Speed13 },
		{ actionId: () => ActionId.fromItemId(7371), value: QuiverBonus.Speed12 },
		{ actionId: () => ActionId.fromItemId(3605), value: QuiverBonus.Speed11 },
		{ actionId: () => ActionId.fromItemId(3573), value: QuiverBonus.Speed10 },
	],
});

export const PetTypeInput = makePetTypeInputConfig(true);

export const PetUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecHunter>({
	fieldName: 'petUptime',
	label: 'Pet Uptime (%)',
	labelTooltip: 'Percent of the fight duration for which your pet will be alive.',
	percent: true,
});

export const PetAttackSpeedInput = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecHunter>({
	fieldName: 'petAttackSpeed',
	label: 'Pet Attack Speed',
	labelTooltip: 'The pets auto attacks speed.',
	values: [
		{ name: '1.0', value: PetAttackSpeed.One },
		{ name: '1.2', value: PetAttackSpeed.OneTwo },
		{ name: '1.3', value: PetAttackSpeed.OneThree },
		{ name: '1.4', value: PetAttackSpeed.OneFour },
		{ name: '1.5', value: PetAttackSpeed.OneFive },
		{ name: '1.6', value: PetAttackSpeed.OneSix },
		{ name: '1.7', value: PetAttackSpeed.OneSeven },
		{ name: '2.0', value: PetAttackSpeed.Two },
		{ name: '2.4', value: PetAttackSpeed.TwoFour },
		{ name: '2.5', value: PetAttackSpeed.TwoFive },
	],
	showWhen: player => player.getSpecOptions().petType != Hunter_Options_PetType.PetNone,
	changeEmitter: (player: Player<Spec.SpecHunter>) => TypedEvent.onAny([player.specOptionsChangeEmitter]),
});

export const HunterRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecHunter>({
			fieldName: 'type',
			label: 'Type',
			values: [
				{ name: 'Single Target', value: RotationType.SingleTarget },
				{ name: 'AOE', value: RotationType.Aoe },
			],
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecHunter>({
			fieldName: 'sting',
			label: 'Sting',
			labelTooltip: 'Maintains the selected Sting on the primary target.',
			values: [
				{ name: 'None', value: StingType.NoSting },
				{ name: 'Scorpid Sting', value: StingType.ScorpidSting },
				{ name: 'Serpent Sting', value: StingType.SerpentSting },
			],
			showWhen: (player: Player<Spec.SpecHunter>) => player.getSimpleRotation().type == RotationType.SingleTarget,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecHunter>({
			fieldName: 'multiDotSerpentSting',
			label: 'Multi-Dot Serpent Sting',
			labelTooltip: 'Casts Serpent Sting on multiple targets',
			changeEmitter: (player: Player<Spec.SpecHunter>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
	],
};
