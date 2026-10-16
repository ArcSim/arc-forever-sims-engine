package sim

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/stats"
	"github.com/wowsims/forever/sim/druid"
	"github.com/wowsims/forever/sim/hunter"
	"github.com/wowsims/forever/sim/rogue"
	"github.com/wowsims/forever/sim/warlock"
	"github.com/wowsims/forever/sim/warrior"
	"google.golang.org/protobuf/encoding/protojson"
)

// What a dot reads when it lands and what it reads at every tick, ability by ability.
//
// Forever's dots tick on the caster's current stats (patch 22, docs/mythicsim-patches.md): a tick adds
// the attack power or spell power share it has then and takes the damage multiplier it has then. Flat
// damage, stack counts and combo point values stay with the application. Each ability is applied, left
// to tick once, given more attack power, spell power or a bigger damage multiplier, and left to tick
// again. The numbers come from the debug line every tick prints before crits and resistances:
// "MAP: .. RAP: .. SP: .. BaseDamage:.. AfterAttackerMods:..".

type dotProbe struct {
	sim    *core.Simulation
	caster *core.Unit
	dot    *core.Dot
}

func probePlayer(t *testing.T, playerJSON string, mainHand ...int32) (*core.Simulation, core.Agent, *core.Unit) {
	t.Helper()
	player := &proto.Player{}
	if err := protojson.Unmarshal([]byte(playerJSON), player); err != nil {
		t.Fatal(err)
	}
	items := make([]*proto.ItemSpec, proto.ItemSlot_ItemSlotRanged+1)
	for i := range items {
		items[i] = &proto.ItemSpec{}
	}
	if len(mainHand) > 0 {
		items[proto.ItemSlot_ItemSlotMainHand] = &proto.ItemSpec{Id: mainHand[0]}
	}
	player.Equipment = &proto.EquipmentSpec{Items: items}
	player.Rotation = &proto.APLRotation{Type: proto.APLRotation_TypeAPL}
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 1},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  core.MakeSingleTargetEncounter(0),
	}, simsignals.CreateSignals())
	sim.Reset()
	agent := sim.Raid.Parties[0].Players[0]
	return sim, agent, agent.GetCharacter().CurrentTarget
}

const (
	catJSON      = `{"name":"cat","class":"ClassDruid","race":"RaceNightElf","talentsString":"-55210032020132012051-05503","feralCatDruid":{"options":{}}}`
	bearJSON     = `{"name":"bear","class":"ClassDruid","race":"RaceNightElf","talentsString":"-50032302120132010501-0550325","feralBearDruid":{"options":{"classOptions":{}}}}`
	balanceJSON  = `{"name":"owl","class":"ClassDruid","race":"RaceNightElf","talentsString":"5532220115001351--505302","balanceDruid":{"options":{"classOptions":{}}}}`
	rogueJSON    = `{"name":"rogue","class":"ClassRogue","race":"RaceHuman","talentsString":"00530310551021051-302303202004","rogue":{"options":{"classOptions":{}}}}`
	warriorJSON  = `{"name":"warr","class":"ClassWarrior","race":"RaceHuman","talentsString":"32305213132515201-5502","dpsWarrior":{"options":{"classOptions":{}}}}`
	warlockJSON  = `{"name":"lock","class":"ClassWarlock","race":"RaceHuman","talentsString":"2535002013521105--05000551","warlock":{"options":{"classOptions":{}}}}`
	hunterJSON   = `{"name":"hunt","class":"ClassHunter","race":"RaceOrc","talentsString":"-3050552301503151-50024001","hunter":{"options":{"classOptions":{"ammo":"Doomshot","quiverBonus":"Speed15"}}}}`
	dotProbeBump = 1000.0
)

type dotCase struct {
	name string
	// Wants, per point of the stat that rose (the debug line's own MAP, RAP or SP): the share a tick
	// reads at the tick. A zero is a dot that takes none of it after landing.
	meleeAP, rangedAP, spellPower float64
	// The part of a tick that stays with the application: the client's flat tick, a combo point or stack
	// count carried with it. Zero leaves it unchecked.
	flat float64
	// The tick spell ignores caster damage modifiers (Attributes[6] 0x20000000), so a multiplier gained
	// mid-dot leaves it alone: Deep Wounds' 412613 (upstream #645).
	ignoresCasterMods bool
	build             func(t *testing.T) dotProbe
}

