// Package netclient connects to the lobby. It hides the goroutines and
// channels from the rest of the client: concurrency is a topic for a later
// workshop.
package netclient

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"

	"arena/internal/protocol"
)

// Errors that Dial returns when the lobby rejects a join.
var (
	ErrNameTaken = errors.New("name is taken")
	ErrLobbyFull = errors.New("lobby is full")
)

// FieldError says which part of your avatar the lobby rejected, and why.
type FieldError struct {
	Field  string
	Reason string
}

func (e *FieldError) Error() string { return e.Field + ": " + e.Reason }

type Conn struct {
	c      *websocket.Conn
	events chan protocol.ServerMsg
	out    chan protocol.ClientMsg
	ctx    context.Context
	cancel context.CancelFunc

	mu  sync.Mutex
	err error
}

// Dial joins the lobby at addr. addr can be a full URL ("https://x.ngrok-free.app",
// "ws://127.0.0.1:8080/ws") or a bare host ("x.ngrok-free.app", "192.168.1.5:8080").
func Dial(ctx context.Context, addr string, join protocol.ClientMsg) (*Conn, protocol.ServerMsg, error) {
	u, err := WebSocketURL(addr)
	if err != nil {
		return nil, protocol.ServerMsg{}, err
	}
	dctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	c, _, err := websocket.Dial(dctx, u, &websocket.DialOptions{
		CompressionMode: websocket.CompressionContextTakeover,
		HTTPHeader: http.Header{
			"ngrok-skip-browser-warning": {"1"},
			"User-Agent":                 {"go-arena"},
		},
	})
	if err != nil {
		return nil, protocol.ServerMsg{}, fmt.Errorf("connect to %s: %w", u, err)
	}
	c.SetReadLimit(1 << 20)
	join.Version = protocol.Version

	if err := wsjson.Write(dctx, c, join); err != nil {
		c.CloseNow()
		return nil, protocol.ServerMsg{}, fmt.Errorf("join %s: %w", u, err)
	}
	var first protocol.ServerMsg
	if err := wsjson.Read(dctx, c, &first); err != nil {
		c.CloseNow()
		return nil, protocol.ServerMsg{}, fmt.Errorf("join %s: %w", u, err)
	}
	if first.Type == protocol.TypeError {
		c.CloseNow()
		return nil, protocol.ServerMsg{}, fmt.Errorf("join %s: %w", u, joinError(first))
	}
	if first.Type != protocol.TypeWelcome {
		c.CloseNow()
		return nil, protocol.ServerMsg{}, fmt.Errorf("join %s: unexpected %q message", u, first.Type)
	}

	cctx, ccancel := context.WithCancel(context.Background())
	conn := &Conn{
		c:      c,
		events: make(chan protocol.ServerMsg, 1024),
		out:    make(chan protocol.ClientMsg, 64),
		ctx:    cctx,
		cancel: ccancel,
	}
	go conn.readLoop()
	go conn.writeLoop()
	return conn, first, nil
}

func joinError(m protocol.ServerMsg) error {
	switch m.Code {
	case protocol.CodeNameTaken:
		return ErrNameTaken
	case protocol.CodeLobbyFull:
		return ErrLobbyFull
	case protocol.CodeOutdated:
		return fmt.Errorf("%s (the server has newer code)", m.Reason)
	case protocol.CodeInvalidProfile:
		return &FieldError{Field: fieldName(m.Field), Reason: m.Reason}
	default:
		return fmt.Errorf("%s: %s", m.Code, m.Reason)
	}
}

// fieldName maps a server field name to the name that learners know.
func fieldName(f string) string {
	switch f {
	case "name":
		return "Name"
	case "art":
		return "Art"
	default:
		return f
	}
}

// WebSocketURL turns what a person types into a WebSocket URL with the /ws path.
func WebSocketURL(addr string) (string, error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return "", errors.New("empty server address")
	}
	if !strings.Contains(addr, "://") {
		host := addr
		if i := strings.IndexAny(host, "/"); i >= 0 {
			host = host[:i]
		}
		if isLocal(host) {
			addr = "ws://" + addr
		} else {
			addr = "wss://" + addr
		}
	}
	u, err := url.Parse(addr)
	if err != nil {
		return "", fmt.Errorf("server address %q: %w", addr, err)
	}
	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	case "ws", "wss":
	default:
		return "", fmt.Errorf("server address %q: unknown scheme %q", addr, u.Scheme)
	}
	if u.Path == "" || u.Path == "/" {
		u.Path = "/ws"
	}
	return u.String(), nil
}

func isLocal(hostport string) bool {
	host := hostport
	if h, _, err := net.SplitHostPort(hostport); err == nil {
		host = h
	}
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && (ip.IsLoopback() || ip.IsPrivate())
}

// Events returns the messages from the server. The channel closes when the
// connection ends; Err then tells why.
func (c *Conn) Events() <-chan protocol.ServerMsg { return c.events }

// Send queues m. When the queue is full, Send drops m.
func (c *Conn) Send(m protocol.ClientMsg) {
	select {
	case c.out <- m:
	default:
	}
}

func (c *Conn) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.err
}

func (c *Conn) Close() {
	c.cancel()
	c.c.Close(websocket.StatusNormalClosure, "bye")
}

func (c *Conn) fail(err error) {
	c.mu.Lock()
	if c.err == nil {
		c.err = err
	}
	c.mu.Unlock()
	c.cancel()
}

func (c *Conn) readLoop() {
	defer close(c.events)
	for {
		var m protocol.ServerMsg
		if err := wsjson.Read(c.ctx, c.c, &m); err != nil {
			c.fail(err)
			return
		}
		select {
		case c.events <- m:
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Conn) writeLoop() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case m := <-c.out:
			wctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
			err := wsjson.Write(wctx, c.c, m)
			cancel()
			if err != nil {
				c.fail(err)
				return
			}
		}
	}
}
