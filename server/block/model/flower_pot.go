package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// FlowerPot is the model of a flower pot: a centred box 3/8 wide and 3/8
// high, as the Bedrock server has it. It has no solid faces.
type FlowerPot struct{}

// BBox ...
func (FlowerPot) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	const inset, height = 0.3125, 0.375
	return []cube.BBox{cube.Box(inset, 0, inset, 1-inset, height, 1-inset)}
}

// FaceSolid ...
func (FlowerPot) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
