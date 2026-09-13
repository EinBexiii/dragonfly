package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// PistonArm is the model of an extended piston's arm: a head plate against the block it pushes, a shaft down
// the middle and a connector reaching back into the piston base behind it.
type PistonArm struct {
	// Facing is the arm's facing_direction state value, used to index pistonArmBoxes numerically rather than
	// by cube.Face name. The piston's own direction table is 0=-Y, 1=+Y, 2=+Z, 3=-Z, 4=+X, 5=-X, which is
	// mirrored against cube.Face on both horizontal axes, so reading this as a named face would turn every
	// horizontal arm around.
	Facing cube.Face
}

// pistonArmBoxes holds the head plate, shaft and rear connector of a fully extended arm for each of the six
// facing_direction values, in that order. The rear connector reaches 0.25 beyond the block on the side the arm
// came from: the Bedrock server emits it that way and the client collides with it, so it is not clipped here.
var pistonArmBoxes = [6][3]cube.BBox{
	{
		cube.Box(0, 0, 0, 1, 0.25, 1),
		cube.Box(0.375, 0.25, 0.375, 0.625, 0.75, 0.625),
		cube.Box(0.3125, 0.75, 0.3125, 0.6875, 1.25, 0.6875),
	},
	{
		cube.Box(0, 0.75, 0, 1, 1, 1),
		cube.Box(0.375, 0.25, 0.375, 0.625, 0.75, 0.625),
		cube.Box(0.3125, -0.25, 0.3125, 0.6875, 0.25, 0.6875),
	},
	{
		cube.Box(0, 0, 0.75, 1, 1, 1),
		cube.Box(0.375, 0.375, 0.25, 0.625, 0.625, 0.75),
		cube.Box(0.3125, 0.3125, -0.25, 0.6875, 0.6875, 0.25),
	},
	{
		cube.Box(0, 0, 0, 1, 1, 0.25),
		cube.Box(0.375, 0.375, 0.25, 0.625, 0.625, 0.75),
		cube.Box(0.3125, 0.3125, 0.75, 0.6875, 0.6875, 1.25),
	},
	{
		cube.Box(0.75, 0, 0, 1, 1, 1),
		cube.Box(0.25, 0.375, 0.375, 0.75, 0.625, 0.625),
		cube.Box(-0.25, 0.3125, 0.3125, 0.25, 0.6875, 0.6875),
	},
	{
		cube.Box(0, 0, 0, 0.25, 1, 1),
		cube.Box(0.25, 0.375, 0.375, 0.75, 0.625, 0.625),
		cube.Box(0.75, 0.3125, 0.3125, 1.25, 0.6875, 0.6875),
	},
}

// BBox returns the three boxes of the arm facing the direction it was built with.
func (p PistonArm) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	boxes := pistonArmBoxes[p.Facing]
	return boxes[:]
}

// FaceSolid always returns false.
func (PistonArm) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
