package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
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

// UseOnBlock grows the cluster out of the face clicked, so it points away from the block carrying it.
func (a AmethystCluster) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, face, used := firstReplaceable(tx, pos, face, a)
	if !used {
		return false
	}
	a.Facing = face
	if !a.canSurvive(pos, tx) {
		return false
	}

	place(tx, pos, a, user, ctx)
	return placed(ctx)
}

// NeighbourUpdateTick ...
func (a AmethystCluster) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if !a.canSurvive(pos, tx) {
		breakBlockNoDrops(a, pos, tx)
	}
}

// canSurvive checks if the block the cluster grows out of is still there.
func (a AmethystCluster) canSurvive(pos cube.Pos, tx *world.Tx) bool {
	return faceSolid(tx, pos.Side(a.Facing.Opposite()), a.Facing)
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
