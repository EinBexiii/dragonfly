package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
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

// UseOnBlock turns the litter towards the player, or adds a fourth leaf at most to the litter already on the block.
func (l LeafLitter) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	if existing, ok := tx.Block(pos).(LeafLitter); ok {
		if existing.AdditionalCount >= 3 {
			return false
		}
		existing.AdditionalCount++

		place(tx, pos, existing, user, ctx)
		return placed(ctx)
	}
	pos, _, used := firstReplaceable(tx, pos, face, l)
	if !used || !faceSolid(tx, pos.Side(cube.FaceDown), cube.FaceUp) {
		return false
	}
	l.AdditionalCount, l.Facing = 0, user.Rotation().Direction().Opposite()

	place(tx, pos, l, user, ctx)
	return placed(ctx)
}

// NeighbourUpdateTick ...
func (l LeafLitter) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if !faceSolid(tx, pos.Side(cube.FaceDown), cube.FaceUp) {
		breakBlockNoDrops(l, pos, tx)
	}
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
