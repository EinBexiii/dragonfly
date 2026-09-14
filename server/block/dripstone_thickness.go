package block

// DripstoneThickness represents the segment a PointedDripstone block forms in the spike it is part of.
type DripstoneThickness struct {
	dripstoneThickness
}

// TipDripstoneThickness is the pointed end of a spike.
func TipDripstoneThickness() DripstoneThickness {
	return DripstoneThickness{0}
}

// FrustumDripstoneThickness is the segment right behind the tip.
func FrustumDripstoneThickness() DripstoneThickness {
	return DripstoneThickness{1}
}

// MiddleDripstoneThickness is a segment in the middle of a spike.
func MiddleDripstoneThickness() DripstoneThickness {
	return DripstoneThickness{2}
}

// BaseDripstoneThickness is the segment a spike grows from.
func BaseDripstoneThickness() DripstoneThickness {
	return DripstoneThickness{3}
}

// MergeDripstoneThickness is the tip of a spike that has grown into a spike coming from the other side.
func MergeDripstoneThickness() DripstoneThickness {
	return DripstoneThickness{4}
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
	switch d {
	case 0:
		return "tip"
	case 1:
		return "frustum"
	case 2:
		return "middle"
	case 3:
		return "base"
	case 4:
		return "merge"
	}
	panic("unknown dripstone thickness")
}
