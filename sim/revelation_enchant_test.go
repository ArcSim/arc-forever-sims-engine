package sim

import (
	"math"
	"testing"
	"time"

	"github.com/wowsims/forever/sim/common/forever"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
	"google.golang.org/protobuf/encoding/protojson"
)

// Enchant Weapon - Revelation (enchant 8217, mythicsim patch 157): a chance on each non-crit
// non-periodic spell hit to gain the class's one-charge buff, +100% crit on the next spell its class
// mask names; a channel holds the charge to its end.

const (
	revelationEnchant = 8217
	ironbarkStaff     = 20069

	revelationMageJSON    = `{"name":"m","class":"ClassMage","race":"RaceHuman","talentsString":"0502105003-13550000130135031-003","mage":{"options":{"classOptions":{}}}}`
	revelationWarlockJSON = `{"name":"l","class":"ClassWarlock","race":"RaceHuman","talentsString":"254-0050203001-155030510310005","warlock":{"options":{"classOptions":{}}}}`
	revelationShamanJSON  = `{"name":"s","class":"ClassShaman","race":"RaceOrc","talentsString":"5501331020011351-052-05305","elementalShaman":{"options":{"classOptions":{}}}}`
	revelationPriestJSON  = `{"name":"p","class":"ClassPriest","race":"RaceHuman","talentsString":"02530301103--50052201221302251","dpsPriest":{"options":{"classOptions":{}}}}`
	revelationDruidJSON   = `{"name":"d","class":"ClassDruid","race":"RaceNightElf","talentsString":"5532220115501351--505","balanceDruid":{"options":{"classOptions":{}}}}`
)

func revelationPlayer(t *testing.T, playerJSON string, enchant int32) *proto.Player {
	t.Helper()
	player := &proto.Player{}
	if err := protojson.Unmarshal([]byte(playerJSON), player); err != nil {
		t.Fatal(err)
	}
	items := make([]*proto.ItemSpec, proto.ItemSlot_ItemSlotRanged+1)
	for i := range items {
		items[i] = &proto.ItemSpec{}
	}
	items[proto.ItemSlot_ItemSlotMainHand] = &proto.ItemSpec{Id: ironbarkStaff, Enchant: enchant}
	player.Equipment = &proto.EquipmentSpec{Items: items}
	player.Consumables = &proto.ConsumesSpec{}
	if player.Rotation == nil {
		player.Rotation = &proto.APLRotation{Type: proto.APLRotation_TypeAPL}
	}
	return player
}

func revelationSim(t *testing.T, player *proto.Player, targets int) (*core.Simulation, *core.Character) {
	t.Helper()
	encounter := core.MakeSingleTargetEncounter(0)
	encounter.Duration = 300
	for len(encounter.Targets) < targets {
		encounter.Targets = append(encounter.Targets, core.NewDefaultTarget())
	}
	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{RandomSeed: 3},
		Raid:       core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
		Encounter:  encounter,
	}, simsignals.CreateSignals())
	sim.Reset()
	return sim, sim.Raid.Parties[0].Players[0].GetCharacter()
}

func revelationAdvance(t *testing.T, sim *core.Simulation, d time.Duration) {
	t.Helper()
	end := sim.CurrentTime + d
	for steps := 0; sim.CurrentTime < end; steps++ {
		if steps > 1_000_000 || sim.Step() {
			t.Fatalf("the sim ended at %s before %s", sim.CurrentTime, end)
		}
	}
}

// The highest registered rank of a spell by its client name.
func revelationSpell(t *testing.T, character *core.Character, name string) *core.Spell {
	t.Helper()
	var best *core.Spell
	for _, spell := range character.Spellbook {
		if spell.ActionID.SpellID == 0 || spelldata.Find(spell.ActionID.SpellID).Name != name {
			continue
		}
		if best == nil || spell.ActionID.SpellID > best.ActionID.SpellID {
			best = spell
		}
	}
	if best == nil {
		t.Fatalf("%s is not registered", name)
	}
	return best
}

func revelationBuff(character *core.Character) *core.Aura {
	return character.GetAura("Revelation")
}

