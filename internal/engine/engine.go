// Package engine draws the arena and handles input. It calls the functions
// that learners write in package avatar.
package engine

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand/v2"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/gofont/goregular"

	"arena/avatar"
	"arena/internal/protocol"
)

const (
	artLineH     = 16 // pixels per row of ASCII art
	cardRange    = 140
	nearbyRange  = 200
	bubbleTime   = 6 * time.Second
	emoteTime    = 2500 * time.Millisecond
	toastTime    = 3500 * time.Millisecond
	bannerTime   = 10 * time.Second
	moveInterval = 66 * time.Millisecond
	maxInput     = 120
	maxBubbleLen = 40
	speed        = 3.0 // pixels per frame
	chatLines    = 8
)

var (
	colFloor    = color.RGBA{36, 40, 52, 255}
	colGrid     = color.RGBA{44, 49, 63, 255}
	colPanel    = color.RGBA{16, 18, 24, 210}
	colText     = color.RGBA{230, 232, 238, 255}
	colDim      = color.RGBA{150, 156, 170, 255}
	colAccent   = color.RGBA{255, 205, 70, 255}
	colBubble   = color.RGBA{250, 250, 252, 255}
	colBubbleTx = color.RGBA{25, 25, 30, 255}
	colError    = color.RGBA{235, 90, 80, 255}
	colGood     = color.RGBA{110, 210, 120, 255}
	colNoColor  = color.RGBA{120, 124, 135, 255}
)

// Conn is the connection to the lobby: a real one, or solo mode.
type Conn interface {
	Events() <-chan protocol.ServerMsg
	Send(protocol.ClientMsg)
}

type Config struct {
	Conn    Conn
	Welcome protocol.ServerMsg
	Me      *avatar.Avatar
	Look    protocol.Look
}

type person struct {
	info        protocol.PlayerInfo
	pos         avatar.Pos
	target      avatar.Pos
	left        bool
	bubble      string
	bubbleUntil time.Time
	emote       string
	emoteUntil  time.Time
}

type chatLine struct {
	from, text string
	system     bool
}

type toast struct {
	text  string
	col   color.Color
	until time.Time
}

type particle struct {
	x, y, vx, vy float64
	life, max    int
	col          color.RGBA
	glyph        string
}

type Game struct {
	cfg    Config
	conn   Conn
	you    string
	me     *avatar.Avatar
	people map[string]*person
	chat   *avatar.Ring[chatLine]

	typing bool
	input  []rune

	toasts      []toast
	banner      string
	bannerUntil time.Time
	particles   []particle
	emoteCD     avatar.Cooldown

	meLeft     bool
	lastSent   avatar.Pos
	lastLeft   bool
	lastSentAt time.Time
	lost       bool
	frame      int

	font, small, bold, mono *text.GoTextFace
	monoW                   float64

	shotPath string // ARENA_SCREENSHOT: save one frame as PNG, then quit
	shotDone bool
}

func New(cfg Config) *Game {
	reg, _ := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	bld, _ := text.NewGoTextFaceSource(bytes.NewReader(gobold.TTF))
	mon, _ := text.NewGoTextFaceSource(bytes.NewReader(gomonobold.TTF))
	g := &Game{
		cfg:     cfg,
		conn:    cfg.Conn,
		you:     cfg.Welcome.You,
		me:      cfg.Me,
		people:  map[string]*person{},
		chat:    avatar.NewRing[chatLine](chatLines),
		emoteCD: avatar.Cooldown{Every: 2 * time.Second},
		font:    &text.GoTextFace{Source: reg, Size: 14},
		small:   &text.GoTextFace{Source: reg, Size: 12},
		bold:    &text.GoTextFace{Source: bld, Size: 14},
		mono:    &text.GoTextFace{Source: mon, Size: 15},

		shotPath: os.Getenv("ARENA_SCREENSHOT"),
	}
	g.monoW, _ = text.Measure("M", g.mono, 0)
	if g.chat == nil {
		g.chat = &avatar.Ring[chatLine]{}
	}
	g.lastSent = g.me.Pos
	self := protocol.PlayerInfo{
		ID:      g.you,
		Profile: protocol.Profile{Name: g.me.Name},
		Look:    cfg.Look,
	}
	g.add(self)
	for _, p := range cfg.Welcome.Players {
		g.add(p)
	}
	return g
}

