# Arc's Forever Sims engine

The simulation engine behind [Arc's Forever Sims](https://arcs-forever-sims.pages.dev), a hunter sim for WoW Forever.
It is a fork of Sanctum's Forever version of the wowsims classic engine (MIT, see LICENSE):
[compare with Sanctum's](https://github.com/Sanctumx/classic/compare/master...ArcSim:arc-forever-sims-engine:master).

- [docs/FOREVER_CHANGES.md](docs/FOREVER_CHANGES.md): every change made for Forever, where it lives in the code, and
  the assumptions behind it
- `sim/`: the sim itself, every class (the hunter is the one tuned for Forever so far)
- `ui/*/apls` and `ui/*/gear_sets`: the rotations and gear sets
- `builds/`: the settings behind every reference number on the site (2H Weave here uses ignite-style Lacerating
  Strikes, as the MythicSim tier list engine does; the site uses a plain refresh)
- `assets/database`: the item database, built by `tools/database`
- `proto/`: the settings and results formats

## Check a number

```
AB_SETTINGS=../../builds/ship_dw_weave.json AB_ITERS=50000 go test --tags=with_db ./sim/hunter/ -run TestForeverSettings -v
```

prints the DPS for that build (the path is relative to `sim/hunter`). A setup exported from the site (Export, JSON)
runs the same way.

## Compare with the MythicSim tier list

`builds/compare/` has the tier list's Survival reference (its gear, talents, rotation, buffs and consumables, plus the
sapper and food other specs get there) and our builds on that same setup (2H Weave on the harder and the light weave,
DW Pure Melee), all with ignite-style Lacerating Strikes:

```
AB_SETTINGS=../../builds/compare/tierlist_survival.json AB_ITERS=50000 go test --tags=with_db ./sim/hunter/ -run TestForeverSettings -v
AB_SETTINGS=../../builds/compare/arc_2h_weave.json AB_ITERS=50000 go test --tags=with_db ./sim/hunter/ -run TestForeverSettings -v
```

## Build and test

Needs Go 1.23 or newer and protoc with protoc-gen-go. Generate the Go code from the protos once:

```
protoc -I=./proto --go_out=./sim/core ./proto/*.proto
```

Then:

```
go test --tags=with_db ./sim/hunter/ -skip TestP1Hunter
```

`TestP1Hunter` compares against stored results from before the Forever changes, so it is skipped.

The WebAssembly build the site runs:

```
GOOS=js GOARCH=wasm go build -o lib.wasm ./sim/wasm/
```
