package avatar

import "errors"

// Nothing to do in this file. It holds the values that the steps share.

const (
	MaxNameLen    = 16
	MaxMessageLen = 80
	MaxStatusLen  = 20

	// Your Art must fit in this box.
	MaxArtRows = 5
	MaxArtCols = 10
)

// Errors that the arena client returns when the lobby rejects a join.
var (
	ErrNameTaken = errors.New("name is taken")
	ErrLobbyFull = errors.New("lobby is full")
)

// Errors for step 8.
var (
	ErrEmptyMessage   = errors.New("message is empty")
	ErrMessageTooLong = errors.New("message is too long")
	ErrTooFar         = errors.New("too far away")
	ErrOnCooldown     = errors.New("on cooldown")
)
