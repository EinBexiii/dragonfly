package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// StoneButton represents vanilla stone button states without redstone mechanics.
type StoneButton struct {
	empty
	transparent

	// Blackstone specifies polished blackstone instead of stone.
	Blackstone bool
	// Facing is the facing_direction state, from 0 to 5.
	Facing cube.Face
	// Pressed is the button_pressed_bit state.
	Pressed bool
}

// EncodeBlock ...
func (b StoneButton) EncodeBlock() (string, map[string]any) {
	name := "stone"
	if b.Blackstone {
		name = "polished_blackstone"
	}
	return "minecraft:" + name + "_button", map[string]any{"facing_direction": int32(b.Facing), "button_pressed_bit": boolByte(b.Pressed)}
}

// allStoneButtons ...
func allStoneButtons() (blocks []world.Block) {
	for _, material := range []bool{false, true} {
		for facing := cube.Face(0); facing < 6; facing++ {
			for _, pressed := range []bool{false, true} {
				blocks = append(blocks, StoneButton{Blackstone: material, Facing: facing, Pressed: pressed})
			}
		}
	}
	return
}
