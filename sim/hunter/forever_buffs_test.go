package hunter

import (
	"math"
	"testing"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
)

const (
	regular  = proto.TristateEffect_TristateEffectRegular
	improved = proto.TristateEffect_TristateEffectImproved
)

// stat difference between two setups (b minus a)
func statGain(t *testing.T, stat stats.Stat, a, b patchTestSetup) float64 {
	_, ha := newPatchTestHunter(t, a)
	_, hb := newPatchTestHunter(t, b)
	return hb.GetStat(stat) - ha.GetStat(stat)
}

// Forever has no talents that improve raid buffs (none of Improved Battle Shout, Improved Blessing of Might or Wisdom,
// Enhancing Totems, Improved Mark of the Wild is in a talent tree, and the client values already hold Classic's
// improved share), so each buff gives its client value and an old "Improved" setting gives the same; Power Word:
// Fortitude is 70 Stamina.
func TestImprovedBuffs(t *testing.T) {
	for _, tc := range []struct {
		name string
		stat stats.Stat
		a, b patchTestSetup
		want float64
	}{
		{"Battle Shout", stats.AttackPower, patchTestSetup{}, patchTestSetup{buffs: &proto.RaidBuffs{BattleShout: improved}}, 139},
		{"Blessing of Might", stats.AttackPower, patchTestSetup{}, patchTestSetup{ibuffs: &proto.IndividualBuffs{BlessingOfMight: improved}}, 133},
		{"Strength of Earth", stats.Strength, patchTestSetup{}, patchTestSetup{buffs: &proto.RaidBuffs{StrengthOfEarthTotem: improved}}, 53},
		{"Grace of Air", stats.Agility, patchTestSetup{}, patchTestSetup{buffs: &proto.RaidBuffs{GraceOfAirTotem: improved}}, 89},
		{"Mark of the Wild", stats.Agility, patchTestSetup{}, patchTestSetup{buffs: &proto.RaidBuffs{GiftOfTheWild: improved}}, 16},
		{"Blessing of Wisdom", stats.MP5, patchTestSetup{ibuffs: &proto.IndividualBuffs{BlessingOfWisdom: regular}}, patchTestSetup{ibuffs: &proto.IndividualBuffs{BlessingOfWisdom: improved}}, 0},
		{"Power Word: Fortitude", stats.Stamina, patchTestSetup{}, patchTestSetup{buffs: &proto.RaidBuffs{PowerWordFortitude: regular}}, 70},
		{"Slip'kik's Savvy", stats.MeleeCrit, patchTestSetup{}, patchTestSetup{ibuffs: &proto.IndividualBuffs{SlipkiksSavvy: true}}, 3},
	} {
		if got := statGain(t, tc.stat, tc.a, tc.b); math.Abs(got-tc.want) > 1e-6 {
			t.Errorf("%s: +%.2f %v, want +%.0f", tc.name, got, tc.stat, tc.want)
		}
	}
}

// Forever has no Improved Hunter's Mark talent: an old "Improved" setting gives the plain Mark (110 ranged AP);
// Expose Armor removes 2250 armor.
func TestForeverDebuffs(t *testing.T) {
	_, marked := newPatchTestHunter(t, patchTestSetup{debuffs: &proto.Debuffs{HuntersMark: improved}})
	if len(marked.CurrentTarget.GetAurasWithTag(core.HuntersMarkAuraTag)) == 0 {
		t.Fatal("no Hunter's Mark on the target")
	}
	if aura := marked.CurrentTarget.GetAura("HuntersMark-110"); aura == nil {
		t.Errorf("an Improved Hunter's Mark setting is not the plain 110")
	}
	_, plain := newPatchTestHunter(t, patchTestSetup{})
	sim, exposed := newPatchTestHunter(t, patchTestSetup{debuffs: &proto.Debuffs{ExposeArmor: regular}})
	// the sim applies it a few seconds into the fight
	exposed.CurrentTarget.GetAura("ExposeArmor").Activate(sim)
	if got := plain.CurrentTarget.GetStat(stats.Armor) - exposed.CurrentTarget.GetStat(stats.Armor); math.Abs(got-2250) > 1e-6 {
		t.Errorf("Expose Armor removes %.0f armor, want 2250", got)
	}
}
