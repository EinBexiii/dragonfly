package item

// Saddle is an item that lets a player ride a mob it is put on.
type Saddle struct{}

// MaxCount always returns 1.
func (Saddle) MaxCount() int {
	return 1
}

// EncodeItem ...
func (Saddle) EncodeItem() (name string, meta int16) {
	return "minecraft:saddle", 0
}
