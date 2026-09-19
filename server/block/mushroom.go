package block

import (
	"github.com/df-mc/dragonfly/server/world"
)

// Mushroom is a small fungus that grows in the dark.
type Mushroom struct {
	transparent
	empty

	// Red specifies if the mushroom is a red mushroom rather than a brown one.
	Red bool
}

// BreakInfo ...
func (m Mushroom) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, oneOf(m))
}

// EncodeItem ...
func (m Mushroom) EncodeItem() (name string, meta int16) {
	if m.Red {
		return "minecraft:red_mushroom", 0
	}
	return "minecraft:brown_mushroom", 0
}

// EncodeBlock ...
func (m Mushroom) EncodeBlock() (string, map[string]any) {
	if m.Red {
		return "minecraft:red_mushroom", nil
	}
	return "minecraft:brown_mushroom", nil
}

// allMushrooms ...
func allMushrooms() []world.Block {
	return []world.Block{Mushroom{}, Mushroom{Red: true}}
}
