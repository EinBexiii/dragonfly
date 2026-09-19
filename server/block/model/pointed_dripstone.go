package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// DripstoneThickness is the segment a pointed dripstone block forms in the spike it is part of. These are Bedrock's
// values and the only definition of them: block.DripstoneThickness derives its ordinals from these constants and
// converts straight to this type, because model cannot import block.
type DripstoneThickness uint8

const (
	// DripstoneTip is the pointed end of a spike.
	DripstoneTip DripstoneThickness = iota
	// DripstoneFrustum is the segment right behind the tip.
	DripstoneFrustum
	// DripstoneMiddle is a segment in the middle of a spike.
	DripstoneMiddle
	// DripstoneBase is the segment a spike grows from.
	DripstoneBase
	// DripstoneMerge is the tip of a spike that has grown into a spike coming from the other side.
	DripstoneMerge
)

// PointedDripstone is a model used by pointed dripstone. Each segment is a square column that narrows towards the tip.
type PointedDripstone struct {
	// Thickness is the segment this block forms in the spike it is part of.
	Thickness DripstoneThickness
	// Hanging specifies if the spike grows down from a ceiling rather than up from a floor.
	Hanging bool
}

// BBox returns the column of the segment. Insets and the shortened tip come from the 2026-09-14 BDS review, batch A
// section 5; only the tip branch reads hanging.
func (d PointedDripstone) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	box := full.Stretch(cube.X, -d.inset()).Stretch(cube.Z, -d.inset())
	if d.Thickness == DripstoneTip {
		// The point itself is 11/16 long and sits at the end of the block the spike grows away from.
		return []cube.BBox{box.ExtendTowards(d.tipFace(), -5.0/16)}
	}
	return []cube.BBox{box}
}

// inset returns how far the sides of the segment are set in from the sides of the block.
func (d PointedDripstone) inset() float64 {
	switch d.Thickness {
	case DripstoneFrustum:
		return 4.0 / 16
	case DripstoneMiddle:
		return 3.0 / 16
	case DripstoneBase:
		return 2.0 / 16
	}
	// A tip and a merged tip are the two narrowest segments.
	return 5.0 / 16
}

// tipFace returns the face of the block the point of a tip segment tapers towards.
func (d PointedDripstone) tipFace() cube.Face {
	if d.Hanging {
		return cube.FaceDown
	}
	return cube.FaceUp
}

// FaceSolid always returns false.
func (PointedDripstone) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
