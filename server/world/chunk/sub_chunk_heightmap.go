package chunk

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// SubChunkHeightMaps prepares a Chunk's surface once for multiple sub-chunk
// entries. Recreate it after editing the Chunk.
type SubChunkHeightMaps struct {
	c               *Chunk
	heights         HeightMap
	lowest, highest int16
}

// NewSubChunkHeightMaps prepares the surface summary for c.
func NewSubChunkHeightMaps(c *Chunk) SubChunkHeightMaps {
	heights := c.HeightMap()
	lowest := c.SubIndex(heights[0])
	highest := lowest
	for _, y := range heights[1:] {
		index := c.SubIndex(y)
		lowest, highest = min(lowest, index), max(highest, index)
	}
	return SubChunkHeightMaps{c: c, heights: heights, lowest: lowest, highest: highest}
}

// At returns the height-map type for index and, when the surface crosses
// that sub-chunk, the height of every column relative to it.
func (m SubChunkHeightMaps) At(index int16) (byte, protocol.Optional[protocol.HeightMap]) {
	switch {
	case index < m.lowest:
		return protocol.HeightMapDataTooHigh, protocol.Optional[protocol.HeightMap]{}
	case index > m.highest:
		return protocol.HeightMapDataTooLow, protocol.Optional[protocol.HeightMap]{}
	}
	var heights protocol.HeightMap
	for z := uint8(0); z < 16; z++ {
		for x := uint8(0); x < 16; x++ {
			y := m.heights.At(x, z)
			columnIndex := m.c.SubIndex(y)
			switch {
			case columnIndex > index:
				heights[z][x] = 16
			case columnIndex < index:
				heights[z][x] = -1
			default:
				heights[z][x] = int8(y - m.c.SubY(columnIndex))
			}
		}
	}
	return protocol.HeightMapDataHasData, protocol.Option(heights)
}
