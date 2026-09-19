package particle

// Custom is a particle identified by name. It may be used to show particle
// effects defined by a resource pack.
type Custom struct {
	particle

	// Name is the identifier of the particle effect.
	Name string
}
