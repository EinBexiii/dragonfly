package block

// BellAttachment represents the way a Bell is attached to the blocks around it.
type BellAttachment struct {
	bellAttachment
}

// StandingBellAttachment is a bell standing on the floor.
func StandingBellAttachment() BellAttachment {
	return BellAttachment{0}
}

// HangingBellAttachment is a bell hanging from the block above it.
func HangingBellAttachment() BellAttachment {
	return BellAttachment{1}
}

// SideBellAttachment is a bell attached to the side of a single block.
func SideBellAttachment() BellAttachment {
	return BellAttachment{2}
}

// MultipleBellAttachment is a bell suspended between two blocks facing each other.
func MultipleBellAttachment() BellAttachment {
	return BellAttachment{3}
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
	switch b {
	case 0:
		return "standing"
	case 1:
		return "hanging"
	case 2:
		return "side"
	case 3:
		return "multiple"
	}
	panic("unknown bell attachment")
}
