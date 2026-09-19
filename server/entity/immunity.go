package entity

import "time"

// AttackImmunity is the invulnerability window an entity has after being
// hurt.
type AttackImmunity struct {
	until  time.Time
	last   float64
	excess bool // the last hit landed inside a running window
}

// tickDuration is one world tick. The window expires a tick early, as a
// vanilla tick-counted window does depending on phase.
const tickDuration = time.Second / 20

// Reduce reduces the damage of a hit by what the window already absorbed.
// Callers deal the returned damage only if it is positive.
func (a *AttackImmunity) Reduce(damage float64) (left float64, immune bool) {
	if !time.Now().Before(a.until) {
		return damage, false
	}
	return damage - a.last, true
}

// Arm starts a window of d for a hit of damage. A hit inside a running
// window raises the remembered damage without restarting the window.
func (a *AttackImmunity) Arm(d time.Duration, damage float64) {
	if time.Now().Before(a.until) {
		a.last, a.excess = max(a.last, damage), true
		return
	}
	a.until, a.last, a.excess = time.Now().Add(max(0, d-tickDuration)), damage, false
}

// Absorbing reports whether the window is running and its last hit landed
// inside it: such a hit deals its excess and nothing else, so no knockback
// and no hurt feedback belong to it.
func (a *AttackImmunity) Absorbing() bool {
	return a.excess && time.Now().Before(a.until)
}
