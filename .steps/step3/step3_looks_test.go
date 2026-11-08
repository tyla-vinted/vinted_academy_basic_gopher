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

func TestStep3_ReverseString(t *testing.T) {
	tests := []struct{ in, want string }{
		{"abc", "cba"},
		{"a", "a"},
		{"", ""},
		{"héllo", "olléh"},
	}
	for _, tt := range tests {
		if got := ReverseString(tt.in); got != tt.want {
			t.Errorf("ReverseString(%q) = %q, want %q", tt.in, got, tt.want)
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
