package avatar

// Step 1: variables, constants, shadowing, functions with multiple returns.
// Check: make check

const (
	WalkSpeed     = 3.0 // pixels per frame
	RunMultiplier = 2

	// MaxBubbleLen is the longest text that fits in a speech bubble.
	MaxBubbleLen = 40
)

// MoveSpeed returns how many pixels per frame your avatar moves.
// Hold Shift in the arena to run.
func MoveSpeed(running bool) float64 {
	speed := WalkSpeed
	if running {
		speed := WalkSpeed * RunMultiplier
		debugf("running at %.1f px/frame", speed)
	}
	return speed
}

// Move returns your position after one frame. dx and dy are the direction
// of the keys you hold, from -1 to 1: left and up are negative, right and
// down are positive, 0 means no key. Move MoveSpeed(running) pixels in that
// direction.
//
//	Move(100, 100, 1, 0, false) → 103, 100   (walking right)
//	Move(100, 100, 0, -1, true) → 100, 94    (running up)
func Move(x, y, dx, dy float64, running bool) (newX, newY float64) {
	return x, y // TODO: you stand still until you write this
}

// Truncate shortens s to at most max characters (runes, not bytes). When it
// cuts s, the last 3 characters of the result are "..." and the second
// return value is true.
//
//	Truncate("hello", 10)      → "hello", false
//	Truncate("hello world", 8) → "hello...", true
func Truncate(s string, max int) (string, bool) {
	// TODO: utf8.RuneCountInString counts characters; []rune(s) lets you cut by character.
	return s, false
}
