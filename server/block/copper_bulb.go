package block

import (
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

// CopperBulb is a light source that a redstone pulse switches on and off. It is a full cube in every state: the
// 2026-09-14 BDS review, batch B section 1, finds all eight variants inheriting the unit-box getter, which reads
// neither lit nor powered_bit.
type CopperBulb struct {
	solid
	bassDrum

	// Oxidation is the level of oxidation of the copper bulb.
	Oxidation OxidationType
	// Waxed bool is whether the copper bulb has been waxed with honeycomb.
	Waxed bool
	// Lit specifies if the bulb is currently giving off light.
	Lit bool
	// Powered specifies if the bulb is currently receiving a redstone signal.
	Powered bool
}

// BreakInfo ...
func (c CopperBulb) BreakInfo() BreakInfo {
	return newBreakInfo(3, func(t item.Tool) bool {
		return t.ToolType() == item.TypePickaxe && t.HarvestLevel() >= item.ToolTierStone.HarvestLevel
	}, pickaxeEffective, oneOf(c)).withBlastResistance(6)
}

// Wax waxes the copper bulb to stop it from oxidising further.
func (c CopperBulb) Wax(cube.Pos, mgl64.Vec3) (world.Block, bool) {
	if c.Waxed {
		return c, false
	}
	c.Waxed = true
	return c, true
}

func (c CopperBulb) Strip() (world.Block, world.Sound, bool) {
	if c.Waxed {
		c.Waxed = false
		return c, sound.WaxRemoved{}, true
	} else if ot, ok := c.Oxidation.Decrease(); ok {
		c.Oxidation = ot
		return c, sound.CopperScraped{}, true
	}
	return c, nil, false
}

func (c CopperBulb) CanOxidate() bool {
	return !c.Waxed
}

func (c CopperBulb) OxidationLevel() OxidationType {
	return c.Oxidation
}

func (c CopperBulb) WithOxidationLevel(o OxidationType) Oxidisable {
	c.Oxidation = o
	return c
}

func (c CopperBulb) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	attemptOxidation(pos, tx, r, c)
}

// EncodeItem ...
func (c CopperBulb) EncodeItem() (name string, meta int16) {
	return copperBlockName("copper_bulb", c.Oxidation, c.Waxed), 0
}

// EncodeBlock ...
func (c CopperBulb) EncodeBlock() (string, map[string]any) {
	return copperBlockName("copper_bulb", c.Oxidation, c.Waxed), map[string]any{
		"lit":         c.Lit,
		"powered_bit": c.Powered,
	}
}

// allCopperBulbs returns a list of all copper bulb variants.
func allCopperBulbs() (c []world.Block) {
	for _, waxed := range []bool{false, true} {
		for _, o := range OxidationTypes() {
			for _, lit := range []bool{false, true} {
				for _, powered := range []bool{false, true} {
					c = append(c, CopperBulb{Oxidation: o, Waxed: waxed, Lit: lit, Powered: powered})
				}
			}
		}
	}
	return
}
