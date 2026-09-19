package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// Cauldron is the model of a cauldron: a floor and four walls, the same
// whatever it holds, as the Bedrock server has it (collision collector
// 0x1469a9630).
type Cauldron struct{}

// BBox ...
func (Cauldron) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	const floorTop, wall = 0.3125, 0.125
	return []cube.BBox{
		cube.Box(0, 0, 0, 1, floorTop, 1),
		cube.Box(0, 0, 0, wall, 1, 1),
		cube.Box(1-wall, 0, 0, 1, 1, 1),
		cube.Box(0, 0, 0, 1, 1, wall),
		cube.Box(0, 0, 1-wall, 1, 1, 1),
	}
}

// FaceSolid ...
func (Cauldron) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
