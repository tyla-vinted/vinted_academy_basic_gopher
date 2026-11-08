package avatar

import "testing"

func TestStep4_Message(t *testing.T) {
	tests := []struct {
		in, want string
		ok       bool
	}{
		{"hello", "hello", true},
		{"  hi all  ", "hi all", true},
		{"", "", false},
		{"   ", "", false},
	}
	for _, tt := range tests {
		got, ok := Message(tt.in)
		if got != tt.want || ok != tt.ok {
			t.Errorf("Message(%q) = %q, %v; want %q, %v", tt.in, got, ok, tt.want, tt.ok)
		}
	}
}

func TestStep4_EmoteFor(t *testing.T) {
	names := map[Emote]string{Wave: "Wave", Dance: "Dance", Confetti: "Confetti", Heart: "Heart", Think: "Think"}
	found := []struct {
		in   string
		want Emote
	}{
		{"/wave", Wave},
		{"/dance", Dance},
		{"/confetti", Confetti},
		{"/heart", Heart},
		{"/think", Think},
		{"  /dance ", Dance},
	}
	for _, tt := range found {
		got, ok := EmoteFor(tt.in)
		if !ok || got != tt.want {
			t.Errorf("EmoteFor(%q) = %s, %v; want %s, true", tt.in, names[got], ok, names[tt.want])
		}
	}
	for _, in := range []string{"/fly", "wave", "/", "hello"} {
		if _, ok := EmoteFor(in); ok {
			t.Errorf("EmoteFor(%q): ok is true; want false", in)
		}
	}
}
