package retribution

import (
	"fmt"
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/paladin"
)

// The four seals whose judgement only puts up a debuff, by the rank the rotation casts.
// The judgement each names is on its second effect, the Crusader's on its third.
var debuffJudgementSeals = []struct {
	name     string
	sealID   int32
	judgeEff int
}{
	{"Seal of Light", paladin.SealOfLightRankMap.Highest().ID, 2},
	{"Seal of Wisdom", paladin.SealOfWisdomRankMap.Highest().ID, 2},
	{"Seal of Justice", 20164, 2},
	{"Seal of the Crusader", paladin.SealOfTheCrusaderRankMap.Highest().ID, 3},
}

func judgementID(sealID int32, effect int) int32 {
	return int32(spelldata.MustFind(sealID).EffectN(effect).BaseValue())
}

// Keeps one seal up and judges it whenever Judgement is ready, with Arcanite Reaper (enchanted with
// enchant, or bare for 0). Returns the player's metrics.
func runSealAndJudge(t *testing.T, sealID, enchant int32, seed int64) *proto.UnitMetrics {
	t.Helper()
	items := make([]*proto.ItemSpec, 15)
	for i := range items {
		items[i] = &proto.ItemSpec{}
	}
	items[14] = &proto.ItemSpec{Id: 12784, Enchant: enchant}
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceHuman,
		Class:         proto.Class_ClassPaladin,
		Equipment:     &proto.EquipmentSpec{Items: items},
		TalentsString: RetTalents,
		Rotation: core.APLRotationFromJsonString(fmt.Sprintf(`{"type":"TypeAPL","priorityList":[
			{"action":{"castSpell":{"spellId":{"spellId":%d}},"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":%d}}}}}}},
			{"action":{"castSpell":{"spellId":{"spellId":20271}}}}
		]}`, sealID, sealID)),
	}, DefaultOptions)
	result := core.RunRaidSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: seed, Iterations: 200},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	})
	if result.Error != nil {
		t.Fatal(result.Error.Message)
	}
	return result.RaidMetrics.Parties[0].Players[0]
}

// Landed hits: every outcome but a miss, dodge or parry.
func countHits(metrics *proto.UnitMetrics, match func(*proto.ActionMetrics) bool) (n float64) {
	for _, action := range metrics.Actions {
		if !match(action) {
			continue
		}
		for _, target := range action.Targets {
			n += float64(target.Hits + target.Crits + target.Glances + target.Blocks + target.BlockedCrits + target.GlanceBlocks + target.Crushes)
		}
	}
	return n
}

// MythicSim patch 169: the judgements of Light, Wisdom, Justice and the Crusader roll weapon procs
// (Lightclub's sheet: "All judgements can proc weapon chance on hit effects"). Fiery Weapon on Arcanite
// Reaper (6 a minute on a 3.8 speed weapon, 0.38 a hit, less the Fireballs the boss resists) rolls as
// often on a judgement as on a landed swing with the switch on, and never with it off.
func TestDebuffJudgementsTriggerWeaponProcs(t *testing.T) {
	if !paladin.DebuffJudgementsTriggerWeaponProcs {
		t.Fatal("DebuffJudgementsTriggerWeaponProcs is off by default; patch 169 turns it on")
	}
	const (
		fieryWeapon = 803
		fireball    = 13897
	)
	isProc := func(a *proto.ActionMetrics) bool { return a.Id.GetSpellId() == fireball }
	isSwing := func(a *proto.ActionMetrics) bool { return a.Id.GetOtherId() == proto.OtherAction_OtherActionAttack }
	for _, seal := range debuffJudgementSeals {
		t.Run(seal.name, func(t *testing.T) {
			judged := judgementID(seal.sealID, seal.judgeEff)
			judgements := func(m *proto.UnitMetrics) float64 {
				return countHits(m, func(a *proto.ActionMetrics) bool { return a.Id.GetSpellId() == judged })
			}
			on := runSealAndJudge(t, seal.sealID, fieryWeapon, 169)
			var off *proto.UnitMetrics
			withSwitch(&paladin.DebuffJudgementsTriggerWeaponProcs, false, func() {
				off = runSealAndJudge(t, seal.sealID, fieryWeapon, 169)
			})
			procsOn, swingsOn, judgedOn := countHits(on, isProc), countHits(on, isSwing), judgements(on)
			procsOff, swingsOff := countHits(off, isProc), countHits(off, isSwing)
			if judgedOn == 0 || swingsOn == 0 || swingsOff == 0 || procsOff == 0 {
				t.Fatalf("judgements %v, swings %v / %v, procs off %v", judgedOn, swingsOn, swingsOff, procsOff)
			}
			// The procs the swings account for, at the rate they give with the switch off.
			perSwing := procsOff / swingsOff
			perJudgement := (procsOn - perSwing*swingsOn) / judgedOn
			t.Logf("%s: %.0f judgements, %.4f extra procs a judgement, %.4f a landed swing", seal.name, judgedOn, perJudgement, perSwing)
			if math.Abs(perJudgement/perSwing-1) > 0.15 {
				t.Errorf("%s: %.4f extra Fiery Weapon procs a judgement, %.4f a landed swing; want the same", seal.name, perJudgement, perSwing)
			}
		})
	}
}

