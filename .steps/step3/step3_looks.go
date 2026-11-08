package avatar

// Step 3: for loops: facing left and right.
// Check: make check

// MirrorRune returns r as it looks in a mirror:
//
//	( ↔ )    [ ↔ ]    { ↔ }    < ↔ >    / ↔ \    b ↔ d    p ↔ q
//
// Every other rune stays the same.
func MirrorRune(r rune) rune {
	// TODO
	return r
}

// ReverseString returns s with its runes in reverse order: "abc" → "cba".
func ReverseString(s string) string {
	// TODO
	return s
}

// Mirror returns art flipped left to right: each row reversed, and each
// rune mirrored with MirrorRune. The arena draws it when you walk left.
func Mirror(rows []string) []string {
	// TODO: use ReverseString and MirrorRune.
	return nil
}
