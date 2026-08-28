package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"testing"
)

func TestLivingEntHurt(t *testing.T) {
	t.Run("subtracts health", func(t *testing.T) {
		w := world.Config{}.New()
		t.Cleanup(func() { _ = w.Close() })

		handle := spawnTestLivingEnt(t, w, newTestLivingBehaviour(20))
		mustDo(t, w, func(tx *world.Tx) {
			e, _ := handle.Entity(tx)
			l := e.(Living)
			if n, vulnerable := l.Hurt(3, testDamageSource{}); n != 3 || !vulnerable {
				t.Fatalf("Hurt(3) = %v, %v, want 3, true", n, vulnerable)
			}
			if l.Health() != 17 {
				t.Fatalf("Health() = %v after 3 damage, want 17", l.Health())
			}
		})
	})

	t.Run("clamps negative damage", func(t *testing.T) {
		w := world.Config{}.New()
		t.Cleanup(func() { _ = w.Close() })

		handle := spawnTestLivingEnt(t, w, newTestLivingBehaviour(20))
		mustDo(t, w, func(tx *world.Tx) {
			e, _ := handle.Entity(tx)
			l := e.(Living)
			if n, _ := l.Hurt(-5, testDamageSource{}); n != 0 {
				t.Fatalf("Hurt(-5) = %v, want 0", n)
			}
			if l.Health() != 20 {
				t.Fatalf("Health() = %v after negative damage, want 20", l.Health())
			}
		})
	})

	t.Run("is immune within HurtImmunity of a hit", func(t *testing.T) {
		w := world.Config{}.New()
		t.Cleanup(func() { _ = w.Close() })

		handle := spawnTestLivingEnt(t, w, newTestLivingBehaviour(20))
		mustDo(t, w, func(tx *world.Tx) {
			e, _ := handle.Entity(tx)
			l := e.(Living)
			l.Hurt(3, testDamageSource{})
			if n, vulnerable := l.Hurt(3, testDamageSource{}); n != 0 || vulnerable {
				t.Fatalf("second Hurt(3) = %v, %v, want 0, false", n, vulnerable)
			}
			if l.Health() != 17 {
				t.Fatalf("Health() = %v after an absorbed hit, want 17", l.Health())
			}
		})
	})
}

func TestLivingEntHandleHurt(t *testing.T) {
	t.Run("may cancel damage", func(t *testing.T) {
		w := world.Config{}.New()
		t.Cleanup(func() { _ = w.Close() })

		b := newTestLivingBehaviour(20)
		b.handleHurt = func(*LivingEnt, *float64, world.DamageSource) bool { return true }
		handle := spawnTestLivingEnt(t, w, b)
		mustDo(t, w, func(tx *world.Tx) {
			e, _ := handle.Entity(tx)
			l := e.(Living)
			if n, vulnerable := l.Hurt(3, testDamageSource{}); n != 0 || vulnerable {
				t.Fatalf("cancelled Hurt(3) = %v, %v, want 0, false", n, vulnerable)
			}
			if l.Health() != 20 {
				t.Fatalf("Health() = %v after a cancelled hit, want 20", l.Health())
			}
		})
	})

	t.Run("may change damage", func(t *testing.T) {
		w := world.Config{}.New()
		t.Cleanup(func() { _ = w.Close() })

		b := newTestLivingBehaviour(20)
		b.handleHurt = func(_ *LivingEnt, damage *float64, _ world.DamageSource) bool {
			*damage /= 2
			return false
		}
		handle := spawnTestLivingEnt(t, w, b)
		mustDo(t, w, func(tx *world.Tx) {
			e, _ := handle.Entity(tx)
			l := e.(Living)
			if n, _ := l.Hurt(6, testDamageSource{}); n != 3 {
				t.Fatalf("halved Hurt(6) = %v, want 3", n)
			}
			if l.Health() != 17 {
				t.Fatalf("Health() = %v after halved damage, want 17", l.Health())
			}
		})
	})
}

