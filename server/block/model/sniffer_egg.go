package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// SnifferEgg is a model used by sniffer eggs.
type SnifferEgg struct{}

// BBox returns the box of the egg. The 2026-09-14 BDS review, batch B section 6, reads one fixed box: full height,
// inset 1/16 on the X axis and 2/16 on the Z axis, unchanged by cracked_state.
func (SnifferEgg) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	return []cube.BBox{cube.Box(1.0/16, 0, 2.0/16, 15.0/16, 1, 14.0/16)}
}

// FaceSolid always returns false.
func (SnifferEgg) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
