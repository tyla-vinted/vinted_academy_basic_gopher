package avatar

import "testing"

func TestStep2_DirectionFor(t *testing.T) {
	names := map[Direction]string{Stay: "Stay", Up: "Up", Down: "Down", Left: "Left", Right: "Right"}
	tests := []struct {
		key  string
		want Direction
	}{
		{"w", Up},
		{"up", Up},
		{"s", Down},
		{"down", Down},
		{"a", Left},
		{"left", Left},
		{"d", Right},
		{"right", Right},
		{"q", Stay},
		{"", Stay},
	}
	for _, tt := range tests {
		if got := DirectionFor(tt.key); got != tt.want {
			t.Errorf("DirectionFor(%q) = %s, want %s", tt.key, names[got], names[tt.want])
		}
	}
}
