# MythicSim downstream patches

MythicSim runs this engine from its fork (`sage3648/mythicsim-forever-engine-go`, branch
`mythicsim/upstream-sync-20261010`). The branch is ElliotWood/Forever master, which is built on the
official wowsims/forever, plus the patches below. The first base was `442076902` (Merge
wowsims/forever master ea5412873). The current base is `071e0cb439` (2026-10-10, client 1.60.1.70334); the 2026-10-10 sync merged #720 to #760, dropped patch 11 and gave patch 5 upstream's move into melee ("Upstream sync 2026-10-10, #720 to #760" at the end of this file). The base before it was `5c115f1725` (2026-10-07, client 1.60.1.70235 with the 2026-10-06 hotfixes); the 2026-10-07 sync merged #677 to #719 and dropped patches 3, 76, 91, 93 and 94 and folded upstream's Flametongue Totem into patches 70 to 72 ("Upstream sync 2026-10-07, #677 to #719" at the end of this file). The base before it was `67f14b04a5` (2026-10-04, client 1.60.1.70205); the 2026-10-05 sync merged #642 to #676 and dropped patches 13, 18, 19, 20, 22, 23, 27, 30 and 31, which upstream now carries ("Upstream sync 2026-10-05, #642 to #676" at the end of this file). The base before it was `f764984d8b` (2026-10-03), merged by "Upstream sync 2026-10-03, #613 to #641", and before that `f4b776b4f4`, and before that `ccfaacb5c3` (2026-10-02, client 1.60.1.70170 with the 2026-10-02 hotfix cache). The 2026-10-01 syncs merged 113 upstream commits and then 28 more ("Upstream sync 2026-10-01, second merge" below); the decisions are in "Upstream sync 2026-10-01". The 2026-10-02 sync merged the 12 commits #602 to #609 ("Upstream sync 2026-10-02, #602 to #609" at the end of this file): patches 50 and 61 are dropped, patch 63 is narrowed to Mystic Mushroom, and the interim 70170 regeneration is replaced by upstream's real one. The previous base was `8dc19a4241` (2026-09-27). It includes form-speed and actual spell cast-time Omen of Clarity proc corrections, life-drain weapon effects, Sword of Zeal, Argent Avenger, Fiery Weapon and Lifestealing enchants, Flurry Axe and Electrified Dagger, the 2026-09-27 client hotfix database, Stinging Viper and eight Classic weapon procs, Mage Scroll of Cryoblast, non-engineer explosives and SAF-T / EZ-Thro bombs, Deep Wounds weapon-only damage with outstanding bleed rollover, Raptor pet Savage Rend, Venomstrike procs, Defias Leather set effects, Barbaric Crossbow, Plaguefang and Wolfsbane weapon procs, the Stormshroud and Volcanic Armor proc chances, item effects below item level 50, the refreshed client database, Druid form Faerie Fire cost and timing, Hunter pet Lightning Breath scaling, Inspiration armor bonuses, the client hotfix databases, Hunter ranged scaling, Rogue Hack and Slash cooldown, Shaman Flametongue and Fire Nova fixes, and the merged Penance timing and cost fixes, Demonic Pact pre-pull sacrifice, Mana Tide Totem party restoration, Frost Mage talent fixes, and rank 4 Trueshot Aura. It also carries client 1.60.1.70009 and the earlier lower-rank spell, aura-cap, and consumable fixes.

Keep the set small. Each patch exists because MythicSim needs something upstream does not do
yet. Drop a patch as soon as upstream covers it; do not keep ours alongside an upstream version.

| # | Commit subject | Why MythicSim needs it |
|---|---|---|
| 1 | `cli: sim --strict rejects unknown fields and enum names` | The worker builds requests in code. Without it, a misspelt field or a race the build does not know is dropped silently and the sim runs a different character. |
| 2 | `core: a player option to disable racials` | The race comparison page sims each character with and without its racials to show what they are worth. |
| 3 | `rotation: Destruction casts Conflagrate for Shadow and Flame` | Dropped 2026-10-07: upstream #716 makes `destruction_conflag` the Shadow and Flame preset. |
| 4 | `hunter: Aspect of the Beast` | Forever made Beast the melee aspect. Upstream models only Hawk, so a melee hunter has no aspect. |
| 5 | `rotation: a melee Survival rotation` | Upstream's melee Survival rotation (#734) sims 19% below ours in melee. MythicSim ranks melee Survival. |
| 6 | `items: Iceblade Hacker and Warblade of Caer Darrow proc from their own hand` | The two hand-written weapon procs fired off both hands, so a main-hand Iceblade Hacker added its Frost damage to every off-hand swing. |
| 7 | `data: inherited stat indices and armor` | Preserve the corrected item stats and armor. |
| 8 | `mage: implement baseline Frostfire Bolt` | Implement the baseline spell and its hybrid rotation. |

## 1. `cli: sim --strict`

- **What it does.** `wowsimcli sim --strict` loads the request with protojson `DiscardUnknown`
  off and exits non-zero with the protojson error. The flag defaults to off, which is
  upstream's behaviour.
- **Files.** `cmd/wowsimcli/cmd/basic_sim.go` (`loadRaidSimRequest`) and `basic_sim_test.go`.
- **Drop it when** upstream's CLI can reject unknown names itself. If upstream's flag has a
  different name, switch the worker to it and drop this patch.
- **Conflicts** can only happen where `simMain` loads the file.

## 2. `core: a player option to disable racials`

- **What it does.** `Player.disable_racials` (field 59, JSON `disableRacials`) skips every racial
  effect and keeps the race's base stats. Besides `applyRaceEffects`, it gates the race-only
  effects that live elsewhere: the night elf priest's Starshards and the undead priest's Dark
  Sacrifice (`sim/priest/priest.go`, patch 85),
  Bloodthistle (`sim/core/consumes.go`) and the racial multipliers the reforge optimizer models
  (`sim/core/reforge_optimizer/model.go`). The rest of the code checks
  `Character.RacialsDisabled()`.
- **Tests.** `sim/core/racials_test.go` (`TestDisableRacials*`) and
  `sim/priest/racials_test.go`.
- **When rebasing**, grep upstream's new code for race checks outside `applyRaceEffects`
  (`\.Race ==`, `Race_Race`, `GetRace()`) and gate any new ones on `RacialsDisabled()`.
- **Field number.** The worker sends protojson, which reads the field by name, so the number
  only matters for binary protos and saved UI links. If upstream takes 59 for something else,
  move ours to the next free number and leave the name alone.
- **Drop it when** upstream has an equivalent option. Point the worker's `disableRacials`
  (`worker/cmd/refresh-forever-races`) at upstream's name first.

## Dropped: `core: Forever races from client 1.60.1.69977`

Dropped in the rebase onto `6cb2603d`. Upstream now models the Forever races itself: the Skyborne
(`RaceHighOrderSkyborne`, `RaceWindshaperSkyborne`), the Forever race/class pairings, and the
client's racials. Blood Elf, Draenei and Bloodthistle are gone. The patch's client tests were run
against upstream's version:
- Blood Fury, Berserking, Elune's Light, Touch of the Grave, Tauren Endurance, the Human Spirit and
  the weapon specializations match.
- The Skyborne match. Upstream applies Wind Blessed as one attack speed multiplier, and it takes
  the Skyborne base stat offsets (±1) from level 1 character sheets where the patch had zero.
- Eureka! is the 70009 client's 10% cost cut on every class. The patch had 69977's per-class
  figures.
- Upstream also models Gnome Expansive Mind, which the patch left out.
- The pairings are the ones MythicSim's race pages offer.

## 3. `rotation: Destruction casts Conflagrate for Shadow and Flame` (dropped, upstream #716)

Dropped in the 2026-10-07 sync: upstream's Shadow and Flame preset and `TestDestruction` now run
`destruction_conflag.apl.json`, which casts Conflagrate whenever Immolate is up, and
`destruction.apl.json` is back to upstream's. On the merged engine `destruction_conflag` scores
194.66 against this patch's rule at 192.34 (TestDestruction Average-Default).


- **What it does.** `ui/specs/warlock/dps/apls/destruction.apl.json` casts Conflagrate (18932)
  after the Immolate refresh, while Immolate is up and either Immolate has under 4 seconds
  left or the warlock knows Shadow and Flame (Shadow, 1293816) and its buff is down. Without
  the talent it only Conflagrates at the end of Immolate. The rule is the one MythicSim's
  previous engine line measured (+5.9% on the 5/5 Destruction reference, neutral for 2/5
  builds). MythicSim copies this rotation into its worker presets.
- **Goldens.** `sim/warlock/TestDestruction.results` (average 410.50 to 420.67 DPS).
- **Drop it when** upstream's Destruction rotation casts Conflagrate.

## 4. `hunter: Aspect of the Beast`

