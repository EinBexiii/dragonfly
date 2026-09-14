package block

import (
	"github.com/df-mc/dragonfly/server/world"
)

// RespawnAnchor is a block that lets a player set their spawn point in the nether. It is a full cube at every charge:
// the 2026-09-14 BDS review, batch B section 5, finds it inheriting the unit-box getter and never reading the charge.
type RespawnAnchor struct {
	solid
	bassDrum

	// Charge is the number of glowstone blocks the anchor has been charged with, from 0 to 4.
	Charge int
}

// EncodeBlock ...
func (r RespawnAnchor) EncodeBlock() (string, map[string]any) {
	return "minecraft:respawn_anchor", map[string]any{"respawn_anchor_charge": int32(r.Charge)}
}

// EncodeItem ...
func (r RespawnAnchor) EncodeItem() (name string, meta int16) {
	return blockItemName(r), 0
}

// allRespawnAnchors returns all respawn anchor states.
func allRespawnAnchors() (anchors []world.Block) {
	for charge := range 5 {
		anchors = append(anchors, RespawnAnchor{Charge: charge})
	}
	return
}
