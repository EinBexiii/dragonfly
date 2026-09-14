package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// LeafLitter is a scattering of fallen leaves on the ground. It has no collision: the 2026-09-14 BDS review, batch B
// section 9, finds the thin box its constructor builds used as an outline only, with the collision getter returning
// the empty sentinel for every amount and direction.
type LeafLitter struct {
	empty
	transparent

	// AdditionalCount is the amount of additional leaves in the block, ranging from 0 to 7.
	AdditionalCount int
	// Facing is the direction the leaves are turned towards.
	Facing cube.Direction
}

// EncodeBlock ...
func (l LeafLitter) EncodeBlock() (string, map[string]any) {
	return "minecraft:leaf_litter", map[string]any{
		"growth":                       int32(l.AdditionalCount),
		"minecraft:cardinal_direction": l.Facing.String(),
	}
}

// EncodeItem ...
func (l LeafLitter) EncodeItem() (name string, meta int16) {
	return blockItemName(l), 0
}

// allLeafLitter returns all leaf litter states.
func allLeafLitter() (litter []world.Block) {
	for count := range 8 {
		for _, d := range cube.Directions() {
			litter = append(litter, LeafLitter{AdditionalCount: count, Facing: d})
		}
	}
	return
}
