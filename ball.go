package gong

import "github.com/hajimehoshi/ebiten"

const (
	ballRadius   = 10.0
	ballVelocity = 5.0
)

type ball struct {
	x, y       float64
	xv, yv     float64
	image      *ebiten.Image
	ghostImage *ebiten.Image
	visible    bool
	fader      xyBuffer
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

func (b *ball) draw(screen *ebiten.Image) {
	if b.visible {
		op1 := &ebiten.DrawImageOptions{}
		op1.GeoM.Translate(b.x+5, b.y)
		screen.DrawImage(b.ghostImage, op1)

		op1.GeoM.Translate(-5, 0)
		screen.DrawImage(b.image, op1)
		for i := 0; i < b.fader.size(); i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(b.fader.x[i], b.fader.y[i])
			op.ColorM.Scale(1.0, 1.0, 1.0, float64(i)/15.0)
			screen.DrawImage(b.image, op)
		}
		b.fader.save(b.x, b.y)
	}
}
