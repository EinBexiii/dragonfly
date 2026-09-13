package block

import (
	"github.com/df-mc/dragonfly/server/world"
)

// FrostedIce is the ice that Frost Walker boots leave behind. It melts away in vanilla as its age rises, but
// the lobby keeps it as permanent decoration, so no melting is implemented.
type FrostedIce struct {
	solid

	// Age is how far the ice has melted, 0-3. 3 is about to disappear.
	Age int
}

// BreakInfo ...
func (FrostedIce) BreakInfo() BreakInfo {
	return newBreakInfo(0.5, pickaxeHarvestable, pickaxeEffective, simpleDrops())
}

// EncodeBlock ...
func (f FrostedIce) EncodeBlock() (string, map[string]any) {
	return "minecraft:frosted_ice", map[string]any{"age": int32(f.Age)}
}

// allFrostedIce ...
func allFrostedIce() (ice []world.Block) {
	for age := 0; age < 4; age++ {
		ice = append(ice, FrostedIce{Age: age})
	}
	return
}
