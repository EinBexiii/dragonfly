package block

// MangroveRoots is the tangle of roots a mangrove tree grows on. It is a full cube: the 2026-09-14 BDS review, batch B
// section 10, finds it inheriting the unit-box getter despite the gaps in its model.
type MangroveRoots struct {
	solid
	transparent
}

// EncodeBlock ...
func (MangroveRoots) EncodeBlock() (string, map[string]any) {
	return "minecraft:mangrove_roots", nil
}

// EncodeItem ...
func (m MangroveRoots) EncodeItem() (name string, meta int16) {
	return blockItemName(m), 0
}

// HangingRoots are the roots that hang from the ceiling of a lush cave. They have no collision: batch B section 10
// finds the dimensions of their constructor box bypassed by the empty collision getter.
type HangingRoots struct {
	empty
	transparent
}

// EncodeBlock ...
func (HangingRoots) EncodeBlock() (string, map[string]any) {
	return "minecraft:hanging_roots", nil
}

// EncodeItem ...
func (h HangingRoots) EncodeItem() (name string, meta int16) {
	return blockItemName(h), 0
}
