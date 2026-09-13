package block

import (
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

// DaylightDetector is a block that emits a redstone signal proportional to the daylight reaching it. Its
// signal is not recalculated here: the saved value is kept as it is.
type DaylightDetector struct {
	transparent

	// Inverted specifies if the detector emits a stronger signal the darker it is instead.
	Inverted bool
	// Signal is the redstone signal the detector currently emits, 0-15.
	Signal int
}

// Model ...
func (DaylightDetector) Model() world.BlockModel {
	return model.DaylightDetector{}
}

// BreakInfo ...
func (DaylightDetector) BreakInfo() BreakInfo {
	return newBreakInfo(0.2, alwaysHarvestable, axeEffective, simpleDrops(item.NewStack(DaylightDetector{}, 1)))
}

// EncodeItem ...
func (DaylightDetector) EncodeItem() (name string, meta int16) {
	return "minecraft:daylight_detector", 0
}

// EncodeBlock ...
func (d DaylightDetector) EncodeBlock() (string, map[string]any) {
	properties := map[string]any{"redstone_signal": int32(d.Signal)}
	if d.Inverted {
		return "minecraft:daylight_detector_inverted", properties
	}
	return "minecraft:daylight_detector", properties
}

// allDaylightDetectors ...
func allDaylightDetectors() (detectors []world.Block) {
	for signal := 0; signal < 16; signal++ {
		detectors = append(detectors, DaylightDetector{Signal: signal})
		detectors = append(detectors, DaylightDetector{Inverted: true, Signal: signal})
	}
	return
}
