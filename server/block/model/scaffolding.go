package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// Scaffolding is a model used by scaffolding.
type Scaffolding struct {
	// Stability is how far the block is from the scaffolding column holding it up, from 0 to 7.
	Stability int
}

// BBox returns a full cube for supported scaffolding and nothing for unsupported scaffolding.
func (s Scaffolding) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	if s.Stability == 7 {
		// The 2026-09-14 BDS review, batch B section 3, finds an explicit rejection of stability 7, the state a block
		// takes on just before it falls.
		return nil
	}
	// Batch B section 3: the native getter returns either this cube or nothing, chosen by a predicate on the collision
	// context that the corpus does not name. That predicate is what lets a player drop down through a column, so
	// until it is identified the supported cube is the state a block has for everything else in the world.
	return []cube.BBox{full}
}

// FaceSolid always returns false.
func (Scaffolding) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
