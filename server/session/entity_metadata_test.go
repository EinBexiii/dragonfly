package session

import (
	"bytes"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// TestRidingMetadataWritable guards the seat fields a riding entity carries
// against a value the protocol cannot write. Entity metadata takes a fixed set
// of Go types and panics inside the packet writer on anything else, which kills
// the connection goroutine rather than failing here.
func TestRidingMetadataWritable(t *testing.T) {
	m := protocol.NewEntityMetadata()
	m.SetFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagRiding)
	m[protocol.EntityDataKeySeatOffset] = mgl32.Vec3{0, 2.471, 0}
	m[protocol.EntityDataKeySeatLockPassengerRotation] = byte(0)
	m[protocol.EntityDataKeySeatLockPassengerRotationDegrees] = float32(0)
	m[protocol.EntityDataKeySeatRotationOffset] = float32(0)
	m[protocol.EntityDataKeySeatRotationOffsetDegrees] = float32(0)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("riding metadata is not writable: %v", r)
		}
	}()
	w := protocol.NewWriter(bytes.NewBuffer(nil), 0)
	w.EntityMetadata(&m)
}
