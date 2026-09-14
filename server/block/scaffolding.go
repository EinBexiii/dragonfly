package block

import (
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

// Scaffolding is a climbable block that can be built out sideways from a supporting column.
type Scaffolding struct {
	transparent

	// Stability is how far the block is from the scaffolding column holding it up, from 0 to 7. At 7 the block is
	// unsupported and falls.
	Stability int
	// StabilityCheck specifies if the block still has to recalculate its stability.
	StabilityCheck bool
}

// Model ...
func (s Scaffolding) Model() world.BlockModel {
	return model.Scaffolding{Stability: s.Stability}
}

// EncodeBlock ...
func (s Scaffolding) EncodeBlock() (string, map[string]any) {
	return "minecraft:scaffolding", map[string]any{
		"stability":       int32(s.Stability),
		"stability_check": s.StabilityCheck,
	}
}

// allScaffolding returns all scaffolding states.
func allScaffolding() (scaffolding []world.Block) {
	for stability := range 8 {
		for _, check := range []bool{false, true} {
			scaffolding = append(scaffolding, Scaffolding{Stability: stability, StabilityCheck: check})
		}
	}
	return
}
