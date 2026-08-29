package session

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// InteractHandler handles the packet.Interact.
type InteractHandler struct{}

// Handle ...
func (h *InteractHandler) Handle(p packet.Packet, s *Session, tx *world.Tx, c Controllable) error {
	pk := p.(*packet.Interact)
	pos := c.Position()

	switch pk.ActionType {
	case packet.InteractActionMouseOverEntity:
		// We don't need this action.
	case packet.InteractActionLeaveVehicle:
		// The player sneaked while riding an entity to get off it.
		tx.Dismount(c)
	case packet.InteractActionOpenInventory:
		if m := c.H().Mount(); m != nil {
			// A player riding a mount opens the mount's own inventory rather
			// than its own, the way a saddled horse opens its tack window.
			// Claiming the inventory latch here would strand it: the client
			// closes a window it opened while riding without saying so.
			if mount, ok := m.Entity(tx); ok && s.OpenEntityInventory(mount) {
				return nil
			}
		}
		if s.invOpened.Load() {
			// When there is latency, this might end up being sent multiple times. If we send a ContainerOpen
			// multiple times, the client crashes.
			return nil
		}
		s.invOpened.Store(true)
		s.writePacket(&packet.ContainerOpen{
			WindowID:                0,
			ContainerType:           0xff,
			ContainerEntityUniqueID: -1,
			ContainerPosition: protocol.BlockPos{
				int32(pos[0]),
				int32(pos[1]),
				int32(pos[2]),
			},
		})
	default:
		return fmt.Errorf("unexpected interact packet action %v", pk.ActionType)
	}
	return nil
}
