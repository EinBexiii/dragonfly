package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// CopperBulb is a light source that a redstone pulse switches on and off. It is a full cube in every state: the
// 2026-09-14 BDS review, batch B section 1, finds all eight variants inheriting the unit-box getter, which reads
// neither lit nor powered_bit.
// Oxidation and Waxed are decoded state and nothing else: this block has no BreakInfo, no waxing or scraping and
// no random oxidation tick, because a loaded world must come back exactly as it was written.
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

// RedstonePowerUpdate toggles the bulb on the rising edge of the signal reaching it: a powered bulb lights, and
// stays lit until it is powered a second time.
func (c CopperBulb) RedstonePowerUpdate(_ cube.Pos, _ *world.Tx, power int) (world.Block, bool) {
	powered := power > 0
	if powered == c.Powered {
		return c, false
	}
	if c.Powered = powered; powered {
		c.Lit = !c.Lit
	}
	return c, true
}

// LightEmissionLevel returns the light a lit bulb gives off, which its oxidation dims. The levels are the vanilla
// ones; the BDS corpus was read for collision, not for light.
func (c CopperBulb) LightEmissionLevel() uint8 {
	if !c.Lit {
		return 0
	}
	switch c.Oxidation {
	case UnoxidisedOxidation():
		return 15
	case ExposedOxidation():
		return 12
	case WeatheredOxidation():
		return 8
	}
	return 4
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
