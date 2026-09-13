package block

import (
	"github.com/df-mc/dragonfly/server/item"
)

// RedstoneLamp is a block that emits light when powered by redstone.
type RedstoneLamp struct {
	solid

	// Lit specifies if the redstone lamp is lit and emitting light.
	Lit bool
}

// LightEmissionLevel ...
func (l RedstoneLamp) LightEmissionLevel() uint8 {
	if l.Lit {
		return 15
	}
	return 0
}

// BreakInfo ...
func (RedstoneLamp) BreakInfo() BreakInfo {
	return newBreakInfo(0.3, alwaysHarvestable, nothingEffective, simpleDrops(item.NewStack(RedstoneLamp{}, 1)))
}

// EncodeItem ...
func (RedstoneLamp) EncodeItem() (name string, meta int16) {
	return "minecraft:redstone_lamp", 0
}

// EncodeBlock ...
func (l RedstoneLamp) EncodeBlock() (string, map[string]any) {
	if l.Lit {
		return "minecraft:lit_redstone_lamp", nil
	}
	return "minecraft:redstone_lamp", nil
}
