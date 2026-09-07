package session

import (
	"fmt"
	"math"

	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
)

// MovementBroadcastConfig bounds how often a Session tells a client about an
// entity's movement, by how far that entity is from the player. A crowd of
// mobs across the world costs a client a packet each per tick otherwise, and
// the arrow it is watching waits behind them.
//
// Enabled is false by default: with it off a Session sends every movement it
// is offered, as it always has.
type MovementBroadcastConfig struct {
	Enabled bool

	// NearDistance is the distance within which movement is always sent at
	// once. MidDistance and FarDistance bound the slower bands.
	NearDistance, MidDistance, FarDistance float64
	// MidInterval, FarInterval and DistantInterval are how many ticks apart
	// movement is sent in each band.
	MidInterval, FarInterval, DistantInterval uint8
	// CombatTicks is how long an entity stays at full rate after a player
	// attacked it or it targeted them.
	CombatTicks int64
}

func (c MovementBroadcastConfig) normalized() MovementBroadcastConfig {
	if c.NearDistance == 0 {
		c.NearDistance = 16
	}
	if c.MidDistance == 0 {
		c.MidDistance = 32
	}
	if c.FarDistance == 0 {
		c.FarDistance = 64
	}
	if c.MidInterval == 0 {
		c.MidInterval = 2
	}
	if c.FarInterval == 0 {
		c.FarInterval = 4
	}
	if c.DistantInterval == 0 {
		c.DistantInterval = 10
	}
	if c.CombatTicks == 0 {
		c.CombatTicks = 40
	}
	return c
}

// Validate reports whether the configuration describes bands that grow with
// distance and intervals this scheduler can hold.
func (c MovementBroadcastConfig) Validate() error {
	n := c.normalized()
	for _, d := range [3]float64{n.NearDistance, n.MidDistance, n.FarDistance} {
		if math.IsNaN(d) || math.IsInf(d, 0) || d <= 0 {
			return fmt.Errorf("movement broadcast: distances must be positive and finite, got %v", d)
		}
	}
	if !(n.NearDistance < n.MidDistance && n.MidDistance < n.FarDistance) {
		return fmt.Errorf("movement broadcast: distances must grow: %v, %v, %v", n.NearDistance, n.MidDistance, n.FarDistance)
	}
	if !(n.MidInterval <= n.FarInterval && n.FarInterval <= n.DistantInterval) {
		return fmt.Errorf("movement broadcast: intervals must grow: %v, %v, %v", n.MidInterval, n.FarInterval, n.DistantInterval)
	}
	if n.MidInterval < 1 || n.DistantInterval > movementWheel-1 {
		return fmt.Errorf("movement broadcast: intervals must be between 1 and %v", movementWheel-1)
	}
	if n.CombatTicks < 0 {
		return fmt.Errorf("movement broadcast: combat ticks may not be negative")
	}
	return nil
}

// movementWheel is one more than the longest interval, so every pending entity
// sits in exactly one slot of the tick it is due.
const movementWheel = 11

// movementTrack is one entity's movement a Session is holding back.
type movementTrack struct {
	h   *world.EntityHandle
	id  uint64
	due int64
	// update is the newest movement offered since the last one sent; velocity
	// is coalesced rather than replayed.
	update      world.EntityMovementUpdate
	pending     bool
	velPending  bool
	pinnedUntil int64
	lastGround  bool
}

// ViewEntityMovementUpdate takes one tick of an entity's movement and either
// sends it now or holds it for the tick it is due, by how far the entity is
// from this player. Nothing is dropped: the newest state a viewer has not been
// told replaces the one before it and is sent by its deadline.
func (s *Session) ViewEntityMovementUpdate(tx *world.Tx, e world.Entity, u world.EntityMovementUpdate) {
	if !s.conf.MovementBroadcast.Enabled {
		s.sendMovementUpdate(e, u)
		return
	}
	id := s.entityRuntimeID(e)
	if (id == selfEntityRuntimeID && s.moving) || s.entityHidden(e) {
		return
	}
	interval := s.movementInterval(tx, e, u, id)
	if interval <= 1 {
		s.clearTrack(e.H())
		s.sendMovementUpdate(e, u)
		return
	}
	s.entityMutex.Lock()
	t, ok := s.movements[e.H()]
	if !ok {
		t = &movementTrack{h: e.H(), id: id, lastGround: u.OnGround}
		s.movements[e.H()] = t
	}
	// A ground state changing, or a jump, is motion a client cannot guess.
	sendNow := u.OnGround != t.lastGround || !u.OnGround
	t.lastGround = u.OnGround
	if !sendNow {
		if !t.pending {
			t.due = u.Tick + int64(movementPhase(e.H().UUID(), s.uuid(), interval))
			t.pending = true
		}
		t.update, t.velPending = u, t.velPending || u.VelocityChanged
		due := t.due
		s.entityMutex.Unlock()
		if u.Tick < due {
			return
		}
		s.flushTrack(tx, e.H())
		return
	}
	t.pending, t.velPending = false, false
	s.entityMutex.Unlock()
	s.sendMovementUpdate(e, u)
}

