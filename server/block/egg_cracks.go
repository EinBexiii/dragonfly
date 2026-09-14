package block

// EggCracks represents how far an egg block has cracked open.
type EggCracks struct {
	eggCracks
}

// NoEggCracks is an egg that has not started hatching yet.
func NoEggCracks() EggCracks {
	return EggCracks{0}
}

// SomeEggCracks is an egg that has started hatching.
func SomeEggCracks() EggCracks {
	return EggCracks{1}
}

// MaxEggCracks is an egg that is about to hatch.
func MaxEggCracks() EggCracks {
	return EggCracks{2}
}

// AllEggCracks returns all egg crack stages.
func AllEggCracks() []EggCracks {
	return []EggCracks{NoEggCracks(), SomeEggCracks(), MaxEggCracks()}
}

type eggCracks uint8

// Uint8 returns the egg cracks as a uint8.
func (e eggCracks) Uint8() uint8 {
	return uint8(e)
}

// String ...
func (e eggCracks) String() string {
	switch e {
	case 0:
		return "no_cracks"
	case 1:
		return "cracked"
	case 2:
		return "max_cracked"
	}
	panic("unknown egg cracks")
}
