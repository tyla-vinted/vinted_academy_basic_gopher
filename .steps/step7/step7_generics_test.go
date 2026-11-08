package avatar

import (
	"slices"
	"testing"
)

func TestStep7_Reverse(t *testing.T) {
	nums := []int{1, 2, 3}
	if got, want := Reverse(nums), []int{3, 2, 1}; !slices.Equal(got, want) {
		t.Errorf("Reverse(%v) = %v, want %v", nums, got, want)
	}
	if !slices.Equal(nums, []int{1, 2, 3}) {
		t.Errorf("Reverse changed its input to %v; return a new slice", nums)
	}
	words := []string{"a", "b"}
	if got, want := Reverse(words), []string{"b", "a"}; !slices.Equal(got, want) {
		t.Errorf("Reverse(%q) = %q, want %q", words, got, want)
	}
	if got := Reverse([]int{}); len(got) != 0 {
		t.Errorf("Reverse of an empty slice = %v, want empty", got)
	}
}

func TestStep7_Filter(t *testing.T) {
	even := Filter([]int{1, 2, 3, 4, 5, 6}, func(n int) bool { return n%2 == 0 })
	if want := []int{2, 4, 6}; !slices.Equal(even, want) {
		t.Errorf("even numbers = %v, want %v", even, want)
	}

	me := Avatar{Name: "Me", Pos: Pos{0, 0}}
	people := []Avatar{
		{Name: "Ada", Pos: Pos{50, 0}},
		{Name: "Linus", Pos: Pos{900, 900}},
		{Name: "Grace", Pos: Pos{0, 80}},
	}
	near := Filter(people, me.CanHighFive)
	if len(near) != 2 || near[0].Name != "Ada" || near[1].Name != "Grace" {
		t.Errorf("Filter(people, me.CanHighFive) = %+v, want Ada and Grace", near)
	}

	if got := Filter([]string{"a"}, func(string) bool { return false }); len(got) != 0 {
		t.Errorf("Filter with no match = %q, want empty", got)
	}
}

func TestStep7_Last(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	tests := []struct {
		n    int
		want []int
	}{
		{2, []int{4, 5}},
		{5, []int{1, 2, 3, 4, 5}},
		{9, []int{1, 2, 3, 4, 5}},
		{0, []int{}},
	}
	for _, tt := range tests {
		if got := Last(nums, tt.n); !slices.Equal(got, tt.want) {
			t.Errorf("Last(%v, %d) = %v, want %v", nums, tt.n, got, tt.want)
		}
	}
	if got := Last([]string{"hi"}, 3); !slices.Equal(got, []string{"hi"}) {
		t.Errorf(`Last([]string{"hi"}, 3) = %q, want ["hi"]`, got)
	}
}
