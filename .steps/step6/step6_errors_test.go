package avatar

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestStep6_ValidationErrorMessage(t *testing.T) {
	var err error = &ValidationError{Field: "Name", Reason: "must not be empty"}
	if got, want := err.Error(), "Name: must not be empty"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestStep6_CheckMessage(t *testing.T) {
	if err := CheckMessage("hello"); err != nil {
		t.Errorf(`CheckMessage("hello") = %v, want nil`, err)
	}
	if err := CheckMessage("   "); !errors.Is(err, ErrEmptyMessage) {
		t.Errorf("CheckMessage(spaces) = %v, want ErrEmptyMessage", err)
	}

	long := strings.Repeat("a", 81)
	err := CheckMessage(long)
	if !errors.Is(err, ErrMessageTooLong) {
		t.Fatalf("CheckMessage(81 chars) = %v, want an error that wraps ErrMessageTooLong (use %%w)", err)
	}
	if got, want := err.Error(), "81 characters, max 80: message is too long"; got != want {
		t.Errorf("message = %q, want %q", got, want)
	}
}

func TestStep6_CheckHighFive(t *testing.T) {
	me := Avatar{Pos: Pos{0, 0}}
	near := Avatar{Name: "Ada", Pos: Pos{10, 0}}
	far := Avatar{Name: "Linus", Pos: Pos{500, 0}}

	if err := CheckHighFive(me, near); err != nil {
		t.Errorf("CheckHighFive(near) = %v, want nil", err)
	}
	err := CheckHighFive(me, far)
	if !errors.Is(err, ErrTooFar) {
		t.Fatalf("CheckHighFive(far) = %v, want an error that wraps ErrTooFar", err)
	}
	if got, want := err.Error(), "high-five with Linus: too far away"; got != want {
		t.Errorf("message = %q, want %q", got, want)
	}
}

func TestStep6_ValidateName(t *testing.T) {
	tests := []struct {
		name      string
		in        string
		wantField string // "" means no error
	}{
		{"valid", "Ada", ""},
		{"16 characters", strings.Repeat("a", 16), ""},
		{"empty", "", "Name"},
		{"spaces only", "   ", "Name"},
		{"17 characters", strings.Repeat("a", 17), "Name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.in)
			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("got %v, want nil", err)
				}
				return
			}
			var ve *ValidationError
			if !errors.As(err, &ve) {
				t.Fatalf("got %v (%T), want a *ValidationError", err, err)
			}
			if ve.Field != tt.wantField {
				t.Errorf("Field = %q, want %q", ve.Field, tt.wantField)
			}
			if ve.Reason == "" {
				t.Error("Reason is empty")
			}
		})
	}
}

func TestStep6_JoinHint(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"name taken", fmt.Errorf("join wss://x: %w", ErrNameTaken), "Someone has this name already. Change Name in avatar/me.go."},
		{"wrapped twice", fmt.Errorf("arena: %w", fmt.Errorf("join: %w", ErrNameTaken)), "Someone has this name already. Change Name in avatar/me.go."},
		{"lobby full", fmt.Errorf("join: %w", ErrLobbyFull), "The lobby is full. Ask the host."},
		{"validation", fmt.Errorf("join: %w", &ValidationError{Field: "Name", Reason: "must not be empty"}), "Fix Name in avatar/me.go: must not be empty."},
		{"other", errors.New("connection refused"), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinHint(tt.err); got != tt.want {
				t.Errorf("JoinHint(%v) = %q, want %q", tt.err, got, tt.want)
			}
		})
	}
}
