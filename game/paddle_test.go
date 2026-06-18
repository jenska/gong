package game

import (
	"math/rand/v2"
	"testing"
)

func TestPaddleIsBallApproaching(t *testing.T) {
	tests := []struct {
		name     string
		player   int
		velocity float64
		want     bool
	}{
		{name: "left approaching", player: leftPlayer, velocity: -1, want: true},
		{name: "left receding", player: leftPlayer, velocity: 1, want: false},
		{name: "right approaching", player: rightPlayer, velocity: 1, want: true},
		{name: "right receding", player: rightPlayer, velocity: -1, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := paddle{player: tt.player}
			b := ball{xVelocity: tt.velocity}
			if got := p.isBallApproaching(&b); got != tt.want {
				t.Fatalf("isBallApproaching() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaddlePredictBallY(t *testing.T) {
	tests := []struct {
		name string
		ball ball
		want float64
	}{
		{
			name: "direct path",
			ball: ball{sprite: sprite{x: 400, y: 100}, xVelocity: 5, yVelocity: 5},
			want: 430,
		},
		{
			name: "reflects from bottom",
			ball: ball{sprite: sprite{x: 400, y: 500}, xVelocity: 5, yVelocity: 5},
			want: 350,
		},
		{
			name: "stationary",
			ball: ball{sprite: sprite{x: 400, y: 123}},
			want: 123,
		},
	}

	p := paddle{player: rightPlayer, sprite: sprite{x: windowWidth - paddleShift - paddleWidth}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := p.predictBallY(&tt.ball); got != tt.want {
				t.Fatalf("predictBallY() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComputerPaddleAcceleratesInsteadOfSnappingToSpeed(t *testing.T) {
	p := paddle{
		player:     rightPlayer,
		sprite:     sprite{x: windowWidth - paddleShift - paddleWidth, y: 100},
		aiTargetY:  500,
		aiCooldown: 10,
		aiTracking: true,
		aiRandom:   testRandom(),
		yVelocity:  0,
	}
	g := Gong{
		ball: &ball{
			sprite:    sprite{x: windowWidth / 2, y: 200},
			xVelocity: ballVelocity,
		},
	}

	p.updateComputer(&g)

	if p.yVelocity != aiAcceleration {
		t.Fatalf("yVelocity = %v, want gradual acceleration of %v", p.yVelocity, aiAcceleration)
	}
	if p.y != 100+aiAcceleration {
		t.Fatalf("y = %v, want %v", p.y, 100+aiAcceleration)
	}
}

func TestComputerPaddleBrakesBeforeChangingDirection(t *testing.T) {
	p := paddle{
		player:     rightPlayer,
		sprite:     sprite{x: windowWidth - paddleShift - paddleWidth, y: 300},
		aiTargetY:  0,
		aiCooldown: 10,
		aiTracking: true,
		aiRandom:   testRandom(),
		yVelocity:  2,
	}
	g := Gong{
		ball: &ball{
			sprite:    sprite{x: windowWidth / 2, y: 200},
			xVelocity: ballVelocity,
		},
	}

	p.updateComputer(&g)

	want := 2 - aiBraking
	if p.yVelocity != want {
		t.Fatalf("yVelocity = %v, want braking to %v", p.yVelocity, want)
	}
}

func TestComputerPaddleReactsMoreSlowlyToDistantBall(t *testing.T) {
	near := paddle{
		player:   rightPlayer,
		sprite:   sprite{x: windowWidth - paddleShift - paddleWidth},
		aiRandom: testRandom(),
	}
	far := near
	far.aiRandom = testRandom()

	nearBall := ball{sprite: sprite{x: near.x - 50}}
	farBall := ball{sprite: sprite{x: 50}}

	nearDelay := near.reactionDelay(&nearBall)
	farDelay := far.reactionDelay(&farBall)
	if farDelay <= nearDelay {
		t.Fatalf("far reaction delay = %d, want greater than near delay %d", farDelay, nearDelay)
	}
}

func TestComputerAimStaysWithinReachableArea(t *testing.T) {
	p := paddle{
		player:   rightPlayer,
		sprite:   sprite{x: windowWidth - paddleShift - paddleWidth},
		aiRandom: testRandom(),
	}
	g := Gong{
		ball: &ball{
			sprite:    sprite{x: 100, y: 10},
			xVelocity: ballVelocity,
			yVelocity: -ballVelocity,
		},
	}

	for range 100 {
		target := p.nextAITargetY(&g)
		if target < 0 || target > windowHeight-paddleHeight {
			t.Fatalf("target = %v, want within [0, %d]", target, windowHeight-paddleHeight)
		}
	}
}

func testRandom() *rand.Rand {
	return rand.New(rand.NewPCG(1, 2))
}
