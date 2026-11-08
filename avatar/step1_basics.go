package avatar

// Step 1: variables, constants, functions with multiple returns.
// Check: make check

const Speed = 3.0 // pixels per frame

// Move returns your position after one frame. dx and dy are the direction
// of the keys you hold, from -1 to 1: left and up are negative, right and
// down are positive, 0 means no key. Move Speed pixels in that direction.
//
//	Move(100, 100, 1, 0)  → 103, 100   (right)
//	Move(100, 100, 0, -1) → 100, 97    (up)
func Move(x, y, dx, dy float64) (newX, newY float64) {
	return x, y // TODO: you stand still until you write this
}
