// Package protocol defines the messages that the arena client and the lobby
// server exchange over a WebSocket, as JSON.
package protocol

import (
	"fmt"
	"image/color"
	"slices"
	"strings"
	"unicode/utf8"
)

// Version is the protocol version. The server and the client each have a
// copy of this package: change Version in both copies when the messages change.
const Version = 2

const (
	RoomW = 1024
	RoomH = 640

	MaxPlayers    = 40
	MaxNameLen    = 16
	MaxMessageLen = 80
	MaxStatusLen  = 20
	MaxArtRows    = 5
	MaxArtCols    = 10

	// HighFiveRange is the largest distance in pixels between two avatar
	// centres that still allows a high-five.
	HighFiveRange = 100
)

// Emotes is the set of emotes that the server accepts.
var Emotes = []string{"dance", "confetti", "heart"}

func IsEmote(name string) bool {
	return slices.Contains(Emotes, name)
}

type Profile struct {
	Name string `json:"name"`
}

// Look is how an avatar appears: ASCII art facing right, the same art
// facing left, and one colour as packed RGBA (0xRRGGBBAA).
type Look struct {
	Art     []string `json:"art,omitempty"`
	ArtLeft []string `json:"art_left,omitempty"`
	Color   uint32   `json:"color,omitempty"`
}

type PlayerInfo struct {
	ID      string  `json:"id"`
	Profile Profile `json:"profile"`
	Look    Look    `json:"look"`
	Status  string  `json:"status,omitempty"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	Left    bool    `json:"left,omitempty"`
}

type Position struct {
	ID   string  `json:"id"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Left bool    `json:"left,omitempty"`
}

// Client to server message types.
const (
	TypeJoin     = "join"
	TypeMove     = "move"
	TypeSay      = "say"
	TypeEmote    = "emote"
	TypeStatus   = "status"
	TypeHighFive = "highfive"
)

// Server to client message types.
const (
	TypeWelcome   = "welcome"
	TypeJoined    = "joined"
	TypeLeft      = "left"
	TypeState     = "state"
	TypeSaid      = "said"
	TypeEmoted    = "emoted"
	TypeStatusOf  = "status_of"
	TypeHighFived = "highfived"
	TypeAnnounce  = "announce"
	TypeError     = "error"
)

// Error codes in a TypeError message.
const (
	CodeNameTaken      = "name_taken"
	CodeLobbyFull      = "lobby_full"
	CodeInvalidProfile = "invalid_profile"
	CodeBadRequest     = "bad_request"
	CodeRateLimited    = "rate_limited"
	CodeMuted          = "muted"
	CodeTooFar         = "too_far"
	CodeKicked         = "kicked"
	CodeOutdated       = "outdated"
)

type ClientMsg struct {
	Type    string   `json:"type"`
	Version int      `json:"version,omitempty"`
	Profile *Profile `json:"profile,omitempty"`
	Look    *Look    `json:"look,omitempty"`
	X       float64  `json:"x,omitempty"`
	Y       float64  `json:"y,omitempty"`
	Left    bool     `json:"left,omitempty"`
	Text    string   `json:"text,omitempty"`
	Emote   string   `json:"emote,omitempty"`
	Status  string   `json:"status,omitempty"`
	To      string   `json:"to,omitempty"`
}

type ServerMsg struct {
	Type      string       `json:"type"`
	You       string       `json:"you,omitempty"`
	Players   []PlayerInfo `json:"players,omitempty"`
	Player    *PlayerInfo  `json:"player,omitempty"`
	Positions []Position   `json:"positions,omitempty"`
	ID        string       `json:"id,omitempty"`
	From      string       `json:"from,omitempty"`
	To        string       `json:"to,omitempty"`
	Text      string       `json:"text,omitempty"`
	Emote     string       `json:"emote,omitempty"`
	Status    string       `json:"status,omitempty"`
	Code      string       `json:"code,omitempty"`
	Field     string       `json:"field,omitempty"`
	Reason    string       `json:"reason,omitempty"`
}

// FieldError tells which profile field the server rejected and why.
type FieldError struct {
	Field  string
	Reason string
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Reason)
}

// ValidateJoin applies the server rules to a profile and a look.
// The client cannot skip these rules.
func ValidateJoin(p Profile, l Look) *FieldError {
	name := strings.TrimSpace(p.Name)
	switch {
	case name == "":
		return &FieldError{"name", "must not be empty"}
	case utf8.RuneCountInString(name) > MaxNameLen:
		return &FieldError{"name", fmt.Sprintf("must be at most %d characters", MaxNameLen)}
	case !printable(name):
		return &FieldError{"name", "must contain only printable characters"}
	}
	for _, art := range [][]string{l.Art, l.ArtLeft} {
		if len(art) > MaxArtRows {
			return &FieldError{"art", fmt.Sprintf("must have at most %d rows", MaxArtRows)}
		}
		for _, row := range art {
			if len(row) > MaxArtCols {
				return &FieldError{"art", fmt.Sprintf("rows must have at most %d characters", MaxArtCols)}
			}
			for _, r := range row {
				if r < ' ' || r > '~' {
					return &FieldError{"art", "must contain only printable ASCII characters"}
				}
			}
		}
	}
	return nil
}

func Pack(c color.RGBA) uint32 {
	return uint32(c.R)<<24 | uint32(c.G)<<16 | uint32(c.B)<<8 | uint32(c.A)
}

func Unpack(v uint32) color.RGBA {
	return color.RGBA{uint8(v >> 24), uint8(v >> 16), uint8(v >> 8), uint8(v)}
}

// CleanText trims s and reports whether it is printable and not empty.
func CleanText(s string, max int) (string, bool) {
	s = strings.TrimSpace(s)
	if s == "" || utf8.RuneCountInString(s) > max || !printable(s) {
		return "", false
	}
	return s, true
}

func printable(s string) bool {
	for _, r := range s {
		if r < 0x20 || r == 0x7f {
			return false
		}
	}
	return true
}