// movementInterval is how many ticks apart this entity's movement may be sent
// to this player. One means every tick, which is what everything gets unless
// it is far away, ordinary, and says it may wait.
func (s *Session) movementInterval(tx *world.Tx, e world.Entity, u world.EntityMovementUpdate, id uint64) uint8 {
	c := s.conf.MovementBroadcast
	if id == selfEntityRuntimeID || u.Driven {
		return 1
	}
	h := e.H()
	if h.Mount() != nil || h.Ridden() {
		return 1
	}
	// Only an entity that says its movement may wait ever waits.
	th, ok := e.(interface{ ThrottleMovement() bool })
	if !ok || !th.ThrottleMovement() {
		return 1
	}
	// Motion the client cannot carry on by itself.
	if u.DeltaPosition.Len() > 0.5 || u.Velocity.Len() > 0.5 || u.DeltaVelocity.Len() > 0.2 {
		return 1
	}
	s.entityMutex.RLock()
	t := s.movements[h]
	s.entityMutex.RUnlock()
	if t != nil && u.Tick < t.pinnedUntil {
		return 1
	}
	self, ok := tx.EntityPosition(s.ent)
	if !ok {
		return 1
	}
	d := u.Position.Sub(self)
	d2 := d[0]*d[0] + d[2]*d[2]
	switch {
	case d2 <= c.NearDistance*c.NearDistance:
		return 1
	case d2 <= c.MidDistance*c.MidDistance:
		return c.MidInterval
	case d2 <= c.FarDistance*c.FarDistance:
		return c.FarInterval
	}
	return c.DistantInterval
}

// FlushEntityMovements sends the movement held back for entities due this
// tick, at the end of the world's own tick.
func (s *Session) FlushEntityMovements(tx *world.Tx) {
	if !s.conf.MovementBroadcast.Enabled {
		return
	}
	tick := tx.CurrentTick()
	s.entityMutex.RLock()
	var due []*world.EntityHandle
	for h, t := range s.movements {
		if t.pending && tick >= t.due {
			due = append(due, h)
		}
	}
	s.entityMutex.RUnlock()
	for _, h := range due {
		s.flushTrack(tx, h)
	}
}

// flushTrack sends what is held for one entity, if this world still holds it.
func (s *Session) flushTrack(tx *world.Tx, h *world.EntityHandle) {
	s.entityMutex.Lock()
	t, ok := s.movements[h]
	if !ok || !t.pending {
		s.entityMutex.Unlock()
		return
	}
	u, vel := t.update, t.velPending
	t.pending, t.velPending = false, false
	s.entityMutex.Unlock()
	if _, ok := tx.EntityPosition(h); !ok {
		return
	}
	e, ok := h.Entity(tx)
	if !ok {
		return
	}
	u.VelocityChanged = vel
	s.sendMovementUpdate(e, u)
}

// PinEntityMovement keeps an entity at full rate for the next ticks, for a
// fight or an interaction the player must see exactly.
func (s *Session) PinEntityMovement(tx *world.Tx, h *world.EntityHandle, ticks int64) {
	if !s.conf.MovementBroadcast.Enabled || h == nil {
		return
	}
	s.entityMutex.Lock()
	t, ok := s.movements[h]
	if !ok {
		t = &movementTrack{h: h, id: s.entityRuntimeIDs[h]}
		s.movements[h] = t
	}
	if until := tx.CurrentTick() + ticks; until > t.pinnedUntil {
		t.pinnedUntil = until
	}
	s.entityMutex.Unlock()
	s.flushTrack(tx, h)
}

// clearTrack drops what is held for an entity, its pin kept: the caller is
// about to send its state itself.
func (s *Session) clearTrack(h *world.EntityHandle) {
	s.entityMutex.Lock()
	if t, ok := s.movements[h]; ok {
		t.pending, t.velPending = false, false
	}
	s.entityMutex.Unlock()
}

// forgetMovement drops an entity's schedule outright, for one leaving the
// player's view.
func (s *Session) forgetMovement(h *world.EntityHandle) {
	s.entityMutex.Lock()
	delete(s.movements, h)
	s.entityMutex.Unlock()
}

// sendMovementUpdate is the plain per-tick send: what the viewer would have
// been told before any of this.
func (s *Session) sendMovementUpdate(e world.Entity, u world.EntityMovementUpdate) {
	if u.PositionChanged || u.RotationChanged || u.Driven {
		s.ViewEntityMovement(e, u.Position, u.Rotation, u.OnGround)
	}
	if u.VelocityChanged && !u.Driven {
		s.ViewEntityVelocity(e, u.Velocity)
	}
}

// movementPhase spreads the entities due to one player across the ticks of
// their interval, so a crowd does not arrive in one burst.
func movementPhase(entity, viewer uuid.UUID, interval uint8) uint8 {
	var h uint8
	for i := 0; i < len(entity); i++ {
		h += entity[i] ^ viewer[i]
	}
	return h%interval + 1
}

// uuid returns the player's own identifier, or the zero one before it is set.
func (s *Session) uuid() uuid.UUID {
	if s.ent == nil {
		return uuid.UUID{}
	}
	return s.ent.UUID()
}
