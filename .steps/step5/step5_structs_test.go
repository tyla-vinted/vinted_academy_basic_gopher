package avatar

import "testing"

func TestStep5_PosAdd(t *testing.T) {
	p := Pos{10, 20}
	if got, want := p.Add(5, -5), (Pos{15, 15}); got != want {
		t.Errorf("Pos{10, 20}.Add(5, -5) = %v, want %v", got, want)
	}
	if p != (Pos{10, 20}) {
		t.Errorf("Add changed p to %v; it must return a new Pos", p)
	}
}

func TestStep5_DistanceTo(t *testing.T) {
	if got := (Pos{0, 0}).DistanceTo(Pos{3, 4}); got != 5 {
		t.Errorf("distance (0,0)→(3,4) = %v, want 5", got)
	}
	if got := (Pos{3, 4}).DistanceTo(Pos{0, 0}); got != 5 {
		t.Errorf("distance (3,4)→(0,0) = %v, want 5", got)
	}
}

func TestStep5_CanHighFive(t *testing.T) {
	a := Avatar{Pos: Pos{100, 100}}
	tests := []struct {
		at   Pos
		want bool
	}{
		{Pos{100, 100}, true},
		{Pos{100 + HighFiveRange, 100}, true},
		{Pos{100 + HighFiveRange + 1, 100}, false},
		{Pos{500, 500}, false},
	}
	for _, tt := range tests {
		b := Avatar{Pos: tt.at}
		if got := a.CanHighFive(b); got != tt.want {
			t.Errorf("CanHighFive from %v to %v = %v, want %v", a.Pos, tt.at, got, tt.want)
		}
	}
}
