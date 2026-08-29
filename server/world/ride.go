package world

import "github.com/go-gl/mathgl/mgl64"

// Rideable is an Entity that is steered by the Entity riding it, such as a
// saddled horse. An Entity that does not implement Rideable can still be
// ridden, but carries its rider along without taking its input.
type Rideable interface {
	Entity
	// Steer passes the movement input of the rider for a single tick. It is
	// called only while rider rides the Rideable.
	Steer(rider Entity, in RideInput)
}

// Seater is an Entity that seats its rider somewhere other than at its own
// position, such as a horse seating a player on its back. An Entity that does
// not implement Seater seats its rider at its feet.
type Seater interface {
	Entity
	// SeatPosition returns the offset from the Entity's own position at which
	// its rider sits.
	SeatPosition() mgl64.Vec3
}

// RideInput is the movement a rider asks of its mount in one tick.
type RideInput struct {
	// Forward and Strafe are the movement axes, each in the range [-1, 1].
	// Forward is positive towards the direction the rider faces and Strafe is
	// positive to its right.
	Forward, Strafe float64
	// Jump reports whether the rider is holding the jump key.
	Jump bool
}

// Mount returns the handle of the Entity that e rides, or nil if e rides
// nothing.
func (e *EntityHandle) Mount() *EntityHandle {
	return e.data.Mount
}

// Rider returns the handle of the Entity riding e, or nil if e carries no
// rider.
func (e *EntityHandle) Rider() *EntityHandle {
	return e.data.Rider
}

// Driven reports whether a rider is steering e, and where that rider's client
// predicts it. The position is where the rider wants e to be: the server still
// resolves it against the blocks around it, and a server that took it as given
// would follow a client straight through the floor.
func (e *EntityHandle) Driven() (mgl64.Vec3, bool) {
	return e.data.DrivenPos, e.data.Driven
}

// Drive records where the rider steering e predicts it. Passing driven false
// gives e back to its own movement.
func (e *EntityHandle) Drive(driven bool, pos mgl64.Vec3) {
	e.data.Driven, e.data.DrivenPos = driven, pos
}

// SeatPosition returns the offset from its mount's position at which e sits.
// It is the zero vector for an Entity that rides nothing.
func (e *EntityHandle) SeatPosition() mgl64.Vec3 {
	return e.data.Seat
}

// Ride seats rider on mount and shows the two linked to every viewer of the
// mount. A rider rides at most one mount and a mount carries at most one
// rider, so a link either of them is part of is broken first. Riding the mount
// already ridden, or riding oneself, does nothing.
func (tx *Tx) Ride(rider, mount Entity) {
	r, m := rider.H(), mount.H()
	if r == m || r.data.Mount == m {
		return
	}
	tx.Dismount(rider)
	if old := m.data.Rider; old != nil {
		if e, ok := old.Entity(tx); ok {
			tx.Dismount(e)
		} else {
			old.data.Mount, m.data.Rider = nil, nil
		}
	}
	r.data.Mount, m.data.Rider = m, r
	if seater, ok := mount.(Seater); ok {
		// The client seats the rider from the rider's own metadata, so the
		// seat travels with it rather than with the mount, and it is measured
		// in the space the rider's own position is sent in.
		r.data.Seat = seater.SeatPosition().Add(riderSeatOffset(rider))
	}
	for _, v := range tx.Viewers(mount.Position()) {
		v.ViewEntityMount(rider, mount, true)
		v.ViewEntityState(rider)
	}
}

// Dismount takes rider off the mount it rides and shows it standing on its own
// again to every viewer of that mount. Dismounting an Entity that rides
// nothing does nothing.
func (tx *Tx) Dismount(rider Entity) {
	r := rider.H()
	m := r.data.Mount
	if m == nil {
		return
	}
	r.data.Mount, m.data.Rider, r.data.Seat = nil, nil, mgl64.Vec3{}
	m.data.Driven, m.data.DrivenPos = false, mgl64.Vec3{}
	mount, ok := m.Entity(tx)
	if !ok {
		// The mount left the world before the rider got off it: there is no
		// position left to broadcast from, and its viewers dropped the link
		// along with the entity itself.
		return
	}
	for _, v := range tx.Viewers(mount.Position()) {
		v.ViewEntityMount(rider, mount, false)
		v.ViewEntityState(rider)
	}
}

// breakRideLinks takes e off its mount and its rider off e, so that neither
// side is left pointing at an Entity that has left the world.
func (tx *Tx) breakRideLinks(e Entity) {
	tx.Dismount(e)
	if r := e.H().data.Rider; r != nil {
		if rider, ok := r.Entity(tx); ok {
			tx.Dismount(rider)
		} else {
			r.data.Mount, e.H().data.Rider = nil, nil
		}
	}
}

// seatedFactor is the share of its standing eye height a seated rider sits at.
// A rider that does not sit upright drops by the remainder of it instead.
const seatedFactor = 0.75308642

// riderSeatOffset returns what a rider adds to the seat its mount names. A
// mount names only where its saddle sits; how high the rider itself rides on it
// follows from the rider's own size, so the two are added.
func riderSeatOffset(rider Entity) mgl64.Vec3 {
	h := rider.H().Type().BBox(rider).Height()
	if rider.H().Type().EncodeEntity() == "minecraft:player" {
		// A player sits at a share of its standing eye height, which is nine
		// tenths of its own height.
		return mgl64.Vec3{0, h * 0.9 * seatedFactor}
	}
	return mgl64.Vec3{0, -h * (1 - seatedFactor)}
}