// Sets a fixed chance per eligible non-crit hit (the crit-scaled chance off) for one test.
func revelationFixedChance(t *testing.T, chance float64) {
	t.Helper()
	savedChance, savedSwitch := forever.RevelationTriggerChance, forever.RevelationChanceFallsWithCrit
	forever.RevelationTriggerChance, forever.RevelationChanceFallsWithCrit = chance, false
	t.Cleanup(func() {
		forever.RevelationTriggerChance, forever.RevelationChanceFallsWithCrit = savedChance, savedSwitch
	})
}

// MythicSim patch 170: the chance per eligible non-crit hit is 0.098 x (1 - spell crit) with
// RevelationChanceFallsWithCrit on (the default), the constant 0.0885 off.
func TestRevelationChanceFallsWithCrit(t *testing.T) {
	if !forever.RevelationChanceFallsWithCrit {
		t.Fatal("RevelationChanceFallsWithCrit is off by default; patch 170 turns it on")
	}
	saved := forever.RevelationChanceFallsWithCrit
	defer func() { forever.RevelationChanceFallsWithCrit = saved }()
	for _, c := range []struct{ crit, on float64 }{
		{-0.1, 0.098}, {0, 0.098}, {0.075, 0.09065}, {0.13, 0.08526}, {0.25, 0.0735}, {1.2, 0},
	} {
		forever.RevelationChanceFallsWithCrit = true
		if got := forever.RevelationChance(c.crit); math.Abs(got-c.on) > 1e-9 {
			t.Errorf("on, crit %.3f: chance %.5f, want %.5f", c.crit, got, c.on)
		}
		forever.RevelationChanceFallsWithCrit = false
		if got := forever.RevelationChance(c.crit); got != 0.0885 {
			t.Errorf("off, crit %.3f: chance %.5f, want 0.0885", c.crit, got)
		}
	}
}

// Each class takes its own buff from the client: 1248808 for the mage, 1323377 to 1323419 for the
// rest; a warrior has none. Switched off, nothing registers.
func TestRevelationRegistersTheClassBuff(t *testing.T) {
	for _, c := range []struct {
		name, json string
		buff       int32
	}{
		{"mage", revelationMageJSON, 1248808},
		{"warlock", revelationWarlockJSON, 1323392},
		{"shaman", revelationShamanJSON, 1323418},
		{"priest", revelationPriestJSON, 1323377},
		{"druid", revelationDruidJSON, 1323400},
		{"warrior", warriorJSON, 0},
	} {
		_, character := revelationSim(t, revelationPlayer(t, c.json, revelationEnchant), 1)
		aura := revelationBuff(character)
		if c.buff == 0 {
			if aura != nil {
				t.Errorf("%s: has a Revelation buff", c.name)
			}
			continue
		}
		if aura == nil {
			t.Fatalf("%s: no Revelation buff", c.name)
		}
		if aura.ActionID.SpellID != c.buff {
			t.Errorf("%s: buff %d, want %d", c.name, aura.ActionID.SpellID, c.buff)
		}
		if aura.Duration != 15*time.Second {
			t.Errorf("%s: buff lasts %s, want 15s", c.name, aura.Duration)
		}

		_, without := revelationSim(t, revelationPlayer(t, c.json, 0), 1)
		if revelationBuff(without) != nil {
			t.Errorf("%s: a Revelation buff without the enchant", c.name)
		}
	}

	forever.RevelationEnchantModelled = false
	defer func() { forever.RevelationEnchantModelled = true }()
	_, character := revelationSim(t, revelationPlayer(t, revelationMageJSON, revelationEnchant), 1)
	if revelationBuff(character) != nil {
		t.Error("switched off, the mage still has a Revelation buff")
	}
}

