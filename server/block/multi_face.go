package block

import (
	"github.com/df-mc/dragonfly/server/world"
)

// MultiFace is a growth that clings to any number of the six faces of the blocks around it. Every kind of it has no
// collision at all: the 2026-09-14 BDS review, batch B section 8, finds all three installing the empty collision
// getter, with the face-dependent outlines used for selection only.
type MultiFace struct {
	empty
	transparent

	// Type is the kind of growth this block is.
	Type MultiFaceType
	// Faces is the bit set of block faces the growth clings to, from 0 to 63.
	Faces int
}

// EncodeBlock ...
func (m MultiFace) EncodeBlock() (string, map[string]any) {
	return "minecraft:" + m.Type.String(), map[string]any{"multi_face_direction_bits": int32(m.Faces)}
}

// allMultiFace returns all sculk vein, glow lichen and resin clump states.
func allMultiFace() (growths []world.Block) {
	for _, t := range MultiFaceTypes() {
		for faces := range 64 {
			growths = append(growths, MultiFace{Type: t, Faces: faces})
		}
	}
	return
}
