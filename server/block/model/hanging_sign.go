package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// HangingSign is a model used by hanging signs. A sign that hangs below a block has no collision at all, while a sign
// attached to the side of one collides only with the bar along the top edge of the block.
type HangingSign struct {
	// Facing is the face of the block the sign is attached by.
	Facing cube.Face
	// Hanging specifies if the sign hangs below a block instead of being attached to the side of one.
	Hanging bool
}

// BBox returns the top bar of a sign attached to a wall and nothing at all for a sign hanging from a chain.
func (s HangingSign) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	if s.Hanging {
		// The chain and the sign panel below it are decoration only: 2026-09-14 BDS review, batch A §1, has the native
		// getter return the empty sentinel for every hanging state.
		return nil
	}
	if s.Facing.Axis() == cube.Z {
		return []cube.BBox{cube.Box(0, 14.0/16, 6.0/16, 1, 1, 10.0/16)}
	}
	// Batch A §1: the branch mask 0x33 covers down, up, west and east, so a vertical facing shares the box of the two
	// horizontal ones on the X axis.
	return []cube.BBox{cube.Box(6.0/16, 14.0/16, 0, 10.0/16, 1, 1)}
}

// FaceSolid always returns false.
func (HangingSign) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
