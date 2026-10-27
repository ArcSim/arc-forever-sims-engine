// Pre-raid gear sweep for one shipped build and play mode. Stat weights shortlist each slot's candidates
// (preraid_items.js decides what counts as pre-raid), every legal single swap is screened against the current set,
// the promising ones are re-simmed at high iterations, and the confirmed swaps are combined and checked together.
// Repeats from the new set until a round confirms nothing.
// usage: node gear_screen.js <build> <God|Pro|Regular> <out.json>
// env: ENGINE_DIR (engine checkout), HUNTER_TEST (built test binary), SCREEN_ITERS, CONFIRM_ITERS, PER_SLOT, WORKERS, ROUNDS
const fs = require('fs');
const os = require('os');
const path = require('path');
const { spawn } = require('child_process');
const PR = require('./preraid_items.js');

const [build, mode, outPath] = process.argv.slice(2);
const ENGINE = process.env.ENGINE_DIR || path.join(__dirname, '..', 'engine');
const EXE = process.env.HUNTER_TEST || path.join(__dirname, '..', 'hunter_test');
const SCREEN = Number(process.env.SCREEN_ITERS) || 6000;
const CONFIRM = Number(process.env.CONFIRM_ITERS) || 30000;
const PER_SLOT = Number(process.env.PER_SLOT) || 8;
const WORKERS = Number(process.env.WORKERS) || os.cpus().length;
const ROUNDS = Number(process.env.ROUNDS) || 2;
// a swap must hold this much at CONFIRM iterations to count
const KEEP = 1.0;
const PLAY = { God: 'WeaveGodMode', Pro: 'WeavePerfectPlay', Regular: 'WeaveNormalPlay' };

const db = JSON.parse(fs.readFileSync(path.join(ENGINE, 'assets', 'database', 'db.json'), 'utf8'));
const byId = new Map(db.items.map((i) => [i.id, i]));
const name = (id) => (byId.get(id) || {}).name || (id ? String(id) : '(empty)');
let settings = JSON.parse(fs.readFileSync(path.join(ENGINE, 'builds', `ship_${build}.json`), 'utf8'));
settings.player.hunter.options.weavePlay = PLAY[mode];
const race = settings.player.race;

const S = { str: 0, agi: 1, int: 3, ap: 17, hit: 18, crit: 19, rap: 27 };
const SLOTS = ['Head', 'Neck', 'Shoulder', 'Back', 'Chest', 'Wrist', 'Hands', 'Waist', 'Legs', 'Feet', 'Finger 1', 'Finger 2', 'Trinket 1', 'Trinket 2', 'Main hand', 'Off hand', 'Ranged'];
const SLOT_TYPE = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 11, 12, 12, 13, 13, 14];
const ENCHANT_1H = 2564; // +15 Agility
const ENCHANT_2H = 2646; // +25 Agility
const rangedOnly = build === 'bm_ranged';

const startMH = byId.get((settings.player.equipment.items[14] || {}).id);
const keepDW = startMH && startMH.handType !== 4 && !!(settings.player.equipment.items[15] || {}).id;
const keep2H = startMH && startMH.handType === 4;

const TMP = fs.mkdtempSync(path.join(os.tmpdir(), 'sweep-'));
let seq = 0;
function sim(s, iters) {
	const file = path.join(TMP, `${seq++}.json`);
	fs.writeFileSync(file, JSON.stringify(s));
	return new Promise((res) => {
		const p = spawn(EXE, ['-test.run', 'TestForeverSettings', '-test.v'], {
			cwd: path.join(ENGINE, 'sim', 'hunter'),
			env: { ...process.env, AB_SETTINGS: file, AB_ITERS: String(iters) },
		});
		let out = '';
		p.stdout.on('data', (d) => (out += d));
		p.on('close', () => {
			fs.rmSync(file, { force: true });
			res(Number((out.match(/raid DPS incl\. pets ([0-9.]+)/) || [])[1]));
		});
	});
}
async function pool(tasks) {
	const out = new Array(tasks.length);
	let next = 0;
	await Promise.all(
		Array.from({ length: WORKERS }, async () => {
			while (next < tasks.length) {
				const i = next++;
				out[i] = await tasks[i]();
			}
		}),
	);
	return out;
}

