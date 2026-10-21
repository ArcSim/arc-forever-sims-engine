import { MageTalents } from '@generated/proto/mage';
import { talentStringToProto } from '@sim/talents/factory';
import { mageTalentsConfig } from '@sim/talents/mage';
import { describe, expect, it } from 'vitest';

import { parseTalentsString, serializeTalentsString, totalPointsSpent, treePointTotal, withTalentPoints, withTreeCleared } from './talents_string';

// A valid Forever mage build: 51 points, Fire left empty so the codec's empty-run
// handling stays covered.
const MAGE_DEFAULT = '2552252231221--25352';

describe('talents string codec', () => {
	it('round-trips a real Forever default byte-identically, empty middle tree and all', () => {
		expect(serializeTalentsString(parseTalentsString(mageTalentsConfig, MAGE_DEFAULT))).toBe(MAGE_DEFAULT);
	});

	it('keeps the empty middle tree as an empty run rather than collapsing the dashes', () => {
		const points = parseTalentsString(mageTalentsConfig, MAGE_DEFAULT);
		expect(treePointTotal(points, 1)).toBe(0);
		expect(
			serializeTalentsString(points)
				.split('-')
				.map(tree => tree.length),
		).toEqual([13, 0, 5]);
	});

	it('widens a short string to the config shape and trims it back', () => {
		const points = parseTalentsString(mageTalentsConfig, '25');
		expect(points[0].length).toBe(mageTalentsConfig[0].talents.length);
		expect(points[1].every(value => value === 0)).toBe(true);
		expect(serializeTalentsString(points)).toBe('25');
	});

	it('trims trailing zeroes and trailing empty trees', () => {
		expect(
			serializeTalentsString([
				[1, 0, 0],
				[0, 0, 0],
				[0, 0, 0],
			]),
		).toBe('1');
		expect(
			serializeTalentsString([
				[0, 0, 0],
				[0, 0, 0],
				[0, 0, 2],
			]),
		).toBe('--002');
	});

	it('clamps a digit above the talent maximum, as setPoints did', () => {
		const firstMax = mageTalentsConfig[0].talents[0].maxPoints;
		expect(parseTalentsString(mageTalentsConfig, '9')[0][0]).toBe(firstMax);
	});

	it('reads a non-digit as unspent', () => {
		expect(parseTalentsString(mageTalentsConfig, 'x5')[0][0]).toBe(0);
		expect(parseTalentsString(mageTalentsConfig, 'x5')[0][1]).toBe(5);
	});

	it('replaces one talent without disturbing the rest', () => {
		const points = parseTalentsString(mageTalentsConfig, MAGE_DEFAULT);
		const next = withTalentPoints(points, 2, 10, 3);
		expect(next[2][10]).toBe(3);
		expect(next[0]).toEqual(points[0]);
		expect(totalPointsSpent(next)).toBe(totalPointsSpent(points) + 3);
	});

	it('clears exactly one tree', () => {
		const points = parseTalentsString(mageTalentsConfig, MAGE_DEFAULT);
		const next = withTreeCleared(points, 0);
		expect(treePointTotal(next, 0)).toBe(0);
		expect(treePointTotal(next, 2)).toBe(treePointTotal(points, 2));
		expect(serializeTalentsString(next)).toBe('--25352');
	});

	it('counts every point in a full default build', () => {
		expect(totalPointsSpent(parseTalentsString(mageTalentsConfig, MAGE_DEFAULT))).toBe(51);
	});

	it('reads a build longer than the tree without throwing', () => {
		// Arcane holds 18 talents: pad to 18 before overflowing, or the digits land inside the tree.
		const overlong = MAGE_DEFAULT.replace('2552252231221', '2552252231221' + '0'.repeat(5) + '555555');
		expect(() => talentStringToProto(MageTalents.create(), overlong, mageTalentsConfig)).not.toThrow();
		expect(talentStringToProto(MageTalents.create(), overlong, mageTalentsConfig)).toEqual(
			talentStringToProto(MageTalents.create(), MAGE_DEFAULT, mageTalentsConfig),
		);
	});

	it('reads a build shorter than the tree, and an empty one', () => {
		expect(() => talentStringToProto(MageTalents.create(), '25', mageTalentsConfig)).not.toThrow();
		expect(talentStringToProto(MageTalents.create(), '', mageTalentsConfig)).toEqual(MageTalents.create());
	});
});
