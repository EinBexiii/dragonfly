package block

import "github.com/df-mc/dragonfly/server/block/model"

// DripstoneThickness represents the segment a PointedDripstone block forms in the spike it is part of. Its ordinals
// are Bedrock's, and model.DripstoneThickness holds the one definition of them: the constructors below derive from
// it, so the conversion at the Model() call site cannot drift.
type DripstoneThickness struct {
	dripstoneThickness
}

// TipDripstoneThickness is the pointed end of a spike.
func TipDripstoneThickness() DripstoneThickness {
	return DripstoneThickness{dripstoneThickness(model.DripstoneTip)}
}

// FrustumDripstoneThickness is the segment right behind the tip.
func FrustumDripstoneThickness() DripstoneThickness {
	return DripstoneThickness{dripstoneThickness(model.DripstoneFrustum)}
}

// MiddleDripstoneThickness is a segment in the middle of a spike.
func MiddleDripstoneThickness() DripstoneThickness {
	return DripstoneThickness{dripstoneThickness(model.DripstoneMiddle)}
}

// BaseDripstoneThickness is the segment a spike grows from.
func BaseDripstoneThickness() DripstoneThickness {
	return DripstoneThickness{dripstoneThickness(model.DripstoneBase)}
}

// MergeDripstoneThickness is the tip of a spike that has grown into a spike coming from the other side.
func MergeDripstoneThickness() DripstoneThickness {
	return DripstoneThickness{dripstoneThickness(model.DripstoneMerge)}
}

// DripstoneThicknesses returns all dripstone thicknesses.
func DripstoneThicknesses() []DripstoneThickness {
	return []DripstoneThickness{TipDripstoneThickness(), FrustumDripstoneThickness(), MiddleDripstoneThickness(), BaseDripstoneThickness(), MergeDripstoneThickness()}
}

type dripstoneThickness uint8

// Uint8 returns the dripstone thickness as a uint8.
func (d dripstoneThickness) Uint8() uint8 {
	return uint8(d)
}

// String ...
func (d dripstoneThickness) String() string {
	switch model.DripstoneThickness(d) {
	case model.DripstoneTip:
		return "tip"
	case model.DripstoneFrustum:
		return "frustum"
	case model.DripstoneMiddle:
		return "middle"
	case model.DripstoneBase:
		return "base"
	case model.DripstoneMerge:
		return "merge"
	}
	panic("unknown dripstone thickness")
}
