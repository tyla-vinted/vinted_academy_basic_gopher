package avatar

import "image/color"

// Step 2: if, for, switch.
// Run: go test ./avatar -run Step2 -v

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

// ColorFor returns the colour for a name, in any case ("Cyan" and "cyan"
// are the same): red, orange, yellow, green, blue, cyan, purple, pink,
// white, grey. For any other name, ok is false.
func ColorFor(name string) (c color.RGBA, ok bool) {
	// TODO: strings.ToLower, then a switch.
	return color.RGBA{}, false
}

// NormalizeArt makes art safe to draw. It returns new rows, where:
//   - each character that is not printable ASCII (' ' to '~') is '?'
//   - every row is padded with spaces to the width of the widest row
//
// It does not change rows.
func NormalizeArt(rows []string) []string {
	// TODO: find the widest row first (count runes, not bytes).
	return nil
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

type Command int

const (
	CommandSay Command = iota
	CommandEmote
	CommandStatus
	CommandHighFive
	CommandUnknown
)

// ParseCommand reads one line from the chat box. Trim spaces first.
//
//	"hello all"       → CommandSay,      "hello all"
//	"/wave"           → CommandEmote,    "wave"      (any name in Emotes)
//	"/status lunch"   → CommandStatus,   "lunch"     ("/status" alone → "")
//	"/highfive"       → CommandHighFive, ""
//	"/fly"            → CommandUnknown,  "fly"
func ParseCommand(input string) (Command, string) {
	// TODO: strings.TrimSpace, strings.HasPrefix and strings.Cut help here.
	// slices.Contains(Emotes, word) tells you if word is an emote.
	return CommandSay, ""
}