// MythicSim patch 169: no aura proc hears those judgements ("'Equip' effects only for those that count
// as melee attacks"), so with nothing that rolls a weapon proc worn the switch changes nothing at all:
// not the seal's own effects (Seal of Light's heal, Seal of Wisdom's mana), not Judgement Refresh, not
// a talent. Seal of Light's rate is switched off so its heal fires on every landed melee attack, which
// the judgement would add to if it counted as one.
func TestDebuffJudgementsReachNoAuraProcs(t *testing.T) {
	withSwitch(&paladin.SealOfLightFifteenPPM, false, func() {
		for _, seal := range debuffJudgementSeals {
			on := runSealAndJudge(t, seal.sealID, 0, 1690)
			var off *proto.UnitMetrics
			withSwitch(&paladin.DebuffJudgementsTriggerWeaponProcs, false, func() {
				off = runSealAndJudge(t, seal.sealID, 0, 1690)
			})
			judged := judgementID(seal.sealID, seal.judgeEff)
			if n := countHits(on, func(a *proto.ActionMetrics) bool { return a.Id.GetSpellId() == judged }); n == 0 {
				t.Fatalf("%s: no judgement (%d) landed", seal.name, judged)
			}
			if on.Dps.Avg != off.Dps.Avg || on.Hps.Avg != off.Hps.Avg {
				t.Errorf("%s: dps %.4f / hps %.4f with the switch on, %.4f / %.4f off", seal.name, on.Dps.Avg, on.Hps.Avg, off.Dps.Avg, off.Hps.Avg)
			}
			counts := func(m *proto.UnitMetrics) map[string]int32 {
				c := map[string]int32{}
				for _, a := range m.Actions {
					for _, target := range a.Targets {
						c[a.Id.String()] += target.Casts + target.Hits + target.Crits
					}
				}
				return c
			}
			cOn, cOff := counts(on), counts(off)
			for id, n := range cOn {
				if cOff[id] != n {
					t.Errorf("%s: action %s %d times with the switch on, %d off", seal.name, id, n, cOff[id])
				}
			}
		}
	})
}

func runRet(t *testing.T, talents string, items []*proto.ItemSpec, apl string, seed int64) *proto.UnitMetrics {
	t.Helper()
	player := core.WithSpec(&proto.Player{
		Race:          proto.Race_RaceHuman,
		Class:         proto.Class_ClassPaladin,
		Equipment:     &proto.EquipmentSpec{Items: items},
		TalentsString: talents,
		Rotation:      core.APLRotationFromJsonString(apl),
	}, DefaultOptions)
	result := core.RunRaidSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: seed, Iterations: 200},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	})
	if result.Error != nil {
		t.Fatal(result.Error.Message)
	}
	return result.RaidMetrics.Parties[0].Players[0]
}

func mainHandOnly(id int32) []*proto.ItemSpec {
	items := make([]*proto.ItemSpec, 15)
	for i := range items {
		items[i] = &proto.ItemSpec{}
	}
	items[14] = &proto.ItemSpec{Id: id}
	return items
}

func actionTotals(m *proto.UnitMetrics, spellID int32) (casts, hits, damage float64) {
	for _, a := range m.Actions {
		if a.Id.GetSpellId() != spellID {
			continue
		}
		for _, target := range a.Targets {
			casts += float64(target.Casts)
			hits += float64(target.Hits + target.Crits)
			damage += target.Damage
		}
	}
	return casts, hits, damage
}

