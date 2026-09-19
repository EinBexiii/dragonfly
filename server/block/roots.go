package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// MangroveRoots is the tangle of roots a mangrove tree grows on. It is a full cube: the 2026-09-14 BDS review, batch B
// section 10, finds it inheriting the unit-box getter despite the gaps in its model.
type MangroveRoots struct {
	solid
	transparent
}

// EncodeBlock ...
func (MangroveRoots) EncodeBlock() (string, map[string]any) {
	return "minecraft:mangrove_roots", nil
}

// EncodeItem ...
func (m MangroveRoots) EncodeItem() (name string, meta int16) {
	return blockItemName(m), 0
}

// HangingRoots are the roots that hang from the ceiling of a lush cave. They have no collision: batch B section 10
// finds the dimensions of their constructor box bypassed by the empty collision getter.
type HangingRoots struct {
	empty
	transparent
}

// UseOnBlock hangs the roots from the block above them.
func (h HangingRoots) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, h)
	if !used || !faceSolid(tx, pos.Side(cube.FaceUp), cube.FaceDown) {
		return false
	}

	place(tx, pos, h, user, ctx)
	return placed(ctx)
}

// NeighbourUpdateTick ...
func (h HangingRoots) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if !faceSolid(tx, pos.Side(cube.FaceUp), cube.FaceDown) {
		breakBlockNoDrops(h, pos, tx)
	}
}

// EncodeBlock ...
func (HangingRoots) EncodeBlock() (string, map[string]any) {
	return "minecraft:hanging_roots", nil
}

// EncodeItem ...
func (h HangingRoots) EncodeItem() (name string, meta int16) {
	return blockItemName(h), 0
}
