package block

import (
	"maps"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

// connected adds the four connection flags the client reads to properties,
// asking connects for each horizontal face.
func connected(properties map[string]any, connects func(face cube.Face) bool) map[string]any {
	out := make(map[string]any, len(properties)+4)
	maps.Copy(out, properties)
	for _, face := range cube.HorizontalFaces() {
		out["minecraft:connection_"+face.String()] = boolByte(connects(face))
	}
	return out
}

// fenceState is the shaped state of a block with a fence model.
func fenceState(b world.Block, pos cube.Pos, src world.BlockSource) (string, map[string]any) {
	name, properties := b.EncodeBlock()
	m := b.Model().(model.Fence)
	return name, connected(properties, func(face cube.Face) bool { return m.Connects(pos, face, src) })
}

// thinState is the shaped state of a block with a thin model.
func thinState(b world.Block, pos cube.Pos, src world.BlockSource) (string, map[string]any) {
	name, properties := b.EncodeBlock()
	m := b.Model().(model.Thin)
	return name, connected(properties, func(face cube.Face) bool { return m.Connects(pos, face, src) })
}

// ShapedState ...
func (w WoodFence) ShapedState(pos cube.Pos, src world.BlockSource) (string, map[string]any) {
	return fenceState(w, pos, src)
}

// ShapedState ...
func (n NetherBrickFence) ShapedState(pos cube.Pos, src world.BlockSource) (string, map[string]any) {
	return fenceState(n, pos, src)
}

// ShapedState ...
func (g GlassPane) ShapedState(pos cube.Pos, src world.BlockSource) (string, map[string]any) {
	return thinState(g, pos, src)
}

// ShapedState ...
func (p StainedGlassPane) ShapedState(pos cube.Pos, src world.BlockSource) (string, map[string]any) {
	return thinState(p, pos, src)
}

// ShapedState ...
func (i IronBars) ShapedState(pos cube.Pos, src world.BlockSource) (string, map[string]any) {
	return thinState(i, pos, src)
}

// ShapedState ...
func (c CopperBars) ShapedState(pos cube.Pos, src world.BlockSource) (string, map[string]any) {
	return thinState(c, pos, src)
}

// ShapedState joins tripwire to the tripwire next to it.
func (s String) ShapedState(pos cube.Pos, src world.BlockSource) (string, map[string]any) {
	name, properties := s.EncodeBlock()
	return name, connected(properties, func(face cube.Face) bool {
		_, ok := src.Block(pos.Side(face)).(String)
		return ok
	})
}

// ShapedState adds the corner the stairs form.
func (s Stairs) ShapedState(pos cube.Pos, src world.BlockSource) (string, map[string]any) {
	name, properties := s.EncodeBlock()
	properties["minecraft:corner"] = s.Model().(model.Stair).Corner(pos, src).String()
	return name, properties
}
