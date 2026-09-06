package session

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestVelocityWorthSending(t *testing.T) {
	kb := mgl32.Vec3{0.3, 0.36, 0.3}
	zero := mgl32.Vec3{}
	cases := []struct {
		name string
		last mgl32.Vec3
		sent bool
		vel  mgl32.Vec3
		want bool
	}{
		{"first", zero, false, zero, true},
		{"knockback", zero, true, kb, true},
		{"same knockback again", kb, true, kb, true},
		{"stop after knockback", kb, true, zero, true},
		{"resting repeats zero", zero, true, zero, false},
		{"near zero repeats", mgl32.Vec3{1e-5, 0, 0}, true, zero, false},
	}
	for _, c := range cases {
		if got := velocityWorthSending(c.last, c.sent, c.vel); got != c.want {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}
