package item

// HorseArmourTier is the material a piece of horse armour is made of, which
// decides how much damage it takes off the horse wearing it.
type HorseArmourTier struct {
	// Name is the material's name in the item identifier.
	Name string
	// Protection is the damage reduction the armour gives the horse.
	Protection float64
}

// HorseArmourTierLeather is the weakest horse armour.
func HorseArmourTierLeather() HorseArmourTier {
	return HorseArmourTier{Name: "leather", Protection: 3}
}

// HorseArmourTierIron is horse armour of iron.
func HorseArmourTierIron() HorseArmourTier {
	return HorseArmourTier{Name: "iron", Protection: 5}
}

// HorseArmourTierGold is horse armour of gold.
func HorseArmourTierGold() HorseArmourTier {
	return HorseArmourTier{Name: "golden", Protection: 7}
}

// HorseArmourTierDiamond is the strongest horse armour.
func HorseArmourTierDiamond() HorseArmourTier {
	return HorseArmourTier{Name: "diamond", Protection: 11}
}

// HorseArmourTiers returns every tier horse armour comes in.
func HorseArmourTiers() []HorseArmourTier {
	return []HorseArmourTier{
		HorseArmourTierLeather(), HorseArmourTierIron(),
		HorseArmourTierGold(), HorseArmourTierDiamond(),
	}
}

// HorseArmour is armour worn by a horse, put on it through its inventory.
type HorseArmour struct {
	// Tier is the material the armour is made of.
	Tier HorseArmourTier
}

// MaxCount always returns 1.
func (HorseArmour) MaxCount() int {
	return 1
}

// EncodeItem ...
func (a HorseArmour) EncodeItem() (name string, meta int16) {
	return "minecraft:" + a.Tier.Name + "_horse_armor", 0
}
