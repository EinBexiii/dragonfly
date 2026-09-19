package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// ChorusFlower is the model of a chorus flower. It is a full cube where the stem below it is not, and it is a
// type of its own rather than a Solid so that ChorusPlant can recognise it as a neighbour to reach out to.
type ChorusFlower struct{}

// BBox returns a full block.
func (ChorusFlower) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	return []cube.BBox{full}
}

// FaceSolid always returns false.
func (ChorusFlower) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
