// Rotations for one shipped build and play mode with melee swings resetting Auto Shot: the shipped rotation with the
// reset off and on, then every rotation in rotations/ with it on; with WEAPONS, each of those main hands under the
// rotation named by ROTATION (the shipped one with the reset off as well).
// usage: node rotation_check.js <build> <God|Pro|Regular> <out.json>
// env: ENGINE_DIR (engine checkout), HUNTER_TEST (built test binary), ITERS, WORKERS, WEAPONS (comma-separated item ids),
// ROTATION (file in rotations/, default arc_reset_weave.apl.json), GEAR (file in gear/ worn in place of the shipped gear)
const fs = require('fs');
const os = require('os');
const path = require('path');
const { spawn } = require('child_process');

const [build, mode, outPath] = process.argv.slice(2);
const ENGINE = process.env.ENGINE_DIR || path.join(__dirname, '..', 'engine');
const EXE = process.env.HUNTER_TEST || path.join(__dirname, '..', 'hunter_test');
const ITERS = Number(process.env.ITERS) || 20000;
const WORKERS = Number(process.env.WORKERS) || os.cpus().length;
const WEAPONS = (process.env.WEAPONS || '').split(',').filter(Boolean).map(Number);
const ROTATION = process.env.ROTATION || 'arc_reset_weave.apl.json';
const PLAY = { God: 'WeaveGodMode', Pro: 'WeavePerfectPlay', Regular: 'WeaveNormalPlay' };
const ENCHANT_1H = 2564; // +15 Agility
const ENCHANT_2H = 2646; // +25 Agility

const db = JSON.parse(fs.readFileSync(path.join(ENGINE, 'assets', 'database', 'db.json'), 'utf8'));
const byId = new Map(db.items.map((i) => [i.id, i]));
const shipped = JSON.parse(fs.readFileSync(path.join(ENGINE, 'builds', `ship_${build}.json`), 'utf8'));
shipped.player.hunter.options.weavePlay = PLAY[mode];
if (process.env.GEAR) shipped.player.equipment = JSON.parse(fs.readFileSync(path.join(__dirname, 'gear', process.env.GEAR), 'utf8'));
const ROT_DIR = path.join(__dirname, 'rotations');
const rotations = fs.readdirSync(ROT_DIR).filter((f) => f.endsWith('.apl.json')).sort();

const TMP = fs.mkdtempSync(path.join(os.tmpdir(), 'rot-'));
let seq = 0;
function sim(s) {
	const file = path.join(TMP, `${seq++}.json`);
	fs.writeFileSync(file, JSON.stringify(s));
	return new Promise((res) => {
		const p = spawn(EXE, ['-test.run', 'TestForeverSettings', '-test.v'], {
			cwd: path.join(ENGINE, 'sim', 'hunter'),
			env: { ...process.env, AB_SETTINGS: file, AB_ITERS: String(ITERS) },
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
function variant({ reset, rotation, weapon }) {
	const v = clone(shipped);
	v.player.hunter.options.meleeResetsAutoShot = reset;
	if (rotation) v.player.rotation = JSON.parse(fs.readFileSync(path.join(ROT_DIR, rotation), 'utf8'));
	if (weapon) {
		const item = byId.get(weapon);
		v.player.equipment.items[14] = { id: weapon, enchant: item.handType === 4 ? ENCHANT_2H : ENCHANT_1H };
		if (item.handType === 4) v.player.equipment.items[15] = {};
	}
	return v;
}

(async () => {
	const t0 = Date.now();
	const runs = [
		{ label: 'shipped rotation, reset off', reset: false },
		{ label: 'shipped rotation, reset on', reset: true },
		...rotations.map((r) => ({ label: `${r}, reset on`, reset: true, rotation: r })),
	];
	const mh = (shipped.player.equipment.items[14] || {}).id;
	for (const weapon of WEAPONS) {
		if (weapon === mh) continue;
		const n = (byId.get(weapon) || {}).name || weapon;
		runs.push({ label: `${n}, shipped rotation, reset off`, reset: false, weapon });
		runs.push({ label: `${n}, ${ROTATION}, reset on`, reset: true, rotation: ROTATION, weapon });
	}
	const dps = await pool(runs.map((r) => () => sim(variant(r))));
	const report = {
		build,
		mode,
		iters: ITERS,
		mainHand: (byId.get(mh) || {}).name,
		results: runs.map((r, i) => ({ run: r.label, dps: +dps[i].toFixed(1) })),
		minutes: +((Date.now() - t0) / 60000).toFixed(1),
	};
	fs.mkdirSync(path.dirname(outPath), { recursive: true });
	fs.writeFileSync(outPath, JSON.stringify(report, null, 1));
	console.log(`${build} ${mode}: ${report.results.map((r) => `${r.run} ${r.dps}`).join('; ')}; ${report.minutes} min`);
	fs.rmSync(TMP, { recursive: true, force: true });
})();
