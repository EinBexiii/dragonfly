package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

// AmethystCluster is a bud growing out of budding amethyst, in one of four stages of growth.
type AmethystCluster struct {
	transparent

	// Size is how far the cluster has grown.
	Size AmethystClusterSize
	// Facing is the face of the neighbouring block the cluster grows on, so the cluster itself points the other way.
	Facing cube.Face
}

// Model ...
func (a AmethystCluster) Model() world.BlockModel {
	return model.AmethystCluster{Height: a.Size.Height(), Inset: a.Size.Inset(), Facing: a.Facing}
}

// EncodeBlock ...
func (a AmethystCluster) EncodeBlock() (string, map[string]any) {
	// The getter reads minecraft:block_face rather than the legacy facing_direction, verified in batch A section 6.
	return "minecraft:" + a.Size.String(), map[string]any{"minecraft:block_face": a.Facing.String()}
}

// EncodeItem ...
func (a AmethystCluster) EncodeItem() (name string, meta int16) {
	return blockItemName(a), 0
}

// allAmethystClusters returns all amethyst bud and cluster states.
func allAmethystClusters() (clusters []world.Block) {
	for _, size := range AmethystClusterSizes() {
		for _, f := range cube.Faces() {
			clusters = append(clusters, AmethystCluster{Size: size, Facing: f})
		}
	}
	return
}
