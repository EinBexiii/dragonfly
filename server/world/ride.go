package world

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/go-gl/mathgl/mgl64"
)

// Seat is a place a rider takes on a Rideable.
type Seat struct {
	// Position is the offset from the Rideable's own position at which the
	// rider sits.
	Position mgl64.Vec3
	// LockRotation keeps a rider facing the way its mount faces, to within
	// MaxAngle degrees either side. A boat locks its riders to 90 degrees.
	LockRotation bool
	MaxAngle     float32
	// Rotate turns a rider Angle degrees from its mount's own facing. A boat
	// seats its riders sideways, at -90.
	Rotate bool
	Angle  float32
}

// Rideable is an Entity that carries riders, such as a horse or a boat.
type Rideable interface {
	Entity
	// Seats returns the seats the Rideable offers, in the order riders take
	// them. The first seat drives. It is read every time a rider takes one, so
	// a Rideable whose seats move with how many it carries may answer from its
	// own Riders.
	Seats() []Seat
}

// Steerable is a Rideable steered by the rider in its driving seat. A Rideable
// that is not Steerable carries its riders without taking their input.
type Steerable interface {
	Rideable
	// Steer passes one tick of the driving rider's input.
	Steer(rider Entity, in RideInput)
}

// Dismounter is a Rideable that answers a rider asking to get off it. One that
// does not implement it lets every rider off that asks.
type Dismounter interface {
	Rideable
	// AllowDismount reports whether rider may get off.
	AllowDismount(rider Entity) bool
}

// RideInput is the movement a driver asks of its mount in one tick.
type RideInput struct {
	// Forward and Strafe are the movement axes, each in the range [-1, 1].
	// Forward is positive towards the direction the rider faces and Strafe is
	// positive to its right.
	Forward, Strafe float64
	// Jump reports whether the rider is holding the jump key.
	Jump bool
}

// NoSeat is the seat of an Entity that rides nothing.
const NoSeat = -1

// ride is an Entity's part in a ride.
type ride struct {
	// mount is the Entity this one rides and seat the seat it rides in, or
	// NoSeat when it rides nothing. offset is where it sits, measured from its
	// mount's position, and lock holds the rotation the seat imposes.
	mount  *EntityHandle
	seat   int
	offset mgl64.Vec3
	lock   Seat
	// riders holds the Entity in each of this one's seats, nil where empty.
	riders []*EntityHandle
	// predicted is where the client driving this Entity reports it, on the
	// tick of its own clock it reported for. It is a request: the server
	// decides where the Entity ends up.
	predicting bool
	predicted  positionState
	tick       uint64
	// corrected holds where a driven Entity really ended up after turning down
	// what its driver reported, waiting to be sent back.
	correcting bool
	corrected  positionState
}

// positionState is a position and rotation, with the ground under it.
type positionState struct {
	pos      mgl64.Vec3
	rot      cube.Rotation
	onGround bool
}

// Mount returns the handle of the Entity that e rides, or nil if e rides
// nothing.
func (e *EntityHandle) Mount() *EntityHandle { return e.data.Ride.mount }

// Seat returns the seat e rides its mount in, or NoSeat if e rides nothing.
func (e *EntityHandle) Seat() int {
	if e.data.Ride.mount == nil {
		return NoSeat
	}
	return e.data.Ride.seat
}

// SeatOffset returns the offset from its mount's position at which e sits, and
// the rotation its seat imposes on it.
func (e *EntityHandle) SeatOffset() (mgl64.Vec3, Seat) {
	return e.data.Ride.offset, e.data.Ride.lock
}

// Rider returns the handle of the Entity riding e in the seat passed, or nil if
// that seat is empty.
func (e *EntityHandle) Rider(seat int) *EntityHandle {
	if seat < 0 || seat >= len(e.data.Ride.riders) {
		return nil
	}
	return e.data.Ride.riders[seat]
}

