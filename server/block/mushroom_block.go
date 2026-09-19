package block

import (
	"github.com/df-mc/dragonfly/server/world"
)

// itemMushroomBits is the huge_mushroom_bits value Bedrock gives the item form of a huge mushroom block: every
// face textured as a cap or stem.
const itemMushroomBits = 14

// MushroomBlock is the large cap or stem block that huge mushrooms are built out of.
type MushroomBlock struct {
	solid

	// Type is the type of huge mushroom block: a brown cap, a red cap or a stem.
	Type MushroomBlockType
	// Bits selects which of the block's faces carry the cap or stem texture, 0-15.
	Bits int
}

// BreakInfo ...
func (m MushroomBlock) BreakInfo() BreakInfo {
	return newBreakInfo(0.2, alwaysHarvestable, axeEffective, oneOf(m))
}

// EncodeItem ...
func (m MushroomBlock) EncodeItem() (name string, meta int16) {
	return "minecraft:" + m.Type.String(), 0
}

// EncodeBlock ...
func (m MushroomBlock) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:" + m.Type.String(), map[string]any{"huge_mushroom_bits": int32(m.Bits)}
}

// allMushroomBlocks ...
func allMushroomBlocks() (blocks []world.Block) {
	for _, t := range MushroomBlockTypes() {
		for bits := 0; bits < 16; bits++ {
			blocks = append(blocks, MushroomBlock{Type: t, Bits: bits})
		}
	}
	return
}
