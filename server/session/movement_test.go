package session

import (
	"testing"

	"github.com/google/uuid"
)

func TestMovementBroadcastValidate(t *testing.T) {
	cases := []struct {
		name string
		c    MovementBroadcastConfig
		ok   bool
	}{
		{"defaults", MovementBroadcastConfig{}, true},
		{"distances must grow", MovementBroadcastConfig{NearDistance: 32, MidDistance: 16, FarDistance: 64}, false},
		{"intervals must grow", MovementBroadcastConfig{MidInterval: 8, FarInterval: 4, DistantInterval: 10}, false},
		{"interval past the wheel", MovementBroadcastConfig{DistantInterval: movementWheel}, false},
		{"negative combat", MovementBroadcastConfig{CombatTicks: -1}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.c.Validate(); (err == nil) != c.ok {
				t.Fatalf("Validate() = %v, want ok=%v", err, c.ok)
			}
		})
	}
	n := MovementBroadcastConfig{}.normalized()
	if n.NearDistance != 16 || n.MidDistance != 32 || n.FarDistance != 64 || n.DistantInterval != 10 || n.CombatTicks != 40 {
		t.Fatalf("defaults = %+v", n)
	}
}

// The phase spreads a crowd across the ticks of its interval rather than
// letting them all come due together.
func TestMovementPhaseSpreads(t *testing.T) {
	viewer := uuid.New()
	const interval = 10
	counts := map[uint8]int{}
	for i := 0; i < 2000; i++ {
		p := movementPhase(uuid.New(), viewer, interval)
		if p < 1 || p > interval {
			t.Fatalf("phase %d outside 1..%d", p, interval)
		}
		counts[p]++
	}
	for slot := uint8(1); slot <= interval; slot++ {
		if counts[slot] < 100 {
			t.Fatalf("slot %d took %d of 2000, want the crowd spread", slot, counts[slot])
		}
	}
	// The same pair always lands on the same tick.
	e := uuid.New()
	if movementPhase(e, viewer, interval) != movementPhase(e, viewer, interval) {
		t.Fatal("the phase of one pair moved between calls")
	}
}
