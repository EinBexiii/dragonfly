package block

import "github.com/df-mc/dragonfly/server/world"

// PoplarPressurePlate represents vanilla poplar pressure plate states without redstone mechanics.
type PoplarPressurePlate struct {
	empty
	transparent

	// Power is the redstone_signal state, from 0 to 15.
	Power int
}

// EncodeBlock ...
func (b PoplarPressurePlate) EncodeBlock() (string, map[string]any) {
	return "minecraft:poplar_pressure_plate", map[string]any{"redstone_signal": int32(b.Power)}
}

// allPoplarPressurePlates ...
func allPoplarPressurePlates() (blocks []world.Block) {
	for power := 0; power < 16; power++ {
		blocks = append(blocks, PoplarPressurePlate{Power: power})
	}
	return
}
