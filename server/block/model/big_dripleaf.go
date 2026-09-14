package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// DripleafTilt is how far the leaf of a BigDripleaf has tipped over. These are Bedrock's values and the only
// definition of them: block.DripleafTilt derives its ordinals from these constants and converts straight to this type,
// because model cannot import block.
type DripleafTilt uint8

const (
	// DripleafTiltNone is a leaf that is not tilted at all.
	DripleafTiltNone DripleafTilt = iota
	// DripleafTiltUnstable is a leaf that has just been stepped on and is about to tilt.
	DripleafTiltUnstable
	// DripleafTiltPartial is a leaf that has tilted halfway down.
	DripleafTiltPartial
	// DripleafTiltFull is a leaf that has tilted all the way down.
	DripleafTiltFull
)

// BigDripleaf is a model used by big dripleaves. Only the leaf on top of the plant has collision, and only while it is
// still holding itself up.
type BigDripleaf struct {
	// Head specifies if this block is the leaf on top of the plant rather than one of the stem blocks below it.
	Head bool
	// Tilt is how far the leaf has tipped over.
	Tilt DripleafTilt
}

// BBox returns the thin plate the leaf is made of. Boxes are from the 2026-09-14 BDS review, batch A section 4; the
// stem geometry the block also carries is not returned by the collision path, and cardinal direction is not read.
func (d BigDripleaf) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	if !d.Head || d.Tilt == DripleafTiltFull {
		// A stem block, and a leaf that has tipped all the way down, drop whatever is standing on them.
		return nil
	}
	if d.Tilt == DripleafTiltPartial {
		return []cube.BBox{cube.Box(0, 11.0/16, 0, 1, 13.0/16, 1)}
	}
	return []cube.BBox{cube.Box(0, 11.0/16, 0, 1, 15.0/16, 1)}
}

// FaceSolid always returns false.
func (BigDripleaf) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
