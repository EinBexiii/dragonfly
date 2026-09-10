package block

import "github.com/df-mc/dragonfly/server/world"

// ActivatorRail represents the vanilla activator_rail states without rail mechanics.
type ActivatorRail struct {
	empty
	transparent

	// Direction is the Bedrock rail direction, from 0 to 5.
	Direction int
	// Powered is the rail_data_bit state.
	Powered bool
}

// EncodeBlock ...
func (r ActivatorRail) EncodeBlock() (string, map[string]any) {
	return "minecraft:activator_rail", map[string]any{"rail_direction": int32(r.Direction), "rail_data_bit": boolByte(r.Powered)}
}

// allActivatorRails ...
func allActivatorRails() (rails []world.Block) {
	for direction := 0; direction < 6; direction++ {
		for _, powered := range []bool{false, true} {
			rails = append(rails, ActivatorRail{Direction: direction, Powered: powered})
		}
	}
	return
}