const clone = (o) => JSON.parse(JSON.stringify(o));
function equip(s, slot, id) {
	const v = clone(s);
	const items = v.player.equipment.items;
	const keep = items[slot] && items[slot].enchant;
	const item = byId.get(id);
	items[slot] = { id };
	if (slot === 14 || slot === 15) {
		items[slot].enchant = item.handType === 4 ? ENCHANT_2H : ENCHANT_1H;
		if (slot === 14 && item.handType === 4) items[15] = {};
	} else if (keep) items[slot].enchant = keep;
	return v;
}
function legal(s) {
	const items = s.player.equipment.items;
	const ids = items.map((i) => i && i.id).filter(Boolean);
	// trinkets count as unique: nearly every Classic trinket is, and the database misses the flag on some
	const unique = (id) => byId.get(id) && (byId.get(id).unique || byId.get(id).type === 12);
	for (const id of ids) if (unique(id) && ids.filter((x) => x === id).length > 1) return false;
	const mh = byId.get(items[14] && items[14].id);
	const oh = byId.get(items[15] && items[15].id);
	if (oh && (!mh || mh.handType === 4)) return false;
	if (oh && oh.handType === 1) return false;
	return true;
}

async function weights(s) {
	const probes = [['agi', S.agi, 40], ['str', S.str, 40], ['int', S.int, 40], ['ap', S.ap, 80], ['rap', S.rap, 80], ['hit', S.hit, 2], ['crit', S.crit, 2]];
	const variants = [s, ...probes.map(([, idx, amount]) => {
		const v = clone(s);
		const stats = ((v.player.bonusStats = v.player.bonusStats || {}).stats || new Array(44).fill(0)).slice();
		stats[idx] += amount;
		v.player.bonusStats.stats = stats;
		return v;
	})];
	const d = await pool(variants.map((v) => () => sim(v, SCREEN * 2)));
	const w = {};
	probes.forEach(([k, , amount], i) => (w[k] = Math.max(0, (d[i + 1] - d[0]) / amount)));
	return w;
}
const ep = (item, w) => {
	const st = item.stats || [];
	return Object.entries(S).reduce((a, [k, idx]) => a + (st[idx] || 0) * (w[k] || 0), 0);
};
const weaponDps = (item) => ((item.weaponDamageMin || 0) + (item.weaponDamageMax || 0)) / 2 / (item.weaponSpeed || 1);
const legalPool = db.items.filter((i) => PR.preRaid(i, race));

function shortlist(slot, w, s) {
	const type = SLOT_TYPE[slot];
	let c = legalPool.filter((i) => i.type === type);
	if (slot === 15) c = keep2H ? [] : c.filter((i) => i.handType !== 1 && i.handType !== 4);
	if (slot === 14) c = c.filter((i) => i.handType !== 3 && !(keepDW && i.handType === 4) && !(keep2H && i.handType !== 4));
	const score = (i) => ep(i, w) + (type === 14 ? weaponDps(i) * 4 : 0) + (type === 13 && !rangedOnly ? weaponDps(i) * 2 : 0);
	const cur = (s.player.equipment.items[slot] || {}).id;
	return c.filter((i) => i.id !== cur).sort((a, b) => score(b) - score(a)).slice(0, type === 12 ? PER_SLOT + 3 : PER_SLOT);
}

