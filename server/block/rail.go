package block

import "github.com/df-mc/dragonfly/server/world"

// Rail represents the vanilla rail states without rail mechanics.
type Rail struct {
	empty
	transparent

	// Direction is the Bedrock rail direction, from 0 to 9.
	Direction int
}

// EncodeBlock ...
func (r Rail) EncodeBlock() (string, map[string]any) {
	return "minecraft:rail", map[string]any{"rail_direction": int32(r.Direction)}
}

// allRails ...
func allRails() (rails []world.Block) {
	for direction := 0; direction < 10; direction++ {
		rails = append(rails, Rail{Direction: direction})
	}
	return
}
