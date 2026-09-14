package block

import "github.com/df-mc/dragonfly/server/block/cube"

// Attachment describes the attachment of a block to another block. It is either of the type WallAttachment, which can
// only have 90 degree facing values, or StandingAttachment, which has more freedom using a cube.Orientation.
type Attachment struct {
	hanging bool
	facing  cube.Direction
	o       cube.Orientation
}

// WallAttachment returns an Attachment to a wall with a facing direction.
func WallAttachment(facing cube.Direction) Attachment {
	return Attachment{hanging: true, facing: facing}
}

// StandingAttachment returns an Attachment to the ground with an orientation.
func StandingAttachment(o cube.Orientation) Attachment {
	return Attachment{o: o}
}

// supportFace returns the face of the block the Attachment is attached to. The second return value is false
// for a wall attachment whose facing has no horizontal equivalent, which is attached to nothing.
func (a Attachment) supportFace() (cube.Face, bool) {
	if !a.hanging {
		return cube.FaceDown, true
	}
	if !horizontal(a.facing) {
		return 0, false
	}
	return a.facing.Opposite().Face(), true
}

// Uint8 returns the Attachment as a uint8.
func (a Attachment) Uint8() uint8 {
	if !a.hanging {
		return 1 | (uint8(a.o) << 1)
	}
	return uint8(a.facing) << 1
}

// FaceUint8 returns the facing of the Attachment as a uint8.
func (a Attachment) FaceUint8() uint8 {
	if !a.hanging {
		return 1
	}
	return uint8(a.facing) << 1
}

// RotateLeft rotates the Attachment the left way around by 90 degrees.
func (a Attachment) RotateLeft() Attachment {
	return Attachment{hanging: a.hanging, facing: a.facing.RotateLeft(), o: a.o.RotateLeft()}
}

// RotateRight rotates the Attachment the right way around by 90 degrees.
func (a Attachment) RotateRight() Attachment {
	return Attachment{hanging: a.hanging, facing: a.facing.RotateLeft(), o: a.o.RotateLeft()}
}

// Rotation returns the rotation of the Attachment, based on the orientation if it's a StandingAttachment, or the
// facing direction if it is a WallAttachment.
func (a Attachment) Rotation() cube.Rotation {
	yaw := a.o.Yaw()
	if a.hanging {
		switch a.facing {
		case cube.West:
			yaw = 90
		case cube.East:
			yaw = -90
		case cube.North, unknownDirection, unknownUpDirection:
			// facing_direction 0 and 1 point down and up, which a wall attachment has no rotation for. The
			// game draws such a block as if it faced north, so give it the same rotation.
			yaw = 180
		}
	}
	return cube.Rotation{yaw}
}