// Riders returns the handle of the Entity in each of e's seats, nil where a
// seat is empty. The slice is indexed by seat and must not be modified.
func (e *EntityHandle) Riders() []*EntityHandle { return e.data.Ride.riders }

// Driver returns the handle of the Entity in e's driving seat, or nil if
// nobody is driving it.
func (e *EntityHandle) Driver() *EntityHandle { return e.Rider(0) }

// Ridden reports whether any Entity rides e.
func (e *EntityHandle) Ridden() bool {
	for _, r := range e.data.Ride.riders {
		if r != nil {
			return true
		}
	}
	return false
}

// FreeSeat returns the first free seat of mount, reporting false if it has
// none. Seats are counted in the order Seats returns them.
func FreeSeat(mount Rideable) (int, bool) {
	for seat := range mount.Seats() {
		if mount.H().Rider(seat) == nil {
			return seat, true
		}
	}
	return NoSeat, false
}

// NearestSeat returns the free seat of mount closest to pos, reporting false if
// it has none. It is how a rider takes the seat it clicked on.
func NearestSeat(mount Rideable, pos mgl64.Vec3) (int, bool) {
	best, found, dist := NoSeat, false, 0.0
	for seat, s := range mount.Seats() {
		if mount.H().Rider(seat) != nil {
			continue
		}
		if d := mount.Position().Add(s.Position).Sub(pos).LenSqr(); !found || d < dist {
			best, found, dist = seat, true, d
		}
	}
	return best, found
}

// Ride seats rider in the first free seat of mount. It reports false if the
// mount has no seat free or offers none at all.
func (tx *Tx) Ride(rider, mount Entity) bool {
	m, ok := mount.(Rideable)
	if !ok {
		return false
	}
	seat, ok := FreeSeat(m)
	if !ok {
		return false
	}
	return tx.RideSeat(rider, mount, seat)
}

// RideSeat seats rider in the seat of mount passed and shows the two linked to
// every viewer of the mount. It reports false if the mount offers no such seat,
// if that seat is taken by another rider, if the rider is the mount itself, or
// if a Handler cancelled it. A ride the rider is already part of ends first.
func (tx *Tx) RideSeat(rider, mount Entity, seat int) bool {
	r, m := rider.H(), mount.H()
	rideable, ok := mount.(Rideable)
	if !ok || r == m {
		return false
	}
	seats := rideable.Seats()
	if seat < 0 || seat >= len(seats) {
		return false
	}
	if held := m.Rider(seat); held != nil && held != r {
		return false
	}
	ctx := &Context{}
	tx.World().Handler().HandleEntityMount(ctx, rider, mount, &seat)
	if ctx.Cancelled() || seat < 0 || seat >= len(seats) {
		return false
	}
	tx.Dismount(rider)

	for len(m.data.Ride.riders) <= seat {
		m.data.Ride.riders = append(m.data.Ride.riders, nil)
	}
	m.data.Ride.riders[seat] = r
	s := seats[seat]
	// The client seats a rider from the rider's own metadata, so the seat
	// travels with the rider and is measured in the space its position is sent
	// in.
	r.data.Ride.mount, r.data.Ride.seat = m, seat
	r.data.Ride.offset, r.data.Ride.lock = s.Position.Add(riderSeatOffset(rider)), s

	for _, v := range tx.Viewers(mount.Position()) {
		v.ViewEntityMount(rider, mount, true)
		v.ViewEntityState(rider)
	}
	if m.Driver() == r {
		for _, v := range tx.Viewers(mount.Position()) {
			v.ViewEntityState(mount)
		}
	}
	return true
}

