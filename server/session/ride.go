package session

import (
	"math"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// ViewEntityMount views one Entity taking a seat on another, or leaving it
// again when mounted is false.
func (s *Session) ViewEntityMount(rider, mount world.Entity, mounted bool) {
	if s.entityHidden(rider) || s.entityHidden(mount) {
		return
	}
	link, ok := s.entityLink(rider.H(), mount.H())
	if !ok {
		return
	}
	if !mounted {
		link.Type = protocol.EntityLinkRemove
	}
	if link.RiderEntityUniqueID == selfEntityRuntimeID {
		// Taking a seat or leaving one drops whatever window the client had
		// open without telling the server. Left claimed, the player's own
		// inventory never opens again.
		s.invOpened = false
	}
	s.writePacket(&packet.SetActorLink{EntityLink: link})
}

// entityLink builds the link seating rider on mount. It reports false if the
// session does not have both Entities in view.
func (s *Session) entityLink(rider, mount *world.EntityHandle) (protocol.EntityLink, bool) {
	s.entityMutex.RLock()
	riderID, riderOK := s.entityRuntimeIDs[rider]
	mountID, mountOK := s.entityRuntimeIDs[mount]
	s.entityMutex.RUnlock()
	if !riderOK || !mountOK {
		return protocol.EntityLink{}, false
	}
	// The driving seat is a rider link and every other seat a passenger one.
	linkType := byte(protocol.EntityLinkPassenger)
	if rider.Seat() == 0 {
		linkType = protocol.EntityLinkRider
	}
	return protocol.EntityLink{
		RiddenEntityUniqueID: int64(mountID),
		RiderEntityUniqueID:  int64(riderID),
		Type:                 linkType,
		RiderInitiated:       riderID == selfEntityRuntimeID,
	}, true
}

// entityLinks returns the links seating e on the Entity it rides and every
// rider of its own on it, for the packet that spawns e. A link is only carried
// by the spawn packet of the side that appears last: the client drops one
// naming an Entity it does not know yet.
func (s *Session) entityLinks(e world.Entity) []protocol.EntityLink {
	var links []protocol.EntityLink
	h := e.H()
	if m := h.Mount(); m != nil {
		if link, ok := s.entityLink(h, m); ok {
			links = append(links, link)
		}
	}
	for _, r := range h.Riders() {
		if r == nil {
			continue
		}
		if link, ok := s.entityLink(r, h); ok {
			links = append(links, link)
		}
	}
	return links
}

// writeSeat writes where a rider sits and how its seat turns it. The client
// reads all of it from the rider rather than from its mount.
func writeSeat(m protocol.EntityMetadata, offset mgl64.Vec3, seat world.Seat) {
	m[protocol.EntityDataKeySeatOffset] = vec64To32(offset)
	m[protocol.EntityDataKeySeatLockPassengerRotation] = boolByte(seat.LockRotation)
	m[protocol.EntityDataKeySeatLockPassengerRotationDegrees] = seat.MaxAngle
	m[protocol.EntityDataKeySeatRotationOffset] = boolByte(seat.Rotate)
	m[protocol.EntityDataKeySeatRotationOffsetDegrees] = seat.Angle
}

// handleRiding moves a riding player with its mount and passes on the movement
// it asks of that mount, reporting whether the player rides at all. The rider
// goes where the mount is rather than where it reports itself.
//
// A driver's client runs the mount itself and reports the mount's position here
// rather than its own. That is taken as a request: whether the mount may go
// there is the mount's own tick to answer.
func (h PlayerAuthInputHandler) handleRiding(pk *packet.PlayerAuthInput, s *Session, tx *world.Tx, c Controllable) bool {
	m := c.H().Mount()
	if m == nil {
		return false
	}
	mount, ok := m.Entity(tx)
	if !ok {
		return false
	}
	driving := m.Driver() == c.H()
	if pos, rot, ok := h.predictedMount(pk, s, mount); ok {
		m.Predict(pos, rot, pk.Tick)
	} else if _, _, predicting := m.Predicted(); predicting && driving {
		// A driving client leaves the prediction out of the odd packet. The
		// mount stays its client's to move and is simply not moved this tick:
		// handing it back to its own movement sets it walking off on its own.
		m.Predict(mount.Position(), mount.Rotation(), pk.Tick)
	} else {
		m.StopPredicting()
	}
	h.correctMount(s, mount)

	s.moving = true
	yaw, pitch := c.Rotation().Elem()
	c.Move(mount.Position().Sub(c.Position()), float64(pk.Yaw)-yaw, float64(pk.Pitch)-pitch)

	if steerable, ok := mount.(world.Steerable); ok && driving {
		steerable.Steer(c, world.RideInput{
			Forward: float64(pk.MoveVector[1]),
			Strafe:  float64(pk.MoveVector[0]),
			Jump:    pk.InputData.Load(packet.InputFlagJumping) || pk.InputData.Load(packet.InputFlagStartJumping),
		})
	}
	return true
}

// correctMount tells a driver's client where its mount really went, for a
// request the mount turned down. The client re-runs its own prediction from the
// position given for the tick named.
func (h PlayerAuthInputHandler) correctMount(s *Session, mount world.Entity) {
	pos, rot, onGround, tick, ok := mount.H().Refusal()
	if !ok {
		return
	}
	s.writePacket(&packet.CorrectPlayerMovePrediction{
		PredictionType: packet.PredictionTypeVehicle,
		Position:       vec64To32(pos.Add(entityOffset(mount))),
		Rotation:       mgl32.Vec2{float32(rot.Pitch()), float32(rot.Yaw())},
		OnGround:       onGround,
		Tick:           tick,
	})
}

// predictedMount returns where a driver's client reports the mount it runs, and
// whether it is running one at all. The position is the mount's rather than the
// rider's, so the mount's own offset comes back off it.
func (h PlayerAuthInputHandler) predictedMount(pk *packet.PlayerAuthInput, s *Session, mount world.Entity) (mgl64.Vec3, cube.Rotation, bool) {
	id, ok := pk.ClientPredictedVehicle.Value()
	if !ok || id == 0 || uint64(id) != s.entityRuntimeID(mount) {
		return mgl64.Vec3{}, cube.Rotation{}, false
	}
	rot, ok := pk.VehicleRotation.Value()
	if !ok {
		return mgl64.Vec3{}, cube.Rotation{}, false
	}
	pos := vec32To64(pk.Position).Sub(entityOffset(mount))
	if math.IsNaN(pos[0]) || math.IsNaN(pos[1]) || math.IsNaN(pos[2]) {
		return mgl64.Vec3{}, cube.Rotation{}, false
	}
	return pos, cube.Rotation{float64(rot[1]), float64(rot[0])}, true
}
