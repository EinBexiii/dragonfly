package block

// MultiFaceType represents the kind of growth a MultiFace block is.
type MultiFaceType struct {
	multiFace
}

// SculkVeinMultiFace is the vein that spreads over blocks near a sculk catalyst.
func SculkVeinMultiFace() MultiFaceType {
	return MultiFaceType{0}
}

// GlowLichenMultiFace is the lichen that grows in caves and gives off light.
func GlowLichenMultiFace() MultiFaceType {
	return MultiFaceType{1}
}

// ResinClumpMultiFace is the resin that creaking hearts leave on nearby blocks.
func ResinClumpMultiFace() MultiFaceType {
	return MultiFaceType{2}
}

// MultiFaceTypes returns all multi face types.
func MultiFaceTypes() []MultiFaceType {
	return []MultiFaceType{SculkVeinMultiFace(), GlowLichenMultiFace(), ResinClumpMultiFace()}
}

type multiFace uint8

// Uint8 returns the multi face type as a uint8.
func (m multiFace) Uint8() uint8 {
	return uint8(m)
}

// String ...
func (m multiFace) String() string {
	switch m {
	case 0:
		return "sculk_vein"
	case 1:
		return "glow_lichen"
	case 2:
		return "resin_clump"
	}
	panic("unknown multi face type")
}
