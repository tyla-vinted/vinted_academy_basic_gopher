package avatar

// Step 5: structs and methods: high-fives.
// Check: make check

type Pos struct {
	X, Y float64
}

type Avatar struct {
	Name   string
	Pos    Pos
	Status string
}

// NewAvatar returns an avatar with the given name, standing at x, y.
func NewAvatar(name string, x, y float64) Avatar {
	return Avatar{} // TODO
}

// DistanceTo returns the straight-line distance from p to q.
// Hint: math.Hypot, https://pkg.go.dev/math#Hypot
func (p Pos) DistanceTo(q Pos) float64 {
	return 0 // TODO
}

// CanHighFive reports whether other is at most 100.0 away
func (a Avatar) CanHighFive(other Avatar) bool {
	return false // TODO
}
