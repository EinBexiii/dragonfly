package world

import "iter"

// entityView groups a World's entity handles by their identifier. Rebuilt
// only when the World's membership changes, it spares a caller looking for
// entities of one kind, players above all, a walk over every handle in the
// chunks around it.
type entityView struct {
	version uint64
	kinds   map[string][]*EntityHandle
}

// invalidateEntityView marks the view stale. Called on every insertion into
// and deletion from w.entities, by the goroutine that owns the World.
func (w *World) invalidateEntityView() {
	w.entitiesVersion++
	w.entityView = nil
}

// entityViewOf returns a view of the current membership, building one if the
// last is stale. A closed World builds a view without keeping it, so work
// deferred to its close cannot repopulate what closing dropped.
func (w *World) entityViewOf(tx *Tx) *entityView {
	if v := w.entityView; v != nil && v.version == w.entitiesVersion {
		return v
	}
	v := &entityView{version: w.entitiesVersion, kinds: map[string][]*EntityHandle{}}
	for h := range tx.EntityHandles() {
		id := h.t.EncodeEntity()
		v.kinds[id] = append(v.kinds[id], h)
	}
	if !w.closed.Load() {
		w.entityView = v
	}
	return v
}

// EntityHandlesOf yields the handles of the entities whose identifier is one
// of ids, in the order the identifiers are given, without opening them. It
// reads a view of the World's membership rebuilt only when that membership
// changes, so a crowd of mobs each looking for the few players costs the
// players, not the crowd.
//
// A handle removed from the World since the view was built is skipped, and
// one added during the iteration appears in the next query rather than this
// one. Positions are not part of the view: read them through EntityPosition.
func (tx *Tx) EntityHandlesOf(ids ...string) iter.Seq[*EntityHandle] {
	w := tx.World()
	return func(yield func(*EntityHandle) bool) {
		if len(ids) == 0 {
			return
		}
		v := w.entityViewOf(tx)
		for i, id := range ids {
			// An identifier named twice yields its handles once.
			dup := false
			for _, seen := range ids[:i] {
				if seen == id {
					dup = true
					break
				}
			}
			if dup {
				continue
			}
			for _, h := range v.kinds[id] {
				if _, here := w.entities[h]; !here {
					continue
				}
				if !yield(h) {
					return
				}
			}
		}
	}
}
