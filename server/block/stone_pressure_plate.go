package block

import "github.com/df-mc/dragonfly/server/world"

// StonePressurePlate represents vanilla stone pressure plate states without redstone mechanics.
type StonePressurePlate struct {
	empty
	transparent

	// Blackstone specifies polished blackstone instead of stone.
	Blackstone bool
	// Power is the redstone_signal state, from 0 to 15.
	Power int
}

// EncodeBlock ...
func (b StonePressurePlate) EncodeBlock() (string, map[string]any) {
	name := "stone"
	if b.Blackstone {
		name = "polished_blackstone"
	}
	return "minecraft:" + name + "_pressure_plate", map[string]any{"redstone_signal": int32(b.Power)}
}

// allStonePressurePlates ...
func allStonePressurePlates() (blocks []world.Block) {
	for _, material := range []bool{false, true} {
		for power := 0; power < 16; power++ {
			blocks = append(blocks, StonePressurePlate{Blackstone: material, Power: power})
		}
	}
	return
}
