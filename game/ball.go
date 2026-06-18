package game

import "github.com/hajimehoshi/ebiten/v2"

const (
	ballRadius       = 10
	ballDiameter     = 2 * ballRadius
	ballVelocity     = 5.0
	ballAcceleration = 0.1
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
	b.xVelocity, b.yVelocity = ballVelocity, ballVelocity
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
