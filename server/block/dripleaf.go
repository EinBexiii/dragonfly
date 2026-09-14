package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
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

// UseOnBlock places the leaf of the plant with its tip towards the player. The stem blocks below it are grown, not
// placed.
func (d BigDripleaf) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, d)
	if !used || !faceSolid(tx, pos.Side(cube.FaceDown), cube.FaceUp) {
		return false
	}
	d.Head, d.Tilt, d.Facing = true, NoneDripleafTilt(), user.Rotation().Direction().Opposite()

	place(tx, pos, d, user, ctx)
	return placed(ctx)
}

// EncodeBlock ...
func (d BigDripleaf) EncodeBlock() (string, map[string]any) {
	return "minecraft:big_dripleaf", map[string]any{
		"big_dripleaf_head":            d.Head,
		"big_dripleaf_tilt":            d.Tilt.String(),
		"minecraft:cardinal_direction": d.Facing.String(),
	}
}

// EncodeItem ...
func (d BigDripleaf) EncodeItem() (name string, meta int16) {
	return blockItemName(d), 0
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

// UseOnBlock places both halves of the plant, with its leaves towards the player.
func (d SmallDripleaf) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, d)
	if !used || !supportsVegetation(d, tx.Block(pos.Side(cube.FaceDown))) {
		return false
	}
	upper := SmallDripleaf{UpperPart: true, Facing: user.Rotation().Direction().Opposite()}
	if !replaceableWith(tx, pos.Side(cube.FaceUp), upper) {
		return false
	}
	d.UpperPart, d.Facing = false, upper.Facing

	place(tx, pos, d, user, ctx)
	place(tx, pos.Side(cube.FaceUp), upper, user, ctx)
	return placed(ctx)
}

// EncodeBlock ...
func (d SmallDripleaf) EncodeBlock() (string, map[string]any) {
	return "minecraft:small_dripleaf_block", map[string]any{
		"minecraft:cardinal_direction": d.Facing.String(),
		"upper_block_bit":              d.UpperPart,
	}
}

// EncodeItem ...
func (d SmallDripleaf) EncodeItem() (name string, meta int16) {
	return blockItemName(d), 0
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