// Dismount takes rider off the mount it rides and shows it standing on its own
// again to every viewer of that mount. Dismounting an Entity that rides nothing
// does nothing.
func (tx *Tx) Dismount(rider Entity) {
	r := rider.H()
	m := r.data.Ride.mount
	if m == nil {
		return
	}
	mount, alive := m.Entity(tx)
	if alive {
		ctx := &Context{}
		tx.World().Handler().HandleEntityDismount(ctx, rider, mount)
		if ctx.Cancelled() {
			return
		}
	}
	driving := m.Driver() == r
	if seat := r.data.Ride.seat; seat < len(m.data.Ride.riders) && m.data.Ride.riders[seat] == r {
		m.data.Ride.riders[seat] = nil
	}
	r.data.Ride = ride{seat: NoSeat}
	if !m.Ridden() {
		m.data.Ride.predicting, m.data.Ride.correcting = false, false
	}
	if !alive {
		// The mount left the world before the rider got off it: its viewers
		// dropped the link along with the entity.
		return
	}
	for _, v := range tx.Viewers(mount.Position()) {
		v.ViewEntityMount(rider, mount, false)
		v.ViewEntityState(rider)
		if driving {
			v.ViewEntityState(mount)
		}
	}
}

// breakRideLinks takes e off its mount and every rider off e, so that no side
// is left pointing at an Entity that has left the world.
func (tx *Tx) breakRideLinks(e Entity) {
	tx.Dismount(e)
	for _, h := range e.H().data.Ride.riders {
		if h == nil {
			continue
		}
		if rider, ok := h.Entity(tx); ok {
			tx.Dismount(rider)
			continue
		}
		h.data.Ride = ride{seat: NoSeat}
	}
	clear(e.H().data.Ride.riders)
}

// Predicted returns where the client driving e reports it, and whether one is
// driving it at all. The position is a request: the server still resolves it
// against the blocks around it.
func (e *EntityHandle) Predicted() (mgl64.Vec3, cube.Rotation, bool) {
	r := &e.data.Ride
	return r.predicted.pos, r.predicted.rot, r.predicting
}

// Predict records where the client driving e reports it, on the tick of its own
// clock it reported for.
func (e *EntityHandle) Predict(pos mgl64.Vec3, rot cube.Rotation, tick uint64) {
	r := &e.data.Ride
	r.predicting, r.predicted.pos, r.predicted.rot, r.tick = true, pos, rot, tick
}

// StopPredicting gives e back to its own movement.
func (e *EntityHandle) StopPredicting() {
	e.data.Ride.predicting, e.data.Ride.predicted = false, positionState{}
}

// Refuse records that e turned down what its driver reported and where it
// really ended up, so the driver can be told on its own tick.
func (e *EntityHandle) Refuse(pos mgl64.Vec3, rot cube.Rotation, onGround bool) {
	e.data.Ride.correcting = true
	e.data.Ride.corrected = positionState{pos: pos, rot: rot, onGround: onGround}
}

// Refusal returns what e turned down and the tick it was reported for, and
// clears it. It reports false if e turned nothing down.
func (e *EntityHandle) Refusal() (pos mgl64.Vec3, rot cube.Rotation, onGround bool, tick uint64, ok bool) {
	r := &e.data.Ride
	if !r.correcting {
		return mgl64.Vec3{}, cube.Rotation{}, false, 0, false
	}
	r.correcting = false
	return r.corrected.pos, r.corrected.rot, r.corrected.onGround, r.tick, true
}

// seatedPlayer is how far above the seat its mount names a player rides.
// Captures of BDS give the same figure from two mounts: a horse
// seating a player at 2.32 on a seat of 1.2, and a boat at 1.02 on a seat of
// -0.1.
const seatedPlayer = 1.12

// riderSeatOffset returns what a rider itself adds to the seat its mount names.
// Only a player is measured; what another Entity adds is not, so it sits where
// its mount says.
func riderSeatOffset(rider Entity) mgl64.Vec3 {
	if rider.H().Type().EncodeEntity() != "minecraft:player" {
		return mgl64.Vec3{}
	}
	return mgl64.Vec3{0, seatedPlayer}
}
