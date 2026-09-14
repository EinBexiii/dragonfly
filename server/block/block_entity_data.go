package block

import "maps"

// blockEntityData preserves the block entity NBT of a block that the server decodes but does not otherwise interact
// with. Before these blocks had concrete types they were unknownBlock, which carried their saved record through
// untouched. A block that is registered but implements no world.NBTer makes the loader drop that record instead, and
// the next write of the chunk removes it from the save for good, so every wave 3 block the game gives a block entity
// embeds this and round-trips its data unchanged.
type blockEntityData struct {
	// nbt points at the block entity NBT exactly as it was read. It is held behind a pointer because a map field
	// would make the block value uncomparable, and blocks are compared and used as map keys throughout the server.
	nbt *map[string]any
}

// storeNBT returns the data with the block entity NBT passed stored in it unchanged.
func (blockEntityData) storeNBT(data map[string]any) blockEntityData {
	stored := maps.Clone(data)
	return blockEntityData{nbt: &stored}
}

// storedNBT returns a copy of the block entity NBT held. The id passed is the one the game knows the block entity by:
// it fills in only when the record carries none, which happens for a block that was in the save without a block entity
// of its own, and never overwrites a saved id. The map is always non-nil because the chunk writer adds the position to
// it.
func (b blockEntityData) storedNBT(id string) map[string]any {
	data := make(map[string]any)
	if b.nbt != nil {
		maps.Copy(data, *b.nbt)
	}
	if _, ok := data["id"]; !ok {
		data["id"] = id
	}
	return data
}
