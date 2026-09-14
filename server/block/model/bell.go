package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// BellAttachment is the way a Bell is attached to the blocks around it. Its values are the ones Bedrock uses, so that
// block.BellAttachment converts to it directly.
type BellAttachment uint8

const (
	// BellStanding is a bell standing on the floor.
	BellStanding BellAttachment = iota
	// BellHanging is a bell hanging from the block above it.
	BellHanging
	// BellSide is a bell attached to the side of a single block.
	BellSide
	// BellMultiple is a bell suspended between two blocks facing each other.
	BellMultiple
)

// Bell is a model used by bells. Every attachment has its own box, mirrored along the direction the bell faces.
type Bell struct {
	// Attachment is the way the bell is attached to the blocks around it.
	Attachment BellAttachment
	// Facing is the direction the mouth of the bell points in.
	Facing cube.Direction
}

// BBox returns the single box the bell collides with. The boxes come from the 2026-09-14 BDS review, batch A section
// 3, which reads one combined AABB per attachment and direction rather than separate bell and support pieces.
func (b Bell) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	axis := b.Facing.Face().Axis()
	switch b.Attachment {
	case BellStanding:
		// The two posts stand on either side of the bell, across the axis it faces.
		return []cube.BBox{cube.Box(0, 0, 0, 1, 13.0/16, 1).Stretch(axis, -4.0/16)}
	case BellHanging:
		return []cube.BBox{cube.Box(4.0/16, 4.0/16, 4.0/16, 12.0/16, 1, 12.0/16)}
	case BellMultiple:
		// The bar spans the full block between the two blocks it is strung between, across the facing axis.
		return []cube.BBox{cube.Box(0, 4.0/16, 0, 1, 1, 1).Stretch(axis.RotateLeft(), -4.0/16)}
	}
	// A bell on the side of a block reaches 13/16 out from the wall behind it.
	return []cube.BBox{cube.Box(0, 4.0/16, 0, 1, 15.0/16, 1).Stretch(axis.RotateLeft(), -4.0/16).ExtendTowards(b.Facing.Face(), -3.0/16)}
}

// FaceSolid always returns false.
func (Bell) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
