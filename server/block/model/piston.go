package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// PistonArm is the model of an extended piston's arm and head.
type PistonArm struct{}

// BBox returns a full block.
// TODO(bds): arm and head boxes from the Bedrock server
func (PistonArm) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	return []cube.BBox{full}
}

// FaceSolid always returns false.
func (PistonArm) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
