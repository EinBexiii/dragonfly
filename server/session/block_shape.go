package session

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// The client used to work out fence connections and stair corners itself. It
// now reads them from the block state, so the session works them out from the
// chunks it has loaded and sends the shaped state in place of the stored one.

// loadedBlocks reads blocks from the chunks the session has loaded. Anything
// outside them is air, and is sent again once the chunk holding it arrives.
type loadedBlocks struct{ s *Session }

func (l loadedBlocks) Block(pos cube.Pos) world.Block {
	col, ok := l.s.chunkLoader.Chunk(world.ChunkPos{int32(pos[0] >> 4), int32(pos[2] >> 4)})
	if r := col.Range(); !ok || pos[1] < r[0] || pos[1] > r[1] {
		return block.Air{}
	}
	return l.s.br.BlockByRuntimeIDOrAir(col.Block(uint8(pos[0]&15), int16(pos[1]), uint8(pos[2]&15), 0))
}

// blockID is the ID the client reads for a block at a position.
func (s *Session) blockID(pos cube.Pos, b world.Block) uint32 {
	if shaped, ok := b.(world.NeighbourShaped); ok {
		return world.NetworkBlockHash(shaped.ShapedState(pos, loadedBlocks{s}))
	}
	return networkBlockID(s.br, b)
}

// shapedBlock reports whether a block's shape depends on its neighbours.
func (s *Session) shapedBlock(rid uint32) bool {
	_, ok := s.br.BlockByRuntimeIDOrAir(rid).(world.NeighbourShaped)
	return ok
}

// subChunkShaper shapes the blocks of the sub-chunk that starts at origin.
type subChunkShaper struct {
	s      *Session
	origin cube.Pos
}

func (sh subChunkShaper) Shaped(rid uint32) bool { return sh.s.shapedBlock(rid) }

func (sh subChunkShaper) Shape(x, y, z byte, rid uint32) uint32 {
	return sh.s.blockID(sh.origin.Add(cube.Pos{int(x), int(y), int(z)}), sh.s.br.BlockByRuntimeIDOrAir(rid))
}

// sendBlock sends the block at pos as the client should see it.
func (s *Session) sendBlock(pos cube.Pos, b world.Block, layer int) {
	s.writePacket(&packet.UpdateBlock{
		Position:          protocol.BlockPos{int32(pos[0]), int32(pos[1]), int32(pos[2])},
		NewBlockRuntimeID: s.blockID(pos, b),
		Flags:             packet.BlockUpdateNetwork,
		Layer:             uint32(layer),
	})
}

// resendShaped sends every shaped block again whose shape a change at pos may
// have moved: two blocks out, since a fence next to a stair follows the corner
// the stair forms with the block beyond it.
func (s *Session) resendShaped(pos cube.Pos) {
	src := loadedBlocks{s}
	for dx := -2; dx <= 2; dx++ {
		for dz := -2; dz <= 2; dz++ {
			if (dx == 0 && dz == 0) || abs(dx)+abs(dz) > 2 {
				continue
			}
			p := pos.Add(cube.Pos{dx, 0, dz})
			if b, ok := src.Block(p).(world.NeighbourShaped); ok {
				s.sendBlock(p, b, 0)
			}
		}
	}
}

// resendBorders sends the shaped blocks in the two columns of each loaded
// neighbour that face the chunk at pos, since they were shaped while pos was
// still missing.
func (s *Session) resendBorders(pos world.ChunkPos) {
	for _, side := range [4][2]int32{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		col, ok := s.chunkLoader.Chunk(world.ChunkPos{pos[0] + side[0], pos[1] + side[1]})
		if !ok {
			continue
		}
		// The two local columns of the neighbour nearest to pos.
		var xs, zs [2]byte
		switch {
		case side[0] == 1:
			xs, zs = [2]byte{0, 1}, [2]byte{0, 15}
		case side[0] == -1:
			xs, zs = [2]byte{15, 14}, [2]byte{0, 15}
		case side[1] == 1:
			xs, zs = [2]byte{0, 15}, [2]byte{0, 1}
		default:
			xs, zs = [2]byte{0, 15}, [2]byte{15, 14}
		}
		base := cube.Pos{int(pos[0]+side[0]) << 4, int(col.Range()[0]), int(pos[1]+side[1]) << 4}
		for ind, sub := range col.Sub() {
			if sub.Empty() || !s.shapedLayer(sub.Layer(0)) {
				continue
			}
			for y := byte(0); y < 16; y++ {
				for _, x := range s.span(xs, side[0] != 0) {
					for _, z := range s.span(zs, side[1] != 0) {
						p := base.Add(cube.Pos{int(x), ind<<4 + int(y), int(z)})
						if b, ok := s.br.BlockByRuntimeIDOrAir(sub.Block(x, y, z, 0)).(world.NeighbourShaped); ok {
							s.sendBlock(p, b, 0)
						}
					}
				}
			}
		}
	}
}

// shapedLayer reports whether a storage holds any shaped block.
func (s *Session) shapedLayer(storage *chunk.PalettedStorage) bool {
	p := storage.Palette()
	for i := 0; i < p.Len(); i++ {
		if s.shapedBlock(p.Value(uint16(i))) {
			return true
		}
	}
	return false
}

// span is the two named columns when the axis crosses the border, or the
// whole edge when it runs along it.
func (*Session) span(v [2]byte, across bool) []byte {
	if across {
		return v[:]
	}
	out := make([]byte, 0, 16)
	for i := v[0]; ; i++ {
		out = append(out, i)
		if i == v[1] {
			return out
		}
	}
}
