package entity

import (
	"time"

	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// LivingBehaviour is a Behaviour for entities that are alive. The LivingState it returns holds the state,
// such as health and effects, that a LivingEnt exposes through the Living interface. Because an entity is
// reopened in every transaction, this state must be part of the Behaviour rather than the entity itself.
type LivingBehaviour interface {
	Behaviour
	// LivingState returns the state of the living entity. It must return the same LivingState for the
	// lifetime of the entity.
	LivingState() *LivingState
}

// livingHurtHandler may be implemented by a LivingBehaviour to handle damage dealt to its entity. It is
// called before any other processing of the damage, so it also runs while the entity is immune to attacks.
type livingHurtHandler interface {
	// HandleHurt handles damage about to be dealt to the entity. The damage may be changed through the
	// pointer passed. If true is returned, the damage is cancelled entirely.
	HandleHurt(e *LivingEnt, damage *float64, src world.DamageSource) bool
}

// livingDeathHandler may be implemented by a LivingBehaviour to react to the death of its entity, for
// example to spawn drops. It is called on the tick after the fatal damage, so a *world.Tx is available.
type livingDeathHandler interface {
	// HandleDeath handles the death of the entity, before it is closed and removed from the world.
	HandleDeath(e *LivingEnt, tx *world.Tx, src world.DamageSource)
}

// LivingState holds the state of a living entity that persists across transactions. A LivingBehaviour
// creates one using NewLivingState and returns it from its LivingState method.
type LivingState struct {
	// Health manages the entity's current and maximum health.
	Health *HealthManager
	// Effects manages the effects active on the entity. They are ticked by LivingEnt.Tick.
	Effects *EffectManager
	// Speed is the movement speed of the entity.
	Speed float64
	// KnockBackResistance is the fraction, between 0 and 1, by which knock back on the entity is reduced.
	KnockBackResistance float64
	// HurtImmunity is the duration for which the entity is immune to damage after being hurt.
	HurtImmunity time.Duration

	immuneUntil time.Duration
	dead        bool
	deathSrc    world.DamageSource
}

// NewLivingState returns a LivingState with the health passed and defaults for the remaining fields: a speed
// of 0.1 and half a second of immunity after being hurt.
func NewLivingState(health float64) *LivingState {
	return &LivingState{
		Health:       NewHealthManager(health, health),
		Effects:      NewEffectManager(),
		Speed:        0.1,
		HurtImmunity: time.Second / 2,
	}
}

// Dead checks if the living entity has taken fatal damage.
func (s *LivingState) Dead() bool {
	return s.dead
}

// LivingEnt is an Ent that is alive. It implements Living on top of the LivingState of its Behaviour, which
// makes the entity part of everything gated on Living, such as the full combat path of a player attack,
// projectile knock back and potion effects. Entity implementations based on a LivingBehaviour return a
// LivingEnt from the Open method of their world.EntityType by calling OpenLiving.
type LivingEnt struct {
	*Ent
}

var _ Living = (*LivingEnt)(nil)

// OpenLiving converts a world.EntityHandle to a LivingEnt in a world.Tx. It panics if the Behaviour of the
// entity does not implement LivingBehaviour.
func OpenLiving(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) *LivingEnt {
	e := &LivingEnt{Ent: Open(tx, handle, data)}
	_ = e.behaviour()
	return e
}

// behaviour returns the Behaviour of the entity as a LivingBehaviour.
func (e *LivingEnt) behaviour() LivingBehaviour {
	return e.Behaviour().(LivingBehaviour)
}

// state returns the LivingState of the entity's Behaviour.
func (e *LivingEnt) state() *LivingState {
	return e.behaviour().LivingState()
}

// Health returns the current health of the entity.
func (e *LivingEnt) Health() float64 {
	return e.state().Health.Health()
}

// MaxHealth returns the maximum health of the entity.
func (e *LivingEnt) MaxHealth() float64 {
	return e.state().Health.MaxHealth()
}

// SetMaxHealth changes the maximum health of the entity to the value passed.
func (e *LivingEnt) SetMaxHealth(v float64) {
	e.state().Health.SetMaxHealth(v)
}

// Dead checks if the entity has taken fatal damage.
func (e *LivingEnt) Dead() bool {
	return e.state().Dead()
}

// Hurt hurts the entity for a given amount of damage. The damage is first passed to the Behaviour's
// HandleHurt, if implemented, which may change or cancel it. Damage is then absorbed if the entity is still
// immune from an earlier hit. If the damage is fatal, the entity is killed on its next tick, so that drops
// spawned by the Behaviour's HandleDeath have a transaction to be added in.
func (e *LivingEnt) Hurt(damage float64, src world.DamageSource) (float64, bool) {
	damage = max(damage, 0)
	s := e.state()
	if s.dead {
		return 0, false
	}
	if h, ok := e.Behaviour().(livingHurtHandler); ok && h.HandleHurt(e, &damage, src) {
		return 0, false
	}
	if e.Age() < s.immuneUntil {
		return 0, false
	}
	s.immuneUntil = e.Age() + s.HurtImmunity

	s.Health.AddHealth(-damage)
	e.PlayAction(HurtAction{})
	if s.Health.Health() <= 0 {
		s.dead, s.deathSrc = true, src
	}
	return damage, true
}

// Heal heals the entity for a given amount of health, up to its maximum health, and returns the amount of
// health that was regenerated. A dead entity is not healed.
func (e *LivingEnt) Heal(health float64, _ world.HealingSource) float64 {
	s := e.state()
	if s.dead || health <= 0 {
		return 0
	}
	before := s.Health.Health()
	s.Health.AddHealth(health)
	return s.Health.Health() - before
}

// KnockBack knocks the entity back away from the source position passed, with the knock back resistance of
// its LivingState applied. A dead entity is not knocked back.
func (e *LivingEnt) KnockBack(src mgl64.Vec3, force, height float64) {
	s := e.state()
	if s.dead {
		return
	}
	velocity := KnockBackVector(e.Position(), src, force, height)
	e.SetVelocity(velocity.Mul(1 - s.KnockBackResistance))
}

// AddEffect adds an effect to the entity. If the effect is instant, it is applied immediately. If not, it is
// applied every tick until it expires.
func (e *LivingEnt) AddEffect(eff effect.Effect) {
	e.state().Effects.Add(eff, e)
}

// RemoveEffect removes any effect of the type passed that is active on the entity.
func (e *LivingEnt) RemoveEffect(t effect.Type) {
	e.state().Effects.Remove(t, e)
}

// Effects returns the effects currently active on the entity.
func (e *LivingEnt) Effects() []effect.Effect {
	return e.state().Effects.Effects()
}

// Speed returns the current speed of the entity.
func (e *LivingEnt) Speed() float64 {
	return e.state().Speed
}

// SetSpeed sets the speed of the entity to a new value.
func (e *LivingEnt) SetSpeed(v float64) {
	e.state().Speed = v
}

// Tick ticks the entity. A dead entity plays its death animation, has the Behaviour's HandleDeath called if
// implemented and is closed. A living entity has its effects ticked before the Ent is ticked as usual.
func (e *LivingEnt) Tick(tx *world.Tx, current int64) {
	s := e.state()
	if s.dead {
		e.PlayAction(DeathAction{})
		if h, ok := e.Behaviour().(livingDeathHandler); ok {
			h.HandleDeath(e, tx, s.deathSrc)
		}
		_ = e.Close()
		return
	}
	s.Effects.Tick(e, tx)
	e.Ent.Tick(tx, current)
}
