package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// AmethystCluster is a model used by amethyst buds and clusters. It is a box growing out of the face of the block the
// cluster is attached to.
type AmethystCluster struct {
	// Height is how far the cluster sticks out of the block it grows on.
	Height float64
	// Inset is how far the sides of the cluster are set in from the sides of the block.
	Inset float64
	// Facing is the face of the neighbouring block the cluster grows on, so the cluster itself points the other way.
	Facing cube.Face
}

// BBox returns the box the cluster grows into. The 2026-09-14 BDS review, batch A section 6, derives every size and
// face from the same pair of constructor parameters, so one formula covers all 24 combinations.
func (a AmethystCluster) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	box := full.ExtendTowards(a.Facing, a.Height-1)
	for _, axis := range []cube.Axis{cube.X, cube.Y, cube.Z} {
		if axis != a.Facing.Axis() {
			box = box.Stretch(axis, -a.Inset)
		}
	}
	return []cube.BBox{box}
}

// FaceSolid always returns false.
func (AmethystCluster) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
