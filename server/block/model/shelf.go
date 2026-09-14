package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// Shelf is a model used by shelves. A shelf is a 5/16 deep slab standing against the wall it faces away from.
type Shelf struct {
	// Facing is the direction the front of the shelf points in.
	Facing cube.Direction
}

// BBox returns a box 5/16 deep against the block face opposite to the direction the shelf faces.
func (s Shelf) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	// 2026-09-14 BDS review, batch A section 2: the native getter starts from a unit box and moves one endpoint by
	// 0.6875 towards the facing direction, leaving 5/16 of depth at the back. powered_bit is not read.
	return []cube.BBox{full.ExtendTowards(s.Facing.Face(), -11.0/16)}
}

// FaceSolid returns true for the face the shelf stands against.
func (s Shelf) FaceSolid(_ cube.Pos, face cube.Face, _ world.BlockSource) bool {
	return face == s.Facing.Face().Opposite()
}
