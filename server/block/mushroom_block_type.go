package block

// MushroomBlockType represents a type of huge mushroom block, being the brown cap, the red cap or the stem.
type MushroomBlockType struct {
	mushroomBlock
}

// BrownMushroomBlock returns the brown mushroom cap type.
func BrownMushroomBlock() MushroomBlockType {
	return MushroomBlockType{0}
}

// RedMushroomBlock returns the red mushroom cap type.
func RedMushroomBlock() MushroomBlockType {
	return MushroomBlockType{1}
}

// MushroomStem returns the mushroom stem type.
func MushroomStem() MushroomBlockType {
	return MushroomBlockType{2}
}

// MushroomBlockTypes returns all huge mushroom block types.
func MushroomBlockTypes() []MushroomBlockType {
	return []MushroomBlockType{BrownMushroomBlock(), RedMushroomBlock(), MushroomStem()}
}

type mushroomBlock uint8

// Uint8 returns the huge mushroom block type as a uint8.
func (m mushroomBlock) Uint8() uint8 {
	return uint8(m)
}

// String returns the huge mushroom block type as a string.
func (m mushroomBlock) String() string {
	switch m {
	case 0:
		return "brown_mushroom_block"
	case 1:
		return "red_mushroom_block"
	case 2:
		return "mushroom_stem"
	}
	panic("should never happen")
}
