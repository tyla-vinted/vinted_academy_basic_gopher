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
	tests := []struct {
		in, want string
		ok       bool
	}{
		{"/wave", "wave", true},
		{"/confetti", "confetti", true},
		{"  /dance ", "dance", true},
		{"/fly", "", false},
		{"wave", "", false},
		{"/", "", false},
		{"hello", "", false},
	}
	for _, tt := range tests {
		got, ok := EmoteFor(tt.in)
		if got != tt.want || ok != tt.ok {
			t.Errorf("EmoteFor(%q) = %q, %v; want %q, %v", tt.in, got, ok, tt.want, tt.ok)
		}
	}
}
