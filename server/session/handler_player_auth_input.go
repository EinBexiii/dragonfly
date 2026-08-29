package session

import (
	"fmt"
	"math"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// PlayerAuthInputHandler handles the PlayerAuthInput packet.
type PlayerAuthInputHandler struct{}

// maxMountDrift is how far a rider's prediction of its mount may run from where
// the server has it before the server sends its own answer. A ride is never
// exact, so a small gap is left alone; anything past it is the prediction
// coming loose.
const maxMountDrift = 1.0

// Handle ...
func (h PlayerAuthInputHandler) Handle(p packet.Packet, s *Session, tx *world.Tx, c Controllable) error {
	pk := p.(*packet.PlayerAuthInput)
	if !h.handleRiding(pk, s, tx, c) {
		if err := h.handleMovement(pk, s, c); err != nil {
			return err
		}
	}
	return h.handleActions(pk, s, tx, c)
}

// handleRiding moves a riding player with its mount and passes on the movement
// it asks of that mount, reporting whether the player rides at all. The rider
// goes wherever the mount is rather than where it reports itself: the two are
// one body while it rides.
//
// A mount the client drives itself reports its position here rather than the
// rider's, in the space the mount is sent in, along with the rotation it turned
// it to. Both are recorded as a request: a client simulates its mount without
// the ground under it, so where the mount really ends up stays the mount's own
// tick to decide. Once the two have drifted apart the server's answer is sent
// on its own, which is what ends the drift; a mount standing its ground
// broadcasts nothing, so without that the client would never hear from it.
func (h PlayerAuthInputHandler) handleRiding(pk *packet.PlayerAuthInput, s *Session, tx *world.Tx, c Controllable) bool {
	m := c.H().Mount()
	if m == nil {
		return false
	}
	mount, ok := m.Entity(tx)
	if !ok {
		return false
	}
	actual := mount.Position()
	if pos, rot, ok := h.predictedMount(pk, s, mount); ok {
		m.Drive(true, pos, rot)
		if actual.Sub(pos).Len() > maxMountDrift {
			s.ViewEntityDisplacement(mount, actual, mount.Rotation(), true)
		}
	} else {
		m.Drive(false, mgl64.Vec3{}, cube.Rotation{})
	}
	s.moving = true
	yaw, pitch := c.Rotation().Elem()
	c.Move(actual.Sub(c.Position()), float64(pk.Yaw)-yaw, float64(pk.Pitch)-pitch)

	if r, ok := mount.(world.Rideable); ok {
		r.Steer(c, world.RideInput{
			Forward: float64(pk.MoveVector[1]),
			Strafe:  float64(pk.MoveVector[0]),
			Jump:    pk.InputData.Load(packet.InputFlagJumping) || pk.InputData.Load(packet.InputFlagStartJumping),
		})
	}
	return true
}

// predictedMount returns where the rider's client predicts the mount it drives,
// and reports whether it drives it at all. The position is the one the mount is
// sent at rather than the rider's, so the mount's own offset comes back off it.
func (h PlayerAuthInputHandler) predictedMount(pk *packet.PlayerAuthInput, s *Session, mount world.Entity) (mgl64.Vec3, cube.Rotation, bool) {
	id, ok := pk.ClientPredictedVehicle.Value()
	if !ok || id == 0 || uint64(id) != s.entityRuntimeID(mount) {
		return mgl64.Vec3{}, cube.Rotation{}, false
	}
	rot, ok := pk.VehicleRotation.Value()
	if !ok {
		// The client has just got on or off and has no rotation to give yet.
		return mgl64.Vec3{}, cube.Rotation{}, false
	}
	pos := vec32To64(pk.Position).Sub(entityOffset(mount))
	if math.IsNaN(pos[0]) || math.IsNaN(pos[1]) || math.IsNaN(pos[2]) {
		return mgl64.Vec3{}, cube.Rotation{}, false
	}
	return pos, cube.Rotation{float64(rot[1]), float64(rot[0])}, true
}

// handleMovement handles the movement part of the packet.PlayerAuthInput.
func (h PlayerAuthInputHandler) handleMovement(pk *packet.PlayerAuthInput, s *Session, c Controllable) error {
	yaw, pitch := c.Rotation().Elem()
	pos := c.Position()

	reference := []float64{pitch, yaw, yaw, pos[0], pos[1], pos[2]}
	for i, v := range [...]*float32{&pk.Pitch, &pk.Yaw, &pk.HeadYaw, &pk.Position[0], &pk.Position[1], &pk.Position[2]} {
		f := float64(*v)
		if math.IsNaN(f) || math.IsInf(f, 1) || math.IsInf(f, 0) {
			// Sometimes, the PlayerAuthInput packet is in fact sent with NaN/INF after being teleported (to another
			// world), see #425. For this reason, we don't actually return an error if this happens, because this will
			// result in the player being kicked. Just log it and replace the NaN value with the one we have tracked
			// server-side.
			s.conf.Log.Debug("process packet: PlayerAuthInput: found nan/inf values. assuming server-side values", "pos", fmt.Sprint(pk.Position), "yaw", pk.Yaw, "head-yaw", pk.HeadYaw, "pitch", pk.Pitch)
			*v = float32(reference[i])
		}
	}

	pk.Position = pk.Position.Sub(mgl32.Vec3{0, 1.62}) // Sub the base offset of players from the pos.

	newPos := vec32To64(pk.Position)
	deltaPos, deltaYaw, deltaPitch := newPos.Sub(pos), float64(pk.Yaw)-yaw, float64(pk.Pitch)-pitch

	// The PlayerAuthInput packet is sent every tick, so don't check for teleport if the position and rotation
	// were unchanged.
	if !mgl64.FloatEqual(deltaPos.Len(), 0) || !mgl64.FloatEqual(deltaYaw, 0) || !mgl64.FloatEqual(deltaPitch, 0) {
		if expected := s.teleportPos.Load(); expected != nil {
			if newPos.Sub(*expected).Len() > 1 {
				// The player has moved before it received the teleport packet. Ignore this movement entirely and
				// wait for the client to sync itself back to the server. Once we get a movement that is close
				// enough to the teleport position, we'll allow the player to move around again.
				return nil
			}
			s.teleportPos.Store(nil)
		}
	}

	s.moving = true
	c.Move(deltaPos, deltaYaw, deltaPitch)
	return nil
}

// handleActions handles the actions with the world that are present in the PlayerAuthInput packet.
func (h PlayerAuthInputHandler) handleActions(pk *packet.PlayerAuthInput, s *Session, tx *world.Tx, c Controllable) error {
	if pk.InputData.Load(packet.InputFlagPerformItemInteraction) {
		data, ok := pk.ItemInteractionData.Value()
		if !ok {
			return fmt.Errorf("item interaction flag set without item interaction data")
		}
		if err := h.handleUseItemData(data, s, c); err != nil {
			return err
		}
	}
	if pk.InputData.Load(packet.InputFlagPerformBlockActions) {
		actions, ok := pk.BlockActions.Value()
		if !ok {
			return fmt.Errorf("block actions flag set without block actions")
		}
		if err := h.handleBlockActions(actions, s, c); err != nil {
			return err
		}
	}
	h.handleInputFlags(pk.InputData, s, c)

	if pk.InputData.Load(packet.InputFlagPerformItemStackRequest) {
		request, ok := pk.ItemStackRequest.Value()
		if !ok {
			return fmt.Errorf("item stack request flag set without item stack request")
		}
		s.inTransaction.Store(true)
		defer s.inTransaction.Store(false)

		// As of 1.18 this is now used for sending item stack requests such as when mining a block.
		sh := s.handlers[packet.IDItemStackRequest].(*ItemStackRequestHandler)
		if err := sh.handleRequest(request, s, tx, c); err != nil {
			// Item stacks being out of sync isn't uncommon, so don't error. Just debug the error and let the
			// revert do its work.
			s.conf.Log.Debug("process packet: PlayerAuthInput: resolve item stack request: " + err.Error())
		}
	}
	return nil
}

// handleInputFlags handles the toggleable input flags set in a PlayerAuthInput packet.
func (h PlayerAuthInputHandler) handleInputFlags(flags protocol.InputFlags, s *Session, c Controllable) {
	if flags.Load(packet.InputFlagStartSprinting) {
		c.StartSprinting()
	}
	if flags.Load(packet.InputFlagStopSprinting) {
		c.StopSprinting()
	}
	if sneaking := flags.Load(packet.InputFlagSneaking); sneaking != c.Sneaking() {
		if sneaking {
			c.StartSneaking()
		} else {
			c.StopSneaking()
		}
	}
	if flags.Load(packet.InputFlagStartSwimming) {
		c.StartSwimming()
	}
	if flags.Load(packet.InputFlagStopSwimming) {
		c.StopSwimming()
	}
	if flags.Load(packet.InputFlagStartGliding) {
		c.StartGliding()
	}
	if flags.Load(packet.InputFlagStopGliding) {
		c.StopGliding()
	}
	if flags.Load(packet.InputFlagStartJumping) {
		c.Jump()
	}
	if flags.Load(packet.InputFlagStartCrawling) {
		c.StartCrawling()
	}
	if flags.Load(packet.InputFlagStopCrawling) {
		c.StopCrawling()
	}
	if flags.Load(packet.InputFlagMissedSwing) {
		s.swingingArm.Store(true)
		defer s.swingingArm.Store(false)
		c.PunchAir()
	}
	if flags.Load(packet.InputFlagStartFlying) {
		if !c.GameMode().AllowsFlying() {
			s.conf.Log.Debug("process packet: PlayerAuthInput: flying flag enabled while unable to fly")
			s.SendAbilities(c)
		} else {
			c.StartFlying()
		}
	}
	if flags.Load(packet.InputFlagStopFlying) {
		c.StopFlying()
	}
}

// handleUseItemData handles the protocol.UseItemTransactionData found in a packet.PlayerAuthInput.
func (h PlayerAuthInputHandler) handleUseItemData(data protocol.UseItemTransactionData, s *Session, c Controllable) error {
	s.swingingArm.Store(true)
	defer s.swingingArm.Store(false)

	held, _ := c.HeldItems()
	if !held.Equal(stackToItem(s.br, data.HeldItem.Stack)) {
		s.conf.Log.Debug("process packet: PlayerAuthInput: UseItemTransaction: mismatch between actual held item and client held item")
		return nil
	}
	pos := cube.Pos{int(data.BlockPosition[0]), int(data.BlockPosition[1]), int(data.BlockPosition[2])}

	// Seems like this is only used for breaking blocks at the moment.
	switch data.ActionType {
	case protocol.UseItemActionBreakBlock:
		c.BreakBlock(pos)
	default:
		return fmt.Errorf("unhandled UseItem ActionType for PlayerAuthInput packet %v", data.ActionType)
	}
	return nil
}

// handleBlockActions handles a slice of protocol.PlayerBlockAction present in a PlayerAuthInput packet.
func (h PlayerAuthInputHandler) handleBlockActions(a []protocol.PlayerBlockAction, s *Session, c Controllable) error {
	for _, action := range a {
		if err := handlePlayerAction(action.Action, action.Face, action.BlockPos, selfEntityRuntimeID, s, c); err != nil {
			return err
		}
	}
	return nil
}
