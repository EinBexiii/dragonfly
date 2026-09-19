package block

// AmethystClusterSize represents how far an AmethystCluster has grown out of the budding amethyst it sits on.
type AmethystClusterSize struct {
	amethystClusterSize
}

// SmallAmethystClusterSize is the first bud to grow out of budding amethyst.
func SmallAmethystClusterSize() AmethystClusterSize {
	return AmethystClusterSize{0}
}

// MediumAmethystClusterSize is the second stage of a growing bud.
func MediumAmethystClusterSize() AmethystClusterSize {
	return AmethystClusterSize{1}
}

// LargeAmethystClusterSize is the third stage of a growing bud.
func LargeAmethystClusterSize() AmethystClusterSize {
	return AmethystClusterSize{2}
}

// FullAmethystClusterSize is a fully grown cluster, which drops amethyst shards when broken.
func FullAmethystClusterSize() AmethystClusterSize {
	return AmethystClusterSize{3}
}

// AmethystClusterSizes returns all amethyst cluster sizes.
func AmethystClusterSizes() []AmethystClusterSize {
	return []AmethystClusterSize{SmallAmethystClusterSize(), MediumAmethystClusterSize(), LargeAmethystClusterSize(), FullAmethystClusterSize()}
}

type amethystClusterSize uint8

// Uint8 returns the amethyst cluster size as a uint8.
func (a amethystClusterSize) Uint8() uint8 {
	return uint8(a)
}

// Height returns how far the cluster sticks out of the block it grows on. This is the H constructor parameter the
// 2026-09-14 BDS review, batch A section 6, reads per size, divided by 16.
func (a amethystClusterSize) Height() float64 {
	switch a {
	case 0:
		return 3.0 / 16
	case 1:
		return 4.0 / 16
	case 2:
		return 5.0 / 16
	case 3:
		return 7.0 / 16
	}
	panic("unknown amethyst cluster size")
}

// Inset returns how far the sides of the cluster are set in from the sides of the block. This is the P constructor
// parameter of batch A section 6, divided by 16.
func (a amethystClusterSize) Inset() float64 {
	if a == 0 {
		return 4.0 / 16
	}
	return 3.0 / 16
}

// String ...
func (a amethystClusterSize) String() string {
	switch a {
	case 0:
		return "small_amethyst_bud"
	case 1:
		return "medium_amethyst_bud"
	case 2:
		return "large_amethyst_bud"
	case 3:
		return "amethyst_cluster"
	}
	panic("unknown amethyst cluster size")
}
