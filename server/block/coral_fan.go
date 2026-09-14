package block

import (
	"github.com/df-mc/dragonfly/server/world"
)

// CoralFan is a fan of coral growing on the floor. Like every coral fan it has no collision at all: the 2026-09-14 BDS
// review, batch A section 7, finds the floor and wall constructors both installing the empty collision getter, with no
// separate implementation for the dead variants.
type CoralFan struct {
	empty
	transparent

	// Type is the type of coral of the block.
	Type CoralType
	// Dead is whether the coral fan is dead.
	Dead bool
	// Direction is the raw coral_fan_direction state, either 0 or 1. The corpus does not resolve what the two values
	// mean and the collision getter does not read them.
	Direction int
}

// EncodeBlock ...
func (c CoralFan) EncodeBlock() (string, map[string]any) {
	return coralName(c.Type, c.Dead, "coral_fan"), map[string]any{"coral_fan_direction": int32(c.Direction)}
}

// EncodeItem ...
func (c CoralFan) EncodeItem() (name string, meta int16) {
	return blockItemName(c), 0
}

// CoralWallFan is a fan of coral growing on the side of a block.
type CoralWallFan struct {
	empty
	transparent

	// Type is the type of coral of the block.
	Type CoralType
	// Dead is whether the coral wall fan is dead.
	Dead bool
	// Direction is the raw coral_direction state, ranging from 0 to 3. The corpus does not map the values onto cardinal
	// directions and the collision getter does not read them.
	Direction int
}

// EncodeBlock ...
func (c CoralWallFan) EncodeBlock() (string, map[string]any) {
	return coralName(c.Type, c.Dead, "coral_wall_fan"), map[string]any{"coral_direction": int32(c.Direction)}
}

// EncodeItem ...
func (c CoralWallFan) EncodeItem() (name string, meta int16) {
	return blockItemName(c), 0
}

// coralName returns the block name of a coral block of the type and suffix passed, dead or alive.
func coralName(t CoralType, dead bool, suffix string) string {
	if dead {
		return "minecraft:dead_" + t.String() + "_" + suffix
	}
	return "minecraft:" + t.String() + "_" + suffix
}

// allCoralFans returns all coral fan and coral wall fan states.
func allCoralFans() (fans []world.Block) {
	for _, t := range CoralTypes() {
		for _, dead := range []bool{false, true} {
			for direction := range 2 {
				fans = append(fans, CoralFan{Type: t, Dead: dead, Direction: direction})
			}
			for direction := range 4 {
				fans = append(fans, CoralWallFan{Type: t, Dead: dead, Direction: direction})
			}
		}
	}
	return
}
