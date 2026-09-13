package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

type (
	// Piston is a block that pushes the blocks in front of it when powered. It does not move anything here:
	// the lobby only uses it as decoration.
	Piston struct {
		solid

		// Sticky specifies if the piston is a sticky piston, which pulls blocks back when it retracts.
		Sticky bool
		// Facing is the face the piston's head pushes towards.
		Facing cube.Face
	}

	// PistonArm is the arm and head an extended Piston leaves in the block in front of it.
	PistonArm struct {
		// Sticky specifies if the arm belongs to a sticky piston.
		Sticky bool
		// Facing is the face the arm points towards.
		Facing cube.Face
	}
)

// Model ...
func (PistonArm) Model() world.BlockModel {
	return model.PistonArm{}
}

// BreakInfo ...
func (p Piston) BreakInfo() BreakInfo {
	return newBreakInfo(1.5, alwaysHarvestable, pickaxeEffective, oneOf(p))
}

// BreakInfo ...
func (PistonArm) BreakInfo() BreakInfo {
	return newBreakInfo(1.5, alwaysHarvestable, pickaxeEffective, simpleDrops())
}

// EncodeItem ...
func (p Piston) EncodeItem() (name string, meta int16) {
	if p.Sticky {
		return "minecraft:sticky_piston", 0
	}
	return "minecraft:piston", 0
}

// EncodeBlock ...
func (p Piston) EncodeBlock() (string, map[string]any) {
	if p.Sticky {
		return "minecraft:sticky_piston", map[string]any{"facing_direction": int32(p.Facing)}
	}
	return "minecraft:piston", map[string]any{"facing_direction": int32(p.Facing)}
}

// EncodeBlock ...
func (p PistonArm) EncodeBlock() (string, map[string]any) {
	if p.Sticky {
		return "minecraft:sticky_piston_arm_collision", map[string]any{"facing_direction": int32(p.Facing)}
	}
	return "minecraft:piston_arm_collision", map[string]any{"facing_direction": int32(p.Facing)}
}

// allPistons ...
func allPistons() (pistons []world.Block) {
	for _, f := range cube.Faces() {
		pistons = append(pistons, Piston{Facing: f})
		pistons = append(pistons, Piston{Sticky: true, Facing: f})
		pistons = append(pistons, PistonArm{Facing: f})
		pistons = append(pistons, PistonArm{Sticky: true, Facing: f})
	}
	return
}
