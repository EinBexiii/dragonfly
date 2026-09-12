package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// trapdoorThickness is the thickness of a trapdoor's collision box as the
// Bedrock server has it, 2.92 of 16. Java's 3/16 leaves a player who lands
// on the trapdoor standing inside the box rather than on it.
const trapdoorThickness = float64(float32(0.1825))

// Trapdoor is a model used for trapdoors. It has no solid faces and a bounding box that changes depending on
// the direction of the trapdoor.
type Trapdoor struct {
	// Facing is the facing direction of the Trapdoor. In addition to the texture, it influences the direction in which
	// the Trapdoor is opened.
	Facing cube.Direction
	// Open and Top specify if the Trapdoor is opened and if it's in the top or bottom part of a block respectively.
	Open, Top bool
}

// BBox returns a physics.BBox that depends on the facing direction of the Trapdoor and whether it is open and in the
// top part of the block.
func (t Trapdoor) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	if t.Open {
		return []cube.BBox{full.ExtendTowards(t.Facing.Face(), trapdoorThickness-1)}
	} else if t.Top {
		return []cube.BBox{cube.Box(0, 1-trapdoorThickness, 0, 1, 1, 1)}
	}
	return []cube.BBox{cube.Box(0, 0, 0, 1, trapdoorThickness, 1)}
}

// FaceSolid returns true if the face is completely filled with the trapdoor.
func (t Trapdoor) FaceSolid(_ cube.Pos, face cube.Face, _ world.BlockSource) bool {
	if t.Open {
		return t.Facing.Face().Opposite() == face
	} else if t.Top {
		return face == cube.FaceUp
	}
	return face == cube.FaceDown
}
