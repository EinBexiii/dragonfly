package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// Hopper is a model used by hoppers: the bowl's floor plate and walls, the
// body under the bowl, and the outlet the hopper points at, as the Bedrock
// server has them (collision collector 0x148ed08d0).
type Hopper struct {
	// Facing is the direction of the outlet, down or a horizontal face.
	Facing cube.Face
}

// BBox ...
func (h Hopper) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	const floorBottom, floorTop, wall = 0.625, 0.6875, 0.125
	return []cube.BBox{
		cube.Box(0, floorBottom, 0, 1, floorTop, 1),
		cube.Box(0, floorTop, 0, wall, 1, 1),
		cube.Box(1-wall, floorTop, 0, 1, 1, 1),
		cube.Box(wall, floorTop, 0, 1-wall, 1, wall),
		cube.Box(wall, floorTop, 1-wall, 1-wall, 1, 1),
		cube.Box(0.25, 0.25, 0.25, 0.75, floorBottom, 0.75),
		hopperOutlet(h.Facing),
	}
}

// hopperOutlet is the spout below the body, or the one reaching out of a
// side for a hopper that feeds sideways.
func hopperOutlet(facing cube.Face) cube.BBox {
	const near, far = 0.375, 0.625
	switch facing {
	case cube.FaceNorth:
		return cube.Box(near, 0.25, 0, far, 0.5, 0.25)
	case cube.FaceSouth:
		return cube.Box(near, 0.25, 0.75, far, 0.5, 1)
	case cube.FaceWest:
		return cube.Box(0, 0.25, near, 0.25, 0.5, far)
	case cube.FaceEast:
		return cube.Box(0.75, 0.25, near, 1, 0.5, far)
	}
	return cube.Box(near, 0, near, far, 0.25, far)
}

// FaceSolid only returns true for the top face of the hopper.
func (Hopper) FaceSolid(_ cube.Pos, face cube.Face, _ world.BlockSource) bool {
	return face == cube.FaceUp
}
