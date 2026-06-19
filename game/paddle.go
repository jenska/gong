package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	paddleWidth        = 20
	paddleHeight       = 80
	paddleShift        = 50
	paddleAcceleration = 1
)

type paddle struct {
	sprite
	yVelocity  float64
	side       Side
	controller Controller
}

func newPaddle(side Side, controller Controller) *paddle {
	p := &paddle{
		side:       side,
		controller: controller,
	}
	if side == LeftSide {
		p.x = paddleShift
	} else {
		p.x = windowWidth - paddleShift - paddleWidth
	}
	p.y = windowHeight/2 - paddleHeight/2
	p.image = ebiten.NewImage(paddleWidth, paddleHeight)
	p.image.Fill(objectColor)
	p.ghostImage = ebiten.NewImage(paddleWidth, paddleHeight)
	p.ghostImage.Fill(ghostColor)
	return p
}

func (p *paddle) update(g *Gong) {
	if g.state == play {
		p.applyControl(p.controller.Control(p.gameView(g)))
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
	p.visible = g.state == play || g.state == interrupt
	p.recordPosition()
}

func (p *paddle) gameView(g *Gong) GameView {
	return GameView{
		Side:           p.side,
		Width:          windowWidth,
		Height:         windowHeight,
		PaddleX:        p.x,
		PaddleY:        p.y,
		PaddleHeight:   paddleHeight,
		PaddleVelocity: p.yVelocity,
		BallX:          g.ball.x,
		BallY:          g.ball.y,
		BallVelocityX:  g.ball.xVelocity,
		BallVelocityY:  g.ball.yVelocity,
	}
}

func (p *paddle) applyControl(control Control) {
	acceleration := max(control.Acceleration, 0)
	if p.yVelocity*control.TargetVelocity < 0 ||
		control.TargetVelocity == 0 && p.yVelocity != 0 {
		acceleration = max(control.Braking, 0)
	}
	p.yVelocity = moveTowards(p.yVelocity, control.TargetVelocity, acceleration)
}

func moveTowards(current, target, maximumDelta float64) float64 {
	if math.Abs(target-current) <= maximumDelta {
		return target
	}
	return current + math.Copysign(maximumDelta, target-current)
}
