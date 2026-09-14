package block

import (
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

// PitcherPlant is the two block tall flower a PitcherCrop grows into. It has no collision: the 2026-09-14 BDS review,
// batch B section 12, finds both halves installing the empty collision getter.
type PitcherPlant struct {
	empty
	transparent

	// UpperPart specifies if this block is the upper half of the plant.
	UpperPart bool
}

// EncodeBlock ...
func (p PitcherPlant) EncodeBlock() (string, map[string]any) {
	return "minecraft:pitcher_plant", map[string]any{"upper_block_bit": p.UpperPart}
}

// EncodeItem ...
func (p PitcherPlant) EncodeItem() (name string, meta int16) {
	return blockItemName(p), 0
}

// PitcherCrop is the crop a pitcher pod grows into before it becomes a PitcherPlant.
type PitcherCrop struct {
	transparent

	// Growth is the stage of growth of the crop, ranging from 0 to 7.
	Growth int
	// UpperPart specifies if this block is the upper half of the plant.
	UpperPart bool
}

// Model ...
func (p PitcherCrop) Model() world.BlockModel {
	return model.PitcherCrop{UpperPart: p.UpperPart, Sprouted: p.Growth > 0}
}

// EncodeBlock ...
func (p PitcherCrop) EncodeBlock() (string, map[string]any) {
	return "minecraft:pitcher_crop", map[string]any{
		"growth":          int32(p.Growth),
		"upper_block_bit": p.UpperPart,
	}
}

// EncodeItem ...
func (p PitcherCrop) EncodeItem() (name string, meta int16) {
	return blockItemName(p), 0
}

// allPitcherPlants returns all pitcher plant and pitcher crop states.
func allPitcherPlants() (plants []world.Block) {
	for _, upper := range []bool{false, true} {
		plants = append(plants, PitcherPlant{UpperPart: upper})
		for growth := range 8 {
			plants = append(plants, PitcherCrop{Growth: growth, UpperPart: upper})
		}
	}
	return
}
