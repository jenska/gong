package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type (
	gameState byte

	// Gong is the game state consumed by Ebitengine.
	Gong struct {
		state                    gameState
		objects                  []gameObject
		ball                     *ball
		score1, score2           int
		isComputer1, isComputer2 bool
		interruptTicks           int
	}

	gameObject interface {
		update(g *Gong)
		draw(screen *ebiten.Image)
	}
)

const (
	windowWidth  = 800
	windowHeight = 600

	maxScore = 10

	interruptDurationTicks = ebiten.DefaultTPS / 2

	leftPlayer  = 0
	rightPlayer = 1

	start gameState = iota
	controls
	play
	pause
	interrupt
	gameOver
)

var (
	ghostColor  = color.RGBA{32, 32, 32, 255}
	screenColor = color.RGBA{60, 60, 64, 255}
	objectColor = color.White
)

// NewGong creates a Gong game.
func NewGong() *Gong {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Gong! - The Go Pong")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetWindowResizable(true)
	g := &Gong{}
	g.reset()
	return g
}

func (g *Gong) reset() {
	g.score1, g.score2 = 0, 0
	g.isComputer1, g.isComputer2 = false, false
	g.interruptTicks = 0
	g.ball = newBall()
	g.objects = []gameObject{
		newPaddle(leftPlayer, &g.score1, &g.isComputer1),
		g.ball,
		newPaddle(rightPlayer, &g.score2, &g.isComputer2),
		newHUD(),
	}
	g.state = start
}

// Layout reports the game's fixed logical screen size.
func (g *Gong) Layout(_, _ int) (int, int) {
	return windowWidth, windowHeight
}

func isMenuSelected(key ebiten.Key) bool {
	if inpututil.IsKeyJustPressed(key) {
		playSound(menuSelect)
		return true
	}
	return false
}

// Update advances the game state by one tick.
func (g *Gong) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	switch g.state {
	case start:
		if isMenuSelected(ebiten.KeyH) {
			g.state = controls
		} else if isMenuSelected(ebiten.KeyA) {
			g.isComputer1, g.isComputer2 = true, false
			g.state = play
		} else if isMenuSelected(ebiten.KeyV) {
			g.isComputer1, g.isComputer2 = false, false
			g.state = play
		} else if isMenuSelected(ebiten.KeyB) {
			g.isComputer1, g.isComputer2 = true, true
			g.state = play
		} else if isMenuSelected(ebiten.KeyF) {
			ebiten.SetFullscreen(!ebiten.IsFullscreen())
		} else if isMenuSelected(ebiten.KeyW) {
			sl.volume = math.Min(1.0, sl.volume+0.1)
		} else if isMenuSelected(ebiten.KeyS) {
			sl.volume = math.Max(0.0, sl.volume-0.1)
		}
	case controls:
		if isMenuSelected(ebiten.KeySpace) {
			g.state = start
		}
	case pause:
		if isMenuSelected(ebiten.KeySpace) {
			g.state = play
		} else if isMenuSelected(ebiten.KeyR) {
			g.reset()
		}
	case play:
		if isMenuSelected(ebiten.KeySpace) {
			g.state = pause
		}
	case gameOver:
		if isMenuSelected(ebiten.KeySpace) {
			g.reset()
		}
	case interrupt:
		g.interruptTicks--
		if g.interruptTicks <= 0 {
			g.state = play
		}
	}

	for _, object := range g.objects {
		object.update(g)
	}
	return nil
}

// Draw renders the current game state.
func (g *Gong) Draw(screen *ebiten.Image) {
	screen.Fill(screenColor)

	for _, object := range g.objects {
		object.draw(screen)
	}
}

func (g *Gong) interrupt() {
	g.ball.reset()
	g.interruptTicks = interruptDurationTicks
	g.state = interrupt
}
