package avatar

import (
	"strings"
	"testing"
	"time"
)

func TestStep6_SetStatus(t *testing.T) {
	var me Avatar
	me.SetStatus("  lunch  ")
	if me.Status != "lunch" {
		t.Fatalf(`after SetStatus("  lunch  "), Status = %q; want "lunch"`, me.Status)
	}
	me.SetStatus(strings.Repeat("z", 30))
	if want := strings.Repeat("z", 20); me.Status != want {
		t.Errorf("after SetStatus of 30 characters, Status = %q; want %q", me.Status, want)
	}
	me.SetStatus("")
	if me.Status != "" {
		t.Errorf(`after SetStatus(""), Status = %q; want ""`, me.Status)
	}
}

func TestStep6_Cooldown(t *testing.T) {
	start := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	c := Cooldown{Every: 2 * time.Second}

	steps := []struct {
		at   time.Duration
		want bool
	}{
		{0, true},
		{500 * time.Millisecond, false},
		{1999 * time.Millisecond, false},
		{2 * time.Second, true},
		{3 * time.Second, false},
		{10 * time.Second, true},
	}
	for _, s := range steps {
		if got := c.Try(start.Add(s.at)); got != s.want {
			t.Errorf("Try at +%v = %v, want %v", s.at, got, s.want)
		}
	}
}