func dotCases() []dotCase {
	return []dotCase{
		// Rake 34 a tick plus 5.26% of attack power, read at the tick (mythicsim patch 99: the client row
		// has no attack power coefficient, the share is fitted to beta logs).
		{name: "Rake", meleeAP: 0.05, flat: 34, build: func(t *testing.T) dotProbe {
			sim, agent, target := probePlayer(t, catJSON)
			cat := agent.(druid.DruidAgent).GetDruid()
			dot := cat.Rake.Dot(target)
			dot.Apply(sim)
			return dotProbe{sim, &cat.Unit, dot}
		}},
		// Rip rank 6: 15 plus 25.5 a combo point, kept from the 5 points it was cast with, and 1% of attack
		// power a point at four points at most.
		{name: "Rip", meleeAP: 0.04, flat: 15 + 25.5*5, build: func(t *testing.T) dotProbe {
			sim, agent, target := probePlayer(t, catJSON)
			cat := agent.(druid.DruidAgent).GetDruid()
			cat.AddComboPoints(sim, 5, cat.Rake.ComboPointMetrics())
			dot := cat.Rip.Dot(target)
			dot.Apply(sim)
			// The cast spends the points after the dot lands; the dot keeps what they were worth.
			cat.SpendComboPoints(sim, cat.Rip.ComboPointMetrics())
			return dotProbe{sim, &cat.Unit, dot}
		}},
		// Lacerate 15 a stack, the three stacks it was applied with, and no attack power share.
		{name: "Lacerate", flat: 15 * 3, build: func(t *testing.T) dotProbe {
			sim, agent, target := probePlayer(t, bearJSON)
			bear := agent.(druid.DruidAgent).GetDruid()
			dot := bear.Lacerate.Dot(target)
			dot.Apply(sim)
			dot.SetStacks(sim, 3)
			dot.TakeSnapshot(sim)
			return dotProbe{sim, &bear.Unit, dot}
		}},
		{name: "Moonfire", spellPower: 0.13, flat: 60, build: func(t *testing.T) dotProbe {
			sim, agent, target := probePlayer(t, balanceJSON)
			owl := agent.(druid.DruidAgent).GetDruid()
			dot := owl.Moonfire.RelatedDotSpell.Dot(target)
			dot.Apply(sim)
			return dotProbe{sim, &owl.Unit, dot}
		}},
		{name: "Insect Swarm", spellPower: 0.158, flat: 31, build: func(t *testing.T) dotProbe {
			sim, agent, target := probePlayer(t, balanceJSON)
			owl := agent.(druid.DruidAgent).GetDruid()
			if owl.InsectSwarm == nil {
				t.Skip("the Balance test build has no Insect Swarm")
			}
			dot := owl.InsectSwarm.Dot(target)
			dot.Apply(sim)
			return dotProbe{sim, &owl.Unit, dot}
		}},
		{name: "Garrote", meleeAP: 0.03, flat: 92, build: func(t *testing.T) dotProbe {
			sim, agent, target := probePlayer(t, rogueJSON)
			r := agent.(rogue.RogueAgent).GetRogue()
			dot := r.Garrote.Dot(target)
			dot.Apply(sim)
			return dotProbe{sim, &r.Unit, dot}
		}},
		// Rupture rank 6: 35 plus 4.73 a combo point, kept from the 5 points it was cast with.
		{name: "Rupture", meleeAP: 0.03, flat: 35 + 4.73*5, build: func(t *testing.T) dotProbe {
			sim, agent, target := probePlayer(t, rogueJSON)
			r := agent.(rogue.RogueAgent).GetRogue()
			r.AddComboPoints(sim, 5, r.Garrote.ComboPointMetrics())
			dot := r.Rupture.Dot(target)
			dot.BaseTickCount = 3 + 5 // Rupture's cast sets its length from the combo points
			dot.Apply(sim)
			r.SpendComboPoints(sim, r.Rupture.ComboPointMetrics())
			return dotProbe{sim, &r.Unit, dot}
		}},
		{name: "Rend", meleeAP: 0.02, flat: 21, build: func(t *testing.T) dotProbe {
			sim, agent, target := probePlayer(t, warriorJSON)
			w := agent.(warrior.WarriorAgent).GetWarrior()
			dot := w.Rend.Dot(target)
			dot.Apply(sim)
			return dotProbe{sim, &w.Unit, dot}
		}},
		{name: "Deep Wounds", ignoresCasterMods: true, build: func(t *testing.T) dotProbe {
			sim, agent, target := probePlayer(t, warriorJSON, 7230) // Smite's Mighty Hammer
			w := agent.(warrior.WarriorAgent).GetWarrior()
			if w.Talents.DeepWounds == 0 {
				t.Skip("the test build has no Deep Wounds")
			}
			w.DeepWounds.Cast(sim, target)
			return dotProbe{sim, &w.Unit, w.DeepWounds.Dot(target)}
		}},
		{name: "Corruption", spellPower: 0.2, flat: 73, build: func(t *testing.T) dotProbe {
			sim, agent, target := probePlayer(t, warlockJSON)
			lock := agent.(warlock.WarlockAgent).GetWarlock()
			dot := lock.Corruption.Dot(target)
			dot.Apply(sim)
			return dotProbe{sim, &lock.Unit, dot}
		}},
		{name: "Serpent Sting", rangedAP: 0.035, spellPower: 0.2, flat: 111, build: func(t *testing.T) dotProbe {
			sim, agent, target := probePlayer(t, hunterJSON)
			h := agent.(hunter.HunterAgent).GetHunter()
			dot := h.SerpentSting.Dot(target)
			dot.Apply(sim)
			return dotProbe{sim, &h.Unit, dot}
		}},
	}
}

