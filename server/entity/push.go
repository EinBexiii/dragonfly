package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// pushSpeed is how fast an entity inside a block rises out of it, in blocks per
// tick. Measured against BDS.
const pushSpeed = 0.21

// pushOutOfBlock lifts an entity that ended up inside a block out of it and
// returns the movement doing so, or nil if it is not in one. It replaces the
// entity's physics for the tick, as collision alone would hold the entity at
// the first block it is not yet inside.
func pushOutOfBlock(e *Ent, tx *world.Tx) *Movement {
	pos := e.Position()
	if !obstructed(tx, e.H().Type().BBox(e).Translate(pos)) {
		return nil
	}
	up, prevVel := mgl64.Vec3{0, pushSpeed, 0}, e.data.Vel
	return &Movement{v: tx.Viewers(pos), e: e,
		pos: pos.Add(up), vel: up, dpos: up, dvel: up.Sub(prevVel),
		rot: e.data.Rot, onGround: false,
	}
}

// obstructed reports whether any block intersects the bounding box passed.
func obstructed(tx *world.Tx, box cube.BBox) bool {
	min, max := cube.PosFromVec3(box.Min()), cube.PosFromVec3(box.Max())
	for x := min[0]; x <= max[0]; x++ {
		for y := min[1]; y <= max[1]; y++ {
			for z := min[2]; z <= max[2]; z++ {
				pos := cube.Pos{x, y, z}
				for _, blockBox := range tx.Block(pos).Model().BBox(pos, tx) {
					if blockBox.Translate(pos.Vec3()).IntersectsWith(box) {
						return true
					}
				}
			}
		}
	}
	return false
}
