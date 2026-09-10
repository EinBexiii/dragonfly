package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// WoodButton represents vanilla wood button states without redstone mechanics.
type WoodButton struct {
	empty
	transparent

	// Wood is the type of wood.
	Wood WoodType
	// Facing is the facing_direction state, from 0 to 5.
	Facing cube.Face
	// Pressed is the button_pressed_bit state.
	Pressed bool
}

// EncodeBlock ...
func (b WoodButton) EncodeBlock() (string, map[string]any) {
	name := b.Wood.String()
	if b.Wood == OakWood() {
		name = "wooden"
	}
	return "minecraft:" + name + "_button", map[string]any{"facing_direction": int32(b.Facing), "button_pressed_bit": boolByte(b.Pressed)}
}

// allWoodButtons ...
func allWoodButtons() (blocks []world.Block) {
	for _, material := range WoodTypes() {
		for facing := cube.Face(0); facing < 6; facing++ {
			for _, pressed := range []bool{false, true} {
				blocks = append(blocks, WoodButton{Wood: material, Facing: facing, Pressed: pressed})
			}
		}
	}
	return
}
