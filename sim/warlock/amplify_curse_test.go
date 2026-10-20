package warlock

import (
	"math"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
)

// Patch 177: Amplify Curse raises Bane of Agony's whole tick, the base and the spell power share
// (wowsims/forever e2f38aef70). With 1000 spell damage, an amplified first tick is exactly 1.5 of a plain one
// with the switch on; off, only the base is amplified, so the share keeps the ratio below 1.5.
func TestAmplifyCurseWholeTick(t *testing.T) {
	ratio := func(on bool) float64 {
		defer func(old bool) { AmplifyCurseWholeTick = old }(AmplifyCurseWholeTick)
		AmplifyCurseWholeTick = on
		player := &proto.Player{
			Name: "lock", Class: proto.Class_ClassWarlock, Race: proto.Race_RaceOrc, TalentsString: AfflictionTalents,
			Equipment:  &proto.EquipmentSpec{},
			BonusStats: &proto.UnitStats{Stats: stats.Stats{stats.SpellDamage: 1000}.ToProtoArray()},
			Spec: &proto.Player_Warlock{Warlock: &proto.Warlock{Options: &proto.Warlock_Options{ClassOptions: &proto.WarlockOptions{
				Summon: proto.WarlockOptions_Succubus, Armor: proto.WarlockOptions_DemonArmor, CurseOptions: proto.WarlockOptions_Elements}}}},
			Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
		}
		sim := core.NewSim(&proto.RaidSimRequest{
			Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
			SimOptions: &proto.SimOptions{Iterations: 1, RandomSeed: 7},
		}, simsignals.CreateSignals())
		sim.Reset()
		warlock := sim.Raid.Parties[0].Players[0].(WarlockAgent).GetWarlock()
		if !warlock.Talents.AmplifyCurse || warlock.AmplifyCurseAura == nil {
			t.Fatal("AfflictionTalents take no Amplify Curse; pick a build that does")
		}
		target := warlock.CurrentTarget
		agony := warlock.CurseOfAgony
		dot := agony.Dot(target)

		land := func(amplified bool) float64 {
			dot.Deactivate(sim)
			if amplified {
				warlock.AmplifyCurseAura.Activate(sim)
			}
			for i := 0; i < 50 && !dot.IsActive(); i++ {
				agony.SkipCastAndApplyEffects(sim, target)
			}
			if !dot.IsActive() {
				t.Fatal("Bane of Agony did not land")
			}
			return agony.ExpectedTickDamageFromCurrentSnapshot(sim, target)
		}
		plain := land(false)
		return land(true) / plain
	}

	if got := ratio(true); math.Abs(got-1.5) > 1e-9 {
		t.Errorf("switch on: an amplified tick is %.6f of a plain one, want 1.5", got)
	}
	if got := ratio(false); got >= 1.5-1e-6 || got <= 1 {
		t.Errorf("switch off: an amplified tick is %.6f of a plain one, want between 1 and 1.5 (base only)", got)
	}
}
