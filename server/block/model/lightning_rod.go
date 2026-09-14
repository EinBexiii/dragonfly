package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// LightningRod is a model used by lightning rods.
type LightningRod struct {
	// Axis is the axis the rod points along.
	Axis cube.Axis
}

// BBox returns the rod. The 2026-09-14 BDS review, batch B section 2, reads a 4/16 wide column spanning the block
// along the facing axis, the same geometry an end rod has, so the two share it.
func (l LightningRod) BBox(pos cube.Pos, s world.BlockSource) []cube.BBox {
	return EndRod{Axis: l.Axis}.BBox(pos, s)
}

// FaceSolid always returns false.
func (LightningRod) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
