package solo

import (
	"testing"

	"arena/internal/protocol"
)

func TestEcho(t *testing.T) {
	c, w := New()
	c.Send(protocol.ClientMsg{Type: protocol.TypeSay, Text: " hi "})
	c.Send(protocol.ClientMsg{Type: protocol.TypeEmote, Emote: "dance"})
	c.Send(protocol.ClientMsg{Type: protocol.TypeEmote, Emote: "fly"})
	c.Send(protocol.ClientMsg{Type: protocol.TypeMove, X: 1})

	if m := <-c.Events(); m.Type != protocol.TypeSaid || m.Text != "hi" || m.From != w.You {
		t.Errorf("first event = %+v, want said hi", m)
	}
	if m := <-c.Events(); m.Type != protocol.TypeEmoted || m.Emote != "dance" {
		t.Errorf("second event = %+v, want emoted dance", m)
	}
	select {
	case m := <-c.Events():
		t.Errorf("unexpected event %+v", m)
	default:
	}
}
