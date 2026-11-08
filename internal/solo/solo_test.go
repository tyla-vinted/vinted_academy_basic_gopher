package solo

import (
	"testing"

	"arena/internal/protocol"
)

func TestEcho(t *testing.T) {
	c, w := New(0, 0, nil)
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

func TestDummyHighFive(t *testing.T) {
	c, w := New(0, 0, &Dummy{Name: "Dummy", X: DummyX, Y: DummyY})
	if len(w.Players) != 1 || w.Players[0].ID != dummy || w.Players[0].Profile.Name != "Dummy" {
		t.Fatalf("welcome players = %+v, want the dummy", w.Players)
	}
	c.Send(protocol.ClientMsg{Type: protocol.TypeHighFive, To: dummy})
	if m := <-c.Events(); m.Type != protocol.TypeError || m.Code != protocol.CodeTooFar {
		t.Errorf("far high-five: got %+v, want too_far", m)
	}
	c.Send(protocol.ClientMsg{Type: protocol.TypeMove, X: DummyX + 50, Y: DummyY})
	c.Send(protocol.ClientMsg{Type: protocol.TypeHighFive, To: dummy})
	if m := <-c.Events(); m.Type != protocol.TypeHighFived || m.To != dummy {
		t.Errorf("close high-five: got %+v, want highfived", m)
	}
	if m := <-c.Events(); m.Type != protocol.TypeSaid || m.From != dummy {
		t.Errorf("after high-five: got %+v, want the dummy to say something", m)
	}
}

func TestNoDummy(t *testing.T) {
	c, w := New(0, 0, nil)
	if len(w.Players) != 0 {
		t.Errorf("welcome players = %+v, want none", w.Players)
	}
	c.Send(protocol.ClientMsg{Type: protocol.TypeHighFive, To: dummy})
	if m := <-c.Events(); m.Type != protocol.TypeError {
		t.Errorf("high-five without dummy: got %+v, want an error", m)
	}
}
