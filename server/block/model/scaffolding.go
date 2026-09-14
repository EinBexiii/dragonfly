package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// scaffoldingSupportTolerance is how far under the top of the block an actor's minimum Y may sit and still be held up.
// The 2026-09-14 BDS review, batch B section 3, reads it out of the binary as 0xb58637bd, -0.0000009999999974752427.
const scaffoldingSupportTolerance = 1e-6

// Scaffolding is a model used by scaffolding.
type Scaffolding struct {
	// Stability is how far the block is from the scaffolding column holding it up, from 0 to 7.
	Stability int
}

// BBox returns the plate supported scaffolding holds an actor up on, and nothing for unsupported scaffolding.
func (s Scaffolding) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	if s.Stability == 7 {
		// Batch B section 3 finds an explicit rejection of stability 7, the state a block takes on just before it
		// falls.
		return nil
	}
	// Batch B section 3: the native getter returns the full cube only when the collision context is non-null, a
	// predicate on it that the corpus does not name returns false, and the actor's box reaches down no further than
	// blockY + 1 - 1e-6; in every other case it returns nothing, which is what lets an actor pass through the
	// interior. A model here has no actor to test, so the equivalent of those two conditions is the top plate the
	// tolerance describes: an actor coming down from above is stopped exactly where the cube would stop it, and one
	// beside or below the block passes through as the game allows.
	return []cube.BBox{cube.Box(0, 1-scaffoldingSupportTolerance, 0, 1, 1, 1)}
}

// FaceSolid always returns false.
func (Scaffolding) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
