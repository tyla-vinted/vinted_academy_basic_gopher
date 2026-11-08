package avatar

import (
	"image/color"
	"slices"
	"testing"
)

func TestStep3_ColorFor(t *testing.T) {
	tests := []struct {
		name   string
		want   color.RGBA
		wantOK bool
	}{
		{"red", Red, true},
		{"orange", Orange, true},
		{"yellow", Yellow, true},
		{"green", Green, true},
		{"blue", Blue, true},
		{"cyan", Cyan, true},
		{"purple", Purple, true},
		{"pink", Pink, true},
		{"white", White, true},
		{"grey", Grey, true},
		{"Cyan", Cyan, true},
		{"GREEN", Green, true},
		{"magenta", color.RGBA{}, false},
		{"", color.RGBA{}, false},
	}
	for _, tt := range tests {
		got, ok := ColorFor(tt.name)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("ColorFor(%q) = %v, %v; want %v, %v", tt.name, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestStep3_NormalizeArt(t *testing.T) {
	in := []string{
		"ab",
		"a\tb",
		"é",
		"",
	}
	want := []string{
		"ab ",
		"a?b",
		"?  ",
		"   ",
	}
	got := NormalizeArt(in)
	if !slices.Equal(got, want) {
		t.Errorf("NormalizeArt(%q)\n got %q\nwant %q", in, got, want)
	}
	if in[1] != "a\tb" {
		t.Error("NormalizeArt changed its input; return new rows")
	}
}

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

func TestStep3_MyColor(t *testing.T) {
	if _, ok := ColorFor(Color); !ok {
		t.Errorf("Color = %q is not a known colour", Color)
	}
}
