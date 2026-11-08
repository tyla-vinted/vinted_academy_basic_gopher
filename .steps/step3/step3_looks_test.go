package avatar

import (
	"slices"
	"testing"
)

func TestStep3_MirrorRune(t *testing.T) {
	pairs := []struct{ in, want rune }{
		{'(', ')'}, {')', '('},
		{'[', ']'}, {']', '['},
		{'{', '}'}, {'}', '{'},
		{'<', '>'}, {'>', '<'},
		{'/', '\\'}, {'\\', '/'},
		{'b', 'd'}, {'d', 'b'},
		{'p', 'q'}, {'q', 'p'},
		{'o', 'o'}, {'|', '|'}, {' ', ' '},
	}
	for _, p := range pairs {
		if got := MirrorRune(p.in); got != p.want {
			t.Errorf("MirrorRune(%q) = %q, want %q", p.in, got, p.want)
		}
	}
}

func TestStep3_Mirror(t *testing.T) {
	in := []string{
		`(o>`,
		`/| `,
	}
	want := []string{
		`<o)`,
		` |\`,
	}
	if got := Mirror(in); !slices.Equal(got, want) {
		t.Errorf("Mirror(%q) = %q, want %q", in, got, want)
	}
	if in[0] != `(o>` {
		t.Error("Mirror changed its input; return new rows")
	}
}