// The +100% crit reaches the spells the class mask names and no others.
func TestRevelationCritFollowsTheClassMask(t *testing.T) {
	for _, c := range []struct {
		json string
		on   []string
		off  []string
	}{
		{revelationMageJSON, []string{"Frostbolt", "Fire Blast", "Scorch", "Arcane Missile", "Arcane Missiles", "Blizzard", "Arcane Explosion", "Flamestrike"}, []string{"Fireball", "Pyroblast"}},
		{revelationWarlockJSON, []string{"Shadow Bolt", "Searing Pain", "Hellfire", "Rain of Fire", "Conflagrate"}, []string{"Immolate", "Corruption"}},
		{revelationShamanJSON, []string{"Lightning Bolt", "Chain Lightning", "Earth Shock", "Frost Shock"}, []string{"Flame Shock"}},
		{revelationPriestJSON, []string{"Mind Blast", "Mind Flay", "Smite"}, []string{"Shadow Word: Pain", "Holy Fire"}},
		{revelationDruidJSON, []string{"Wrath", "Starfire", "Hurricane"}, []string{"Moonfire", "Insect Swarm"}},
	} {
		sim, character := revelationSim(t, revelationPlayer(t, c.json, revelationEnchant), 1)
		before := map[*core.Spell]float64{}
		for _, spell := range character.Spellbook {
			before[spell] = spell.BonusCritPercent
		}
		revelationBuff(character).Activate(sim)
		check := func(name string, want float64) {
			for _, spell := range character.Spellbook {
				if spell.ActionID.SpellID == 0 || spelldata.Find(spell.ActionID.SpellID).Name != name {
					continue
				}
				if got := spell.BonusCritPercent - before[spell]; got != want {
					t.Errorf("%s %s (%d): Revelation adds %v crit, want %v", character.Class, name, spell.ActionID.SpellID, got, want)
				}
			}
			revelationSpell(t, character, name)
		}
		for _, name := range c.on {
			check(name, 100)
		}
		for _, name := range c.off {
			check(name, 0)
		}
		revelationBuff(character).Deactivate(sim)
		for spell, crit := range before {
			if spell.BonusCritPercent != crit {
				t.Errorf("%s %s: crit %v after the buff fell, want %v", character.Class, spell.ActionID, spell.BonusCritPercent, crit)
			}
		}
	}
}

// Which hits trigger it, with the chance set to 1 and crits suppressed so every hit is a non-crit.
func TestRevelationTriggersOnNonPeriodicSpellHits(t *testing.T) {
	revelationFixedChance(t, 1)

	for _, c := range []struct {
		json  string
		spell string
		wait  time.Duration
		want  bool
	}{
		{revelationMageJSON, "Frostbolt", 2 * time.Second, true},
		{revelationMageJSON, "Scorch", time.Second, true},
		{revelationMageJSON, "Arcane Explosion", time.Second, true},
		// A spell with a dot part: Immolate's and Flame Shock's direct hits triggered 1 time in 45.
		{revelationMageJSON, "Fireball", 4 * time.Second, false},
		{revelationWarlockJSON, "Shadow Bolt", 2 * time.Second, true},
		{revelationWarlockJSON, "Searing Pain", time.Second, true},
		{revelationWarlockJSON, "Immolate", 6 * time.Second, false},
		// Dot ticks and periodic channels never trigger it.
		{revelationWarlockJSON, "Corruption", 10 * time.Second, false},
		{revelationWarlockJSON, "Drain Life", 3 * time.Second, false},
		// Hellfire's ticks are Hellfire Effect hits in the game.
		{revelationWarlockJSON, "Hellfire", 3 * time.Second, true},
		{revelationShamanJSON, "Lightning Bolt", 2 * time.Second, true},
		{revelationShamanJSON, "Flame Shock", 8 * time.Second, false},
		{revelationPriestJSON, "Mind Flay", 2500 * time.Millisecond, false},
		{revelationPriestJSON, "Mind Blast", time.Second, true},
	} {
		sim, character := revelationSim(t, revelationPlayer(t, c.json, revelationEnchant), 1)
		character.AddStatDynamic(sim, stats.SpellCritPercent, -1000)
		spell := revelationSpell(t, character, c.spell)
		spell.SkipCastAndApplyEffects(sim, character.CurrentTarget)
		revelationAdvance(t, sim, c.wait)
		if got := revelationBuff(character).IsActive(); got != c.want {
			t.Errorf("%s: Revelation up %v after %s, want %v", c.spell, got, c.wait, c.want)
		}
	}
}