- **What it does.** Registers Aspect of the Beast at its top rank (1299447): 110 melee attack
  power, exclusive with Aspect of the Hawk. With Deadly Aspects, a landed melee auto attack has
  Deadly Aspects' second effect as its chance (2% a rank) to trigger Quick Strikes (1299448), 30%
  melee haste for 12 sec. Beast states no proc chance of its own, so there is no Quick Strikes
  without the talent. Raptor Strike takes a swing's place as a special attack and does not proc it
  (Beast's proc flags are melee auto attacks, 0x4).
- **Files.** `sim/hunter/aspects.go`, the `AspectOfTheBeast` fields and spell mask in
  `sim/hunter/hunter.go`, and `sim/hunter/aspects_test.go`.
- **Tests.** `TestAspectOfTheBeastMeleeAttackPower`, `TestQuickStrikesNeedsDeadlyAspects`. No
  golden moves: the ranged rotations cast Hawk.
- **Drop it when** upstream models Beast. Check that its melee attack power and Quick Strikes
  match before dropping.

## 5. `rotation: a melee Survival rotation`

- **What it does.** `ui/specs/hunter/dps/apls/sv_melee.apl.json` is a melee rotation: Aspect of
  the Beast before the pull, then Raptor Strike queued on cooldown, Mongoose Bite whenever Expose
  Prey opens it, Summon Hawk, Strider Kick, Immolation Trap, and Wing Clip in any global where
  Raptor Strike and Strider Kick are more than half a second away. Wing Clip is a landed melee
  special, so it can proc Expose Prey (1310532's proc flags include melee specials, 0x10), and
  more Wing Clips mean more Mongoose Bites. Spells a build lacks are skipped, so the file serves
  any talents. `presets.ts` publishes it as `SurvivalMeleeRotation`, which MythicSim's
  `scripts/build-forever-gear.mjs` requires before it copies the file into the worker presets as
  `hunter_survival`.
- **Measured** on MythicSim's Beast Mastery reference character moved to 5 yards, with wowtbc.gg's
  5/10/35 build, at 10,000 iterations. Raptor Strike, Mongoose Bite and Strider Kick alone sim 380
  DPS. Adding Wing Clip takes it to 423, and Wing Clip plus Immolation Trap to 466. Explosive Trap
  sims 1 to 2% behind Immolation Trap. Upstream's `sv.apl.json`
  is untouched; it is the ranged Survival rotation.
- **Tests.** `sim/hunter/survival_melee_test.go` (`TestSurvivalMelee`, 5 yards,
  `SurvivalMeleeTalents` = wowtbc.gg's 5/10/35 with the spare point in Focused Fire) and its golden
  `TestSurvivalMelee.results` (average 398.57 DPS on the suite's weapons-only gear).
- **Upstream's rotation (2026-10-10 sync).** Upstream #734 ships its own `sv_melee.apl.json` (walks into
  melee, Mongoose Bite, Strider Kick, Serpent Sting above 30% mana, Immolation Trap, Raptor Strike last,
  with Explosive Trap and Volley on 2+ targets from #742 and #743) and `TestSurvival` runs it from 30
  yards on 0/15/36. Ours keeps the file and takes upstream's first line, a move into melee when the
  target is more than 5 yards away, so a sim that starts at range closes in. Average DPS on the merged
  engine: `TestSurvivalMelee` (5 yards, 5/10/35) ours 396.59, upstream's 322.70, ours with the move
  396.59; `TestSurvival` (30 yards, 0/15/36) ours 184.58 (it never closes), upstream's 290.27, ours
  with the move 352.19. Upstream's AoE lines stay out: MythicSim's worker adds Explosive Trap and
  Volley at 2+ targets itself.
- **Drop it when** upstream's melee Survival rotation sims at least as high as ours on both goldens.

## 6. `items: Iceblade Hacker and Warblade of Caer Darrow proc from their own hand`

- **What it does.** The two weapon procs in `sim/common/classic/items_store_gaps.go` ("Melee
  attacks with this weapon deal 41 / 28 Frost damage") get a
  `NewDynamicLegacyProcForWeapon(item, 0, 1)` proc manager, as every generated weapon proc has.
  Their proc masks alone named both hands' autos and specials. A main-hand Iceblade Hacker also
  procced off every off-hand swing, which was worth 12.5% of a dual-wielding melee hunter's damage.
- **Tests.** `sim/rogue/weapon_proc_hand_test.go` (`TestIcebladeHackerProcsOnlyFromItsHand`). The
  AllItems rows for the two weapons move in `TestFury`, `TestArms`, `TestProtectionWarrior` and
  `TestRetribution`: the harness equips Warblade in the off hand beside a main-hand weapon, where
  it used to proc off main-hand hits too (112 procs a fight on Arms, 31 now). No other row moves.
- **Drop it when** upstream's generator carries these procs (the file says to remove an entry
  then), or upstream scopes them to their hand. A generated `CreateWeaponCoHProcDamage` already
  does.

## Rebasing onto a newer upstream

1. Fetch ElliotWood/Forever master. Rebase the patches onto it:
   `git rebase --onto <new upstream> <old base> <fork branch>`.
2. Resolve `*.results` conflicts by taking upstream's side. During a rebase `--ours` is the
   upstream side and `--theirs` is the patch being replayed, so use
   `git checkout --ours -- <file>`. Never hand-merge a golden.
3. Regenerate the protos (the `-I=/usr/include` is libprotobuf-dev's `descriptor.proto`),
   build, and run the patches' own tests:

   ```
   protoc -I=./proto -I=/usr/include --go_opt=Mgoogle/protobuf/descriptor.proto=google.golang.org/protobuf/types/descriptorpb --go_out=./sim/core ./proto/*.proto
   go build --tags=with_db ./sim/... ./cmd/...
   go test ./cmd/wowsimcli/...
   go test --tags=with_db ./sim/core -run 'DisableRacials|Racial|Skyborne'
   go test --tags=with_db ./sim/priest -run 'Starshards|Arena'
   go test --tags=with_db ./sim/hunter -run 'Aspect|QuickStrikes|SurvivalMelee'
   go test --tags=with_db ./sim/rogue -run IcebladeHacker
   ```

4. After an inherited-stat database correction, verify offensive stat indices below 30 are unchanged, run the importer and armor regression tests, and regenerate affected result snapshots from a completed full test run. The correction of armor, health, mana, resistance and false Physical Damage changes item and mitigation cases across specs.

5. Run the full suite, `go test --tags=with_db $(go list ./sim/... | grep -v sim/web)`. Only
   `sim/warlock/TestDestruction` should fail, from patch 3's Conflagrate, and
   `sim/hunter/TestSurvivalMelee` wherever upstream moved hunter numbers (re-bless it into patch 5).
   If upstream's AllItems rows for Iceblade Hacker or Warblade of Caer Darrow move, re-bless them
   into patch 6. Re-bless it (copy
   `TestDestruction.results.tmp` over `TestDestruction.results`) and fold it into patch 3, so
   the patch carries the golden it moves. Any other failure is upstream's or the rebase's, not a
   golden to re-bless.
6. Check each patch's "drop it when" condition above, and update this file when a patch goes.

## 7. Legacy stat indices and armor

The inherited snapshot was imported before Resilience was removed from the Stat
proto. Its armor and later stats were one index too high. The importer now resolves
stat names against the current enum. The guarded migration corrects inherited
items, enchants and suffixes, preserves distinct client stat maps, restores missing
armor, and takes base/bonus armor from the pinned Forever planner. No conditional
attack power or spell effects are copied from the planner. Full generation uses
the same armor source. JSON and embedded binaries are updated together.

Validation: Python importer regressions and `go test -tags with_db ./sim/core
-run TestForeverReferenceArmor` check all 12 affected reference armor pieces and
Stoneskin Gargoyle Cape's 43 base armor plus 50 bonus armor. The one-time migration
refuses to run against an already migrated or different source snapshot.

## 8. Frostfire Bolt

Registers the three baseline Mage ranks from the client store and adds an experimental Frostfire
hybrid rotation. All ranks have a direct hit and a periodic effect. The spell counts as both Fire
and Frost for talent mods and direct-hit procs. Hot Streak, Fingers of Frost and Missile Barrage
explicitly include it. Improved Fireball reduces its cast time; Improved Frostbolt does not.

Core Frostfire handling takes the lower Fire/Frost resistance and the larger school spell-power,
hit and school multiplier values, avoiding duplicate Curse of the Elements or Elemental Precision.
This is limited to the Frostfire school so other mixed-school implementations are unchanged.

The pinned client store provides the damage, coefficients, mana, speed and periodic-crit flag.
Binary resistance follows Frostbolt and the original WoWSims SoD Frostfire Bolt, and still needs
Forever combat-log confirmation. The existing Fingers of Frost in-flight rule and Missile Barrage
proc chance remain the engine's beta assumptions. See `docs/frostfire-bolt.md` for evidence and tests.

Drop this patch once upstream implements all ranks, the relevant talent hooks and dual-school
calculations. Compare actual casts and damage before switching.

## 9. Light's Vigil refund reporting

The resource serializer merges rows sharing an ActionID. Light's Vigil used the
same ID for its cost and 75% refund, producing one net-negative row and hiding
its refund from resource-source breakdowns. Give only the refund metric tag 1.
The spell's cost, damage, refund calculation and timing remain unchanged.

Evidence: the shared Shockadin report `4161cca1-f878-498b-8aa4-8386c5f85331`.
`TestLightsVigilReportsCostAndRefundSeparately` exercises all three ranks. A
10,000-iteration replay on seed 42 produces identical DPS, threat, damage taken,
healing, time out of mana, actions and auras; the two resource rows sum to the
original net row. Costs and refunds can now be shown independently, including
mana lost to the cap.

Drop this patch when upstream gives the refund a distinct resource action ID
or exposes separate positive and negative resource totals in the result schema.

## 10. Innervate resource attribution

Self-cast Innervate recorded a synthetic 0.2 mana per use. External Innervate
recorded an estimate on top of the real mana ticks, double-counting regeneration.
Both now use one driver that attributes each tick's incremental regeneration to
Innervate and subtracts it from ordinary regeneration. Baseline regeneration gets
mana-cap space first; the bonus records only the remaining actual gain.

The reporting baseline follows other spirit effects, including additive Evocation,
form changes, current Spirit/MP5 and regen speed. Live arithmetic, tick timing and
cast decisions remain unchanged. This does not establish or change the game's
Evocation/Innervate stacking rules. Source-attributed passive regeneration does not
create mana-gain threat; the old synthetic entries incorrectly did. Only TPS values
changed in the Balance and Feral Cat golden fixtures.

Validation: core tests cover casting, noncasting, caps, dynamic stats, regen speed,
overlap, expiry and reset. Class tests exercise self/external Innervate and both
Evocation activation orders. All core, Druid and Mage suites pass with the embedded
database. Paired 10,000-iteration, seed-42, 120/300-second Balance and Arcane replays
have identical damage actions, auras, DPS and OOM time. Resource totals match after
removing the old synthetic row. The 300-second Balance build now attributes
2,965.21236 actual mana per fight to Innervate instead of 0.2.

Drop this patch when upstream reports actual Innervate gains without also counting
them as ordinary regeneration, including cap losses and passive-regen threat rules.

## 11. Holy Nova healing crit defense type (dropped, upstream #751)

Dropped in the 2026-10-10 sync: upstream #751 declares the heal's defense type from its client row
(`heal.DefenseTypeCore()`, magic for every rank). `TestHolyNovaHealingCritAllRanks` stays as coverage and passes.

Holy Nova's triggered party heal used `OutcomeHealingCrit` without declaring a
DefenseType. Its first healing crit panicked in `CritDamageMultiplier`, failing
the simulation. Declare `core.DefenseTypeMagic` on the shared heal configuration,
which gives all six ranks the existing 1.5 base healing crit multiplier.

`TestHolyNovaHealingCritAllRanks` casts each rank with a forced healing crit and
checks positive healing, a recorded crit, and the multiplier. A 100-iteration
Docker reproduction failed before the patch and completed with 309 healing crits
afterward. Keep the engine's strict missing-defense-type guard.

Drop this patch when upstream declares the same defense type for every rank of
the triggered Holy Nova heal and the regression passes.

## 12. Unset spell field guards

A field left at its Go zero value can mean "never set" rather than a real value.
Three of them were read without a check, and each fails silently or only by
chance:

- **DefenseType.** `CritDamageMultiplier` panics on `DefenseTypeNone`, but only
  once a crit lands. A spell with a low crit chance passes every test and fails in
  a user's sim, which is how the Holy Nova heal (patch 11) shipped. The three crit
  chance functions (`PhysicalCritChance`, `SpellCritChance`, `HealingCritChance`)
  now call `requireDefenseType`, so a missing DefenseType fails on the spell's
  first crit roll. Every crit path, including the attack-table and expected-damage
  helpers, goes through one of them. The strict panic in `CritDamageMultiplier`
  stays, as patch 11 asks.
- **SnapshotAttackerMultiplier.** An expired dot resets its snapshot to 0. A custom
  `OnSnapshot` that sets `SnapshotBaseDamage` but not the multiplier leaves that 0
  in place, and every tick deals or heals nothing. `TakeSnapshot` marks the
  multiplier with NaN before the callback, panics if a base amount was set without
  a multiplier, and otherwise restores the previous value. A multiplier the
  callback computes as 0 is kept.
- **Cast.CD / Cast.SharedCD.** A Timer without a Duration was already rejected. A
  Duration without a Timer is never read, so the spell silently had no cooldown;
  `RegisterSpell` now rejects it too.

The snapshot guard found one live case: the Bloodfang 8-piece heal set only
`SnapshotBaseDamage = 50` and had no `DamageMultiplier`, so every tick healed for
50 x 0. It now snapshots through `SnapshotHeal` with a multiplier of 1.

Separately, `gen_db` decoded the pinned Forever planner's `stats` as
`map[string]float64`. 10,287 rows carry an `appearances` object and 3 a `skillBuff`
object there, so `make db` panicked in `fillPlannerArmor` on the committed
`wowhead_forever_gearplanner.txt`. The map is now decoded loosely and only `armor`
and `armorbonus` are typed, with a panic naming the item if either is ever not a
number. A full `make db DB2TOOL_FLAGS="--cdn --dbcache ..."` run then completes.
The regenerated database is not part of this patch: its armor for the 14 items
whose values move matches the planner on 13 of them (the committed database on 1),
but the effect generator also swaps Skullflame Shield's active proc from
Flamestrike to Drain Life, which needs its own look.

Validation: `sim/core/unset_field_guards_test.go` covers a damage and a healing crit
roll without a DefenseType (panics at 0% crit), the same spell with one (casts), a
snapshot that sets a base amount but no multiplier (panics), a computed 0
multiplier (allowed), and both cooldown cases. The full suite,
`go test --tags=with_db $(go list ./sim/... | grep -v sim/web)`, has the same
results with and without the patch: every class suite passes, so no registered
spell trips a guard, and only the stale `sim/common/TestRegisteredEffects` baseline
fails, as it does on the base branch.

Drop this patch when upstream checks DefenseType at the crit roll, rejects an
unset snapshot multiplier and a cooldown Duration without a Timer, and decodes the
planner stats without assuming every value is a number.

## 13. Windfury Totem procs in Cat and Bear Form (dropped, upstream #669)

Classic's Windfury Totem enchanted the held weapon: the totem's periodic aura
(8515) cast 8514, which applied temporary enchant 1783. A shapeshifted druid does
not swing its weapon, so the Classic sim stripped `WindfuryTotem` from the party
buffs of both feral specs.

Forever changed the totem. In client 1.60.1.70009, 8515 is an area party aura
(effect 35, aura 42) with a 20% chance on melee autos and melee specials
(`ProcTypeMask` 0x14) that triggers 8516, the extra attack and its attack power.
Nothing in it names a weapon, so a cat or bear procs it like anyone else.
`feralcat.AddPartyBuffs` and `feralbear.AddPartyBuffs` no longer clear
`WindfuryTotem`. The shared driver already swings the current main-hand weapon,
which in form is the paw.

Validation: `sim/druid/feralcat/windfury_test.go` and
`sim/druid/feralbear/windfury_test.go` run full sims with a 3.5 s two-hander
equipped. With the totem they record Windfury extra attacks (tag 25584) at a
plausible rate per auto, and each extra attack hits for about what the form's own
auto does plus the proc's attack power, not like a swing of the two-hander.
Without the totem they record none. Both fail on the unpatched source. The Cat and
Bear suite goldens are regenerated because the suite's party buffs include the
totem; no other class's results move.

Drop this patch when upstream stops stripping Windfury Totem from the feral specs.

## 14. Forever Shadow crits: Devouring Plague ticks and Shadow Word: Death

Two Shadow crit behaviours on the live Forever server are not in beta client
1.60.1.70009, which the engine otherwise follows. Both were reported by players
running the Shadow benchmark rather than read from client data:

- **Devouring Plague ticks crit.** Every rank's row leaves Periodic Can Crit
  (`ATTR_EX_8_PERIODIC_CAN_CRIT`) off, so `priestTickOutcome(rank.PeriodicCanCrit())`
  never rolled. The benchmark report this fork produced shows 0 crit ticks in 72,232.
  The ticks now always roll for a crit; the heal still equals the damage dealt.
- **Shadowform's +100% critical strike damage bonus covers Shadow Word: Death.**
  The client's mask (41984016) names Mind Blast, Mind Flay, Shadow Word: Pain,
  Devouring Plague and Mana Burn. `PriestSpellShadowWordDeath` is added to the
  modifier's class mask, so Shadow Word: Death crits for 2.0 in Shadowform like the
  rest instead of 1.5.

The spell notes in `ui/sim/spells/priest.json` say the same. Measured on the current
tier Shadow benchmark (Undead, 120 s, 100,000 iterations): the preset rotation goes
from 581.10 to 588.66 DPS, about +2.0 from Devouring Plague and +5.6 from Shadow
Word: Death.

Validation: `sim/priest/shadow_crits_regression_test.go` forces a Devouring Plague
tick at 100% crit for every rank and checks the Shadow Word: Death and Mind Blast
crit multipliers with and without Shadowform. Both tests fail without the patch.
Only `TestShadowPriest` goldens move (re-blessed); `TestSmitePriest` is unchanged,
and the rest of the full suite matches the base branch.

Drop this patch when the client data sets Periodic Can Crit on the Devouring Plague
ranks and Shadowform's crit damage mask names Shadow Word: Death, or if a combat log
shows either behaviour is not live.

## 15. Rage log lines name the real maximum

Boundless Rage raises a Warrior's maximum rage by 10, 20 or 30 (`warrior.go`), and a
Gnome's Expansive Mind raises it 5% more (`racials.go`). Both were applied, but
`rageBar.AddRage` and `SpendRage` logged `of 100 total` whatever the bar held, unlike
energy, which logs its own maximum. MythicSim's timeline reads the maximum from those
lines, so every Warrior's rage row said "max 100". Both lines now log `rb.maxRage`.
Nothing the sim computes changes.

Validation: `sim/warrior/dps/max_rage_log_test.go` checks the log's maximum for a
Human without the talent (100), with Boundless Rage 3/3 (130) and a Gnome with it
(136.5). It fails on the unpatched source.

Drop this patch when upstream logs the rage bar's maximum.

## 16. Faerie Fire and Curse of Recklessness share their armor reduction (narrowed, upstream #671)

In client 1.60.1.70009 both take 505 armor (9907 and 11717), and Wowhead's Forever
class guides state they no longer stack. The client rows carry nothing that says so,
and the buff manifest gave each its own category, so a target with both lost 1010.
Both now sit in `MinorArmorReduction` with the new manifest option `PerStat`: each
stat the aura attaches bids alone in the category, so only the armor competes and
each keeps anything else it does. `gen_buffs.go` renders `PerStat: true` into the
`buffs.Meta`, and `newDebuff` bids with `spelldata.ExclusivePerStat` for it. The two
generated rows in `sim/core/buffs/debuffs_auto_gen.go` were edited to match what the
generator writes, since `gen_spelldata` needs the client database; the synthetic
render fixture covers the new field.

Validation: `sim/core/buffs/armor_reduction_test.go` applies each alone (-505) and
both together (-505, not -1010), and fails on the unpatched source. Every suite that
puts both on the target moves: all physical specs, pets, and one physical weapon
proc (Everlook Pathcarver) in the caster suites. Healer suites do not.

Drop this patch when the client data puts the two in one exclusive category, or if a
combat log shows them stacking on live.

## 17. A hard-cast Shaman spell holds the melee swing

Tested on the Forever beta: a Lightning Bolt cast between swings resets the swing
timer as it completes; a swing that comes due during the cast waits at zero and
lands as it completes; a bolt Maelstrom Weapon makes instant leaves the swing alone.
The engine only acted when the cast would end after the next swing, and then pushed
that swing a full swing past the cast, so a bolt that fit between swings cost no
melee at all. That made weaving 3 to 4 stack bolts between swings look free.

`AutoAttacks.HoldMeleeForCast` models the tested behaviour, and Lightning Bolt, Chain
Lightning and Lava Burst call it for any cast longer than zero. A held swing lands one
nanosecond after the cast completes: at the same instant it would wake the rotation
before the cast's completion runs, and a new cast would replace it, losing the
bolt's damage and cost.

Validation: `sim/shaman/enhancement/swing_hold_test.go` chain-casts an unhasted
2H Orc until it is out of mana, then melees, and checks from the log that no swing
lands inside a cast, that swings due during one land as it completes, that a swing
after a cast that completed first comes at least a full swing later, and that every
bolt started also completes. It fails on the unpatched source. The Enhancement suite
goldens move.

Drop this patch when upstream models a hard cast resetting the swing timer. Superseded in part by patch 105 (2026-10-11): a swing due mid-cast no longer lands as the cast completes.

## 18. Windfury Totem leaves the main-hand stone or oil alone (dropped, upstream #665)

Classic's Windfury Totem enchanted the held weapon, so `applyConsumeEffects` skipped
a main-hand stone or oil whenever the party had the totem. Forever's totem is a party
aura that procs extra attacks and names no weapon (see patch 13), so the main-hand
imbue now applies beside it. Rogue poisons were never affected: they register in
`sim/rogue/poisons.go`, not here.

Validation: `sim/warrior/dps/windfury_imbue_test.go` checks that a main-hand
Elemental Sharpening Stone adds 2% crit with and without the totem. It fails on the
unpatched source. No suite golden moves, since none pairs a main-hand stone with the
totem.

Drop this patch when upstream stops displacing the main-hand imbue under Windfury
Totem.

## 19. Rockbiter Weapon (dropped, upstream #673)

`ShamanImbue_RockbiterWeapon` was in the proto and the talent code named it, but nothing
registered it, so choosing it simmed the same as no imbue. Client 1.60.1.70009 defines
Rockbiter Weapon rank 7 (16316) as an imbue whose passive (16313) is a permanent melee
attack power aura: the tooltip reads "increasing melee attack power by 653 and allowing melee
attacks to cause additional threat when using that weapon", lasting 5 minutes. Elemental
Weapons raises it by 7, 13 and 20% (its effect 1, already named for Rockbiter's spell mask).

`RegisterRockbiterImbue` adds one permanent aura, "Rockbiter Weapon", with that attack power
and the talent's share, when the main or off hand carries the imbue. It is not in Windfury
Totem's exclusive category, so it stacks with the totem. A second Rockbiter weapon does not
add it twice, since the passive is one spell. The threat half and the 5 minute duration are
not modelled, the second because no imbue models its duration.

Measured on the Enhancement reference (120 +/- 15 s, 10,000 iterations, board seed), main
hand imbue with the party's Windfury Totem and without it: Windfury Weapon 597.2 and 597.0,
Rockbiter 624.7 and 545.4, Frostbrand 572.5 and 501.2, no imbue 513.3 and 450.9.

Validation: `sim/shaman/enhancement/rockbiter_test.go` checks that Rockbiter adds exactly 653
attack power in either or both hands and that Windfury Weapon adds none standing. It fails
on the unpatched source. No suite golden moves, since none imbues Rockbiter.

Drop this patch when upstream registers Rockbiter Weapon.

## 20. Flametongue Weapon keeps Windfury Totem's procs (dropped, upstream #674)

`RegisterFlametongueImbue` put a main-hand Flametongue Weapon in Windfury Totem's exclusive category at a
higher priority, copied from Classic, where the totem enchanted the main-hand weapon and any main-hand imbue
displaced it. Forever's totem is a party aura that names no weapon (patches 13 and 18), and the beta describes
only Windfury Weapon as disabling it: "When applied to mainhand, disables any benefit you personally benefit
from Windfury Totem". So with Flametongue in the main hand the totem stood up 96% of the fight and never
proced: the "Windfury Totem (proc)" aura showed 0% and the extra attack dealt nothing (reported by Kerani on
the Discord, 2026-09-29). Frostbrand and Rockbiter were never in the category.

The Flametongue block is removed; Windfury Weapon keeps its own.

Measured on the Enhancement reference (120 +/- 15 s, 10,000 iterations, board seed), Flametongue main hand
with the party's Windfury Totem: 510.5 before, 583.9 after, and the totem's proc aura is up 119.8 s of 120.
Every other imbue is unchanged.

Validation: `TestOnlyWindfuryWeaponDisplacesWindfuryTotem` in `sim/shaman/enhancement/rockbiter_test.go`
checks which main-hand imbues leave the totem's "Windfury Totem Trigger" active. It fails on the unpatched
source for Flametongue. No suite golden moves.

Drop this patch when upstream models the totem and Flametongue Weapon the same way.

## 21. Warlock demons inherit hit, crit, spell power and attack power

Upstream gives every warlock demon an empty stat inheritance ("Forever's demons inherit nothing
from the warlock"). The client has a hidden passive aura, Warlock Pet Scaling (416189, flagged
Owner Power Scaling), whose effects include hit, crit, spell haste, attack power, damage done and
more, with placeholder values in the database. The beta's pet stats show 100% of the warlock's hit
and crit, 10% of its spell power as the demon's spell power and 17% of it as attack power.

The patch inherits exactly those, dynamically, so a spell power proc reaches the demon at the
same rates. The warlock's spell hit and crit also fill the demon's melee hit and crit. Hit from a
talent (Suppression) was seen reaching the pet on the level 20 beta; hit from gear was not
testable there. Haste, Intellect, health, resistances and mana regen are also on the aura and are
not modelled. The demon keeps its own Agility and Intellect crit on top of the inherited crit.

Validation: `sim/warlock/pets_test.go` checks the mapping and that it is linear. On the
Demonology reference the Succubus's Lash of Pain goes from 17.1% miss and no crits to 10.1% miss
and 13.9% crit, its DPS from 61.1 to 84.2 and the warlock's total from 576.7 to 599.8. The
Affliction golden moves (naked, average 229.7 to 239.3 DPS); Destruction sacrifices its demon and
does not.

Drop this patch when upstream models the aura, and compare the demon's stats before switching.
The numbers are Sage's and the beta testers' readings, not client data.

## 22. Dots tick on the caster's current stats (dropped, upstream #668)

A dot built by `Dot.Snapshot` stored the caster's spell power share and whole attacker damage multiplier
when it landed, and every tick dealt those stored numbers. Forever's dots do not snapshot: the beta log of a
Gnome priest (foreverlogs.gg report 2668, encounter 4607) has Shadow Word: Pain ticking 34, 34, 34 and then
38, 38, 38, 39 on the same application when Eureka! is cast after it landed, with Eureka! the only change on
the priest (eight applications in the log agree, 37.8 inside the window against 34.2 outside). Redfall
reported the developers confirming it (Discord, 2026-09-29). Upstream's Rend already reads attack power at
each tick (#524: Battle Shout gained mid-bleed raises the next tick).

`core.Dot` now remembers what `Snapshot` and `SnapshotPhysical` folded in and swaps it at every tick: the
spell power share (`BonusCoefficient` times the caster's current spell power) and the attacker multiplier
(talents, personal buffs, Eureka!, school and target-table multipliers, `PeriodicDamageMultiplier`) are read
when the tick lands. Flat damage stays with the application: combo point values, a Deadly Poison stack count,
Bane of Agony's ramp, the T5 set bonus's edit of `SnapshotBaseDamage`. Crit chance and target-side modifiers
were already read at the tick. A dot that writes `SnapshotBaseDamage` by hand, and every snapshot heal,
ticks on the stored numbers as before.

A spell whose base damage includes a share of attack power declares it with `Dot.SnapshotAttackPowerShare`
(share, melee or ranged) right after the snapshot, and the tick reads that share from the attack power it has
then: Rip (1% a combo point, four at most), Rupture (its per point share, Hemorrhage scaling both halves),
Garrote (3%) and Serpent Sting (3.5% ranged). `CopyDotAndApply` and `SaveState`/`RestoreState` carry the two
remembered shares.

Measured on the 26 reference builds at pin `2e2fa64cd`, before patch 21's demon scaling (board seed, 10,000 iterations, empty trinket slots): Affliction Warlock
474.21 to 470.65 (its Gnome Eureka! used to lock +10% into the DoTs it cast at the pull), Destruction 507.63
to 510.76 (Improved Shadow Bolt now reaches running DoTs), Shadow Priest 615.16 to 617.98 (Shadow Weaving),
Feral Druid -0.18, Fire Mage -0.06, Frostfire Mage -0.07, Enhancement Shaman -0.04, Subtlety Rogue +0.08; the
other 18 are identical to four decimals.

Validation: `TestDotTicksOnCurrentSpellPower`, `TestDotTicksOnCurrentMultiplier` and
`TestDotTicksOnCurrentAttackPower` in `sim/core/dot_test.go` (the first replaces `TestDotSnapshotSpellDamage`,
which asserted the old rule), `TestDotTicksReadTheDamageMultiplierAtTheTick` in `sim/warlock` (Corruption's
tick multiplier steps up while Eureka! is up and back down after) and `TestSerpentStingTickReadsAttackPowerAtTheTick`
in `sim/hunter`. The warlock and hunter tests fail with the current-stat read disabled. Goldens moved:
Balance, Feral Cat, Beast Mastery, Marksmanship, Survival, Fire, Shadow Priest, Smite Priest, Subtlety,
Elemental, Affliction and Destruction.

Not covered, still on stored numbers: hunter pet abilities, Lacerating Strikes, two item-proc dots
(`sim/common/classic/items_weapons.go`, `sim/rogue/items.go`) and the healing hots. Rake, Lacerate and
Deep Wounds have no attack power share to read.

Drop this patch when upstream's dots tick on current stats the same way.

## 23. Rain of Fire (dropped, upstream #672)

Rain of Fire had client data (`spellData.RainOfFire`, 5740 to 11678, and Forever's
`RainOfFireTriggered`, 1282380 to 1282385) and a class mask (`WarlockSpellRainOfFire`, already in the
fire and Destruction groups), but no spell was registered, so a rotation line naming it was dropped
with "does not know spell" and Warlocks had no ground AoE in multi-target sims.

Forever's Rain of Fire is an area trigger, like its Blizzard: rank 4 (11678) is an 8 second channel
costing 1185 mana whose periodic dummy fires every 2 seconds and casts 1282385, a direct Fire hit for
221 (0.083 spell power coefficient) on every enemy in the area. `registerRainOfFire`
(`sim/warlock/rain_of_fire.go`) reads all of that from the client data: the channel is an AoE dot
with the dummy's period, each tick casts the damage spell with `CalcAndDealAoeDamage`, and the tick
rolls crit unless the client marks the damage spell as unable to (it does not). Both spells carry
`WarlockSpellRainOfFire`, so the existing fire and Destruction talent modifiers apply.

Measured with a gearless Warlock against three targets (30 s, 200 iterations): 2 channels a fight
(mana-bound), 4 ticks per channel on every target, about 211 per landed tick after partial
resists.

Validation: `TestRainOfFireRainsOnEveryTarget` in `sim/warlock/rain_of_fire_test.go`. On the
unpatched source the channel is never cast. No suite golden moves: no preset rotation casts it.

Drop this patch when upstream registers Rain of Fire.

### Patch 24: Wolfshead Helm cooldown resources

Forever item 8345 (effect 17768) grants 20 extra Energy from Tiger's Fury and
5 extra Rage from Enrage. It no longer grants resources on entering Cat or
Bear Form. Verified against the current Forever item tooltip on 30 September
2026: https://www.wowhead.com/forever/item=8345/wolfshead-helm . The in-game
screenshot posted in the MythicSim Discord shows the same wording.

Resource gains are applied at cast time while the item aura is active. The
Cat and Bear regression tests verify both the cooldown gains and absence of
shift gains, with and without the helm. Both fail on the previous pin. The
Wolfshead Trophy enchant is a separate effect and has not been changed.

Client 1.60.1.70170 rewords the item to "an additional 20 Energy from activating Shifting Power" (the Rage half is
unchanged), and patch 41 moves the Energy onto that spell.

## Upstream sync 2026-10-01

Merged ElliotWood/Forever `d91d4afe40` (113 commits since `8dc19a4241`, client 1.60.1.70124) into the
fork tip with a merge commit, so the patch history stays intact. Where upstream and a patch met:

- **Database.** Upstream's 15 changed item rows, 10 new items and 13 new spell icons were merged
  per field into our `db.json` (no conflicts inside a row), `db.bin` rebuilt with
  `go run ./tools/sync_db_binary`, and the Forever planner armor re-applied to the 11 rows the new
  planner file changes (what `fillPlannerArmor` does on a full regeneration).
- **Frostfire Bolt (patch 8) kept.** Upstream #521 registers the spell and wires the talent masks, but
  not the Frostfire school handling (lower resistance, no doubled school bonuses), binary resistance or
  the hybrid rotation. Upstream's duplicate registration was dropped; its mask wiring is used. Upstream's
  `TestFrostfireBolt` moved to seed 2, because binary resistance shifts which seeds hit.
- **Innervate attribution (patch 10) kept.** Upstream #554 stops the double count but does not
  attribute the regen to Innervate. Upstream #576 made the druid's own cast spend mana under
  `{29166, Tag: druid.Index}`, which collides with the attribution row for the first druid (Index 0), and
  the serializer merges rows sharing an ActionID. The attribution row now uses tag -2.
- **Demonic Brand (PR #15) kept.** Upstream #535 makes the Imp's brand Fire and lets pet spells spend it,
  but keeps one charge pool on the pet. Our per-target charges and `ProcMaskDirect` spending stay, and
  upstream's `TestDemonicBrandImpSpendsWithFirebolt` was removed because `TestDemonicBrandPetHits`
  covers the same two cases against the per-target model.
- **Rupture** keeps `ruptureAttackPowerShare`; upstream's change there was a comment and the
  `PointsPerResource` read, both taken.
- **Test leak fixed.** Upstream's `TestPetStrikes` zeroed the shared default target's armor, which changed
  every later sim in the hunter package (`TestSurvivalMelee` read 206 instead of 130 DPS). It clones the
  target now.
- **Goldens** were regenerated on the merged tree. Survival Melee moved 1.7% on average; the other specs
  track upstream's own golden movement (glancing blow rolls, dodge from agility, client mana costs, pet
  abilities).


## 26. A cast of an on-next-swing ability without its queue tag queues it

Heroic Strike, Cleave and Maul register their hit untagged and the APL queue under tag 1; Raptor
Strike's queue is tag 3. A rotation that names the spell without the tag, which is what a spell
picker produces, fell back to the untagged hit (`GetAPLSpell` prefers an APL-flagged spell, then takes
any spell with the same ActionID). That hit fired instantly, off the global cooldown and on top of
the white swing. A Heroic Strike only Fury rotation reported in the MythicSim Discord on 1 October
read 1,253 DPS, against 799 for the reference preset on the same engine; 80 of the 145 Heroic
Strikes in its first iteration landed within 250 ms of the previous one.

`GetAPLCastSpell` resolves a cast action: when the named ActionID exists but is not APL-flagged and
an APL-flagged spell shares its IDs under another tag, the cast presses that one. Only cast actions
use it. Value lookups (`spellTimeToReady`, `spellIsReady`) keep `GetAPLSpell`, because they read the
real spell's cooldown; redirecting those too changed the upstream Survival rotation's golden.

Validation: `TestUntaggedHeroicStrikeQueuesOntoTheSwing` in `sim/warrior/dps`. Unpatched, a Heroic
Strike only rotation records 83.4 white swings and 110.3 Heroic Strikes a fight against 83.7 swings
with no rotation; patched, Heroic Strike replaces the swing. No suite golden moves: the presets name
the queue tags already.

Drop this patch when upstream resolves untagged casts of queued abilities.

## 27. Feral forms swing the equipped weapon's DPS (dropped, upstream #670)

Forever's Druid class deep dive (worldofwarcraft.blizzard.com/en-us/news/24301515, 30 September
2026): "While Shapeshifted, the Druid's melee auto attack DPS ... in Bear Form, Cat Form, or Dire
Bear Form is now the same as the DPS of the Druid's equipped weapon. But the speed of attack is
changed to 1.0 or 2.5 seconds, depending on the form. Abilities that deal weapon damage likewise use
these damage values."

The paw was a fixed level 60 weapon (Cat 43.84 to 65.76 at 1.0 s, Bear 109 to 165 at 2.5 s), so an
equipped weapon's damage did nothing. `formWeapon` rescales the equipped main hand's damage range,
weapon damage enchants and bonus DPS to the form's swing: the paw keeps the weapon's DPS and spread.
With nothing equipped it is the unarmed fist. A player's paired character sheet measurements
(unarmed, a 9.0 DPS mace and a 20.4 DPS mace in Cat Form, posted to the MythicSim Discord on 1
October) agree with this to the sheet's whole-number rounding.

Validation: `TestFormPawCarriesTheEquippedWeaponDPS` in `sim/druid/feralcat`. The Cat and Bear
goldens move because their suite gear is naked: the unarmed paw now deals only attack power damage.

## 28. Elemental preset casts Fire Nova

`ui/specs/shaman/elemental/apls/forever.apl.json` casts Fire Nova (408345) above 30% mana, after
Flame Shock and before the Lightning Bolt filler, as submitted by a MythicSim player on 1 October.
On the MythicSim Elemental reference (10,000 iterations): 393.18 to 396.33 DPS on the board encounter,
393.06 to 394.50 at a fixed 120 seconds and 357.88 to 363.78 at 300 seconds. The Elemental golden
moves with it. This is a preset change, not an engine one: drop it if upstream's preset adopts Fire
Nova or measures it worse.

## Upstream sync 2026-10-01, second merge

Merged upstream `723ea18f32` (28 commits since `d91d4afe40`). No patch needed changing; the overlaps
were all clean text merges:

- Frostfire Bolt: upstream's range roll (#595) applies to our spell as is.
- Serpent Sting and Rip: upstream's comment and client-row reads (#592, #598) merge with our code.
- Database: merged per record with `scripts/forever-merge-db.py` (the 1 October hotfix delta: 13 changed
  lines in `db.json`, one in `leftover_db.json`).
- Goldens: Balance, Mage, Priest, Warlock, Shaman and Retribution moved with the nuke damage rolls (#595,
  #596) and two-handed Seal of Righteousness (#591); Protection Warrior with Thunder Clap (#584). The
  melee Survival golden (ours, no upstream equivalent) falls 6 to 8% with no buffs because Expose Prey
  now needs Hunter's Mark on the target (#594) and moves under 1% with full buffs.

## 29. `core: player options for the race comparison's weapons`

- **What it does.** Three `Player` fields. `disable_weapon_specialization` (60, JSON
  `disableWeaponSpecialization`) withholds the race's weapon specialization (Human sword +2% crit,
  Dwarf mace +1%, Orc axe +1%) even with its weapon equipped, in `applyWeaponSpecialization`.
  `weapon_type_override` (61) retypes the melee weapons in the main hand and off hand, keeping
  their stats, in `NewCharacter` (`Equipment.overrideWeaponTypes`); `weapon_type_override_off_hand_only`
  (62) leaves the main hand alone. Shields, off-hand items and ranged weapons are untouched.
- **Why.** The race comparison holds one gear set for every race, so a race's weapon
  specialization was active or not by accident of the preset's weapons (Orc ranked 9th of 10 on
  Fury, whose reference wields a mace and a sword). Retyping the weapons gives each race its own
  weapon type, and the rules that key on weapon type follow it: the specialization, the warrior's
  Weaponmaster, the rogue's Hack and Slash, Backstab and Ambush.
- **Tests.** `TestWeaponTypeOverrideGivesTheRaceItsWeapon` in `sim/core/disable_racials_test.go`.
- **Default.** All three are off, which changes no existing result.
- **Drop it when** upstream has an equivalent option. Point the worker's fields
  (`worker/cmd/refresh-forever-races`) at upstream's names first.

## 30. One air totem per party (dropped, upstream #675)

Redfall and Kerani (Discord, 2 October 2026) found a Windfury Totem from the party buffs and a Grace of Air
the rotation casts both up in one Enhancement sim. The engine had no air totem rule beyond the shaman's own
slot: `Shaman.AirTotemAura` replaces the shaman's previous cast totem, but the party's Windfury Totem
(`buffs.driveWindfuryTotem`, category `WindfuryTotem`), the party's Grace of Air (`buffs.driveGraceOfAirTotem`,
category `GraceOfAirTotemAgilityAdd`) and the two cast totems (`sim/shaman/totems.go`) were four independent
auras. On Redfall's build (sim a2f1441c) with a Rockbiter main hand, Windfury Totem from the party and Grace of Air
cast: 641.7 DPS, both up all fight, against 620.0 for the Windfury Totem alone and 558.6 for the cast Grace of
Air alone. Client build 70009 made Windfury Totem a party aura and its 2026-09-24 patch notes allow one air
totem per party (the `fcea407ab9` commit already built the melee and ranged presets around it).

- **What it does.** `buffs.AirTotemCategory` is a single-aura exclusive category (`sim/core/buffs/air_totem.go`).
  The party's Windfury Totem, the party's Grace of Air and the shaman's two cast totems each join it with a
  bid: party Grace of Air 1, party Windfury Totem 2, cast Grace of Air 3, cast Windfury Totem 4. The higher bid
  replaces the air totem that is up and the lower one is refused, so exactly one stands.
- **Precedence.** A totem the shaman casts replaces the one the party buffs assume. The cast is the later,
  deliberate placement (the rotation spends a global and mana on it), the party buff is a standing assumption
  that models someone else's totem, and the shaman's own slot already works the same way (a new air totem
  replaces the old one). When both party buffs are set, Windfury Totem outbids Grace of Air; MythicSim's worker
  never sends both. A party totem a cast has replaced is not restored if the cast totem later expires (the
  rotation recasts at once, so the gap is under a global). Totem twisting (`PartyBuffs.TotemTwisting`) keeps
  its own timing and stays out of the category, since it is by definition two air totems alternating; nothing
  in MythicSim sets it.
- **Windfury Weapon is not a member.** The beta describes it as disabling only "any benefit you personally
  benefit from Windfury Totem", so it stays in `WindfuryTotemCategory` (patch 20 and `weapon_imbues.go`) and
  leaves Grace of Air alone. A Windfury Weapon works under either air totem; the report that it "cannot be used
  with Grace of Air" was not the engine (Windfury Weapon with a cast Grace of Air: 9.9 procs a fight and Grace
  of Air up 119.9 s of 120 before this patch, the same after).
- **Uptime.** `Aura.Deactivate` no longer books a negative uptime for an aura that ends before the pull. A party
  totem that a pre-pull cast replaces ended at -3 s and reported an uptime of -3.0 s, as did any cast totem
  replaced in the pre-pull.
- **Tests.** `sim/shaman/enhancement/air_totem_test.go`: a cast Grace of Air replaces the party's Windfury Totem
  (no extra attacks from it), a cast Windfury Totem replaces the party's Grace of Air, the two cast totems
  replace each other, a request with both party totems holds Windfury Totem, and Windfury Weapon works beside a
  party or cast Grace of Air and still disables the Windfury Totem. The replacement tests and the negative
  uptime check fail on the unpatched source; the Windfury Weapon with Grace of Air test documents behaviour that
  was already right.
- **Default.** A request with at most one air totem is bit-identical: all 29 reference builds' default requests
  (MythicSim's one-click requests, with Windfury Totem, Grace of Air or neither) give the same seeded DPS to the
  last digit before and after, the Enhancement preset (Windfury Weapon, a cast Grace of Air, no party air totem)
  included, and every other spec's golden passes unchanged. The Enhancement suite golden
  (`TestEnhancement.results`) moves 168.69 to 153.42 (-9.0%) because upstream's Enhancement APL casts Grace of Air
  while the suite's party carries Windfury Totem and the suite wields no weapon, so it used to hold both.
- **Drop it when** upstream makes the air totems exclusive. Keep the uptime clamp either way.

## 31. Windfury Weapon's extra attacks have their own row (dropped, superseded by upstream #644)

Kerani could not see Windfury Weapon proc: with the imbue the sim report had no "Melee (extra attack)" row, which
the party's Windfury Totem has. The imbue did proc (8.5 to 9.9 times a fight on Redfall's build and the reference, the
`16361` attack power aura), but its two extra attacks were main-hand swings booked under the plain Melee row, so
nothing named them.

- **What it does.** `AutoAttacks.ExtraMHAttacksFrom(sim, count, spell)` is `ExtraMHAttacks` with the swings
  booked to `spell`, a copy of the main-hand swing under the granting spell's tag. Windfury Weapon grants its
  two attacks through it with the copy tagged 16361, so the report lists "Melee (Windfury Weapon)" and the plain
  Melee row counts only the swings of the timer. The totem's extra attack (tags 25584 party, 10610 cast) already
  worked this way.
- **No result moves.** The copy is the same swing: the same config, hit table, procs and timing. Seeded runs of
  Redfall's build in all 24 combinations of imbue, party air totem and cast air totem give the same DPS
  to the last digit as before the patch. Only the row the swings land in changes. `ExtraMHAttacks` itself is unchanged for its other
  callers (Hand of Justice, Hack and Slash, Ironfoe).
- **Tests.** `TestWindfuryWeaponExtraAttacksHaveTheirOwnRow` (two swings a proc under tag 16361; fails on the
  unpatched source) and `TestWindfuryWeaponGrantsExtraSwings`, which now counts the Melee row and the new row
  together.
- **Drop it when** upstream books the imbue's extra attacks under their own row.

## 32. `druid: a Feral Cat knows Moonfire`

- **What it does.** `RegisterFeralCatSpells` registers Moonfire (9835) and its dot, as `RegisterBalanceSpells` does for
  Balance. The client row (Moonkin stance, castable in caster form, not castable while shapeshifted) makes it a
  caster-form spell, so a cast from Cat Form leaves the form first, like every other caster spell in the form
  masks of `sim/druid/druid.go`. An APL that cancels Cat Form, casts Moonfire and shifts back works the same way.
- **Why.** A player's custom Cat APL weaves Moonfire to refresh the dot while powershifting. The Cat agent never
  registered the spell, so the engine dropped every line naming it ("does not know spell"), and the report
  adapter (`mythicsim-forever/cli/report_rotation.go`) then showed those lines as skipped and cut the Moonfire dot
  out of the conditions that mentioned it. The cancel-form line lost its "no Moonfire dot" test and powershifted
  every time Energy fell to 20: about 20 shifts a fight, 427 DPS for a character that makes 522 with no Moonfire
  lines. With the spell registered the same weave makes 534 (10,000 iterations, seed 42, the bare Moonfire line
  removed; 521.6 for the no-Moonfire rotation).
- **Tests.** `TestCatCanCastMoonfireFromCatForm` and `TestMoonfireWeaveRunsFromAnAPL` in
  `sim/druid/feralcat/moonfire_test.go`; both fail without the registration.
- **Default.** No preset casts Moonfire, so no existing result moves (`TestFeralCat` golden unchanged). Feral Bear
  and the other caster spells (Wrath, Starfire, Insect Swarm without the talent) are still unknown to the feral
  agents.
- **Drop it when** upstream registers Moonfire for the Feral Cat agent.

## 33. `core: an empty weapon slot never swings`

- **What it does.** `addWeaponAttack` skips a weapon whose swing speed is zero and leaves it disabled, so an
  empty slot (`Weapon{}` from `WeaponFromRanged` and `WeaponFromOffHand`) never joins the sim's swing list.
- **Why.** A hunter with no ranged weapon has `AutoSwingRanged` on and a zero-speed bow. Its auto shot scheduled
  itself again at the same instant, `advanceWeaponAttacks` never moved the clock, and one iteration ran until the
  process was killed (exit 137 at the worker's 1 GiB cap, in about 8 seconds). Two Quick sims of a Build Lab
  Survival Hunter with an almost empty gear list failed this way on 2 October 2026, and a single iteration
  reproduces it.
- **Tests.** `TestHunterWithNoRangedWeaponFinishes` in `sim/hunter/no_weapon_test.go` (no gear, a two-hander and no
  bow, a bow and no melee weapon); the first two hang on the unpatched source.
- **Default.** Every build with a weapon in each swinging slot behaves exactly as before; the melee, ranged and
  caster suites pass unchanged.
- **Drop it when** upstream guards zero-speed weapons in the swing loop.

## 40. `tools: gen_spelldata files a talent-granted ability on a ladder of its own`

- **What it does.** `discoverTraitLadders` skips a one-rank talent node whose spell is not passive, because such a node
  usually grants an ability the game teaches (Hemorrhage, Water Shield) and the ability's own ranks are the ladder.
  Shifting Power is a node on an ability only the talent grants, so it fell through and the Druid file had no ladder
  for it. `overrides.TalentGrantedAbilities` names spells that are the exception; the generator then takes the skill
  line row that grants the spell (AcquireMethod 3), as it does for Cat Form, and writes `ShiftingPower:
  spelldata.Ranked(1322605)`. The list is hand kept so no other class file moves.
- **Why.** `spellData.ShiftingPower` is how the new spell reads its cost, cooldown and energy off the client row.
- **Tests.** `TestShippedOverridesAreWellFormed` checks the entry states a reason and a source; the generated file is
  covered by `gen_spelldata -check` and the store by `TestStoreRegeneratesFromTheCommittedInputs`.
- **Default.** The only generated change is the one Druid ladder.
- **Drop it when** upstream's generator files one-rank talent abilities itself.

## 41. `druid: Shifting Power replaces Tiger's Fury`

- **What it does.** Client 1.60.1.70170 took Tiger's Fury out of the spell book and King of the Jungle out of the
  tree, and added Shifting Power (1322605, one rank) and Improved Shifting Power (1322670, two ranks). The sim drops
  the Tiger's Fury spell, aura and Cat Form exit hook (`sim/druid/tigers_fury.go`) and registers Shifting Power for a
  Cat that took the talent (`sim/druid/shifting_power.go`). The client rows state 55% of base mana (SpellPower, the
  same share Cat Form costs), a 40 Energy energize effect, a 16 second cooldown and a 1 second global cooldown
  (SpellCooldowns), and Cat Form as the only form (caster aura 768, shapeshift mask 1); the druid stays in Cat Form.
  Improved Shifting Power takes 4 and 8 seconds off the cooldown (its curve states -4000 and -8000 ms) through a
  `SpellMod_Cooldown_Flat` on the `DruidSpellShiftingPower` bit, which takes Tiger's Fury's slot. The tooltip says its
  cost is "reduced by effects that reduce the cost of Shapeshifting": Natural Shapeshifter's class mask (word 0,
  0xE0000000) holds the family bit Shifting Power carries (0x20000000) beside Cat, Bear and Moonkin Form, so
  `applyNaturalShapeshifter` names it too and the talent cuts both costs by the same 10% a rank. Clearcasting's mask
  (16870) does not name it. Wolfshead Helm (patch 24) now reads "an additional 20 Energy from activating Shifting
  Power", so its 20 Energy moved from Tiger's Fury to this spell: a cast gives 60 Energy with the helm and its 5 Rage
  from Enrage is unchanged. `tools/database/overrides/extra_spells.go` no longer pins Tiger's Fury (5217, 417045), so
  the store lost those two rows; an APL line that names 5217 is now dropped like any spell the engine does not know.
- **Why.** Without it the Cat kept casting a spell the client removed (+3.6% damage and a Wolfshead Energy gain, 482.56
  on the Cat reference at 70170 with King of the Jungle gone) and had nothing to spend the two new talents on.
- **Tests.** `sim/druid/feralcat/shifting_power_test.go`: a cast spends 55% of base mana and gives 40 Energy, costs
  what Cat Form costs, takes a 1 second global cooldown and leaves the form; the cooldown is 16, 12 and 8 seconds at 0,
  1 and 2 ranks of Improved Shifting Power; Natural Shapeshifter ranks 0 to 3 cut Shifting Power and Cat Form alike;
  the spell needs Cat Form and the talent; `TestTigersFuryIsGone`. `TestWolfsheadResourcesComeFromCooldowns` checks 40
  and 60 Energy, `TestClearcastingSpentByNextCostedAbility` uses Shifting Power as the spell outside the mask.
  `sim/core/spelldata` pins its caster aura and stance mask. Every one fails on the previous source.
- **Not modelled.** Howling Idol (item 272427, effect 1291059) and the Tier 1 Feral 5 piece (1301247) now reduce
  Shifting Power's cooldown, as they reduced Tiger's Fury's. Neither was ever applied.
- **Drop it when** upstream implements Shifting Power and removes Tiger's Fury.

## 42. `druid: the Feral Cat default rotation shifts with Shifting Power`

- **What it does.** `ui/specs/druid/feralcat/apls/default.apl.json` loses the Tiger's Fury lines (prepull and
  priority) and the two powershift lines (`cancelAura` of Cat Form), gains a Shifting Power line after Shred and
  Claw (Energy at or under 60, so it fires when nothing else can), and moves the out-of-form lines (Goblin Sapper,
  Demonic Rune, Major Mana Potion, Innervate, then Cat Form) to the top. Three in-form lines open that window: a
  rune, a potion or Innervate that is ready, with room in the mana pool (the original lines' 1500, 2250 and 40%),
  Furor known and Energy at or under 30. Furor's carry-over keeps the Energy through the shift, so the window costs the
  Cat Form recast (684 Mana, 1.5 s of global cooldown) and nothing else. The Furor and Energy conditions keep a build
  without Furor from throwing its Energy away. The two Cat presets in `presets.ts` and the default talents spend the
  three points King of the Jungle left on Shifting Power and Improved Shifting Power 2/2.
- **Why.** Shifting Power makes mana the Cat's second resource. On the Cat reference (10,000 iterations, board seed,
  Improved Shifting Power 2/2) the ShP line alone runs out of mana for 37 s of 120 and 210 s of 300, and the three
  consumables fix it: 546.31 to 563.51 on the board, 545.05 to 564.47 at 120 s and 510.81 to 530.72 at 300 s. The
  powershift lines the default shipped with burn that mana on 15 Energy a shift: with them kept, Shifting Power on top
  makes 531.72, against 563.51 without them. A Moonfire weave is no longer worth it either (465.74).
- **Tests.** `TestFeralCat` golden re-blessed (212.08 to 261.81 on naked gear, Average-Default; the 212.08 already lacked King of the
  Jungle's Energy, so the rise is Shifting Power, its two talents and the new lines). No other golden moves.
- **Default.** The reference Cat (`050022-55000032121032212051-052`): 563.51 board, 564.47 at 120 s, 530.72 at 300 s,
  against 519.25 before client 70170.
- **Drop it when** upstream ships its own Shifting Power rotation, and keep the Furor-gated shift lines in the
  presets that MythicSim's `scripts/lib/forever-rotation-adjustments.mjs` generates.

## Client 1.60.1.70170 (2 October 2026), interim

**Superseded on 2026-10-02 by upstream's own `[DB] Update to 1.60.1.70170 with hotfixes as of 2026-10-02` (#608).** The sync took upstream's
`db.json` (merged per record with the fork's stat-index and armor corrections), `db.bin`, `leftover_db.*`, spell store inputs,
store, class spell data, talent trees, protos and item procs, and `interimHotfixItemSpells` is gone (the real hotfix overlay
roots those items again). The notes below stay as the history of what the regeneration changed.

The branch `mythicsim/client-70170` merges upstream `696a6c4040` (4 commits: Master of Elements refunds once
per cast, the 2026-10-01 Wowhead data refresh, Warrior Dual Wield Specialization no longer raises off-hand
rage) and then regenerates the client data for build 1.60.1.70170. Upstream had not merged its own
`[DB] Update to 1.60.1.70170` when this was cut, and its Update DB workflow would fail on this build
at `gen_spelldata` for the reasons below.

**The regeneration is interim.** `db2tool --cdn` ran without a hotfix cache (`DBCache.bin`), because the
workflow's download was not available to the run. Without it the CDN tables carry no ItemSparse rows for
about 70 spells' worth of items that `db.json` holds, so the store lost 77 item roots, and `gen_db` would
rewrite the item procs, enchants and `leftover_db` with those items gone. So only what does not depend on
item rows was kept: the spell store (`spells_auto_gen.go`, `spell_store_inputs.json`), the class
`spell_data_auto_gen.go` files, the talent trees and `proto/*.proto`. `db.json`, `db.bin`, `leftover_db.*`,
`enchants_auto_gen.go`, `stat_bonus_*_auto_gen.go`, `missing_effects_auto_gen.ts` and
`forever_client_build.txt` are the 70124 ones. The 77 lost roots are pinned in
`tools/database/overrides/extra_spells.go` (`interimHotfixItemSpells`). To finish: download the mirror the
workflow uses into `tools/db2tool/caches/DBCache.bin`, run `make db DB2TOOL_FLAGS="--cdn --dbcache
tools/db2tool/caches/DBCache.bin"`, then `make basestats`, `go run ./tools/database/gen_spelldata`,
`make proto`, `go run ./tools/database/gen_db -gen=go-to-ts`, drop `interimHotfixItemSpells`, and re-bless.
`tools/database` `TestEveryReachableProcIsSupportedOrListed` fails until then: the capture lacks the pinned
roots, and Rotmending (1321587) and Siphon the Grave (1321560) are new unsupported procs.

What the regeneration changed that the sim read by position or by name:

- **Renamed talents.** The client renamed Hot Streak to Heating Up and Soul Harvesting to Soul Harvest, so the
  proto fields, the generated ladders (`HeatingUp`, `SoulHarvest`) and the tree `fieldName`s changed.
- **Druid tree.** King of the Jungle is gone, Shifting Power (1322605) and Improved Shifting Power (1322670, 2
  ranks) are in, Shredding Attacks moved up a row and Predatory Instincts to (4,3). The Feral tab is 20
  talents, not 19, and every proto field number after Shredding Attacks moved. Talent strings read by position,
  so `TalentTreeSizes` in `sim/druid/druid.go` is now `{16, 20, 16}` (a stale size shifts every Restoration
  talent by one), and the Feral section of every druid string in the engine was translated by talent name:
  `-5521002023132213051-05503` is `-55210032020132012051-05503`, `-5003232120132010501-0550325` is
  `-50032302120132010501-0550325`, `050022-5500002123032213051-052` is
  `050022-55000032120032012051-052`. King of the Jungle's 3 points are unspent in the two Feral presets.
- **Tiger's Fury stays registered.** The client dropped it from the spell book and tree but kept the rows, so
  `extra_spells.go` pins 5217 and 417045 and `sim/druid/tigers_fury.go` reads 5217 by id. Its King of the
  Jungle energy (20 a rank) is gone with the talent, which is the whole -9.9% on the Feral Cat golden. The
  Shifting Power work should replace it, Wolfshead Helm (patch 24) included. Patches 40 to 42 did: the spell,
  the pins and the file are gone.
- **Flametongue Weapon ladder (patches 20 and 31).** The 70170 tooltips name a Flametongue Attack spell (10444,
  29469, 29470) beside each rank's proc, so the generated `FlametongueWeaponTriggered` ladder is nine spells in
  id order and `Highest()` became a spell with no damage. `sim/shaman/weapon_imbues.go` takes 16344 by id.
  Any ladder built from `triggeredSpells` can shift this way when a tooltip gains a `$id` reference.
- **Insight and Increased Spirit.** Insight's buff 1299796 and Increased Spirit 1248751 (Mystic Mushroom) went
  from `A_MOD_PERCENT_STAT` on Spirit to `A_MOD_TOTAL_STAT_PERCENTAGE` with no stat named, which the engine
  read as Strength while the tooltips say Spirit. Resolved by patch 63: the parser reads both rows as Spirit,
  `TestInsightMultipliesSpirit` runs again, and the Mystic Mushroom and Insight rows in the goldens that moved
  with Strength move back.
- **Hard-coded mirrors that ignore the row.** Gnome Eureka! (`sim/core/racials.go`) applies to non-periodic
  abilities only: resolved by patch 62, which states each class's three lists from the client rows. Elemental
  Focus's Clearcasting mask (`sim/shaman/talents_elemental.go`) lost two class bits (word 1 0x40000, word 3
  0x40000000) and gained Fire Nova's: not a change to the spells it names (see "Caster rows that needed no
  engine change" below), and `TestClearcastingNamesTheSpellsTheClientRowDoes` holds the sim's list to the row.
- **Fork patches.** Patch 14's Devouring Plague half is now in the client data (Periodic Can Crit set on 2944 and
  19276 to 19280), so patch 60 reads it from the rank instead of rolling unconditionally; the Shadowform mask half
  still stands. Patch 27 reads the form passives only by hard-coded behaviour. Patch 24
  names Tiger's Fury, which the client replaces with Shifting Power.
- **Tooltips.** Holy Shield's block chance is 30 again (was 20) and Sniper Shot is 8 to 45 yd with a range bonus
  on the next 3 shots; both `ui/sim/spells` rows and the tooltip manifest's allow list follow.

Golden movements on the regenerated data (all explained, none from a silent misread):

| Golden | Average-Default DPS | Cause |
|---|---|---|
| TestFeralCat | 235.46 to 212.08 (-9.9%) | King of the Jungle removed (restoring its 60 energy gives the old number exactly) |
| TestFire | 189.89 to 185.79 (-2.2%) | Combustion charges 4 to 3 (restoring the row gives the old number) |
| TestRetribution | 408.98 to 391.64 (-4.2%) | Champion of the Light 33/66/100 to 20/40/60 (restoring the curve gives the old number) |
| TestProtection | 258.72 to 260.78 (+0.8%) | Holy Shield block 20 to 30 (+6.0) and Redoubt 6..30 to 4..20 (-4.1) |
| TestSurvival | 290.26 unchanged | Deflection parry 10 to 5 only |
| TestBalance, FeralBear, FeralCat, SurvivalMelee, Arcane, Frost, Shadow, Smite, Elemental, Enhancement, Affliction, Destruction, Arms | unchanged | one to three AllItems rows each: Searing Dagger (Sear 1291568), Mystic Mushroom (Increased Spirit), Plaguefang (Poison 1309315), Insight enchant |

## Client 1.60.1.70170: the Warrior and Bear Rage items (patches 50 to 52)

Numbered from 50 so they do not collide with the other 70170 branches. Each patch states the reading it
took of the patch notes, because for all three the client rows carry nothing the sim could read.

## 50. A critical auto attack gives 75% more Rage (dropped, upstream #609)

- **Dropped in the 2026-10-02 sync, upstream #609.** Upstream multiplies a critical auto attack's Rage by
  `core.CritRageMultiplier` (1.75) inside the shared rage bar, which every Rage user takes: Warrior, Feral Bear and
  Feral Cat's Bear Form (the Cat's bar is only ever used in Bear Form). This patch made it an opt-in
  `RageBarOptions.CritRageBonus` (0.75) that those three set, so with every user opted in the option protected
  nothing and upstream's constant is the one mechanism. The tests stay, and each pays the bonus exactly once (a ratio
  of 1.75, never 3.06): `sim/core/rage_crit_test.go`, `sim/warrior/dps/crit_rage_test.go`,
  `sim/druid/feralbear/crit_rage_test.go` and the new `sim/druid/feralcat/crit_rage_test.go` (Bear Form aura on the
  Cat). The merge moved the Fury golden only for #606's changes: with crit Rage off on both sides the merge moves Fury
  +6.88% and Arms 0.00% (upstream's own range without #609: +6.9% and 0.0%), and #609 alone is +5.1% Fury and +7.1%
  Arms. The text below is the history of the reading.
- **What it did.** `RageBarOptions.CritRageBonus` (0.75, `core.CritAutoAttackRageBonus`) multiplied the Rage a
  swing gives by 1.75 when it crits. Warrior, Feral Bear and Feral Cat set it.
- **The notes.** Warrior: "Players now generate 75% increased Rage when landing a critical strike with a basic
  attack." Druid: "Bear Form and Dire Bear Form now generate 75% increased Rage when landing a Critical Strike."
- **What the client says.** Nothing numeric. The 70124 to 70170 diff has no row with 75 on Rage. It does add the
  hooks the server needs: Bear Form (Passive) 1178 and Dire Bear Form (Passive) 9635 gain a SpellAuraOptions row
  (ProcChance 100, ProcTypeMask 0x4 = a melee auto attack landed) with no effect that uses it; a new Warrior
  passive "Rule of Rage (DND)" 1322574 (class set 4, label 25) has the same 100% / 0x4 proc and an A_DUMMY effect
  of 10; and 1313291, the energize that Dual Wield Specialization used for its off-hand Rage, is renamed "Rule
  of Rage (DND)" and loses its class set. Both warrior rows are generated into `spellData.RuleOfRage` and read by
  nothing. The proc mask is 0x4 on all three, so the client's own wording of "basic attack" is the auto attack
  of either hand: no ability proc flag (0x10 and up) is on any of them. The crit condition is a server hit mask
  that the DB2 tables do not carry, which is why no row names it. The dummy's 10 is not 75 and is left unread.
- **Reading.** A crit auto attack, main hand or off hand, one-hander or two-hander, pays 1.75 times what the
  same swing pays as a hit. The engine's rage formula does not depend on the damage dealt (3.46 a second of
  weapon speed, 4.5 for a two-hander, half for the off hand, from `f9f9f21883`), so a crit that deals twice a
  hit's damage pays 1.75 times, not 2. A glancing blow is not a crit. Abilities (Heroic Strike, Cleave, Mortal
  Strike, Maul) pay no Rage on a hit in this engine, so "basic attack" against the Bear's "Critical Strike" does
  not separate them here. Rage from damage taken is a separate rule and does not change.
