package avatar

// Step 4: if, strings and maps: chat and emotes.
// Check: make check

// Message returns what you typed, without the spaces around it. ok is false
// when the input is empty or only spaces.
func Message(input string) (text string, ok bool) {
	return "", false // TODO
}

// Emote is an animation that everyone around you sees.
type Emote int

const (
	Wave Emote = iota
	Dance
	Confetti
	Heart
	Think
)

// emotes maps each emote command to its emote.
var emotes = map[string]Emote{
	"/wave": Wave,
	// TODO: "/dance", "/confetti", "/heart" and "/think"
}

// EmoteFor returns the emote for a command like "/wave". ok is false when
// the command is not in emotes.
func EmoteFor(input string) (emote Emote, ok bool) {
	return 0, false // TODO
}
