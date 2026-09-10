package block

import "github.com/df-mc/dragonfly/server/world"

// WeightedPressurePlate represents vanilla weighted pressure plate states without redstone mechanics.
type WeightedPressurePlate struct {
	empty
	transparent

	// Heavy specifies an iron plate instead of a gold plate.
	Heavy bool
	// Power is the redstone_signal state, from 0 to 15.
	Power int
}

// EncodeBlock ...
func (b WeightedPressurePlate) EncodeBlock() (string, map[string]any) {
	name := "light_weighted"
	if b.Heavy {
		name = "heavy_weighted"
	}
	return "minecraft:" + name + "_pressure_plate", map[string]any{"redstone_signal": int32(b.Power)}
}

// allWeightedPressurePlates ...
func allWeightedPressurePlates() (blocks []world.Block) {
	for _, material := range []bool{false, true} {
		for power := 0; power < 16; power++ {
			blocks = append(blocks, WeightedPressurePlate{Heavy: material, Power: power})
		}
	}
	return
}
