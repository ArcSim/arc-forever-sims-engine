package core

import (
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

func init() {
	RegisterAgentFactory(
		proto.Player_DpsWarrior{},
		proto.Spec_SpecDpsWarrior,
		NewFakeRageWarrior,
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_DpsWarrior)
			if !ok {
				panic("Invalid spec value for Dps Warrior!")
			}
			player.Spec = playerSpec
		},
	)
}

const (
	fakeMHSwingSpeed = 2.6
	fakeOHSwingSpeed = 1.8
)

type FakeRageWarrior struct {
	Character
}

func (fw *FakeRageWarrior) GetCharacter() *Character { return &fw.Character }

func (fw *FakeRageWarrior) Initialize() {
	fw.registerFakeThreatSpell()
	fw.registerFakeCategorySpells()
}
func (fw *FakeRageWarrior) ApplyTalents()                  {}
func (fw *FakeRageWarrior) Reset(_ *Simulation)            {}
func (fw *FakeRageWarrior) OnGCDReady(_ *Simulation)       {}
func (fw *FakeRageWarrior) OnEncounterStart(_ *Simulation) {}

func NewFakeRageWarrior(char *Character, _ *proto.Player, _ *proto.Raid) Agent {
	fw := &FakeRageWarrior{
		Character: *char,
	}

	fw.EnableRageBar(RageBarOptions{
		MaxRage:            100,
		BaseRageMultiplier: 1,
		StartingRage:       0,
	})

	fw.EnableAutoAttacks(fw, AutoAttackOptions{
		MainHand:       Weapon{SwingSpeed: fakeMHSwingSpeed},
		OffHand:        Weapon{SwingSpeed: fakeOHSwingSpeed},
		AutoSwingMelee: true,
	})

	return fw
}

func SetupFakeRageSim() *Simulation {
	sim := NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{
			RandomSeed: 100,
		},
		Raid: &proto.Raid{
			Parties: []*proto.Party{
				{
					Players: []*proto.Player{
						{
							Name:      "Warrior",
							Class:     proto.Class_ClassWarrior,
							Buffs:     &proto.IndividualBuffs{},
							Spec:      &proto.Player_DpsWarrior{},
							Equipment: &proto.EquipmentSpec{},
							// AddRage() pokes the rotation, so it needs to be non-nil.
							Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
						},
					},
					Buffs: &proto.PartyBuffs{},
				},
			},
		},
		Encounter: &proto.Encounter{
			Targets: []*proto.Target{
				{Name: "target", Level: 63, MobType: proto.MobType_MobTypeHumanoid},
			},
			Duration: 180,
		},
	}, simsignals.CreateSignals())
	sim.Reset()

	return sim
}

// Feeds a hand-crafted auto attack result to the Rage bar and returns the Rage gained from it.
// PostArmorAndResistanceMultiplier is always populated, because the sim fills it in before the
// outcome is applied, even for swings which end up doing no damage.
func rageFromAutoAttack(sim *Simulation, fw *FakeRageWarrior, spell *Spell, outcome HitOutcome, preOutcomeDamage float64) float64 {
	result := &SpellResult{
		Target:                           sim.Encounter.ActiveTargetUnits[0],
		Outcome:                          outcome,
		PostArmorAndResistanceMultiplier: preOutcomeDamage,
	}
	if outcome.Matches(OutcomeLanded) {
		result.Damage = preOutcomeDamage
	}

	rageBefore := fw.CurrentRage()
	fw.Unit.OnSpellHitDealt(sim, spell, result)

	return fw.CurrentRage() - rageBefore
}

func TestAutoAttackRageGeneration(t *testing.T) {
	// A one-hand MH swing at 2.6 speed: 2.6 * 3.46 = 8.996, whatever it dealt; a crit pays it twice.
	const swingDamage = 500.0

	tests := []struct {
		name     string
		outcome  HitOutcome
		wantRage float64
	}{
		{
			name:     "hit",
			outcome:  OutcomeHit,
			wantRage: 8.996,
		},
		{
			name:     "crit",
			outcome:  OutcomeCrit,
			wantRage: 8.996 * 2,
		},
		{
			name:     "glance",
			outcome:  OutcomeGlance,
			wantRage: 8.996,
		},
		{
			name:     "dodge",
			outcome:  OutcomeDodge,
			wantRage: 0,
		},
		{
			name:     "parry",
			outcome:  OutcomeParry,
			wantRage: 0,
		},
		{
			name:     "miss",
			outcome:  OutcomeMiss,
			wantRage: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sim := SetupFakeRageSim()
			fw := sim.Raid.Parties[0].Players[0].(*FakeRageWarrior)

			gained := rageFromAutoAttack(sim, fw, fw.AutoAttacks.MHAuto(), test.outcome, swingDamage)

			if !WithinToleranceFloat64(test.wantRage, gained, 0.01) {
				t.Fatalf("Incorrect Rage generated on %s: Expected: %0.3f, Actual: %0.3f", test.name, test.wantRage, gained)
			}
		})
	}
}