- **Alternatives, measured** (10,000 iterations on the reference builds, DPS against this patch, 120 s / 300 s).
  (a) The bonus also multiplies Unbridled Wrath's Rage when it procs off a crit white hit: Fury +0.4% / +0.3%,
  Arms +0.4% / +0.2%, Fury-Protection +0.2%, Protection +0.1%. (b) The bonus also multiplies the Bear's Blood
  Frenzy Rage (5 Rage on a crit): Feral Bear +3.8% / +3.9% DPS, +5.6% / +5.6% TPS, because the Bear is
  rage-starved. Not taken: the client's proc mask is a melee auto attack, and Blood Frenzy's Rage is its own
  energize.
- **Tests.** `sim/core/rage_crit_test.go` (1.75 times a hit in each hand and for a two-hander, 15.743 against
  8.996 for a 2.6 speed main hand, damage taken unchanged; the "bar without the bonus" case went with the option),
  `sim/core/rage_test.go` (the crit rows now pay 15.743), `sim/warrior/dps/crit_rage_test.go`,
  `sim/druid/feralbear/crit_rage_test.go`. The first and the crit rows of the second fail on the unpatched source.
- **Default.** Warrior and Bear goldens move up (Average-Default DPS: Fury +6.6%, Arms +7.3%, Protection +0.3%,
  Feral Bear +1.1%); Feral Cat is unchanged.
- **Drop it when** (now moot) upstream models the client's Rule of Rage, or measures a different factor: then change
  `core.CritRageMultiplier`.

## 51. Swipe gains 3% of attack power

- **What it does.** `sim/druid/swipe.go` adds `0.03 * attack power` to each target's base Swipe damage, before
  Feral Instinct's and the other damage mods.
- **The note.** "Fixed a bug causing Swipe to not scale with Attack Power. It will now correctly gain 3% of the
  Druid's attack power" and, in the same line, that the tooltip will not update.
- **What the client says.** Nothing: none of the five Swipe ranks (779, 780, 769, 9754, 9908) has a
  BonusCoefficientFromAP, and the 70170 rows are unchanged, which is what "tooltip will not update" means. So the
  coefficient is stated in the engine (`swipeAttackPowerCoefficient`) and not read by row. Only the Bear has
  Swipe (family 7, stance mask 0x90); the Cat has none, and the Family 9 "Swipe" rows (1264494 to 1264502) are a
  Hunter pet's.
- **Tests.** `sim/druid/feralbear/swipe_test.go` casts the same roll at two attack powers and solves for the
  coefficient, without naming the modifiers: it reads 0.03 and fails on the unpatched source (the hit does not move).
- **Default.** The default Bear rotation casts no Swipe on one target, so no golden moves. On the reference Bear
  with a Swipe line before Maul on two or more targets (the preset has none): 2 targets +4.1% / +4.4% DPS
  (+4.5% / +4.8% TPS), 3 targets +5.0% / +5.4% DPS (+5.3% / +5.6% TPS) at 120 s / 300 s.
- **Drop it when** the client carries the coefficient (then read it by row).

## 52. Spearing Strike requires Battle Stance

- **What it does.** `registerSpearingStrike` (`sim/warrior/talents_arms.go`) casts only in Battle Stance.
- **The note.** "Arms: Spearing Strike no longer requires a 2handed weapon. Spearing Strike requires Battle Stance."
- **What the client says.** Spell 1310222 gains a SpellShapeshift row (ShapeshiftMask 0x10000, Battle Stance) and its
  SpellEquippedItems subclass mask goes from 1378 (two-handed axes, maces, swords, polearms, staves) to 173555
  (every melee weapon type), plus a ProcChance 101 row. The weapon half needs no change: the sim never checked the
  weapon, and the swing uses the main hand's own normalized damage. The stance half is new and stated in the
  engine like every other warrior stance requirement (`StanceMatches`), since the warrior abilities do not read
  StanceMask.
- **Tests.** `sim/warrior/dps/spearing_strike_test.go`: castable in Battle Stance with a two-hander or a one-hander,
  not in Berserker or Defensive Stance; the two refusals fail on the unpatched source.
- **Default.** The Arms preset swaps to Berserker Stance only in the execute phase, where it never reaches
  Spearing Strike, so the reference Arms request is bit-identical before and after this patch. The engine's own
  Arms golden moves -0.28% Average-Default (-2.85% to +4.13% across rows): its `dps_reck` and `dps_no_reck` APLs
  fight in Berserker Stance throughout and never cast Spearing Strike now.
- **Drop it when** upstream gates Spearing Strike by the client's stance mask.

### The rest of the Warrior and Druid notes: nothing to change

Checked against the 70124 to 70170 diff and the engine; tests pin the current behaviour so a later change shows.

- **Faerie Fire no longer resets your swing timer.** The row loses its SpellInterrupts entry (InterruptFlags 8) on
  770, 778, 9749 and 9907. The engine never moved the swing for Faerie Fire (an instant on the global cooldown;
  nothing in `sim/druid/faerie_fire.go` or the cast path calls the swing functions), so it already matches.
  `sim/druid/feralbear/faerie_fire_swing_test.go` and `sim/druid/feralcat/faerie_fire_swing_test.go`.
- **Queueing Heroic Strike no longer increases off-hand hit chance.** The Warrior engine already lifts the dual wield
  miss penalty for the queued swing alone (`calcQueuedSwing`), so an off-hand auto against a queued Heroic Strike
  rolls the same table as without it. `sim/warrior/dps/off_hand_hit_test.go`. The Hunter's Raptor Strike does hold
  the penalty off for the whole queue (`sim/hunter/raptor_strike.go`, `makeRaptorStrikeQueueSpell`), so a
  dual-wielding Hunter's off hand still gets the better table while one is queued: left alone, as the note names
  Heroic Strike only.
- **Unbridled Rage proc chance.** The talent is Unbridled Wrath, 12322. The build changes the spell-level ProcChance
  60 to 100 and nothing else; the rank values are the trait curve on its effect (12/24/36/48/60, the tooltip's
  `$m1%`), which did not change, and the engine rolls the curve (`ProcChanceEffectN`, `registerUnbridledWrath`)
  and never the row. Read as "the real chance was the row times the curve, now the curve", the engine already
  rolled the intended chance. Read as a flat 100% at any rank, a 5/5 warrior gains Fury +0.7% / +0.5%, Arms +0.7% /
  +0.6%, Fury-Protection +0.5% / +0.4%, Protection +0.4% / +0.35% (120 s / 300 s). Not taken.
  `sim/warrior/unbridled_wrath_test.go`.
- **Dual Wield Specialization.** The rage half went with upstream #601 and came back in #606 as a dummy effect, 10% a
  rank on the off hand's Rage. The hit half was aura 54 on effect 1, applied as an off-hand-only mod (the intended
  2/4/6/8/10%, not the both-hands behaviour the build lists as a known issue); #606 moved it to Furious Precision
  (4/7/10 at 1/2/3 ranks), which is the same off-hand-only mod. `TestFuriousPrecisionHitChanceIsOffHandOnly`.
- **Rend and Sunder Armor tap enemies instantly.** Attributes_6 0x800000 (TAPS_IMMEDIATELY) on the Rend and Sunder
  Armor ranks. Tagging has no effect in a sim with one boss.

## 60. `priest: Inner Focus's crit follows the client's non-periodic list, Devouring Plague reads its crit from the row`

