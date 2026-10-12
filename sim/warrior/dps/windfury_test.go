package dps

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/buffs"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

// Windfury Totem's proc is a party aura since build 70009 (10612): a melee hit
// that lands and deals damage can fire it, auto or special (the row's
// ProcTypeMask 0x14); nothing ranged and no spell can. With
// buffs.WindfuryTotemMainHandOnly on (the tooltip's "Each main hand hit", the default) no
// off-hand weapon's hit can, though a shield strike (an off-hand special with a
// shield) still can: the client swings it as a main-hand attack. That is the
// default; off, either hand rolls it as the row alone has it. An auto that procs it has spent the first
// charge itself and leaves one for the extra attack; a special leaves both.
// Its internal cooldown is the row's 100 ms, not Era's 1.5 sec.
func TestWindfuryTotemProcsOffMeleeHitsWithTheChargesTheTriggerLeaves(t *testing.T) {
	setupWith := func(t *testing.T, gear *proto.EquipmentSpec) (*core.Simulation, *core.Unit, *core.Aura, *[]int32) {
		t.Helper()

		player := &proto.Player{
			Name:          "Windfury",
			Race:          proto.Race_RaceOrc,
			Class:         proto.Class_ClassWarrior,
			Equipment:     gear,
			TalentsString: FuryTalents,
			Spec:          DefaultOptions,
			Rotation:      core.GetAplRotation("../../../ui/specs/warrior/dps/apls", "dps_reck").Rotation,
		}
		sim := core.NewSim(&proto.RaidSimRequest{
			Raid: core.SinglePlayerRaidProto(player, &proto.PartyBuffs{WindfuryTotem: true}, &proto.RaidBuffs{},
				&proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
			SimOptions: &proto.SimOptions{RandomSeed: 1},
		}, simsignals.CreateSignals())
		sim.Reset()

		war := sim.Raid.Parties[0].Players[0].(*DpsWarrior)
		trigger := war.GetAura("Windfury Totem Trigger")
		proc := war.GetAura("Windfury Totem (External)")
		if trigger == nil || proc == nil {
			t.Fatalf("trigger %v, proc %v; want both registered", trigger, proc)
		}
		if !trigger.IsActive() {
			t.Fatal("the totem's trigger is not up at the pull")
		}

		var started []int32
		proc.ApplyOnStacksChange(func(_ *core.Aura, _ *core.Simulation, oldStacks int32, newStacks int32) {
			if oldStacks == 0 {
				started = append(started, newStacks)
			}
		})
		return sim, sim.Encounter.AllTargetUnits[0], trigger, &started
	}
	setup := func(t *testing.T) (*core.Simulation, *core.Unit, *core.Aura, *[]int32) {
		t.Helper()
		return setupWith(t, DualWieldGear.GearSet)
	}

	// Offers the trigger up to 200 landed hits of one kind and answers whether any of them procced.
	offer := func(sim *core.Simulation, target *core.Unit, trigger *core.Aura, started *[]int32, procMask core.ProcMask, damage float64) bool {
		spell := &core.Spell{ProcMask: procMask}
		for i := 0; i < 200 && len(*started) == 0; i++ {
			trigger.OnSpellHitDealt(trigger, sim, spell, &core.SpellResult{Target: target, Outcome: core.OutcomeHit, Damage: damage})
		}
		return len(*started) > 0
	}

	for _, row := range []struct {
		name     string
		procMask core.ProcMask
		charges  int32
	}{
		{"a main-hand auto leaves one charge", core.ProcMaskMeleeMHAuto, 1},
		{"a main-hand special leaves both charges", core.ProcMaskMeleeMHSpecial, 2},
	} {
		t.Run(row.name, func(t *testing.T) {
			sim, target, trigger, started := setup(t)
			if !offer(sim, target, trigger, started, row.procMask, 100) {
				t.Fatal("200 landed hits never procced the totem")
			}
			if got := (*started)[0]; got != row.charges {
				t.Errorf("the proc started with %d charges, want %d", got, row.charges)
			}
		})
	}

	for _, row := range []struct {
		name     string
		procMask core.ProcMask
	}{
		{"an off-hand auto", core.ProcMaskMeleeOHAuto},
		{"an off-hand special", core.ProcMaskMeleeOHSpecial},
		{"a ranged auto", core.ProcMaskRangedAuto},
		{"a ranged special", core.ProcMaskRangedSpecial},
		{"a spell", core.ProcMaskSpellDamage},
	} {
		t.Run(row.name+" never procs it", func(t *testing.T) {
			if row.procMask&core.ProcMaskMeleeOH != 0 {
				defer func(old bool) { buffs.WindfuryTotemMainHandOnly = old }(buffs.WindfuryTotemMainHandOnly)
				buffs.WindfuryTotemMainHandOnly = true
			}
			sim, target, trigger, started := setup(t)
			if offer(sim, target, trigger, started, row.procMask, 100) {
				t.Errorf("%s procced the totem", row.name)
			}
		})
	}

	// buffs.WindfuryTotemMainHandOnly set false: the off hand rolls it as the
	// row alone would have it.
	for _, row := range []struct {
		name     string
		procMask core.ProcMask
		charges  int32
	}{
		{"an off-hand auto leaves one charge", core.ProcMaskMeleeOHAuto, 1},
		{"an off-hand special leaves both charges", core.ProcMaskMeleeOHSpecial, 2},
	} {
		t.Run("either hand: "+row.name, func(t *testing.T) {
			defer func(old bool) { buffs.WindfuryTotemMainHandOnly = old }(buffs.WindfuryTotemMainHandOnly)
			buffs.WindfuryTotemMainHandOnly = false
			sim, target, trigger, started := setup(t)
			if !offer(sim, target, trigger, started, row.procMask, 100) {
				t.Fatal("200 landed hits never procced the totem")
			}
			if got := (*started)[0]; got != row.charges {
				t.Errorf("the proc started with %d charges, want %d", got, row.charges)
			}
		})
	}

	t.Run("a shield strike leaves both charges", func(t *testing.T) {
		sim, target, trigger, started := setupWith(t, weaponsOnly(15240, 18756)) // Dreadguard's Protector
		if !offer(sim, target, trigger, started, core.ProcMaskMeleeOHSpecial, 100) {
			t.Fatal("200 landed shield strikes never procced the totem")
		}
		if got := (*started)[0]; got != 2 {
			t.Errorf("the proc started with %d charges, want 2", got)
		}
	})

	// buffs.WindfuryProcsFromSpecials set false: only auto attacks roll it.
	t.Run("autos only: a main-hand special never procs it", func(t *testing.T) {
		defer func(old bool) { buffs.WindfuryProcsFromSpecials = old }(buffs.WindfuryProcsFromSpecials)
		buffs.WindfuryProcsFromSpecials = false
		sim, target, trigger, started := setup(t)
		if offer(sim, target, trigger, started, core.ProcMaskMeleeMHSpecial, 100) {
			t.Error("a main-hand special procced the totem")
		}
		sim, target, trigger, started = setup(t)
		if !offer(sim, target, trigger, started, core.ProcMaskMeleeMHAuto, 100) {
			t.Error("200 landed main-hand autos never procced the totem")
		}
	})

	t.Run("the internal cooldown is 100 ms", func(t *testing.T) {
		_, _, trigger, _ := setup(t)
		if trigger.Icd == nil || trigger.Icd.Duration != 100*time.Millisecond {
			t.Errorf("the trigger's internal cooldown is %v, want 100ms", trigger.Icd)
		}
	})
}
