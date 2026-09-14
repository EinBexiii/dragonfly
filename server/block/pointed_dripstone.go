package block

import (
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

// PointedDripstone is one segment of a stalactite or stalagmite growing from a Dripstone block.
type PointedDripstone struct {
	transparent

	// Thickness is the segment this block forms in the spike it is part of.
	Thickness DripstoneThickness
	// Hanging specifies if the spike grows down from a ceiling rather than up from a floor.
	Hanging bool
}

// Model ...
func (d PointedDripstone) Model() world.BlockModel {
	return model.PointedDripstone{Thickness: model.DripstoneThickness(d.Thickness.Uint8()), Hanging: d.Hanging}
}

// EncodeBlock ...
func (d PointedDripstone) EncodeBlock() (string, map[string]any) {
	return "minecraft:pointed_dripstone", map[string]any{
		"dripstone_thickness": d.Thickness.String(),
		"hanging":             d.Hanging,
	}
}

// EncodeItem ...
func (d PointedDripstone) EncodeItem() (name string, meta int16) {
	return blockItemName(d), 0
}

// allPointedDripstone returns all pointed dripstone states.
func allPointedDripstone() (dripstone []world.Block) {
	for _, t := range DripstoneThicknesses() {
		for _, hanging := range []bool{false, true} {
			dripstone = append(dripstone, PointedDripstone{Thickness: t, Hanging: hanging})
		}
	}
	return
}
