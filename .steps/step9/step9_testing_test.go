package avatar

import "testing"

func TestStep9_Greet(t *testing.T) {
	me := Avatar{Name: "Ada"}

	tests := []struct {
		name  string
		other Avatar
		want  string
	}{
		{"same name", Avatar{Name: "Ada"}, "Hey, we are both Ada!"},
		// TODO: one row for each other rule in the Greet comment,
		// then the edge cases. Two of the edge cases find bugs.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Greet(me, tt.other); got != tt.want {
				t.Errorf("Greet(%q, %q) = %q, want %q", me.Name, tt.other.Name, got, tt.want)
			}
		})
	}

	t.Fatal("add your rows to the table, then delete this line")
}
