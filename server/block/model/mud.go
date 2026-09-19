package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// Mud is the seven-eighths-high collision model shared by mud and soul sand.
// These blocks retain their solid faces despite their lower collision surface.
type Mud struct {
	Solid
}

// BBox returns a box spanning the full footprint with a height of 7/8.
func (Mud) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	return []cube.BBox{cube.Box(0, 0, 0, 1, 0.875, 1)}
}
