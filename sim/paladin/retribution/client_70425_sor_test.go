package retribution

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
	"github.com/wowsims/forever/sim/paladin"
)

const (
	judgementOfRighteousnessRank8 = 20286
	holyStrikeRank8               = 10333
	consecrationRank5             = 20924
)

// MythicSim patch 148: Seal of Righteousness and Holy Strike against Ashreon's bugs forum post of
// 17 October 2026, and Consecration's first four targets.
//
// Each sim seals, judges and Holy Strikes with only a weapon and the given spell power, so a ratio
// of two runs cancels every talent and debuff multiplier. Crits are left out.
type patch148Hits struct {
	seal, judgement, holyStrike float64
	targetDamage                []float64 // Consecration's damage on each target
}

func patch148Run(t *testing.T, weapon int32, spellPower float64, targets int, rotation string) patch148Hits {
	t.Helper()
	items := make([]*proto.ItemSpec, 15)
	for i := range items {
		items[i] = &proto.ItemSpec{}
	}
	items[14] = &proto.ItemSpec{Id: weapon}
	var bonus stats.Stats
	bonus[stats.SpellDamage] = spellPower
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceHuman,
		Class:         proto.Class_ClassPaladin,
		Equipment:     &proto.EquipmentSpec{Items: items},
		BonusStats:    &proto.UnitStats{Stats: bonus.ToProtoArray()},
		TalentsString: RetTalents,
		Rotation:      core.APLRotationFromJsonString(rotation),
	}, DefaultOptions)
	encounter := core.MakeSingleTargetEncounter(0)
	for len(encounter.Targets) < targets {
		encounter.Targets = append(encounter.Targets, core.NewDefaultTarget())
	}
	result := core.RunRaidSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 148, Iterations: 40},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  encounter,
	})
	if result.Error != nil {
		t.Fatal(result.Error.Message)
	}
	out := patch148Hits{targetDamage: make([]float64, targets)}
	nonCrit := func(id int32) float64 {
		var damage float64
		var hits int32
		for _, action := range result.RaidMetrics.Parties[0].Players[0].Actions {
			if action.Id.GetSpellId() != id {
				continue
			}
			for _, target := range action.Targets {
				damage += target.Damage - target.CritDamage
				hits += target.Hits - target.Crits
			}
		}
		if hits == 0 {
			return 0
		}
		return damage / float64(hits)
	}
	out.seal = nonCrit(sealOfRighteousnessHit8)
	out.judgement = nonCrit(judgementOfRighteousnessRank8)
	out.holyStrike = nonCrit(holyStrikeRank8)
	for _, action := range result.RaidMetrics.Parties[0].Players[0].Actions {
		if action.Id.GetSpellId() != consecrationRank5 {
			continue
		}
		for _, target := range action.Targets {
			if i := int(target.UnitIndex); i < targets {
				out.targetDamage[i] += target.Damage
			}
		}
	}
	return out
}

const sealJudgeStrike = `{"type":"TypeAPL","priorityList":[
	{"action":{"castSpell":{"spellId":{"spellId":20293}},"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":20293}}}}}}},
	{"action":{"castSpell":{"spellId":{"spellId":20271}}}},
	{"action":{"castSpell":{"spellId":{"spellId":10333}}}}
]}`

func withSwitch(sw *bool, value bool, run func()) {
	old := *sw
	*sw = value
	defer func() { *sw = old }()
	run()
}

