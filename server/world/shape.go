package world

import "github.com/df-mc/dragonfly/server/block/cube"

// NeighbourShaped is a block whose look the client cannot work out on its own,
// such as how a fence connects or which corner a stair forms. The server sends
// that shape as part of the block's state, worked out from the neighbours at
// the moment the block goes out; the world itself does not store it.
type NeighbourShaped interface {
	Block
	// ShapedState is EncodeBlock plus the properties the shape adds.
	ShapedState(pos cube.Pos, src BlockSource) (name string, properties map[string]any)
}

// NetworkBlockHash is the network hash of a block state, registered or not.
func NetworkBlockHash(name string, properties map[string]any) uint32 {
	h, _ := networkBlockHash(name, properties, nil)
	return h
}