type tickReading struct{ meleeAP, rangedAP, spellPower, base, multiplier float64 }

// tick runs the probe until its dot has ticked n times and returns the debug reading of the last one.
// change runs once, after the first tick.
func (p dotProbe) ticks(t *testing.T, n int, change func()) []tickReading {
	t.Helper()
	line := regexp.MustCompile(regexp.QuoteMeta(p.dot.Spell.ActionID.String()) +
		` \[DEBUG\] MAP: (-?[0-9.]+), RAP: (-?[0-9.]+), SP: (-?[0-9.]+), BaseDamage:(-?[0-9.]+), AfterAttackerMods:(-?[0-9.]+)`)
	var readings []tickReading
	p.sim.Log = func(format string, args ...interface{}) {
		if m := line.FindStringSubmatch(fmt.Sprintf(format, args...)); m != nil {
			f := func(i int) float64 { v, _ := strconv.ParseFloat(m[i], 64); return v }
			r := tickReading{meleeAP: f(1), rangedAP: f(2), spellPower: f(3), base: f(4)}
			r.multiplier = f(5) / r.base
			readings = append(readings, r)
		}
	}
	changed := false
	for steps := 0; len(readings) < n; steps++ {
		if steps > 1_000_000 || p.sim.Step() {
			t.Fatalf("%s ticked %d times of %d", p.dot.Spell.ActionID, len(readings), n)
		}
		if len(readings) >= 1 && !changed && change != nil {
			changed = true
			change()
		}
	}
	return readings
}

func TestDotsReadStatsAtTheTick(t *testing.T) {
	for _, c := range dotCases() {
		t.Run(c.name, func(t *testing.T) {
			// Nothing changes: two ticks of the same dot deal the same base damage.
			still := c.build(t).ticks(t, 2, nil)
			if math.Abs(still[1].base-still[0].base) > 0.11 || math.Abs(still[1].multiplier-still[0].multiplier) > 1e-3 {
				t.Fatalf("a dot with nothing changed ticked %.2f then %.2f (multiplier %.4f then %.4f)", still[0].base, still[1].base, still[0].multiplier, still[1].multiplier)
			}

			// The flat part plus the shares of the stats the caster had when the first tick landed.
			if c.flat != 0 {
				first := still[0]
				want := c.flat + c.meleeAP*first.meleeAP + c.rangedAP*first.rangedAP + c.spellPower*first.spellPower
				if math.Abs(first.base-want) > 0.3 {
					t.Errorf("the first tick's base damage is %.2f, want %.2f (flat %.2f plus the shares of the caster's stats)", first.base, want, c.flat)
				}
			}

			// More attack power and spell power after the dot landed.
			p := c.build(t)
			power := p.ticks(t, 2, func() {
				p.caster.AddStatDynamic(p.sim, stats.AttackPower, dotProbeBump)
				p.caster.AddStatDynamic(p.sim, stats.RangedAttackPower, dotProbeBump)
				p.caster.AddStatDynamic(p.sim, stats.SpellDamage, dotProbeBump)
			})
			t.Logf("%s: base %.2f -> %.2f after +%.0f attack power, ranged attack power and spell power (MAP %.0f -> %.0f, RAP %.0f -> %.0f, SP %.0f -> %.0f)",
				c.name, power[0].base, power[1].base, dotProbeBump, power[0].meleeAP, power[1].meleeAP, power[0].rangedAP, power[1].rangedAP, power[0].spellPower, power[1].spellPower)
			// The three stats rose together, so what a tick adds is the sum of the shares that apply.
			wantRise := (c.meleeAP + c.rangedAP + c.spellPower) * dotProbeBump
			if rise := power[1].base - power[0].base; math.Abs(rise-wantRise) > 0.5+0.002*wantRise {
				t.Errorf("the tick rose by %.2f after the stats rose by %.0f, want %.2f (melee AP share %.3f, ranged AP share %.3f, spell power share %.3f)",
					rise, dotProbeBump, wantRise, c.meleeAP, c.rangedAP, c.spellPower)
			}

			// A bigger damage multiplier after the dot landed.
			p = c.build(t)
			multiplier := p.ticks(t, 2, func() { p.caster.PseudoStats.DamageDealtMultiplier *= 1.5 })
			t.Logf("%s: multiplier %.4f -> %.4f after a 1.5x damage multiplier gained mid-dot", c.name, multiplier[0].multiplier, multiplier[1].multiplier)
			wantMultiplier := 1.5
			if c.ignoresCasterMods {
				wantMultiplier = 1
			}
			if got := multiplier[1].multiplier / multiplier[0].multiplier; math.Abs(got-wantMultiplier) > 0.01 {
				t.Errorf("the tick's multiplier rose %.3fx, want %.1fx", got, wantMultiplier)
			}
			if math.Abs(multiplier[1].base-multiplier[0].base) > 0.11 {
				t.Errorf("a multiplier change moved the base damage %.2f -> %.2f", multiplier[0].base, multiplier[1].base)
			}
		})
	}
}
