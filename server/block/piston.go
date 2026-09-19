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
		// The base keeps a full cube whether it is retracted or extended: the Bedrock server never shortens
		// it to 0.75 on its facing side, so neither does this.
		solid

		// Sticky specifies if the piston is a sticky piston, which pulls blocks back when it retracts.
		Sticky bool
		// Facing is Bedrock's facing_direction value as the world stores it. Pistons use their own table
		// (0=-Y 1=+Y 2=+Z 3=-Z 4=+X 5=-X), whose horizontal values are mirrored against cube.Face, so the
		// number is never read as a face.
		Facing cube.Face
	}

	// PistonArm is the arm and head an extended Piston leaves in the block in front of it.
	PistonArm struct {
		// Sticky specifies if the arm belongs to a sticky piston.
		Sticky bool
		// Facing is the piston's facing_direction value, with the same table as Piston.Facing.
		Facing cube.Face
	}
)

// Model ...
func (p PistonArm) Model() world.BlockModel {
	return model.PistonArm{Facing: p.Facing}
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
	properties := map[string]any{"facing_direction": int32(p.Facing)}
	if p.Sticky {
		return "minecraft:sticky_piston", properties
	}
	return "minecraft:piston", properties
}

// EncodeBlock ...
func (p PistonArm) EncodeBlock() (string, map[string]any) {
	properties := map[string]any{"facing_direction": int32(p.Facing)}
	if p.Sticky {
		return "minecraft:sticky_piston_arm_collision", properties
	}
	return "minecraft:piston_arm_collision", properties
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
