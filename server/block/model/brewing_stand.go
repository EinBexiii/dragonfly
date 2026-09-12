package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// BrewingStand is a model used by brewing stands.
type BrewingStand struct{}

// BBox ...
func (b BrewingStand) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	return []cube.BBox{
		full.ExtendTowards(cube.FaceUp, -0.875),
		// The stem stands on the base and stops an eighth below the top
		// (BDS collector 0x141ef62a0).
		cube.Box(0.4375, 0, 0.4375, 0.5625, 0.875, 0.5625),
	}
}

// FaceSolid ...
func (b BrewingStand) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
