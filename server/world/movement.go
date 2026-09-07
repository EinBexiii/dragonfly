package world

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/go-gl/mathgl/mgl64"
)

// EntityMovementUpdate is one entity's movement as it happened, offered whole
// to a viewer that schedules what it sends rather than sending every tick.
type EntityMovementUpdate struct {
	Tick int64

	Position, Velocity mgl64.Vec3
	Rotation           cube.Rotation
	OnGround           bool

	// DeltaPosition and DeltaVelocity are the change this tick, which decides
	// whether the motion is one a viewer may be told about later.
	DeltaPosition, DeltaVelocity mgl64.Vec3

	PositionChanged bool
	RotationChanged bool
	VelocityChanged bool
	// Driven is set for an entity a client predicts itself, which is told
	// where it is every tick and its velocity never.
	Driven bool
}

// MovementViewer is a Viewer that takes an entity's movement whole and decides
// itself when to pass it on. A Viewer without it keeps the plain per-tick
// movement and velocity calls.
type MovementViewer interface {
	Viewer
	// ViewEntityMovementUpdate offers one tick of an entity's movement.
	ViewEntityMovementUpdate(tx *Tx, e Entity, update EntityMovementUpdate)
	// FlushEntityMovements is called at the end of a world tick, so movement
	// held back reaches the viewer even once the entity stops moving.
	FlushEntityMovements(tx *Tx)
}
