package game

import (
	"image"
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type (
	gameState byte

	// Gong is the game state consumed by Ebitengine.
	Gong struct {
		state          gameState
		objects        []gameObject
		ball           *ball
		leftPaddle     *paddle
		rightPaddle    *paddle
		hud            *hud
		courtImage     *ebiten.Image
		score1, score2 int
		aiLevel        aiLevel
		interruptTicks int
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
	g.interruptTicks = 0
	g.ball = newBall()
	g.leftPaddle = newPaddle(
		LeftSide,
		NewKeyboardController(ebiten.KeyW, ebiten.KeyS),
	)
	g.rightPaddle = newPaddle(
		RightSide,
		NewKeyboardController(ebiten.KeyUp, ebiten.KeyDown),
	)
	g.hud = newHUD()
	g.courtImage = newCourtImage()
	g.objects = []gameObject{
		g.leftPaddle,
		g.ball,
		g.rightPaddle,
		g.hud,
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
		} else if isMenuSelected(ebiten.Key1) {
			g.aiLevel = beginnerLevel
		} else if isMenuSelected(ebiten.Key2) {
			g.aiLevel = humanLikeLevel
		} else if isMenuSelected(ebiten.Key3) {
			g.aiLevel = perfectLevel
		} else if isMenuSelected(ebiten.KeyA) {
			g.startSelectedMatch(true, false)
		} else if isMenuSelected(ebiten.KeyV) {
			g.startSelectedMatch(false, false)
		} else if isMenuSelected(ebiten.KeyB) {
			g.startSelectedMatch(true, true)
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

	g.leftPaddle.update(g)
	g.rightPaddle.update(g)
	g.ball.update(g)
	if g.state == play {
		g.resolveBall()
	}
	g.hud.update(g)
	return nil
}

// Draw renders the current game state.
func (g *Gong) Draw(screen *ebiten.Image) {
	screen.Fill(screenColor)
	screen.DrawImage(g.courtImage, nil)

	for _, object := range g.objects {
		object.draw(screen)
	}
}

func (g *Gong) interrupt(serveToward Side) {
	verticalDirection := -1.0
	if rand.IntN(2) == 1 {
		verticalDirection = 1
	}
	g.ball.serve(serveToward, verticalDirection)
	g.interruptTicks = interruptDurationTicks
	g.state = interrupt
}

// StartMatch resets the game and immediately starts a match using the supplied
// controllers. It can be called after NewGong to run custom controllers.
func (g *Gong) StartMatch(left, right Controller) {
	if left == nil || right == nil {
		panic("game: controllers must not be nil")
	}
	g.reset()
	g.leftPaddle.controller = left
	g.rightPaddle.controller = right
	g.randomServe()
	g.state = play
}

func (g *Gong) startSelectedMatch(computerLeft, computerRight bool) {
	if computerLeft {
		g.leftPaddle.controller = g.newAIController()
	}
	if computerRight {
		g.rightPaddle.controller = g.newAIController()
	}
	g.randomServe()
	g.state = play
}

func (g *Gong) newAIController() Controller {
	switch g.aiLevel {
	case beginnerLevel:
		return NewBeginnerAI()
	case perfectLevel:
		return NewPerfectAI()
	default:
		return NewHumanLikeAI()
	}
}

func (g *Gong) randomServe() {
	toward := LeftSide
	if rand.IntN(2) == 1 {
		toward = RightSide
	}
	verticalDirection := -1.0
	if rand.IntN(2) == 1 {
		verticalDirection = 1
	}
	g.ball.serve(toward, verticalDirection)
}

func (g *Gong) resolveBall() {
	for _, paddle := range []*paddle{g.leftPaddle, g.rightPaddle} {
		if g.ball.approaching(paddle.side) && paddle.intersects(&g.ball.sprite) {
			g.ball.bounceFrom(paddle)
			playSound(pong)
			return
		}
	}

	switch {
	case g.ball.x > windowWidth:
		g.pointScored(LeftSide)
	case g.ball.x+ballDiameter < 0:
		g.pointScored(RightSide)
	}
}

func (g *Gong) pointScored(side Side) {
	score := &g.score1
	serveToward := RightSide
	if side == RightSide {
		score = &g.score2
		serveToward = LeftSide
	}
	*score++
	if *score >= maxScore {
		playSound(win)
		g.state = gameOver
		return
	}
	playSound(lost)
	g.interrupt(serveToward)
}

func newCourtImage() *ebiten.Image {
	court := ebiten.NewImage(windowWidth, windowHeight)
	lineColor := color.RGBA{R: 255, G: 255, B: 255, A: 45}
	const (
		lineWidth  = 4
		dashHeight = 20
		dashGap    = 14
	)
	x := windowWidth/2 - lineWidth/2
	for y := 0; y < windowHeight; y += dashHeight + dashGap {
		bottom := min(y+dashHeight, windowHeight)
		court.SubImage(image.Rect(x, y, x+lineWidth, bottom)).(*ebiten.Image).Fill(lineColor)
	}
	return court
}
