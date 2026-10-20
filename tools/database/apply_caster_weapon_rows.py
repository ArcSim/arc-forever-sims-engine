#!/usr/bin/env python3
"""Apply the caster weapon rows of assets/db_inputs/forever_caster_weapon_rows.json to assets/database/db.json.

MythicSim patch 180, ported from wowsims/forever: caster weapons the client flags REDUCED_CASTER_DPS deal reduced
melee damage, and CASTER_SPELL_POWER / CASTER_HEALING_POWER weapons carry built-in spell power and healing. Our
item generator has neither rule, so the rows are taken from upstream's database. Run it from the fork's root after
any item regeneration, then `go run ./tools/sync_db_binary`. It rewrites only the listed items' lines.
"""
import json

ROWS = "assets/db_inputs/forever_caster_weapon_rows.json"
DB = "assets/database/db.json"
HEALING, SPELL_DAMAGE = "4", "5"  # stats.HealingPower, stats.SpellDamage

rows = {r["id"]: r for r in json.load(open(ROWS))["items"]}
lines = open(DB).read().split("\n")
changed = 0
section = None
for i, line in enumerate(lines):
    if line.startswith('"') and line.endswith(":["):
        section = line.split('"')[1]
        continue
    # Other sections (enchants, gems, ...) carry ids of their own that can equal an item's.
    if section != "items" or not line.startswith('{"id":'):
        continue
    trailing = "," if line.endswith(",") else ""
    item = json.loads(line[: len(line) - len(trailing)])
    row = rows.get(item["id"])
    if row is None:
        continue
    opts = item.setdefault("scalingOptions", {}).setdefault("0", {})
    opts["weaponDamageMin"], opts["weaponDamageMax"] = row["weaponDamageMin"], row["weaponDamageMax"]
    stats = opts.setdefault("stats", {})
    for key, value in ((HEALING, row["healingPower"]), (SPELL_DAMAGE, row["spellDamage"])):
        if value:
            stats[key] = value
        else:
            stats.pop(key, None)
    new = json.dumps(item, separators=(",", ":"), ensure_ascii=False) + trailing
    if new != line:
        lines[i] = new
        changed += 1
open(DB, "w").write("\n".join(lines))
print(f"applied {len(rows)} caster weapon rows, {changed} lines changed")
