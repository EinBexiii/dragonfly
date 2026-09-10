package block

import (
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

// SnowLayer is snow with one to eight layers.
type SnowLayer struct {
	transparent

	// Layers is the number of layers, from 1 to 8.
	Layers int
	// Covered is the covered_bit state.
	Covered bool
}

// Model ...
func (s SnowLayer) Model() world.BlockModel {
	return model.SnowLayer{Layers: s.Layers}
}

// EncodeBlock ...
func (s SnowLayer) EncodeBlock() (string, map[string]any) {
	return "minecraft:snow_layer", map[string]any{"height": int32(s.Layers - 1), "covered_bit": boolByte(s.Covered)}
}

// allSnowLayers ...
func allSnowLayers() (layers []world.Block) {
	for n := 1; n <= 8; n++ {
		for _, covered := range []bool{false, true} {
			layers = append(layers, SnowLayer{Layers: n, Covered: covered})
		}
	}
	return
}
