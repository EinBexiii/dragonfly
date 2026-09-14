package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// HangingSign is a sign that hangs from a chain below a block or from a bar attached to the side of one.
type HangingSign struct {
	transparent
	blockEntityData

	// Wood is the type of wood of the sign. This field must have one of the values found in the material package.
	Wood WoodType
	// Facing is the face of the block the sign is attached by.
	Facing cube.Face
	// Attached specifies if the chain of the sign is attached to the block above rather than to a bar.
	Attached bool
	// Hanging specifies if the sign hangs below a block instead of being attached to the side of one.
	Hanging bool
	// GroundDirection is the rotation of a sign hanging below a block, ranging from 0 to 15.
	GroundDirection int
}

// Model ...
func (s HangingSign) Model() world.BlockModel {
	return model.HangingSign{Facing: s.Facing, Hanging: s.Hanging}
}

// UseOnBlock hangs the sign below the block clicked underneath, turned the sixteenth of a circle the player faces,
// or bolts it to the wall whose side was clicked. A sign cannot stand on top of a block.
func (s HangingSign) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, face, used := firstReplaceable(tx, pos, face, s)
	if !used || face == cube.FaceUp {
		return false
	}
	if face == cube.FaceDown {
		s.Hanging = true
		s.GroundDirection = int(user.Rotation().Orientation().Opposite()) % 16
	} else {
		s.Facing = wallSignFacing(face, user)
	}
	if !s.canSurvive(pos, tx) {
		return false
	}

	place(tx, pos, s, user, ctx)
	return placed(ctx)
}

// NeighbourUpdateTick ...
func (s HangingSign) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if !s.canSurvive(pos, tx) {
		breakBlockNoDrops(s, pos, tx)
	}
}

// wallSignFacing returns the face a sign put against the wall clicked shows. The bar it hangs from runs along the
// wall, so the sign faces across it, turned to the player of the two ways it can.
func wallSignFacing(wall cube.Face, user item.User) cube.Face {
	facing := user.Rotation().Direction().Opposite().Face()
	if facing.Axis() == wall.Axis() {
		// The player is looking along the bar rather than at the sign, so the wall alone decides.
		return wall.RotateRight()
	}
	return facing
}

// canSurvive checks if the block the sign hangs from is still there. A sign below a block hangs from that block or
// from another sign; a sign on a wall hangs from a bar whose ends rest in the blocks beside it.
func (s HangingSign) canSurvive(pos cube.Pos, tx *world.Tx) bool {
	if s.Hanging {
		above := pos.Side(cube.FaceUp)
		if _, ok := tx.Block(above).(HangingSign); ok {
			return true
		}
		return faceSolid(tx, above, cube.FaceDown)
	}
	for _, end := range []cube.Face{s.Facing.RotateRight(), s.Facing.RotateLeft()} {
		if faceSolid(tx, pos.Side(end), end.Opposite()) {
			return true
		}
	}
	return false
}

// EncodeBlock ...
func (s HangingSign) EncodeBlock() (string, map[string]any) {
	return "minecraft:" + s.Wood.String() + "_hanging_sign", map[string]any{
		"attached_bit":          s.Attached,
		"facing_direction":      int32(s.Facing),
		"ground_sign_direction": int32(s.GroundDirection),
		"hanging":               s.Hanging,
	}
}

// EncodeItem ...
func (s HangingSign) EncodeItem() (name string, meta int16) {
	return blockItemName(s), 0
}

// EncodeNBT ...
func (s HangingSign) EncodeNBT() map[string]any {
	return s.storedNBT("HangingSign")
}

// DecodeNBT ...
func (s HangingSign) DecodeNBT(data map[string]any) any {
	s.blockEntityData = s.storeNBT(data)
	return s
}

// allHangingSigns returns all hanging sign states.
func allHangingSigns() (signs []world.Block) {
	for _, w := range WoodTypes() {
		for _, f := range cube.Faces() {
			for _, attached := range []bool{false, true} {
				for _, hanging := range []bool{false, true} {
					for rotation := range 16 {
						signs = append(signs, HangingSign{Wood: w, Facing: f, Attached: attached, Hanging: hanging, GroundDirection: rotation})
					}
				}
			}
		}
	}
	return
}
