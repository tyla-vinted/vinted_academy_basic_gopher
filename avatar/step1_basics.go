package avatar

// Step 1: variables, constants, functions with multiple returns.
// Check: make check

const Speed = 3.0 // pixels per frame

// Move returns your position after one key press. The arena calls it once
// per frame for each key that you hold, with the key name:
//
//	"w" or "up"      → Speed pixels up    (y gets smaller)
//	"s" or "down"    → Speed pixels down  (y gets bigger)
//	"a" or "left"    → Speed pixels left  (x gets smaller)
//	"d" or "right"   → Speed pixels right (x gets bigger)
//	any other key    → no change
//
// Hint: a switch case can list several values: case "w", "up":
func Move(x, y float64, key string) (newX, newY float64) {
	return x, y // TODO: you stand still until you write this
}