func TestTwoHandAutoAttackRageGeneration(t *testing.T) {
	// A two-hand MH swing at 2.6 speed: 2.6 * 4.5 = 11.7.
	const swingDamage = 500.0

	sim := SetupFakeRageSim()
	fw := sim.Raid.Parties[0].Players[0].(*FakeRageWarrior)
	fw.AutoAttacks.MH().NormalizedSwingSpeed = TwoHandNormalizedSwingSpeed

	hitRage := rageFromAutoAttack(sim, fw, fw.AutoAttacks.MHAuto(), OutcomeHit, swingDamage)
	if !WithinToleranceFloat64(11.7, hitRage, 0.001) {
		t.Fatalf("Incorrect Rage generated on 2H hit: Expected: %0.3f, Actual: %0.3f", 11.695, hitRage)
	}

	// A crit adds the swing's rage before the two-hand factor: 11.7 + 2.6 * 3.46 = 20.696 (Magey's
	// beta logs: 14.85 + 11.42 on a 3.3 sec two-hander).
	critRage := rageFromAutoAttack(sim, fw, fw.AutoAttacks.MHAuto(), OutcomeCrit, swingDamage)
	if !WithinToleranceFloat64(20.696, critRage, 0.001) {
		t.Fatalf("Incorrect Rage generated on 2H crit: Expected: %0.3f, Actual: %0.3f", 20.696, critRage)
	}
}

func TestOffHandAutoAttackRageGeneration(t *testing.T) {
	// An OH swing at 1.8 speed: 1.8 * 3.46 * 0.5 = 3.114.
	const swingDamage = 500.0

	sim := SetupFakeRageSim()
	fw := sim.Raid.Parties[0].Players[0].(*FakeRageWarrior)
	ohAuto := fw.AutoAttacks.OHAuto()

	hitRage := rageFromAutoAttack(sim, fw, ohAuto, OutcomeHit, swingDamage)
	if !WithinToleranceFloat64(3.114, hitRage, 0.01) {
		t.Fatalf("Incorrect Rage generated on OH hit: Expected: %0.3f, Actual: %0.3f", 3.114, hitRage)
	}

	dodgeRage := rageFromAutoAttack(sim, fw, ohAuto, OutcomeDodge, swingDamage)
	if dodgeRage != 0 {
		t.Fatalf("Dodged OH swing generated %0.3f Rage, expected none", dodgeRage)
	}
}

func TestOffHandRageMultiplier(t *testing.T) {
	const swingDamage = 500.0

	sim := SetupFakeRageSim()
	fw := sim.Raid.Parties[0].Players[0].(*FakeRageWarrior)

	unmodified := rageFromAutoAttack(sim, fw, fw.AutoAttacks.OHAuto(), OutcomeHit, swingDamage)
	if !WithinToleranceFloat64(3.114, unmodified, 0.01) {
		t.Fatalf("Incorrect Rage generated on OH hit: Expected: %0.3f, Actual: %0.3f", 3.114, unmodified)
	}

	fw.SetOffHandRageMultiplier(1.5)

	modified := rageFromAutoAttack(sim, fw, fw.AutoAttacks.OHAuto(), OutcomeHit, swingDamage)
	if !WithinToleranceFloat64(4.671, modified, 0.01) {
		t.Fatalf("Incorrect Rage generated on multiplied OH hit: Expected: %0.3f, Actual: %0.3f", 4.671, modified)
	}

	mhRage := rageFromAutoAttack(sim, fw, fw.AutoAttacks.MHAuto(), OutcomeHit, swingDamage)
	if !WithinToleranceFloat64(8.996, mhRage, 0.01) {
		t.Fatalf("Incorrect Rage generated on MH hit: Expected: %0.3f, Actual: %0.3f", 8.996, mhRage)
	}
}
