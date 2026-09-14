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
	d.Hanging = face == cube.FaceDown
	if !d.canSurvive(pos, tx) {
		return false
	}
	d.Thickness = d.thickness(pos, tx)

	place(tx, pos, d, user, ctx)
	return placed(ctx)
}

// NeighbourUpdateTick breaks a segment that has nothing to grow from, and otherwise lets it take the shape its place
// in the spike now gives it. Writing only a changed shape keeps the segments from updating each other forever.
func (d PointedDripstone) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if !d.canSurvive(pos, tx) {
		breakBlockNoDrops(d, pos, tx)
		return
	}
	if thickness := d.thickness(pos, tx); thickness != d.Thickness {
		d.Thickness = thickness
		tx.SetBlock(pos, d, nil)
	}
}

// thickness returns the segment this block forms in its spike: the tip at the growing end, the base at the block it
// grows from, the frustum just above a tip and the middle anywhere between. Two spikes meeting tip to tip merge.
func (d PointedDripstone) thickness(pos cube.Pos, tx *world.Tx) DripstoneThickness {
	beyond, ok := tx.Block(pos.Side(d.growth())).(PointedDripstone)
	if ok && beyond.Hanging != d.Hanging {
		return MergeDripstoneThickness()
	}
	if !ok || beyond.Hanging != d.Hanging {
		return TipDripstoneThickness()
	}
	if beyond.Thickness == TipDripstoneThickness() || beyond.Thickness == MergeDripstoneThickness() {
		return FrustumDripstoneThickness()
	}
	if behind, ok := tx.Block(pos.Side(d.growth().Opposite())).(PointedDripstone); !ok || behind.Hanging != d.Hanging {
		return BaseDripstoneThickness()
	}
	return MiddleDripstoneThickness()
}

// growth returns the face the spike grows towards, away from the block carrying it.
func (d PointedDripstone) growth() cube.Face {
	if d.Hanging {
		return cube.FaceDown
	}
	return cube.FaceUp
}

// canSurvive checks if the spike still grows out of a block or out of the segment before it.
func (d PointedDripstone) canSurvive(pos cube.Pos, tx *world.Tx) bool {
	grownFrom := d.growth().Opposite()
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
