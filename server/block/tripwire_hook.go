package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// TripwireHook is the block a tripwire is strung between. It is decoration here: no redstone logic is run.
type TripwireHook struct {
	transparent
	empty

	// Facing is the direction the hook points away from the wall it hangs on.
	Facing cube.Direction
	// Attached is whether a tripwire is strung to the hook.
	Attached bool
	// Powered is whether the tripwire attached to the hook is being stepped on.
	Powered bool
}

// UseOnBlock hangs the hook on the wall face clicked, pointing away from it.
func (t TripwireHook) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, face, used := firstReplaceable(tx, pos, face, t)
	if !used || face.Axis() == cube.Y {
		return false
	}
	t.Facing = face.Direction()
	if !t.supported(pos, tx) {
		return false
	}
	place(tx, pos, t, user, ctx)
	return placed(ctx)
}

// NeighbourUpdateTick drops the hook when its wall goes.
func (t TripwireHook) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if !t.supported(pos, tx) {
		breakBlock(t, pos, tx)
	}
}

// supported reports whether the block behind the hook has a solid face to hang on.
func (t TripwireHook) supported(pos cube.Pos, tx *world.Tx) bool {
	wall := pos.Side(t.Facing.Face().Opposite())
	return tx.Block(wall).Model().FaceSolid(wall, t.Facing.Face(), tx)
}

// BreakInfo ...
func (t TripwireHook) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, oneOf(t))
}

// EncodeItem ...
func (TripwireHook) EncodeItem() (name string, meta int16) {
	return "minecraft:tripwire_hook", 0
}

// EncodeBlock ...
func (t TripwireHook) EncodeBlock() (string, map[string]any) {
	return "minecraft:tripwire_hook", map[string]any{
		"direction":    int32(horizontalDirection(t.Facing)),
		"attached_bit": t.Attached,
		"powered_bit":  t.Powered,
	}
}

// allTripwireHooks ...
func allTripwireHooks() (hooks []world.Block) {
	for _, d := range cube.Directions() {
		hooks = append(hooks, TripwireHook{Facing: d})
		hooks = append(hooks, TripwireHook{Facing: d, Attached: true})
		hooks = append(hooks, TripwireHook{Facing: d, Powered: true})
		hooks = append(hooks, TripwireHook{Facing: d, Attached: true, Powered: true})
	}
	return
}
