package block

import "github.com/df-mc/dragonfly/server/block/model"

// BellAttachment represents the way a Bell is attached to the blocks around it. Its ordinals are Bedrock's, and
// model.BellAttachment holds the one definition of them: the constructors below derive from it, so the conversion
// at the Model() call site cannot drift.
type BellAttachment struct {
	bellAttachment
}

// StandingBellAttachment is a bell standing on the floor.
func StandingBellAttachment() BellAttachment {
	return BellAttachment{bellAttachment(model.BellStanding)}
}

// HangingBellAttachment is a bell hanging from the block above it.
func HangingBellAttachment() BellAttachment {
	return BellAttachment{bellAttachment(model.BellHanging)}
}

// SideBellAttachment is a bell attached to the side of a single block.
func SideBellAttachment() BellAttachment {
	return BellAttachment{bellAttachment(model.BellSide)}
}

// MultipleBellAttachment is a bell suspended between two blocks facing each other.
func MultipleBellAttachment() BellAttachment {
	return BellAttachment{bellAttachment(model.BellMultiple)}
}

// BellAttachments returns all bell attachments.
func BellAttachments() []BellAttachment {
	return []BellAttachment{StandingBellAttachment(), HangingBellAttachment(), SideBellAttachment(), MultipleBellAttachment()}
}

type bellAttachment uint8

// Uint8 returns the bell attachment as a uint8.
func (b bellAttachment) Uint8() uint8 {
	return uint8(b)
}

// String ...
func (b bellAttachment) String() string {
	switch model.BellAttachment(b) {
	case model.BellStanding:
		return "standing"
	case model.BellHanging:
		return "hanging"
	case model.BellSide:
		return "side"
	case model.BellMultiple:
		return "multiple"
	}
	panic("unknown bell attachment")
}
