package block

import (
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

// TurtleEgg is a clutch of one to four eggs that hatches into turtles.
type TurtleEgg struct {
	transparent

	// Count is the number of eggs in the clutch, from 1 to 4.
	Count int
	// Cracks is how far the eggs have cracked open.
	Cracks EggCracks
}

// Model ...
func (t TurtleEgg) Model() world.BlockModel {
	return model.TurtleEgg{}
}

// EncodeBlock ...
func (t TurtleEgg) EncodeBlock() (string, map[string]any) {
	return "minecraft:turtle_egg", map[string]any{
		"turtle_egg_count": turtleEggCounts[t.Count-1],
		"cracked_state":    t.Cracks.String(),
	}
}

// EncodeItem ...
func (t TurtleEgg) EncodeItem() (name string, meta int16) {
	return blockItemName(t), 0
}

// turtleEggCounts holds the turtle_egg_count state value for a clutch of one to four eggs.
var turtleEggCounts = [...]string{"one_egg", "two_egg", "three_egg", "four_egg"}

// SnifferEgg is an egg that hatches into a sniffer.
type SnifferEgg struct {
	transparent

	// Cracks is how far the egg has cracked open.
	Cracks EggCracks
}

// Model ...
func (s SnifferEgg) Model() world.BlockModel {
	return model.SnifferEgg{}
}

// EncodeBlock ...
func (s SnifferEgg) EncodeBlock() (string, map[string]any) {
	return "minecraft:sniffer_egg", map[string]any{"cracked_state": s.Cracks.String()}
}

// EncodeItem ...
func (s SnifferEgg) EncodeItem() (name string, meta int16) {
	return blockItemName(s), 0
}

// allEggs returns all turtle egg and sniffer egg states.
func allEggs() (eggs []world.Block) {
	for _, c := range AllEggCracks() {
		eggs = append(eggs, SnifferEgg{Cracks: c})
		for count := 1; count <= len(turtleEggCounts); count++ {
			eggs = append(eggs, TurtleEgg{Count: count, Cracks: c})
		}
	}
	return
}
