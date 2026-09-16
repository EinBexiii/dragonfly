package session

import (
	"maps"

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
	if !ok {
		return block.Air{}
	}
	if r := col.Range(); pos[1] < r[0] || pos[1] > r[1] {
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
	s.writeBlock(pos, s.blockID(pos, b), layer)
}

// writeBlock sends a block ID for pos.
func (s *Session) writeBlock(pos cube.Pos, id uint32, layer int) {
	s.writePacket(&packet.UpdateBlock{
		Position:          protocol.BlockPos{int32(pos[0]), int32(pos[1]), int32(pos[2])},
		NewBlockRuntimeID: id,
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

// chunkless reads blocks as they were sent before the chunk at pos arrived:
// air there, the loaded chunks elsewhere. That is what the client holds, with
// one exception: a neighbour that left the loader and changed while it was
// away is compared as it is now, so a block shaped by that change waits for
// the next update near it.
type chunkless struct {
	src world.BlockSource
	pos world.ChunkPos
}

func (c chunkless) Block(p cube.Pos) world.Block {
	if (world.ChunkPos{int32(p[0] >> 4), int32(p[2] >> 4)}) == c.pos {
		return block.Air{}
	}
	return c.src.Block(p)
}

// resendBorders sends the shaped blocks of each loaded neighbour again whose
// shape the chunk at pos changes: those within two of it were shaped while
// pos was still missing, as if it held air, and most of them look the same
// either way. Diagonal neighbours count: a fence follows a stair whose
// corner follows a block in the chunk beyond.
func (s *Session) resendBorders(pos world.ChunkPos) {
	now := loadedBlocks{s}
	before := chunkless{now, pos}
	for dx := int32(-1); dx <= 1; dx++ {
		for dz := int32(-1); dz <= 1; dz++ {
			if dx == 0 && dz == 0 {
				continue
			}
			at := world.ChunkPos{pos[0] + dx, pos[1] + dz}
			col, ok := s.chunkLoader.Chunk(at)
			if !ok {
				continue
			}
			xs, zs := edge(-dx), edge(-dz)
			base := cube.Pos{int(at[0]) << 4, int(col.Range()[0]), int(at[1]) << 4}
			for ind, sub := range col.Sub() {
				if sub.Empty() || !s.shapedLayer(sub.Layer(0)) {
					continue
				}
				for y := byte(0); y < 16; y++ {
					for _, x := range xs {
						for _, z := range zs {
							b, ok := s.br.BlockByRuntimeIDOrAir(sub.Block(x, y, z, 0)).(world.NeighbourShaped)
							if !ok {
								continue
							}
							p := base.Add(cube.Pos{int(x), ind<<4 + int(y), int(z)})
							name, properties := b.ShapedState(p, now)
							if _, sent := b.ShapedState(p, before); !maps.Equal(properties, sent) {
								s.writeBlock(p, world.NetworkBlockHash(name, properties), 0)
							}
						}
					}
				}
			}
		}
	}
}

// edge is the local coordinates of a neighbour chunk nearest to the chunk on
// the given side: the two rows on that side, or the whole chunk when the side
// is not along this axis.
func edge(side int32) []byte {
	switch side {
	case 1:
		return []byte{15, 14}
	case -1:
		return []byte{0, 1}
	}
	return whole[:]
}

var whole = [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

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

// encodeChunk encodes a whole chunk for the network with its blocks shaped.
func (s *Session) encodeChunk(pos world.ChunkPos, c *chunk.Chunk) chunk.SerialisedData {
	d := chunk.SerialisedData{SubChunks: make([][]byte, len(c.Sub()))}
	for i := range d.SubChunks {
		origin := cube.Pos{int(pos[0]) << 4, int(c.Range()[0]) + i<<4, int(pos[1]) << 4}
		d.SubChunks[i] = chunk.EncodeSubChunkShaped(c, chunk.NetworkEncoding, i, subChunkShaper{s: s, origin: origin})
	}
	d.Biomes = chunk.EncodeBiomes(c, chunk.NetworkEncoding)
	return d
}
