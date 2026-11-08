package avatar

// Step 2: flow control and enums: walking.
// Check: make check

type Direction int

const (
	Stay Direction = iota
	Up
	Down
	Left
	Right
)

// DirectionFor returns the direction for a key.
func DirectionFor(key string) Direction {
	return Stay // TODO: you stand still until you write this
}
