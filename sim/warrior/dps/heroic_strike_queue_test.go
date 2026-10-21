package dps

import (
	"strings"
	"testing"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

// A queued Heroic Strike replaces one main hand swing and is gone (Ambiwar, MythicSim Discord #bugs,
// 21 October 2026). The strike used to drop its queue only after its damage, so an extra attack its
// own hit procs at once (Windfury Totem here, Sword Specialization alike) went through
// MaybeReplaceMHSwing with the queue still up and became a second Heroic Strike from the one queue,
// paying the rage again. With core.QueuedSwingConsumedBeforeDamage each queue is spent by one strike;
// off, the old order spends some queues twice, which shows the check can see it.
func TestQueuedHeroicStrikeIsSpentOnce(t *testing.T) {
	strikesPerQueue := func(t *testing.T) (queues int, multi int, extraAfterStrike int) {
		t.Helper()
		player := &proto.Player{
			Name:          "Heroic Strike",
			Race:          proto.Race_RaceOrc,
			Class:         proto.Class_ClassWarrior,
			Equipment:     DualWieldGear.GearSet,
			TalentsString: FuryTalents,
			Spec:          DefaultOptions,
			Rotation:      core.GetAplRotation("../../../ui/specs/warrior/dps/apls", "dps_reck").Rotation,
		}
		result := core.RunRaidSim(&proto.RaidSimRequest{
			Raid: core.SinglePlayerRaidProto(player, &proto.PartyBuffs{WindfuryTotem: true}, &proto.RaidBuffs{},
				&proto.Debuffs{}),
			Encounter:  core.MakeSingleTargetEncounter(0),
			SimOptions: &proto.SimOptions{Iterations: 30, RandomSeed: 3, Debug: true},
		})
		if result.Error != nil {
			t.Fatalf("sim failed: %s", result.Error.Message)
		}
		casts, open := 0, false
		lines := strings.Split(result.Logs, "\n")
		for i, line := range lines {
			switch {
			case strings.Contains(line, "Aura gained: {SpellID: 25286, Tag: 1}"):
				casts, open = 0, true
			case strings.Contains(line, "Aura faded: {SpellID: 25286, Tag: 1}"):
				if open {
					queues++
					if casts > 1 {
						multi++
					}
				}
				open = false
			case open && strings.Contains(line, "Casting {SpellID: 25286} "):
				casts++
			case strings.Contains(line, "{SpellID: 25286} Hit for") || strings.Contains(line, "{SpellID: 25286} Crit for"):
				// The strike's hit procs Windfury and the extra attack is the totem's white swing.
				for _, next := range lines[i+1 : min(i+8, len(lines))] {
					if strings.Contains(next, "Casting {OtherID: 3, Tag: 25584}") {
						extraAfterStrike++
						break
					}
				}
			}
		}
		return queues, multi, extraAfterStrike
	}

	defer func(old bool) { core.QueuedSwingConsumedBeforeDamage = old }(core.QueuedSwingConsumedBeforeDamage)

	core.QueuedSwingConsumedBeforeDamage = false
	if _, multi, _ := strikesPerQueue(t); multi == 0 {
		t.Fatal("with the old order no queue was spent twice: the fight never procs an extra attack off Heroic Strike")
	}

	core.QueuedSwingConsumedBeforeDamage = true
	queues, multi, extra := strikesPerQueue(t)
	if queues == 0 {
		t.Fatal("the rotation never queued Heroic Strike")
	}
	if multi > 0 {
		t.Errorf("%d of %d Heroic Strike queues were spent by more than one strike", multi, queues)
	}
	if extra == 0 {
		t.Error("no Heroic Strike hit was followed by Windfury's white extra attack")
	}
}
