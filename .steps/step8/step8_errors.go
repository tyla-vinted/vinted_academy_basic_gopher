package avatar

import (
	"errors"
	"time"
)

// Step 8: errors: why a high-five failed.
// Check: make check

// Errors that a high-five check can return.
var (
	ErrTooFar     = errors.New("too far away")
	ErrOnCooldown = errors.New("on cooldown")
)

// checkDistance returns ErrTooFar when your CanHighFive (step 5) says no.
func checkDistance(me, other Avatar) error {
	return nil // TODO
}

// checkCooldown returns ErrOnCooldown when your cooldown (step 6) says no.
func checkCooldown(cd *Cooldown, now time.Time) error {
	return nil // TODO
}

// HighFive checks a high-five before the arena sends it. It runs
// checkDistance, then checkCooldown. At the first error it stops and returns
// that error wrapped with fmt.Errorf and %w:
// "high-five with <other's name>: <the error>". Otherwise it returns nil.
func HighFive(me, other Avatar, cd *Cooldown, now time.Time) error {
	return nil // TODO
}
