package avatar

import "testing"

func TestStep1_MeIsFilledIn(t *testing.T) {
	if Name == "" {
		t.Error("set Name in avatar/me.go")
	}
}

func TestStep1_Move(t *testing.T) {
	tests := []struct {
		key          string
		wantX, wantY float64
	}{
		{"w", 100, 100 - Speed},
		{"up", 100, 100 - Speed},
		{"s", 100, 100 + Speed},
		{"down", 100, 100 + Speed},
		{"a", 100 - Speed, 100},
		{"left", 100 - Speed, 100},
		{"d", 100 + Speed, 100},
		{"right", 100 + Speed, 100},
		{"q", 100, 100},
		{"", 100, 100},
	}
	for _, tt := range tests {
		x, y := Move(100, 100, tt.key)
		if x != tt.wantX || y != tt.wantY {
			t.Errorf("Move(100, 100, %q) = %v, %v; want %v, %v", tt.key, x, y, tt.wantX, tt.wantY)
		}
	}
}
