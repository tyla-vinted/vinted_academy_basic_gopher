package avatar

// Step 4: if, strings and maps: chat and emotes.
// Check: make check

// Message returns what you typed, without the spaces around it. ok is false
// when the input is empty or only spaces.
func Message(input string) (text string, ok bool) {
	return "", false // TODO
}

// emotes maps each emote command to the emote that it plays.
var emotes = map[string]string{
	"/wave": "wave",
	// TODO: "/dance", "/confetti", "/heart" and "/think"
}

// EmoteFor returns the emote for a command like "/wave". ok is false when
// the command is not in emotes.
func EmoteFor(input string) (emote string, ok bool) {
	return "", false // TODO
}
