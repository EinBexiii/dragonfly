package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// LightningRod is a block that attracts lightning strikes to the point it is placed at.
// Oxidation and Waxed are decoded state and nothing else: this block has no BreakInfo, no waxing or scraping and
// no random oxidation tick, because a loaded world must come back exactly as it was written.
type LightningRod struct {
	transparent
	bassDrum

	// Oxidation is the level of oxidation of the lightning rod.
	Oxidation OxidationType
	// Waxed bool is whether the lightning rod has been waxed with honeycomb.
	Waxed bool
	// Facing is the face of the neighbouring block the rod is attached to, so the rod itself points the other way.
	Facing cube.Face
	// Powered specifies if the rod is currently giving off a redstone signal after being struck.
	Powered bool
}

// Model ...
func (l LightningRod) Model() world.BlockModel {
	return model.LightningRod{Axis: l.Facing.Axis()}
}

// UseOnBlock points the rod out of the face of the block it is placed against, so one placed on the ground stands
// upright and one placed against a wall sticks out of it.
func (l LightningRod) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, face, used := firstReplaceable(tx, pos, face, l)
	if !used {
		return false
	}
	l.Facing, l.Powered = face, false

	place(tx, pos, l, user, ctx)
	return placed(ctx)
}

// EncodeItem ...
func (l LightningRod) EncodeItem() (name string, meta int16) {
	return copperBlockName("lightning_rod", l.Oxidation, l.Waxed), 0
}

// EncodeBlock ...
func (l LightningRod) EncodeBlock() (string, map[string]any) {
	return copperBlockName("lightning_rod", l.Oxidation, l.Waxed), map[string]any{
		"facing_direction": int32(l.Facing),
		"powered_bit":      l.Powered,
	}
}

// allLightningRods returns a list of all lightning rod variants.
func allLightningRods() (rods []world.Block) {
	for _, waxed := range []bool{false, true} {
		for _, o := range OxidationTypes() {
			for _, f := range cube.Faces() {
				for _, powered := range []bool{false, true} {
					rods = append(rods, LightningRod{Oxidation: o, Waxed: waxed, Facing: f, Powered: powered})
				}
			}
		}
	}
	return
}
