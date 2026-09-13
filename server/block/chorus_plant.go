package block

import (
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

// ChorusPlant is the stem a chorus plant grows out of in the end.
type ChorusPlant struct{}

// Model ...
func (ChorusPlant) Model() world.BlockModel {
	return model.ChorusPlant{}
}

// BreakInfo ...
func (c ChorusPlant) BreakInfo() BreakInfo {
	return newBreakInfo(0.4, alwaysHarvestable, axeEffective, oneOf(c))
}

// EncodeItem ...
func (ChorusPlant) EncodeItem() (name string, meta int16) {
	return "minecraft:chorus_plant", 0
}

// EncodeBlock ...
func (ChorusPlant) EncodeBlock() (string, map[string]any) {
	return "minecraft:chorus_plant", nil
}
