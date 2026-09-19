package block

import (
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

// FlowerPot is an empty flower pot, registered for its collision; the
// plant it may hold lives in its block entity and is not modelled here.
type FlowerPot struct {
	transparent

	// Update is the update_bit state.
	Update bool
}

// Model ...
func (FlowerPot) Model() world.BlockModel {
	return model.FlowerPot{}
}

// EncodeBlock ...
func (f FlowerPot) EncodeBlock() (string, map[string]any) {
	return "minecraft:flower_pot", map[string]any{"update_bit": boolByte(f.Update)}
}

// EncodeItem ...
func (f FlowerPot) EncodeItem() (name string, meta int16) {
	name, _ = f.EncodeBlock()
	return name, 0
}

// allFlowerPots ...
func allFlowerPots() []world.Block {
	return []world.Block{FlowerPot{}, FlowerPot{Update: true}}
}
