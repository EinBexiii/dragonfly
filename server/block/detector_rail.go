package block

import "github.com/df-mc/dragonfly/server/world"

// DetectorRail represents the vanilla detector_rail states without rail mechanics.
type DetectorRail struct {
	empty
	transparent

	// Direction is the Bedrock rail direction, from 0 to 5.
	Direction int
	// Powered is the rail_data_bit state.
	Powered bool
}

// EncodeBlock ...
func (r DetectorRail) EncodeBlock() (string, map[string]any) {
	return "minecraft:detector_rail", map[string]any{"rail_direction": int32(r.Direction), "rail_data_bit": boolByte(r.Powered)}
}

// allDetectorRails ...
func allDetectorRails() (rails []world.Block) {
	for direction := 0; direction < 6; direction++ {
		for _, powered := range []bool{false, true} {
			rails = append(rails, DetectorRail{Direction: direction, Powered: powered})
		}
	}
	return
}
