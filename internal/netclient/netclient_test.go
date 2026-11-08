package netclient

import "testing"

func TestWebSocketURL(t *testing.T) {
	tests := map[string]string{
		"https://abc.ngrok-free.app":  "wss://abc.ngrok-free.app/ws",
		"abc.ngrok-free.app":          "wss://abc.ngrok-free.app/ws",
		"https://abc.ngrok-free.app/": "wss://abc.ngrok-free.app/ws",
		"http://localhost:8080":       "ws://localhost:8080/ws",
		"192.168.1.5:8080":            "ws://192.168.1.5:8080/ws",
		"ws://127.0.0.1:9000/ws":      "ws://127.0.0.1:9000/ws",
	}
	for in, want := range tests {
		if got, err := WebSocketURL(in); err != nil || got != want {
			t.Errorf("WebSocketURL(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
}
