package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// PoplarButton represents vanilla poplar button states without redstone mechanics.
type PoplarButton struct {
	empty
	transparent

	// Facing is the facing_direction state, from 0 to 5.
	Facing cube.Face
	// Pressed is the button_pressed_bit state.
	Pressed bool
}

// EncodeBlock ...
func (b PoplarButton) EncodeBlock() (string, map[string]any) {
	return "minecraft:poplar_button", map[string]any{"facing_direction": int32(b.Facing), "button_pressed_bit": boolByte(b.Pressed)}
}

// allPoplarButtons ...
func allPoplarButtons() (blocks []world.Block) {
	for facing := cube.Face(0); facing < 6; facing++ {
		for _, pressed := range []bool{false, true} {
			blocks = append(blocks, PoplarButton{Facing: facing, Pressed: pressed})
		}
	}
	return
}
