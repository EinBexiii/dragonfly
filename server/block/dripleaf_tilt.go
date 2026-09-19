package block

import "github.com/df-mc/dragonfly/server/block/model"

// DripleafTilt represents how far the leaf of a BigDripleaf has tipped over. Its ordinals are Bedrock's, and
// model.DripleafTilt holds the one definition of them: the constructors below derive from it, so the conversion at
// the Model() call site cannot drift.
type DripleafTilt struct {
	dripleafTilt
}

// NoneDripleafTilt is a leaf that is not tilted at all.
func NoneDripleafTilt() DripleafTilt {
	return DripleafTilt{dripleafTilt(model.DripleafTiltNone)}
}

// UnstableDripleafTilt is a leaf that has just been stepped on and is about to tilt.
func UnstableDripleafTilt() DripleafTilt {
	return DripleafTilt{dripleafTilt(model.DripleafTiltUnstable)}
}

// PartialDripleafTilt is a leaf that has tilted halfway down.
func PartialDripleafTilt() DripleafTilt {
	return DripleafTilt{dripleafTilt(model.DripleafTiltPartial)}
}

// FullDripleafTilt is a leaf that has tilted all the way down.
func FullDripleafTilt() DripleafTilt {
	return DripleafTilt{dripleafTilt(model.DripleafTiltFull)}
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
	switch model.DripleafTilt(d) {
	case model.DripleafTiltNone:
		return "none"
	case model.DripleafTiltUnstable:
		return "unstable"
	case model.DripleafTiltPartial:
		return "partial_tilt"
	case model.DripleafTiltFull:
		return "full_tilt"
	}
	panic("unknown dripleaf tilt")
}
