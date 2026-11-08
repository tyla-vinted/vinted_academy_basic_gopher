package avatar

// This file is you. Change the constants and the art, then run the arena
// again to see the result.

const (
	Name = "" // TODO: your name, 1-16 characters. It must be unique in the lobby.

	// Color is the colour of your art: red, orange, yellow, green, blue,
	// cyan, purple, pink, white or grey.
	Color = "cyan"
)

// Art is your avatar as ASCII art, facing right. At most 5 rows of 10
// characters. Use only the characters on a US keyboard. The arena mirrors it
// when you walk left. In backticks, \ needs no escape.
var Art = []string{
	`  ,---.  `,
	` ( o o)> `,
	`  \ - /  `,
	`  /|_|\  `,
	`   / \   `,
}
