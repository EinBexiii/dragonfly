package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
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

// EncodeBlock ...
func (b Bell) EncodeBlock() (string, map[string]any) {
	return "minecraft:bell", map[string]any{
		"attachment": b.Attachment.String(),
		"direction":  int32(horizontalDirection(b.Facing)),
		"toggle_bit": b.Toggled,
	}
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
