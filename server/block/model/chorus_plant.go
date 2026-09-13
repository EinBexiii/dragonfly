package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// ChorusPlant is the model of a chorus plant stem.
type ChorusPlant struct{}

// BBox returns a full block.
// TODO(bds): real boxes from the Bedrock server
func (ChorusPlant) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	return []cube.BBox{full}
}

// FaceSolid always returns false.
func (ChorusPlant) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
