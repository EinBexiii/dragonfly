package block

import (
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

// LightningRod is a block that attracts lightning strikes to the point it is placed at.
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

// BreakInfo ...
func (l LightningRod) BreakInfo() BreakInfo {
	return newBreakInfo(3, func(t item.Tool) bool {
		return t.ToolType() == item.TypePickaxe && t.HarvestLevel() >= item.ToolTierStone.HarvestLevel
	}, pickaxeEffective, oneOf(l)).withBlastResistance(6)
}

// Wax waxes the lightning rod to stop it from oxidising further.
func (l LightningRod) Wax(cube.Pos, mgl64.Vec3) (world.Block, bool) {
	if l.Waxed {
		return l, false
	}
	l.Waxed = true
	return l, true
}

func (l LightningRod) Strip() (world.Block, world.Sound, bool) {
	if l.Waxed {
		l.Waxed = false
		return l, sound.WaxRemoved{}, true
	} else if ot, ok := l.Oxidation.Decrease(); ok {
		l.Oxidation = ot
		return l, sound.CopperScraped{}, true
	}
	return l, nil, false
}

func (l LightningRod) CanOxidate() bool {
	return !l.Waxed
}

func (l LightningRod) OxidationLevel() OxidationType {
	return l.Oxidation
}

func (l LightningRod) WithOxidationLevel(o OxidationType) Oxidisable {
	l.Oxidation = o
	return l
}

func (l LightningRod) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	attemptOxidation(pos, tx, r, l)
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
