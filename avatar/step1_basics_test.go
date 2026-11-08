package avatar

import (
	"slices"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestStep1_Name(t *testing.T) {
	name := Name()
	if strings.TrimSpace(name) == "" {
		t.Error("Name returns an empty name")
	}
	if n := utf8.RuneCountInString(name); n > 16 {
		t.Errorf("Name returns %q: %d characters; want at most 16", name, n)
	}
}

func TestStep1_Art(t *testing.T) {
	art := Art()
	if len(art) == 0 || len(art) > 5 {
		t.Errorf("Art returns %d rows; want 1-5", len(art))
	}
	for i, row := range art {
		if n := utf8.RuneCountInString(row); n > 10 {
			t.Errorf("Art()[%d] has %d characters; want at most 10", i, n)
		}
		for _, r := range row {
			if r < ' ' || r > '~' {
				t.Errorf("Art()[%d] has %q; use only printable ASCII", i, r)
			}
		}
	}
}

func TestStep1_Color(t *testing.T) {
	known := []string{"red", "orange", "yellow", "green", "blue", "cyan", "purple", "pink", "white", "grey"}
	if c := Color(); !slices.Contains(known, strings.ToLower(c)) {
		t.Errorf("Color returns %q; want one of %q", c, known)
	}
}

func TestStep1_Me(t *testing.T) {
	name, color, art := Me()
	if name == "" || color == "" || art == nil {
		t.Fatalf("Me returns %q, %q, %q; want your name, your colour and your art", name, color, art)
	}
	if name != Name() {
		t.Errorf("Me returns the name %q; want %q, the same as Name", name, Name())
	}
	if color != Color() {
		t.Errorf("Me returns the colour %q; want %q, the same as Color", color, Color())
	}
	if !slices.Equal(art, Art()) {
		t.Errorf("Me returns other art than Art")
	}
}
