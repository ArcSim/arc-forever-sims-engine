package main

import (
	"slices"
	"testing"

	"github.com/wowsims/forever/sim/core/spelldata"
)

func TestEffectLines(t *testing.T) {
	cases := []struct {
		name    string
		id      int32
		effect  int
		human   string
		literal string
	}{
		{
			name:    "a bleed tick",
			id:      11574,
			effect:  1,
			human:   "21 physical damage every 3 s to the enemy (7 ticks)",
			literal: "E_APPLY_AURA A_PERIODIC_DAMAGE base=21 period=3000ms target=[6,0]",
		},
		{
			name:    "school damage that rolls and scales",
			id:      116,
			effect:  2,
			human:   "18–20 frost damage to the enemy (at level 60, +0.407 spell power)",
			literal: "E_SCHOOL_DAMAGE base=19 ppl=0.2 variance=0.10526316 sp=0.407 target=[6,0]",
		},
		{
			name:    "normalised weapon damage over an area",
			id:      1680,
			effect:  1,
			human:   "normalised weapon damage to every enemy around the caster",
			literal: "E_NORMALIZED_WEAPON_DMG base=0 sp=1 target=[22,15]",
		},
		{
			name:    "normalised weapon damage with a flat bonus",
			id:      19434,
			effect:  1,
			human:   "normalised weapon damage +20 to the enemy",
			literal: "E_NORMALIZED_WEAPON_DMG base=20 sp=1 target=[6,0]",
		},
		{
			name:    "a proc that names its triggered spell",
			id:      21838,
			effect:  1,
			human:   "casts 29478 Battlegear of Might on its trigger",
			literal: "E_APPLY_AURA A_PROC_TRIGGER_SPELL base=0 trigger=29478 target=[1,0]",
		},
		{
			name:   "a percentage modifier naming its spells by class mask",
			id:     16493,
			effect: 1,
			human:  "SPELLMOD_CRIT_DAMAGE_BONUS +10% on family 4 mask 0xee604cee,0x8140,0x1,0x10000400",
			literal: "E_APPLY_AURA A_ADD_PCT_MODIFIER base=10 sp=1 misc=15 " +
				"family=4 mask=0xee604cee,0x8140,0x1,0x10000400 target=[1,0]",
		},
		{
			name:    "a flat cost modifier, which the client states on a 0-1000 bar",
			id:      12282,
			effect:  1,
			human:   "SPELLMOD_COST -10 (-1 on a rage bar) on family 4 mask 0x40",
			literal: "E_APPLY_AURA A_ADD_FLAT_MODIFIER base=-10 misc=14 family=4 mask=0x40 target=[1,0]",
		},
		{
			name:    "a talent the client states nothing about",
			id:      12295,
			effect:  1,
			human:   "a dummy aura holding 15",
			literal: "E_APPLY_AURA A_DUMMY base=15 target=[1,0]",
		},
		{
			name:    "a dummy tick, which states no amount at all",
			id:      10,
			effect:  3,
			human:   "a dummy aura every 1 s (8 ticks)",
			literal: "E_APPLY_AURA A_PERIODIC_DUMMY base=0 period=1000ms target=[1,0]",
		},
		{
			name:    "a heal that rolls",
			id:      20007,
			effect:  2,
			human:   "75–125 healing to the caster",
			literal: "E_HEAL base=100 variance=0.5 target=[1,0]",
		},
		{
			name:    "a heal that rolls and scales",
			id:      2050,
			effect:  1,
			human:   "46–56 healing to the friendly target (at level 60, +0.429 spell power)",
			literal: "E_HEAL base=51 ppl=0.3 variance=0.19607843 sp=0.429 amplitude=1 target=[21,0]",
		},
		{
			name:    "rage off the client's 0-1000 bar",
			id:      100,
			effect:  2,
			human:   "restores 9 rage to the caster",
			literal: "E_ENERGIZE base=90 misc=1 target=[1,0]",
		},
		{
			name:    "a buff the caster hands to its raid (client 70425; party before)",
			id:      25289,
			effect:  1,
			human:   "+139 attack power to the raid around the caster (at level 60)",
			literal: "E_APPLY_AURA A_MOD_ATTACK_POWER base=139 ppl=0.6 radius=20 target=[56,0]",
		},
		{
			name:    "a shape no branch words",
			id:      348,
			effect:  3,
			human:   unrecognised,
			literal: "E_SCRIPT_EFFECT base=100 target=[6,0]",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			lines := effectLines(spelldata.MustFind(c.id))
			if len(lines) < c.effect {
				t.Fatalf("spell %d has %d effects, want at least %d", c.id, len(lines), c.effect)
			}
			line := lines[c.effect-1]
			if line.Human != c.human {
				t.Errorf("human = %q, want %q", line.Human, c.human)
			}
			if line.Literal != c.literal {
				t.Errorf("literal = %q, want %q", line.Literal, c.literal)
			}
		})
	}
}

