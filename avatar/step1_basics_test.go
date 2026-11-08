package avatar

import "testing"

func TestStep1_MeIsFilledIn(t *testing.T) {
	if Name == "" {
		t.Error("set Name in avatar/me.go")
	}
}

func TestStep1_Move(t *testing.T) {
	tests := []struct {
		name         string
		dx, dy       float64
		wantX, wantY float64
	}{
		{"stand still", 0, 0, 100, 100},
		{"right", 1, 0, 100 + Speed, 100},
		{"left", -1, 0, 100 - Speed, 100},
		{"down", 0, 1, 100, 100 + Speed},
		{"up", 0, -1, 100, 100 - Speed},
	}
	for _, tt := range tests {
		x, y := Move(100, 100, tt.dx, tt.dy)
		if x != tt.wantX || y != tt.wantY {
			t.Errorf("%s: Move(100, 100, %v, %v) = %v, %v; want %v, %v",
				tt.name, tt.dx, tt.dy, x, y, tt.wantX, tt.wantY)
		}
	}
}
