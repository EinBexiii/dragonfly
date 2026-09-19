package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// ChorusRoot is implemented by the one block a chorus plant reaches its bottom bound down to besides another
// chorus plant or flower: end stone, and only directly below the plant.
type ChorusRoot interface {
	ChorusRoot()
}

// chorusInset is how far every unconnected bound of a chorus plant sits inside the block.
const chorusInset = 0.125

// ChorusPlant is the model of a chorus plant stem.
type ChorusPlant struct{}

// BBox returns a single box, not a core with separate arms: Bedrock moves each of the six bounds out to the
// edge of the block on the sides that connect, so connections on different axes widen the whole box between
// them.
func (c ChorusPlant) BBox(pos cube.Pos, s world.BlockSource) []cube.BBox {
	box := cube.Box(chorusInset, chorusInset, chorusInset, 1-chorusInset, 1-chorusInset, 1-chorusInset)
	for _, f := range cube.Faces() {
		if c.connects(pos, f, s) {
			box = box.ExtendTowards(f, chorusInset)
		}
	}
	return []cube.BBox{box}
}

// FaceSolid always returns false.
func (ChorusPlant) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}

// connects checks if the chorus plant reaches the edge of the block on the face passed.
func (ChorusPlant) connects(pos cube.Pos, f cube.Face, s world.BlockSource) bool {
	side := s.Block(pos.Side(f))
	switch side.Model().(type) {
	case ChorusPlant, ChorusFlower:
		return true
	}
	if f != cube.FaceDown {
		return false
	}
	_, root := side.(ChorusRoot)
	return root
}
