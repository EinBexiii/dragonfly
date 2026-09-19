package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
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
