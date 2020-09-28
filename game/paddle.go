package game

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	paddleWidth        = 20
	paddleHeight       = 80
	paddleShift        = 50
	paddleAcceleration = 1
)

type paddle struct {
	sprite
	yVelocity      float64
	player         int
	isComputer     *bool
	score          *int
	upKey, downKey ebiten.Key
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
	if g.state == play {
		if *p.isComputer {
			p.y = g.ball.y - paddleHeight/2
		} else {
			if inpututil.KeyPressDuration(p.upKey) > 0 {
				p.yVelocity -= paddleAcceleration
			}
			if inpututil.KeyPressDuration(p.downKey) > 0 {
				p.yVelocity += paddleAcceleration
			}

			p.y += p.yVelocity

			if p.y < 0 {
				p.y = 1.0
				p.yVelocity = 0
				playSound(ping)
			} else if p.y+paddleHeight > windowHeight {
				p.y = windowHeight - paddleHeight - 1.0
				p.yVelocity = 0
				playSound(ping)
			}

		}

		// inelastic collision
		if p.intersects(&g.ball.sprite) {
			playSound(pong)
			if p.player == leftPlayer {
				g.ball.x = p.x + paddleWidth/2 + ballRadius
			} else {
				g.ball.x = p.x - paddleWidth/2 - ballRadius
			}
			g.ball.xVelocity = -g.ball.xVelocity - ballAcceleration
			if g.ball.yVelocity > 0 {
				g.ball.yVelocity += ballAcceleration
			} else {
				g.ball.yVelocity -= ballAcceleration
			}
		}

		// scored
		if (g.ball.x < 0 && p.player == rightPlayer) || (g.ball.x > windowWidth && p.player == leftPlayer) {
			*p.score++
			g.state = interrupt
			if *p.score > maxScore {
				playSound(win)
				g.state = gameOver
			} else {
				playSound(lost)
			}
		}
	}
	p.visible = g.state == play || g.state == interrupt
}
