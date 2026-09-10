package entity

import (
	"testing"

	"github.com/df-mc/dragonfly/server/block/cube"
)

// An entity settling onto a fence or a wall from just above it has its box
// entirely in the block above the barrier; the search must still reach the
// barrier's block, whose box stands half a block into it.
func TestCollisionSearchReachesABarrierBelow(t *testing.T) {
	box := cube.Box(0.2, 1.5, 0.2, 0.8, 2, 0.8)
	low, _ := searchRange(box)
	if low[1] > 0 {
		t.Fatalf("search starts at y=%d, the barrier is at y=0", low[1])
	}
}