// MythicSim patch 169: with Improved Seals talented, the part of each Seal of Command proc that the
// paladin's own Judgement of the Crusader adds takes the talent, so the procs deal more with the switch
// on, the same procs landing (the switch changes no roll). Without the talent the switch changes nothing.
func TestSealOfCommandBonusTakesImprovedSeals(t *testing.T) {
	crusader := paladin.SealOfTheCrusaderRankMap.Highest()
	judgedCrusader := int32(crusader.EffectN(3).BaseValue())
	command := paladin.SealOfCommandRankMap.Highest().ID
	const commandProc = 20424
	apl := fmt.Sprintf(`{"type":"TypeAPL","priorityList":[
		{"action":{"castSpell":{"spellId":{"spellId":20271}},"condition":{"auraIsActive":{"auraId":{"spellId":%d}}}}},
		{"action":{"castSpell":{"spellId":{"spellId":%d}},"condition":{"not":{"val":{"auraIsActive":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":%d}}}}}}},
		{"action":{"castSpell":{"spellId":{"spellId":%d}},"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":%d}}}}}}}
	]}`, crusader.ID, crusader.ID, judgedCrusader, command, command)
	// RetTalents has 3/3 Improved Seals (the Holy tree's fifth talent).
	noImprovedSeals := "550000022001--05225231001330321"
	for _, c := range []struct {
		talents string
		more    bool
	}{{RetTalents, true}, {noImprovedSeals, false}} {
		on := runRet(t, c.talents, WeaponOnly.Items, apl, 1691)
		var off *proto.UnitMetrics
		withSwitch(&paladin.SealOfCommandBonusTakesImprovedSeals, false, func() {
			off = runRet(t, c.talents, WeaponOnly.Items, apl, 1691)
		})
		_, hitsOn, damageOn := actionTotals(on, commandProc)
		_, hitsOff, damageOff := actionTotals(off, commandProc)
		t.Logf("talents %s: %.0f procs, %.1f damage a proc with the switch on, %.1f off", c.talents, hitsOn, damageOn/hitsOn, damageOff/hitsOff)
		if hitsOn == 0 || hitsOn != hitsOff {
			t.Fatalf("talents %s: %v procs landed with the switch on, %v off", c.talents, hitsOn, hitsOff)
		}
		if c.more && !(damageOn > damageOff) {
			t.Errorf("talents %s: Seal of Command dealt %.1f with the switch on, %.1f off; want more on", c.talents, damageOn, damageOff)
		}
		if !c.more && damageOn != damageOff {
			t.Errorf("talents %s (no Improved Seals): Seal of Command dealt %.1f with the switch on, %.1f off; want the same", c.talents, damageOn, damageOff)
		}
	}
}

// MythicSim patch 169: Seal of Light heals at 15 procs a minute (Lightclub's sheet). With a 2.0 speed
// weapon (Quel'Serrar) that is 0.5 a landed swing; off, every landed swing heals.
func TestSealOfLightFifteenPPM(t *testing.T) {
	if !paladin.SealOfLightFifteenPPM {
		t.Fatal("SealOfLightFifteenPPM is off by default; patch 169 turns it on")
	}
	light := paladin.SealOfLightRankMap.Highest()
	heal := light.Refs()[0].ID
	apl := fmt.Sprintf(`{"type":"TypeAPL","priorityList":[
		{"action":{"castSpell":{"spellId":{"spellId":%d}},"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":%d}}}}}}}
	]}`, light.ID, light.ID)
	isSwing := func(a *proto.ActionMetrics) bool { return a.Id.GetOtherId() == proto.OtherAction_OtherActionAttack }
	perSwing := func(m *proto.UnitMetrics) float64 {
		_, heals, _ := actionTotals(m, heal)
		return heals / countHits(m, isSwing)
	}
	on := perSwing(runRet(t, RetTalents, mainHandOnly(18348), apl, 1692))
	var off float64
	withSwitch(&paladin.SealOfLightFifteenPPM, false, func() {
		off = perSwing(runRet(t, RetTalents, mainHandOnly(18348), apl, 1692))
	})
	t.Logf("Seal of Light heals a landed swing: %.4f at 15 a minute, %.4f on every hit", on, off)
	if off < 0.97 || off > 1.0001 {
		t.Errorf("switch off: %.4f heals a landed swing, want 1 (the seal's few down moments aside)", off)
	}
	if math.Abs(on/off-0.5) > 0.04 {
		t.Errorf("switch on: %.4f heals a landed swing against %.4f off, ratio %.3f, want 0.5", on, off, on/off)
	}
}
