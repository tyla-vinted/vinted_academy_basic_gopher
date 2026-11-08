package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"arena/avatar"
	"arena/internal/engine"
	"arena/internal/netclient"
	"arena/internal/protocol"
	"arena/internal/solo"
)

func main() {
	server := flag.String("server", "", "lobby address, e.g. https://abc.ngrok-free.app (empty: solo mode)")
	flag.Parse()
	online := strings.TrimSpace(*server) != ""

	name, colorName, rawArt := avatar.Me()
	name = strings.TrimSpace(name)
	if online {
		if err := avatar.ValidateName(name); err != nil {
			fail(err)
		}
	}

	x := 120 + rand.Float64()*(protocol.RoomW-240)
	y := 160 + rand.Float64()*(protocol.RoomH-280)
	a := avatar.NewAvatar(name, x, y)
	if a == (avatar.Avatar{}) {
		a = avatar.Avatar{Name: name, Pos: avatar.Pos{X: x, Y: y}} // NewAvatar is not written yet
	}
	me := &a
	art := engine.NormalizeArt(rawArt)
	c := engine.ColorFor(colorName)
	look := protocol.Look{Art: art, ArtLeft: avatar.Mirror(art), Color: protocol.Pack(c)}

	var conn engine.Conn
	var welcome protocol.ServerMsg
	title := "Go Arena"
	if name != "" {
		title += " · " + name
	}
	if online {
		nc, w, err := netclient.Dial(context.Background(), *server, protocol.ClientMsg{
			Type:    protocol.TypeJoin,
			Profile: &protocol.Profile{Name: name},
			Look:    &look,
			X:       me.Pos.X,
			Y:       me.Pos.Y,
		})
		if err != nil {
			fail(err)
		}
		defer nc.Close()
		conn, welcome = nc, w
	} else {
		// The learner's NewAvatar (step 5) creates the dummy; until it is
		// written, it returns an empty Avatar and the room stays empty.
		var d *solo.Dummy
		if a := avatar.NewAvatar("Dummy", solo.DummyX, solo.DummyY); a != (avatar.Avatar{}) {
			d = &solo.Dummy{Name: a.Name, X: a.Pos.X, Y: a.Pos.Y}
		}
		conn, welcome = solo.New(me.Pos.X, me.Pos.Y, d)
		title += " (solo)"
	}

	ebiten.SetWindowTitle(title)
	ebiten.SetWindowSize(protocol.RoomW, protocol.RoomH)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	step9, _ := filepath.Glob("avatar/step9_*.go")
	g := engine.New(engine.Config{Conn: conn, Welcome: welcome, Me: me, Look: look, Greetings: len(step9) > 0})
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "Could not join:", err)
	if hint := avatar.JoinHint(err); hint != "" {
		fmt.Fprintln(os.Stderr, "Hint:", hint)
	}
	os.Exit(1)
}
