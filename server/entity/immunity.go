package entity

import "time"

// AttackImmunity is the invulnerability window an entity has after being
// hurt.
type AttackImmunity struct {
	until time.Time
	last  float64
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
		a.last = max(a.last, damage)
		return
	}
	a.until, a.last = time.Now().Add(max(0, d-tickDuration)), damage
}
