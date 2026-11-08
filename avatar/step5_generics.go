package avatar

import "cmp"

// Step 5: generics.
// Run: go test ./avatar -run Step5 -v

// ClampInt and ClampFloat are the same function twice. Write Clamp once, for
// any ordered type, so that both of them can go.
func ClampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func ClampFloat(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// Clamp returns v limited to the range [lo, hi].
func Clamp[T cmp.Ordered](v, lo, hi T) T {
	return v // TODO
}

// Filter returns the elements of xs for which keep returns true, in order.
// The arena uses it for the "Nearby" list.
func Filter[T any](xs []T, keep func(T) bool) []T {
	return nil // TODO
}

// Ring keeps the last N values that you push. The arena uses a Ring for the
// chat history.
type Ring[T any] struct {
	buf   []T
	start int // index of the oldest value
	size  int // number of values in buf
}

// NewRing returns an empty Ring that keeps up to n values.
func NewRing[T any](n int) *Ring[T] {
	return &Ring[T]{} // TODO: make room for n values
}

// Push adds v. When the ring is full, Push drops the oldest value.
func (r *Ring[T]) Push(v T) {
	// TODO
}

// Len returns the number of values in the ring.
func (r *Ring[T]) Len() int {
	return 0 // TODO
}

// Items returns the values from oldest to newest, in a new slice.
func (r *Ring[T]) Items() []T {
	return nil // TODO
}
