package avatar

// Step 6: errors. The sentinel errors are in provided.go.
// Check: make check

// ValidationError tells which field in avatar/me.go is wrong and why.
type ValidationError struct {
	Field  string
	Reason string
}

// Error returns "<Field>: <Reason>", for example "Name: must not be empty".
func (e *ValidationError) Error() string {
	return "" // TODO
}

// CheckMessage checks a chat message before the arena sends it.
//   - empty after trimming spaces → ErrEmptyMessage
//   - more than MaxMessageLen characters → an error that wraps
//     ErrMessageTooLong, with the text "<n> characters, max 80: message is too long"
//   - otherwise nil
func CheckMessage(msg string) error {
	return nil // TODO
}

// CheckHighFive returns an error that wraps ErrTooFar when me cannot
// high-five other, with the text "high-five with <name>: too far away".
// Use your CanHighFive from step 3.
func CheckHighFive(me, other Avatar) error {
	return nil // TODO
}

// ValidateName returns a *ValidationError for the first rule that name breaks:
//   - empty after trimming spaces → Field "Name", Reason "must not be empty"
//   - longer than MaxNameLen      → Field "Name", Reason "must be at most 16 characters"
//
// Otherwise it returns nil.
func ValidateName(name string) error {
	return nil // TODO
}

// JoinHint turns a join error into advice for the terminal. The error can be
// wrapped any number of times.
//   - ErrNameTaken       → "Someone has this name already. Change Name in avatar/me.go."
//   - ErrLobbyFull       → "The lobby is full. Ask the host."
//   - *ValidationError   → "Fix <Field> in avatar/me.go: <Reason>."
//   - anything else      → ""
func JoinHint(err error) string {
	return "" // TODO
}
