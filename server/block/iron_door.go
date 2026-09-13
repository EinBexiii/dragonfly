package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// IronDoor is a 1x2 barrier a player cannot open by hand. The game opens it by redstone; this fork has no
// redstone hook for it yet, so it stays as the world stores it.
type IronDoor struct {
	transparent
	bass
	sourceWaterDisplacer

	// Facing is the direction that the door opens towards. When closed, the door sits on the side of its
	// block on the opposite direction.
	Facing cube.Direction
	// Open is whether the door is open.
	Open bool
	// Top is whether the block is the top or bottom half of a door.
	Top bool
	// Right is whether the door hinge is on the right side.
	Right bool
}

// Model ...
func (d IronDoor) Model() world.BlockModel {
	return model.Door{Facing: d.Facing, Open: d.Open, Right: d.Right, Top: d.Top, Kind: "iron"}
}

// NeighbourUpdateTick ...
func (d IronDoor) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if d.Top {
		if _, ok := tx.Block(pos.Side(cube.FaceDown)).(IronDoor); !ok {
			breakBlockNoDrops(d, pos, tx)
		}
	} else if solid := tx.Block(pos.Side(cube.FaceDown)).Model().FaceSolid(pos.Side(cube.FaceDown), cube.FaceUp, tx); !solid {
		breakBlock(d, pos, tx)
	} else if _, ok := tx.Block(pos.Side(cube.FaceUp)).(IronDoor); !ok {
		breakBlockNoDrops(d, pos, tx)
	}
}

// UseOnBlock handles the directional placing of doors
func (d IronDoor) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	return placeDoor(pos, face, tx, user, ctx, d, func(facing cube.Direction, right bool) (lower, upper world.Block) {
		return IronDoor{Facing: facing, Right: right}, IronDoor{Facing: facing, Right: right, Top: true}
	})
}

// BreakInfo ...
func (d IronDoor) BreakInfo() BreakInfo {
	return newBreakInfo(5, pickaxeHarvestable, pickaxeEffective, oneOf(d))
}

// SideClosed ...
func (d IronDoor) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

// EncodeItem ...
func (d IronDoor) EncodeItem() (name string, meta int16) {
	return "minecraft:iron_door", 0
}

// EncodeBlock ...
func (d IronDoor) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:iron_door", map[string]any{"minecraft:cardinal_direction": d.Facing.RotateRight().String(), "door_hinge_bit": d.Right, "open_bit": d.Open, "upper_block_bit": d.Top}
}

// allIronDoors returns every state of the iron door.
func allIronDoors() (doors []world.Block) {
	for i := cube.Direction(0); i <= 3; i++ {
		for _, open := range []bool{false, true} {
			for _, top := range []bool{false, true} {
				for _, right := range []bool{false, true} {
					doors = append(doors, IronDoor{Facing: i, Open: open, Top: top, Right: right})
				}
			}
		}
	}
	return
}
