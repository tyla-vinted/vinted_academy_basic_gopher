package avatar

// Step 7: generics: the same code for any type.
// Check: make check

// Reverse returns a new slice with the elements of xs in reverse order. Your
// ReverseString from step 3 only works on runes; Reverse works on any slice.
// When it works, you can use it in ReverseString.
func Reverse[T any](xs []T) []T {
	return xs // TODO
}

// Filter returns the elements of xs for which keep returns true, in order.
// For example, Filter(nums, func(n int) bool { return n > 0 }) keeps the
// positive numbers. The arena lists who you can high-five with
// Filter(others, me.CanHighFive).
func Filter[T any](xs []T, keep func(T) bool) []T {
	return nil // TODO
}

// Last returns the last n elements of xs, or all of them when xs is shorter.
// The arena shows Last(chat, 8) as the chat history.
func Last[T any](xs []T, n int) []T {
	return nil // TODO
}
