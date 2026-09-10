package block

import "github.com/df-mc/dragonfly/server/world"

// WoodPressurePlate represents vanilla wood pressure plate states without redstone mechanics.
type WoodPressurePlate struct {
	empty
	transparent

	// Wood is the type of wood.
	Wood WoodType
	// Power is the redstone_signal state, from 0 to 15.
	Power int
}

// EncodeBlock ...
func (b WoodPressurePlate) EncodeBlock() (string, map[string]any) {
	name := b.Wood.String()
	if b.Wood == OakWood() {
		name = "wooden"
	}
	return "minecraft:" + name + "_pressure_plate", map[string]any{"redstone_signal": int32(b.Power)}
}

// allWoodPressurePlates ...
func allWoodPressurePlates() (blocks []world.Block) {
	for _, material := range WoodTypes() {
		for power := 0; power < 16; power++ {
			blocks = append(blocks, WoodPressurePlate{Wood: material, Power: power})
		}
	}
	return
}
