package block

// FrogSpawn is the clutch of eggs a frog lays on the surface of water. It has no collision: the 2026-09-14 BDS review,
// batch B section 11, finds the thin box of its constructor used as an outline only.
type FrogSpawn struct {
	empty
	transparent
}

// EncodeBlock ...
func (FrogSpawn) EncodeBlock() (string, map[string]any) {
	return "minecraft:frog_spawn", nil
}

// EncodeItem ...
func (f FrogSpawn) EncodeItem() (name string, meta int16) {
	return blockItemName(f), 0
}
