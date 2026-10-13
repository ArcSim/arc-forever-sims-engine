#!/usr/bin/env python3
"""Merge Forever Cache tooltips into assets/database/db.json.

- Item present in Forever: update stats / weapon numbers from the tooltip.
- Item only in Forever: add it.
- Item only in the old db: leave it untouched.
"""

import json
import re
import urllib.request
from pathlib import Path

CACHE_URL = "https://raw.githubusercontent.com/Napalmsteak/Classic-WoW-Item-Caches/main/Forever%20Cache.json"
DB_PATH = Path("assets/database/db.json")

STAT_MARKERS = {
    "stat4": 0,   # Strength
    "stat3": 1,   # Agility
    "stat7": 2,   # Stamina
    "stat5": 3,   # Intellect
    "stat6": 4,   # Spirit
    "rtg41": 5,   # Spell Power
    "rtg43": 12,  # MP5
    "rtg38": 17,  # Attack Power
    "amr": 26,    # Armor
}
TEXT_STATS = [
    (r"chance to hit by ([\d.]+)%", 18),
    (r"critical strike by ([\d.]+)%", 19),
    (r"chance to Dodge an attack by ([\d.]+)%", 31),
    (r"chance to Parry an attack by ([\d.]+)%", 32),
]

SLOT_TO_TYPE = {
    "Head": 1, "Neck": 2, "Shoulder": 3, "Back": 4, "Chest": 5,
    "Wrist": 6, "Hands": 7, "Waist": 8, "Legs": 9, "Feet": 10,
    "Finger": 11, "Trinket": 12, "One-Hand": 13, "Two-Hand": 13,
    "Main Hand": 13, "Off Hand": 13, "Held In Off-hand": 13,
    "Ranged": 13, "Relic": 13,
}
ARMOR_TYPE = {"Cloth": 1, "Leather": 2, "Mail": 3, "Plate": 4}
WEAPON_TYPE = {
    "Axe": 1, "Dagger": 2, "Mace": 4, "Staff": 5, "Polearm": 6,
    "Sword": 9, "Fist Weapon": 11, "Bow": 18, "Crossbow": 18, "Gun": 18,
}
HAND_TYPE = {"One-Hand": 2, "Main Hand": 1, "Off Hand": 2, "Two-Hand": 4, "Held In Off-hand": 5}
QUALITY = {"q0": 0, "q1": 1, "q2": 2, "q3": 3, "q4": 4, "q5": 5}

STAT_LEN = 48


def grab(html, pattern):
    m = re.search(pattern, html, re.I)
    return m.group(1) if m else None


def parse_tooltip(html):
    stats = [0] * STAT_LEN
    for marker, idx in STAT_MARKERS.items():
        m = re.search(rf"<!--{marker}-->\+?([\d.]+)", html)
        if m:
            stats[idx] = int(float(m.group(1)))
    for pattern, idx in TEXT_STATS:
        m = re.search(pattern, html, re.I)
        if m:
            stats[idx] = int(float(m.group(1)))

    dmg = re.search(r"<!--dmg-->(\d+)\s*-\s*(\d+)", html)
    spd = re.search(r"<!--spd-->([\d.]+)", html)
    ilvl = grab(html, r"Item Level.*?(\d+)")
    quality = 1
    for cls, q in QUALITY.items():
        if f'class="{cls}"' in html or f"class='{cls}'" in html:
            quality = q
            break

    slot = grab(html, r"<td>(Head|Neck|Shoulder|Back|Chest|Wrist|Hands|Waist|Legs|Feet|Finger|Trinket|One-Hand|Two-Hand|Main Hand|Off Hand|Held In Off-hand|Ranged|Relic)</td>")
    armor = grab(html, r">(Cloth|Leather|Mail|Plate)</span>")
    weapon = grab(html, r">(Axe|Dagger|Mace|Staff|Polearm|Sword|Fist Weapon|Bow|Crossbow|Gun)</span>")

    return {
        "stats": stats,
        "ilvl": int(ilvl) if ilvl else 0,
        "quality": quality,
        "type": SLOT_TO_TYPE.get(slot, 0),
        "armorType": ARMOR_TYPE.get(armor, 0),
        "weaponType": WEAPON_TYPE.get(weapon, 0),
        "handType": HAND_TYPE.get(slot, 0),
        "weaponDamageMin": int(dmg.group(1)) if dmg else 0,
        "weaponDamageMax": int(dmg.group(2)) if dmg else 0,
        "weaponSpeed": float(spd.group(1)) if spd else 0,
    }


def main():
    print("Downloading Forever Cache...")
    with urllib.request.urlopen(CACHE_URL) as r:
        forever = json.load(r)

    db = json.loads(DB_PATH.read_text(encoding="utf-8"))
    by_id = {item["id"]: item for item in db["items"]}

    added = updated = skipped = 0
    for entry in forever:
        parsed = parse_tooltip(entry.get("tooltipHtml", ""))
        if parsed["type"] == 0 and parsed["weaponDamageMax"] == 0:
            skipped += 1
            continue

        item_id = entry["itemId"]
        existing = by_id.get(item_id)

        if existing is None:
            new_item = {
                "id": item_id,
                "name": entry["name"],
                "icon": entry.get("iconName") or "inv_misc_questionmark",
                "type": parsed["type"],
                "stats": parsed["stats"],
                "weaponSkills": [0] * 16,
                "ilvl": parsed["ilvl"],
                "phase": 1,
                "quality": parsed["quality"],
            }
            if parsed["armorType"]:
                new_item["armorType"] = parsed["armorType"]
            if parsed["weaponType"]:
                new_item["weaponType"] = parsed["weaponType"]
            if parsed["handType"]:
                new_item["handType"] = parsed["handType"]
            if parsed["weaponDamageMax"]:
                new_item["weaponDamageMin"] = parsed["weaponDamageMin"]
                new_item["weaponDamageMax"] = parsed["weaponDamageMax"]
                new_item["weaponSpeed"] = parsed["weaponSpeed"]
            db["items"].append(new_item)
            by_id[item_id] = new_item
            added += 1
        else:
            existing["stats"] = parsed["stats"]
            if parsed["ilvl"]:
                existing["ilvl"] = parsed["ilvl"]
            if parsed["weaponDamageMax"]:
                existing["weaponDamageMin"] = parsed["weaponDamageMin"]
                existing["weaponDamageMax"] = parsed["weaponDamageMax"]
                existing["weaponSpeed"] = parsed["weaponSpeed"]
            updated += 1

    DB_PATH.write_text(json.dumps(db, separators=(",", ":")), encoding="utf-8")
    print(f"Added {added}, updated {updated}, skipped non-gear {skipped}")
    print("Items only in the old db were left alone.")


if __name__ == "__main__":
    main()
