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
| `attack.go` | Swing retry (`ForeverSwingRetryBeat`): a main hand that is ready while out of melee swings on the first half-second step at or after the moment the server sees you arrive (arrival + 50 ms); the steps count from the timer's end or the melee attack turning on, whichever is later (`retrySwingAt`, `swingOnArrival`). Arriving before the timer ends swings at its end. Walking out leaves the timer alone, and `ForceMeleeSwing` does nothing with the retry on. Measured with Arc's swing probe and 13 beta combat logs (TestSwingRetryOnArrival). Each weapon keeps `lastAttackAt` (the last real swing or Auto Shot) for the rotation value `autoTimeSinceLast`. Why a swing or shot waited is logged (`Swing wait <reason> from A to B`: range, retry, left, close, moving, windup, shotdelay) for the timeline (TestSwingWaitLogs). With the option `melee_resets_auto_shot` (the beta bug, off by default) a white swing or a queued ability restarts the Auto Shot timer once, from the swing. When the melee attack turns on, the main hand swings when ready and the off hand waits half its swing time (`OffHandWait`, `OffHandWaitsForAttackStart`); this only applies if an Auto Shot turned the melee attack off. A walk into melee starts swinging at 5 yd, not at arrival (`StartMeleeMidMove`). The walk in turns the melee attack on as it starts, or at an earlier Start Attack action (`StartAttack`, `attackOnForWalkIn`); an Auto Shot turns it off again. With the rotation switch `manual_start_attack` only Start Attack turns it on, and without one the arrival does: the main hand swings as the server sees you (arrival + 50 ms, no retry), the off hand half its swing time after arrival (TestStartAttack). Off by default; every shipped build comes out byte for byte the same. `RangedCycleCut` shortens each Auto Shot timer (the Auto Shot spam bug option). Normalized instants (Strider Kick, Mongoose Bite) take the plain normalized speed, slow weapons too (Arc's tests: a 3.8 two-hander's kicks land at 3.3; a 2.8 and a 1.7 main hand in one session hit 24.7 / 24.0 apart, their base gap, not the 33 a kept 2.8 adds; TestNormalizedStrikeSpeed). The earlier "higher of the two" reading was wrong; knob AB_NORM_MAX=1 brings it back for comparisons. Every swing timer start and haste move is logged (`startTimer`, `moveReadyAt`, under the weapon's own id from `timerID`) for the timeline's swing timer bars; results do not change. The rotation is looked at again right after each Auto Shot, like after a melee swing, so a move into melee can start during a GCD (TestMoveAfterShotDuringGcd); the builds not on Arc Weave come out the same with or without it. The pull (`startPull`) leaves the attacks to a walk already under way: the first Auto Shot fires before the pull and can start a walk in (TestMoveAtPull). |
| `movement.go` | Auto Shot fires 0.6 s after you stop moving (`AutoShotStandStill`; Arc's logs 0.57-0.73 s). Base run speed comes from `Unit.BaseRunSpeed` (0 means 7 yd/s). Speed bonuses normally compete (strongest counts); `MultiplyStackingMoveSpeed` is for ones that stack. `Unit.MoveDelay` is the move reaction time, a wait before each step into or out of melee (the site's builds and default use 0.1 s); `PendingMove` stops a rotation from issuing a second move during it, and no spell with a cast time can start during it, since the walk would cancel it (`castBlockedByMove`; a Multi-Shot started there used to finish mid-walk; TestNoCastDuringMoveReaction). A walk out of melee keeps both hands swinging on their running timers until its step onto 5 yd, where a continuous walk leaves melee range (`KeepMeleeWhileWalkingOut`, switch `ForeverMeleeWhileWalkingOut`; every walk used to stop the swings at its start; TestMeleeSwingsWhileWalkingOut). A melee swing that lands mid-walk (melee starts at 5 yd while the walk goes on deeper) cannot start the walk back, so the walk's last step runs the rotation (`swungWhileMoving`, set in `attack.go`); a walk with no swing on the way is looked at again within 50 ms anyway (`ForeverMovePoll` in `apl.go`). Not in the move aura's expiry, which also runs when a fight ends mid-walk. Distances are whole yards (`WholeYards`: move spots, in `MoveTo` and the APL move action, and the starting distance in `character.go`): a walk steps a yard at a time, and a 4.5 yd spot was never reached, freezing the rotation. The walk back out can be decided on the way in once inside 5 yd (`APLActionMove.IsReady`, `queuedMove` in `MoveTo`): the move reaction runs during the rest of the walk and the walk out starts on arrival at the melee spot, not a reaction after it (TestWalkOutDecidedOnTheWayIn). Only that case: a walk out still decides at its end, since the next Auto Shot needs standing still (allowing every mid-walk decision cost the hard weaves and MythicSim's own Survival rotation up to 94). A run speed change during a walk (Strider Kick's +30% landing just after a walk out began, or running out mid-walk) re-paces the rest of it: the yard under way keeps the share already walked and the rest go at the new speed (`repace`); a walk took its pace from its start before. Walks with no change step as before, byte for byte (TestRunSpeedChangeMidWalk). Range checks read where you are mid-yard (`Unit.DistanceNow`): a walk off the 5 yd edge is out of melee as it starts, and a walk in from 11 yd is out of shot range as it starts; `DistanceFromTarget` itself still changes a yard at a time. It feeds the hunter's melee and shot range checks, the trap reach, Dragonfire and Caustic Smog, and the rotation's distance and walk time values (TestRangeChecksMidWalk). Before, a melee ability could land up to a yard's time (0.14 s) after a walk out began, and a shot after a walk in began. |
| `constants.go` | Ranged minimum distance is 11 yd center to center (`MinRangedAttackDistance`; the client's 8 yd is hitbox to hitbox, and in game shots stop inside about 10.5-11 yd of a mob's center). |
| `unit.go` | New fields: `MoveDelay`, `BaseRunSpeed`, `SkyborneConvergence`. |
| `character.go` | `MainHandCritBonus`: crit from a main-hand stone applies to melee attacks only, not ranged. `FillTalentsProto` skips trees past the third and digits past the last talent instead of crashing. `CapTalentsString` reads a talent string the way the site's talent picker shows it: each tree keeps only its own talents, each capped at its rank count (TestCapTalentsString, TestFillTalentsProtoSkipsWhatDoesNotFit). |
| `apl_actions_misc.go` | The Move action waits while a delayed move is pending. Start Attack (`APLActionStartAttack`): ready while the melee attack is off, so once per Auto Shot (TestStartAttackActionOncePerAutoShot). |
| `apl.go` | `ForeverMovePoll`: a rotation with a Move action is looked at every 50 ms while the GCD rolls too, as the shared core already did while the GCD sits free; before, a walk timed to a swing timer during a GCD waited for the GCD's end (TestMovePollDuringGcd, TestNoArrivalRecheckWithoutSwing). A hardcast ending this instant completes before the rotation decides, so a pre-pull Aimed Shot lands at the pull (TestPrepullAimedShotCompletesAtPull). |
| `apl_values_auto_attacks.go` | `autoTimeSinceLast`: time since that weapon's last swing or Auto Shot. Swing time minus time to the next swing reads short while a swing retry holds the next one (TestTimeSinceLastSwing). |
| `spell_resistances.go`, `spell.go` | No level-based partial resists for spells against a higher-level target (neither wowsims Forever nor MythicSim has them); resistances still apply. |
| `apl_values_operators.go` | A duration times or divided by a plain number scales correctly (fractions were truncated). |

Racials, consumables, buffs, debuffs, attack table

| File | Forever-specific behaviour |
| --- | --- |
| `racials.go` | Racials follow the client data: weapon specialisations give crit (no Bow/Gun/Thrown spec), Ruthless and Savage crit, Earth Fury and Endurance crit damage, Blood Fury (flat +10% attack power and spell power), Berserking (spell 20554), Elune's Light, Big Game Hunter, Beast Slaying, Quickness (stacks with other speed bonuses), Gnome Eureka!, Undead Touch of the Grave, and the Forever-only Skyborne race (Wind Blessed, Elemental Insight, Skysight, Read Leyline). |
| `base_stats.go` | Base stat offsets for Skyborne (Alliance and Horde). |
| `consumes.go` | Natural flasks: +60 Stamina plus a second effect that only works in the Hyjal zones (see `OutsideHyjal`). Shadow Oil's and Frost Oil's bolts use Forever's spells (350026, 350025): 48 to 56 and 31 to 35, since a client variance field is the whole spread, half each way (the same reading gives the potions, explosives, pet abilities and weapon procs their ranges). New foods and elixirs: Soaring Pamplona (+20 Intellect, the well fed buff it shares with Bat Hachee and Prehistoric Pulled Raptor), Greater Mageblood Elixir (20 mana per 5), Elixirs of Sages and the Whale (+25 Spirit, +2% crit). Elixirs take one slot per stat, the stat their buff raises first, as an elixir of one stat never stacks with another of it: Agility (Mongoose, Cunning), Strength (Juju Power, Giants, Grizzly, Ferocity), Attack Power (Juju Might, Winterfall Firewater), Intellect (`intellect_elixir`: the Owl), Spirit (`spirit_elixir`: Sages, the Whale) and Mana Regeneration (the Magebloods). Settings saved before move each elixir to its slot, or drop it if that is taken (`migrateElixirSlots`, and the site's twin in `ui/core/proto_utils/elixir_slots.ts`). The +2% crit of the Mongoose, Grizzly, Owl, Sages and Whale is one effect and counts once, for the first in slot order (`elixirCrit`). Explosives: Crystal Charge (instant) and the EZ-Thro / SAF-T line need no profession; the sapper, Dark Iron Bomb and the old dynamite and grenades need Engineering; all share one minute. The Holy Hand Grenade of Alterac does Holy damage. Scroll of Greater Cryoblast: a Mage item, 448 to 568 Frost (400 +-14.9% plus 12 a level from 46 to its cap of 55), a 0.5 sec cast. Every explosive shares the explosives' minute (spell category 24); only the Goblin Sapper Charge has a cooldown of its own on top, so the sapper and one other explosive are the only pair worth picking. `imbueFitsWeapon`: sharpening stones only on axes, daggers, polearms and swords; weightstones only on fist weapons, maces and staves; a misfit stone adds nothing. Damage potions on the shared 2 min potion cooldown: Dragonfire Potion (280 Fire at once and 35 every 2 sec for 8 sec, 12 yd cone, always hits, never crits) and Caustic Smog Potion (thrown at 10 yd/s from up to 20 yd, a 12 sec cloud of 33 Nature +-5% a second, always hits); neither adds spell power (TestForeverDamagePotions). Potion of Beast Slaying: +180 melee attack power against Beasts for 5 min (its aura has no ranged twin), never drunk without a Beast or while it lasts (TestPotionOfBeastSlaying). Goblin Sapper Charge hits its user through a spell of its own kept out of the results; the shared spell put the self hit on the first target's row (TestSapperSelfDamageIsNotDamageDone). |
| `buffs.go` | Raid buffs give their Forever client values, which already hold Classic's improved share (Gift of the Wild 16 stats, Grace of Air 89); no Forever talent tree has Improved Mark of the Wild, Enhancing Totems, Improved Blessing of Might or Wisdom, or Improved Battle Shout, so an old "Improved" setting gives the plain buff. Mana Spring keeps its Improved setting: Restorative Totems is in the Forever shaman tree with 5 ranks (16187, 5% a rank), so Improved is +25%; the armor and stamina buffs are unchanged. Windfury Totem (attack power buff lasts 1 s, 100 ms between procs, ranks 95/179/246) and Flametongue Totem (a party aura adding fire damage to main-hand white hits) are player-only, and Windfury wins if both are asked for. Windfury's extra swing is a white main-hand swing that leaves a queued Raptor Strike for the next normal swing (TestWindfuryExtraSwingLeavesRaptorQueued). Pets take no raid buffs. Grace of Air stacks with Windfury. |
| `debuffs.go` | Values updated to Forever data (Expose Armor, Demoralizing Shout, Thunder Clap, Curse of Recklessness and others). Caster debuffs such as Improved Shadow Bolt, Shadow Weaving, Improved Scorch, Winter's Chill and Stormstrike only raise damage taken from the caster who applied them. School resistance debuffs lower each of their schools by 75 (the inherited code lowered every resistance once per school listed: 375 each for Curse of the Elements; TestCurseOfElementsResistance). There is no Improved Hunter's Mark: the setting means the plain Mark, up the whole fight (the inherited code registered it without ever putting it up; TestPrepullTrueshotAndMark). |
| `target.go` | Attack table: no 1% hit suppression against +3 targets. `OutsideHyjal` flag on the encounter. `BonusMeleeAttackPowerTaken`: a creature type bonus to melee attack power only (read by `MeleeAttackPower` in `spell_result.go`). |
| `spell_outcome.go`, `stats/stats.go` | `PseudoStats.SpellCritMultiplier`: Earth Fury and Endurance raise spell crit damage too. |
| `focus.go`, `mana.go` | `MultiplyFocusRegen` scales the whole focus regeneration rate (Bestial Discipline uses it on top of the doubled pet rate). Mana regeneration while casting never exceeds the not-casting rate. `CastingRegenSource`: a regen-while-casting buff (Resourcefulness, Rapid Recuperation's sting buff) refreshes the cached mana ticks when it goes up or down, and its share of each casting tick is its own mana row. Mana per 5 ticks in a row per source (`NameMP5`, `AddStatsNamed`; items by their id, Blessing of Wisdom, Mana Spring Totem, mana oils, Mageblood, and any aura whose stats carry it); mana per 5 nobody names goes to "Mana per 5 (unnamed source)" (TestManaPerFiveRowsBySource), and the spirit rows hold spirit regen only. |
| `metrics_aggregator.go` | Auto Shot is counted as a cast in the results, not as a passive action. APM (`UnitMetrics.apm`, averaged per iteration and merged by weight in `sim_concurrent.go`): `countPress` counts each successful cast of a spell a rotation can use (`SpellFlagAPL`: shots, strikes, traps, stings, Raptor Strike queues, and every major cooldown: potions, explosives, trinkets, racials) in `Spell.Cast`, each walk the rotation starts (`MoveTo`, once per accepted walk) and each Start Attack, from the pull on; auto attacks and procs are not presses (TestAPMCountsPresses). |
| `sim_concurrent.go` | When worker results are merged, auras are matched by id, not by position. A worker that registers an aura mid-fight no longer breaks the merge; missing auras count as zero uptime. |

Encounter options and protos (`proto/`)

| File | Forever-specific content |
| --- | --- |
| `common.proto` | Skyborne races, Natural flasks, new elixirs and foods, the Intellect and Spirit elixir slots (`intellect_elixir`, `spirit_elixir`; the Elixir of Ferocity moved to `StrengthBuff`), `flametongue_totem`, `Encounter.outside_hyjal` (the "Hyjal zone" option; unset counts as inside). |
| `hunter.proto` | Pet families (Fox, Core Hound, Bird of Prey), per-ability pet switches, new ammo, Lacerating Strikes mode, pet totem switches, move delay, run speed, Auto Shot spam bug, melee swings reset Auto Shot, shots delay Auto Shot, Mortal Shots skips Auto Shot, Explosive Trap burn scaling, Skyborne convergence. |
| `apl.proto` | `APLSimpleView`, `APLSimpleBox`, `APLSimpleValue`: a plain-words view of a rotation line, written by the rotation's author. `APLValueAutoTimeSinceLast` ("Time Since Last Swing" in the editor). `APLAction.start_attack` and `APLRotation.manual_start_attack` (the Start Attack action and the walk-in switch). `APLRotation.recheck_after_auto_shot`: no longer read (the recheck is always on), kept so rotations saved with it still load. `APLValueWalkTime`. |
| `ui.proto` | `forever_changed_only` item filter. |

### Item data and effects (`sim/common`, `tools/database`)

| File | Content |
| --- | --- |
| `sim/common/enchant_effects.go` | Two-hand Impact at +6 weapon damage (Forever effect id). Grand Crusader and Grand Inquisitor strength procs. |
| `sim/common/item_effects.go` | Weapon proc damage ranges for 18 weapons, Hand of Justice at 1%, Nat Pagle's Broken Reel cooldown, a Forever version of the Blackhand's Breadth trinket, a weapon fire proc with its Forever spell id. Also declares the shared armor debuff group name `MinorArmorReduction`. |
| `sim/common/item_sets/dungeon_set_1.go` | Dungeon Set 1 bonuses at 2 to 6 pieces. Left out: the 4-piece when-struck effects and the caster 5-piece procs. |
| `sim/common/item_sets/crafted.go` | Devilsaur Armor 2-piece: +2% hit and +2% spell hit, per the client data (TestForeverSetBonuses). |
| `sim/common/item_sets/crafted.go` | Zul'Gurub crafted set bonuses need their profession at 300, as Forever's ItemSet rows say (the pieces need none): Bloodvine (Tailoring), Blood Tiger and Primal Batskin (Leatherworking; Primal Batskin's 3-piece is 8% run speed that does not stack with other speed bonuses), The Darksoul (Blacksmithing). |
| `sim/common/item_sets/item_sets_pve.go` | Cadaverous Garb bonuses changed to the Forever values; Gladiator set bonuses at 3, 4 and 5 pieces. |
| `tools/database/database.go` | `MergeExistingItems` (a feed never adds bare items), weapon skills replaced on merge, icon fallbacks for Forever-only ids. `ToUIProto` sorts enchants on their whole key, so ones sharing an effect id and slot (a weapon and a two-hand enchant) come out in the same order every run (TestEnchantOrderIsStable). |
| `tools/database/enchant_overrides.go` | Forever's enchant list: changed values, new effect ids, armor kits, new agility and strength enchants. |
| `tools/database/forever_overrides.go` | Generated: stat and weapon changes for 1114 items from the Forever item feed. |
| `tools/database/forever_new_items.go` | Generated: four Forever-only items a hunter can wear. |
| `tools/database/forever_icons.go` | Generated: icons for Forever-only items and spells. |
| `tools/database/overrides.go` | Adds the new foods, elixirs, potions and flasks to the list of consumable item ids the generator keeps. `ForeverProfessionFreeItems`: the 11 Zul'Gurub crafted pieces (Bloodvine, Primal Batskin, Blood Tiger, The Darksoul) anyone can wear in Forever, cleared by `ClearRequiredProfessions` in `gen_db`. |
| `tools/database/gen_db/main.go` | Applies the overrides. `FOREVER_OVERRIDES=off` builds the unmodified baseline. |

The three generated files name generator scripts (`scripts/gen_forever_*.js`) that are not in
this repository. Treat them as data; edit by hand only with care.

### Hunter class (`sim/hunter`)

| File | Forever-specific behaviour |
| --- | --- |
| `hunter.go` | Wires the new options (move delay, run speed, Skyborne convergence, Auto Shot spam cut). Swiftfeather Arrow and Swiftstrike Shot ammo. Mortal Shots applies to Auto Shot (Arc's log: Auto Shot crits at 2.30x a hit with 5/5 Mortal Shots); the old `mortal_shots_skips_auto_shot` option is no longer read (TestMortalShotsOnAutoShot). Option `shots_delay_auto_shot` (the beta's Auto Shot delay, off by default; beta logs show a 0.19 s median): the first shot cast in an Auto Shot cycle pushes the next Auto Shot back 0.2 s (TestShotsDelayAutoShot). `rangedHastedCast`: cast times take all ranged haste, the quiver included. Auto Shot flies at 40 yd/s. The talent string is capped at the Forever ranks before it is read (`TalentMaxRanks`, checked talent by talent against the site's tree by TestTalentRanksMatchSiteTree): a Classic string imported into the Forever tree gave Improved Stings and Barrage 5 points (3 ranks each) and stopped the sim on their rank tables; it now reads as the talent picker shows it, a 5 on a one-point talent included (TestTalentsCappedAtForeverRanks). |
| `aimed_shot.go`, `arcane_shot.go`, `sniper_shot.go`, `multi_shot.go`, `volley.go` | Aimed Shot 2.0 s cast, rank damage tables, Sniper Shot rank 2, Barrage bonus (3/7/10% on Multi-Shot, Aimed Shot, Volley). Aimed, Multi and Sniper Shot casts are sped up by all ranged haste, the quiver included (logs: Aimed 1.81 s, Multi-Shot 0.45 s with an 11% quiver, 0.32 s under Rapid Fire; TestShotCastsTakeRangedHaste). Missile speeds from the client data: 40 yd/s, Multi-Shot 30 (TestShotMissileSpeeds). Arcane Shot adds 0.10 of ranged attack power (players' logs hit for base + 0.10-0.12, and an aspect swap moved hits about 0.10 a point; TestArcaneShotScaling). |
| `serpent_sting.go` | Forever tick damage per rank from the client data (rank 9: 111). The client gives no attack power or spell power share, but in game each tick adds 0.030 of ranged attack power (rank 4: 27 at 174, 32 at 361 before rounding) and spell power changes nothing (TestSerpentStingScaling). A tick crits for double damage (65 crit ticks in logs, all from hunters without Mortal Shots), and Mortal Shots raises it like a shot's: its crit damage modifier's class mask (0x67a01) holds the sting's (0x4000) in the client data (TestMortalShotsOnSerpentStingTicks; a 5/5 Mortal Shots log with the sting up would confirm it, 2.3x). |
| `rapid_fire.go` | Also +40% melee attack speed. Rapid Killing takes time off the cooldown per rank. |
| `raptor_strike.go` | Forever rank damage. Queues at any range and waits for the next normal melee swing (not Windfury's extra one). Its cooldown is readable by rotations. No Mortal Shots crit bonus. |
| `mongoose_bite.go` | Mongoose Bite is ready while Defensive State (5302, dodges and parries; nothing attacks the hunter in a raid sim) or Expose Prey (1310726, its own 10 s buff from the talent) is up. Any bite spends Expose Prey, a miss too (38 missed bites in the logs); only a landed one spends Defensive State (TestMissedBiteSpendsExposePrey). Short marker buffs count what became of each Expose Prey (used by a Mongoose Bite, run out unused, refreshed by a new proc) for the buff table and timeline (TestExposePreyMarkers). Rotations check 1310726 for a ready bite. The buff table's group row shows the buff's own procs (`AuraMetrics.merge` in `ui/core/proto_utils/sim_result.ts`). |
| `hunters_mark.go` | Hunter's Mark is castable, so a rotation can put it up before the pull. It is the same aura as the raid debuff setting, so with that on it only refreshes it (TestPrepullTrueshotAndMark). |
| `strider_kick.go` | Also +30% run speed for 3 s. Yellow-only proc mask. |
| `aspects.go` | Quick Strikes buff, and Aspect of the Cheetah (+30% run speed, shares the aspect slot, Pathfinding). |
| `immolation_trap.go`, `explosive_trap.go` | Forever spell ids, 15 s burn in 3 s ticks, magic spell (ticks crit for 1.5x, no melee procs), shared 30 s cooldown, attack power share from T.N.T. (`trapAPBonus`, higher of melee and ranged attack power). Explosive Trap's burn only scales with attack power if the option is on.. `trapReach`: how far from the target a trap can be dropped and still catch it, 6 yd (it lands about a yard in front of you; in game it hits from 6 yd, not from 8); Immolation Trap cannot be dropped beyond it and both do nothing beyond it (TestTrapReach; AB_TRAP_REACH for comparisons). |
| `summon_hawk.go` | Hawk attacks follow the hunter's ranged attack speed and resolve as white melee (glancing, dodges, misses) by the level table. A hawk's own attack power is 80% of the hunter's ranged attack power (1,245 logged swings: 264 of 264 beside an owner swing at 0.800; the owner's logged value moves 50 with Aspect of the Hawk); a swing is 0.0815 of the creature dps over 2.5 s (Unleashed Fury's 20% taken out), a third either side, plus 2.5 / 134.4 a point of that attack power, times Unleashed Fury. It crits for double damage at half the hunter's crit plus Ferocity (an owner's pet critting 22.2% and hawk 16.2% solve to 12% owner crit and 10% Ferocity). TestHawkDamage, TestHawkCritAndUnleashedFury. A hawk's swing is the hawk's attack: it sets off none of the hunter's on-hit effects (no proc mask), so Judgement of Wisdom's mana goes to the hawk, not the hunter; swings that counted as the hunter's ranged attacks had about doubled its mana (TestHawkSwingsGiveNoHunterMana). The dive scales with ranged attack power and Unleashed Fury. Not modeled for lack of evidence: hit or the hunter's damage multipliers carried over to the hawk. |
| `pet.go` | A pet's swing is the creature dps for its level (ExpectedStat CreatureAutoAttackDps, 60.45 at 60), 20% either side, over its swing time, plus its attack power: a level 30 pet in the logs dealt 43.9 raw a second = 1.25 (happy) x (24.69 + 146 attack power / 14) (TestPetSwingFromCreatureDps; AB_PET_CLASSIC_BASE=1 brings back the inherited Classic 18.17-27.66). No crit from the pet's own Agility: its crit is its owner's (a pet critting 9.9% of 5067 swings beside its owner's 11.5% melee crit). Spell hit is 17/8 of the inherited hit (a Forever dev, forever-bugs #394). The pet's hit is its owner's ranged hit, so ranged-only hit (the Biznicks 247x128 Accurascope's +3%, `PseudoStats.BonusRangedHitRating`) is added to it once in `ApplyTalents`, 17/8 of it as spell hit (TestPetInheritsRangedOnlyHit). The pet's attack power is 10% of the higher of its owner's melee and ranged attack power, worked out from the owner's totals at every change (`Pet.addOwnerStats` in `sim/core/pet.go`): from the change alone, a buff that ended never came off the pet and its attack power climbed all fight (TestPetAttackPowerFollowsOwner). Pet attack speed is one setting per pet (default: the family's fastest tamable speed). Focus regeneration doubled. Family table with each family's ability list and health, armor and damage scalars from its Tamed Pet Passive in the client data (the Fox carries the Wolf's; TestPetFamilyScalars); the Bat keeps Sonic Blast, whose Forever ranks sit on the Bat skill line (653). No Pet Aggression or Cobra Reflexes: Forever pets learn neither. |
| `pet_abilities.go` | Top-rank pet abilities per the client data: Bite, Claw, Lightning Breath, Scorpid Poison, Savage Rend, Furious Howl (buffs the hunter), Thunderstomp, Dismember, Pinch, Tendon Rip, Web, Swipe, Mine!, Charge (opener), Demoralizing Screech, Lava Breath, Dust Cloud, Trickster's Dance, Sonic Blast. The pet rotation handles any number of abilities and honours the per-ability switches. |
| `talents.go` | Forever per-rank values (Lightning Reflexes, Unleashed Fury, Efficiency, Barrage, Survival Tactics, Ferocity, Resourcefulness, Bestial Discipline), Intimidation, Frenzy buff, Lacerating Strikes with three bleed modes (refresh, strongest stays, rolling). Trueshot Aura is also a castable spell for a rotation's pull (no cost); every hunter knows it, so a rotation shared with untalented builds gives no warning, and it only goes out with the talent. |
| `items.go` | Item effects that buff the pet's crit now apply to every pet ability, not two fixed ones. |

### Tests

| File | Covers |
| --- | --- |
| `sim/hunter/forever_patch_test.go` | The bulk of the hunter checks (52 tests): totems, racials, Auto Shot stand-still, spam bug and melee restart, off-hand wait, mid-walk swings, Raptor Strike queueing, traps (attack power, crit, cooldown, cast metrics), shared armor debuff group, Aspect of the Cheetah, run speed setting, imbue fit, foods, elixirs, weapon enchants, Natural flasks, Skyborne racials, set bonuses, Expose Prey counts, pull casts, damage potions, Potion of Beast Slaying, crafted set bonuses needing their profession, range checks mid-walk. |
| `sim/hunter/forever_buffs_test.go` | Raid buffs at their client values (an old "Improved" setting gives the same) and Forever debuff values. |
| `sim/hunter/forever_hawk_test.go` | Hawk swing interval, damage and outcomes. |
| `sim/hunter/bm_talents_test.go` | Bestial Discipline, casting regeneration cap, Intimidation. |
| `sim/hunter/forever_ab_test.go` | Local A/B harness (dual wield vs two-hand under matched settings); prints a table, almost no assertions. |
| `sim/hunter/forever_settings_test.go` | Runs a settings export from the site (`AB_SETTINGS`) and prints the per-action table; skipped without it. |
| `sim/hunter/forever_stats_test.go` | Prints final character stats of a settings export to check buff claims; skipped without `AB_SETTINGS`. |
| `sim/core/forever_racials_test.go` | Touch of the Grave proc chance by class. |
| `sim/core/talents_test.go`, `sim/hunter/talent_ranks_test.go` | Talent strings: capped at each talent's ranks as the talent picker shows them, no crash on extra trees or digits, the hunter rank table against `ui/core/talents/trees/hunter.json`, an imported Classic string. |
| `tools/database/enchant_order_test.go` | The item database generator writes enchants in a fixed order. |
| `tools/database/profession_free_test.go` | The shipped item database asks no profession for the Zul'Gurub crafted pieces Forever unlocked, and keeps the locks Forever keeps. |
| `sim/core/sim_concurrent_test.go` | Merging worker results matches auras by id. |
| `sim/forever_racials_test.go` | Crit-damage racials, Blood Fury spell power, Eureka!, Touch of the Grave drain, caster-only debuffs, each on other classes; the Mage-only Scroll of Greater Cryoblast. Package `sim` is not in the test command below; run it with `go test --tags=with_db ./sim/`. |

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
is skipped. Without the skip it fails on those numbers but no longer crashes the test binary
(its Classic talent string did, at the first subtest, so nothing after it could run). The `lib.wasm` files under `dist/classic`, `public` and `ui/public` are committed
build outputs.

## Assumptions

Values the code or the site marks as assumed, estimated or not yet confirmed in game:

- Pets follow the creature dps table from level 30 (where logs match it) to 60; no level 60 logs exist.
- Hawk swing damage at level 60 scales its base from level 30 by the creature damage table, and
  assumes the hawk is the hunter's level (`sim/hunter/summon_hawk.go`).
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
- Clever Traps raises the whole trap damage, the attack power share included. The 0.132 logs had Clever Traps
  2/2 (Counterattack PvP spec), so the share is read off them with the 1.3 taken out, and every build with Clever
  Traps gets the same damage under either rule. Only a build without it would differ: if Clever Traps scales the
  base only, the share is 0.172 there (BM Weave +2.8 DPS, AB_TRAP_AP=0.1716). No log without Clever Traps exists.
- Only the strongest run speed bonus counts, except Night Elf Quickness, which is the only one
  known to stack (`AB_SKYSIGHT_STACKS=1` lets Skysight stack, for comparison).
- Core Hound and Fox use neutral family scalars (damage 1.00); nothing else is listed for them
  in `sim/hunter/pet.go`.
- Shipped builds assume no extra move time and the game's run speed of 7 yd/s; the step
  distances in the rotations carry the footwork (`ui/hunter/presets.ts`).
- Caustic Smog Potion ticks 12 times (its cloud lasts 12 sec; an 11th or 12th tick at the edge is not confirmed in
  game). Every target stands in the cloud and in Dragonfire's cone, and the cone reaches 12 yd center to center.
- Shadow Oil's and Frost Oil's bolts carry points per level in the client data (0.9 and 0.7; max level 17 and 12,
  spell level 12 and 8, base level 0). The sim keeps their base (52 and 33). At level 60 the server formula would add
  4.5 and 2.8, and counting from base level 0 (as Wowhead's tooltips do) 15.3 and 8.4; a combat log of oil procs settles it.
- Elixirs stack unless they raise the same stat first (one slot per stat, as WoWSims Forever groups them); the client
  data has no stacking table, so this is not confirmed in game. The elixirs' +2% crit counting once comes from testers,
  not a log.
- Left out on purpose: Aspect of the Cheetah's daze, the self-heal of Touch of the Grave, and
  caster-only effects of dungeon sets.
