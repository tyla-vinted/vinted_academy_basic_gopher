package avatar

import "testing"

func TestStep4_ParseCommand(t *testing.T) {
	tests := []struct {
		in      string
		want    Command
		wantArg string
	}{
		{"hello all", CommandSay, "hello all"},
		{"  hi  ", CommandSay, "hi"},
		{"/wave", CommandEmote, "wave"},
		{"/confetti", CommandEmote, "confetti"},
		{"/status lunch break", CommandStatus, "lunch break"},
		{"/status", CommandStatus, ""},
		{"/highfive", CommandHighFive, ""},
		{"/fly", CommandUnknown, "fly"},
	}
	for _, tt := range tests {
		got, arg := ParseCommand(tt.in)
		if got != tt.want || arg != tt.wantArg {
			t.Errorf("ParseCommand(%q) = %v, %q; want %v, %q", tt.in, got, arg, tt.want, tt.wantArg)
		}
	}
}
