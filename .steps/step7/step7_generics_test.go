package avatar

import (
	"slices"
	"testing"
)

func TestStep7_Clamp(t *testing.T) {
	if got := Clamp(15, 0, 10); got != 10 {
		t.Errorf("Clamp(15, 0, 10) = %v, want 10", got)
	}
	if got := Clamp(-2.5, 0.0, 1.0); got != 0 {
		t.Errorf("Clamp(-2.5, 0, 1) = %v, want 0", got)
	}
	if got := Clamp("m", "a", "z"); got != "m" {
		t.Errorf(`Clamp("m", "a", "z") = %q, want "m"`, got)
	}
}

func TestStep7_Filter(t *testing.T) {
	even := Filter([]int{1, 2, 3, 4, 5, 6}, func(n int) bool { return n%2 == 0 })
	if want := []int{2, 4, 6}; !slices.Equal(even, want) {
		t.Errorf("even numbers = %v, want %v", even, want)
	}

	people := []Avatar{
		{Name: "Ada", Pos: Pos{0, 0}},
		{Name: "Linus", Pos: Pos{900, 900}},
		{Name: "Grace", Pos: Pos{10, 10}},
	}
	near := Filter(people, func(a Avatar) bool { return a.Pos.DistanceTo(Pos{}) < 100 })
	if len(near) != 2 || near[0].Name != "Ada" || near[1].Name != "Grace" {
		t.Errorf("nearby = %+v, want Ada and Grace", near)
	}

	if got := Filter([]string{"a"}, func(string) bool { return false }); len(got) != 0 {
		t.Errorf("Filter with no match = %q, want empty", got)
	}
}

func TestStep7_Ring(t *testing.T) {
	r := NewRing[string](3)
	if r.Len() != 0 || len(r.Items()) != 0 {
		t.Fatalf("new ring: Len() = %d, Items() = %q; want empty", r.Len(), r.Items())
	}

	r.Push("a")
	r.Push("b")
	if want := []string{"a", "b"}; !slices.Equal(r.Items(), want) {
		t.Errorf("after a, b: Items() = %q, want %q", r.Items(), want)
	}

	r.Push("c")
	r.Push("d")
	r.Push("e")
	if want := []string{"c", "d", "e"}; !slices.Equal(r.Items(), want) {
		t.Errorf("after a..e: Items() = %q, want %q", r.Items(), want)
	}
	if r.Len() != 3 {
		t.Errorf("Len() = %d, want 3", r.Len())
	}

	items := r.Items()
	if len(items) == 0 {
		return
	}
	items[0] = "changed"
	if r.Items()[0] != "c" {
		t.Error("changing the slice from Items() changed the ring; return a new slice")
	}
}

func TestStep7_RingOfStructs(t *testing.T) {
	type line struct{ from, text string }
	r := NewRing[line](2)
	r.Push(line{"Ada", "hi"})
	r.Push(line{"Linus", "hello"})
	r.Push(line{"Grace", "hey"})
	got := r.Items()
	if len(got) != 2 || got[0].from != "Linus" || got[1].from != "Grace" {
		t.Errorf("Items() = %+v, want Linus then Grace", got)
	}
}
