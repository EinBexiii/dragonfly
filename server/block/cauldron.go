package block

import (
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

// CauldronLiquid is what a cauldron holds.
type CauldronLiquid uint8

const (
	CauldronWater CauldronLiquid = iota
	CauldronLava
	CauldronPowderSnow
)

// Uint8 ...
func (l CauldronLiquid) Uint8() uint8 {
	return uint8(l)
}

// String ...
func (l CauldronLiquid) String() string {
	switch l {
	case CauldronLava:
		return "lava"
	case CauldronPowderSnow:
		return "powder_snow"
	}
	return "water"
}

// Cauldron is registered for its collision only; filling and emptying it
// are not implemented.
type Cauldron struct {
	transparent

	// Liquid is what the cauldron holds; meaningless while Level is 0.
	Liquid CauldronLiquid
	// Level is the fill level, from 0 to 6.
	Level int
}

// Model ...
func (Cauldron) Model() world.BlockModel {
	return model.Cauldron{}
}

// EncodeBlock ...
func (c Cauldron) EncodeBlock() (string, map[string]any) {
	return "minecraft:cauldron", map[string]any{"cauldron_liquid": c.Liquid.String(), "fill_level": int32(c.Level)}
}

// EncodeItem ...
func (c Cauldron) EncodeItem() (name string, meta int16) {
	name, _ = c.EncodeBlock()
	return name, 0
}

// allCauldrons ...
func allCauldrons() (cauldrons []world.Block) {
	for _, liquid := range []CauldronLiquid{CauldronWater, CauldronLava, CauldronPowderSnow} {
		for level := 0; level <= 6; level++ {
			cauldrons = append(cauldrons, Cauldron{Liquid: liquid, Level: level})
		}
	}
	return
}
