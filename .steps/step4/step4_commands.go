package avatar

// Step 4: if and strings: chat commands.
// Check: make check

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
