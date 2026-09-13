package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// Dispenser is a block that shoots or places the items it holds. The lobby only uses it as decoration, so it
// holds no inventory and cannot be activated.
type Dispenser struct {
	solid

	// Facing is the face the dispenser's barrel points towards.
	Facing cube.Face
	// Triggered is whether the dispenser is currently powered by redstone.
	Triggered bool
}

// BreakInfo ...
func (d Dispenser) BreakInfo() BreakInfo {
	return newBreakInfo(3.5, pickaxeHarvestable, pickaxeEffective, oneOf(d))
}

// EncodeItem ...
func (Dispenser) EncodeItem() (name string, meta int16) {
	return "minecraft:dispenser", 0
}

// EncodeBlock ...
func (d Dispenser) EncodeBlock() (string, map[string]any) {
	return "minecraft:dispenser", map[string]any{"facing_direction": int32(d.Facing), "triggered_bit": d.Triggered}
}

// allDispensers ...
func allDispensers() (dispensers []world.Block) {
	for _, f := range cube.Faces() {
		dispensers = append(dispensers, Dispenser{Facing: f})
		dispensers = append(dispensers, Dispenser{Facing: f, Triggered: true})
	}
	return
}
