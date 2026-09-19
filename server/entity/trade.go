package entity

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

// Trader is an Entity that trades with a player, such as a villager.
type Trader interface {
	world.Entity
	// TradeOffers returns the offers the Entity makes in the order the window
	// shows them, the trading tier it has reached and the name the window
	// carries.
	TradeOffers() (offers []Offer, tier int, name string)
	// Trade completes the offer at the index passed for the customer passed.
	// It returns false if the offer cannot be made, in which case nothing
	// about the Trader changes.
	Trade(i int, customer world.Entity) bool
}

// Offer is a single trade a Trader makes: up to two items asked for in
// exchange for one item given.
type Offer struct {
	// Wanted holds the items the Trader asks for. The second item is empty in
	// an offer that asks for a single item.
	Wanted [2]item.Stack
	// Given is the item the Trader hands to the customer in exchange for
	// Wanted.
	Given item.Stack
	// Uses is the number of times the offer was made since the Trader last
	// restocked it.
	Uses int
	// MaxUses is the number of times the offer may be made before it runs out
	// of stock. An offer that has no uses left is shown as unavailable.
	MaxUses int
	// Experience is the experience the Trader itself earns by making the
	// offer, counting towards its next trading tier.
	Experience int
	// RewardsExperience specifies if making the offer rewards the customer
	// with experience.
	RewardsExperience bool
	// Tier is the trading tier the Trader must have reached for the offer to
	// be available. It is zero-based, like the tier TradeOffers returns.
	Tier int
	// PriceMultiplier is the fraction of the count of the first wanted item
	// that every point of Demand adds to the price of the offer.
	PriceMultiplier float64
	// Demand measures how often the offer was made recently. It raises the
	// price of the first wanted item by PriceMultiplier of its count for each
	// point of demand.
	Demand int
}

// TradeWatcher is a Trader that is told while a customer stands at its window,
// so that it can attend to that customer instead of wandering off.
type TradeWatcher interface {
	Trader
	// TradeOpened is called when a customer opens the window, and TradeClosed
	// once that window is closed again.
	TradeOpened(customer world.Entity)
	TradeClosed()
}