func TestHeader(t *testing.T) {
	cases := []struct {
		name string
		id   int32
		want []field
	}{
		{
			name: "a bleed with a stance and a weapon requirement",
			id:   11574,
			want: []field{
				{"school", "physical"},
				{"defense", "melee"},
				{"mechanic", "bleed"},
				{"duration", "21 s"},
				{"gcd", "1.5 s"},
				{"cost", "10 rage"},
				{"range", "5 yd"},
				{"stance", "Battle, Defensive"},
				{"equip", "needs a melee weapon"},
				{"attrs", "refund on miss, periodic can crit"},
				{"labels", "25"},
			},
		},
		{
			name: "an area ability with a target cap",
			id:   1680,
			want: []field{
				{"school", "physical"},
				{"defense", "melee"},
				{"cooldown", "10 s (category)"},
				{"gcd", "1.5 s"},
				{"cost", "25 rage"},
				{"stance", "Berserker"},
				{"equip", "needs a melee weapon"},
				{"targets", "up to 4"},
				{"labels", "25"},
			},
		},
		{
			name: "charges and a shield",
			id:   2565,
			want: []field{
				{"school", "physical"},
				{"duration", "7 s"},
				{"cooldown", "5 s"},
				{"cost", "10 rage"},
				{"stance", "Defensive"},
				{"equip", "needs a shield"},
				{"charges", "2"},
				{"labels", "25"},
			},
		},
		{
			name: "an enchant naming the slots it takes",
			id:   410021,
			want: []field{
				{"school", "physical"},
				{"cast", "3 s"},
				{"equip", "needs cloth, leather, mail or plate in the chest or robe"},
				{"labels", "3071"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := headerFields(spelldata.MustFind(c.id))
			if !slices.Equal(got, c.want) {
				t.Errorf("header =\n%q\nwant\n%q", got, c.want)
			}
		})
	}

	if got := title(spelldata.MustFind(11574)); got != "11574 Rend (Rank 7)" {
		t.Errorf("title = %q", got)
	}
}

func TestProcSummary(t *testing.T) {
	cases := []struct {
		name string
		id   int32
		want string
	}{
		{
			name: "a chance the column really states",
			id:   21838,
			want: "hears PROC_FLAG_TAKE_ANY_DAMAGE; 20% chance",
		},
		{
			name: "a 100 the tooltip states no chance beside",
			id:   12834,
			want: "hears PROC_FLAG_DEAL_MELEE_SWING, PROC_FLAG_DEAL_MELEE_ABILITY, " +
				"PROC_FLAG_DEAL_RANGED_ATTACK, PROC_FLAG_DEAL_RANGED_ABILITY, " +
				"PROC_FLAG_DEAL_HARMFUL_ABILITY, PROC_FLAG_DEAL_HARMFUL_SPELL; " +
				"no roll, it fires whenever its condition is met; tooltip says on a crit",
		},
		{
			name: "a chance an effect states, not the column",
			id:   11103,
			want: "hears PROC_FLAG_DEAL_MELEE_ABILITY, PROC_FLAG_DEAL_RANGED_ATTACK, " +
				"PROC_FLAG_DEAL_RANGED_ABILITY, PROC_FLAG_DEAL_HELPFUL_ABILITY, " +
				"PROC_FLAG_DEAL_HARMFUL_ABILITY, PROC_FLAG_DEAL_HELPFUL_SPELL, " +
				"PROC_FLAG_DEAL_HARMFUL_SPELL; 10% chance, from effect 1; " +
				"tooltip says one named ability",
		},
		{
			name: "the internal cooldown is a header line, not part of the proc",
			id:   324,
			want: "hears PROC_FLAG_TAKE_MELEE_SWING, PROC_FLAG_TAKE_MELEE_ABILITY, " +
				"PROC_FLAG_TAKE_RANGED_ATTACK, PROC_FLAG_TAKE_RANGED_ABILITY, " +
				"PROC_FLAG_TAKE_HARMFUL_ABILITY, PROC_FLAG_TAKE_HARMFUL_SPELL; " +
				"no roll, it fires whenever its condition is met; tooltip says heals count",
		},
		{
			name: "the heartbeat bit alone is not a proc",
			id:   21849,
			want: "",
		},
		{
			name: "no proc columns at all",
			id:   11574,
			want: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := procSummary(spelldata.MustFind(c.id)); got != c.want {
				t.Errorf("procSummary = %q, want %q", got, c.want)
			}
		})
	}

	if got := refList(spelldata.MustFind(12834)); !slices.Equal(got, []string{"412609 Deep Wound"}) {
		t.Errorf("refList = %q", got)
	}
	if got := refList(spelldata.MustFind(11574)); len(got) != 0 {
		t.Errorf("refList = %q, want none", got)
	}
}

