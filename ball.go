package gong

import "github.com/hajimehoshi/ebiten"

const (
	ballRadius   = 10.0
	ballVelocity = 5.0
)

type ball struct {
	sprite
	xv, yv float64
}

func newBall() *ball {
	b := &ball{}
	b.reset()
	b.xv, b.yv = ballVelocity, ballVelocity
	b.image, _ = ebiten.NewImage(ballRadius*2, ballRadius*2, ebiten.FilterDefault)
	b.image.Fill(objectColor)
	b.ghostImage, _ = ebiten.NewImage(ballRadius*2, ballRadius*2, ebiten.FilterDefault)
	b.ghostImage.Fill(ghostColor)
	return b
}

func (b *ball) reset() {
	b.x, b.y = windowWidth/2, windowHeight/2
}

func (b *ball) update(g *Gong) {
	b.visible = g.state == play || g.state == interrupt
	if g.state == play {
		b.x += b.xv
		b.y += b.yv

		if b.y-ballRadius > windowHeight {
			b.yv = -b.yv
			b.y = windowHeight - ballRadius
		} else if b.y+ballRadius < 0 {
			b.yv = -b.yv
			b.y = ballRadius
		}
	}
}
