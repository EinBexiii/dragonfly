package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
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

// Face checks if the growth clings to the face of its own block passed.
func (m MultiFace) Face(f cube.Face) bool {
	return m.Faces&(1<<f) != 0
}

// WithFace returns the growth with the face passed added to or removed from the faces it clings to.
func (m MultiFace) WithFace(f cube.Face, clinging bool) MultiFace {
	if clinging {
		m.Faces |= 1 << f
	} else {
		m.Faces &^= 1 << f
	}
	return m
}

// UseOnBlock makes the growth cling to the face that was clicked, adding that face to a growth of the same kind
// already there rather than replacing it.
func (m MultiFace) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	if !faceSolid(tx, pos, face) {
		return false
	}
	// A growth covers the face of the block carrying it that points at the growth: the opposite of the face clicked.
	clings, growthPos := face.Opposite(), pos.Side(face)
	if existing, ok := tx.Block(growthPos).(MultiFace); ok && existing.Type == m.Type {
		if existing.Face(clings) {
			return false
		}
		place(tx, growthPos, existing.WithFace(clings, true), user, ctx)
		return placed(ctx)
	}
	if !replaceableWith(tx, growthPos, m) {
		return false
	}

	place(tx, growthPos, m.WithFace(clings, true), user, ctx)
	return placed(ctx)
}

// NeighbourUpdateTick drops the faces of the growth whose block is gone, and the growth itself with its last face.
func (m MultiFace) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	updated := m
	for _, f := range cube.Faces() {
		if updated.Face(f) && !faceSolid(tx, pos.Side(f), f.Opposite()) {
			updated = updated.WithFace(f, false)
		}
	}
	if updated.Faces == m.Faces {
		return
	}
	if updated.Faces == 0 {
		breakBlockNoDrops(m, pos, tx)
		return
	}
	tx.SetBlock(pos, updated, nil)
}

// EncodeBlock ...
func (m MultiFace) EncodeBlock() (string, map[string]any) {
	return "minecraft:" + m.Type.String(), map[string]any{"multi_face_direction_bits": int32(m.Faces)}
}

// EncodeItem ...
func (m MultiFace) EncodeItem() (name string, meta int16) {
	return blockItemName(m), 0
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