func TestLivingEntDeath(t *testing.T) {
	w := world.Config{}.New()
	t.Cleanup(func() { _ = w.Close() })

	var (
		deathSrc world.DamageSource
		died     bool
	)
	b := newTestLivingBehaviour(5)
	b.handleDeath = func(_ *LivingEnt, _ *world.Tx, src world.DamageSource) {
		died, deathSrc = true, src
	}
	handle := spawnTestLivingEnt(t, w, b)

	mustDo(t, w, func(tx *world.Tx) {
		e, _ := handle.Entity(tx)
		l := e.(Living)
		if n, vulnerable := l.Hurt(10, testDamageSource{}); n != 10 || !vulnerable {
			t.Fatalf("fatal Hurt(10) = %v, %v, want 10, true", n, vulnerable)
		}
		if !l.Dead() {
			t.Fatal("Dead() = false after fatal damage")
		}
		if _, ok := handle.Entity(tx); !ok {
			t.Fatal("entity removed within Hurt: removal must be deferred to the next tick")
		}
	})
	mustDo(t, w, func(tx *world.Tx) {
		e, _ := handle.Entity(tx)
		e.(world.TickerEntity).Tick(tx, 0)
	})
	mustDo(t, w, func(tx *world.Tx) {
		if _, ok := handle.Entity(tx); ok {
			t.Fatal("entity still in the world after its death tick")
		}
	})
	if !died {
		t.Fatal("HandleDeath was not called on the death tick")
	}
	if _, ok := deathSrc.(testDamageSource); !ok {
		t.Fatalf("HandleDeath source = %T, want testDamageSource", deathSrc)
	}
}

// spawnTestLivingEnt spawns an entity backed by the behaviour passed into the world and returns its handle.
func spawnTestLivingEnt(t *testing.T, w *world.World, b *testLivingBehaviour) *world.EntityHandle {
	t.Helper()
	handle := world.EntitySpawnOpts{Position: mgl64.Vec3{0.5, 1, 0.5}}.New(testLivingEntType{}, testLivingConfig{b: b})
	mustDo(t, w, func(tx *world.Tx) {
		tx.AddEntity(handle)
	})
	return handle
}

type testDamageSource struct{}

func (testDamageSource) ReducedByArmour() bool     { return false }
func (testDamageSource) ReducedByResistance() bool { return false }
func (testDamageSource) Fire() bool                { return false }
func (testDamageSource) IgnoreTotem() bool         { return false }

type testLivingConfig struct {
	b *testLivingBehaviour
}

func (c testLivingConfig) Apply(data *world.EntityData) {
	data.Data = c.b
}

type testLivingBehaviour struct {
	BaseBehaviour

	state       *LivingState
	handleHurt  func(e *LivingEnt, damage *float64, src world.DamageSource) bool
	handleDeath func(e *LivingEnt, tx *world.Tx, src world.DamageSource)
}

func newTestLivingBehaviour(health float64) *testLivingBehaviour {
	return &testLivingBehaviour{BaseBehaviour: NewBaseBehaviour(), state: NewLivingState(health)}
}

func (b *testLivingBehaviour) Tick(*Ent, *world.Tx) *Movement { return nil }

func (b *testLivingBehaviour) LivingState() *LivingState { return b.state }

func (b *testLivingBehaviour) HandleHurt(e *LivingEnt, damage *float64, src world.DamageSource) bool {
	if b.handleHurt == nil {
		return false
	}
	return b.handleHurt(e, damage, src)
}

func (b *testLivingBehaviour) HandleDeath(e *LivingEnt, tx *world.Tx, src world.DamageSource) {
	if b.handleDeath != nil {
		b.handleDeath(e, tx, src)
	}
}

type testLivingEntType struct{}

func (testLivingEntType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return OpenLiving(tx, handle, data)
}

func (testLivingEntType) EncodeEntity() string { return "minecraft:test_living_ent" }
func (testLivingEntType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.8, 0.3)
}
func (testLivingEntType) DecodeNBT(map[string]any, *world.EntityData) {}
func (testLivingEntType) EncodeNBT(*world.EntityData) map[string]any  { return nil }
