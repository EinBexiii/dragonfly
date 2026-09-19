package server

import (
	"encoding/base64"
	"image/color"
	"strconv"
	"strings"

	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
)

// parsePersona preserves the appearance metadata carried alongside the baked
// texture. Login encodes piece types as names and colours as RGB/ARGB strings;
// skin.Skin uses the protocol representation shared with live skin updates.
func parsePersona(data login.ClientData, s *skin.Skin) {
	s.CapeID = data.CapeID
	s.Premium = data.PremiumSkin
	s.PersonaCapeOnClassicSkin = data.CapeOnClassicSkin
	s.ProfileHash = data.ProfileHash
	s.GeometryDataEngineVersion, _ = base64.StdEncoding.DecodeString(data.SkinGeometryVersion)
	s.AnimationData, _ = base64.StdEncoding.DecodeString(data.SkinAnimationData)
	if data.ArmSize == "wide" {
		s.ArmSize = protocol.ArmSizeWide
	}
	s.SkinColour = parseSkinColour(data.SkinColour, false)
	s.PersonaPieces = make([]protocol.PersonaPiece, len(data.PersonaPieces))
	for i, p := range data.PersonaPieces {
		packID, _ := uuid.Parse(p.PackID)
		s.PersonaPieces[i] = protocol.PersonaPiece{
			PieceID: p.PieceID, PieceType: personaPieceType(p.PieceType), PackID: packID,
			Default: p.Default, ProductID: p.ProductID,
		}
	}
	s.PieceTintColours = make([]protocol.PersonaPieceTintColour, len(data.PieceTintColours))
	for i, p := range data.PieceTintColours {
		s.PieceTintColours[i].PieceType = p.PieceType
		for j, c := range p.Colours {
			s.PieceTintColours[i].Colours[j] = parseSkinColour(c, true)
		}
	}
}

// SkinColor is RGB; piece tints are ARGB, including the transparent "#0".
// Keep the login parser's existing tolerance for malformed optional metadata.
func parseSkinColour(s string, argb bool) color.RGBA {
	n, err := strconv.ParseUint(strings.TrimPrefix(s, "#"), 16, 32)
	if err != nil {
		return color.RGBA{}
	}
	c := color.RGBA{R: uint8(n >> 16), G: uint8(n >> 8), B: uint8(n), A: 255}
	if argb {
		c.A = uint8(n >> 24)
	}
	return c
}

func personaPieceType(name string) uint32 {
	for t, n := range personaPieceNames {
		if name == n {
			return uint32(t)
		}
	}
	return protocol.PieceTypeUnsupported
}

var personaPieceNames = [...]string{
	protocol.PieceTypeUnknown:       "persona_unknown",
	protocol.PieceTypeSkeleton:      "persona_skeleton",
	protocol.PieceTypeBody:          "persona_body",
	protocol.PieceTypeSkin:          "persona_skin",
	protocol.PieceTypeBottom:        "persona_bottom",
	protocol.PieceTypeFeet:          "persona_feet",
	protocol.PieceTypeDress:         "persona_dress",
	protocol.PieceTypeTop:           "persona_top",
	protocol.PieceTypeHighPants:     "persona_high_pants",
	protocol.PieceTypeHands:         "persona_hand",
	protocol.PieceTypeOuterwear:     "persona_outerwear",
	protocol.PieceTypeFacialHair:    "persona_facial_hair",
	protocol.PieceTypeMouth:         "persona_mouth",
	protocol.PieceTypeEyes:          "persona_eyes",
	protocol.PieceTypeHair:          "persona_hair",
	protocol.PieceTypeHood:          "persona_hood",
	protocol.PieceTypeBack:          "persona_back",
	protocol.PieceTypeFaceAccessory: "persona_face_accessory",
	protocol.PieceTypeHead:          "persona_head",
	protocol.PieceTypeLegs:          "persona_legs",
	protocol.PieceTypeLeftLeg:       "persona_left_leg",
	protocol.PieceTypeRightLeg:      "persona_right_leg",
	protocol.PieceTypeArms:          "persona_arms",
	protocol.PieceTypeLeftArm:       "persona_left_arm",
	protocol.PieceTypeRightArm:      "persona_right_arm",
	protocol.PieceTypeCapes:         "persona_capes",
	protocol.PieceTypeClassicSkin:   "persona_classic_skin",
	protocol.PieceTypeEmote:         "persona_emote",
	protocol.PieceTypeUnsupported:   "unsupported",
}
