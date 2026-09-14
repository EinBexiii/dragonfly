package block

// GrindstoneAttachment represents a type of attachment for a Grindstone.
type GrindstoneAttachment struct {
	grindstoneAttachment
}

// StandingGrindstoneAttachment is a type of attachment for a standing Grindstone.
func StandingGrindstoneAttachment() GrindstoneAttachment {
	return GrindstoneAttachment{0}
}

// HangingGrindstoneAttachment is a type of attachment for a hanging Grindstone.
func HangingGrindstoneAttachment() GrindstoneAttachment {
	return GrindstoneAttachment{1}
}

// WallGrindstoneAttachment is a type of attachment for a wall Grindstone.
func WallGrindstoneAttachment() GrindstoneAttachment {
	return GrindstoneAttachment{2}
}

// MultipleGrindstoneAttachment is a legacy type of attachment that the game still stores but never places. It
// is given the same model and support as a standing Grindstone, which is what the game draws it as.
func MultipleGrindstoneAttachment() GrindstoneAttachment {
	return GrindstoneAttachment{3}
}

// GrindstoneAttachments returns all possible GrindstoneAttachments.
func GrindstoneAttachments() []GrindstoneAttachment {
	return []GrindstoneAttachment{StandingGrindstoneAttachment(), HangingGrindstoneAttachment(), WallGrindstoneAttachment(), MultipleGrindstoneAttachment()}
}

type grindstoneAttachment uint8

// Uint8 returns the GrindstoneAttachment as a uint8.
func (g grindstoneAttachment) Uint8() uint8 {
	return uint8(g)
}

// String returns the GrindstoneAttachment as a string.
func (g grindstoneAttachment) String() string {
	switch g {
	case 0:
		return "standing"
	case 1:
		return "hanging"
	case 2:
		return "side"
	case 3:
		return "multiple"
	}
	panic("should never happen")
}
