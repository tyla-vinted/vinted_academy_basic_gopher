// Package solo stands in for the lobby when no server is given: it sends
// your own chat, emotes and status back to you, so you can see them.
package solo

import (
	"arena/internal/protocol"
)

const you = "you"

type Conn struct {
	events chan protocol.ServerMsg
}

// New returns a solo connection and the welcome message for it.
func New() (*Conn, protocol.ServerMsg) {
	return &Conn{events: make(chan protocol.ServerMsg, 64)}, protocol.ServerMsg{Type: protocol.TypeWelcome, You: you}
}

func (c *Conn) Events() <-chan protocol.ServerMsg { return c.events }

func (c *Conn) Send(m protocol.ClientMsg) {
	switch m.Type {
	case protocol.TypeSay:
		if text, ok := protocol.CleanText(m.Text, protocol.MaxMessageLen); ok {
			c.push(protocol.ServerMsg{Type: protocol.TypeSaid, From: you, Text: text})
		}
	case protocol.TypeEmote:
		if protocol.IsEmote(m.Emote) {
			c.push(protocol.ServerMsg{Type: protocol.TypeEmoted, From: you, Emote: m.Emote})
		}
	case protocol.TypeStatus:
		c.push(protocol.ServerMsg{Type: protocol.TypeStatusOf, From: you, Status: m.Status})
	}
}

func (c *Conn) Close() {}

func (c *Conn) push(m protocol.ServerMsg) {
	select {
	case c.events <- m:
	default:
	}
}
