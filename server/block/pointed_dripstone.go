package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
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

// UseOnBlock hangs the spike from the block above when it is placed under one, and stands it on the block below
// otherwise. A single segment is the tip of its spike.
func (d PointedDripstone) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, face, used := firstReplaceable(tx, pos, face, d)
	if !used {
		return false
	}
	d.Thickness, d.Hanging = TipDripstoneThickness(), face == cube.FaceDown
	if !d.canSurvive(pos, tx) {
		return false
	}

	place(tx, pos, d, user, ctx)
	return placed(ctx)
}

// NeighbourUpdateTick ...
func (d PointedDripstone) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if !d.canSurvive(pos, tx) {
		breakBlockNoDrops(d, pos, tx)
	}
}

// canSurvive checks if the spike still grows out of a block or out of the segment before it.
func (d PointedDripstone) canSurvive(pos cube.Pos, tx *world.Tx) bool {
	grownFrom := cube.FaceDown
	if d.Hanging {
		grownFrom = cube.FaceUp
	}
	side := pos.Side(grownFrom)
	if other, ok := tx.Block(side).(PointedDripstone); ok {
		return other.Hanging == d.Hanging
	}
	return faceSolid(tx, side, grownFrom.Opposite())
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