func TestLadderRefs(t *testing.T) {
	cases := []struct {
		name string
		id   int32
		want []string
	}{
		{
			name: "the top rank of a Ranked family",
			id:   11574,
			want: []string{"warrior spellData.Rend.Highest()"},
		},
		{
			name: "a rank below the top",
			id:   11572,
			want: []string{"warrior spellData.Rend.ByID(11572)"},
		},
		{
			name: "a trait-tree talent, whose ranks are one spell",
			id:   12282,
			want: []string{"warrior spellData.ImprovedHeroicStrike.Rank(n), n up to 3"},
		},
		{
			name: "an id two families carry",
			id:   12295,
			want: []string{
				"warrior spellData.ImprovedTacticalMastery.Rank(n), n up to 5",
				"warrior spellData.TacticalMasteryTriggered.Highest()",
			},
		},
		{
			name: "an item set bonus, which no class file names",
			id:   21838,
			want: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := ladderRefs(c.id); !slices.Equal(got, c.want) {
				t.Errorf("ladderRefs(%d) = %q, want %q", c.id, got, c.want)
			}
		})
	}
}

func TestHighestRank(t *testing.T) {
	if s := highestRank(spelldata.ByName("Rend")); s.ID != 11574 {
		t.Errorf("Rend picks %d %s, want 11574 Rank 7", s.ID, s.Rank)
	}
	if s := highestRank(spelldata.ByName("Whirlwind")); s.ID != 1680 {
		t.Errorf("Whirlwind picks %d, want the warrior ability 1680", s.ID)
	}
}

func TestParseArgs(t *testing.T) {
	for _, args := range [][]string{{"11574", "-json"}, {"-json", "11574"}} {
		opts, err := parseArgs(args)
		if err != nil {
			t.Fatalf("parseArgs(%q): %v", args, err)
		}
		if opts.mode != "spell" || opts.arg != "11574" || !opts.json {
			t.Errorf("parseArgs(%q) = %+v", args, opts)
		}
	}
	if _, err := parseArgs(nil); err == nil {
		t.Error("parseArgs with no arguments wants a usage error")
	}
}
