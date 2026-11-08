package avatar

import "image/color"

// Step 3: maps and for loops: how you look.
// Check: make check

var (
	Red    = color.RGBA{235, 80, 70, 255}
	Orange = color.RGBA{245, 150, 50, 255}
	Yellow = color.RGBA{245, 210, 60, 255}
	Green  = color.RGBA{100, 200, 100, 255}
	Blue   = color.RGBA{90, 150, 245, 255}
	Cyan   = color.RGBA{70, 210, 220, 255}
	Purple = color.RGBA{180, 120, 240, 255}
	Pink   = color.RGBA{245, 130, 190, 255}
	White  = color.RGBA{240, 240, 240, 255}
	Grey   = color.RGBA{150, 155, 165, 255}
)

// Color is the colour of your art: red, orange, yellow, green, blue, cyan,
// purple, pink, white or grey.
const Color = "cyan"

// palette maps each colour name to its colour.
var palette = map[string]color.RGBA{
	"red": Red,
	// TODO: the other colours
}

// ColorFor returns the colour for a name, in any case ("Cyan" and "cyan" are
// the same). For a name that is not in palette, ok is false.
func ColorFor(name string) (c color.RGBA, ok bool) {
	// TODO: look the name up in palette. A map lookup can return two values.
	return color.RGBA{}, false
}

// MirrorRune returns r as it looks in a mirror:
//
//	( ↔ )    [ ↔ ]    { ↔ }    < ↔ >    / ↔ \    b ↔ d    p ↔ q
//
// Every other rune stays the same.
func MirrorRune(r rune) rune {
	// TODO
	return r
}

// Mirror returns art flipped left to right: each row reversed, and each
// rune mirrored with MirrorRune. The arena draws it when you walk left.
func Mirror(rows []string) []string {
	// TODO
	return nil
}
