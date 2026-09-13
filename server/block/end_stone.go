package block

// EndStone is a block found in The End.
type EndStone struct {
	solid
	bassDrum
}

// ChorusRoot marks end stone as the block a chorus plant standing on it reaches its bottom bound down to.
func (EndStone) ChorusRoot() {}

// BreakInfo ...
func (e EndStone) BreakInfo() BreakInfo {
	return newBreakInfo(3, pickaxeHarvestable, pickaxeEffective, oneOf(e)).withBlastResistance(9)
}

// EncodeItem ...
func (EndStone) EncodeItem() (name string, meta int16) {
	return "minecraft:end_stone", 0
}

// EncodeBlock ...
func (EndStone) EncodeBlock() (string, map[string]any) {
	return "minecraft:end_stone", nil
}
