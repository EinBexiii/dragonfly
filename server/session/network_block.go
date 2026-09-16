package session

import "github.com/df-mc/dragonfly/server/world"

// networkBlockID is the ID a block carries over the network where it has no
// position: the hash of its state, in the shape it has on its own. The client
// resolves it against its own palette, so it means the same block on every
// client version, unlike a palette index. A block at a position is sent with
// Session.blockID, which shapes it by its neighbours.
func networkBlockID(br world.BlockRegistry, b world.Block) uint32 {
	h, _ := br.RuntimeIDToHash(br.BlockRuntimeID(b))
	return h
}

// blockByNetworkID is the inverse of networkBlockID.
func blockByNetworkID(br world.BlockRegistry, id uint32) (world.Block, bool) {
	rid, ok := br.HashToRuntimeID(id)
	if !ok {
		return nil, false
	}
	return br.BlockByRuntimeID(rid)
}
