package server

import (
	"encoding/base64"
	"testing"

	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
)

func TestLoginSkinAnimationNone(t *testing.T) {
	data := login.ClientData{
		DeviceOS: 1, GameVersion: "1.26.45", LanguageCode: "en_US", ServerAddress: "127.0.0.1:19132",
		SkinID: "builtin", SkinImageWidth: 64, SkinImageHeight: 64,
		SkinData:          base64.StdEncoding.EncodeToString(make([]byte, 64*64*4)),
		SkinResourcePatch: base64.StdEncoding.EncodeToString([]byte(`{"geometry":{"default":"geometry.humanoid"}}`)),
	}
	for typ := range 4 {
		data.AnimatedImageData = append(data.AnimatedImageData, login.SkinAnimation{Type: typ})
	}
	if err := data.Validate(); err != nil {
		t.Fatal(err)
	}
	got := (&Server{}).parseSkin(data)
	if len(got.Model) != 0 {
		t.Fatal("empty geometry changed")
	}
	for i, typ := range []skin.AnimationType{skin.AnimationNone, skin.AnimationHead, skin.AnimationBody32x32, skin.AnimationBody128x128} {
		if got.Animations[i].Type() != typ {
			t.Errorf("animation %d: got %d, want %d", i, got.Animations[i].Type(), typ)
		}
	}
}
