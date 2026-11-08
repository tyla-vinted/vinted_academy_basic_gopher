package avatar

// Step 1: functions with multiple return values.
// Check: make check

// Me returns your name (1-16 characters) and your ASCII art: at most 5 rows
// of 10 characters, facing right, only the characters on a US keyboard.
// Write each row in backticks, so that \ needs no escape. You can start
// from this gopher:
//
//	 n____n
//	( O  O )
//	(  ''  )/
//	(______)
//	  d  b
func Me() (name string, art []string) {
	return "", nil // TODO
}

func Direction(key string) (dx, dy int) {
	return 0, 0 // TODO: you stand still until you write this
}
