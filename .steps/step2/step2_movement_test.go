package avatar

import "testing"

func TestStep2_Direction(t *testing.T) {
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