// Rank 8's judgement and seal hit come from the rows' base points, not the level 60 gain: 170
// against 178 on the judgement (162 to 178, as Ashreon's post gives), 1786 against 1880 a hundred
// weapon speed on the seal's hit (its tooltip's 21 to 71). A run with 1000 more spell power gives
// the multiplier each hit carries (the paladin also has a little spell power of its own).
func TestSealOfRighteousnessBasesAtSpellLevel(t *testing.T) {
	var off patch148Hits
	withSwitch(&paladin.SealOfRighteousnessBasesAtSpellLevel, false, func() {
		off = patch148Run(t, massOfMcGowan, 0, 1, sealJudgeStrike)
	})
	on := patch148Run(t, massOfMcGowan, 0, 1, sealJudgeStrike)
	power := patch148Run(t, massOfMcGowan, 1000, 1, sealJudgeStrike)
	sealMultiplier := (power.seal - on.seal) / (1000 * paladin.SealOfRighteousnessHitCoefficient(0.20, false))
	judgementMultiplier := (power.judgement - on.judgement) / (1000 * paladin.JudgementOfRighteousnessSpellPowerCoefficient)
	t.Logf("seal hit %.2f -> %.2f (multiplier %.3f), judgement %.2f -> %.2f (multiplier %.3f)",
		off.seal, on.seal, sealMultiplier, off.judgement, on.judgement, judgementMultiplier)
	// Mass of McGowan, 2.8 speed one-hander: (1880 - 1786) / 100 x 2.8 x 0.85 = 2.24 a hit.
	if got, want := off.seal-on.seal, 0.94*2.8*0.85*sealMultiplier; math.Abs(got-want) > 0.01 {
		t.Errorf("the seal's hit lost %.3f, want %.3f", got, want)
	}
	// The judgement rolls, so its averages carry noise: 178 - 170 = 8 before multipliers.
	if got, want := off.judgement-on.judgement, 8*judgementMultiplier; math.Abs(got-want) > 1.5 {
		t.Errorf("the judgement lost %.2f, want about %.2f", got, want)
	}
}

// sealShare is what 1000 spell power adds to a seal hit with the weapon, under the switches as set.
func sealShare(t *testing.T, weapon int32) float64 {
	return patch148Run(t, weapon, 1000, 1, sealJudgeStrike).seal - patch148Run(t, weapon, 0, 1, sealJudgeStrike).seal
}

// MythicSim patch 158: the seal hit takes half the judgement dummy's share, 0.10 with a one-hander and
// 0.12 with a two-hander, against the pre-fix 0.20 and 0.24.
func TestSealOfRighteousnessSpellPowerFixed(t *testing.T) {
	if !paladin.SealOfRighteousnessSpellPowerFixed {
		t.Fatal("SealOfRighteousnessSpellPowerFixed is off by default; patch 158 turns it on")
	}
	for _, c := range []struct {
		name   string
		weapon int32
	}{
		{"one-hander", massOfMcGowan},
		{"two-hander", arcaniteReaper},
	} {
		fixed := sealShare(t, c.weapon)
		var bugged float64
		withSwitch(&paladin.SealOfRighteousnessSpellPowerFixed, false, func() { bugged = sealShare(t, c.weapon) })
		t.Logf("%s: 1000 spell power adds %.2f a hit, %.2f before the fix", c.name, fixed, bugged)
		if got := fixed / bugged; math.Abs(got-0.5) > 1e-6 {
			t.Errorf("%s: the fixed share over the pre-fix one is %.5f, want 0.5", c.name, got)
		}
	}
}

