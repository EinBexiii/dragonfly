package session

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"reflect"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func liveSkinFixture() protocol.Skin {
	return protocol.Skin{
		SkinID: "live-selection", PlayFabID: "abcd",
		SkinImageWidth: 64, SkinImageHeight: 64, SkinData: make([]byte, 64*64*4),
		SkinGeometry:      []byte(`{}`),
		SkinResourcePatch: []byte(`{"geometry":{"default":"geometry.humanoid.custom"}}`),
	}
}

// Validate the equivalent login data with the actual pinned login validator.
func liveSkinLoginData(sk protocol.Skin) login.ClientData {
	enc := base64.StdEncoding.EncodeToString
	data := login.ClientData{
		DeviceOS: 1, GameVersion: "1.26.45", LanguageCode: "en_US", ServerAddress: "127.0.0.1:19132",
		SkinID: sk.SkinID, PlayFabID: sk.PlayFabID,
		SkinImageWidth: int(sk.SkinImageWidth), SkinImageHeight: int(sk.SkinImageHeight), SkinData: enc(sk.SkinData),
		CapeImageWidth: int(sk.CapeImageWidth), CapeImageHeight: int(sk.CapeImageHeight), CapeData: enc(sk.CapeData),
		SkinGeometry: enc(sk.SkinGeometry), SkinResourcePatch: enc(sk.SkinResourcePatch),
	}
	for _, a := range sk.Animations {
		data.AnimatedImageData = append(data.AnimatedImageData, login.SkinAnimation{
			ImageWidth: int(a.ImageWidth), ImageHeight: int(a.ImageHeight), Image: enc(a.ImageData),
			Type: int(a.AnimationType), Frames: float64(a.FrameCount), AnimationExpression: int(a.ExpressionType),
		})
	}
	return data
}

func liveSkinWireRoundTrip(sk protocol.Skin) protocol.Skin {
	var buf bytes.Buffer
	out := &packet.PlayerSkin{Skin: sk}
	out.Marshal(protocol.NewWriter(&buf, 0))
	var in packet.PlayerSkin
	in.Marshal(protocol.NewReader(&buf, 0, true))
	return in.Skin
}

func TestLiveSkinLoginAcceptedFields(t *testing.T) {
	cases := []struct {
		name   string
		change func(*protocol.Skin)
	}{
		{"empty geometry", func(s *protocol.Skin) { s.SkinGeometry = nil }},
		{"null geometry", func(s *protocol.Skin) { s.SkinGeometry = []byte(`null`) }},
		{"custom geometry", func(s *protocol.Skin) { s.SkinGeometry = []byte(`{"minecraft:geometry":[]}`) }},
		{"empty images", func(s *protocol.Skin) { s.SkinImageWidth, s.SkinImageHeight, s.SkinData = 0, 0, nil }},
		{"zero width", func(s *protocol.Skin) { s.SkinImageWidth, s.SkinData = 0, nil }},
		{"empty optional fields", func(s *protocol.Skin) { s.PlayFabID, s.FullID, s.CapeID = "", "", "" }},
	}
	for _, patch := range []string{`null`, `{}`, `{"geometry":null}`, `{"geometry":{}}`, `{"geometry":{"default":"","animated_face":""}}`, `{"geometry":{"default":null,"animated_face":null}}`, `{"unknown":[],"geometry":{"default":"geometry.humanoid","unknown":42}}`} {
		cases = append(cases, struct {
			name   string
			change func(*protocol.Skin)
		}{"patch " + patch, func(s *protocol.Skin) { s.SkinResourcePatch = []byte(patch) }})
	}
	for typ := uint32(0); typ <= 3; typ++ {
		cases = append(cases, struct {
			name   string
			change func(*protocol.Skin)
		}{fmt.Sprintf("animation %d", typ), func(s *protocol.Skin) {
			s.Animations = []protocol.SkinAnimation{{AnimationType: typ}, {
				AnimationType: typ, ImageWidth: 32, ImageHeight: 32, ImageData: make([]byte, 32*32*4), FrameCount: 1,
			}}
		}})
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sk := liveSkinFixture()
			tc.change(&sk)
			if err := liveSkinLoginData(sk).Validate(); err != nil {
				t.Fatalf("login rejected fixture: %v", err)
			}
			sk = liveSkinWireRoundTrip(sk)
			got, err := protocolToSkin(sk)
			if err != nil {
				t.Fatalf("live rejected login-accepted skin: %v", err)
			}
			if !bytes.Equal(got.Model, sk.SkinGeometry) {
				t.Fatal("geometry changed during decode")
			}
			broadcast := liveSkinWireRoundTrip(skinToProtocol(got))
			if !reflect.DeepEqual(broadcast.Animations, sk.Animations) {
				t.Fatalf("animation changed: got %+v, want %+v", broadcast.Animations, sk.Animations)
			}
			if broadcast.SkinImageWidth != sk.SkinImageWidth || broadcast.SkinImageHeight != sk.SkinImageHeight || !bytes.Equal(broadcast.SkinData, sk.SkinData) {
				t.Fatal("skin image changed")
			}
			if broadcast.SkinResourcePatch != nil && string(broadcast.SkinResourcePatch) != string(got.ModelConfig.Encode()) {
				t.Fatal("resource patch changed")
			}
		})
	}
}

