package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

// StainedGlassPane is a transparent block that can be used as a more efficient alternative to glass blocks.
type StainedGlassPane struct {
	transparent
	thin
	clicksAndSticks
	sourceWaterDisplacer

	// Colour specifies the colour of the block.
	Colour item.Colour
	// Connections holds the sides that the glass pane connects to.
	Connections Connections
	// Hardened specifies if the pane is the hardened variant, which is only obtainable through commands.
	Hardened bool
}

// NeighbourUpdateTick ...
func (p StainedGlassPane) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if connections := calculateConnections(p.Model().(connector), tx, pos); connections != p.Connections {
		p.Connections = connections
		tx.SetBlock(pos, p, nil)
	}
}

// DeriveState ...
func (p StainedGlassPane) DeriveState(pos cube.Pos, src world.BlockSource) world.Block {
	p.Connections = calculateConnections(p.Model().(connector), src, pos)
	return p
}

// SideClosed ...
func (p StainedGlassPane) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

// BreakInfo ...
func (p StainedGlassPane) BreakInfo() BreakInfo {
	return newBreakInfo(0.3, alwaysHarvestable, nothingEffective, silkTouchOnlyDrop(p))
}

// EncodeItem ...
func (p StainedGlassPane) EncodeItem() (name string, meta int16) {
	if p.Hardened {
		return "minecraft:hard_" + p.Colour.String() + "_stained_glass_pane", 0
	}
	return "minecraft:" + p.Colour.String() + "_stained_glass_pane", 0
}

// EncodeBlock ...
func (p StainedGlassPane) EncodeBlock() (name string, properties map[string]any) {
	if p.Hardened {
		return "minecraft:hard_" + p.Colour.String() + "_stained_glass_pane", p.Connections.properties()
	}
	return "minecraft:" + p.Colour.String() + "_stained_glass_pane", p.Connections.properties()
}

// allStainedGlassPane returns stained-glass panes with all possible colours.
func allStainedGlassPane() []world.Block {
	b := make([]world.Block, 0, 512)
	for _, colour := range item.Colours() {
		for _, c := range allConnections() {
			b = append(b, StainedGlassPane{Colour: colour, Connections: c})
			b = append(b, StainedGlassPane{Colour: colour, Connections: c, Hardened: true})
		}
	}
	return b
}
