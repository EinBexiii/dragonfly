package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// DaylightDetector is the model of a daylight detector: a block six pixels high, lower than a slab.
type DaylightDetector struct{}

// BBox returns a BBox 0.375 of a block high.
func (DaylightDetector) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	return []cube.BBox{cube.Box(0, 0, 0, 1, 0.375, 1)}
}

// FaceSolid always returns false.
func (DaylightDetector) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
