package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// LilyPad is a model for the lily pad block.
type LilyPad struct{}

// BBox returns the pad, 3/32 high and inset a sixteenth, as Bedrock has
// it (BDS constructor 0xc6840da); the client lands on it at that height.
func (LilyPad) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	return []cube.BBox{cube.Box(0.0625, 0, 0.0625, 0.9375, 0.09375, 0.9375)}
}

// FaceSolid ...
func (LilyPad) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
