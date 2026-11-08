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
