package avatar

import (
	"strings"
	"time"
)

// Step 6: pointers vs values.
// Check: make check

// SetStatus sets the status shown under your name, as in "/status lunch":
// s without the spaces around it, cut to MaxStatusLen characters.
// An empty s clears the status.
func (a Avatar) SetStatus(s string) {
	s = strings.TrimSpace(s)
	if r := []rune(s); len(r) > MaxStatusLen {
		s = string(r[:MaxStatusLen])
	}
	a.Status = s
}

// Cooldown allows an action at most once per Every.
type Cooldown struct {
	Every time.Duration
	last  time.Time
}

// Try reports whether the action may run at now. When it may, Try records
// now as the last use. The first call always returns true.
func (c *Cooldown) Try(now time.Time) bool {
	return true // TODO
}
