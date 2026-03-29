package game

import (
	"math"
	"math/rand"
	"time"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	paddleWidth        = 20
	paddleHeight       = 80
	paddleShift        = 50
	paddleAcceleration = 1

	aiMaxSpeed       = 6.0
	aiDeadZone       = 6.0
	aiAimError       = 18.0
	aiIdleAimError   = 10.0
	aiReactionMin    = 4
	aiReactionJitter = 6
)

type paddle struct {
	sprite
	yVelocity      float64
	player         int
	isComputer     *bool
	score          *int
	upKey, downKey ebiten.Key

	aiTargetY  float64
	aiCooldown int
	aiRandom   *rand.Rand
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
	p.aiTargetY = p.y
	p.aiRandom = rand.New(rand.NewSource(time.Now().UnixNano() + int64(player*101)))
	p.image = ebiten.NewImage(paddleWidth, paddleHeight)
	p.image.Fill(objectColor)
	p.ghostImage = ebiten.NewImage(paddleWidth, paddleHeight)
	p.ghostImage.Fill(ghostColor)
	return p
}

func (p *paddle) update(g *Gong) {
	if g.state == play {
		if *p.isComputer {
			p.updateComputer(g)
		} else {
			if inpututil.KeyPressDuration(p.upKey) > 0 {
				p.yVelocity -= paddleAcceleration
			}
			if inpututil.KeyPressDuration(p.downKey) > 0 {
				p.yVelocity += paddleAcceleration
			}
			p.y += p.yVelocity
		}

		if p.y < 0 {
			p.y = 1.0
			p.yVelocity = 0
			playSound(ping)
		} else if p.y+paddleHeight > windowHeight {
			p.y = windowHeight - paddleHeight - 1.0
			p.yVelocity = 0
			playSound(ping)
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

func (p *paddle) updateComputer(g *Gong) {
	if p.aiCooldown > 0 {
		p.aiCooldown--
	}

	if p.aiCooldown == 0 {
		p.aiTargetY = p.nextAITargetY(g)
		p.aiCooldown = aiReactionMin + p.aiRandom.Intn(aiReactionJitter+1)
	}

	delta := p.aiTargetY - p.y
	if math.Abs(delta) < aiDeadZone {
		p.yVelocity *= 0.7
		p.y += p.yVelocity
		return
	}

	p.yVelocity = clamp(delta*0.20, -aiMaxSpeed, aiMaxSpeed)
	p.y += p.yVelocity
}

func (p *paddle) nextAITargetY(g *Gong) float64 {
	if p.isBallApproaching(g.ball) {
		return p.predictBallY(g.ball) - paddleHeight/2 + p.randomError(aiAimError)
	}
	idleTarget := float64(windowHeight/2-paddleHeight/2) + p.randomError(aiIdleAimError)
	return idleTarget
}

func (p *paddle) isBallApproaching(b *ball) bool {
	if p.player == leftPlayer {
		return b.xVelocity < 0
	}
	return b.xVelocity > 0
}

func (p *paddle) predictBallY(b *ball) float64 {
	targetX := p.x + paddleWidth/2
	if p.player == leftPlayer {
		targetX += ballRadius
	} else {
		targetX -= ballRadius
	}

	if b.xVelocity == 0 {
		return b.y
	}
	timeToReach := (targetX - b.x) / b.xVelocity
	if timeToReach <= 0 {
		return b.y
	}

	predictedY := b.y + b.yVelocity*timeToReach
	minY := float64(ballRadius)
	maxY := float64(windowHeight - ballRadius)

	for predictedY < minY || predictedY > maxY {
		if predictedY < minY {
			predictedY = minY + (minY - predictedY)
		} else {
			predictedY = maxY - (predictedY - maxY)
		}
	}

	return predictedY
}

func (p *paddle) randomError(max float64) float64 {
	return (p.aiRandom.Float64()*2 - 1) * max
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
