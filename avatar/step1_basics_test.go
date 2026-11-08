package avatar

import "testing"

func TestStep1_MeIsFilledIn(t *testing.T) {
	if Name == "" {
		t.Error("set Name in avatar/me.go")
	}
}

func TestStep1_MoveSpeed(t *testing.T) {
	if got := MoveSpeed(false); got != WalkSpeed {
		t.Errorf("MoveSpeed(false) = %v, want %v", got, WalkSpeed)
	}
	if got, want := MoveSpeed(true), WalkSpeed*RunMultiplier; got != want {
		t.Errorf("MoveSpeed(true) = %v, want %v (look closely at the := inside the if)", got, want)
	}
}

func TestStep1_Move(t *testing.T) {
	tests := []struct {
		name         string
		dx, dy       float64
		running      bool
		wantX, wantY float64
	}{
		{"stand still", 0, 0, false, 100, 100},
		{"walk right", 1, 0, false, 100 + WalkSpeed, 100},
		{"walk left", -1, 0, false, 100 - WalkSpeed, 100},
		{"walk down", 0, 1, false, 100, 100 + WalkSpeed},
		{"walk up", 0, -1, false, 100, 100 - WalkSpeed},
		{"run right", 1, 0, true, 100 + WalkSpeed*RunMultiplier, 100},
	}
	for _, tt := range tests {
		x, y := Move(100, 100, tt.dx, tt.dy, tt.running)
		if x != tt.wantX || y != tt.wantY {
			t.Errorf("%s: Move(100, 100, %v, %v, %v) = %v, %v; want %v, %v",
				tt.name, tt.dx, tt.dy, tt.running, x, y, tt.wantX, tt.wantY)
		}
	}
}

func TestStep1_Truncate(t *testing.T) {
	tests := []struct {
		in      string
		max     int
		want    string
		wantCut bool
	}{
		{"hello", 10, "hello", false},
		{"hello", 5, "hello", false},
		{"hello world", 8, "hello...", true},
		{"", 5, "", false},
		{"ąžuolas ąžuolas", 10, "ąžuolas...", true},
	}
	for _, tt := range tests {
		got, cut := Truncate(tt.in, tt.max)
		if got != tt.want || cut != tt.wantCut {
			t.Errorf("Truncate(%q, %d) = %q, %v; want %q, %v", tt.in, tt.max, got, cut, tt.want, tt.wantCut)
		}
	}
}
