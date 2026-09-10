package model

import (
	"testing"

	"github.com/df-mc/dragonfly/server/block/cube"
)

// A wall is a block and a half tall to an entity, whatever it looks like:
// the collision box is what keeps a player from jumping over it and what a
// player standing on it stands on.
func TestWallCollisionIsABlockAndAHalf(t *testing.T) {
	for _, w := range []Wall{{}, {Post: true}, {NorthConnection: 0.8125, SouthConnection: 1}} {
		for _, box := range w.BBox(cube.Pos{}, nil) {
			if box.Max().Y() != wallCollisionHeight {
				t.Fatalf("%+v: box %v is %v tall", w, box, box.Max().Y())
			}
		}
	}
}
