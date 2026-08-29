package session

import (
	"bytes"
	"testing"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// TestRidingMetadataWritable guards the seat fields a riding entity carries
// against a value the protocol cannot write. Entity metadata takes a fixed set
// of Go types and panics inside the packet writer on anything else, which kills
// the connection goroutine rather than failing here.
func TestRidingMetadataWritable(t *testing.T) {
	m := protocol.NewEntityMetadata()
	m.SetFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagRiding)
	writeSeat(m, mgl64.Vec3{0, 2.371, -0.2})

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("riding metadata is not writable: %v", r)
		}
	}()
	w := protocol.NewWriter(bytes.NewBuffer(nil), 0)
	w.EntityMetadata(&m)
}
