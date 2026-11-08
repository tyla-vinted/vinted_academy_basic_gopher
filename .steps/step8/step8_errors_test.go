package avatar

import (
	"errors"
	"testing"
	"time"
)

func TestStep8_checkDistance(t *testing.T) {
	me := Avatar{Pos: Pos{0, 0}}
	if err := checkDistance(me, Avatar{Pos: Pos{50, 0}}); err != nil {
		t.Errorf("close by: got %v, want nil", err)
	}
	if err := checkDistance(me, Avatar{Pos: Pos{500, 0}}); err != ErrTooFar {
		t.Errorf("far away: got %v, want ErrTooFar", err)
	}
}

func TestStep8_checkCooldown(t *testing.T) {
	start := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	cd := Cooldown{Every: time.Second}
	if err := checkCooldown(&cd, start); err != nil {
		t.Errorf("first try: got %v, want nil", err)
	}
	if err := checkCooldown(&cd, start.Add(100*time.Millisecond)); err != ErrOnCooldown {
		t.Errorf("second try after 100ms: got %v, want ErrOnCooldown", err)
	}
}

func TestStep8_HighFive(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	me := Avatar{Name: "Me", Pos: Pos{0, 0}}
	near := Avatar{Name: "Ada", Pos: Pos{50, 0}}
	far := Avatar{Name: "Linus", Pos: Pos{500, 0}}

	tests := []struct {
		name     string
		other    Avatar
		wantText string // "" means no error
		wantErr  error  // the error that must be inside, for errors.Is
	}{
		{"close by", near, "", nil},
		{"too far", far, "high-five with Linus: too far away", ErrTooFar},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cd := Cooldown{Every: time.Second}
			err := HighFive(me, tt.other, &cd, now)
			if tt.wantText == "" {
				if err != nil {
					t.Fatalf("got %v, want nil", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("got nil, want %q", tt.wantText)
			}
			if err.Error() != tt.wantText {
				t.Errorf("got %q, want %q", err.Error(), tt.wantText)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("errors.Is(err, %v) is false: wrap the error with %%w", tt.wantErr)
			}
		})
	}

	t.Run("on cooldown", func(t *testing.T) {
		cd := Cooldown{Every: time.Second}
		HighFive(me, near, &cd, now)
		err := HighFive(me, near, &cd, now.Add(100*time.Millisecond))
		if !errors.Is(err, ErrOnCooldown) {
			t.Fatalf("second high-five after 100ms: got %v, want ErrOnCooldown inside", err)
		}
		if want := "high-five with Ada: on cooldown"; err.Error() != want {
			t.Errorf("got %q, want %q", err.Error(), want)
		}
	})

	t.Run("a high-five that is too far does not use the cooldown", func(t *testing.T) {
		cd := Cooldown{Every: time.Second}
		HighFive(me, far, &cd, now)
		if err := HighFive(me, near, &cd, now); err != nil {
			t.Errorf("after a refused high-five: got %v, want nil (check the distance first)", err)
		}
	})
}
