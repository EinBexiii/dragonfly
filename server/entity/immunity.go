package entity

import "time"

// AttackImmunity is the invulnerability window an entity has after being
// hurt.
type AttackImmunity struct {
	until time.Time
	last  float64
}

// tickDuration is one world tick. A vanilla window of n ticks expires
// somewhere between n−1 and n ticks after the hit, depending on where in
// a tick the hit landed. Wall time is measured at processing time, after
// the proxy's flush, the backend's flush and the wait for the world
// goroutine, so two hits the client sent a window apart may reach here
// slightly closer together. Expiring one tick early covers both, as
// vanilla's own phase does.
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
// window (a stronger one, or one a handler let through) raises the damage
// the window remembers but does not restart it: restarting would let a
// strong second hit shield the target from a third that an expired window
// should have accepted.
func (a *AttackImmunity) Arm(d time.Duration, damage float64) {
	if time.Now().Before(a.until) {
		a.last = max(a.last, damage)
		return
	}
	a.until, a.last = time.Now().Add(max(0, d-tickDuration)), damage
}
