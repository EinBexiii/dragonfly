package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// Wall is a model used by all wall types.
type Wall struct {
	// NorthConnection is the height of the connection for the north direction.
	NorthConnection float64
	// EastConnection is the height of the connection for the east direction.
	EastConnection float64
	// SouthConnection is the height of the connection for the south direction.
	SouthConnection float64
	// WestConnection is the height of the connection for the west direction.
	WestConnection float64
	// Post is if the wall is the full height of a block or not.
	Post bool
}

// wallCollisionHeight is how tall a wall is to an entity: like a fence, a
// wall cannot be jumped over, so its collision box reaches a block and a
// half whatever its post and connections look like. A player standing on a
// wall stands at that height; a shorter box has the server see them float.
const wallCollisionHeight = 1.5

// BBox returns the wall's collision boxes. The post and connection heights
// describe the shape the client draws; the boxes an entity collides with
// are all wallCollisionHeight tall.
func (w Wall) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	boxes := []cube.BBox{cube.Box(0.25, 0, 0.25, 0.75, wallCollisionHeight, 0.75)}
	if w.NorthConnection > 0 {
		boxes = append(boxes, cube.Box(0.25, 0, 0, 0.75, wallCollisionHeight, 0.25))
	}
	if w.EastConnection > 0 {
		boxes = append(boxes, cube.Box(0.75, 0, 0.25, 1, wallCollisionHeight, 0.75))
	}
	if w.SouthConnection > 0 {
		boxes = append(boxes, cube.Box(0.25, 0, 0.75, 0.75, wallCollisionHeight, 1))
	}
	if w.WestConnection > 0 {
		boxes = append(boxes, cube.Box(0, 0, 0.25, 0.25, wallCollisionHeight, 0.75))
	}
	return boxes
}

// FaceSolid returns true if the face is in the Y axis.
func (w Wall) FaceSolid(_ cube.Pos, face cube.Face, _ world.BlockSource) bool {
	return face.Axis() == cube.Y
}
