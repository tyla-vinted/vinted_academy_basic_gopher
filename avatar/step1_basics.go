package avatar

// Step 1: functions, single and multiple return values.
// Check: make check

// Name returns your name: 1-16 characters.
func Name() string {
	return "Evaldas" // TODO
}

// Art returns your gopher avatar as ASCII art: at most 5 rows of 10 characters.
// You can start from this gopher:
//
//	 n____n
//	( O  O )
//	(  ''  )/
//	(______)
//	  d  b
func Art() []string {
	return []string{
		`  n____n  `,
		` ( O  O ) `,
		` (  ''  )/`,
		` (______) `,
		`   d  b   `,
	}
}

// Color returns the colour of your art: red, orange, yellow, green, blue,
// cyan, purple, pink, white or grey.
func Color() string {
	return "cyan" // TODO
}

// Me returns your name, your colour and your art.
func Me() (name, color string, art []string) {
	return Name(), Color(), Art() // TODO
}