- **What it does.** Inner Focus (14751) raises the crit of the next spell by 25. Its effect 1 names the spells by class
  mask, and client 1.60.1.70170 ("non-periodic" in its tooltip) took Devouring Plague and Shadow Word: Pain off the
  mask and put Mind Flay and Starshards on it; Shadow Word: Death was never on it. `applyInnerFocus`
  (`sim/priest/talents_discipline.go`) now leaves out Shadow Word: Death, Devouring Plague and Shadow Word: Pain
  (it left out Mind Flay, Shadow Word: Death and Starshards). Mind Flay and Starshards gain nothing in practice: the
  buff is spent by the cast that starts their channel, before the first tick. The Devouring Plague half of patch
  14 is retired: `priestTickOutcome(rank.PeriodicCanCrit(), dot)` replaces the unconditional roll, because every
  rank now carries Periodic Can Crit (2944 and 19276 to 19280). The two are bit-identical (the same output at the
  same seed on the Shadow and Smite references), and `TestEveryDevouringPlagueRankCarriesPeriodicCanCrit` fails the
  day a rank loses the attribute. The Shadowform half of patch 14 (Shadow Word: Death's crit damage) stands.
- **Why.** The Shadow rotation casts Inner Focus ahead of a Mind Blast hard cast, and the buff stays up until that
  cast completes, so the Shadow Word: Pain and Devouring Plague ticks that land in that window used to roll 25%
  more crit.
- **Measured** (same seed, 200,000 iterations, 120 s and 300 s): Shadow Priest 613.58 to 613.03 and 611.77 to
  611.29 (-0.09% and -0.08%); Smite Priest 440.17 to 439.74 and 416.16 to 415.84 (-0.10% and -0.08%). The
  Devouring Plague line alone moves nothing.
- **Tests.** `sim/priest/client_70170_test.go`: `TestInnerFocusCritIsNonPeriodic`,
  `TestEveryDevouringPlagueRankCarriesPeriodicCanCrit`; `TestDevouringPlagueTicksCrit` still passes.
- **Drop it when** never: it follows the row. Patch 14's Shadowform half is the part that stays.

## 61. `warlock: Hellfire's ticks can crit` (dropped, upstream #605)

- **Dropped in the 2026-10-02 sync, upstream #605.** Upstream reads the same flag off the Hellfire Effect row
  (`burnCanCrit`, from `Cannot Crit` on the triggered spell, so a hotfix that flips it back is followed) and checks the
  warlock's death against the base tick it burns instead of the first target's post-crit hit, which is the better rule.
  The fork keeps `sim/warlock/client_70170_test.go` (`TestHellfireTicksCanCrit`), which passes on upstream's code.
  The original reading follows.
- **What it did.** Hellfire Effect (5857, 11681, 11682) lost Cannot Crit (Attributes_2 0x20000000) and gained Periodic
  Can Crit (Attributes_8 0x200) in client 70170. The sim deals Hellfire's hits as ticks of the channel, so
  `sim/warlock/hellfire.go` rolls `OutcomeTickMagicHitAndCrit` on each target instead of a plain hit. The warlock
  still burns the base tick, before any crit, as the client's self damage is its own spell.
- **Measured.** No reference casts Hellfire (none of the presets or the engine's APLs name 1949, 11683 or 11684),
  so no reference moves. A Hellfire tick gains the crit chance times half again its damage.
- **Tests.** `sim/warlock/client_70170_test.go` (`TestHellfireTicksCanCrit`, which fails without the change).
- **Drop it when** never: it follows the row.

## 62. `core: Eureka! states each class's lists, and leaves periodic effects out`

- **What it does.** Gnome Eureka! (rows 1259812 Rogue, 1259813 Warrior, 1259817 Mage, 1259821 Warlock, 1259823 Priest)
  gives the next three casts of its listed abilities -10% cost and +10% damage. Client 70170 ("no longer
  benefits periodic effects at all. Channeled spells do not count as periodics") rewrote the lists:
  effect 0 (cost) and effect 1 (damage) lost every dot, effect 2 (the periodic bonus) kept only the channels
  (priest: Mind Flay, Penance, Starshards; warlock: Drain Life, Drain Soul, Wrack; rogue and warrior: a dummy now;
  mage: no mask at all), and the lists gained spells (mage: Pyroblast, Frost Nova, the Arcane Missile tick;
  warlock: Hellfire, Haunt; warrior: Intercept, Pummel, Revenge, Shield Bash, Spearing Strike; rogue: Hemorrhage;
  priest: Shadow Word: Death). Core used to take every class ability that deals damage and could not see the
  masks. `core.EurekaSpells` is now the three lists as class masks, which each class states in `eureka.go` next to
  its spells (`sim/mage`, `sim/warlock`, `sim/priest`, `sim/rogue`, `sim/warrior`), and `applyEureka` builds from
  them: the cost cut on the cost list, +10% on the damage and tick lists, the dots of a spell on the damage list
  alone give the 10% back (a spell's multiplier covers its hits and its ticks), and a charge is spent by a cast of
  any spell on any of the lists (assumed to be the client's own proc rule, a charge following the union of the
  effect masks as a talent proc's class mask does; alternative 1 below prices the other reading). A cast of Corruption,
  Curse of Agony, Siphon Life, Shadow Word: Pain or Rend now neither costs less nor spends a charge, and no tick of
  them gains anything.
- **The sim's lists match the rows.** The racial rows are not in the spell store, so each class test copies the
  70170 masks (family and four words) from the client and checks every spell the class registers against them,
  so a hotfix that moves a list fails there. The few differences are the sim's own shape and are named in the tests:
  Immolate's dot is its own spell, Hellfire's hits are ticks of the channel, Penance's bolts come from a channel
  spell the sim folds into the cast.
- **Measured** (same seed, 200,000 iterations, 120 s and 300 s, base is the regenerated 70170 data):

  | Reference (Gnome) | 120 s before | 120 s after | 300 s before | 300 s after |
  |---|---|---|---|---|
  | Affliction Warlock | 542.22 | 540.53 (-0.31%) | 534.35 | 533.21 (-0.21%) |
  | Demonology Warlock | 611.78 | 610.30 (-0.24%) | 598.46 | 597.69 (-0.13%) |
  | Fire Mage | 588.90 | 588.84 (-0.01%) | 575.09 | 574.97 (-0.02%) |
  | Frostfire Mage | 596.41 | 596.37 (-0.01%) | 567.53 | 567.38 (-0.03%) |

  Those four are the only Gnome references, and Arcane, Frost and the Destruction Warlock do not race Gnome.
  Goldens: `TestFire` Average-Default 185.791 to 185.769 (-0.01%); the Gnome rows of `TestAffliction` and
  `TestDestruction` (naked, no gear to spend) move between -0.06% and +3.9% (Affliction without buffs, short fight,
  169.11 to 175.75), because the charges are no longer spent on Corruption and the curses.
- **Alternative readings.**
  (1) Spend a charge on any class spell, as the old code did: identical at 120 s and 0.08% (Affliction) to 0.14%
  (Demonology) lower at 300 s, where a second Eureka! window opens.
  (2) The Mage's third effect keeps aura 108 and misc 22 with an empty mask. If an empty mask meant the whole
  family, as TrinityCore reads it, every Mage dot would gain 10% again, against the patch note. The Rogue and
  Warrior rows turn the same effect into a dummy, which settles the intent; the Mage mask is read as empty.
  (3) The client lists the channels (Mind Flay, Drain Life) on both the damage and the periodic effect where the
  damage one only matters for a hit of a triggered spell; the sim gives them the 10% once.
- **Tests.** `TestEurekaListsMatchTheClientRow` in `sim/mage`, `sim/warlock`, `sim/priest`, `sim/rogue` and
  `sim/warrior/dps`, `TestEurekaRaisesHitsButNotDots` (mage), `TestEurekaPassesOverTheDots` (warlock, priest),
  `TestEurekaChargesGoToListedSpellsOnly` (mage). `TestDotTicksReadTheDamageMultiplierAtTheTick`
  (`sim/warlock`, patch 22) used Eureka! as its multiplier; it now raises Corruption's own damage multiplier, which
  is the number Eureka! raised. Patch 22's log evidence (a Gnome priest's Shadow Word: Pain stepping from 34 to 38) is
  client 70009's and no longer happens, while dots still tick on current stats.
- **Merge decision, 2026-10-02 (upstream #602 and #607 do the same job): kept ours.** Upstream's `applyEureka` mirrors
  the same three masks in core and tests each sim spell's client row (`ClientClassFlags`, installed by `spelldata`)
  against them: cost on effect 0, +10% on the hit for effect 1, dots scaled back out for a spell on the direct list
  alone. Its masks equal the ones copied into the class tests here, and for every spell the sim registers the two
  mechanisms agree except three. Hellfire (11684) and Penance (1316995) are channels whose ticks are the hits of the
  client's Hellfire Effect and Penance's bolts: the client names them on the damage list, so ours boosts their ticks and
  upstream's, which scales a direct-list spell's dots back out, gives them nothing. Immolate's dot (its own spell on the
  hit's id, tag 1) nets out to the same in both. No Gnome reference casts Hellfire or Penance (the Smite rotations cast Penance and
  the Smite reference is Undead), so no reference moves between the two; it matters to a Gnome priest or warlock who
  does. `ClientClassFlags` stays in `core` and the store
  (it is upstream's, and nothing here reads it) so the files match upstream's. To retire this patch: take upstream's
  `applyEureka` and add a per-class "dealt as ticks" mask for Hellfire and Penance, then convert the class tests to
  behaviour checks. Upstream's `TestGnomeWarrior` additions are not merged (they assume a core without the class
  lists); the core file is the fork's.
- **Drop it when** never: it follows the rows. Re-check the masks in the tests when the client changes them.

## 63. `spelldata: Insight and Increased Spirit read as Spirit` (narrowed to Mystic Mushroom)

- **Narrowed in the 2026-10-02 sync.** Upstream #606 puts Insight's 70124 row back in SQL (`tools/database/overrides/2.sql`,
  which `gen_spelldata` now runs before it reads the tables), so the store has Insight's buff as Spirit again and
  `percentStatRowsNamingNoStat` lost its 1299796 entry (the test said so: the client names a stat now). Mystic Mushroom's
  Increased Spirit (1248751) has no override upstream, so its entry stays and the fork's parser still reads it as
  Spirit. `TestInsightMultipliesSpirit` passes on upstream's row; `TestPercentStatRowsNamingNoStatReadAsSpirit` now
  covers the Mushroom only. The original reading follows.
- **What it did.** Client 70170 moved Insight's buff 1299796 (Enchant Weapon - Insight, "Increases your Spirit by
  100%" for 10 s) and Mystic Mushroom's Increased Spirit 1248751 ("Increases Spirit by 5%") from
  `A_MOD_PERCENT_STAT` on Spirit (misc 4) to `A_MOD_TOTAL_STAT_PERCENTAGE` with MiscValue_0 0 and no stat mask in
  MiscValue_1. Every other row of that aura names its stat in one of the two (Spirit Tap: misc 4 and mask 16;
  Arcane Mind: misc 0 and mask 8), and the parser read these two as Strength. Both tooltips still say Spirit, and
  Coward! (422978) lost its all-stats -1 in the same way, which looks like a data error rather than a change of
  stat, so `percentStatRow` (`sim/core/spelldata/parse_effects_table.go`) reads exactly these two rows as Spirit
  while they name no stat. A row that gains a stat, or goes back to `A_MOD_PERCENT_STAT`, reads as it states.
- **Measured.** No reference wields Insight or the Mystic Mushroom (a Warrior's ranged slot item), so no reference
  moves. The Insight enchant on the main hand of four caster references, same seed, 100,000 iterations (standard
  error 0.05 to 0.15), DPS with Insight minus DPS without it on the same engine:

  | Reference | 120 s, Strength read | 120 s, Spirit read | 300 s, Strength read | 300 s, Spirit read |
  |---|---|---|---|---|
  | Shadow Priest | -0.10 | +2.51 (+0.41%) | -0.07 | +0.14 (+0.02%) |
  | Smite Priest | +0.05 | +1.45 (+0.33%) | -0.04 | +3.32 (+0.80%) |
  | Fire Mage | +0.43 | +0.16 (+0.03%) | +0.06 | +2.41 (+0.42%) |
  | Affliction Warlock | +0.04 | +1.88 (+0.35%) | +0.10 | +2.76 (+0.52%) |

  So the Strength read made Insight worth nothing to a caster (the +0.43 is within two standard errors), and the
  Spirit read makes it worth up to 0.8%, which is what the enchant's tooltip sells. The goldens' Insight and Mystic
  Mushroom rows (`TestBalance`, `TestFeralBear`, `TestFeralCat`, `TestSurvivalMelee`, `TestProtection`,
  `TestRetribution`) move back to their 70124 numbers.
- **Alternative reading.** Follow the row: Insight and the Mushroom raise Strength. Nothing in the tooltips,
  enchant text or the earlier client supports it, so it is read as a client error. If a hotfix names the stat,
  the entry stops applying (the parser checks the misc values), and `TestPercentStatRowsNamingNoStatReadAsSpirit`
  fails to say the table can lose the entry.
- **Tests.** `TestInsightMultipliesSpirit` (re-enabled), `TestPercentStatRowsNamingNoStatReadAsSpirit` and
  `TestPercentStatRowsThatNameAStatAreUnchanged` in `sim/common/shared`.
- **Drop it when** the client names Spirit in the two rows.

## Caster rows of client 1.60.1.70170 that needed no engine change

Each was read against the engine; the tests below pin the behaviour so a later hotfix shows up.

- **Shadow Word: Death and Early Demise.** The talent (1310076) states 30 crit at rank 2 on effect 1 and a health
  threshold of 20 on effect 2. `shadow_word_death.go` already adds the crit only in the 20% execute phase, so the
  bug the notes describe is not in the sim. `TestEarlyDemiseOnlyBelowItsHealthThreshold` measures 3.2% crits above
  20% health and 32.8% at or below.
- **Shadow Weaving.** The stack aura 15258 gained Always Hit: the debuff no longer rolls to land. The sim applies the
  stack to the priest on every landed Shadow hit the talent's own chance picks (33, 67 or 100%), with no second
  roll. `TestShadowWeavingAppliesOnEveryHitAtRankThree`.
- **Improved Scorch and Winter's Chill.** Fire Vulnerability (22959) and Winter's Chill (12579) gained Always Hit
  the same way. Both stacks are applied after the talent's chance and nothing else
  (`TestImprovedScorchRollsOnlyTheTalentChance`, `TestWintersChillStacksOnEveryLandedFrostHit`).
- **Heating Up.** The client renamed Hot Streak to Heating Up and reworded it ("reduce the cast time of your next
  Pyroblast cast within 20 sec"); no row's numbers changed, and the sim already adds a stack for every non-periodic
  crit, so nothing depends on a streak. The aura label, the log name and the `ui/sim/spells/mage.json` entry
  follow the name. Rotations name 400625 by id, which still resolves.
- **Combustion.** Three charges are read from the row (`ProcCharges` 4 to 3), so the engine already has them:
  Fire Mage 598.67 to 588.84 at 120 s and 582.80 to 574.97 at 300 s (-1.6% and -1.3%, 200,000 iterations);
  Frostfire has no Combustion. `TestCombustionEndsAfterThreeCrits`.
- **Soul Harvest.** Renamed from Soul Harvesting, and its first effect moved from aura 379 to
  `A_MOD_POWER_REGEN_PERCENT` on mana, which is the "now correctly grants" fix. The buff starts on a kill under
  Drain Soul, which no encounter has, so the sim never applies it and has nothing to change.
- **Elemental Focus's Clearcasting mask.** Of the two class bits it lost, word 1 0x40000 was Fire Nova's old bit
  (Fire Nova now sits on word 0 bit 27, beside Fire Nova Totem, and the mask names that instead) and word 3
  0x40000000 was Molten Blast (425339), a Season of Discovery rune with no talent in the sim. The spells the mask
  names are unchanged: Lightning Bolt, Chain Lightning, Lava Burst, the three shocks and Fire Nova.
  `TestClearcastingNamesTheSpellsTheClientRowDoes` (`sim/shaman/elemental`) holds the sim's consumption list to
  the row.
- **No sim effect**, as the notes say: scrolls from comprehension cannot be cast while moving, and pets in
  aggressive mode.

## Client 1.60.1.70170 adopted (2 October 2026)

`assets/db_inputs/forever_client_build.txt` reads 1.60.1.70170, and since the 2026-10-02 sync everything generated is upstream's
real regeneration with the hotfix cache (#608): the item database (`db.json`, `db.bin`, `leftover_db.*`), the spell store
and its inputs, the class spell data, the talent trees, the protos, the enchant and item proc files. Patches 40 to 42, 51,
52, 60, 62 and 63 implement what the client rows do not carry (Shifting Power, Swipe attack power, Spearing Strike's stance,
Inner Focus, Eureka!, Mystic Mushroom); the Rage on crits (50) and Hellfire (61) patches went to upstream.

The alternate "Simple Vaelastrasz" Feral Cat rotation (`ui/specs/druid/feralcat/apls/simple_vael.apl.json`) cast and
refreshed Tiger's Fury (9846), which 70170 removed; those two lines are gone.

## Upstream sync 2026-10-02, #602 to #609

Merged `ccfaacb5c3` into `mythicsim/release-20261002` (7325f41ecb), merge base `696a6c4040`: #602 Gnome Eureka! only affects the
spells the client names, #603 and #604 data, #605 Hellfire's area hits crit, #606 talent trees for client 70170 and the
hotfixes (Shifting Power, Furious Precision, Heating Up, the Warrior tree reshape, Bloodthirst 45%, Dual Wield
Specialization's dummy Rage, Flametongue's attack spells, Insight's override), #607 Eureka! hit against DoT, #608 the
`[DB] Update to 1.60.1.70170 with hotfixes as of 2026-10-02` regeneration, #609 critical auto attacks pay 75% more Rage and
Champion of the Light drops healing, and the changelog bot commits.

**Data.** `db.json` and `leftover_db.json` are a three-way merge (base `696a6c4040`, fork tip, upstream) with the fork's stat
index and planner armor corrections kept; `scripts/forever-merge-db.py` needed two extensions for this range, both applied to a
patched copy: a section only upstream changed is taken whole (`consumables`: 29 rows moved to category 2593 and Major Troll's
Blood Elixir went), and an upstream deletion of a record the fork never touched is applied (items 274978, 276765, 285326,
263411 and 263412 in the leftover file, and five spell icons). A stats map the base lacked merges field by field (the three
Rotmender's items gained real stats upstream and the fork had added their armor). `db.bin` and `leftover_db.bin` come from
`go run ./tools/sync_db_binary`. The store, `spell_store_inputs.json`, the talent trees, protos, enchant and proc files are
upstream's; the store was re-rendered from the committed inputs with the fork's extras and differs from upstream's only by
Tiger's Fury (5217), which `extra_spells.go` no longer pins (patch 41), and the Druid class file only by the Shifting Power
ladder (patch 40). `interimHotfixItemSpells` and the interim regeneration notes are gone: the real overlay roots those items.

**Decisions per patch.**

| Patch | Decision | Why |
|---|---|---|
| 50 Rage on crits | dropped, upstream #609 | One constant in the shared bar, paid by Warrior, Bear and the Cat's Bear Form exactly once; the opt-in option protected nothing. Tests kept, a Cat Bear Form test added. |
| 61 Hellfire crits | dropped, upstream #605 | Same behaviour, read off the row and with the better death check. Our test kept. |
| 62 Eureka! lists | kept | Equal masks, but upstream's mask handling gives Hellfire's and Penance's ticks no bonus; ours does what the client lists. |
| 63 Insight and Increased Spirit | narrowed | #606 restores Insight's row in SQL; Mystic Mushroom has no override. |
| Flametongue ladder pin (`ByID(16344)`) | equal to upstream's | Same line. |
| 40 talent-granted ladder, 41 Shifting Power, 42 Cat APL | kept | Upstream also registers Shifting Power but keeps Tiger's Fury, has no Wolfshead Energy move and no tests; see the table below. |
| 30 to 33, 51, 52, 60 | kept | Unrelated to the range. |

Upstream's Shifting Power (`sim/druid/shifting_power.go`) and the fork's name the same spell: the merged file is the fork's
(reads the cost and cooldown the same way, gives Wolfshead Helm's 20 Energy, drops Tiger's Fury) and the fork's Improved
Shifting Power. Natural Shapeshifter's mask change was identical on both sides.

**Cat rotation.** Upstream's default APL (Tiger's Fury first, the powershift lines, then Shifting Power at 60 Energy or less
after the finishers) against the fork's, on the Cat reference (`050022-55000032121032212051-052`, 10,000 iterations, seed
4242, the worker's request with the rotation swapped; "adjusted" is what `forever-rotation-adjustments.mjs` leaves, with the
powershift lines and Rune of Metamorphosis removed):

| Engine | Rotation | 120 s | 300 s |
|---|---|---|---|
| merged | fork's | 564.29 | 530.65 |
| merged | upstream's, adjusted | 544.45 | 510.43 |
| merged | upstream's, raw | 531.98 | 506.63 |
| upstream master | upstream's, adjusted | 539.81 | 502.56 |
| upstream master | upstream's, raw | 523.82 | 495.95 |
| upstream master | fork's | 542.05 | 507.49 |

The fork's rotation on the fork's engine is 3.6% and 4.0% ahead of the best upstream pairing, so it stays.

**Golden movements** (Average-Default DPS, each explained; upstream's own goldens over the same range are the check):

| Golden | Fork before to after | Upstream over the range | Cause |
|---|---|---|---|
| TestFury | 293.66 to 309.28 (+5.3%) | +12.4% | #606: Furious Precision, Dual Wield Specialization's Rage, Bloodthirst 45%. The fork already had Rage on crits, so only +6.9% of upstream's +12.4% applies: with crit Rage off on both engines the merge moves Fury +6.88% (upstream +6.9%), and crit Rage on top is +5.1% (upstream's #609 alone is +5.1%). |
| TestArms | 253.00 unchanged | +7.3% | All of upstream's is #609, which the fork had. |
| TestProtectionWarrior | DPS and TPS unchanged, DTPS +2.6% | +0.3% (#609) | Toughness (armor) left the tree and the preset spends Iron Will instead: armor 4118.8 to 3872 (upstream's golden has no armor to move). |
| TestRetribution | unchanged DPS, CharacterStats row moved | -4.1% | Champion of the Light 60% of Intellect and no healing (#609); the fork already had the damage cut, so only the healing power stat moved. |
| TestFeralCat, TestFeralBear, TestFire, TestProtection, Balance, Hunter, Mage, Priest, Shaman, Warlock suites | one or two AllItems rows (Searing Dagger 1291568, Plaguefang 1309315, Vile Protector, Mystic Mushroom dropped from the enumeration) | the same rows | #608's item data. Where the fork's gear matches upstream's the new numbers are equal (Balance, Arcane, Frost, Fire Searing Dagger rows). |

Nothing moved without an upstream cause.

**Tests.** `go test --tags=with_db` over every `./sim/...` package except `sim/web`, `./tools/database/...` (green: the
hotfix cache that failed it before is in the data now), `./tools/spelldata`, `./tools/db2tool/...` and `./cmd/wowsimcli/...`
pass. Tests this sync changed: the off-hand hit test (now Furious Precision), the max Rage log test (Gnome's Expansive Mind
is the only source above 100 now), the crit Rage tests (Dual Wield Specialization pays the off hand 1.5 times at 5/5), the
percent-stat test (Mystic Mushroom only), and the talent strings of the Warrior tests that were not upstream's.

## 70. `shaman: Flametongue Totem` (narrowed 2026-10-07, upstream #677 and #682)

Upstream now registers the cast (#677, ported from this patch) and lands its hit with no spell power
coefficient and no talent mask (#682: the hit logs as Flametongue Attack 16368, beta log 2713, whose
class mask Elemental Fury and Elemental Weapons do not name). The fork keeps its own registration
(`sim/core/buffs/flametongue_totem.go` and `registerFlametongueTotemSpell`) because the party flag
(patch 71) shares the trigger and the hit, and takes #682's rules: the shaman's traits now carry the
spell flag only, and `TestFlametongueTotemHitTakesNoElementalFuryOrElementalWeapons` replaces the
talent test. Since patch 98 Windfury Totem switches it off as upstream's does, so the fork and upstream
agree on the rules; the registration and the party flag are what remain ours. The measurements below
predate all three changes; since patch 98 rows (b) and (c) get no totem hits at all.


Redfall (Discord, 2 October 2026) asked for a Windfury Weapon main hand with Windfury Totem and a Flametongue Totem in
place of Searing Totem ("ftt (for yourself, which is a dps increase over searing totem) while keeping wf imbue up"). The
engine had the rank ladders (`spellData.FlametongueTotem` 8227, 8249, 10526, 16387, and the triggered ladder 8230 to
16389) and the proto enum value `FireTotem_FlametongueTotem`, but nothing registered the cast, so a rotation line naming
16387 was dropped ("does not know spell") and the fire slot only held Searing Totem, Magma Totem and Fire Nova.

**What the client says** (build 1.60.1.70170, wago rows and Wowhead's Forever tooltip):

- Rank 4 (16387): 275 mana, a 1 sec global (SpellCooldowns StartRecoveryTime 1000), 5 min (DurationIndex 5, 300000 ms),
  summons a totem with 5 health. Level 58.
- The totem's aura is 15036, an area party aura (effect 35, 30 yards) with aura 42 (proc trigger spell): ProcChance 100,
  ProcTypeMask 4 (a landed melee auto attack), triggering 16389. It names no weapon in its row; the aura text reads "Main
  hand attacks deal an additional (1363 / 77) to (1363 / 25) Fire damage" and the cast text "Each main hand hit causes
  (1363 / 77 * mult - 1) to (1363 / 25 * mult) additional Fire damage, based on the speed of the weapon".
- 16389, "Flametongue Totem Proc": one dummy effect, `EffectBasePointsF` 1363, no per-level gain, school fire, class mask
  bit 34. `mult` is Improved Weapon Totems (29192, 29193: +6% and +12% on that mask), which Forever's talent trees do not
  carry, so it is 1.
- The same shape as Flametongue Weapon's proc (16344: 2810 at 60 from 2498 + 78 a level). Both are hundredths of damage
  per second of weapon speed with the speed held to 1.3 to 4.0 (the tooltip's "/ 77" and "/ 25"): the totem is 13.63 a
  second of speed, 17.719 at the floor, 54.52 at the cap, 35.438 for a 2.6 speed mace, 51.794 for Arcanite Reaper (3.8).

**What it does.**

- `Shaman.registerFlametongueTotemSpell` casts 16387 like Searing Totem: the rank's cost and global, 300 s, the fire slot
  (`cancelFireTotems` now ends it, and its cast ends Searing Totem and Magma Totem, and `TotemExpirations[FireTotem]`
  follows it). It joins `SpellMaskTotem`, so Totemic Focus's cost cut reaches it. A rotation line naming 16387 works.
- The shaman's totem aura ("Flametongue Totem (Self)", 300 s) switches on a trigger, "Flametongue Totem Trigger", built
  from the party aura's own row (`spelldata.ProcTrigger` on 15036: 100%, melee autos, damage dealt) narrowed to the main
  hand (`ProcMaskMeleeMHAuto`). Windfury's extra attacks are main-hand autos, so they add the hit too; off-hand swings,
  specials, ranged and spell hits do not. Everything lives in `sim/core/buffs/flametongue_totem.go` so the party flag
  (patch 71) and the cast share one trigger and one hit.
- The hit keeps the totem's id (16389) so a report lists it apart from the imbue. It is the imbue's spell with the totem's
  base damage: Magic fire (spell hit and crit tables, partial resists, 1.5 on a crit), 13.63 a second of weapon speed,
  and a 0.1 spell power coefficient. A weapon with no speed adds nothing, as the imbue's.
- **An inference, not a row.** Neither dummy deals damage. The only client fire spells of the family with a hit table are
  the three "Flametongue Attack" spells (10444, 29469, 29470: Magic in SpellCategories, a 0.1 `EffectBonusCoefficient`,
  class mask bit 21), which the imbue's dummies feed; 16389 has no damage row, no SpellCategories row and its own class
  mask. The earlier note in `applyElementalFury` ("the same attack granted by Flametongue Totem crits for 1.5x on other
  players") says the totem's hit is that attack. The engine therefore gives a shaman's totem hit the imbue's class mask
  and spell flag (`buffs.SetFlametongueAttackTraits`, set in `weapon_imbues.go`), so Elemental Fury (2.0 on a crit),
  Elemental Weapons (+5% a point) and Natural Grace's threat cut reach it as they reach the imbue's hit, and the coefficient
  is read from 10444. Another class's totem hit has the coefficient and no talents. If the beta shows the hit as its own
  spell with no coefficient, set `ClassSpellMask`, `Flags` traits and `BonusCoefficient` aside in `FlametongueTotemAttack`:
  the numbers in the table below say what that is worth (the "no coefficient" rows).
- **Main-hand autos only, as the row states.** ProcTypeMask 4 is melee autos; the Windfury Totem aura beside it states 0x14
  (autos and specials), so the difference looks deliberate. If the totem also hears main-hand specials (Stormstrike) the
  hit count rises by about a fifth on the Enhancement preset: see the table.

**Tests.** `sim/shaman/fire_totems_test.go` (`TestFlametongueTotemRank4`, `TestFlametongueTotemBaseDamage`: the rank's
cost, global and duration, the 1363 value and the exact damage at 1.0 to 4.5 speed) and
`sim/shaman/enhancement/flametongue_totem_test.go`: `TestFlametongueTotemHitDamage` (a hit is exactly 35.438 for a 2.6
speed mace, 51.794 for a 3.8 speed axe and 20.445 for a 1.5 speed dagger, and a crit 1.5 times that, read off clean hits
of a 40 iteration run), `TestFlametongueTotemHitScalesWithSpellDamage` (a +35 elixir adds exactly 3.5 to a hit),
`TestFlametongueTotemHitTakesElementalFuryAndElementalWeapons` (2.0 on a crit at 5/5, +15% at 3/3),
`TestFlametongueTotemHitsOnlyMainHandAutoAttacks` (one hit per landed main-hand swing, none for the more numerous
off-hand ones), `TestFlametongueTotemNeedsTheCast`, `TestFlametongueTotemReplacesSearingTotem` (both orders) and, since
patch 98, `TestWindfuryTotemSwitchesFlametongueTotemOff` and `TestGraceOfAirLeavesFlametongueTotemAlone`. None of them can build or pass without the cast and trigger.

**Default.** Nothing sets the cast or the flag, so every request that does not name 16387 is bit-identical: the Enhancement
reference's one-click request gives 619.6935 (120 s) and 596.4838 (300 s) before and after, and every suite golden passes
unchanged.

**Measured** on the Enhancement reference (Redfall's talents, 10,000 iterations, board seed 1179746067, fixed 120 s and
300 s, DPS / TPS), each row replacing only what it names; the preset is Rockbiter Weapon, the party's Windfury Totem and
Searing Totem:

| Variant | 120 s DPS | 300 s DPS | 120 s TPS | 300 s TPS |
| --- | ---: | ---: | ---: | ---: |
| (a) preset | 619.69 | 596.48 | 644.28 | 623.00 |
| (b) Rockbiter + Windfury Totem + Flametongue Totem | **628.48 (+1.42%)** | **608.18 (+1.96%)** | 653.05 | 634.79 |
| (c) Windfury Weapon + Windfury Totem + Flametongue Totem | 586.70 (-5.33%) | 567.65 (-4.83%) | 612.41 | 595.36 |
| (d) Windfury Weapon + Grace of Air + Flametongue Totem | 609.10 (-1.71%) | 589.34 (-1.20%) | 634.94 | 617.17 |
| (e) Flametongue Weapon + Flametongue Totem | 560.73 (-9.51%) | 541.82 (-9.17%) | 585.31 | 568.44 |
| (f) no fire totem | 598.23 (-3.46%) | 576.96 (-3.28%) | 622.88 | 603.56 |

Controls: Windfury Weapon + Windfury Totem + Searing 572.29 / 551.82; Windfury Weapon + Grace of Air + Searing 594.46 /
573.34; Flametongue Weapon + Searing 584.70 / 562.23; Flametongue Weapon alone 563.41 / 542.59 ((e) is 2.7 DPS under it
because the totem adds no hit beside the imbue and its 275 mana costs half a Fire Nova a fight); Rockbiter + Grace of Air
with Flametongue Totem 563.79 / 543.79 and with Searing 561.21 / 537.88. Against Searing Totem in the same build the
Flametongue Totem is +1.4% / +2.0% with Rockbiter + Windfury Totem, +2.5% / +2.8% with Windfury Weapon + Grace of Air,
+2.5% / +2.9% with Windfury Weapon + Windfury Totem and +0.5% / +1.1% with Rockbiter + Grace of Air, and -4.1% / -3.6% with
Flametongue Weapon (the totem is switched off and Searing Totem is gone). Its hit is 91 on average (51.8 base, 1.15 from
Elemental Weapons, 0.1 of about 200 spell damage that Mental Quickness gives, crits) against Searing Totem's 60.

How much the answer depends on the two inferences above, build (b) against the preset (120 s / 300 s):

| Model of the totem's hit | (b) DPS | vs preset |
| --- | ---: | ---: |
| Shipped: the imbue's spell (0.1 coefficient, talents), main-hand autos | 628.48 / 608.18 | +1.42% / +1.96% |
| Same, and main-hand specials too | 637.39 / 617.36 | +2.85% / +3.50% |
| No coefficient and no talents, main-hand autos | 614.90 / 594.89 | -0.77% / -0.27% |

**Drop it when** upstream registers Flametongue Totem. Its hit may then be a different spell; keep the exclusivity (patch 72).

## 71. `core: PartyBuffs.FlametongueTotem`

A party buff flag (`flametongue_totem`, field 20, JSON `flametongueTotem`) so a melee spec can assume another shaman's
Flametongue Totem, as `windfuryTotem` does. It is a row of the buff manifest (`tools/database/buffmanifest/buffs.go`, a
manual driver on 15036 in the "FlametongueTotem" category), so `proto/buffs.proto`, `sim/core/buffs/buffs_auto_gen.go` and
`ui/features/settings/model/buffs_debuffs_auto_gen.ts` come from `go run ./tools/gen_buffs_proto` and the buff render.
`driveFlametongueTotem` keeps a permanent "Flametongue Totem" aura that switches the shared trigger on. A Windfury
Totem switches it off (patch 98); Grace of Air does not. The party's totem and the shaman's own cast are the same
effect, so both together add one hit per swing. No spec's defaults set it: every default
request is unchanged.

**Tests.** `TestPartyAndCastFlametongueTotemAddOneHit` and `TestWindfuryTotemSwitchesFlametongueTotemOff` (all four
pairings of cast and party totems, patch 98) in the Enhancement package, and `TestPartyFlametongueTotemHitsOnMainHandAutoAttacksOnly`
in `sim/warrior/dps` (a dual-wielding Fury Warrior: a main-hand auto adds the hit, an off-hand auto, a main-hand or
off-hand special, a ranged auto and a spell add none).

**Drop it when** upstream carries a Flametongue Totem party flag of its own.

## 72. `shaman: a main-hand Flametongue Weapon disables the personal Flametongue Totem`

The beta's Flametongue Weapon tooltip: "When applied to main hand, disables any benefit you personally receive from
Flametongue Totem" (and Windfury Weapon's, "...Windfury Totem": patches 20 and 30). `RegisterFlametongueImbue` puts the
main-hand imbue's trigger aura in the "FlametongueTotem" exclusive category at twice the totem's bid
(`buffs.DisableFlametongueTotem`), so with a main-hand Flametongue Weapon the totem stands (the aura stays up, the mana is
spent) and adds no hit, from the shaman's own cast and from the party flag; the totem takes over again when the imbue leaves
the main hand. An off-hand Flametongue Weapon, Windfury Weapon, Frostbrand Weapon and Rockbiter Weapon leave it alone, since
the client says this only of Flametongue. A main-hand Windfury Weapon turns off only Windfury Totem's benefit, a
Flametongue Weapon only Flametongue Totem's. Windfury Totem itself switches Flametongue Totem off since patch 98.

**Tests.** `TestOnlyAMainHandFlametongueWeaponDisablesFlametongueTotem`: six imbue setups by cast and party totem, the
imbue's own hit counted beside; it fails if the exclusive effect is removed.

**Drop it when** upstream models the two weapon texts the same way.

## Upstream sync 2026-10-02, #610 to #612

Merged ElliotWood/Forever `f4b776b4f41d5c7799b8141697a2c9e67c89d426` on top of release `0ebf4100ae`.
Arcane Power and Arcane Instability now boost Frostfire Bolt's direct hit alone. Nature's Grace
also shortens the affected spells' global cooldown. The item database includes seven changed
items and four additions, merged per field while retaining the fork's stat indices and planner armor.
Only the Balance golden changes, matching upstream's Nature's Grace movement.

## 73. `cli: statweights`

The application worker requests character stat weights through `statweights --strict --infile ...
--outfile ...`, but the CLI only exposed raid simulations. The missing command caused optional
weights to fail while the DPS simulation completed. The new command calls the existing
`core.StatWeights` API with a `StatWeightsRequest` and writes its `StatWeightsResult`. Strict mode
rejects unknown fields and enum names. Incomplete requests fail before calculation; structured
engine errors are saved and also return a nonzero exit status.

Tests: `TestStatWeightsCommand` calculates a real Mage spell damage weight and checks positive,
finite weight and reference EP of one. `TestStatWeightsRequestValidation` checks strict decoding
and missing fields, iterations and stats.

Drop this adapter when upstream exposes the equivalent CLI command.

## 74. `core: concurrent sims honour GOMAXPROCS`

The concurrent runner used `runtime.NumCPU()` even when a caller bounded Go concurrency.
A tank stat-weights calculation therefore allocated twelve simulations inside the application's
one GiB container and was killed for memory use. The runner now caps the split count at
`min(runtime.NumCPU(), runtime.GOMAXPROCS(0))`. The stat-weights worker supplies `GOMAXPROCS=2`
and `GOMEMLIMIT=512MiB`. Iterations, seeds, encounter and weight calculations stay the same.
Test-mode simulations keep their existing three splits. The normal default remains the host's
available Go concurrency. The real caster and tank image checks cover the worker budget.

Drop this when upstream respects the caller's concurrency limit.

## 75. `core: measure applied exclusive effect uptime`

Aura uptime includes time an aura remains active while another effect in its category wins.
Each exclusive effect now measures only its applied intervals after the pull. Aura metrics
export these intervals by category, including a measured zero when fully suppressed.
Concurrent results combine the averages with the same iteration weights as aura uptime.
Selection, callbacks and damage calculations are unchanged.

Tests cover equal-strength suppression, resuming after the winner expires, iteration resets,
prepull time, lazy expiration and unequal concurrent result weights.

Drop this when upstream exports equivalent per-effect uptime metrics.

## Upstream sync 2026-10-03, #613 to #641

Merged ElliotWood/Forever `f764984d8b` (55 commits, 27 of them data or changelog) into the
patch 75 release `f43215d07e`. No patch is dropped. This supersedes the earlier attempt on
`codex/forever-upstream-oct03` (#613 to #624 onto patch 74), whose two resolutions are reused.

Behaviour adopted: Whirlwind strikes with both weapons, and Raging Blows also cuts its Rage
cost (#613); Booming Voice takes 5% a point off the shouts (#614); Piercing Ice gives Frostfire
Bolt's DoT a flat 2% (#615); an APL aura the character lacks reads as absent (#622) and a strict
sequence with an unknown step is dropped (#625); Mutilate's hand strikes stop re-rolling miss,
dodge and parry (#632); Slam is cast when Improved Slam makes it fast (#629); client
1.60.1.70205 item data. Optional rotations: destruction_conflag, forever_flameshock, dps_battle,
dps_dance, Ghostly Strike rogue variants, Smite with Mind Blast and Shadow Word: Death; Auto
picks the warrior stance rotations (#641).

Resolutions:

- `sim/warrior/talents_arms.go`: both sides made Spearing Strike require Battle Stance (#626);
  kept ours, which also records the 70170 weapon change.
- `simple_vael.apl.json`: took upstream's guarded Shifting Power (1322605). Both sides removed
  Tiger's Fury.
- `db.json`: `scripts/forever-merge-db.py` with the merged planner. Mantle of Woe (7750) takes
  upstream's quality 3 stats with the planner's 38 armor and 50 bonus armor; planner armor was
  re-applied to 54 items, including Vile Protector (7747, 1051 to 1078 armor).
- `frostfire_bolt_test.go`: checks the hit (1.14) and DoT (1.10) multipliers separately,
  since #615 moved Piercing Ice's DoT share off the additive multiplier.

Goldens regenerated for Marksmanship, Survival, Assassination, Affliction, Destruction, Arms,
Fury and both Protection specs. Each moved the same way as upstream's own goldens over the
range (Survival median +11.6% here, +10.8% upstream; Assassination +3.9% both; Fury -0.3% both).
Arms (+6.0% vs +7.4%) and Destruction (+3.6% vs +5.3%) move about 1 to 1.5% less on average
because our Warrior and Warlock patches change the baseline; every row moves the same
direction except naked Destruction, where our default rotation keeps Conflagrate. Protection
Paladin and Warrior move only on the Vile Protector item row (armor). Whole suite green.

## 84. `druid: a Dense stone counts in cat and bear form`

Forever's forms rebuild the paw from the equipped weapon (`formWeapon`), rescaling its damage range
to the form's 1.0 or 2.5 second swing. Flat weapon damage is not rescaled: it lands in full on
every paw hit. Hameru's beta character sheet (MythicSim Discord, 5 October 2026) shows Heavyhammer
(73 to 110 at 3.3 s) in cat form at 46 to 58, and 49 to 61 with a +3 weightstone, so the stone
adds +3 per paw hit rather than 3 x 1.0 / 3.3. Hameru also reports a +5 weapon-damage enchant adds 5.
The client stores both as the same enchantment kind (SpellItemEnchantment effect 2: "Weighted +8",
"Sharpened +8" and "Weapon Damage +9" alike), and the server adds that kind to the hand's damage as
one flat amount, so stones and weapon-damage enchants are treated the same.

Before this, `registerStaticImbue` added a Dense Sharpening Stone's or Dense Weightstone's +8 to the
humanoid weapon only, which the form replaces, so the stone did nothing in form; and
`newWeaponFromItem` folded `Enchant.WeaponDamage` into the range that `formWeapon` rescales, so +9
Superior Impact gave a cat on a 2.9 second staff about +3.1 a swing. `formWeapon` now rescales the
weapon's own range and bonus DPS, then adds the enchant's and the stone's flat damage. Crusader,
+25 Agility and the Elemental Sharpening Stone are not weapon damage and are unchanged.

MythicSim's references, board seed, 10,000 iterations, empty trinkets, default length (300 s in
brackets):

| Main hand | Before | After |
| --- | ---: | ---: |
| Cat, Crusader, Elemental stone | 599.92 (563.43) | 599.92 (563.43) |
| Cat, +9 Superior Impact, Elemental stone | 594.51 (557.02) | 606.76 (568.53) |
| Cat, Crusader, Dense Weightstone | 590.16 (553.70) | 606.50 (569.07) |
| Cat, +9 Superior Impact, Dense Weightstone | 584.75 (547.43) | 613.15 (574.12) |
| Bear, Crusader, Elemental stone | 432.80 (435.85) | 432.80 (435.85) |
| Bear, +9 Superior Impact, Elemental stone | 425.74 (427.82) | 426.36 (428.43) |
| Bear, Crusader, Dense Weightstone | 422.68 (425.57) | 429.87 (432.73) |

In cat form +9 weapon damage now beats Crusader and a Dense stone beats the Elemental one; the
bear's 2.5 second swing keeps Crusader and the Elemental stone ahead.

`TestFormPawCarriesFlatWeaponDamage` checks both Dense stones, Superior Impact and the two together
in both forms, and that Crusader and the Elemental Sharpening Stone add no weapon damage. Feral
goldens are unchanged: the suite carries no weapon-damage enchant or Dense stone.

Drop this when upstream adds flat weapon damage to the form weapon after the rescale.

## 76. `hunter: Hawk follow-up attacks roll melee avoidance` (dropped, upstream #703)

Dropped in the 2026-10-07 sync: upstream's #703 fits the hawk to 534 beta swings (every 2.5 sec
hasted, about 0.35 of the dive bomb base, `OutcomeMeleeSpecialHit`: misses, dodges and parries, no
crits) and its `TestSummonHawkTwoHawks` replaces ours.


Forever Logs report 2701 records Trapz's Hawk auto-attacks missing, being dodged,
being parried and critically hitting Saltspine. The previous follow-up outcome
was `OutcomeTickPhysicalCrit`, which could only hit or crit. Follow-up attacks
now use the existing melee special hit/crit table, without the Hunter's
dual-wield miss penalty or an assumed guardian glancing multiplier. The initial
dive remains always-hit from the client attribute.

Sources: [damage breakdown](https://foreverlogs.gg/reports/2701/encounters/damage-done?source=30826)
and [attack events](https://foreverlogs.gg/reports/2701/encounters/damage-done?source=30826&spells=-1&view=events).
The first event page records a dodge at 00:08.172, an 18-damage crit at 00:09.062,
a miss at 00:15.932 and a parry at 00:34.782, all against Saltspine.

This is an incremental correction to the existing scheduled attack model.
Guardian damage, three-second attack interval, owner-based hit/expertise/crit
and position remain approximations, not fitted results. The displayed
low-level damage cannot establish a level-60 formula. The merged Hawk names
also cannot establish an individual guardian's base interval without GUIDs and
attack-speed buffs. Do not calibrate damage to the displayed average of seven.

`TestSummonHawkTwoHawks` checks both positions: two active Hawks attack and crit,
follow-ups can miss and be dodged, parries occur only in front, and the dive
never misses, dodges, parries or blocks.

## Upstream sync 2026-10-05, #642 to #676

Merged ElliotWood/Forever `67f14b04a5` (69 commits, 35 of them data, changelog or leaderboard) into
the 2026-10-03 sync `e5c9ba5ea3` together with patch 84 (flat stone and weapon-damage enchant damage on every paw hit, above) and patch
76 (Hawk follow-up outcomes). This supersedes the unreleased `e5c9ba5ea3` pin.

Behaviour adopted: Arcane Concentration skips triggered spells and Blizzard rolls on its cast
(#643); Windfury Weapon is two 439440 special hits with the rank's attack power and no swing reset
(#644, beta log 2708); Deep Wounds ignores caster damage modifiers (#645); Ignite, Elemental
Devastation and Flurry read Can Proc From Procs from the client (#646); Blizzard ticks no longer
stack Winter's Chill (#647); the Demonic Brand hit cannot miss and can crit (#648); trap burns,
Scorpid Poison and Volley ticks can crit (#650); Ravage cannot be dodged (#651); Life Tap, Dark
Sacrifice and Charge make no threat (#652, #653); Hack and Slash, Furor and Improved Sayaad read the
effects their tooltips name (#654, #655: the Succubus' Lash of Pain bonus read Seduction's
duration); Presence of Mind and Combustion share a cooldown (#656); traps read their own cooldown
categories (#657); per-tree mage launch sets (#658); Lightning Overload reads its own rows (#659);
the mage, rogue and warrior audits (#660, #665, #667: Blizzard and Flamestrike ticks crit, Thistle
Tea restores a flat 100, Hemorrhage's Rupture bonus is a damage-taken effect read per tick); Arcane
Missiles spends Arcane Blast stacks as it starts (#661); Missile Barrage rolls on a landed hit.

Patches upstream now carries, dropped here:

| Patch | Upstream | Notes |
|---|---|---|
| 13 Windfury Totem in Cat and Bear | #669 | Same code and test. #676 also shows the checkbox. |
| 18 stone or oil under Windfury Totem | #665 | Same rule. |
| 19 Rockbiter Weapon | #673 | Upstream's also applies Spirit Weapons' threat half. |
| 20 Flametongue keeps Windfury Totem | #674 | Same; our Flametongue Totem disable (patch 72) stays on top. Upstream's `TestOnlyWindfuryWeaponDisplacesWindfuryTotem` replaces ours. |
| 22 dots read stats at the tick | #668 | Same `dot_test.go`. Our unset-multiplier guard (patch 12) and `sim/dot_rules_test.go` stay; Deep Wounds is now pinned to ignore a mid-dot multiplier (#645). |
| 23 Rain of Fire | #672 | Upstream's adds the tick proc flag and the cast's dummy on every enemy. |
| 27 form weapon DPS | #670 | Our zero-speed guard (patch 33) and patch 84 stay in `formWeapon`. |
| 30 one air totem | #675 | Upstream joins the same category with the same priorities; `sim/core/buffs/air_totem.go` is gone. Our extra cases moved to `air_totem_cases_test.go` and count Windfury Weapon procs from 439440 hits. Upstream's uptime clamp is equivalent to ours. |
| 31 Windfury Weapon row | superseded by #644 | The procs are now their own 439440 spell, so the white-swing row and its test are gone. |

Kept against upstream:

- **16 Faerie Fire and Curse of Recklessness**: upstream (#671) puts both in one single-aura
  category. Ours shares the category per stat, so only the armor competes and the curse keeps its
  other effects; `armor_reduction_test.go` is ours and checks that.
- **Demonic Brand**: our target-scoped charges stay; the hit takes #648's `OutcomeMagicCrit`.
  `TestDemonicBrandPetHits` now expects six landed hits (hits or crits) and a crit at 100% pet spell
  crit. Upstream's pet-aura `TestDemonicBrandImpSpendsWithFirebolt` does not apply to the
  target-scoped brand and is not taken.
- **Rupture**: upstream's (#665) replaces our snapshot-time Hemorrhage multiplier; keeping both
  would count it twice.

`db.json`: `scripts/forever-merge-db.py` from base `f764984d8b`, 12 new upstream items, planner
armor re-applied to 9. Emissary Cuffs (9455) and Ebony Boneclub (10571) were field conflicts where
both sides moved the same stat index; upstream's refreshed value (6, with its new item level) is
taken.

Goldens: every movement is upstream's for a commit new to the fork. Arcane -0.28% (#643 -1.17%,
#660 +0.90%), Fire +5.60%, Frost +0.61%, Arms -1.75% (#645, #667), Fury -0.08%, Assassination
-1.07% (#665 -1.02%), Subtlety +1.49% (#665 +1.31%; the rest is the per-tick Hemorrhage against
our old snapshot one), Enhancement weapon rows +2.4 to +4.2% (#644, within 0.1% of upstream's on
each row), Survival Melee +1.05% (trap crits and categories). Feral Cat moves +0.57% (Short naked
+1.4%), all from Ravage no longer being dodged (#651): with that one flag removed the golden is
unchanged; our rotation opens with Ravage, upstream's barely uses it. Specs whose upstream goldens
moved for patches the fork already had (Feral -22%, hunters -5 to -6% from #670 and #671) do not
move here. Whole suite green.

## 85. `priest: Dark Sacrifice is Undead only`

Dark Sacrifice (1277324 to 1277328) is the undead priest's race ability: the client's
SkillLineAbility rows for its five ranks carry race mask 16 (Undead), as Starshards' carry 8 (Night
Elf). The fork registered it for every priest and added it as a major cooldown, so Autocast Other
Cooldowns and APL lines cast it for every race, behind `spellIsKnown` too (Zwuggel and Kerani on the
MythicSim Discord, 5 October 2026). `Initialize` now registers it like Starshards: Undead, racials
on, so a race comparison with racials disabled loses it as well (patch 2).

`TestDarkSacrificeIsUndeadOnly` checks the spell and its major cooldown for Undead with and without
racials and for four other races. TestShadowPriest and TestSmitePriest run Troll and Night Elf
priests, whose goldens drop without the spell (Average-Default: Shadow 347.12 to 342.89, Smite 174.58
to 155.68); Undead is unchanged, and so are MythicSim's Shadow and Smite references, which are
Undead.

Drop this when upstream gates Dark Sacrifice on race.

## 86. `items: Dragon's Call whelp honours its 45 sec cooldown`

Dragon's Call (10847) summons an Emerald Dragon Whelp with spell 13049, which carries a 45 sec
category cooldown in the Forever client (SpellCooldowns 54877, category 23), retuned from Classic
Era's 60 sec. The port rolled 1 PPM with no cooldown, so a proc while the 15 sec whelp was out
refreshed it: about 50% uptime and 11.4 Acid Spits a fight on MythicSim's melee Survival Hunter
reference (reported by Bae on the MythicSim Discord). The proc now waits out the cooldown
(`WeaponProcTrigger.ICD`, new), and the summon's cooldown and duration are read from the 13049 row:
25% uptime, at most one 15 sec whelp per 45 sec, 5.3 spits a fight.

Measured on `20b551c6b` (10,000 iterations, seed 42), the Survival reference falls from 760.9 to
735.3 DPS, and Dragon's Call drops from first to fourth of the main hand swords, behind Teebu's
Blazing Longsword, Dal'Rend's Sacred Charge and Krol Blade. Goldens here: TestSurvivalMelee
Dragon's Call 491.70 to 483.90, TestArms Dragon's Call 328.73 to 325.12; nothing else moves.

`TestDragonsCallWhelpProcCooldown` checks that summons are 45 sec apart and each whelp lives
exactly 15 sec.

Drop this when upstream gives the Dragon's Call proc its spell's cooldown.

## 87. `hunter: Lacerating Strikes keeps a refreshed bleed`

`procLaceratingStrikes` wrote the bleed (40% of the Mongoose Bite over 7 ticks, multiplier 1) and
then cast it. `Dot.Apply` deactivates a running dot before it snapshots, and the expiry zeroes
`SnapshotBaseDamage` and `SnapshotAttackerMultiplier`, so a bite that landed with the last bleed
still up (Mongoose Bite comes round about every 8 sec, the bleed lasts 21) left a bleed that ticked
for 0 for its whole duration. Only the first bleed of a fight did damage, until the next bite. A
first-iteration log of MythicSim's survival-hunter reference had 35 of 37 ticks at 0. The bleed is
now written after the cast; a new bite still replaces the running bleed rather than rolling it over,
as before. Ignite and Deep Wounds already write theirs after the cast.

The bleed also reports under its own id, 1310536 (Lacerating Strikes), instead of Mongoose Bite's
with tag 1, so results stop listing a second Mongoose Bite.

MythicSim's survival-hunter reference, 10,000 iterations, 120 s with 15 s variation:

| Seed | Before | After | Lacerating Strikes |
| --- | ---: | ---: | ---: |
| 1179722310 | 768.17 | 785.58 | 1.9 to 19.3 DPS |
| 1180722310 | 766.94 | 784.35 | 1.9 to 19.3 DPS |

Rolling the unticked damage into the new bleed (as Ignite does) would give 812 (+5.9%). The talent
needs 31 Survival points, out of reach at the level 30 beta cap, so no log can settle it yet.

`TestLaceratingStrikesRefreshKeepsTheNewBleed` lands two bites back to back and checks the second
bleed carries its own 40%; it fails on the old order. `TestSurvivalMelee` goldens move.

Drop this when upstream writes the bleed after the cast.

## 88. `racials: weapon specializations leave the ranged auto attack alone`

Human Sword, Orc Axe and Dwarf Mace Specialization added their crit as global physical and spell
crit, so a Human hunter with a sword shot Auto Shot with 2% more crit (an Orc with an axe 1%). In the
game they do not touch the ranged auto attack: on forever-bugs #91 an Orc with an axe, 2% melee and
1% ranged crit on the sheet, took 0 crits from over 500 Auto Shots against a level 21 target, where
the racial would have left 1% after the level suppression, and the character select text names
physical abilities. Melee swings keep it (the same test had 1% melee crit over 101 swings), and so do
ranged abilities and spells. The stat buffs stay global; the aura now also carries a
`SpellMod_BonusCrit_Percent` of minus the bonus on `ProcMaskRangedAuto`, the mask Mortal Shots and
Ranged Weapon Specialization already use for Auto Shot. A warrior's or rogue's Shoot is the same
ranged auto attack and loses it too. The character sheet's ranged crit still shows the bonus.

`TestWeaponRacialsSkipTheRangedAutoAttack` compares Human (sword), Orc (axe) and Dwarf (mace)
hunters with and without their weapon racial: Auto Shot crit does not move, Aimed Shot, Multi-Shot
and the melee swing move by 2%, 1% and 1%; a Night Elf moves nowhere. It fails on the previous
commit. Goldens: the hunter suites run an Orc with Arcanite Reaper (Average-Default: Beast Mastery
424.38 to 423.02, Marksmanship 312.53 to 311.01, Survival 324.29 to 322.65); melee Survival does not
move.

MythicSim references, 10,000 iterations, empty trinkets, board seed (s1) and board seed + 1,000,000
(s2), 120 s with 15 s variation and 300 s without:

| Reference | 120 s s1 | 120 s s2 | 300 s s1 | 300 s s2 |
| --- | ---: | ---: | ---: | ---: |
| hunter (Beast Mastery, Human, Barbarous Blade) | 599.59 to 594.45 | 599.86 to 594.79 | 592.58 to 587.57 | 592.41 to 587.43 |
| marksmanship-hunter (Human, Barbarous Blade) | 675.52 to 669.10 | 675.29 to 668.76 | 618.05 to 611.55 | 617.60 to 611.17 |

No other reference moves; survival-hunter (Orc, axes) fights in melee. On the race board (seed
1179603525, 120 s) Night Elf passes Human for both ranged builds: Beast Mastery Human 599.85 to
594.87, Orc with an axe 597.11 to 594.61, Night Elf 598.20; Marksmanship Human 675.83 to 669.29, Orc
with an axe 670.62 to 667.37, Night Elf 672.40. Dwarf hunters cannot wield a mace and Troll and
Night Elf have no weapon racial, so they do not move.

Drop this when upstream keeps the weapon racials off the ranged auto attack.

## 89. `warrior: a refreshed Deep Wounds keeps its tick timer`

Zirene's answer on forever-bugs #234 (6 October 2026): Deep Wounds scales with weapon damage and not
attack power, rolls over its damage when refreshed, and does not reset its tick timer when refreshed.
Upstream (#506) already had the first two. A refresh deactivated the bleed and applied it again, so
every crit on a running bleed restarted its 3 sec tick phase and pushed the owed damage further out.

`Dot.ApplyKeepingTickTimer` (new, `sim/core/dot.go`) applies a dot without touching a running one's
pending tick: the duration starts over from now, the next tick lands when it was due, and
`RemainingTicks` becomes the ticks that fit between that tick and the new expiry (4 for Deep Wounds,
or 5 when a tick is due the same instant). Deep Wounds reads `OutstandingDmg` before the call and
spreads it with the new crit's share over `RemainingTicks` after it. A first application is a plain
`Apply`. Whether the server also carries the partial tick into the new duration is not known; this
keeps the 12 sec duration.

The Impale connection Zirene mentions ships in the next client build, so it is not here.

`TestDeepWoundsRefreshKeepsItsTickTimer` refreshes a 3 point bleed at 4.5 sec: the next tick stays
at 6 sec, the bleed runs out at 16.5 sec with 4 ticks of (3 x the first tick + the new share) / 4,
and it ticks at 3, 6, 9, 12 and 15 sec. It fails on the previous commit. Goldens (Average-Default):
TestArms 262.95 to 264.65, TestFury 310.72 to 313.02.

MythicSim references, same runs as patch 88:

| Reference | 120 s s1 | 120 s s2 | 300 s s1 | 300 s s2 |
| --- | ---: | ---: | ---: | ---: |
| warrior (Human, Fury) | 843.32 to 851.99 | 842.40 to 851.07 | 826.96 to 831.49 | 827.17 to 831.71 |
| arms-warrior (Human) | 693.33 to 694.81 | 691.05 to 692.52 | 684.16 to 684.53 | 683.51 to 683.87 |

The gain is the bleed no longer pushed past the end of the fight, so it is about a constant amount
of damage and shrinks with fight length. protection-warrior and fury-protection-warrior take no
Deep Wounds and do not move.

Drop this when upstream keeps the Deep Wounds tick timer on a refresh.

## 90. `rotations: hunters stop timing shots around Auto Shot`

The Beast Mastery, Marksmanship and ranged Survival rotations cast Aimed Shot only with the next Auto
Shot over 1 sec away, Multi-Shot and Arcane Shot over 0.5 sec, Summon Hawk and Sniper Shot over 1 sec,
and Rapid Fire just before an Auto Shot. That avoided clipping the shot, but Forever's Auto Shot
cannot be clipped: casts never hold it, only moving and a melee swing do (`swing` in
`sim/core/attack.go`, upstream since e4fd251171), and a hunter at range never swings. Sanctum made
the point on the MythicSim Discord. The conditions only held casts back, so the three rotations
drop every `autoTimeToNext` comparison; Rapid Fire keeps its wait for Aimed Shot. The arena's
`bm_arcane` and the melee Survival rotation are unchanged.

MythicSim builds its presets from these files. Same runs as patch 88, on this branch's engine:

| Reference | 120 s s1 | 120 s s2 | 300 s s1 | 300 s s2 |
| --- | ---: | ---: | ---: | ---: |
| hunter (Beast Mastery) | 594.45 to 597.36 | 594.79 to 597.65 | 587.57 to 589.62 | 587.43 to 589.96 |
| marksmanship-hunter | 669.10 to 671.58 | 668.76 to 671.56 | 611.55 to 612.93 | 611.17 to 612.70 |

The Marksmanship row keeps the condition on the Summon Hawk line MythicSim adds to this rotation.
That one is worth keeping: without it the reference loses 2.8 DPS, as Summon Hawk shares Arcane
Shot's cooldown and the condition is what leaves room for Arcane Shot, not anything to do with
clipping. The ranged Survival rotation on the two ranged references' gear and talents gains 1.1 to
3.9 DPS. Goldens (Average-Default): TestBeastMastery 423.02 to 425.18, TestMarksmanship 311.01 to
310.17, TestSurvival 322.65 to 325.47.

Drop this when upstream's hunter rotations drop the Auto Shot timing.

## 91. `buffs: Flametongue Totem's hit takes no spell power` (dropped, upstream #682)

Upstream's #682 takes the coefficient off too, and also the talents; see patch 70.


Hameru tested Flametongue Totem on the beta (MythicSim Discord #contributors, 6 October 2026; Kerani
and Lazyshadow agree): its hit does not scale with spell power, and in Cat Form it is sized by the
speed of the weapon in the main hand, not by the 1.0 s paw. Patch 70 gave the hit Flametongue
Attack's 0.1 coefficient as an inference; `FlametongueTotemAttack` now has none. Flametongue Weapon's
own hit keeps its coefficient, and the talents that name Flametongue Attack (Elemental Fury,
Elemental Weapons) still reach the totem's hit, which the beta has not been asked about. The hit
still fires on landed main-hand auto attacks only.

The form rule needed no change: the hit already read `Character.MainHand().SwingSpeed`, the equipped
item, which a form leaves alone while its paw swings at 1.0 or 2.5 s. Bear Form is assumed to follow
the same rule as Cat Form; nobody has tested it.

Hameru's rank 4 tooltip reads "18.825 to 61.062", which would be a dummy of about 1526 against the
engine's 1363. It is 1363 with the tooltip's $mult of 1.12: 1363 / 77 x 1.12 - 1 = 18.825 and
1363 / 25 x 1.12 = 61.062. The description (16387) multiplies by $mult, and SpellDescriptionVariables
860 sets it to 1.12 when the reader knows Improved Weapon Totems rank 2 (29193), 1.06 for rank 1
(29192) and 1 otherwise. The 16389 dummy is 1363 on builds 1.60.1.70205 and 1.60.1.70235 (SpellEffect
694279, wago.tools). Forever's talent trees have no Improved Weapon Totems, so the engine keeps 1363.
The 29193 row still carries its +12% Flametongue Totem dummy (SpellEffect 705635, class mask bit 34),
so a beta log of the largest plain hit decides it: a 3.8 s weapon hits for 51.79 at 1363 and 58.01
with the 12%.

`TestFlametongueTotemInFormUsesTheEquippedWeaponAndNoSpellPower` (sim/) puts a Cat and a Bear with
Soulkeeper (3.8 s) under the party totem: the largest plain hit is 3.8 x 13.63 through the druid's
damage multipliers, and 1000 spell damage leaves the total unchanged on the same seed. It fails on the
previous commit. The enhancement test of the old coefficient becomes
`TestFlametongueTotemHitIgnoresSpellDamage`. No golden moves.

Drop this when upstream gives the totem's hit no coefficient.

## 92. `metrics: each action reports the spread of its hits, crits and ticks`

Lazyshadow asked on the MythicSim Discord for the average, smallest and largest damage of an
ability's normal hits and of its crits, for APL work. The result had totals (`damage`, `crit_damage`,
`tick_damage`, `crit_tick_damage`) and counts whose kinds do not line up one for one with them, and no
smallest or largest. `TargetedActionMetrics` gains four `DamageRange` fields (field numbers 37 to 40,
additive, so an older reader skips them):

| Field | Events |
| --- | --- |
| `hit_range` | direct hits that are not critical, glancing, blocked or crushing |
| `crit_range` | direct critical strikes (blocked crits excluded) |
| `tick_range` | periodic ticks that are not critical |
| `crit_tick_range` | critical periodic ticks |

`DamageRange` is `count`, `total`, `min` and `max` over every iteration of the run, as the other
totals are; the average is `total / count`. Partial resists count in their kind. A landed result that
deals no damage is the application of a dot or a debuff (Deep Wounds' trigger, say) and is left out.
A kind with no landed damage is left unset (null in the CLI's JSON). Together with `glance_damage`, `block_damage`, `blocked_crit_damage` and
`crush_damage` the four totals add up to `damage`. JSON names are `hitRange`, `critRange`, `tickRange`
and `critTickRange`.

`SpellMetrics` records each landed damage event in `dealDamageInternal` with a few comparisons and no
allocation, `doneIteration` merges the iteration into the action's target metrics, and
`CombineConcurrentSimResults` merges the ranges of concurrent sims (counts and totals add, the
extremes are kept).

`TestActionMetricsCarryDamageRanges` (sim/) runs an Arms warrior and checks each range is ordered and
holds no zero-damage event,
that the totals add up to the action's damage, and that hits, crits and ticks all appear.
`TestConcurrentResultsCombineDamageRanges` combines two runs and checks the merge. No golden moves.

## 93. `mage: Ignite ignores hits on a unit that is not an enemy` (dropped, upstream #699/#702)

Upstream's fix is the same condition; its test in `sim/mage/mage_test.go` replaces
`sim/mage/ignite_test.go`.


Since 8fb1a2d75a the half of a Goblin Sapper Charge that goes off in the thrower's face is its own
spell (`newSapperSelfDamageSpell` in `sim/core/consumes.go`), with the spell damage proc mask and
the Fire school, so that a listener on spell damage taken hears it. Ignite's trigger in
`sim/mage/talents_fire.go` listens for Fire spell damage crits the mage deals, so a crit of that half
reached it with the mage as the target. The mage carries no Ignite dot, `Ignite.Dot` returned nil,
and the handler panicked on `IsActive`: any Fire Mage with Ignite and the charge crashed the sim as
soon as the self hit crit. MythicSim's `fire-mage-goblin-sapper` request failed every run before
this and now sims 593.67 DPS (3,000 iterations).

The trigger now also requires the target to be an enemy unit. The sapper's proc mask is left alone:
the self hit is a harmful spell landing on the character, which is why 8fb1a2d75a gave it that mask.

`TestIgniteIgnoresSapperCritOnTheMage` (`sim/mage/ignite_test.go`) crits the self damage spell on the
mage, which panicked before the fix, then on the target, which must still ignite it. No golden moves:
the mage suites carry no sapper.

Drop this when upstream's Ignite (or the sapper's self hit) keeps the self hit out of Ignite.

## 94. `core: pushback only pushes back a cast still in progress` (dropped, upstream #700/#701)

Upstream's handler returns on the same `Hardcast.Expires <= sim.CurrentTime` check (its test is
`sim/core/pushback_test.go`). The fork keeps `sim/druid/feralbear/pushback_test.go` as extra coverage
(the 10, 5 and 1 ms cases and a cast still in progress); it passes on upstream's code.


The pushback trigger in `sim/core/character.go` checks that a hardcast is running when the hit
lands, but its handler runs one spell batch window (10 ms) later and did not check again. If the
cast completed in between, `Hardcast.pushBack` moved the finished cast (whose end is then
`startingCDTime`, so it took the full 500 ms) and `newHardcastAction` scheduled it again: the cast
completed a second time at once and its effect landed twice. A hit that rolled at the very instant
the cast ended ran before the completion and moved the cast a full 500 ms. MythicSim's
`feral-bear-druid-boomerang-pushback-after-cast` request shows it: the target's swing lands at 0.49 s
into Linken's Boomerang's 0.5 s cast, and the log has "Completed cast {ItemID: 11905}" twice at 0.50
and two Boomerang hits. With this patch it has one of each and no pushback line.

The handler now returns, before the pushback roll, when the cast has ended
(`Hardcast.Expires <= sim.CurrentTime`), which is when the hardcast action itself counts a cast as
complete. The same return covers a channel that ended inside the window, which the channeled branch
would otherwise have given a new end at the current time and completed again. A hit in the last 10 ms of a cast now never
pushes it back, as the batch window already implied.

`TestPushbackLeavesACompletedCastAlone` (`sim/druid/feralbear/pushback_test.go`) has a tanking Bear
hardcast Linken's Boomerang and takes a swing 10, 5 and 1 ms before it completes. Before the fix the
first completed once at 1 s and the other two completed twice; now each completes once at 0.5 s.
`TestPushbackStillDelaysACastInProgress` checks that a swing at 0.2 s still moves the cast to 0.71 s.
No golden moves: no suite run reaches the handler after its cast has ended.

Drop this when upstream rechecks the cast in the pushback handler, or rolls pushback when the hit
lands.

## 95. `items: weapons from master's database leave out Classic's bonus damage roll`

Bae asked on the MythicSim Discord why Iceblade Hacker's Frost hits looked far too big. The Frost
hits are right: Forever's client gives the axe an equip spell, "Melee attacks with this weapon deal
41 Frost damage" (1298413 triggering 1298414, 40.7 flat in every 1.60.1 build from 69876 to 70235).
What was wrong is the axe's swing. Classic gives some weapons a second damage roll of another
school, "+ 1 - 5 Frost Damage" on Iceblade Hacker. The Forever client has no such roll: it drops the
line or rebuilds it as an equip spell, and Wowhead Forever lists those weapons with their base damage
only (Shadowfang 29 - 55, Torturing Poker 22 - 45, Thunderfury 44 - 115).

The items `mergeForeverSimDB` fills in, the ones the beta client does not ship yet, come from master's
database, which took Wowhead Classic's `damageMinAll` and `damageMaxAll`: the base roll plus the bonus
roll, simmed as Physical. So Iceblade Hacker swung for 58 - 111 and the equip spell's Frost came on top.

`dropClassicBonusDamage` (`tools/database/gen_db/classic_bonus_damage.go`) puts each filled-in weapon
whose range equals Classic's base-plus-bonus total back on its base roll, read from the same Wowhead
Classic planner (`assets/db_inputs/wowhead_gearplannerdb.txt`, `dmgmin1` and `dmgmax1`). A bonus roll is
at least 1 at each end, so a half-point rounding gap between the planner's base and total (Blade of
Eternal Darkness, 33.5 - 69.5 against 34 - 70) is not one. Client rows are left alone.

Three weapons move, all level 57 and up, which the beta client does not ship yet:

| Item | Before | After |
| --- | --- | --- |
| Iceblade Hacker (13952) | 58 - 111 | 57 - 106 |
| Warblade of Caer Darrow (13982) | 143 - 236 | 142 - 214 |
| Ta'Kierthan Songblade (16039) | 130 - 214 | 129 - 194 |

`assets/database/db.json` carries the same three edits, and `go run ./tools/sync_db_binary` rebuilt
`db.bin` from it (run on the unedited JSON first, it reproduced the old binary byte for byte). A full
`make db` should write the same three rows through `dropClassicBonusDamage`; it was not run here, as it
needs `tools/database/wowsims.db`.

Tests: `TestDropClassicBonusDamage` and `TestDropClassicBonusDamageLeavesClientRows`
(`tools/database/gen_db`), and `TestForeverWeaponsLeaveOutClassicBonusDamage` (`sim/core`, with_db), which
reads the embedded binary and fails on the old one. Goldens move only on AllItems rows for Iceblade
Hacker (`TestSurvivalMelee`, `TestProtection`, `TestFury`, `TestArms`, `TestProtectionWarrior`,
`TestEnhancement`) and Warblade of Caer Darrow (`TestSurvivalMelee`, `TestArms`), each by 0.4% to 0.9%.

Drop it when the client ships these items, or master's database carries base damage.

## 96. `hunter: a weaver's rotation acts before the swing that lands on arrival in melee range, so a queued Raptor Strike replaces it`

Sanctum asked on the MythicSim Discord why Raptor Strike waited after his melee-weaving Hunter stepped
in (sim `cb8a969e`). The hunter reached 5 yards at 3.71 s with his main-hand swing ready since the pull.
`UpdatePosition` turned the swing on (`EnableMeleeSwing`) and it landed at once as a white hit. The
rotation's next check was on its 100 ms grid, after the swing, so the Raptor Strike it queued at 3.71
waited a full swing timer and landed at 5.99. `swing()` already runs the rotation before a swing to let
a last-moment Heroic Strike or Raptor Strike in, but `DoNextAction` returns while the rotation timer is
not ready, which is the case for any arrival between two checks.

`AutoAttacks.holdArrivalSwingForRotation` (`sim/core/attack.go`), called from `UpdatePosition` when the
main hand comes into range: for a unit that swings both ranged and melee and has a swing replacer (a
Hunter), a swing already due moves `heldSwingLag` (1 ns) later and the rotation wakes now
(`ReactToEvent`), so it runs first and a Raptor Strike queued on arrival takes the swing. A rotation held
by a Wait or WaitUntil is left alone and the swing lands as before. Arrivals with the swing not yet due
are unchanged: the rotation already queued Raptor Strike before that swing.

Sanctum's request (3000 iterations) goes from 744.6 to 739.2 DPS: the white hit on arrival is gone, and
his rotation steps back out as soon as Raptor Strike lands, so he spends about 1.5 s less in melee per
trip on those arrivals.

`TestRaptorStrikeTakesTheSwingOnArrival` (`sim/hunter/raptor_strike_arrival_test.go`) walks a Hunter
6 yards in 0.857 s, between two 100 ms checks, with Raptor Strike first in its rotation. Before the fix
the arrival swing was a white hit and the Raptor Strike never landed in the 3 s fight; now the arrival
swing is the Raptor Strike. No golden moves: no suite rotation steps into melee.

Drop this when upstream runs the rotation before a swing that comes due on arrival in range.

## Upstream sync 2026-10-07, #677 to #719

Merged ElliotWood/Forever `5c115f1725` (76 commits, #677 to #719 plus data, changelog and arena
commits) into the live pin `cd7d44aec7` (patches 1 to 96). Client data is now 1.60.1.70235 with the
2026-10-06 hotfixes (#706).

Behaviour adopted: Flametongue Totem (#677, #682, see patch 70); hard-cast bolts and Lava Burst hold the
swing (#681, #684, see patch 17); Malediction leaves Hellfire alone (#680); Penance rank 3 for Smite and
Holy Precision / Holy Specialization leave Chastise out (#678, #679); Shadowform refuses Holy Nova and
Chastise (#686); Demonic Embrace keeps its -1% Spirit (#687); balance audit and gear sets (#688, #689);
Cold Blood crits both Mutilate hands and is cast before Mutilate (#690, #691); Unbridled Wrath from white
autos only (#692); Scorpid poison stacks to 5 (#694); Searing Totem every 2.43 sec (#697); Imp Firebolt
every ~0.4 sec (#698); Sanctified Judgement refunds (#705); hunter attack power shares at 24 to 30
(#707); Wowhead data (#708); energy refills smoothly at 10 a second (#709); shadow casts Devouring Plague
on its 1 min cooldown (#710); cat powershifts only with Wolfshead Helm (#711); no Heroic Strike / Cleave
in execute below 40 rage (#712); Intimidation is not an auto-cast cooldown (#713); arcane falls back to
Frostbolt (#714); Arcane Missiles keep Arcane Blast stacks at 15% each (#715); Shadow and Flame preset
casts Conflagrate (#716); Marksmanship paces Sniper Shot and holds Arcane Shot for Aimed Shot (#717,
#718); Power Infusion keeps Shadowform (#719).

Patches upstream now carries, dropped here:

| Patch | Upstream | Notes |
|---|---|---|
| 3 Destruction Conflagrate | #716 | `destruction_conflag` is the Shadow and Flame preset; 194.66 against our rule's 192.34 on TestDestruction. |
| 76 Hawk avoidance | #703 | Upstream's log-fitted hawk (2.5 sec hasted, ~0.35 of the dive bomb base, no crits, can miss and be dodged). |
| 91 Flametongue Totem coefficient | #682 | Also drops the talents; see patch 70. |
| 93 Ignite and the sapper | #699/#702 | Same condition. |
| 94 Pushback after the cast | #700/#701 | Same check; our feral bear test stays as coverage. |

Kept against upstream:

- **17 hard casts hold the swing.** Upstream's #681/#684 restart the swing a full weapon speed after
  every hard cast (`StopMeleeUntil`). Ours (`HoldMeleeForCast`) also lands a swing that came due during
  the cast as the cast completes, as tested on the beta; a cast that completes first resets the timer the
  same way upstream's does. Lightning Bolt, Chain Lightning and Lava Burst call ours; upstream's
  `holdSwingDuringCast` is removed. Upstream's `TestHardCastRestartsTheSwingTimer` passes on ours.
- **70 to 72 Flametongue Totem.** One registration (ours), with #682's hit. A main-hand Flametongue
  Weapon still disables it (patch 72, same rule as #677). Windfury Totem (party flag or the shaman's own)
  switches it off as #677 has it: the merge first kept both, and patch 98 (same day) follows upstream.
  Upstream's `TestFlametongueTotem` is taken whole, its Windfury case included.
- **71 party Flametongue Totem.** Upstream has no party flag.
- **Feral Cat default rotation.** #711 adds a Wolfshead Helm guard to the powershift rows; the fork's
  rotation (patch 42) has none, so it is unchanged.

`db.json`: `scripts/forever-merge-db.py` from base `67f14b04a5`, 13 new upstream items and 5 spell
icons, planner armor re-applied to 15. Four field conflicts, all armor: upstream turned Revelosh's Boots,
Armguards and Spaulders (9387 to 9389) and Ironaya's Bracers (9409) into random-property items and
removed their stat rows; the planner's armor (patch 7) is kept. Patch 95's three weapon rows survive.
`go run ./tools/sync_db_binary` rebuilt `db.bin` and `leftover_db.bin`.

Goldens, Average-Default, fork old to new against upstream's own move over the same range. Every
movement was traced to an upstream commit by reverting that commit alone on the merged tree, which put
the suite back on the old golden exactly (Shadow, Survival, Survival Melee, Marksmanship, Feral Cat,
Elemental):

| Suite | Fork | Upstream | Cause |
|---|---|---|---|
| Arcane | +3.45% | +3.45% | #715 |
| Balance | -1.91% | -1.91% | #688 |
| Feral Cat | -0.08% | +3.31% | #709 only; #711 changes rows the fork's rotation does not have |
| Beast Mastery | 425.18 to 406.29 (-4.44%) | -7.48% | #703, #704, #713. The fork's old hawk already rolled avoidance (patch 76); end value is 0.27% over upstream's, the size of patches 88 and 90 |
| Marksmanship | +7.55% | +6.54% | #704, #717, #718 together; the rest is patch 90's missing Auto Shot gates |
| Survival | +3.41% | +3.45% | #704 |
| Survival Melee (fork only) | 385.47 to 396.59 (+2.89%) | n/a | #704 |
| Retribution | +0.34% | +0.34% | |
| Shadow | +1.12% | +2.75% | #710. Smaller here because the suite's Troll and Night Elf priests have no Dark Sacrifice (patch 85), so the extra plagues are mana-bound |
| Smite | +7.38% | +7.79% | #678, #679; same Dark Sacrifice gap |
| Assassination / Combat / Subtlety | +0.88 / -0.05 / -0.17% | same | #690, #691, #709 |
| Elemental | -2.22% | -2.11% | #697 only. Rows spread differently (median +0.57% against -2.64%) because the fork's preset casts Fire Nova above 30% mana (patch 28) |
| Enhancement | -1.67% | -1.67% | #697 |
| Destruction | 192.50 to 194.66 | 182.02 to 194.66 | #716; patch 3 dropped, so the suite runs upstream's rotation |
| Arms / Fury | +4.77 / +2.35% | +4.80 / +2.38% | #692, #712 |

Single item rows that move differently: Warblade of Caer Darrow and Dragon's Call, which the fork's own
item patches (95, 86) changed, and two Balance rows not traced further (Shard of the Gods -2.35% against
-3.20%, Cenarion Raiment +0.03% against -0.53%; 2 of 162 rows). Whole suite green after the goldens were
regenerated.

## 97. `shaman: Totem of Thunder adds 1% crit to Lightning Bolt`

Bae asked on the MythicSim Discord (7 October 2026) why Totem of Thunder (228176, the crafted shaman
relic) did nothing in the sim. Its equip spell 461295 reads "Increases the critical strike chance of
Lightning Bolt by 1%": one A_ADD_FLAT_MODIFIER effect, misc 7 (crit chance), base points 1, on family 11
class mask word 0 bit 0. Nothing implemented it; the item generator listed it as a TODO in
`sim/common/forever/stat_bonus_procs_auto_gen.go`.

`sim/shaman/items.go` registers it the way Totem of the Storm is: a permanent aura with a
`SpellMod_BonusCrit_Percent` mod, read from the row (1), registered with `ItemSwap.RegisterProc`.
The mask follows the client's class flags: every Lightning Bolt rank (403 to 15208) and every Lightning
Overload Lightning Bolt row (408439 to 408477) carries bit 0, so the mod names
`SpellMaskLightningBolt | SpellMaskLightningBoltOverload`. Chain Lightning and its overload rows carry
bit 1 and are left out. The TODO block is removed from the generated file by hand; the generator skips
items with an effect (`core.HasItemEffect`), so the next DB run writes the same file.

`TestTotemOfThunderAddsOnePercentLightningBoltCrit` (`sim/shaman/elemental`) checks the row's mask
against the client rows, then compares every registered spell's bonus crit with and without the relic:
all 10 Lightning Bolt ranks and all 10 overloads gain exactly 1%, nothing else moves. It fails with the
item effect removed. Goldens: one new AllItems row in each of `TestElemental` (120.64 DPS) and
`TestEnhancement` (153.88 DPS); nothing else moves, since no preset equips the relic.

Drop it when upstream implements Totem of Thunder.

## 98. `shaman: Windfury Totem switches Flametongue Totem off`

Follows upstream #677. The Forever beta development notes say Flametongue Totem "no longer stacks" with
Windfury Totem; they do not say which totem holds. Upstream has Windfury hold, and the players on the
MythicSim Discord (7 October 2026) treat the two as either/or, so the fork now does the same. Until this
patch the merge kept both: the client rows (70235) give no exclusivity, the two are party proc auras in
different totem slots, and no log has shown them failing to stack. If a log shows both proccing on the
same swings, drop this patch.

`buffs.WindfuryTotemDisablesFlametongueTotem` joins an aura to `FlametongueTotemCategory` at 1.5 times
the totem's bid (2044.5, between the totem's 1363 and a main-hand Flametongue Weapon's 2726), the
place upstream's `FlametongueTotemWindfuryTotem` holds between `FlametongueTotemCast` and
`FlametongueTotemMainHandImbue`. Both Windfury Totem auras call it, as upstream's do:
`driveWindfuryTotem`'s "Windfury Totem" (the party flag, joined after the air slot so a cast Grace of
Air refuses it first) and `registerWindfuryTotemSpell`'s "Windfury Totem (Self)". Both Flametongue
Totems, the party flag (patch 71) and the shaman's cast, are in the category, so all four pairings
switch off. The Flametongue Totem stays down and its mana is spent; it adds no hits. Grace of Air is
not in the category. A cast Grace of Air that takes the air slot from the party's Windfury Totem lets
the Flametongue Totem back on. As upstream, a main-hand Windfury Weapon beside a Windfury Totem leaves
the totem standing, so Flametongue Totem stays off there too: for the whole fight with a cast Windfury
Totem, and from 5 s in with the party's, whose aura the imbue knocks off at the pull until its next 5 s
refresh (`driveWindfuryTotem`'s `OnExpire`, same code upstream). That gap gives about 1.6 totem hits a
fight on a 3.8 speed weapon; not worth a fork-only change to the driver.

Tests (`sim/shaman/enhancement/flametongue_totem_test.go`): `TestWindfuryTotemSwitchesFlametongueTotemOff`
replaces `TestFlametongueTotemAndWindfuryTotemDoNotInteract` (both cast in either order, party Windfury
with cast Flametongue, cast Windfury with party Flametongue, both party, and party Windfury with party
Grace of Air: extra attacks still come, no totem hits, the cast totem stays up).
`TestGraceOfAirLeavesFlametongueTotemAlone` checks one totem hit per landed main-hand swing beside a cast
or party Grace of Air, and after a cast Grace of Air replaces the party's Windfury Totem. Upstream's
Windfury case is back in `TestFlametongueTotem`. All three fail with the two joins removed.

Goldens: none move; no suite preset puts a Flametongue Totem beside a Windfury Totem. Whole suite green.

Drop it when upstream changes the rule, or keep it as long as upstream carries the same rule (then it is
only the party-flag half, since upstream has no party Flametongue Totem).

## 99. `druid: Rake adds 5.26% of attack power to its hit and every tick`

The client rows at the pinned build carry no BonusCoefficientFromAP on Rake, and Blizzard's Druid deep
dive (30 September 2026) only says it "now gains increased damage from Attack Power". The share is a
fit to beta combat logs (Hameru, MythicSim Discord #bugs, 7 October 2026, level 30, non-crits on many
mobs): at 324 attack power Rake hit 40 and ticked 33, 33, 33; at 225 it hit 35 and ticked 28, 28, 27.
That rank's base is 23 on the hit and 16 a tick, so hit = 23 + 0.0526 x AP (40.0, 34.8) and tick =
16 + 0.0526 x AP (33.0, 27.8), every logged number within rounding.

`sim/druid/rake.go`: `rakeAttackPowerShare = 0.0526` rides on whatever rank the engine casts, which keeps
its own base from client data (rank 4, 9904: 61 on the hit, 34 a tick). `rakeHitDamage` adds the share
to the initial hit in `ApplyEffects`; `rakeTickDamage` goes into the snapshot, followed by
`dot.SnapshotAttackPowerShare`, as Rip does, so a tick reads the share from the attack power the druid
has then (`currentTickInputs` swaps the snapshotted share for the current one, nothing counted twice).
The non-snapshot `ExpectedTickDamage` path (the projection `UpdateBleedPower` stores, which nothing in
the feral rotation reads) uses the same formula.

Tests: `TestRakeAddsAttackPowerShareToHitAndTick` (`sim/druid/rake_test.go`) checks the formula over
Hameru's 99 attack power gap. `TestRakeScalesWithAttackPower` (`sim/druid/feralcat/rake_test.go`) drives
the live spell on a naked Cat with crit removed: the non-crit hit (each on a target that is not
bleeding, out of Rend and Tear's reach) and a tick, each divided by the multipliers on top of its base,
gain 0.0526 x 99 from 99 more attack power, on a running dot snapshotted at the old attack power and on
a fresh one.

Goldens: none move. Neither Feral Cat APL (`default`, `simple_vael`) casts Rake, so `TestFeralCat` is
unchanged; Bear, Balance and Restoration do not use it. A player APL that casts Rake gains about 20
damage on the hit and on each tick at 391 attack power.

Drop it when a client build carries Rake's attack power coefficient (then read it from the row), or
refit if a log at another level disagrees.

Follow-up: `TestDotsReadStatsAtTheTick` (`sim/dot_rules_test.go`) listed Rake with no share; its row now
carries `meleeAP: 0.0526`.

## 100. `core: a temporary enchantment beside the weapon's imbue`

Blizzard's class deep dives (Mage and Shaman, 7 October 2026; Rogue and Warlock, 8 October) let a
temporary enchantment (Sharpening Stones, Weightstones, Wizard Oil) sit on a weapon beside a class
imbue: a Rogue's poisons, a Shaman's weapon imbue, a Warlock's Firestone or Spellstone. The engine had
one consumable imbue per hand (`mhImbue_id`, `ohImbue_id`), and a Rogue's poisons live there, so a
Rogue could not carry a stone. A Shaman's imbue is a class option and already stacks.

`ConsumesSpec` gains `mh_temp_enchant_id` (27) and `oh_temp_enchant_id` (28). `applyConsumeEffects`
registers each through `registerStaticImbue`, as the imbue slots are, so a stone or oil there gives
the same stats and flat weapon damage, and `MHImbueFlatWeaponDamage` counts it for the classes that
rebuild the main hand. Nothing else reads the slots: a poison stays in the imbue slot and keeps its
hand (`getPoisonProcMask`).

Tests: `TestTempEnchantStacksBesideTheImbue` (`sim/core/temp_enchant_test.go`): an Elemental Sharpening
Stone beside Wizard Oil adds 2% melee crit, Wizard Oil beside a stone 24 spell damage, an off-hand
Dense Sharpening Stone 8 to the off hand only, and a main-hand Dense Weightstone counts in
`MHImbueFlatWeaponDamage`. `TestPoisonsStackWithStones` (`sim/rogue/poison_stone_test.go`): Instant
and Deadly Poison keep their hands beside an Elemental and a Dense Sharpening Stone, which add 2% crit
and 8 off-hand damage.

Goldens: none move; no suite request sets the new fields.

Drop it if upstream adds the same slots (then map to theirs).

## 101. `warlock: Firestone and Spellstone as weapon imbues`

The same deep dive makes Firestone and Spellstone weapon imbues that stack with a temporary
enchantment such as Wizard Oil, instead of held off-hand items. `WarlockOptions.weapon_stone`
(`NoWeaponStone`, `Firestone`, `Spellstone`) chooses one. `registerWeaponStone`
(`sim/warlock/weapon_stones.go`) reads the Create spell's triggered aura at its highest rank from
client data and applies it as a permanent aura named after it: Firestone (23483) is A_MOD_DAMAGE_DONE
21 on mask 4 (Fire) and A_MOD_SPELL_CRIT_CHANCE 2; Spellstone (1237165) is A_MOD_DAMAGE_DONE 21 on mask
36 (Fire and Shadow) and A_MOD_CASTING_SPEED_NOT_STACK 2. The deep dive names Shadow damage only for
the Spellstone; the client row's mask also carries Fire, and the row is what is applied.

Tests: `TestWeaponStonesStackWithWizardOil` (`sim/warlock/weapon_stones_test.go`): +21 Fire damage and
+2% spell crit for the Firestone, +21 Fire and Shadow damage and x1.02 cast speed for the Spellstone,
neither touching the other's stats, both beside Brilliant Wizard Oil's 36 spell damage.

Goldens: none move; no suite request sets the option.

Drop it when upstream models the stones (then compare its values with the client rows above).

## Upstream sync 2026-10-10, #720 to #760

Merged ElliotWood/Forever `071e0cb439` (69 commits, #720 to #760 plus data, changelog and arena commits)
into the live pin `66837240ec` (patches 1 to 102). Client data is now 1.60.1.70334 (#760), with the
2026-10-08 hotfixes (#752) and Wowhead data through 2026-10-10.

Behaviour adopted: client spell ranges and the range talents for Mage, Warlock, Priest, Shaman and
Balance (#728, #729, #731, #732, #733); point-blank AoE needs 10 yards (#725, #728); multi-target lines in
the Fire, Arcane and Frost (#724, #725, #747, #749), Warlock (#735), Balance (#738), Hunter Volley (#739,
#743), Warrior Sweeping Strikes and Protection Cleave (#740, #748), Bear Swipe (#741), Retribution
Consecration (#744), Protection paladin Hammer of the Righteous (#754) and Enhancement Fire Nova (#737,
#750) rotations; Ice Lance coefficient 0.1 (#726); Retribution Aura scales with the holder's spell power
(#736); Protection paladin judges on cooldown (#746); the suites run the rotations the site picks (#721,
#722, #723).

Patches upstream now carries, dropped here:

| Patch | Upstream | Notes |
|---|---|---|
| 11 Holy Nova healing crit | #751 | Same defense type; our regression test passes on upstream's code. |

Kept against upstream:

- **5 melee Survival.** See patch 5: ours plus upstream's move into melee.

Conflicts: `db.json` and `leftover_db.json` by `forever-merge-db.py` (seven new upstream items in
`db.json`, Tabard of the Scarlet Crusade in `leftover_db.json`; Revelosh's Gloves, which upstream turned
into a fixed rare without armor, keeps our 38 armor), then `sync_db_binary`. `presets.ts` takes
upstream's `Survival (melee)` label.

Goldens: each regenerated suite moves with upstream's own goldens over the same range (Balance +9.3%,
Arms +16.3% against upstream's +16.6%, Protection paladin +12.7%, Enhancement +106.7% against +106.4%,
Frost +47.6%), except where a patch meets a new upstream line: Bear multi-target moves 7% past upstream
because our Swipe attack power share (client 70170) now counts once #741 swipes on 3+ targets; Marksmanship, Beast
Mastery, Affliction and Destruction multi-target rows differ from upstream by up to 5% where patch 90's
shot timing meets #739's Volley; `TestSurvival` follows patch 5. Every other test passes.

## Upstream sync 2026-10-11, #761 to #763

Merged `ElliotWood/Forever` master `f7598154eb` (#761, #762, #763) into the live pin `8361dcca54`
(patches 1 to 102). Client data moves from 1.60.1.70245 to 1.60.1.70338 with the 2026-10-10 hotfixes.
The previous sync's pin comment said client 70334, but #757 and #760 only changed
`forever_client_build.txt`; the spell store was last regenerated at 70245 (#752). This is the first pin
that simulates the 70291 class changes: Penance (bolt coefficient 0.285 to 0.19, rank 3 bolt 180 to 71 at
level 60, rank 4 131 to 92, higher mana costs), low-rank rescaling (Smite ranks 1 to 4, Frostbolt rank 1
0.5 to 0.2 a level), Thorns 22 to 18, Retribution Aura 30 to 20, and Predatory Instincts renamed Natural
Instinct (#762).

Conflicts: `talents_feral_combat.go` keeps our `applyImprovedShiftingPower` call and takes upstream's
`applyNaturalInstinct`. `db.json` and `leftover_db.json` by `forever-merge-db.py` (new upstream items and
the Expose Prey icon; planner armor re-applied to 18 and 19 items), then `sync_db_binary`.

`tools/spelldata/describe_test.go` pins Frostbolt rank 1 and a heal from the client store; both follow
70338 (ppl 0.5 to 0.2 and 0.9 to 0.3).

Goldens move with upstream's own over the same range: Smite -12.1% (upstream -13.7%; upstream's Smite APL
also moves to Penance rank 4), Bear -2.7%, Protection paladin -1.6%, Protection warrior -1.6%, Elemental
+1.2% (upstream +1.15%), Survival +0.9% against +0.2% (patch 5's melee rotation meets the Immolation Trap
rename). Every other test passes.

## 103. `druid: Howling Idol takes 1 second off Shifting Power's cooldown`

Item 272427 carries buff 1291059, "Reduces the cooldown of your Shifting Power ability by 1 sec.":
`A_ADD_FLAT_MODIFIER`, `SPELLMOD_COOLDOWN`, -1000 ms, class mask word 1 `0x80000000`, which is Shifting
Power's. The generated proc table only had it as a commented "manual implementation required" block, so
the idol did nothing (reported by Hameru on the MythicSim Discord).

`sim/druid/items.go` registers it as a flat cooldown mod on `DruidSpellShiftingPower`, so it stacks with
Improved Shifting Power. Test: `TestHowlingIdolShortensShiftingPower` (`sim/druid/feralcat`), 16 s
without the idol and 15 s with it. Goldens: new `AllItems-HowlingIdol-272427` rows in Balance, Bear and
Cat.

Drop it when the generator emits the effect.

## 104. `paladin: Hammer of Wrath and Holy Wrath hold the swing like a Shaman bolt`

Both used `StopMeleeUntil`, which pushes a swing due mid-cast to a full swing after the cast. The Forever
rule measured on the beta for Shaman bolts (patch 17, #681, #684) is `HoldMeleeForCast`: a swing due
mid-cast lands as the cast completes, and a cast that completes first restarts the timer. Both now hold
the swing on a hard cast only, so an instant Hammer of Wrath (Instrument of Law) leaves the timer alone,
as before. Raised by Redfall on the MythicSim Discord (spell casts and the swing timer).

Goldens: none move (Retribution and Protection paladin read the same with and without it).

Paladin heals (`flash_of_light.go`, `holy_light.go`, `lights_vigil.go`) still have no swing hook. No
preset heals while meleeing; add the same hold if one does.

## 105. `core: a completed hard cast restarts the melee swing timer`

Patch 17 held a swing that came due during a hard cast and let it land as the cast completed. Upstream
#681 (beta logs 2701 and 2712 on foreverlogs.gg) shows the next swing one full weapon speed after every
completed hard cast, even when a swing came due during it, which is also the Classic rule. Redfall
(MythicSim Discord, 10 and 11 October) saw the Elemental with melee buffs "autoing right after spells
finished casting" in the combat log and expected the reset.

`AutoAttacks.HoldMeleeForCast` now puts the next swing at cast end plus a full swing. A swing landing at
the instant the cast starts still lands, and `castResetAt` pushes the one after it past the cast.
Callers: Shaman Lightning Bolt, Chain Lightning and Lava Burst (patch 17), Paladin Hammer of Wrath and
Holy Wrath (patch 104), and now Light's Vigil, Flash of Light and Holy Light, and druid cast-time spells
through `RegisterSpell`. Hunter shots have a minimum range, so a hunter never hard-casts in melee. Instant
casts leave the timer alone.

Test: `TestLightningBoltRestartsTheSwing` (`sim/shaman/enhancement/swing_hold_test.go`) reads the log of an
Orc chain-casting Lightning Bolt and checks that no swing lands inside a cast and that the first swing after
one comes at least a full swing after it completes. Goldens: `TestProtection` (paladin), -0.0% to -1.5%.

References (default path, 10,000 iterations, board seed, board length): Elemental with Melee buffs 531.27 to
461.03 (-13.2%); Shockadin 885.03 to 867.93, Protection Paladin 560.82 to 555.76, Holy-Protection 567.00 to
559.96, each back to the live pin's value (patch 104 had raised them); every other reference unchanged.

Patch 17's "held swing" no longer applies; its test is replaced by this one.

## 106. `hunter: Mongoose Bite is a main-hand attack`

Mongoose Bite carried `ProcMaskMeleeSpecial` (main-hand and off-hand special). Predator's Edge's "+50%
Off Hand weapon damage" (1310627, `A_MOD_OFFHAND_DAMAGE_PCT`) is a `ProcMaskMeleeOH` damage mod, so a
dual-wielding hunter's bite, which deals main-hand weapon damage, took the +50% too and rolled the
off-hand weapon's procs. A two-hander was unaffected (the off-hand mod needs dual wield). Now
`ProcMaskMeleeMHSpecial`, like Raptor Strike, Strider Kick and Wing Clip. Raised by ShinShimon on the
MythicSim Discord (Mongoose Bite averaging 851 a hit against 300 for the main-hand swing).

Survival reference (dual wield, 10,000 iterations, board seed, board length): 880.75 to 835.88 DPS
(-5.1%); Mongoose Bite 847 to 568 a cast, the Lacerating Strikes bleed (40% of the bite) with it.

Tests: `TestPredatorsEdgeOffHandBonusSkipsMongooseBite` and `TestSavageStrikesRank2OnEveryMeleeAbility`
(`sim/hunter/dual_wield_melee_test.go`, `with_db`). The second pins Savage Strikes rank 2 at +4% crit on
Raptor Strike, Mongoose Bite, Strider Kick, Wing Clip and the bleed: a Discord report read it as +2%
from the per-attempt crit column, where a yellow hit's separate crit roll only runs on a landed hit, so
the boss's 6.5% dodge dilutes it. Goldens: `TestSurvival` and `TestSurvivalMelee` `AllItems-*` rows for
one-handed weapons (dual wield) move down; two-hander rows do not.

## 107. `items: Forever's caster weapon spell power and master's missing healing`

Redfall on the MythicSim Discord: swapping a Warlock's Rod of the Ogre Magi (23 spell damage) for
Crackling Staff, which shows 99 spell damage in game, gained 2 spell damage and 25 healing in the sim.
Crackling Staff's row here was Intellect 25 and spell power 25 (`ItemSparse` 45). Wowhead's Forever
planner (`assets/db_inputs/wowhead_forever_gearplanner.txt`) has `splpwr` 25 plus `spldmg` 74 and
`splheal` 74, and the three add up: the client rows that carry them agree (Earthen Silk Slippers,
`splpwr` 18 and `spldmg` 6, is 24 SpellDamage and 18 HealingPower here; Ogre Loincloth, `splpwr` 6
and `splheal` 12, is 6 and 18). The 74 is Forever's caster weapon spell power: it goes with the lower
caster weapon damage (102 - 154 in the planner), scales with item level and quality (74 on an item
level 65 rare, 94 on a 65 epic, 106 on a 78 epic) and a healer's weapon gets a third of it as damage and
three times that as healing (Benediction 65 and 195). Neither the client rows nor master's database
carry it, so 159 weapons had no weapon spell power at all: Jin'do's Judgement, Ironbark Staff,
Sageblade, the Grand Marshal's and High Warlord's caster weapons, Benediction, the AQ staves.

Separately, master's database (`forever_sim_db.json`, the items `mergeForeverSimDB` fills in) states
"damage and healing" as SpellDamage alone, and the engine heals off HealingPower only
(`Spell.HealingPower`), so those items gave casters no healing: Rod of the Ogre Magi, the Tier 1 and 2
caster sets, Spellweaver's Turban.

`tools/database/gen_db/planner_spell_power.go`, after `fillPlannerArmor`:

- `fillPlannerSpellPower` raises each item's SpellDamage to `splpwr + spldmg`, HealingPower to
  `splpwr + splheal` and the school stats to `arcsplpwr` and the rest, where the planner's is higher.
  It never lowers a stat. Items with area stats are skipped, as the planner adds those in (Jewel of the
  Guard Captain: 8 everywhere and 16 more in the open world, planner 24). 205 items move: the 159
  caster weapons; the four Atiesh variants, which had little or no spell power; Soul Harvester's 73
  Shadow; 36 armor pieces and held items with healing that master's rows lacked (Orb of Lorica among
  them); and five school rows (Bad Mojo Mask, Robes of the Lich, Nature's Embrace, Robes of Servitude,
  Pyric Caduceus). On spell power the planner never has less than the database.
- `fillClassicHealing` gives the filled-in items the Forever planner does not list their healing from
  the Wowhead Classic planner's `splheal`, when they have none: 141 items, Rod of the Ogre Magi 23.

Spell crit and hit need nothing: Forever states gear crit and hit as the generic rating, which the
database reads as melee (Ironbark Staff's 2% is 28 melee crit rating, as on all 512 crit items in the
planner), and `UnifyGearHitAndCrit` (`sim/core/forever_rules.go`) pays it into the spell pool too.
Planner and database agree on crit and hit for every item but the Atiesh variants. Redfall's 1% crit
loss is Rod of the Ogre Magi's own spell crit (13: 14).

Not changed: these weapons keep Classic's weapon damage (Crackling Staff 138 - 207, the planner 102 -
154), which only matters to a caster who melees.

`assets/database/db.json` carries the same 346 edits, made by the two functions above run on the
committed JSON (the filled-in set taken as the items in `forever_sim_db.json` whose stats match it and
that client 69893's `ItemSparse` lacks), applied line by line so no other record moves; a rerun changes
nothing. `go run ./tools/sync_db_binary` rebuilt `db.bin` (run on the unedited JSON first, it reproduced
the old binary byte for byte). `leftover_db` has no affected item. A full `make db` was not run, as it
needs `tools/database/wowsims.db`. The site's `scripts/forever-merge-db.py` raises planner spell power
on a merge the way it re-applies planner armor.

Tests: `TestFillPlannerSpellPower`, `TestFillPlannerSpellPowerKeepsHigherValues` and
`TestFillClassicHealing` (`tools/database/gen_db`); `TestForeverItemSpellPower` and
`TestIronbarkStaffCritCountsForSpells` (`sim/core`, with_db, on the embedded binary; the first fails on
the old one). Goldens: one AllItems row each in Balance and Bear (Symbols of Unending Life), Arcane, Fire
and Frost (Coldflame Saber), Elemental and Enhancement (Gift of the Gathering Storm).

Drop it when the client rows or upstream's database carry Forever's weapon spell power, and master's
rows their healing.

## 108. `paladin: Light's Vigil adds its strike to Holy Shock instead of replacing it` (reverted)

Reverted on 11 October 2026: a vigiled Holy Shock again deals only the vigil's strike (upstream's behaviour). Redfall asked on the MythicSim Discord (#mythic) not to change Light's Vigil until it is seen in game, and the only requests came from the discussion he named; there is no beta log either way. Redo this from a combat log. The text below records what the patch did.

A Holy Shock on a target carrying the paladin's Light's Vigil dealt only the vigil's strike (395 at
rank 3) and none of its own damage (334 to 362 at rank 4). The client ships the strike as a damage
spell of its own (1310914, 1311592, 1311598, 0.429 coefficient), as the friendly path's party heal
(1310912, 1311591, 1311596) is a heal of its own beside Holy Shock's, and 1311595's tooltip reads as
an addition: the next Holy Shock "triggers no cooldown and causes ... enemy targets to suffer
$1311598s1 Holy damage and refund $s2% of Light's Vigil's Mana cost". The behaviour is server
scripted (Holy Shock 1311606 and the vigil cast are `E_DUMMY`), so no client row settles it
directly; no beta log has been checked. Raised by ShinShimon on the MythicSim Discord
(#tier-list-feedback, 10 October).

`holy_shock.go` now deals Holy Shock's damage, then consumes the vigil (strike and refund) and resets
the cooldown, as before. Test: `TestLightsVigilAddsToHolyShock` (`sim/paladin/retribution`), all
three ranks: every Holy Shock cast has a damage outcome, and each vigiled one a strike.

Goldens: none move (no preset casts Light's Vigil).

Light's Vigil is a 1.5 sec hard cast (`CastTimeMs: 1500`, Spiritual Focus covers its pushback) with no
swing hook yet; a hold or restart would go in its `Cast.ModifyCast`, as in `holy_wrath.go`.

## 109. `rogue: one switch for Cold Blood and Mutilate`

No behaviour change. Upstream #690 modelled what the beta logs (2701/2702) and the live client show:
Mutilate's hand strikes take Cold Blood's crit without spending the charge, so Cold Blood stays up
through every Mutilate until a Sinister Strike, Backstab, Ambush or Eviscerate. Players confirm it is a
live bug (Ambiwar, Redfall and Dyst on the MythicSim Discord, 10 October), and the sim keeps modelling
it until Blizzard fixes it.

`ColdBloodSpentByMutilate` in `sim/rogue/talents_assassination.go` (default false, the live bug) is the
one switch to flip when they do: set true, the Mutilate's off-hand strike spends Cold Blood, so both
hands of one Mutilate keep the crit and the next Mutilate does not. Test:
`TestColdBloodSpentByMutilateWhenFixed` (`sim/rogue/mutilate_test.go`); `TestColdBloodLastsThroughMutilate`
still covers the default. Goldens: none move.

Measured on the assassination-rogue reference (Venom 38/11/2, Ambiwar's rotation, which casts Cold
Blood on cooldown and Eviscerate only in the last 10 s): with the bug Cold Blood is up 115 of 120 s and
98.4% of Mutilate hand strikes crit; switched to the fix, 47% crit and the reference loses 8.9%.

## 110. `hunter: Dust Cloud yields to Sunder Armor and Expose Armor`

Players report that in Forever the physical armor reductions (Sunder Armor, Expose Armor, Dust Cloud)
don't stack and only the strongest applies, as the magic ones (Faerie Fire, Curse of Recklessness,
patch 16) don't. No combat log confirms it and the client rows carry no stacking rule. Sunder and Expose
already share the single-aura `MajorArmorReduction` category; Dust Cloud had its own and stacked its 505
on top of them.

`DustCloudYieldsToSunderAndExpose` in `sim/hunter/pet_abilities.go` (default true): the Tallstrider
holds Dust Cloud while anything in `MajorArmorReduction` is on the target. Both outweigh 505 (Sunder 450
a stack from the raid's first global cooldown, 2250 at five; Expose 2250), so only the stronger applies
and the pet keeps the focus for Bite. Dust Cloud does not join the category itself: the raid's Sunder
ramps from 0 stacks and a single-aura category refuses a weaker newcomer, so a Dust Cloud cast first
would shut Sunder out for the fight. False restores stacking.

Test: `TestDustCloudYieldsToSunderAndExpose` (`sim/hunter/dust_cloud_test.go`): Dust Cloud never goes up
beside Sunder Armor or Expose Armor, stays up beside Faerie Fire and Curse of Recklessness, and stacks
again when switched off. Goldens: none move (no suite runs a Tallstrider beside Sunder or Expose). No
reference build uses a Tallstrider.

Drop it if a combat log shows Dust Cloud stacking with Sunder Armor or Expose Armor.

## 111. `druid: a potion no longer drops Bear or Cat Form`

`registerFormBreakingConsumes` (upstream 0677650f4a) dropped a druid out of Bear or Cat Form after any
potion, conjured item or explosive, and the automatic use of each waited for caster form. In Forever a
druid drinks a potion in form and stays in it (Redfall, MythicSim Discord, 11 October 2026), so potions
leave the list: a bear's Mighty Rage Potion lands its rage and Strength in Bear Form, and the automatic
potion use no longer waits for caster form. Conjured items and explosives still drop the form (the cat's
Goblin Sapper Charge relies on it).

Test: `TestBearDrinksMightyRageInForm` (`sim/druid/feralbear`). Goldens: none move.

## 112. `items: Iceblade Hacker and Warblade of Caer Darrow Frost always hits`

The two hand-written weapon procs in `sim/common/classic/items_store_gaps.go` rolled their Frost
(1298414 and 1298499) on the full melee special table: miss, dodge, parry, block and crit. They only
fire off a swing that already landed, so every swing was rolled for avoidance twice. A player on the
MythicSim Discord (12 October 2026) asked why Iceblade Hacker's Frost showed dodges like a swing. Both
spells carry the client's always-hit attribute (`ATTR_EX_3_ALWAYS_HIT`, 0x40000 in attribute word 3,
`Attr: {3: 0x40070000}`), which the hand-written entries ignored.

The server honours that attribute on a physical attack: Forever Logs report 2701 has Trapz's Summon
Hawk dive (1293241, always-hit, ranged defense) landing 20 of 20 while the same hunter's Auto Shot
missed 7 of 75 in the same fights. No log has either weapon yet; both drop at level 57 and up.

The Frost now uses `shared.OutcomeMeleeCritOnly`, a new `OutcomeType` mapped to
`spell.OutcomeMeleeSpecialCritOnly`: it never misses, is dodged, parried or blocked, and still rolls
to crit, as Summon Hawk's dive does. Before the fix the rogue test below saw 171 of 1,165 procs
avoided (14.7%).

Test: `TestIcebladeHackerFrostAlwaysHits` (`sim/rogue/weapon_proc_hand_test.go`). Goldens: only the
AllItems rows for the two weapons move.

| Row | Before | After |
| --- | ---: | ---: |
| TestArms Iceblade Hacker | 288.39 | 299.38 |
| TestArms Warblade of Caer Darrow | 377.48 | 378.15 |
| TestFury Iceblade Hacker | 347.77 | 355.45 |
| TestProtectionWarrior Iceblade Hacker | 264.23 | 274.89 |
| TestProtection Iceblade Hacker | 288.45 | 298.55 |
| TestEnhancement Iceblade Hacker | 324.67 | 323.22 |
| TestSurvivalMelee Iceblade Hacker | 387.38 | 388.83 |
| TestSurvivalMelee Warblade of Caer Darrow | 471.40 | 472.24 |
| TestSurvival Iceblade Hacker | 329.36 | 331.15 |
| TestSurvival Warblade of Caer Darrow | 414.70 | 410.47 |

The two small drops are the random stream shifting: the proc no longer draws its avoidance roll.

Drop this when the store generates these procs and honours the always-hit attribute.

## 113. `core: a critical auto attack gives twice the rage`

Client 70170 brought back extra rage on a critical auto attack at 75% (WoW: Forever Podcast episode 2), and
`CritRageMultiplier` in `sim/core/rage.go` was 1.75. The beta development notes raise it to 100% for both
classes that use rage: the Warrior notes of 1 October 2026 (https://us.forums.blizzard.com/en/wow/t/2360696)
and the Druid notes of 8 October 2026 (https://eu.forums.blizzard.com/en/wow/t/631316/5) both read
"Additional Rage generated from landing Critical Strikes increased to 100% increased Rage (was 75%)". The
number is server side; no client row carries it.

`CritRageMultiplier` is now a variable, 2.0 by default; 1.75 restores the 70170 value. Every rage bar reads
it: Warriors, Bear Form, and the Cat's bar (fed only if a custom rotation puts the Cat in Bear Form, so the
Cat reference doesn't move). Đegausser (MythicSim Discord, #bugs) saw 8.30 x 1.75 on crits: 8.30 is
3.46 x 2.4, a 2.4 speed one-hander's main-hand swing, so that was either a log from before the 1 October
build or a MythicSim log from the 1.75 engine; it says nothing against the notes.

Tests: `TestCritAutoAttackRageIsTwiceAHit` and `TestCritRageMultiplierSwitchesBackToSeventyFivePercent`
(`sim/core/rage_crit_test.go`), `TestAutoAttackRageGeneration` (`sim/core/rage_test.go`),
`TestCritAutoAttackGivesAWarriorTwiceTheRage` and `TestCritAutoAttackGivesATwoHanderTwiceTheRage`
(`sim/warrior/dps`), `TestBearCritAutoAttackGivesTwiceTheRage` (`sim/druid/feralbear`) and
`TestCatsBearFormCritAutoAttackGivesTwiceTheRage` (`sim/druid/feralcat`). `TestFlurryAxeAndElectrifiedDagger`
now runs 5 iterations: its single iteration's only Electrified Dagger bolt was resisted once the rage stream
shifted.

Goldens: every row of TestFury (Average-Default 320.36 to 325.59), TestArms (329.21 to 331.77),
TestProtectionWarrior (252.06 to 252.24) and TestFeralBear (177.43 to 178.44).

Measured on the references (10,000 iterations, empty trinket slots, production default rotation), against
the branch base `bde2043eab` (patch 112 moves Arms's Warblade of Caer Darrow, so the live pin reads 741.15):

| Reference | Board seed, board length | Board seed, 120 s | Alt seed, board length | Alt seed, 120 s |
| --- | ---: | ---: | ---: | ---: |
| arms-warrior | 743.93 → 747.20 (+0.44%) | 738.25 → 741.80 (+0.48%) | 742.42 → 746.72 (+0.58%) | 737.87 → 741.69 (+0.52%) |
| warrior (Fury) | 950.76 → 961.37 (+1.12%) | 950.44 → 962.87 (+1.31%) | 950.38 → 962.07 (+1.23%) | 948.46 → 959.21 (+1.13%) |
| protection-warrior | 479.21 → 480.68 (+0.31%) | 476.33 → 478.34 (+0.42%) | 478.61 → 480.11 (+0.31%) | 476.39 → 477.23 (+0.18%) |
| fury-protection-warrior | 516.86 → 518.92 (+0.40%) | 509.42 → 512.20 (+0.55%) | 516.81 → 519.89 (+0.60%) | 509.65 → 512.45 (+0.55%) |
| feral-bear-druid | 505.52 → 507.75 (+0.44%) | 499.52 → 502.54 (+0.60%) | 506.02 → 508.33 (+0.46%) | 499.29 → 502.81 (+0.71%) |
| feral-druid | 724.47 → 724.47 (+0.00%) | 723.61 → 723.61 (+0.00%) | 724.57 → 724.57 (+0.00%) | 724.20 → 724.20 (+0.00%) |

## 114. `buffs: Windfury Totem procs off the main hand only`

**Default true, the tooltip as the source of truth (decided 12 October 2026).** The 12 October review
turned it off for a few hours: this client marks a main-hand-only proc on the aura itself (Bloodthrill's
1289682 carries ATTR_EX_3_REQ_MAIN_HAND, Attributes[3] 0x400; Dual Wield Specialization 23584 marks its
off-hand proc with flag 0x800000 and Attributes[3] 0x1000000), and 10612 and 10610 carry neither and no
SpellEquippedItems row. Blizzard added "main hand" to all three tooltips in build 70009, when the totem became
the party aura (Era's read "Each hit"), so the tooltip is taken over the missing attribute until a combat log
of a dual-wielder shows otherwise. Yellow hits keep proccing: beta log 2708 (ElliotWood/Forever #644) has 16
of 62 Windfury Weapon procs off Stormstrike, on the same 0x14 mask.

Windfury Totem's party aura 10612 rolls 20% on ProcTypeMask 0x14, melee autos (0x4) and melee specials
(0x10), with a 100 ms internal cooldown. The mask names no hand and no SpellEquippedItems row restricts it,
so the engine let off-hand autos and specials roll it. Every Forever totem tooltip names the main hand: "Each
main hand hit has a 20% chance of granting the attacker 1 extra attack" (8512, 10613, 10614), as Era's
main-hand enchant 564 did, and players report off-hand procs as a bug (MythicSim Discord, 12 October 2026).

`WindfuryTotemMainHandOnly` in `sim/core/buffs/drivers.go` (default true) narrows the party trigger and the
shaman's own cast totem (`sim/shaman/totems.go`) to main-hand autos and specials. A shield strike still rolls
it: Shield Slam and Shield Bash are off-hand specials in the engine so a main-hand weapon's chance on hit
stays off them (upstream b6ce30c972), but their client rows carry no off-hand attribute (Attributes[3]
0x1000000, which the off-hand Windfury Weapon strike 439441 has), so the server swings them as main-hand
attacks; with no off-hand weapon, an off-hand special is one of those. False restores both hands.

Windfury Weapon is an enchant on the weapon (439431 through SpellItemEnchantment), so only that weapon's hits
roll it; Forever shamans cannot dual wield, so that is the main hand already and nothing changes.

Yellow hits are a separate switch, `WindfuryProcsFromSpecials` (default true, unchanged behaviour). Both
client rows include melee specials: 0x14 on the totem's 10612 and on Windfury Weapon's 439431. Against it is
only word of mouth: a WoW: Forever discussion that Windfury no longer procs off yellow hits and a player's
test in which Stormstrike did not trigger Windfury Weapon (Redfall, MythicSim Discord, 12 October 2026). No
log shows it either way at level 60 (beta log 2708, level 30, is the source of the two-strike model in
`weapon_imbues.go`). False leaves only auto attacks rolling the totem and the imbue.

Tests: `TestWindfuryTotemProcsOffMeleeHitsWithTheChargesTheTriggerLeaves` (`sim/warrior/dps/windfury_test.go`:
main-hand autos and specials proc, off-hand autos and specials with an off-hand weapon don't, a shield strike
does, either hand with `WindfuryTotemMainHandOnly` false, autos only with `WindfuryProcsFromSpecials` false)
and `TestWindfuryWeaponProcsOffSpecialsUnlessSwitchedOff` (`sim/shaman/enhancement/windfury_test.go`).

Goldens: TestFury (Average-Default 325.59 to 302.51), TestAssassination (333.58 to 315.49), TestCombat
(319.95 to 303.41), TestSubtlety (288.48 to 275.02), and the dual-wield rows of TestArms (29), TestSurvival
and TestSurvivalMelee (25 each) and TestProtectionWarrior (2).

Measured on the references (10,000 iterations, after patch 113). Every reference not listed is unchanged in
all four cells (Arms, the two Protection Warriors, Enhancement, Elemental with Melee buffs, both ferals and the
four paladins: one weapon or a shield):

| Reference | Board seed, board length | Board seed, 120 s | Alt seed, board length | Alt seed, 120 s |
| --- | ---: | ---: | ---: | ---: |
| warrior | 961.37 → 867.58 (-9.76%) | 962.87 → 866.48 (-10.01%) | 962.07 → 866.76 (-9.91%) | 959.21 → 866.28 (-9.69%) |
| assassination-rogue | 645.75 → 595.15 (-7.84%) | 645.93 → 596.05 (-7.72%) | 645.38 → 595.62 (-7.71%) | 645.69 → 595.40 (-7.79%) |
| combat-rogue | 682.26 → 640.21 (-6.16%) | 683.07 → 640.90 (-6.17%) | 683.13 → 639.97 (-6.32%) | 682.51 → 640.80 (-6.11%) |
| subtlety-rogue | 666.54 → 618.50 (-7.21%) | 666.45 → 618.57 (-7.18%) | 666.57 → 618.87 (-7.16%) | 666.56 → 619.02 (-7.13%) |
| survival-hunter | 841.87 → 781.36 (-7.19%) | 839.22 → 778.67 (-7.22%) | 840.88 → 781.66 (-7.04%) | 838.32 → 778.74 (-7.11%) |

`WindfuryProcsFromSpecials` false, on top of the default (not shipped, for the decision):

| Reference | Board seed, board length | Board seed, 120 s | Alt seed, board length | Alt seed, 120 s |
| --- | ---: | ---: | ---: | ---: |
| arms-warrior | 747.20 → 661.11 (-11.52%) | 741.80 → 656.75 (-11.47%) | 746.72 → 660.43 (-11.56%) | 741.69 → 655.59 (-11.61%) |
| warrior | 867.58 → 822.32 (-5.22%) | 866.48 → 822.24 (-5.11%) | 866.76 → 823.62 (-4.98%) | 866.28 → 821.79 (-5.14%) |
| protection-warrior | 480.68 → 459.63 (-4.38%) | 478.34 → 457.46 (-4.37%) | 480.11 → 459.58 (-4.28%) | 477.23 → 457.63 (-4.11%) |
| fury-protection-warrior | 518.92 → 493.35 (-4.93%) | 512.20 → 485.78 (-5.16%) | 519.89 → 493.12 (-5.15%) | 512.45 → 484.98 (-5.36%) |
| assassination-rogue | 595.15 → 583.76 (-1.91%) | 596.05 → 583.81 (-2.05%) | 595.62 → 583.70 (-2.00%) | 595.40 → 584.16 (-1.89%) |
| combat-rogue | 640.21 → 625.19 (-2.35%) | 640.90 → 625.96 (-2.33%) | 639.97 → 625.54 (-2.25%) | 640.80 → 625.77 (-2.35%) |
| subtlety-rogue | 618.50 → 592.27 (-4.24%) | 618.57 → 592.57 (-4.20%) | 618.87 → 592.61 (-4.24%) | 619.02 → 592.57 (-4.27%) |
| enhancement-shaman | 712.68 → 691.26 (-3.01%) | 708.55 → 688.84 (-2.78%) | 711.50 → 691.27 (-2.84%) | 709.63 → 688.60 (-2.96%) |
| feral-bear-druid | 507.75 → 430.16 (-15.28%) | 502.54 → 426.07 (-15.22%) | 508.33 → 430.33 (-15.34%) | 502.81 → 426.01 (-15.27%) |
| feral-druid | 724.47 → 707.64 (-2.32%) | 723.61 → 706.72 (-2.33%) | 724.57 → 707.10 (-2.41%) | 724.20 → 706.51 (-2.44%) |
| holy-protection-paladin | 563.30 → 547.71 (-2.77%) | 559.73 → 544.80 (-2.67%) | 563.44 → 547.67 (-2.80%) | 560.43 → 544.64 (-2.82%) |
| protection-paladin | 551.38 → 530.69 (-3.75%) | 547.86 → 527.17 (-3.78%) | 551.48 → 530.84 (-3.74%) | 546.95 → 526.60 (-3.72%) |
| ret-protection-paladin | 569.42 → 546.31 (-4.06%) | 563.74 → 541.48 (-3.95%) | 568.71 → 546.04 (-3.99%) | 564.07 → 541.41 (-4.02%) |
| retribution-paladin | 791.81 → 737.03 (-6.92%) | 789.44 → 736.31 (-6.73%) | 791.34 → 737.49 (-6.80%) | 789.60 → 735.79 (-6.81%) |
| shockadin-paladin | 927.14 → 894.24 (-3.55%) | 926.89 → 894.99 (-3.44%) | 927.40 → 894.60 (-3.54%) | 926.91 → 894.37 (-3.51%) |
| survival-hunter | 781.36 → 718.33 (-8.07%) | 778.67 → 715.69 (-8.09%) | 781.66 → 718.20 (-8.12%) | 778.74 → 715.44 (-8.13%) |

Elemental with Melee buffs does not move in either table: it swings a staff (the totem still procs off its
autos) and has no melee special.

## 115. `items: switches for Ironfoe's proc rate and hand` (default unchanged)

A player report (MythicSim Discord, 12 October 2026): Ironfoe procs about 6% of hits, twice Era's 0.8 PPM
(3.2% a swing on its 2.4 sec speed), and procs off off-hand hits. That is what the engine does, and it is what
the client says. In Forever Ironfoe (11684) no longer has a chance-on-hit spell: its ItemEffect is
TriggerType 1 (on equip) for Fury of Forgewright 1301046, an A_PROC_TRIGGER_SPELL aura for 15494 ("Grants 2
extra attacks on your next swing") with ProcChance 6, ProcsPerMinuteID 0, a 100 ms ProcCategoryRecovery and
ProcTypeMask 0x14 (melee autos and specials). No SpellEquippedItems row ties it to a weapon, so the off hand
and specials roll it too. Its tooltip: "Attacks have a chance to grant 2 extra attacks on your next swing.
Attacks against Orcs are 2 times as likely to activate this effect." Nothing in the client says main hand or
0.8 PPM; those are Era's chance-on-hit rules and the Classic sim's number.

So the default stays the client's. Two switches in `sim/common/classic/items_weapons.go` give the report's
reading for Sage to decide: `IronfoeMainHandOnly` (default false) drops the off-hand weapon's rolls (a shield
strike still rolls, as for Windfury Totem in patch 114), and `IronfoePPM` (default 0) replaces the 6% with a
rate per minute (0.8 for Era's).

Test: `TestIronfoeProcRateAndHands` (`sim/warrior/dps/ironfoe_test.go`): 6% on main-hand autos, off-hand
autos and main-hand specials; switched, 3.2% on main-hand autos and none off the off hand. Goldens: none move.

Measured on every reference that wields Ironfoe (10,000 iterations, after patch 114), main hand only at 6%:

| Reference | Board seed, board length | Board seed, 120 s | Alt seed, board length | Alt seed, 120 s |
| --- | ---: | ---: | ---: | ---: |
| warrior | 867.58 → 823.48 (-5.08%) | 866.48 → 823.80 (-4.93%) | 866.76 → 823.84 (-4.95%) | 866.28 → 822.90 (-5.01%) |
| protection-warrior | 480.68 → 480.68 (+0.00%) | 478.34 → 478.34 (+0.00%) | 480.11 → 480.11 (+0.00%) | 477.23 → 477.23 (+0.00%) |
| fury-protection-warrior | 518.92 → 518.92 (+0.00%) | 512.20 → 512.20 (+0.00%) | 519.89 → 519.89 (+0.00%) | 512.45 → 512.45 (+0.00%) |
| subtlety-rogue | 618.50 → 598.03 (-3.31%) | 618.57 → 597.18 (-3.46%) | 618.87 → 598.74 (-3.25%) | 619.02 → 597.97 (-3.40%) |
| protection-paladin | 551.38 → 551.38 (+0.00%) | 547.86 → 547.86 (+0.00%) | 551.48 → 551.48 (+0.00%) | 546.95 → 546.95 (+0.00%) |
| ret-protection-paladin | 569.42 → 569.42 (+0.00%) | 563.74 → 563.74 (+0.00%) | 568.71 → 568.71 (+0.00%) | 564.07 → 564.07 (+0.00%) |

Main hand only at 0.8 PPM:

| Reference | Board seed, board length | Board seed, 120 s | Alt seed, board length | Alt seed, 120 s |
| --- | ---: | ---: | ---: | ---: |
| warrior | 867.58 → 799.64 (-7.83%) | 866.48 → 799.18 (-7.77%) | 866.76 → 799.63 (-7.74%) | 866.28 → 798.86 (-7.78%) |
| protection-warrior | 480.68 → 465.53 (-3.15%) | 478.34 → 463.44 (-3.11%) | 480.11 → 465.77 (-2.99%) | 477.23 → 462.36 (-3.12%) |
| fury-protection-warrior | 518.92 → 502.54 (-3.16%) | 512.20 → 493.85 (-3.58%) | 519.89 → 502.58 (-3.33%) | 512.45 → 494.72 (-3.46%) |
| subtlety-rogue | 618.50 → 583.53 (-5.65%) | 618.57 → 583.17 (-5.72%) | 618.87 → 583.88 (-5.65%) | 619.02 → 583.60 (-5.72%) |
| protection-paladin | 551.38 → 535.16 (-2.94%) | 547.86 → 531.32 (-3.02%) | 551.48 → 535.83 (-2.84%) | 546.95 → 531.70 (-2.79%) |
| ret-protection-paladin | 569.42 → 553.27 (-2.84%) | 563.74 → 548.17 (-2.76%) | 568.71 → 552.99 (-2.76%) | 564.07 → 547.92 (-2.86%) |

## 116. `warrior: Deep Wounds from an off-hand crit uses the off-hand weapon`

`registerDeepWounds` sized every crit's bleed off the main hand's average weapon damage, an off-hand crit's
too. Forever's beta development notes: "Fixed a bug where Deep Wounds was scaling with attack power for your
Offhand weapon", and the developers' note "Deep Wounds does not scale with Attack Power but is a rolling DoT
compared to Classic/Vanilla". So the server sizes an off-hand crit's bleed off the off-hand weapon (the bug
was that its off-hand figure carried attack power, as the TrinityCore lineage's UNIT_FIELD_MINOFFHANDDAMAGE
does), and since the fix neither hand adds attack power.

`DeepWoundsUsesStrikingWeapon` in `sim/warrior/talents_arms.go` (default true): an off-hand crit (auto or
special, while dual wielding) takes its share of the off-hand weapon's average damage; a main-hand crit, a
shield strike and anything else keep the main hand. False sizes every crit off the main hand, as before.
Attack power is unchanged: the engine already leaves it out of both hands (beta logs 33/35 and 2705: a tick
is share x the weapon's average damage / 4). Whether the off-hand figure carries the off-hand 50% penalty is
not stated; the tooltip reads "your melee weapon's average damage", so it is the weapon's own average here.

Test: `TestDeepWoundsFromAnOffHandCritUsesTheOffHand` (`sim/warrior/dps/deep_wounds_test.go`). Goldens:
TestFury (244 of 245 rows, Average-Default 302.51 to 302.16) and the 28 dual-wield AllItems rows of TestArms.

Measured (10,000 iterations, after patch 115); only Fury dual wields:

| Reference | Board seed, board length | Board seed, 120 s | Alt seed, board length | Alt seed, 120 s |
| --- | ---: | ---: | ---: | ---: |
| warrior | 867.58 → 857.39 (-1.17%) | 866.48 → 856.38 (-1.17%) | 866.76 → 856.56 (-1.18%) | 866.28 → 856.16 (-1.17%) |
| arms-warrior | 747.20 → 747.20 (+0.00%) | 741.80 → 741.80 (+0.00%) | 746.72 → 746.72 (+0.00%) | 741.69 → 741.69 (+0.00%) |
| protection-warrior | 480.68 → 480.68 (+0.00%) | 478.34 → 478.34 (+0.00%) | 480.11 → 480.11 (+0.00%) | 477.23 → 477.23 (+0.00%) |
| fury-protection-warrior | 518.92 → 518.92 (+0.00%) | 512.20 → 512.20 (+0.00%) | 519.89 → 519.89 (+0.00%) | 512.45 → 512.45 (+0.00%) |

## 118. `hunter: a switch for pets pooling focus for Bite and their other cooldowns` (off)

Kragtor's Survival sim (MythicSim Discord, 12 October 2026, 6/10/35 with a Cat, 240 seconds) cast Bite
once. The Cat's rotation (`DefaultPetConfigs`, `pet.go`) casts Bite when it can and Claw otherwise, and
Claw (client 3009: 25 focus, no cooldown) fires the moment focus reaches 25. Focus regenerates 10 a
second, so after the opener it never climbs to Bite's 35 (17261: 35 focus, 10 sec category cooldown)
while Bite is off cooldown: the one Bite is the opener, off the full bar the pet starts with. The same
starves every pet ability that costs more than the pet's focus dump: the Raptor's Savage Rend (50), the
Gorilla's Thunderstomp (60), the Crab's Pinch (50). Upstream's `TestPetStrikes` says beta logs show the
same ("long-cooldown strikes land about once a fight, Claw spending the focus first"), which is the
game's autocast with both abilities on.

`PetPoolsFocusForCooldowns` (`sim/hunter/pet.go`, default **false**, so nothing changes) makes the pet
hold its focus dump while the focus left after it, plus what it regenerates before a higher priority
ability with a cooldown is ready, would not pay for that ability. Abilities with no cooldown (Dust Cloud)
or that could not be cast for another reason (Swipe on fewer than three targets, Furious Howl while it
is up) are not pooled for. It models a Hunter who manages Claw's autocast around Bite. Bite and Claw do
about the same damage per focus (90 for 35, 51 for 25), so it is worth little: on the 12 October
references (10,000 iterations, four cells) the Cat Bites about 11 times in 120 seconds instead of once
and the Hunter gains 0.5 to 2.7 DPS with a Cat or a Raptor, nothing with the other pets
(docs/forever-oct12-hunter.md in the site repo).

Adds `CurrentFocus` and `FocusRegenPerSecond` to the focus bar (`sim/core/focus.go`).

Test: `TestCatBitesOnCooldown` (`sim/hunter/pet_pool_test.go`): one Bite a fight unpooled, at least one
every 11 seconds pooled, and the pet still spends about 10 focus a second. Goldens: none move (off).

## 119. `hunter: the Wolf casts Furious Howl`

The Wolf only Bit; its Furious Howl, the reason to pick a Wolf for a melee Hunter (Kragtor, 12 October
2026), was not registered. In Forever it is an attack power aura on the party around the pet: client
row 24597 (rank 4) is `A_MOD_ATTACK_POWER` 82 plus 0.4 a level from 56, so 84 at 60, implicit target
20 (party within 15 yards), 60 sec duration, 30 sec cooldown, 60 focus. The 24 September beta notes cut
it ("Furious Howl Attack Power bonus reduced 40%"); the 70334 row is the cut value. It carries no
`SPELL_ATTR3_ONLY_TARGET_PLAYERS`, so unlike Battle Shout it reaches pets too.

The Wolf (`ExtraAbility: FuriousHowl`) howls whenever the buff is not up on it; the Hunter takes the
same aura while within 15 yards of the target the Wolf stands at. Melee attack power only, so it pays a
melee Hunter and the Wolf, not Auto Shot. Other party members are not simulated. Switch:
`WolfCastsFuriousHowl` (`sim/hunter/pet_abilities.go`, default true).

Test: `TestWolfFuriousHowl` (`sim/hunter/pet_pool_test.go`): the buff is up all fight on a melee Hunter,
never on one 30 yards out, and the melee Hunter does more damage than beside a Bite-only Wolf.
Goldens: none move (no golden tames a Wolf).

## 120. `hunter: Strider Kick speeds the Hunter up 30% for 3 sec`

Strider Kick (1317257) registered only its damage. Its third effect is `A_MOD_INCREASE_SPEED` 30 on the
caster for the spell's 3 sec, and the 24 September beta notes say so ("Strider Kick now increases
movement speed by 30% for 3 seconds."). The kick now activates a "Strider Kick" aura (action ID 1317257,
so a rotation can check `auraIsActive` on it) with an active movement speed effect of 30%. It only matters
to a rotation that moves after a kick, such as a melee weave stepping out to shoot. Switch:
`StriderKickSpeedsUp` (`sim/hunter/strider_kick.go`, default true).

Test: `TestStriderKickSpeedsUp` (`sim/hunter/pet_pool_test.go`). Goldens: with patch 121 none move. On
its own this patch moved `TestSurvival` (176 of 193 rows, -4.4 to +3.6, Average 352.97 to 352.96): that
suite walks in from 30 yards, the first kick lands at 5 yards mid-walk, and the speed change hit the
movement bug patch 121 fixes.

## 121. `core: a speed change mid-move keeps the move's destination`

`MultiplyMovementSpeed` re-plans a move in progress when the unit's speed changes, but passed `MoveTo`
the move's length (speed times planned duration) where `MoveTo` takes the distance from the target to
end at. A Hunter walking from 5 to 12 yards who sped up mid-walk stopped at 7, inside the Auto Shot dead
zone; one walking in from 30 to 0 went to -30. Nothing in the fork changed speed mid-move until patch
120, so it never showed. The re-plan now ends where the move was going (start plus the move's distance).
Switch: `MidMoveSpeedChangeKeepsDestination` (`sim/core/movement.go`, default true).

Effect, Survival reference with the Experimental melee weave (patches 118 to 120, 10,000 iterations,
board seed): 692.8 before patch 120, 664.7 with it, back to the value in docs/forever-oct12-hunter.md
(site repo) with this patch.

Test: `TestSpeedChangeMidMoveKeepsDestination` (`sim/core/movement_speed_test.go`): a walk from 5 to 12
yards sped up 30% half a second in stops at 12, and at 7 with the switch off. Goldens: `TestSurvival`
returns to its pre-120 values; the full `sim/...` suite passes.

## 125. `core: a switch for the boss's 1% hit suppression`

No behaviour change by default. Against a level 63 target the attack table subtracts
`HitSuppression` (1%) from melee and ranged hit before the 8% base miss, so a single-handed yellow
attack needs 9% hit to never miss (`NewAttackTable`, `sim/core/target.go`). Magey's Forever warrior
test (github.com/magey/forever-warrior issue 1, 4,259 attacks, 25.08% observed against 24.78%
expected without it) concludes the suppression is gone in Forever and the yellow cap is 8%; the 1.8%
crit suppression remains. Đegausser raised it on the MythicSim Discord (#bugs, 12 October 2026). The
developers have not confirmed it, so the default keeps the beta's 9%.

`Encounter.no_hit_suppression` (`proto/common.proto`, field 13; `Encounter.NoHitSuppression` in
`sim/core/target.go`) is the switch, read per request so the site can offer it as an Advanced choice:
true zeroes `HitSuppression` on every attack table against an enemy (`setupAttackTables`,
`sim/core/environment.go`). The dual-wield 19% penalty is unchanged. Spells never had a suppression:
the spell table is a 17% base miss against a level 63 target with spell hit counted in full (16% at
most), and the switch does not touch it.

Test: `TestNoHitSuppressionSwitch` (`sim/core/hit_suppression_test.go`): 8% hit leaves a 1% miss with
the suppression and none without; the spell base miss stays 17%. Goldens: none move (no test sets the
field).

Measured on the physical reference builds (MythicSim site, 12 October 2026, 10,000 iterations, board
and alt seed at the board length and 120 s, switch off against on, same binary): +0.45% to +1.46% for
the references below the 9% cap or dual wielding (Protection Warrior +1.46%, Bear +0.98%, Subtlety
+0.98%, Warrior +0.95%, Marksmanship +0.89%); 0.00% for the four already at 9% with one weapon
(Arms, Feral, Enhancement, Protection Paladin). Details in the site's `docs/forever-oct12-gear.md`.

Drop the switch (and make it the only behaviour) once Blizzard confirms the 8% cap.

## 126. `items: Frost Oil as a weapon imbue`

Frost Oil (item 3829, use spell 3595) puts temporary enchantment 26 on a melee weapon, and the engine
registered nothing for it: an imbue ID of 3595 or 3829 did nothing. The client's enchant row casts
350025, a Frostbolt that rolls 33 + 0.7 a level (levels 8 to 12, 12% variance) with 0.269 of spell
power, at the 10% chance the client states (`EnchantChances` 350025: 10% for enchant 26; tooltip "a
10% chance of casting Frostbolt at the opponent when it hits"). It is a flat chance on each landed hit
of the oiled weapon, not a rate per minute. MythicSim offers it beside the stones for the Shockadin
(Redfall, MythicSim Discord, 12 October 2026), whose spell power makes the bolt worth ~200 a hit.

`registerStaticImbue` now takes the hand's mask and runs `core.StaticImbueProcs[imbueID]`, a hook
`sim/common/forever/frost_oil.go` fills for 3595 and 3829: a weapon proc (`IsWeaponProc`) on that hand's
landed hits that casts 350025 through `shared.SpellDataProcDamageSpell`. Like every other weapon proc
it also rolls on a Seal of Righteousness hit (patch notes in `seal_of_righteousness.go`): the Shockadin
reference casts about 15.8 bolts in 120 seconds, 10% of its swings and seal hits.

Switch: `forever.EnableFrostOil` (default on). Test: `TestFrostOilProcsOffItsHand` (`sim/rogue`): 1,306
bolts off 13,322 main-hand hits (9.8%), none off the off hand, none with the switch off. Goldens: none
move (no default carries Frost Oil).

| Shockadin reference (10,000 iterations) | Elemental Sharpening Stone | Frost Oil |
| --- | ---: | ---: |
| board seed, board length | 927.14 | 942.68 |
| board seed, 120 s | 926.89 | 941.65 |
| alt seed, board length | 927.40 | 942.27 |
| alt seed, 120 s | 926.91 | 941.21 |

Drop it when the store registers temporary enchantment combat procs itself.

## 127. `core: Party.battle_squawks, Battle Squawk stacks from the party's engineers`

Battle Squawk (23060) is the buff a Gnomish Battle Chicken (item 10725, summon 13166) gives its party:
"Attack speed increased by 5%" for 4 minutes. Forever's client states it as melee haste (SpellEffect
699978, client 1.60.1.69893: aura 319, 5 points; Era 1.15.9 had aura 138, 4 + 1 die), the aura MCP's
Haste (13494) uses. Several engineers' chickens stack, one aura each, which multiply as separate
haste auras do. The engine had no way to assume it.

`Party` (api.proto, hand-written, not the buff manifest) gains `int32 battle_squawks = 3`, 0 to 5.
`registerBattleSquawk` (`sim/core/battle_squawk.go`) gives each of the party's players a "Battle Squawk"
aura (23060) at that many stacks from the pull, `MultiplyMeleeSpeed(1.05)` a stack, lasting its 4
minutes: a fight longer than that runs its end without it. It is off unless a request sets the field.

Switch: `core.EnableBattleSquawk` (default on). Test: `TestBattleSquawkStacksMeleeSpeed` (`sim/rogue`): 1
and 5 stacks give 1.05x and 1.05^5 the main-hand swings over 120 seconds, 9 stacks are held to 5, an
8 minute fight gets the bonus for half of it, and the switch off changes nothing. Goldens: none move.

| 10,000 iterations, 5 stacks | Without | With |
| --- | ---: | ---: |
| Shockadin, board seed, board length | 927.14 | 1033.43 |
| Shockadin, board seed, 120 s | 926.89 | 1024.56 |
| Shockadin, alt seed, board length | 927.40 | 1034.05 |
| Shockadin, alt seed, 120 s | 926.91 | 1025.50 |
| Fury Warrior, board seed, board length | 950.76 | 1126.78 |
| Fury Warrior, board seed, 120 s | 950.44 | 1124.86 |
| Fury Warrior, alt seed, board length | 950.38 | 1128.12 |
| Fury Warrior, alt seed, 120 s | 948.46 | 1125.59 |

Drop it when upstream carries a Battle Squawk party option.

## 128. `core: a crit's extra rage leaves out the two-hander's factor`

Patch 113 doubled a critical auto attack's whole rage. Magey's beta logs on client 70170 (after the 1
October notes; github.com/magey/forever-warrior/issues/3) fit `crit = swing + 3.46 x speed x hand`: the
extra leaves out the two-hander's 1.3 (`TwoHandRageMultiplier`), so a one-hander doubles and a two-hander gets
about 1.77 times its swing (2H warrior bonus 11.3 to 11.5 measured, 11.42 predicted). Yellow crits add
nothing, as the engine already has it. The client's hook, Rule of Rage (DND) 1322574, is an A_DUMMY of 100
on melee auto attacks (mask 0x4).

`CritRageBonusSkipsTwoHand` in `sim/core/rage.go` (default true) pays the crit's extra on the swing's rage
before the two-hand factor. False multiplies the whole swing by `CritRageMultiplier`, as patch 113 did.
Bear Form's weapon is normalized to 2.5 and never gets the two-hand factor in the engine, so the Bear is
unchanged. (Magey measured the bear swing at 3.46 x 2.5 x 1.3, about 23% above the engine's; not changed
here.)

Tests: `TestCritAutoAttackRageIsTwiceAHit` (`sim/core/rage_crit_test.go`, two-hand case 11.7 to 20.696) and
`TestCritRageBonusOnTheWholeTwoHandSwingWhenSwitchedOff`, `TestCritAutoAttackGivesATwoHanderItsSwingPlusTheOneHandRage`
(`sim/warrior/dps/crit_rage_test.go`). Goldens: the two-hand rows of TestArms, TestFury and
TestProtectionWarrior.

## 129. `paladin: Seal of Righteousness takes its coefficient from the weapon in the hand`

Seal of Righteousness's hit is the value per hundred times the swing speed times 0.85 for a one-hander
or 1.2 for a two-hander, plus spell power at the judgement dummy's 0.185, times 1.2 for a two-hander
(see the comment in `seal_of_righteousness.go`). The weapon-speed part read the main hand at each hit;
the spell power share was fixed when the seal registered, from the weapon worn at the pull. A fight
that swaps a two-hander for a one-hander (MythicSim's Manual Crowd Pummeler opener: MCP's Haste for
30 s, then the player's own weapon through `Player.item_swap`) kept the two-hander's 0.222 on the
one-hander for the rest of the fight, and the reverse kept 0.185 on a two-hander.

`SealOfRighteousnessFollowsSwappedWeapon` in `sim/paladin/seal_of_righteousness.go` (default true) sets
the coefficient from the main hand at each hit. False keeps the registration-time value. Nothing
changes for a fight without a swap.

Checked for the same problem across a two-hand/one-hand swap, nothing else needed: the paladin and
warrior weapon specializations (`applyWeaponSpecialization`, `registerTwoHandedWeaponSpecialization`,
Weaponmaster) and the racial weapon specializations re-read the weapon in an item swap callback; rage
reads `AutoAttacks.MH()`/`OH()` and the weapon's normalized speed at each swing (`core/rage.go`);
Seal of Command's PPM manager (`NewLegacyPPMManager`) rebuilds on a swap and its hit reads
`MHWeaponDamage`; Hammer of the Righteous and Holy Strike read the main hand at the hit; Windfury
Totem's proc is per attack; weapon-effect and enchant PPMs (`NewDynamicLegacyProcFor*`) rebuild on a
swap; the shaman imbues follow `imbueMhSwap`. Not fixed, so MythicSim only offers the opener to a
main hand with nothing in the off hand: Whirlwind registers its off-hand strike only when the off hand
is equipped at the pull (`sim/warrior/whirlwind.go`), so a dual wielder who opened with a two-hander
would never get it back.

Test: `TestSealOfRighteousnessFollowsASwappedWeapon` (`sim/paladin/retribution`): with 1000 spell
power, Arcanite Reaper swapped for Mass of McGowan at the pull hits for 289.3 a non-critical hit
against 289.1 for Mass of McGowan worn all fight; switched off, 336.9. Goldens: none move (no suite
swaps weapons).

## Note on 127 (Battle Squawk), 12 October review

In Forever, 23060's effect targets the pet's master and the caster (ImplicitTarget 27 and 1, builds 69893
and 70338); Era's targeted the caster's party (20). A Battle Chicken buffs only its engineer, so one stack on
the engineer is what the game can give. `Party.battle_squawks` stays 0 to 5 in the engine; MythicSim offers
one stack on the player.


## 130. `warrior: Last Stand reads its 30% from client 70291's health aura`

Client 1.60.1.70291 moved Last Stand's buff (12976) from `A_MOD_MAX_HEALTH` (230) to
`A_MOD_INCREASE_HEALTH` (34), still with 30, and the buff tooltip still reads "Health increased by
30% of maximum". `registerLastStand` asked for aura 230 as the buff was gained, so from `c91f22924a`
(the fork's merge of client 70338) every fight where a warrior gains Last Stand panicked with "spell
12976 has no effect with aura 230". Five of the Rust engine's 49 Protection Warrior fixtures reach it;
the reference builds do not.

`registerLastStand` now reads the share once at registration, from aura 230 when the client states it
and from aura 34 otherwise, and the buff grants 30% of maximum health as before.

Test: `TestLastStandGrantsThirtyPercentOfMaximumHealth` (`sim/warrior/protection`). Goldens: none move
(the fights that reached Last Stand panicked).

## 131. `buffs: party auras without ONLY_TARGET_PLAYERS reach the party's pets`

A pet got no party or raid buff at all: `Pet.AddPartyBuffs` is empty and `applyBuffEffects` runs on
players only. In Forever a party aura is an area aura on the party (SpellEffect type 35), and the
server puts it on the party's pets too unless the spell carries Attributes[3] 0x100
(`ATTR_EX_3_ONLY_TARGET_PLAYERS`). Beta logs show the rule: in foreverlogs.gg report 2695 (2 October)
a Warlock's Imp and a Hunter's pet carry Retribution Aura (7294, no flag) 90% and 96% of the run,
while Devotion Aura (643, flagged) is on the players 83% to 90% and never on a pet; in report 2748
(8 October) two pets carry Retribution Aura and none of its 35 pets carries Devotion Aura.

Flagged in client 70338, so nothing changes for a pet: Arcane Brilliance (23028), Prayer of Spirit
(27681), Gift of the Wild (21850), Prayer of Fortitude (21564), the Greater Blessings (25898, 25916,
25918), Devotion Aura, Battle Shout, Leader of the Pack, Moonkin Aura, Strength of Earth and Grace of
Air. Not flagged, now on a summoned pet (`sim/core/buffs/pet_buffs.go`): Mana Spring Totem (10494,
10 mana every 2 sec, kept while casting), Trueshot Aura and the four Atiesh auras. Left out: Mana
Tide Totem (the sim casts it as the player's external cooldown), Retribution Aura (pets take no
hits), Blood Pact (health only; the casting Imp does not carry its own, report 2695) and the weapon
totems. Guardians are not party members and get nothing.

`PartyAurasReachPets` (default true). The core hook is `BuffHooks.ApplyPetBuffs`, called from
`Raid.applyCharacterEffects` for each non-guardian pet.

Effect: an Imp with Mana Spring Totem regains about 440 mana a fight while casting, which it did
not before. Demonology Warlock reference with the Imp kept out instead of the Succubus (Discord,
14 October): +0.6% to +0.7% (Imp +4.3 to +4.6 DPS), 10,000 iterations, every cell. No MythicSim
reference moves (their Succubus is never short of mana; no reference has a pet beside Trueshot
Aura). Goldens: `TestAffliction` (Succubus with Mana Spring Totem) +0.5%.

Tests: `TestPartyAuraPetFlags` and `TestImpGetsManaSpringTotemOnly` (`sim/warlock`).

At integration onto `d2409ccd0` (14 October), `TestShadowPriest` and `TestSmitePriest` move too: the Shadowfiend is a summoned pet, so it now carries Trueshot Aura and Mana Spring Totem from the test raid. Only DPS lines change (naked Shadow +1.3% short fights, +0.4% long).

## 132. `hunter: Lacerating Strikes rolls a refreshed bleed over and keeps its tick timer`

Lacerating Strikes (1310533) bleeds the target for 40% of the Mongoose Bite over 21 sec (1310536,
seven 3 sec ticks). Patch 87 let a new bite replace the running bleed. Mongoose Bite comes round every
5 to 13 sec, so a bleed was replaced after 3 or 4 of its 7 ticks: on Sanctum's weave
(`1a573b3d`, 3,000 iterations) 8.99 bleeds a fight dealt 28.2 ticks, 3.1 a bleed, 45% of the 40%,
and the report showed 17.6 DPS where 40% of every bite would be about 34 (raised by Sanctum on the
MythicSim Discord, 13 October). The client cannot say which the server does. Forever's Deep Wounds,
the other "X% as a bleed" effect, rolls over and keeps its tick timer (Zirene on forever-bugs #234,
patch 89), so Lacerating Strikes now does the same: `procLaceratingStrikes` reads `OutstandingDmg`
before the cast, the bleed is applied with `ApplyKeepingTickTimer`, and the owed damage and the new
bite's 40% are spread over `RemainingTicks`. No log at level 60 has settled it yet.

`LaceratingStrikesRollsOver` (default true) switches it; false is the old replace and tick restart.

`TestLaceratingStrikesRefreshRollsOver` (replaces `TestLaceratingStrikesRefreshKeepsTheNewBleed`)
bites for 1000 at 0 sec and 1400 at 4.5 sec: rolled over, the next tick stays at 6 sec, the bleed runs
out at 25.5 sec with 7 ticks of (6 x 400/7 + 560) / 7 and ticks at 3, 6, ..., 24 sec; switched off, 7
ticks of 560/7 at 7.5 to 25.5 sec. Goldens (Average-Default): TestSurvival 352.97 to 359.47,
TestSurvivalMelee 397.50 to 404.01; Beast Mastery and Marksmanship do not move.

10,000 iterations, live pin `944a791060` against this patch:

| Build | Before | After | Lacerating Strikes |
| --- | ---: | ---: | ---: |
| survival-hunter reference, board seed, board length | 781.36 | 804.10 | 16.8 to 38.7 DPS |
| survival-hunter reference, board seed, 120 s | 778.67 | 800.31 | 16.8 to 38.6 DPS |
| survival-hunter reference, alt seed, board length | 781.66 | 803.60 | 16.8 to 38.8 DPS |
| survival-hunter reference, alt seed, 120 s | 778.74 | 800.49 | 16.8 to 38.6 DPS |
| Sanctum's weave `1a573b3d` (seed 1179722310) | 900.54 | 914.97 | 17.6 to 32.4 DPS |
| Sanctum's weave `f7751f61` (seed 1179722310) | 945.88 | 960.74 | 19.0 to 34.3 DPS |

Drop this when upstream rolls the bleed over, or set the switch false if a log shows a refresh
replaces it.


## 133. `shaman, core: Rockbiter, Flametongue and the main-hand stone follow the weapon in the hand`

With an item swap (MythicSim's Manual Crowd Pummeler opener: the Pummeler worn at the pull, the
player's own weapon swapped in after 30 seconds), only Windfury and Frostbrand Weapon followed the
weapon in the hand. Rockbiter's attack power aura was built only for the weapon worn at the pull and
was permanent: Rockbiter on the Pummeler alone counted all fight, on the own weapon alone never.
Flametongue's proc aura was built only for the weapon worn at the pull, so Flametongue on the own
weapon alone never procced. The main-hand stone or oil is a per-hand consumable
(`Consumables.MhImbueId`): an Elemental Sharpening Stone or a wizard oil is a character stat for the
whole fight, Frost Oil a main-hand proc for the whole fight, and a Dense stone is added to the swing
built at the pull and lost at the first reset when the main hand swaps.

- `RegisterRockbiterImbue` (`sim/shaman/weapon_imbues.go`): with an item swap where only one of the
  two sets carries Rockbiter (`imbueMh` or `imbueMhSwap`), the aura is no longer permanent; an item
  swap callback keeps it up while a weapon with enchant 3021 is in the hand (the item swap's reset
  sets it for the pull's weapons), and Spirit Weapons' x1.86 threat moves into its OnGain/OnExpire.
  Both sets with Rockbiter, or no item swap: permanent as before.
- `RegisterFlametongueImbue`: `imbueMhSwap` (or `imbueOhSwap`) Flametongue without it on the pull's
  weapon now builds the proc aura on the slot without marking the pull's weapon; the existing swap
  callback in `makeFTProcTriggerAura` switches it off at the reset and on at the swap.
- `ConsumesSpec.mh_imbue_swap_id` (field 29, `sim/core/imbue_swap.go`): set, `mhImbue_id` counts only
  while the pull's weapon is in the main hand and `mh_imbue_swap_id` only while the swap's weapon is;
  -1 is nothing on the swap's weapon, the same id on both counts all fight. Elemental Sharpening
  Stone, Brilliant Wizard Oil, Wizard Oil and Mana Oil become a stat aura, Frost Oil (and any
  `StaticImbueProcs` entry) its proc aura switched on and off, and a Dense stone's 8 damage is added to
  the main-hand swing rebuilt at each swap. Blessed Wizard Oil and Consecrated Sharpening Stone (undead
  only) still count all fight. A Feral Druid's paw keeps `MHImbueFlatWeaponDamage` as before. Unset (0),
  or no main-hand swap weapon, everything is as before.

Switches: `shaman.PerWeaponImbues` and `core.EnablePerWeaponMainHandImbue` (default true).

Effect, MythicSim request path (`claude/forever-oct14-integration`, defaults with the opener, the
swap fields written into the request by hand), 10,000 iterations, board seed and alt seed, board
length and 120 s, against today's Enhancement default (Rockbiter on both, Frost Oil):

| Enhancement | board seed, board | board seed, 120 s | alt seed, board | alt seed, 120 s |
| --- | ---: | ---: | ---: | ---: |
| Rockbiter both, Frost Oil (default) | 730.71 | 727.18 | 729.70 | 726.97 |
| Pummeler Flametongue, own Rockbiter, Frost Oil both | 731.42 (+0.10%) | 728.98 (+0.25%) | 730.37 (+0.09%) | 728.04 (+0.15%) |
| Pummeler Flametongue + Frost Oil, own Rockbiter + Elemental Sharpening Stone | 732.10 (+0.19%) | 730.10 (+0.40%) | 731.79 (+0.29%) | 728.99 (+0.28%) |
| Pummeler Flametongue + Frost Oil, own Rockbiter + Brilliant Wizard Oil | 733.26 (+0.35%) | 730.97 (+0.52%) | 732.35 (+0.36%) | 729.68 (+0.37%) |

The last split also wins at 300 s (+0.17%, +0.19%) and on three targets (+0.67% to +0.73%). Elehance
(two and three targets, Flametongue and Brilliant Wizard Oil with the opener) gains from no split:
Frost Oil on the Pummeler ties on two targets (-0.04% to +0.04%) and loses on three (-0.15% to
-0.27%); Rockbiter on either weapon loses 0.7% to 1.9% (3,000 iteration screen).

Tests: `TestRockbiterFollowsTheWeaponInTheHand`, `TestRockbiterThreatFollowsTheWeaponInTheHand`,
`TestFlametongueOnTheSwapWeaponProcsOnlyAfterTheSwap` and `TestStoneAndOilFollowTheWeaponInTheHand`
(`sim/shaman/enhancement`). Goldens: none move (no suite swaps weapons); MythicSim's default
Enhancement and Elehance requests give the same DPS to the hundredth of a point as `944a791060`.

## 134. `items: Forever's revised items take their client equip spells' hit and crit`

Finding B9 of MythicSim's 14 October review: the items `mergeForeverSimDB` fills in from master's
database (`forever_sim_db.json`, the Classic-era items client 69893's `ItemSparse` does not ship) kept
Classic's hit and crit, stated as rating, even where Forever revised the item. Their client
`ItemXItemEffect` rows link new equip spells, but `dbc.MergeItemEffectsForAllStates` skipped the equip
spells of every filled-in item, as master's stats already include Classic's. `UnifyGearHitAndCrit`
then paid the Classic rating into every pool.

Seven filled-in items have a client equip spell that states hit or crit (the client spells are in the
70338 spell store's item roots, so the 70338 item tables still link them):

| Item | Client equip spell | Before (master, unified) | After |
| --- | --- | --- | --- |
| Blackhand's Breadth 13965 | 1318954 "Increased Critical 1 - Melee and Ranged", aura 52, 1 | 28 crit rating: 2% melee, ranged and spell | 1% melee and ranged crit |
| Eye of the Beast 13968 | 1318956 "Increased Critical 1 - Spell", aura 57, 1 | 28 spell crit rating: 2% to every pool | 1% spell crit |
| Trueaim Gauntlets 13255 | 1301127 "Increased Hit Chance - Ranged - 03", aura 54, 3, bows, guns and crossbows | 10 hit rating: 1% to every pool | 3% ranged hit |
| Sash of the Grand Hunt 22207, Stoneshatter 18388, Marksman Bands 18296, Infernal Trickster Leggings 17754 | 1294769 "Increased Hit Chance - Ranged - 01", aura 54, 1, the same weapons | none | 1% ranged hit |

**Type.** The client keeps two kinds of hit and crit apart. Item rows state them only as the generic
ratings, `ItemSparse` 31 and 32 (184 and 521 of client 69893's rows; none of the typed ratings 16 to 21
occurs), and the Forever planner lists them as generic `hitpct` and `critstrkpct`; that is the rule
`UnifyGearHitAndCrit` models, and it stays. The revised equip spells name their type, in the spell name
("- Melee and Ranged", "- Spell", "- Ranged"), the description ("with Melee and Ranged attacks", "with
Spells", "with ranged weapons"), the aura and the equipped-item class. So these follow the spell:
`dbc.AddEquipSpellStats` already turns the auras into typed percent pseudo stats (the path client items
and enchants take: Biznick's Accurascope, Presence of Sight), and the sim does not unify pseudo stats.
`sim/core/forever_rules.go` is unchanged.

`dbc.TakeClientHitAndCrit` (called from `MergeItemEffectsForAllStates` for a filled-in item): where an
equip spell states hit, the item's hit rating (melee and spell) is dropped and the spell's hit pseudo
stats are added; crit the same way. Classic states all hit and crit on equip spells, never in the item
row, so the dropped rating is the Classic spell's. Any other stat a filled-in item's equip spell states
is still skipped. Items whose client spells state no hit or crit are unchanged.

Separately, `fillPlannerLegendaryStats` (`tools/database/gen_db/planner_legendary_stats.go`, after
`fillPlannerSpellPower`): `RandPropPoints` has no legendary column, so `GetScaledStat` priced every stat
of a legendary client row at 0, and the epic column does not fit (Atiesh's allocations come out at a
budget of 59 at item level 90, epic 64; Sulfuras' at 54 at item level 80, epic 57). The Forever planner
states the priced stats, each the client allocation at that budget (Atiesh, mage: 5254, 5424 and 4068 of
59 are 31 Stamina, 32 Intellect and 24 Spirit, 3390 is 20 hit rating). A legendary item takes its primary
stats, resistances and generic hit and crit from the planner where it has none:

| Item | Added |
| --- | --- |
| Atiesh, mage 22589 | 31 Sta, 32 Int, 24 Spi, 20 hit rating (2%, `ItemSparse` 31) |
| Atiesh, warlock 22630 | 30 Sta, 29 Int, 28 crit rating (2%, `ItemSparse` 32) |
| Atiesh, priest 22631 and druid 22632 | 28 Sta, 28 Int, 27 Spi |
| Thunderfury 19019 | 5 Agi, 8 Sta, 8 Fire and 9 Nature resistance |
| Sulfuras 17182 | 12 Str, 12 Sta, 30 Fire resistance |

The mage staff's planner `splcritstrkpct` 2 is Power of the Guardian (28142), a party aura the item
already carries as an effect, and is not read. Ashbringer and Shard of the Defiler are not in the
planner and are unchanged.

Not changed, for want of client data: Skyfury Helm 20134 (Classic 1% crit and 2% spell crit, unified to
3%; no client row or effect), the 21 reference-build items the planner does not list (Truestrike
Shoulders, Tarnished Elven Ring, Mark of Fordring, Clutch of Andros, Heliotrope Cloak and the rest; none
has a client `ItemXItemEffect` row), and 11 filled-in items whose client equip spells state something else
(Savage Gladiator Chain, Kindling Stave and Death Grips among them; the planner keeps their crit and hit,
1.5%, 1% and 0.5%).

`assets/database/db.json` carries the 13 edits, made by the two functions run on the committed JSON (the
client spells' pseudo stats read off client 69893's `ItemXItemEffect` and `ItemEffect` in
`tools/data_watch/.cache` and the spell store's effects and equipped classes), each edited item line
re-marshalled the way `WriteProtoArrayToBuffer` writes it, so no other record moves. `go run
./tools/sync_db_binary` rebuilt `db.bin` (run on the unedited JSON first, it reproduced the old binary
byte for byte). `leftover_db` has no affected item. A full `make db` was not run (no
`tools/database/wowsims.db`).

Tests: `TestTakeClientHitAndCrit` (`tools/database/dbc`), `TestFillPlannerLegendaryStats`
(`tools/database/gen_db`), `TestForeverRevisedItemsTakeClientHitAndCrit` and
`TestForeverLegendaryItemStats` (`sim/core`, with_db, on the embedded binary; both fail on the old one).
Goldens: AllItems rows only. Sulfuras gains its Strength and Stamina in `TestFeralBear` (+2.1%),
`TestFeralCat` (+1.0%), `TestRetribution` (+0.9%), `TestEnhancement` (+0.9%) and `TestArms` (+1.2%);
Thunderfury its Agility and Stamina in `TestBeastMastery`, `TestMarksmanship`, `TestSurvival` and
`TestSurvivalMelee` (+0.3% to +0.5%) and `TestArms` (-0.4%: the Agility moves the random stream, and the suite's
few iterations read that as a loss). None of the seven revised items is in an AllItems list.

Measured (MythicSim's 31 reference builds, board seed, board length, 10,000 iterations, the worker's
default requests, base `d2409ccd0` with the old binary): every reference is bit-identical, as none wears
an affected item. On the trinket board, Blackhand's Breadth loses 0.8% to 1.2% (Fury Warrior 908.0 to
897.9, Survival weave 951.0 to 941.5, Enhancement 744.8 to 735.5, Marksmanship 797.1 to 789.1,
Retribution 915.5 to 907.4, Feral 751.6 to 745.8, Combat 654.6 to 648.9, Protection Warrior 574.6 to
570.1), about half its value over no trinket. Eye of the Beast loses 0.8% to 1.0% (Fire 726.1 to 719.4,
Destruction 726.2 to 719.3, Balance 664.2 to 657.8, Shadow 725.1 to 719.1, Elemental 524.8 to 520.5).

Downstream: the items' hit and crit now live in `pseudoStats`, which a gear summary that reads only
`stats` (MythicSim's `web/src/lib/forever-gear-stats.ts`) does not show.

Drop it when the client ships these items' `ItemSparse` rows (the generator then reads their equip
spells as for any client item) and `RandPropPoints` prices legendary items.
