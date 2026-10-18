# Forever changes: code map

This repository is a fork of a WoWSims "classic" fork (MIT licence). Upstream ends at commit
`f9980ef`; everything after it is this fork's work (about 200 commits). The fork adapts the
hunter simulator to WoW Forever, Blizzard's re-release with changed spells, items, racials,
enchants and consumables. Values come from the Forever client data (build 1.60.1.x) where it
has them, from in-game logs where it does not, and are listed under "Assumptions" where
neither settles it. The hunter is the only sim the site lists. Racials and caster-debuff rules
are shared by every class and are tested on mage, warrior and warlock too.

To see the full change set: `git diff f9980ef..HEAD -- . ':!assets' ':!*.pb.go' ':!ui/core/proto'`.

## Conventions used in the code

- Engine switches named `Forever*` in `sim/core/attack.go` default to true. Setting one to
  false restores the older rule, for A/B comparisons only.
- Environment variables starting with `AB_` are local comparison hooks read with `os.Getenv`
  (for example `AB_MIN_RANGED`, `AB_NO_STAND_STILL`, `AB_NORM_MAX`, `AB_NO_SK_SPEED`,
  `AB_TRAP_AP`, `AB_LAC_MODE`, `AB_OH_SPEED`, `AB_MH_SPEED`, `AB_HASTE_PCT`). All are unset
  in normal use. The test harnesses also use `AB_SETTINGS`,
  `AB_ITERS`, `AB_TALENTS`, `AB_APL`, `AB_RACE`, `AB_WORLDBUFFS`.
- Log-only markers: some rows on the timeline (off-hand wait, extra move time) are written to
  the sim log as a "cast" with no effects, so the UI can draw them. They change no results.

## Where things live

### Engine core (`sim/core`)

Attack and movement

