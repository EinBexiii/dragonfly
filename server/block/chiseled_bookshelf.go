package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// ChiseledBookshelf is a bookshelf with six slots that books can be put into and taken out of. It is a full cube in
// every state: the 2026-09-14 BDS review, batch B section 7, finds it inheriting the unit-box getter.
type ChiseledBookshelf struct {
	solid
	bass
	blockEntityData

	// BooksStored is the bit set of occupied slots, from 0 to 63.
	BooksStored int
	// Facing is the direction the front of the bookshelf points in.
	Facing cube.Direction
}

// EncodeBlock ...
func (c ChiseledBookshelf) EncodeBlock() (string, map[string]any) {
	return "minecraft:chiseled_bookshelf", map[string]any{
		"books_stored": int32(c.BooksStored),
		"direction":    int32(horizontalDirection(c.Facing)),
	}
}

// EncodeItem ...
func (c ChiseledBookshelf) EncodeItem() (name string, meta int16) {
	return blockItemName(c), 0
}

// EncodeNBT ...
func (c ChiseledBookshelf) EncodeNBT() map[string]any {
	return c.storedNBT("ChiseledBookshelf")
}

// DecodeNBT ...
func (c ChiseledBookshelf) DecodeNBT(data map[string]any) any {
	c.blockEntityData = c.storeNBT(data)
	return c
}

// allChiseledBookshelves returns all chiseled bookshelf states.
func allChiseledBookshelves() (bookshelves []world.Block) {
	for books := range 64 {
		for _, d := range cube.Directions() {
			bookshelves = append(bookshelves, ChiseledBookshelf{BooksStored: books, Facing: d})
		}
	}
	return
}
