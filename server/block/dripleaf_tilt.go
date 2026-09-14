package block

// DripleafTilt represents how far the leaf of a BigDripleaf has tipped over.
type DripleafTilt struct {
	dripleafTilt
}

// NoneDripleafTilt is a leaf that is not tilted at all.
func NoneDripleafTilt() DripleafTilt {
	return DripleafTilt{0}
}

// UnstableDripleafTilt is a leaf that has just been stepped on and is about to tilt.
func UnstableDripleafTilt() DripleafTilt {
	return DripleafTilt{1}
}

// PartialDripleafTilt is a leaf that has tilted halfway down.
func PartialDripleafTilt() DripleafTilt {
	return DripleafTilt{2}
}

// FullDripleafTilt is a leaf that has tilted all the way down.
func FullDripleafTilt() DripleafTilt {
	return DripleafTilt{3}
}

// DripleafTilts returns all dripleaf tilts.
func DripleafTilts() []DripleafTilt {
	return []DripleafTilt{NoneDripleafTilt(), UnstableDripleafTilt(), PartialDripleafTilt(), FullDripleafTilt()}
}

type dripleafTilt uint8

// Uint8 returns the dripleaf tilt as a uint8.
func (d dripleafTilt) Uint8() uint8 {
	return uint8(d)
}

// String ...
func (d dripleafTilt) String() string {
	switch d {
	case 0:
		return "none"
	case 1:
		return "unstable"
	case 2:
		return "partial_tilt"
	case 3:
		return "full_tilt"
	}
	panic("unknown dripleaf tilt")
}
