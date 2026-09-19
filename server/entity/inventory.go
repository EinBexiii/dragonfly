package entity

import (
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/world"
)

// InventoryCarrier is an Entity that carries an inventory a player can open,
// such as a saddled horse holding its saddle and armour.
type InventoryCarrier interface {
	world.Entity
	// CarriedInventory returns the inventory the Entity carries, together with
	// the equipment slots the client reserves in the window that shows it. A
	// nil inventory means the Entity has no window to open right now.
	CarriedInventory() (*inventory.Inventory, []InventorySlot)
}

// InventorySlot is one equipment slot of an InventoryCarrier's window. The
// client outlines Icon in the slot while it is empty and refuses to put
// anything in it that is not in Accepts.
type InventorySlot struct {
	// Slot is the index of the slot in the carried inventory.
	Slot int
	// Icon is the item outlined in the slot while it holds nothing.
	Icon world.Item
	// Accepts holds the only items the slot takes.
	Accepts []world.Item
}
