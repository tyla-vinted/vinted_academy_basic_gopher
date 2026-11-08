package avatar

// Step 7: generics: the same code for any type.
// Check: make check

// Reverse returns a new slice with the items in reverse order. Your
// ReverseString from step 3 only works on runes; Reverse works on any slice.
//
// TODO: when Reverse works, use it in your ReverseString (step 3), then
// run make check: the step 3 tests show that nothing broke.
func Reverse[T any](items []T) []T {
	return items // TODO
}

// Filter returns the items for which keep returns true, in order.
// For example, Filter(nums, func(n int) bool { return n > 0 }) keeps the
// positive numbers.
func Filter[T any](items []T, keep func(T) bool) []T {
	return nil // TODO
}

// Last returns the last n items, or all of them when there are fewer.
func Last[T any](items []T, n int) []T {
	return nil // TODO
}