// The per-hit spell power share from client 70425's tooltip, 0.0765 with a one-hander and 0.12
// with a two-hander, against the judgement dummy's 0.20 and 0.24. Off by default; when on it
// replaces patch 158's half of the dummy as well.
func TestSealOfRighteousnessPerHitTooltipCoefficient(t *testing.T) {
	if paladin.SealOfRighteousnessPerHitTooltipCoefficient {
		t.Fatal("SealOfRighteousnessPerHitTooltipCoefficient is on by default; patch 148 leaves it off")
	}
	for _, c := range []struct {
		name   string
		weapon int32
		want   float64
	}{
		{"one-hander", massOfMcGowan, 0.0765 / 0.20},
		{"two-hander", arcaniteReaper, 0.12 / 0.24},
	} {
		var dummy, tooltip, tooltipUnfixed float64
		withSwitch(&paladin.SealOfRighteousnessSpellPowerFixed, false, func() {
			dummy = sealShare(t, c.weapon)
			withSwitch(&paladin.SealOfRighteousnessPerHitTooltipCoefficient, true, func() { tooltipUnfixed = sealShare(t, c.weapon) })
		})
		withSwitch(&paladin.SealOfRighteousnessPerHitTooltipCoefficient, true, func() { tooltip = sealShare(t, c.weapon) })
		t.Logf("%s: 1000 spell power adds %.2f a hit with the dummy's share, %.2f with the tooltip's", c.name, dummy, tooltip)
		if got := tooltip / dummy; math.Abs(got-c.want) > 1e-6 {
			t.Errorf("%s: tooltip share over the dummy's %.5f, want %.5f", c.name, got, c.want)
		}
		if tooltip != tooltipUnfixed {
			t.Errorf("%s: with the tooltip switch on, patch 158's switch moved the share from %.4f to %.4f", c.name, tooltipUnfixed, tooltip)
		}
	}
}

// The other reading of "Holy Strike 40-53": the flat amount before the 50% is 40 to 53, half the
// row's 81 to 105, so each hit loses a quarter of 93 (23.25) before multipliers. Off by default.
func TestHolyStrikeFlatBeforeWeaponPercentIs40To53(t *testing.T) {
	if paladin.HolyStrikeFlatBeforeWeaponPercentIs40To53 {
		t.Fatal("HolyStrikeFlatBeforeWeaponPercentIs40To53 is on by default; patch 148 leaves it off")
	}
	base := patch148Run(t, massOfMcGowan, 0, 1, sealJudgeStrike).holyStrike
	withPower := patch148Run(t, massOfMcGowan, 1000, 1, sealJudgeStrike).holyStrike
	var halved float64
	withSwitch(&paladin.HolyStrikeFlatBeforeWeaponPercentIs40To53, true, func() {
		halved = patch148Run(t, massOfMcGowan, 0, 1, sealJudgeStrike).holyStrike
	})
	// 1000 spell power at Holy Strike's 0.286 gives the multiplier the hit carries.
	multiplier := (withPower - base) / (1000 * paladin.HolyStrikeSpellPowerCoefficient)
	t.Logf("Holy Strike %.2f a hit, %.2f with the flat halved (multiplier %.3f)", base, halved, multiplier)
	if got, want := base-halved, 0.25*93*multiplier; math.Abs(got-want) > 1.5 {
		t.Errorf("halving the flat amount takes %.2f a hit, want about %.2f", got, want)
	}
}

// Consecration: every target takes the base tick, only the first four the bonus and its spell power
// share ("The first 4 enemies who enter the area will take an additional ..."; Ashreon: 33% for the
// first 4 targets). The engine already did; this pins it.
func TestConsecrationBonusOnFirstFourTargets(t *testing.T) {
	consecrate := `{"type":"TypeAPL","priorityList":[{"action":{"castSpell":{"spellId":{"spellId":20924}}}}]}`
	hits := patch148Run(t, massOfMcGowan, 500, 5, consecrate)
	t.Logf("Consecration damage by target: %.0f", hits.targetDamage)
	first := hits.targetDamage[0]
	if first == 0 {
		t.Fatal("no Consecration damage")
	}
	for i := 1; i < 4; i++ {
		if math.Abs(hits.targetDamage[i]/first-1) > 0.1 {
			t.Errorf("target %d took %.0f, target 0 %.0f; the first four should take the same", i, hits.targetDamage[i], first)
		}
	}
	// The fifth: 12 a tick against 12 + 27 + 0.042 x 500 = 60, before crits and multipliers.
	if got := hits.targetDamage[4] / first; got > 0.3 {
		t.Errorf("the fifth target took %.2f of the first's damage, want about 0.2 (no bonus)", got)
	}
}