func TestLiveSkinRejectsMalformedFields(t *testing.T) {
	cases := []struct {
		name   string
		change func(*protocol.Skin)
	}{
		{"empty ID", func(s *protocol.Skin) { s.SkinID = "" }},
		{"geometry syntax", func(s *protocol.Skin) { s.SkinGeometry = []byte(`{`) }},
		{"geometry whitespace", func(s *protocol.Skin) { s.SkinGeometry = []byte(` `) }},
		{"geometry array", func(s *protocol.Skin) { s.SkinGeometry = []byte(`[]`) }},
		{"empty patch", func(s *protocol.Skin) { s.SkinResourcePatch = nil }},
		{"patch syntax", func(s *protocol.Skin) { s.SkinResourcePatch = []byte(`{`) }},
		{"patch array", func(s *protocol.Skin) { s.SkinResourcePatch = []byte(`[]`) }},
		{"animation type", func(s *protocol.Skin) { s.Animations = []protocol.SkinAnimation{{AnimationType: 4}} }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sk := liveSkinFixture()
			tc.change(&sk)
			if liveSkinLoginData(sk).Validate() == nil {
				t.Fatal("login unexpectedly accepted malformed fixture")
			}
			if _, err := protocolToSkin(sk); err == nil {
				t.Fatal("accepted malformed skin")
			}
		})
	}
}

// Login checks a generic object and parseSkin discards typed decoding errors.
// These mismatches are malformed model references, not empty optional fields.
func TestLiveSkinRetainsResourcePatchTypes(t *testing.T) {
	for _, patch := range []string{`{"geometry":[]}`, `{"geometry":""}`, `{"geometry":0}`, `{"geometry":false}`, `{"geometry":{"default":0}}`, `{"geometry":{"default":[]}}`, `{"geometry":{"default":{}}}`, `{"geometry":{"animated_face":false}}`, `{"Geometry":{"Default":42}}`} {
		t.Run(patch, func(t *testing.T) {
			sk := liveSkinFixture()
			sk.SkinResourcePatch = []byte(patch)
			if err := liveSkinLoginData(sk).Validate(); err != nil {
				t.Fatal(err)
			}
			if _, err := protocolToSkin(liveSkinWireRoundTrip(sk)); err == nil {
				t.Fatal("accepted invalid model reference types")
			}
		})
	}
}

func TestLiveSkinImageLengths(t *testing.T) {
	for _, field := range []string{"skin", "cape", "animation"} {
		for _, size := range []struct {
			name string
			w, h uint32
			data []byte
		}{
			{"short", 2, 2, make([]byte, 12)}, {"long", 2, 2, make([]byte, 20)},
			{"partial pixel", 2, 2, make([]byte, 15)}, {"zero with data", 0, 2, make([]byte, 4)},
			{"uint32 overflow", 1 << 16, 1 << 16, nil}, {"byte count overflow", 1 << 30, 1, nil},
			{"uint64 byte overflow", 1 << 31, 1 << 31, nil}, {"maximum dimensions", ^uint32(0), ^uint32(0), nil},
		} {
			t.Run(field+"/"+size.name, func(t *testing.T) {
				sk := liveSkinFixture()
				switch field {
				case "skin":
					sk.SkinImageWidth, sk.SkinImageHeight, sk.SkinData = size.w, size.h, size.data
				case "cape":
					sk.CapeImageWidth, sk.CapeImageHeight, sk.CapeData = size.w, size.h, size.data
				case "animation":
					sk.Animations = []protocol.SkinAnimation{{AnimationType: protocol.SkinAnimationHead, ImageWidth: size.w, ImageHeight: size.h, ImageData: size.data}}
				}
				// These dimensions wrap to zero in the wire codec's uint32
				// check. Confirm the live decoder rejects the decoded packet too.
				if size.name == "uint32 overflow" || size.name == "byte count overflow" || size.name == "uint64 byte overflow" {
					sk = liveSkinWireRoundTrip(sk)
				}
				if _, err := protocolToSkin(sk); err == nil {
					t.Fatal("accepted inconsistent image")
				}
			})
		}
	}
}
