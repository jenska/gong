package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ballRadius       = 10
	ballDiameter     = 2 * ballRadius
	ballVelocity     = 5.0
	ballAcceleration = 0.35
	ballMaxXVelocity = 9.0
	ballMaxYVelocity = 7.0
	ballMinYVelocity = 1.0
	paddleDeflection = 5.5
	paddleSpin       = 0.45
	serveYVelocity   = 3.25
)

type ball struct {
	sprite
	xVelocity, yVelocity float64
}

func newBall() *ball {
	b := &ball{}
	b.reset()
	b.image = ebiten.NewImage(ballDiameter, ballDiameter)
	b.image.Fill(objectColor)
	b.ghostImage = ebiten.NewImage(ballDiameter, ballDiameter)
	b.ghostImage.Fill(ghostColor)
	return b
}

func (b *ball) reset() {
	b.x, b.y = windowWidth/2, windowHeight/2
	b.xVelocity, b.yVelocity = ballVelocity, serveYVelocity
	b.trailCount, b.trailNext = 0, 0
}

func (b *ball) serve(toward Side, verticalDirection float64) {
	b.reset()
	if toward == LeftSide {
		b.xVelocity = -ballVelocity
	}
	b.yVelocity = math.Copysign(serveYVelocity, verticalDirection)
}

func (b *ball) update(g *Gong) {
	if g.state == play {
		b.x += b.xVelocity
		b.y += b.yVelocity

		if b.y+ballDiameter >= windowHeight {
			playSound(ping)
			b.yVelocity = -b.yVelocity
			b.y = windowHeight - ballDiameter
		} else if b.y <= 0 {
			playSound(ping)
			b.yVelocity = -b.yVelocity
			b.y = 0
		}
	}
	b.visible = g.state == play || g.state == interrupt
	b.recordPosition()
}

func (b *ball) approaching(side Side) bool {
	if side == LeftSide {
		return b.xVelocity < 0
	}
	return b.xVelocity > 0
}

func (b *ball) bounceFrom(p *paddle) {
	ballCenter := b.y + ballRadius
	paddleCenter := p.y + paddleHeight/2
	impact := min(max((ballCenter-paddleCenter)/(paddleHeight/2), -1), 1)

	xSpeed := min(math.Abs(b.xVelocity)+ballAcceleration, ballMaxXVelocity)
	if p.side == LeftSide {
		b.x = p.x + paddleWidth
		b.xVelocity = xSpeed
	} else {
		b.x = p.x - ballDiameter
		b.xVelocity = -xSpeed
	}

	yVelocity := b.yVelocity*0.25 + impact*paddleDeflection + p.yVelocity*paddleSpin
	yVelocity = min(max(yVelocity, -ballMaxYVelocity), ballMaxYVelocity)
	if math.Abs(yVelocity) < ballMinYVelocity {
		yVelocity = math.Copysign(ballMinYVelocity, yVelocity)
	}
	b.yVelocity = yVelocity
}
