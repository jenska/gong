package gong

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	paddleWidth  = 20
	paddleHeight = 100
	paddleShift  = 50
	paddleSpeed  = 10.0
)

type paddle struct {
	sprite
	player         int
	isComputer     *bool
	score          *int
	upKey, downKey ebiten.Key
	up, down       bool
}

func newPaddle(player int, score *int, mode *bool) *paddle {
	p := &paddle{}
	p.player = player
	p.score = score
	p.isComputer = mode
	if player == leftPlayer {
		p.x = paddleShift
		p.upKey = ebiten.KeyW
		p.downKey = ebiten.KeyS
	} else {
		p.x = windowWidth - paddleShift - paddleWidth
		p.upKey = ebiten.KeyUp
		p.downKey = ebiten.KeyDown
	}
	p.y = windowHeight/2 - paddleHeight/2
	p.image, _ = ebiten.NewImage(paddleWidth, paddleHeight, ebiten.FilterDefault)
	p.image.Fill(objectColor)
	p.ghostImage, _ = ebiten.NewImage(paddleWidth, paddleHeight, ebiten.FilterDefault)
	p.ghostImage.Fill(ghostColor)
	return p
}

func (p *paddle) update(g *Gong) {
	p.visible = g.state == play || g.state == interrupt
	if g.state == play {
		if *p.isComputer {
			p.y = g.ball.y - paddleHeight/2
		} else {

			if inpututil.IsKeyJustPressed(p.upKey) {
				p.up, p.down = true, false
			} else if inpututil.IsKeyJustReleased(p.upKey) || !ebiten.IsKeyPressed(p.upKey) {
				p.up = false
			}
			if inpututil.IsKeyJustPressed(p.downKey) {
				p.down, p.up = true, false
			} else if inpututil.IsKeyJustReleased(p.downKey) || !ebiten.IsKeyPressed(p.downKey) {
				p.down = false
			}

			if p.up {
				p.y -= paddleSpeed
			} else if p.down {
				p.y += paddleSpeed
			}

			if p.y < 0 {
				p.y = 1.0
			} else if p.y+paddleHeight > windowHeight {
				p.y = windowHeight - paddleHeight - 1.0
			}

		}

		// bounce ball off paddle
		if p.intersects(&g.ball.sprite) {
			if p.player == leftPlayer {
				g.ball.x = p.x + paddleWidth/2 + ballRadius
				g.ball.xv = -g.ball.xv
			} else {
				g.ball.x = p.x - paddleWidth/2 - ballRadius
				g.ball.xv = -g.ball.xv
			}
		}

		// scored
		if (g.ball.x < 0 && p.player == rightPlayer) || (g.ball.x > windowWidth && p.player == leftPlayer) {
			*p.score++
			g.state = interrupt
		}

		if *p.score > maxScore {
			g.state = gameOver
		}
	}
}
