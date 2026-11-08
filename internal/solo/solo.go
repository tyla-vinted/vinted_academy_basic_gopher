// Package solo stands in for the lobby when no server is given: it sends
// your own chat, emotes and status back to you, so you can see them. A
// dummy can stand in the room, so you can try high-fives.
package solo

import (
	"math"

	"arena/internal/protocol"
)

const (
	you   = "you"
	dummy = "dummy"
)

// Where the arena asks NewAvatar to put the dummy.
const (
	DummyX = protocol.RoomW * 0.7
	DummyY = protocol.RoomH / 2
)

var (
	dummyArt = []string{` [o_o] `, `  /|\  `, `   |   `, `  / \  `}
	replies  = []string{"Up top!", "Nice one!", "Again!"}
)

// Dummy is an avatar that stands still in the solo room.
type Dummy struct {
	Name string
	X, Y float64
}

type Conn struct {
	events  chan protocol.ServerMsg
	x, y    float64
	dummy   *Dummy
	replied int
}

// New returns a solo connection and the welcome message for it. x and y are
// where you start. When d is not nil, the dummy stands in the room.
func New(x, y float64, d *Dummy) (*Conn, protocol.ServerMsg) {
	c := &Conn{events: make(chan protocol.ServerMsg, 64), x: x, y: y, dummy: d}
	w := protocol.ServerMsg{Type: protocol.TypeWelcome, You: you}
	if d != nil {
		w.Players = []protocol.PlayerInfo{{
			ID:      dummy,
			Profile: protocol.Profile{Name: d.Name},
			Look:    protocol.Look{Art: dummyArt, ArtLeft: dummyArt, Color: 0xf0f0f0ff},
			X:       d.X,
			Y:       d.Y,
		}}
	}
	return c, w
}

func (c *Conn) Events() <-chan protocol.ServerMsg { return c.events }

func (c *Conn) Send(m protocol.ClientMsg) {
	switch m.Type {
	case protocol.TypeMove:
		c.x, c.y = m.X, m.Y
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
	case protocol.TypeHighFive:
		c.highFive(m.To)
	}
}

// highFive applies the same distance rule as the lobby server.
func (c *Conn) highFive(to string) {
	if c.dummy == nil || to != dummy {
		c.push(protocol.ServerMsg{Type: protocol.TypeError, Code: protocol.CodeBadRequest, Text: "nobody to high-five"})
		return
	}
	if math.Hypot(c.x-c.dummy.X, c.y-c.dummy.Y) > protocol.HighFiveRange+24 {
		c.push(protocol.ServerMsg{Type: protocol.TypeError, Code: protocol.CodeTooFar, Text: "too far away from " + c.dummy.Name})
		return
	}
	c.push(protocol.ServerMsg{Type: protocol.TypeHighFived, From: you, To: dummy})
	c.push(protocol.ServerMsg{Type: protocol.TypeSaid, From: dummy, Text: replies[c.replied%len(replies)]})
	c.replied++
}

func (c *Conn) Close() {}

func (c *Conn) push(m protocol.ServerMsg) {
	select {
	case c.events <- m:
	default:
	}
}
