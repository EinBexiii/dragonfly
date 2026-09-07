package server

import (
	"encoding/base64"
	"image/color"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
)

func TestParseSkinPreservesDefaultPersona(t *testing.T) {
	enc := base64.StdEncoding.EncodeToString
	packID := uuid.MustParse("48bf1558-1e1e-4b65-b530-53aac6a21726")
	data := login.ClientData{
		SkinID: "default-persona", PersonaSkin: true, PlayFabID: "1234abcd", CapeID: "persona-cape",
		SkinImageWidth: 64, SkinImageHeight: 64, SkinData: enc(make([]byte, 64*64*4)),
		SkinResourcePatch: enc([]byte(`{"geometry":{"default":"geometry.persona"}}`)),
		SkinGeometry:      enc([]byte(`{}`)), SkinGeometryVersion: enc([]byte("1.21.0")),
		SkinAnimationData: enc([]byte(`{"animations":[]}`)),
		ArmSize:           "wide", SkinColour: "#b37b62", PremiumSkin: true, CapeOnClassicSkin: true,
		ProfileHash: "default-persona-profile",
		PersonaPieces: []login.PersonaPiece{
			{PieceID: "body", PieceType: "persona_body", PackID: packID.String(), Default: true},
			{PieceID: "hands", PieceType: "persona_hand", PackID: packID.String(), ProductID: "product"},
		},
		PieceTintColours: []login.PersonaPieceTintColour{{PieceType: "persona_eyes",
			Colours: [4]string{"#ff010203", "#80000000", "#0", "#ff280000"}}},
	}
	got := new(Server).parseSkin(data)
	wantPieces := []protocol.PersonaPiece{
		{PieceID: "body", PieceType: protocol.PieceTypeBody, PackID: packID, Default: true},
		{PieceID: "hands", PieceType: protocol.PieceTypeHands, PackID: packID, ProductID: "product"},
	}
	wantTints := []protocol.PersonaPieceTintColour{{PieceType: "persona_eyes",
		Colours: [4]color.RGBA{{R: 1, G: 2, B: 3, A: 255}, {A: 128}, {}, {R: 40, A: 255}}}}
	if !reflect.DeepEqual(got.PersonaPieces, wantPieces) || !reflect.DeepEqual(got.PieceTintColours, wantTints) {
		t.Errorf("login persona lost its pieces/tints: %v / %v", got.PersonaPieces, got.PieceTintColours)
	}
	if !got.Persona || got.ArmSize != protocol.ArmSizeWide || got.SkinColour != (color.RGBA{R: 0xb3, G: 0x7b, B: 0x62, A: 255}) {
		t.Errorf("login persona appearance changed: persona=%v arms=%v colour=%v", got.Persona, got.ArmSize, got.SkinColour)
	}
	if got.CapeID != data.CapeID || got.FullID != data.SkinID || got.PlayFabID != data.PlayFabID ||
		got.ProfileHash != data.ProfileHash || !got.Premium || !got.PersonaCapeOnClassicSkin ||
		string(got.GeometryDataEngineVersion) != "1.21.0" || string(got.AnimationData) != `{"animations":[]}` {
		t.Error("login persona metadata changed")
	}
}

func TestParseSkinColour(t *testing.T) {
	for _, tc := range []struct {
		name, value string
		argb        bool
		want        color.RGBA
	}{
		{"RGB", "#b37b62", false, color.RGBA{R: 0xb3, G: 0x7b, B: 0x62, A: 255}},
		{"ARGB", "#803aafd9", true, color.RGBA{R: 0x3a, G: 0xaf, B: 0xd9, A: 128}},
		{"transparent", "#0", true, color.RGBA{}},
		{"malformed", "invalid", true, color.RGBA{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := parseSkinColour(tc.value, tc.argb); got != tc.want {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}
