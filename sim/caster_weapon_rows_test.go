//go:build with_db

package sim

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

// Patch 180: caster weapons take upstream's rows (wowsims/forever 7c9528c140). The client's REDUCED_CASTER_DPS
// flag cuts a caster weapon's melee damage, and CASTER_SPELL_POWER / CASTER_HEALING_POWER give it built-in spell
// power and healing; assets/db_inputs/forever_caster_weapon_rows.json carries the 158 weapons.
func TestCasterWeaponRows(t *testing.T) {
	for _, c := range []struct {
		id                    int32
		name                  string
		min, max, sp, healing float64
	}{
		{20069, "Ironbark Staff", 114, 172, 88, 88},
		{2243, "Hand of Edward the Odd", 33, 62, 90, 90},
	} {
		if core.GetItemByID(c.id) == nil {
			t.Fatalf("%s (%d) is not in the database", c.name, c.id)
		}
		item := core.NewItem(core.ItemSpec{ID: c.id})
		if item.WeaponDamageMin != c.min || item.WeaponDamageMax != c.max ||
			item.Stats[stats.SpellDamage] != c.sp || item.Stats[stats.HealingPower] != c.healing {
			t.Errorf("%s: %v-%v damage, %v spell damage, %v healing; want %v-%v, %v, %v", c.name,
				item.WeaponDamageMin, item.WeaponDamageMax, item.Stats[stats.SpellDamage], item.Stats[stats.HealingPower],
				c.min, c.max, c.sp, c.healing)
		}
	}
}
