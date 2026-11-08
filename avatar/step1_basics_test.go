package avatar

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestStep1_Me(t *testing.T) {
	name, art := Me()
	if strings.TrimSpace(name) == "" {
		t.Error("Me returns an empty name")
	}
	if n := utf8.RuneCountInString(name); n > MaxNameLen {
		t.Errorf("name %q has %d characters; want at most %d", name, n, MaxNameLen)
	}
	if len(art) == 0 || len(art) > MaxArtRows {
		t.Errorf("art has %d rows; want 1-%d", len(art), MaxArtRows)
	}
	for i, row := range art {
		if n := utf8.RuneCountInString(row); n > MaxArtCols {
			t.Errorf("art[%d] has %d characters; want at most %d", i, n, MaxArtCols)
		}
		for _, r := range row {
			if r < ' ' || r > '~' {
				t.Errorf("art[%d] has %q; use only printable ASCII", i, r)
			}
		}
	}
}

func TestStep1_Direction(t *testing.T) {
	tests := []struct {
		key    string
		dx, dy int
	}{
		{"w", 0, -1},
		{"up", 0, -1},
		{"s", 0, 1},
		{"down", 0, 1},
		{"a", -1, 0},
		{"left", -1, 0},
		{"d", 1, 0},
		{"right", 1, 0},
		{"q", 0, 0},
		{"", 0, 0},
	}
	for _, tt := range tests {
		dx, dy := Direction(tt.key)
		if dx != tt.dx || dy != tt.dy {
			t.Errorf("Direction(%q) = %d, %d; want %d, %d", tt.key, dx, dy, tt.dx, tt.dy)
		}
	}
}
