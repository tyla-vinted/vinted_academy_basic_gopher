package avatar

// Step 1: functions, single and multiple return values.
// Check: make check

// Name returns your name: 1-16 characters.
func Name() string {
	return "" // TODO
}

// Art returns your avatar as ASCII art: at most 5 rows of 10 characters,
// facing right, only the characters on a US keyboard. Write each row in
// backticks, so that \ needs no escape. It is a gopher now: draw your own
// if you like.
func Art() []string {
	return []string{
		`  n____n  `,
		` ( O  O ) `,
		` (  ''  )/`,
		` (______) `,
		`   d  b   `,
	}
}

// Me returns your name and your art. Use Name and Art.
func Me() (name string, art []string) {
	return "", nil // TODO
}
