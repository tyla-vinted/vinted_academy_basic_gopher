package avatar

// Step 3: structs and methods.
// Run: go test ./avatar -run Step3 -v

type Pos struct {
	X, Y float64
}

type Avatar struct {
	Name   string
	Pos    Pos
	Status string
}

// Add returns p moved by dx and dy.
func (p Pos) Add(dx, dy float64) Pos {
	return Pos{p.X + dx, p.Y + dy}
}

// DistanceTo returns the straight-line distance from p to q.
// Hint: math.Hypot.
func (p Pos) DistanceTo(q Pos) float64 {
	return 0 // TODO
}

// CanHighFive reports whether other is at most HighFiveRange away from a.
func (a Avatar) CanHighFive(other Avatar) bool {
	return false // TODO
}

// CardLines returns the text of the card that others see when they walk
// close to you: your name, then "Status: <status>" when you have a status.
//
//	["Ada", "Status: lunch"]    or ["Ada"] with no status
func (a Avatar) CardLines() []string {
	return nil // TODO
}
