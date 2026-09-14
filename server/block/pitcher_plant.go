package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// PitcherPlant is the two block tall flower a PitcherCrop grows into. It has no collision: the 2026-09-14 BDS review,
// batch B section 12, finds both halves installing the empty collision getter.
type PitcherPlant struct {
	empty
	transparent

	// UpperPart specifies if this block is the upper half of the plant.
	UpperPart bool
}

// UseOnBlock places both halves of the plant, which is two blocks tall.
func (p PitcherPlant) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, p)
	if !used || !supportsVegetation(p, tx.Block(pos.Side(cube.FaceDown))) {
		return false
	}
	upper := PitcherPlant{UpperPart: true}
	if !replaceableWith(tx, pos.Side(cube.FaceUp), upper) {
		return false
	}
	p.UpperPart = false

	place(tx, pos, p, user, ctx)
	place(tx, pos.Side(cube.FaceUp), upper, user, ctx)
	return placed(ctx)
}

// NeighbourUpdateTick breaks the half of the plant whose other half is gone, and the plant itself when its soil is.
func (p PitcherPlant) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if !p.canSurvive(pos, tx) {
		breakBlockNoDrops(p, pos, tx)
	}
}

// canSurvive checks that the plant still has both of its halves and, under the lower one, its soil.
func (p PitcherPlant) canSurvive(pos cube.Pos, tx *world.Tx) bool {
	other := cube.FaceUp
	if p.UpperPart {
		other = cube.FaceDown
	}
	half, ok := tx.Block(pos.Side(other)).(PitcherPlant)
	if !ok || half.UpperPart == p.UpperPart {
		return false
	}
	if p.UpperPart {
		return true
	}
	return supportsVegetation(p, tx.Block(pos.Side(cube.FaceDown)))
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
