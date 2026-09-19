package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// SnowLayer is the ordinary collision model of layered snow.
type SnowLayer struct {
	// Layers is the number of layers, from 1 to 8.
	Layers int
}

// BBox returns a box one layer shorter than the visible snow, or no box for one layer.
func (s SnowLayer) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	if s.Layers <= 1 {
		return nil
	}
	return []cube.BBox{cube.Box(0, 0, 0, 1, float64(s.Layers-1)/8, 1)}
}

// FaceSolid returns true for the bottom face when the snow has collision.
func (s SnowLayer) FaceSolid(_ cube.Pos, face cube.Face, _ world.BlockSource) bool {
	return s.Layers > 1 && face == cube.FaceDown
}