// An instant spell's own non-crit hit does not spend the charge it just gained (the buff comes a batch
// window after the hit, after the cast completes); the next spell the mask names takes it.
func TestRevelationIsSpentByTheNextCast(t *testing.T) {
	revelationFixedChance(t, 1)

	sim, character := revelationSim(t, revelationPlayer(t, revelationMageJSON, revelationEnchant), 1)
	buff := revelationBuff(character)
	target := character.CurrentTarget
	character.AddStatDynamic(sim, stats.SpellCritPercent, -1000)

	fireBlast := revelationSpell(t, character, "Fire Blast")
	if !fireBlast.Cast(sim, target) {
		t.Fatal("Fire Blast was not cast")
	}
	revelationAdvance(t, sim, 2*time.Second)
	if !buff.IsActive() {
		t.Fatal("a non-crit Fire Blast left no Revelation")
	}

	// Fireball is not on the mage's mask: it neither takes nor spends the charge.
	fireball := revelationSpell(t, character, "Fireball")
	if !fireball.Cast(sim, target) {
		t.Fatal("Fireball was not cast")
	}
	revelationAdvance(t, sim, 4*time.Second)
	if !buff.IsActive() {
		t.Fatal("Fireball spent Revelation")
	}

	// Frostbolt is: its cast takes the charge.
	character.AddStatDynamic(sim, stats.SpellCritPercent, 1000)
	frostbolt := revelationSpell(t, character, "Frostbolt")
	buff.Activate(sim) // refresh: 15 s from here
	critsBefore := frostbolt.SpellMetrics[target.UnitIndex].Crits
	if !frostbolt.Cast(sim, target) {
		t.Fatal("Frostbolt was not cast")
	}
	revelationAdvance(t, sim, 4*time.Second)
	if buff.IsActive() {
		t.Error("Frostbolt did not spend Revelation")
	}
	if frostbolt.SpellMetrics[target.UnitIndex].Crits != critsBefore+1 {
		t.Error("the Frostbolt cast with Revelation up did not crit")
	}
}

// A channel cast with Revelation up crits on every tick and the buff falls when the channel ends.
func TestRevelationHoldsThroughAChannel(t *testing.T) {
	revelationFixedChance(t, 0)

	for _, c := range []struct {
		json, channel, tick string
		targets             int
	}{
		{revelationMageJSON, "Blizzard", "Blizzard", 3},
		{revelationMageJSON, "Arcane Missiles", "Arcane Missiles", 1},
		{revelationWarlockJSON, "Hellfire", "Hellfire", 3},
	} {
		sim, character := revelationSim(t, revelationPlayer(t, c.json, revelationEnchant), c.targets)
		buff := revelationBuff(character)
		var channel *core.Spell
		for _, spell := range character.Spellbook {
			if spell.Flags.Matches(core.SpellFlagChanneled) && spelldata.Find(spell.ActionID.SpellID).Name == c.channel &&
				(channel == nil || spell.ActionID.SpellID > channel.ActionID.SpellID) {
				channel = spell
			}
		}
		if channel == nil {
			t.Fatalf("no %s channel", c.channel)
		}
		var ticks []*core.Spell
		for _, spell := range character.Spellbook {
			if spell.ActionID.SpellID != 0 && spelldata.Find(spell.ActionID.SpellID).Name == c.tick && (spell == channel || !spell.Flags.Matches(core.SpellFlagChanneled)) {
				ticks = append(ticks, spell)
			}
		}
		crits := func() (hits, crits int32) {
			for _, spell := range ticks {
				for _, m := range spell.SpellMetrics {
					hits += m.Hits + m.Ticks
					crits += m.Crits + m.CritTicks
				}
			}
			return
		}

		buff.Activate(sim)
		if !channel.Cast(sim, character.CurrentTarget) {
			t.Fatalf("%s was not cast", c.channel)
		}
		revelationAdvance(t, sim, 2500*time.Millisecond)
		if !buff.IsActive() {
			t.Fatalf("%s spent Revelation at the start of the channel", c.channel)
		}
		for character.ChanneledDot != nil {
			revelationAdvance(t, sim, 100*time.Millisecond)
		}
		revelationAdvance(t, sim, 100*time.Millisecond)
		if buff.IsActive() {
			t.Errorf("%s: Revelation still up after the channel ended", c.channel)
		}
		hits, critCount := crits()
		// Hellfire's channel is 15 sec, as long as the buff: its last tick lands as the buff runs out.
		if c.channel == "Hellfire" && hits <= int32(c.targets) {
			hits = 0
		}
		if critCount == 0 || hits != 0 {
			t.Errorf("%s: %d ticks crit and %d did not, want every tick to crit", c.channel, critCount, hits)
		}
	}
}

