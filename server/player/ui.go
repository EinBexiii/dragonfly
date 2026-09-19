package player

import (
	"github.com/df-mc/dragonfly/server/player/bossbar"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/player/title"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// UI is the screen of a player: one-way display elements that exist on the client only and carry no
// state in the world. Interactive surfaces that the client reports back on, such that the player remains
// their receiver, stay on Player. It is accessed through Player.UI.
type UI struct {
	p *Player
}

// UI returns the UI of the player.
func (p *Player) UI() UI {
	return UI{p: p}
}

// SendPopup sends a formatted popup to the player. The popup is shown above the hotbar of the player and
// overwrites/is overwritten by the name of the item equipped.
// The popup is formatted following the rules of fmt.Sprintln without a newline at the end.
func (u UI) SendPopup(a ...any) {
	u.p.session().SendPopup(format(a))
}

// SendTip sends a tip to the player. The tip is shown in the middle of the screen of the player.
// The tip is formatted following the rules of fmt.Sprintln without a newline at the end.
func (u UI) SendTip(a ...any) {
	u.p.session().SendTip(format(a))
}

// SendToast sends a toast to the player. This toast is shown at the top of the screen, similar to achievements or pack
// loading.
func (u UI) SendToast(title, message string) {
	u.p.session().SendToast(title, message)
}

// SendTitle sends a title to the player. The title may be configured to change the duration it is displayed
// and the text it shows.
// If non-empty, the subtitle is shown in a smaller font below the title. The same counts for the action text
// of the title, which is shown in a font similar to that of a tip/popup.
func (u UI) SendTitle(t title.Title) {
	p := u.p
	p.session().SetTitleDurations(t.FadeInDuration(), t.Duration(), t.FadeOutDuration())
	if t.Text() != "" || t.Subtitle() != "" {
		p.session().SendTitle(t.Text())
		if t.Subtitle() != "" {
			p.session().SendSubtitle(t.Subtitle())
		}
	}
	if t.ActionText() != "" {
		p.session().SendActionBarMessage(t.ActionText())
	}
}

// SendScoreboard sends a scoreboard to the player. The scoreboard will be present indefinitely until removed
// by the caller.
// SendScoreboard may be called at any time to change the scoreboard of the player.
func (u UI) SendScoreboard(scoreboard *scoreboard.Scoreboard) {
	u.p.session().SendScoreboard(scoreboard)
}

// RemoveScoreboard removes any scoreboard currently present on the screen of the player. Nothing happens if
// the player has no scoreboard currently active.
func (u UI) RemoveScoreboard() {
	u.p.session().RemoveScoreboard()
}

// SendBossBar sends a boss bar to the player, so that it will be shown indefinitely at the top of the
// player's screen.
// The boss bar may be removed by calling UI.RemoveBossBar.
func (u UI) SendBossBar(bar bossbar.BossBar) {
	u.p.session().SendBossBar(bar.Text(), bar.Colour().Uint8(), bar.HealthPercentage())
}

// RemoveBossBar removes any boss bar currently active on the player's screen. If no boss bar is currently
// present, nothing happens.
func (u UI) RemoveBossBar() {
	u.p.session().RemoveBossBar()
}

// ShowCoordinates enables the vanilla coordinates for the player.
func (u UI) ShowCoordinates() {
	u.p.session().EnableCoordinates(true)
}

// HideCoordinates disables the vanilla coordinates for the player.
func (u UI) HideCoordinates() {
	u.p.session().EnableCoordinates(false)
}

// EnableInstantRespawn enables the vanilla instant respawn for the player.
func (u UI) EnableInstantRespawn() {
	u.p.session().EnableInstantRespawn(true)
}

// DisableInstantRespawn disables the vanilla instant respawn for the player.
func (u UI) DisableInstantRespawn() {
	u.p.session().EnableInstantRespawn(false)
}

// HideEntity hides a world.Entity from the Player so that it can under no circumstance see it. Hidden entities can be
// made visible again through a call to ShowEntity.
func (u UI) HideEntity(e world.Entity) {
	p := u.p
	if p.session() != session.Nop && p.H() != e.H() {
		p.session().StopShowingEntity(e)
	}
}

// ShowEntity shows a world.Entity previously hidden from the Player using HideEntity. It does nothing if the entity
// wasn't currently hidden.
func (u UI) ShowEntity(e world.Entity) {
	p := u.p
	if p.session() != session.Nop {
		p.session().StartShowingEntity(e)
	}
}

// ShowParticle shows a particle that only this Player can see. Unlike World.AddParticle, it is not broadcast
// to players around it.
func (u UI) ShowParticle(pos mgl64.Vec3, particle world.Particle) {
	u.p.session().ViewParticle(pos, particle)
}

// ViewNameTag overrides the public name tag of the entity for this player.
func (u UI) ViewNameTag(entity world.Entity, nameTag string) {
	u.p.session().ViewNameTag(entity, nameTag)
}

// ViewPublicNameTag removes the name tag override of the entity for this player.
func (u UI) ViewPublicNameTag(entity world.Entity) {
	u.p.session().ViewPublicNameTag(entity)
}

// ViewAlwaysShowNameTag overrides whether the entity's name tag is shown at all distances for this player.
func (u UI) ViewAlwaysShowNameTag(entity world.Entity, alwaysShow bool) {
	u.p.session().ViewAlwaysShowNameTag(entity, alwaysShow)
}

// ViewPublicAlwaysShowNameTag removes the always-show name tag override of the entity for this player.
func (u UI) ViewPublicAlwaysShowNameTag(entity world.Entity) {
	u.p.session().ViewPublicAlwaysShowNameTag(entity)
}

// ViewScoreTag overrides the public score tag of the entity for this player.
func (u UI) ViewScoreTag(entity world.Entity, scoreTag string) {
	u.p.session().ViewScoreTag(entity, scoreTag)
}

// ViewPublicScoreTag removes the score tag override of the entity for this player.
func (u UI) ViewPublicScoreTag(entity world.Entity) {
	u.p.session().ViewPublicScoreTag(entity)
}

// ViewVisibility overrides the public visibility of the entity for this player.
func (u UI) ViewVisibility(entity world.Entity, level world.VisibilityLevel) {
	u.p.session().ViewVisibility(entity, level)
}

// RemoveViewLayer removes all view-layer overrides of the entity for this player.
func (u UI) RemoveViewLayer(entity world.Entity) {
	u.p.session().RemoveViewLayer(entity)
}
