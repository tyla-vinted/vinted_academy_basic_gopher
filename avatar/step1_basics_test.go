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
