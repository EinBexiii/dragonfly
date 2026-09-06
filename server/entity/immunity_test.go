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
	// A weaker hit afterwards does not lower what is remembered.
	a.Arm(time.Second, 3)
	if a.last != 6 {
		t.Fatalf("remembered damage lowered to %v", a.last)
	}
	// Once expired, the next hit starts a fresh window.
	a.until = time.Now().Add(-time.Millisecond)
	if _, immune := a.Reduce(4); immune {
		t.Fatal("expired window still immune")
	}
	a.Arm(time.Second, 4)
	if !time.Now().Before(a.until) || a.last != 4 {
		t.Fatalf("fresh window: until %v last %v", a.until, a.last)
	}
}
