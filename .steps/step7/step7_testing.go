package avatar

import "strings"

// Step 7: testing. Greet has no tests yet, and it has bugs.
// Write a table-driven test in step7_testing_test.go, find the bugs, fix them.
// Check: make check

// Greet returns the line that you see when you walk up to other.
//   - same name, in any case → "Hey, we are both <name>!"
//   - other has a status     → "Hi <name> (<status>)!"
//   - otherwise              → "Nice to meet you, <name>!"
func Greet(me, other Avatar) string {
	switch {
	case strings.Contains(other.Name, me.Name):
		return "Hey, we are both " + other.Name + "!"
	case other.Status != "":
		return "Hi " + other.Name + " (" + other.Status + ")!"
	default:
		return "Nice to meet you, " + other.Name + "!"
	}
}
