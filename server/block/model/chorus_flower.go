package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// ChorusFlower is the model of a chorus flower.
type ChorusFlower struct{}

// BBox returns a full block.
// TODO(bds): real boxes from the Bedrock server
func (ChorusFlower) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	return []cube.BBox{full}
}

// FaceSolid always returns false.
func (ChorusFlower) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
