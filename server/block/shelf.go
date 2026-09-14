package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

// Shelf is a block that holds up to three item stacks on display.
type Shelf struct {
	transparent

	// Wood is the type of wood of the shelf. This field must have one of the values found in the material package.
	Wood WoodType
	// Facing is the direction the front of the shelf points in.
	Facing cube.Direction
	// Powered specifies if the shelf is currently receiving a redstone signal.
	Powered bool
	// PoweredShelfType is the raw powered_shelf_type state, ranging from 0 to 3. The BDS corpus does not name its
	// values and the native collision getter does not read it, so it is carried through verbatim.
	PoweredShelfType int
}

// Model ...
func (s Shelf) Model() world.BlockModel {
	return model.Shelf{Facing: s.Facing}
}

// EncodeBlock ...
func (s Shelf) EncodeBlock() (string, map[string]any) {
	return "minecraft:" + s.Wood.String() + "_shelf", map[string]any{
		"minecraft:cardinal_direction": s.Facing.String(),
		"powered_bit":                  s.Powered,
		"powered_shelf_type":           int32(s.PoweredShelfType),
	}
}

// allShelves returns all shelf states.
func allShelves() (shelves []world.Block) {
	for _, w := range WoodTypes() {
		for _, d := range cube.Directions() {
			for _, powered := range []bool{false, true} {
				for shelfType := range 4 {
					shelves = append(shelves, Shelf{Wood: w, Facing: d, Powered: powered, PoweredShelfType: shelfType})
				}
			}
		}
	}
	return
}
