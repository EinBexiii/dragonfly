package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// Fence is a model used by fences of any type. It can attach to blocks with solid faces and other fences of the same
// type and has a model height just slightly over 1.
type Fence struct {
	// Wood specifies if the Fence is made from wood. This field is used to check if two fences are able to attach to
	// each other.
	Wood bool
}

// BarrierHeight is how tall a fence or a wall is to an entity: a block and a
// half, so neither can be jumped over, whatever their drawn shape. Entity
// movement searches that far below a box for blocks that reach into it.
const BarrierHeight = 1.5

const fenceInset = 0.375

// BBox returns multiple physics.BBox depending on how many connections it has with the surrounding blocks.
func (f Fence) BBox(pos cube.Pos, s world.BlockSource) []cube.BBox {
	boxes := make([]cube.BBox, 0, 2)

	connectWest, connectEast := f.Connects(pos, cube.FaceWest, s), f.Connects(pos, cube.FaceEast, s)
	connectNorth, connectSouth := f.Connects(pos, cube.FaceNorth, s), f.Connects(pos, cube.FaceSouth, s)

	// Check if we have any connections on the Z axis
	if connectWest || connectEast {
		sideBox := cube.Box(0, 0, 0, 1, BarrierHeight, 1).Stretch(cube.Z, -fenceInset)
		if connectWest {
			boxes = append(boxes, sideBox.ExtendTowards(cube.FaceEast, -fenceInset))
		}
		if connectEast {
			boxes = append(boxes, sideBox.ExtendTowards(cube.FaceWest, -fenceInset))
		}
	}

	// Check if we have any connections on the X axis
	if connectNorth || connectSouth {
		sideBox := cube.Box(0, 0, 0, 1, BarrierHeight, 1).Stretch(cube.X, -fenceInset)
		if connectNorth {
			boxes = append(boxes, sideBox.ExtendTowards(cube.FaceSouth, -fenceInset))
		}
		if connectSouth {
			boxes = append(boxes, sideBox.ExtendTowards(cube.FaceNorth, -fenceInset))
		}
	}

	// If no connections, create a center post box
	if len(boxes) == 0 {
		boxes = append(boxes, cube.Box(fenceInset, 0, fenceInset, 1-fenceInset, BarrierHeight, 1-fenceInset))
	}

	return boxes
}

// FaceSolid returns true if the face is cube.FaceDown or cube.FaceUp.
func (f Fence) FaceSolid(_ cube.Pos, face cube.Face, _ world.BlockSource) bool {
	return face == cube.FaceDown || face == cube.FaceUp
}

// Connects reports whether the fence at pos joins the block on the given face.
func (f Fence) Connects(pos cube.Pos, face cube.Face, src world.BlockSource) bool {
	sidePos := pos.Side(face)
	sideBlock := src.Block(sidePos)
	if fence, ok := sideBlock.Model().(Fence); ok && fence.Wood == f.Wood {
		return true
	}
	if sideBlock.Model().FaceSolid(sidePos, face.Opposite(), src) {
		return true
	}
	// A gate joins a fence only at the ends of its bar: Bedrock and Java both
	// leave the faces the gate swings through unconnected, open or closed.
	gate, ok := sideBlock.Model().(FenceGate)
	return ok && gate.Facing.Face().Axis() != face.Axis()
}
