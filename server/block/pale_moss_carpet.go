package block

import (
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

// PaleMossCarpet is a carpet of pale moss that climbs up the blocks around it. Only its lower half collides, with a
// layer 1/16 high: the 2026-09-14 BDS review, batch B section 9, has the native getter return that layer for
// upper_block_bit false regardless of the side heights, and nothing at all for true.
type PaleMossCarpet struct {
	transparent

	// UpperPart specifies if this block is the part of the carpet climbing the sides of a block above it.
	UpperPart bool
	// North, East, South and West are how far the carpet climbs the block on each of its sides.
	North, East, South, West MossCarpetSide
}

// Model ...
func (p PaleMossCarpet) Model() world.BlockModel {
	if p.UpperPart {
		return model.Empty{}
	}
	return model.Carpet{}
}

// EncodeBlock ...
func (p PaleMossCarpet) EncodeBlock() (string, map[string]any) {
	return "minecraft:pale_moss_carpet", map[string]any{
		"upper_block_bit":             p.UpperPart,
		"pale_moss_carpet_side_north": p.North.String(),
		"pale_moss_carpet_side_east":  p.East.String(),
		"pale_moss_carpet_side_south": p.South.String(),
		"pale_moss_carpet_side_west":  p.West.String(),
	}
}

// allPaleMossCarpets returns all pale moss carpet states.
func allPaleMossCarpets() (carpets []world.Block) {
	for _, upper := range []bool{false, true} {
		for _, north := range MossCarpetSides() {
			for _, east := range MossCarpetSides() {
				for _, south := range MossCarpetSides() {
					for _, west := range MossCarpetSides() {
						carpets = append(carpets, PaleMossCarpet{UpperPart: upper, North: north, East: east, South: south, West: west})
					}
				}
			}
		}
	}
	return
}