(async () => {
	const t0 = Date.now();
	const report = { build, mode, race, screenIters: SCREEN, confirmIters: CONFIRM, notPreRaid: [], rounds: [] };
	settings.player.equipment.items.forEach((it, slot) => {
		if (it && it.id && byId.get(it.id) && !PR.preRaid(byId.get(it.id), race)) report.notPreRaid.push(`${SLOTS[slot]}: ${name(it.id)}`);
	});
	const shipped = clone(settings);
	const w = await weights(settings);
	report.weights = w;
	report.shippedDps = await sim(shipped, CONFIRM);
	for (let round = 1; round <= ROUNDS; round++) {
		const tries = [];
		for (let slot = 0; slot < SLOTS.length; slot++)
			for (const cand of shortlist(slot, w, settings)) {
				const v = equip(settings, slot, cand.id);
				if (legal(v)) tries.push({ slot, id: cand.id, s: v });
			}
		const base = await sim(settings, SCREEN);
		const screened = await pool(tries.map((t) => () => sim(t.s, SCREEN)));
		tries.forEach((t, i) => (t.screen = screened[i] - base));
		// the best three per slot that did not lose at screening iterations go to the high-iteration check
		const promising = [];
		for (let slot = 0; slot < SLOTS.length; slot++)
			promising.push(...tries.filter((t) => t.slot === slot && t.screen > -0.5).sort((a, b) => b.screen - a.screen).slice(0, 3));
		const confirmBase = await sim(settings, CONFIRM);
		const confirmed = await pool(promising.map((t) => () => sim(t.s, CONFIRM)));
		promising.forEach((t, i) => (t.confirm = confirmed[i] - confirmBase));
		const keep = promising.filter((t) => t.confirm >= KEEP).sort((a, b) => b.confirm - a.confirm);
		// combine the best confirmed swap per slot, skipping any that would break a unique item or the hands
		let combined = clone(settings);
		const taken = [];
		for (const t of keep) {
			if (taken.some((x) => x.slot === t.slot || ((x.slot === 14 || x.slot === 15) && (t.slot === 14 || t.slot === 15)))) continue;
			const v = equip(combined, t.slot, t.id);
			if (!legal(v)) continue;
			combined = v;
			taken.push(t);
		}
		let combinedDps = taken.length ? await sim(combined, CONFIRM) : confirmBase;
		// keep the set only if it beats the best single swap; otherwise take that swap alone
		if (taken.length > 1 && combinedDps < confirmBase + taken[0].confirm) {
			combined = taken[0].s;
			combinedDps = confirmBase + taken[0].confirm;
			taken.splice(1);
		}
		report.rounds.push({
			round,
			startDps: confirmBase,
			screened: tries.length,
			checked: promising.map((t) => ({ slot: SLOTS[t.slot], from: name((settings.player.equipment.items[t.slot] || {}).id), to: name(t.id), screen: +t.screen.toFixed(1), confirm: +t.confirm.toFixed(1) })),
			taken: taken.map((t) => ({ slot: SLOTS[t.slot], from: name((settings.player.equipment.items[t.slot] || {}).id), to: name(t.id), confirm: +t.confirm.toFixed(1) })),
			endDps: combinedDps,
		});
		if (!taken.length) break;
		settings = combined;
	}
	report.finalDps = await sim(settings, CONFIRM);
	report.gain = +(report.finalDps - report.shippedDps).toFixed(1);
	report.changes = settings.player.equipment.items
		.map((it, slot) => ({ slot, it, was: shipped.player.equipment.items[slot] || {} }))
		.filter(({ it, was }) => (it.id || 0) !== (was.id || 0) || (it.enchant || 0) !== (was.enchant || 0))
		.map(({ slot, it, was }) => `${SLOTS[slot]}: ${name(was.id)} -> ${name(it.id)}`);
	report.finalGear = settings.player.equipment;
	report.minutes = +((Date.now() - t0) / 60000).toFixed(1);
	fs.mkdirSync(path.dirname(outPath), { recursive: true });
	fs.writeFileSync(outPath, JSON.stringify(report, null, 1));
	console.log(`${build} ${mode}: shipped ${report.shippedDps.toFixed(1)}, final ${report.finalDps.toFixed(1)} (${report.gain >= 0 ? '+' : ''}${report.gain}); ${report.changes.join('; ') || 'no change'}; ${report.minutes} min`);
	fs.rmSync(TMP, { recursive: true, force: true });
})();