func (g *Game) add(info protocol.PlayerInfo) {
	pos := avatar.Pos{X: info.X, Y: info.Y}
	if info.ID == g.you {
		pos = g.me.Pos
	}
	g.people[info.ID] = &person{info: info, pos: pos, target: pos, left: info.Left}
}

var unknownArt = []string{"[?]", `/|\`, `/ \`}

// art returns the rows to draw for p, facing the way p walks.
func (p *person) art() []string {
	if p.left && len(p.info.Look.ArtLeft) > 0 {
		return p.info.Look.ArtLeft
	}
	if len(p.info.Look.Art) > 0 {
		return p.info.Look.Art
	}
	return unknownArt
}

func (p *person) color() color.RGBA {
	if c := protocol.Unpack(p.info.Look.Color); c.A != 0 {
		return c
	}
	return colNoColor
}

func (g *Game) artSize(rows []string) (w, h float64) {
	cols := 0
	for _, r := range rows {
		cols = max(cols, len([]rune(r)))
	}
	return float64(cols) * g.monoW, float64(len(rows) * artLineH)
}

func (g *Game) Layout(int, int) (int, int) { return protocol.RoomW, protocol.RoomH }

// ---------------------------------------------------------------- update

func (g *Game) Update() error {
	if g.shotDone {
		return ebiten.Termination
	}
	now := time.Now()
	g.frame++
	g.drain(now)
	if g.typing {
		g.updateTyping(now)
	} else {
		g.updateKeys(now)
	}
	for id, p := range g.people {
		if id == g.you {
			p.pos, p.left = g.me.Pos, g.meLeft
			continue
		}
		p.pos.X += (p.target.X - p.pos.X) * 0.25
		p.pos.Y += (p.target.Y - p.pos.Y) * 0.25
	}
	g.updateParticles()
	if !g.lost && (g.me.Pos != g.lastSent || g.meLeft != g.lastLeft) && now.Sub(g.lastSentAt) >= moveInterval {
		g.conn.Send(protocol.ClientMsg{Type: protocol.TypeMove, X: g.me.Pos.X, Y: g.me.Pos.Y, Left: g.meLeft})
		g.lastSent, g.lastLeft, g.lastSentAt = g.me.Pos, g.meLeft, now
	}
	return nil
}

func (g *Game) drain(now time.Time) {
	for {
		select {
		case m, ok := <-g.conn.Events():
			if !ok {
				if !g.lost {
					g.lost = true
					g.showBanner("Disconnected from the lobby. Close the window and join again.", now)
				}
				return
			}
			g.handle(m, now)
		default:
			return
		}
	}
}

func (g *Game) handle(m protocol.ServerMsg, now time.Time) {
	switch m.Type {
	case protocol.TypeState:
		for _, ps := range m.Positions {
			if p, ok := g.people[ps.ID]; ok && ps.ID != g.you {
				p.target = avatar.Pos{X: ps.X, Y: ps.Y}
				p.left = ps.Left
			}
		}
	case protocol.TypeJoined:
		if m.Player != nil {
			g.add(*m.Player)
			g.system(m.Player.Profile.Name + " joined")
		}
	case protocol.TypeLeft:
		if p, ok := g.people[m.ID]; ok {
			g.system(p.info.Profile.Name + " left")
			delete(g.people, m.ID)
		}
	case protocol.TypeSaid:
		p, ok := g.people[m.From]
		if !ok {
			return
		}
		p.bubble, p.bubbleUntil = truncate(m.Text, maxBubbleLen), now.Add(bubbleTime)
		g.chat.Push(chatLine{from: p.info.Profile.Name, text: m.Text})
	case protocol.TypeEmoted:
		p, ok := g.people[m.From]
		if !ok {
			return
		}
		p.emote, p.emoteUntil = m.Emote, now.Add(emoteTime)
		g.emoteParticles(p, m.Emote)
	case protocol.TypeStatusOf:
		if p, ok := g.people[m.From]; ok {
			p.info.Status = m.Status
		}
	case protocol.TypeHighFived:
		a, aok := g.people[m.From]
		b, bok := g.people[m.To]
		if !aok || !bok {
			return
		}
		g.burst((a.pos.X+b.pos.X)/2, (a.pos.Y+b.pos.Y)/2-20, 18, "*", colAccent)
		g.system(a.info.Profile.Name + " high-fived " + b.info.Profile.Name)
		var other *person
		switch g.you {
		case m.From:
			other = b
		case m.To:
			other = a
		}
		if other != nil {
			g.toast("High-five with "+other.info.Profile.Name+"!", colGood, now)
		}
	case protocol.TypeAnnounce:
		g.showBanner(m.Text, now)
		g.chat.Push(chatLine{from: "HOST", text: m.Text, system: true})
	case protocol.TypeError:
		g.toast(m.Text, colError, now)
	}
}

func (g *Game) updateKeys(now time.Time) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) {
		g.typing, g.input = true, g.input[:0]
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		g.typing, g.input = true, []rune{'/'}
		return
	}
	if g.lost {
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		g.highFive(now)
	}
	for i, key := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4, ebiten.Key5} {
		if inpututil.IsKeyJustPressed(key) {
			g.emote(avatar.Emotes[i], now)
		}
	}

	x, y := g.me.Pos.X, g.me.Pos.Y
	for _, k := range moveKeys {
		if ebiten.IsKeyPressed(k.key) {
			switch avatar.DirectionFor(k.name) {
			case avatar.Up:
				y -= speed
			case avatar.Down:
				y += speed
			case avatar.Left:
				x -= speed
			case avatar.Right:
				x += speed
			}
		}
	}
	if x == g.me.Pos.X && y == g.me.Pos.Y {
		return
	}
	if x != g.me.Pos.X {
		g.meLeft = x < g.me.Pos.X
	}
	next := avatar.Pos{X: x, Y: y}
	w, h := g.artSize(g.people[g.you].art())
	next.X = math.Max(w/2, math.Min(protocol.RoomW-w/2, next.X))
	next.Y = math.Max(h/2+40, math.Min(protocol.RoomH-h/2-30, next.Y))
	g.me.Pos = next
}

// moveKeys are the keys that the arena passes to avatar.DirectionFor, by name.
var moveKeys = []struct {
	key  ebiten.Key
	name string
}{
	{ebiten.KeyW, "w"}, {ebiten.KeyA, "a"}, {ebiten.KeyS, "s"}, {ebiten.KeyD, "d"},
	{ebiten.KeyArrowUp, "up"}, {ebiten.KeyArrowDown, "down"},
	{ebiten.KeyArrowLeft, "left"}, {ebiten.KeyArrowRight, "right"},
}

func (g *Game) updateTyping(now time.Time) {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyEscape):
		g.typing = false
		return
	case inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter):
		g.typing = false
		if s := strings.TrimSpace(string(g.input)); s != "" {
			g.submit(s, now)
		}
		return
	}
	if d := inpututil.KeyPressDuration(ebiten.KeyBackspace); len(g.input) > 0 && (d == 1 || (d > 25 && d%3 == 0)) {
		g.input = g.input[:len(g.input)-1]
	}
	for _, r := range ebiten.AppendInputChars(nil) {
		if unicode.IsPrint(r) && len(g.input) < maxInput {
			g.input = append(g.input, r)
		}
	}
}

func (g *Game) submit(line string, now time.Time) {
	if g.lost {
		return
	}
	kind, arg := avatar.ParseCommand(line)
	if kind == avatar.CommandSay && arg == "" {
		arg = line // ParseCommand is not written yet: treat everything as chat
	}
	switch kind {
	case avatar.CommandSay:
		if err := avatar.CheckMessage(arg); err != nil {
			g.toast(err.Error(), colError, now)
			return
		}
		g.conn.Send(protocol.ClientMsg{Type: protocol.TypeSay, Text: arg})
	case avatar.CommandEmote:
		g.emote(arg, now)
	case avatar.CommandStatus:
		g.me.SetStatus(arg)
		g.conn.Send(protocol.ClientMsg{Type: protocol.TypeStatus, Status: g.me.Status})
	case avatar.CommandHighFive:
		g.highFive(now)
	default:
		g.toast(fmt.Sprintf("Unknown command /%s. Try /%s, /status <text>, /highfive.", arg, strings.Join(avatar.Emotes, " /")), colError, now)
	}
}

func (g *Game) emote(name string, now time.Time) {
	if !g.emoteCD.Try(now) {
		g.toast(fmt.Sprintf("/%s: %v", name, avatar.ErrOnCooldown), colError, now)
		return
	}
	g.conn.Send(protocol.ClientMsg{Type: protocol.TypeEmote, Emote: name})
}

func (g *Game) highFive(now time.Time) {
	var best *person
	bestD := math.Inf(1)
	for id, p := range g.people {
		if id == g.you {
			continue
		}
		if d := math.Hypot(p.pos.X-g.me.Pos.X, p.pos.Y-g.me.Pos.Y); d < bestD {
			best, bestD = p, d
		}
	}
	if best == nil || bestD > 3*avatar.HighFiveRange {
		g.toast("Nobody nearby. Walk up to someone first.", colError, now)
		return
	}
	if err := avatar.CheckHighFive(*g.me, g.asAvatar(best)); err != nil {
		g.toast(err.Error(), colError, now)
		return
	}
	g.conn.Send(protocol.ClientMsg{Type: protocol.TypeHighFive, To: best.info.ID})
}

func (g *Game) asAvatar(p *person) avatar.Avatar {
	return avatar.Avatar{
		Name:   p.info.Profile.Name,
		Pos:    p.pos,
		Status: p.info.Status,
	}
}

func (g *Game) toast(s string, c color.Color, now time.Time) {
	if s == "" {
		return
	}
	g.toasts = append(g.toasts, toast{s, c, now.Add(toastTime)})
	if len(g.toasts) > 4 {
		g.toasts = g.toasts[len(g.toasts)-4:]
	}
}

func (g *Game) showBanner(s string, now time.Time) {
	g.banner, g.bannerUntil = s, now.Add(bannerTime)
}

func (g *Game) system(s string) {
	g.chat.Push(chatLine{text: s, system: true})
}

// ---------------------------------------------------------------- particles

var confettiColors = []color.RGBA{
	{220, 50, 47, 255}, {245, 210, 50, 255}, {70, 170, 80, 255},
	{50, 110, 220, 255}, {150, 80, 200, 255}, {240, 120, 180, 255},
}

func (g *Game) emoteParticles(p *person, emote string) {
	switch emote {
	case "confetti":
		for range 40 {
			g.particles = append(g.particles, particle{
				x: p.pos.X, y: p.pos.Y - 30,
				vx: rand.Float64()*6 - 3, vy: -rand.Float64()*5 - 1,
				life: 70, max: 70, col: confettiColors[rand.IntN(len(confettiColors))],
			})
		}
	case "heart":
		g.burst(p.pos.X, p.pos.Y-30, 6, "<3", color.RGBA{240, 90, 130, 255})
	}
}

func (g *Game) burst(x, y float64, n int, glyph string, c color.RGBA) {
	for range n {
		a := rand.Float64() * 2 * math.Pi
		s := 1 + rand.Float64()*2
		g.particles = append(g.particles, particle{
			x: x, y: y, vx: math.Cos(a) * s, vy: math.Sin(a)*s - 1,
			life: 50, max: 50, col: c, glyph: glyph,
		})
	}
}

func (g *Game) updateParticles() {
	alive := g.particles[:0]
	for _, p := range g.particles {
		p.x += p.vx
		p.y += p.vy
		if p.glyph == "" {
			p.vy += 0.15
		} else {
			p.vx *= 0.95
			p.vy *= 0.95
		}
		p.life--
		if p.life > 0 {
			alive = append(alive, p)
		}
	}
	g.particles = alive
}

// ---------------------------------------------------------------- draw

func (g *Game) Draw(screen *ebiten.Image) {
	now := time.Now()
	screen.Fill(colFloor)
	for x := 0; x < protocol.RoomW; x += 64 {
		vector.DrawFilledRect(screen, float32(x), 0, 1, protocol.RoomH, colGrid, false)
	}
	for y := 0; y < protocol.RoomH; y += 64 {
		vector.DrawFilledRect(screen, 0, float32(y), protocol.RoomW, 1, colGrid, false)
	}

	order := make([]*person, 0, len(g.people))
	for _, p := range g.people {
		order = append(order, p)
	}
	sort.Slice(order, func(i, j int) bool { return order[i].pos.Y < order[j].pos.Y })
	for _, p := range order {
		g.drawPerson(screen, p, now)
	}
	for _, p := range order {
		g.drawOverhead(screen, p, now)
	}
	for _, p := range g.particles {
		alpha := float32(p.life) / float32(p.max)
		if p.glyph != "" {
			g.textAt(screen, p.glyph, g.bold, p.x, p.y, withAlpha(p.col, alpha), true)
		} else {
			vector.DrawFilledRect(screen, float32(p.x), float32(p.y), 4, 4, withAlpha(p.col, alpha), false)
		}
	}

	g.drawHUD(screen, now)
	g.drawChat(screen)
	g.drawCard(screen)
	g.drawInput(screen)
	g.drawToasts(screen, now)

	if g.shotPath != "" && g.frame == 360 && !g.shotDone {
		saveScreenshot(screen, g.shotPath)
		g.shotDone = true
	}
}

func saveScreenshot(screen *ebiten.Image, path string) {
	b := screen.Bounds()
	img := image.NewRGBA(b)
	screen.ReadPixels(img.Pix)
	if f, err := os.Create(path); err == nil {
		_ = png.Encode(f, img)
		f.Close()
	}
}

func (g *Game) drawPerson(screen *ebiten.Image, p *person, now time.Time) {
	x, y := p.pos.X, p.pos.Y
	me := p.info.ID == g.you
	rows := p.art()
	w, h := g.artSize(rows)
	top := y - h/2
	if p.emote == "dance" && now.Before(p.emoteUntil) {
		top -= math.Abs(math.Sin(float64(g.frame)/4)) * 8
	}

	vector.DrawFilledRect(screen, float32(x-w/2), float32(y+h/2+1), float32(w), 4, color.RGBA{0, 0, 0, 80}, true)
	if me {
		vector.DrawFilledRect(screen, float32(x-w/2), float32(y+h/2+6), float32(w), 2, colAccent, false)
	}
	c := p.color()
	for i, row := range rows {
		g.textAt(screen, row, g.mono, x-w/2, top+float64(i*artLineH), c, false)
	}

	nameCol := colText
	if me {
		nameCol = colAccent
	}
	g.textAt(screen, p.info.Profile.Name, g.bold, x, y+h/2+10, nameCol, true)
	if p.info.Status != "" {
		g.textAt(screen, "["+p.info.Status+"]", g.small, x, y+h/2+27, colDim, true)
	}
}

var emoteLabel = map[string]string{
	"wave":     "o/  *waves*",
	"dance":    "\\o/ *dances*",
	"confetti": "*confetti!*",
	"heart":    "<3",
	"think":    "hmm...",
}

func (g *Game) drawOverhead(screen *ebiten.Image, p *person, now time.Time) {
	_, h := g.artSize(p.art())
	x, top := p.pos.X, p.pos.Y-h/2-4
	if now.Before(p.bubbleUntil) && p.bubble != "" {
		w, h := text.Measure(p.bubble, g.font, 0)
		bw, bh := float32(w+16), float32(h+8)
		bx, by := float32(x)-bw/2, float32(top)-bh-6
		vector.DrawFilledRect(screen, bx, by, bw, bh, colBubble, true)
		vector.DrawFilledRect(screen, float32(x)-4, by+bh, 8, 5, colBubble, false)
		g.textAt(screen, p.bubble, g.font, x, float64(by+4), colBubbleTx, true)
		top = float64(by) - 4
	}
	if now.Before(p.emoteUntil) {
		label := emoteLabel[p.emote]
		if label == "" {
			label = "*" + p.emote + "*"
		}
		g.textAt(screen, label, g.bold, x, top-18, colAccent, true)
	}
}

func (g *Game) drawHUD(screen *ebiten.Image, now time.Time) {
	others := make([]avatar.Avatar, 0, len(g.people))
	for id, p := range g.people {
		if id != g.you {
			others = append(others, g.asAvatar(p))
		}
	}
	near := avatar.Filter(others, func(a avatar.Avatar) bool {
		return g.me.Pos.DistanceTo(a.Pos) <= nearbyRange
	})
	names := make([]string, 0, len(near))
	for _, a := range near {
		names = append(names, a.Name)
	}
	sort.Strings(names)
	nearby := "nobody"
	if len(names) > 5 {
		nearby = strings.Join(names[:5], ", ") + fmt.Sprintf(" +%d", len(names)-5)
	} else if len(names) > 0 {
		nearby = strings.Join(names, ", ")
	}
	g.textAt(screen, "Nearby: "+nearby, g.small, 10, 8, colDim, false)

	if now.Before(g.bannerUntil) && g.banner != "" {
		w, h := text.Measure(g.banner, g.bold, 0)
		bx := float32(protocol.RoomW/2 - w/2 - 14)
		vector.DrawFilledRect(screen, bx, 58, float32(w+28), float32(h+14), colAccent, true)
		g.textAt(screen, g.banner, g.bold, protocol.RoomW/2, 65, colBubbleTx, true)
	}
}

func (g *Game) drawChat(screen *ebiten.Image) {
	lines := g.chat.Items()
	const lh, w = 18, 420
	h := float32(len(lines)*lh + 10)
	y0 := float32(protocol.RoomH - 44)
	if len(lines) == 0 {
		return
	}
	vector.DrawFilledRect(screen, 8, y0-h, w, h, colPanel, true)
	for i, l := range lines {
		y := float64(y0-h) + 5 + float64(i*lh)
		s := truncate(l.text, 56)
		if l.system {
			if l.from != "" {
				s = l.from + ": " + s
			}
			g.textAt(screen, s, g.small, 16, y+1, colDim, false)
			continue
		}
		if l.from == "" {
			g.textAt(screen, s, g.font, 16, y, colText, false)
			continue
		}
		g.textAt(screen, l.from+":", g.bold, 16, y, colAccent, false)
		nw, _ := text.Measure(l.from+": ", g.bold, 0)
		g.textAt(screen, s, g.font, 16+nw, y, colText, false)
	}
}

func (g *Game) drawCard(screen *ebiten.Image) {
	var best *person
	bestD := math.Inf(1)
	for id, p := range g.people {
		if id == g.you {
			continue
		}
		if d := math.Hypot(p.pos.X-g.me.Pos.X, p.pos.Y-g.me.Pos.Y); d < bestD {
			best, bestD = p, d
		}
	}
	if best == nil || bestD > cardRange {
		return
	}
	other := g.asAvatar(best)
	lines := other.CardLines()
	if len(lines) == 0 {
		lines = []string{other.Name}
	}
	greet := avatar.Greet(*g.me, other)

	const w, lh = 330, 20
	h := float32((len(lines)+1)*lh + 16)
	x, y := float32(protocol.RoomW-w-8), float32(protocol.RoomH-44)-h
	vector.DrawFilledRect(screen, x, y, w, h, colPanel, true)
	vector.StrokeRect(screen, x, y, w, h, 1, colAccent, true)
	for i, l := range lines {
		face := g.font
		if i == 0 {
			face = g.bold
		}
		s := truncate(l, 44)
		g.textAt(screen, s, face, float64(x+12), float64(y+8)+float64(i*lh), colText, false)
	}
	g.textAt(screen, greet, g.font, float64(x+12), float64(y+8)+float64(len(lines)*lh)+4, colAccent, false)
}

func (g *Game) drawInput(screen *ebiten.Image) {
	y := float32(protocol.RoomH - 36)
	if !g.typing {
		return
	}
	vector.DrawFilledRect(screen, 8, y, protocol.RoomW-16, 28, colPanel, true)
	vector.StrokeRect(screen, 8, y, protocol.RoomW-16, 28, 1, colAccent, true)
	s := "> " + string(g.input)
	if (g.frame/30)%2 == 0 {
		s += "_"
	}
	g.textAt(screen, s, g.font, 16, float64(y+5), colText, false)
}

func (g *Game) drawToasts(screen *ebiten.Image, now time.Time) {
	y := 100.0
	alive := g.toasts[:0]
	for _, t := range g.toasts {
		if now.After(t.until) {
			continue
		}
		alive = append(alive, t)
		w, h := text.Measure(t.text, g.bold, 0)
		vector.DrawFilledRect(screen, float32(protocol.RoomW/2-w/2-12), float32(y), float32(w+24), float32(h+10), colPanel, true)
		g.textAt(screen, t.text, g.bold, protocol.RoomW/2, y+5, t.col, true)
		y += h + 16
	}
	g.toasts = alive
}

func (g *Game) textAt(dst *ebiten.Image, s string, face *text.GoTextFace, x, y float64, c color.Color, center bool) {
	op := &text.DrawOptions{}
	if center {
		w, _ := text.Measure(s, face, 0)
		x -= w / 2
	}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(c)
	text.Draw(dst, s, face, op)
}

func withAlpha(c color.RGBA, a float32) color.Color {
	return color.NRGBA{c.R, c.G, c.B, uint8(float32(255) * max(0, min(1, a)))}
}

// truncate cuts s to at most max runes, ending with "...".
func truncate(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max-3]) + "..."
}

// NormalizeArt makes art safe to draw: each rune that is not printable ASCII
// becomes '?', and every row is padded with spaces to the widest row, so
// that a mirrored row lines up with the others.
func NormalizeArt(rows []string) []string {
	width := 0
	for _, row := range rows {
		width = max(width, len([]rune(row)))
	}
	out := make([]string, 0, len(rows))
	for _, row := range rows {
		line := make([]rune, 0, width)
		for _, r := range row {
			if r < ' ' || r > '~' {
				r = '?'
			}
			line = append(line, r)
		}
		for len(line) < width {
			line = append(line, ' ')
		}
		out = append(out, string(line))
	}
	return out
}

var palette = map[string]color.RGBA{
	"red":    {235, 80, 70, 255},
	"orange": {245, 150, 50, 255},
	"yellow": {245, 210, 60, 255},
	"green":  {100, 200, 100, 255},
	"blue":   {90, 150, 245, 255},
	"cyan":   {70, 210, 220, 255},
	"purple": {180, 120, 240, 255},
	"pink":   {245, 130, 190, 255},
	"white":  {240, 240, 240, 255},
	"grey":   {150, 155, 165, 255},
}

// ColorFor returns the colour for a name in any case. An unknown name gives
// a transparent colour, which the arena draws grey.
func ColorFor(name string) color.RGBA {
	return palette[strings.ToLower(strings.TrimSpace(name))]
}
