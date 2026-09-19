package session

import (
	"bytes"
	"image/color"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TestPersonaSkinSurvivesDecodeAndBroadcast(t *testing.T) {
	original := personaSkinFixture()
	stored, err := protocolToSkin(original)
	if err != nil {
		t.Fatal(err)
	}
	for range 2 {
		broadcast := skinToProtocol(stored)
		// Preserve the existing per-broadcast SkinID generation while checking
		// that the character's pieces and appearance metadata survive.
		original.SkinID = broadcast.SkinID
		if !reflect.DeepEqual(broadcast, original) {
			t.Errorf("persona did not round-trip: pieces=%v, tints=%v, arm=%v, colour=%v, cape=%q, full=%q, engine=%q, profile=%q",
				broadcast.PersonaPieces, broadcast.PieceTintColours, broadcast.ArmSize, broadcast.SkinColour,
				broadcast.CapeID, broadcast.FullID, broadcast.GeometryDataEngineVersion, broadcast.ProfileHash)
		}
		// Exercise the player-list wire path used for an observer joining.
		out := &packet.PlayerList{Entries: []protocol.PlayerListEntry{{
			ActionType: protocol.PlayerListActionAdd, UUID: uuid.New(), Skin: broadcast,
		}}}
		var buf bytes.Buffer
		out.Marshal(protocol.NewWriter(&buf, 0))
		encoded := bytes.Clone(buf.Bytes())
		var in packet.PlayerList
		in.Marshal(protocol.NewReader(&buf, 0, false))
		var again bytes.Buffer
		in.Marshal(protocol.NewWriter(&again, 0))
		if !bytes.Equal(encoded, again.Bytes()) {
			t.Fatal("player-list encoding lost persona metadata")
		}
	}
}

func TestClassicSkinStillRoundTrips(t *testing.T) {
	original := personaSkinFixture()
	original.PersonaSkin = false
	original.PersonaPieces, original.PieceTintColours = nil, nil
	stored, err := protocolToSkin(original)
	if err != nil {
		t.Fatal(err)
	}
	got := skinToProtocol(stored)
	original.SkinID = got.SkinID
	if !reflect.DeepEqual(got, original) {
		t.Fatal("classic skin changed in conversion")
	}
}

func personaSkinFixture() protocol.Skin {
	return protocol.Skin{
		SkinID: "selected-persona", PlayFabID: "1234abcd", FullID: "full-persona", CapeID: "cape-id",
		SkinImageWidth: 64, SkinImageHeight: 64, SkinData: make([]byte, 64*64*4),
		SkinGeometry: []byte(`{}`), SkinResourcePatch: []byte(`{"geometry":{"default":"geometry.persona"}}`),
		GeometryDataEngineVersion: []byte("1.21.0"), AnimationData: []byte(`{"animations":[]}`),
		PersonaSkin: true, PremiumSkin: true, PersonaCapeOnClassicSkin: true, PrimaryUser: true,
		ArmSize: protocol.ArmSizeWide, SkinColour: color.RGBA{R: 0xb3, G: 0x7b, B: 0x62, A: 255},
		PersonaPieces: []protocol.PersonaPiece{{PieceID: "body", PieceType: protocol.PieceTypeBody,
			PackID: uuid.MustParse("48bf1558-1e1e-4b65-b530-53aac6a21726"), Default: true, ProductID: "product"}},
		PieceTintColours: []protocol.PersonaPieceTintColour{{PieceType: "persona_eyes",
			Colours: [4]color.RGBA{{R: 1, G: 2, B: 3, A: 255}, {A: 128}, {}, {R: 40, A: 255}}}},
		ProfileHash: "persona-profile", Trusted: true, OverrideAppearance: true,
	}
}

func TestPersonaSkinWithoutFullIDKeepsStableIdentity(t *testing.T) {
	original := personaSkinFixture()
	original.FullID = ""
	stored, err := protocolToSkin(original)
	if err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if got := skinToProtocol(stored).FullID; got != original.SkinID {
			t.Fatalf("FullID = %q, want client SkinID %q", got, original.SkinID)
		}
	}
}
