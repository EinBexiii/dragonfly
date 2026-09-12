package block

import (
	"math"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

// IronTrapdoor is a trapdoor that only redstone opens, registered for its
// collision: the Bedrock server gives it the wooden trapdoor's shape.
type IronTrapdoor struct {
	transparent
	bass

	// Facing is the direction the trapdoor opens towards.
	Facing cube.Direction
	// Open is whether the trapdoor is open.
	Open bool
	// Top is whether the trapdoor sits in the top half of its block.
	Top bool
}

// Model ...
func (t IronTrapdoor) Model() world.BlockModel {
	return model.Trapdoor{Facing: t.Facing, Top: t.Top, Open: t.Open}
}

// EncodeBlock ...
func (t IronTrapdoor) EncodeBlock() (string, map[string]any) {
	return "minecraft:iron_trapdoor", map[string]any{"direction": int32(math.Abs(float64(t.Facing) - 3)), "open_bit": t.Open, "upside_down_bit": t.Top}
}

// allIronTrapdoors ...
func allIronTrapdoors() (trapdoors []world.Block) {
	for facing := cube.Direction(0); facing <= 3; facing++ {
		for _, open := range []bool{false, true} {
			for _, top := range []bool{false, true} {
				trapdoors = append(trapdoors, IronTrapdoor{Facing: facing, Open: open, Top: top})
			}
		}
	}
	return
}
