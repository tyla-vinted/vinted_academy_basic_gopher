package avatar

import (
	"errors"
	"log"
	"os"
)

// Nothing to do in this file. It holds the values that the steps share.

const (
	MaxNameLen    = 16
	MaxMessageLen = 80
	MaxStatusLen  = 20

	// HighFiveRange is the largest distance in pixels between two avatars
	// that still allows a high-five.
	HighFiveRange = 100

	// Your Art must fit in this box.
	MaxArtRows = 5
	MaxArtCols = 10
)

// Emotes lists the emotes that the lobby accepts, as in "/wave".
var Emotes = []string{"wave", "dance", "confetti", "heart", "think"}

// Errors that the arena client returns when the lobby rejects a join.
var (
	ErrNameTaken = errors.New("name is taken")
	ErrLobbyFull = errors.New("lobby is full")
)

// Errors for step 6.
var (
	ErrEmptyMessage   = errors.New("message is empty")
	ErrMessageTooLong = errors.New("message is too long")
	ErrTooFar         = errors.New("too far away")
	ErrOnCooldown     = errors.New("on cooldown")
)

var debug = os.Getenv("ARENA_DEBUG") != ""

func debugf(format string, args ...any) {
	if debug {
		log.Printf(format, args...)
	}
}
