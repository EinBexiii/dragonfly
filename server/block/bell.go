package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Bell is a block that rings when it is hit or used.
type Bell struct {
	transparent
	blockEntityData

	// Attachment is the way the bell is attached to the blocks around it.
	Attachment BellAttachment
	// Facing is the direction the mouth of the bell points in.
	Facing cube.Direction
	// Toggled specifies if the bell is currently swinging.
	Toggled bool
}

// Model ...
func (b Bell) Model() world.BlockModel {
	return model.Bell{Attachment: model.BellAttachment(b.Attachment.Uint8()), Facing: b.Facing}
}

// UseOnBlock attaches the bell the way the face clicked allows: one placed on top of a block stands on it, one
// placed under a block hangs from it, and one placed against a wall is bolted to it, spanning both walls if the
// block on the far side carries it too.
func (b Bell) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, face, used := firstReplaceable(tx, pos, face, b)
	if !used {
		return false
	}
	switch face {
	case cube.FaceUp:
		b.Attachment, b.Facing = StandingBellAttachment(), user.Rotation().Direction().Opposite()
	case cube.FaceDown:
		b.Attachment, b.Facing = HangingBellAttachment(), user.Rotation().Direction().Opposite()
	default:
		b.Attachment, b.Facing = SideBellAttachment(), face.Direction()
		if faceSolid(tx, pos.Side(face), face.Opposite()) {
			b.Attachment = MultipleBellAttachment()
		}
	}
	if !b.canSurvive(pos, tx) {
		return false
	}

	place(tx, pos, b, user, ctx)
	return placed(ctx)
}

// NeighbourUpdateTick ...
func (b Bell) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if !b.canSurvive(pos, tx) {
		breakBlockNoDrops(b, pos, tx)
	}
}

// canSurvive checks if the block the bell's attachment needs is still there.
func (b Bell) canSurvive(pos cube.Pos, tx *world.Tx) bool {
	switch b.Attachment {
	case StandingBellAttachment():
		return faceSolid(tx, pos.Side(cube.FaceDown), cube.FaceUp)
	case HangingBellAttachment():
		return faceSolid(tx, pos.Side(cube.FaceUp), cube.FaceDown)
	}
	// A side bell hangs off the wall its mouth points away from, and a multiple bell needs that wall as well as the
	// one behind it.
	wall := b.Facing.Opposite().Face()
	if !faceSolid(tx, pos.Side(wall), wall.Opposite()) {
		return false
	}
	if b.Attachment == MultipleBellAttachment() {
		return faceSolid(tx, pos.Side(b.Facing.Face()), b.Facing.Face().Opposite())
	}
	return true
}

// EncodeBlock ...
func (b Bell) EncodeBlock() (string, map[string]any) {
	return "minecraft:bell", map[string]any{
		"attachment": b.Attachment.String(),
		"direction":  int32(horizontalDirection(b.Facing)),
		"toggle_bit": b.Toggled,
	}
}

// EncodeItem ...
func (b Bell) EncodeItem() (name string, meta int16) {
	return blockItemName(b), 0
}

// EncodeNBT ...
func (b Bell) EncodeNBT() map[string]any {
	return b.storedNBT("Bell")
}

// DecodeNBT ...
func (b Bell) DecodeNBT(data map[string]any) any {
	b.blockEntityData = b.storeNBT(data)
	return b
}

// allBells returns all bell states.
func allBells() (bells []world.Block) {
	for _, a := range BellAttachments() {
		for _, d := range cube.Directions() {
			for _, toggled := range []bool{false, true} {
				bells = append(bells, Bell{Attachment: a, Facing: d, Toggled: toggled})
			}
		}
	}
	return
}
