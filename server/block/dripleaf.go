package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

// BigDripleaf is a plant with a leaf that tips over when a player stands on it. It is either the leaf itself or one of
// the stem blocks holding it up.
type BigDripleaf struct {
	transparent

	// Head specifies if this block is the leaf on top of the plant rather than one of the stem blocks below it.
	Head bool
	// Tilt is how far the leaf has tipped over.
	Tilt DripleafTilt
	// Facing is the direction the leaf points away from its stem.
	Facing cube.Direction
}

// Model ...
func (d BigDripleaf) Model() world.BlockModel {
	return model.BigDripleaf{Head: d.Head, Tilt: model.DripleafTilt(d.Tilt.Uint8())}
}

// EncodeBlock ...
func (d BigDripleaf) EncodeBlock() (string, map[string]any) {
	return "minecraft:big_dripleaf", map[string]any{
		"big_dripleaf_head":            d.Head,
		"big_dripleaf_tilt":            d.Tilt.String(),
		"minecraft:cardinal_direction": d.Facing.String(),
	}
}

// SmallDripleaf is the two block tall plant that a BigDripleaf grows from.
type SmallDripleaf struct {
	transparent
	empty

	// UpperPart specifies if this block is the upper half of the plant.
	UpperPart bool
	// Facing is the direction the leaves of the plant point in.
	Facing cube.Direction
}

// EncodeBlock ...
func (d SmallDripleaf) EncodeBlock() (string, map[string]any) {
	return "minecraft:small_dripleaf_block", map[string]any{
		"minecraft:cardinal_direction": d.Facing.String(),
		"upper_block_bit":              d.UpperPart,
	}
}

// allDripleaves returns all big and small dripleaf states.
func allDripleaves() (dripleaves []world.Block) {
	for _, d := range cube.Directions() {
		for _, upper := range []bool{false, true} {
			dripleaves = append(dripleaves, SmallDripleaf{UpperPart: upper, Facing: d})
		}
		for _, head := range []bool{false, true} {
			for _, t := range DripleafTilts() {
				dripleaves = append(dripleaves, BigDripleaf{Head: head, Tilt: t, Facing: d})
			}
		}
	}
	return
}
