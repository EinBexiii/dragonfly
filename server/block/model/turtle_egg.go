package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// TurtleEgg is a model used by turtle eggs.
type TurtleEgg struct{}

// BBox returns the single box the eggs share. The 2026-09-14 BDS review, batch B section 6, reads one fixed box for
// every turtle_egg_count and cracked_state, whose bounds are these exact binary32 values rather than round fractions.
func (TurtleEgg) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	return []cube.BBox{cube.Box(0.19999998807907104, 0, 0.19999998807907104, 0.800000011920929, 0.45000001788139343, 0.800000011920929)}
}

// FaceSolid always returns false.
func (TurtleEgg) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
