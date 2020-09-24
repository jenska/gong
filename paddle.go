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
	player         int
	x, y           float64
	upKey, downKey ebiten.Key
	up, down       bool
	ball           *ball
	score          int
	image          *ebiten.Image
	ghostImage     *ebiten.Image
	visible        bool
	fader          xyBuffer
}

func newPaddle(player int, ball *ball) *paddle {
	p := &paddle{}
	p.player = player
	p.ball = ball
	if player == leftPlayer {
		p.x = paddleShift
		p.y = windowHeight / 2
		p.upKey = ebiten.KeyW
		p.downKey = ebiten.KeyS
	} else {
		p.x = windowWidth - paddleShift - paddleWidth
		p.y = windowHeight / 2
		p.upKey = ebiten.KeyUp
		p.downKey = ebiten.KeyDown
	}
	p.image, _ = ebiten.NewImage(paddleWidth, paddleHeight, ebiten.FilterDefault)
	p.image.Fill(objectColor)
	p.ghostImage, _ = ebiten.NewImage(paddleWidth, paddleHeight, ebiten.FilterDefault)
	p.ghostImage.Fill(ghostColor)
	return p
}

func (p *paddle) update(g *Gong) {
	p.visible = g.state == play || g.state == interrupt
	if g.state == play {
		if g.mode[p.player] {
			p.y = p.ball.y
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

			if p.y-paddleHeight/2.0 < 0 {
				p.y = 1.0 + paddleHeight/2
			} else if p.y+paddleHeight/2 > windowHeight {
				p.y = windowHeight - paddleHeight/2 - 1.0
			}

		}

		// bounce ball off paddle
		bx, by := p.ball.x, p.ball.y
		if p.player == leftPlayer {
			if bx-ballRadius < p.x+paddleWidth/2 && by > p.y-paddleHeight/2 && by < p.y+paddleHeight/2 {
				p.ball.x = p.x + paddleWidth/2 + ballRadius
				p.ball.xv = -p.ball.xv
			}
		} else {
			if bx+ballRadius > p.x-paddleWidth/2 && by > p.y-paddleHeight/2 && by < p.y+paddleHeight/2 {
				p.ball.x = p.x - paddleWidth/2 - ballRadius
				p.ball.xv = -p.ball.xv
			}
		}

		// scored
		if bx < 0 || bx > windowWidth {
			g.score[p.player]++
			p.ball.reset()
			g.state = interrupt
		}

		if g.score[p.player] > maxScore {
			g.state = gameOver
		}
	}
}

func (p *paddle) draw(screen *ebiten.Image) {
	if p.visible {
		var op *ebiten.DrawImageOptions

		for i := 0; i < p.fader.size(); i++ {
			op = &ebiten.DrawImageOptions{}
			op.GeoM.Translate(p.fader.x[i], p.fader.y[i])
			op.ColorM.Scale(1.0, 1.0, 1.0, float64(i)/10.0)
			screen.DrawImage(p.image, op)
		}
		// draw player's paddle
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(p.x+5, p.y-paddleHeight/2)
		screen.DrawImage(p.ghostImage, op)
		op.GeoM.Translate(-5, 0)
		screen.DrawImage(p.image, op)
		p.fader.save(p.x, p.y-paddleHeight/2)
	}
}
