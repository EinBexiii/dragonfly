package block

// MossCarpetSide represents how far a PaleMossCarpet climbs up the block on one of its sides.
type MossCarpetSide struct {
	mossCarpetSide
}

// NoMossCarpetSide is a side the carpet does not climb at all.
func NoMossCarpetSide() MossCarpetSide {
	return MossCarpetSide{0}
}

// ShortMossCarpetSide is a side the carpet climbs halfway up.
func ShortMossCarpetSide() MossCarpetSide {
	return MossCarpetSide{1}
}

// TallMossCarpetSide is a side the carpet climbs all the way up.
func TallMossCarpetSide() MossCarpetSide {
	return MossCarpetSide{2}
}

// MossCarpetSides returns all moss carpet side heights.
func MossCarpetSides() []MossCarpetSide {
	return []MossCarpetSide{NoMossCarpetSide(), ShortMossCarpetSide(), TallMossCarpetSide()}
}

type mossCarpetSide uint8

// Uint8 returns the moss carpet side as a uint8.
func (m mossCarpetSide) Uint8() uint8 {
	return uint8(m)
}

// String ...
func (m mossCarpetSide) String() string {
	switch m {
	case 0:
		return "none"
	case 1:
		return "short"
	case 2:
		return "tall"
	}
	panic("unknown moss carpet side")
}
