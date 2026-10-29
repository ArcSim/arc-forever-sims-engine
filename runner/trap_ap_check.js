// One shipped build and play mode with two trap attack power rules, same seed: the engine's 0.132 x AP, and the line
// fitted to the beta's Immolation Trap logs, 0.1052 x AP + 11.1 a trap, given to the engine as 0.1052 + 11.1 / AP at
// the build's final attack power (the higher of melee and ranged).
// usage: node trap_ap_check.js <build> <God|Pro|Regular> <out.json>
// env: ENGINE_DIR (engine checkout), HUNTER_TEST (built test binary), ITERS (default 3000), ENGINE_REF (for the
// output only), SETTINGS (a file in settings/ run in place of the engine's builds/ship_<build>.json)
const fs = require('fs');
const os = require('os');
const path = require('path');
const { spawnSync } = require('child_process');

const [build, mode, outPath] = process.argv.slice(2);
const ENGINE = process.env.ENGINE_DIR || path.join(__dirname, '..', 'engine');
const EXE = process.env.HUNTER_TEST || path.join(__dirname, '..', 'hunter_test');
const ITERS = Number(process.env.ITERS) || 3000;
const PLAY = { God: 'WeaveGodMode', Pro: 'WeavePerfectPlay', Regular: 'WeaveNormalPlay' };
const SLOPE = 0.1052;
const PER_TRAP = 11.1;

const source = process.env.SETTINGS ? path.join(__dirname, 'settings', process.env.SETTINGS) : path.join(ENGINE, 'builds', `ship_${build}.json`);
const settings = JSON.parse(fs.readFileSync(source, 'utf8'));
settings.player.hunter.options.weavePlay = PLAY[mode];
const file = path.join(fs.mkdtempSync(path.join(os.tmpdir(), 'trap-')), 'settings.json');
fs.writeFileSync(file, JSON.stringify(settings));

function run(test, env) {
	const r = spawnSync(EXE, ['-test.run', `${test}$`, '-test.v'], {
		cwd: path.join(ENGINE, 'sim', 'hunter'),
		env: { ...process.env, AB_SETTINGS: file, ...env },
		encoding: 'utf8',
		maxBuffer: 64 << 20,
	});
	return { out: r.stdout || '', err: r.stderr || '' };
}
function fail(r) {
	console.error(r.out.slice(-2000), r.err);
	process.exit(1);
}

const stats = run('TestForeverStats', {});
const fin = stats.out.match(/final\s.*?AP\s+(\d+)\s+RAP\s+(\d+)/);
if (!fin) fail(stats);
const ap = Math.max(Number(fin[1]), Number(fin[2]));
const coefficient = SLOPE + PER_TRAP / ap;

function sim(env) {
	const r = run('TestForeverSettings', { AB_ITERS: String(ITERS), ...env });
	const dps = Number((r.out.match(/raid DPS incl\. pets ([0-9.]+)/) || [])[1]);
	if (!dps) fail(r);
	const traps = {};
	for (const m of r.out.matchAll(/^\s+((?:Immolation|Explosive) Trap)\s+([0-9.]+)%.*?casts\s+([0-9.]+).*dmg\/fight\s+(\d+)/gm)) {
		traps[m[1]] = { share: Number(m[2]), casts: Number(m[3]), damage: Number(m[4]) };
	}
	return { dps, stdev: Number((r.out.match(/stdev ([0-9.]+)/) || [])[1]), traps };
}
const engine = sim({});
const fitted = sim({ AB_TRAP_AP: coefficient.toFixed(5) });
const result = {
	build, mode, engine: process.env.ENGINE_REF || 'master', iters: ITERS, ap, coefficient: Number(coefficient.toFixed(5)),
	dps: { engine: engine.dps, fitted: fitted.dps, change: Number((fitted.dps - engine.dps).toFixed(1)) },
	traps: { engine: engine.traps, fitted: fitted.traps },
};
fs.mkdirSync(path.dirname(outPath), { recursive: true });
fs.writeFileSync(outPath, JSON.stringify(result, null, 1));
console.log(JSON.stringify(result));
