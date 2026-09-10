package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// Wall is a model used by all wall types.
type Wall struct {
	// NorthConnection is the height of the connection for the north direction.
	NorthConnection float64
	// EastConnection is the height of the connection for the east direction.
	EastConnection float64
	// SouthConnection is the height of the connection for the south direction.
	SouthConnection float64
	// WestConnection is the height of the connection for the west direction.
	WestConnection float64
	// Post is if the wall is the full height of a block or not.
	Post bool
}

const (
	// wallPostInset is the post's edge from the block's edge.
	wallPostInset = 0.25
	// wallArmInset is an arm's edge from the block's edge, across the arm.
	wallArmInset = 0.3125
)

// BBox returns the wall's collision box, one box the way the Bedrock server
// computes it, BarrierHeight tall whatever the post and connection heights
// say about the drawn shape: a straight run without a post is as narrow as
// its arms, anything else is the post grown along each arm.
func (w Wall) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	north, east, south, west := w.NorthConnection > 0, w.EastConnection > 0, w.SouthConnection > 0, w.WestConnection > 0
	if !w.Post {
		switch {
		case north && south && !east && !west:
			return []cube.BBox{cube.Box(wallArmInset, 0, 0, 1-wallArmInset, BarrierHeight, 1)}
		case east && west && !north && !south:
			return []cube.BBox{cube.Box(0, 0, wallArmInset, 1, BarrierHeight, 1-wallArmInset)}
		}
	}
	box := cube.Box(wallPostInset, 0, wallPostInset, 1-wallPostInset, BarrierHeight, 1-wallPostInset)
	for face, connected := range map[cube.Face]bool{cube.FaceNorth: north, cube.FaceEast: east, cube.FaceSouth: south, cube.FaceWest: west} {
		if connected {
			box = box.ExtendTowards(face, wallPostInset)
		}
	}
	return []cube.BBox{box}
}

// FaceSolid returns true if the face is in the Y axis.
func (w Wall) FaceSolid(_ cube.Pos, face cube.Face, _ world.BlockSource) bool {
	return face.Axis() == cube.Y
}
