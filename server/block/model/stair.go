package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Stair is a model for stair-like blocks. These have different solid sides depending on the direction the
// stairs are facing, the surrounding blocks and whether it is upside down or not.
type Stair struct {
	// Facing specifies the direction that the full side of the Stair faces.
	Facing cube.Direction
	// UpsideDown turns the Stair upside-down, meaning the full side of the Stair is turned to the top side of the
	// block.
	UpsideDown bool
}

// BBox returns a slice of physics.BBox depending on if the Stair is upside down and which direction it is facing.
// Additionally, these BBoxs depend on the Stair blocks surrounding this one, which can influence the model.
func (s Stair) BBox(pos cube.Pos, bs world.BlockSource) []cube.BBox {
	b := []cube.BBox{cube.Box(0, 0, 0, 1, 0.5, 1)}
	if s.UpsideDown {
		b[0] = cube.Box(0, 0.5, 0, 1, 1, 1)
	}
	t := s.Corner(pos, bs)

	face, oppositeFace := s.Facing.Face(), s.Facing.Opposite().Face()
	switch t {
	case CornerNone, CornerInnerRight, CornerInnerLeft:
		b = append(b, cube.Box(0.5, 0.5, 0.5, 0.5, 1, 0.5).
			ExtendTowards(face, 0.5).
			Stretch(s.Facing.RotateRight().Face().Axis(), 0.5))
	}
	switch t {
	case CornerOuterLeft:
		b = append(b, cube.Box(0.5, 0.5, 0.5, 0.5, 1, 0.5).
			ExtendTowards(face, 0.5).
			ExtendTowards(s.Facing.RotateLeft().Face(), 0.5))
	case CornerOuterRight:
		b = append(b, cube.Box(0.5, 0.5, 0.5, 0.5, 1, 0.5).
			ExtendTowards(face, 0.5).
			ExtendTowards(s.Facing.RotateRight().Face(), 0.5))
	case CornerInnerRight:
		b = append(b, cube.Box(0.5, 0.5, 0.5, 0.5, 1, 0.5).
			ExtendTowards(oppositeFace, 0.5).
			ExtendTowards(s.Facing.RotateRight().Face(), 0.5))
	case CornerInnerLeft:
		b = append(b, cube.Box(0.5, 0.5, 0.5, 0.5, 1, 0.5).
			ExtendTowards(oppositeFace, 0.5).
			ExtendTowards(s.Facing.RotateLeft().Face(), 0.5))
	}
	if s.UpsideDown {
		for i := range b[1:] {
			b[i+1] = b[i+1].Translate(mgl64.Vec3{0, -0.5})
		}
	}
	return b
}

// FaceSolid returns true for all faces of the Stair that are completely filled.
func (s Stair) FaceSolid(pos cube.Pos, face cube.Face, bs world.BlockSource) bool {
	// Stairs always have a closed side at the top or bottom based on their orientation.
	if (face == cube.FaceUp && s.UpsideDown) || (face == cube.FaceDown && !s.UpsideDown) {
		return true
	}

	switch t := s.Corner(pos, bs); t {
	case CornerOuterLeft, CornerOuterRight:
		// Small corner blocks, they do not block water flowing out horizontally.
		return false
	case CornerNone:
		// Not a corner, so only block directly behind the stairs.
		return s.Facing.Face() == face
	case CornerInnerRight:
		return face == s.Facing.RotateRight().Face() || face == s.Facing.Face()
	default:
		return face == s.Facing.RotateLeft().Face() || face == s.Facing.Face()
	}
}

// StairCorner is the corner stairs form with the stairs around them, named
// for the side of the stairs the corner is on.
type StairCorner uint8

const (
	CornerNone StairCorner = iota
	CornerInnerRight
	CornerInnerLeft
	CornerOuterLeft
	CornerOuterRight
)

// String is the name the client knows the corner by.
func (c StairCorner) String() string {
	switch c {
	case CornerInnerRight:
		return "inner_right"
	case CornerInnerLeft:
		return "inner_left"
	case CornerOuterLeft:
		return "outer_left"
	case CornerOuterRight:
		return "outer_right"
	}
	return "none"
}

// Corner returns the corner the stairs form with the stairs around them, or CornerNone if they form none with any
// other stairs.
func (s Stair) Corner(pos cube.Pos, bs world.BlockSource) StairCorner {
	rotatedFacing := s.Facing.RotateRight()
	if closedSide, ok := bs.Block(pos.Side(s.Facing.Face())).Model().(Stair); ok && closedSide.UpsideDown == s.UpsideDown {
		if closedSide.Facing == rotatedFacing {
			return CornerOuterRight
		} else if closedSide.Facing == rotatedFacing.Opposite() {
			// This will only form a corner if there is not a stair on the right of this one with the same
			// direction.
			if side, ok := bs.Block(pos.Side(s.Facing.RotateRight().Face())).Model().(Stair); !ok || side.Facing != s.Facing || side.UpsideDown != s.UpsideDown {
				return CornerOuterLeft
			}
			return CornerNone
		}
	}
	if openSide, ok := bs.Block(pos.Side(s.Facing.Opposite().Face())).Model().(Stair); ok && openSide.UpsideDown == s.UpsideDown {
		if openSide.Facing == rotatedFacing {
			// This will only form a corner if there is not a stair on the right of this one with the same
			// direction.
			if side, ok := bs.Block(pos.Side(s.Facing.RotateRight().Face())).Model().(Stair); !ok || side.Facing != s.Facing || side.UpsideDown != s.UpsideDown {
				return CornerInnerRight
			}
		} else if openSide.Facing == rotatedFacing.Opposite() {
			return CornerInnerLeft
		}
	}
	return CornerNone
}
