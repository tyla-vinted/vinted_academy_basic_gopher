package avatar

// Step 2: flow control and enums: walking.
// Check: make check

// Direction is where you walk. It is an enum: a type with a named constant
// for each value. iota counts up from 0, so Stay is 0, Up is 1, and so on.
type Direction int

const (
	Stay Direction = iota
	Up
	Down
	Left
	Right
)

// DirectionFor returns the direction for a key. The arena calls it for each
// key that you hold:
//
//	"w" or "up"    → Up
//	"s" or "down"  → Down
//	"a" or "left"  → Left
//	"d" or "right" → Right
//	any other key  → Stay
func DirectionFor(key string) Direction {
	return Stay // TODO: you stand still until you write this
}
