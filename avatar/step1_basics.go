package avatar

// Step 1: variables, constants, shadowing, functions with multiple returns.
// Run: go test ./avatar -run Step1 -v

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
