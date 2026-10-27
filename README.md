# Gear sweeps

Pre-raid gear checks for the builds of Arc's Forever Sims, run on GitHub's machines with the engine on `master`.

- `runner/job.json` lists the builds and play modes to check; pushing a change to it starts a run.
- `runner/gear_screen.js` tries every short-listed pre-raid item in every slot, re-checks the promising ones at 30,000
  fights, and combines the swaps that hold up.
- `runner/preraid_items.js` decides what counts as pre-raid gear.

Each run's results are attached to it for a day.
