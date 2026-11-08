package avatar

import (
	"strings"
	"time"
)

// Step 6: pointers: your status and emote cooldown.
// Check: make check

// SetStatus sets the status shown under your name, as in "/status lunch":
//
// TODO: "/status lunch" does nothing. Find out why.
func (a Avatar) SetStatus(s string) {
	maxLen := 20
	s = strings.TrimSpace(s)
	if r := []rune(s); len(r) > maxLen {
		s = string(r[:maxLen])
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
//
// TODO: the cooldown never starts. Find out why.
func (c Cooldown) Try(now time.Time) bool {
	if !c.last.IsZero() && now.Sub(c.last) < c.Every {
		return false
	}
	c.last = now
	return true
}
