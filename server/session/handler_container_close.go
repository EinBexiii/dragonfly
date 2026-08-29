package session

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// ContainerCloseHandler handles the ContainerClose packet.
type ContainerCloseHandler struct{}

// Handle ...
func (h *ContainerCloseHandler) Handle(p packet.Packet, s *Session, tx *world.Tx, c Controllable) error {
	pk := p.(*packet.ContainerClose)

	c.MoveItemsToInventory()

	// Whatever the client closed, its screen is free again, so the latch that
	// keeps a second inventory from being opened over the first is released
	// here rather than only for the inventory's own window. A window opened
	// without a ContainerOpen of its own, such as a trade or a mount's, is
	// closed with an ID the inventory case never sees, and leaving the latch
	// claimed shuts the player out of its own inventory for good.
	s.invOpened.Store(false)

	// A window the client holds no ID for is closed with 0xff, and the answer
	// carries the ID the server knows it by, which is what a real Bedrock
	// server replies with.
	windowID := pk.WindowID

	var containerType byte
	switch pk.WindowID {
	case 0:
		// Closing of the normal inventory.
	case byte(s.openedWindowID.Load()):
		containerType = byte(s.openedContainerID.Load())
		s.closeCurrentContainer(tx, true)
	case 0xff:
		// The client closes a window it holds no ID for with 0xff, which is
		// how it closes one opened alongside chat and how it closes a trade.
		// It is a close it asked for like any other, so it is answered with
		// the ID it named: answering about a window it does not know leaves it
		// waiting for its own close for good, and it stops asking to open
		// anything at all after that.
		containerType = pk.ContainerType
		if s.containerOpened.Load() {
			containerType, windowID = byte(s.openedContainerID.Load()), byte(s.openedWindowID.Load())
			s.closeCurrentContainer(tx, true)
		}
	default:
		containerType = pk.ContainerType
	}
	s.writePacket(&packet.ContainerClose{
		WindowID:      windowID,
		ContainerType: containerType,
	})
	return nil
}
