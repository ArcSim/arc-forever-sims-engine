# Gear sweeps

Pre-raid gear checks for the builds of Arc's Forever Sims, run on GitHub's machines with the engine on `master`.

- `runner/job.json` lists the builds and play modes to check; pushing a change to it starts a run. Each job has a `tag`
  that names its results.
- `runner/gear_screen.js` tries every short-listed pre-raid item in every slot, re-checks the promising ones at 30,000
  fights, and combines the swaps that hold up. A job's `exclude` items are never tried, and a worn one is replaced first.
  With `reset`, melee swings reset Auto Shot; `rotation` and `gear` start it from a file in `runner/rotations/` or
  `runner/gear/`.
- `runner/rotation_check.js` (jobs with `"script": "rotation_check.js"`) sims the build with melee swings resetting Auto
  Shot: the shipped rotation with the reset off and on, every rotation in `runner/rotations/`, and the job's `weapons`.
- `runner/preraid_items.js` decides what counts as pre-raid gear.

Each run's results are attached to it for a day.
