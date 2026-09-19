package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// PitcherCrop is a model used by pitcher crops. Only the lower half of the plant collides, with a pod that widens once
// the crop has grown past its first stage.
type PitcherCrop struct {
	// UpperPart specifies if this block is the upper half of the plant.
	UpperPart bool
	// Sprouted specifies if the crop has grown past growth stage 0.
	Sprouted bool
}

// BBox returns the pod of the crop. Boxes are from the 2026-09-14 BDS review, batch B section 12, which distinguishes
// growth 0 from any positive growth rather than every stage. The y minimum of -1/16 is what the binary holds: the pod
// really does reach down into the farmland below.
func (p PitcherCrop) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	if p.UpperPart {
		return nil
	}
	if p.Sprouted {
		return []cube.BBox{cube.Box(3.0/16, -1.0/16, 3.0/16, 13.0/16, 5.0/16, 13.0/16)}
	}
	return []cube.BBox{cube.Box(5.0/16, -1.0/16, 5.0/16, 11.0/16, 3.0/16, 11.0/16)}
}

// FaceSolid always returns false.
func (PitcherCrop) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
