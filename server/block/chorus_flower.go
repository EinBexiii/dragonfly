package block

import (
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

// ChorusFlower is the bud on top of a chorus plant. It does not grow here: the lobby only uses it as
// decoration.
type ChorusFlower struct {
	// Age is how far the flower has grown, 0-5. 5 is dead.
	Age int
}

// Model ...
func (ChorusFlower) Model() world.BlockModel {
	return model.ChorusFlower{}
}

// BreakInfo ...
func (c ChorusFlower) BreakInfo() BreakInfo {
	return newBreakInfo(0.4, alwaysHarvestable, axeEffective, oneOf(c))
}

// EncodeItem ...
func (ChorusFlower) EncodeItem() (name string, meta int16) {
	return "minecraft:chorus_flower", 0
}

// EncodeBlock ...
func (c ChorusFlower) EncodeBlock() (string, map[string]any) {
	return "minecraft:chorus_flower", map[string]any{"age": int32(c.Age)}
}

// allChorusFlowers ...
func allChorusFlowers() (flowers []world.Block) {
	for age := 0; age < 6; age++ {
		flowers = append(flowers, ChorusFlower{Age: age})
	}
	return
}
