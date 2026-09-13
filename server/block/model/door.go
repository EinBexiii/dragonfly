package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// doorThickness is the leaf as the Bedrock server has it, the same float32
// as a trapdoor; Java's 3/16 leaves a sliver the client walks into.
const doorThickness = float64(float32(0.1825))

// Door is a model used for doors. It has no solid faces and a bounding box that changes depending on
// the direction of the door, whether it is open, and the side of its hinge.
type Door struct {
	// Facing is the direction that the door is facing when closed.
	Facing cube.Direction
	// Open specifies if the Door is open. The direction it opens towards depends on the Right field.
	Open bool
	// Right specifies the attachment side of the door and, with that, the direction it opens in.
	Right bool
	// Top specifies the upper half of the door.
	Top bool
}

// BBox returns a physics.BBox that depends on if the Door is open, what direction it is facing and whether it is
// attached to the right/left side of a block. Bedrock reads the direction and open state from the lower half
// and the hinge from the upper half, whatever each half stores, so both halves get the same leaf.
func (d Door) BBox(pos cube.Pos, s world.BlockSource) []cube.BBox {
	lower, upper := d.halves(pos, s)
	if !lower.Open {
		return []cube.BBox{full.ExtendTowards(lower.Facing.Face(), doorThickness-1)}
	}
	if upper.Right {
		return []cube.BBox{full.ExtendTowards(lower.Facing.RotateLeft().Face(), doorThickness-1)}
	}
	return []cube.BBox{full.ExtendTowards(lower.Facing.RotateRight().Face(), doorThickness-1)}
}

// halves returns the lower and upper half of the door at pos. A half without a door as its partner is
// judged as Bedrock does: closed, hinge left, facing the first cardinal direction.
func (d Door) halves(pos cube.Pos, s world.BlockSource) (lower, upper Door) {
	partnerFace := cube.FaceUp
	if d.Top {
		partnerFace = cube.FaceDown
	}
	partner, ok := s.Block(pos.Side(partnerFace)).Model().(Door)
	if !ok {
		orphan := Door{Facing: cube.South.RotateLeft()}
		return orphan, orphan
	}
	if d.Top {
		return partner, d
	}
	return d, partner
}

// FaceSolid always returns false.
func (d Door) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
