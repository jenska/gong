package gong

import (
	"image/color"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// Gong game object
type (
	gameState byte
	Gong      struct {
		state   gameState
		sprites []sprite
		score   []int
		mode    []bool
	}

	sprite interface {
		update(g *Gong)
		draw(screen *ebiten.Image)
	}

	xyBuffer struct {
		x []float64
		y []float64
	}
)

const (
	windowWidth  = 800
	windowHeight = 600

	maxScore = 10

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

// NewGong creates a new gong object
func NewGong() *Gong {
	g := &Gong{}
	g.reset()
	return g
}

func (g *Gong) reset() {
	ball := newBall()
	g.sprites = []sprite{
		newPaddle(leftPlayer, ball),
		ball,
		newPaddle(rightPlayer, ball),
		newHUD(),
	}
	g.score = []int{0, 0}
	g.mode = []bool{false, false}
	g.state = start

}

// Layout sets the screen layout
func (g *Gong) Layout(outsideWidth, outsideHeight int) (int, int) {
	return windowWidth, windowHeight
}

// Update game state and sprites
func (g *Gong) Update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	switch g.state {
	case start:
		if inpututil.IsKeyJustPressed(ebiten.KeyH) {
			g.state = controls
		} else if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			g.mode[leftPlayer] = true
			g.state = play
		} else if inpututil.IsKeyJustPressed(ebiten.KeyV) {
			g.mode[leftPlayer] = false
			g.state = play
		} else if inpututil.IsKeyJustPressed(ebiten.KeyB) {
			g.mode[leftPlayer] = true
			g.mode[rightPlayer] = true
			g.state = play
		}
	case controls:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = start
		}
	case pause:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = play
		} else if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.reset()
		}
	case play:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = pause
		}
	case gameOver:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.reset()
		}
	case interrupt:
		time.Sleep(time.Second / 2)
		g.state = play
	}

	for _, object := range g.sprites {
		object.update(g)
	}
	return nil
}

// Draw updates the game screen elements drawn
func (g *Gong) Draw(screen *ebiten.Image) {
	screen.Fill(screenColor)
	for _, object := range g.sprites {
		object.draw(screen)
	}
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
}

func (b *xyBuffer) save(x, y float64) {
	if b.size() < 3 {
		b.x = append(b.x, x)
		b.y = append(b.y, y)
	} else {
		b.x = append(b.x[1:], x)
		b.y = append(b.y[1:], y)
	}
}
func (b *xyBuffer) size() int {
	if b.x == nil {
		b.x = make([]float64, 0)
		b.y = make([]float64, 0)
	}
	return len(b.x)
}