| File | Forever-specific behaviour |
| --- | --- |
| `attack.go` | A melee swing (Raptor Strike included) restarts the Auto Shot timer. When the melee attack turns on, the main hand swings when ready and the off hand waits half its swing time (`OffHandWait`, `OffHandWaitsForAttackStart`); this only applies if an Auto Shot turned the melee attack off. A walk into melee starts swinging at 5 yd, not at arrival (`StartMeleeMidMove`). `RangedCycleCut` shortens each Auto Shot timer (the Auto Shot spam bug option). Normalized instants (Strider Kick, Mongoose Bite) use the higher of real and normalized weapon speed. Every swing timer start and haste move is logged (`startTimer`, `moveReadyAt`, under the weapon's own id from `timerID`) for the timeline's swing timer bars; results do not change. The rotation is looked at again right after each Auto Shot, like after a melee swing, so a move into melee can start during a GCD (TestMoveAfterShotDuringGcd); the builds not on Arc Weave come out the same with or without it. The pull (`startPull`) leaves the attacks to a walk already under way: the first Auto Shot fires before the pull and can start a walk in (TestMoveAtPull). |
| `movement.go` | Auto Shot needs 0.5 s of standing still (`AutoShotStandStill`). Base run speed comes from `Unit.BaseRunSpeed` (0 means 7 yd/s). Speed bonuses normally compete (strongest counts); `MultiplyStackingMoveSpeed` is for ones that stack. `Unit.MoveDelay` is a wait before each step into or out of melee; `PendingMove` stops a rotation from issuing a second move during it. |
| `constants.go` | Ranged minimum distance is 8 yd (`MinRangedAttackDistance`, was 12). |
| `unit.go` | New fields: `MoveDelay`, `BaseRunSpeed`, `SkyborneConvergence`. |
| `character.go` | `MainHandCritBonus`: crit from a main-hand stone applies to melee attacks only, not ranged. |
| `apl_actions_misc.go` | The Move action waits while a delayed move is pending. |
| `apl_values_operators.go` | A duration times or divided by a plain number scales correctly (fractions were truncated). |

Racials, consumables, buffs, debuffs, attack table

| File | Forever-specific behaviour |
| --- | --- |
| `racials.go` | Racials follow the client data: weapon specialisations give crit (no Bow/Gun/Thrown spec), Ruthless and Savage crit, Earth Fury and Endurance crit damage, Blood Fury (flat +10% attack power and spell power), Berserking (spell 20554), Elune's Light, Big Game Hunter, Beast Slaying, Quickness (stacks with other speed bonuses), Gnome Eureka!, Undead Touch of the Grave, and the Forever-only Skyborne race (Wind Blessed, Elemental Insight, Skysight, Read Leyline). |
| `base_stats.go` | Base stat offsets for Skyborne (Alliance and Horde). |
| `consumes.go` | Natural flasks: +60 Stamina plus a second effect that only works in the Hyjal zones (see `OutsideHyjal`). Shadow Oil's bolt uses Forever's spell and range. New foods and elixirs. Explosives: Crystal Charge (instant) and the EZ-Thro / SAF-T line need no profession; the sapper, Dark Iron Bomb and the old dynamite and grenades need Engineering; all share one minute. The Holy Hand Grenade of Alterac does Holy damage. `imbueFitsWeapon`: sharpening stones only on axes, daggers, polearms and swords; weightstones only on fist weapons, maces and staves; a misfit stone adds nothing. |
| `buffs.go` | Talented raid buffs use Forever values (Improved Mark of the Wild, Battle Shout, Blessing of Might, Blessing of Wisdom, Devotion Aura, Enhancing, Restorative and Guardian Totems). Windfury Totem (attack power buff lasts 1 s, 100 ms between procs, ranks 95/179/246) and Flametongue Totem (a party aura adding fire damage to main-hand white hits) are player-only, and Windfury wins if both are asked for. Pets take no raid buffs. Grace of Air stacks with Windfury. |
| `debuffs.go` | Values updated to Forever data (Expose Armor, Demoralizing Shout, Thunder Clap, Curse of Recklessness and others). Caster debuffs such as Improved Shadow Bolt, Shadow Weaving, Improved Scorch, Winter's Chill and Stormstrike only raise damage taken from the caster who applied them. There is no Improved Hunter's Mark: the setting means the plain Mark. |
| `target.go` | Attack table: no 1% hit suppression against +3 targets. `OutsideHyjal` flag on the encounter. |
| `spell_outcome.go`, `stats/stats.go` | `PseudoStats.SpellCritMultiplier`: Earth Fury and Endurance raise spell crit damage too. |
| `focus.go`, `mana.go` | `MultiplyFocusRegen` scales the whole focus regeneration rate (Bestial Discipline uses it on top of the doubled pet rate). Mana regeneration while casting never exceeds the not-casting rate. `CastingRegenSource`: a regen-while-casting buff (Resourcefulness, Rapid Recuperation's sting buff) refreshes the cached mana ticks when it goes up or down, and its share of each casting tick is its own mana row. |
| `metrics_aggregator.go` | Auto Shot is counted as a cast in the results, not as a passive action. |
| `sim_concurrent.go` | When worker results are merged, auras are matched by id, not by position. A worker that registers an aura mid-fight no longer breaks the merge; missing auras count as zero uptime. |

Encounter options and protos (`proto/`)

| File | Forever-specific content |
| --- | --- |
| `common.proto` | Skyborne races, Natural flasks, new elixirs and foods, `flametongue_totem`, `Encounter.outside_hyjal` (the "Hyjal zone" option; unset counts as inside). |
| `hunter.proto` | Pet families (Fox, Core Hound, Bird of Prey), per-ability pet switches, new ammo, Lacerating Strikes mode, pet totem switches, move delay, run speed, Auto Shot spam bug, Explosive Trap burn scaling, Skyborne convergence. |
| `apl.proto` | `APLSimpleView`, `APLSimpleBox`, `APLSimpleValue`: a plain-words view of a rotation line, written by the rotation's author. `APLRotation.recheck_after_auto_shot`: no longer read (the recheck is always on), kept so rotations saved with it still load. `APLValueWalkTime`. |
| `ui.proto` | `forever_changed_only` item filter. |

### Item data and effects (`sim/common`, `tools/database`)

| File | Content |
| --- | --- |
| `sim/common/enchant_effects.go` | Two-hand Impact at +6 weapon damage (Forever effect id). Grand Crusader and Grand Inquisitor strength procs. |
| `sim/common/item_effects.go` | Weapon proc damage ranges for 18 weapons, Hand of Justice at 1%, Nat Pagle's Broken Reel cooldown, a Forever version of the Blackhand's Breadth trinket, a weapon fire proc with its Forever spell id. Also declares the shared armor debuff group name `MinorArmorReduction`. |
| `sim/common/item_sets/dungeon_set_1.go` | Dungeon Set 1 bonuses at 2 to 6 pieces. Left out: the 4-piece when-struck effects and the caster 5-piece procs. |
| `sim/common/item_sets/item_sets_pve.go` | Cadaverous Garb bonuses changed to the Forever values; Gladiator set bonuses at 3, 4 and 5 pieces. |
| `tools/database/database.go` | `MergeExistingItems` (a feed never adds bare items), weapon skills replaced on merge, icon fallbacks for Forever-only ids. |
| `tools/database/enchant_overrides.go` | Forever's enchant list: changed values, new effect ids, armor kits, new agility and strength enchants. |
| `tools/database/forever_overrides.go` | Generated: stat and weapon changes for 1114 items from the Forever item feed. |
| `tools/database/forever_new_items.go` | Generated: four Forever-only items a hunter can wear. |
| `tools/database/forever_icons.go` | Generated: icons for Forever-only items and spells. |
| `tools/database/overrides.go` | Adds the new foods, elixirs, potion and flasks to the list of consumable item ids the generator keeps. |
| `tools/database/gen_db/main.go` | Applies the overrides. `FOREVER_OVERRIDES=off` builds the unmodified baseline. |

The three generated files name generator scripts (`scripts/gen_forever_*.js`) that are not in
this repository. Treat them as data; edit by hand only with care.

### Hunter class (`sim/hunter`)

| File | Forever-specific behaviour |
| --- | --- |
| `hunter.go` | Wires the new options (move delay, run speed, Skyborne convergence, Auto Shot spam cut). Swiftfeather Arrow and Swiftstrike Shot ammo. Mortal Shots no longer applies to Auto Shot. |
| `aimed_shot.go`, `arcane_shot.go`, `sniper_shot.go`, `multi_shot.go`, `volley.go` | Aimed Shot 2.0 s cast, rank damage tables, Sniper Shot rank 2, Barrage bonus (3/7/10% on Multi-Shot, Aimed Shot, Volley). |
| `serpent_sting.go` | Forever tick damage per rank from the client data (rank 9: 111). The client gives no attack power or spell power share, but in game each tick adds 0.032 of ranged attack power (rank 4: 27 at 174, 33 at 361) and spell power changes nothing (TestSerpentStingScaling). |
| `rapid_fire.go` | Also +40% melee attack speed. Rapid Killing takes time off the cooldown per rank. |
| `raptor_strike.go` | Forever rank damage. Queues at any range and waits for the next melee swing. Its cooldown is readable by rotations. No Mortal Shots crit bonus. |
| `strider_kick.go` | Also +30% run speed for 3 s. Yellow-only proc mask. |
| `aspects.go` | Quick Strikes buff, and Aspect of the Cheetah (+30% run speed, shares the aspect slot, Pathfinding). |
| `immolation_trap.go`, `explosive_trap.go` | Forever spell ids, 15 s burn in 3 s ticks, magic spell (ticks crit for 1.5x, no melee procs), shared 30 s cooldown, attack power share from T.N.T. (`trapAPBonus`, higher of melee and ranged attack power). Explosive Trap's burn only scales with attack power if the option is on.. `trapReach`: how far from the target a trap can be dropped and still catch it, 6 yd (it lands about a yard in front of you; in game it hits from 6 yd, not from 8); Immolation Trap cannot be dropped beyond it and both do nothing beyond it (TestTrapReach; AB_TRAP_REACH for comparisons). |
| `summon_hawk.go` | Hawk attacks follow the hunter's ranged attack speed, resolve as white melee (glancing, dodges, misses), crit about 1% for double damage; the dive scales with ranged attack power and Unleashed Fury. |
| `pet.go` | Pet attack speed is one setting per pet (default: the family's fastest tamable speed). Focus regeneration doubled. Family table with each family's ability list and health, armor and damage scalars. |
| `pet_abilities.go` | Top-rank pet abilities per the client data: Bite, Claw, Lightning Breath, Scorpid Poison, Savage Rend, Furious Howl (buffs the hunter), Thunderstomp, Dismember, Pinch, Tendon Rip, Web, Swipe, Mine!, Charge (opener), Demoralizing Screech, Lava Breath, Dust Cloud, Trickster's Dance, Sonic Blast. The pet rotation handles any number of abilities and honours the per-ability switches. |
| `talents.go` | Forever per-rank values (Lightning Reflexes, Unleashed Fury, Efficiency, Barrage, Survival Tactics, Ferocity, Resourcefulness, Bestial Discipline), Intimidation, Frenzy buff, Lacerating Strikes with three bleed modes (refresh, strongest stays, rolling). |
| `items.go` | Item effects that buff the pet's crit now apply to every pet ability, not two fixed ones. |

### Tests

| File | Covers |
| --- | --- |
| `sim/hunter/forever_patch_test.go` | The bulk of the hunter checks (33 tests): totems, racials, Auto Shot stand-still, spam bug and melee restart, off-hand wait, mid-walk swings, Raptor Strike queueing, traps (attack power, crit, cooldown, cast metrics), shared armor debuff group, Aspect of the Cheetah, run speed setting, imbue fit, foods, elixirs, weapon enchants, Natural flasks, Skyborne racials, set bonuses. |
| `sim/hunter/forever_buffs_test.go` | Talented raid buffs and Forever debuff values. |
| `sim/hunter/forever_hawk_test.go` | Hawk swing interval, damage and outcomes. |
| `sim/hunter/bm_talents_test.go` | Bestial Discipline, casting regeneration cap, Intimidation. |
| `sim/hunter/forever_ab_test.go` | Local A/B harness (dual wield vs two-hand under matched settings); prints a table, almost no assertions. |
| `sim/hunter/forever_settings_test.go` | Runs a settings export from the site (`AB_SETTINGS`) and prints the per-action table; skipped without it. |
| `sim/hunter/forever_stats_test.go` | Prints final character stats of a settings export to check buff claims; skipped without `AB_SETTINGS`. |
| `sim/core/forever_racials_test.go` | Touch of the Grave proc chance by class. |
| `sim/core/sim_concurrent_test.go` | Merging worker results matches auras by id. |
| `sim/forever_racials_test.go` | Crit-damage racials, Blood Fury spell power, Eureka!, Touch of the Grave drain, caster-only debuffs, each on other classes. Package `sim` is not in the test command below; run it with `go test --tags=with_db ./sim/`. |

## How to build and test

Needs Go 1.23.x and Node 20 or 22 (see `README.md`).

Protos:

```
protoc -I=./proto --go_out=./sim/core ./proto/*.proto
npx protoc --ts_opt generate_dependencies --ts_out ui/core/proto --proto_path proto proto/api.proto
```

Database:

```
go run ./tools/database/gen_db -outDir=assets -gen=db
```

Tests:

```
go test --tags=with_db ./sim/core/... ./sim/hunter/ ./sim/common/... -skip TestP1Hunter
```

`TestP1Hunter` is upstream's golden-results suite. Its stored results file
(`sim/hunter/TestP1Hunter.results`) is upstream's and was not regenerated for this fork, so it
is skipped. The `lib.wasm` files under `dist/classic`, `public` and `ui/public` are committed
build outputs.

## Assumptions

Values the code or the site marks as assumed, estimated or not yet confirmed in game:

- Hawk swing damage at level 60 (36 per attack) is scaled from level 30 logs by the creature
  damage table, and assumes the hawk is the hunter's level (`sim/hunter/summon_hawk.go`).
- Lacerating Strikes bleed ticks can crit; not confirmed in game (listed in the site's known issues).
- Glancing blow damage of 35% is not confirmed in game (listed in the site's known issues).
- Explosive Trap's burn does not scale with attack power by default; whether it does is not
  confirmed (option `explosive_trap_burn_ap`, off by default).
- Dust Cloud (tallstrider) is assumed to stack with Sunder Armor and similar effects. It shares
  an armor group with Faerie Fire and Curse of Recklessness (`ui/hunter/pet_families.ts`).
- Grand Crusader and Grand Inquisitor: the client has no proc rate, so they use 1 PPM like Crusader.
- Gladiator set 4-piece (Roar of the Crowd): the client has no proc rate, so 1 PPM is used.
- Dungeon Set 1 three-piece attack power is applied as the client data says, though it is
  reported not to apply in game.
- Auto Shot spam bug cut of 0.45 s is the best logged gap; a melee restart is not shortened
  because no log shows it (`AutoShotSpamCut`).
- Immolation Trap's attack power share (0.132) comes from live logs, not the client's 25%.
- Only the strongest run speed bonus counts, except Night Elf Quickness, which is the only one
  known to stack (`AB_SKYSIGHT_STACKS=1` lets Skysight stack, for comparison).
- Core Hound and Fox use neutral family scalars (damage 1.00); nothing else is listed for them
  in `sim/hunter/pet.go`.
- Shipped builds assume no extra move time and the game's run speed of 7 yd/s; the step
  distances in the rotations carry the footwork (`ui/hunter/presets.ts`).
- Left out on purpose: Aspect of the Cheetah's daze, the self-heal of Touch of the Grave, and
  caster-only effects of dungeon sets.