// The rate, end to end: Frostbolt only, gains per non-crit Frostbolt hit near the chance in force
// (0.098 x (1 - the mage's spell crit) with RevelationChanceFallsWithCrit on, 0.0885 off), and every
// gain makes the next Frostbolt crit.
func TestRevelationTriggerRate(t *testing.T) {
	saved := forever.RevelationChanceFallsWithCrit
	defer func() { forever.RevelationChanceFallsWithCrit = saved }()
	_, mage := revelationSim(t, revelationPlayer(t, revelationMageJSON, revelationEnchant), 1)
	crit := mage.GetStat(stats.SpellCritPercent) / 100
	for _, on := range []bool{false, true} {
		forever.RevelationChanceFallsWithCrit = on
		p := forever.RevelationChance(crit)
		if on && math.Abs(p-0.098*(1-crit)) > 1e-12 || !on && p != 0.0885 {
			t.Fatalf("switch %v: chance %v at %.4f crit", on, p, crit)
		}
		t.Run(map[bool]string{false: "constant", true: "falls-with-crit"}[on], func(t *testing.T) {
			revelationTriggerRate(t, p)
		})
	}
}

func revelationTriggerRate(t *testing.T, p float64) {
	const iterations = 1000
	run := func(enchant int32) (casts, hits, crits, procs float64) {
		p := revelationPlayer(t, revelationMageJSON, enchant)
		p.Rotation = &proto.APLRotation{Type: proto.APLRotation_TypeAPL, PriorityList: []*proto.APLListItem{{
			Action: &proto.APLAction{Action: &proto.APLAction_CastSpell{CastSpell: &proto.APLActionCastSpell{
				SpellId: &proto.ActionID{RawId: &proto.ActionID_SpellId{SpellId: 116}},
			}}},
		}}}
		encounter := core.MakeSingleTargetEncounter(0)
		encounter.Duration = 60
		result := core.RunRaidSim(&proto.RaidSimRequest{
			SimOptions: &proto.SimOptions{Iterations: iterations, RandomSeed: 11},
			Raid:       core.SinglePlayerRaidProto(p, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{}),
			Encounter:  encounter,
		})
		if result.Error != nil {
			t.Fatal(result.Error.Message)
		}
		metrics := result.RaidMetrics.Parties[0].Players[0]
		for _, action := range metrics.Actions {
			if action.Id.GetSpellId() != 116 {
				continue
			}
			for _, tm := range action.Targets {
				casts += float64(tm.Casts)
				hits += float64(tm.Hits)
				crits += float64(tm.Crits)
			}
		}
		for _, aura := range metrics.Auras {
			if aura.Id.GetSpellId() == 1248808 {
				procs = aura.ProcsAvg * iterations
			}
		}
		return
	}
	casts, hits, crits, procs := run(revelationEnchant)
	if casts == 0 {
		t.Fatal("Frostbolt rank 1 was never cast")
	}
	t.Logf("chance %.5f: %.0f casts, %.0f non-crit hits, %.0f crits, %.0f gains", p, casts, hits, crits, procs)
	if sd := math.Sqrt(hits * p * (1 - p)); math.Abs(procs-hits*p) > 4*sd {
		t.Errorf("%.0f Revelation gains in %.0f non-crit Frostbolt hits, want about %.0f", procs, hits, hits*p)
	}

	baseCasts, baseHits, baseCrits, baseProcs := run(0)
	if baseProcs != 0 {
		t.Errorf("%.0f Revelation gains without the enchant", baseProcs)
	}
	// Each gain is spent by the next Frostbolt, which crits if it lands (a miss spends it too); the
	// rest crit at the base rate. A gain off a fight's last Frostbolt goes unspent.
	baseRate, landRate := baseCrits/baseCasts, (baseHits+baseCrits)/baseCasts
	want := baseRate*(casts-procs) + landRate*procs
	if sd := math.Sqrt(casts * baseRate * (1 - baseRate)); crits > want+4*sd || crits < want-4*sd-0.1*iterations {
		t.Errorf("%.0f Frostbolt crits in %.0f casts with %.0f gains, want about %.0f", crits, casts, procs, want)
	}
}
