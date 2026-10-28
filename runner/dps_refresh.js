// One shipped build and play mode as the site ships it, on the engine the job checked out (the job's "engine" ref):
// its DPS and spread over ITERS fights, for before / after comparisons of an engine change.
// usage: node dps_refresh.js <build> <God|Pro|Regular> <out.json>
// env: ENGINE_DIR (engine checkout), HUNTER_TEST (built test binary), ITERS, ENGINE_REF (for the output only),
// SETTINGS (a file in settings/ run in place of the engine's builds/ship_<build>.json), EP (the stat weights instead)
const fs = require('fs');
const os = require('os');
const path = require('path');
const { spawnSync } = require('child_process');

const [build, mode, outPath] = process.argv.slice(2);
const ENGINE = process.env.ENGINE_DIR || path.join(__dirname, '..', 'engine');
const EXE = process.env.HUNTER_TEST || path.join(__dirname, '..', 'hunter_test');
const ITERS = Number(process.env.ITERS) || 20000;
const PLAY = { God: 'WeaveGodMode', Pro: 'WeavePerfectPlay', Regular: 'WeaveNormalPlay' };

const source = process.env.SETTINGS ? path.join(__dirname, 'settings', process.env.SETTINGS) : path.join(ENGINE, 'builds', `ship_${build}.json`);
const settings = JSON.parse(fs.readFileSync(source, 'utf8'));
settings.player.hunter.options.weavePlay = PLAY[mode];
const file = path.join(fs.mkdtempSync(path.join(os.tmpdir(), 'dps-')), 'settings.json');
fs.writeFileSync(file, JSON.stringify(settings));
const r = spawnSync(EXE, ['-test.run', 'TestForeverSettings', '-test.v'], {
	cwd: path.join(ENGINE, 'sim', 'hunter'),
	env: { ...process.env, AB_SETTINGS: file, AB_ITERS: String(ITERS), AB_EP: process.env.EP ? '1' : '' },
	encoding: 'utf8',
	maxBuffer: 64 << 20,
});
const out = r.stdout || '';
if (process.env.EP) {
	// the stat weight lines as builds/ep/ship_<build>.txt holds them
	const ep = out.split('\n').filter((l) => l.startsWith('EP '));
	if (!ep.length) {
		console.error(out.slice(-2000), r.stderr);
		process.exit(1);
	}
	fs.mkdirSync(path.dirname(outPath), { recursive: true });
	fs.writeFileSync(outPath, JSON.stringify({ build, mode, engine: process.env.ENGINE_REF || 'master', iters: ITERS, ep }, null, 1));
	console.log(ep.join('\n'));
	process.exit(0);
}
const dps = Number((out.match(/raid DPS incl\. pets ([0-9.]+)/) || [])[1]);
const stdev = Number((out.match(/stdev ([0-9.]+)/) || [])[1]);
if (!dps) {
	console.error(out.slice(-2000), r.stderr);
	process.exit(1);
}
const result = { build, mode, engine: process.env.ENGINE_REF || 'master', iters: ITERS, dps, stdev };
fs.mkdirSync(path.dirname(outPath), { recursive: true });
fs.writeFileSync(outPath, JSON.stringify(result, null, 1));
console.log(JSON.stringify(result));
