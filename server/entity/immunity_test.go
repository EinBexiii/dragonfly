package entity

import (
	"testing"
	"time"
)

func TestAttackImmunityDoesNotRestartOnAStrongerHit(t *testing.T) {
	var a AttackImmunity
	a.Arm(time.Second, 4)
	if left, immune := a.Reduce(4); !immune || left != 0 {
		t.Fatalf("same hit inside the window: %v %v", left, immune)
	}
	// A stronger hit inside the window lands with the difference and raises
	// the remembered damage, but keeps the original expiry.
	until := a.until
	if left, immune := a.Reduce(6); !immune || left != 2 {
		t.Fatalf("stronger hit: %v %v", left, immune)
	}
	a.Arm(time.Second, 6)
	if a.until != until || a.last != 6 {
		t.Fatalf("window restarted: until moved %v, last %v", a.until != until, a.last)
	}
	a.Arm(time.Second, 3)
	if a.last != 6 {
		t.Fatalf("remembered damage lowered to %v", a.last)
	}
	a.until = time.Now().Add(-time.Millisecond)
	if _, immune := a.Reduce(4); immune {
		t.Fatal("expired window still immune")
	}
	a.Arm(time.Second, 4)
	if !time.Now().Before(a.until) || a.last != 4 {
		t.Fatalf("fresh window: until %v last %v", a.until, a.last)
	}
}

// The default half-second window expires a tick early, where vanilla's
// tick phase would already have released it.
func TestAttackImmunityExpiresATickEarly(t *testing.T) {
	var a AttackImmunity
	before := time.Now()
	a.Arm(time.Second/2, 1)
	if got := a.until.Sub(before); got > 450*time.Millisecond+time.Millisecond || got < 449*time.Millisecond {
		t.Fatalf("window is %v, want 450ms", got)
	}
	a.Arm(time.Millisecond, 1)
	a.until = time.Time{}
	a.Arm(10*time.Millisecond, 1)
	if time.Now().Before(a.until) {
		t.Fatal("a window shorter than a tick did not expire at once")
	}
}
